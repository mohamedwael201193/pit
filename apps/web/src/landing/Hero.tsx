import { useRef } from "react";
import { Link } from "react-router-dom";
import { motion, useReducedMotion, useScroll, useSpring, useTransform } from "motion/react";
import { WireTurn } from "../diagrams/WireTurn";

export function Hero() {
  const reduce = useReducedMotion();
  const pinRef = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({ target: pinRef, offset: ["start start", "end start"] });
  const rawRotate = useTransform(scrollYProgress, [0, 1], [0, 240]);
  const rotate = useSpring(rawRotate, { stiffness: 40, damping: 18 });
  const wireScale = useTransform(scrollYProgress, [0, 0.7], [1, 1.12]);
  const titleY = useTransform(scrollYProgress, [0, 0.55], [0, -28]);

  return (
    <div ref={pinRef} className="guide-pin">
      <section className="guide-pin__sticky guide-coral isolate">
        <div className="guide-grain" aria-hidden="true" />
        <div className="container-pit relative flex min-h-[100dvh] flex-col justify-center pt-20 pb-12 md:pt-24">
          <div className="grid items-center gap-8 lg:grid-cols-[1.15fr_0.85fr] lg:gap-8">
            <motion.div style={reduce ? undefined : { y: titleY }}>
              <motion.p
                className="guide-kicker"
                initial={reduce ? false : { opacity: 0, y: 24 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
              >
                PIT
              </motion.p>
              <motion.h1
                className="mt-8 max-w-[18ch] text-[1.7rem] leading-[1.2] font-semibold tracking-[-0.03em] text-black sm:text-[2.05rem] sm:leading-[1.18] lg:text-[2.25rem]"
                initial={reduce ? false : { opacity: 0, y: 16 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.75, delay: 0.08, ease: [0.16, 1, 0.3, 1] }}
              >
                It hunts in private. You authorize on this computer.
              </motion.h1>
              <motion.div
                className="mt-10 flex flex-wrap items-center gap-3"
                initial={reduce ? false : { opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.7, delay: 0.16, ease: [0.16, 1, 0.3, 1] }}
              >
                <Link to="/radar" className="pill pill-ink">
                  Explore live PIT
                </Link>
                <Link to="/download" className="pill pill-ghost">
                  Download PIT Desktop
                </Link>
              </motion.div>
            </motion.div>
            <motion.div
              className="flex justify-end"
              style={reduce ? undefined : { scale: wireScale }}
              initial={reduce ? false : { opacity: 0, scale: 0.92 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 1, delay: 0.12, ease: [0.16, 1, 0.3, 1] }}
            >
              <WireTurn className="size-48 sm:size-64 lg:size-[22rem]" rotate={reduce ? 0 : rotate} />
            </motion.div>
          </div>
        </div>
      </section>
    </div>
  );
}
