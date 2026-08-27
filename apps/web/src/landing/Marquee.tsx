import { PIPELINE } from "../diagrams/pipeline";

export function Marquee() {
  const items = PIPELINE.map((label) => label);
  const loop = [...items, ...items];

  return (
    <section className="overflow-hidden border-y border-[rgb(240_231_212/0.25)] py-16 md:py-20">
      <div className="container-pit mb-10">
        <h2 className="guide-display max-w-[12ch]">
          Already
          <br />
          honest.
        </h2>
        <p className="mt-5 max-w-[40ch] text-[1.125rem] leading-8 text-[rgb(240_231_212/0.72)]">
          The beats are the product. None of them invent a fill, a TEE, or a score.
        </p>
      </div>
      <div className="relative">
        <div className="guide-marquee gap-20 px-6">
          {loop.map((t, i) => (
            <span
              key={`${t}-${i}`}
              className="shrink-0 text-[1.65rem] font-bold tracking-[-0.035em] text-[var(--guide-cream)] whitespace-nowrap md:text-[2.25rem]"
            >
              {t}
            </span>
          ))}
        </div>
      </div>
    </section>
  );
}
