import { LINKS } from "./links";

export function explainStopHref(code: string | null): { href: string; label: string } | null {
  if (!code) return null;
  if (
    code === "DIRECT_CREDIT_INSUFFICIENT" ||
    code === "DIRECT_PROVIDER_TIMEOUT" ||
    code === "DIRECT_PROVIDER_UNAVAILABLE" ||
    code === "direct_ledger" ||
    code === "SPONSOR_QUOTA"
  ) {
    return { href: LINKS.pcAdvanced, label: "Open 0G Private Compute" };
  }
  if (code === "approveAgent_required" || code === "AGENT_NOT_APPROVED" || code === "SESSION_EXPIRED") {
    return { href: LINKS.hlAPI, label: "Open Hyperliquid API" };
  }
  return null;
}

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
      title: "TEE verification failed",
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
      title: "Sealed response could not be opened",
      body: "The sealed response could not be opened. PIT did not fall back to Router or plaintext. No order was placed.",
    };
  }
  if (code === "research_cancelled") {
    return {
      title: "Research stopped",
      body: "You cancelled this sealed request. No order was placed. No funds moved.",
    };
  }
  if (code === "research_cancelled" || code === "CANCELED_BY_USER") {
    return {
      title: "You cancelled",
      body: "You cancelled this sealed request. No order was placed. No funds moved.",
    };
  }
  if (code === "DIRECT_RATE_LIMITED") {
    return {
      title: "Too many sealed requests",
      body: "The sealed provider rate-limited this workspace. Wait, then retry. That is not a TEE failure. No order was placed.",
    };
  }
  if (code === "WRONG_NETWORK") {
    return {
      title: "Wrong network",
      body: "This workspace is bound to one world. Mixing MAINNET compute with TESTNET venue is refused.",
    };
  }
  if (code === "POLICY_DENIED" || code === "POLICY_REJECTED") {
    return {
      title: "Policy blocked this",
      body: "The pinned policy rejected this market or the kill switch is on. No order was placed.",
    };
  }
  if (code === "MARKET_DENIED") {
    return {
      title: "Market not usable",
      body: "PIT will not invent a book. Watch still works. No order was placed.",
    };
  }
  if (code === "READY_STOOD_DOWN") {
    return {
      title: "Committee stood down",
      body: "Three roles verified. The committee did not propose a trade. That is a verified result, not a crash.",
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
      title: "Local PIT is not reachable",
      body: "Launch PIT Desktop on this computer. A missed status check does not cancel a sealed job. Open Research to reload the last verified result. No order was placed.",
    };
  }
  if (code === "POLL_FAILED") {
    return {
      title: "Connection check missed",
      body: "Research is still running on this computer. PIT did not place an order.",
    };
  }
  if (code === "JOB_CRASHED") {
    return {
      title: "Research process stopped",
      body: "The sealed job on this computer did not finish. PIT did not place an order. You can start a new research when private compute is ready.",
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
      title: "Private provider is not reachable",
      body: "The sealed provider did not accept this request. Earlier roles may already be verified. PIT did not fall back to Router. No order was placed.",
    };
  }
  if (code === "DIRECT_CREDIT_INSUFFICIENT") {
    return {
      title: "Not enough private compute",
      body: "PIT will not start a 3-role sealed committee without Direct capacity for this wallet. Open 0G Private Compute with the same wallet. No order was placed.",
    };
  }
  if (code === "DIRECT_PROVIDER_TIMEOUT") {
    return {
      title: "Private provider timed out",
      body: "The sealed request did not finish in time. That is a provider delay, not a TEE failure. Earlier roles may already be verified. Retry. No order was placed.",
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
      title: "TEE verification failed",
      body: "We could not verify that the AI result came from the provider PIT expected. No order was placed. No funds moved.",
    };
  }
  if (code === "COMMITTEE_INCOMPLETE") {
    return {
      title: "Research ended before all committee roles completed",
      body: "Researcher, challenger, and risk must all finish over the same sealed book. PIT did not invent a result. No order was placed.",
    };
  }
  if (code === "risk_killed" || code === "RISK_KILLED") {
    return {
      title: "Committee stood down — risk",
      body: "Researcher, challenger, and risk all verified. Risk killed the idea. That is a verified result, not a TEE failure. No order was placed.",
    };
  }
  if (code === "challenger_killed" || code === "CHALLENGER_KILLED") {
    return {
      title: "Committee stood down — challenger",
      body: "The challenger did not survive the thesis. That is a verified result. No order was placed.",
    };
  }
  if (code === "no_side" || code === "NO_SIDE") {
    return {
      title: "Committee stood down — no side",
      body: "The sealed committee did not propose a side. Host will not invent a trade. No order was placed.",
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
      title: "Unexpected TEE signer",
      body: "The recovered signer did not match the on-chain TEE signer PIT expected. No order was placed.",
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
  if (code === "approveAgent_required") {
    return {
      title: "Approve PIT on Hyperliquid",
      body: "This computer has an order/cancel session. Hyperliquid does not list it yet. Open Hyperliquid API and approve the agent shown on Security. PIT still cannot withdraw. No order was placed.",
    };
  }
  if (code === "preview_required" || code === "need_exact_AUTHORIZE" || code === "preview_hash_mismatch") {
    return {
      title: "Exact preview required",
      body: "Type AUTHORIZE on the card this computer generated. Piped yes is never enough. No order was placed.",
    };
  }
  return {
    title: "Research incomplete",
    body: "PIT halted this step. No order was placed. No funds moved.",
  };
}
