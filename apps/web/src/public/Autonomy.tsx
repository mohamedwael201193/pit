import { windowsInstallerUrl } from "./facts";
import { PageHead } from "../ui/PageHead";

const BEATS = [
  "WATCH",
  "PRIVATE RESEARCH",
  "0G VERIFY",
  "BOUNDED AUTONOMY",
  "REAL EXECUTION",
  "PROOF",
  "MEMORY",
] as const;

export function AutonomyPage() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Proof-carrying Sleep Missions"
        lede="The website discovers and proves. It never arms a mission, never holds a session key, and never handles a Direct token. Arming happens only on PIT Desktop."
      />
      <p className="mt-8 max-w-[52ch] text-[1.05rem] leading-7 text-[rgb(240_231_212/0.7)]">
        A Sleep Mission is bounded host execution on this computer. You set duration, max trade, max mission notional, max autonomous trades, max daily loss, max open positions, allowed assets, and max data age. Every setting is clamped to the pinned policy ceiling. This computer must stay awake for the bound. If it sleeps, the mission stops. That gap is not backfilled.
      </p>
      <ol className="intel-pipe mt-10">
        {BEATS.map((s, i, arr) => (
          <li key={s}>
            <span>{s}</span>
            {i < arr.length - 1 ? <span className="intel-pipe-arrow">↓</span> : null}
          </li>
        ))}
      </ol>
      <div className="mt-10 border border-[rgb(240_231_212/0.14)] px-5 py-6">
        <p className="text-[0.6875rem] tracking-[0.16em] text-[rgb(240_231_212/0.45)]">THIS BROWSER</p>
        <p className="mt-2 text-[1.25rem] font-semibold">Cannot arm, authorize, pin, or execute</p>
        <p className="mt-2 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.6)]">
          Private strategy remains on desktop. Public pages show hashes, OIDs when public-safe, and named no-trades. They never show prompts, memory, or session keys.
        </p>
        <a href={windowsInstallerUrl()} className="intel-cta mt-6 inline-flex">
          Download PIT Desktop
        </a>
      </div>
    </div>
  );
}
