import { defineChain, http, isHex, createPublicClient, encodeFunctionData, keccak256, stringToBytes } from "viem";
import { ARISTOTLE_ID, ARISTOTLE_RPC, DESK_ID_CONTRACT, IDENTITY_8004, PIT_AGENT, REPUTATION_8004 } from "./facts";

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
  | {
      ok: true;
      hash: `0x${string}`;
      from: string;
      to: string | null;
      blockNumber: string;
      status: "success" | "reverted" | "pending";
    }
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
    let status: "success" | "reverted" | "pending" = "pending";
    try {
      const rec = await client.getTransactionReceipt({ hash: h as `0x${string}` });
      status = rec.status === "success" ? "success" : "reverted";
    } catch {
      status = "pending";
    }
    return {
      ok: true,
      hash: tx.hash,
      from: tx.from,
      to: tx.to,
      blockNumber: tx.blockNumber?.toString() ?? "pending",
      status,
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

const AUTH_ABI = [
  {
    type: "function",
    name: "isAuthorized",
    stateMutability: "view",
    inputs: [
      { name: "tokenId", type: "uint256" },
      { name: "user", type: "address" },
    ],
    outputs: [{ name: "", type: "bool" }],
  },
] as const;

export type AgentRead = { ok: true; owner: string } | { ok: false; reason: string };

export type DeskRead =
  | { ok: true; owner: string; ownerAuthorized: boolean }
  | { ok: false; reason: string };

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

export async function readDesk(tokenId: bigint): Promise<DeskRead> {
  try {
    const client = aristotleClient();
    const owner = await client.readContract({
      address: DESK_ID_CONTRACT,
      abi: OWNER_ABI,
      functionName: "ownerOf",
      args: [tokenId],
    });
    let ownerAuthorized = false;
    try {
      ownerAuthorized = await client.readContract({
        address: DESK_ID_CONTRACT,
        abi: AUTH_ABI,
        functionName: "isAuthorized",
        args: [tokenId, owner],
      });
    } catch {
      ownerAuthorized = false;
    }
    return { ok: true, owner, ownerAuthorized };
  } catch {
    return { ok: false, reason: "PitDeskID did not return ownerOf. This page will not invent a mint." };
  }
}

const LAST_INDEX_ABI = [
  {
    type: "function",
    name: "getLastIndex",
    stateMutability: "view",
    inputs: [
      { name: "agentId", type: "uint256" },
      { name: "clientAddress", type: "address" },
    ],
    outputs: [{ name: "", type: "uint64" }],
  },
] as const;

const READ_FEEDBACK_ABI = [
  {
    type: "function",
    name: "readFeedback",
    stateMutability: "view",
    inputs: [
      { name: "agentId", type: "uint256" },
      { name: "clientAddress", type: "address" },
      { name: "feedbackIndex", type: "uint64" },
    ],
    outputs: [
      { name: "value", type: "int128" },
      { name: "valueDecimals", type: "uint8" },
      { name: "tag1", type: "string" },
      { name: "tag2", type: "string" },
      { name: "isRevoked", type: "bool" },
    ],
  },
] as const;

const TOKEN_URI_ABI = [
  {
    type: "function",
    name: "tokenURI",
    stateMutability: "view",
    inputs: [{ name: "tokenId", type: "uint256" }],
    outputs: [{ name: "", type: "string" }],
  },
] as const;

export type ReputationRead =
  | { ok: true; index: string; tag1: string; tag2: string; revoked: boolean; uri: string }
  | { ok: false; reason: string };

export async function read8004Reputation(agentId: bigint, client: `0x${string}`): Promise<ReputationRead> {
  try {
    const clientRpc = aristotleClient();
    const uri = await clientRpc.readContract({
      address: IDENTITY_8004,
      abi: TOKEN_URI_ABI,
      functionName: "tokenURI",
      args: [agentId],
    });
    const index = await clientRpc.readContract({
      address: REPUTATION_8004,
      abi: LAST_INDEX_ABI,
      functionName: "getLastIndex",
      args: [agentId, client],
    });
    if (index === 0n) {
      return { ok: true, index: "0", tag1: "", tag2: "", revoked: false, uri };
    }
    const fb = await clientRpc.readContract({
      address: REPUTATION_8004,
      abi: READ_FEEDBACK_ABI,
      functionName: "readFeedback",
      args: [agentId, client, index],
    });
    return {
      ok: true,
      index: index.toString(),
      tag1: fb[2],
      tag2: fb[3],
      revoked: fb[4],
      uri,
    };
  } catch {
    return { ok: false, reason: "Reputation registry did not return feedback for that client." };
  }
}

export function encodeAuthorizeUsage(tokenId: bigint, user: `0x${string}`): `0x${string}` {
  return encodeFunctionData({
    abi: [
      {
        type: "function",
        name: "authorizeUsage",
        inputs: [
          { name: "tokenId", type: "uint256" },
          { name: "user", type: "address" },
        ],
      },
    ],
    functionName: "authorizeUsage",
    args: [tokenId, user],
  });
}

export function encodeRevokeAuthorization(tokenId: bigint, user: `0x${string}`): `0x${string}` {
  return encodeFunctionData({
    abi: [
      {
        type: "function",
        name: "revokeAuthorization",
        inputs: [
          { name: "tokenId", type: "uint256" },
          { name: "user", type: "address" },
        ],
      },
    ],
    functionName: "revokeAuthorization",
    args: [tokenId, user],
  });
}

export function encodeSetAgentURI(agentId: bigint, uri: string): `0x${string}` {
  return encodeFunctionData({
    abi: [
      {
        type: "function",
        name: "setAgentURI",
        inputs: [
          { name: "agentId", type: "uint256" },
          { name: "newURI", type: "string" },
        ],
      },
    ],
    functionName: "setAgentURI",
    args: [agentId, uri],
  });
}

export function encodeMint(to: `0x${string}`, uri: string, dataHash: `0x${string}`, desc: string): `0x${string}` {
  return encodeFunctionData({
    abi: [
      {
        type: "function",
        name: "mint",
        inputs: [
          { name: "to", type: "address" },
          { name: "uri", type: "string" },
          { name: "dataHash", type: "bytes32" },
          { name: "dataDescription", type: "string" },
        ],
      },
    ],
    functionName: "mint",
    args: [to, uri, dataHash, desc],
  });
}

export function publicDataHash(uri: string): `0x${string}` {
  return keccak256(stringToBytes(uri));
}

type Eth = { request: (args: { method: string; params?: unknown[] }) => Promise<string> };

export async function sendWalletCall(from: string, to: `0x${string}`, data: `0x${string}`): Promise<string> {
  const eth = (window as unknown as { ethereum?: Eth }).ethereum;
  if (!eth?.request) throw new Error("Connect a wallet. PIT never asks for a seed phrase.");
  try {
    await eth.request({ method: "wallet_switchEthereumChain", params: [{ chainId: "0x4115" }] });
  } catch {
    try {
      await eth.request({
        method: "wallet_addEthereumChain",
        params: [
          {
            chainId: "0x4115",
            chainName: "0G Aristotle",
            nativeCurrency: { name: "0G", symbol: "0G", decimals: 18 },
            rpcUrls: [ARISTOTLE_RPC],
            blockExplorerUrls: ["https://chainscan.0g.ai"],
          },
        ],
      });
    } catch {
      throw new Error("Switch the wallet to Aristotle 16661.");
    }
  }
  return eth.request({
    method: "eth_sendTransaction",
    params: [{ from, to, data, chainId: "0x4115" }],
  });
}

export const PIT_AGENT_ADDR = PIT_AGENT.address as `0x${string}`;
