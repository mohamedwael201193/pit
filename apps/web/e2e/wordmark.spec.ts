// MOCK TEST HARNESS — public UI copy only. Never stub AUTHORIZE.

export function assertHomeHref(href: string) {
  if (href !== "/" && href !== "http://127.0.0.1:5173/") {
    throw new Error("home");
  }
}
