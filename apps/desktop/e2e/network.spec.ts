// MOCK TEST HARNESS — public UI copy only. Never stub VerifyE2EE or AUTHORIZE.

export function assertNetworkToggle(net: "mainnet" | "testnet") {
  if (net !== "mainnet" && net !== "testnet") {
    throw new Error("network");
  }
}
