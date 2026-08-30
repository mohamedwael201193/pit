import { CAPABILITY } from "../src/CapabilityMatrix";
import { explainCommittee } from "../src/committee";
import { deskHeadline } from "../src/deskCopy";
import { nearestPolicyClip, nearestVenueMin } from "../src/format";
import { researchCardTitle } from "../src/honesty";

export function assertCapitalFloorIgnoresDust() {
  const coins = [
    { coin: "MEW", minNotional: 10, policyClip: 10, eligible: false, policyEligible: false, executionFeasible: false },
    { coin: "DOGE", minNotional: 10.02, policyClip: 10, eligible: true, policyEligible: true, executionFeasible: false },
    { coin: "ETH", minNotional: 10.08, policyClip: 10, eligible: true, policyEligible: true, executionFeasible: false },
  ];
  if (nearestVenueMin(coins) !== 10.02) throw new Error("policy-eligible min must ignore MEW");
  if (nearestPolicyClip(coins) !== 10) throw new Error("policy clip");
  const exec = [{ minNotional: 10.52, eligible: true, policyEligible: true, executionFeasible: true }, ...coins];
  if (nearestVenueMin(exec) !== 10.52) throw new Error("executable min wins");
}

export function assertDeskRanksExecutableFirst() {
  const coins = [
    { coin: "ETH", eligible: true, executionFeasible: false, previewReady: false },
    { coin: "BTC", eligible: true, executionFeasible: true, previewReady: false },
  ];
  const ranked = [...coins]
    .filter((c) => c.executionFeasible || c.previewReady || c.eligible)
    .sort((a, b) => Number(Boolean(b.executionFeasible)) - Number(Boolean(a.executionFeasible)))
    .map((c) => c.coin);
  if (ranked[0] !== "BTC") throw new Error("executable first");
}

export function assertPolicyClipTightCopy() {
  const clip = explainCommittee("policy_clip_tight");
  if (!/too tight/i.test(clip.title)) throw new Error("clip title");
  if (/Hyperliquid \$10 minimum/.test(clip.body)) throw new Error("generic $10");
  const below = explainCommittee("below_min_notional");
  if (/Host sized under the Hyperliquid \$10/.test(below.body)) throw new Error("generic $10 below");
}

export function assertAgenticIdPartial() {
  const row = CAPABILITY.find((r) => r.id === "agentic");
  if (!row) throw new Error("agentic row");
  if (row.mainnet !== "partial") throw new Error("agentic must not be OFF while mint is live");
  if (!/iTransfer/.test(row.note) || !/UNAVAILABLE/.test(row.note)) throw new Error("iTransfer unavailable");
  if (!/mint/i.test(row.note)) throw new Error("mint live");
}

export function assertDeskHeadlinePrefersExecutable() {
  const open = deskHeadline({
    doing: "Manual. Waiting for you.",
    ready: true,
    canOpen: true,
    execN: 6,
    researchKind: "READY_STOOD_DOWN",
    attentionTitle: "Desk is ready",
  });
  if (open !== "6 books can open") throw new Error(open);
  const tight = deskHeadline({
    doing: "Idle",
    ready: false,
    canOpen: false,
    execN: 0,
    researchKind: "READY_STOOD_DOWN",
    attentionTitle: "Policy cap is too tight",
    policyTight: true,
  });
  if (tight !== "Policy cap is too tight") throw new Error(tight);
  if (researchCardTitle("POLL_FAILED") !== "CONNECTION CHECK MISSED") throw new Error("poll taxonomy");
}
