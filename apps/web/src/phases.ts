export function Known(s: string): boolean {
  return [
    "CONNECTING",
    "AUTHENTICATING",
    "SEALING",
    "RESEARCHING",
    "CHALLENGING",
    "ASSESSING_RISK",
    "SCORING",
    "POLICY_CHECK",
    "WAITING_FOR_USER",
    "SIGNING",
    "SUBMITTING",
    "CONFIRMING",
    "EXECUTED",
    "VERIFYING",
    "RESOLVED",
    "CALIBRATED",
    "FAILED",
  ].includes(s);
}
