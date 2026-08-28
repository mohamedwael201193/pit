import { useEffect, useMemo, useState } from "react";

export type PaletteAction = { id: string; label: string; run: () => void };

export function CommandPalette({
  open,
  onClose,
  actions,
}: {
  open: boolean;
  onClose: () => void;
  actions: PaletteAction[];
}) {
  const [q, setQ] = useState("");
  const [i, setI] = useState(0);
  const rows = useMemo(() => {
    const n = q.trim().toLowerCase();
    return n ? actions.filter((a) => a.label.toLowerCase().includes(n)) : actions;
  }, [actions, q]);

  useEffect(() => {
    if (!open) {
      setQ("");
      setI(0);
    }
  }, [open]);

  useEffect(() => {
    if (i >= rows.length) setI(0);
  }, [rows.length, i]);

  if (!open) return null;
  return (
    <div
      className="palette"
      role="dialog"
      aria-label="Command palette"
      onClick={onClose}
    >
      <div className="palette-box" onClick={(e) => e.stopPropagation()}>
        <input
          autoFocus
          aria-label="Filter commands"
          placeholder="Open Desk, Markets, Research…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Escape") onClose();
            if (e.key === "ArrowDown") {
              e.preventDefault();
              setI((n) => Math.min(n + 1, Math.max(0, rows.length - 1)));
            }
            if (e.key === "ArrowUp") {
              e.preventDefault();
              setI((n) => Math.max(0, n - 1));
            }
            if (e.key === "Enter") {
              e.preventDefault();
              rows[i]?.run();
              onClose();
            }
          }}
        />
        <ul>
          {rows.map((a, n) => (
            <li key={a.id} className={n === i ? "on" : ""}>
              <button
                type="button"
                onClick={() => {
                  a.run();
                  onClose();
                }}
              >
                {a.label}
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
