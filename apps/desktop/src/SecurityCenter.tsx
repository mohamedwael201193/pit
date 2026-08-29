import { hyperliquidAPI } from "./links";
import { ExternalLink } from "./ExternalLink";
import { BrandMark } from "./BrandMark";

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
    <section>
      <p>Each domain is READY, NEEDS ACTION, or BLOCKED. Color is never the only signal.</p>
      <table className="desk-table">
        <thead>
          <tr>
            <th>Domain</th>
            <th>State</th>
            <th>Why</th>
            <th>Do</th>
          </tr>
        </thead>
        <tbody>
          {domains.map((d) => (
            <tr key={d.id}>
              <td>
                <span className="asset">
                  {d.id === "hyperliquid" ? <BrandMark symbol="HL" /> : null}
                  {d.id}
                </span>
              </td>
              <td>{d.state || "NEEDS ACTION"}</td>
              <td>{d.why || "—"}</td>
              <td>
                <div className="cta-row" style={{ margin: 0 }}>
                  {d.id === "session" ? (
                    <button type="button" className="linkish" onClick={onSession} disabled={busy}>
                      Create session
                    </button>
                  ) : null}
                  {d.id === "policy" ? (
                    <button type="button" className="linkish" onClick={onPolicy} disabled={busy}>
                      Pin policy
                    </button>
                  ) : null}
                  {d.id === "hyperliquid" ? (
                    <ExternalLink className="linkish" href={hyperliquidAPI(net)}>
                      Open official page
                    </ExternalLink>
                  ) : null}
                  {d.href ? (
                    <ExternalLink className="linkish" href={d.href}>
                      {d.hrefLabel || "Open official page"}
                    </ExternalLink>
                  ) : null}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <button type="button" className="linkish" onClick={onCheck} disabled={busy}>
        Check again
      </button>
    </section>
  );
}
