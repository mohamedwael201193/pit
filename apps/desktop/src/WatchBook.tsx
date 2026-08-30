import { useMemo, useState } from "react";
import { BrandMark } from "./BrandMark";
import { accountSizeGate, compactNum, compactUsd, marketSizeGate, nearestVenueMin, pctFunding, powerSourceLabel } from "./format";
import { ExternalLink } from "./ExternalLink";

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

function researchTitle(pinned?: boolean, computeReady?: boolean) {
  if (!pinned) {
    return "Research privately — blocked: policy is not pinned. Protect and compute are not the issue if they are ready. Pin on Security.";
  }
  if (!computeReady) {
    return "Research privately — blocked: Protect or private compute is not ready.";
  }
  return "Start sealed private research for this book.";
}

function layerChip(c: MarketCoin) {
  if (c.previewReady) return { t: "Preview", k: "ok" as const };
  if (c.executionFeasible) return { t: "Can open", k: "ok" as const };
  if (c.policyEligible || c.eligible) return { t: "Policy", k: "pass" as const };
  if (c.researchEligible) return { t: "Research", k: "wait" as const };
  return { t: "Blocked", k: "blocked" as const };
}

export function WatchBook({
  coins,
  bestWhy,
  scanned,
  execWhy,
  computeReady,
  researchBusy,
  capitalNote,
  buyingPower,
  powerSource,
  fundHref,
  pinned,
  onPin,
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
  fundHref?: string;
  pinned?: boolean;
  onPin?: () => void;
  onResearch: (coin: string) => void;
}) {
  const [sel, setSel] = useState("");
  const [q, setQ] = useState("");
  const [filter, setFilter] = useState<"all" | "pass" | "exec" | "research" | "blocked">("pass");
  const filtered = useMemo(() => {
    const n = q.trim().toLowerCase();
    return coins.filter((c) => {
      if (filter === "exec" && !c.executionFeasible) return false;
      if (filter === "pass" && !c.eligible) return false;
      if (filter === "research" && !c.researchEligible) return false;
      if (filter === "blocked" && (c.executionFeasible || c.eligible)) return false;
      if (!n) return true;
      return c.coin.toLowerCase().includes(n) || (c.why || "").toLowerCase().includes(n);
    });
  }, [coins, q, filter]);
  const execN = coins.filter((c) => c.executionFeasible).length;
  const passN = coins.filter((c) => c.eligible).length;
  const venueMin = nearestVenueMin(coins);
  const have = typeof buyingPower === "number" ? buyingPower : coins[0]?.availableMargin;
  const gate = accountSizeGate(have, venueMin, execN);
  const row = filtered.find((c) => c.coin === sel);
  const useTiles = false;
  const counts = {
    all: coins.length,
    pass: passN,
    exec: execN,
    research: coins.filter((c) => c.researchEligible).length,
    blocked: coins.filter((c) => !c.executionFeasible && !c.eligible).length,
  };

  return (
    <main className="page dense watch-terminal">
      <div className="page-head">
        <div>
          <p className="eyebrow">Markets</p>
          <h1>Opportunity terminal</h1>
        </div>
        <p className="fine" style={{ margin: 0 }}>
          Live Hyperliquid books. Host ranks size for this account. Side is not decided here. Chat cannot AUTHORIZE.
        </p>
      </div>

      <section className={`capital-gate ${gate.canOpen ? "open" : "short"}`} role="status">
        <div className="gate-copy">
          <p className="label">{gate.canOpen ? "Executable" : "Capital gate"}</p>
          <h2>{gate.headline}</h2>
          <p>{gate.detail}</p>
          {capitalNote ? <p className="fine">{capitalNote}</p> : null}
          {execWhy && execWhy !== gate.detail ? <p className="fine">{execWhy}</p> : null}
          {bestWhy && execN > 0 ? <p className="fine">{bestWhy}</p> : null}
          <p className="fine">
            {scanned ? `${scanned} scanned` : `${coins.length} books`} · {passN} policy eligible · {execN} executable
            {powerSource ? ` · ${powerSourceLabel(powerSource)}` : ""}
          </p>
        </div>
        <div className="gate-stats">
          <div>
            <span>This account</span>
            <strong>{compactUsd(gate.have)}</strong>
          </div>
          <div>
            <span>This market min</span>
            <strong>{compactUsd(gate.min)}</strong>
          </div>
          <div>
            <span>{gate.canOpen ? "Headroom" : "Shortfall"}</span>
            <strong>{compactUsd(gate.canOpen ? Math.max(0, gate.have - gate.min) : gate.shortfall)}</strong>
          </div>
          <div className="gate-meter" aria-label={`${compactUsd(gate.have)} of ${compactUsd(gate.min)}`}>
            <span style={{ width: `${Math.min(100, gate.min > 0 ? (gate.have / gate.min) * 100 : 0)}%` }} />
          </div>
          <div className="cta-row">
            {!pinned && onPin ? (
              <button type="button" className="primary" onClick={onPin}>
                Pin a trading policy
              </button>
            ) : execN === 0 && fundHref ? (
              <ExternalLink className="primary" href={fundHref}>
                Fund this Hyperliquid account
              </ExternalLink>
            ) : row ? (
              <button
                type="button"
                className="primary"
                disabled={researchBusy || !computeReady || !pinned}
                title={researchTitle(pinned, computeReady)}
                onClick={() => onResearch(row.coin)}
              >
                Research {row.coin} privately
              </button>
            ) : null}
          </div>
        </div>
      </section>

      <div className="market-tools">
        <input aria-label="Search markets" placeholder="Search asset" value={q} onChange={(e) => setQ(e.target.value)} />
        <div className="filter-row" role="tablist" aria-label="Market filters">
          {(
            [
              ["pass", `Policy ${counts.pass}`],
              ["exec", `Executable ${counts.exec}`],
              ["research", `Research ${counts.research}`],
              ["blocked", `Blocked ${counts.blocked}`],
              ["all", `All ${counts.all}`],
            ] as const
          ).map(([id, label]) => (
            <button
              key={id}
              type="button"
              role="tab"
              aria-selected={filter === id}
              className={filter === id ? "filter-chip on" : "filter-chip"}
              onClick={() => setFilter(id)}
            >
              {label}
            </button>
          ))}
        </div>
      </div>

      {filtered.length === 0 ? (
        <p className="empty">
          {coins.length === 0
            ? "No opportunities match your policy. Empty is the honest state until live books arrive."
            : "No markets match that filter."}
        </p>
      ) : useTiles ? (
        <ul className="book-grid" aria-label="Markets">
          {filtered.map((c) => {
            const chip = layerChip(c);
            const gap = marketSizeGate(c.coin, c.availableMargin ?? have, c.minNotional || venueMin, c.executionFeasible);
            return (
              <li key={c.coin}>
                <button
                  type="button"
                  className={c.coin === sel ? "book-tile on" : "book-tile"}
                  onClick={() => setSel((cur) => (cur === c.coin ? "" : c.coin))}
                >
                  <span className="tile-head">
                    <BrandMark symbol={c.coin} size={16} />
                    <strong>{c.coin}</strong>
                    <span className={`layer-chip ${chip.k}`}>{chip.t}</span>
                  </span>
                  <span className="tile-mark">{compactNum(c.mark)}</span>
                  <span className="tile-meta">
                    {pctFunding(c.funding)} funding · {gap.chip}
                  </span>
                </button>
              </li>
            );
          })}
        </ul>
      ) : (
        <ul className="book-list" aria-label="Markets">
          {filtered.map((c) => {
            const chip = layerChip(c);
            const gap = marketSizeGate(c.coin, c.availableMargin ?? have, c.minNotional || venueMin, c.executionFeasible);
            return (
              <li key={c.coin} className={c.coin === sel ? "on" : ""}>
                <button type="button" className="book-row" onClick={() => setSel((cur) => (cur === c.coin ? "" : c.coin))}>
                  <span className="asset">
                    <BrandMark symbol={c.coin} size={14} />
                    <strong>{c.coin}</strong>
                  </span>
                  <span className="tile-mark">{compactNum(c.mark)}</span>
                  <span className={`layer-chip ${chip.k}`}>{chip.t}</span>
                  <span className="fine" style={{ margin: 0 }}>
                    {gap.chip}
                  </span>
                </button>
              </li>
            );
          })}
        </ul>
      )}

      {row ? (
        <InspectCard
          coin={row}
          have={row.availableMargin ?? have}
          venueMin={row.minNotional || venueMin}
          pinned={pinned}
          computeReady={computeReady}
          researchBusy={researchBusy}
          onResearch={onResearch}
        />
      ) : filtered.length > 0 ? (
        <p className="fine inspect-hint">Select a book for venue facts. Rank is host order, not a model score.</p>
      ) : null}
    </main>
  );
}

function InspectCard({
  coin,
  have,
  venueMin,
  pinned,
  computeReady,
  researchBusy,
  onResearch,
}: {
  coin: MarketCoin;
  have?: number;
  venueMin: number;
  pinned?: boolean;
  computeReady: boolean;
  researchBusy: boolean;
  onResearch: (coin: string) => void;
}) {
  const gap = marketSizeGate(coin.coin, have, venueMin, coin.executionFeasible);
  const chip = layerChip(coin);
  return (
    <article className="inspect-card">
      <div className="inspect-head">
        <div>
          <p className="label">
            {coin.coin} · {chip.t}
          </p>
          <h2>
            <BrandMark symbol={coin.coin} size={16} /> {coin.coin}
            <span className="tile-mark">{compactNum(coin.mark)}</span>
          </h2>
          <p>{gap.detail}</p>
        </div>
        <button
          type="button"
          className="primary"
          disabled={researchBusy || !coin.eligible || !computeReady || !pinned}
          title={researchTitle(pinned, computeReady)}
          onClick={() => onResearch(coin.coin)}
        >
          Research privately
        </button>
      </div>
      <p className="fine">
        {coin.whyRanked || "Host rank of venue facts. Side is not decided here."} · {coin.venue || "hyperliquid"} ·{" "}
        {coin.provenance || "hyperliquid.info"}
      </p>
      <details className="facts-disclosure">
        <summary>Venue facts</summary>
        <ul className="opp-facts">
          <li>
            Research {coin.researchEligible ? "yes" : "no"} · Policy {coin.policyEligible || coin.eligible ? "yes" : "no"} · Execution{" "}
            {coin.executionFeasible ? "feasible" : "blocked"} · Preview {coin.previewReady ? "ready" : "no"}
          </li>
          <li>
            Oracle {coin.oracle ? compactNum(coin.oracle) : "—"} · funding {pctFunding(coin.funding)} · OI {compactNum(coin.openInterest)} ·
            volume {compactUsd(coin.volume)}
          </li>
          <li>
            This market min {compactUsd(coin.minNotional || 10)} · required {compactUsd(coin.requiredMargin || 0)} · available{" "}
            {compactUsd(coin.availableMargin || 0)}
          </li>
          <li>Policy clip {compactUsd(coin.policyClip || 0)} · host-sized {coin.hostSz ? `${coin.hostSz} / ${compactUsd(coin.hostNotional || 0)}` : "unsized"}</li>
          <li>{coin.estimatedSlippage || "Slippage ceiling is host policy, not a live L2 estimate."}</li>
          <li>Chat cannot AUTHORIZE. Host sizes after sealed research.</li>
        </ul>
        <p className="fine">
          freshness {coin.freshness || "live"} · host rank {coin.rank ?? "—"} · skills {(coin.skillIds || []).join(", ") || "none"} ·{" "}
          {coin.timestamp || ""}
        </p>
      </details>
    </article>
  );
}
