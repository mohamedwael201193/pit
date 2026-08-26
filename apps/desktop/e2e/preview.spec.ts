// MOCK TEST HARNESS — public UI copy only. Never stub a live preview hash.

export function assertPreviewMutationCopy(copy: string) {
  if (!copy.includes("mutation invalidates")) {
    throw new Error("preview copy");
  }
}
