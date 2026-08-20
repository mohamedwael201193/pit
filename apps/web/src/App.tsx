import { usePrivy } from "@privy-io/react-auth";
import { useMemo, useState } from "react";

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

const STEPS = [
  { id: "connect", title: "YOUR WALLET", copy: "Connect. PIT never asks for a seed phrase." },
  { id: "network", title: "YOUR NETWORK", copy: "One workspace is mainnet or testnet. Never both." },
  { id: "hl", title: "YOUR TRADING ACCOUNT", copy: "Identify your Hyperliquid account. Spot USDC counts as funded." },
  { id: "policy", title: "YOUR POLICY", copy: "You set clip, leverage, and kill. The model cannot raise them." },
  { id: "session", title: "YOUR SESSION", copy: "Order and cancel only. Created on desktop or CLI — never in this browser." },
  { id: "ready", title: "YOUR DESK", copy: "Ask, watch, authorize the exact preview, or do nothing." },
] as const;

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
  const { ready, authenticated, login, logout, user } = usePrivy();
  const [net, setNet] = useState<Net>("mainnet");
  const [hash, setHash] = useState("");
  const [root, setRoot] = useState("");
  const addr = user?.wallet?.address;

  const step = useMemo(() => {
    if (!authenticated) return 0;
    return 1;
  }, [authenticated]);

  const explorer =
    net === "mainnet" ? "https://chainscan.0g.ai" : "https://chainscan-galileo.0g.ai";

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
          <ol className="onboard">
            {STEPS.map((s, i) => (
              <li key={s.id} className={i <= step ? "live" : ""}>
                <strong>{s.title}</strong>
                <span>{s.copy}</span>
              </li>
            ))}
          </ol>
        </section>
        <section className="right">
          <ol className="ring">
            {RING.map((s) => (
              <li key={s}>{s}</li>
            ))}
          </ol>
          <p className="fine">
            Web can connect and inspect. Signing Hyperliquid orders happens on desktop or CLI. Session keys never enter this browser.
          </p>
          <div className="verify">
            <p className="label">VERIFY A RECEIPT</p>
            <input
              aria-label="preview hash"
              placeholder="preview hash 0x…"
              value={hash}
              onChange={(e) => setHash(e.target.value)}
            />
            <input
              aria-label="storage root"
              placeholder="storage root 0x…"
              value={root}
              onChange={(e) => setRoot(e.target.value)}
            />
            <a className="ghost" href={`${explorer}`} target="_blank" rel="noreferrer">
              Open {net} explorer
            </a>
          </div>
        </section>
      </main>
    </div>
  );
}
