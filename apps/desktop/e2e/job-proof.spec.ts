import { committeeVerified, oidBelongsToPreview } from "../src/honesty";
import {
  collectJobReceipts,
  evidenceObjectForJob,
  receiptBelongsToJob,
  venueOrderState,
} from "../src/jobProof";
import type { ActivityEvent } from "../src/companion";

export function assertJobProofCorrelation() {
  const job = "job-live";
  const old: ActivityEvent = {
    job_id: "job-old",
    root: "0xoldroot",
    tx: "0xoldtx",
    tx_link: "https://chainscan.0g.ai/tx/0xoldtx",
    digest: "0xolddigest",
    market: "AVAX",
    ts: 1,
  };
  const live: ActivityEvent = {
    job_id: job,
    root: "0xnewroot",
    tx: "0xnewtxdeadbeef",
    tx_link: "https://chainscan.0g.ai/tx/0xnewtxdeadbeef",
    digest: "0xnewdigest",
    market: "HYPE",
    ts: 2,
  };
  const rows = collectJobReceipts([old, live], job, "mainnet");
  if (rows.some((r) => r.text.includes("old") || (r.href || "").includes("oldtx"))) {
    throw new Error("stale historical receipt cannot render");
  }
  if (!rows.some((r) => r.jobId === job && (r.href || "").includes("newtx"))) {
    throw new Error("current job receipt missing");
  }
  const other = collectJobReceipts([old, live], "job-other", "mainnet");
  if (other.length) throw new Error("two concurrent jobs cannot cross-contaminate receipts");
  if (collectJobReceipts([old, live], "", "mainnet").length) throw new Error("empty job must render nothing");
  if (receiptBelongsToJob(job, "")) throw new Error("blank receipt job is not a match");
  if (receiptBelongsToJob(job, job) !== true) throw new Error("exact job must match");

  const ev = evidenceObjectForJob({ job_id: "job-old", tx: "0xoldtx" }, job);
  if (ev) throw new Error("evidence blob from another job must drop");
  if (evidenceObjectForJob({ tx: "0xorphan" }, job)) throw new Error("evidence without jobId is not a fallback");

  if (oidBelongsToPreview("0xoldhash", "0xnewhash")) throw new Error("old OID cannot render in current preview");
  if (venueOrderState({ oid: "1", posted: true, status: "resting" }) === "filled") {
    throw new Error("RESTING cannot become FILLED");
  }
  if (venueOrderState({ oid: "1", status: "filled" }) !== "filled") throw new Error("fill must stay fill");
  if (venueOrderState({ oid: "1", posted: true, status: "open" }) !== "resting") throw new Error("posted without fill is resting");
  if (venueOrderState({ oid: "1", posted: true, lifecycle: "reconciled", status: "filled" }) !== "filled") {
    throw new Error("reconciled fill must stay FILLED");
  }

  const three = [
    { role: "researcher", verify_e2ee: "OK" },
    { role: "challenger", verify_e2ee: "OK" },
    { role: "risk", verify_e2ee: "OK" },
  ];
  if (!committeeVerified(three)) throw new Error("committee");
}
