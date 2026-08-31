import { motion, useReducedMotion } from "motion/react";
import { WireIris, WireTurn } from "../diagrams/WireTurn";
import { HeroCtas } from "./HeroCtas";

export function Hero() {
  const reduce = useReducedMotion();

  return (
    <section className="guide-coral isolate relative min-h-[100dvh] overflow-x-clip">
      <div className="guide-grain" aria-hidden="true" />
      <div className="relative z-[1] grid min-h-[100dvh] grid-rows-[1fr_auto_1fr]">
        <div aria-hidden="true" />
        <div className="container-pit w-full py-6">
        <div className="hero-board grid w-full items-center gap-8 lg:grid-cols-[minmax(0,1.15fr)_minmax(12rem,22rem)] lg:gap-12">
          <motion.div
            initial={reduce ? false : { opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.75, ease: [0.16, 1, 0.3, 1] }}
          >
            <p className="guide-kicker">PRIVATE AI TRADING DESK</p>
            <h1 className="mt-5 text-[clamp(1.85rem,4.2vw,3.15rem)] leading-[1.14] font-semibold tracking-[-0.035em] text-black md:mt-6">
              0G seals the book.
              <br />
              Hyperliquid fills the order.
            </h1>
            <p className="mt-4 max-w-[48ch] text-[1.25rem] leading-8 text-black/85 md:text-[1.35rem] md:leading-8">
              Private intelligence → Direct TeeML → Researcher, Challenger, Risk → VerifyE2EE → policy preview → you authorize → a real Hyperliquid order, then 0G proof.
            </p>
            <HeroCtas />
          </motion.div>
          <motion.div
            className="hero-wire relative mx-auto grid size-[min(42vw,18rem)] place-items-center lg:size-[min(46vh,22rem)] lg:justify-self-end"
            initial={reduce ? false : { opacity: 0, scale: 0.94 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.9, delay: 0.08, ease: [0.16, 1, 0.3, 1] }}
          >
            <div className={`size-full ${reduce ? "" : "guide-spin"}`}>
              <WireTurn className="size-full" rotate={0} />
            </div>
            <WireIris className={`pointer-events-none absolute inset-[22%] ${reduce ? "" : "guide-spin-rev"}`} />
          </motion.div>
        </div>
        </div>
        <div aria-hidden="true" />
      </div>
    </section>
  );
}
