import { useMemo, useState } from "react";
import { BrandMark } from "./BrandMark";
import { compactNum, compactUsd, pctFunding } from "./format";

export type MarketCoin = {
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
  volume?: number;
  timestamp?: string;
  venue?: string;
  policyFit?: string;
  riskFlags?: string[];
  provenance?: string;
  block?: string;
};

export function WatchBook({
  coins,
  bestWhy,
  scanned,
  computeReady,
  researchBusy,
  onResearch,
}: {
  coins: MarketCoin[];
  bestWhy?: string;
  scanned?: number;
  computeReady: boolean;
  researchBusy: boolean;
  onResearch: (coin: string) => void;
}) {
  const [sel, setSel] = useState(coins.find((c) => c.eligible)?.coin || coins[0]?.coin || "ETH");
  const [q, setQ] = useState("");
  const [onlyPass, setOnlyPass] = useState(true);
  const filtered = useMemo(() => {
    const n = q.trim().toLowerCase();
    return coins.filter((c) => {
      if (onlyPass && !c.eligible) return false;
      if (!n) return true;
      return c.coin.toLowerCase().includes(n) || (c.why || "").toLowerCase().includes(n);
    });
  }, [coins, q, onlyPass]);
  const best = coins.find((c) => c.eligible);
  const row = filtered.find((c) => c.coin === sel) || filtered[0] || best;
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Markets</p>
          <h1>Discovery universe</h1>
        </div>
        <p className="fine" style={{ margin: 0 }}>
          Live Hyperliquid books. Host rank uses venue facts. Side is not decided here.
        </p>
      </div>
      {best ? (
        <section className="best-strip">
          <div>
            <p className="label">Best opportunity right now</p>
            <h2>
              <BrandMark symbol={best.coin} /> {best.coin} · {compactNum(best.mark)}
            </h2>
            <p>{best.why}</p>
            <p className="fine">
              {bestWhy || "Highest host rank among policy-eligible live books. Not a model score."} · {best.venue || "hyperliquid"} · {best.provenance || "hyperliquid.info"}
            </p>
          </div>
          <button type="button" className="primary" disabled={researchBusy || !computeReady} onClick={() => onResearch(best.coin)}>
            Research privately
          </button>
        </section>
      ) : (
        <p className="empty">No opportunities match your policy. Empty is the honest state until live books arrive.</p>
      )}
      <div className="market-tools">
        <input
          aria-label="Search markets"
          placeholder="Search asset"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />
        <label className="fine">
          <input type="checkbox" checked={onlyPass} onChange={(e) => setOnlyPass(e.target.checked)} /> Policy PASS only
        </label>
        <p className="fine" style={{ margin: 0 }}>
          {scanned ? `${scanned} scanned` : `${coins.length} books`} · {coins.filter((c) => c.eligible).length} PASS
        </p>
      </div>
      {filtered.length === 0 ? (
        <p className="empty">No markets match that filter.</p>
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
          <ul className="market-rows" aria-label="Markets">
            {filtered.map((c) => (
              <li key={c.coin} className={c.coin === sel ? "on" : ""} onClick={() => setSel(c.coin)}>
                <span className="asset">
                  <BrandMark symbol={c.coin} />
                  <strong>{c.coin}</strong>
                </span>
                <span className="mark-num">{compactNum(c.mark)}</span>
                <span>{c.trend || "—"}</span>
                <span>{pctFunding(c.funding)}</span>
                <span>{compactNum(c.openInterest)}</span>
                <span>{c.policyFit || (c.eligible ? "PASS" : "BLOCKED")}</span>
                <span className="why-cell">{c.why || (c.eligible ? "In policy universe." : c.block || "Outside policy.")}</span>
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
              {row.coin} · freshness {row.freshness || "live"} · host rank {row.rank ?? "—"} · volume {compactUsd(row.volume)} ·{" "}
              {(row.riskFlags || []).join(", ") || "no extra risk flags"} · {row.timestamp || ""}
            </p>
          ) : null}
        </>
      )}
    </main>
  );
}
