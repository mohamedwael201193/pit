import { useRef } from "react";
import { motion, useReducedMotion, useScroll, useSpring, useTransform } from "motion/react";
import { WireIris, WireTurn } from "../diagrams/WireTurn";
import { HeroCtas } from "./HeroCtas";

export function Hero() {
  const reduce = useReducedMotion();
  const pinRef = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({ target: pinRef, offset: ["start start", "end start"] });
  const rawRotate = useTransform(scrollYProgress, [0, 1], [0, 160]);
  const rotate = useSpring(rawRotate, { stiffness: 40, damping: 18 });
  const wireScale = useTransform(scrollYProgress, [0, 0.8], [1, 1.08]);
  const titleY = useTransform(scrollYProgress, [0, 0.7], [0, -18]);

  return (
    <div ref={pinRef} className="guide-pin">
      <section className="guide-pin__sticky guide-coral isolate">
        <div className="guide-grain" aria-hidden="true" />
        <div className="container-pit relative flex min-h-[100dvh] flex-col justify-center pt-20 pb-12 md:pt-24">
          <div className="grid items-center gap-10 lg:grid-cols-[1.1fr_0.9fr] lg:gap-6">
            <motion.div style={reduce ? undefined : { y: titleY }}>
              <motion.p
                className="guide-kicker"
                initial={reduce ? false : { opacity: 0, y: 24 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
              >
                PRIVATE ALPHA OS
              </motion.p>
              <motion.h1
                className="mt-8 max-w-5xl text-[clamp(1.85rem,4.4vw,3.35rem)] leading-[1.12] font-semibold tracking-[-0.035em] text-black"
                initial={reduce ? false : { opacity: 0, y: 16 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.75, delay: 0.08, ease: [0.16, 1, 0.3, 1] }}
              >
                It hunts while you sleep.
                <br />
                Your keys never leave your machine.
              </motion.h1>
              <motion.p
                className="mt-6 max-w-[42ch] text-[1.15rem] leading-7 text-black/75"
                initial={reduce ? false : { opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.7, delay: 0.12, ease: [0.16, 1, 0.3, 1] }}
              >
                PIT watches live markets, researches your private strategy through 0G, challenges every idea, enforces your policy on the host, and when you explicitly arm a bounded Sleep Mission, can act within those limits while you are away.
              </motion.p>
              <motion.p
                className="mt-4 max-w-[42ch] text-[1rem] leading-7 text-black/70"
                initial={reduce ? false : { opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.65, delay: 0.18, ease: [0.16, 1, 0.3, 1] }}
              >
                This computer must stay awake for the bound. If it sleeps, the mission stops. That gap is not backfilled.
              </motion.p>
              <HeroCtas />
            </motion.div>
            <motion.div
              className="relative mx-auto flex size-[min(78vw,22rem)] items-center justify-center lg:size-[28rem] lg:justify-self-end"
              style={reduce ? undefined : { scale: wireScale }}
              initial={reduce ? false : { opacity: 0, scale: 0.92 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 1, delay: 0.12, ease: [0.16, 1, 0.3, 1] }}
            >
              <WireTurn className="size-full" rotate={reduce ? 0 : rotate} />
              <WireIris className={`pointer-events-none absolute inset-[22%] ${reduce ? "" : "guide-spin-rev"}`} />
            </motion.div>
          </div>
        </div>
      </section>
    </div>
  );
}
