import { Link } from "react-router-dom";
import { HISTORICAL_FILL } from "../public/facts";
import { MarketHead, MarketRow } from "../public/MarketRow";
import { actionableCoins, eligibleCoins, useWatch } from "../public/Watch";
import { Reveal } from "../ui/Reveal";

export function LiveTape() {
  const { watch, health, error, loading } = useWatch();
  const eligible = eligibleCoins(watch);
  const actionable = actionableCoins(watch);

  return (
    <section id="live" className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit">
        <Reveal>
          <p className="text-[1.25rem] text-[var(--guide-cream)]">Live PIT</p>
          <h2 className="guide-display mt-4 max-w-[14ch]">What is happening right now?</h2>
          <p className="mt-6 max-w-[46ch] text-[1.2rem] leading-8 text-[rgb(240_231_212/0.78)]">
            Public Hyperliquid marks. Spread and book depth are not invented. Private thesis is sealed.
          </p>
        </Reveal>

        {error ? (
          <p className="mt-10 max-w-[48ch] text-[1.0625rem] text-[#ff8a8a]">{error}</p>
        ) : (
          <dl className="mt-12 grid grid-cols-2 gap-px border border-[rgb(240_231_212/0.25)] bg-[rgb(240_231_212/0.18)] md:grid-cols-4">
            <Stat k="Markets scanned" v={loading ? "…" : String(watch?.scanned ?? 0)} />
            <Stat k="Policy-eligible" v={loading ? "…" : String(watch?.count ?? 0)} />
            <Stat k="Research candidates" v={loading ? "…" : String(eligible.length)} />
            <Stat k="Actionable (public)" v={String(actionable.length)} note="Account size is not on this feed" />
            <Stat k="Verified proofs" v="0 live" note="No public receipt object on this site" />
            <Stat k="Live execution" v="0" note={`Historical fill OID ${HISTORICAL_FILL.oid} is not live`} />
            <Stat k="Autonomy" v="desktop" note="This site cannot arm a Sleep Mission" />
            <Stat k="Health" v={health?.version ?? (loading ? "…" : "—")} />
          </dl>
        )}

        <div className="intel-table mt-10">
          <MarketHead />
          {eligible.slice(0, 6).map((c) => (
            <MarketRow key={c.coin} c={c} />
          ))}
          {!loading && eligible.length === 0 ? (
            <p className="px-4 py-6 text-[0.9375rem] text-[rgb(240_231_212/0.5)]">No policy-eligible books on the public feed.</p>
          ) : null}
        </div>

        <Link to="/radar" className="pill pill-line mt-10">
          Full radar
        </Link>
      </div>
    </section>
  );
}

function Stat({ k, v, note }: { k: string; v: string; note?: string }) {
  return (
    <div className="bg-[#1a1a1a] px-4 py-5">
      <dt className="text-[0.6875rem] tracking-[0.12em] text-[rgb(240_231_212/0.45)] uppercase">{k}</dt>
      <dd className="mt-2 text-[1.35rem] font-semibold tracking-[-0.03em] text-[var(--guide-cream)]">{v}</dd>
      {note ? <p className="mt-2 text-[0.75rem] leading-4 text-[rgb(240_231_212/0.42)]">{note}</p> : null}
    </div>
  );
}
