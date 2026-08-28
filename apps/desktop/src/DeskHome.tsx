import type { ComponentProps } from "react";
import { EmptyHome } from "./EmptyHome";
import { HyperliquidCard } from "./HyperliquidCard";
import { LINKS } from "./links";
import { prettyCode } from "./companion";
import type { NextFix } from "./nextFix";
import type { Probe } from "./readiness";

type Coin = { coin: string; reason: string; mark: number; eligible?: boolean; funding?: number; openInterest?: number };

export function DeskHome({
  ready,
  items,
  attention,
  coins,
  net,
  code,
  companionUp,
  sessionAlive,
  computeReady,
  policyPinned,
  hlApproved,
  onResearch,
  onGo,
  hl,
}: {
  ready: boolean;
  items: Probe[];
  attention: NextFix;
  coins: Coin[];
  net: string;
  code: string;
  companionUp: boolean;
  sessionAlive: boolean;
  computeReady: boolean;
  policyPinned: boolean;
  hlApproved: boolean;
  onResearch: (coin: string) => void;
  onGo: (view: "watch" | "research" | "security" | "settings") => void;
  hl: ComponentProps<typeof HyperliquidCard>;
}) {
  const eth = coins.find((c) => c.coin === "ETH");
  const showPair = !items.find((p) => p.id === "wallet") || items.find((p) => p.id === "wallet")?.state !== "ok";
  return (
    <main className="page dense">
      <p className="eyebrow">Desk</p>
      <h1>{ready ? "Desk ready" : attention.title}</h1>
      <p className="lead">What is happening, what needs you, what is ready. Chat cannot AUTHORIZE.</p>
      <div className="action-center">
        <article className="card">
          <p className="label">Action center</p>
          <dl className="status-grid">
            <dt>Private research</dt>
            <dd>{computeReady ? "READY" : "NEEDS ACTION"}</dd>
            <dt>Trading session</dt>
            <dd>{sessionAlive ? "READY" : "NEEDS ACTION"}</dd>
            <dt>Hyperliquid</dt>
            <dd>{hlApproved ? "APPROVED" : "NEEDS APPROVAL"}</dd>
            <dt>Policy</dt>
            <dd>{policyPinned ? "ACTIVE" : "UNPINNED"}</dd>
            <dt>Compute</dt>
            <dd>{computeReady ? "READY" : "NEEDS FUNDS"}</dd>
          </dl>
          <p className="fine">{attention.why}</p>
          {eth && ready ? (
            <>
              <p>
                ETH matches your policy at {eth.mark}. Research ETH privately?
              </p>
              <button type="button" className="primary" onClick={() => onResearch("ETH")}>
                Research ETH
              </button>
            </>
          ) : (
            <div className="cta-row">
              {attention.href ? (
                <a className="linkish" href={attention.href} target="_blank" rel="noreferrer">
                  {attention.hrefLabel || "Open"}
                </a>
              ) : null}
              {attention.go ? (
                <button type="button" className="linkish" onClick={() => onGo(attention.go!)}>
                  {attention.goLabel || "Open"}
                </button>
              ) : null}
            </div>
          )}
        </article>
        <article className="card">
          <p className="label">Next</p>
          <EmptyHome count={coins.filter((c) => c.eligible).length} next={attention} onGo={onGo} />
        </article>
      </div>
      {showPair ? (
        <article className="card pair-card">
          <p className="label">Pair this computer</p>
          <p className="pair-code" aria-label="pairing code">
            {code ? prettyCode(code) : companionUp ? "rotating…" : "waiting for local PIT"}
          </p>
          <p className="fine">The website never receives a session key. The code expires in two minutes.</p>
          <a className="linkish" href={LINKS.pair} target="_blank" rel="noreferrer">
            Open pairing
          </a>
        </article>
      ) : null}
      <div className="desk-grid" style={{ marginTop: 12 }}>
        {coins.slice(0, 3).map((c) => (
          <article key={c.coin} className="card">
            <p className="label">{c.coin}</p>
            <p className="mark-num">{c.mark || "—"}</p>
            <p className="fine">
              Funding {c.funding ?? "—"} · OI {c.openInterest ? Math.round(c.openInterest) : "—"} · Policy{" "}
              {c.eligible ? "PASS" : "BLOCKED"}
            </p>
            <p className="fine">{c.reason}</p>
            <button type="button" className="linkish" onClick={() => onResearch(c.coin)} disabled={!c.eligible}>
              Research
            </button>
          </article>
        ))}
      </div>
      <HyperliquidCard {...hl} />
    </main>
  );
}
