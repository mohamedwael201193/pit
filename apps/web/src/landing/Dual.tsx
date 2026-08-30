import { Reveal } from "../ui/Reveal";
import { DiagramWideBanner } from "../diagrams/pitGuide";
import { TokenRow } from "./TokenRow";

export function Dual() {
  return (
    <section className="border-t border-[rgb(240_231_212/0.25)]">
      <div className="container-pit flex flex-col items-center py-14 text-center md:py-16">
        <Reveal>
          <p className="text-[1.125rem] text-[rgb(240_231_212/0.7)]">The product is</p>
          <h2 className="guide-display mt-3">MAINNET</h2>
          <p className="mx-auto mt-4 max-w-[40ch] text-[1.05rem] leading-7 text-[rgb(240_231_212/0.62)]">
            The laboratory exists for CI and developers, not for the public desk.
          </p>
        </Reveal>
        <TokenRow />
      </div>
      <div className="container-pit pb-20 md:pb-28">
        <div className="bg-[#d82f2f] p-8 text-black md:p-12">
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
