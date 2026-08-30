import { Link } from "react-router-dom";
import { HISTORICAL_FILL } from "./facts";
import { PageHead } from "../ui/PageHead";

export function MissionsPage() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Public-safe PIT missions"
        lede="The public site only lists missions that have a public-safe receipt. Private research stays on desktop. This page will not invent Mission IDs, OIDs, or fills."
      />

      <div className="mt-10 border border-[rgb(240_231_212/0.14)] px-5 py-6">
        <p className="text-[0.6875rem] tracking-[0.16em] text-[rgb(240_231_212/0.45)]">LIVE</p>
        <p className="mt-2 text-[1.25rem] font-semibold">No live public mission</p>
        <p className="mt-2 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.6)]">
          Empty is honest. Desktop missions remain private unless a public-safe proof is published.
        </p>
      </div>

      <div className="mt-6 border border-[rgb(240_231_212/0.14)] px-5 py-6">
        <p className="text-[0.6875rem] tracking-[0.16em] text-[#d82f2f]">HISTORICAL</p>
        <p className="mt-2 text-[1.25rem] font-semibold">
          {HISTORICAL_FILL.market} · OID {HISTORICAL_FILL.oid}
        </p>
        <p className="mt-2 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.6)]">
          Size {HISTORICAL_FILL.sz} @ {HISTORICAL_FILL.px}. {HISTORICAL_FILL.note}
        </p>
        <ol className="intel-pipe mt-6">
          {["DISCOVERED", "PRIVATE RESEARCH", "RISK PASSED", "POLICY PASSED", "EXECUTED", "FILLED", "PROOF ANCHORED"].map(
            (s, i, arr) => (
              <li key={s}>
                <span>{s}</span>
                {i < arr.length - 1 ? <span className="intel-pipe-arrow">↓</span> : null}
              </li>
            ),
          )}
        </ol>
        <p className="mt-4 text-[0.8125rem] text-[rgb(240_231_212/0.45)]">
          Stages after fill are labeled HISTORICAL. Proof anchored is not claimed without a public storage root.
        </p>
        <Link to={`/missions/${HISTORICAL_FILL.id}/replay`} className="intel-cta mt-6 inline-flex">
          Open replay
        </Link>
      </div>
    </div>
  );
}
