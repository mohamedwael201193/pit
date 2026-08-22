// MOCK TEST HARNESS — public UI copy only. Never stub VerifyE2EE success.
// Do not count this file as live Hyperliquid or TeeML evidence.

export function assertConnectCopy(text: string) {
  if (!text.includes("seed phrase")) {
    throw new Error("seed copy missing");
  }
  if (text.toLowerCase().includes("withdraw") && text.toLowerCase().includes("allowed")) {
    throw new Error("withdraw must stay denied");
  }
}
