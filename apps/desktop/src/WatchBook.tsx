import { useState } from "react";

type Coin = {
  coin: string;
  reason: string;
  mark: number;
  eligible?: boolean;
  oracle?: number;
  funding?: number;
  openInterest?: number;
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
      <p className="eyebrow">Watch</p>
      <h1>Market book</h1>
      <p className="lead">Public Hyperliquid marks only. This surface cannot place an order. No invented scores.</p>
      {coins.length === 0 ? (
        <p className="fine">Empty Watch is the honest state until live books arrive.</p>
      ) : (
        <div className="desk-grid">
          <ul className="market-rows" aria-label="Markets">
            {coins.map((c) => (
              <li key={c.coin} className={c.coin === sel ? "on" : ""} onClick={() => setSel(c.coin)}>
                <strong>{c.coin}</strong>
                <span className="mark-num">{c.mark}</span>
                <span>{c.oracle ?? "—"}</span>
                <span>{c.funding ?? "—"}</span>
                <span>{c.openInterest ? Math.round(c.openInterest) : "—"}</span>
                <span>{c.eligible ? "PASS" : "BLOCKED"}</span>
                <span>{computeReady ? "Ready" : "Needs compute"}</span>
                <button
                  type="button"
                  className="linkish"
                  disabled={researchBusy || !c.eligible}
                  onClick={(e) => {
                    e.stopPropagation();
                    onResearch(c.coin);
                  }}
                >
                  Research
                </button>
              </li>
            ))}
          </ul>
          {row ? (
            <article className="card">
              <p className="label">Selected</p>
              <h2 style={{ margin: "4px 0 8px" }}>{row.coin}</h2>
              <p className="mark-num">{row.mark}</p>
              <p>Oracle {row.oracle ?? "—"}</p>
              <p>Funding {row.funding ?? "—"}</p>
              <p>Open interest {row.openInterest ? Math.round(row.openInterest) : "—"}</p>
              <p>Policy {row.eligible ? "PASS" : "BLOCKED"}</p>
              <p>Research {computeReady ? "Ready" : "Needs private compute"}</p>
              <p className="fine">{row.reason}</p>
              <button type="button" className="primary" onClick={() => onResearch(row.coin)} disabled={researchBusy || !row.eligible}>
                Research {row.coin}
              </button>
            </article>
          ) : null}
        </div>
      )}
      <p className="fine">Confidence NOT ENOUGH DATA. Side is not decided on this surface.</p>
    </main>
  );
}
