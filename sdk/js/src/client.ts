import { COMPANION_DEFAULT, HEALTH_DEFAULT } from "./constants.js";

export async function getJson(url: string): Promise<Record<string, unknown>> {
  const r = await fetch(url, { method: "GET" });
  if (!r.ok) {
    throw new Error(`http_${r.status}`);
  }
  const body = (await r.json()) as Record<string, unknown>;
  if (body && typeof body === "object") {
    if (body.sign === true) throw new Error("sign_refused");
    if (body.trade === true) throw new Error("trade_refused");
  }
  return body;
}

function trim(base: string): string {
  return base.replace(/\/$/, "");
}

/** Public health. Never holds a session. */
export async function publicHealth(base = HEALTH_DEFAULT): Promise<Record<string, unknown>> {
  return getJson(`${trim(base)}/health`);
}

/** Public Hyperliquid watch. Never places an order. */
export async function publicWatch(network = "mainnet", base = HEALTH_DEFAULT): Promise<Record<string, unknown>> {
  return getJson(`${trim(base)}/watch?network=${encodeURIComponent(network)}`);
}

export async function publicRelease(base = HEALTH_DEFAULT): Promise<Record<string, unknown>> {
  return getJson(`${trim(base)}/release`);
}

/** Loopback companion health. Does not authorize. */
export async function companionHealth(base = COMPANION_DEFAULT): Promise<Record<string, unknown>> {
  return getJson(`${trim(base)}/health`);
}

/** Loopback desk status. GET only. Cannot AUTHORIZE. */
export async function companionStatus(base = COMPANION_DEFAULT): Promise<Record<string, unknown>> {
  return getJson(`${trim(base)}/local/status`);
}
