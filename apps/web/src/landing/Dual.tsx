import { useState } from "react";
import { cn } from "../lib/cn";
import { Reveal } from "../ui/Reveal";
import { DiagramWideBanner } from "../diagrams/pitGuide";
import { NetworkBanner } from "../NetworkBanner";

type Net = "mainnet" | "testnet";

export function Dual() {
  const [net, setNet] = useState<Net>("mainnet");

  return (
    <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit">
        <Reveal>
          <p className="text-[1.125rem] text-[rgb(240_231_212/0.7)]">See the difference between</p>
          <h2 className="guide-display mt-3 max-w-[14ch]">
            lab
            <br />
            and production
          </h2>
        </Reveal>

        <div className="mt-10 flex gap-2">
          {(
            [
              ["mainnet", "MAINNET"],
              ["testnet", "TESTNET"],
            ] as const
          ).map(([id, label]) => (
            <button
              key={id}
              type="button"
              onClick={() => setNet(id)}
              className={cn(
                "rounded-full border px-5 py-2.5 text-[0.9375rem] font-medium transition-colors",
                net === id
                  ? "border-[var(--guide-cream)] bg-[var(--guide-cream)] text-black"
                  : "border-[rgb(240_231_212/0.35)] text-[var(--guide-cream)] hover:border-[var(--guide-cream)]",
              )}
            >
              {label}
            </button>
          ))}
        </div>

        <NetworkBanner net={net} />

        <div className="mt-12 grid gap-0 border border-[rgb(240_231_212/0.35)] md:grid-cols-2">
          <div
            className={cn(
              "min-h-[16rem] p-8 transition-colors duration-300 md:p-12",
              net === "mainnet" ? "bg-[#d82f2f] text-black" : "bg-transparent text-[rgb(240_231_212/0.45)]",
            )}
          >
            <p className="text-[0.8125rem] tracking-[0.16em] uppercase opacity-70">Production</p>
            <p className="mt-6 text-[1.75rem] leading-9 font-bold tracking-[-0.03em] md:text-[2.35rem] md:leading-10">
              Aristotle 16661. Hyperliquid mainnet. Direct glm-5.2. Transfer of Agentic ID is not live.
            </p>
          </div>
          <div
            className={cn(
              "min-h-[16rem] p-8 transition-colors duration-300 md:p-12",
              net === "testnet" ? "bg-[#f0e7d4] text-black" : "bg-transparent text-[rgb(240_231_212/0.45)]",
            )}
          >
            <p className="text-[0.8125rem] tracking-[0.16em] uppercase opacity-70">Laboratory</p>
            <p className="mt-6 text-[1.75rem] leading-9 font-bold tracking-[-0.03em] md:text-[2.35rem] md:leading-10">
              Galileo 16602. Hyperliquid testnet. Sealed ask stays off until VerifyE2EE is proven. Different model
              catalog than production.
            </p>
          </div>
        </div>

        <figure className="mt-14 overflow-hidden border border-[rgb(240_231_212/0.3)]">
          <DiagramWideBanner className="aspect-[21/9] w-full" />
        </figure>
      </div>
    </section>
  );
}
