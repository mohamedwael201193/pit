import { useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { ArrowLeft } from "@phosphor-icons/react";
import { PageHead } from "../ui/PageHead";
import { Button, ButtonLink } from "../ui/Button";
import { NetworkBanner } from "../NetworkBanner";
import { PolicyPanel } from "../PolicyPanel";
import { namedState } from "../namedStates";
import { cn } from "../lib/cn";

type Net = "mainnet" | "testnet";

const STEPS = [
  {
    id: 1,
    title: "Connect wallet",
    why: "Your address is the desk identity. PIT never asks for a seed phrase.",
    action: "Connect with Privy. You already did if you can see this screen.",
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
    action: "Open Hyperliquid in your wallet. PIT reads public state only from the web.",
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
    action: "Approve on Hyperliquid. extraAgents must list PIT and the workspace.",
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
  const navigate = useNavigate();
  const [step, setStep] = useState(authenticated ? 2 : 1);
  const [net, setNet] = useState<Net>("mainnet");
  const current = useMemo(() => STEPS.find((s) => s.id === step) ?? STEPS[0], [step]);

  return (
    <div className="mx-auto max-w-[80rem]">
      <button
        type="button"
        onClick={() => (step <= 1 ? navigate("/app") : setStep((s) => Math.max(1, s - 1)))}
        className="mb-6 inline-flex items-center gap-2 text-[0.875rem] text-[rgb(240_231_212/0.6)] hover:text-[var(--guide-cream)]"
      >
        <ArrowLeft size={16} aria-hidden="true" />
        Back
      </button>
      <PageHead title={current.title} lede={current.why} />

      {current.id === 2 ? (
        <div className="mt-10 grid gap-4 sm:grid-cols-2">
          <ChoiceCard
            title="MAINNET"
            body="Production. Aristotle and Hyperliquid mainnet. Transfer of Agentic ID is not live."
            onClick={() => setNet("mainnet")}
            active={net === "mainnet"}
          />
          <ChoiceCard
            title="TESTNET"
            body="Protocol laboratory. Galileo and Hyperliquid testnet. Sealed ask stays off until proven."
            onClick={() => setNet("testnet")}
            active={net === "testnet"}
          />
        </div>
      ) : null}

      {current.id === 2 ? <NetworkBanner net={net} /> : null}

      {current.id === 5 ? (
        <div className="mt-10 grid gap-4 sm:grid-cols-2">
          <ChoiceCard title="10 USD clip" body="Default until you raise it on purpose." onClick={() => undefined} active />
          <ChoiceCard
            title="Raise later"
            body="Change clip on desktop or CLI. The model cannot raise it."
            onClick={() => undefined}
          />
        </div>
      ) : null}

      {current.id === 6 ? <PolicyPanel /> : null}

      {current.id !== 2 && current.id !== 5 && current.id !== 6 ? (
        <div className="mt-10 max-w-[40rem] rounded-2xl border border-[rgb(240_231_212/0.25)] bg-[#141414] p-6">
          <p className="text-[0.8125rem] text-[rgb(240_231_212/0.55)]">
            {current.id} of 12
          </p>
          <p className="mt-4 text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.8)]">
            <strong className="text-[var(--guide-cream)]">Do this. </strong>
            {current.action}
          </p>
          <p className="mt-4 text-[0.875rem] text-[rgb(240_231_212/0.6)]">
            If it fails: {current.fail.title}. {current.fail.next}
          </p>
        </div>
      ) : (
        <p className="mt-8 max-w-[48ch] text-[0.9375rem] text-[rgb(240_231_212/0.65)]">
          Do this: {current.action} If it fails: {current.fail.title}. {current.fail.next}
        </p>
      )}

      <div className="mt-10 flex flex-wrap gap-3">
        <Button variant="secondary" type="button" disabled={step <= 1} onClick={() => setStep((s) => Math.max(1, s - 1))}>
          Retry previous
        </Button>
        {step < 12 ? (
          <Button type="button" trailingArrow onClick={() => setStep((s) => Math.min(12, s + 1))}>
            Continue
          </Button>
        ) : (
          <ButtonLink as={Link} to="/app" trailingArrow>
            Open Home
          </ButtonLink>
        )}
      </div>
    </div>
  );
}

function ChoiceCard({
  title,
  body,
  onClick,
  active,
}: {
  title: string;
  body: string;
  onClick: () => void;
  active?: boolean;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "rounded-2xl border bg-[#141414] p-6 text-left transition-colors",
        "hover:border-[#d82f2f]/50 active:scale-[0.99]",
        active ? "border-[#d82f2f]" : "border-[rgb(240_231_212/0.25)]",
      )}
    >
      <h3 className="text-[1.25rem] font-semibold tracking-[-0.03em] text-[var(--guide-cream)]">{title}</h3>
      <p className="mt-2 text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">{body}</p>
    </button>
  );
}
