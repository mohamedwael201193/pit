import { useMemo, useState } from "react";
import { ExternalLink } from "./ExternalLink";
import { compactUsd, pctFunding, powerSourceLabel } from "./format";
import { committeeVerified, oidBelongsToPreview, researchCardTitle } from "./honesty";
import { explainStop } from "./explain";
import { researchWhyCopy } from "./researchWhy";
import { LINKS } from "./links";
import type { ActivityEvent, BindResult, DirectModel } from "./companion";
import type { MarketCoin } from "./WatchBook";

export const CHAT_AGENT_COPY = {
  cannotAuthorize: "The model cannot AUTHORIZE",
  acceptOnDesk: "TRADE NOW on this computer submits the exact preview through the host.",
};

const PIPE = [
  { id: "DISCOVERY", label: "Universe" },
  { id: "SEALING_PRIVATE_BOOK", label: "Private book" },
  { id: "RESEARCHER", label: "Researcher" },
  { id: "CHALLENGER", label: "Challenger" },
  { id: "RISK", label: "Risk" },
  { id: "VERIFYING_TEE_SIGNATURE", label: "TEE" },
  { id: "DETERMINISTIC_ENGINE", label: "Host engine" },
  { id: "POLICY", label: "Policy" },
  { id: "PREVIEW", label: "Decision" },
] as const;

type Role = {
  role?: string;
  verify_e2ee?: string;
  proposed_side?: string;
  survives?: boolean;
  kill?: boolean;
};

type LastOrder = NonNullable<import("./companion").LocalStatus["lastOrder"]>;

function roleOk(roles: Role[], name: string) {
  return roles.some((r) => String(r.role || "").toLowerCase() === name && String(r.verify_e2ee || "").toUpperCase() === "OK");
}

function roleOf(roles: Role[], name: string) {
  return roles.find((r) => String(r.role || "").toLowerCase() === name);
}

function stageIndex(stage: string, roles: Role[], busy: boolean) {
  const s = (stage || "").toUpperCase();
  if (!busy && !s) return -1;
  if (s === "READING_MARKET" || s === "DISCOVERY") return 0;
  if (s.includes("SEAL") || s.includes("CONTACT") || s.includes("RECEIVING")) return 1;
  if (s === "RESEARCHER" || (roleOk(roles, "researcher") && !roleOk(roles, "challenger") && busy)) return 2;
  if (s === "CHALLENGER") return 3;
  if (s === "RISK" || s.startsWith("RISK_")) return 4;
  if (s.includes("TEE") || s.includes("VERIFY")) return 5;
  if (s === "DETERMINISTIC_ENGINE") return 6;
  if (s === "POLICY") return 7;
  if (s === "PREVIEW" || s === "READY") return 8;
  if (roleOk(roles, "risk")) return 6;
  if (roleOk(roles, "challenger")) return 4;
  if (roleOk(roles, "researcher")) return 3;
  return busy ? 0 : -1;
}

function pipeState(i: number, current: number, done: boolean, failed: boolean) {
  if (failed && i === current) return "fail";
  if (done || i < current) return "done";
  if (i === current) return "on";
  return "";
}

function orderKind(order?: LastOrder | null) {
  if (!order?.oid) return "";
  const life = String(order.lifecycle || order.status || "").toLowerCase();
  if (order.cancelled) return "cancelled";
  if (life.includes("fail") || life.includes("reject")) return "failed";
  if (life.includes("fill")) return "filled";
  if (life.includes("rest") || life.includes("open") || order.posted) return "resting";
  return String(order.status || "submitted");
}

function ageLabel(ms: number | undefined) {
  if (ms == null || !Number.isFinite(ms) || ms < 0) return "";
  if (ms < 2000) return "just now";
  if (ms < 60000) return `${Math.round(ms / 1000)}s ago`;
  return `${Math.round(ms / 60000)}m ago`;
}

function shortHash(v?: string) {
  const s = String(v || "");
  if (s.length < 12) return s || "—";
  return `${s.slice(0, 8)}…${s.slice(-4)}`;
}

function mark(ok: boolean) {
  return ok ? "✓" : "○";
}

export function AgentRun({
  busy,
  coin,
  stage,
  elapsedMs,
  jobId,
  pollMiss,
  updatedAt,
  roles,
  kind,
  coins,
  scanned,
  buyingPower,
  powerSource,
  watchAgeMs,
  bestWhy,
  preview,
  previewHash,
  lastOrder,
  lastOid,
  activity,
  evidence,
  researchNote,
  researchStop,
  huntRejected,
  huntSurvived,
  pinned,
  sessionAlive,
  autonomy,
  researchSku,
  authBusy,
  authErr,
  onAsk,
  onOpenPreview,
  onOpenPolicy,
  onOpenAutomation,
  onOpenActivity,
  onStop,
  onTradeNow,
}: {
  busy: boolean;
  coin: string;
  stage: string;
  elapsedMs: number;
  jobId: string;
  pollMiss: boolean;
  updatedAt?: number;
  roles: Role[];
  kind: string;
  coins: MarketCoin[];
  scanned: number;
  buyingPower?: number;
  powerSource?: string;
  watchAgeMs?: number;
  bestWhy?: string;
  preview: BindResult["preview"] | null;
  previewHash?: string;
  lastOrder?: LastOrder | null;
  lastOid?: string;
  activity: ActivityEvent[];
  evidence?: unknown;
  researchNote?: string | null;
  researchStop?: string | null;
  huntRejected: string[];
  huntSurvived: string;
  pinned: boolean;
  sessionAlive: boolean;
  autonomy: string;
  researchSku?: DirectModel | null;
  authBusy?: boolean;
  authErr?: string | null;
  onAsk: (q: string) => void;
  onOpenPreview: () => void;
  onOpenPolicy: () => void;
  onOpenAutomation: () => void;
  onOpenActivity: () => void;
  onStop: () => void;
  onTradeNow: () => void;
}) {
  const [openStep, setOpenStep] = useState<string | null>(null);
  const [tech, setTech] = useState(false);
  const [pick, setPick] = useState("");
  const [whyOpen, setWhy] = useState(false);
  const [sleepOpen, setSleep] = useState(false);
  const verified = committeeVerified(roles);
  const current = stageIndex(stage, roles, busy);
  const executable = coins.filter((c) => c.executionFeasible);
  const eligible = coins.filter((c) => c.eligible || c.policyEligible);
  const scannedN = scanned || coins.length;
  const best = executable[0] || eligible[0] || coins[0];
  const focus = coins.find((c) => c.coin === (pick || coin)) || best;
  const title = researchCardTitle(kind, verified);
  const deny = String(preview?.deny || "");
  const policyBlock = kind === "POLICY_DENIED" || deny.includes("policy");
  const capitalBlock = kind === "MARKET_DENIED" || deny.includes("min_notional") || deny.includes("margin");
  const noTrade = kind === "READY_STOOD_DOWN" || deny === "no_side";
  const ready = kind === "READY_ELIGIBLE" && Boolean(preview?.eligible && (preview.hash || previewHash));
  const fail = Boolean(researchStop) && !noTrade && !ready;
  const stop = explainStop(researchStop || "");
  const researcher = roleOf(roles, "researcher");
  const challenger = roleOf(roles, "challenger");
  const risk = roleOf(roles, "risk");
  const side = String(preview?.side || researcher?.proposed_side || "").toUpperCase();
  const hash = String(previewHash || preview?.hash || "");
  const alreadyPosted = Boolean(lastOrder?.posted && lastOrder.hash && hash && lastOrder.hash === hash);
  const canTrade = ready && Boolean(hash) && sessionAlive && pinned && !authBusy && !alreadyPosted;
  const proof = useMemo(() => {
    const fromAct = activity.find((e) => (e.job_id && e.job_id === jobId) || e.root || e.tx);
    const ev = evidence && typeof evidence === "object" ? (evidence as Record<string, unknown>) : null;
    const root = String(fromAct?.root || ev?.root || ev?.storage_root || "");
    const tx = String(fromAct?.tx || ev?.tx || "");
    const txLink = String(fromAct?.tx_link || ev?.tx_link || "");
    const digest = String(fromAct?.digest || ev?.digest || "");
    if (!root && !tx && !jobId) return null;
    return { root, tx, txLink, digest, jobId };
  }, [activity, evidence, jobId]);
  const orderState = oidBelongsToPreview(lastOrder?.hash, previewHash, hash) ? orderKind(lastOrder) : "";
  const elapsed = elapsedMs > 0 ? `${(elapsedMs / 1000).toFixed(1)}s` : "";
  const beat = ageLabel(updatedAt ? Date.now() - updatedAt : undefined);
  const why = researchWhyCopy({
    coin: coin || focus?.coin || "",
    kind,
    note: researchNote,
    stop: researchStop,
    deny,
    eligible: Boolean(preview?.eligible),
    roles,
    snap: focus ? { mark: focus.mark, reason: focus.reason, why: focus.why } : undefined,
  });
  const status = busy ? "RESEARCHING" : ready ? "READY" : noTrade ? "NO-TRADE" : policyBlock ? "BLOCKED" : fail ? "STOPPED" : verified ? "VERIFIED" : "LIVE";

  const follow = ready
    ? [
        ["Review details", "Prepare the exact trade"],
        ["Compare candidates", "Compare top opportunities"],
        ["Why this book?", "__why"],
        ["Sleep Mission", "Trade this while I sleep"],
      ]
    : noTrade
      ? [
          ["Research next", "Find the best opportunity"],
          ["Show why", "__why"],
          ["Show alternatives", "Compare top opportunities"],
          ["Sleep Mission", "__sleep"],
        ]
      : policyBlock
        ? [
            ["Open Policy", "Explain my policy"],
            ["What is executable?", "What is executable?"],
          ]
        : capitalBlock
          ? [
              ["What can I trade?", "What can I trade now?"],
              ["Scan again", "Scan all markets"],
            ]
          : busy
            ? [
                ["Why this book?", "__why"],
                ["Stop research", "Stop research"],
              ]
            : [
                ["Find best", "Find the best opportunity"],
                ["What can I trade?", "What can I trade now?"],
                ["Compare", "Compare top opportunities"],
              ];

  return (
    <section className="cockpit" aria-label="PIT agent live">
      <header className="cockpit-status">
        <div>
          <p className="cockpit-kicker">PIT AGENT</p>
          <h2 className="cockpit-title">Live market → private book → 0G Direct → Hyperliquid</h2>
        </div>
        <p className={`cockpit-pill ${busy ? "on" : ready ? "ok" : ""}`}>{status}</p>
      </header>

      <ul className="cockpit-facts">
        <li>Policy {pinned ? "pinned" : "draft"}</li>
        <li>Session {sessionAlive ? "live" : "off"}</li>
        <li>{compactUsd(buyingPower)} {powerSourceLabel(powerSource)}</li>
        <li>{autonomy || "manual"}</li>
        <li>Books {ageLabel(watchAgeMs) || "live"}</li>
        {beat ? <li>Job {beat}</li> : null}
      </ul>

      {!busy ? (
        <div className="cockpit-funnel" aria-label="scan funnel">
          <div><span>{scannedN || "—"}</span><em>Scanned</em></div>
          <div><span>{eligible.length || "—"}</span><em>Policy</em></div>
          <div><span>{executable.length || "—"}</span><em>Executable</em></div>
          <div><span>{huntRejected.length || (noTrade ? 1 : 0)}</span><em>Rejected</em></div>
          <div><span>{huntSurvived || best?.coin || "—"}</span><em>Best</em></div>
          <div><span>{ready ? "READY" : noTrade ? "NO-TRADE" : status}</span><em>Result</em></div>
        </div>
      ) : (
        <p className="cockpit-scanline">
          Scanning {scannedN || "live"} markets · {eligible.length} policy · {executable.length} executable · {coin || best?.coin || "selecting"}
        </p>
      )}

      {focus ? (
        <article className="cockpit-quote">
          <p className="cockpit-kicker">{busy ? "BEST OPPORTUNITY" : ready ? "CANDIDATE" : "LIVE BOOK"}</p>
          <header>
            <h3>{focus.coin}</h3>
            <p>{compactUsd(focus.mark)}</p>
            <span>Hyperliquid</span>
          </header>
          <dl>
            <div><dt>Venue min</dt><dd>{compactUsd(focus.minNotional)}</dd></div>
            <div><dt>Host clip</dt><dd>{compactUsd(focus.hostNotional || focus.policyClip)}</dd></div>
            <div><dt>Funding</dt><dd>{pctFunding(focus.funding)}</dd></div>
          </dl>
          <p className="fine">Mark is the price. Venue min is the order notional. They are not the same number.</p>
          {bestWhy ? <p className="fine">{bestWhy}</p> : null}
        </article>
      ) : null}

      {!busy && (executable.length || eligible.length) ? (
        <div className="cockpit-strip" role="list">
          {(executable.length ? executable : eligible).slice(0, 6).map((c) => (
            <button
              key={c.coin}
              type="button"
              role="listitem"
              className={c.coin === (pick || coin) ? "on" : ""}
              onClick={() => setPick(c.coin)}
            >
              <strong>{c.coin}</strong>
              <em>{c.executionFeasible ? "READY" : c.eligible ? "POLICY" : "WATCH"}</em>
            </button>
          ))}
        </div>
      ) : null}

      {busy || stage ? (
        <article className="cockpit-work" aria-live="polite">
          <header>
            <div>
              <p className="cockpit-kicker">PRIVATE 0G RESEARCH</p>
              <p className="cockpit-work-title">{coin || focus?.coin || "—"} · {elapsed}{pollMiss ? " · reconnecting" : ""}{beat ? ` · ${beat}` : ""}</p>
            </div>
            {busy ? (
              <button type="button" className="ghost" onClick={onStop}>Stop</button>
            ) : null}
          </header>
          <ol className="cockpit-pipe">
            {PIPE.map((step, i) => {
              const markState = pipeState(i, current, verified || ready || noTrade, fail);
              return (
                <li key={step.id} className={markState}>
                  <button type="button" onClick={() => setOpenStep(openStep === step.id ? null : step.id)}>
                    <em aria-hidden>{markState === "done" ? "✓" : markState === "on" ? "●" : markState === "fail" ? "✕" : "○"}</em>
                    <strong>{step.label}</strong>
                    <span>{markState === "on" ? (stage || "live") : ""}</span>
                  </button>
                  {openStep === step.id ? (
                    <p className="fine">
                      {step.id === "RESEARCHER" ? (researcher ? `verify ${researcher.verify_e2ee || "—"} · ${researcher.proposed_side || "no side"}` : "Waiting for researcher.") : null}
                      {step.id === "CHALLENGER" ? (challenger ? `verify ${challenger.verify_e2ee || "—"} · survives ${String(challenger.survives)}` : "Waiting for challenger.") : null}
                      {step.id === "RISK" ? (risk ? `verify ${risk.verify_e2ee || "—"}` : "Waiting for risk.") : null}
                      {step.id === "PREVIEW" ? (hash ? `Preview ${shortHash(hash)}` : "No exact preview yet.") : null}
                      {step.id === "POLICY" ? (preview?.eligible ? "Policy pass on this computer." : deny || "Policy has not sized a preview.") : null}
                      {step.id === "DISCOVERY" ? `${scannedN} live books. ${executable.length} executable.` : null}
                      {step.id === "VERIFYING_TEE_SIGNATURE" ? (verified ? "VerifyE2EE OK on named roles." : "TEE is not claimed until VerifyE2EE succeeds.") : null}
                    </p>
                  ) : null}
                </li>
              );
            })}
          </ol>
        </article>
      ) : null}

      {ready && preview ? (
        <article className="cockpit-card ready">
          <p className="cockpit-kicker">READY TO TRADE</p>
          <h3>{preview.market || coin} · {side || "SIDE"}</h3>
          <p className="cockpit-num">{compactUsd(preview.notionalUsd)}</p>
          <p>1x leverage · Hyperliquid · nothing executed yet.</p>
          <ul className="cockpit-checks">
            <li>asset allowed</li>
            <li>clip allowed</li>
            <li>leverage allowed</li>
            <li>position available</li>
            <li>venue minimum satisfied</li>
            <li>slippage allowed</li>
          </ul>
          <p className="fine">Account {compactUsd(buyingPower)} · preview {shortHash(hash)}</p>
          {!sessionAlive ? <p className="fine">Create a live session on Security before TRADE NOW.</p> : null}
          {!pinned ? <p className="fine">Pin policy on this computer first.</p> : null}
          {alreadyPosted ? <p className="fine">This preview already produced OID {lastOrder?.oid}. PIT will not double-submit.</p> : null}
          {authErr ? <p className="fine" role="alert">{authErr}</p> : null}
          <div className="cta-row">
            <button
              type="button"
              className="primary"
              aria-label="TRADE NOW"
              disabled={!canTrade}
              onClick={onTradeNow}
            >
              {authBusy ? "Submitting…" : "TRADE NOW"}
            </button>
            <button type="button" className="ghost" onClick={() => onAsk("Do not trade")}>REJECT</button>
            <button type="button" className="ghost" onClick={onOpenPreview}>REVIEW DETAILS</button>
          </div>
          <p className="fine">{CHAT_AGENT_COPY.cannotAuthorize}. {CHAT_AGENT_COPY.acceptOnDesk}</p>
        </article>
      ) : null}

      {noTrade ? (
        <article className="cockpit-card stand">
          <p className="cockpit-kicker">NO TRADE, VERIFIED</p>
          <h3>{coin || focus?.coin || "Candidate"}</h3>
          <p>{title}. Nothing was executed. This is a valid result.</p>
          <ul className="cockpit-checks">
            <li>Researcher {researcher?.proposed_side ? `proposed ${researcher.proposed_side}` : "no candidate proposed"}</li>
            <li>Challenger {challenger?.survives === false || challenger?.kill ? "no surviving side" : challenger ? "checked" : "not required"}</li>
            <li>Risk {risk ? "checked" : "not required"}</li>
            <li>Execution no order created</li>
          </ul>
          {huntRejected.length ? (
            <p className="fine">
              {scannedN} scanned · {executable.length} executable · {huntRejected.length} researched/rejected
              {huntSurvived ? ` · survived ${huntSurvived}` : ""}.
            </p>
          ) : null}
        </article>
      ) : null}

      {policyBlock ? (
        <article className="cockpit-card block">
          <p className="cockpit-kicker">POLICY BLOCK</p>
          <h3>{title}</h3>
          <p>{deny || "Host law refused this size. The model cannot raise clip or leverage."}</p>
          <div className="cta-row">
            <button type="button" className="primary" onClick={onOpenPolicy}>Open Policy</button>
          </div>
        </article>
      ) : null}

      {capitalBlock ? (
        <article className="cockpit-card">
          <p className="cockpit-kicker">INSUFFICIENT CAPITAL</p>
          <h3>{title}</h3>
          <p>Buying power {compactUsd(buyingPower)}. Venue min {compactUsd(focus?.minNotional)}. PIT will not invent size.</p>
        </article>
      ) : null}

      {fail && stop ? (
        <article className="cockpit-card">
          <p className="cockpit-kicker">STOPPED</p>
          <h3>{stop.title}</h3>
          <p>{stop.body}</p>
        </article>
      ) : null}

      {proof && (proof.root || proof.tx) && !busy ? (
        <article className="cockpit-card og">
          <p className="cockpit-kicker">0G PROOF</p>
          <h3>{verified ? "TEE verified" : "Not claimed verified"}</h3>
          {proof.jobId ? <p>Job {shortHash(proof.jobId)}</p> : null}
          {proof.root ? <p>Root {shortHash(proof.root)}</p> : <p className="fine">No storage root on this computer yet.</p>}
          {proof.tx ? <p>Chain {shortHash(proof.tx)}</p> : <p className="fine">No chain transaction until a root is filed.</p>}
          <div className="cta-row">
            {proof.txLink ? <ExternalLink href={proof.txLink}>Open 0G explorer</ExternalLink> : null}
            {proof.root ? <button type="button" className="ghost" onClick={onOpenActivity}>Verify on 0G</button> : null}
          </div>
        </article>
      ) : null}

      {orderState && lastOrder ? (
        <article className={`cockpit-card ${orderState === "filled" ? "ready" : ""}`}>
          <p className="cockpit-kicker">
            {orderState === "filled" ? "FILLED" : orderState === "resting" ? "ORDER SUBMITTED" : orderState === "cancelled" ? "CANCELLED" : orderState === "failed" ? "FAILED" : "ORDER"}
          </p>
          <h3>{lastOrder?.market || "—"} {String(lastOrder?.side || "").toUpperCase()}</h3>
          <p>OID {lastOrder?.oid || lastOid}</p>
          <p>Status {String(lastOrder?.lifecycle || lastOrder?.status || orderState).toUpperCase()}</p>
          {orderState === "resting" ? <p className="fine">RESTING is not a fill.</p> : null}
          {orderState === "failed" ? <p className="fine">No execution occurred.</p> : null}
          <div className="cta-row">
            <ExternalLink href={LINKS.hl}>Open Hyperliquid</ExternalLink>
            <button type="button" className="ghost" onClick={onOpenActivity}>Open Activity</button>
            {proof?.root || proof?.tx ? <button type="button" className="ghost" onClick={onOpenActivity}>Verify on 0G</button> : null}
          </div>
        </article>
      ) : null}

      <div className="cockpit-follow">
        {follow.map(([label, q]) => (
          <button
            key={label}
            type="button"
            className="chip-btn"
            onClick={() => {
              if (q === "Stop research") {
                onStop();
                return;
              }
              if (q === "__why") {
                setWhy(true);
                return;
              }
              if (q === "__sleep" || q === "Review Sleep Mission") {
                setSleep(true);
                return;
              }
              onAsk(q);
            }}
          >
            {label}
          </button>
        ))}
      </div>

      {whyOpen ? (
        <article className="cockpit-card">
          <p className="cockpit-kicker">WHY</p>
          {why.map((row) => (
            <p key={row.q} className="fine"><strong>{row.q}</strong> {row.a}</p>
          ))}
        </article>
      ) : null}

      {sleepOpen && !busy ? (
        <article className="cockpit-sleep">
          <div>
            <p className="cockpit-kicker">SLEEP MISSION</p>
            <p>Bound by current policy. Research, challenge, risk, policy, then execute. This Agent may prepare it. Desktop still arms.</p>
          </div>
          <button type="button" className="ghost" onClick={onOpenAutomation}>Open Sleep Mission</button>
        </article>
      ) : null}

      <details className="cockpit-tech" open={tech} onToggle={(e) => setTech((e.target as HTMLDetailsElement).open)}>
        <summary>Technical details</summary>
        <p>Direct TeeML · {researchSku?.model || "not this chat stream"} · {researchSku?.proven_e2ee ? "sealed" : "privacy not claimed here"}</p>
        <p>Job {jobId || "—"} · Kind {kind || "—"} · Stage {stage || "—"}</p>
        <p>Committee {mark(roleOk(roles, "researcher"))} researcher {mark(roleOk(roles, "challenger"))} challenger {mark(roleOk(roles, "risk"))} risk · TEE {verified ? "VerifyE2EE" : "not claimed"}</p>
        <p>Private book never uses Router. Catalog listing is not an inference path.</p>
      </details>
    </section>
  );
}
