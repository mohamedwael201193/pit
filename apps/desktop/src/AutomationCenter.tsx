import { useEffect, useState } from "react";
import { BrandMark } from "./BrandMark";
import { remainLabel } from "./format";
import { Select } from "./Select";
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
  mode: string;
  running: boolean;
  host?: string;
}): string {
  if (args.phase === "enabling") return "ENABLING";
  if (args.phase === "stopping") return "STOPPING";
  if (args.error) return "FAILED";
  if (args.kill) return "BLOCKED";
  if (args.host === "ACTIVE" || (args.mode === "guarded" && args.running) || (args.mode === "research_only" && args.running)) {
    return args.host === "BLOCKED" ? "BLOCKED" : "ACTIVE";
  }
  if (args.host === "STOPPED" || args.host === "FAILED" || args.host === "BLOCKED") return args.host;
  if (args.mode === "guarded" && !args.running) return "STOPPED";
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
  const running = Boolean(mission.running || (m.running && mode !== "manual"));
  const status = lifeOf({
    kill,
    phase,
    error: err || mission.error,
    mode,
    running,
    host: mission.status,
  });
  const live = status === "ACTIVE" && mode === "guarded";
  const deadline = m.guarded_until_unix || m.deadline_unix;
  const remain = remainLabel(deadline, mission.now || now);

  async function enable() {
    setPhase("enabling");
    setErr(null);
    try {
      const r = await onEnable(ENABLE_TOKEN, hours);
      if (r?.error) {
        setErr(r.error);
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
          <h1>Mission</h1>
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
          <button type="button" className="danger" onClick={() => void stop()} disabled={busy || phase !== "idle"}>
            Stop autonomy
          </button>
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

      <p className="label">Mission status</p>
      <dl className="mission-grid">
        <div>
          <dt>Running</dt>
          <dd>{status}</dd>
        </div>
        <div>
          <dt>Started</dt>
          <dd>{m.guarded_enabled_unix ? new Date(m.guarded_enabled_unix * 1000).toLocaleString() : "—"}</dd>
        </div>
        <div>
          <dt>Next scan</dt>
          <dd>{m.next_scan_unix ? new Date(m.next_scan_unix * 1000).toLocaleTimeString() : "on cadence"}</dd>
        </div>
        <div>
          <dt>Current opportunity</dt>
          <dd>
            {m.best_coin ? (
              <span className="asset">
                <BrandMark symbol={m.best_coin} /> {m.best_coin}
              </span>
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
          <dt>Trades today</dt>
          <dd>{m.trades_today || 0}</dd>
        </div>
        <div>
          <dt>Current exposure</dt>
          <dd>{m.current_position || "venue"}</dd>
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
          <dt>Stop conditions</dt>
          <dd>{m.last_stop ? humanStop(m.last_stop) : "deadline, kill, session, policy, open positions, daily loss"}</dd>
        </div>
      </dl>
      {m.best_why ? <p className="fine">Why it acted: {m.best_why}</p> : null}

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
  if (s === "user_stop") return "Stopped on this computer";
  if (s === "kill_switch") return "Kill switch";
  if (s === "deadline") return "Duration ended";
  if (s === "session_expired") return "Session ended";
  if (s === "max_open_positions") return "Open position ceiling";
  if (s === "daily_loss") return "Daily loss";
  return s.replaceAll("_", " ");
}
