import { BindNote } from "./BindNote";
import { EmptyHome } from "./EmptyHome";
import { NAMED } from "./namedStates";
import { NetworkBanner } from "./NetworkBanner";
import { PermissionsCard } from "./Permissions";
import { PolicyLaw } from "./PolicyLaw";
import { Progress } from "./Progress";

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
  const net: Net = "mainnet";
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
          <BindNote />
          <PermissionsCard />
          <PolicyLaw />
          <EmptyHome />
          <Progress current="WAITING_FOR_USER" />
          <NetworkBanner net={net} />
          <p className="fine">Network: {net}. Transfer of Agentic ID is not live on mainnet.</p>
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
