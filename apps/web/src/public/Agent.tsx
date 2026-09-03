import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import {
  AGENT_8004_ID,
  AGENT_CARD_URL,
  ARISTOTLE_EXPLORER,
  DESK_ID_CONTRACT,
  DESK_TOKEN_ID,
  HL_API,
  IDENTITY_8004,
  PIT_AGENT,
  REPUTATION_8004,
  REPUTATION_CLIENT,
} from "./facts";
import { PageHead } from "../ui/PageHead";
import { shortAddr } from "./format";
import {
  encodeAuthorizeUsage,
  encodeMint,
  encodeRevokeAuthorization,
  encodeSetAgentURI,
  publicDataHash,
  read8004Owner,
  read8004Reputation,
  readDesk,
  sendWalletCall,
  type AgentRead,
  type DeskRead,
  type ReputationRead,
} from "./chain";

export function AgentPage() {
  const { ready, authenticated, login, user } = usePrivy();
  const addr = user?.wallet?.address || "";
  const [reg, setReg] = useState<AgentRead | null>(null);
  const [desk, setDesk] = useState<DeskRead | null>(null);
  const [rep, setRep] = useState<ReputationRead | null>(null);
  const [tx, setTx] = useState<string | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    void read8004Owner(AGENT_8004_ID).then(setReg);
    void readDesk(DESK_TOKEN_ID).then(setDesk);
    void read8004Reputation(AGENT_8004_ID, REPUTATION_CLIENT).then(setRep);
  }, []);

  const isOwner =
    Boolean(addr) &&
    ((desk?.ok && addr.toLowerCase() === desk.owner.toLowerCase()) ||
      (reg?.ok && addr.toLowerCase() === reg.owner.toLowerCase()));

  const run = async (label: string, to: `0x${string}`, data: `0x${string}`) => {
    setErr(null);
    setTx(null);
    if (!addr) {
      setErr("Connect the owner wallet. PIT never asks for a seed phrase.");
      return;
    }
    setBusy(true);
    try {
      const hash = await sendWalletCall(addr, to, data);
      setTx(`${label} ${hash}`);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Wallet declined the signature.");
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title={PIT_AGENT.name}
        lede="Wallet is the human. Desk ID is the agent. The Hyperliquid API wallet is a permission, not identity. The session key never enters the browser. Mint and authorizeUsage are signed by the connected owner wallet on Aristotle."
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
              ? `Live chain read owner ${shortAddr(reg.owner)}. Agent ${AGENT_8004_ID.toString()}.`
              : reg?.reason ?? "Reading Aristotle…"}
          </p>
        </div>
        <div className="intel-metric">
          <dt>Reputation</dt>
          <dd>{rep?.ok ? (rep.index === "0" ? "no client index" : `index ${rep.index}`) : "unverified"}</dd>
          <p>
            {rep?.ok
              ? `URI ${rep.uri || "—"}. ${rep.tag1 ? `${rep.tag1} / ${rep.tag2}` : "Waiting for on-chain feedback readback."}`
              : rep?.reason ?? "Reading Aristotle…"}
          </p>
        </div>
      </dl>

      <section className="intel-section">
        <h2 className="intel-kicker">Owner wallet</h2>
        <p className="mt-3 max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">
          These calls go to Aristotle from the connected wallet. They cannot place a Hyperliquid order. The session
          agent cannot sign them.
        </p>
        {!ready ? (
          <p className="mt-4">Loading wallet…</p>
        ) : !authenticated || !addr ? (
          <button type="button" className="intel-cta mt-4" onClick={() => void login()}>
            Connect owner wallet
          </button>
        ) : (
          <div className="mt-4 flex flex-wrap gap-2">
            {isOwner ? (
              <>
                <button
                  type="button"
                  className="intel-cta"
                  disabled={busy}
                  onClick={() =>
                    void run(
                      "authorizeUsage",
                      DESK_ID_CONTRACT,
                      encodeAuthorizeUsage(DESK_TOKEN_ID, PIT_AGENT.address as `0x${string}`),
                    )
                  }
                >
                  Authorize desk usage
                </button>
                <button
                  type="button"
                  className="intel-ghost"
                  disabled={busy}
                  onClick={() =>
                    void run(
                      "revokeAuthorization",
                      DESK_ID_CONTRACT,
                      encodeRevokeAuthorization(DESK_TOKEN_ID, PIT_AGENT.address as `0x${string}`),
                    )
                  }
                >
                  Revoke desk usage
                </button>
                <button
                  type="button"
                  className="intel-ghost"
                  disabled={busy}
                  onClick={() => void run("setAgentURI", IDENTITY_8004, encodeSetAgentURI(AGENT_8004_ID, AGENT_CARD_URL))}
                >
                  Publish agent URI
                </button>
              </>
            ) : (
              <button
                type="button"
                className="intel-cta"
                disabled={busy}
                onClick={() =>
                  void run(
                    "mint",
                    DESK_ID_CONTRACT,
                    encodeMint(addr as `0x${string}`, AGENT_CARD_URL, publicDataHash(AGENT_CARD_URL), "pit-desk"),
                  )
                }
              >
                Mint desk ID
              </button>
            )}
          </div>
        )}
        {tx ? (
          <p className="mt-3 break-all text-[0.875rem] text-[var(--guide-cream)]">
            Submitted {tx}. Confirm on ChainScan.
          </p>
        ) : null}
        {err ? (
          <p className="mt-3 text-[0.875rem] text-[#ff7a7a]" role="alert">
            {err}
          </p>
        ) : null}
      </section>

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
          Identity {shortAddr(IDENTITY_8004)} token {AGENT_8004_ID.toString()}. Reputation registry{" "}
          {shortAddr(REPUTATION_8004)}. Owner self-feedback reverts. This page reads on-chain feedback. It does not
          invent verified case counts, PnL, or leaderboards.
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
