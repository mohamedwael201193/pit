import { Link } from "react-router-dom";
import { HISTORICAL_FILL } from "./facts";
import { MarketHead, MarketRow } from "./MarketRow";
import { actionableCoins, eligibleCoins, useWatch } from "./Watch";

export function HomeLanding() {
  const { watch, health, error, loading } = useWatch();
  const eligible = eligibleCoins(watch);
  const actionable = actionableCoins(watch);

  return (
    <div>
      <section className="intel-hero">
        <p className="intel-kicker">PIT · PRIVATE AUTONOMOUS TRADING OS</p>
        <h1 className="intel-display">
          Your trading agent doesn't sleep.
          <br />
          Your keys don't leave your machine.
        </h1>
        <p className="intel-lede">
          PIT privately researches live markets with 0G, enforces the policy you set, and can act inside those limits —
          while the sensitive execution surface stays on your computer.
        </p>
        <div className="mt-8 flex flex-wrap gap-2.5">
          <Link to="/radar" className="intel-cta">
            Explore live PIT
          </Link>
          <Link to="/download" className="intel-secondary">
            Download PIT Desktop
          </Link>
          <Link to="/proof" className="intel-ghost">
            Verify a mission
          </Link>
        </div>
        <p className="mt-6 max-w-[52ch] text-[0.875rem] leading-6 text-[rgb(240_231_212/0.5)]">
          The web discovers and proves. The desktop protects and acts. MAINNET only: Aristotle 16661 and Hyperliquid
          mainnet. The laboratory exists for CI and developers, not for the public desk.
        </p>
      </section>

      <section className="intel-section" aria-labelledby="live-pit">
        <h2 id="live-pit" className="intel-kicker">
          Live PIT
        </h2>
        {error ? (
          <p className="mt-4 text-[0.9375rem] text-[#ff8a8a]">{error}</p>
        ) : (
          <dl className="intel-metrics">
            <Metric k="Markets scanned" v={loading ? "…" : String(watch?.scanned ?? 0)} />
            <Metric k="Policy-eligible" v={loading ? "…" : String(watch?.count ?? 0)} />
            <Metric k="Research candidates" v={loading ? "…" : String(eligible.length)} />
            <Metric k="Actionable (public)" v={String(actionable.length)} note="Account size is not on this feed" />
            <Metric k="Verified proofs" v="0 live" note="No public receipt object on this site" />
            <Metric k="Live execution" v="0" note={`Historical fill OID ${HISTORICAL_FILL.oid} is not live`} />
            <Metric k="Autonomy" v="desktop" note="This site cannot read or enable Guarded Autonomy" />
            <Metric k="Health" v={health?.version ?? (loading ? "…" : "—")} />
          </dl>
        )}
      </section>

      <section className="intel-section" aria-labelledby="now">
        <div className="flex items-end justify-between gap-4">
          <h2 id="now" className="intel-title">
            What is happening right now?
          </h2>
          <Link to="/radar" className="intel-ghost">
            Full radar
          </Link>
        </div>
        <p className="mt-2 max-w-[56ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.6)]">
          Public Hyperliquid marks on the default research list. Spread and book depth are not invented. Private thesis
          is sealed.
        </p>
        <div className="intel-table mt-6">
          <MarketHead />
          {eligible.slice(0, 6).map((c) => (
            <MarketRow key={c.coin} c={c} />
          ))}
          {!loading && eligible.length === 0 ? (
            <p className="px-3 py-6 text-[0.875rem] text-[rgb(240_231_212/0.5)]">No policy-eligible books on the public feed.</p>
          ) : null}
        </div>
      </section>

      <section className="intel-section" aria-labelledby="why">
        <h2 id="why" className="intel-title">
          Why PIT
        </h2>
        <ol className="intel-steps">
          <li>
            <strong>Private research with 0G.</strong> Direct TeeML. No Router fallback for the private book.
          </li>
          <li>
            <strong>Host-enforced policy.</strong> Clip, assets, and kill are law the model cannot raise.
          </li>
          <li>
            <strong>Bounded autonomous execution.</strong> Optional, desktop-only, phrase-gated.
          </li>
          <li>
            <strong>Real Hyperliquid execution.</strong> Real OIDs. This site cannot place an order.
          </li>
          <li>
            <strong>Verifiable proof.</strong> What was checked, and how — not a badge that says Verified.
          </li>
        </ol>
      </section>

      <Split />

      <section className="intel-section" id="og" aria-labelledby="og-title">
        <h2 id="og-title" className="intel-title">
          0G is the private OS
        </h2>
        <dl className="intel-grid-2">
          <Pair k="0G Compute" v="Private intelligence. Sealed Direct TeeML. The website never sees the prompt." />
          <Pair k="0G Storage" v="Durable evidence and memory roots when a public-safe receipt is published." />
          <Pair k="0G Chain" v="Verifiable proof. Read the transaction from Aristotle RPC in this browser." />
          <Pair k="Agentic ID" v="Agent identity. iTransfer / iClone is not live on Aristotle mainnet." />
          <Pair k="ERC-8004" v="Reputation and feedback when a registry record exists. No invented ranking." />
          <Pair k="Not claimed" v="DA, mainnet iTransfer, Windows Authenticode, and macOS/Linux packages." />
        </dl>
      </section>

      <section className="intel-section border-t border-[rgb(240_231_212/0.12)]" aria-labelledby="final-cta">
        <h2 id="final-cta" className="intel-display max-w-[16ch]">
          Let PIT watch. Keep execution on your machine.
        </h2>
        <div className="mt-8 flex flex-wrap gap-2.5">
          <Link to="/download" className="intel-cta">
            Download PIT Desktop
          </Link>
          <Link to="/radar" className="intel-secondary">
            Explore live PIT
          </Link>
        </div>
      </section>
    </div>
  );
}

function Metric({ k, v, note }: { k: string; v: string; note?: string }) {
  return (
    <div className="intel-metric">
      <dt>{k}</dt>
      <dd>{v}</dd>
      {note ? <p>{note}</p> : null}
    </div>
  );
}

function Pair({ k, v }: { k: string; v: string }) {
  return (
    <div className="intel-pair">
      <dt>{k}</dt>
      <dd>{v}</dd>
    </div>
  );
}

function Split() {
  return (
    <section className="intel-section" aria-labelledby="split">
      <h2 id="split" className="intel-title">
        0G lets PIT prove the machine without exposing the intelligence.
      </h2>
      <div className="mt-8 grid gap-8 md:grid-cols-2">
        <div>
          <p className="intel-kicker text-[#7dffb3]">Public</p>
          <ul className="intel-list">
            <li>Market</li>
            <li>Decision state</li>
            <li>Policy outcome</li>
            <li>TEE status</li>
            <li>Storage proof</li>
            <li>Chain transaction</li>
            <li>OID</li>
            <li>Fill</li>
            <li>Safe performance metadata</li>
          </ul>
        </div>
        <div>
          <p className="intel-kicker">Private</p>
          <ul className="intel-list">
            <li>Strategy</li>
            <li>Prompt</li>
            <li>Portfolio context</li>
            <li>Private memory</li>
            <li>Sealed reasoning</li>
          </ul>
          <p className="mt-4 text-[0.875rem] text-[rgb(240_231_212/0.5)]">Private reasoning is sealed with 0G Direct.</p>
        </div>
      </div>
    </section>
  );
}
