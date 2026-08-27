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
  "SIGNING",
  "SUBMITTING",
  "CONFIRMING",
  "EXECUTED",
  "VERIFYING",
  "RESOLVED",
  "CALIBRATED",
  "FAILED",
];

export function ProgressStrip({ current }: { current: string }) {
  const active = Known(current) ? current : "FAILED";
  const idx = STEPS.indexOf(active);
  return (
    <ol className="mt-8 grid gap-0 border-y border-[rgb(240_231_212/0.22)]">
      {STEPS.map((s, i) => (
        <li
          key={s}
          className={
            s === active
              ? "border-b border-[rgb(240_231_212/0.12)] bg-[#d82f2f] px-4 py-3 font-mono text-[0.8125rem] text-black"
              : i < idx
                ? "border-b border-[rgb(240_231_212/0.12)] px-4 py-3 font-mono text-[0.8125rem] text-[rgb(240_231_212/0.45)]"
                : "border-b border-[rgb(240_231_212/0.12)] px-4 py-3 font-mono text-[0.8125rem] text-[rgb(240_231_212/0.28)]"
          }
        >
          {s}
        </li>
      ))}
    </ol>
  );
}
