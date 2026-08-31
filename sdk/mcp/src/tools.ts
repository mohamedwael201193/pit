export const PROTOCOL = "2024-11-05";
export const HEALTH = "https://pit-health.onrender.com";
export const COMPANION = "http://127.0.0.1:17373";

const FORBIDDEN = new Set([
  "authorize",
  "order",
  "cancel",
  "export_session",
  "arm",
  "pit_authorize",
  "pit_order",
  "pit_export",
]);

export function listTools() {
  return [
    { name: "pit_health", description: "Public PIT health. Never signs.", inputSchema: { type: "object", properties: {} } },
    { name: "pit_watch", description: "Public Hyperliquid watch. Never trades.", inputSchema: { type: "object", properties: { network: { type: "string" } } } },
    { name: "pit_release", description: "Latest installer metadata. SHA and filename only.", inputSchema: { type: "object", properties: {} } },
    { name: "pit_companion_health", description: "Loopback companion health if PIT Desktop is running.", inputSchema: { type: "object", properties: {} } },
    { name: "pit_status", description: "Loopback desk status. GET only.", inputSchema: { type: "object", properties: {} } },
    { name: "pit_activity", description: "Loopback activity log. GET only.", inputSchema: { type: "object", properties: {} } },
    { name: "pit_positions", description: "Loopback Hyperliquid positions for the bound wallet.", inputSchema: { type: "object", properties: {} } },
    { name: "pit_research_status", description: "Current research job if one is running. Cannot start research.", inputSchema: { type: "object", properties: {} } },
    { name: "pit_proofs", description: "This-machine 0G proof index. GET only.", inputSchema: { type: "object", properties: {} } },
    { name: "pit_security", description: "Capability denial card.", inputSchema: { type: "object", properties: {} } },
  ];
}

async function getJson(url: string): Promise<Record<string, unknown>> {
  const r = await fetch(url, { method: "GET" });
  if (!r.ok) throw new Error(`http_${r.status}`);
  const body = (await r.json()) as Record<string, unknown>;
  if (body?.sign === true) throw new Error("sign_refused");
  if (body?.trade === true) throw new Error("trade_refused");
  return body;
}

function text(obj: unknown, isError = false) {
  return { content: [{ type: "text", text: JSON.stringify(obj) }], isError };
}

export async function callTool(name: string, args: Record<string, unknown>) {
  if (FORBIDDEN.has(name) || name.includes("authorize") || name.includes("order")) {
    return text({ ok: false, error: "authorize_denied", sign: false, trade: false }, true);
  }
  try {
    switch (name) {
      case "pit_health":
        return text(await getJson(`${HEALTH}/health`));
      case "pit_watch": {
        const net = typeof args.network === "string" ? args.network : "mainnet";
        return text(await getJson(`${HEALTH}/watch?network=${encodeURIComponent(net)}`));
      }
      case "pit_release":
        return text(await getJson(`${HEALTH}/release`));
      case "pit_companion_health":
        return text(await getJson(`${COMPANION}/health`));
      case "pit_status":
        return text(await getJson(`${COMPANION}/local/status`));
      case "pit_activity":
        return text(await getJson(`${COMPANION}/local/activity`));
      case "pit_positions":
        return text(await getJson(`${COMPANION}/local/positions`));
      case "pit_research_status":
        return text(await getJson(`${COMPANION}/local/research/status`));
      case "pit_proofs":
        return text(await getJson(`${COMPANION}/local/proofs`));
      case "pit_security":
        return text({
          authorize: false,
          sign: false,
          trade: false,
          arm: false,
          export_session: false,
          copy: "PIT MCP is read-only. AUTHORIZE and orders stay on PIT Desktop.",
        });
      default:
        return text({ ok: false, error: "unknown_tool", sign: false, trade: false }, true);
    }
  } catch (err) {
    return text({ ok: false, error: err instanceof Error ? err.message : "mcp_error", sign: false, trade: false }, true);
  }
}
