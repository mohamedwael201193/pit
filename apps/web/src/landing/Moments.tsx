import type { ComponentType } from "react";
import { motion, useReducedMotion } from "motion/react";
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
    body: "Verified cases only. Until enough outcomes: NOT ENOUGH DATA. The model did not learn this.",
  },
];

export function Moments() {
  const reduce = useReducedMotion();
  return (
    <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit">
        <Reveal>
          <h2 className="guide-display max-w-[10ch]">
            Must-see
            <br />
            Beats
          </h2>
          <p className="mt-6 max-w-[46ch] text-[1.2rem] leading-8 text-[rgb(240_231_212/0.78)]">
            Market in. You in the middle. Proof out. Scroll sideways through the desk that keeps the promise.
          </p>
        </Reveal>
      </div>
      <div className="guide-track mt-12 pl-[max(1.25rem,calc((100vw-72rem)/2+1.5rem))] pr-6">
        {CARDS.map(({ Diagram, title, body }) => (
          <motion.article
            key={title}
            className="flex flex-col border border-[rgb(240_231_212/0.3)] bg-[#141414]"
            whileHover={reduce ? undefined : { y: -6 }}
            transition={{ type: "spring", stiffness: 260, damping: 22 }}
          >
            <div className="overflow-hidden border-b border-[rgb(240_231_212/0.25)]">
              <Diagram className="aspect-[4/3] w-full" />
            </div>
            <div className="flex flex-1 flex-col gap-3 p-5">
              <h3 className="text-[1.65rem] font-bold tracking-[-0.035em] text-[var(--guide-cream)]">{title}</h3>
              <p className="text-[1rem] leading-7 text-[rgb(240_231_212/0.7)]">{body}</p>
            </div>
          </motion.article>
        ))}
      </div>
    </section>
  );
}
