import { useMemo, useState } from "react";
import { BrandMark } from "./BrandMark";
import { compactNum, compactUsd, pctFunding } from "./format";

export type MarketCoin = {
  coin: string;
  reason: string;
  why?: string;
  trend?: string;
  rank?: number;
  rankGroup?: number;
  freshness?: string;
  mark: number;
  eligible?: boolean;
  oracle?: number;
  funding?: number;
  openInterest?: number;
  volume?: number;
  szDecimals?: number;
  timestamp?: string;
  venue?: string;
  policyFit?: string;
  researchEligible?: boolean;
  policyEligible?: boolean;
  executionFeasible?: boolean;
  previewReady?: boolean;
  layer?: string;
  riskFlags?: string[];
  provenance?: string;
  block?: string;
  execGate?: string;
  execWhy?: string;
  minNotional?: number;
  requiredMargin?: number;
  availableMargin?: number;
  policyClip?: number;
  hostNotional?: number;
  hostSz?: number;
  estimatedSlippage?: string;
  whyExecutable?: string;
  expectedEdge?: string;
  invalidation?: string;
  whyRanked?: string;
  skillIds?: string[];
};

export function WatchBook({
  coins,
  bestWhy,
  scanned,
  execGate,
  execWhy,
  computeReady,
  researchBusy,
  capitalNote,
  buyingPower,
  powerSource,
  onResearch,
}: {
  coins: MarketCoin[];
  bestWhy?: string;
  scanned?: number;
  execGate?: string;
  execWhy?: string;
  computeReady: boolean;
  researchBusy: boolean;
  capitalNote?: string;
  buyingPower?: number;
  powerSource?: string;
  onResearch: (coin: string) => void;
}) {
  const [sel, setSel] = useState(coins.find((c) => c.previewReady || c.executionFeasible)?.coin || coins.find((c) => c.eligible)?.coin || coins[0]?.coin || "");
  const [q, setQ] = useState("");
  const [onlyPass, setOnlyPass] = useState(true);
  const [onlyExec, setOnlyExec] = useState(false);
  const filtered = useMemo(() => {
    const n = q.trim().toLowerCase();
    return coins.filter((c) => {
      if (onlyExec && !c.executionFeasible) return false;
      if (onlyPass && !c.eligible) return false;
      if (!n) return true;
      return c.coin.toLowerCase().includes(n) || (c.why || "").toLowerCase().includes(n);
    });
  }, [coins, q, onlyPass, onlyExec]);
  const best = coins.find((c) => c.previewReady || c.executionFeasible) || coins.find((c) => c.eligible);
  const row = filtered.find((c) => c.coin === sel) || filtered[0] || best;
  const execN = coins.filter((c) => c.executionFeasible).length;
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Markets</p>
          <h1>Discovery universe</h1>
        </div>
        <p className="fine" style={{ margin: 0 }}>
          Live Hyperliquid books. Host ranks executable size first. Side is not decided here.
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
            <OpportunityFacts coin={best} />
            {execGate ? (
              <p className="fine" role="status">
                Execution blocked: {execGate.replaceAll("_", " ")}. {execWhy || best.execWhy} Research can still run.
              </p>
            ) : null}
            {capitalNote ? <p className="fine">{capitalNote}</p> : null}
            <p className="fine">
              {bestWhy || "Highest host rank among books this account can size. Not a model score."} · {best.venue || "hyperliquid"} · {best.provenance || "hyperliquid.info"}
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
          <input type="checkbox" checked={onlyPass} onChange={(e) => setOnlyPass(e.target.checked)} /> Policy PASS
        </label>
        <label className="fine">
          <input type="checkbox" checked={onlyExec} onChange={(e) => setOnlyExec(e.target.checked)} /> Executable now
        </label>
        <p className="fine" style={{ margin: 0 }}>
          {scanned ? `${scanned} scanned` : `${coins.length} books`} · {coins.filter((c) => c.eligible).length} PASS · {execN} executable
          {typeof buyingPower === "number" ? ` · buying power ${compactUsd(buyingPower)}` : ""}
          {powerSource ? ` · ${powerSource.replaceAll("_", " ")}` : ""}
        </p>
      </div>
      {filtered.length === 0 ? (
        <p className="empty">No markets match that filter.</p>
      ) : (
        <>
          <div className="market-head watch-head">
            <span>Asset</span>
            <span>Mark</span>
            <span>Layer</span>
            <span>Funding</span>
            <span>Host $</span>
            <span>Policy</span>
            <span>Why executable</span>
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
                <span>{(c.layer || (c.eligible ? "policy-eligible" : "blocked")).replaceAll("-", " ")}</span>
                <span>{pctFunding(c.funding)}</span>
                <span>{c.hostNotional ? compactUsd(c.hostNotional) : "—"}</span>
                <span>{c.policyFit || (c.eligible ? "PASS" : "BLOCKED")}</span>
                <span className="why-cell">{c.whyExecutable || c.why || c.block || "—"}</span>
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
            <article className="card" style={{ marginTop: 12 }}>
              <p className="label">{row.coin} card</p>
              <OpportunityFacts coin={row} />
              <p className="fine">
                freshness {row.freshness || "live"} · host rank {row.rank ?? "—"} · skills {(row.skillIds || []).join(", ") || "none"} · {row.timestamp || ""}
              </p>
            </article>
          ) : null}
        </>
      )}
    </main>
  );
}

function OpportunityFacts({ coin }: { coin: MarketCoin }) {
  return (
    <ul className="opp-facts">
      <li>Oracle {coin.oracle ? compactNum(coin.oracle) : "—"} · funding {pctFunding(coin.funding)} · OI {compactNum(coin.openInterest)} · volume {compactUsd(coin.volume)}</li>
      <li>Min notional {compactUsd(coin.minNotional || 10)} · required {compactUsd(coin.requiredMargin || 0)} · available {compactUsd(coin.availableMargin || 0)}</li>
      <li>Policy clip {compactUsd(coin.policyClip || 0)} · host-sized {coin.hostSz ? `${coin.hostSz} / ${compactUsd(coin.hostNotional || 0)}` : "unsized"}</li>
      <li>{coin.estimatedSlippage || "Slippage ceiling is host policy, not a live L2 estimate."}</li>
      <li>{coin.whyExecutable || coin.execWhy || "Not executable for this account right now."}</li>
      <li>{coin.whyRanked || coin.expectedEdge || "Host rank of venue facts."}</li>
      {coin.invalidation ? <li>Invalidation: {coin.invalidation}</li> : null}
      <li>Side is not decided here. Host sizes after sealed research. Chat cannot AUTHORIZE.</li>
    </ul>
  );
}
