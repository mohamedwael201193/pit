import type { ComponentType } from "react";
import { useEffect, useLayoutEffect, useRef, useState } from "react";
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

function BeatCard({
  Diagram,
  title,
  body,
}: {
  Diagram: ComponentType<{ className?: string }>;
  title: string;
  body: string;
}) {
  return (
    <article className="flex w-[min(84vw,24rem)] shrink-0 flex-col border border-[rgb(240_231_212/0.3)] bg-[#141414] snap-start">
      <div className="overflow-hidden border-b border-[rgb(240_231_212/0.25)]">
        <Diagram className="aspect-[4/3] w-full" />
      </div>
      <div className="flex flex-1 flex-col gap-3 p-6">
        <h3 className="text-[1.75rem] font-bold tracking-[-0.035em] text-[var(--guide-cream)]">{title}</h3>
        <p className="text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.72)]">{body}</p>
      </div>
    </article>
  );
}

function BeatHeader() {
  return (
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
  );
}

export function Moments() {
  const reduce = useReducedMotion();
  const wrap = useRef<HTMLElement>(null);
  const track = useRef<HTMLDivElement>(null);
  const [travel, setTravel] = useState(0);
  const [wide, setWide] = useState(false);
  const { scrollYProgress } = useScroll({ target: wrap, offset: ["start start", "end end"] });
  const rawX = useTransform(scrollYProgress, [0, 1], [0, -travel]);
  const x = useSpring(rawX, { stiffness: 70, damping: 24, restDelta: 0.4 });

  useEffect(() => {
    const mq = window.matchMedia("(min-width: 768px)");
    const sync = () => setWide(mq.matches);
    sync();
    mq.addEventListener("change", sync);
    return () => mq.removeEventListener("change", sync);
  }, []);

  useLayoutEffect(() => {
    const el = track.current;
    if (!el) return;
    const measure = () => {
      const view = el.parentElement?.clientWidth ?? window.innerWidth;
      setTravel(Math.max(0, el.scrollWidth - view));
    };
    measure();
    const ro = new ResizeObserver(measure);
    ro.observe(el);
    window.addEventListener("resize", measure);
    return () => {
      ro.disconnect();
      window.removeEventListener("resize", measure);
    };
  }, [wide, reduce]);

  const cards = CARDS.map((c) => <BeatCard key={c.title} {...c} />);

  if (reduce || !wide) {
    return (
      <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
        <BeatHeader />
        <div className="guide-track mt-12 pl-[max(1.25rem,calc((100vw-72rem)/2+1.5rem))] pr-6">{cards}</div>
      </section>
    );
  }

  return (
    <section
      ref={wrap}
      className="relative border-t border-[rgb(240_231_212/0.25)]"
      style={{ height: `calc(100dvh + ${Math.max(travel, 1)}px)` }}
    >
      <div className="sticky top-0 flex min-h-[100dvh] flex-col justify-center overflow-hidden py-16">
        <BeatHeader />
        <motion.div
          ref={track}
          className="mt-12 flex gap-4 px-[max(1.25rem,calc((100vw-72rem)/2+1.5rem))] will-change-transform"
          style={{ x }}
        >
          {cards}
        </motion.div>
      </div>
    </section>
  );
}
