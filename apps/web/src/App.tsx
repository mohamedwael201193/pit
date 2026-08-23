import { usePrivy } from "@privy-io/react-auth";
import { useEffect, useMemo, useState } from "react";
import { EmptyWatch } from "./EmptyWatch";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PolicyPanel } from "./PolicyPanel";
import { ProgressStrip } from "./ProgressStrip";
import { Ring } from "./Ring";
import { SiweBind } from "./SiweBind";
import { VerifyForm } from "./VerifyForm";

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
  { id: "session", title: "YOUR SESSION", copy: "Order and cancel only. Created on desktop or CLI. Never in this browser." },
  { id: "ready", title: "YOUR DESK", copy: "Ask, watch, authorize the exact preview, or do nothing." },
] as const;

export function App() {
  const { ready, authenticated, login, logout, user } = usePrivy();
  const [net, setNet] = useState<Net>("mainnet");
  const [hash, setHash] = useState("");
  const [root, setRoot] = useState("");
  const [route, setRoute] = useState(() => (typeof window === "undefined" ? "" : window.location.hash));
  const addr = user?.wallet?.address;
  const step = useMemo(() => (authenticated ? 1 : 0), [authenticated]);
  const explorer = net === "mainnet" ? "https://chainscan.0g.ai" : "https://chainscan-galileo.0g.ai";

  useEffect(() => {
    const onHash = () => setRoute(window.location.hash);
    window.addEventListener("hashchange", onHash);
    return () => window.removeEventListener("hashchange", onHash);
  }, []);

  const verifyOnly = route === "#verify";

  return (
    <div className="relative min-h-[100dvh]">
      <div className="grain" aria-hidden />
      <header className="flex justify-between px-10 py-7">
        <div className="text-[42px] tracking-[-0.06em] text-coral">PIT.</div>
        <a className="text-[11px] uppercase tracking-[0.18em] text-cream no-underline opacity-80" href="#verify">
          Verify
        </a>
      </header>
      <main className="grid min-h-[calc(100dvh-96px)] grid-cols-1 lg:grid-cols-[1.1fr_0.9fr]">
        <section className="px-8 py-12 lg:px-14">
          <p className="text-[11px] tracking-[0.22em] text-coral">YOUR DESK</p>
          <h1 className="mt-3 max-w-3xl text-4xl leading-[0.95] tracking-[-0.04em] lg:text-5xl">
            {verifyOnly ? "Verify a receipt on the matching explorer." : "Private research. Controlled execution. A desk that learns."}
          </h1>
          <p className="mt-5 max-w-[36rem] text-base leading-relaxed opacity-85">
            PIT never asks for a seed phrase. Your wallet stays yours. Your session cannot withdraw.
          </p>
          {!ready ? (
            <p className="mt-6">Loading wallet connect</p>
          ) : !authenticated ? (
            <button
              className="mt-6 rounded-full bg-coral px-7 py-3.5 font-semibold text-white"
              onClick={login}
              type="button"
            >
              Connect your wallet
            </button>
          ) : (
            <div className="mt-6 rounded-[20px] border border-[#222] p-5">
              <p className="text-[11px] tracking-[0.16em]">YOUR WALLET</p>
              <p className="font-mono break-all">{addr}</p>
              <button
                className="mt-3 rounded-full border border-[#333] px-7 py-3 font-semibold"
                onClick={logout}
                type="button"
              >
                Disconnect
              </button>
            </div>
          )}
          <SiweBind connected={authenticated} />
          <NetworkToggle net={net} onChange={setNet} />
          <NetworkBanner net={net} />
          {!verifyOnly && (
            <>
              <PolicyPanel />
              <EmptyWatch />
              <ProgressStrip current={authenticated ? "AUTHENTICATING" : "CONNECTING"} />
            </>
          )}
          <ul className="mt-6 list-disc pl-5">
            {CAP[net].map((line) => (
              <li key={line}>{line}</li>
            ))}
          </ul>
          <ol className="mt-8 list-none p-0">
            {STEPS.map((s, i) => (
              <li key={s.id} className={`border-t border-[#1d1f24] py-2.5 ${i <= step ? "opacity-100" : "opacity-45"}`}>
                <strong className="block text-[11px] tracking-[0.12em]">{s.title}</strong>
                <span className="text-sm opacity-80">{s.copy}</span>
              </li>
            ))}
          </ol>
        </section>
        <section className="bg-coral px-8 py-14 text-ink lg:px-14">
          <Ring />
          <p className="mt-10 max-w-[28rem]">
            Web can connect and inspect. Signing Hyperliquid orders happens on desktop or CLI. Session keys never enter
            this browser.
          </p>
          <VerifyForm hash={hash} root={root} explorer={explorer} net={net} onHash={setHash} onRoot={setRoot} />
        </section>
      </main>
    </div>
  );
}
