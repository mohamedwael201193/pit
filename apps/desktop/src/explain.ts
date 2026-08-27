export function explainStop(code: string | null): { title: string; body: string } | null {
  if (!code) return null;
  if (code === "direct_token_required" || code === "sealer_not_wired") {
    return {
      title: "Private research is not armed",
      body: "Pair the browser and sign Protect my strategy. PIT never asks you to edit an env file. Watch still works. No order was placed.",
    };
  }
  if (code === "direct_token_expired") {
    return {
      title: "Sealed-path signature expired",
      body: "Sign Protect my strategy again. The previous token lasted 24 hours. No order was placed.",
    };
  }
  if (code === "TEE_VERIFY_FAIL") {
    return {
      title: "Research stopped",
      body: "We could not verify that the AI result came from the provider PIT expected. No order was placed. No funds moved.",
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
  if (code === "signature_mismatch") {
    return {
      title: "Wrong wallet signed",
      body: "The bound wallet must sign Protect my strategy. No token was stored. No order was placed.",
    };
  }
  if (code === "empty_envelope") {
    return {
      title: "Not enough market data",
      body: "PIT will not invent a book. Watch still works. No order was placed.",
    };
  }
  if (code === "direct_ledger") {
    return {
      title: "0G Direct is not funded for this wallet",
      body: "The sealed-path signature was accepted on this computer. The provider still needs Direct credit at pc.0g.ai Advanced with the same wallet. Watch still works. No order was placed.",
    };
  }
  if (code === "TEE_OPEN_FAIL") {
    return {
      title: "Research stopped",
      body: "The sealed response could not be opened. PIT did not fall back to Router or plaintext. No order was placed.",
    };
  }
  return {
    title: "Stopped",
    body: "PIT halted this step. No order was placed. No funds moved.",
  };
}
