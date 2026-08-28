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
      body: "The desktop companion on this computer did not start. No order was placed. Launch PIT Desktop again.",
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
  if (code === "research_cancelled") {
    return {
      title: "Research stopped",
      body: "You cancelled this sealed request. No order was placed. No funds moved.",
    };
  }
  if (code === "companion_http") {
    return {
      title: "Local PIT dropped a status poll",
      body: "The sealed job may still be running. Retry. PIT did not place an order.",
    };
  }
  if (code === "COMPANION_NOT_RUNNING") {
    return {
      title: "Local PIT dropped status polls",
      body: "A sealed job may still finish on this computer. Open Research again. PIT did not place an order. Do not launch a second PIT window.",
    };
  }
  if (code === "FAILED") {
    return {
      title: "Research failed closed",
      body: "The sealed path stopped before a verified committee. PIT did not fall back to Router. No order was placed.",
    };
  }
  if (code === "WORKSPACE_NOT_BOUND") {
    return {
      title: "This computer is not bound",
      body: "Bind your public wallet on this computer first. No order was placed.",
    };
  }
  if (code === "DIRECT_NOT_AUTHORIZED") {
    return {
      title: "Private research is not armed",
      body: "Pair the browser and sign Protect my strategy. No order was placed.",
    };
  }
  if (code === "SPONSOR_QUOTA") {
    return {
      title: "Sponsored research is paused for today",
      body: "This workspace used its daily sponsored sealed-research cap. Fund your own Direct sub-account at pc.0g.ai Advanced, or wait until tomorrow. No order was placed.",
    };
  }
  if (code === "DIRECT_PROVIDER_UNAVAILABLE") {
    return {
      title: "0G Direct credit ran out",
      body: "Researcher may already be verified. Challenger and risk did not finish because the provider rejected the next sealed call (insufficient locked Direct balance). Top up the same wallet at pc.0g.ai Advanced. PIT did not fall back to Router. No order was placed.",
    };
  }
  if (code === "HL_MARKET_UNAVAILABLE") {
    return {
      title: "Not enough market data",
      body: "PIT will not invent a book. Watch still works. No order was placed.",
    };
  }
  if (code === "TEE_SIGNATURE_INVALID" || code === "TEE_RESPONSE_INVALID") {
    return {
      title: "Research stopped",
      body: "We could not verify that the AI result came from the provider PIT expected. No order was placed. No funds moved.",
    };
  }
  if (code === "POLICY_REJECTED") {
    return {
      title: "Policy blocked this research",
      body: "The pinned policy rejected this market or the kill switch is on. No order was placed.",
    };
  }
  if (code === "TEE_SIGNER_MISMATCH") {
    return {
      title: "Research stopped",
      body: "The recovered signer did not match the on-chain teeSigner PIT expected. No order was placed.",
    };
  }
  if (code === "RESEARCHER_FAILED" || code === "CHALLENGER_FAILED" || code === "RISK_FAILED") {
    return {
      title: "Committee did not complete",
      body: "Researcher, challenger, and risk must all finish over the same sealed book. No order was placed.",
    };
  }
  if (code === "asset_not_allowed") {
    return {
      title: "Market is outside policy",
      body: "PIT will not research a coin that is not in the pinned universe. No order was placed.",
    };
  }
  if (code === "kill_switch") {
    return {
      title: "Kill switch is on",
      body: "Local execution is halted. No order was placed.",
    };
  }
  return {
    title: "Research stopped",
    body: "PIT halted this step. No order was placed. No funds moved.",
  };
}
