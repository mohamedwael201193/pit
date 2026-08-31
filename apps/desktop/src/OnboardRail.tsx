import type { OnboardStep } from "./onboard";

export function OnboardRail({ steps }: { steps: OnboardStep[] }) {
  return (
    <ol className="onboard-rail" aria-label="Setup steps">
      {steps.map((s) => (
        <li
          key={s.id}
          className={`${s.done ? "done" : ""} ${s.current ? "current" : ""} ${s.locked ? "locked" : ""}`}
          aria-current={s.current ? "step" : undefined}
        >
          <span className="onboard-n">{s.n}</span>
          <strong>{s.title}</strong>
          <span className={`onboard-state ${s.state.replace(/\s+/g, "-").toLowerCase()}`}>{s.state}</span>
        </li>
      ))}
    </ol>
  );
}
