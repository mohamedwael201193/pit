import { hyperliquidAPI } from "./links";

export type SecurityDomain = {
  id?: string;
  state?: string;
  why?: string;
  means?: string;
  do?: string;
  href?: string;
  hrefLabel?: string;
};

export function SecurityCenter({
  domains,
  net,
  onSession,
  onPolicy,
  onCheck,
  busy,
}: {
  domains: SecurityDomain[];
  net: string;
  onSession: () => void;
  onPolicy: () => void;
  onCheck: () => void;
  busy?: boolean;
}) {
  return (
    <article className="card">
      <p className="label">SECURITY CENTER</p>
      <p>Each domain is READY, NEEDS ACTION, or BLOCKED. Color is never the only signal.</p>
      <ul className="security-domains">
        {domains.map((d) => (
          <li key={d.id}>
            <strong>{d.state || "NEEDS ACTION"}</strong>
            <span>{d.id}</span>
            <p>
              <em>Why.</em> {d.why || "—"}
            </p>
            <p>
              <em>What it means.</em> {d.means || "—"}
            </p>
            <p>
              <em>What to do.</em> {d.do || "—"}
            </p>
            <div className="cta-row">
              {d.id === "session" ? (
                <button type="button" className="linkish" onClick={onSession} disabled={busy}>
                  Create / refresh secure session
                </button>
              ) : null}
              {d.id === "policy" ? (
                <button type="button" className="linkish" onClick={onPolicy} disabled={busy}>
                  Pin policy
                </button>
              ) : null}
              {d.id === "hyperliquid" ? (
                <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
                  Open Hyperliquid API
                </a>
              ) : null}
              {d.href ? (
                <a className="linkish" href={d.href} target="_blank" rel="noreferrer">
                  {d.hrefLabel || "Open"}
                </a>
              ) : null}
              <button type="button" className="linkish" onClick={onCheck} disabled={busy}>
                Check again
              </button>
            </div>
          </li>
        ))}
      </ul>
    </article>
  );
}
