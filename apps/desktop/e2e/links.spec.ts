import { LINKS, explorerAddress, explorerTx, hyperliquidAPI, hyperliquidApp, hyperliquidTrade } from "../src/links";
import { isAllowedHttps } from "../src/allowedUrl";

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
  if (!isAllowedHttps(LINKS.pair)) throw new Error("allow pair");
  if (!isAllowedHttps(LINKS.pcAdvanced)) throw new Error("allow pc");
  if (!isAllowedHttps(hyperliquidTrade("mainnet", "ETH"))) throw new Error("allow trade");
  if (!isAllowedHttps("https://chainscan.0g.ai/tx/0xabc")) throw new Error("allow tx");
  if (!isAllowedHttps(LINKS.releases)) throw new Error("allow releases");
  if (isAllowedHttps("http://app.hyperliquid.xyz")) throw new Error("http denied");
  if (isAllowedHttps("https://evil.example")) throw new Error("host denied");
  if (isAllowedHttps("https://github.com/other/repo")) throw new Error("github denied");
}
