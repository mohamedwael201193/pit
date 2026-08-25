// MOCK TEST HARNESS — public UI copy only. Never stub AUTHORIZE.

export function assertMainnetProduction(net: string, explorer: string) {
  if (net !== "mainnet") {
    throw new Error("mainnet");
  }
  if (!explorer.includes("chainscan.0g.ai") || explorer.includes("galileo")) {
    throw new Error("explorer");
  }
}
