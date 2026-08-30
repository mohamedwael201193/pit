import type { PublicCoin } from "./types";

export const HYPERLIQUID_PERP_FLOOR = 10;

export function perpMinNotionalUSD(mark: number, szDecimals: number): number {
  if (!Number.isFinite(mark) || mark <= 0) return HYPERLIQUID_PERP_FLOOR;
  let d = szDecimals;
  if (d < 0) d = 0;
  if (d > 8) d = 8;
  const pow = 10 ** d;
  const sz = Math.ceil((HYPERLIQUID_PERP_FLOOR / mark) * pow) / pow;
  if (sz <= 0) return HYPERLIQUID_PERP_FLOOR;
  const got = sz * mark;
  return got < HYPERLIQUID_PERP_FLOOR ? HYPERLIQUID_PERP_FLOOR : got;
}

export function coinMin(c: Pick<PublicCoin, "mark" | "minNotional" | "szDecimals">): number {
  if (typeof c.minNotional === "number" && c.minNotional >= HYPERLIQUID_PERP_FLOOR) return c.minNotional;
  return perpMinNotionalUSD(c.mark, c.szDecimals ?? 0);
}
