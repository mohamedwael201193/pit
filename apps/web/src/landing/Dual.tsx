import { Reveal } from "../ui/Reveal";
import { DiagramWideBanner } from "../diagrams/pitGuide";

export function Dual() {
  return (
    <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit">
        <Reveal>
          <p className="text-[1.125rem] text-[rgb(240_231_212/0.7)]">The product is</p>
          <h2 className="guide-display mt-3 max-w-[16ch]">MAINNET only</h2>
          <p className="mt-6 max-w-[46ch] text-[1.2rem] leading-8 text-[rgb(240_231_212/0.78)]">
            Aristotle 16661 and Hyperliquid mainnet. Direct TeeML for the private book. Transfer of Agentic ID is not live.
            The laboratory exists for CI and developers, not for the public desk.
          </p>
        </Reveal>
        <div className="mt-12 min-h-[16rem] bg-[#d82f2f] p-8 text-black md:p-12">
          <p className="text-[0.8125rem] tracking-[0.16em] uppercase opacity-70">Production</p>
          <p className="mt-6 text-[1.75rem] leading-9 font-bold tracking-[-0.03em] md:text-[2.35rem] md:leading-10">
            Live Hyperliquid. Private 0G research. Host-enforced policy. Proof you can verify. Execution stays on
            desktop.
          </p>
        </div>
        <figure className="mt-14 overflow-hidden border border-[rgb(240_231_212/0.3)]">
          <DiagramWideBanner className="aspect-[21/9] w-full" />
        </figure>
      </div>
    </section>
  );
}
