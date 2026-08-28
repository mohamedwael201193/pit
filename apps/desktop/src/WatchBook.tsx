import { useState } from "react";
import { BrandMark } from "./BrandMark";

type Coin = {
  coin: string;
  reason: string;
  why?: string;
  trend?: string;
  rank?: number;
  freshness?: string;
  mark: number;
  eligible?: boolean;
  oracle?: number;
  funding?: number;
  openInterest?: number;
  timestamp?: string;
};

export function WatchBook({
  coins,
  computeReady,
  researchBusy,
  onResearch,
}: {
  coins: Coin[];
  computeReady: boolean;
  researchBusy: boolean;
  onResearch: (coin: string) => void;
}) {
  const [sel, setSel] = useState(coins[0]?.coin || "ETH");
  const row = coins.find((c) => c.coin === sel) || coins[0];
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Watch</p>
          <h1>Discovery</h1>
        </div>
        <p className="fine" style={{ margin: 0 }}>
          Live Hyperliquid marks. No invented scores. Side is not decided here.
        </p>
      </div>
      {coins.length === 0 ? (
        <p className="empty">No opportunities match your policy. Empty is the honest state until live books arrive.</p>
      ) : (
        <>
          <div className="market-head watch-head">
            <span>Asset</span>
            <span>Mark</span>
            <span>Trend</span>
            <span>Funding</span>
            <span>OI</span>
            <span>Policy</span>
            <span>Why</span>
            <span></span>
          </div>
          <ul className="market-rows" aria-label="Opportunities">
            {coins.map((c) => (
              <li key={c.coin} className={c.coin === sel ? "on" : ""} onClick={() => setSel(c.coin)}>
                <span className="asset">
                  <BrandMark symbol={c.coin} />
                  <strong>{c.coin}</strong>
                </span>
                <span className="mark-num">{c.mark}</span>
                <span>{c.trend || "—"}</span>
                <span>{c.funding ?? "—"}</span>
                <span>{c.openInterest ? Math.round(c.openInterest) : "—"}</span>
                <span>{c.eligible ? "PASS" : "BLOCKED"}</span>
                <span className="why-cell">{c.why || "In policy universe."}</span>
                <button
                  type="button"
                  className="primary"
                  disabled={researchBusy || !c.eligible || !computeReady}
                  onClick={(e) => {
                    e.stopPropagation();
                    onResearch(c.coin);
                  }}
                >
                  Research privately
                </button>
              </li>
            ))}
          </ul>
          {row ? (
            <p className="fine">
              {row.coin} · freshness {row.freshness || "live"} · host rank {row.rank ?? "—"} (venue facts, not a model)
              {row.timestamp ? ` · ${row.timestamp}` : ""}.{" "}
              {computeReady ? "Private compute is funded." : "Protect and fund private compute before sealed research."}
            </p>
          ) : null}
        </>
      )}
    </main>
  );
}
