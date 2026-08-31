import { explorerTx } from "./links";
import type { ActivityEvent } from "./companion";

export type JobProofRow = {
  label: string;
  text: string;
  href?: string;
  jobId: string;
  market?: string;
  ts?: number;
};

export function receiptBelongsToJob(jobId: string, got?: string | null): boolean {
  const want = String(jobId || "").trim();
  const have = String(got || "").trim();
  return want !== "" && have === want;
}

export function shortProof(v?: string) {
  const s = String(v || "");
  if (s.length < 12) return s || "";
  return `${s.slice(0, 8)}…${s.slice(-4)}`;
}

export function collectJobReceipts(activity: ActivityEvent[], jobId: string, venue: string): JobProofRow[] {
  const want = String(jobId || "").trim();
  if (!want) return [];
  const rows: JobProofRow[] = [];
  const seen = new Set<string>();
  function add(label: string, text: string, href: string | undefined, e: ActivityEvent) {
    const key = `${label}:${text}:${href || ""}`;
    if (!text || seen.has(key)) return;
    seen.add(key);
    rows.push({
      label,
      text,
      href,
      jobId: want,
      market: e.market,
      ts: typeof e.ts === "number" ? e.ts : undefined,
    });
  }
  for (const e of activity) {
    if (!receiptBelongsToJob(want, e.job_id)) continue;
    const og = Boolean(e.root || e.tx || e.tx_link || e.digest);
    if (!og) continue;
    if (e.root) add("0G storage root", shortProof(String(e.root)), e.tx_link || undefined, e);
    if (e.tx) add("0G transaction", shortProof(String(e.tx)), e.tx_link || explorerTx(String(e.tx), venue) || undefined, e);
    if (e.tx_link && !e.tx) add("0G explorer", "Open", e.tx_link, e);
    if (e.digest) add("0G digest", shortProof(String(e.digest)), e.tx_link || undefined, e);
  }
  return rows.slice(0, 8);
}

export function evidenceObjectForJob(ev: Record<string, unknown> | null, jobId: string): ActivityEvent | null {
  if (!ev) return null;
  const evJob = String(ev.job_id || ev.jobId || ev.id || "");
  if (!receiptBelongsToJob(jobId, evJob)) return null;
  const root = String(ev.root || ev.storage_root || "");
  const tx = String(ev.tx || "");
  const txLink = String(ev.tx_link || "");
  const digest = String(ev.digest || "");
  if (!root && !tx && !txLink && !digest) return null;
  return {
    job_id: evJob,
    market: String(ev.market || ev.coin || ""),
    root: root || undefined,
    tx: tx || undefined,
    tx_link: txLink || undefined,
    digest: digest || undefined,
    ts: typeof ev.ts === "number" ? Number(ev.ts) : typeof ev.filed_at === "number" ? Number(ev.filed_at) : undefined,
    kind: "evidence.filed",
  };
}

export function venueOrderState(order?: {
  oid?: string;
  lifecycle?: string;
  status?: string;
  cancelled?: boolean;
  posted?: boolean;
} | null): string {
  if (!order?.oid) return "";
  const status = String(order.status || "").toLowerCase();
  const life = String(order.lifecycle || "").toLowerCase();
  const blob = `${status} ${life}`;
  if (order.cancelled || blob.includes("cancel")) return "cancelled";
  if (status.includes("fill") || life.includes("fill")) return "filled";
  if (blob.includes("fail") || blob.includes("reject")) return "failed";
  if (status.includes("rest") || status.includes("open") || life.includes("rest") || life.includes("open")) return "resting";
  if (order.posted && !status.includes("fill")) return "resting";
  return String(order.status || "submitted");
}
