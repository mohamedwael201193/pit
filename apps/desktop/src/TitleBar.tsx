import { nativeInvoke } from "./companion";

export function TitleBar({ status }: { status: string }) {
  return (
    <header className="titlebar" role="banner">
      <div className="titlebar-drag" data-tauri-drag-region>
        <span className="word">PIT</span>
        <span className="titlebar-sub">Private Trading Desk · {status}</span>
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
