import { useEffect } from "react";
import { isAllowedHttps } from "./allowedUrl";

export { isAllowedHttps } from "./allowedUrl";

type TauriWindow = Window & {
  __TAURI_INTERNALS__?: { invoke: (cmd: string, args?: unknown) => Promise<unknown> };
  __TAURI__?: { core?: { invoke: (cmd: string, args?: unknown) => Promise<unknown> } };
};

function isTauri(): boolean {
  const w = window as TauriWindow;
  return Boolean(w.__TAURI_INTERNALS__ || w.__TAURI__);
}

function invokeOpen(url: string): Promise<boolean> | null {
  const w = window as TauriWindow;
  const inv = w.__TAURI_INTERNALS__?.invoke ?? w.__TAURI__?.core?.invoke;
  if (typeof inv !== "function") return null;
  return inv("open_url", { url }) as Promise<boolean>;
}

export async function openExternal(url: string): Promise<boolean> {
  const href = String(url || "").trim();
  if (!isAllowedHttps(href)) return false;
  if (isTauri()) {
    try {
      const p = invokeOpen(href);
      if (!p) return false;
      return Boolean(await p);
    } catch {
      return false;
    }
  }
  const w = window.open(href, "_blank", "noopener,noreferrer");
  return Boolean(w);
}

export function useNativeExternalLinks() {
  useEffect(() => {
    const onClick = (e: MouseEvent) => {
      const el = e.target instanceof Element ? e.target.closest("a[href]") : null;
      if (!el) return;
      const href = el.getAttribute("href") || "";
      if (!href.startsWith("https://")) return;
      if (!isTauri()) return;
      e.preventDefault();
      void openExternal(href);
    };
    window.addEventListener("click", onClick, true);
    window.addEventListener("auxclick", onClick, true);
    return () => {
      window.removeEventListener("click", onClick, true);
      window.removeEventListener("auxclick", onClick, true);
    };
  }, []);
}
