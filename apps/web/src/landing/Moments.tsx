import type { ComponentType } from "react";
import { useEffect, useRef, useState } from "react";
import { motion, useReducedMotion, useScroll, useSpring, useTransform } from "motion/react";
import {
  DiagramAuthorize,
  DiagramLearn,
  DiagramPolicy,
  DiagramPrivate,
  DiagramSealed,
} from "../diagrams/pitGuide";
import { Reveal } from "../ui/Reveal";

const CARDS: {
  Diagram: ComponentType<{ className?: string }>;
  title: string;
  body: string;
}[] = [
  {
    Diagram: DiagramPrivate,
    title: "Private book",
    body: "Your thesis is sealed to Direct TeeML. The Router never sees the envelope.",
  },
  {
    Diagram: DiagramSealed,
    title: "Three envelopes",
    body: "Researcher, Challenger, and Risk. Same provider is labeled as role separation.",
  },
  {
    Diagram: DiagramPolicy,
    title: "Policy is law",
    body: "Clip, assets, kill. The model cannot raise size, leverage, or permissions.",
  },
  {
    Diagram: DiagramAuthorize,
    title: "You authorize",
    body: "Type AUTHORIZE on the exact preview. This browser cannot hold the session.",
  },
  {
    Diagram: DiagramLearn,
    title: "Then it remembers",
    body: "Verified cases only. Until enough outcomes: NOT ENOUGH DATA.",
  },
];

export function Moments() {
  const reduce = useReducedMotion();
  const wrap = useRef<HTMLElement>(null);
  const track = useRef<HTMLDivElement>(null);
  const [travel, setTravel] = useState(0);
  const { scrollYProgress } = useScroll({ target: wrap, offset: ["start start", "end start"] });
  const rawX = useTransform(scrollYProgress, [0.08, 0.92], [0, -travel]);
  const x = useSpring(rawX, { stiffness: 70, damping: 22, restDelta: 0.4 });

  useEffect(() => {
    const measure = () => {
      if (!track.current) return;
      const extra = Math.max(0, track.current.scrollWidth - window.innerWidth + 48);
      setTravel(extra);
    };
    measure();
    window.addEventListener("resize", measure);
    return () => window.removeEventListener("resize", measure);
  }, []);

  const cards = CARDS.map(({ Diagram, title, body }) => (
    <article
      key={title}
      className="flex w-[min(82vw,26rem)] shrink-0 flex-col border border-[rgb(240_231_212/0.3)] bg-[#141414]"
    >
      <div className="overflow-hidden border-b border-[rgb(240_231_212/0.25)]">
        <Diagram className="aspect-[4/3] w-full" />
      </div>
      <div className="flex flex-1 flex-col gap-3 p-6">
        <h3 className="text-[1.75rem] font-bold tracking-[-0.035em] text-[var(--guide-cream)]">{title}</h3>
        <p className="text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.72)]">{body}</p>
      </div>
    </article>
  ));

  if (reduce) {
    return (
      <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
        <div className="container-pit">
          <h2 className="guide-display max-w-[10ch]">
            Must-see
            <br />
            Beats
          </h2>
          <p className="mt-6 max-w-[36ch] text-[1.25rem] leading-8 text-[rgb(240_231_212/0.78)]">
            Market in. You in the middle. Proof out.
          </p>
        </div>
        <div className="guide-track mt-12 pl-[max(1.25rem,calc((100vw-72rem)/2+1.5rem))] pr-6">{cards}</div>
      </section>
    );
  }

  return (
    <section ref={wrap} className="relative border-t border-[rgb(240_231_212/0.25)]" style={{ height: "240vh" }}>
      <div className="sticky top-0 flex min-h-[100dvh] flex-col justify-center overflow-hidden py-16">
        <div className="container-pit">
          <Reveal>
            <h2 className="guide-display max-w-[10ch]">
              Must-see
              <br />
              Beats
            </h2>
            <p className="mt-6 max-w-[36ch] text-[1.25rem] leading-8 text-[rgb(240_231_212/0.78)]">
              Market in. You in the middle. Proof out.
            </p>
          </Reveal>
        </div>
        <motion.div ref={track} className="mt-12 flex gap-4 px-[max(1.25rem,calc((100vw-72rem)/2+1.5rem))]" style={{ x }}>
          {cards}
        </motion.div>
      </div>
    </section>
  );
}
