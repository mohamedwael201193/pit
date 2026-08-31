import { Link } from "react-router-dom";

const STEPS = [
  { n: 1, label: "Pair", to: "/pair" },
  { n: 2, label: "Protect", to: "/protect" },
  { n: 3, label: "Hyperliquid", to: "/protect" },
  { n: 4, label: "Policy", to: "/protect" },
  { n: 5, label: "Ready", to: "/protect" },
] as const;

export function OnboardRail({ current, paired }: { current: 1 | 2 | 3 | 4 | 5; paired?: boolean }) {
  return (
    <ol className="mb-8 flex flex-wrap gap-2" aria-label="Setup steps">
      {STEPS.map((s) => {
        const done = s.n === 1 ? Boolean(paired) : s.n < current && Boolean(paired);
        const here = s.n === current;
        return (
          <li key={s.n}>
            <Link
              to={s.n <= 2 ? s.to : "/pair"}
              className={`inline-flex items-center gap-2 rounded-full border px-3 py-1.5 text-[0.75rem] tracking-[0.08em] no-underline ${
                here
                  ? "border-[#d82f2f] text-[#f0e7d4]"
                  : done
                    ? "border-[rgb(125_255_179/0.35)] text-[#7dffb3]"
                    : "border-[rgb(240_231_212/0.2)] text-[rgb(240_231_212/0.45)]"
              }`}
              aria-current={here ? "step" : undefined}
            >
              <span>{s.n}</span>
              {s.label}
              {done ? " ✓" : ""}
            </Link>
          </li>
        );
      })}
    </ol>
  );
}
