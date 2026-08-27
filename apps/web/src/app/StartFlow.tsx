import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { NetworkToggle } from "../NetworkToggle";
import { NetworkBanner } from "../NetworkBanner";
import { PolicyPanel } from "../PolicyPanel";
import { namedState } from "../namedStates";

type Net = "mainnet" | "testnet";

const STEPS = [
  {
    id: 1,
    title: "Connect wallet",
    why: "Your address is the desk identity. PIT never asks for a seed phrase.",
    action: "Connect with Privy.",
    fail: namedState("SIGNATURE_DECLINED"),
  },
  {
    id: 2,
    title: "Select network",
    why: "One workspace is MAINNET production or TESTNET lab. Never both.",
    action: "Pick MAINNET or TESTNET and stay there.",
    fail: namedState("WRONG_NETWORK"),
  },
  {
    id: 3,
    title: "Create workspace",
    why: "Policy, memory, and ledger bind to this workspace, not to a global master address.",
    action: "Confirm the workspace created on connect.",
    fail: namedState("WRONG_NETWORK"),
  },
  {
    id: 4,
    title: "Connect Hyperliquid",
    why: "PIT needs YOUR trading account. Spot USDC counts as funded.",
    action: "Open Hyperliquid in your wallet's network. PIT reads public state only from the web.",
    fail: namedState("BACKEND_UNREACHABLE"),
  },
  {
    id: 5,
    title: "Choose capital",
    why: "You decide what this desk may touch. PIT does not size from a model.",
    action: "Keep clip at 10 USD until you raise it on purpose.",
    fail: namedState("WRONG_NETWORK"),
  },
  {
    id: 6,
    title: "Set first policy",
    why: "The law is readable. The model cannot raise clip, leverage, or permissions.",
    action: "Review the cards. Change them on desktop or CLI if needed.",
    fail: namedState("WRONG_NETWORK"),
  },
  {
    id: 7,
    title: "Create local session",
    why: "The one-hour agent lives in the OS keychain. Not here.",
    action: "Run pit session on desktop or CLI.",
    fail: namedState("SESSION_EXPIRED"),
  },
  {
    id: 8,
    title: "Approve scoped agent",
    why: "Your wallet must approveAgent the printed address. Withdraw stays denied.",
    action: "Approve on Hyperliquid. extraAgents must list PIT-{workspace}.",
    fail: namedState("SESSION_EXPIRED"),
  },
  {
    id: 9,
    title: "Permissions review",
    why: "Order yes. Cancel yes. Withdraw no. Leverage no.",
    action: "Read the card. Walk away if it is wrong.",
    fail: namedState("AUTHORIZE_WEB_DENIED"),
  },
  {
    id: 10,
    title: "Research-only test",
    why: "Watch and ask can run without sending an order.",
    action: "Open Home. Empty Watch is honest if nothing matches.",
    fail: namedState("BACKEND_UNREACHABLE"),
  },
  {
    id: 11,
    title: "Optional tiny test trade",
    why: "A dust fill is YOUR AUTHORIZE on desktop. This browser cannot sign it.",
    action: "Open PIT desktop. Type AUTHORIZE on the exact preview.",
    fail: namedState("AUTHORIZE_WEB_DENIED"),
  },
  {
    id: 12,
    title: "Ready",
    why: "The desk hunts. You authorize. Receipts verify.",
    action: "Go to Home.",
    fail: namedState("BACKEND_UNREACHABLE"),
  },
] as const;

export function StartFlow() {
  const { authenticated } = usePrivy();
  const [step, setStep] = useState(authenticated ? 2 : 1);
  const [net, setNet] = useState<Net>("mainnet");
  const current = useMemo(() => STEPS.find((s) => s.id === step) ?? STEPS[0], [step]);

  return (
    <div className="mx-auto max-w-[80rem]">
      <Link to="/app" className="mb-6 inline-block text-[0.875rem] text-[rgb(240_231_212/0.6)]">
        Back
      </Link>
      <p className="text-[11px] tracking-[0.18em] text-[#d82f2f]">GET STARTED</p>
      <h1 className="mt-2 text-4xl tracking-[-0.04em]">Get started</h1>
      <p className="mt-3 max-w-[48ch] text-[1.05rem] leading-7 text-[rgb(240_231_212/0.75)]">
        Twelve beats. Each one has a reason, an action, and a recovery. None of them ask for a seed phrase.
      </p>
      <ol className="mt-10 grid gap-3 sm:grid-cols-2">
        {STEPS.map((s) => (
          <li key={s.id}>
            <button
              type="button"
              onClick={() => setStep(s.id)}
              className={`w-full rounded-2xl border p-5 text-left ${
                s.id === step ? "border-[#d82f2f] bg-[#141414]" : "border-[rgb(240_231_212/0.22)]"
              }`}
            >
              <p className="font-semibold">{s.title}</p>
              <p className="mt-1 text-[0.875rem] text-[rgb(240_231_212/0.65)]">{s.why}</p>
            </button>
          </li>
        ))}
      </ol>
      <article className="mt-10 max-w-[40rem] border border-[rgb(240_231_212/0.22)] p-6">
        <p className="text-[11px] tracking-[0.16em] text-[#d82f2f]">STEP {current.id} OF 12</p>
        <h2 className="mt-2 text-2xl">{current.title}</h2>
        <p className="mt-3 text-[rgb(240_231_212/0.8)]">{current.why}</p>
        <p className="mt-4">
          <strong>Do this:</strong> {current.action}
        </p>
        <p className="mt-3 text-[0.875rem] text-[rgb(240_231_212/0.6)]">
          If it fails: {current.fail.title}. {current.fail.next}
        </p>
        {current.id === 2 ? (
          <div className="mt-4">
            <NetworkToggle net={net} onChange={setNet} />
            <NetworkBanner net={net} />
          </div>
        ) : null}
        {current.id === 6 ? <PolicyPanel /> : null}
        <div className="mt-6 flex flex-wrap gap-3">
          <button
            className="pill pill-line"
            type="button"
            disabled={step <= 1}
            onClick={() => setStep((s) => Math.max(1, s - 1))}
          >
            Retry previous
          </button>
          {step < 12 ? (
            <button className="pill pill-coral" type="button" onClick={() => setStep((s) => Math.min(12, s + 1))}>
              Continue
            </button>
          ) : (
            <Link to="/app" className="pill pill-coral">
              Open Home
            </Link>
          )}
        </div>
      </article>
    </div>
  );
}
