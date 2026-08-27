import { useState } from "react";
import { BindNote } from "./BindNote";
import { EmptyHome } from "./EmptyHome";
import { RecoverNote } from "./RecoverNote";
import { IsolateNote } from "./IsolateNote";
import { KillNote } from "./KillNote";
import { NAMED } from "./namedStates";
import { NoSession } from "./NoSession";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PermissionsCard } from "./Permissions";
import { PolicyLaw } from "./PolicyLaw";
import { Progress } from "./Progress";
import { AuthorizeGate } from "./AuthorizeGate";
import { LocalSign } from "./LocalSign";
import { SessionNote } from "./SessionNote";
import { PreviewNote } from "./PreviewNote";
import { LedgerNote } from "./LedgerNote";
import { CancelNote } from "./CancelNote";
import { SignedNote } from "./SignedNote";
import { LinkedNote } from "./LinkedNote";
import { StatusNote } from "./StatusNote";
import { PostedNote } from "./PostedNote";

type Net = "mainnet" | "testnet";

const RING = [
  "PRIVATE_BOOK",
  "SEALING",
  "TEE",
  "TEE_SIGNATURE",
  "ONCHAIN_SIGNER",
  "STORAGE",
  "RECEIPT",
  "CALIBRATION",
];

const PIPELINE = [
  "MARKET",
  "PRIVATE",
  "SEALED",
  "RESEARCH",
  "CHALLENGE",
  "RISK",
  "POLICY",
  "AUTHORIZE",
  "EXECUTE",
  "PROVE",
  "LEARN",
];

const BEATS = [
  { title: "Connect wallet", body: "PIT never asks for a seed phrase." },
  { title: "Select network", body: "MAINNET production or TESTNET lab. Never mix." },
  { title: "Workspace", body: "Two wallets never share a workspace." },
  { title: "Hyperliquid", body: "Spot USDC counts as funded." },
  { title: "Capital", body: "You set the clip. The model cannot raise it." },
  { title: "Policy", body: "Readable cards. Kill switch is yours." },
  { title: "Local session", body: "One hour. Order and cancel only. On this machine." },
  { title: "Approve agent", body: "Your wallet approves the printed address." },
  { title: "Permissions", body: "Withdraw denied. Leverage denied." },
  { title: "Research only", body: "Watch and ask without sending an order." },
  { title: "Tiny test trade", body: "Type AUTHORIZE on the exact preview." },
  { title: "Ready", body: "The desk hunts. You authorize. Receipts verify." },
];

export function App() {
  const [net, setNet] = useState<Net>("mainnet");
  const [beat, setBeat] = useState(0);
  return (
    <div className="shell">
      <header className="top">
        <div className="word">PIT.</div>
        <nav>
          <span className="kicker">Desktop · session lives here</span>
        </nav>
      </header>
      <main className="split">
        <section className="left">
          <p className="eyebrow">YOUR DESK</p>
          <h1>Authorize on this machine. Never in the browser.</h1>
          <p className="lead">{NAMED.SEED_FORBIDDEN}</p>
          <p className="lead">Your wallet stays yours. Your trading session cannot withdraw.</p>
          <div className="start">
            {BEATS.slice(0, 2).map((c) => (
              <article key={c.title}>
                <p className="label">{c.title.toUpperCase()}</p>
                <p>{c.body}</p>
              </article>
            ))}
          </div>
          <ol className="beats">
            {BEATS.map((b, i) => (
              <li key={b.title}>
                <button type="button" className={i === beat ? "on" : ""} onClick={() => setBeat(i)}>
                  {b.title}
                </button>
              </li>
            ))}
          </ol>
          <p className="fine">
            <strong>{BEATS[beat].title}.</strong> {BEATS[beat].body} Retry from this list. Recovery is on this machine.
          </p>
          <BindNote />
          <AuthorizeGate sessionAlive={false} />
          <LocalSign />
          <SessionNote />
          <PreviewNote />
          <LedgerNote />
          <CancelNote />
          <SignedNote />
          <LinkedNote />
          <StatusNote />
          <PostedNote />
          <NoSession />
          <RecoverNote />
          <PermissionsCard />
          <PolicyLaw />
          <EmptyHome />
          <Progress current="WAITING_FOR_USER" />
          <NetworkToggle net={net} onChange={setNet} />
          <NetworkBanner net={net} />
          <IsolateNote />
          <KillNote />
          <p className="fine">{NAMED.TWO_WALLETS}</p>
          <p className="fine">
            Network: {net}. {NAMED.TRANSFER_NOT_LIVE}
          </p>
        </section>
        <section className="right">
          <svg viewBox="0 0 240 240" className="mark" aria-hidden="true">
            {Array.from({ length: 32 }, (_, i) => {
              const a = (i / 32) * Math.PI * 2 - Math.PI / 2;
              return (
                <line
                  key={i}
                  x1={120 + 70 * Math.cos(a)}
                  y1={120 + 70 * Math.sin(a)}
                  x2={120 + 108 * Math.cos(a)}
                  y2={120 + 108 * Math.sin(a)}
                  stroke="#0a0a0a"
                  strokeWidth="2.4"
                  strokeLinecap="round"
                />
              );
            })}
            <circle cx="120" cy="120" r="26" fill="none" stroke="#0a0a0a" strokeWidth="2.2" />
            <circle cx="120" cy="120" r="8" fill="#0a0a0a" />
          </svg>
          <ol className="pipe">
            {PIPELINE.map((s) => (
              <li key={s}>{s}</li>
            ))}
          </ol>
          <ol className="ring">
            {RING.map((s) => (
              <li key={s}>{s}</li>
            ))}
          </ol>
          <p className="fine">{NAMED.TEE_VERIFY_FAIL}</p>
        </section>
      </main>
    </div>
  );
}
