import type { PathStep } from "./setupPath";

export function WelcomePath({
  steps,
  onGo,
}: {
  steps: PathStep[];
  onGo?: (view: NonNullable<PathStep["go"]>) => void;
}) {
  return (
    <article className="card">
      <p className="label">WELCOME TO PIT</p>
      <ol className="setup-path">
        {steps.map((s, i) => (
          <li key={s.id}>
            <span className={`tone ${s.tone.toLowerCase()}`}>{s.tone}</span>
            <div>
              <strong>
                {i + 1}. {s.title}
              </strong>
              <p>{s.why}</p>
              <div className="cta-row">
                {s.href ? (
                  <a className="linkish" href={s.href} target="_blank" rel="noreferrer">
                    {s.hrefLabel || "Open"}
                  </a>
                ) : null}
                {s.go && onGo ? (
                  <button type="button" className="linkish" onClick={() => onGo(s.go!)}>
                    {s.goLabel || "Open"}
                  </button>
                ) : null}
              </div>
            </div>
          </li>
        ))}
      </ol>
    </article>
  );
}
