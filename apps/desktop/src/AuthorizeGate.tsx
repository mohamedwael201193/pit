import { FormEvent, useState } from "react";
import { confirmAuthorize } from "./authorize";
import { NAMED } from "./namedStates";

export function AuthorizeGate({ sessionAlive }: { sessionAlive: boolean }) {
  const [typed, setTyped] = useState("");
  const [err, setErr] = useState<string | null>(null);

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    setErr(confirmAuthorize(typed, sessionAlive));
  }

  return (
    <form className="card auth" onSubmit={onSubmit}>
      <p className="label">YOUR PREVIEW</p>
      <p>{NAMED.AUTHORIZE_EXACT}</p>
      <input
        aria-label="authorize token"
        autoComplete="off"
        value={typed}
        onChange={(ev) => setTyped(ev.target.value)}
        placeholder="type the token"
      />
      <button type="submit" disabled={!sessionAlive}>
        Authorize
      </button>
      {!sessionAlive ? (
        <p className="err" role="alert">
          {NAMED.SESSION_EXPIRED}
        </p>
      ) : null}
      {err && sessionAlive ? (
        <p className="err" role="alert">
          {err === "session_expired" ? NAMED.SESSION_EXPIRED : NAMED.AUTHORIZE_REFUSED}
        </p>
      ) : null}
    </form>
  );
}
