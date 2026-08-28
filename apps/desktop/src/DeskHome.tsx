import { LINKS } from "./links";
import { prettyCode } from "./companion";
import type { NextFix } from "./nextFix";
import type { Probe } from "./readiness";

type Coin = { coin: string; reason: string; mark: number; eligible?: boolean; funding?: number; openInterest?: number; oracle?: number };

function tone(ok: boolean) {
  return ok ? "ok" : "bad";
}

export function DeskHome({
  ready,
  items,
  attention,
  coins,
  code,
  companionUp,
  sessionAlive,
  computeReady,
  protectedOk,
  policyPinned,
  hlApproved,
  onResearch,
  onGo,
}: {
  ready: boolean;
  items: Probe[];
  attention: NextFix;
  coins: Coin[];
  code: string;
  companionUp: boolean;
  sessionAlive: boolean;
  computeReady: boolean;
  protectedOk: boolean;
  policyPinned: boolean;
  hlApproved: boolean;
  onResearch: (coin: string) => void;
  onGo: (view: "watch" | "research" | "security" | "settings") => void;
}) {
  const eth = coins.find((c) => c.coin === "ETH");
  const showPair = items.find((p) => p.id === "wallet")?.state !== "ok";
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Desk</p>
          <h1>{ready ? "Ready" : attention.title}</h1>
        </div>
        <p className="fine" style={{ margin: 0 }}>
          Chat cannot AUTHORIZE.
        </p>
      </div>
      <dl className="ready-strip">
        <div>
          <dt>Research</dt>
          <dd className={tone(protectedOk)}>{protectedOk ? "PROTECTED" : "NEEDS PROTECT"}</dd>
        </div>
        <div>
          <dt>Compute</dt>
          <dd className={tone(computeReady)}>{computeReady ? "FUNDED" : "NEEDS FUNDS"}</dd>
        </div>
        <div>
          <dt>Session</dt>
          <dd className={tone(sessionAlive)}>{sessionAlive ? "LIVE" : "NONE"}</dd>
        </div>
        <div>
          <dt>Hyperliquid</dt>
          <dd className={tone(hlApproved)}>{hlApproved ? "APPROVED" : "NEEDS APPROVAL"}</dd>
        </div>
        <div>
          <dt>Policy</dt>
          <dd className={tone(policyPinned)}>{policyPinned ? "ACTIVE" : "UNPINNED"}</dd>
        </div>
      </dl>
      <section className="next-row">
        <div>
          <p className="label">Next</p>
          <h2>{attention.title}</h2>
          <p className="fine" style={{ margin: 0 }}>
            {attention.why}
          </p>
        </div>
        <div className="cta-row">
          {ready && eth ? (
            <button type="button" className="primary" onClick={() => onResearch("ETH")}>
              Research ETH
            </button>
          ) : null}
          {attention.href ? (
            <a className="linkish" href={attention.href} target="_blank" rel="noreferrer">
              {attention.hrefLabel || "Open official page"}
            </a>
          ) : null}
          {attention.go ? (
            <button type="button" className="linkish" onClick={() => onGo(attention.go!)}>
              {attention.goLabel || "Open"}
            </button>
          ) : null}
        </div>
      </section>
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
      <p className="label" style={{ marginTop: 14 }}>
        Markets
      </p>
      <div className="market-head">
        <span>Asset</span>
        <span>Mark</span>
        <span>Oracle</span>
        <span>Funding</span>
        <span>OI</span>
        <span>Policy</span>
        <span>Ready</span>
        <span></span>
      </div>
      <ul className="market-rows" aria-label="Policy markets">
        {coins.slice(0, 3).map((c) => (
          <li key={c.coin}>
            <strong>{c.coin}</strong>
            <span className="mark-num">{c.mark || "—"}</span>
            <span>{c.oracle ?? "—"}</span>
            <span>{c.funding ?? "—"}</span>
            <span>{c.openInterest ? Math.round(c.openInterest) : "—"}</span>
            <span>{c.eligible ? "PASS" : "BLOCKED"}</span>
            <span>{computeReady ? "Ready" : "Needs compute"}</span>
            <button type="button" className="linkish" onClick={() => onResearch(c.coin)} disabled={!c.eligible}>
              Research
            </button>
          </li>
        ))}
      </ul>
      <p className="fine">Public Hyperliquid marks. Side is not decided here.</p>
    </main>
  );
}
