import { NAMED, PERMISSIONS } from "./namedStates";

type Net = "mainnet" | "testnet";

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
          <div className="card">
            <p className="label">YOUR SESSION</p>
            <ul className="perms">
              {PERMISSIONS.map((p) => (
                <li key={p.k}>
                  {p.k} {p.ok ? "allowed" : "denied"}
                </li>
              ))}
            </ul>
          </div>
          <p className="fine">Network: {net}. Transfer of Agentic ID is not live on mainnet.</p>
        </section>
        <section className="right">
          <ol className="ring">
            {["PRIVATE_BOOK", "SEALING", "TEE", "TEE_SIGNATURE", "ONCHAIN_SIGNER", "STORAGE", "RECEIPT", "CALIBRATION"].map(
              (s) => (
                <li key={s}>{s}</li>
              ),
            )}
          </ol>
          <p className="fine">{NAMED.TEE_VERIFY_FAIL}</p>
        </section>
      </main>
    </div>
  );
}
