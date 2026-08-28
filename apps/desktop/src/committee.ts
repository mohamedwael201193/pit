export const COMMITTEE_DENY = new Set([
  "risk_killed",
  "challenger_killed",
  "no_side",
  "below_min_notional",
  "policy_denied",
  "kill_switch",
  "coin_not_allowed",
  "leverage_above_policy",
  "CHALLENGER_KILLED",
  "RISK_KILLED",
  "NO_SIDE",
]);

export function committeeDeny(code?: string | null): boolean {
  if (!code) return false;
  return COMMITTEE_DENY.has(code) || COMMITTEE_DENY.has(code.toLowerCase());
}

export function explainCommittee(code: string): { title: string; body: string } {
  const c = code.toLowerCase();
  if (c === "risk_killed") {
    return {
      title: "Committee stood down — risk",
      body: "Researcher, challenger, and risk all verified. Risk killed the idea. That is a verified result, not a TEE failure. No order was placed.",
    };
  }
  if (c === "challenger_killed") {
    return {
      title: "Committee stood down — challenger",
      body: "Researcher, challenger, and risk all verified. The challenger did not survive the thesis. That is a verified result, not a crash. No order was placed.",
    };
  }
  if (c === "no_side") {
    return {
      title: "Committee stood down — no side",
      body: "The sealed committee did not propose a side. Host will not invent a trade. No order was placed.",
    };
  }
  if (c === "below_min_notional") {
    return {
      title: "Below venue minimum",
      body: "Host sized under the Hyperliquid $10 minimum. PIT will not pad a fake size. No order was placed.",
    };
  }
  if (c === "kill_switch") {
    return {
      title: "Kill switch is on",
      body: "Local execution is halted. No order was placed.",
    };
  }
  if (c === "coin_not_allowed" || c === "policy_denied" || c === "leverage_above_policy") {
    return {
      title: "Policy blocked this preview",
      body: "The pinned policy rejected this market, size, or leverage. The model cannot raise clip. No order was placed.",
    };
  }
  return {
    title: "Committee did not size a trade",
    body: "Verified sealed research produced no eligible preview. No order was placed.",
  };
}
