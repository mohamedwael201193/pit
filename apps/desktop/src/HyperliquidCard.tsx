import { useState } from "react";
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
  const stage = !sessionAlive
    ? "ACTION REQUIRED"
    : approved
      ? "VERIFIED"
      : "WAITING FOR APPROVAL";
  const [copied, setCopied] = useState("");

  function copy(label: string, value: string) {
    if (!value) return;
    void navigator.clipboard.writeText(value).then(() => setCopied(label));
  }

  return (
    <section className="hl-card">
      <p className="label">
        <span className="asset">
          <BrandMark symbol="HL" />
          Hyperliquid
        </span>
        <span className={`onboard-state ${stage.replace(/\s+/g, "-").toLowerCase()}`}>{stage}</span>
      </p>
      <div className="hl-ids">
        <p>
          <span>Your wallet</span>
          <strong>{account || "not bound"}</strong>
        </p>
        <p>
          <span>PIT Agent</span>
          <strong>
            {agentName || "not created"}
            {agent ? ` · ${agent}` : ""}
          </strong>
        </p>
      </div>
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
      <div className="hl-can">
        <p>
          <strong>PIT can</strong> place orders and cancel them after Hyperliquid lists this agent.
        </p>
        <p>
          <strong>PIT cannot</strong> withdraw, transfer, change leverage, or run account admin.
        </p>
      </div>
      {ttl ? <p className="fine">Listed until {ttl}. PIT never uses the agent address as your account for info queries.</p> : null}
      <div className="cta-row">
        {!sessionAlive ? (
          <button type="button" className="primary" onClick={onCreateSession} disabled={busy}>
            Create PIT Agent on this computer
          </button>
        ) : !approved ? (
          <ExternalLink className="primary" href={hyperliquidAPI(net)}>
            Approve PIT on Hyperliquid
          </ExternalLink>
        ) : (
          <ExternalLink className="primary" href={hyperliquidApp(net)}>
            Open Hyperliquid
          </ExternalLink>
        )}
        {agent ? (
          <button type="button" className="linkish" onClick={() => copy("agent", agent)}>
            {copied === "agent" ? "Copied PIT Agent address" : "Copy PIT Agent address"}
          </button>
        ) : null}
        <button type="button" className="linkish" onClick={onRefreshApproval || onCheck} disabled={busy}>
          Check approval
        </button>
        {onRevoke && agent ? (
          <button type="button" className="linkish" onClick={onRevoke} disabled={busy}>
            Revoke local session
          </button>
        ) : null}
      </div>
      {approvedDetail ? <p className="fine">{approvedDetail}</p> : null}
      <p className="fine">
        Approval is verified from live Hyperliquid extraAgents on your master wallet. A button click is not proof.
        Compute money lives at <ExternalLink href={LINKS.pcAdvanced}>0G Private Compute</ExternalLink>.
      </p>
    </section>
  );
}
