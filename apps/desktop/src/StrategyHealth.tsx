export function StrategyHealth({
  copy,
  n,
  need,
  enough,
  skills,
}: {
  copy?: string;
  n?: number;
  need?: number;
  enough?: boolean;
  skills?: Array<{ id?: string; title?: string; version?: string; n?: number; copy?: string }>;
}) {
  const rows = skills || [];
  const ok = Boolean(enough);
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Strategy Health</p>
          <h1>Resolved observations only</h1>
        </div>
        <span className={`status-chip ${ok ? "active" : "blocked"}`}>{ok ? "READY" : "NOT ENOUGH DATA"}</span>
      </div>
      <p className="lead">
        PIT never invents skill performance. A skill with no resolved outcomes stays NOT ENOUGH DATA. The model did not “learn this.”
      </p>
      <section className="card">
        <p className="label">Calibration</p>
        <p>{copy || "NOT ENOUGH DATA"}</p>
        <p className="fine">
          Sample {n ?? 0} / need {need ?? 0}. Brier, ECE, and drift stay unpublished until that floor.
        </p>
      </section>
      <p className="label">Skills</p>
      {rows.length === 0 ? (
        <p className="empty">NOT ENOUGH DATA</p>
      ) : (
        <ul className="skill-list">
          {rows.map((s) => (
            <li key={s.id || s.title}>
              <strong>{s.title || s.id}</strong>
              <span>{s.version || "1.0.0"}</span>
              <span>{typeof s.n === "number" && s.n > 0 ? `${s.n} resolved observations` : "NOT ENOUGH DATA"}</span>
            </li>
          ))}
        </ul>
      )}
    </main>
  );
}
