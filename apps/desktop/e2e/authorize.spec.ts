// MOCK TEST HARNESS — public UI copy only. Never stub a live order or VerifyE2EE.

import { confirmAuthorize } from "../src/authorize";

export function assertAuthorizeClosed(typed: string, sessionAlive: boolean) {
  const err = confirmAuthorize(typed, sessionAlive);
  if (sessionAlive && typed === "AUTHORIZE") {
    if (err) {
      throw new Error("exact token should pass when the session is alive");
    }
    return;
  }
  if (!err) {
    throw new Error("authorize must fail closed");
  }
}
