// MOCK TEST HARNESS — public UI copy only. Never stub VerifyE2EE success.
// Do not count this file as live Hyperliquid or TeeML evidence.

import { NAMED, PERMISSIONS } from "../src/namedStates";

export function assertDesktopCopy() {
  if (!NAMED.SEED_FORBIDDEN.includes("seed")) {
    throw new Error("seed copy");
  }
  const withdraw = PERMISSIONS.find((p) => p.k === "withdraw");
  if (!withdraw || withdraw.ok) {
    throw new Error("withdraw must be denied");
  }
}
