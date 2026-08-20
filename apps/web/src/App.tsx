import { usePrivy } from "@privy-io/react-auth";
import { useState } from "react";

type Net = "mainnet" | "testnet";

const CAP: Record<Net, string[]> = {
  mainnet: [
    "Production",
    "Direct TeeML on Aristotle",
    "Desk mint / authorize / revoke",
    "Transfer of Agentic ID is not live",
    "Hyperliquid mainnet",
  ],
  testnet: [
    "Full test environment",
    "Galileo chain + Hyperliquid testnet",
    "Experimental transfer only if proven",
    "Different model catalog than production",
  ],
};

export function App() {
  const { ready, authenticated, login, logout, user } = usePrivy();
  const [net, setNet] = useState<Net>("mainnet");
  const addr = user?.wallet?.address;

  return (
    <div className="shell">
      <header className="top">
        <div className="word">PIT.</div>
        <nav>
          <span className="kicker">Private Alpha OS</span>
        </nav>
      </header>
      <main className="split">
        <section className="left">
          <p className="eyebrow">YOUR DESK</p>
          <h1>Private research. Controlled execution. A desk that learns.</h1>
          <p className="lead">
            PIT never asks for a seed phrase. Your wallet stays yours. Your session cannot withdraw.
          </p>
          {!ready ? (
            <p className="state">Loading wallet connect…</p>
          ) : !authenticated ? (
            <button className="cta" onClick={login} type="button">
              Connect your wallet
            </button>
          ) : (
            <div className="card">
              <p className="label">YOUR WALLET</p>
              <p className="mono">{addr}</p>
              <button className="ghost" onClick={logout} type="button">
                Disconnect
              </button>
            </div>
          )}
          <div className="nets">
            <button type="button" className={net === "mainnet" ? "on" : ""} onClick={() => setNet("mainnet")}>
              MAINNET
            </button>
            <button type="button" className={net === "testnet" ? "on" : ""} onClick={() => setNet("testnet")}>
              TESTNET
            </button>
          </div>
          <ul className="caps">
            {CAP[net].map((line) => (
              <li key={line}>{line}</li>
            ))}
          </ul>
        </section>
        <section className="right">
          <ol className="ring">
            {[
              "Market",
              "Seal",
              "Research",
              "Challenge",
              "Risk",
              "Policy",
              "Authorize",
              "Execute",
              "Prove",
              "Learn",
            ].map((s) => (
              <li key={s}>{s}</li>
            ))}
          </ol>
          <p className="fine">
            Web can connect and inspect. Signing Hyperliquid orders happens on desktop or CLI. Session keys never enter this browser.
          </p>
        </section>
      </main>
    </div>
  );
}
