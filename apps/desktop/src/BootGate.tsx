type Row = { id: string; label: string; state: "ok" | "wait" | "fail" };

export function BootGate({
  open,
  rows,
  stuck,
}: {
  open: boolean;
  rows: Row[];
  stuck?: string;
}) {
  if (!open) return null;
  return (
    <div className="boot" role="status" aria-live="polite">
      <div className="boot-card">
        <p className="eyebrow">PIT</p>
        <h1>Starting local desk</h1>
        <ul className="boot-list">
          {rows.map((r) => (
            <li key={r.id} className={r.state}>
              <span aria-hidden="true">{r.state === "ok" ? "✓" : r.state === "fail" ? "✕" : "·"}</span>
              <span>{r.label}</span>
            </li>
          ))}
        </ul>
        {stuck ? <p className="fine">{stuck}</p> : <p className="fine">Each line is a live check on this computer. Nothing is a timer.</p>}
      </div>
    </div>
  );
}
