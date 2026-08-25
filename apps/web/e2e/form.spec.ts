// MOCK TEST HARNESS — public UI copy only. Never stub a live receipt or AUTHORIZE.

export function assertVerifyFields(hash: string, root: string) {
  if (!hash.startsWith("0x") || !root.startsWith("0x")) {
    throw new Error("verify fields");
  }
}
