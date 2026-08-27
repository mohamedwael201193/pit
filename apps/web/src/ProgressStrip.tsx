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
  const nearby = STEPS.slice(Math.max(0, idx - 1), Math.min(STEPS.length, idx + 3));

  return (
    <div>
      <p className="text-[0.75rem] tracking-[0.16em] text-[rgb(240_231_212/0.5)] uppercase">Now</p>
      <p className="mt-2 font-mono text-[1.35rem] text-[#d82f2f]">{active}</p>
      <ol className="mt-6 flex flex-col gap-2">
        {nearby.map((s) => (
          <li
            key={s}
            className={
              s === active
                ? "font-mono text-[0.875rem] text-[var(--guide-cream)]"
                : "font-mono text-[0.8125rem] text-[rgb(240_231_212/0.4)]"
            }
          >
            {s}
          </li>
        ))}
      </ol>
    </div>
  );
}
