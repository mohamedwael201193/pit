import { LINKS } from "./links";
import type { DoctorCheck } from "./companion";
import { checkNamed } from "./companion";

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
    <article className="card">
      <p className="label">PRIVATE COMPUTE</p>
      <p>
        Protected {protectedOk ? "yes" : "no"} · Direct {ready ? "Ready" : "Needs funds or Protect signature"} · Estimated ~3 0G
        locked for one sealed committee.
      </p>
      <p>{credit?.detail || auth?.detail || "Sign Protect my strategy, then fund 0G Private Compute if asked."}</p>
      {sponsor ? (
        <p className="fine">
          PIT is paying this sealed run from a shared compute account. This is not your 0G. Trading capital is never used.
        </p>
      ) : (
        <p className="fine">
          Compute money lives at 0G Private Compute. Hyperliquid trading capital never pays inference. Delayed settlement can drop
          the ledger in a lump — that is not theft.
        </p>
      )}
      <div className="cta-row">
        <a className="linkish" href={LINKS.app} target="_blank" rel="noreferrer">
          Protect my strategy
        </a>
        <a className="linkish" href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
          Open 0G Private Compute
        </a>
        <button type="button" className="linkish" onClick={onCheck}>
          Check again
        </button>
      </div>
    </article>
  );
}
