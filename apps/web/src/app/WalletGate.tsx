import { useState } from "react";
import { Link } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { PipelineCard } from "../diagrams/PipelineCard";
import { classifyError } from "../namedStates";

export function WalletGate() {
  const { ready, authenticated, login, logout, user } = usePrivy();
  const reduce = useReducedMotion();
  const [error, setError] = useState<string | null>(null);
  const addr = user?.wallet?.address;

  const run = async () => {
    setError(null);
    try {
      await login();
    } catch (err) {
      const named = classifyError(err instanceof Error ? err.message : "connect failed");
      setError(`${named.title}. ${named.body} ${named.next}`);
    }
  };

  return (
    <div className="guide-shell guide-app grid min-h-[100dvh] lg:grid-cols-[minmax(0,1fr)_minmax(0,0.9fr)]">
      <div className="flex flex-col justify-center px-5 py-16 sm:px-10 lg:px-16">
        <div className="w-full max-w-[28rem]">
          <Link to="/" className="text-[1.5rem] font-bold tracking-[-0.06em] text-[#d82f2f] no-underline" aria-label="PIT home">
            PIT.
          </Link>
          <p className="mt-10 text-[0.8125rem] font-bold tracking-[0.18em] text-[#d82f2f] uppercase">YOUR DESK</p>
          <h1 className="guide-display mt-3 !text-[clamp(2.25rem,5vw,3.5rem)]">
            {authenticated ? "Wallet connected" : "Sign in with your wallet"}
          </h1>
          <AnimatePresence mode="wait" initial={false}>
            <motion.div
              key={authenticated ? "in" : "out"}
              initial={reduce ? false : { opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={reduce ? { opacity: 1 } : { opacity: 0, y: -8 }}
              transition={{ duration: 0.32, ease: [0.16, 1, 0.3, 1] }}
            >
              <p className="mt-5 max-w-[46ch] text-[1.125rem] leading-8 text-[rgb(240_231_212/0.75)]">
                PIT never asks for a seed phrase. Connect the wallet you already use. Session keys stay on desktop.
              </p>
              {!ready ? (
                <p className="mt-8">Loading wallet connect</p>
              ) : authenticated ? (
                <div className="mt-8 rounded-none border border-[rgb(240_231_212/0.28)] p-5">
                  <p className="text-[11px] tracking-[0.16em]">YOUR WALLET</p>
                  <p className="font-mono break-all">{addr}</p>
                  <div className="mt-4 flex flex-wrap gap-3">
                    <Link to="/app/start" className="pill pill-coral">
                      Continue
                    </Link>
                    <button className="pill pill-line" type="button" onClick={logout}>
                      Disconnect
                    </button>
                  </div>
                </div>
              ) : (
                <div className="mt-8 flex flex-col gap-3">
                  <button className="pill pill-coral h-12 w-full" type="button" onClick={() => void run()}>
                    Connect your wallet
                  </button>
                  <Link to="/" className="pill pill-line h-12 w-full">
                    Browse first
                  </Link>
                </div>
              )}
              {error ? (
                <p role="alert" className="mt-6 text-[0.875rem] leading-6 text-[#ff7a7a]">
                  {error}
                </p>
              ) : null}
            </motion.div>
          </AnimatePresence>
        </div>
      </div>
      <aside className="relative hidden items-center justify-center border-l border-[rgb(240_231_212/0.25)] bg-[#d82f2f] lg:flex">
        <div className="guide-grain opacity-30" aria-hidden="true" />
        <div className="relative w-full max-w-[26rem] px-10">
          <PipelineCard className="w-full shadow-none" />
          <p className="mt-8 text-center text-[1rem] leading-7 font-medium text-black">
            Your EVM address is the login. Orders stay on desktop. Transfer of Agentic ID is not live on mainnet.
          </p>
        </div>
      </aside>
    </div>
  );
}
