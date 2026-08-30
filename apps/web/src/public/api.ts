import { HEALTH_DEFAULT } from "./facts";
import type { HealthView, WatchView } from "./types";

export function healthBase(): string {
  const raw = import.meta.env.VITE_HEALTH_URL || HEALTH_DEFAULT;
  return raw.replace(/\/$/, "");
}

export class UnsafeWatchError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "UnsafeWatchError";
  }
}

export async function fetchHealth(signal?: AbortSignal): Promise<HealthView> {
  const r = await fetch(`${healthBase()}/health`, { signal });
  if (!r.ok) throw new Error(`health_${r.status}`);
  const body = (await r.json()) as HealthView;
  if (body.sign) throw new UnsafeWatchError("health_sign");
  return body;
}

export async function fetchWatch(network = "mainnet", signal?: AbortSignal): Promise<WatchView> {
  const r = await fetch(`${healthBase()}/watch?network=${network}`, { signal });
  if (!r.ok) throw new Error(`watch_${r.status}`);
  const body = (await r.json()) as WatchView;
  if (body.trade || body.sign) {
    throw new UnsafeWatchError("watch_refused_trade_or_sign");
  }
  return {
    ...body,
    coins: Array.isArray(body.coins) ? body.coins : [],
    count: typeof body.count === "number" ? body.count : 0,
    scanned: typeof body.scanned === "number" ? body.scanned : 0,
  };
}
