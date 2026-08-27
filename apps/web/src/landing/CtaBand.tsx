import { Link } from "react-router-dom";
import { Reveal } from "../ui/Reveal";

export function CtaBand() {
  return (
    <section className="relative isolate overflow-hidden border-t border-[rgb(240_231_212/0.25)]">
      <div className="absolute inset-0 bg-[#1a1a1a]" aria-hidden="true" />
      <div className="guide-grain opacity-40" aria-hidden="true" />
      <div className="container-pit relative py-28 md:py-36">
        <Reveal>
          <h2 className="guide-display mt-5 max-w-[14ch]">
            Built for the desk,
            <br />
            signed on the machine.
          </h2>
          <p className="mt-7 max-w-[40ch] text-[1.3rem] leading-8 text-[rgb(240_231_212/0.82)]">
            Connect a wallet you already own. Inspect here. Authorize on desktop. Prove it later.
          </p>
          <div className="mt-12 flex flex-wrap gap-3">
            <Link
              to="/signin"
              className="inline-flex h-12 items-center rounded-full bg-[#d82f2f] px-6 text-base font-medium text-black no-underline"
            >
              Get started
            </Link>
            <a
              href="#story"
              className="inline-flex h-12 items-center rounded-full border border-[var(--guide-cream)] px-6 text-base font-medium text-[var(--guide-cream)] no-underline"
            >
              How it works
            </a>
          </div>
        </Reveal>
      </div>
    </section>
  );
}
