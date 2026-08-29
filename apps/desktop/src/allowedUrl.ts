const ALLOWED_HOSTS = new Set([
  "pit0g.vercel.app",
  "app.hyperliquid.xyz",
  "app.hyperliquid-testnet.xyz",
  "pc.0g.ai",
  "0g.ai",
  "www.0g.ai",
  "docs.0g.ai",
  "chainscan.0g.ai",
  "chainscan-galileo.0g.ai",
  "hyperliquid.info",
  "github.com",
]);

export function isAllowedHttps(url: string): boolean {
  const u = String(url || "").trim();
  if (!u.startsWith("https://") || u.length > 2048) return false;
  if (/[\s\\"<>`|]/.test(u)) return false;
  let parsed: URL;
  try {
    parsed = new URL(u);
  } catch {
    return false;
  }
  if (parsed.protocol !== "https:") return false;
  if (parsed.username || parsed.password) return false;
  if (parsed.port) return false;
  const host = parsed.hostname.toLowerCase();
  if (!ALLOWED_HOSTS.has(host)) return false;
  if (host === "github.com" && !parsed.pathname.toLowerCase().startsWith("/mohamedwael201193/pit")) {
    return false;
  }
  return true;
}
