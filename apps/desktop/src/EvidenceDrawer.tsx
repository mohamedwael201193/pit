import type { BindResult } from "./companion";
import { committeeVerified } from "./honesty";

export function EvidenceDrawer({
  jobId,
  roles,
  preview,
  previewHash,
  kind,
  deny,
  evidence,
  open,
  onToggle,
}: {
  jobId?: string;
  roles: Array<{ role?: string; verify_e2ee?: string; pubkey_signer?: string; proposed_side?: string; survives?: boolean; kill?: boolean }>;
  preview?: BindResult["preview"] | null;
  previewHash?: string;
  kind?: string;
  deny?: string;
  evidence?: string;
  open?: boolean;
  onToggle?: () => void;
}) {
  const verified = committeeVerified(roles);
  return (
    <section style={{ marginTop: 14 }}>
      <p className="label">Evidence</p>
      <p className="fine">
        Job {jobId || "none"} · {kind || "idle"} · {verified ? "three roles verified" : "committee not complete"}. The sealed prompt stays private.
      </p>
      <button type="button" className="linkish" onClick={onToggle}>
        {open ? "Hide technical evidence" : "View technical evidence"}
      </button>
      {open ? (
        <dl className="evidence-list">
          <dt>Committee</dt>
          <dd>
            {roles.length
              ? roles
                  .map(
                    (r) =>
                      `${r.role}: ${r.verify_e2ee || "—"}${r.proposed_side ? ` ${r.proposed_side}` : ""}${r.kill ? " kill" : ""}${r.survives === false ? " stood down" : ""}`,
                  )
                  .join(" · ")
              : "No roles yet."}
          </dd>
          <dt>Policy</dt>
          <dd>{deny || (preview?.eligible ? "Allowed this clip." : "No preview.")}</dd>
          <dt>Preview hash</dt>
          <dd>{previewHash || preview?.hash || "none"}</dd>
          <pre className="pipe evidence">{evidence || "No local result file. Router fallback is impossible."}</pre>
        </dl>
      ) : null}
    </section>
  );
}
