export type PublicCoin = {
  coin: string;
  venue?: string;
  reason?: string;
  why?: string;
  trend?: string;
  rank?: number;
  freshness?: string;
  mark: number;
  oracle?: number;
  funding?: number;
  openInterest?: number;
  volume?: number;
  szDecimals?: number;
  timestamp?: string;
  provenance?: string;
  source?: string;
  network?: string;
  eligible?: boolean;
  policyFit?: string;
  researchEligible?: boolean;
  policyEligible?: boolean;
  executionFeasible?: boolean;
  block?: string;
  execGate?: string;
  execWhy?: string;
  minNotional?: number;
  riskFlags?: string[];
};

export type WatchView = {
  ok?: boolean;
  sign?: boolean;
  trade?: boolean;
  count?: number;
  scanned?: number;
  copy?: string;
  coins?: PublicCoin[];
  best?: PublicCoin | null;
  bestWhy?: string;
  source?: string;
  network?: string;
  version?: string;
};

export type HealthView = {
  ok?: boolean;
  sign?: boolean;
  service?: string;
  time?: string;
  version?: string;
  requestId?: string;
};

export type DesktopProbe = {
  present: boolean;
  version?: string;
  refused: boolean;
};

export type SimKind = "TRADE" | "WAIT" | "HOLD" | "LIQUIDITY";

export type SimRow = {
  coin: string;
  mark: number;
  minNotional: number;
  kind: SimKind;
  why: string;
};

export type EvidenceKind = "LIVE" | "HISTORICAL" | "REPLAY" | "ABSENT";
