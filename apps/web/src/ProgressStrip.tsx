import { Known } from "./phases";

const STEPS = [
  "CONNECTING",
  "AUTHENTICATING",
  "SEALING",
  "RESEARCHING",
  "CHALLENGING",
  "ASSESSING_RISK",
  "SCORING",
  "POLICY_CHECK",
  "WAITING_FOR_USER",
];

export function ProgressStrip({ current }: { current: string }) {
  const active = Known(current) ? current : "FAILED";
  return (
    <ol className="mt-8 space-y-1">
      {STEPS.map((s) => (
        <li key={s} className={s === active ? "text-coral" : "opacity-40"}>
          {s}
        </li>
      ))}
    </ol>
  );
}
