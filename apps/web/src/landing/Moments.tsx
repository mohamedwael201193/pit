import type { ComponentType } from "react";
import { Reveal } from "../ui/Reveal";
import {
  DiagramAuthorize,
  DiagramLearn,
  DiagramPolicy,
  DiagramPrivate,
  DiagramSealed,
} from "../diagrams/pitGuide";

const CARDS: {
  n: string;
  Diagram: ComponentType<{ className?: string }>;
  title: string;
  body: string;
}[] = [
  {
    n: "01",
    Diagram: DiagramPrivate,
    title: "Private book",
    body: "Your thesis is sealed to Direct TeeML. The Router never sees the envelope.",
  },
  {
    n: "02",
    Diagram: DiagramSealed,
    title: "Three envelopes",
    body: "Researcher, Challenger, and Risk. Same provider is labeled as role separation.",
  },
  {
    n: "03",
    Diagram: DiagramPolicy,
    title: "Policy is law",
    body: "Clip, assets, kill. The model cannot raise size, leverage, or permissions.",
  },
  {
    n: "04",
    Diagram: DiagramAuthorize,
    title: "You authorize",
    body: "Type AUTHORIZE on the exact preview. This browser cannot hold the session.",
  },
  {
    n: "05",
    Diagram: DiagramLearn,
    title: "Then it remembers",
    body: "Verified cases only. Until enough outcomes: NOT ENOUGH DATA.",
  },
];

export function Moments() {
  return (
    <section className="border-t border-[rgb(240_231_212/0.25)] py-14 md:py-16">
      <div className="container-pit">
        <Reveal>
          <h2 className="guide-display max-w-[10ch]">
            Must-see
            <br />
            Beats
          </h2>
          <p className="mt-5 max-w-[36ch] text-[1.2rem] leading-7 text-[rgb(240_231_212/0.78)]">
            Market in. You in the middle. Proof out.
          </p>
        </Reveal>
        <div className="mt-10 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-6">
          {CARDS.map((c, i) => (
            <article
              key={c.title}
              className={`flex min-h-0 flex-col overflow-hidden border border-[rgb(240_231_212/0.3)] bg-[#141414] ${
                i < 2 ? "lg:col-span-3" : "lg:col-span-2"
              }`}
            >
              <div className="relative h-[13.25rem] overflow-hidden bg-[#f0e7d4] sm:h-[14.5rem]">
                <span className="absolute left-3 top-3 z-10 grid size-8 place-items-center rounded-full bg-black text-[0.7rem] font-bold tracking-wide text-[#f0e7d4]">
                  {c.n}
                </span>
                <c.Diagram className="h-full w-full" />
              </div>
              <div className="flex flex-1 flex-col gap-2 p-4">
                <h3 className="text-[1.2rem] font-bold tracking-[-0.03em] text-[var(--guide-cream)] md:text-[1.3rem]">
                  {c.title}
                </h3>
                <p className="text-[0.98rem] leading-6 text-[rgb(240_231_212/0.78)]">{c.body}</p>
              </div>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
