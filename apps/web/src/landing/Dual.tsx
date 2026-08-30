import { Reveal } from "../ui/Reveal";
import { DiagramWideBanner } from "../diagrams/pitGuide";
import { TokenRow } from "./TokenRow";

export function Dual() {
  return (
    <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit">
        <Reveal>
          <p className="text-[1.125rem] text-[rgb(240_231_212/0.7)]">The product is</p>
          <h2 className="guide-display mt-3 max-w-[16ch]">MAINNET only</h2>
          <p className="mt-4 max-w-[40ch] text-[1.05rem] leading-7 text-[rgb(240_231_212/0.62)]">
            The laboratory exists for CI and developers, not for the public desk.
          </p>
        </Reveal>
        <TokenRow />
        <div className="mt-12 min-h-[14rem] bg-[#d82f2f] p-8 text-black md:p-12">
          <p className="text-[1.75rem] leading-9 font-bold tracking-[-0.03em] md:text-[2.35rem] md:leading-10">
            Live Hyperliquid. Private 0G research. Host-enforced policy. Bounded Sleep Missions on desktop. Proof you can verify. Execution stays on your machine.
          </p>
          <p className="mt-8 text-[1.05rem] leading-7 opacity-85">
            Watch. Private research. 0G verify. Bounded autonomy. Real execution. Proof. Memory.
          </p>
        </div>
        <figure className="mt-14 overflow-hidden border border-[rgb(240_231_212/0.3)]">
          <DiagramWideBanner className="aspect-[21/9] w-full" />
        </figure>
      </div>
    </section>
  );
}
