// MOCK TEST HARNESS — public UI copy only.

export function assertTwoWalletsIsolated(copy: string) {
  if (!copy.includes("never share a workspace")) {
    throw new Error("isolation copy");
  }
}
