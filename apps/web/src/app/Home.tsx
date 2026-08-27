import { useState } from "react";
import { Link } from "react-router-dom";
import { PageHead } from "../ui/PageHead";
import { EmptyWatch } from "../EmptyWatch";
import { NetworkToggle } from "../NetworkToggle";
import { NetworkBanner } from "../NetworkBanner";
import { ButtonLink } from "../ui/Button";
import { cn } from "../lib/cn";

type Net = "mainnet" | "testnet";

export function Home() {
  const [net, setNet] = useState<Net>("mainnet");
  return (
    <div className="mx-auto flex max-w-[80rem] flex-col gap-10">
      <PageHead
        title="What needs your attention"
        lede="Live books only. PIT may research and notify. PIT may not sign automatically. Authorize lives on desktop."
      />
      <NetworkToggle net={net} onChange={setNet} />
      <NetworkBanner net={net} />
      <EmptyWatch network={net} />
      <div className="grid gap-4 sm:grid-cols-2">
        <Link
          to="/app/start"
          className={cn(
            "rounded-2xl border border-[rgb(240_231_212/0.25)] bg-[#141414] p-6 text-left no-underline transition-colors",
            "hover:border-[#d82f2f]/50 active:scale-[0.99]",
          )}
        >
          <h3 className="text-[1.25rem] font-semibold tracking-[-0.03em] text-[var(--guide-cream)]">Finish setup</h3>
          <p className="mt-2 text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">
            Twelve beats. Wallet, network, policy, then a session on the machine.
          </p>
        </Link>
        <Link
          to="/app/policy"
          className={cn(
            "rounded-2xl border border-[rgb(240_231_212/0.25)] bg-[#141414] p-6 text-left no-underline transition-colors",
            "hover:border-[#d82f2f]/50 active:scale-[0.99]",
          )}
        >
          <h3 className="text-[1.25rem] font-semibold tracking-[-0.03em] text-[var(--guide-cream)]">Read the law</h3>
          <p className="mt-2 text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">
            Clip, assets, kill. The model cannot raise them.
          </p>
        </Link>
      </div>
      <p>
        <ButtonLink as={Link} to="/app/start" trailingArrow size="lg">
          Resume onboarding
        </ButtonLink>
      </p>
    </div>
  );
}
