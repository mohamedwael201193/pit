// MOCK TEST HARNESS — public UI copy only. Never stub a signed exchange payload.

export function assertWebCannotPost(copy: string) {
  if (copy.toLowerCase().includes("posted to the venue from the browser")) {
    throw new Error("web post");
  }
}
