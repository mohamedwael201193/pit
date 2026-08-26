// MOCK TEST HARNESS — public UI copy only. Never stub VerifyE2EE success.

import { NAMED } from "../src/namedStates";

export function assertNamedErrors() {
  for (const k of ["SIGNATURE_DECLINED", "WRONG_NETWORK", "SESSION_EXPIRED", "POLICY_BLOCK", "TEE_VERIFY_FAIL", "AUTHORIZE_EXACT"] as const) {
    if (!NAMED[k]) {
      throw new Error(k);
    }
  }
}
