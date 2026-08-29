import { prettyCode } from "./companion";
import { LINKS } from "./links";
import { ExternalLink } from "./ExternalLink";

export function PairingDock({
  code,
  expires,
  companionUp,
  paired,
  devices,
  busy,
  onRotate,
}: {
  code: string;
  expires?: string;
  companionUp: boolean;
  paired?: boolean;
  devices?: number;
  busy?: boolean;
  onRotate?: () => void;
}) {
  const display = code ? prettyCode(code) : companionUp ? "rotating…" : "waiting for local PIT";
  const left = expires ? remaining(expires) : "";
  return (
    <section className="pairing-dock" aria-label="Browser pairing">
      <div>
        <p className="label">Pair this browser</p>
        <p className="pair-chip" aria-label="pairing code">
          {display}
        </p>
        <p className="fine" style={{ margin: 0 }}>
          {left ? `Expires ${left}. ` : ""}One-time. The website never receives a session key.
        </p>
        <p className="fine" style={{ margin: 0 }}>
          Desktop {companionUp ? "connected" : "offline"} · Browser {paired ? `connected (${devices || 1})` : "unpaired"}
        </p>
      </div>
      <div className="cta-row">
        {onRotate ? (
          <button type="button" className="linkish" onClick={onRotate} disabled={busy || !companionUp}>
            Regenerate
          </button>
        ) : null}
        <ExternalLink className="primary" href={LINKS.pair}>
          Open pairing
        </ExternalLink>
      </div>
    </section>
  );
}

function remaining(expires: string) {
  const t = Date.parse(expires);
  if (Number.isNaN(t)) return expires;
  const sec = Math.max(0, Math.round((t - Date.now()) / 1000));
  if (sec <= 0) return "now — regenerate";
  return `${sec}s`;
}
