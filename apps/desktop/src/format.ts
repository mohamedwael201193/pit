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

export function nearestVenueMin(
  coins: Array<{ minNotional?: number; eligible?: boolean; policyEligible?: boolean; executionFeasible?: boolean }>,
): number {
  const actionable = coins.filter((c) => c.executionFeasible);
  const policy = coins.filter((c) => c.policyEligible || c.eligible);
  const src = actionable.length ? actionable : policy.length ? policy : coins;
  const mins = src.map((c) => c.minNotional).filter((n): n is number => typeof n === "number" && n > 0);
  if (!mins.length) return 10;
  return Math.min(...mins);
}

export function nearestPolicyClip(
  coins: Array<{ policyClip?: number; eligible?: boolean; policyEligible?: boolean }>,
): number {
  const policy = coins.filter((c) => c.policyEligible || c.eligible);
  const src = policy.length ? policy : coins;
  const clips = src.map((c) => c.policyClip).filter((n): n is number => typeof n === "number" && n > 0);
  if (!clips.length) return 10;
  return clips[0];
}

export type AccountSizeGate = {
  have: number;
  min: number;
  shortfall: number;
  policyGap: number;
  canOpen: boolean;
  cta: "open" | "policy" | "fund";
  headline: string;
  detail: string;
};

export function accountSizeGate(opts: {
  have?: number;
  venueMin?: number;
  executable?: number;
  execGate?: string;
  execWhy?: string;
  policyClip?: number;
}): AccountSizeGate {
  const h = typeof opts.have === "number" && Number.isFinite(opts.have) ? opts.have : 0;
  const m = opts.venueMin && opts.venueMin > 0 ? opts.venueMin : 10;
  const execN = opts.executable ?? 0;
  const clip = opts.policyClip && opts.policyClip > 0 ? opts.policyClip : 10;
  const fundShort = Math.max(0, Number((m - h).toFixed(2)));
  const policyGap = Math.max(0, Number((m - clip).toFixed(2)));
  const namedTight = opts.execGate === "policy_clip_tight";
  const inferredTight = execN === 0 && h + 1e-9 >= m && h > 0 && clip + 1e-9 < m;
  if (execN > 0) {
    return {
      have: h,
      min: m,
      shortfall: 0,
      policyGap: 0,
      canOpen: true,
      cta: "open",
      headline: `${execN} book${execN === 1 ? "" : "s"} can open`,
      detail: `This account has ${compactUsd(h)}. Host sized executable books without inventing size.`,
    };
  }
  if (namedTight || inferredTight) {
    return {
      have: h,
      min: m,
      shortfall: 0,
      policyGap,
      canOpen: false,
      cta: "policy",
      headline: "Policy cap is too tight",
      detail:
        opts.execWhy ||
        `Policy cap is ${compactUsd(policyGap)} too tight for a ${compactUsd(m)} venue minimum. This account has ${compactUsd(h)}. Raise max trade on Security, preview, then pin. PIT will not invent size.`,
    };
  }
  return {
    have: h,
    min: m,
    shortfall: fundShort,
    policyGap: 0,
    canOpen: false,
    cta: "fund",
    headline: "Nothing can open yet",
    detail: `Hyperliquid needs ${compactUsd(m)} to open a position. This account has ${compactUsd(h)}${fundShort > 0 ? ` — ${compactUsd(fundShort)} short of that floor` : ""}. PIT will not invent size.`,
  };
}

export function marketSizeGate(
  coin: string,
  have?: number,
  min?: number,
  feasible?: boolean,
  execGate?: string,
  policyClip?: number,
) {
  const m = min && min > 0 ? min : 10;
  const h = typeof have === "number" && Number.isFinite(have) ? have : 0;
  const clip = policyClip && policyClip > 0 ? policyClip : 10;
  const shortfall = Math.max(0, Number((m - h).toFixed(2)));
  const policyGap = Math.max(0, Number((m - clip).toFixed(2)));
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
  if (execGate === "policy_clip_tight" || clip + 1e-9 < m) {
    return {
      chip: `Policy ${compactUsd(policyGap)} tight`,
      detail: `Policy cap is ${compactUsd(policyGap)} too tight for ${coin} min ${compactUsd(m)}. This account has ${compactUsd(h)}. Raise max trade on Security, preview, then pin. PIT will not invent size.`,
    };
  }
  return {
    chip: "Blocked",
    detail: `${coin} is not executable for this account right now. PIT will not invent size.`,
  };
}
