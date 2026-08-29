import { LINKS, explorerAddress, explorerTx, hyperliquidAPI, hyperliquidApp, hyperliquidTrade } from "../src/links";

export function assertOfficialLinks() {
  if (LINKS.pair !== "https://pit0g.vercel.app/pair") throw new Error("pair");
  if (LINKS.app !== "https://pit0g.vercel.app/app") throw new Error("app");
  if (LINKS.pcAdvanced !== "https://pc.0g.ai/sdk/dashboard/funds") throw new Error("pc");
  if (LINKS.og !== "https://0g.ai") throw new Error("og");
  if (LINKS.hl !== "https://app.hyperliquid.xyz") throw new Error("hl");
  if (LINKS.hlAPI !== "https://app.hyperliquid.xyz/API") throw new Error("hl api");
  if (LINKS.hlTestnet !== "https://app.hyperliquid-testnet.xyz") throw new Error("hl testnet");
  if (LINKS.hlAPITestnet !== "https://app.hyperliquid-testnet.xyz/API") throw new Error("hl api testnet");
  if (LINKS.releases !== "https://github.com/mohamedwael201193/pit/releases/latest") throw new Error("releases");
  if (LINKS.explorer !== "https://chainscan.0g.ai") throw new Error("explorer");
  if (hyperliquidAPI("mainnet") !== LINKS.hlAPI) throw new Error("hl api net");
  if (hyperliquidAPI("testnet") !== LINKS.hlAPITestnet) throw new Error("hl api testnet net");
  if (hyperliquidApp("mainnet") !== LINKS.hl) throw new Error("hl app");
  if (hyperliquidApp("testnet") !== LINKS.hlTestnet) throw new Error("hl testnet app");
  if (hyperliquidTrade("mainnet", "BTC") !== "https://app.hyperliquid.xyz/trade/BTC") throw new Error("hl trade");
  if (!explorerAddress("0xabc").startsWith("https://chainscan.0g.ai/address/")) throw new Error("explorer addr");
  if (!explorerTx("0xabc", "mainnet").startsWith("https://chainscan.0g.ai/tx/")) throw new Error("explorer tx");
}
