import { Link, useParams } from "react-router-dom";
import { HISTORICAL_FILL } from "./facts";
import { PageHead } from "../ui/PageHead";

export function MissionDetailPage() {
  const { id = "" } = useParams();
  const historical = id === HISTORICAL_FILL.id;

  return (
    <div className="mx-auto max-w-[80rem]">
      <Link to="/missions" className="intel-ghost">
        ← Missions
      </Link>
      <PageHead
        title={historical ? `${HISTORICAL_FILL.market} · OID ${HISTORICAL_FILL.oid}` : id || "Mission"}
        lede={
          historical
            ? `${HISTORICAL_FILL.note} Private strategy remains on desktop.`
            : "No public-safe mission with that id. PIT will not invent a timeline, fill, or proof."
        }
      />
      <div className="mt-8 border border-[rgb(240_231_212/0.14)] px-5 py-6">
        <p className="text-[0.6875rem] tracking-[0.16em] text-[rgb(240_231_212/0.45)]">PUBLIC-SAFE</p>
        <p className="mt-2 text-[1.125rem] font-semibold">{historical ? "HISTORICAL fill" : "Redacted"}</p>
        <p className="mt-2 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.6)]">
          {historical
            ? `Size ${HISTORICAL_FILL.sz} @ ${HISTORICAL_FILL.px}. This site does not fetch another account's book.`
            : "Desktop missions remain private unless a public-safe proof is published."}
        </p>
        <p className="mt-4 text-[0.8125rem] text-[rgb(240_231_212/0.45)]">Private strategy remains on desktop.</p>
        <div className="mt-6 flex flex-wrap gap-3">
          {historical ? (
            <Link to={`/missions/${HISTORICAL_FILL.id}/replay`} className="intel-cta inline-flex">
              Open replay
            </Link>
          ) : null}
          <Link to="/download" className="intel-ghost inline-flex">
            Open PIT Desktop
          </Link>
        </div>
      </div>
    </div>
  );
}
