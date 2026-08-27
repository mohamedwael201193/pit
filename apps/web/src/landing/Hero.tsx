import { useRef } from "react";
import { usePrivy } from "@privy-io/react-auth";
import { motion, useReducedMotion, useScroll, useSpring, useTransform } from "motion/react";
import { WireTurn } from "../diagrams/WireTurn";
import { DiagramHeroPostcard } from "../diagrams/pitGuide";

export function Hero() {
  const { ready, authenticated, login } = usePrivy();
  const reduce = useReducedMotion();
  const pinRef = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({ target: pinRef, offset: ["start start", "end start"] });
  const rawRotate = useTransform(scrollYProgress, [0, 1], [0, 220]);
  const rotate = useSpring(rawRotate, { stiffness: 45, damping: 18 });
  const postcardY = useTransform(scrollYProgress, [0, 1], [0, 56]);
  const megaScale = useTransform(scrollYProgress, [0, 0.55], [1, 0.92]);

  return (
    <div ref={pinRef} className="guide-pin">
      <section className="guide-pin__sticky guide-coral isolate">
        <div className="guide-grain" aria-hidden="true" />
        <div className="container-pit relative flex min-h-[100svh] flex-col pt-20 pb-8 md:pt-24">
          <div className="grid flex-1 items-start gap-6 pt-4 lg:grid-cols-[1.25fr_0.75fr] lg:gap-10">
            <div>
              <motion.h1
                className="guide-kicker max-w-6xl"
                initial={reduce ? false : { opacity: 0, y: 28 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
              >
                PIT
              </motion.h1>
              <motion.p
                className="mt-8 max-w-[22ch] text-[1.35rem] leading-8 font-medium text-black sm:text-[1.55rem] sm:leading-9"
                initial={reduce ? false : { opacity: 0, y: 18 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.75, delay: 0.1, ease: [0.16, 1, 0.3, 1] }}
              >
                Private research. Controlled execution. A desk that learns.
              </motion.p>
              <motion.p
                className="mt-4 max-w-[36ch] text-[1.0625rem] leading-7 text-black/75"
                initial={reduce ? false : { opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.7, delay: 0.16, ease: [0.16, 1, 0.3, 1] }}
              >
                PIT never asks for a seed phrase. Your wallet stays yours. Your trading session cannot withdraw.
              </motion.p>
            </div>
            <div className="flex justify-end self-start lg:pt-2">
              <WireTurn className="size-36 sm:size-44 lg:size-[11.5rem]" rotate={reduce ? 0 : rotate} />
            </div>
          </div>

          <hr className="guide-rule mt-6 border-black" />

          <div className="grid items-end gap-8 py-8 md:grid-cols-[1fr_auto] md:gap-14 md:py-10">
            <motion.h2
              className="guide-mega origin-left"
              {...(!reduce ? { style: { scale: megaScale } } : {})}
              initial={reduce ? false : { opacity: 0, y: 36 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.9, delay: 0.16, ease: [0.16, 1, 0.3, 1] }}
            >
              PIT.
            </motion.h2>
            <motion.figure
              className="w-full max-w-[16rem] justify-self-end border border-black bg-black/5 md:w-52"
              {...(!reduce ? { style: { y: postcardY } } : {})}
            >
              <DiagramHeroPostcard className="aspect-[4/3] w-full" />
              <figcaption className="border-t border-black bg-[#f0e7d4] px-3 py-2 text-[0.75rem] leading-4 text-black/75">
                <span className="mr-1.5 rounded-sm border border-black/25 px-1 py-px text-[0.625rem] font-semibold tracking-wide uppercase">
                  Illustration
                </span>
                Example desk. One lit seat is yours to sign.
              </figcaption>
            </motion.figure>
          </div>

          <hr className="guide-rule border-black" />

          <div className="flex flex-wrap items-center gap-3 pt-6 pb-2">
            {!ready ? (
              <p className="text-black/80">Loading wallet connect</p>
            ) : authenticated ? (
              <a href="/app" className="pill pill-ink">
                Open your desk
              </a>
            ) : (
              <button className="pill pill-ink" type="button" onClick={login}>
                Connect your wallet
              </button>
            )}
            <a href="https://github.com/mohamedwael201193/pit/releases/latest" className="pill pill-ghost">
              Download PIT
            </a>
            <a href="#story" className="pill pill-ghost">
              Read the story
            </a>
          </div>
        </div>
      </section>
    </div>
  );
}
