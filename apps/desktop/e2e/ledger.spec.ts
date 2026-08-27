// MOCK TEST HARNESS — public UI copy only. Never stub a signed exchange payload.

export function assertLedgerUnsignedCopy(copy: string) {
  if (!copy.includes("unsigned") || !copy.includes("venue")) {
    throw new Error("ledger copy");
  }
}
