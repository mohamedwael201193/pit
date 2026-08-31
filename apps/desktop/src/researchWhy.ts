import { committeeVerified } from "./honesty";

type Role = {
  role?: string;
  verify_e2ee?: string;
  proposed_side?: string;
  survives?: boolean;
  kill?: boolean;
};

function named(roles: Role[], name: string) {
  return roles.find((r) => String(r.role || "").toLowerCase() === name);
}

function sideOf(roles: Role[], name: string) {
  const row = named(roles, name);
  if (!row) return "This role has not finished.";
  const side = String(row.proposed_side || "").trim();
  if (row.kill || row.survives === false) {
    return side ? `Stood down after proposing ${side}.` : "Stood down. No side survived.";
  }
  if (String(row.verify_e2ee || "").toUpperCase() === "OK") {
    return side ? `Verified. Proposed ${side}.` : "Verified. No side recorded.";
  }
  return "Not verified yet. Incomplete work is not a committee result.";
}

export function researchWhyCopy(input: {
  coin: string;
  kind?: string;
  note?: string | null;
  stop?: string | null;
  deny?: string;
  eligible?: boolean;
  roles: Role[];
  snap?: {
    mark?: number;
    reason?: string;
    why?: string;
    whyRanked?: string;
    invalidation?: string;
    expectedEdge?: string;
  };
}): { q: string; a: string }[] {
  const verified = committeeVerified(input.roles);
  const stood = input.kind === "READY_STOOD_DOWN" || input.deny === "no_side" || input.stop === "READY_STOOD_DOWN";
  const blocked = input.kind === "POLICY_DENIED" || input.kind === "POLICY_REJECTED";
  const accepted = verified && input.eligible && !stood && !blocked;
  const found = accepted
    ? `Verified ${input.coin || "this book"} reached an exact host preview. Host sized under the pinned clip.`
    : input.snap?.whyRanked ||
      input.snap?.why ||
      input.note ||
      (input.coin ? `Public ${input.coin} book under policy.` : "No completed pass yet.");
  let change = "Protect private compute, pin policy, and run a sealed pass on an eligible market.";
  if (stood) change = "Checking the next eligible book, a different thesis, or more venue margin. Host will not invent a side.";
  if (blocked) change = "Pin a policy that allows this market, or pick a coin that already passes.";
  if (input.stop === "DIRECT_CREDIT_INSUFFICIENT") change = "Fund Direct with the same wallet. That is compute money, not trading capital.";
  if (input.stop === "DIRECT_PROVIDER_TIMEOUT") change = "Retry when the private provider answers. A timeout is not a TEE failure.";
  if (input.stop === "TEE_VERIFY_FAIL" || input.stop === "TEE_SIGNATURE_INVALID") change = "Do not accept this result. Start a new sealed pass after the provider verifies.";
  const riskRow = named(input.roles, "risk");
  const riskLine =
    riskRow?.kill || riskRow?.survives === false
      ? "Risk killed the idea after verifying the sealed book."
      : blocked
        ? "Host policy blocked it. Risk is not a substitute for policy."
        : stood
          ? "No side survived. Engine did not size."
          : accepted
            ? "Risk allowed the thesis to reach the host engine. Host still sizes the clip."
            : "No order was placed.";
  const engine =
    accepted
      ? "Host sized an exact preview under policy. TRADE NOW on this computer submits it through the host."
      : stood
        ? "Host did not size. A verified stand-down is the result, not a crash."
        : blocked
          ? "Policy blocked sizing. The model cannot raise clip."
          : input.stop
            ? `Stopped: ${input.stop.replaceAll("_", " ")}. Host did not size.`
            : "Host did not size a trade.";
  return [
    { q: "What did PIT find?", a: found },
    {
      q: "Why is it interesting?",
      a: verified
        ? "Three named roles finished over the same private book. Host still sizes. One verified role is never a committee result."
        : "The committee has not completed. Incomplete work is not verified.",
    },
    { q: "What did the researcher think?", a: sideOf(input.roles, "researcher") },
    { q: "What did the challenger attack?", a: sideOf(input.roles, "challenger") },
    { q: "What did risk reject or allow?", a: riskLine },
    {
      q: "What did policy enforce?",
      a: blocked
        ? "Host policy blocked this market, clip, or permission. Chat cannot raise it."
        : "Host clip, universe, leverage, kill switch, and slippage stay on this computer.",
    },
    { q: "Why did the engine size or reject it?", a: engine },
    {
      q: "What would need to change?",
      a: input.snap?.invalidation
        ? `Live host invalidation: ${input.snap.invalidation}`
        : change,
    },
  ];
}
