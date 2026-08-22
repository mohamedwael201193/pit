// MOCK TEST HARNESS — public UI copy only. Never stub VerifyE2EE success.

export function assertNetworkCopy(text: string, net: "mainnet" | "testnet") {
  if (net === "testnet" && !text.toLowerCase().includes("lab") && !text.toLowerCase().includes("testnet")) {
    throw new Error("lab copy");
  }
  if (net === "mainnet" && text.toLowerCase().includes("galileo") && text.toLowerCase().includes("production")) {
    throw new Error("mixed networks");
  }
}
