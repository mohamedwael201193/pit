// MOCK TEST HARNESS — public UI copy only. Never stub two live wallets.

export function assertDistinctWorkspaces(idA: string, idB: string) {
  if (!idA || !idB || idA === idB) {
    throw new Error("same workspace");
  }
}
