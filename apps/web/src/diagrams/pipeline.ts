export const PIPELINE = [
  "MARKET",
  "PRIVATE",
  "SEALED",
  "RESEARCH",
  "CHALLENGE",
  "RISK",
  "POLICY",
  "AUTHORIZE",
  "EXECUTE",
  "PROVE",
  "LEARN",
] as const;

export type PipelineBeat = (typeof PIPELINE)[number];
