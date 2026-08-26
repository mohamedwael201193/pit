// MOCK TEST HARNESS — public UI copy only. Never create a session in the browser.

export function assertWebCannotHoldSession(copy: string) {
  if (!copy.toLowerCase().includes("never") || !copy.toLowerCase().includes("session")) {
    throw new Error("web session copy");
  }
}
