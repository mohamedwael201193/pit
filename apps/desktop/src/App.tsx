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
import { StartCards } from "./StartCards";
import { AuthorizeGate } from "./AuthorizeGate";
import { LocalSign } from "./LocalSign";

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

export function App() {
  const [net, setNet] = useState<Net>("mainnet");
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
          <StartCards />
          <BindNote />
          <AuthorizeGate sessionAlive={false} />
          <LocalSign />
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
          <p className="fine">Network: {net}. {NAMED.TRANSFER_NOT_LIVE}</p>
        </section>
        <section className="right">
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
