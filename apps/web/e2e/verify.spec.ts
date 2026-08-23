// MOCK TEST HARNESS — public UI copy only. Never stub VerifyE2EE success.

export function assertVerifyRoute(hash: string) {
  if (hash !== "#verify") {
    throw new Error("verify route");
  }
}
