export const LINKS = {
  pair: "https://pit0g.vercel.app/pair",
  app: "https://pit0g.vercel.app/app",
  pcAdvanced: "https://pc.0g.ai/sdk/dashboard/funds",
  hl: "https://app.hyperliquid.xyz",
  hlTestnet: "https://app.hyperliquid-testnet.xyz",
  hlAPI: "https://app.hyperliquid.xyz/API",
  hlAPITestnet: "https://app.hyperliquid-testnet.xyz/API",
  releases: "https://github.com/mohamedwael201193/pit/releases/latest",
  explorer: "https://chainscan.0g.ai",
} as const;

export function explorerAddress(addr: string) {
  return `${LINKS.explorer}/address/${addr}`;
}

export function hyperliquidAPI(net: string) {
  return net === "testnet" ? LINKS.hlAPITestnet : LINKS.hlAPI;
}

export function hyperliquidApp(net: string) {
  return net === "testnet" ? LINKS.hlTestnet : LINKS.hl;
}
