import { LINKS } from "./links";
import type { DoctorCheck } from "./companion";
import { checkNamed } from "./companion";

export function ComputeCard({
  checks,
  onCheck,
  onEvidence,
}: {
  checks: DoctorCheck[];
  onCheck: () => void;
  onEvidence?: () => void;
}) {
  const auth = checkNamed(checks, "direct_auth");
  const credit = checkNamed(checks, "direct_credit");
  const sponsor = Boolean(credit?.ok && credit.detail.toLowerCase().includes("sponsor"));
  const protectedOk = Boolean(auth?.ok);
  const ready = Boolean(credit?.ok);
  return (
    <section className="compute-strip">
      <p className="label" style={{ margin: 0 }}>
        Private compute
      </p>
      <span>
        Protected <strong>{protectedOk ? "yes" : "no"}</strong>
      </span>
      <span>
        Direct <strong>{ready ? "Ready" : "Needs action"}</strong>
      </span>
      <span>{credit?.detail || "Sign Protect my strategy, then fund if asked."}</span>
      <span className="pile">~3 0G locked per sealed committee · not trading capital</span>
      {sponsor ? <span className="fine">Sponsored compute. Not your 0G. Never trading capital.</span> : null}
      <div className="cta-row" style={{ marginTop: 0 }}>
        {!protectedOk ? (
          <a className="primary" href={LINKS.app} target="_blank" rel="noreferrer">
            Protect my strategy
          </a>
        ) : (
          <a className="linkish" href={LINKS.app} target="_blank" rel="noreferrer">
            Protect my strategy
          </a>
        )}
        <a className="linkish" href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
          Open 0G Private Compute
        </a>
        <button type="button" className="linkish" onClick={onCheck}>
          Check again
        </button>
        {onEvidence ? (
          <button type="button" className="linkish" onClick={onEvidence}>
            View technical evidence
          </button>
        ) : null}
      </div>
    </section>
  );
}
