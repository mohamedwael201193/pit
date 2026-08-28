import type { ReactNode } from "react";
import { LINKS, hyperliquidAPI } from "./links";
import { NAMED } from "./namedStates";

export function AuthorizeGate({
  sessionAlive,
  agent,
  agentName,
  net,
  busy,
  onCreateSession,
  children,
}: {
  sessionAlive: boolean;
  agent: string;
  agentName?: string;
  net: string;
  busy?: boolean;
  onCreateSession: () => void;
  children?: ReactNode;
}) {
  return (
    <div className="card auth">
      {!sessionAlive ? (
        <>
          <p className="label">YOUR SESSION</p>
          <p>
            Create a secure order/cancel session on this computer, then approve that agent on Hyperliquid API. If
            extraAgents already lists this address, PIT reuses it. PIT cannot withdraw.
          </p>
          <button type="button" className="linkish" onClick={onCreateSession} disabled={busy}>
            Create secure PIT session
          </button>
          {agent ? <p className="fine">Agent {agent}</p> : null}
          {agentName ? <p className="fine">API wallet name {agentName}. Must be under 17 characters on Hyperliquid.</p> : null}
          <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
            Open Hyperliquid API
          </a>
          <p className="fine">{NAMED.SESSION_EXPIRED}</p>
        </>
      ) : (
        <>
          <p className="label">YOUR SESSION</p>
          <p>Order and cancel only. Type AUTHORIZE on the exact preview on Research. This page does not place an order.</p>
          {agent ? <p className="fine">Agent {agent}</p> : null}
          {agentName ? <p className="fine">API wallet name {agentName}. Must be under 17 characters.</p> : null}
          <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
            Open Hyperliquid API
          </a>
          <p className="fine">{NAMED.AUTHORIZE_EXACT}</p>
        </>
      )}
      {children}
      <p className="fine">
        Direct credit lives at{" "}
        <a href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
          pc.0g.ai Advanced funds
        </a>
        . That page is provider credit, not a Hyperliquid balance.
      </p>
    </div>
  );
}
