export function explainStop(code: string | null): { title: string; body: string } | null {
  if (!code) return null;
  if (code === "TEE_VERIFY_FAIL") {
    return {
      title: "Research stopped",
      body: "We could not verify that the AI result came from the provider PIT expected. No order was placed. No funds moved.",
    };
  }
  if (code === "direct_token_required" || code === "sealer_not_wired") {
    return {
      title: "Private research is not armed",
      body: "This computer does not have a Direct TeeML token yet. Watch still works. No order was placed.",
    };
  }
  if (code === "companion_down") {
    return {
      title: "Local PIT is not running",
      body: "The desktop companion on this computer did not start. No order was placed. Reinstall PIT or run pit companion from a terminal.",
    };
  }
  if (code === "galileo_e2ee_unproven") {
    return {
      title: "TESTNET sealed research is off",
      body: "Galileo TeeML is not enabled until VerifyE2EE is proven. Switch to MAINNET for production research, or keep TESTNET as the lab.",
    };
  }
  return {
    title: "Stopped",
    body: "PIT halted this step. No order was placed. No funds moved.",
  };
}
