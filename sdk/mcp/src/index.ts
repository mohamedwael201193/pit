#!/usr/bin/env node
import { createInterface } from "node:readline";
import { callTool, listTools, PROTOCOL } from "./tools.js";

type Rpc = { jsonrpc?: string; id?: unknown; method?: string; params?: Record<string, unknown> };

function reply(id: unknown, result: unknown) {
  process.stdout.write(JSON.stringify({ jsonrpc: "2.0", id, result }) + "\n");
}

function fail(id: unknown, message: string) {
  process.stdout.write(JSON.stringify({ jsonrpc: "2.0", id, error: { code: -32000, message } }) + "\n");
}

async function handle(msg: Rpc) {
  const id = msg.id;
  const method = String(msg.method || "");
  if (method === "initialize") {
    reply(id, {
      protocolVersion: PROTOCOL,
      capabilities: { tools: {} },
      serverInfo: { name: "pit-mcp", version: "0.9.12" },
    });
    return;
  }
  if (method === "notifications/initialized" || method === "initialized") {
    return;
  }
  if (method === "tools/list") {
    reply(id, { tools: listTools() });
    return;
  }
  if (method === "tools/call") {
    const name = String(msg.params?.name || "");
    const args = (msg.params?.arguments || {}) as Record<string, unknown>;
    const out = await callTool(name, args);
    reply(id, out);
    return;
  }
  if (method === "ping") {
    reply(id, {});
    return;
  }
  fail(id, "method_not_found");
}

const rl = createInterface({ input: process.stdin });
rl.on("line", (line) => {
  const raw = line.trim();
  if (!raw) return;
  let msg: Rpc;
  try {
    msg = JSON.parse(raw) as Rpc;
  } catch {
    return;
  }
  void handle(msg).catch((err: unknown) => {
    fail(msg.id, err instanceof Error ? err.message : "mcp_error");
  });
});
