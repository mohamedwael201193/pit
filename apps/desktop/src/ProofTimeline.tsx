import { useCallback, useEffect, useState } from "react";
import { fetchProofs, verifyProof, type FiledReceipt, type ProofsView, type VerifiedProof } from "./companion";

function shortHash(v?: string) {
  const s = String(v || "");
  if (s.length <= 18) return s;
  return `${s.slice(0, 10)}...${s.slice(-6)}`;
}

function kindLabel(kind?: string) {
  switch (String(kind || "")) {
    case "research":
      return "Research verdict";
    case "order":
      return "Venue order";
    case "cancel":
      return "Venue cancel";
    default:
      return "Evidence";
  }
}

function filedTime(v?: string) {
  if (!v) return "";
  const t = Date.parse(v);
  if (Number.isNaN(t)) return v;
  return new Date(t).toLocaleString();
}

function subject(row: FiledReceipt) {
  const bits: string[] = [];
  if (row.market) bits.push(row.market);
  if (row.side) bits.push(row.side);
  if (typeof row.size === "number" && row.size !== 0) bits.push(String(row.size));
  if (row.verdict) bits.push(row.verdict.replaceAll("_", " ").toLowerCase());
  if (row.oid) bits.push(`OID ${row.oid}`);
  return bits.join(" · ");
}

function blockLabel(hex?: string) {
  const s = String(hex || "");
  if (!s.startsWith("0x")) return s;
  const n = Number.parseInt(s, 16);
  return Number.isNaN(n) ? s : `block ${n.toLocaleString()}`;
}

function VerifyReport({ got }: { got: VerifiedProof }) {
  const rows: Array<[string, boolean, string]> = [
    ["Merkle proof checked by the official client", Boolean(got.proof_validated), "the bytes match the root"],
    ["Digest recomputed from downloaded bytes", Boolean(got.digest_match), shortHash(got.recomputed)],
    ["Record carries no key material", Boolean(got.public_safe), "safe to publish"],
    ["Sealed roles verified", Boolean(got.roles_verified), "envelope OK and signer matched"],
  ];
  if (got.tx) {
    rows.push([
      "0G Chain transaction commits this exact root",
      Boolean(got.anchor_bound),
      blockLabel(got.anchor?.block_number),
    ]);
  }
  const nodes = got.nodes || [];
  return (
    <div className="proof-report">
      <ul className="proof-checks">
        {rows.map(([label, pass, note]) => (
          <li key={label} className={pass ? "pass" : "fail"}>
            <span className="proof-check-mark">{pass ? "OK" : "NO"}</span>
            <span className="proof-check-label">{label}</span>
            <span className="proof-check-note hash">{note}</span>
          </li>
        ))}
      </ul>
      <p className="fine">
        Finalized on {got.finalized_nodes || 0} of {nodes.length} storage nodes the indexer named.
        {got.anchor?.flow ? ` The root is an indexed topic in a successful call to the storage flow contract ${got.anchor.flow}.` : ""}
      </p>
      {nodes.length ? (
        <ul className="proof-nodes">
          {nodes.map((n) => (
            <li key={n.node} className={n.finalized && !n.error ? "pass" : "fail"}>
              <span className="hash">{n.node}</span>
              <span>{n.finalized && !n.error ? "finalized" : n.error || "not finalized"}</span>
            </li>
          ))}
        </ul>
      ) : null}
      {got.failure ? <p className="proof-failure">Verification failed: {got.failure.replaceAll("_", " ")}.</p> : null}
      {!got.failure ? <p className="proof-pass">Independently verified against the live 0G network.</p> : null}
    </div>
  );
}

export function ProofTimeline() {
  const [view, setView] = useState<ProofsView | null>(null);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState<string>("");
  const [checking, setChecking] = useState<string>("");
  const [results, setResults] = useState<Record<string, VerifiedProof>>({});

  const load = useCallback(async () => {
    const got = await fetchProofs();
    setView(got);
    setLoading(false);
  }, []);

  useEffect(() => {
    void load();
    const t = window.setInterval(() => void load(), 15000);
    return () => window.clearInterval(t);
  }, [load]);

  const runVerify = useCallback(async (root: string) => {
    setChecking(root);
    setOpen(root);
    try {
      const got = await verifyProof(root);
      setResults((prev) => ({ ...prev, [root]: got }));
    } finally {
      setChecking("");
    }
  }, []);

  const rows = view?.receipts || [];

  if (loading) {
    return (
      <section className="proofs">
        <ul className="proof-skeleton">
          <li />
          <li />
          <li />
        </ul>
      </section>
    );
  }

  return (
    <section className="proofs">
      <div className="proof-head">
        <div>
          <p className="label">0G evidence</p>
          <h2>Published proof trail</h2>
          <p className="fine">
            Every research verdict and every posted order is written to 0G Storage in the clear and anchored by a 0G Chain
            transaction. Verify recomputes the digest from bytes the network hands back.
          </p>
        </div>
        {view?.network ? (
          <p className="proof-net hash">
            {view.network} · chain {view.chain_id}
          </p>
        ) : null}
      </div>

      {!view?.ready ? (
        <p className="proof-blocked">
          Filing is off: {String(view?.blocked || "unknown").replaceAll("_", " ")}. Bind the 0G payer account with
          <span className="hash"> pit evidence bind-payer </span>
          and the trail starts on the next verdict.
        </p>
      ) : null}

      {rows.length === 0 ? (
        <p className="empty">
          Nothing published yet. Run research and the verdict lands here with a storage root and a chain transaction.
        </p>
      ) : null}

      <ul className="proof-list">
        {rows.map((row) => {
          const root = String(row.root || "");
          const got = results[root];
          const isOpen = open === root;
          return (
            <li key={root + String(row.filed_at)}>
              <div className="proof-row">
                <div className="proof-when">
                  <time>{filedTime(row.filed_at)}</time>
                  <strong>{kindLabel(row.kind)}</strong>
                </div>
                <div className="proof-what">
                  <p>{subject(row) || "no subject"}</p>
                  <p className="hash proof-root" title={root}>
                    root {shortHash(root)}
                  </p>
                  <p className="hash proof-root" title={row.digest}>
                    digest {shortHash(row.digest)}
                  </p>
                </div>
                <div className="proof-act">
                  {row.tx_link ? (
                    <a href={row.tx_link} target="_blank" rel="noreferrer noopener">
                      Chain transaction
                    </a>
                  ) : (
                    <span className="fine">{row.duplicate ? "bytes already stored" : "no transaction recorded"}</span>
                  )}
                  <button
                    type="button"
                    className="proof-verify"
                    onClick={() => void runVerify(root)}
                    disabled={checking === root}
                  >
                    {checking === root ? "Verifying" : got ? "Verify again" : "Verify on 0G"}
                  </button>
                </div>
              </div>
              {isOpen && checking === root ? (
                <p className="fine proof-progress">
                  Downloading the record through the official client with merkle proof checking, then asking storage nodes
                  whether the root is finalized.
                </p>
              ) : null}
              {isOpen && got && checking !== root ? <VerifyReport got={got} /> : null}
            </li>
          );
        })}
      </ul>
    </section>
  );
}
