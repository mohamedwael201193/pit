export type RoleLike = { role?: string; verify_e2ee?: string };

export function committeeVerified(roles: RoleLike[] | undefined, hostVerify?: boolean): boolean {
  if (hostVerify === true) {
    if (!roles || roles.length < 3) return false;
  }
  if (!roles || roles.length < 3) return false;
  const ok = new Set<string>();
  for (const r of roles) {
    if (String(r.verify_e2ee).toUpperCase() === "OK") {
      ok.add(String(r.role || "").toLowerCase());
    }
  }
  return ok.has("researcher") && ok.has("challenger") && ok.has("risk");
}

export function oidBelongsToPreview(orderHash?: string | null, previewHash?: string | null, previewHashAlt?: string | null): boolean {
  const o = String(orderHash || "");
  if (!o) return false;
  return o === String(previewHash || "") || o === String(previewHashAlt || "");
}

export function researchCardTitle(kind?: string | null, verified?: boolean): string {
  if (kind === "READY_ELIGIBLE" || (verified && !kind)) return "RESEARCH COMPLETE";
  if (kind === "READY_STOOD_DOWN") return "NO TRADE SURVIVED CHALLENGE";
  if (kind === "COMMITTEE_INCOMPLETE") return "RESEARCH INCOMPLETE";
  if (kind === "CANCELED_BY_USER") return "YOU CANCELLED";
  if (kind === "POLICY_DENIED") return "POLICY BLOCKED THIS";
  if (kind === "MARKET_DENIED") return "MARKET NOT USABLE";
  if (kind === "DIRECT_NOT_AUTHORIZED") return "PRIVATE RESEARCH IS NOT ARMED";
  if (kind === "DIRECT_CREDIT_INSUFFICIENT") return "PRIVATE COMPUTE NEEDS FUNDS";
  if (kind === "DIRECT_PROVIDER_TIMEOUT") return "PROVIDER TIMEOUT";
  if (kind === "DIRECT_PROVIDER_UNAVAILABLE") return "PRIVATE PROVIDER UNREACHABLE";
  if (kind === "DIRECT_RATE_LIMITED") return "TOO MANY SEALED REQUESTS";
  if (kind === "JOB_CRASHED") return "RESEARCH PROCESS STOPPED";
  if (kind === "POLL_FAILED") return "CONNECTION CHECK MISSED";
  if (kind === "SPONSOR_QUOTA") return "SPONSORED RESEARCH PAUSED";
  if (kind === "WRONG_NETWORK") return "WRONG NETWORK";
  if (kind === "COMPANION_NOT_RUNNING" || kind === "companion_down") return "COMPANION FAILURE";
  if (kind === "TEE_SIGNATURE_INVALID" || kind === "TEE_SIGNER_MISMATCH" || kind === "TEE_VERIFY_FAIL") {
    return "TEE VERIFICATION FAILED";
  }
  if (kind) return "RESEARCH STATUS";
  return "RESEARCH";
}
