import { LINKS, hyperliquidAPI, hyperliquidApp } from "./links";

export function HyperliquidCard({
  net,
  agent,
  agentName,
  sessionAlive,
  sessionExpires,
  approved,
  approvedDetail,
  busy,
  onCreateSession,
  onConnectionPreview,
  onCheck,
  onRevoke,
  onRefreshApproval,
}: {
  net: string;
  agent: string;
  agentName?: string;
  sessionAlive?: boolean;
  sessionExpires?: number;
  approved: boolean;
  approvedDetail?: string;
  busy?: boolean;
  onCreateSession: () => void;
  onConnectionPreview?: () => void;
  onCheck?: () => void;
  onRevoke?: () => void;
  onRefreshApproval?: () => void;
}) {
  const session = sessionAlive ? "Active" : agent ? "Expired" : "Not created";
  const status = approved && sessionAlive ? "Ready" : !agent ? "Not connected" : approved ? "Session expired" : "Needs approval";
  const next = !agent
    ? "Create a secure PIT session on this computer."
    : !approved
      ? "Open Hyperliquid API. Authorize API Wallet with the name and address below. PIT still cannot withdraw."
      : !sessionAlive
        ? "Create a secure PIT session. Hyperliquid already lists this agent, so PIT will reuse it."
        : "Your trading account is ready. Research, then type AUTHORIZE on the exact preview.";
  const ttl =
    sessionExpires && sessionExpires > 0 ? new Date(sessionExpires).toISOString().replace(".000Z", "Z") : "";
  return (
    <article className="card">
      <p className="label">HYPERLIQUID</p>
      <p>
        Connected {agent ? "yes" : "no"} · Account is your wallet. PIT Agent is a scoped API wallet. Withdraw, transfer,
        leverage, and account admin stay denied.
      </p>
      <p>
        <strong>PIT AGENT</strong> {agent || "none"} {agentName ? `· ${agentName}` : ""}
      </p>
      <p>
        <strong>SESSION</strong> {session}
        {ttl ? ` — until ${ttl}` : ""}
      </p>
      <p>
        <strong>PERMISSION</strong> Order · Cancel · Withdraw denied · Transfer denied · Leverage denied
      </p>
      <p>
        <strong>STATUS</strong> {status}
        {approvedDetail ? ` — ${approvedDetail}` : ""}
      </p>
      <dl className="status-grid" style={{ marginTop: 8 }}>
        <dt>Order</dt>
        <dd>{sessionAlive && approved ? "yes" : "waiting"}</dd>
        <dt>Cancel</dt>
        <dd>{sessionAlive && approved ? "yes" : "waiting"}</dd>
        <dt>Withdraw</dt>
        <dd>no</dd>
        <dt>Transfer</dt>
        <dd>no</dd>
        <dt>Leverage</dt>
        <dd>no</dd>
      </dl>
      <p className="fine">{next}</p>
      <div className="cta-row">
        <a className="linkish" href={hyperliquidApp(net)} target="_blank" rel="noreferrer">
          Connect Hyperliquid
        </a>
        <button type="button" className="linkish" onClick={onCreateSession} disabled={busy}>
          {sessionAlive ? "Refresh secure session" : "Create secure session"}
        </button>
        <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
          Approve PIT
        </a>
        <button
          type="button"
          className="linkish"
          onClick={onRefreshApproval || onCheck}
          disabled={busy || (!onRefreshApproval && !onCheck)}
        >
          Refresh approval
        </button>
        <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
          Open Hyperliquid API
        </a>
        {onCheck ? (
          <button type="button" className="linkish" onClick={onCheck} disabled={busy}>
            Check again
          </button>
        ) : null}
        {onRevoke ? (
          <button type="button" className="linkish" onClick={onRevoke} disabled={busy || !agent}>
            Revoke PIT
          </button>
        ) : null}
        {onConnectionPreview && sessionAlive && approved ? (
          <button type="button" className="linkish" onClick={onConnectionPreview} disabled={busy}>
            Prepare connection-test preview
          </button>
        ) : null}
      </div>
      <p className="fine">
        Private compute credit lives at{" "}
        <a href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
          0G Private Compute
        </a>
        . That is not a Hyperliquid balance. Revoke PIT deletes the local session, then you remove the agent on Hyperliquid
        API.
      </p>
    </article>
  );
}
