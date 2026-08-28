export type NamedId =
  | "SIGNATURE_DECLINED"
  | "WRONG_NETWORK"
  | "SESSION_EXPIRED"
  | "BACKEND_UNREACHABLE"
  | "SEED_FORBIDDEN"
  | "AUTHORIZE_WEB_DENIED";

export type NamedState = { id: NamedId; title: string; body: string; next: string };

const STATES: Record<NamedId, NamedState> = {
  SIGNATURE_DECLINED: {
    id: "SIGNATURE_DECLINED",
    title: "You declined",
    body: "You declined the bind message. PIT did not collect a key.",
    next: "Connect again when you are ready.",
  },
  WRONG_NETWORK: {
    id: "WRONG_NETWORK",
    title: "Wrong network",
    body: "This workspace is bound to one chain. Mixing Galileo with Aristotle is refused.",
    next: "Pick MAINNET or TESTNET and stay there.",
  },
  SESSION_EXPIRED: {
    id: "SESSION_EXPIRED",
    title: "Session expired",
    body: "The local session expired. If Hyperliquid still lists this PIT agent, create a new local session and PIT reuses the address. This browser never held the key.",
    next: "Create a new session on desktop or CLI.",
  },
  BACKEND_UNREACHABLE: {
    id: "BACKEND_UNREACHABLE",
    title: "Watch unreachable",
    body: "The public book feed did not answer. PIT will not invent opportunities.",
    next: "Retry. Empty Watch is the honest state.",
  },
  SEED_FORBIDDEN: {
    id: "SEED_FORBIDDEN",
    title: "No seed phrase",
    body: "PIT never asks for a seed phrase.",
    next: "Connect a wallet you already control.",
  },
  AUTHORIZE_WEB_DENIED: {
    id: "AUTHORIZE_WEB_DENIED",
    title: "Web cannot authorize",
    body: "Hyperliquid orders are signed on desktop or CLI. This browser cannot hold a session.",
    next: "Open PIT desktop after the preview is bound.",
  },
};

export function namedState(id: NamedId): NamedState {
  return STATES[id];
}

export function classifyError(message: string): NamedState {
  const m = message.toLowerCase();
  if (m.includes("declin") || m.includes("reject") || m.includes("denied")) return STATES.SIGNATURE_DECLINED;
  if (m.includes("network") || m.includes("chain")) return STATES.WRONG_NETWORK;
  if (m.includes("failed to fetch") || m.includes("unreachable")) return STATES.BACKEND_UNREACHABLE;
  return { ...STATES.BACKEND_UNREACHABLE, body: message || STATES.BACKEND_UNREACHABLE.body };
}
