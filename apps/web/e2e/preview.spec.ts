// MOCK TEST HARNESS — public UI copy only. Never stub a live preview hash.

export function assertWebCannotAuthorizePreview(copy: string) {
  if (!copy.includes("never") || copy.toLowerCase().includes("authorize this preview in the browser")) {
    throw new Error("web preview");
  }
}
