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
