import { LINKS } from "./links";
import { prettyCode } from "./companion";
import type { NextFix } from "./nextFix";
import type { Probe } from "./readiness";

const LOOP = [
  "Watch",
  "Research privately",
  "Challenge",
  "Explain",
  "Policy check",
  "Exact preview",
  "You approve",
  "Execute",
  "Prove",
  "Learn",
] as const;

function tone(ok: boolean) {
  return ok ? "ok" : "bad";
}

function loopIndex(ready: boolean, attention: NextFix): number {
  const t = attention.title.toLowerCase();
  if (!ready && t.includes("wallet")) return 0;
  if (t.includes("protect")) return 1;
  if (t.includes("session") || t.includes("approve")) return 6;
  if (t.includes("policy")) return 4;
  if (t.includes("fund") || t.includes("research")) return 1;
  if (ready) return 0;
  return 0;
}

export function DeskHome({
  ready,
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
  onResearch,
  onGo,
}: {
  ready: boolean;
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
  onResearch: (coin: string) => void;
  onGo: (view: "watch" | "research" | "security" | "settings") => void;
}) {
  const showPair = items.find((p) => p.id === "wallet")?.state !== "ok";
  let lit = loopIndex(ready, attention);
  if (researchBusy) lit = 1;
  if (awaitingAuth) lit = 6;
  return (
    <main className="page dense desk-story">
      <div className="page-head">
        <div>
          <p className="eyebrow">Desk</p>
          <h1>{ready ? "Private trading desk" : attention.title}</h1>
        </div>
        <p className="fine" style={{ margin: 0 }}>
          Chat cannot AUTHORIZE.
        </p>
      </div>
      <ol className="loop" aria-label="How PIT works">
        {LOOP.map((step, i) => (
          <li key={step} className={i === lit ? "on" : ""}>
            <span>{String(i + 1).padStart(2, "0")}</span>
            {step}
          </li>
        ))}
      </ol>
      <p className="lead">
        PIT watches the public book under your policy, researches your private strategy in sealed compute, challenges
        itself, then waits for you to approve the exact order.
      </p>
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
          {ready ? (
            <button type="button" className="primary" onClick={() => onGo("watch")}>
              Open Watch
            </button>
          ) : null}
          {ready ? (
            <button type="button" className="linkish" onClick={() => onResearch("ETH")}>
              Research ETH privately
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
    </main>
  );
}
