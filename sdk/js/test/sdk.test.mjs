import test from "node:test";
import assert from "node:assert/strict";
import {
  attention,
  canArm,
  canExecute,
  canPostExchange,
  canSign,
  explorer,
  refuseArm,
  refuseAuthorize,
  refuseSessionExport,
  refuseUnsignedPost,
} from "../dist/index.js";

test("never signs or executes", () => {
  assert.equal(canSign, false);
  assert.equal(canExecute, false);
  assert.equal(canPostExchange, false);
  assert.equal(canArm, false);
  assert.equal(attention(0), "No opportunities match your policy.");
  assert.ok(explorer("mainnet").includes("chainscan.0g.ai"));
  assert.ok(explorer("testnet").includes("galileo"));
});

test("refuses authorize export arm post", () => {
  assert.throws(() => refuseSessionExport(), /session_export_denied/);
  assert.throws(() => refuseAuthorize(), /authorize_denied/);
  assert.throws(() => refuseArm(), /mission_arm_denied/);
  assert.throws(() => refuseUnsignedPost(), /exchange_unsigned/);
});
