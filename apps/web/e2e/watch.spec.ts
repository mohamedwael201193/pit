// MOCK TEST HARNESS — public UI copy only. Never stub opportunities.

export function assertWatchCopy(copy: string, count: number) {
  if (count <= 0 && copy !== "No opportunities match your policy.") {
    throw new Error("empty watch");
  }
  if (count > 0 && !copy.includes("opportunities match your policy")) {
    throw new Error("count");
  }
}
