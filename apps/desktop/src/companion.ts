export const COMPANION = "http://127.0.0.1:17373";

export type LocalStatus = {
  sessionAlive?: boolean;
  agent?: string;
  network?: string;
  kill?: boolean;
  sign?: boolean;
  trade?: boolean;
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

export async function localStatus(): Promise<LocalStatus | null> {
  try {
    const r = await fetch(`${COMPANION}/local/status`);
    if (!r.ok) return null;
    const body = (await r.json()) as LocalStatus;
    if (body.sign || body.trade) return null;
    return body;
  } catch {
    return null;
  }
}

export async function pairCode(): Promise<PairCode | null> {
  try {
    const r = await fetch(`${COMPANION}/local/code`);
    if (!r.ok) return null;
    const body = (await r.json()) as PairCode;
    if (body.sign || body.trade) return null;
    return body;
  } catch {
    return null;
  }
}

export async function doctor(): Promise<DoctorCheck[]> {
  try {
    const r = await fetch(`${COMPANION}/local/doctor`);
    if (!r.ok) return [];
    const body = (await r.json()) as { checks?: DoctorCheck[]; sign?: boolean };
    if (body.sign) return [];
    return Array.isArray(body.checks) ? body.checks : [];
  } catch {
    return [];
  }
}

export function prettyCode(code: string) {
  const raw = code.replace(/[^A-Za-z0-9]/g, "").toUpperCase();
  if (raw.length !== 8) return raw;
  return `${raw.slice(0, 4)}-${raw.slice(4)}`;
}
