import { LINKS } from "./links";
import type { DoctorCheck } from "./companion";
import { checkNamed } from "./companion";
import { BrandMark } from "./BrandMark";

export function ComputeCard({
  checks,
  onCheck,
}: {
  checks: DoctorCheck[];
  onCheck: () => void;
}) {
  const auth = checkNamed(checks, "direct_auth");
  const credit = checkNamed(checks, "direct_credit");
  const sponsor = Boolean(credit?.ok && credit.detail.toLowerCase().includes("sponsor"));
  const protectedOk = Boolean(auth?.ok);
  const ready = Boolean(credit?.ok);
  return (
    <section className="compute-strip">
      <span className="asset">
        <BrandMark symbol="0G" />
        <p className="label" style={{ margin: 0 }}>
          Private compute
        </p>
      </span>
      <span>
        Model <strong>glm-5.2</strong>
      </span>
      <span>
        Protected <strong>{protectedOk ? "yes" : "no"}</strong>
      </span>
      <span>
        Status <strong>{ready ? "Ready" : "Needs action"}</strong>
      </span>
      <span>{credit?.detail || "Sign Protect my strategy, then fund if asked."}</span>
      <span className="pile">Available for research · ~3 0G estimated per sealed committee · not trading capital</span>
      {sponsor ? <span className="fine">Sponsored compute. Not your 0G. Never trading capital.</span> : null}
      {!protectedOk ? (
        <p className="fine" style={{ margin: 0 }}>
          Missing: wallet signature for sealed research. Open the paired site, then Check again.
        </p>
      ) : !ready ? (
        <p className="fine" style={{ margin: 0 }}>
          Missing: Direct credit for this wallet at 0G Private Compute. That is compute money, not Hyperliquid.
        </p>
      ) : null}
      <div className="cta-row" style={{ marginTop: 0 }}>
        <a className={protectedOk ? "linkish" : "primary"} href={LINKS.app} target="_blank" rel="noreferrer">
          Protect my strategy
        </a>
        <a className="linkish" href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
          Open 0G Private Compute
        </a>
        <button type="button" className="linkish" onClick={onCheck}>
          Check again
        </button>
      </div>
    </section>
  );
}
