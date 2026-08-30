import { useMemo, useState } from "react";
import { actionableCoins, blockedCoins, eligibleCoins, useWatch, watchCoins } from "./Watch";
import { MarketHead, MarketRow } from "./MarketRow";
import type { PublicCoin } from "./types";
import { PageHead } from "../ui/PageHead";

const TABS = ["ACTIONABLE", "RESEARCH", "WATCH", "BLOCKED"] as const;
type Tab = (typeof TABS)[number];

export function RadarPage() {
  const { watch, error, loading } = useWatch();
  const [tab, setTab] = useState<Tab>("RESEARCH");

  const groups = useMemo(() => {
    return {
      ACTIONABLE: actionableCoins(watch),
      RESEARCH: eligibleCoins(watch),
      WATCH: watchCoins(watch),
      BLOCKED: blockedCoins(watch),
    };
  }, [watch]);

  const rows: PublicCoin[] = groups[tab];

  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="What is happening right now?"
        lede="Public market intelligence from Hyperliquid via PIT health. This browser cannot size your account. Actionable here means the public feed marked executionFeasible. Website origins do not receive buying power."
      />
      {error ? <p className="mt-4 text-[#ff8a8a]">{error}</p> : null}
      <div className="mt-8 flex flex-wrap gap-1" role="tablist" aria-label="Radar tabs">
        {TABS.map((t) => (
          <button
            key={t}
            type="button"
            role="tab"
            aria-selected={tab === t}
            className={tab === t ? "intel-cta" : "intel-chip"}
            onClick={() => setTab(t)}
          >
            {t}
            <span className="ml-2 text-[0.6875rem] text-[rgb(240_231_212/0.55)]">{loading ? "…" : groups[t].length}</span>
          </button>
        ))}
      </div>
      {tab === "ACTIONABLE" ? (
        <p className="mt-4 text-[0.875rem] leading-6 text-[rgb(240_231_212/0.55)]">
          Public health does not attach capital. If this list is empty, that is honest — not a fake opportunity.
        </p>
      ) : null}
      <div className="intel-table mt-6">
        <MarketHead />
        {rows.map((c) => (
          <MarketRow key={c.coin} c={c} />
        ))}
        {!loading && rows.length === 0 ? (
          <p className="px-3 py-8 text-[0.875rem] text-[rgb(240_231_212/0.5)]">No rows in this tab.</p>
        ) : null}
      </div>
    </div>
  );
}

