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

export function healthOrigin(): string {
  const raw = (import.meta.env.VITE_HEALTH_URL as string | undefined) || HEALTH_DEFAULT;
  return String(raw).replace(/\/$/, "");
}

/** First-party installer. Health 302s to the latest NSIS asset, not the GitHub Releases HTML page. */
export function windowsInstallerUrl(): string {
  return `${healthOrigin()}/windows`;
}

export function windowsChecksumsUrl(): string {
  return `${healthOrigin()}/checksums`;
}

export const DESK_ID_CONTRACT = "0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35" as const;
export const DESK_TOKEN_ID = 1n;
export const IDENTITY_8004 = "0x8004A169FB4a3325136EB29fA0ceB6D2e539a432" as const;
export const REPUTATION_8004 = "0x8004BAa17C55a88189AE136b182e5fdA19dE9b63" as const;
export const AGENT_8004_ID = 3489333n;
export const REPUTATION_CLIENT = "0xaAE3EAC0d6665832fe0E5036d61CE2DBC6ECAC2a" as const;
export const AGENT_CARD_URL = "https://pit0g.vercel.app/.well-known/agent-card.json" as const;

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
  note: "Older Hyperliquid fill on this desk. Not flattened. Not a live public mission stream.",
};

export const VERIFIED_FILL = {
  id: "recorded-hype",
  market: "HYPE",
  oid: "531667200134",
  sz: "0.16",
  px: "80.909",
  job: "4a1d45ec-8c3f-4883-a162-19739accb9cf",
  researchTx: "0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b",
  orderTx: "0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb",
  researchRoot: "0x9fd42770545ecaacbfff12e3ef7a537b564e31c9ef5515b3a820fd276c22f72e",
  orderRoot: "0x8c94ec8e643c90fe69276ff20f50a0bc3121f007d611e10e6ab9f24d26f2ff66",
  kind: "RECORDED" as const,
  note: "Same-desk READY → TRADE NOW → FILLED. Not a live public mission stream. This site does not fetch another account's book.",
} as const;

/** Recorded Aristotle storage filings from this desk. Roots from Flow event topics. */
export const STORAGE_PROOFS = [
  {
    label: "HYPE research · job 4a1d45ec",
    tx: "0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b",
    root: "0x9fd42770545ecaacbfff12e3ef7a537b564e31c9ef5515b3a820fd276c22f72e",
  },
  {
    label: "HYPE order evidence · job 4a1d45ec",
    tx: "0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb",
    root: "0x8c94ec8e643c90fe69276ff20f50a0bc3121f007d611e10e6ab9f24d26f2ff66",
  },
  {
    label: "ETH research receipt",
    tx: "0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80",
    root: "0x07238aa66936340f7ea9fa59f279a8e2313b0bb839699c805b91cb30ccb7741d",
  },
  {
    label: "BTC research receipt",
    tx: "0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a",
    root: "0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2",
  },
  {
    label: "S6 encrypted storage roundtrip",
    tx: "",
    root: "0x3b4b3b772aae7195109ded219ca861da7eb3ca51776538e3486f7084b4ef193a",
  },
] as const;

export const RESEARCH_TXS = [
  { label: "HYPE hunt 9761cbd5", tx: "0x266c45cbd35cb8b9e856d7f3c850e5ce72d34fb33251bba616345e34cd04cb78" },
  { label: "Long hunt HYPE 7d28f3e3", tx: "0x6009ede35278fc6157507792388b87e2f0c7173494a32095fb0692bc65c77ff4" },
  { label: "HYPE none-hypothesis 3c9f2e96", tx: "0x30df71b929e05a4feca6d4683bbe86af97750b70807a28957bcc54e2d99aa4ed" },
  { label: "DOGE job b4ed73ce", tx: "0x28f0f7474760ec88c8c2a76f9959e136756eb5dd8ccfd530eb43d38c10f7277c" },
  { label: "HYPE storage 8c8b78e8", tx: "0x8c8b78e8add46c79983d344ac571bcb8e6fd1d6c2ae072add00147f2ede1151d" },
  { label: "Research 0xcc02a780", tx: "0xcc02a780b12ed2a884d3aa845f486acb89c60f1e8c306f0773e147f5311b4438" },
  { label: "Research 0xd682aa45", tx: "0xd682aa45aea64a26d1ab7a18d9867260a38502b086b9730010a394011ef6114c" },
  { label: "Research 0x2a7a5838", tx: "0x2a7a58381ef4507174a777fb2f9a65d826d9988ce22610fc16b4d9e1fcd54b9d" },
  { label: "Research 0x2045c98a", tx: "0x2045c98a69aae505ee5be36eaa1cf05c5d93c2662d90b5d7b07dc8452d537711" },
  { label: "BTC challenger_killed", tx: "0x6abe43772f1b953e2c6debec31dba1d64b77a7f8c3b6f83cf950f18f11e263e4" },
  { label: "SOL no_side", tx: "0x7e7f85aaf4aacd29129b8697cbc5de7e8f6d56745754897807a262e2d31b21ef" },
  { label: "ETH job 78617f6c", tx: "0xdf4f8f95cbee81f99402754455915635bbc3f4623861318f5fc171da631f8ae0" },
] as const;

export const HISTORICAL_TEE_SIGNER = "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9";
export const EXPECTED_TEE_SIGNER = "0x041a09E5bEF30fd776D66Bb892d18B97637C7C7c";

export const CAPITAL_PRESETS = [10, 25, 50, 100, 500, 1000] as const;

export const PUBLIC_NAV = [
  { to: "/", label: "PIT", end: true },
  { to: "/radar", label: "Radar", end: false },
  { to: "/autonomy", label: "Autonomy", end: false },
  { to: "/missions", label: "Missions", end: false },
  { to: "/proof", label: "Proof", end: false },
  { to: "/agent", label: "Agent", end: false },
  { to: "/capital", label: "Capital", end: false },
  { to: "/how-it-works", label: "How it works", end: false },
] as const;
