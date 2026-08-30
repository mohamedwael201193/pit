import { useMemo, useState } from "react";
import { usePrivy } from "@privy-io/react-auth";
import { CAPITAL_PRESETS } from "./facts";
import { markLabel, usd } from "./format";
import type { SimRow } from "./types";
import { coinMin } from "./venue";
import { PageHead } from "../ui/PageHead";
import { eligibleCoins, useWatch } from "./Watch";

export function CapitalPage() {
  const { watch, error, loading } = useWatch();
  const { authenticated, user } = usePrivy();
  const [cap, setCap] = useState(100);
  const eligible = eligibleCoins(watch);

  const rows: SimRow[] = useMemo(() => {
    if (!eligible.length) return [];
    return eligible.map((c) => {
      const min = coinMin(c);
      if (cap + 1e-9 >= min) {
        return {
          coin: c.coin,
          mark: c.mark,
          minNotional: min,
          kind: "TRADE",
          why: `SIMULATION: ${usd(cap)} meets venue min ${usd(min)} on ${c.coin} at mark ${markLabel(c.mark)}. Not a fill. Not your account.`,
        };
      }
      return {
        coin: c.coin,
        mark: c.mark,
        minNotional: min,
        kind: "WAIT",
        why: `SIMULATION: ${usd(cap)} is below venue min ${usd(min)} on ${c.coin}. PIT will not invent size.`,
      };
    });
  }, [eligible, cap]);

  const tradeN = rows.filter((r) => r.kind === "TRADE").length;
  const waitN = rows.filter((r) => r.kind === "WAIT").length;

  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="What could PIT do with this capital?"
        lede="This is a host-style check of public Hyperliquid venue floors against the number you type. It is not a return forecast. It does not execute. It does not use Zia pools or APR."
      />
      <p className="mt-3 inline-block border border-[#d82f2f]/50 px-2 py-1 text-[0.6875rem] font-semibold tracking-[0.16em] text-[#d82f2f]">
        SIMULATION
      </p>

      <div className="mt-8 flex flex-wrap gap-1.5">
        {CAPITAL_PRESETS.map((n) => (
          <button key={n} type="button" className={cap === n ? "intel-cta" : "intel-chip"} onClick={() => setCap(n)}>
            ${n}
          </button>
        ))}
      </div>

      {error ? <p className="mt-4 text-[#ff8a8a]">{error}</p> : null}

      <dl className="intel-metrics mt-8">
        <div className="intel-metric">
          <dt>TRADE</dt>
          <dd>{loading ? "…" : String(tradeN)}</dd>
          <p>Would meet venue min in this simulation</p>
        </div>
        <div className="intel-metric">
          <dt>WAIT</dt>
          <dd>{loading ? "…" : String(waitN)}</dd>
          <p>Below this market’s min</p>
        </div>
        <div className="intel-metric">
          <dt>HOLD</dt>
          <dd>{eligible.length === 0 ? "yes" : "—"}</dd>
          <p>No public-eligible book to simulate</p>
        </div>
        <div className="intel-metric">
          <dt>LIQUIDITY</dt>
          <dd>unavailable</dd>
          <p>No live PIT LP/swap route. Hyperliquid perps only.</p>
        </div>
      </dl>

      <div className="intel-table mt-8">
        <div className="intel-row intel-row-head">
          <span>Asset</span>
          <span>Mark</span>
          <span>Min</span>
          <span className="hidden sm:inline">Kind</span>
          <span className="justify-self-end">Result</span>
        </div>
        {rows.map((r) => (
          <div key={r.coin} className="intel-row">
            <span className="font-semibold">{r.coin}</span>
            <span className="intel-num">${markLabel(r.mark)}</span>
            <span className="intel-num">{usd(r.minNotional)}</span>
            <span className="hidden sm:inline text-[0.75rem] tracking-[0.12em] text-[#d82f2f]">{r.kind}</span>
            <span className="col-span-full text-[0.75rem] leading-5 text-[rgb(240_231_212/0.55)] md:col-span-1 md:justify-self-end">
              {r.why}
            </span>
          </div>
        ))}
      </div>

      <section className="intel-section">
        <h2 className="intel-kicker">Account truth</h2>
        {authenticated ? (
          <p className="mt-3 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">
            Connected as {user?.wallet?.address ? `${user.wallet.address.slice(0, 8)}…${user.wallet.address.slice(-4)}` : "wallet"}.
            That is login identity. This browser does not receive Hyperliquid buying power. Open PIT Desktop for
            account truth.
          </p>
        ) : (
          <p className="mt-3 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.55)]">
            No wallet connected. Simulation still runs on public venue floors. Connecting does not unlock execution
            here.
          </p>
        )}
      </section>
    </div>
  );
}
