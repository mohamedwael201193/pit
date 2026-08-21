import { NAMED } from "./namedStates";

const STEPS = ["CONNECTING", "SEALING", "RESEARCHING", "CHALLENGING", "ASSESSING_RISK", "POLICY_CHECK", "WAITING_FOR_USER"];

export function Progress({ current }: { current: string }) {
  return (
    <ol className="perms">
      {STEPS.map((s) => (
        <li key={s} className={s === current ? "hot" : undefined}>
          {s}
        </li>
      ))}
      <li className="fine">{NAMED.TEE_VERIFY_FAIL}</li>
    </ol>
  );
}
