import { useEffect, useState } from "react";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { PIPELINE } from "../diagrams/pipeline";
import { cn } from "../lib/cn";

const COPY: Record<(typeof PIPELINE)[number], string> = {
  MARKET: "Live Hyperliquid public books. Empty Watch is real.",
  PRIVATE: "Your thesis stays in the envelope. Not in a chat log.",
  SEALED: "Direct TeeML HPKE. Router keys are refused.",
  RESEARCH: "Researcher reads the book you attached.",
  CHALLENGE: "Challenger attacks the thesis in a separate envelope.",
  RISK: "Risk scores what remains after the challenge.",
  POLICY: "Your clip, assets, and kill. The model cannot raise them.",
  AUTHORIZE: "Desktop or CLI. This browser cannot sign the order.",
  EXECUTE: "Order or cancel only. Hyperliquid must list the PIT agent.",
  PROVE: "0G Storage with the official Go client --proof.",
  LEARN: "Brier and ECE after outcomes. No fake 72 percent.",
};

export function PipelineRing() {
  const [active, setActive] = useState(7);
  const reduce = useReducedMotion();
  const label = PIPELINE[active]!;
  const r = 42;

  useEffect(() => {
    if (reduce) return;
    const t = window.setInterval(() => {
      setActive((i) => (i + 1) % PIPELINE.length);
    }, 2600);
    return () => window.clearInterval(t);
  }, [reduce]);

  return (
    <section className="overflow-hidden border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit">
        <p className="text-[1.25rem] text-[var(--guide-cream)]">The desk</p>
        <h2 className="guide-display mt-4 max-w-[12ch]">
          The
          <br />
          Universal
          <br />
          Pipeline.
        </h2>
        <p className="mt-6 max-w-[36ch] text-[1.2rem] leading-8 text-[rgb(240_231_212/0.75)]">
          Eleven beats. One lit seat is yours.
        </p>

        <div className="relative mx-auto mt-16 aspect-square w-full max-w-[30rem]">
          <div className="absolute inset-[12%] rounded-full border border-[rgb(240_231_212/0.28)]" aria-hidden="true" />
          <div
            className="absolute inset-[28%] rounded-full border border-dashed border-[rgb(240_231_212/0.2)]"
            aria-hidden="true"
          />
          {PIPELINE.map((item, i) => {
            const a = (i / PIPELINE.length) * Math.PI * 2 - Math.PI / 2;
            const x = 50 + r * Math.cos(a);
            const y = 50 + r * Math.sin(a);
            const on = i === active;
            return (
              <button
                key={item}
                type="button"
                onClick={() => setActive(i)}
                className={cn(
                  "absolute grid size-12 -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full border text-[0.625rem] font-bold transition-all duration-300 sm:size-14 sm:text-[0.7rem]",
                  on
                    ? "scale-110 border-[#d82f2f] bg-[#d82f2f] text-black"
                    : "border-[rgb(240_231_212/0.45)] bg-[#1a1a1a] text-[var(--guide-cream)] hover:border-[var(--guide-cream)]",
                )}
                style={{ left: `${x}%`, top: `${y}%` }}
                aria-pressed={on}
                aria-label={`pipeline seat ${i + 1}`}
              >
                {i + 1}
              </button>
            );
          })}
          <div className="absolute inset-0 grid place-items-center px-16 text-center">
            <AnimatePresence mode="wait">
              <motion.div
                key={label}
                initial={reduce ? false : { opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                {...(!reduce ? { exit: { opacity: 0, y: -10 } } : {})}
                transition={{ duration: 0.35 }}
              >
                <p className="text-[0.75rem] tracking-[0.18em] text-[rgb(240_231_212/0.55)] uppercase">{label}</p>
                <p className="mt-3 text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.8)]">{COPY[label]}</p>
              </motion.div>
            </AnimatePresence>
          </div>
        </div>
      </div>
    </section>
  );
}
