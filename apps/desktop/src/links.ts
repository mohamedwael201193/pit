export const LINKS = {
  pair: "https://pit0g.vercel.app/pair",
  app: "https://pit0g.vercel.app/app",
  pcAdvanced: "https://pc.0g.ai/sdk/dashboard/funds",
  hlAPI: "https://app.hyperliquid.xyz/API",
  hlAPITestnet: "https://app.hyperliquid-testnet.xyz/API",
  releases: "https://github.com/mohamedwael201193/pit/releases/latest",
} as const;

export function hyperliquidAPI(net: string) {
  return net === "testnet" ? LINKS.hlAPITestnet : LINKS.hlAPI;
}
