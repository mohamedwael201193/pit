export const NAMED = {
  SIGNATURE_DECLINED: "Signature declined. PIT did not collect a key.",
  WRONG_NETWORK: "Wrong network. Pick MAINNET or TESTNET and stay there.",
  SESSION_EXPIRED: "Your session expired. Approve a new agent on this machine.",
  HL_UNFUNDED: "Your trading account shows no USDC. Spot USDC still counts as funded.",
  POLICY_BLOCK: "Your policy blocked this preview.",
  TEE_VERIFY_FAIL: "TeeML signature did not match the on-chain signer. Stopped.",
  SEED_FORBIDDEN: "PIT never asks for a seed phrase.",
  TWO_WALLETS: "Two wallets never share a workspace, session, or memory key.",
  TRANSFER_NOT_LIVE: "Transfer of Agentic ID is not live on mainnet.",
  AUTHORIZE_EXACT: "Type AUTHORIZE on the exact preview. Piped yes is never enough.",
  AUTHORIZE_REFUSED: "Authorization refused. The token did not match.",
} as const;

export const PERMISSIONS = [
  { k: "order", ok: true },
  { k: "cancel", ok: true },
  { k: "withdraw", ok: false },
  { k: "leverage", ok: false },
] as const;
