import { FormEvent } from "react";
import { ComputeCard } from "./ComputeCard";
import { EvidenceDrawer } from "./EvidenceDrawer";
import { PreviewNote } from "./PreviewNote";
import { committeeVerified, oidBelongsToPreview, researchCardTitle } from "./honesty";
import { explainStop, explainStopHref } from "./explain";
import { hyperliquidAPI } from "./links";
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

function roleState(roles: ResearchRole[], name: string, stage: string) {
  const mark = stageMark(name.toUpperCase(), stage, roles);
  if (roleVerified(roles, name)) return "success";
  if (mark === "lit") return "running";
  if (mark === "done") return "success";
  return "pending";
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
  eligible: Array<{ coin: string; mark: number; reason: string }>;
  net: string;
  onResearch: (coin?: string) => void;
  onCancel: () => void;
  onAuthorize: (e: FormEvent) => void;
  onCancelBound: (e: FormEvent) => void;
  onCheck: () => void;
}) {
  const shown = canonicalResearchStage(researchStage, researchRoles);
  const explained = explainStop(researchStop);
  const verified = committeeVerified(researchRoles);
  return (
    <main className="page dense">
      <p className="eyebrow">Research</p>
      <h1>{coin || "ETH"}</h1>
      <p className="lead">Private committee evaluation. Host sizes. Chat cannot AUTHORIZE.</p>
      <article className="card">
        <p className="label">Thesis</p>
        <p className="fine">Sealed into the private book. The committee may still stand down.</p>
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
        </div>
      </article>
      <ComputeCard checks={checks} onCheck={onCheck} />
      <div className="cta-row">
        <button
          type="button"
          className="primary"
          onClick={() => onResearch(coin)}
          disabled={researchBusy || !checks.find((c) => c.name === "direct_credit")?.ok}
        >
          Start research
        </button>
      </div>
      {researchBusy ? (
        <article className="card" role="status">
          <p className="label">Live sealed request</p>
          <p>
            {shown.replaceAll("_", " ")} · {coin} · {(researchElapsed / 1000).toFixed(1)}s elapsed. This is a live Direct
            round-trip, not a timer.
          </p>
          {pollMiss ? <p role="status">Live view delayed — research is still running.</p> : null}
          <div className="committee">
            {(["researcher", "challenger", "risk"] as const).map((name) => (
              <div key={name} className="role-card">
                <p className="label">{name}</p>
                <p className="state">{roleState(researchRoles, name, researchStage)}</p>
                <p className="fine">{(researchElapsed / 1000).toFixed(1)}s</p>
              </div>
            ))}
            <div className="role-card">
              <p className="label">TEE</p>
              <p className="state">{roleVerified(researchRoles, "researcher") ? "signer checked" : "pending"}</p>
              <p className="fine">Direct · glm-5.2</p>
            </div>
          </div>
          <ol className="pipe stages">
            {RESEARCH_STAGES.map((name) => {
              const mark = stageMark(name, researchStage, researchRoles);
              const prefix = mark === "done" ? "✓ " : mark === "lit" ? "● " : "○ ";
              return (
                <li key={name} className={mark}>
                  {prefix}
                  {name.replaceAll("_", " ")}
                </li>
              );
            })}
          </ol>
          <button type="button" className="linkish" onClick={onCancel}>
            Cancel
          </button>
        </article>
      ) : null}
      {researchNote && !explained && !researchBusy ? (
        <article className="card">
          <p className="label">{researchCardTitle(researchKind, verified)}</p>
          <p>{researchNote}</p>
          {researchJobId ? <p className="fine">Job {researchJobId}</p> : null}
          <ul className="doctor">
            {researchRoles.map((role) => (
              <li key={role.role}>
                <strong>{role.verify_e2ee}</strong> {role.role}
                {role.proposed_side ? ` side ${role.proposed_side}` : ""}
                {role.survives === false ? " stood down" : ""}
                {role.kill ? " kill" : ""}
              </li>
            ))}
          </ul>
        </article>
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
      {explained && !researchBusy ? (
        <article className="card stop" role="alert">
          <p className="label">{researchCardTitle(researchKind || researchStop, verified)}</p>
          <h2>{explained.title}</h2>
          <p>{explained.body}</p>
          {explainStopHref(researchStop) ? (
            <a className="linkish" href={explainStopHref(researchStop)?.href} target="_blank" rel="noreferrer">
              {explainStopHref(researchStop)?.label}
            </a>
          ) : null}
          <button type="button" className="linkish" onClick={() => setTechOpen((v) => !v)}>
            {techOpen ? "Hide technical evidence" : "View technical evidence"}
          </button>
          {techOpen ? (
            <pre className="pipe evidence">
              Code {researchStop}
              {"\n"}
              {researchEvidenceText || "Verification is fail-closed. Router fallback is impossible."}
            </pre>
          ) : null}
          <button type="button" onClick={() => onResearch(coin)} disabled={researchBusy}>
            Retry
          </button>
        </article>
      ) : null}
      {!researchBusy && !researchNote && !explained ? (
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
      />
      {eligible.length ? (
        <div className="cta-row">
          {eligible.map((c) => (
            <button key={c.coin} type="button" className="linkish" onClick={() => onResearch(c.coin)}>
              Research {c.coin}
            </button>
          ))}
        </div>
      ) : (
        <p className="fine">No policy-eligible market is waiting. Watch does not invent cards.</p>
      )}
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
            <th>Size</th>
            <td>
              {preview.sz} {preview.market}
            </td>
          </tr>
          <tr>
            <th>Limit</th>
            <td>{preview.limitPx}</td>
          </tr>
          <tr>
            <th>Venue</th>
            <td>Hyperliquid</td>
          </tr>
          <tr>
            <th>Session</th>
            <td>PIT Agent · {sessionAlive ? "Active" : "None"}</td>
          </tr>
          <tr>
            <th>Compute</th>
            <td>Verified private research</td>
          </tr>
          <tr>
            <th>Hash</th>
            <td className="hash">{preview.hash || previewHash}</td>
          </tr>
        </tbody>
      </table>
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
        <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
          Open Hyperliquid API
        </a>
      ) : null}
    </article>
  );
}
