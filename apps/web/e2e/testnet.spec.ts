// MOCK TEST HARNESS — public UI copy only. Never stub a live order.

export function assertTestnetLab(net: string, explorer: string) {
  if (net !== "testnet") {
    throw new Error("testnet");
  }
  if (!explorer.includes("chainscan-galileo.0g.ai")) {
    throw new Error("explorer");
  }
}
