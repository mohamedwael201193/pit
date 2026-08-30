import { useRef } from "react";
import { motion, useReducedMotion, useScroll, useTransform } from "motion/react";
import { DiagramPrivate } from "../diagrams/pitGuide";
import { Reveal } from "../ui/Reveal";

const BEATS = [
  {
    title: "Sealed,",
    body: "Your book goes into Direct TeeML. The Router never sees it.",
  },
  {
    title: "Challenged,",
    body: "Researcher, Challenger, and Risk run as three envelopes. Same provider is labeled honestly.",
  },
  {
    title: "Authorized.",
    body: "You type AUTHORIZE on the exact preview. Chat and this website cannot.",
  },
] as const;

export function Story() {
  const reduce = useReducedMotion();
  const pinRef = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({ target: pinRef, offset: ["start end", "end start"] });
  const y1 = useTransform(scrollYProgress, [0.2, 0.55], [24, 0]);
  const y2 = useTransform(scrollYProgress, [0.28, 0.62], [32, 0]);
  const y3 = useTransform(scrollYProgress, [0.36, 0.7], [40, 0]);
  const ys = [y1, y2, y3];

  return (
    <section id="story" className="relative border-t border-[rgb(240_231_212/0.25)]">
      <div className="container-pit py-16 md:py-20">
        <Reveal>
          <p className="max-w-[28ch] text-[1.65rem] leading-9 font-medium tracking-[-0.03em] text-[var(--guide-cream)] md:text-[2rem] md:leading-10">
            The web discovers. The desktop acts.
          </p>
        </Reveal>
      </div>
      <div ref={pinRef} className="relative">
        <div className="flex min-h-[100dvh] items-center bg-[#1a1a1a]">
          <div className="container-pit w-full py-16">
            <div className="flex flex-col gap-2 md:gap-3">
              {BEATS.map((b, i) => (
                <motion.h2
                  key={b.title}
                  className="guide-display"
                  style={reduce ? undefined : { y: ys[i] }}
                >
                  {b.title}
                </motion.h2>
              ))}
            </div>
            <div className="mt-14 grid gap-8 border-t border-[rgb(240_231_212/0.25)] pt-10 md:grid-cols-3">
              {BEATS.map((b) => (
                <p key={b.title} className="max-w-[32ch] text-[1.125rem] leading-7 text-[rgb(240_231_212/0.78)]">
                  {b.body}
                </p>
              ))}
            </div>
            <figure className="mt-12 max-w-xl overflow-hidden border border-[rgb(240_231_212/0.35)]">
              <DiagramPrivate className="aspect-[16/10] w-full" />
            </figure>
          </div>
        </div>
      </div>
    </section>
  );
}
