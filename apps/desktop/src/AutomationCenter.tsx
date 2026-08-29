import { useEffect, useState } from "react";
import { BrandMark } from "./BrandMark";
import { elapsedLabel, nextScanLabel, remainLabel } from "./format";
import { Select } from "./Select";
import { explorerAddress, explorerTx, hyperliquidAPI, hyperliquidApp, hyperliquidTrade, LINKS } from "./links";
import { ExternalLink } from "./ExternalLink";
import type { AutoPrefs, MissionPublic } from "./companion";

const ENABLE_TOKEN = "ENABLE GUARDED AUTONOMY";

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
    { k: "Kill switch", v: limits.kill_switch ? "On — new orders halted" : "Off" },
    { k: "Withdraw", v: "Forbidden" },
    { k: "Transfer", v: "Forbidden" },
    { k: "Policy mutation", v: "Forbidden" },
  ];
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
  const running = Boolean(mission.status === "ACTIVE" && (mode === "guarded" || mode === "research_only"));
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
      const r = await onEnable(ENABLE_TOKEN, hours);
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

  return (
    <main className="page dense mission-page">
      <div className="page-head">
        <div>
          <p className="eyebrow">Automation</p>
          <h1>{live ? "Autonomy running" : status === "STOPPED" ? "Autonomy stopped" : "Autonomy"}</h1>
        </div>
        <span className={`status-chip ${status.toLowerCase()}`}>{status}</span>
      </div>
      <p className="lead">
        Host-enforced. The model cannot change these limits. Chat cannot enable Guarded Autonomy.
      </p>

      <p className="label">Mode</p>
      <div className="mode-grid">
        <button type="button" className={mode === "manual" ? "on" : ""} disabled={busy || phase !== "idle"} onClick={() => onMode("manual")}>
          <p className="label">Manual</p>
          <p>Scan and research on your command. Every order waits for AUTHORIZE.</p>
        </button>
        <button type="button" className={mode === "research_only" ? "on" : ""} disabled={busy || phase !== "idle"} onClick={() => onMode("research_only")}>
          <p className="label">Research only</p>
          <p>Scan, research, notify, and prepare. Never execute.</p>
        </button>
        <button
          type="button"
          className={mode === "guarded" && running ? "on" : ""}
          disabled={busy || phase !== "idle"}
          onClick={() => {
            if (live) return;
            setConfirm(true);
            setReviewed(false);
          }}
        >
          <p className="label">Guarded Autonomy</p>
          <p>Research and execute only inside the pinned policy after you confirm on this computer.</p>
        </button>
      </div>

      {live ? (
        <section className="live-banner" role="status">
          <p className="label">LIVE AUTONOMY</p>
          <h2>Guarded Autonomy is on for {remain}</h2>
          <p>
            PIT may research the best eligible book and, after host gates, execute only inside the pinned policy. Chat
            cannot AUTHORIZE. Withdraw and transfer stay impossible.
          </p>
          {mission.block_reason || mission.why_not ? (
            <p className="err" role="status">
              Why PIT did not trade: {mission.why_not || humanStop(mission.block_reason || "")}. {mission.block_explain || "The mission stays alive. Scan and research continue. Existing positions are not flattened."}
            </p>
          ) : null}
          <button type="button" className="kill-switch" onClick={() => void stop()} disabled={busy || phase !== "idle"}>
            Stop autonomy now
          </button>
        </section>
      ) : m.last_stop ? (
        <section className="stop-banner" role="status">
          <p className="label">AUTONOMY STOPPED</p>
          <h2>Reason: {humanStop(m.last_stop)}</h2>
          <p>{mission.explain || "PIT will not place further orders until you enable Guarded Autonomy again on this computer."}</p>
        </section>
      ) : (
        <section className="next-row">
          <div>
            <p className="label">{status}</p>
            <h2>{mode === "research_only" ? "Research Only" : mode === "guarded" ? "Guarded Autonomy" : "Manual"}</h2>
            <p className="fine">
              Confirm Guarded Autonomy on this computer after you review the limits. The enable phrase is sent by this
              window, not by chat.
            </p>
          </div>
          <div className="cta-row">
            <button
              type="button"
              className="primary"
              disabled={busy || phase !== "idle" || kill}
              onClick={() => {
                setConfirm(true);
                setReviewed(false);
              }}
            >
              Enable Guarded Autonomy
            </button>
            <button type="button" className="danger" onClick={() => void stop()} disabled={busy || phase !== "idle" || mode === "manual"}>
              Stop autonomy
            </button>
          </div>
        </section>
      )}

      {err || mission.error ? (
        <p className="err" role="alert">
          {err || mission.error}
        </p>
      ) : null}
      {kill ? <p className="err">Kill switch is on. New orders are halted.</p> : null}

      <AwayBoard away={mission.away} whyNot={mission.why_not} whyCode={mission.why_not_code || mission.block_reason} />

      <p className="label">Mission status</p>
      <dl className="mission-grid">
        <div>
          <dt>Stage</dt>
          <dd>{humanStage(mission.stage || m.stage || (mission.research_running ? "researching" : status))}</dd>
        </div>
        <div>
          <dt>Running</dt>
          <dd>{status}</dd>
        </div>
        <div>
          <dt>Elapsed</dt>
          <dd>{elapsedLabel(m.guarded_enabled_unix, mission.now || now)}</dd>
        </div>
        <div>
          <dt>Started</dt>
          <dd>{m.guarded_enabled_unix ? new Date(m.guarded_enabled_unix * 1000).toLocaleString() : "—"}</dd>
        </div>
        <div>
          <dt>Next scan</dt>
          <dd>{nextScanLabel(mission.next_scan_unix || m.next_scan_unix, mission.now || now)}</dd>
        </div>
        <div>
          <dt>Current opportunity</dt>
          <dd>
            {m.best_coin ? (
              <ExternalLink className="asset" href={hyperliquidTrade(net || "mainnet", m.best_coin)}>
                <BrandMark symbol={m.best_coin} /> {m.best_coin}
              </ExternalLink>
            ) : (
              "none"
            )}
          </dd>
        </div>
        <div>
          <dt>Current action</dt>
          <dd>{m.last_action || "none"}</dd>
        </div>
        <div>
          <dt>Result</dt>
          <dd>{m.last_result || mission.research_stage || "—"}</dd>
        </div>
        <div>
          <dt>Universe</dt>
          <dd>
            {m.scanned || 0} scanned · {m.eligible || 0} pass policy
            {m.scan_count ? ` · ${m.scan_count} ticks` : ""}
          </dd>
        </div>
        <div>
          <dt>Trades today</dt>
          <dd>{m.trades_today || 0}</dd>
        </div>
        <div>
          <dt>Current exposure</dt>
          <dd>{m.current_position || (m.open_positions ? `${m.open_positions} open` : "none")}</dd>
        </div>
        <div>
          <dt>P&L</dt>
          <dd>Venue account</dd>
        </div>
        <div>
          <dt>Remaining risk budget</dt>
          <dd>${String(mission.remaining_risk_usd ?? limits.daily_loss_usd ?? 50)} daily · {String(mission.remaining_consecutive_losses ?? limits.max_consecutive_losses ?? 3)} losses left</dd>
        </div>
        <div>
          <dt>Stop reason</dt>
          <dd>{m.last_stop ? humanStop(m.last_stop) : live ? "none — halt only on deadline, kill, session, policy, max trades, daily loss, consecutive losses" : "deadline, kill, session, policy, max trades, daily loss, consecutive losses"}</dd>
        </div>
        <div>
          <dt>Exec gate</dt>
          <dd>{mission.block_reason ? humanStop(mission.block_reason) : "clear"}</dd>
        </div>
      </dl>
      {mission.explain && m.last_stop ? <p className="fine">{mission.explain}</p> : null}
      {m.best_why ? <p className="fine">Why it acted: {m.best_why}</p> : null}

      <p className="label">Official links</p>
      <div className="mission-links">
        <ExternalLink className="linkish" href={hyperliquidApp(net || "mainnet")}>
          Hyperliquid
        </ExternalLink>
        <ExternalLink className="linkish" href={hyperliquidAPI(net || "mainnet")}>
          Hyperliquid API
        </ExternalLink>
        {m.best_coin ? (
          <ExternalLink className="linkish" href={hyperliquidTrade(net || "mainnet", m.best_coin)}>
            {m.best_coin} book
          </ExternalLink>
        ) : null}
        {wallet ? (
          <ExternalLink className="linkish" href={explorerAddress(wallet, net || "mainnet")}>
            0G explorer
          </ExternalLink>
        ) : (
          <ExternalLink className="linkish" href={LINKS.explorer}>
            0G explorer
          </ExternalLink>
        )}
        <ExternalLink className="linkish" href={LINKS.og}>
          0G
        </ExternalLink>
        <ExternalLink className="linkish" href={LINKS.pcAdvanced}>
          Private compute
        </ExternalLink>
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

      <p className="label">Policy — immutable by the model</p>
      <div className="policy-grid">
        {limitRows(limits).map((row) => (
          <div className="policy-cell" key={row.k}>
            <p className="label">{row.k}</p>
            <strong>{row.v}</strong>
          </div>
        ))}
      </div>
      <p className="host-enforced">The model can never modify these.</p>

      <section className="cadence-row">
        <Select
          id="scan-cadence"
          label="Scan cadence"
          value={String(prefs.cadence_minutes || 15)}
          disabled={busy}
          options={CADENCE}
          onChange={(v) => onSavePrefs({ ...prefs, cadence_minutes: Number(v) })}
        />
        <p className="fine">Compute money is not trading capital. Best eligible books move into Research automatically in Research Only and Guarded Autonomy.</p>
      </section>

      {confirm ? (
        <div className="overlay" role="dialog" aria-modal="true" aria-labelledby="enable-title">
          <div className="confirm">
            <p className="label">Confirm</p>
            <h2 id="enable-title">Enable Guarded Autonomy</h2>
            <p>Review the host limits. Then confirm exactly what this window may do.</p>
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
            </ul>
            <Select label="Duration" value={String(hours)} options={HOURS} onChange={(v) => setHours(Number(v))} disabled={busy || phase !== "idle"} />
            <label className="check-row">
              <input type="checkbox" checked={reviewed} onChange={(e) => setReviewed(e.target.checked)} />
              I reviewed these limits and understand Guarded Autonomy stays inside them.
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
                {phase === "enabling" ? "Enabling…" : "Enable Guarded Autonomy"}
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
  if (s === "user_stop") return "user_stop";
  if (s === "kill_switch") return "Kill switch";
  if (s === "deadline") return "Duration ended";
  if (s === "session_expired") return "Hyperliquid session permissions no longer match this desk.";
  if (s === "max_open_positions") return "Open position ceiling";
  if (s === "daily_loss") return "Daily loss";
  if (s === "policy_changed" || s === "need_pin") return "Your policy changed. Re-pin it before trading.";
  if (s === "max_trades") return "Trade ceiling";
  if (s === "consecutive_loss_limit") return "Consecutive loss ceiling";
  if (s === "duplicate_preview") return "Duplicate preview";
  if (s === "preview_before_guarded") return "Preview started before enable";
  if (s === "insufficient_margin") return "Not enough available trading capital for this market.";
  if (s === "below_min_notional") return "This account cannot size a clip at the $10 Hyperliquid minimum.";
  if (s === "no_opportunity") return "Nothing qualifies under your law right now.";
  return s.replaceAll("_", " ");
}

function humanStage(s: string) {
  const t = String(s || "").replaceAll("_", "-");
  if (t === "scanning" || t === "starting") return "Watching live markets";
  if (t === "researching") return "Researching privately";
  if (t === "waiting" || t === "waiting after research") return "Waiting for the next scan";
  if (t === "eligible") return "Ready to trade";
  if (t === "execution-blocked" || t === "exec blocked") return "Ready to trade — host refused";
  if (t === "executing" || t === "executed") return "Submitting to Hyperliquid";
  if (t === "resting") return "Resting on the venue";
  if (t === "cooldown") return "Cooldown";
  if (t === "stopped") return "Autonomy stopped";
  if (t === "empty") return "Watching — nothing executable";
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
      <h2>What PIT did without asking again</h2>
      <p className="fine">
        {whyNot || (whyCode ? humanStop(whyCode) : "No named refusal yet.")} Guarded Autonomy never raises your limits.
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

function coinFromMarket(market?: string) {
  const m = String(market || "");
  const parts = m.split(":");
  return parts[parts.length - 1] || m;
}
