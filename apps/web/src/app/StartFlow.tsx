import { useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { ArrowLeft } from "@phosphor-icons/react";
import type { ComponentType } from "react";
import { PageHead } from "../ui/PageHead";
import { Button, ButtonLink } from "../ui/Button";
import { ChoiceCard } from "../ui/ChoiceCard";
import { Bezel } from "../ui/Surface";
import { NetworkBanner } from "../NetworkBanner";
import { PolicyPanel } from "../PolicyPanel";
import { DirectSign } from "../DirectSign";
import { namedState } from "../namedStates";
import { cn } from "../lib/cn";
import {
  DiagramAuthorize,
  DiagramHeroPostcard,
  DiagramHyperliquid,
  DiagramLearn,
  DiagramMainnet,
  DiagramPolicy,
  DiagramPrivate,
  DiagramSealed,
  DiagramSession,
  DiagramTestnet,
} from "../diagrams/pitGuide";

type Net = "mainnet" | "testnet";

const STEPS = [
  {
    id: 1,
    title: "Connect wallet",
    why: "Your address is the desk identity. PIT never asks for a seed phrase.",
    action: "You already did this if you can see this screen.",
    fail: namedState("SIGNATURE_DECLINED"),
    Diagram: DiagramHeroPostcard,
  },
  {
    id: 2,
    title: "Get started",
    why: "Pick MAINNET production or TESTNET lab. You'll set policy next.",
    action: "Tap a card. One workspace. Never both.",
    fail: namedState("WRONG_NETWORK"),
    Diagram: DiagramMainnet,
  },
  {
    id: 3,
    title: "Protect my strategy",
    why: "Pair this browser at /pair, then sign. PIT sends your private strategy only through 0G’s verified sealed path. The website never receives the token.",
    action: "Open pairing, then sign the Direct message from the bound wallet.",
    fail: namedState("SIGNATURE_DECLINED"),
    Diagram: DiagramSealed,
  },
  {
    id: 4,
    title: "Connect Hyperliquid",
    why: "PIT needs YOUR trading account. Spot USDC counts as funded.",
    action: "Open Hyperliquid in your wallet. PIT reads public state only from the web.",
    fail: namedState("BACKEND_UNREACHABLE"),
    Diagram: DiagramHyperliquid,
  },
  {
    id: 5,
    title: "Choose capital",
    why: "You decide what this desk may touch. PIT does not size from a model.",
    action: "Keep clip at 10 USD until you raise it on purpose.",
    fail: namedState("WRONG_NETWORK"),
    Diagram: DiagramPolicy,
  },
  {
    id: 6,
    title: "Set first policy",
    why: "The law is readable. The model cannot raise clip, leverage, or permissions.",
    action: "Read the cards. Change them on desktop or CLI if needed.",
    fail: namedState("WRONG_NETWORK"),
    Diagram: DiagramPolicy,
  },
  {
    id: 7,
    title: "Create local session",
    why: "The one-hour agent lives in the OS keychain. Not here.",
    action: "Create the session in PIT Desktop, then approve the printed agent on Hyperliquid API (app.hyperliquid.xyz/API).",
    fail: namedState("SESSION_EXPIRED"),
    Diagram: DiagramSession,
  },
  {
    id: 8,
    title: "Approve scoped agent",
    why: "Your wallet must approveAgent the printed address. Withdraw stays denied.",
    action: "Approve on Hyperliquid API. extraAgents must list PIT. Withdraw stays denied.",
    fail: namedState("SESSION_EXPIRED"),
    Diagram: DiagramAuthorize,
  },
  {
    id: 9,
    title: "Permissions review",
    why: "Order yes. Cancel yes. Withdraw no. Leverage no.",
    action: "Read the card. Walk away if it is wrong.",
    fail: namedState("AUTHORIZE_WEB_DENIED"),
    Diagram: DiagramAuthorize,
  },
  {
    id: 10,
    title: "Research-only test",
    why: "Watch and ask can run without sending an order.",
    action: "Open Home. Empty Watch is honest if nothing matches.",
    fail: namedState("BACKEND_UNREACHABLE"),
    Diagram: DiagramPrivate,
  },
  {
    id: 11,
    title: "Optional tiny test trade",
    why: "A dust fill is YOUR AUTHORIZE on desktop. This browser cannot sign it.",
    action: "Open PIT desktop. Type AUTHORIZE on the exact preview.",
    fail: namedState("AUTHORIZE_WEB_DENIED"),
    Diagram: DiagramAuthorize,
  },
  {
    id: 12,
    title: "Ready",
    why: "The desk hunts. You authorize. Receipts verify.",
    action: "Go to Home.",
    fail: namedState("BACKEND_UNREACHABLE"),
    Diagram: DiagramLearn,
  },
] as const;

export function StartFlow() {
  const { authenticated } = usePrivy();
  const navigate = useNavigate();
  const reduce = useReducedMotion();
  const [step, setStep] = useState(authenticated ? 2 : 1);
  const [net, setNet] = useState<Net>("mainnet");
  const current = useMemo(() => STEPS.find((s) => s.id === step) ?? STEPS[0], [step]);
  const Diagram = current.Diagram as ComponentType<{ className?: string }>;

  const back = () => (step <= 1 ? navigate("/app") : setStep((s) => Math.max(1, s - 1)));
  const next = () => setStep((s) => Math.min(12, s + 1));

  return (
    <div className="mx-auto max-w-[80rem]">
      <button
        type="button"
        onClick={back}
        className="mb-6 inline-flex items-center gap-2 text-[0.875rem] text-[rgb(240_231_212/0.6)] hover:text-[var(--guide-cream)]"
      >
        <ArrowLeft size={16} aria-hidden="true" />
        Back
      </button>

      <ol className="mb-8 flex flex-wrap gap-1.5" aria-label="Onboarding progress">
        {STEPS.map((s) => (
          <li
            key={s.id}
            className={cn(
              "h-1.5 w-6",
              s.id < step ? "bg-[#d82f2f]" : s.id === step ? "bg-[var(--guide-cream)]" : "bg-[rgb(240_231_212/0.2)]",
            )}
          />
        ))}
      </ol>

      <AnimatePresence mode="wait">
        <motion.div
          key={current.id}
          initial={reduce ? false : { opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          exit={reduce ? undefined : { opacity: 0, y: -12 }}
          transition={{ duration: 0.35, ease: [0.16, 1, 0.3, 1] }}
        >
          <PageHead title={current.title} lede={current.why} />

          {current.id === 2 ? (
            <>
              <div className="mt-10 grid gap-4 sm:grid-cols-2">
                <ChoiceCard
                  title="MAINNET"
                  body="Production. Aristotle and Hyperliquid mainnet. Transfer of Agentic ID is not live."
                  Diagram={DiagramMainnet}
                  active={net === "mainnet"}
                  onClick={() => {
                    setNet("mainnet");
                    next();
                  }}
                />
                <ChoiceCard
                  title="TESTNET"
                  body="Protocol laboratory. Galileo and Hyperliquid testnet. Sealed ask stays off until proven."
                  Diagram={DiagramTestnet}
                  active={net === "testnet"}
                  onClick={() => {
                    setNet("testnet");
                    next();
                  }}
                />
              </div>
              <NetworkBanner net={net} />
            </>
          ) : null}

          {current.id === 3 ? (
            <div className="mt-10">
              <ButtonLink as={Link} to="/pair" size="lg">
                Pair PIT Desktop
              </ButtonLink>
              <DirectSign />
              <div className="mt-8">
                <Button type="button" trailingArrow size="lg" onClick={next}>
                  Continue
                </Button>
              </div>
            </div>
          ) : null}

          {current.id === 5 ? (
            <div className="mt-10 grid gap-4 sm:grid-cols-2">
              <ChoiceCard
                title="10 USD clip"
                body="Default until you raise it on purpose."
                Diagram={DiagramPolicy}
                active
                onClick={next}
              />
              <ChoiceCard
                title="Raise later"
                body="Change clip on desktop or CLI. The model cannot raise it."
                onClick={next}
              />
            </div>
          ) : null}

          {current.id === 6 ? (
            <div className="mt-10">
              <PolicyPanel />
            </div>
          ) : null}

          {current.id !== 2 && current.id !== 3 && current.id !== 5 ? (
            <figure className="mt-10 max-w-[36rem] overflow-hidden border border-[rgb(240_231_212/0.28)]">
              <Diagram className="aspect-[16/10] w-full" />
            </figure>
          ) : null}

          {current.id !== 2 && current.id !== 3 && current.id !== 5 ? (
            <Bezel className="mt-8 max-w-[36rem]">
              <p className="text-[0.8125rem] text-[rgb(240_231_212/0.55)]">{current.id} of 12</p>
              <p className="mt-4 text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.8)]">
                <strong className="text-[var(--guide-cream)]">Do this. </strong>
                {current.action}
              </p>
              <p className="mt-4 text-[0.875rem] text-[rgb(240_231_212/0.6)]">
                If it fails: {current.fail.title}. {current.fail.next}
              </p>
            </Bezel>
          ) : null}

          {current.id !== 2 && current.id !== 3 && current.id !== 5 ? (
            <div className="mt-10">
              {step < 12 ? (
                <Button type="button" trailingArrow size="lg" onClick={next}>
                  Continue
                </Button>
              ) : (
                <ButtonLink as={Link} to="/app" trailingArrow size="lg">
                  Open Home
                </ButtonLink>
              )}
            </div>
          ) : null}
        </motion.div>
      </AnimatePresence>
    </div>
  );
}
