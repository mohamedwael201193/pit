import { Link, useParams } from "react-router-dom";
import { HISTORICAL_FILL, VERIFIED_FILL, windowsInstallerUrl } from "./facts";
import { PageHead } from "../ui/PageHead";

export function MissionDetailPage() {
  const { id = "" } = useParams();
  const historical = id === HISTORICAL_FILL.id;
  const recorded = id === VERIFIED_FILL.id;
  const known = historical || recorded;
  const fill = recorded ? VERIFIED_FILL : HISTORICAL_FILL;

  return (
    <div className="mx-auto max-w-[80rem]">
      <Link to="/missions" className="intel-ghost">
        ← Missions
      </Link>
      <PageHead
        title={known ? `${fill.market} · OID ${fill.oid}` : id || "Mission"}
        lede={
          recorded
            ? `${VERIFIED_FILL.note} Private strategy remains on desktop.`
            : historical
              ? `${HISTORICAL_FILL.note} Private strategy remains on desktop.`
              : "No public-safe mission with that id. PIT will not invent a timeline, fill, or proof."
        }
      />
      <div className="mt-8 border border-[rgb(240_231_212/0.14)] px-5 py-6">
        <p className="text-[0.6875rem] tracking-[0.16em] text-[rgb(240_231_212/0.45)]">PUBLIC-SAFE</p>
        <p className="mt-2 text-[1.125rem] font-semibold">
          {recorded ? "RECORDED fill" : historical ? "HISTORICAL fill" : "Redacted"}
        </p>
        <p className="mt-2 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.6)]">
          {known
            ? `Size ${fill.sz} @ ${fill.px}. This site does not fetch another account's book.`
            : "Desktop missions remain private unless a public-safe proof is published."}
        </p>
        <p className="mt-4 text-[0.8125rem] text-[rgb(240_231_212/0.45)]">Private strategy remains on desktop.</p>
        <div className="mt-6 flex flex-wrap gap-3">
          {known ? (
            <Link to={`/missions/${fill.id}/replay`} className="intel-cta inline-flex">
              Open replay
            </Link>
          ) : null}
          <a href={windowsInstallerUrl()} className="intel-ghost inline-flex">
            Download PIT Desktop
          </a>
        </div>
      </div>
    </div>
  );
}
