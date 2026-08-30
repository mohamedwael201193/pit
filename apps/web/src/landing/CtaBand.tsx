import { Link } from "react-router-dom";
import { Reveal } from "../ui/Reveal";

export function CtaBand() {
  return (
    <section className="relative isolate overflow-hidden border-t border-[rgb(240_231_212/0.25)]">
      <div className="absolute inset-0 bg-[#1a1a1a]" aria-hidden="true" />
      <div className="guide-grain opacity-40" aria-hidden="true" />
      <div className="container-pit relative py-28 md:py-36">
        <Reveal>
          <h2 className="guide-display mt-5 max-w-[16ch]">Let PIT watch. Keep execution on your machine.</h2>
          <p className="mt-7 max-w-[40ch] text-[1.3rem] leading-8 text-[rgb(240_231_212/0.82)]">
            Explore live books here. Verify proof here. Authorize only on desktop.
          </p>
          <div className="mt-12 flex flex-wrap gap-3">
            <Link to="/download" className="pill pill-coral">
              Download PIT Desktop
            </Link>
            <Link to="/radar" className="pill pill-line">
              Explore live PIT
            </Link>
            <Link to="/proof" className="pill pill-line">
              Verify a mission
            </Link>
          </div>
        </Reveal>
      </div>
    </section>
  );
}
