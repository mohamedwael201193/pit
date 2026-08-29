import { LINKS, hyperliquidAPI, hyperliquidApp } from "./links";
import { ExternalLink } from "./ExternalLink";
import { BrandMark } from "./BrandMark";

export function HyperliquidCard({
  net,
  agent,
  agentName,
  sessionAlive,
  sessionExpires,
  approved,
  approvedDetail,
  busy,
  tradingCapital,
  onCreateSession,
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
  tradingCapital?: string;
  onCreateSession: () => void;
  onConnectionPreview?: () => void;
  onCheck?: () => void;
  onRevoke?: () => void;
  onRefreshApproval?: () => void;
}) {
  const session = sessionAlive ? "Active" : agent ? "Needs session" : "Not created";
  const ttl =
    sessionExpires && sessionExpires > 0 ? new Date(sessionExpires).toISOString().replace(".000Z", "Z") : "";
  return (
    <section>
      <p className="label">
        <span className="asset">
          <BrandMark symbol="HL" />
          Hyperliquid
        </span>
      </p>
      <dl className="status-grid hl-grid">
        <dt>Connected account</dt>
        <dd>{agent ? "yes" : "no"}</dd>
        <dt>Trading capital</dt>
        <dd>{tradingCapital || "—"}</dd>
        <dt>PIT Agent</dt>
        <dd>
          {agentName || "none"}
          {agent ? ` · ${agent.slice(0, 6)}…${agent.slice(-4)}` : ""}
        </dd>
        <dt>Session</dt>
        <dd>
          {session}
          {ttl ? ` · until ${ttl}` : ""}
        </dd>
        <dt>Expiration</dt>
        <dd>{ttl || "—"}</dd>
        <dt>Order</dt>
        <dd>{sessionAlive && approved ? "yes" : "no"}</dd>
        <dt>Cancel</dt>
        <dd>{sessionAlive && approved ? "yes" : "no"}</dd>
        <dt>Withdraw</dt>
        <dd>no</dd>
        <dt>Approval</dt>
        <dd>{approved ? "Approved" : "Needs approval"}</dd>
      </dl>
      {approvedDetail ? <p className="fine">{approvedDetail}</p> : null}
      <p className="fine">Account is your wallet. PIT Agent can order and cancel only. Trading capital is not private compute.</p>
      <div className="cta-row">
        <ExternalLink className="primary" href={hyperliquidApp(net)}>
          Connect Hyperliquid
        </ExternalLink>
        <ExternalLink className="linkish" href={hyperliquidAPI(net)}>
          Open Hyperliquid API
        </ExternalLink>
        <ExternalLink className="linkish" href={hyperliquidAPI(net)}>
          Approve PIT
        </ExternalLink>
        {!agent || !sessionAlive ? (
          <button type="button" className="linkish" onClick={onCreateSession} disabled={busy}>
            Create secure session
          </button>
        ) : null}
        <button type="button" className="linkish" onClick={onRefreshApproval || onCheck} disabled={busy}>
          Refresh status
        </button>
        {onRevoke && agent ? (
          <button type="button" className="linkish" onClick={onRevoke} disabled={busy}>
            Revoke
          </button>
        ) : null}
      </div>
      <p className="fine">
        Compute money lives at{" "}
        <ExternalLink href={LINKS.pcAdvanced}>0G Private Compute</ExternalLink>
        . That is not Hyperliquid.
      </p>
    </section>
  );
}
