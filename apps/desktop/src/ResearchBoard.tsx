import { FormEvent } from "react";
import { ComputeCard } from "./ComputeCard";
import { EvidenceDrawer } from "./EvidenceDrawer";
import { PreviewNote } from "./PreviewNote";
import { committeeVerified, oidBelongsToPreview, researchCardTitle } from "./honesty";
import { explainStop, explainStopHref } from "./explain";
import { researchWhyCopy } from "./researchWhy";
import { hyperliquidAPI } from "./links";
import { ExternalLink } from "./ExternalLink";
import type { BindResult, DoctorCheck, LocalStatus } from "./companion";

const RESEARCH_STAGES = [
  "READING_MARKET",
  "SEALING_PRIVATE_BOOK",
  "CONTACTING_PRIVATE_PROVIDER",
  "RECEIVING_SEALED_RESPONSE",
  "VERIFYING_TEE_SIGNATURE",
  "RESEARCHER",
  "CHALLENGER",
  "RISK",
  "DETERMINISTIC_ENGINE",
  "POLICY",
  "PREVIEW",
  "READY",
] as const;

type ResearchRole = {
  role?: string;
  verify_e2ee?: string;
  pubkey_signer?: string;
  proposed_side?: string;
  survives?: boolean;
  kill?: boolean;
  elapsed_ms?: number;
};

function roleVerified(roles: ResearchRole[], name: string) {
  return roles.some(
    (r) => String(r.role || "").toLowerCase() === name && String(r.verify_e2ee || "").toUpperCase() === "OK",
  );
}

function canonicalResearchStage(stage: string, roles: ResearchRole[]) {
  const s = (stage || "").toUpperCase();
  if (s === "RISK_START" || s === "RISK_HTTP_REQUEST" || s === "RISK_HTTP_RESPONSE" || s === "RISK_E2EE_VERIFY") {
    return "RISK";
  }
  if (s.endsWith("_VERIFIED")) {
    const base = s.slice(0, -"_VERIFIED".length);
    if (base === "RESEARCHER") return "CHALLENGER";
    if (base === "CHALLENGER") return "RISK";
    if (base === "RISK") return "DETERMINISTIC_ENGINE";
  }
  if (s.endsWith("_FAILED")) return s.slice(0, -"_FAILED".length);
  if (s === "CONTACTING_PRIVATE_PROVIDER" || s === "RECEIVING_SEALED_RESPONSE" || s === "VERIFYING_TEE_SIGNATURE") {
    if (roleVerified(roles, "challenger")) return "RISK";
    if (roleVerified(roles, "researcher")) return "CHALLENGER";
    return s;
  }
  if ((RESEARCH_STAGES as readonly string[]).includes(s)) return s;
  return s;
}

function stageMark(name: string, stage: string, roles: ResearchRole[]) {
  const s = canonicalResearchStage(stage, roles);
  if (name === "RESEARCHER") {
    if (roleVerified(roles, "researcher") || ["CHALLENGER", "RISK", "DETERMINISTIC_ENGINE", "POLICY", "PREVIEW", "READY"].includes(s)) {
      return "done";
    }
    if (s === "RESEARCHER") return "lit";
    return "";
  }
  if (name === "CHALLENGER") {
    if (roleVerified(roles, "challenger") || ["RISK", "DETERMINISTIC_ENGINE", "POLICY", "PREVIEW", "READY"].includes(s)) {
      return "done";
    }
    if (s === "CHALLENGER") return "lit";
    return "";
  }
  if (name === "RISK") {
    if (roleVerified(roles, "risk") || ["DETERMINISTIC_ENGINE", "POLICY", "PREVIEW", "READY"].includes(s)) return "done";
    if (s === "RISK") return "lit";
    return "";
  }
  const current = RESEARCH_STAGES.indexOf(s as (typeof RESEARCH_STAGES)[number]);
  const i = RESEARCH_STAGES.indexOf(name as (typeof RESEARCH_STAGES)[number]);
  if (current < 0 || i < 0) return "";
  if (i < current) return "done";
  if (i === current) return "lit";
  return "";
}

function teeState(roles: ResearchRole[], busy: boolean, stop: string | null) {
  const teeFail =
    stop === "TEE_VERIFY_FAIL" ||
    stop === "TEE_SIGNATURE_INVALID" ||
    stop === "TEE_SIGNER_MISMATCH" ||
    stop === "TEE_RESPONSE_INVALID" ||
    stop === "TEE_OPEN_FAIL";
  if (teeFail) return "failure";
  if (committeeVerified(roles)) return "success";
  if (busy) return "running";
  return "pending";
}

function engineState(stage: string, roles: ResearchRole[], busy: boolean) {
  const s = canonicalResearchStage(stage, roles);
  if (["POLICY", "PREVIEW", "READY"].includes(s) && !busy) return "success";
  if (s === "DETERMINISTIC_ENGINE" || s === "POLICY" || s === "PREVIEW") return busy ? "running" : "success";
  if (busy && roleVerified(roles, "risk")) return "running";
  return "pending";
}

function roleState(roles: ResearchRole[], name: string, stage: string) {
  const row = roles.find((r) => String(r.role || "").toLowerCase() === name);
  if (row?.kill || row?.survives === false) return "stand-down";
  const mark = stageMark(name.toUpperCase(), stage, roles);
  if (roleVerified(roles, name)) return "success";
  if (mark === "lit") return "running";
  if (mark === "done") return "success";
  return "pending";
}

function namedRoleLabel(st: string, stop: string | null, kind: string) {
  if (stop === "POLICY_DENIED" || stop === "POLICY_REJECTED" || kind === "POLICY_DENIED") return "POLICY BLOCKED";
  if (stop === "DIRECT_PROVIDER_TIMEOUT" || kind === "DIRECT_PROVIDER_TIMEOUT") return "PROVIDER TIMEOUT";
  if (stop === "DIRECT_PROVIDER_UNAVAILABLE" || kind === "DIRECT_PROVIDER_UNAVAILABLE") return "PROVIDER UNAVAILABLE";
  if (stop === "COMPANION_NOT_RUNNING" || stop === "companion_down") return "COMPANION FAILURE";
  if (
    stop === "TEE_VERIFY_FAIL" ||
    stop === "TEE_SIGNATURE_INVALID" ||
    stop === "TEE_SIGNER_MISMATCH" ||
    stop === "TEE_RESPONSE_INVALID" ||
    stop === "TEE_OPEN_FAIL"
  ) {
    return "TEE FAILURE";
  }
  if (st === "stand-down") return "STOOD DOWN";
  if (st === "success") return "VERIFIED";
  if (st === "running") return "RUNNING";
  if (st === "failure") return "TEE FAILURE";
  return "PENDING";
}

function roleReason(roles: ResearchRole[], name: string) {
  const row = roles.find((r) => String(r.role || "").toLowerCase() === name);
  if (!row) return "Waiting.";
  if (row.kill || row.survives === false) return row.proposed_side ? `Stood down after ${row.proposed_side}` : "Stood down.";
  if (String(row.verify_e2ee || "").toUpperCase() === "OK") return row.proposed_side ? `Proposed ${row.proposed_side}` : "Verified.";
  return "Not finished.";
}

export function ResearchBoard({
  coin,
  hypothesis,
  setHypothesis,
  researchBusy,
  researchStage,
  researchElapsed,
  researchRoles,
  pollMiss,
  researchJobId,
  researchKind,
  researchNote,
  researchStop,
  preview,
  previewHash,
  authTyped,
  setAuthTyped,
  authBusy,
  authErr,
  lastOid,
  status,
  sessionAlive,
  checks,
  techOpen,
  setTechOpen,
  researchEvidenceText,
  eligible,
  coins,
  net,
  onResearch,
  onCancel,
  onAuthorize,
  onCancelBound,
  onCheck,
}: {
  coin: string;
  hypothesis: "none" | "long" | "short";
  setHypothesis: (h: "none" | "long" | "short") => void;
  researchBusy: boolean;
  researchStage: string;
  researchElapsed: number;
  researchRoles: ResearchRole[];
  pollMiss?: boolean;
  researchJobId: string;
  researchKind: string;
  researchNote: string | null;
  researchStop: string | null;
  preview: NonNullable<BindResult["preview"]> | null;
  previewHash: string;
  authTyped: string;
  setAuthTyped: (v: string) => void;
  authBusy: boolean;
  authErr: string | null;
  lastOid: string;
  status: LocalStatus | null;
  sessionAlive: boolean;
  checks: DoctorCheck[];
  techOpen: boolean;
  setTechOpen: (v: boolean | ((c: boolean) => boolean)) => void;
  researchEvidenceText: string;
  eligible: Array<{ coin: string; mark: number; reason: string; why?: string; oracle?: number; funding?: number; openInterest?: number }>;
  coins?: Array<{ coin: string; mark: number; oracle?: number; funding?: number; openInterest?: number; reason?: string; why?: string }>;
  net: string;
  onResearch: (coin?: string) => void;
  onCancel: () => void;
  onAuthorize: (e: FormEvent) => void;
  onCancelBound: (e: FormEvent) => void;
  onCheck: () => void;
}) {
  const shown = canonicalResearchStage(researchStage, researchRoles);
  const whyCode =
    researchStop ||
    (researchKind && researchKind !== "READY_ELIGIBLE" ? researchKind : "") ||
    preview?.deny ||
    null;
  const explained = explainStop(whyCode || null);
  const verified = committeeVerified(researchRoles);
  const snap = (coins || []).find((c) => c.coin === (coin || "ETH")) || eligible.find((c) => c.coin === (coin || "ETH"));
  const title = researchCardTitle(researchKind || researchStop, verified);
  const whyRows = !researchBusy
    ? researchWhyCopy({
        coin: coin || "ETH",
        kind: researchKind,
        note: researchNote,
        stop: researchStop,
        deny: preview?.deny,
        eligible: Boolean(preview?.eligible),
        roles: researchRoles,
        snap: snap,
      })
    : [];
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Research</p>
          <h1>{coin || "ETH"}</h1>
        </div>
        <p className="fine" style={{ margin: 0 }}>
          Private committee. Host sizes. Chat cannot AUTHORIZE.
        </p>
      </div>
      <ol className="life-strip" aria-label="Research lifecycle">
        {[
          ["DISCOVERED", "READING_MARKET"],
          ["PRIVATE BOOK SEALED", "SEALING_PRIVATE_BOOK"],
          ["RESEARCHER", "RESEARCHER"],
          ["CHALLENGER", "CHALLENGER"],
          ["RISK", "RISK"],
          ["TEE VERIFICATION", "VERIFYING_TEE_SIGNATURE"],
          ["HOST ENGINE", "DETERMINISTIC_ENGINE"],
          ["POLICY", "POLICY"],
          ["DECISION", "PREVIEW"],
        ].map(([label, key]) => (
          <li key={label} className={stageMark(key, researchStage, researchRoles) || (researchKind && key === "PREVIEW" ? "done" : "")}>
            {label}
          </li>
        ))}
      </ol>
      {snap ? (
        <div className="snap">
          <div>
            <span>Mark</span>
            <strong>{snap.mark}</strong>
          </div>
          <div>
            <span>Oracle</span>
            <strong>{"oracle" in snap ? snap.oracle ?? "—" : "—"}</strong>
          </div>
          <div>
            <span>Funding</span>
            <strong>{"funding" in snap ? snap.funding ?? "—" : "—"}</strong>
          </div>
          <div>
            <span>OI</span>
            <strong>{"openInterest" in snap && snap.openInterest ? Math.round(Number(snap.openInterest)) : "—"}</strong>
          </div>
        </div>
      ) : null}
      <p className="fine">{snap && "why" in snap && snap.why ? snap.why : "Selected from Markets because it is in your policy universe with a live venue book."}</p>
      <p className="label" style={{ marginTop: 12 }}>
        Thesis
      </p>
      <div className="cta-row">
        {(["none", "long", "short"] as const).map((h) => (
          <button
            key={h}
            type="button"
            className={hypothesis === h ? "linkish on" : "linkish off"}
            onClick={() => setHypothesis(h)}
            disabled={researchBusy}
          >
            {h === "none" ? "No bias" : h === "long" ? "Consider long" : "Consider short"}
          </button>
        ))}
        <button
          type="button"
          className="primary"
          onClick={() => onResearch(coin)}
          disabled={researchBusy || !checks.find((c) => c.name === "direct_credit")?.ok}
        >
          Research privately
        </button>
        {researchBusy ? (
          <button type="button" className="linkish" onClick={onCancel}>
            Stop
          </button>
        ) : null}
      </div>
      <ComputeCard checks={checks} onCheck={onCheck} />
      <p className="label" style={{ marginTop: 12 }}>
        Committee
      </p>
      {pollMiss ? <p role="status">Live view delayed — research is still running.</p> : null}
      {researchBusy ? (
        <p role="status">
          {shown.replaceAll("_", " ")} · {(researchElapsed / 1000).toFixed(1)}s elapsed. Live Direct round-trip, not a timer.
        </p>
      ) : null}
      <table className="desk-table">
        <thead>
          <tr>
            <th>Role</th>
            <th>Status</th>
            <th>Duration</th>
            <th>Reason</th>
          </tr>
        </thead>
        <tbody>
          {(["researcher", "challenger", "risk"] as const).map((name) => {
            const st = roleState(researchRoles, name, researchStage);
            const row = researchRoles.find((r) => String(r.role || "").toLowerCase() === name);
            const ms = row?.elapsed_ms || (researchBusy && st === "running" ? researchElapsed : 0);
            return (
              <tr key={name}>
                <td>{name}</td>
                <td className={`state ${st}`}>{namedRoleLabel(st, researchStop, researchKind)}</td>
                <td>{ms ? `${(ms / 1000).toFixed(1)}s` : "—"}</td>
                <td>{roleReason(researchRoles, name)}</td>
              </tr>
            );
          })}
          <tr>
            <td>TEE</td>
            <td className={`state ${teeState(researchRoles, researchBusy, researchStop)}`}>
              {namedRoleLabel(teeState(researchRoles, researchBusy, researchStop), researchStop, researchKind)}
            </td>
            <td>{researchBusy ? `${(researchElapsed / 1000).toFixed(1)}s` : "—"}</td>
            <td>Direct TeeML. Recovered signer must equal on-chain teeSigner.</td>
          </tr>
          <tr>
            <td>Engine</td>
            <td className={`state ${engineState(researchStage, researchRoles, researchBusy)}`}>
              {namedRoleLabel(engineState(researchStage, researchRoles, researchBusy), researchStop, researchKind)}
            </td>
            <td>{researchBusy && roleVerified(researchRoles, "risk") ? `${(researchElapsed / 1000).toFixed(1)}s` : "—"}</td>
            <td>Host sizes. Model cannot raise clip.</td>
          </tr>
          <tr>
            <td>Policy</td>
            <td className={`state ${preview ? (preview.eligible ? "success" : "blocked") : "pending"}`}>
              {preview ? (preview.eligible ? "PASS" : "BLOCKED") : researchBusy ? "PENDING" : "PENDING"}
            </td>
            <td>—</td>
            <td>Host policy. The model cannot mutate it.</td>
          </tr>
          <tr>
            <td>Decision</td>
            <td className={`state ${preview ? (preview.eligible ? "success" : "stand-down") : "pending"}`}>
              {preview ? (preview.eligible ? "EXACT PREVIEW" : "STOOD DOWN") : "PENDING"}
            </td>
            <td>—</td>
            <td>{preview?.eligible ? "Exact preview after verification." : "A stand-down is a successful research outcome."}</td>
          </tr>
        </tbody>
      </table>
      {!researchBusy && (explained || researchKind) ? (
        <section className="why-banner" role="status">
          <p className="label">{title}</p>
          <h2>{explained?.title || title}</h2>
          <p>{explained?.body || researchNote}</p>
          {researchJobId ? <p className="fine">Job {researchJobId}</p> : null}
          {explainStopHref(whyCode) ? (
            <ExternalLink className="linkish" href={explainStopHref(whyCode)?.href || ""}>
              {explainStopHref(whyCode)?.label}
            </ExternalLink>
          ) : null}
        </section>
      ) : null}
      {!researchBusy && whyRows.length ? (
        <section className="why-list" aria-label="Why">
          {whyRows.map((row) => (
            <div key={row.q}>
              <p className="label">{row.q}</p>
              <p>{row.a}</p>
            </div>
          ))}
        </section>
      ) : null}
      {!researchBusy && preview ? (
        <PreviewContract
          preview={preview}
          previewHash={previewHash}
          sessionAlive={sessionAlive}
          authTyped={authTyped}
          setAuthTyped={setAuthTyped}
          authBusy={authBusy}
          authErr={authErr}
          lastOid={lastOid}
          status={status}
          net={net}
          onAuthorize={onAuthorize}
          onCancelBound={onCancelBound}
        />
      ) : null}
      {!researchBusy && !researchNote && !explained && !preview ? (
        <p className="fine">Private research has not been run on this machine in this session.</p>
      ) : null}
      <PreviewNote />
      <EvidenceDrawer
        jobId={researchJobId}
        roles={researchRoles}
        preview={preview}
        previewHash={previewHash}
        kind={researchKind}
        deny={preview?.deny}
        evidence={researchEvidenceText}
        open={techOpen}
        onToggle={() => setTechOpen((v) => !v)}
      />
    </main>
  );
}

export function PreviewContract({
  preview,
  previewHash,
  sessionAlive,
  authTyped,
  setAuthTyped,
  authBusy,
  authErr,
  lastOid,
  status,
  net,
  onAuthorize,
  onCancelBound,
}: {
  preview: NonNullable<BindResult["preview"]>;
  previewHash: string;
  sessionAlive: boolean;
  authTyped: string;
  setAuthTyped: (v: string) => void;
  authBusy: boolean;
  authErr: string | null;
  lastOid: string;
  status: LocalStatus | null;
  net: string;
  onAuthorize: (e: FormEvent) => void;
  onCancelBound: (e: FormEvent) => void;
}) {
  if (!preview.eligible) {
    return (
      <article className="contract">
        <p className="label">Exact preview</p>
        <p>
          {preview.deny === "no_side"
            ? "Committee stood down. The committee did not propose a side. This is a verified result, not a crash. No order was placed."
            : `Host did not size a trade (${preview.deny || "no_side"}). The model cannot raise clip. No order was placed.`}
        </p>
      </article>
    );
  }
  return (
    <article className="contract">
      <p className="label">Exact preview</p>
      <h2>
        {preview.market} {preview.side}
      </h2>
      {preview.kind === "connection_test" ? (
        <p className="fine">Connection test. This is not a research recommendation. Host sized a policy clip.</p>
      ) : preview.kind === "reduce_only_close" ? (
        <p className="fine">Reduce-only close. Type AUTHORIZE to send it. PIT cannot withdraw.</p>
      ) : (
        <p className="fine">This is the exact thing that will happen. If anything changes, this authorization no longer applies.</p>
      )}
      <table className="kv">
        <tbody>
          <tr>
            <th>Asset</th>
            <td>{preview.market}</td>
          </tr>
          <tr>
            <th>Side</th>
            <td>{preview.side}</td>
          </tr>
          <tr>
            <th>Exact size</th>
            <td>
              {preview.sz} {preview.market}
            </td>
          </tr>
          <tr>
            <th>Venue</th>
            <td>Hyperliquid</td>
          </tr>
          <tr>
            <th>Price</th>
            <td>{preview.limitPx}</td>
          </tr>
          <tr>
            <th>Policy</th>
            <td>{preview.reasons?.length ? preview.reasons.join(" · ") : "Host clip and venue minimum already applied."}</td>
          </tr>
          <tr>
            <th>Session</th>
            <td>PIT Agent · {sessionAlive ? "Active" : "None"}</td>
          </tr>
          <tr>
            <th>Estimated cost</th>
            <td>Trading capital at the venue after AUTHORIZE. Private compute was already spent on research.</td>
          </tr>
          <tr>
            <th>Preview hash</th>
            <td className="hash">{preview.hash || previewHash}</td>
          </tr>
          <tr>
            <th>Will happen</th>
            <td>PIT will send this exact order after you type AUTHORIZE on this card.</td>
          </tr>
          <tr>
            <th>Will not happen</th>
            <td>Chat cannot authorize. Size cannot change. Withdraw, transfer, and leverage stay denied.</td>
          </tr>
        </tbody>
      </table>
      <p className="fine">AUTHORIZE approves only this exact preview. Any mutation invalidates it.</p>
      {sessionAlive ? (
        <form onSubmit={onAuthorize}>
          <input
            aria-label="type AUTHORIZE"
            autoComplete="off"
            value={authTyped}
            onChange={(ev) => setAuthTyped(ev.target.value)}
            placeholder="Type AUTHORIZE"
          />
          <button type="submit" disabled={authBusy || !previewHash}>
            AUTHORIZE
          </button>
        </form>
      ) : (
        <p>Create a local session, then type AUTHORIZE here.</p>
      )}
      {oidBelongsToPreview(status?.lastOrder?.hash, previewHash, preview.hash) && lastOid ? (
        status?.lastOrder?.status !== "filled" && !status?.lastOrder?.cancelled ? (
          <form onSubmit={onCancelBound}>
            <p className="fine">OID {lastOid} is resting for this preview. Type AUTHORIZE again to cancel. PIT cannot withdraw.</p>
            <button type="submit" disabled={authBusy}>
              Cancel this order
            </button>
          </form>
        ) : status?.lastOrder?.status === "filled" ? (
          <p className="fine">
            OID {lastOid} FILLED for this preview. Flatten only with a reduce-only close that YOU authorize. PIT cannot withdraw.
          </p>
        ) : null
      ) : null}
      {authErr ? (
        <p className="err" role="alert">
          {authErr === "approveAgent_required"
            ? "Approve this agent on Hyperliquid before PIT will send an order."
            : authErr}
        </p>
      ) : null}
      {authErr === "approveAgent_required" ? (
        <ExternalLink className="linkish" href={hyperliquidAPI(net)}>
          Open Hyperliquid API
        </ExternalLink>
      ) : null}
    </article>
  );
}
