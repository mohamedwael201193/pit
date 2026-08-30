import { useEffect, useState } from "react";
import { BrandMark } from "./BrandMark";
import { elapsedLabel, nextScanLabel, remainLabel } from "./format";
import { Select } from "./Select";
import { explorerAddress, explorerTx, hyperliquidAPI, hyperliquidApp, hyperliquidTrade, LINKS } from "./links";
import { ExternalLink } from "./ExternalLink";
import type { AutoPrefs, MissionPublic } from "./companion";

const ARM_TOKEN = "ARM SLEEP MISSION";

const CADENCE = [
  { value: "1", label: "Every minute" },
  { value: "5", label: "Every 5 minutes" },
  { value: "15", label: "Every 15 minutes" },
  { value: "60", label: "Every hour" },
];

const HOURS = [
  { value: "1", label: "1 hour" },
  { value: "8", label: "8 hours" },
  { value: "24", label: "24 hours" },
  { value: "72", label: "72 hours" },
];

type Phase = "idle" | "enabling" | "stopping";

function limitRows(limits: Record<string, unknown>): { k: string; v: string }[] {
  const assets = Array.isArray(limits.allowed_assets) ? (limits.allowed_assets as string[]).join(" ") : "ETH BTC SOL HYPE DOGE AVAX";
  const venues = Array.isArray(limits.allowed_venues) ? (limits.allowed_venues as string[]).join(" ") : "hyperliquid";
  return [
    { k: "Allowed assets", v: assets },
    { k: "Allowed venues", v: venues },
    { k: "Max trade", v: `$${String(limits.max_trade_usd ?? limits.max_clip_usd ?? 10)}` },
    { k: "Max position", v: `$${String(limits.max_position_usd ?? limits.max_clip_usd ?? 10)}` },
    { k: "Max leverage", v: `${String(limits.max_leverage ?? 1)}x` },
    { k: "Daily loss", v: `$${String(limits.daily_loss_usd ?? 50)}` },
    { k: "Consecutive loss limit", v: String(limits.max_consecutive_losses ?? 3) },
    { k: "Max open positions", v: String(limits.max_open_positions ?? 1) },
    { k: "Slippage", v: `${String(limits.max_slippage_bps ?? 80)} bps` },
    { k: "Liquidity", v: limits.min_liquidity_usd ? `$${String(limits.min_liquidity_usd)}` : "No extra floor" },
    { k: "Cooldown", v: `${String(limits.cooldown_seconds ?? 0)}s` },
    { k: "Uncertainty threshold", v: String(limits.max_uncertainty ?? 1) },
    { k: "Session expiry", v: `${String(limits.session_ttl_seconds ?? limits.session_expiry ?? 3600)}s` },
    { k: "Kill switch", v: limits.kill_switch ? "On - new orders halted" : "Off" },
    { k: "Withdraw", v: "Forbidden" },
    { k: "Transfer", v: "Forbidden" },
    { k: "Policy mutation", v: "Forbidden" },
  ];
}

function envelopeGroups(limits: Record<string, unknown>) {
  const rows = limitRows(limits);
  const pick = (keys: string[]) => rows.filter((r) => keys.includes(r.k));
  return {
    size: pick(["Allowed assets", "Allowed venues", "Max trade", "Max position", "Max leverage"]),
    risk: pick(["Daily loss", "Consecutive loss limit", "Max open positions", "Slippage", "Liquidity", "Cooldown", "Uncertainty threshold"]),
    halt: pick(["Session expiry", "Kill switch"]),
    forbidden: pick(["Withdraw", "Transfer", "Policy mutation"]),
  };
}

function EnvelopeGroup({ title, rows }: { title: string; rows: { k: string; v: string }[] }) {
  return (
    <div className="envelope-group">
      <p className="label">{title}</p>
      <dl>
        {rows.map((row) => (
          <div key={row.k}>
            <dt>{row.k}</dt>
            <dd>{row.v}</dd>
          </div>
        ))}
      </dl>
    </div>
  );
}

function lifeOf(args: {
  kill?: boolean;
  phase: Phase;
  error?: string;
  host?: string;
}): string {
  if (args.phase === "enabling") return "ENABLING";
  if (args.phase === "stopping") return "STOPPING";
  if (args.error) return "FAILED";
  if (args.host === "ACTIVE") return "ACTIVE";
  if (args.host === "BLOCKED") return "BLOCKED";
  if (args.host === "STOPPED" || args.host === "FAILED") return args.host;
  if (args.kill) return "BLOCKED";
  return args.host || "READY";
}

export function AutomationCenter({
  mission,
  prefs,
  busy,
  kill,
  openConfirm,
  confirmHours,
  onMode,
  onEnable,
  onStop,
  onSavePrefs,
  onConfirmConsumed,
  onOpenHistory,
  net,
  wallet,
  execGate,
  execWhy,
}: {
  mission: MissionPublic;
  prefs: AutoPrefs;
  busy?: boolean;
  kill?: boolean;
  openConfirm?: boolean;
  confirmHours?: number;
  onMode: (mode: "manual" | "research_only") => void;
  onEnable: (typed: string, hours: number) => Promise<MissionPublic | void> | void;
  onStop: () => Promise<unknown> | void;
  onSavePrefs: (p: AutoPrefs) => void;
  onConfirmConsumed?: () => void;
  onOpenHistory?: () => void;
  net?: string;
  wallet?: string;
  execGate?: string;
  execWhy?: string;
}) {
  const [hours, setHours] = useState(confirmHours || 8);
  const [confirm, setConfirm] = useState(false);
  const [reviewed, setReviewed] = useState(false);
  const [phase, setPhase] = useState<Phase>("idle");
  const [err, setErr] = useState<string | null>(null);
  const [now, setNow] = useState(() => Math.floor(Date.now() / 1000));

  useEffect(() => {
    const t = window.setInterval(() => setNow(Math.floor(Date.now() / 1000)), 1000);
    return () => window.clearInterval(t);
  }, []);

  useEffect(() => {
    if (openConfirm) {
      setHours(confirmHours || 8);
      setConfirm(true);
      setReviewed(false);
      onConfirmConsumed?.();
    }
  }, [openConfirm, confirmHours, onConfirmConsumed]);

  const m = mission.mission || {};
  const mode = mission.mode || m.mode || "manual";
  const limits = mission.limits || {};
  const status = lifeOf({
    kill,
    phase,
    error: err || mission.error,
    host: mission.status,
  });
  const live = status === "ACTIVE" && mode === "guarded";
  const deadline = m.guarded_until_unix || m.deadline_unix;
  const remain = remainLabel(deadline, mission.now || now);
  const oid = String(m.last_oid || mission.last_order?.oid || "").trim();
  const proofHash = "";

  async function enable() {
    setPhase("enabling");
    setErr(null);
    try {
      const r = await onEnable(ARM_TOKEN, hours);
      if (r?.error) {
        setErr(r.explain || (r.error === "need_pin" ? "Your policy changed. Re-pin it before trading." : r.error));
        setPhase("idle");
        return;
      }
      setConfirm(false);
      setReviewed(false);
      setPhase("idle");
    } catch {
      setErr("companion_http");
      setPhase("idle");
    }
  }

  async function stop() {
    setPhase("stopping");
    setErr(null);
    try {
      await onStop();
    } finally {
      setPhase("idle");
    }
  }

  const env = envelopeGroups(limits);
  const gateLine = mission.block_reason
    ? humanStop(mission.block_reason)
    : mission.why_not
      ? mission.why_not
      : execGate
        ? `${humanStop(execGate)}${execWhy ? ` - ${execWhy}` : ""}`
        : "clear";

  const sleepState = String(mission.sleep_state || m.sleep_state || (live ? "WATCHING" : "STOPPED"));

  return (
    <main className="page mission-page">
      <div className="page-head">
        <div>
          <p className="eyebrow">Automation</p>
          <h1>{live ? "Sleep Mission armed" : status === "STOPPED" ? "Sleep Mission stopped" : "Sleep Mission"}</h1>
          <p className="lead">Arm once. PIT hunts within your rules. Chat cannot arm or AUTHORIZE.</p>
        </div>
        <div className="halt-rail">
          <span className={`status-chip ${status.toLowerCase()}`}>{status}</span>
          <button
            type="button"
            className="kill-switch compact"
            onClick={() => void stop()}
            disabled={busy || phase !== "idle" || (mode === "manual" && !live)}
          >
            {live ? "STOP" : "Kill switch"}
          </button>
        </div>
      </div>

      {mission.search_note || m.search_note || m.last_result ? (
        <p className="search-note" role="status">
          {mission.search_note || m.search_note || m.last_result}
        </p>
      ) : null}

      <div className="mode-switch" role="tablist" aria-label="Desk mode">
        <button type="button" role="tab" aria-selected={mode === "manual"} className={mode === "manual" ? "on" : ""} disabled={busy || phase !== "idle"} onClick={() => onMode("manual")}>
          <strong>Manual</strong>
          <span>Every order waits for AUTHORIZE.</span>
        </button>
        <button type="button" role="tab" aria-selected={mode === "research_only"} className={mode === "research_only" ? "on" : ""} disabled={busy || phase !== "idle"} onClick={() => onMode("research_only")}>
          <strong>Research</strong>
          <span>Scan and prepare. Never execute.</span>
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={mode === "guarded"}
          className={mode === "guarded" ? "on" : ""}
          disabled={busy || phase !== "idle"}
          onClick={() => {
            if (live) return;
            setConfirm(true);
            setReviewed(false);
          }}
        >
          <strong>Sleep Mission</strong>
          <span>Hunt while this computer stays awake.</span>
        </button>
      </div>

      <div className="mission-desk">
        <section className="mission-primary">
          {live ? (
            <div className="mission-status" role="status">
              <p className="label">Remaining</p>
              <h2>{remain}</h2>
              <dl className="mission-strip">
                <div>
                  <dt>State</dt>
                  <dd>{sleepState}</dd>
                </div>
                <div>
                  <dt>Gate</dt>
                  <dd>{gateLine}</dd>
                </div>
                <div>
                  <dt>Stage</dt>
                  <dd>{humanStage(mission.stage || m.stage || (mission.research_running ? "researching" : status))}</dd>
                </div>
              </dl>
              {mission.block_reason || mission.why_not ? (
                <p className="err" role="status">
                  Why PIT did not trade: {mission.why_not || humanStop(mission.block_reason || "")}. {mission.block_explain || "Scan continues. Positions are not flattened."}
                </p>
              ) : (
                <p className="fine">This computer must stay awake for the bound. If it sleeps, the mission stops.</p>
              )}
            </div>
          ) : (
            <div className="mission-status">
              <p className="label">{m.last_stop ? "Last stop" : "Ready to arm"}</p>
              <h2>{m.last_stop ? humanStop(m.last_stop) : mode === "research_only" ? "Research only" : "Manual"}</h2>
              <p className="fine">
                {m.last_stop
                  ? mission.explain || "No further autonomous orders until you arm again on this computer."
                  : "Chat cannot send the arm phrase. Review limits, then arm on this computer."}
              </p>
            </div>
          )}

          {err || mission.error ? (
            <p className="err" role="alert">
              {err || mission.error}
            </p>
          ) : null}
          {kill ? <p className="err">Kill switch is on. New orders are halted.</p> : null}

          <div className="mission-actions">
            <button
              type="button"
              className="primary"
              onClick={() => {
                setConfirm(true);
                setReviewed(false);
              }}
              disabled={busy || phase !== "idle" || live || kill}
            >
              ARM SLEEP MISSION
            </button>
            <button
              type="button"
              className="linkish"
              onClick={() => {
                setConfirm(true);
                setReviewed(false);
              }}
              disabled={busy || phase !== "idle" || live}
            >
              REVIEW LIMITS
            </button>
          </div>

          <section className="envelope">
            <p className="label">Host envelope</p>
            <ul className="envelope-chips">
              {env.size.slice(2, 5).map((row) => (
                <li key={row.k}>
                  <span>{row.k}</span>
                  <strong>{row.v}</strong>
                </li>
              ))}
              <li>
                <span>Kill switch</span>
                <strong>{kill || limits.kill_switch ? "On" : "Off"}</strong>
              </li>
            </ul>
            <details>
              <summary>Size, risk, universe, forbidden</summary>
              <div className="envelope-groups">
                <EnvelopeGroup title="Size" rows={env.size} />
                <EnvelopeGroup title="Risk" rows={env.risk} />
                <EnvelopeGroup title="Session" rows={env.halt} />
                <EnvelopeGroup title="Forbidden" rows={env.forbidden} />
              </div>
              <p className="host-enforced">The model can never modify these.</p>
            </details>
            <div className="cadence-row">
              <Select
                id="scan-cadence"
                label="Scan cadence"
                value={String(prefs.cadence_minutes || 15)}
                disabled={busy}
                options={CADENCE}
                onChange={(v) => onSavePrefs({ ...prefs, cadence_minutes: Number(v) })}
              />
            </div>
          </section>

          <details className="card">
            <summary>Mission internals</summary>
        <dl className="mission-grid">
          <div>
            <dt>Stage</dt>
            <dd>{humanStage(mission.stage || m.stage || (mission.research_running ? "researching" : status))}</dd>
          </div>
          <div>
            <dt>Next scan</dt>
            <dd>{nextScanLabel(mission.next_scan_unix || m.next_scan_unix, mission.now || now)}</dd>
          </div>
          <div>
            <dt>Elapsed</dt>
            <dd>{elapsedLabel(m.guarded_enabled_unix, mission.now || now)}</dd>
          </div>
          <div>
            <dt>Universe</dt>
            <dd>
              {m.scanned || 0} scanned · {m.eligible || 0} pass policy
              {m.scan_count ? ` · ${m.scan_count} ticks` : ""}
            </dd>
          </div>
          <div>
            <dt>Current book</dt>
            <dd>
              {m.best_coin ? (
                <ExternalLink className="asset" href={hyperliquidTrade(net || "mainnet", m.best_coin)}>
                  <BrandMark symbol={m.best_coin} size={14} /> {m.best_coin}
                </ExternalLink>
              ) : (
                "none"
              )}
            </dd>
          </div>
          <div>
            <dt>Exec gate</dt>
            <dd>{gateLine}</dd>
          </div>
          <div>
            <dt>Trades today</dt>
            <dd>{m.trades_today || 0}</dd>
          </div>
          <div>
            <dt>Risk left</dt>
            <dd>
              ${String(mission.remaining_risk_usd ?? limits.daily_loss_usd ?? 50)} daily ·{" "}
              {String(mission.remaining_consecutive_losses ?? limits.max_consecutive_losses ?? 3)} losses
            </dd>
          </div>
        </dl>
        {m.best_why ? <p className="fine">Why it acted: {m.best_why}</p> : null}
        <div className="mission-links">
          <ExternalLink className="linkish" href={hyperliquidApp(net || "mainnet")}>
            Hyperliquid
          </ExternalLink>
          <ExternalLink className="linkish" href={hyperliquidAPI(net || "mainnet")}>
            Hyperliquid API
          </ExternalLink>
          {wallet ? (
            <ExternalLink className="linkish" href={explorerAddress(wallet, net || "mainnet")}>
              0G explorer
            </ExternalLink>
          ) : (
            <ExternalLink className="linkish" href={LINKS.explorer}>
              0G explorer
            </ExternalLink>
          )}
          {oid ? (
            <ExternalLink className="linkish" href={hyperliquidTrade(net || "mainnet", coinFromMarket(mission.last_order?.market || m.best_coin))}>
              OID {oid}
            </ExternalLink>
          ) : null}
          {proofHash ? (
            <ExternalLink className="linkish" href={explorerTx(proofHash, net || "mainnet")}>
              Storage proof
            </ExternalLink>
          ) : null}
          {onOpenHistory ? (
            <button type="button" className="linkish" onClick={onOpenHistory}>
              Mission history
            </button>
          ) : null}
        </div>
      </details>
        </section>

        <aside className="mission-log">
          <GoodMorning
            live={live}
            elapsed={elapsedLabel(m.guarded_enabled_unix, mission.now || now)}
            events={mission.events}
            stop={m.last_stop}
          />
          <AwayBoard away={mission.away} whyNot={mission.why_not} whyCode={mission.why_not_code || mission.block_reason} />
          <NightReplay events={mission.events} />
        </aside>
      </div>

      {confirm ? (
        <div className="overlay" role="dialog" aria-modal="true" aria-labelledby="enable-title">
          <div className="confirm">
            <p className="label">Confirm</p>
            <h2 id="enable-title">Arm Sleep Mission</h2>
            <p>Review the host limits. The mission cannot weaken pinned policy. Then confirm exactly what this window may do.</p>
            <div className="policy-grid compact">
              {limitRows(limits)
                .slice(0, 8)
                .map((row) => (
                  <div className="policy-cell" key={row.k}>
                    <p className="label">{row.k}</p>
                    <strong>{row.v}</strong>
                  </div>
                ))}
            </div>
            <p className="label">What it can do</p>
            <ul className="can-list">
              <li>Scan the live Hyperliquid universe on cadence.</li>
              <li>Move the best eligible book into private research automatically.</li>
              <li>After verified research and host gates, execute only inside these limits.</li>
            </ul>
            <p className="label">What it cannot do</p>
            <ul className="can-list">
              <li>Execute from chat or change policy, clip, leverage, or permissions.</li>
              <li>Withdraw, transfer, or bypass kill switch, session, preview binding, or duplicate-order protection.</li>
              <li>Keep running if this computer sleeps. The mission stops. That gap is not backfilled. Kill does not flatten.</li>
            </ul>
            <Select label="Duration" value={String(hours)} options={HOURS} onChange={(v) => setHours(Number(v))} disabled={busy || phase !== "idle"} />
            <label className="check-row">
              <input type="checkbox" checked={reviewed} onChange={(e) => setReviewed(e.target.checked)} />
              I reviewed these limits. The Sleep Mission stays inside them. It cannot withdraw, transfer, or raise policy.
            </label>
            {err ? (
              <p className="err" role="alert">
                {err}
              </p>
            ) : null}
            <div className="cta-row">
              <button
                type="button"
                className="linkish"
                onClick={() => {
                  setConfirm(false);
                  setReviewed(false);
                  setErr(null);
                }}
                disabled={phase === "enabling"}
              >
                Cancel
              </button>
              <button type="button" className="primary" disabled={!reviewed || busy || phase === "enabling" || kill} onClick={() => void enable()}>
                {phase === "enabling" ? "Arming…" : "ARM SLEEP MISSION"}
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </main>
  );
}

function humanStop(s: string) {
  if (s === "chat_stop") return "Stopped from Chat";
  if (s === "user_stop") return "You stopped it";
  if (s === "kill_switch") return "Kill switch";
  if (s === "deadline" || s === "autonomy_expired") return "Duration ended";
  if (s === "session_expired") return "Hyperliquid session permissions no longer match this desk.";
  if (s === "max_open_positions") return "Open position ceiling";
  if (s === "daily_loss") return "Daily loss";
  if (s === "policy_changed" || s === "need_pin") return "Your policy changed. Re-pin it before trading.";
  if (s === "max_trades") return "Trade ceiling";
  if (s === "consecutive_loss_limit") return "Consecutive loss ceiling";
  if (s === "duplicate_preview") return "Duplicate preview";
  if (s === "preview_before_guarded") return "Preview started before enable";
  if (s === "insufficient_margin") return "This account cannot clear that book's Hyperliquid floor. PIT will not invent size.";
  if (s === "policy_clip_tight") return "Pinned max trade cannot meet this book's rounded Hyperliquid minimum. The account can. Raise clip, preview, pin. PIT will not invent size.";
  if (s === "below_min_notional") return "This book's rounded Hyperliquid minimum is above this account and policy. PIT will not invent size.";
  if (s === "no_opportunity") return "Nothing qualifies under your law right now.";
  return s.replaceAll("_", " ");
}

function humanStage(s: string) {
  const t = String(s || "").replaceAll("_", "-");
  if (t === "scanning" || t === "starting") return "Watching live markets";
  if (t === "researching") return "Researching privately";
  if (t === "waiting" || t === "waiting after research") return "Waiting for the next scan";
  if (t === "eligible") return "Ready to trade";
  if (t === "execution-blocked" || t === "exec blocked") return "Ready to trade - host refused";
  if (t === "executing" || t === "executed") return "Submitting to Hyperliquid";
  if (t === "resting") return "Resting on the venue";
  if (t === "cooldown") return "Cooldown";
  if (t === "stopped") return "Autonomy stopped";
  if (t === "empty") return "Watching — nothing executable";
  if (t === "searching") return "Checking the next market";
  if (t === "ranked") return "Ranked a candidate";
  if (!t) return "Idle";
  return t.replaceAll("-", " ");
}

function AwayBoard({
  away,
  whyNot,
  whyCode,
}: {
  away?: {
    since_unix?: number;
    detected?: number;
    researched?: number;
    rejected?: number;
    traded?: number;
    filled?: number;
    events?: Array<{ unix?: number; kind?: string; coin?: string; why?: string; human?: string; oid?: string }>;
  };
  whyNot?: string;
  whyCode?: string;
}) {
  const events = [...(away?.events || [])].reverse().slice(0, 12);
  return (
    <section className="away-board">
      <p className="label">While you were away</p>
      <h2>What the desk recorded</h2>
      <p className="fine">
        {whyNot || (whyCode ? humanStop(whyCode) : "No named refusal yet.")} A Sleep Mission never raises your limits.
      </p>
      <dl className="mission-grid">
        <div>
          <dt>Detected</dt>
          <dd>{away?.detected ?? 0}</dd>
        </div>
        <div>
          <dt>Researched</dt>
          <dd>{away?.researched ?? 0}</dd>
        </div>
        <div>
          <dt>Refused</dt>
          <dd>{away?.rejected ?? 0}</dd>
        </div>
        <div>
          <dt>Autonomous trades</dt>
          <dd>{away?.traded ?? 0}</dd>
        </div>
        <div>
          <dt>Fills</dt>
          <dd>{away?.filled ?? 0}</dd>
        </div>
      </dl>
      {events.length === 0 ? (
        <p className="empty">Empty is honest until this computer records a scan, research, or refusal.</p>
      ) : (
        <ul className="away-list">
          {events.map((ev, i) => (
            <li key={`${ev.unix}-${ev.kind}-${i}`}>
              <strong>{String(ev.kind || "event").replaceAll("_", " ")}</strong>
              {ev.coin ? ` ${ev.coin}` : ""}
              {ev.oid ? ` · OID ${ev.oid}` : ""}
              {ev.human || ev.why ? ` · ${ev.human || humanStop(ev.why || "")}` : ""}
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}

function GoodMorning({
  live,
  elapsed,
  events,
  stop,
}: {
  live: boolean;
  elapsed: string;
  events?: MissionLog;
  stop?: string;
}) {
  if (live || (!stop && !(events?.events || []).length)) return null;
  return (
    <section className="away-board" role="status">
      <p className="label">GOOD MORNING</p>
      <h2>Your desk ran for {elapsed === "—" ? "the armed window" : elapsed}.</h2>
      <p className="fine">
        {!(events?.executions || events?.fills)
          ? "No trade was a successful night. Every count is a persisted mission event."
          : "Every count is a persisted mission event. Empty is honest."}
      </p>
      <dl className="mission-grid">
        <div>
          <dt>Opportunities</dt>
          <dd>{events?.opportunities_detected ?? 0}</dd>
        </div>
        <div>
          <dt>Private researches</dt>
          <dd>{events?.private_researches ?? 0}</dd>
        </div>
        <div>
          <dt>Challenger rejects</dt>
          <dd>{events?.challenger_rejects ?? 0}</dd>
        </div>
        <div>
          <dt>Risk rejects</dt>
          <dd>{events?.risk_rejects ?? 0}</dd>
        </div>
        <div>
          <dt>Policy blocks</dt>
          <dd>{events?.policy_blocks ?? 0}</dd>
        </div>
        <div>
          <dt>Executions</dt>
          <dd>{events?.executions ?? 0}</dd>
        </div>
        <div>
          <dt>Fills</dt>
          <dd>{events?.fills ?? 0}</dd>
        </div>
        <div>
          <dt>Proofs</dt>
          <dd>{events?.proofs ?? 0}</dd>
        </div>
        <div>
          <dt>Lessons</dt>
          <dd>{events?.lessons ?? 0}</dd>
        </div>
      </dl>
      <a className="linkish" href="#night-replay">
        VIEW NIGHT
      </a>
    </section>
  );
}

function NightReplay({ events }: { events?: MissionLog }) {
  const rows = events?.events || [];
  return (
    <section className="away-board" id="night-replay">
      <p className="label">Night Replay</p>
      <h2>Committed timeline</h2>
      <p className="fine">A no-trade is success when the thesis did not survive challenge. Replay is not a new trade.</p>
      {rows.length === 0 ? (
        <p className="empty">Empty is honest until this computer records a watch, research, or refusal.</p>
      ) : (
        <ol className="night-replay">
          {rows.map((ev, i) => (
            <li key={`${ev.unix}-${ev.node}-${i}`}>
              <strong>{ev.node || "NODE"}</strong>
              {ev.no_trade || ev.status === "NO-TRADE" ? " NO-TRADE" : ` ${ev.status || ""}`}
              {ev.coin ? ` · ${ev.coin}` : ""}
              {ev.oid ? ` · OID ${ev.oid}` : ""}
              {ev.job_id ? ` · job ${ev.job_id}` : ""}
              {ev.reason ? ` · ${ev.human || ev.reason}` : ""}
            </li>
          ))}
        </ol>
      )}
    </section>
  );
}

type MissionLog = {
  opportunities_detected?: number;
  private_researches?: number;
  challenger_rejects?: number;
  risk_rejects?: number;
  policy_blocks?: number;
  executions?: number;
  fills?: number;
  proofs?: number;
  lessons?: number;
  events?: Array<{
    unix?: number;
    node?: string;
    status?: string;
    reason?: string;
    human?: string;
    job_id?: string;
    oid?: string;
    coin?: string;
    no_trade?: boolean;
  }>;
};

function coinFromMarket(market?: string) {
  const m = String(market || "");
  const parts = m.split(":");
  return parts[parts.length - 1] || m;
}
