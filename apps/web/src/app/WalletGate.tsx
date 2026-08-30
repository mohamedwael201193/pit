import { useState } from "react";
import { Link, Navigate, useLocation } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { PitMark } from "../brand/PitMark";
import { DiagramHeroPostcard } from "../diagrams/pitGuide";
import { classifyError } from "../namedStates";
import { Button, ButtonLink } from "../ui/Button";

export function WalletGate() {
  const { ready, authenticated, login, logout, user } = usePrivy();
  const location = useLocation();
  const reduce = useReducedMotion();
  const [error, setError] = useState<string | null>(null);
  const addr = user?.wallet?.address;

  if (ready && authenticated && location.pathname === "/signin") {
    return <Navigate to="/protect" replace />;
  }

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
          <Link to="/" className="inline-flex" aria-label="PIT home">
            <PitMark />
          </Link>
          <h1 className="guide-display mt-10 !text-[clamp(2.25rem,5vw,3.5rem)]">
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
                Connect the wallet you already use. No seed field exists. Session keys stay on desktop.
              </p>
              {!ready ? (
                <p className="mt-8">Loading wallet connect</p>
              ) : authenticated ? (
                <div className="mt-8 border border-[rgb(240_231_212/0.28)] p-5">
                  <p className="text-[0.75rem] tracking-[0.16em] text-[rgb(240_231_212/0.55)]">YOUR WALLET</p>
                  <p className="mt-2 font-mono break-all text-[0.9375rem]">{addr}</p>
                  <div className="mt-5 flex flex-wrap gap-3">
                    <ButtonLink as={Link} to="/protect" trailingArrow size="lg">
                      Protect my strategy
                    </ButtonLink>
                    <ButtonLink as={Link} to="/app" variant="secondary" size="lg">
                      Overview
                    </ButtonLink>
                    <Button variant="secondary" type="button" onClick={logout}>
                      Disconnect
                    </Button>
                  </div>
                </div>
              ) : (
                <div className="mt-8 flex flex-col gap-3">
                  <button
                    className="inline-flex h-12 w-full items-center justify-center rounded-full bg-[#d82f2f] px-6 text-base font-medium text-[#f0e7d4]"
                    type="button"
                    onClick={() => void run()}
                  >
                    Connect your wallet
                  </button>
                  <ButtonLink as={Link} to="/radar" variant="ghost" size="lg" className="w-full">
                    Browse radar first
                  </ButtonLink>
                </div>
              )}
              {error ? (
                <p role="alert" className="mt-6 border-t border-[rgb(240_231_212/0.25)] pt-6 text-[0.875rem] leading-6 text-[#ff7a7a]">
                  Connection issue: {error} Try again in a moment.
                </p>
              ) : null}
            </motion.div>
          </AnimatePresence>
        </div>
      </div>
      <aside className="relative hidden items-center justify-center border-l border-[rgb(240_231_212/0.25)] bg-[#d82f2f] lg:flex">
        <div className="guide-grain opacity-30" aria-hidden="true" />
        <div className="relative w-full max-w-[26rem] px-10">
          <DiagramHeroPostcard className="w-full border border-black" />
          <p className="mt-8 text-center text-[1rem] leading-7 font-medium text-black">
            Your EVM address is the login. Orders stay on desktop. Transfer of Agentic ID is not live on mainnet.
          </p>
        </div>
      </aside>
    </div>
  );
}
