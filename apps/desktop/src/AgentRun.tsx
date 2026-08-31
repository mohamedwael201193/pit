import { useMemo, useState } from "react";
import { ExternalLink } from "./ExternalLink";
import { compactNum, compactUsd, pctFunding } from "./format";
import { committeeVerified, oidBelongsToPreview } from "./honesty";
import { explainStop } from "./explain";
import { researchWhyCopy } from "./researchWhy";
import { hyperliquidTrade } from "./links";
import { collectJobReceipts, evidenceObjectForJob, shortProof, venueOrderState } from "./jobProof";
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
  { id: "DECISION", label: "Decision" },
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
  if (s === "DECISION") return 8;
  if (s === "PREVIEW" || s === "READY") return 9;
  if (roleOk(roles, "risk")) return 8;
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

function ResearchStages({ current, done, fail }: { current: number; done: boolean; fail: boolean }) {
  return (
    <>
      <ol className="agent-track" aria-label="Research stage bar">
        {PIPE.map((step, i) => {
          const markState = pipeState(i, current, done, fail);
          return (
            <li key={step.id} className={markState} title={step.label}>
              <span className="sr-only">{step.label}</span>
            </li>
          );
        })}
      </ol>
      <ol className="agent-pipe" aria-label="Research stages">
        {PIPE.map((step, i) => {
          const markState = pipeState(i, current, done, fail);
          return (
            <li key={step.id} className={markState}>
              <em>{markState === "done" ? "✓" : markState === "on" ? "●" : markState === "fail" ? "×" : "○"}</em>
              <strong>{step.label}</strong>
              {markState === "on" ? <span>live</span> : markState === "done" ? <span>done</span> : <span />}
            </li>
          );
        })}
      </ol>
    </>
  );
}

function shortHash(v?: string) {
  return shortProof(v);
}

function roleMark(ok: boolean, waiting: boolean) {
  if (ok) return "✓";
  if (waiting) return "·";
  return "○";
}

function committeeSide(role?: Role) {
  const raw = String(role?.proposed_side || "").trim().toLowerCase();
  if (role?.kill || role?.survives === false) {
    return raw && raw !== "none" ? `rejected ${raw}` : "rejected none";
  }
  if (!raw || raw === "none") return "none";
  return raw;
}

function LiveFacts({ coin, buyingPower }: { coin: MarketCoin; buyingPower?: number }) {
  return (
    <dl className="agent-metrics">
      <div>
        <dt>Mark</dt>
        <dd>{compactUsd(coin.mark)}</dd>
      </div>
      <div>
        <dt>Oracle</dt>
        <dd>{coin.oracle ? compactUsd(coin.oracle) : "—"}</dd>
      </div>
      <div>
        <dt>Funding</dt>
        <dd>{pctFunding(coin.funding)}</dd>
      </div>
      <div>
        <dt>Open interest</dt>
        <dd>{compactNum(coin.openInterest)}</dd>
      </div>
      <div>
        <dt>Venue min</dt>
        <dd>{compactUsd(coin.minNotional)}</dd>
      </div>
      <div>
        <dt>Host notional</dt>
        <dd>{compactUsd(coin.hostNotional)}</dd>
      </div>
      <div>
        <dt>Policy clip</dt>
        <dd>{compactUsd(coin.policyClip || coin.hostNotional)}</dd>
      </div>
      <div>
        <dt>Capital</dt>
        <dd>{compactUsd(buyingPower ?? coin.availableMargin)}</dd>
      </div>
    </dl>
  );
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
  net,
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
  net?: string;
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
  const cancelled =
    kind === "CANCELED_BY_USER" ||
    kind === "CANCELED" ||
    researchStop === "research_cancelled" ||
    researchStop === "CANCELED_BY_USER";
  const noTrade = !busy && kind === "READY_STOOD_DOWN" && !cancelled;
  const ready = !busy && kind === "READY_ELIGIBLE" && Boolean(preview?.eligible && (preview.hash || previewHash));
  const fail = !busy && (Boolean(researchStop) || cancelled) && !noTrade && !ready;
  const stop = explainStop(researchStop || "");
  const researcher = roleOf(roles, "researcher");
  const challenger = roleOf(roles, "challenger");
  const risk = roleOf(roles, "risk");
  const side = String(preview?.side || researcher?.proposed_side || "").toUpperCase();
  const hash = String(previewHash || preview?.hash || "");
  const alreadyPosted = Boolean(lastOrder?.posted && lastOrder.hash && hash && lastOrder.hash === hash);
  const canTrade = ready && Boolean(hash) && sessionAlive && pinned && !authBusy && !alreadyPosted;
  const current = stageIndex(stage, roles, busy, scannedN);
  const venue = net || "mainnet";
  const proofRows = useMemo(() => {
    const ev = evidence && typeof evidence === "object" ? (evidence as Record<string, unknown>) : null;
    const extra = evidenceObjectForJob(ev, jobId);
    const events = extra ? [...activity, extra] : activity;
    return collectJobReceipts(events, jobId, venue);
  }, [activity, evidence, jobId, venue]);
  const orderState = oidBelongsToPreview(lastOrder?.hash, previewHash, hash) ? venueOrderState(lastOrder) : "";
  const orderFiledAt = activity.find((e) => e.oid && e.oid === (lastOrder?.oid || lastOid))?.ts;
  const elapsed = elapsedMs > 0 ? `${(elapsedMs / 1000).toFixed(1)}s` : "";
  const why = researchWhyCopy({
    coin: coin || focus?.coin || "",
    kind,
    note: researchNote,
    stop: researchStop,
    deny,
    eligible: Boolean(preview?.eligible),
    roles,
    snap: focus
      ? {
          mark: focus.mark,
          reason: focus.reason,
          why: focus.why,
          whyRanked: focus.whyRanked,
          invalidation: focus.invalidation,
          expectedEdge: focus.expectedEdge,
        }
      : undefined,
  });
  const asset = (busy ? coin : preview?.market) || coin || focus?.coin || "";
  const verdict = ready && (side === "LONG" || side === "SHORT") ? side : noTrade || policyBlock || capitalBlock ? "NO TRADE" : fail ? "STOPPED" : "";
  const reason =
    researchNote ||
    deny ||
    (noTrade ? `${asset || "This book"} did not survive the private challenge.` : "") ||
    (policyBlock ? "Host policy blocked this size." : "") ||
    (capitalBlock ? `Buying power ${compactUsd(buyingPower)} is below venue min ${compactUsd(focus?.minNotional)}.` : "") ||
    stop?.body ||
    "";

  const huntDone = !busy && huntRejected.length > 0 && executable.length > 0 && huntRejected.length >= executable.length;
  const follow = ready
    ? []
    : huntDone
      ? [
          ["Scan again", "Find the best opportunity"],
          ["Show why", "__why"],
        ]
      : noTrade || policyBlock || capitalBlock
      ? [
          ["Research next", "Find the next opportunity"],
          ["Scan again", "Find the best opportunity"],
          ["Show why", "__why"],
        ]
      : fail
        ? [
            ["Research next", "Find the next opportunity"],
            ["Show why", "__why"],
            ["Technical details", "__tech"],
          ]
        : busy
          ? [["Stop research", "Stop research"]]
          : [];

  const showBook = Boolean(focus && (busy || noTrade || ready || policyBlock || capitalBlock));
  const showVerdict = !busy && Boolean(verdict) && !ready && !noTrade && !cancelled;
  const showPipe = busy || noTrade || ready || fail || showVerdict;
  const liveStep = current >= 0 ? PIPE[Math.min(current, PIPE.length - 1)] : PIPE[0];
  const rejectedSide = [committeeSide(researcher), committeeSide(challenger), committeeSide(risk)].filter((s, i, all) => all.indexOf(s) === i).join(" · ");
  const lead = busy
    ? `Researching ${asset || coin || "live books"} on 0G Direct…`
    : ready
      ? "Opportunity found"
      : huntDone
        ? "Checked every executable book. No side survived."
        : noTrade
          ? `NO TRADE · ${asset || "this book"}`
          : fail
            ? "Research stopped"
            : "Live Hyperliquid books are on this computer.";

  return (
    <div className="agent-mission" aria-label="PIT agent turn">
      <p className="agent-lead">
        {lead}
        {busy && elapsed ? <span className="agent-elapsed">{elapsed}{pollMiss ? " reconnecting" : ""}</span> : null}
      </p>

      {busy ? (
        <article className="agent-live" aria-live="polite">
          <p className="agent-kicker">PRIVATE 0G RESEARCH</p>
          <p className="agent-live-coin">{asset || coin || "markets"}</p>
          <p className="agent-live-step">{liveStep.label}{elapsed ? ` · ${elapsed}` : ""}</p>
          <ResearchStages current={current} done={false} fail={fail} />
          <p className="agent-note">
            Researcher {researcher?.proposed_side || (roleOk(roles, "researcher") ? "verified" : "waiting")}
            {" · "}Challenger {challenger?.survives === false ? "killed" : challenger?.proposed_side || (roleOk(roles, "challenger") ? "verified" : "waiting")}
            {" · "}Risk {risk?.kill ? "kill" : risk?.proposed_side || (roleOk(roles, "risk") ? "verified" : "waiting")}
          </p>
          <p className="agent-note">
            {scannedN ? `${scannedN} live Hyperliquid perps scanned` : "Waiting for live books"}
            {eligible.length ? ` · ${eligible.length} eligible` : ""}
            {executable.length ? ` · ${executable.length} executable` : ""}
            {jobId ? ` · TeeML ${shortHash(jobId)}` : ""}
            {huntRejected.length ? ` · checked ${huntRejected.join(", ")}` : ""}
          </p>
          {proofRows.length ? (
            <ul className="agent-live-proof">
              {proofRows.map((row) => (
                <li key={`live-${row.label}-${row.text}`}>
                  {row.href ? (
                    <ExternalLink className="ghost" href={row.href}>
                      {row.label} {row.text}
                    </ExternalLink>
                  ) : (
                    <span>{row.label} {row.text}</span>
                  )}
                </li>
              ))}
            </ul>
          ) : (
            <p className="agent-note">Waiting for this research run’s 0G receipt…</p>
          )}
          <button type="button" className="ghost" onClick={onStop}>
            Stop
          </button>
        </article>
      ) : scannedN ? (
        <p className="agent-note">{scannedN} scanned · {eligible.length} eligible · {executable.length} executable{huntRejected.length ? ` · checked ${huntRejected.join(", ")}` : ""}</p>
      ) : null}

      {!busy && showPipe ? (
        <article className="agent-live done">
          <p className="agent-kicker">{fail ? "COMMITTEE STOPPED" : ready ? "COMMITTEE PASSED" : "COMMITTEE"}</p>
          <ResearchStages current={fail ? Math.max(current, 2) : PIPE.length} done={!fail} fail={fail} />
        </article>
      ) : null}

      {showBook && focus ? (
        <article className="agent-card book">
          <p className="agent-kicker">LIVE MARKET</p>
          <header className="agent-book-head">
            <h3>{focus.coin}</h3>
            <p>{compactUsd(focus.mark)}</p>
          </header>
          <LiveFacts coin={focus} buyingPower={buyingPower} />
          {focus.whyRanked || bestWhy ? (
            <p className="agent-note">
              <strong>Why ranked. </strong>
              {focus.whyRanked || bestWhy} This is host rank of live venue facts, not a forecast.
            </p>
          ) : null}
          {focus.expectedEdge ? (
            <p className="agent-note">
              <strong>Host edge note. </strong>
              {focus.expectedEdge} Not a committee forecast.
            </p>
          ) : null}
          {focus.invalidation ? (
            <p className="agent-note">
              <strong>Live invalidation. </strong>
              {focus.invalidation}
            </p>
          ) : (
            <p className="agent-note">Research horizon: sealed snapshot only. Host did not publish a timed forecast.</p>
          )}
        </article>
      ) : null}

      {showVerdict && !ready ? (
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
              <h4>Rejected side</h4>
              <p>{rejectedSide || "none"}</p>
            </section>
            <section>
              <h4>Risk</h4>
              <p>{why[4]?.a}</p>
            </section>
            <section>
              <h4>Policy</h4>
              <p>{why[5]?.a}</p>
            </section>
            <section>
              <h4>Invalidation</h4>
              <p>{why[7]?.a}</p>
            </section>
          </div>
        </article>
      ) : null}

      {ready && preview ? (
        <article className="agent-card ready">
          <p className="agent-kicker">OPPORTUNITY FOUND</p>
          <dl className="agent-found">
            <div><dt>Asset</dt><dd>{preview.market || asset}</dd></div>
            <div><dt>Side</dt><dd>{side || "—"}</dd></div>
            <div><dt>Live mark</dt><dd>{compactUsd(focus?.mark)}</dd></div>
            <div><dt>Oracle</dt><dd>{focus?.oracle ? compactUsd(focus.oracle) : "—"}</dd></div>
            <div><dt>Venue min</dt><dd>{compactUsd(focus?.minNotional)}</dd></div>
            <div><dt>Host size</dt><dd>{compactUsd(preview.notionalUsd)}</dd></div>
            <div><dt>Policy clip</dt><dd>{compactUsd(focus?.hostNotional || focus?.policyClip)}</dd></div>
            <div><dt>Available capital</dt><dd>{compactUsd(buyingPower)}</dd></div>
            <div><dt>Leverage</dt><dd>1x</dd></div>
            <div><dt>Exact size</dt><dd>{preview.sz ?? "—"}</dd></div>
            <div><dt>Limit</dt><dd>{preview.limitPx || "—"}</dd></div>
            <div><dt>Policy</dt><dd>{pinned ? "pinned" : "not pinned"}</dd></div>
          </dl>
          <div className="agent-sections">
            <section>
              <h4>Thesis</h4>
              <p>{why[0]?.a || `${side} survived the private committee on ${asset}.`}</p>
            </section>
            <section>
              <h4>Committee forecast</h4>
              <p>
                Expected direction is the verified committee side {side || "none"}. Host did not publish a price target or confidence on this preview.
              </p>
            </section>
            <section>
              <h4>Invalidation</h4>
              <p>{focus?.invalidation || why[7]?.a}</p>
            </section>
            <section>
              <h4>Risk</h4>
              <p>{why[4]?.a} Risk/reward is not shown unless the engine computed it. This preview has size {preview.sz ?? "—"} at {preview.limitPx || "the host limit"}.</p>
            </section>
          </div>
          <ul className="agent-gates">
            <li>Researcher {roleMark(roleOk(roles, "researcher"), false)}</li>
            <li>Challenger {roleMark(roleOk(roles, "challenger"), false)}</li>
            <li>Risk {roleMark(roleOk(roles, "risk"), false)}</li>
            <li>TEE {verified ? "✓" : "○"}</li>
            <li>Policy {pinned ? "✓" : "open Policy first"}</li>
          </ul>
          <p className="agent-kicker">PREVIEW READY</p>
          {!sessionAlive ? <p className="agent-note">Create a live session on Security before TRADE NOW.</p> : null}
          {!pinned ? <p className="agent-note">Pin policy on this computer first.</p> : null}
          {alreadyPosted ? <p className="agent-note">This preview already produced OID {lastOrder?.oid}.</p> : null}
          {authErr ? <p className="agent-note" role="alert">{authErr}</p> : null}
          <div className="cta-row trade-now-row">
            <button type="button" className="ghost" onClick={onOpenPreview}>
              REVIEW
            </button>
            <button type="button" className="primary" aria-label="TRADE NOW" disabled={!canTrade} onClick={onTradeNow}>
              {authBusy ? "Submitting…" : "TRADE NOW"}
            </button>
            <button type="button" className="ghost" onClick={() => onAsk("Do not trade")}>
              REJECT
            </button>
          </div>
        </article>
      ) : null}

      {noTrade ? (
        <article className="agent-card stand">
          <p className="agent-kicker">NO TRADE</p>
          <h3>{asset || "Candidate"}</h3>
          <p>No side survived the private challenge.</p>
          <ul className="agent-gates">
            <li>Researcher {roleMark(roleOk(roles, "researcher"), busy)}</li>
            <li>Challenger {roleMark(roleOk(roles, "challenger"), busy)}</li>
            <li>Risk {roleMark(roleOk(roles, "risk"), busy)}</li>
            <li>Policy {policyBlock ? "blocked" : "✓"}</li>
          </ul>
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
              <h4>Rejected side</h4>
              <p>{rejectedSide || "none"}</p>
            </section>
            <section>
              <h4>Reason</h4>
              <p>
                {reason}
                {huntRejected.length ? ` Checked ${huntRejected.join(", ")}.` : ""}
              </p>
            </section>
            <section>
              <h4>Policy</h4>
              <p>{why[5]?.a}</p>
            </section>
            <section>
              <h4>Risk</h4>
              <p>{why[4]?.a}</p>
            </section>
            <section>
              <h4>Current 0G proof</h4>
              <p>{jobId ? `This job ${shortHash(jobId)} only. Historical receipts stay off this card.` : "Waiting for this research run’s 0G receipt…"}</p>
            </section>
          </div>
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
          <dl className="agent-found">
            <div><dt>OID</dt><dd>{lastOrder.oid || lastOid}</dd></div>
            <div>
              <dt>Status</dt>
              <dd>
                {orderState === "resting"
                  ? "RESTING"
                  : orderState === "filled"
                    ? "FILLED"
                    : orderState === "cancelled"
                      ? "CANCELED"
                      : orderState === "failed"
                        ? "FAILED"
                        : orderState.toUpperCase()}
              </dd>
            </div>
            <div><dt>Side</dt><dd>{String(lastOrder.side || side || "").toUpperCase() || "—"}</dd></div>
            <div><dt>Size</dt><dd>{lastOrder.sz ?? preview?.sz ?? "—"}</dd></div>
            <div><dt>Price</dt><dd>{preview?.limitPx || "—"}</dd></div>
            <div><dt>Time</dt><dd>{orderFiledAt ? new Date(Number(orderFiledAt)).toLocaleTimeString() : "—"}</dd></div>
          </dl>
          <div className="cta-row">
            <ExternalLink className="ghost" href={hyperliquidTrade(venue, String(lastOrder.market || asset).split(":").pop())}>Open Hyperliquid</ExternalLink>
            <button type="button" className="ghost" onClick={onOpenActivity}>
              Open Activity
            </button>
          </div>
        </article>
      ) : null}

      {jobId && (busy || noTrade || ready || fail) ? (
        <article className="agent-receipts">
          <p className="agent-kicker">{verified ? "0G PROOF" : "0G TRAIL"}</p>
          {proofRows.length ? (
            <ul>
              {proofRows.map((row) => (
                <li key={`${row.label}-${row.text}`}>
                  <span>{row.label}</span>
                  {row.href ? (
                    <ExternalLink className="ghost" href={row.href}>
                      {row.text}
                    </ExternalLink>
                  ) : (
                    <strong>{row.text}</strong>
                  )}
                  <em>
                    {row.ts ? new Date(row.ts).toLocaleTimeString() : ""} · job {shortHash(row.jobId)}
                    {row.market ? ` · ${row.market}` : asset ? ` · ${asset}` : ""}
                  </em>
                </li>
              ))}
            </ul>
          ) : (
            <p className="agent-note">Waiting for this research run’s 0G receipt…</p>
          )}
          <p className="agent-note">job {shortHash(jobId)}{asset ? ` · ${asset}` : ""}</p>
        </article>
      ) : null}

      {follow.length ? (
        <div className="agent-follow">
          {follow.map(([label, q], i) => (
            <button
              key={label}
              type="button"
              className={i === 0 && !ready && !busy ? "chip-btn lead" : "chip-btn"}
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
