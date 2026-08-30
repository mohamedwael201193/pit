/** Public-safe constants. Never put keys, Direct tokens, or private prompts here. */

export const RELEASES = "https://github.com/mohamedwael201193/pit/releases/latest";
export const REPO = "https://github.com/mohamedwael201193/pit";
export const HEALTH_DEFAULT = "https://pit-health.onrender.com";
export const COMPANION = "http://127.0.0.1:17373";
export const ARISTOTLE_RPC = "https://evmrpc.0g.ai";
export const ARISTOTLE_EXPLORER = "https://chainscan.0g.ai";
export const ARISTOTLE_ID = 16661;
export const HL_INFO = "https://api.hyperliquid.xyz/info";
export const HL_APP = "https://app.hyperliquid.xyz";
export const HL_API = "https://app.hyperliquid.xyz/API";
export const PC_0G = "https://pc.0g.ai";
export const STORAGE_INDEXER = "https://indexer-storage-turbo.0g.ai";

export const DESK_ID_CONTRACT = "0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35" as const;
export const DESK_TOKEN_ID = 1n;
export const IDENTITY_8004 = "0x8004A169FB4a3325136EB29fA0ceB6D2e539a432" as const;
export const REPUTATION_8004 = "0x8004BAa17C55a88189AE136b182e5fdA19dE9b63" as const;
export const AGENT_8004_ID = 3489333n;

export const PIT_AGENT = {
  name: "PIT-4bbee556",
  address: "0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52",
} as const;

export const HISTORICAL_FILL = {
  id: "historical-eth",
  market: "ETH",
  oid: "529167222216",
  sz: "0.0041",
  px: "2489.7",
  kind: "HISTORICAL" as const,
  note: "Real Hyperliquid fill recorded on this desk. Not a live public mission stream. This site does not fetch another account's book.",
};

export const HISTORICAL_TEE_SIGNER = "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9";
export const EXPECTED_TEE_SIGNER = HISTORICAL_TEE_SIGNER;

export const CAPITAL_PRESETS = [10, 25, 50, 100, 500, 1000] as const;

export const PUBLIC_NAV = [
  { to: "/", label: "PIT", end: true },
  { to: "/radar", label: "Radar", end: false },
  { to: "/missions", label: "Missions", end: false },
  { to: "/proof", label: "Proof", end: false },
  { to: "/agent", label: "Agent", end: false },
  { to: "/capital", label: "Capital", end: false },
  { to: "/how-it-works", label: "How it works", end: false },
  { to: "/download", label: "Download", end: false },
] as const;
