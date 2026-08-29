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

## Install (Windows first)

Ordinary path:

1. Open [pit0g.vercel.app](https://pit0g.vercel.app)
2. Download the Windows installer from [GitHub Releases](https://github.com/mohamedwael201193/pit/releases)
3. Verify SHA256 against `SHA256SUMS` in that release
4. Install and launch PIT
5. Pair the browser at [/pair](https://pit0g.vercel.app/pair) with the one-time code shown on the machine
6. Connect **your wallet**. PIT never asks for a seed phrase.
7. Sign **Protect my strategy** in the paired browser. The Direct token stays on this computer.
8. Pick MAINNET or TESTNET. Connect **your Hyperliquid account**. Set **your policy**. Approve the printed agent (order and cancel only).

macOS and Linux: source build only until those installers are packaged and tested. Do not claim they are production-ready.

The installer is **not Authenticode-signed** until a code-signing certificate exists. Windows SmartScreen may warn. Verify the checksum. Do not skip that step.

Health (public Watch, `sign: false`): [pit-health.onrender.com](https://pit-health.onrender.com/health)

### How to verify your PIT download

```powershell
Get-FileHash .\PIT_0.6.1_x64-setup.exe -Algorithm SHA256
```

Compare with `SHA256SUMS` on the same GitHub Release. The source commit is on the release tag.

---

## 20 seconds

1. Connect **your wallet**. PIT never asks for a seed phrase.
2. Pick **MAINNET** or **TESTNET**. One workspace is exactly one pair. Never mixed.
3. Connect **your Hyperliquid account**. Unified accounts use spot USDC as buying power. PIT will not invent a perp deposit or pad size below the $10 venue minimum.
4. Set **your policy** (clip, leverage, kill, allowlist).
5. Create **your session** on desktop or CLI (order + cancel, 24 hours; reused if Hyperliquid still lists the PIT agent).
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

Build the native sealer (Go 1.24) and point `PIT_COMMITTEE_BIN` at it. A `.py` gate script is not the product path.

```powershell
cd sealer
go test ./...
go build -o pit-sealer .
# then set PIT_COMMITTEE_BIN to that file
```

From the repository root you can also run `make sealer` then `make pit`. If `PIT_COMMITTEE_BIN` is empty, PIT looks for `sealer/pit-sealer` next to the working directory. A missing binary still returns `sealer_not_wired`.

Normal users sign a Direct challenge from the paired browser (`pit direct` or Protect my strategy). PIT stores the wallet-signed `app-sk-` token in the OS keychain. `PIT_DIRECT_AUTH_FILE` is an operator override only. A Router dashboard `sk-` / `mk-` key is refused. The sealer writes evidence without the prompt or the authorization header.

The product committee is three sealed roles in order: researcher, challenger, risk. Role and envelope are split. This is not three independent providers unless the auth files actually differ.

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
  sdk/               Typed Go client; CanSign is always false
sdk/js               Browser package; canSign is always false
contracts/           Foundry 0.8.24 — PitDeskID, PitPolicy, PitReceipts,
                     PitForecasts, PitMemory
sealer/              Native Direct TeeML binary (HPKE Seal + VerifyE2EE)
apps/web             Vite + Privy connect (no HL session)
apps/desktop         Local authorize surface (session stays on the machine)
.env.example         Public RPCs and addresses only
```

---

## Prerequisites

- Go **1.22+** for `pit/`
- Go **1.24+** to build `sealer/`
- Foundry (`forge`)
- Node **20+** (web / desktop)
- A wallet (Rabby / MetaMask). No seed phrase is ever collected.
- Optional: official 0G storage Go client for `--proof` uploads
- Direct sealer beside `pit.exe` or at `PIT_COMMITTEE_BIN`. Normal users sign Direct in the paired browser. `PIT_DIRECT_AUTH_FILE` is an operator override only.

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

cd ..\sealer
go test ./...
go build -o pit-sealer .

cd ..\contracts
forge install foundry-rs/forge-std OpenZeppelin/openzeppelin-contracts
forge test -vvv
```

Optional live Hyperliquid public books (mainnet and testnet info endpoints; no orders):

```powershell
cd pit
go test -tags live ./internal/hl -count=1
```

CI runs the same Go and Foundry suites on every push (`.github/workflows/ci.yml`). Foundry currently reports **26** tests. Re-run `go test ./...` for the current passing count. Do not invent a number.

The Go health process (`pit/cmd/health`) exposes `GET /health` and `GET /watch`. Both always report `sign: false`. `/watch` returns live venue books or a real empty list. It never places an order and never includes a private book. In product mode the process refuses `PIT_SESSION_KEY`, `HL_SECRET`, and `PIT_MEMORY_KEY` in the environment.

### CLI

```powershell
cd pit
go run ./cmd/pit init --network testnet --wallet 0xYourAddress
go run ./cmd/pit session
go run ./cmd/pit login
go run ./cmd/pit policy
go run ./cmd/pit status
go run ./cmd/pit opportunities
go run ./cmd/pit ask --market market.json --book book.json
go run ./cmd/pit forecast
go run ./cmd/pit preview --market ETH --side buy --forecast <id>
go run ./cmd/pit verify --preview 0x... --root 0x... --network mainnet --workspace <id>
go run ./cmd/pit kill
go run ./cmd/pit revoke
```

Full command set:

`init` `login` `wallet` `pair` `network` `policy` `session` `companion` `approve` `hyperliquid` `agent` `compute` `direct` `ask` `research` `watch` `opportunities` `chat` `positions` `health` `forecast` `calibration` `preview` `authorize` `execute` `orders` `cancel` `status` `resolve` `card` `verify` `proof` `kill` `revoke` `doctor` `activity` `receipt` `logout` `version` `update`

Every command accepts `--json`. `pit version` prints `PIT 0.6.1`. `pit doctor` probes version, wallet, network, OS keychain, memory-key hazard, Hyperliquid, 0G RPC, companion, sealer, Direct token (keychain, operator file, or sponsored compute), Direct credit, TEE evidence, storage client, registry, session, Hyperliquid agent, and policy. It never prints secrets. A global `PIT_MEMORY_KEY` is a doctor failure. The desktop can bind a public wallet, pin policy, and mint a session without a terminal. `pit research [ETH] --hypothesis none|long|short` seals a user hypothesis into the private book. `pit pair`, `pit approve`, `pit execute`, `pit mcp`, `pit scan`, `pit mission`, and `pit calibration` are first-class commands. `pit direct` issues the official wallet-signed Direct challenge and stores the token in the keychain. Discovery ranks books this account can actually size: policy PASS is not the same as execution-feasible. Spot USDC is not silently treated as perp margin unless Hyperliquid reports unified account mode.

Official storage client (not the TypeScript SDK): `upload --url --file --key --encryption-key` and `download --proof --root --file --encryption-key`. `pit proof` requires `--key-file` per workspace and refuses a global memory key.

`session` creates a 24-hour order/cancel agent in the OS keychain (or `PIT_KEYRING=file` for tests). It prints the agent address. It never prints the key. Hyperliquid API wallet names must be under 17 characters (`PIT-` plus eight hex). Your wallet must `approveAgent` that address.

`companion` listens on `127.0.0.1:17373` only. Pairing is a one-time code. The browser receives a device token, never the session key. Foreign origins and non-loopback clients are refused.

`proof` requires `--root`, `--out`, and `--key-file`. It never uses `PIT_MEMORY_KEY`.

`status` prints the bound workspace, live session address, Hyperliquid PIT-agent list, ledger row, and whether the bound cloid is on the venue. It never prints a session key and never places an order.

`preview` prints the exact bound fields from a **live mark** and the host sizer. It requires `--market`, `--side`, `--forecast`, and a live session. The model cannot set size. Any later mutation of the card invalidates `AUTHORIZE`.

`proof` uses the official Go storage client with `--proof`. It does not use a global memory key.

`authorize` requires a TTY, `--i-understand`, the exact word `AUTHORIZE`, a live session, and a bound preview. Piped `yes` is rejected. A matching token records `authorized` on the local ledger for that workspace and clientOrderId. A second click is `duplicate_click`. The CLI signs an L1 envelope locally from the session key, queries the Hyperliquid PIT-agent list on the matching venue, and posts only when that agent is present and unexpired. An unsigned or unlinked payload never reaches the venue.

`opportunities` reads live venue books for the policy allowlist and never places an order. Empty Watch is a real empty state.

`policy` prints the law and pins the hash for the bound workspace. A later mutation of clip, assets, or kill fails closed.

Restarting the process keeps a previewed action. A second click does not apply twice.

Web refresh cannot sign. Desktop recovers the exact preview from the local ledger. Two wallets never share a workspace.

`cancel` requires a live session and an authorized preview. It queries `frontendOpenOrders` for the bound cloid, builds a cancel wire from the live asset index, signs locally, and checks the Hyperliquid PIT-agent list. Missing venue state is `not_on_venue` or `query_exchange_first`. A linked cancel posts only when the cloid is on the venue. Cancel cannot withdraw.

`ask` requires `--market` and `--book` files. Missing files return `empty_envelope`. Missing sealer binary, missing Direct token (keychain or operator file), Galileo VerifyE2EE-unproven, or a Router URL all stop the operation. There is no fallback. MCP cannot export the token or invoke the sealer.

### MCP (read-only)

```powershell
cd pit
go run ./cmd/mcp
```

Tools: `market` `opportunities` `forecast` `status` `card` `verify` `preview`  
Not tools: `authorize` `order` `cancel` `export_session`

`preview` is prepare-only. MCP cannot authorize. Live `opportunities` are public Watch cards with `trade: false`.

### SDK

`pit/sdk` is for apps that must not hold a session secret. `Status().CanSign` is always `false`. Explorer URLs follow the selected network.

### Web (connect and inspect)

```powershell
cd apps\web
copy .env.example .env
npm install
npm run dev
```

Routes: `/` landing, `/signin` wallet gate, `/pair` local pairing, `/app` Watch home, `/app/start` onboarding, `/app/activity` named states, `/app/policy` clip cards, `/app/account` verify, `/verify` receipt form. Open http://localhost:3000 or `#verify`.

Web can connect, bind, inspect Watch, policy, and receipts. It cannot create a Hyperliquid session or present Authorize.

Production Watch reads `VITE_HEALTH_URL` (https://pit-health.onrender.com). That process never signs.

MAINNET shows production copy and the Aristotle explorer. TESTNET shows the integration lab and the Galileo explorer. They never share a workspace.

CI builds the web app, runs Playwright copy tests, and builds the desktop frontend on every push. Windows NSIS packaging runs from `.github/workflows/release.yml` on version tags. Session secrets stay in the OS keychain (`PIT_KEYRING=os`) or a local file store in tests. They never enter the web bundle.

### Desktop (authorize locally)

```powershell
cd apps\desktop
npm install
npm run dev
```

Open http://localhost:3001

The desk shows a twelve-beat start, the pipeline, permissions, and an authorize field. The authorize control stays disabled until a live session exists on this machine. Type `AUTHORIZE` on the exact preview. Piped yes is never enough.

Session permissions card: order and cancel allowed; withdraw and leverage denied.

### Playwright (labeled harness)

The Playwright specs live in `apps/web/playwright/`. They assert public copy only. They never stub TeeML, never type `AUTHORIZE`, and never place an order.

```powershell
cd apps\web
npm install
npx playwright install chromium
$env:VITE_PRIVY_APP_ID="cmtafcijw02av0cl1ay81om7m"
npm run build
npm run preview
# in another shell, or let Playwright start preview on :4173
npm run test:e2e
```

Copy-only harness files in `apps/web/e2e/` and `apps/desktop/e2e/` are not Playwright runners.

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
- Restart keeps a previewed action. A duplicate click does not apply twice.
- A mismatched storage encryption key fails closed.
- MCP and SDK cannot authorize or export a session.
- Expired or revoked sessions cannot order or cancel.
- Provider, model, and teeSigner must match the catalog for that network. A tampered E2EE payload fails closed.
- Logs never include session keys, memory keys, or the private book.
- Two wallets never share a workspace, session, policy, or memory key.
- A huge model size still cannot exceed clip. Receipt hashes cannot be filed twice.
- Switching a bound workspace from mainnet to testnet is denied.
- Owner wallets cannot self-report 8004 feedback. A stranger cannot report.
- A daily loss halt and a kill switch stop new orders. The model cannot flip either.
- Mock market sources are denied.
- Mainnet RPC cannot be paired with the Galileo chain id. SIWE chain must match the workspace.
- Open interest must be finite. Slippage above policy fails closed.
- Mark price must be finite. Thin liquidity and cooldown fail closed. Script sealers (`.ts` / `.js` / `.mjs` / `.cjs`) are refused.
- Leverage above policy, a foreign venue, `sendAsset`, and `approveAgent` fail closed. A timeout never blindly reposts.
- Session mint is 24 hours. If Hyperliquid still lists the PIT agent, PIT reuses it. A stale preview cannot authorize.
- Calibration below the floor fails. SOL is outside the default universe. Galileo sealed ask stays disabled until VerifyE2EE is proven.
- Impact prices must be finite. MCP forecasts never carry `sizeUsd`. Health JSON cannot include a session.
- Strategy health needs 30 resolved samples. A TypeScript storage client is refused. MCP cannot export a session. Transfer of Agentic ID is not live on Aristotle.
- Session keys cannot be exported as JSON. The fixture master address cannot be the product user. Spot USDC counts as funded. Signing never happens in the browser.
- SIWE nonces cannot replay. Preview nonces must match. MCP cannot place orders. A Galileo desk address cannot be used on Aristotle.
- Watch never places orders. MCP opportunities never trade. `GET /watch` cannot include a private book.
- CLI never prints session secrets.
- Web bundle must not contain session private-key types.
- An expired session cannot type `AUTHORIZE`. A ledger record from another workspace is `wrong_workspace`.
- Playwright specs never stub TeeML success and never place an order.
- `pit authorize` still fails closed without a live session, even after the exact token is typed.
- The browser SDK cannot authorize. MCP cannot bind another user's workspace.
- An unsigned exchange payload never posts. MCP cannot `post`. The browser SDK cannot post.

If a Direct request fails, PIT **stops**. It does not retry on the Router.

---

## Honest limitations

- Foundation Agentic ID **transfer is not live on Aristotle**. The UI must say so.
- Galileo iTransfer is an official path; PIT has **not** shown a transfer tx that changes `ownerOf`.
- Galileo sealed committee is **not** the mainnet glm-5.2 committee. VerifyE2EE on Omni is still required before enabling testnet sealed ask.
- Live Direct Seal + VerifyE2EE needs the native sealer and a wallet-signed Direct token in the OS keychain (or an operator `PIT_DIRECT_AUTH_FILE`). Python and TypeScript sealers are refused. A missing binary returns `sealer_not_wired`. A missing Direct token returns `direct_token_required`. A Router `sk-` key returns `router_api_key_denied`. Plaintext sealer output returns `TEE_VERIFY_FAIL`. The three roles share one provider unless the auth files actually differ. The provider Ledger still requires the user wallet to have Direct credit; PIT does not invent that credit.
- Live Hyperliquid dust (`order` then `cancel`) requires a funded account on the matching venue, a user `AUTHORIZE`, a resolved asset index, a signed exchange payload, and an unexpired PIT-agent row for the session agent. An unsigned or unlinked authorized ledger row is not a fill.
- Desktop is the local authorize shell; full OS keychain packaging (Tauri / stronghold) can still be tightened.
- Do not claim hardware quotes unless the verifier is wired.

---

## Demo path

Connect wallet → select MAINNET or TESTNET → read the capability list for that network → (desktop/CLI) connect your trading account → set your policy → create your session → first private analysis → inspect the exact preview → `AUTHORIZE` or walk away → verify the receipt hash on the matching explorer.

---

## License

MIT for contracts (`SPDX-License-Identifier: MIT`). Application code in this repository is provided for the PIT desk as published.
