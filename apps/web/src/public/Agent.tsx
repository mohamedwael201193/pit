import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import {
  AGENT_8004_ID,
  ARISTOTLE_EXPLORER,
  DESK_ID_CONTRACT,
  DESK_TOKEN_ID,
  HL_API,
  IDENTITY_8004,
  PIT_AGENT,
  REPUTATION_8004,
} from "./facts";
import { PageHead } from "../ui/PageHead";
import { shortAddr } from "./format";
import { read8004Owner, readDesk, type AgentRead, type DeskRead } from "./chain";

export function AgentPage() {
  const [reg, setReg] = useState<AgentRead | null>(null);
  const [desk, setDesk] = useState<DeskRead | null>(null);

  useEffect(() => {
    void read8004Owner(AGENT_8004_ID).then(setReg);
    void readDesk(DESK_TOKEN_ID).then(setDesk);
  }, []);

  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title={PIT_AGENT.name}
        lede="Wallet is the human. Desk ID is the agent. The Hyperliquid API wallet is a permission, not identity. The session key never enters the browser. Mint and authorizeUsage stay on desktop."
      />

      <dl className="intel-metrics mt-8">
        <div className="intel-metric">
          <dt>Desk ID · ERC-7857</dt>
          <dd>{desk?.ok ? `token ${DESK_TOKEN_ID.toString()}` : "unverified"}</dd>
          <p>
            {desk?.ok
              ? `ownerOf ${shortAddr(desk.owner)}. isAuthorized(owner)=${desk.ownerAuthorized ? "true" : "false"}. This browser read Aristotle RPC.`
              : desk?.reason ?? "Reading Aristotle…"}
          </p>
        </div>
        <div className="intel-metric">
          <dt>iTransfer / iClone</dt>
          <dd>NOT LIVE ON MAINNET</dd>
          <p>Attestor is not on Aristotle. No transfer controls on this page.</p>
        </div>
        <div className="intel-metric">
          <dt>ERC-8004</dt>
          <dd>{reg?.ok ? "ownerOf returned" : "unverified"}</dd>
          <p>
            {reg?.ok
              ? `Live chain read owner ${shortAddr(reg.owner)}. Not a ranking. IDs are not portable 16661↔16602.`
              : reg?.reason ?? "Reading Aristotle…"}
          </p>
        </div>
        <div className="intel-metric">
          <dt>HL session agent</dt>
          <dd className="!text-[1rem]">{PIT_AGENT.name}</dd>
          <p>
            {shortAddr(PIT_AGENT.address)} · order/cancel only. Withdraw denied. Compare on Hyperliquid API. This site
            does not fetch another account’s extraAgents.
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
          Reputation registry {shortAddr(REPUTATION_8004)}. Owner self-feedback reverts. This page does not invent
          verified case counts, PnL, or leaderboards.
        </p>
        <div className="mt-4 flex flex-wrap gap-2">
          <a
            className="intel-ghost"
            href={`${ARISTOTLE_EXPLORER}/address/${DESK_ID_CONTRACT}`}
            target="_blank"
            rel="noreferrer"
          >
            Open Desk ID on ChainScan
          </a>
          <a
            className="intel-ghost"
            href={`${ARISTOTLE_EXPLORER}/address/${IDENTITY_8004}`}
            target="_blank"
            rel="noreferrer"
          >
            Open identity on ChainScan
          </a>
          <a className="intel-ghost" href={HL_API} target="_blank" rel="noreferrer">
            Open Hyperliquid API
          </a>
        </div>
      </section>

      <Link to="/proof" className="intel-cta mt-8 inline-flex">
        Verify
      </Link>
    </div>
  );
}
