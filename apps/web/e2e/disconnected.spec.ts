// MOCK TEST HARNESS — public UI copy only.

export function assertDisconnectedVerify(hash: string, connected: boolean) {
  if (hash !== "#verify") {
    throw new Error("verify");
  }
  if (connected) {
    return;
  }
}
