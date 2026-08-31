import { computeOnboard, onboardInput } from "../src/onboard";
import { nextFix } from "../src/nextFix";
import type { DoctorCheck, LocalStatus } from "../src/companion";
import { readFileSync } from "node:fs";

const base = {
  companionUp: true,
  paired: false,
  walletOk: false,
  protectOk: false,
  sessionAlive: false,
  hlApproved: false,
  policyPinned: false,
};

export function assertOnboardPairFirst() {
  const fresh = computeOnboard(base);
  if (fresh.current !== "pair") throw new Error("new user must start at pair");
  if (fresh.ready) throw new Error("new user is not ready");
  if (fresh.steps[1].locked !== true) throw new Error("protect must lock until paired");
  if (fresh.steps[2].locked !== true) throw new Error("hyperliquid must lock until protect");
}

export function assertOnboardDoesNotInventReady() {
  const almost = computeOnboard({ ...base, paired: true, protectOk: true, sessionAlive: true, policyPinned: true });
  if (almost.ready) throw new Error("missing Hyperliquid approval is not ready");
  if (almost.current !== "hyperliquid") throw new Error("current must wait for Hyperliquid");
  if (almost.steps[2].state !== "WAITING FOR APPROVAL") throw new Error("hl state");
}

export function assertOnboardReadyRequiresAllGates() {
  const all = computeOnboard({
    ...base,
    paired: true,
    walletOk: true,
    protectOk: true,
    sessionAlive: true,
    hlApproved: true,
    policyPinned: true,
  });
  if (!all.ready) throw new Error("all gates should be ready");
  if (all.current !== "ready") throw new Error("current ready");
  const k = computeOnboard({
    ...base,
    paired: true,
    walletOk: true,
    protectOk: true,
    sessionAlive: true,
    hlApproved: true,
    policyPinned: true,
    kill: true,
  });
  if (k.ready) throw new Error("kill is not ready");
}

export function assertNextFixPairsFirst() {
  const checks: DoctorCheck[] = [];
  const status: LocalStatus = { sign: false, paired: false };
  const n = nextFix(true, status, checks, [], false, "mainnet");
  if (!/Pair this browser/i.test(n.title)) throw new Error(`expected pair first, got ${n.title}`);
}

export function assertNextFixDoesNotSkipPairWhenSessionExists() {
  const checks: DoctorCheck[] = [
    { name: "wallet", ok: true, detail: "bound" },
    { name: "direct_auth", ok: true, detail: "ok" },
    { name: "hl_agent", ok: true, detail: "listed" },
    { name: "policy", ok: true, detail: "pinned" },
  ];
  const status: LocalStatus = { sign: false, paired: false, wallet: "0xabc" };
  const n = nextFix(true, status, checks, [], true, "mainnet");
  if (!/Pair this browser/i.test(n.title)) throw new Error(`live session still needs pair first, got ${n.title}`);
}

export function assertOnboardInputFromDoctor() {
  const got = onboardInput(
    true,
    { sign: false, paired: true, pairingDevices: 1, wallet: "0xab" },
    [{ name: "direct_auth", ok: true, detail: "ok" }],
    false,
  );
  if (!got.paired || !got.protectOk || !got.walletOk) throw new Error("onboardInput");
}

export function assertPolicyEditorStaysAfterReady() {
  const src = readFileSync(new URL("../src/SecurityCenter.tsx", import.meta.url), "utf8");
  if (/current\.id === ["']policy["'] \?/.test(src)) {
    throw new Error("policy editor must stay on Security after pin");
  }
  if (!src.includes("Re-pin anytime")) throw new Error("re-pin copy");
  if (!src.includes("id=\"policy\"")) throw new Error("policy section");
}
