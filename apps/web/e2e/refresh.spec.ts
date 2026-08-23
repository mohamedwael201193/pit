// MOCK TEST HARNESS — public UI copy only.

export function assertRefreshCannotSign(canSign: boolean) {
  if (canSign) {
    throw new Error("browser cannot sign after refresh");
  }
}
