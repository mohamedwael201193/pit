# PIT — Implementation Plan

**Do not start product UI until Phase 0 checklist is signed in history.**  
**Do not implement from `_gate` Python as production** — promote `hlcore` Go and `_gate/committee` Go, then rewrite Python out.  
**Every phase: fail closed. No mocks. No env-gating.**

**M17 hardening (2026-08-26):** Phase 6 committee-over-book and Phase 14 DeskID mint/authz are **proven in `_gate`**. Promote, do not re-research from zero. Official attestor remains Galileo `16602`. Aristotle `16661` iTransfer is **NOT LIVE**.

**M18 revalidation (2026-08-26):** Dual **Mainnet + Testnet**. Router dashboard `sk-` is **not** required for Direct. Galileo catalog ≠ Aristotle catalog. Implementation agent **must read the External resources links before coding each module**.

**M19 implementation:** Phase 0 revalidated live. Phase 1 workspace/SIWE libraries landed. Core host engine, session allowlist, durable ledger, compute deny-router, contracts (DeskID + Policy), CLI/MCP/SDK, wallet-connect web shell in progress.

**M21 implementation (2026-08-26):** Direct runner args + dual-network SKU isolation (Galileo Omni unproven; no glm-5.2 copy). HL order/cancel wires; approveAgent is master-wallet only. Exchange POST rejects mocks. Policy pin files per workspace. Desktop shell with session permissions card. Storage/workspace IDOR tests. **Still open:** live Seal+VerifyE2EE in-process, live dust order, Storage upload, Tauri packaging, Chrome E2E against live venues.

Submit: **2026-08-30T15:00:00Z**. Wave 3 in-window git only.

Shared core: `pit/` (Go) + `contracts/` (Foundry). Surfaces: desktop / CLI / MCP-RO / SDK / web-no-session. One policy engine. One executor.

**Competitor consult (always):** `D:\route\0g\03_COMPETITORS_WINNERS_AND_CODE_REVERSE_ENGINEERING.md` to EOF + local clones under `D:\route\0g\research\clones`. KEPT: `D:\route\Flare\kept` + https://kept-ruby-five.vercel.app/ — **visual only**.

---

## Mandatory read-before-implement

| Before coding | Read these first (current official + local clone) |
|---|---|
| Compute / committee | https://docs.0g.ai/developer-hub/building-on-0g/compute-network/inference.md · router/{overview,authentication,comparison,privacy}.md · `D:\route\0g\research\0g-pc-e2ee` · `0g-compute-ts-sdk` `constants.ts` + `inference/broker/base.ts` · `0g-serving-broker` `doc/api-token.md` · live `GET https://router-api.0g.ai/v1/models` |
| Agentic ID / 7857 | https://docs.0g.ai/ai-context?_highlight=itransfer#agentic-id-formerly-inft · https://agenticid.0g.ai/config · `D:\route\0g\research\0g-agentic-id` IERC7857 + `DEPLOYMENT.md` · Lumen `OracleNotLive` |
| Storage | official Go client `D:\route\0g\research\0g-storage-client` · docs storage SDK (TS **forbidden** for `--proof`) |
| Hyperliquid | https://hyperliquid.gitbook.io/hyperliquid-docs · exchange-endpoint · nonces-and-api-wallets · account-abstraction-modes |
| Frontend | `D:\route\Flare\kept` (visual only) + https://kept-ruby-five.vercel.app/ |
| Env | `ENVIRONMENT.md` A–F this folder |

---

## Dual-environment contract (every phase inherits this)

| Field | Rule |
|---|---|
| **Testnet behavior** | Galileo `16602` + HL testnet. Use D-class env vars. Do not assume glm-5.2 exists. |
| **Mainnet behavior** | Aristotle `16661` + HL mainnet. Use E-class env vars. iTransfer remains disabled. |
| **Fallback** | None. Direct fail → stop. Missing Galileo TeeML chat funding → disable sealed ask on testnet, do not Router-downgrade. |
| **Kill** | Showing a testnet-only control as mainnet-live. Mixing networks in one workspace. |

---

## Feature capability matrix (UI source of truth)

Status: ✅ live+PIT-proven · ⚠️ live official / PIT-unproven · ❌ absent or forbidden · 🧪 implement in product (not yet)

| Feature | Testnet | Mainnet | 0G component | External dependency | Real evidence | User action | Status |
|---|---|---|---|---|---|---|---|
| 0G Chain | ✅ 16602 | ✅ 16661 | RPC | Official RPCs | eth_chainId M18 | Select network | 🧪 |
| Direct TeeML | ⚠️ Omni-7b ack; no glm-5.2 | ✅ glm-5.2 | InferenceServing + Ledger | User Ledger top-up | M17 committee; M18 Galileo getAllServices | None (infra) | 🧪 |
| E2EE | ⚠️ protocol same, unproven | ✅ VerifyE2EE | 0g-pc-e2ee | HPKE provider pubkey | M17 scheme `zg-sig-v1/e2ee-ct` | None | 🧪 |
| Storage | ⚠️ addrs live | ✅ `--proof` | Flow + indexer | User 0G for fees | S6 mainnet | None | 🧪 |
| Storage proof | ⚠️ not PIT-run on Galileo | ✅ | Go client | `--encryption-key` `0x` | S6 | None | 🧪 |
| Agentic ID mint | ❌ PitDeskID code 0; Foundation mint possible | ✅ PitDeskID | ERC-7857 | User wallet gas | M17 mint tx | YOU mint | 🧪 |
| authorizeUsage | ⚠️ after mint | ✅ | IERC7857Authorize | Owner wallet | M17 auth tx | YOU authorize | 🧪 |
| revoke | ⚠️ after mint | ✅ | same | Owner wallet | M17 revoke tx | YOU revoke | 🧪 |
| iTransfer | ⚠️ official attestor; no PIT tx | ❌ attestor code 0 | Foundation AgenticID | Attestor service | config 16602; Aristotle code 0 | Lab only | **do not fake** |
| iClone | ⚠️ same | ❌ | IERC7857Cloneable | Attestor | expected revert 16661 | Lab only | **do not fake** |
| ERC-8004 | ✅ registries | ✅ | Identity/Reputation | Reporter ≠ owner | M15 register | YOU register | 🧪 |
| Hyperliquid | ✅ testnet UI | ✅ mainnet | HL API | YOUR account | Chrome both apps M18 | Connect | 🧪 |
| Hyperliquid session | 🧪 keychain | 🧪 keychain | extraAgents | OS keychain | M15 allowlist | YOU approve | 🧪 |
| Researcher | ⚠️ Omni unproven | ✅ glm-5.2 | Direct TeeML | Funded provider | M17 S1–S3 | Ask / Watch | 🧪 |
| Challenger | ⚠️ same SKU | ✅ same provider | Direct TeeML | Envelope split | M17 | Automatic | 🧪 |
| Risk | ⚠️ same SKU | ✅ glm-5.2; 0GM-SIA unfunded | Direct TeeML | Optional 0GM fund | M17 | Automatic | 🧪 |
| Opportunity Watch | 🧪 | 🧪 | Host + HL public | Policy | Spec only | Enable | 🧪 |
| Policy engine | 🧪 | 🧪 | Host + later contract | Pin tx | Spec | YOU set | 🧪 |
| Durable replay ledger | 🧪 | 🧪 | SQLite | Workspace file | Spec | None | 🧪 |
| Memory | 🧪 | 🧪 | Storage + keychain | Per-ws key | S6 pattern | None | 🧪 |
| Calibration | 🧪 | 🧪 | Stats + 8004 scores | Reporter key | Spec | None | 🧪 |
| MCP | 🧪 RO | 🧪 RO | pit MCP | No session | Spec | None | 🧪 |
| SDK | 🧪 | 🧪 | pit SDK | No browser secrets | Spec | None | 🧪 |
| CLI | 🧪 | 🧪 | pit CLI | Keychain | Spec | Commands | 🧪 |
| Desktop | 🧪 | 🧪 | Tauri | KEPT visual | Spec | Flagship | 🧪 |
| Web | 🧪 no sign | 🧪 no sign | wagmi | No session | Spec | View/verify | 🧪 |
| Verify | 🧪 | 🧪 | /verify | Explorer | Spec | Open proof | 🧪 |
| Router `sk-` | ❌ book | ❌ book | Router | Dashboard key | Docs + Chrome API Keys page | **Do not create for PIT** | Forbidden |

---

## Phase template (every phase below uses this)

Goal · User-visible · Files · Dependencies · Env vars · External resources · Contracts/APIs · Commands · Tests · Security tests · Real network tests · Testnet behavior · Mainnet behavior · Expected · Evidence artifacts · Explorer · Failure · Fallback · Kill · Rollback · Consulted · Judge axis

---

## Phase 0 — Research revalidation

- **Goal:** Re-fetch live facts. Do not implement product.
- **User-visible:** none.
- **Files:** append `history.md` only if facts moved.
- **Deps:** RPC, Router, HL info, agenticid config, clone HEAD.
- **Env:** `PIT_CHAIN_RPC_URL`, `PIT_TESTNET_RPC_URL`, `PIT_ROUTER_URL`, `PIT_TESTNET_ROUTER_URL`, `PIT_AGENTICID_CONFIG_URL`. No `sk-`.
- **Resources (READ FIRST):** docs.0g.ai, build.0g.ai, https://pc.0g.ai/models, https://pc.testnet.0g.ai/models, https://router-api.0g.ai/v1/models, https://agenticid.0g.ai/config, https://hyperliquid.gitbook.io/hyperliquid-docs, 03 file EOF.
- **Commands:**
  - `cast chain-id --rpc-url https://evmrpc.0g.ai` → 16661
  - `cast chain-id --rpc-url https://evmrpc-testnet.0g.ai` → 16602
  - `curl -s https://router-api.0g.ai/v1/models`
  - `curl -s https://router-api-testnet.integratenetwork.work/v1/models`
  - `curl -s https://agenticid.0g.ai/config`
  - HL `meta` + `spotMeta` testnet+mainnet
- **Testnet:** Galileo catalog may be 2 Router models; Direct TeeML chat = Omni. Do not copy mainnet SKUs.
- **Mainnet:** Private TeeML 5; glm-5.2 still required for committee.
- **Tests:** none product.
- **Explorer:** chainscan + chainscan-galileo `getCode` attestor.
- **Fallback:** none.
- **Kill:** if TeeML glm-5.2 gone on **mainnet**; if Storage indexer dead; if HL exchange 500 for >1h; if Galileo Omni ack flipped false **and** no other TeeML chat — disable testnet sealed ask, do not Router.
- **Rollback:** n/a.
- **Consulted:** 03 judge pattern §1; Lumen honesty; Axiom mock penalty.
- **Judge:** real primitives still exist.
- **Expected:** DA still epoch 0; attestor still 16602 (**M18 reconfirmed**); HL min $10.

**Exit:** FACT table unchanged or history note. Then Phase 1.

---

## Phase 1 — Multi-user identity model

- **Goal:** `Workspace` + SIWE bind. Fixtures ≠ users.
- **User-visible:** none yet (library).
- **Files:** `pit/internal/identity/`, `pit/internal/workspace/`, schema `workspaces(id, evm, created)`.
- **Deps:** none on-chain yet.
- **Resources:** EIP-4361 SIWE; ENVIRONMENT.md classes A–F.
- **Commands:** `go test ./internal/identity/... ./internal/workspace/...`
- **Tests:** two workspaces cannot `Get(memory)` each other; UUID guess returns not-found not decrypt-error leak.
- **Security:** tenant ID injection; bearer mixup; object ID enumeration.
- **Network:** none.
- **Expected:** `workspaceId` UUID v4; all queries require `ws`.
- **Evidence:** `_gate/evidence/p1_identity.txt`
- **Kill:** any global `PIT_MASTER_ADDRESS` read in product path.
- **Rollback:** drop tables.
- **Consulted:** NeoSoul tenant isolation (concept); 4lpha single-vault anti-pattern.
- **Judge:** concurrency / not a one-wallet demo.

---

## Phase 2 — Wallet onboarding (library + later UI)

- **Goal:** Connect + SIWE. Never collect keys.
- **User-visible (when UI exists):** “PIT never asks for your private key.”
- **Files:** `pit/internal/wallet/`, desktop `apps/desktop` later; web wagmi.
- **Resources:** viem/wagmi; Rabby; EIP-1193; KEPT `WalletGate.tsx` **structure only**.
- **Commands:** unit tests with mocked provider **only in tests**, not product.
- **Tests:** declined signature → named state `SIGNATURE_DECLINED`; wrong chain → `WRONG_NETWORK`.
- **Security:** no mnemonic field in any form.
- **Network:** 16661 `eth_chainId`.
- **Kill:** seed input.
- **Consulted:** KEPT WalletGate; Hanami connect.
- **Judge:** browser-verifiable identity.

---

## Phase 3 — Hyperliquid connection

- **Goal:** Identify HL account; show **spot** USDC; pick testnet/mainnet **bound to `PIT_NETWORK`**.
- **Files:** `pit/internal/hl/account.go` (promote `_gate/hl_*.py` knowledge, not the Python).
- **Env:** `PIT_HL_INFO_URL` / `PIT_HL_EXCHANGE_URL` XOR testnet twins. Never `PIT_MASTER_ADDRESS` as product user.
- **Resources (READ FIRST):** https://hyperliquid.gitbook.io/hyperliquid-docs ; account-abstraction-modes; info `clearinghouseState`, `spotClearinghouseState`, `userFees`.
- **APIs:** `https://api.hyperliquid.xyz/info` and testnet twin.
- **Testnet:** Chrome `app.hyperliquid-testnet.xyz` live; faucet `/drip`. Min $10 still applies to test notional.
- **Mainnet:** Chrome `app.hyperliquid.xyz` live.
- **Commands:** `go test ./internal/hl/...`
- **Tests:** unified `accountValue=0` + spot USDC >0 → UI state `FUNDED_SPOT` not `UNFUNDED`.
- **Network:** fixture wallet optional; **also** empty account path.
- **Fallback:** none (no mock book).
- **Kill:** treating clearinghouse 0 as unfunded when spot has USDC; mixing HL testnet with Aristotle in the same workspace UI.
- **Consulted:** HL docs; M15 FACT.
- **Judge:** real venue, not a wrapper that lies about balances.

---

## Phase 4 — Session generation

- **Goal:** Keygen in OS keychain; user-signed `approveAgent`; TTL ≤1h; name `PIT-{ws8}`.
- **Files:** `pit/internal/session/`, Tauri stronghold / keytar.
- **Resources:** gitbook nonces-and-api-wallets; extraAgents.
- **APIs:** exchange `approveAgent`.
- **Commands:** `go test ./internal/session/...` (no print of keys).
- **Tests:** export-to-JSON denied; web build has zero session types.
- **Security:** key never in logs; never `.env`; never localStorage.
- **Network:** testnet approve + `extraAgents` contains name; expire path.
- **Expected:** permissions card data: order✓ cancel✓ withdraw✗ leverage✗.
- **Kill:** session in Vercel/env.
- **Consulted:** NeoTrade signing-gateway (structure); M15 protocol hole disclosure.
- **Judge:** security.

---

## Phase 5 — Policy engine

- **Goal:** Human cards → executable checks; pin hash; re-check after async.
- **Files:** `pit/internal/policy/`, later `PitPolicy.sol`.
- **Resources:** FINAL policy table.
- **Commands:** `go test ./internal/policy/...`
- **Tests:** each card has a failing fixture; cooldown; daily loss; calibration floor; liquidity; uncertainty.
- **Security:** LLM JSON cannot raise maxClip.
- **Kill:** policy only in prompt.
- **Consulted:** 4lpha PolicyVault **do not copy**; keep host-side.
- **Judge:** security + explainability.

---

## Phase 6 — Sealed committee

- **Goal:** 3 Direct TeeML roles over **book JSON**, not PONG. **M17 gate did this on Aristotle.** Product wires the same path into `pit/internal/compute/`.
- **Files:** `pit/internal/compute/` wrapping vendored `0g-pc-e2ee` HEAD `0ba9fb0`. Promote `_gate/committee/main.go` (`SealRequest` + `OpenResponse` + `VerifyE2EE`; deny `router-api`; `-tamper-ciphertext`; `-wrong-signer`).
- **Env:** serving+ledger for selected network; **no** `PIT_PC_ROUTER_API_KEY`. User Ledger subaccount (prod) or fixture (gate).
- **Resources (READ FIRST):** inference.md Direct; router comparison + privacy + authentication (to know what to **refuse**); `0g-pc-e2ee`; `0g-serving-broker` `doc/api-token.md`; SDK `getHeader` / `createApiKey`; live catalog both networks.
- **APIs:** Direct `-onchain` getService; HPKE; `VerifyE2EE`. `Allow-Fallbacks: false`. `transferFund` serviceName `inference-v1.0`. Ledger min reserve ~1.0 0G + unsettled fees or POST 400.
- **Testnet:** Galileo TeeML chat = `qwen/qwen2.5-omni-7b` provider `0xa48f…7836` signer `0x83df…08cF` (teeAck true). **No glm-5.2.** PIT has not VerifyE2EE’d Omni — first product test on Galileo must prove Seal+Verify before enabling sealed ask. Router Omni is labeled TeeTLS — ignore Router label.
- **Mainnet:** glm-5.2 provider `0x7DCF…87D` teeSigner `0xA46E…46B9`. Optional later: 0GM / 0GM-SIA **after** Ledger subaccounts exist.
- **Commands:** `go test ./internal/compute/...`; live: `_gate/m17_committee.py` + `committee.exe`.
- **Tests:** Router path **compile-fail** on book; missing VerifyE2EE → abort. Engine: LLM cannot set sz; `challenger_killed` / `below_min_notional` / `kill_switch` / `leverage_above_policy`.
- **Explorer:** InferenceServing `getService` teeSigner vs recovered signer.
- **Fallback:** **none**. Hanami-style Router fallback is a kill.
- **Expected:** scheme `zg-sig-v1/e2ee-ct`; signer = teeSigner. S1 may honestly deny. Do not retry until a model says yes.
- **Evidence:** `PIT/_gate/evidence/m17_committee/` (hashes, zg-res-keys, sanitized role JSON; prompts redacted).
- **Kill:** plaintext Router; labeled fallback; faking independent models on one provider; copying glm-5.2 onto Galileo.
- **Consulted:** Knole E2EE; Lumen honesty; Hanami Direct-off; Axiom software TEE anti-pattern.
- **Judge:** real 0G + attestation.

---

## Phase 7 — Deterministic forecast engine

- **Goal:** Structured forecast; size from sizer not LLM.
- **Files:** `pit/internal/engine/`, promote `_gate/hlcore` sizer.
- **Commands:** `go test ./internal/engine/...`
- **Tests:** LLM `sizeUsd` ignored; invalidation required; `market` discriminator `hyperliquid:perp:*` | `0g:univ3:*` | `none`.
- **Kill:** AI probability as p.
- **Consulted:** NeoSoul rubric (concept); TRAIDE journal.
- **Judge:** agent quality.

---

## Phase 8 — Opportunity discovery (Watch)

- **Goal:** Policy-eligible candidates; no auto-sign.
- **Files:** `pit/internal/watch/`
- **APIs:** HL `meta` + funding + `l2Book`/`impactPxs` for allowlisted coins only (ETH, BTC, optional SOL).
- **Commands:** `go test ./internal/watch/...`
- **Tests:** candidate failing policy never appears; no ghost cards; closed app = no cloud watch.
- **Network:** live books.
- **Kill:** fake opportunities; hosted watcher with session key.
- **Consulted:** Axiom-trade terminal (do not clone UI junk).
- **Judge:** frequency + utility.

---

## Phase 9 — Execution gateway

- **Goal:** Fail-closed `order`|`cancel` only; preview bind; tick round.
- **Files:** `pit/internal/exec/` from `_gate/hlcore`.
- **APIs:** HL exchange.
- **Commands:** `go test ./internal/exec/...`
- **Tests:** `updateLeverage`/`withdraw3`/`sendAsset`/`approveAgent` denied (S3 matrix).
- **Network:** testnet dust then mainnet dust (min $10, cancel).
- **Kill:** any extra action.
- **Consulted:** M15; 4lpha allowlists.
- **Judge:** real txs + security.

---

## Phase 10 — Replay / idempotency persistence

- **Goal:** Exactly-once durable SQLite per workspace.
- **Files:** `pit/internal/ledger/`
- **Consulted:** 4lpha `usedActionHashes`; NeoSoul `insertIfAbsent` — **implement PIT-native**.
- **Tests (required):**
  - same action twice
  - concurrent same action
  - restart between preview and sign
  - crash after sign
  - crash after receipt
  - retry after timeout
  - duplicate click
  - refresh (desktop recovers; web cannot sign)
- **Kill:** in-memory map in production binary.
- **Judge:** concurrency.

---

## Phase 11 — 0G Storage

- **Goal:** Official client encrypt+upload+download `--proof`.
- **Files:** `pit/internal/storage/` wrapping `0g-storage-client`.
- **Env:** `PIT_STORAGE_INDEXER` + `PIT_STORAGE_FLOW` or testnet twins. Per-user memory key in keychain, never global `PIT_MEMORY_KEY` in product.
- **Resources (READ FIRST):** `D:\route\0g\research\0g-storage-client` HEAD `3d953af`; docs storage SDK (do not use TS for proofs).
- **Testnet:** Flow `0x22E0…5296`; indexer `https://indexer-storage-testnet-turbo.0g.ai`. PIT `--proof` on Galileo = not yet run.
- **Mainnet:** Flow `0x62D4…526`; indexer turbo. S6 proven.
- **Trap:** `--encryption-key` must be `0x`-prefixed; TS SDK **forbidden** for proofs.
- **Tests:** match; wrong root fail; wrong key fail; **cross-workspace key fail**.
- **Fallback:** none.
- **Kill:** custom AES; TS proof claims.
- **Judge:** Storage primitive.

---

## Phase 12 — Encrypted memory

- **Goal:** Typed memory kinds; namespace `ws/{id}/`.
- **Files:** `pit/internal/memory/`
- **Tests:** kind tags; LLM never receives memory key; A cannot list B prefixes.
- **Kill:** one global `PIT_MEMORY_KEY` in product.
- **Judge:** memory + privacy.

---

## Phase 13 — ERC-8004

- **Goal:** User-owned Identity register; **separate** feedback client.
- **Contracts:** Mainnet Identity `0x8004A169FB4a3325136EB29fA0ceB6D2e539a432`; Reputation `0x8004BAa17C55a88189AE136b182e5fdA19dE9b63`. Testnet Identity `0x8004A818BFB912233c491871b3d84c89A494BD9e`; Reputation `0x8004B663056A597Dffe9eCcC1965A193B7388713` (M18 code 130/130).
- **Resources (READ FIRST):** https://docs.0g.ai/developer-hub/building-on-0g/agentic-id/erc8004
- **Trap:** owner `giveFeedback` reverts (M15). Reporter ≠ owner.
- **Do not use** NeoSoul `0x8004Ae53…`.
- **Testnet / Mainnet:** both registries live; IDs are **not** portable across chains.
- **Tests:** Foundry fork + live register from **user** wallet (or second fixture, not “the PIT owner”).
- **Kill:** shared 8004 agent across users.
- **Judge:** 8004.

---

## Phase 14 — ERC-7857

- **Goal:** `PitDeskID` production interface IDs `0x2afbede9` / `0xdf597d99` / `0x74f8628b` / `0x80ac58cd`. **M17 gate deployed** `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` on **16661 only**. Promote `_gate/desk7857` into `contracts/` (verify source; optional CREATE2). Deploy a **separate** Galileo DeskID if testnet mint is required.
- **Env:** `PIT_DESK_ID_CONTRACT` (pin live address — currently empty in fixture `.env`) vs `PIT_TESTNET_DESK_ID_CONTRACT`.
- **Resources (READ FIRST):** https://docs.0g.ai/ai-context?_highlight=itransfer#agentic-id-formerly-inft ; `D:\route\0g\research\0g-agentic-id` HEAD `b8f4845`; SDK `CHAIN_ID = 16602`; `GET https://agenticid.0g.ai/config`; Lumen OracleNotLive; Chrome `https://agenticid.0g.ai/config/`.
- **Mainnet iTransfer:** **IMPOSSIBLE** — revert `AttestorNotOnAristotle`. UI: “Transfer not live on mainnet.” Do **not** upgrade to LIVE without official repo + explorer + config `chain_id:16661` + a tx that changes `ownerOf`.
- **Testnet iTransfer:** Foundation AgenticID + verifier **bytecode live**. Official `iTransferFrom` path exists. **PIT has not executed a transfer.** UI: ⚠️ until `ownerOf` changes in a PIT tx. Do not build a fake attestor. ERC-721 fallback is **not** official for production DeskID (`ERC7857UseITransferFrom`). Foundation §2.4 seal-bound agents **disable** iTransfer even on Galileo.
- **Lifecycle to implement in product:** mint → ownerOf → encrypted URI/hash → authorizeUsage → isAuthorized → revoke. Host must check `isAuthorized` before sealed inference for that desk.
- **Tests:** mint/own/authorizeUsage; unauthorized authorize reverts; iTransfer/iClone expected revert on 16661; transferFrom reverts `ERC7857UseITransferFrom`; supportsInterface. Optional Galileo: one real `iTransferFrom` **or** honest revert — record either. Gate: Foundry 7/7 + live txs (see RESOURCES / history M17).
- **Fallback:** none. No wrapper labeled as iTransfer.
- **Kill:** fake iTransfer success; Talos always-valid verifier; Knole custom IDs; claiming Galileo transfer LIVE without PIT tx.
- **Consulted:** Lumen `OracleNotLive()`; Foundation `DEPLOYMENT.md` §2.4.
- **Judge:** standards-verifiable 7857, honest.

---

## Phase 15 — Calibration + Strategy Health

- **Goal:** Brier/ECE; Health card; empty state if N<30.
- **Files:** `pit/internal/calib/`
- **Tests:** never display fake 72%; overconfidence flag when ECE exceeds pin.
- **Kill:** LLM-written calibration.
- **Judge:** forecasting + outcome learning.

---

## Phase 16 — MCP (read-only)

- **Goal:** markets, opportunities, forecast, status, card, verify.
- **Files:** `pit/cmd/mcp/`
- **Tests:** tools named `authorize`/`order` **do not exist**; schema tests.
- **Kill:** session export tool.
- **Judge:** MCP without being a trading backdoor.

---

## Phase 17 — SDK

- **Goal:** Typed TS/Go client; browser build strips exec.
- **Files:** `sdk/`
- **Tests:** `esbuild` browser bundle grep fails on `approveAgent` secret types.
- **Kill:** duplicate executor in JS.
- **Judge:** SDK.

---

## Phase 18 — CLI

- **Goal:** `pit init|login|policy|ask|opportunities|forecast|preview|authorize|cancel|status|resolve|card|verify|kill`
- **Files:** `pit/cmd/pit/`
- **Tests:** `authorize` requires TTY confirm; piped yes is not enough without `--i-understand`.
- **Kill:** default yes.
- **Judge:** CLI.

---

## Phase 19 — Contracts

- **Goal:** Foundry: Policy, Forecasts, Receipts, Memory, DeskID.
- **Files:** `contracts/`
- **Commands:** `forge test -vv`; `forge coverage`; fuzz replay.
- **Network:** deploy 16661 from **deployer fixture**; users interact with **their** wallets. DeskID: start from `_gate/desk7857` (already live semantics); do not invent a second interface set.
- **Evidence:** chainscan txs.
- **Kill:** deploy-only without tests (Axiom).
- **Judge:** Solidity tests + mainnet.

---

## Phase 20 — Desktop (flagship)

- **Goal:** Tauri; KEPT tokens; authorize locally.
- **Files:** `apps/desktop/` — read KEPT `guide.css`, `Hero.tsx`, `WalletGate.tsx`, `StartFlow.tsx`, `namedStates.ts`, `keptGuide.tsx` **before** any screen.
- **Resources:** `D:\route\Flare\kept`; live https://kept-ruby-five.vercel.app/
- **User-visible:** landing coral; split sign-in; get-started cards; PIT ring SVG; Watch home; permissions card; real phase labels.
- **Tests:** reduced-motion; no fake spinner timers.
- **Kill:** generic dashboard; fake charts.
- **Judge:** genuine product + UX.

---

## Phase 21 — Web

- **Goal:** Onboarding, research, cards, policy txs, `/verify`. **No HL session.**
- **Files:** `apps/web/`
- **Tests:** bundle has no session private key type; lighthouse later optional.
- **Kill:** trading from Vercel.
- **Judge:** browser attestations.

---

## Phase 22 — Playwright

- **Goal:** Deterministic UI vs stubbed **public** APIs where required; **never** stub VerifyE2EE success.
- **Files:** `apps/desktop/e2e/`
- **Tests:** onboarding 1–11; Watch empty; named errors.
- **Kill:** tests that click AUTHORIZE against a mock exchange in CI **without** a labeled `// MOCK TEST HARNESS` and not counting as live evidence.
- **Judge:** tests.

---

## Phase 23 — Real Chrome E2E

- **Goal:** Real wallet (or fixture) through Chrome DevTools MCP / Playwright against live testnet.
- **Files:** evidence screenshots + traces gitignored if they contain addresses only — OK.
- **Tests:** connect → policy → session → research-only **and** dust order.
- **Kill:** claiming E2E without UI.
- **Judge:** browser + real txs.

---

## Phase 24 — Multi-user E2E

- **Goal:** Two fixtures (A and B). A cannot see B.
- **Commands:** two OS users or two keychain namespaces in tests.
- **Tests:** IDOR on forecast id; storage prefix; session; websocket.
- **Kill:** shared SQLite without `workspaceId`.
- **Judge:** concurrency + security.

---

## Phase 25 — Security red-team

- **Goal:** Automated matrix from FINAL §red team.
- **Files:** `pit/internal/redteam/`
- **Tests:** S3 allowlist regression; preview bind; kill switch; expired session; tenant mixup.
- **Kill:** known hole left undocumented.
- **Judge:** security.

---

## Phase 26 — Dual-network deployment

- **Goal:** Production contracts on **16661**. Laboratory DeskID (+ optional policy/receipts) on **16602**. HL dust on the matching venue.
- **Testnet:** Deploy Galileo PitDeskID; do **not** point mainnet UI at it. Optional Foundation iTransfer experiment stays behind TESTNET badge.
- **Mainnet:** Verify PitDeskID source; iTransfer remains disabled in UI.
- **Env:** both address sets in ENVIRONMENT.md; `PIT_NETWORK` selects one.
- **Commands:** `forge script` + verify source on chainscan **and** chainscan-galileo.
- **Evidence:** addresses in RESOURCES + history.
- **Fallback:** none.
- **Kill:** unverified bytecode; one address used on both chains; testnet transfer control visible on mainnet.
- **Judge:** mainnet + honest testnet lab.

---

## Phase 27 — Real-user onboarding

- **Goal:** A wallet **not** `0xBDfC…0034` completes the **14-beat** journey (includes network select). Copy: YOU / YOUR WALLET / YOUR TRADING ACCOUNT / YOUR SESSION / YOUR POLICY / YOUR MONEY.
- **Testnet first** then **mainnet dust** for the same user identity (separate workspaces).
- **Kill:** hardcoded workspace IDs; preloaded DB required; asking for Router API key or seed.
- **Judge:** genuine in-window product.

---

## Phase 28 — Demo rehearsal

- **Goal:** Spine in FINAL “Demo spine” timed; 0G copy honest.
- **Kill:** fake progress to hit a timebox.
- **Judge:** reproducible proof.

---

## Phase 29 — Final verification

- **Goal:** Gate report in `history.md`. Score honest. Competitor checks re-run.
- **Kill:** UNKNOWN marked PASS.
- **Judge:** all.

---

## Dependency order (do not skip)

`0 → 1 → 2,3,5 (parallel after 1) → 4 (needs 3) → 6,7,10,11 (parallel) → 8 (needs 5+7) → 9 (needs 4+5+7+10) → 12 (needs 11) → 13,14,19 → 15 → 16,17,18 → 20,21 → 22 → 23,24,25 → 26 → 27,28,29`

---

## Exact start command (Phase 0 then 1)

```powershell
cd D:\route\0g\PIT
# Phase 0
cast chain-id --rpc-url https://evmrpc.0g.ai
cast chain-id --rpc-url https://evmrpc-testnet.0g.ai
curl -s https://router-api.0g.ai/v1/models
curl -s https://router-api-testnet.integratenetwork.work/v1/models
curl -s https://agenticid.0g.ai/config
# Phase 1 — first product code
mkdir -p pit/internal/identity pit/internal/workspace
go test ./...
```

Do **not** run `npm create` dashboards first.

---

## Per-phase competitor map

| Phase | Consult (local clone or 03) | Take | Refuse |
|---|---|---|---|
| 1,24 | NeoSoul isolation notes | tenant key | public economy |
| 5,9 | 4lpha | hash-once idea | PolicyVault, MockUSDC |
| 6 | Knole, Lumen, Hanami | E2EE, honesty | Router fallback (`OG_DIRECT_ENABLED`) |
| 10 | 4lpha, NeoSoul gateway | insertIfAbsent idea | copy their bytecode |
| 14 | Lumen | OracleNotLive tone | fake transfer |
| 20–21 | KEPT | visual system | circles/XRPL |
| 8 | Orchestra Watcher | background scan idea | Groq, AUTO_EXECUTE, Uniswap book |
| 21 `/verify` | Hanami | recompute UX | their stack |
| 9 | Maneki | HL connector caution | 0G sticker |
| — | Axiom, Talos, diversifi | — | mocks, Groq, env-gate |
| — | TRAIDE, NEXUS, cognivern | merkle/receipt density | Galileo-as-mainnet |
| — | OptimIEra | — | unused unless needed |

---

## Failure / rollback (global)

- Any mock in `pit/` production: revert commit.
- `PIT_ALLOW_FALLBACKS=true`: process exit.
- Lost session key: user re-approves; funds stay on master.
- Bad contract: pause pin; do not migrate vaults (there are none).

---

## Evidence root

`PIT/_gate/evidence/` (gitignored secrets) + `history.md` (public facts, truncated addresses).
