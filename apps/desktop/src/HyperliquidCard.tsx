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
  account,
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
  account?: string;
  onCreateSession: () => void;
  onConnectionPreview?: () => void;
  onCheck?: () => void;
  onRevoke?: () => void;
  onRefreshApproval?: () => void;
}) {
  const ttl =
    sessionExpires && sessionExpires > 0 ? new Date(sessionExpires).toISOString().replace(".000Z", "Z") : "";
  const canOrder = Boolean(sessionAlive && approved);
  return (
    <section className="hl-card">
      <p className="label">
        <span className="asset">
          <BrandMark symbol="HL" />
          Hyperliquid
        </span>
      </p>
      <div className="sec-metrics">
        <div>
          <span>Capital</span>
          <strong>{tradingCapital || "—"}</strong>
        </div>
        <div>
          <span>Session</span>
          <strong>{sessionAlive ? "live" : "none"}</strong>
        </div>
        <div>
          <span>Order / cancel</span>
          <strong>{canOrder ? "yes" : "no"}</strong>
        </div>
        <div>
          <span>Withdraw</span>
          <strong>no</strong>
        </div>
      </div>
      <p className="sec-meta">
        {account ? `${account.slice(0, 6)}…${account.slice(-4)}` : "no wallet"}
        {agentName ? ` · ${agentName}` : ""}
        {agent ? ` · ${agent.slice(0, 6)}…${agent.slice(-4)}` : ""}
        {ttl ? ` · until ${ttl}` : ""}
      </p>
      <div className="cta-row">
        {!sessionAlive ? (
          <button type="button" className="primary" onClick={onCreateSession} disabled={busy}>
            Create session
          </button>
        ) : !approved ? (
          <ExternalLink className="primary" href={hyperliquidAPI(net)}>
            Approve PIT
          </ExternalLink>
        ) : (
          <ExternalLink className="primary" href={hyperliquidApp(net)}>
            Open Hyperliquid
          </ExternalLink>
        )}
        <ExternalLink className="linkish" href={hyperliquidAPI(net)}>
          Open Hyperliquid API
        </ExternalLink>
        <button type="button" className="linkish" onClick={onRefreshApproval || onCheck} disabled={busy}>
          Refresh
        </button>
        {onRevoke && agent ? (
          <button type="button" className="linkish" onClick={onRevoke} disabled={busy}>
            Revoke
          </button>
        ) : null}
      </div>
      {approvedDetail ? <p className="fine">{approvedDetail}</p> : null}
      <p className="fine">
        Compute money lives at <ExternalLink href={LINKS.pcAdvanced}>0G Private Compute</ExternalLink>. That is not
        Hyperliquid.
      </p>
    </section>
  );
}
