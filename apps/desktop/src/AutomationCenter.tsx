import { useState } from "react";
import type { AutoPrefs, MissionPublic } from "./companion";

export function AutomationCenter({
  mission,
  prefs,
  busy,
  kill,
  onMode,
  onEnable,
  onStop,
  onSavePrefs,
}: {
  mission: MissionPublic;
  prefs: AutoPrefs;
  busy?: boolean;
  kill?: boolean;
  onMode: (mode: "manual" | "research_only") => void;
  onEnable: (typed: string, hours: number) => void;
  onStop: () => void;
  onSavePrefs: (p: AutoPrefs) => void;
}) {
  const [typed, setTyped] = useState("");
  const [hours, setHours] = useState(24);
  const m = mission.mission || {};
  const mode = mission.mode || m.mode || "manual";
  const limits = mission.limits || {};
  const running = Boolean(mission.running || (m.running && mode !== "manual"));
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Automation</p>
          <h1>Missions</h1>
        </div>
        <p className="fine" style={{ margin: 0 }}>
          Host-enforced. The model cannot change these limits. Chat cannot enable Guarded Autonomy.
        </p>
      </div>
      <div className="mode-grid">
        <button type="button" className={mode === "manual" ? "on" : ""} disabled={busy} onClick={() => onMode("manual")}>
          <p className="label">Manual</p>
          <p>Scan and research on your command. Every order waits for AUTHORIZE.</p>
        </button>
        <button type="button" className={mode === "research_only" ? "on" : ""} disabled={busy} onClick={() => onMode("research_only")}>
          <p className="label">Research Only</p>
          <p>Scan, research, notify, and prepare. Never execute.</p>
        </button>
        <article className={mode === "guarded" && running ? "on" : ""}>
          <p className="label">Guarded Autonomy</p>
          <p>Research and execute only inside the pinned policy. Type the enable phrase after you review the limits.</p>
        </article>
      </div>
      <section className="next-row">
        <div>
          <p className="label">{running ? "RUNNING" : "IDLE"}</p>
          <h2>{mode === "guarded" && running ? "Guarded Autonomy" : mode === "research_only" ? "Research Only" : "Manual"}</h2>
          <p className="fine">
            Next scan {m.next_scan_unix ? new Date(m.next_scan_unix * 1000).toLocaleTimeString() : "on cadence"} · Current opportunity{" "}
            {m.best_coin || "none"} · Current position {m.current_position || "venue"} · Trades today {m.trades_today || 0} · Last action{" "}
            {m.last_action || "none"}
          </p>
          {m.last_stop ? <p>Stopped because: {m.last_stop}</p> : null}
          {kill ? <p>Kill switch is on. New orders are halted.</p> : null}
        </div>
      </section>
      <section>
        <p className="label">Host limits (immutable by the model)</p>
        <table className="desk-table">
          <tbody>
            <tr>
              <td>Allowed assets</td>
              <td>{Array.isArray(limits.allowed_assets) ? (limits.allowed_assets as string[]).join(" ") : "ETH BTC SOL HYPE DOGE AVAX"}</td>
            </tr>
            <tr>
              <td>Allowed venues</td>
              <td>{Array.isArray(limits.allowed_venues) ? (limits.allowed_venues as string[]).join(" ") : "hyperliquid"}</td>
            </tr>
            <tr>
              <td>Max clip / leverage / open</td>
              <td>
                ${String(limits.max_clip_usd ?? 10)} · {String(limits.max_leverage ?? 1)}x · {String(limits.max_open_positions ?? 1)} open
              </td>
            </tr>
            <tr>
              <td>Daily loss / consecutive losses</td>
              <td>
                ${String(limits.daily_loss_usd ?? 50)} · {String(limits.max_consecutive_losses ?? 3)} streak
              </td>
            </tr>
            <tr>
              <td>Withdraw / transfer / policy mutation</td>
              <td>Forbidden</td>
            </tr>
          </tbody>
        </table>
      </section>
      <section>
        <p className="label">Enable Guarded Autonomy</p>
        <p>Review the limits. Type ENABLE GUARDED AUTONOMY. PIT will not execute from chat.</p>
        <form
          className="auth-form"
          onSubmit={(e) => {
            e.preventDefault();
            onEnable(typed, hours);
            setTyped("");
          }}
        >
          <input value={typed} onChange={(e) => setTyped(e.target.value)} autoComplete="off" spellCheck={false} placeholder="ENABLE GUARDED AUTONOMY" disabled={busy} />
          <select value={hours} onChange={(e) => setHours(Number(e.target.value))} disabled={busy} aria-label="Hours">
            <option value={1}>1 hour</option>
            <option value={8}>8 hours</option>
            <option value={24}>24 hours</option>
            <option value={72}>72 hours</option>
          </select>
          <button type="submit" className="primary" disabled={busy || typed.trim() !== "ENABLE GUARDED AUTONOMY"}>
            Enable
          </button>
        </form>
        <button type="button" className="linkish" onClick={onStop} disabled={busy || mode === "manual"}>
          Stop autonomy
        </button>
      </section>
      <section>
        <p className="label">Scan cadence</p>
        <label className="fine">
          Cadence
          <select
            value={String(prefs.cadence_minutes || 15)}
            disabled={busy}
            onChange={(e) => onSavePrefs({ ...prefs, cadence_minutes: Number(e.target.value) })}
          >
            <option value="1">Every minute</option>
            <option value="5">Every 5 minutes</option>
            <option value="15">Every 15 minutes</option>
            <option value="60">Every hour</option>
          </select>
        </label>
        <p className="fine">
          Why it acted: {m.best_why || m.last_action || "No mission action yet."} Compute money is not trading capital.
        </p>
      </section>
    </main>
  );
}
