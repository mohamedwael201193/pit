import { Link } from "react-router-dom";
import type { PublicCoin } from "./types";
import { compact, fundingLabel, markLabel, usd } from "./format";
import { coinMin } from "./venue";

export function statusOf(c: PublicCoin): { label: string; tone: "go" | "mid" | "stop" } {
  if (c.executionFeasible) return { label: "ACTIONABLE", tone: "go" };
  if (c.block) return { label: "BLOCKED", tone: "stop" };
  if (c.eligible) return { label: "READY FOR RESEARCH", tone: "mid" };
  return { label: "WATCH", tone: "mid" };
}

export function MarketRow({ c, dense }: { c: PublicCoin; dense?: boolean }) {
  const st = statusOf(c);
  const min = coinMin(c);
  return (
    <Link
      to={`/radar/${encodeURIComponent(c.coin)}`}
      className="intel-row group"
    >
      <span className="font-semibold tracking-tight">{c.coin}</span>
      <span className="intel-num">${markLabel(c.mark)}</span>
      <span className="intel-num text-[rgb(240_231_212/0.7)]">{fundingLabel(c.funding)}</span>
      <span className="intel-num hidden sm:inline">{compact(c.openInterest)}</span>
      <span className="intel-num hidden md:inline">{usd(min)}</span>
      <span className={`justify-self-end text-[0.6875rem] font-semibold tracking-[0.12em] ${st.tone === "go" ? "text-[#7dffb3]" : st.tone === "stop" ? "text-[#ff8a8a]" : "text-[#d82f2f]"}`}>
        {st.label}
      </span>
      {dense ? null : <span className="col-span-full text-[0.75rem] text-[rgb(240_231_212/0.45)] sm:col-span-1 sm:hidden">{c.why || c.reason || ""}</span>}
    </Link>
  );
}

export function MarketHead() {
  return (
    <div className="intel-row intel-row-head">
      <span>Asset</span>
      <span>Price</span>
      <span>Funding</span>
      <span className="hidden sm:inline">OI</span>
      <span className="hidden md:inline">Min</span>
      <span className="justify-self-end">Status</span>
    </div>
  );
}
