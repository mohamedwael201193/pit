import { nativeInvoke } from "./companion";

export type StatusChip = { id: string; label: string; value: string; ok?: boolean };

export function TitleBar({ chips }: { chips: StatusChip[] }) {
  return (
    <header className="titlebar" role="banner">
      <div className="titlebar-drag" data-tauri-drag-region>
        <span className="word">PIT</span>
        <ul className="status-island" aria-label="Desk status">
          {chips.map((c) => (
            <li key={c.id} className={c.ok === false ? "fail" : c.ok ? "ok" : ""} title={`${c.label} ${c.value}`}>
              <span>{c.label}</span>
              <strong>{c.value}</strong>
            </li>
          ))}
        </ul>
      </div>
      <div className="win-controls" aria-label="Window">
        <button type="button" aria-label="Minimize" onClick={() => void nativeInvoke("window_min")}>
          –
        </button>
        <button type="button" aria-label="Maximize" onClick={() => void nativeInvoke("window_toggle_max")}>
          □
        </button>
        <button type="button" className="win-close" aria-label="Close" onClick={() => void nativeInvoke("window_close")}>
          ×
        </button>
      </div>
    </header>
  );
}
