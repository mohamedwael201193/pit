// MOCK TEST HARNESS — desktop shell copy. Never stub VerifyE2EE.

import { MARKET_FILTERS } from "../src/WatchBook";

export function assertShellFilters() {
  if (MARKET_FILTERS.join(" ") !== "Actionable Research Watch Blocked") {
    throw new Error("market filters");
  }
}
