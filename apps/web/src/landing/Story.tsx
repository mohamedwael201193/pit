import { useRef } from "react";
import { motion, useReducedMotion, useScroll, useTransform } from "motion/react";

const BEATS = [
  {
    title: "Sealed,",
    body: "Your private book is sealed to 0G Direct TeeML. The Router never sees it.",
  },
  {
    title: "Challenged,",
    body: "Researcher, Challenger, and Risk run as three envelopes. Same provider is labeled honestly.",
  },
  {
    title: "Authorized.",
    body: "You type AUTHORIZE on the exact preview. A one-hour agent can order or cancel. It cannot withdraw.",
  },
] as const;

export function Story() {
  const reduce = useReducedMotion();
  const pinRef = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({ target: pinRef, offset: ["start end", "end start"] });
  const y1 = useTransform(scrollYProgress, [0.15, 0.55], [40, 0]);
  const y2 = useTransform(scrollYProgress, [0.25, 0.65], [56, 0]);
  const y3 = useTransform(scrollYProgress, [0.35, 0.75], [72, 0]);
  const opacities = [
    useTransform(scrollYProgress, [0.1, 0.35], [0.15, 1]),
    useTransform(scrollYProgress, [0.2, 0.45], [0.15, 1]),
    useTransform(scrollYProgress, [0.3, 0.55], [0.15, 1]),
  ];
  const ys = [y1, y2, y3];

  return (
    <section id="story" className="relative border-t border-[rgb(240_231_212/0.25)]">
      <div className="container-pit py-20 md:py-28">
        <p className="max-w-[52ch] text-[1.35rem] leading-9 text-[rgb(240_231_212/0.85)]">
          You already have a terminal. PIT is the desk that hunts without spending, seals the book, and waits for you.
        </p>
      </div>
      <div ref={pinRef} className="relative min-h-[140vh]">
        <div className="sticky top-0 flex min-h-[100svh] items-center overflow-hidden bg-[#1a1a1a]">
          <div className="container-pit w-full py-16">
            <div className="flex flex-col gap-2 md:gap-3">
              {BEATS.map((b, i) => (
                <motion.h2
                  key={b.title}
                  className="guide-display"
                  style={reduce ? undefined : { y: ys[i], opacity: opacities[i] }}
                >
                  {b.title}
                </motion.h2>
              ))}
            </div>
            <div className="mt-16 grid gap-8 border-t border-[rgb(240_231_212/0.25)] pt-10 md:grid-cols-3">
              {BEATS.map((b) => (
                <p key={b.title} className="max-w-[36ch] text-[1.05rem] leading-7 text-[rgb(240_231_212/0.75)]">
                  {b.body}
                </p>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
