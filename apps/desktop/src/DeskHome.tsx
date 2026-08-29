import { BrandMark } from "./BrandMark";
import { LINKS } from "./links";
import { prettyCode } from "./companion";
import { compactNum } from "./format";
import type { NextFix } from "./nextFix";
import type { Probe } from "./readiness";

type Coin = {
  coin: string;
  why?: string;
  trend?: string;
  mark: number;
  eligible?: boolean;
};

function Chip({ ok, label, value }: { ok: boolean; label: string; value: string }) {
  return (
    <span className={ok ? "chip ok" : "chip fail"}>
      {label} {value}
    </span>
  );
}

export function DeskHome({
  ready,
  doing,
  items,
  attention,
  code,
  companionUp,
  sessionAlive,
  computeReady,
  protectedOk,
  policyPinned,
  hlApproved,
  researchBusy,
  awaitingAuth,
  coins,
  lastEvent,
  mode,
  exposure,
  onResearch,
  onGo,
}: {
  ready: boolean;
  doing: string;
  items: Probe[];
  attention: NextFix;
  code: string;
  companionUp: boolean;
  sessionAlive: boolean;
  computeReady: boolean;
  protectedOk: boolean;
  policyPinned: boolean;
  hlApproved: boolean;
  researchBusy?: boolean;
  awaitingAuth?: boolean;
  coins: Coin[];
  lastEvent?: string;
  mode?: string;
  exposure?: string;
  onResearch: (coin: string) => void;
  onGo: (view: "markets" | "research" | "security" | "chat" | "automation" | "portfolio") => void;
}) {
  const showPair = items.find((p) => p.id === "wallet")?.state !== "ok";
  const best = coins.find((c) => c.eligible);
  const modeLabel = mode === "guarded" ? "Guarded Autonomy" : mode === "research_only" ? "Research Only" : "Manual";
  return (
    <main className="page dense desk-home">
      <div className="page-head">
        <div>
          <p className="eyebrow">Desk</p>
          <h1>{researchBusy ? doing : awaitingAuth ? "Waiting for you" : ready ? "Ready to discover" : attention.title}</h1>
        </div>
      </div>
      <p className="lead">{doing}</p>
      <div className="chip-row" aria-label="Readiness">
        <Chip ok={protectedOk} label="Research" value={protectedOk ? "protected" : "needs protect"} />
        <Chip ok={computeReady} label="Compute" value={computeReady ? "funded" : "needs funds"} />
        <Chip ok={sessionAlive} label="Session" value={sessionAlive ? "live" : "none"} />
        <Chip ok={hlApproved} label="Hyperliquid" value={hlApproved ? "approved" : "needs approval"} />
        <Chip ok={policyPinned} label="Policy" value={policyPinned ? "pinned" : "unpinned"} />
      </div>
      <dl className="metrics">
        <div>
          <dt>What PIT is doing</dt>
          <dd>{researchBusy ? doing : awaitingAuth ? "Preview waiting" : runningCopy(mode)}</dd>
        </div>
        <div>
          <dt>Best opportunity</dt>
          <dd>{best ? `${best.coin} ${compactNum(best.mark)}` : "none"}</dd>
        </div>
        <div>
          <dt>Current exposure</dt>
          <dd>{exposure || "—"}</dd>
        </div>
        <div>
          <dt>Autonomy</dt>
          <dd>{modeLabel}</dd>
        </div>
      </dl>
      <section className="next-row">
        <div>
          <p className="label">Needs you</p>
          <h2>{attention.title}</h2>
          <p className="fine" style={{ margin: 0 }}>
            {attention.why}
          </p>
        </div>
        <div className="cta-row">
          {ready && !researchBusy && !awaitingAuth ? (
            <button type="button" className="primary" onClick={() => onGo("markets")}>
              Open Markets
            </button>
          ) : null}
          {awaitingAuth ? (
            <button type="button" className="primary" onClick={() => onGo("research")}>
              Open preview
            </button>
          ) : null}
          {researchBusy ? (
            <button type="button" className="primary" onClick={() => onGo("research")}>
              Open Research
            </button>
          ) : null}
          <button type="button" className="linkish" onClick={() => onGo("chat")}>
            Ask PIT
          </button>
          <button type="button" className="linkish" onClick={() => onGo("automation")}>
            Automation
          </button>
          {attention.href ? (
            <a className="linkish" href={attention.href} target="_blank" rel="noreferrer">
              {attention.hrefLabel || "Open official page"}
            </a>
          ) : null}
        </div>
      </section>
      {best ? (
        <section>
          <p className="label">Best opportunity</p>
          <ul className="desk-ops">
            <li>
              <BrandMark symbol={best.coin} />
              <strong>{best.coin}</strong>
              <span className="mark-num">{compactNum(best.mark)}</span>
              <span className="fine" style={{ margin: 0 }}>
                {best.why || best.trend || "In policy universe."}
              </span>
              <button type="button" className="primary" disabled={researchBusy || !computeReady} onClick={() => onResearch(best.coin)}>
                Research privately
              </button>
            </li>
          </ul>
        </section>
      ) : (
        <p className="empty">No opportunities match your policy yet. Empty is honest.</p>
      )}
      {lastEvent ? <p className="fine">Recently: {lastEvent}</p> : null}
      {showPair ? (
        <section className="next-row">
          <div>
            <p className="label">Pair this computer</p>
            <p className="pair-chip" aria-label="pairing code">
              {code ? prettyCode(code) : companionUp ? "rotating…" : "waiting for local PIT"}
            </p>
            <p className="fine" style={{ margin: 0 }}>
              The website never receives a session key.
            </p>
          </div>
          <a className="primary" href={LINKS.pair} target="_blank" rel="noreferrer">
            Open pairing
          </a>
        </section>
      ) : null}
    </main>
  );
}

function runningCopy(mode?: string) {
  if (mode === "guarded") return "Guarded Autonomy is live inside your policy.";
  if (mode === "research_only") return "Research Only — scan and prepare, never execute.";
  return "Manual. Waiting for you.";
}
