// MOCK TEST HARNESS — public UI copy only. Never stub opportunities or AUTHORIZE.

export function assertEmptyHome(copy: string) {
  if (copy !== "No opportunities match your policy.") {
    throw new Error("empty home");
  }
}
