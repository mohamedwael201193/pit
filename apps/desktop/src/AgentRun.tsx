import { useMemo, useState } from "react";
import { ExternalLink } from "./ExternalLink";
import { compactNum, compactUsd, pctFunding } from "./format";
import { committeeVerified, oidBelongsToPreview } from "./honesty";
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
  { id: "DISCOVERY", label: "Scanning markets" },
  { id: "RANKING", label: "Ranking" },
  { id: "SEALING_PRIVATE_BOOK", label: "Private 0G research" },
  { id: "RESEARCHER", label: "Researcher" },
  { id: "CHALLENGER", label: "Challenger" },
  { id: "RISK", label: "Risk" },
  { id: "VERIFYING_TEE_SIGNATURE", label: "TEE verification" },
  { id: "POLICY", label: "Policy" },
  { id: "PREVIEW", label: "Preview" },
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

function stageIndex(stage: string, roles: Role[], busy: boolean, scannedN: number) {
  const s = (stage || "").toUpperCase();
  if (!busy && !s) return -1;
  if (s === "READING_MARKET" || s === "DISCOVERY") return scannedN > 0 ? 1 : 0;
  if (s.includes("RANK")) return 1;
  if (s.includes("SEAL") || s.includes("CONTACT") || s.includes("RECEIVING")) return 2;
  if (s === "RESEARCHER" || (roleOk(roles, "researcher") && !roleOk(roles, "challenger") && busy)) return 3;
  if (s === "CHALLENGER") return 4;
  if (s === "RISK" || s.startsWith("RISK_")) return 5;
  if (s.includes("TEE") || s.includes("VERIFY")) return 6;
  if (s === "DETERMINISTIC_ENGINE" || s === "POLICY") return 7;
  if (s === "PREVIEW" || s === "READY") return 8;
  if (roleOk(roles, "risk")) return 7;
  if (roleOk(roles, "challenger")) return 5;
  if (roleOk(roles, "researcher")) return 4;
  return busy ? (scannedN > 0 ? 1 : 0) : -1;
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

function shortHash(v?: string) {
  const s = String(v || "");
  if (s.length < 12) return s || "";
  return `${s.slice(0, 8)}…${s.slice(-4)}`;
}

function roleMark(ok: boolean, waiting: boolean) {
  if (ok) return "✓";
  if (waiting) return "·";
  return "○";
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
  pinned,
  sessionAlive,
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
  const [whyOpen, setWhy] = useState(false);
  const [sleepOpen, setSleep] = useState(false);
  const [techOpen, setTech] = useState(false);
  const verified = committeeVerified(roles);
  const executable = coins.filter((c) => c.executionFeasible);
  const eligible = coins.filter((c) => c.eligible || c.policyEligible);
  const scannedN = scanned || coins.length;
  const best = executable[0] || eligible[0] || coins[0];
  const focus = coins.find((c) => c.coin === coin) || best;
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
  const current = stageIndex(stage, roles, busy, scannedN);
  const done = verified || ready || noTrade;
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
  const asset = preview?.market || coin || focus?.coin || "";
  const verdict = ready && (side === "LONG" || side === "SHORT") ? side : noTrade || policyBlock || capitalBlock ? "NO TRADE" : fail ? "STOPPED" : "";
  const reason =
    researchNote ||
    deny ||
    (noTrade ? `${asset || "This book"} did not survive the private challenge.` : "") ||
    (policyBlock ? "Host policy blocked this size." : "") ||
    (capitalBlock ? `Buying power ${compactUsd(buyingPower)} is below venue min ${compactUsd(focus?.minNotional)}.` : "") ||
    stop?.body ||
    "";

  const follow = ready
    ? [
        ["Review", "__review"],
        ["Compare candidates", "Compare top opportunities"],
        ["Show why", "__why"],
        ["Sleep Mission", "__sleep"],
        ["Technical details", "__tech"],
      ]
    : noTrade || policyBlock || capitalBlock
      ? [
          ["Research next", "Find the best opportunity"],
          ["Show why", "__why"],
          ["Compare candidates", "Compare top opportunities"],
          ["Sleep Mission", "__sleep"],
          ["Technical details", "__tech"],
        ]
      : fail
        ? [
            ["Research next", "Find the best opportunity"],
            ["Show why", "__why"],
            ["Technical details", "__tech"],
          ]
        : busy
          ? [["Stop research", "Stop research"]]
          : [];

  const showBook = Boolean(focus && (busy || ready));
  const showPipe = busy;
  const showVerdict = !busy && Boolean(verdict);

  return (
    <div className="agent-mission" aria-label="PIT agent turn">
      <p className="agent-lead">
        {busy
          ? "Scanning live Hyperliquid markets…"
          : showVerdict
            ? "Research complete"
            : "Live Hyperliquid books are on this computer."}
        {busy && elapsed ? <span className="agent-elapsed">{elapsed}{pollMiss ? " reconnecting" : ""}</span> : null}
      </p>

      {busy ? (
      <ul className="agent-checks">
        <li className={scannedN ? "done" : "on"}>
          {scannedN ? `${scannedN} markets scanned` : "Waiting for live books"}
        </li>
        <li className={eligible.length ? "done" : scannedN ? "on" : ""}>
          {eligible.length ? `policy filtered · ${eligible.length} eligible` : "policy filter pending"}
        </li>
        <li className={executable.length ? "done" : eligible.length ? "on" : ""}>
          {executable.length ? `capital checked · ${executable.length} executable` : "capital check pending"}
        </li>
        <li className={executable.length ? "done" : ""}>
          {executable.length ? `${executable.length} executable candidates` : "no executable candidate yet"}
        </li>
      </ul>
      ) : scannedN ? (
        <p className="agent-note">{scannedN} scanned · {eligible.length} eligible · {executable.length} executable</p>
      ) : null}

      {showBook && focus ? (
        <article className="agent-card book">
          <p className="agent-kicker">{busy ? "Best opportunity" : ready ? "Live preview" : "Candidate"}</p>
          <header className="agent-book-head">
            <h3>{focus.coin}</h3>
            <p>{compactUsd(focus.mark)}</p>
          </header>
          <dl className="agent-metrics">
            <div>
              <dt>min</dt>
              <dd>{compactUsd(focus.minNotional)}</dd>
            </div>
            <div>
              <dt>clip</dt>
              <dd>{compactUsd(focus.hostNotional || focus.policyClip)}</dd>
            </div>
            <div>
              <dt>funding</dt>
              <dd>{pctFunding(focus.funding)}</dd>
            </div>
            <div>
              <dt>OI</dt>
              <dd>{compactNum(focus.openInterest)}</dd>
            </div>
          </dl>
          {bestWhy ? <p className="agent-note">{bestWhy}</p> : null}
        </article>
      ) : null}

      {showPipe ? (
        <article className="agent-card pipe" aria-live="polite">
          <p className="agent-kicker">PRIVATE 0G RESEARCH</p>
          <ol className="agent-pipe">
            {PIPE.map((step, i) => {
              const markState = pipeState(i, current, done, fail);
              return (
                <li key={step.id} className={markState}>
                  <em aria-hidden>
                    {markState === "done" ? "✓" : markState === "on" ? "●" : markState === "fail" ? "✕" : "○"}
                  </em>
                  <strong>{step.label}</strong>
                  {markState === "on" ? <span>{stage.replaceAll("_", " ").toLowerCase() || "live"}</span> : null}
                  {markState === "fail" && stop ? <span>{stop.title}</span> : null}
                </li>
              );
            })}
          </ol>
          {busy ? (
            <button type="button" className="ghost" onClick={onStop}>
              Stop
            </button>
          ) : null}
        </article>
      ) : null}

      {showVerdict ? (
        <article className="agent-card verdict">
          <p className="agent-kicker">Research complete</p>
          <p className="agent-verdict">
            Verdict: {verdict}
          </p>
          <div className="agent-sections">
            <section>
              <h4>Thesis</h4>
              <p>{why[0]?.a || reason || "No thesis survived."}</p>
            </section>
            <section>
              <h4>Evidence</h4>
              <p>
                {why[2]?.a} {why[3]?.a}
              </p>
            </section>
            <section>
              <h4>Risk</h4>
              <p>{why[4]?.a}</p>
            </section>
            <section>
              <h4>Why this {ready ? "passed" : "failed"} policy</h4>
              <p>{why[5]?.a}</p>
            </section>
            <section>
              <h4>What would invalidate it</h4>
              <p>{why[7]?.a}</p>
            </section>
          </div>
        </article>
      ) : null}

      {ready && preview ? (
        <article className="agent-card ready">
          <p className="agent-kicker">READY</p>
          <header className="agent-trade-head">
            <h3>{preview.market || asset}</h3>
            <p>{side || "SIDE"}</p>
            <p>{compactUsd(preview.notionalUsd)}</p>
            <p>1x</p>
            <p>Hyperliquid</p>
          </header>
          <ul className="agent-gates">
            <li>Policy {pinned ? "✓" : "open Policy first"}</li>
            <li>Capital {Number.isFinite(buyingPower) ? "✓" : "unknown"}</li>
            <li>Venue min ✓</li>
            <li>Slippage ✓</li>
            <li>Session {sessionAlive ? "✓" : "create a session"}</li>
          </ul>
          {!sessionAlive ? <p className="agent-note">Create a live session on Security before TRADE NOW.</p> : null}
          {!pinned ? <p className="agent-note">Pin policy on this computer first.</p> : null}
          {alreadyPosted ? <p className="agent-note">This preview already produced OID {lastOrder?.oid}.</p> : null}
          {authErr ? <p className="agent-note" role="alert">{authErr}</p> : null}
          <div className="cta-row">
            <button type="button" className="primary" aria-label="TRADE NOW" disabled={!canTrade} onClick={onTradeNow}>
              {authBusy ? "Submitting…" : "TRADE NOW"}
            </button>
            <button type="button" className="ghost" onClick={onOpenPreview}>
              REVIEW
            </button>
            <button type="button" className="ghost" onClick={() => onAsk("Do not trade")}>
              REJECT
            </button>
          </div>
          <p className="agent-note">{CHAT_AGENT_COPY.cannotAuthorize}. {CHAT_AGENT_COPY.acceptOnDesk}</p>
        </article>
      ) : null}

      {noTrade ? (
        <article className="agent-card stand">
          <p className="agent-kicker">NO TRADE</p>
          <h3>{asset || "Candidate"}</h3>
          <p>{asset ? `${asset} did not survive the private challenge.` : "No side survived the private challenge."}</p>
          <ul className="agent-gates">
            <li>Researcher {roleMark(roleOk(roles, "researcher"), busy)}</li>
            <li>Challenger {roleMark(roleOk(roles, "challenger"), busy)}</li>
            <li>Risk {roleMark(roleOk(roles, "risk"), busy)}</li>
            <li>Policy {policyBlock ? "blocked" : "✓"}</li>
          </ul>
          {reason ? (
            <p className="agent-reason">
              <strong>Reason</strong>
              {reason}
              {huntRejected.length ? ` Checked ${huntRejected.join(", ")}.` : ""}
            </p>
          ) : null}
        </article>
      ) : null}

      {policyBlock && !noTrade ? (
        <article className="agent-card stand">
          <p className="agent-kicker">NO TRADE</p>
          <h3>{asset || "Policy"}</h3>
          <p>Host policy blocked this size. The model cannot raise clip.</p>
          <div className="cta-row">
            <button type="button" className="ghost" onClick={onOpenPolicy}>
              Open Policy
            </button>
          </div>
        </article>
      ) : null}

      {capitalBlock && !noTrade && !policyBlock ? (
        <article className="agent-card stand">
          <p className="agent-kicker">NO TRADE</p>
          <h3>{asset || "Capital"}</h3>
          <p>
            Buying power {compactUsd(buyingPower)}. Venue min {compactUsd(focus?.minNotional)}.
          </p>
        </article>
      ) : null}

      {fail && stop ? (
        <article className="agent-card stand">
          <p className="agent-kicker">Stopped</p>
          <h3>{stop.title}</h3>
          <p>{stop.body}</p>
        </article>
      ) : null}

      {orderState && lastOrder ? (
        <article className="agent-card order">
          <p className="agent-kicker">ORDER SUBMITTED</p>
          <p>OID {lastOrder.oid || lastOid}</p>
          <p>{String(lastOrder.lifecycle || lastOrder.status || orderState).toUpperCase()}</p>
          <p>
            {orderState === "resting"
              ? "RESTING"
              : orderState === "filled"
                ? "FILLED"
                : orderState === "cancelled"
                  ? "CANCELLED"
                  : orderState === "failed"
                    ? "FAILED"
                    : orderState.toUpperCase()}
          </p>
          <div className="cta-row">
            <ExternalLink className="ghost" href={LINKS.hl}>Open Hyperliquid</ExternalLink>
            <button type="button" className="ghost" onClick={onOpenActivity}>
              Open Activity
            </button>
            {proof?.root || proof?.tx ? (
              <button type="button" className="ghost" onClick={onOpenActivity}>
                Verify on 0G
              </button>
            ) : null}
          </div>
        </article>
      ) : null}

      {proof && (proof.root || proof.tx) && !busy ? (
        <p className="agent-proof">
          {verified ? "TEE verified" : "Proof on file, not claimed verified"}
          {proof.root ? ` · root ${shortHash(proof.root)}` : ""}
          {proof.txLink ? (
            <>
              {" "}
              <ExternalLink className="ghost" href={proof.txLink}>Open 0G explorer</ExternalLink>
            </>
          ) : null}
        </p>
      ) : null}

      {follow.length ? (
        <div className="agent-follow">
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
                  setWhy((v) => !v);
                  return;
                }
                if (q === "__review") {
                  onOpenPreview();
                  return;
                }
                if (q === "__sleep") {
                  setSleep(true);
                  return;
                }
                if (q === "__tech") {
                  setTech((v) => !v);
                  return;
                }
                onAsk(q);
              }}
            >
              {label}
            </button>
          ))}
        </div>
      ) : null}

      {whyOpen ? (
        <article className="agent-card why">
          {why.map((row) => (
            <p key={row.q}>
              <strong>{row.q}</strong> {row.a}
            </p>
          ))}
        </article>
      ) : null}

      {sleepOpen && !busy ? (
        <article className="agent-card">
          <p>Sleep Mission stays bound by current policy. This Agent may prepare it. Desktop still arms.</p>
          <button type="button" className="ghost" onClick={onOpenAutomation}>
            Open Automation
          </button>
        </article>
      ) : null}

      {techOpen && !busy ? (
        <article className="agent-card agent-tech">
          <p>
            Direct TeeML · {researchSku?.model || "not this chat stream"} · {researchSku?.proven_e2ee ? "sealed" : "privacy not claimed here"}
          </p>
          <p>
            Job {jobId ? shortHash(jobId) : "none"} · {kind || "idle"}
          </p>
          <p>Private book never uses Router.</p>
        </article>
      ) : null}
    </div>
  );
}
