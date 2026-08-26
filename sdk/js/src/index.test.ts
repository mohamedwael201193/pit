import { attention, canExecute, canPostExchange, canSign, explorer, refuseAuthorize, refuseSessionExport, refuseUnsignedPost } from "./index.ts";

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
try {
  refuseAuthorize();
  throw new Error("authorize");
} catch (e) {
  if (!(e instanceof Error) || e.message !== "authorize_denied") throw e;
}
assert(canPostExchange === false, "post");
try {
  refuseUnsignedPost();
  throw new Error("unsigned");
} catch (e) {
  if (!(e instanceof Error) || e.message !== "exchange_unsigned") throw e;
}
