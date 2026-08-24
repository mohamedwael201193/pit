// MOCK TEST HARNESS — public UI copy only.

export function assertKillCopy(copy: string) {
  if (!copy.toLowerCase().includes("kill switch")) {
    throw new Error("kill copy");
  }
}
