export const canSign = false;
export const canExecute = false;
export const canReadAuthFile = false;

export type Network = "mainnet" | "testnet";

export { HEALTH_DEFAULT, COMPANION_DEFAULT } from "./constants.js";
export { getJson, publicHealth, publicWatch, publicRelease, companionHealth, companionStatus } from "./client.js";

export function explorer(network: Network): string {
  return network === "testnet" ? "https://chainscan-galileo.0g.ai" : "https://chainscan.0g.ai";
}

export function attention(count: number): string {
  if (count <= 0) {
    return "No opportunities match your policy.";
  }
  return `${count} opportunities match your policy.`;
}

export function refuseSessionExport(): never {
  throw new Error("session_export_denied");
}

export { refuseAuthorize } from "./authorize.js";
export { canPostExchange, refuseUnsignedPost } from "./post.js";
export { refuseArm, canArm } from "./mission.js";
