import { committeeVerified } from "./honesty";

type Role = {
  role?: string;
  verify_e2ee?: string;
  proposed_side?: string;
  survives?: boolean;
  kill?: boolean;
};

export function researchWhyCopy(input: {
  coin: string;
  kind?: string;
  note?: string | null;
  stop?: string | null;
  deny?: string;
  eligible?: boolean;
  roles: Role[];
  snap?: { mark?: number; reason?: string; why?: string };
}): { q: string; a: string }[] {
  const verified = committeeVerified(input.roles);
  const sides = input.roles
    .map((r) => `${String(r.role || "")}: ${r.proposed_side || "none"}`)
    .filter((s) => !s.startsWith(":"));
  const disagree = input.roles.some((r) => r.kill || r.survives === false)
    ? "Challenger or risk did not survive the thesis."
    : sides.length
      ? `Roles proposed ${sides.join("; ")}.`
      : "The committee did not record a disagreement yet.";
  const stood = input.kind === "READY_STOOD_DOWN" || input.deny === "no_side" || input.stop === "READY_STOOD_DOWN";
  const blocked = input.kind === "POLICY_DENIED" || input.kind === "POLICY_REJECTED";
  const accepted = verified && input.eligible && !stood && !blocked;
  const found = input.snap?.why || input.note || (input.coin ? `Public ${input.coin} book under policy.` : "No completed pass yet.");
  let change = "Protect private compute, pin policy, and run a sealed pass on an eligible market.";
  if (stood) change = "A different thesis, market, or time. Host will not invent a side.";
  if (blocked) change = "Pin a policy that allows this market, or pick a coin that already passes.";
  if (input.stop === "DIRECT_CREDIT_INSUFFICIENT") change = "Fund Direct with the same wallet. That is compute money, not trading capital.";
  if (input.stop === "DIRECT_PROVIDER_TIMEOUT") change = "Retry when the private provider answers. A timeout is not a TEE failure.";
  if (input.stop === "TEE_VERIFY_FAIL" || input.stop === "TEE_SIGNATURE_INVALID") change = "Do not accept this result. Start a new sealed pass after the provider verifies.";
  return [
    { q: "What did PIT find?", a: found },
    { q: "Why did it matter?", a: verified ? "Three sealed roles finished over the same private book. Host still sizes." : "The committee has not completed. Incomplete work is not verified." },
    { q: "What did the committee disagree with?", a: disagree },
    {
      q: "Why was the trade accepted or rejected?",
      a: accepted
        ? "Accepted for preview only. You still type AUTHORIZE on the exact card."
        : stood
          ? "Rejected by the committee. That is a verified stand-down, not a crash."
          : blocked
            ? "Policy blocked it. The model cannot override host law."
            : input.stop
              ? `Stopped: ${input.stop.replaceAll("_", " ")}.`
              : "No trade was accepted. Chat cannot AUTHORIZE.",
    },
    {
      q: "What risk stopped it?",
      a: input.roles.some((r) => String(r.role).toLowerCase() === "risk" && (r.kill || r.survives === false))
        ? "Risk killed the idea after verifying the sealed book."
        : blocked
          ? "Host policy."
          : stood
            ? "No side survived."
            : "No order was placed.",
    },
    { q: "What would need to change?", a: change },
  ];
}
