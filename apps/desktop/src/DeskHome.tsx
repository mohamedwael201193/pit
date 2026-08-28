import { BrandMark } from "./BrandMark";
import { LINKS } from "./links";
import { prettyCode } from "./companion";
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
  onResearch: (coin: string) => void;
  onGo: (view: "watch" | "research" | "security" | "settings" | "chat" | "policy") => void;
}) {
  const showPair = items.find((p) => p.id === "wallet")?.state !== "ok";
  const top = coins.filter((c) => c.eligible).slice(0, 2);
  return (
    <main className="page dense desk-home">
      <div className="page-head">
        <div>
          <p className="eyebrow">Desk</p>
          <h1>{researchBusy ? doing : awaitingAuth ? "Waiting for you" : ready ? "Ready to discover" : attention.title}</h1>
        </div>
      </div>
      <p className="lead">
        {researchBusy
          ? "Private research is running. Chat cannot AUTHORIZE."
          : awaitingAuth
            ? "An exact preview is waiting on Research. Type AUTHORIZE there."
            : "Discover a market, research it privately, then approve the exact order on this computer."}
      </p>
      <div className="chip-row" aria-label="Readiness">
        <Chip ok={protectedOk} label="Research" value={protectedOk ? "protected" : "needs protect"} />
        <Chip ok={computeReady} label="Compute" value={computeReady ? "funded" : "needs funds"} />
        <Chip ok={sessionAlive} label="Session" value={sessionAlive ? "live" : "none"} />
        <Chip ok={hlApproved} label="Hyperliquid" value={hlApproved ? "approved" : "needs approval"} />
        <Chip ok={policyPinned} label="Policy" value={policyPinned ? "pinned" : "unpinned"} />
      </div>
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
            <button type="button" className="primary" onClick={() => onGo("watch")}>
              Open Watch
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
          {attention.href ? (
            <a className="linkish" href={attention.href} target="_blank" rel="noreferrer">
              {attention.hrefLabel || "Open official page"}
            </a>
          ) : null}
          {attention.go && attention.go !== "watch" ? (
            <button type="button" className="linkish" onClick={() => onGo(attention.go!)}>
              {attention.goLabel || "Open"}
            </button>
          ) : null}
        </div>
      </section>
      {top.length ? (
        <section>
          <p className="label">Interesting now</p>
          <ul className="desk-ops">
            {top.map((c) => (
              <li key={c.coin}>
                <BrandMark symbol={c.coin} />
                <strong>{c.coin}</strong>
                <span className="mark-num">{c.mark}</span>
                <span className="fine" style={{ margin: 0 }}>
                  {c.why || c.trend || "In policy universe."}
                </span>
                <button type="button" className="primary" disabled={researchBusy || !computeReady} onClick={() => onResearch(c.coin)}>
                  Research privately
                </button>
              </li>
            ))}
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
