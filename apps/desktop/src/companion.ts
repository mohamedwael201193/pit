export const COMPANION = "http://127.0.0.1:17373";

export type LocalStatus = {
  sessionAlive?: boolean;
  agent?: string;
  network?: string;
  kill?: boolean;
  sign?: boolean;
  trade?: boolean;
  version?: string;
  workspace?: string;
};

export type PairCode = {
  code?: string;
  expires?: string;
  sign?: boolean;
  trade?: boolean;
};

export type DoctorCheck = {
  name: string;
  ok: boolean;
  detail: string;
};

type TauriWindow = Window & {
  __TAURI_INTERNALS__?: { invoke: (cmd: string, args?: unknown) => Promise<unknown> };
  __TAURI__?: { core?: { invoke: (cmd: string, args?: unknown) => Promise<unknown> } };
};

function nativeInvoke<T>(cmd: string): Promise<T> | null {
  const w = window as TauriWindow;
  const inv = w.__TAURI_INTERNALS__?.invoke ?? w.__TAURI__?.core?.invoke;
  if (typeof inv !== "function") return null;
  return inv(cmd) as Promise<T>;
}

async function nativeJson<T extends object>(cmd: string): Promise<T | null> {
  const p = nativeInvoke<T>(cmd);
  if (!p) return null;
  try {
    return await p;
  } catch {
    return null;
  }
}

async function fetchJson<T extends object>(path: string): Promise<T | null> {
  try {
    const r = await fetch(`${COMPANION}${path}`);
    if (!r.ok) return null;
    return (await r.json()) as T;
  } catch {
    return null;
  }
}

function rejectSecrets<T extends { sign?: boolean; trade?: boolean }>(body: T | null): T | null {
  if (!body || body.sign || body.trade) return null;
  return body;
}

export async function wakeCompanion(): Promise<boolean> {
  const p = nativeInvoke<boolean>("ensure_companion");
  if (p) {
    try {
      return Boolean(await p);
    } catch {
      return false;
    }
  }
  const health = await fetchJson<{ ok?: boolean }>("/health");
  return Boolean(health);
}

export async function localStatus(): Promise<LocalStatus | null> {
  const native = rejectSecrets(await nativeJson<LocalStatus>("local_status"));
  if (native) return native;
  return rejectSecrets(await fetchJson<LocalStatus>("/local/status"));
}

export async function pairCode(): Promise<PairCode | null> {
  const native = rejectSecrets(await nativeJson<PairCode>("local_code"));
  if (native) return native;
  return rejectSecrets(await fetchJson<PairCode>("/local/code"));
}

export async function doctor(): Promise<DoctorCheck[]> {
  const native = await nativeJson<{ checks?: DoctorCheck[]; sign?: boolean }>("local_doctor");
  if (native && !native.sign) {
    return Array.isArray(native.checks) ? native.checks : [];
  }
  const fetched = await fetchJson<{ checks?: DoctorCheck[]; sign?: boolean }>("/local/doctor");
  if (!fetched || fetched.sign) return [];
  return Array.isArray(fetched.checks) ? fetched.checks : [];
}

export function prettyCode(code: string) {
  const raw = code.replace(/[^A-Za-z0-9]/g, "").toUpperCase();
  if (raw.length !== 8) return raw;
  return `${raw.slice(0, 4)}-${raw.slice(4)}`;
}

export function checkNamed(checks: DoctorCheck[], name: string): DoctorCheck | undefined {
  return checks.find((c) => c.name === name);
}
