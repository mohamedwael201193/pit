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
  return (
    <section id="story" className="relative border-t border-[rgb(240_231_212/0.25)]">
      <div className="container-pit py-16 md:py-20">
        <Reveal>
          <p className="max-w-[22ch] text-[1.75rem] leading-9 font-medium tracking-[-0.03em] text-[var(--guide-cream)] md:text-[2.15rem] md:leading-10">
            The web discovers. The desktop acts.
          </p>
          <p className="mt-5 max-w-[46ch] text-[1.125rem] leading-7 text-[rgb(240_231_212/0.78)]">
            This computer must stay awake for the bound. If it sleeps, the mission stops. That gap is not backfilled.
          </p>
        </Reveal>

        <ul className="mt-12 grid gap-8 md:gap-10">
          {BEATS.map((b) => (
            <li
              key={b.title}
              className="grid items-end gap-3 border-t border-[rgb(240_231_212/0.2)] pt-6 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)] lg:gap-10"
            >
              <h2 className="story-display">{b.title}</h2>
              <p className="max-w-[40ch] text-[1.125rem] leading-7 text-[rgb(240_231_212/0.8)]">{b.body}</p>
            </li>
          ))}
        </ul>

        <figure className="mt-12 max-w-lg overflow-hidden border border-[rgb(240_231_212/0.35)]">
          <DiagramPrivate className="aspect-[16/10] w-full" />
        </figure>
      </div>
    </section>
  );
}
