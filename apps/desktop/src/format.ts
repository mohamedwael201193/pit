export function compactNum(n: number | undefined | null): string {
  if (n == null || !Number.isFinite(n)) return "—";
  const a = Math.abs(n);
  if (a >= 1e9) return `${(n / 1e9).toFixed(2)}B`;
  if (a >= 1e6) return `${(n / 1e6).toFixed(2)}M`;
  if (a >= 1000) return n.toLocaleString(undefined, { maximumFractionDigits: 0 });
  if (a >= 1) return n.toLocaleString(undefined, { maximumFractionDigits: 4 });
  return n.toLocaleString(undefined, { maximumFractionDigits: 6 });
}

export function compactUsd(n: number | undefined | null): string {
  if (n == null || !Number.isFinite(n)) return "—";
  const a = Math.abs(n);
  if (a >= 1e9) return `$${(n / 1e9).toFixed(2)}B`;
  if (a >= 1e6) return `$${(n / 1e6).toFixed(2)}M`;
  if (a >= 1000) return `$${n.toLocaleString(undefined, { maximumFractionDigits: 0 })}`;
  return `$${n.toLocaleString(undefined, { maximumFractionDigits: 2 })}`;
}

export function pctFunding(n: number | undefined | null): string {
  if (n == null || !Number.isFinite(n)) return "—";
  return `${(n * 100).toFixed(4)}%`;
}

export function remainLabel(deadlineUnix?: number, nowUnix?: number): string {
  const now = nowUnix && nowUnix > 0 ? nowUnix : Math.floor(Date.now() / 1000);
  const until = deadlineUnix || 0;
  if (until <= 0) return "—";
  const s = until - now;
  if (s <= 0) return "expired";
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  if (h >= 24) return `${Math.floor(h / 24)}d ${h % 24}h`;
  if (h > 0) return `${h}h ${m}m`;
  return `${m}m`;
}

export function elapsedLabel(startUnix?: number, nowUnix?: number): string {
  const now = nowUnix && nowUnix > 0 ? nowUnix : Math.floor(Date.now() / 1000);
  const start = startUnix || 0;
  if (start <= 0) return "—";
  const s = Math.max(0, now - start);
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const sec = s % 60;
  if (h > 0) return `${h}h ${m}m`;
  if (m > 0) return `${m}m ${sec}s`;
  return `${sec}s`;
}

export function nextScanLabel(nextUnix?: number, nowUnix?: number): string {
  const now = nowUnix && nowUnix > 0 ? nowUnix : Math.floor(Date.now() / 1000);
  const next = nextUnix || 0;
  if (next <= 0) return "due now";
  if (next <= now) return "due now";
  return new Date(next * 1000).toLocaleTimeString();
}

export function powerSourceLabel(src?: string) {
  if (!src) return "";
  return src.replaceAll("_", " ");
}

export function nearestVenueMin(coins: Array<{ minNotional?: number }>): number {
  const mins = coins.map((c) => c.minNotional).filter((n): n is number => typeof n === "number" && n > 0);
  if (!mins.length) return 10;
  return Math.min(...mins);
}

export function accountSizeGate(have?: number, min = 10, executable = 0) {
  const h = typeof have === "number" && Number.isFinite(have) ? have : 0;
  const m = min > 0 ? min : 10;
  const shortfall = Math.max(0, Number((m - h).toFixed(2)));
  const canOpen = executable > 0 && h >= m;
  if (canOpen) {
    return {
      have: h,
      min: m,
      shortfall: 0,
      canOpen: true as const,
      headline: `${executable} book${executable === 1 ? "" : "s"} can open`,
      detail: `This account has ${compactUsd(h)}. Hyperliquid's open minimum is ${compactUsd(m)}. Host sized executable books without inventing size.`,
    };
  }
  return {
    have: h,
    min: m,
    shortfall,
    canOpen: false as const,
    headline: "Nothing can open yet",
    detail: `Hyperliquid needs ${compactUsd(m)} to open a position. This account has ${compactUsd(h)}${shortfall > 0 ? ` — ${compactUsd(shortfall)} short of that floor` : ""}. PIT will not invent size.`,
  };
}

export function marketSizeGate(coin: string, have?: number, min?: number, feasible?: boolean) {
  const m = min && min > 0 ? min : 10;
  const h = typeof have === "number" && Number.isFinite(have) ? have : 0;
  const shortfall = Math.max(0, Number((m - h).toFixed(2)));
  if (feasible) {
    return {
      chip: "Can open",
      detail: `${coin} can be sized with ${compactUsd(h)} against a ${compactUsd(m)} floor.`,
    };
  }
  if (h > 0 && h < m) {
    return {
      chip: `${compactUsd(shortfall)} short`,
      detail: `${coin} needs ${compactUsd(m)} to open on Hyperliquid. This account has ${compactUsd(h)} — ${compactUsd(shortfall)} short. PIT will not invent size.`,
    };
  }
  if (h <= 0) {
    return {
      chip: "No margin",
      detail: `No available venue margin. ${coin} still needs ${compactUsd(m)} to open. PIT will not invent size.`,
    };
  }
  return {
    chip: "Blocked",
    detail: `${coin} is not executable for this account right now. PIT will not invent size.`,
  };
}
