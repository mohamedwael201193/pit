import { LINKS } from "./links";
import { ExternalLink } from "./ExternalLink";
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
  const tee = checkNamed(checks, "tee");
  const storage = checkNamed(checks, "storage");
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
        Direct <strong>{protectedOk ? (ready ? "ready" : "protected") : "needs Protect"}</strong>
      </span>
      <span>
        Balance <strong>{credit?.detail || "unread"}</strong>
      </span>
      <span>
        Next research <strong>~3 0G</strong>
      </span>
      <span>
        Funding <strong>{ready ? "funded" : "needs funds"}</strong>
      </span>
      <span>
        Provider <strong>Direct TeeML</strong>
      </span>
      <span>
        TEE <strong>{tee?.ok ? "verified roles" : tee?.detail || "unproven this session"}</strong>
      </span>
      <span>
        Storage <strong>{storage?.ok ? storage.detail : storage?.detail || "no proof until a root exists"}</strong>
      </span>
      <span className="pile">Available for research · not trading capital</span>
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
        <ExternalLink className={protectedOk ? "linkish" : "primary"} href={LINKS.app}>
          Protect my strategy
        </ExternalLink>
        {protectedOk && !ready && !sponsor && !(credit?.detail || "").toLowerCase().includes("unread") ? (
          <ExternalLink className="linkish" href={LINKS.pcAdvanced}>
            Open 0G Direct funds
          </ExternalLink>
        ) : null}
        <button type="button" className="linkish" onClick={onCheck}>
          Check again
        </button>
      </div>
    </section>
  );
}
