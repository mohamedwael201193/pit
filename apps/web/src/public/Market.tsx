import { Link, useParams } from "react-router-dom";
import { compact, fundingLabel, markLabel, usd } from "./format";
import { coinMin } from "./venue";
import { findCoin, useWatch } from "./Watch";
import { statusOf } from "./MarketRow";

export function MarketPage() {
  const { coin = "" } = useParams();
  const { watch, loading, error } = useWatch();
  const c = findCoin(watch, coin);

  if (loading) return <p className="text-[rgb(240_231_212/0.55)]">Loading public book…</p>;
  if (error) return <p className="text-[#ff8a8a]">{error}</p>;
  if (!c) {
    return (
      <div>
        <p className="intel-kicker">Public market</p>
        <h1 className="intel-title mt-2">{coin.toUpperCase() || "Unknown"}</h1>
        <p className="intel-lede">This asset is not in the current public watch payload. PIT will not invent a mark.</p>
        <Link to="/radar" className="intel-ghost mt-6 inline-flex">
          Back to radar
        </Link>
      </div>
    );
  }

  const min = coinMin(c);
  const st = statusOf(c);

  return (
    <div>
      <Link to="/radar" className="intel-ghost">
        ← Radar
      </Link>
      <p className="intel-kicker mt-6">Public market</p>
      <h1 className="intel-title mt-2">{c.coin}</h1>
      <p className="mt-2 text-[0.75rem] tracking-[0.14em] text-[#d82f2f]">{st.label}</p>

      <dl className="intel-metrics mt-8">
        <Item k="Price" v={`$${markLabel(c.mark)}`} />
        <Item k="Oracle" v={c.oracle ? `$${markLabel(c.oracle)}` : "—"} />
        <Item k="Funding" v={fundingLabel(c.funding)} />
        <Item k="Open interest" v={compact(c.openInterest)} />
        <Item k="24h notional" v={compact(c.volume)} note="Venue day notional, not book depth" />
        <Item k="Spread" v="—" note="Not on this public feed" />
        <Item k="Minimum notional" v={usd(min)} />
        <Item k="Venue rules" v="Hyperliquid perp · $10 floor, then size tick" />
      </dl>

      <section className="intel-section">
        <h2 className="intel-kicker">PIT observation</h2>
        <ul className="intel-list mt-4">
          <li>Why PIT is watching: {c.why || c.reason || "live venue fields on the default research list"}</li>
          <li>Policy compatibility: {c.eligible ? "PASS on the public default policy" : "not on the public eligible list"}</li>
          <li>Capital constraint: not attached on the website origin</li>
          <li>Research readiness: {c.eligible ? "ready for private research on desktop" : "not research-eligible here"}</li>
          <li>Public research state: {c.freshness === "live" ? "LIVE marks" : c.freshness || "unknown"}</li>
        </ul>
      </section>

      <section className="intel-seal">
        <p className="intel-kicker">Private reasoning</p>
        <p className="mt-3 text-[1.25rem] font-semibold">Sealed</p>
        <p className="mt-2 max-w-[46ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">
          Private reasoning is sealed with 0G Direct. This page will never show the prompt, thesis, or private memory.
        </p>
      </section>
    </div>
  );
}

function Item({ k, v, note }: { k: string; v: string; note?: string }) {
  return (
    <div className="intel-metric">
      <dt>{k}</dt>
      <dd>{v}</dd>
      {note ? <p>{note}</p> : null}
    </div>
  );
}
