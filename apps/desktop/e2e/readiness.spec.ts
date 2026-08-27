// MOCK TEST HARNESS — public UI copy only. Never stub VerifyE2EE success.

import { probes } from "../src/readiness";

export function assertTeeNeverGreenIdle() {
  const items = probes([], null, true);
  const tee = items.find((p) => p.id === "tee");
  if (!tee || tee.state === "ok") {
    throw new Error("tee must stay waiting without a verified run");
  }
}

export function assertDirectNotGreenWithoutAuth() {
  const items = probes(
    [{ name: "direct_sealer", ok: true, detail: "binary present" }],
    null,
    true,
  );
  const direct = items.find((p) => p.id === "direct");
  if (!direct || direct.state === "ok") {
    throw new Error("direct must not be green without auth file");
  }
}
