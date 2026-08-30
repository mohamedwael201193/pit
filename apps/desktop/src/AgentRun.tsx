import { useMemo, useState } from "react";
import { ExternalLink } from "./ExternalLink";
import { compactNum, compactUsd, pctFunding, powerSourceLabel } from "./format";
import { committeeVerified, researchCardTitle } from "./honesty";
import { explainStop } from "./explain";
import { researchWhyCopy } from "./researchWhy";
import { LINKS } from "./links";
import type { ActivityEvent, BindResult, DirectModel } from "./companion";
import type { MarketCoin } from "./WatchBook";

export const CHAT_AGENT_COPY = {
  cannotAuthorize: "Chat cannot AUTHORIZE",
  acceptOnDesk: "Review the exact preview on this computer, then type AUTHORIZE there.",
};

const PIPE = [
  { id: "DISCOVERY", label: "Discovery" },
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

export type AgentView = {
  kind?: string;
  executive?: string;
  scanned?: number;
  eligible?: number;
  executable?: number;
  researched?: number;
  rejected?: number;
  best?: string;
  why?: string;
  buying_power?: number;
  power_source?: string;
  min_notional?: number;
  host_notional?: number;
  mark?: number;
  funding?: number;
  open_interest?: number;
  freshness?: string;
  coins?: Array<{
    coin: string;
    mark?: number;
    min_notional?: number;
    host_notional?: number;
    funding?: number;
    open_interest?: number;
    execution_feasible?: boolean;
    eligible?: boolean;
    why?: string;
    exec_why?: string;
    block?: string;
    trend?: string;
  }>;
};

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
  if (s === "RESEARCHER" || roleOk(roles, "researcher") && !roleOk(roles, "challenger") && busy) return 2;
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
  if (ms < 2000) return "Updated just now";
  if (ms < 60000) return `Last updated ${Math.round(ms / 1000)}s ago`;
  return `Last updated ${Math.round(ms / 60000)}m ago`;
}

function shortHash(v?: string) {
  const s = String(v || "");
  if (s.length < 12) return s || "—";
  return `${s.slice(0, 8)}…${s.slice(-4)}`;
}

export function AgentRun({
  busy,
  coin,
  stage,
  elapsedMs,
  jobId,
  pollMiss,
  roles,
  kind,
  coins,
  scanned,
  buyingPower,
  powerSource,
  watchAgeMs,
  bestWhy,
  preview,
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
  onAsk,
  onOpenPreview,
  onOpenPolicy,
  onOpenAutomation,
  onOpenActivity,
  onStop,
}: {
  busy: boolean;
  coin: string;
  stage: string;
  elapsedMs: number;
  jobId: string;
  pollMiss: boolean;
  roles: Role[];
  kind: string;
  coins: MarketCoin[];
  scanned: number;
  buyingPower?: number;
  powerSource?: string;
  watchAgeMs?: number;
  bestWhy?: string;
  preview: BindResult["preview"] | null;
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
  onAsk: (q: string) => void;
  onOpenPreview: () => void;
  onOpenPolicy: () => void;
  onOpenAutomation: () => void;
  onOpenActivity: () => void;
  onStop: () => void;
}) {
  const [openStep, setOpenStep] = useState<string | null>(null);
  const [tech, setTech] = useState(false);
  const [pick, setPick] = useState("");
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
  const ready = kind === "READY_ELIGIBLE" && Boolean(preview?.eligible && preview.hash);
  const fail = Boolean(researchStop) && !noTrade && !ready;
  const stop = explainStop(researchStop || "");
  const researcher = roleOf(roles, "researcher");
  const challenger = roleOf(roles, "challenger");
  const risk = roleOf(roles, "risk");
  const side = String(preview?.side || researcher?.proposed_side || "").toUpperCase();
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
  const orderState = orderKind(lastOrder);
  const elapsed = elapsedMs > 0 ? `${(elapsedMs / 1000).toFixed(1)}s` : "";
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
  const status = busy ? "RUNNING" : ready ? "COMPLETE" : noTrade ? "NO-TRADE" : policyBlock ? "BLOCKED" : fail ? "REJECTED" : verified ? "VERIFIED" : "LIVE";

  const follow = ready
    ? [
        ["Prepare exact trade", "Prepare the exact trade"],
        ["Compare alternatives", "Compare top opportunities"],
        ["Why this book?", `Why ${coin || focus?.coin || "this"}?`],
        ["Set up Sleep Mission", "Review Sleep Mission"],
      ]
    : noTrade
      ? [
          ["Research next", "Find the best opportunity"],
          ["Show why", "Why didn't you trade?"],
          ["Show alternatives", "Compare top opportunities"],
          ["Watch automatically", "Review Sleep Mission"],
        ]
      : policyBlock
        ? [
            ["Open Policy", "Explain my policy"],
            ["What is executable?", "What is executable?"],
          ]
        : capitalBlock
          ? [
              ["Show minimum", "What can I trade now?"],
              ["Scan again", "Scan all markets"],
            ]
          : [
              ["Find best", "Find the best opportunity"],
              ["Scan all", "Scan all markets"],
              ["Why no trade?", "Why didn't you trade?"],
            ];

  return (
    <section className="cockpit" aria-label="PIT agent">
      <header className="cockpit-status">
        <div>
          <p className="cockpit-kicker">PIT AGENT</p>
          <h2 className="cockpit-title">LIVE · MAINNET · 0G PRIVATE RESEARCH · HYPERLIQUID</h2>
        </div>
        <p className={`cockpit-pill ${busy ? "on" : ready ? "ok" : ""}`}>{status}</p>
      </header>
      <ul className="cockpit-facts">
        <li>Policy {pinned ? "pinned" : "draft"}</li>
        <li>Session {sessionAlive ? "live" : "off"}</li>
        <li>Buying power {compactUsd(buyingPower)}</li>
        <li>Autonomy {autonomy || "manual"}</li>
        <li>{ageLabel(watchAgeMs) || (focus?.freshness ? `Book ${focus.freshness}` : "Polling live books")}</li>
      </ul>

      <div className="cockpit-funnel">
        <div><span>{scannedN || "—"}</span><em>Scanned</em></div>
        <div><span>{eligible.length || "—"}</span><em>Policy eligible</em></div>
        <div><span>{executable.length || "—"}</span><em>Executable</em></div>
        <div><span>{huntRejected.length || (noTrade ? 1 : 0)}</span><em>Rejected</em></div>
        <div><span>{huntSurvived || best?.coin || "—"}</span><em>Best</em></div>
        <div><span>{ready ? "READY" : noTrade ? "NO-TRADE" : busy ? "WORKING" : status}</span><em>Decision</em></div>
      </div>

      {executable.length || eligible.length ? (
        <div className="cockpit-strip" role="list">
          <p className="cockpit-meta">BEST NOW</p>
          <div>
            {(executable.length ? executable : eligible).slice(0, 8).map((c) => (
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
        </div>
      ) : null}

      {focus ? (
        <article className="cockpit-quote">
          <header>
            <h3>{focus.coin}</h3>
            <p>{compactUsd(focus.mark)}</p>
            <span>Hyperliquid</span>
          </header>
          <dl>
            <div><dt>Venue min</dt><dd>{compactUsd(focus.minNotional)}</dd></div>
            <div><dt>Host clip</dt><dd>{compactUsd(focus.hostNotional || focus.policyClip)}</dd></div>
            <div><dt>Funding</dt><dd>{pctFunding(focus.funding)}</dd></div>
            <div><dt>OI</dt><dd>{compactNum(focus.openInterest)}</dd></div>
          </dl>
          <p className="fine">Mark is the price. Venue min is the order notional. They are not the same number.</p>
          {bestWhy ? <p className="fine">{bestWhy}</p> : null}
        </article>
      ) : null}

      {busy || stage ? (
        <article className="cockpit-work" aria-live="polite">
          <header>
            <p className="cockpit-kicker">{busy ? "PIT is working" : "Committee"}</p>
            <p>{coin || focus?.coin || "—"} · {elapsed} {pollMiss ? " · checking connection" : ""}</p>
            {busy ? (
              <button type="button" className="ghost" onClick={onStop}>Stop</button>
            ) : null}
          </header>
          <ol className="cockpit-pipe">
            {PIPE.map((step, i) => {
              const mark = pipeState(i, current, verified || ready || noTrade, fail);
              return (
                <li key={step.id} className={mark}>
                  <button type="button" onClick={() => setOpenStep(openStep === step.id ? null : step.id)}>
                    <span>{String(i + 1).padStart(2, "0")}</span>
                    <strong>{step.label}</strong>
                    <em>{mark === "done" ? "✓" : mark === "on" ? "●" : mark === "fail" ? "✕" : "○"}</em>
                  </button>
                  {openStep === step.id ? (
                    <p className="fine">
                      {step.id === "RESEARCHER" ? (researcher ? `verify ${researcher.verify_e2ee || "—"} · ${researcher.proposed_side || "no side"}` : "Waiting for researcher.") : null}
                      {step.id === "CHALLENGER" ? (challenger ? `verify ${challenger.verify_e2ee || "—"} · survives ${String(challenger.survives)}` : "Waiting for challenger.") : null}
                      {step.id === "RISK" ? (risk ? `verify ${risk.verify_e2ee || "—"}` : "Waiting for risk.") : null}
                      {step.id === "PREVIEW" ? (preview?.hash ? `Preview hash ${shortHash(preview.hash)}` : "No exact preview yet.") : null}
                      {step.id === "POLICY" ? (preview?.eligible ? "Policy pass on this computer." : deny || "Policy has not sized a preview.") : null}
                      {step.id === "DISCOVERY" ? `${scannedN} live books. ${executable.length} executable.` : null}
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
          <p className="cockpit-kicker">READY TO AUTHORIZE</p>
          <h3>{preview.market || coin} {side}</h3>
          <p className="cockpit-num">{compactUsd(preview.notionalUsd)}</p>
          <p>1x · Hyperliquid · nothing has been executed yet.</p>
          <ul className="cockpit-checks">
            <li>Asset allowed</li>
            <li>Clip within host law</li>
            <li>Leverage 1x</li>
            <li>Venue minimum satisfied</li>
          </ul>
          <p className="fine">Account {compactUsd(buyingPower)} {powerSourceLabel(powerSource)}</p>
          <div className="cta-row">
            <button type="button" className="primary" onClick={onOpenPreview}>Open exact preview on desktop</button>
            <button type="button" className="ghost" onClick={() => onAsk("Do not trade")}>Do not trade</button>
          </div>
          <p className="fine">{CHAT_AGENT_COPY.cannotAuthorize}. {CHAT_AGENT_COPY.acceptOnDesk}</p>
        </article>
      ) : null}

      {noTrade ? (
        <article className="cockpit-card stand">
          <p className="cockpit-kicker">NO TRADE — VERIFIED</p>
          <h3>{title}</h3>
          <p>{coin || focus?.coin || "A ranked book"} was researched privately. Nothing was executed.</p>
          <ul className="cockpit-checks">
            <li>Researcher {researcher?.proposed_side ? `proposed ${researcher.proposed_side}` : "no candidate proposed"}</li>
            <li>Challenger {challenger?.survives === false || challenger?.kill ? "no surviving side" : challenger ? "checked" : "not required"}</li>
            <li>Risk {risk ? "checked" : "not required"}</li>
            <li>Policy not the blocker</li>
            <li>No order created</li>
          </ul>
          {huntRejected.length ? <p className="fine">{huntRejected.length} candidate{huntRejected.length === 1 ? "" : "s"} rejected{huntSurvived ? ` · survived ${huntSurvived}` : ""}.</p> : null}
          <div className="cta-row">
            <button type="button" className="primary" onClick={() => onAsk("Find the best opportunity")}>Research next opportunity</button>
            <button type="button" className="ghost" onClick={() => onAsk("Why didn't you trade?")}>Show why</button>
          </div>
        </article>
      ) : null}

      {policyBlock ? (
        <article className="cockpit-card block">
          <p className="cockpit-kicker">POLICY BLOCK</p>
          <h3>{title}</h3>
          <p>{deny || "Host law refused this size. Chat cannot raise clip or leverage."}</p>
          <div className="cta-row">
            <button type="button" className="primary" onClick={onOpenPolicy}>Open Policy</button>
            <button type="button" className="ghost" onClick={() => onAsk("What is executable?")}>Show required limit</button>
          </div>
        </article>
      ) : null}

      {capitalBlock ? (
        <article className="cockpit-card">
          <p className="cockpit-kicker">INSUFFICIENT CAPITAL</p>
          <h3>{title}</h3>
          <p>Buying power {compactUsd(buyingPower)}. Venue min {compactUsd(focus?.minNotional)}. PIT will not invent size.</p>
          <div className="cta-row">
            <button type="button" className="ghost" onClick={() => onAsk("What can I trade now?")}>Show minimum</button>
          </div>
        </article>
      ) : null}

      {fail && stop ? (
        <article className="cockpit-card">
          <p className="cockpit-kicker">STOPPED</p>
          <h3>{stop.title}</h3>
          <p>{stop.body}</p>
        </article>
      ) : null}

      {(ready || noTrade || verified) && coin ? (
        <article className="cockpit-card og">
          <p className="cockpit-kicker">PRIVATE RESEARCH</p>
          <h3>{coin} {side || "committee"}</h3>
          <p>{researchNote || title}</p>
          {why.slice(0, 3).map((row) => (
            <p key={row.q} className="fine"><strong>{row.q}</strong> {row.a}</p>
          ))}
          <div className="cta-row">
            {ready ? <button type="button" className="primary" onClick={onOpenPreview}>Review</button> : null}
            <button type="button" className="ghost" onClick={() => onAsk("Compare top opportunities")}>Compare</button>
            <button type="button" className="ghost" onClick={() => onAsk("Why didn't you trade?")}>Reject path</button>
          </div>
        </article>
      ) : null}

      {proof && (proof.root || proof.tx || proof.jobId) ? (
        <article className="cockpit-card og">
          <p className="cockpit-kicker">PRIVATE RESEARCH PROOF</p>
          <h3>0G Direct · {verified ? "TEE verified" : "not claimed verified"}</h3>
          {proof.jobId ? <p>Research job {shortHash(proof.jobId)}</p> : null}
          {proof.root ? <p>Stored evidence {shortHash(proof.root)}</p> : <p className="fine">No storage root on this computer yet.</p>}
          {proof.tx ? <p>0G Chain {shortHash(proof.tx)}</p> : <p className="fine">No chain transaction until a root is filed.</p>}
          <div className="cta-row">
            {proof.txLink ? <ExternalLink href={proof.txLink}>Open explorer</ExternalLink> : null}
            {proof.root ? <button type="button" className="ghost" onClick={onOpenActivity}>Verify on 0G</button> : null}
          </div>
        </article>
      ) : null}

      {orderState ? (
        <article className={`cockpit-card ${orderState === "filled" ? "ready" : ""}`}>
          <p className="cockpit-kicker">{orderState === "filled" ? "TRADE EXECUTED" : orderState === "resting" ? "ORDER RESTING" : orderState === "failed" ? "ORDER FAILED" : "ORDER"}</p>
          <h3>{lastOrder?.market || "—"} {String(lastOrder?.side || "").toUpperCase()}</h3>
          <p>OID {lastOrder?.oid || lastOid}</p>
          <p>Status {String(lastOrder?.lifecycle || lastOrder?.status || "").toUpperCase()}</p>
          {orderState === "resting" ? <p className="fine">No fill yet. PIT will not call this a fill.</p> : null}
          {orderState === "failed" ? <p className="fine">No execution occurred.</p> : null}
          <div className="cta-row">
            <ExternalLink href={LINKS.hl}>Open Hyperliquid</ExternalLink>
            <button type="button" className="ghost" onClick={onOpenActivity}>Open Activity</button>
          </div>
        </article>
      ) : null}

      <article className="cockpit-card">
        <p className="cockpit-kicker">SLEEP MISSION</p>
        <p>I can prepare this for your desktop. Chat cannot arm or authorize it.</p>
        <p className="fine">Bound to the pinned host law. Research → challenge → risk → policy → execute only after ARM SLEEP MISSION on this computer.</p>
        <div className="cta-row">
          <button type="button" className="ghost" onClick={onOpenAutomation}>Open Sleep Mission</button>
        </div>
      </article>

      <div className="cockpit-follow">
        {follow.map(([label, q]) => (
          <button key={label} type="button" className="chip-btn" onClick={() => onAsk(q)}>{label}</button>
        ))}
      </div>

      <details className="cockpit-tech" open={tech} onToggle={(e) => setTech((e.target as HTMLDetailsElement).open)}>
        <summary>Technical details</summary>
        <p>Provider: Direct TeeML · Model: {researchSku?.model || "not this chat stream"} · Privacy: {researchSku?.proven_e2ee ? "sealed" : "not claimed private"}</p>
        <p>Job {jobId || "—"} · Kind {kind || "—"} · Stage {stage || "—"}</p>
        <p>Catalog listing is not an inference path. {LINKS.pcAdvanced ? "Private book never uses Router." : ""}</p>
      </details>
    </section>
  );
}
