import { motion, useReducedMotion } from "motion/react";

const FLOW =
  "private intelligence → 0G Direct → Researcher → Challenger → Risk → VerifyE2EE → policy → preview → AUTHORIZE → Hyperliquid order → OID / fill → 0G proof";

export function WatchFilm() {
  const reduce = useReducedMotion();

  return (
    <section
      id="watch"
      className="watch-film relative overflow-x-clip border-t border-black/10 scroll-mt-24"
      aria-labelledby="watch-pit-heading"
    >
      <div className="watch-film-wash" aria-hidden="true" />
      <div className="container-pit relative z-[1] py-10 md:py-12">
        <motion.div
          initial={reduce ? false : { opacity: 0, y: 14 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, amount: 0.35 }}
          transition={{ duration: 0.7, ease: [0.16, 1, 0.3, 1] }}
          className="flex flex-wrap items-end justify-between gap-x-8 gap-y-3"
        >
          <div>
            <p className="text-[0.72rem] font-semibold tracking-[0.22em] text-[#d82f2f]">REAL HYPERLIQUID</p>
            <h2
              id="watch-pit-heading"
              className="mt-2 max-w-[18ch] text-[clamp(1.65rem,3.4vw,2.35rem)] font-semibold tracking-[-0.04em] leading-[1.12] text-black"
            >
              Watch PIT in action
            </h2>
            <p className="mt-2 max-w-[40ch] text-[1.02rem] leading-6 text-black/72">
              Recorded fill OID 531667200134 — buy 0.16 HYPE. Matching 0G research and order proofs on Aristotle.
            </p>
          </div>
          <p className="max-w-[34ch] pb-1 text-[0.9rem] leading-6 text-black/62">{FLOW}</p>
        </motion.div>

        <motion.div
          className="watch-film-stage mt-6 md:mt-8"
          initial={reduce ? false : { opacity: 0, y: 18 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, amount: 0.2 }}
          transition={{ duration: 0.85, delay: 0.06, ease: [0.16, 1, 0.3, 1] }}
        >
          <div className="watch-film-glow" aria-hidden="true" />
          <div className="watch-film-frame">
            <iframe
              className="absolute inset-0 size-full border-0"
              src="https://www.youtube-nocookie.com/embed/zYgxDTI7jIk?rel=0"
              title="Watch PIT in action"
              allow="accelerometer; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
              referrerPolicy="strict-origin-when-cross-origin"
              allowFullScreen
            />
          </div>
        </motion.div>

        <p className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 text-[0.82rem] tracking-wide text-black/55">
          <a className="font-medium text-[#d82f2f] underline-offset-4 hover:underline" href="https://youtu.be/zYgxDTI7jIk">
            Open the film
          </a>
          <span>0G Direct · VerifyE2EE · Storage --proof · Hyperliquid FILLED</span>
        </p>
      </div>
    </section>
  );
}
