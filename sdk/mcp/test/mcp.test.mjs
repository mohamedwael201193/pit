import test from "node:test";
import assert from "node:assert/strict";
import { callTool, listTools } from "../dist/tools.js";

test("tools are named and authorize is absent", () => {
  const names = listTools().map((t) => t.name);
  assert.ok(names.includes("pit_health"));
  assert.ok(names.includes("pit_watch"));
  assert.ok(!names.includes("authorize"));
  assert.ok(!names.includes("pit_authorize"));
});

test("authorize tool names are denied", async () => {
  const denied = await callTool("authorize", {});
  assert.equal(denied.isError, true);
  assert.match(JSON.stringify(denied), /authorize_denied/);
  const order = await callTool("pit_order", {});
  assert.equal(order.isError, true);
});

test("security card never signs", async () => {
  const out = await callTool("pit_security", {});
  const body = JSON.parse(out.content[0].text);
  assert.equal(body.authorize, false);
  assert.equal(body.sign, false);
  assert.equal(body.trade, false);
  assert.equal(body.arm, false);
});
