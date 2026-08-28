import { FormEvent, useState } from "react";
import { confirmAuthorize } from "./authorize";
import { LINKS, hyperliquidAPI } from "./links";
import { NAMED } from "./namedStates";

export function AuthorizeGate({
  sessionAlive,
  agent,
  net,
  busy,
  onCreateSession,
}: {
  sessionAlive: boolean;
  agent: string;
  net: string;
  busy?: boolean;
  onCreateSession: () => void;
}) {
  const [typed, setTyped] = useState("");
  const [err, setErr] = useState<string | null>(null);

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    setErr(confirmAuthorize(typed, sessionAlive));
  }

  return (
    <div className="card auth">
      {!sessionAlive ? (
        <>
          <p className="label">YOUR SESSION</p>
          <p>Create an order/cancel session on this computer, then approve that agent on Hyperliquid. PIT cannot withdraw.</p>
          <button type="button" className="linkish" onClick={onCreateSession} disabled={busy}>
            Create local session
          </button>
          {agent ? <p className="fine">Agent {agent}</p> : null}
          <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
            Open Hyperliquid
          </a>
          <p className="fine">{NAMED.SESSION_EXPIRED}</p>
        </>
      ) : (
        <form onSubmit={onSubmit}>
          <p className="label">YOUR PREVIEW</p>
          <p>{NAMED.AUTHORIZE_EXACT}</p>
          <input
            aria-label="type AUTHORIZE"
            autoComplete="off"
            value={typed}
            onChange={(ev) => setTyped(ev.target.value)}
            placeholder="Type AUTHORIZE"
          />
          <button type="submit">Authorize</button>
          {err ? (
            <p className="err" role="alert">
              {err === "session_expired" ? NAMED.SESSION_EXPIRED : NAMED.AUTHORIZE_REFUSED}
            </p>
          ) : null}
        </form>
      )}
      <p className="fine">
        Direct credit lives at{" "}
        <a href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
          pc.0g.ai
        </a>
        . Switch to Advanced. That page is provider credit, not a Hyperliquid balance.
      </p>
    </div>
  );
}
