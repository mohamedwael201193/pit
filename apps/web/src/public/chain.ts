import { defineChain, http, isHex, createPublicClient } from "viem";
import { ARISTOTLE_ID, ARISTOTLE_RPC, IDENTITY_8004 } from "./facts";

export const aristotleChain = defineChain({
  id: ARISTOTLE_ID,
  name: "0G Aristotle",
  nativeCurrency: { name: "0G", symbol: "0G", decimals: 18 },
  rpcUrls: { default: { http: [ARISTOTLE_RPC] } },
});

export function aristotleClient() {
  return createPublicClient({ chain: aristotleChain, transport: http(ARISTOTLE_RPC) });
}

export type TxRead =
  | { ok: true; hash: `0x${string}`; from: string; to: string | null; blockNumber: string }
  | { ok: false; reason: string };

export async function readAristotleTx(hash: string): Promise<TxRead> {
  const h = hash.trim();
  if (!isHex(h) || h.length !== 66) {
    return { ok: false, reason: "Need a 0x transaction hash, 32 bytes." };
  }
  try {
    const client = aristotleClient();
    const tx = await client.getTransaction({ hash: h as `0x${string}` });
    if (!tx) return { ok: false, reason: "RPC returned no transaction for that hash." };
    return {
      ok: true,
      hash: tx.hash,
      from: tx.from,
      to: tx.to,
      blockNumber: tx.blockNumber?.toString() ?? "pending",
    };
  } catch {
    return { ok: false, reason: "Aristotle RPC did not return that transaction." };
  }
}

const OWNER_ABI = [
  {
    type: "function",
    name: "ownerOf",
    stateMutability: "view",
    inputs: [{ name: "tokenId", type: "uint256" }],
    outputs: [{ name: "", type: "address" }],
  },
] as const;

export type AgentRead = { ok: true; owner: string } | { ok: false; reason: string };

export async function read8004Owner(tokenId: bigint): Promise<AgentRead> {
  try {
    const client = aristotleClient();
    const owner = await client.readContract({
      address: IDENTITY_8004,
      abi: OWNER_ABI,
      functionName: "ownerOf",
      args: [tokenId],
    });
    return { ok: true, owner };
  } catch {
    return { ok: false, reason: "Identity registry did not return ownerOf for that id. This page will not invent a ranking." };
  }
}
