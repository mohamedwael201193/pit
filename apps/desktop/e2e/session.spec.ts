// MOCK TEST HARNESS — public UI copy only. Never print a session key.

export function assertSessionCopy(copy: string) {
  if (!copy.includes("never prints the key") && !copy.includes("PIT never prints the session key")) {
    throw new Error("session copy");
  }
}
