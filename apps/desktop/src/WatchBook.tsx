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

type EventRow = { ts?: number; kind?: string; market?: string; status?: string; job_id?: string };

export function WatchBook({
  coins,
  computeReady,
  researchBusy,
  recent,
  onResearch,
}: {
  coins: Coin[];
  computeReady: boolean;
  researchBusy: boolean;
  recent?: EventRow[];
  onResearch: (coin: string) => void;
}) {
  const [sel, setSel] = useState(coins[0]?.coin || "ETH");
  const row = coins.find((c) => c.coin === sel) || coins[0];
  const related = (recent || []).filter((e) => !row || String(e.market || "").includes(row.coin)).slice(-4);
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Watch</p>
          <h1>Market book</h1>
        </div>
      </div>
      <p className="lead">Public Hyperliquid marks only. This surface cannot place an order. No invented scores.</p>
      {coins.length === 0 ? (
        <p className="fine">Empty Watch is the honest state until live books arrive.</p>
      ) : (
        <div className="book-cols">
          <div>
            <div className="market-head">
              <span>Asset</span>
              <span>Mark</span>
              <span>Oracle</span>
              <span>Funding</span>
              <span>OI</span>
              <span>Policy</span>
              <span>Research</span>
              <span></span>
            </div>
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
          </div>
          {row ? (
            <aside className="inspect" aria-label={`${row.coin} snapshot`}>
              <p className="label">Selected</p>
              <h2>{row.coin}</h2>
              <p className="mark-num">{row.mark}</p>
              <dl className="status-grid">
                <dt>Oracle</dt>
                <dd>{row.oracle ?? "—"}</dd>
                <dt>Funding</dt>
                <dd>{row.funding ?? "—"}</dd>
                <dt>Open interest</dt>
                <dd>{row.openInterest ? Math.round(row.openInterest) : "—"}</dd>
                <dt>Policy</dt>
                <dd>{row.eligible ? "PASS" : "BLOCKED"}</dd>
                <dt>Research</dt>
                <dd>{computeReady ? "Ready" : "Needs private compute"}</dd>
              </dl>
              <p className="fine">{row.reason}</p>
              <button type="button" className="primary" onClick={() => onResearch(row.coin)} disabled={researchBusy || !row.eligible}>
                Research {row.coin}
              </button>
              {related.length ? (
                <>
                  <p className="label" style={{ marginTop: 14 }}>
                    Recent
                  </p>
                  <ul className="timeline">
                    {related.map((e, i) => (
                      <li key={`${e.ts}-${i}`}>
                        {e.kind || "event"} {e.status || ""}
                      </li>
                    ))}
                  </ul>
                </>
              ) : (
                <p className="fine">No activity for this market yet.</p>
              )}
            </aside>
          ) : null}
        </div>
      )}
      <p className="fine">Confidence NOT ENOUGH DATA. Side is not decided on this surface.</p>
    </main>
  );
}
