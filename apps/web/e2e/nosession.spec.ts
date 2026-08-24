// MOCK TEST HARNESS — public UI copy only.

export function assertNoSessionCopy(copy: string) {
  if (!copy.toLowerCase().includes("never holds a hyperliquid session")) {
    throw new Error("session copy");
  }
}
