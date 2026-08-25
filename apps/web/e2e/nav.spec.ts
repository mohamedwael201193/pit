// MOCK TEST HARNESS — public UI copy only. Never stub iTransfer or AUTHORIZE.

export function assertVerifyNav(href: string) {
  if (href !== "#verify") {
    throw new Error("verify nav");
  }
}
