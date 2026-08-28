import { LINKS, hyperliquidAPI } from "./links";

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
}: {
  net: string;
  agent: string;
  agentName?: string;
  sessionAlive: boolean;
  sessionExpires?: number;
  approved: boolean;
  approvedDetail?: string;
  busy?: boolean;
  onCreateSession: () => void;
  onConnectionPreview?: () => void;
}) {
  const approval = approved ? "Approved" : agent ? "Waiting" : "None";
  const session = sessionAlive ? "Active" : agent ? "Expired" : "Not created";
  const next = !agent
    ? "Create a secure PIT session on this computer."
    : !approved
      ? "Open Hyperliquid API. Authorize API Wallet with the name and address below. PIT still cannot withdraw."
      : !sessionAlive
        ? "Create a secure PIT session. extraAgents already lists this agent, so PIT will reuse it."
        : "Session is live. Research, then type AUTHORIZE on the exact preview.";
  const ttl =
    sessionExpires && sessionExpires > 0
      ? new Date(sessionExpires).toISOString().replace(".000Z", "Z")
      : "";
  return (
    <article className="card">
      <p className="label">CONNECT HYPERLIQUID</p>
      <p>
        PIT needs your Hyperliquid account, a scoped PIT API agent, and a local session. Withdraw, transfer, leverage,
        and account admin stay denied.
      </p>
      <p>
        <strong>PIT AGENT</strong> {agent || "none"}
      </p>
      {agentName ? <p>Name {agentName} (must be under 17 characters on Hyperliquid)</p> : null}
      <p>
        <strong>APPROVAL</strong> {approval}
        {approvedDetail ? ` — ${approvedDetail}` : ""}
      </p>
      <p>
        <strong>SESSION</strong> {session}
        {ttl ? ` — expires ${ttl}` : ""}
      </p>
      <p>
        ORDER {sessionAlive && approved ? "✅" : "waiting"} · CANCEL {sessionAlive && approved ? "✅" : "waiting"} ·
        WITHDRAW ❌ · TRANSFER ❌ · LEVERAGE ❌
      </p>
      <p className="fine">{next}</p>
      <div className="cta-row">
        <button type="button" className="linkish" onClick={onCreateSession} disabled={busy}>
          {sessionAlive ? "Refresh secure session" : "Create secure PIT session"}
        </button>
        <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
          Open Hyperliquid API
        </a>
        {onConnectionPreview && sessionAlive && approved ? (
          <button type="button" className="linkish" onClick={onConnectionPreview} disabled={busy}>
            Prepare connection-test preview
          </button>
        ) : null}
      </div>
      <p className="fine">
        Private compute credit lives at{" "}
        <a href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
          pc.0g.ai Advanced funds
        </a>
        . That is not a Hyperliquid balance.
      </p>
    </article>
  );
}
