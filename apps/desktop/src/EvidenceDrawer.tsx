import { useState } from "react";
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
}: {
  jobId?: string;
  roles: Array<{ role?: string; verify_e2ee?: string; pubkey_signer?: string; proposed_side?: string; survives?: boolean; kill?: boolean }>;
  preview?: BindResult["preview"] | null;
  previewHash?: string;
  kind?: string;
  deny?: string;
  evidence?: string;
}) {
  const [open, setOpen] = useState(false);
  const verified = committeeVerified(roles);
  return (
    <article className="card">
      <p className="label">EVIDENCE</p>
      <p>Job {jobId || "none"} · {kind || "idle"} · {verified ? "three roles verified" : "committee not complete"}</p>
      <p className="fine">The sealed prompt stays private. This drawer never shows it.</p>
      <dl className="evidence-list">
        <div>
          <dt>What PIT saw</dt>
          <dd>Public Hyperliquid mark, oracle, funding, and open interest. Not your private thesis.</dd>
        </div>
        <div>
          <dt>Committee</dt>
          <dd>
            {roles.length
              ? roles.map((r) => `${r.role}: ${r.verify_e2ee || "—"}${r.proposed_side ? ` ${r.proposed_side}` : ""}${r.kill ? " kill" : ""}${r.survives === false ? " stood down" : ""}`).join(" · ")
              : "No roles yet."}
          </dd>
        </div>
        <div>
          <dt>Host size</dt>
          <dd>Host sized to policy clip and venue minimum. Model size was ignored.</dd>
        </div>
        <div>
          <dt>Policy</dt>
          <dd>{deny || (preview?.eligible ? "Allowed this clip." : "No preview.")}</dd>
        </div>
        <div>
          <dt>Exact action</dt>
          <dd>
            {preview?.eligible
              ? `${preview.market} ${preview.side} ${preview.sz} @ ${preview.limitPx}`
              : "None. AUTHORIZE is not available."}
          </dd>
        </div>
        <div>
          <dt>Preview hash</dt>
          <dd>
            {previewHash || preview?.hash || "none"}. This hash binds the exact order PIT will send. If anything
            changes, the hash changes, and AUTHORIZE will not apply.
          </dd>
        </div>
      </dl>
      <button type="button" className="linkish" onClick={() => setOpen((v) => !v)}>
        {open ? "Hide provenance" : "Open evidence"}
      </button>
      {open ? (
        <pre className="pipe evidence">
          {evidence || "No local result file. Router fallback is impossible."}
        </pre>
      ) : null}
    </article>
  );
}
