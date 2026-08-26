import { attention, canExecute, canSign, explorer, refuseSessionExport } from "./index.ts";

function assert(cond: unknown, msg: string) {
  if (!cond) throw new Error(msg);
}

assert(canSign === false, "sign");
assert(canExecute === false, "exec");
assert(attention(0) === "No opportunities match your policy.", "empty");
assert(explorer("mainnet").includes("chainscan.0g.ai"), "main");
assert(explorer("testnet").includes("galileo"), "test");
try {
  refuseSessionExport();
  throw new Error("export");
} catch (e) {
  if (!(e instanceof Error) || e.message !== "session_export_denied") throw e;
}
