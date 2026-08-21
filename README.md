# PIT — Private Alpha OS

Private research. Controlled execution. A desk that learns.

PIT is a **multi-user trading desk** for 0G + Hyperliquid. It seals **your** private book into **0G Direct TeeML**, runs researcher / challenger / risk over that book, sizes the order on the **host** (the model cannot set size), then waits for **you** to authorize an exact preview. A short-lived session can `order` or `cancel`. It **cannot withdraw**.

This repository is the product: Go core, Foundry contracts, CLI, read-only MCP, typed SDK, wallet-connect web, and a local desktop authorize surface.

---

## Why this exists

Pasting a trading book into a public chat API is how alpha leaks. Giving a bot a withdraw key is how accounts die.

PIT splits the job:

| Layer | What it does | What it cannot do |
|---|---|---|
| **Your wallet** | Connect, SIWE bind, mint / authorize Desk ID, register 8004, pin policy | Never collected as a seed |
| **0G Direct TeeML** | Sealed inference over the private book + public market | Router `sk-` path is forbidden for the book |
| **Host engine** | Size, policy, preview hash, kill switch | LLM JSON cannot raise clip or leverage |
| **Your session** | Hyperliquid `order` and `cancel` only | Withdraw, leverage change, `approveAgent`, transfers |
| **0G Storage** | Encrypted objects + `--proof` download | TypeScript SDK is not used for proofs |
| **Calibration** | Brier / ECE when N is large enough | Invented accuracy when the sample is empty |

Copy in the product is **YOU / YOUR WALLET / YOUR TRADING ACCOUNT / YOUR SESSION / YOUR POLICY / YOUR MONEY**.

---

## 20 seconds

1. Connect **your wallet**. PIT never asks for a seed phrase.
2. Pick **MAINNET** or **TESTNET**. One workspace is exactly one pair. Never mixed.
3. Connect **your Hyperliquid account**. Spot USDC counts as funded.
4. Set **your policy** (clip, leverage, kill, allowlist).
5. Create **your session** on desktop or CLI (order + cancel, TTL ≤ 1h).
6. Ask, or wait for a policy-eligible opportunity.
7. Type `AUTHORIZE` on the exact preview — or do nothing.

Web can connect and inspect. **Signing Hyperliquid orders happens on desktop or CLI.** Session keys never enter the browser, Vercel, or a server env.

---

## Architecture

```
YOUR PRIVATE BOOK + PUBLIC MARKET
  → HPKE SealRequest
  → Direct 0G TeeML (not router-api.0g.ai)
  → VerifyE2EE  scheme zg-sig-v1/e2ee-ct
  → recovered signer == on-chain teeSigner
  → host sizer + policy
  → YOU AUTHORIZE (TTY, exact token AUTHORIZE)
  → session order | cancel
  → 0G Storage --proof + on-chain receipt
  → calibration
```

Progress labels in the UI map to real backend steps. They are not spinner timers:

`PRIVATE_BOOK` → `SEALING` → `TEE` → `TEE_SIGNATURE` → `ONCHAIN_SIGNER` → `STORAGE` → `RECEIPT` → `CALIBRATION`

Committee independence today is **prompt-envelope / same funded provider** on mainnet glm-5.2. Do not claim three independent Direct providers until each role has its own funded TeeML SKU.

---

## Dual environment (do not mix)

`PIT_NETWORK` selects **0G chain and Hyperliquid venue together**.

| | TESTNET | MAINNET |
|---|---|---|
| 0G chain | Galileo `16602` | Aristotle `16661` |
| RPC | `https://evmrpc-testnet.0g.ai` | `https://evmrpc.0g.ai` |
| Explorer | [chainscan-galileo.0g.ai](https://chainscan-galileo.0g.ai) | [chainscan.0g.ai](https://chainscan.0g.ai) |
| Hyperliquid | [api.hyperliquid-testnet.xyz](https://api.hyperliquid-testnet.xyz) | [api.hyperliquid.xyz](https://api.hyperliquid.xyz) |
| Direct chat SKU | `qwen/qwen2.5-omni-7b` (TeeML ack; PIT VerifyE2EE **not yet proven**) | `glm-5.2` Direct (proven Seal + VerifyE2EE path) |
| glm-5.2 | **Absent** — do not copy the mainnet SKU | Required for the sealed committee |
| Desk ID | Not deployed in this repo yet | `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` |
| Agentic ID transfer | Official attestor lives on 16602; PIT has not executed a transfer | **Disabled** — attestor bytecode is not on Aristotle |

Kill switches:

- Mixing Aristotle compute with Hyperliquid testnet in one workspace
- Falling back from Direct to the Router when TeeML fails
- Showing a testnet-only control as mainnet-live
- Setting `PIT_ALLOW_FALLBACKS` to anything other than `false`

---

## Direct TeeML vs Router

Two official 0G inference paths exist. **The private book uses Direct only.**

| | Router | Direct |
|---|---|---|
| Endpoint | `router-api.0g.ai` / testnet twin | Provider URL from on-chain `getService` |
| Credential | Dashboard `sk-` / `mk-` | Wallet-signed `app-sk-` (tokenId 255 ephemeral) |
| Book visibility | Router sees prompts to route and bill | No gateway in the middle |
| Attestation | `verify_tee` is a Router boolean | HPKE + `VerifyE2EE` + on-chain `teeSigner` |
| PIT private book | **Forbidden** | **Required** |

Catalog URLs are allowed **only to list models**. They are never the inference URL for a sealed ask.

Mainnet glm-5.2 (proven Direct path):

- Provider `0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D`
- teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`
- Serving `0x47340d900bdFec2BD393c626E12ea0656F938d84`
- Ledger `0x2dE54c845Cd948B72D2e32e39586fe89607074E3`
- `transferFund` service name is `inference-v1.0`

Galileo TeeML chat (ack’d, not claimed live for PIT until VerifyE2EE succeeds):

- Model `qwen/qwen2.5-omni-7b`
- Provider `0xa48f01287233509FD694a22Bf840225062E67836`
- teeSigner `0x83df4B8EbA7c0B3B740019b8c9a77ffF77D508cF`

Set `PIT_COMMITTEE_BIN` to the local Direct sealer binary. A `.py` gate script is not the product path.

---

## Hyperliquid session

1. Identify **your** HL account (`clearinghouseState` + `spotClearinghouseState`). Spot USDC with perp value 0 is `FUNDED_SPOT`, not unfunded.
2. Generate an agent key in the **OS keychain** / local keyring. Name `PIT-{workspace8}`. TTL ≤ 1 hour.
3. **Your wallet** signs `approveAgent`. The session key never signs `approveAgent`.
4. Permissions: order ✓  cancel ✓  withdraw ✗  leverage ✗
5. Every `order` is bound to a preview hash. Timeouts query the exchange; PIT does not blindly repost.

Minimum notional is **$10**. Tick rounding uses 5 significant figures.

Denied session actions include `withdraw3`, `updateLeverage`, `sendAsset`, `approveAgent`, `usdSend`, `spotSend`, `vaultTransfer`, `twapOrder`, `modify`, and the rest of the explicit deny list in `pit/internal/session`.

---

## Identity, policy, memory

- **Workspace** — UUID v4 bound to your EVM address. Guessing another UUID returns `not found`, not a decrypt error.
- **SIWE** — bind without storing a private key.
- **ERC-8004** — user-owned Identity + separate feedback reporter. Owner `giveFeedback` reverts. IDs are **not** portable across 16661 and 16602.
- **ERC-7857 Desk ID** — mint / `authorizeUsage` / `isAuthorized` / revoke on Aristotle. Host must check `isAuthorized` before sealed inference for that desk. `iTransferFrom` / `iCloneFrom` revert on mainnet (`AttestorNotOnAristotle`). ERC-721 `transferFrom` is disabled.
- **Policy** — host cards: clip, daily loss, leverage, asset allowlist, venue, calibration floor, cooldown, uncertainty, slippage, liquidity, kill. Hash is pinned per workspace. The model cannot raise `maxClip`.
- **Memory** — typed kinds under `{network}/ws/{workspaceId}/...`. Encryption key is `0x` + 32 bytes in the keychain. A global `PIT_MEMORY_KEY` is forbidden in product mode. The key must never appear in a prompt.
- **Ledger** — durable SQLite per workspace + network. Exactly-once `cloid`. Recover never reposts a signed or timed-out action without an exchange view.
- **Calibration** — Brier and ECE. If N < 30 the health card says there are not enough samples. It will not invent a win rate.

---

## Repository layout

```
pit/                 Go module github.com/mohamedwael201193/pit
  cmd/pit            CLI
  cmd/mcp            Read-only MCP stdio
  internal/          identity, workspace, wallet, siwe, policy, session,
                     hl, engine, exec, compute, ledger, storage, memory,
                     chain8004, deskid, calib, watch, verify, cli, ui
  mcp/               Allowed tool list (no authorize / order / key export)
  sdk/               Typed client; CanSign is always false
contracts/           Foundry 0.8.24 — PitDeskID, PitPolicy, PitReceipts,
                     PitForecasts, PitMemory
apps/web             Vite + Privy connect (no HL session)
apps/desktop         Local authorize surface (session stays on the machine)
.env.example         Public RPCs and addresses only
```

---

## Prerequisites

- Go **1.22+**
- Foundry (`forge`)
- Node **20+** (web / desktop)
- A wallet (Rabby / MetaMask). No seed phrase is ever collected.
- Optional: official 0G storage Go client for `--proof` uploads
- Optional: Direct committee binary at `PIT_COMMITTEE_BIN`

---

## Environment

```powershell
copy .env.example .env
```

**Do not commit `.env`.**

| Kind | Examples | Rule |
|---|---|---|
| Public infra | RPCs, contract addresses, catalog URLs, Privy **app id** | Safe in `.env.example` |
| User secrets | HL session key, memory key, wallet | OS keychain / local only |
| Server-only | Privy app secret | Never in the web bundle |
| Forbidden | Router `sk-` / `mk-` for the private book | Process must not use them |
| Kill | `PIT_ALLOW_FALLBACKS=true` | Process exits |

Public Privy app id (wallet connect only): `cmtafcijw02av0cl1ay81om7m`

Add `localhost:3000` (web) and `localhost:3001` (desktop) in the Privy dashboard allowed origins.

---

## Commands

### Tests

```powershell
cd pit
go test ./...

cd ..\contracts
forge install foundry-rs/forge-std OpenZeppelin/openzeppelin-contracts
forge test -vvv
```

Optional live Hyperliquid public book:

```powershell
cd pit
go test -tags live ./internal/hl -count=1
```

CI runs the same Go and Foundry suites on every push (`.github/workflows/ci.yml`).

### CLI

```powershell
cd pit
go run ./cmd/pit init --network testnet --wallet 0xYourAddress
go run ./cmd/pit login
go run ./cmd/pit policy
go run ./cmd/pit status
go run ./cmd/pit verify --preview 0x... --root 0x... --network mainnet --workspace <id>
go run ./cmd/pit kill
```

Full command set:

`init` `login` `policy` `ask` `opportunities` `forecast` `preview` `authorize` `cancel` `status` `resolve` `card` `verify` `kill`

`authorize` requires a TTY, `--i-understand`, and the exact word `AUTHORIZE`. Piped `yes` is rejected.

### MCP (read-only)

```powershell
cd pit
go run ./cmd/mcp
```

Tools: `market` `opportunities` `forecast` `status` `card` `verify`  
Not tools: `authorize` `order` `cancel` `export_session`

### SDK

`pit/sdk` is for apps that must not hold a session secret. `Status().CanSign` is always `false`. Explorer URLs follow the selected network.

### Web (connect and inspect)

```powershell
cd apps\web
copy .env.example .env
npm install
npm run dev
```

Open http://localhost:3000

### Desktop (authorize locally)

```powershell
cd apps\desktop
npm install
npm run dev
```

Open http://localhost:3001

Session permissions card: order and cancel allowed; withdraw and leverage denied.

---

## Contract addresses

### Mainnet (Aristotle 16661)

| Contract | Address |
|---|---|
| PitDeskID | [`0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35`](https://chainscan.0g.ai/address/0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35) |
| InferenceServing | `0x47340d900bdFec2BD393c626E12ea0656F938d84` |
| LedgerManager | `0x2dE54c845Cd948B72D2e32e39586fe89607074E3` |
| Storage Flow | `0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526` |
| Storage indexer | `https://indexer-storage-turbo.0g.ai` |
| ERC-8004 Identity | `0x8004A169FB4a3325136EB29fA0ceB6D2e539a432` |
| ERC-8004 Reputation | `0x8004BAa17C55a88189AE136b182e5fdA19dE9b63` |

PitReceipts / PitForecasts / PitMemory are in `contracts/src` and tested; production addresses are filled in `.env` after deploy.

Desk ID production interface IDs: `0x2afbede9` / `0xdf597d99` / `0x74f8628b` / `0x80ac58cd`.

### Testnet (Galileo 16602)

| Contract | Address |
|---|---|
| InferenceServing | `0xa79F4c8311FF93C06b8CfB403690cc987c93F91E` |
| LedgerManager | `0xE70830508dAc0A97e6c087c75f402f9Be669E406` |
| Storage Flow | `0x22E03a6A89B950F1c82ec5e74F8eCa321a105296` |
| Storage indexer | `https://indexer-storage-testnet-turbo.0g.ai` |
| ERC-8004 Identity | `0x8004A818BFB912233c491871b3d84c89A494BD9e` |
| ERC-8004 Reputation | `0x8004B663056A597Dffe9eCcC1965A193B7388713` |
| PitDeskID | not deployed from this repo yet |

---

## 0G Storage

Use the **official Go client** only for proofs.

- Encryption key must be `0x` + 32 bytes (66 characters).
- Download includes `--proof`.
- A `.ts` SDK path is rejected.
- Object keys are `{network}/ws/{workspaceId}/{kind}/{name}`. Workspace A cannot assert B’s prefix.

---

## Security model

- Fail closed. No mock TeeML success. No mock exchange.
- No global `PIT_MASTER_ADDRESS` as the product user.
- No mnemonic field in any form (`SEED_FORBIDDEN`).
- Named wallet states: `SIGNATURE_DECLINED`, `WRONG_NETWORK`, `SESSION_EXPIRED`, `HL_UNFUNDED`, `POLICY_BLOCK`, `TEE_VERIFY_FAIL`.
- Preview hash mismatch, expired preview, `cloid` replay, kill switch, and policy version mismatch all deny.
- MCP and SDK cannot authorize or export a session.
- Web bundle must not contain session private-key types.

If a Direct request fails, PIT **stops**. It does not retry on the Router.

---

## Honest limitations

- Foundation Agentic ID **transfer is not live on Aristotle**. The UI must say so.
- Galileo iTransfer is an official path; PIT has **not** shown a transfer tx that changes `ownerOf`.
- Galileo sealed committee is **not** the mainnet glm-5.2 committee. VerifyE2EE on Omni is still required before enabling testnet sealed ask.
- Live Direct Seal + VerifyE2EE inside the product binary needs `PIT_COMMITTEE_BIN` pointed at the sealer.
- Live Hyperliquid dust (`order` then `cancel`) still requires a funded account on the matching venue and a user `AUTHORIZE`.
- Desktop is the local authorize shell; full OS keychain packaging (Tauri / stronghold) can still be tightened.
- Do not claim hardware quotes unless the verifier is wired.

---

## Demo path

Connect wallet → select MAINNET or TESTNET → read the capability list for that network → (desktop/CLI) connect your trading account → set your policy → create your session → first private analysis → inspect the exact preview → `AUTHORIZE` or walk away → verify the receipt hash on the matching explorer.

---

## License

MIT for contracts (`SPDX-License-Identifier: MIT`). Application code in this repository is provided for the PIT desk as published.
