import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { oidBelongsToPreview } from "../src/honesty";
import {
  collectJobReceipts,
  evidenceObjectForJob,
  receiptBelongsToJob,
  venueOrderState,
} from "../src/jobProof";
import type { ActivityEvent } from "../src/companion";
import { assertJobProofCorrelation } from "./job-proof.spec.ts";

export function assertLiveAgentPipeline() {
  assertJobProofCorrelation();

  const jobA = "job-a";
  const jobB = "job-b";
  const a: ActivityEvent = {
    job_id: jobA,
    root: "0xrootA",
    tx: "0xtxaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    tx_link: "https://chainscan.0g.ai/tx/0xtxaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    digest: "0xdigesta",
    market: "AVAX",
    ts: 100,
  };
  const b: ActivityEvent = {
    job_id: jobB,
    root: "0xrootB",
    tx: "0xtxbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
    tx_link: "https://chainscan.0g.ai/tx/0xtxbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
    digest: "0xdigestb",
    market: "HYPE",
    ts: 200,
  };
  const rowsA = collectJobReceipts([a, b], jobA, "mainnet");
  const rowsB = collectJobReceipts([a, b], jobB, "mainnet");
  if (rowsA.some((r) => (r.href || "").includes("bbbb") || r.jobId !== jobA)) {
    throw new Error("two concurrent jobs cannot cross-contaminate receipts");
  }
  if (rowsB.some((r) => (r.href || "").includes("aaaa") || r.jobId !== jobB)) {
    throw new Error("two concurrent jobs cannot cross-contaminate receipts");
  }
  if (!rowsA.length || !rowsB.length) throw new Error("each job must keep its own receipt");
  if (collectJobReceipts([a, b], "", "mainnet").length) throw new Error("empty job must wait, not fallback");
  if (evidenceObjectForJob({ job_id: jobB, tx: "0xtxbbbb" }, jobA)) {
    throw new Error("stale historical receipt cannot render");
  }
  if (evidenceObjectForJob({ tx: "0xlatest" }, jobA)) {
    throw new Error("latest receipt without jobId is not a fallback");
  }
  if (!receiptBelongsToJob(jobA, jobA)) throw new Error("receipt.jobId === currentAgentRun.jobId");

  if (oidBelongsToPreview("0xold-preview", "0xnew-preview")) throw new Error("old OID cannot render in current preview");
  if (venueOrderState({ oid: "529167222216", posted: true, status: "resting" }) === "filled") {
    throw new Error("RESTING cannot become FILLED");
  }

  const here = dirname(fileURLToPath(import.meta.url));
  const app = readFileSync(join(here, "../src/App.tsx"), "utf8");
  const run = readFileSync(join(here, "../src/AgentRun.tsx"), "utf8");
  const chat = readFileSync(join(here, "../src/CommandChat.tsx"), "utf8");
  if (!app.includes('authorizePreview("AUTHORIZE"')) throw new Error("TRADE NOW must use existing desktop authorize path");
  if (run.includes("authorizePreview") || chat.includes("authorizePreview")) {
    throw new Error("no second signing system");
  }
  if (!run.includes("ready && preview")) throw new Error("no READY result means no TRADE NOW button");
  if (!run.includes('kind === "READY_ELIGIBLE"')) throw new Error("TRADE NOW gated on READY_ELIGIBLE");
  if (!run.includes("cancelled") || !run.includes("READY_STOOD_DOWN")) {
    throw new Error("cancelled job cannot become NO TRADE");
  }
  if (!app.includes("openResearchStream")) throw new Error("live research stream");
  if (!app.includes("startResearch")) throw new Error("Find the best must start a real job");
  if (!chat.includes("Find best opportunity") && !chat.includes("Find the best opportunity")) {
    throw new Error("user can type Find the best opportunity");
  }
}
