import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { AGENT_8004_ID, ARISTOTLE_EXPLORER, IDENTITY_8004, PIT_AGENT, REPUTATION_8004 } from "./facts";
import { shortAddr } from "./format";
import { read8004Owner, type AgentRead } from "./chain";

export function AgentPage() {
  const [reg, setReg] = useState<AgentRead | null>(null);

  useEffect(() => {
    void read8004Owner(AGENT_8004_ID).then(setReg);
  }, []);

  return (
    <div>
      <p className="intel-kicker">PIT agent passport</p>
      <h1 className="intel-title mt-2">{PIT_AGENT.name}</h1>
      <p className="intel-lede">
        PIT is an identity, not a wallet account on this page. The session key never enters the browser.
      </p>

      <dl className="intel-metrics mt-8">
        <div className="intel-metric">
          <dt>Agent</dt>
          <dd>{PIT_AGENT.name}</dd>
        </div>
        <div className="intel-metric">
          <dt>API agent</dt>
          <dd className="!text-[1rem]">{shortAddr(PIT_AGENT.address)}</dd>
          <p>{PIT_AGENT.address}</p>
        </div>
        <div className="intel-metric">
          <dt>Agentic ID / ERC-7857</dt>
          <dd>NOT LIVE ON MAINNET</dd>
          <p>iTransfer / iClone unavailable on Aristotle. No fake transfer controls.</p>
        </div>
        <div className="intel-metric">
          <dt>ERC-8004</dt>
          <dd>{reg?.ok ? "ownerOf returned" : "unverified"}</dd>
          <p>
            {reg?.ok
              ? `Live chain read owner ${shortAddr(reg.owner)}. Not a ranking.`
              : reg?.reason ?? "Reading Aristotle…"}
          </p>
        </div>
      </dl>

      <section className="intel-section">
        <h2 className="intel-kicker">Capabilities</h2>
        <ul className="intel-list mt-4">
          <li>Private research</li>
          <li>TEE compute</li>
          <li>Hyperliquid execution</li>
          <li>Bounded autonomy (desktop phrase)</li>
          <li>Memory (encrypted workspace)</li>
          <li>Proof publication when a public-safe receipt exists</li>
        </ul>
      </section>

      <section className="intel-section">
        <h2 className="intel-kicker">Restrictions</h2>
        <ul className="intel-list mt-4">
          <li>Withdraw</li>
          <li>Transfer of Agentic ID</li>
          <li>Policy mutation from the model or this website</li>
          <li>Unbounded leverage</li>
        </ul>
      </section>

      <section className="intel-section">
        <h2 className="intel-kicker">Verifiable reputation</h2>
        <p className="mt-3 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">
          No public ranking is claimed. Identity {shortAddr(IDENTITY_8004)} token {AGENT_8004_ID.toString()}.
          Reputation registry {shortAddr(REPUTATION_8004)}. This page does not invent verified case counts, PnL, or
          leaderboards.
        </p>
        <a
          className="intel-ghost mt-4 inline-flex"
          href={`${ARISTOTLE_EXPLORER}/address/${IDENTITY_8004}`}
          target="_blank"
          rel="noreferrer"
        >
          Open identity on ChainScan
        </a>
      </section>

      <Link to="/proof" className="intel-cta mt-8 inline-flex">
        Verify
      </Link>
    </div>
  );
}
