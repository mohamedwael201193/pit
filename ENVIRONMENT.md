# PIT — Environment

**Six classes. Mixing them is a kill.**

| Class | Lives where | Product meaning |
|---|---|---|
| **A. GLOBAL INFRASTRUCTURE** | Operator `.env` on the deploy / CI machine only | Public RPCs, official contract addresses, catalog URLs |
| **B. USER-SCOPED CREDENTIALS** | OS keychain + encrypted workspace DB | Wallet, HL session, memory key, per-user Ledger spend |
| **C. DEVELOPMENT FIXTURES** | `PIT/.env` gitignored | Gate wallet only. Not “the owner of PIT.” |
| **D. TESTNET ONLY** | Selected when `PIT_NETWORK=testnet` | Galileo `16602` + Hyperliquid testnet |
| **E. MAINNET ONLY** | Selected when `PIT_NETWORK=mainnet` | Aristotle `16661` + Hyperliquid mainnet |
| **F. OPTIONAL** | Off by default | Router keys, x402, labeled 0G dust, Privy |

Fail loud if a required class is missing. **Never** env-gate to mocks (`PIT_ALLOW_FALLBACKS` must be `false` or the process exits).

**Do not print private keys.** This file lists **names and public addresses only.** Never paste Router `sk-` / `mk-` values, Direct `app-sk-` tokens, HL session keys, or hex private keys into chat, git, screenshots, `history.md`, or reports.

**Network selector (product):** `PIT_NETWORK` is `mainnet` or `testnet`. It selects the **0G chain + matching Hyperliquid venue together**. Do not silently mix Aristotle compute with HL testnet in the UI. CI may run a labeled fixture mix; the product never hides it.

---

## Identity mapping (product)

```
workspaceId     = uuid; bound to user EVM address via SIWE (desktop/web) or CLI login
evmAddress      = connected wallet (user-owned) — YOUR WALLET
hlMaster        = the HL account YOU prove by approveAgent — YOUR TRADING ACCOUNT
hlSessionId     = ephemeral agent address; key in OS keychain under
                  pit/<network>/<workspaceId>/session/<sessionId> — YOUR SESSION
memoryKey       = 32-byte; OS keychain pit/<network>/<workspaceId>/memory
policyId        = hash of YOUR POLICY; on-chain pin uses this workspace
deskId (7857)   = minted to evmAddress on the selected 0G chain
agentId (8004)  = registered by evmAddress (owner ≠ feedback client)
```

User A cannot open User B’s keychain namespace. Object IDs are UUID + workspace prefix. Storage object keys = `{network}/ws/{workspaceId}/...`. Guessing another UUID must decrypt-fail.

**No global `PIT_MASTER_ADDRESS` as the product user.** No shared HL session. No shared memory. User #1, #2, and #1000 each get their own workspace.

---

## Compute authentication — FACT (M18)

Two official 0G paths. PIT’s confidential book uses **only Direct**.

| Path | Request goes to | Credential | Router sees book? | E2EE / caller-byte bind | On-chain teeSigner | PIT private book |
|---|---|---|---|---|---|---|
| **A. Router** | `https://router-api.0g.ai/v1` (mainnet) or `https://router-api-testnet.integratenetwork.work/v1` (testnet) | Dashboard `sk-` (inference) + optional `mk-` (admin) | **Yes** — docs: prompts handled in memory to route/bill. Privacy mode still trusts the Router. `verify_tee:true` is a Router-asserted boolean, not VerifyE2EE. | No raw signature to host | No | **Forbidden** |
| **B. Direct / Advanced** | Provider URL from `getService` | Wallet-signed `Bearer app-sk-<base64(rawMessage\|signature)>`. Ephemeral tokenId **255** (24h, SDK `getHeader`) or persistent tokenId **0–254** (`createApiKey`) | **No intermediary** (pc.0g.ai/sdk: “Pick the route yourself… no gateway”) | HPKE SealRequest + VerifyE2EE scheme `zg-sig-v1/e2ee-ct` | Yes — recovered signer must equal on-chain `teeSigner` | **Required** |

**pc.0g.ai Dashboard → API Keys** creates **Router `sk-` keys only** (Chrome M18: “Generate keys for programmatic inference access via the OpenAI-compatible API.” Trust modes Standard / Verified / Private are **Router routing tiers**, not Direct). Management Keys are `mk-` for account admin. Chrome Console showed existing truncated `sk-` prefixes on the logged-in wallet; **PIT did not create a key this pass and must not copy those values into `PIT/.env`.**

**Hanami anti-pattern:** Direct is optional (`OG_DIRECT_ENABLED=false`) and **falls back to Router**. PIT must compile-fail / exit, never fall back.

**Does PIT need a 0G Private Computer API key?** See decision at the bottom. Short: **Direct private-book path does not.** Do not ask the user for `sk-`/`mk-` unless a future **unlabeled public catalog** tool is explicitly scoped — never for the book.

---

## A. GLOBAL INFRASTRUCTURE

Public, non-secret, shared by all users. Stored in operator `.env` / config files. Rotate: change URL if official docs move.

| Name | Required? | Environment | Who owns it? | How obtained? | What it grants | Where stored | Rotation | Never expose? |
|---|---|---|---|---|---|---|---|---|
| `PIT_NETWORK` | Yes | A | Operator / user toggle | Product UI: Testnet or Mainnet | Selects D vs E address set | local DB + process env | N/A | No |
| `PIT_ALLOW_FALLBACKS` | Yes | A | Operator | Must be `false` | Process exits if not false | `.env` | N/A | No |
| `PIT_CHAIN_RPC_URL` | Yes | E default | 0G | Official docs | Aristotle JSON-RPC | `.env` | If RPC moves | No |
| `PIT_CHAIN_ID` | Yes | E default | 0G | Official | Guard `16661` | `.env` | N/A | No |
| `PIT_TESTNET_RPC_URL` | Yes if testnet | D | 0G | docs testnet-overview | Galileo JSON-RPC | `.env` | If RPC moves | No |
| `PIT_TESTNET_CHAIN_ID` | Yes if testnet | D | 0G | Official | Guard `16602` | `.env` | N/A | No |
| `PIT_EXPLORER_BASE` | Yes | E | 0G | Official | Mainnet links | `.env` | N/A | No |
| `PIT_TESTNET_EXPLORER` | Yes if testnet | D | 0G | Official | Galileo links | `.env` | N/A | No |
| `PIT_ROUTER_URL` | Yes | E | 0G | Official | **Catalog only** — never book inference | `.env` | N/A | No |
| `PIT_TESTNET_ROUTER_URL` | Yes if testnet | D | 0G | Official | Galileo catalog only | `.env` | N/A | No |
| `PIT_SERVING_CONTRACT` | Yes | E | 0G SDK | `0g-compute-ts-sdk` `constants.ts` mainnet.inference | Direct `getService` | `.env` | If 0G redeploys | No |
| `PIT_LEDGER_CONTRACT` | Yes | E | 0G SDK | `constants.ts` mainnet.ledger | Per-provider credit | `.env` | If 0G redeploys | No |
| `PIT_TESTNET_SERVING_CONTRACT` | Yes if testnet | D | 0G SDK | `constants.ts` testnet.inference | Galileo Direct | `.env` | If 0G redeploys | No |
| `PIT_TESTNET_LEDGER_CONTRACT` | Yes if testnet | D | 0G SDK | `constants.ts` testnet.ledger | Galileo credit | `.env` | If 0G redeploys | No |
| `PIT_STORAGE_INDEXER` | Yes | E | 0G | mainnet-overview | Turbo indexer | `.env` | N/A | No |
| `PIT_STORAGE_FLOW` | Yes | E | 0G / prior FACT | Mainnet Flow | Uploads | `.env` | N/A | No |
| `PIT_TESTNET_STORAGE_INDEXER` | Yes if testnet | D | Official clones + docs | Galileo turbo indexer | Uploads | `.env` | N/A | No |
| `PIT_TESTNET_STORAGE_FLOW` | Yes if testnet | D | testnet-overview | Galileo Flow | Uploads | `.env` | N/A | No |
| `PIT_STORAGE_CLI` | Yes for proofs | A | Operator | Official Go client binary | `--proof` download | path | Rebuild | No |
| `PIT_8004_IDENTITY` | Yes | E | 0G | Official 8004 docs | Register desk | `.env` | N/A | No |
| `PIT_8004_REPUTATION` | Yes | E | 0G | Official 8004 docs | Feedback | `.env` | N/A | No |
| `PIT_TESTNET_8004_IDENTITY` | Yes if testnet | D | 0G docs | Galileo Identity | Register | `.env` | N/A | No |
| `PIT_TESTNET_8004_REPUTATION` | Yes if testnet | D | 0G docs | Galileo Reputation | Feedback | `.env` | N/A | No |
| `PIT_POLICY_CONTRACT` | After Phase 19 | selected chain | PIT | Deploy | Policy pin | `.env` | Redeploy | No |
| `PIT_RECEIPTS_CONTRACT` | After Phase 19 | selected chain | PIT | Deploy | Receipts | `.env` | Redeploy | No |
| `PIT_FORECASTS_CONTRACT` | After Phase 19 | selected chain | PIT | Deploy | Forecasts | `.env` | Redeploy | No |
| `PIT_DESK_ID_CONTRACT` | Yes on mainnet after gate | E | PIT | Gate-issued | ERC-7857 DeskID | `.env` | CREATE2 later | No |
| `PIT_TESTNET_DESK_ID_CONTRACT` | After Galileo deploy | D | PIT | Not issued (Galileo code len 0 at gate address) | Testnet DeskID | `.env` | Deploy | No |
| `PIT_HL_INFO_URL` | Yes | E | Hyperliquid | Public | Books | `.env` | N/A | No |
| `PIT_HL_EXCHANGE_URL` | Yes | E | Hyperliquid | Public | Orders | `.env` | N/A | No |
| `PIT_HL_TESTNET_INFO` | Yes if testnet | D | Hyperliquid | Public | Test books | `.env` | N/A | No |
| `PIT_HL_TESTNET_EXCHANGE` | Yes if testnet | D | Hyperliquid | Public | Test orders | `.env` | N/A | No |
| `PIT_AGENTICID_CONFIG_URL` | Yes | A | 0G | `https://agenticid.0g.ai/config` | Live attestor chain | `.env` | N/A | No |
| `PIT_COINGECKO_URL` | No | F | CoinGecko | Public | Spot context | `.env` | N/A | No |
| `PIT_DEFILLAMA_URL` | No | F | DefiLlama | Public | TVL context | `.env` | N/A | No |

**Pinned public values (M18 eth_getCode + official SDK):**

| Item | Mainnet `16661` | Galileo `16602` |
|---|---|---|
| RPC | `https://evmrpc.0g.ai` | `https://evmrpc-testnet.0g.ai` |
| Explorer | `https://chainscan.0g.ai` | `https://chainscan-galileo.0g.ai` |
| PC UI | `https://pc.0g.ai` | `https://pc.testnet.0g.ai` |
| Router API | `https://router-api.0g.ai/v1` | `https://router-api-testnet.integratenetwork.work/v1` |
| InferenceServing | `0x47340d900bdFec2BD393c626E12ea0656F938d84` (code 502) | `0xa79F4c8311FF93C06b8CfB403690cc987c93F91E` (code 502) |
| LedgerManager | `0x2dE54c845Cd948B72D2e32e39586fe89607074E3` (code 502) | `0xE70830508dAc0A97e6c087c75f402f9Be669E406` (code 502) |
| Storage Flow | `0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526` | `0x22E03a6A89B950F1c82ec5e74F8eCa321a105296` (code 295) |
| Storage indexer | `https://indexer-storage-turbo.0g.ai` | `https://indexer-storage-testnet-turbo.0g.ai` |
| ERC-8004 Identity | `0x8004A169FB4a3325136EB29fA0ceB6D2e539a432` | `0x8004A818BFB912233c491871b3d84c89A494BD9e` (code 130) |
| ERC-8004 Reputation | `0x8004BAa17C55a88189AE136b182e5fdA19dE9b63` | `0x8004B663056A597Dffe9eCcC1965A193B7388713` (code 130) |
| Foundation AgenticID | **code 0** | `0x34493302287308f565cf3409daadedf4c8895648` (code 295) |
| TEE verifier | **code 0** | `0x9D48FcCE51b4b39fcB6e4bd0840f75A987cEf980` (code 295) |
| PitDeskID (gate) | `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` (code 7071) | **code 0** — not deployed |
| Faucet | n/a | `https://faucet.0g.ai` (0.1 0G/day) |

Startup: if `PIT_ALLOW_FALLBACKS!=false` → exit. If selected chainId mismatches RPC → exit.

---

## B. USER-SCOPED CREDENTIALS

Created by `pit init` / desktop onboarding. Namespaced by `workspaceId` **and** `PIT_NETWORK`. **Never** in server `.env`, localStorage, Vercel, MCP, or git.

| Name | Required? | Environment | Who owns it? | How obtained? | What it grants | Where stored | Rotation | Never expose? |
|---|---|---|---|---|---|---|---|---|
| Wallet address | Yes | B | User | Connect + SIWE | YOUR WALLET identity | local DB (public) | Reconnect | No |
| HL master address | Yes to trade | B | User | User proves `approveAgent` | Queries + extraAgents | local DB | Re-link | No |
| HL network | Yes | B | User | Same as `PIT_NETWORK` in product | Testnet vs mainnet book | local DB | Toggle | No |
| Optional HL subaccount | No | B | User | HL UI | Isolated capital | local DB | N/A | No |
| Session agent key | Yes to trade | B | User | Generated locally; user signs approve | `order`/`cancel` ~1h; **cannot withdraw** | **OS keychain** | New session | **Yes** |
| Memory encryption key | Yes if memory | B | User | Generated locally | Decrypt YOUR memory | **OS keychain** | Rotate = new namespace | **Yes** |
| Policy JSON + version | Yes to trade | B | User | Cards + pin tx | YOUR POLICY | local DB + chain | New pin | Semi-private |
| Skill weights | After fills | B | User | Calibration | YOUR history only | encrypted memory | N/A | Ciphertext |
| 8004 agentId | After register | B | User | User-signed register | Desk identity | local DB + chain | N/A | No |
| 8004 feedback client key | If feedback | B | User | Separate local EOA | Reporter ≠ owner | OS keychain | Rotate | **Yes** |
| Direct ephemeral `app-sk-` | Each inference | B | User wallet | SDK `getHeader` tokenId 255 | 24h Direct provider auth | Memory / OS keychain | Auto 24h | **Yes** — not Router `sk-` |
| Direct persistent token | Optional | B | User | SDK `createApiKey` tokenId 0–254 | Long-lived Direct auth **per provider** | OS keychain | Revoke on-chain | **Yes** |
| User 0G spend | Yes for Direct | B | User | User-signed `transferFund` to **their** Ledger subaccount | Pays TeeML | On-chain Ledger | Top up | Key never in PIT |

User funds **their own** 0G native for LedgerManager. PIT does **not** spend a global payer for user inference in production.

**Router `sk-` / `mk-` are not user-scoped PIT credentials.** They belong to whoever created them on pc.0g.ai and spend the **Router Payment Layer**, which is a **separate contract** from LedgerManager subaccounts (official comparison.md). They cannot be used to “see” a user’s private book unless PIT sends that book to the Router — which is forbidden.

---

## C. DEVELOPMENT FIXTURES (`PIT/.env` — gitignored)

Used **only** for `_gate` scripts. **Must not** be wired as “the PIT user.”

| Name | Required? | Who owns it? | How obtained? | What it grants | Where stored | Rotation | Never expose? |
|---|---|---|---|---|---|---|---|
| `PIT_OG_PAYER_KEY` | Gate only | Fixture | Local wallet | 0G spend + historical HL master **fixture** `0xBDfCeE82Bd42FEfA58ee850B3709636a8B6b0034` | gitignored `.env` | Rotate fixture | **Yes** |
| `PIT_DEPLOYER_KEY` | Gate / CI deploy | Fixture | Local wallet | Gate deploy + 8004 client `0xaAE3…C2a` | gitignored `.env` | Rotate | **Yes** |
| `PIT_MEMORY_KEY` | Gate S6 | Fixture | Random 32-byte | Gate encrypt only | gitignored `.env` | Per-run ok | **Yes** |
| `PIT_8004_AGENT_ID` | Gate | Fixture | Live register | `3489333` owner = fixture payer | `.env` | Re-register | No |
| `PIT_DESK_ID_CONTRACT` | **Should be pinned** | PIT | Gate deploy | `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` | `.env` | CREATE2 | No |
| `PIT_MASTER_ADDRESS` | Must stay empty in product | — | — | Product reads HL master from **user connect** | — | — | — |
| `PIT_HL_NETWORK` | Gate leftover | Fixture | Currently `Testnet` in file | **Do not** use as product default; UI owns the toggle | `.env` gate only | — | No |
| `PIT_PROVIDER_URL` | Optional override | Gate | Empty in current file | Direct URL if set; else `getService` | `.env` | N/A | No |
| `PIT_MODEL_RESEARCH` / `RISK` / `VISION` | Gate hints | Gate | SKU names | Must still pass TeeML + funded provider | `.env` | Re-fetch catalog | No |

**M18 fixture note:** `PIT_DESK_ID_CONTRACT` is pinned to the live Aristotle address (public). `PIT_PROVIDER_URL` may stay empty (resolve via `getService`). No Router `sk-`/`mk-` belongs in this file.

---

## D. TESTNET ONLY (`PIT_NETWORK=testnet`)

Purpose: **full protocol laboratory** — CI, Chrome E2E, HL test fills, Galileo 8004, Foundation Agentic ID attestor path, experimental ERC-7857 transfer **if** a PIT-signed `iTransferFrom` tx is later proven.

| Capability | Status | Evidence |
|---|---|---|
| Galileo chain | LIVE | `eth_chainId` 16602 |
| Galileo compute contracts | LIVE | serving+ledger code 502 |
| Galileo Direct TeeML **chat** | **Available, PIT-unproven** | `qwen/qwen2.5-omni-7b` teeAck true, provider `0xa48f01287233509FD694a22Bf840225062E67836`, signer `0x83df4B8EbA7c0B3B740019b8c9a77ffF77D508cF` |
| Galileo glm-5.2 | **Absent** | Router catalog N=2; Chrome Private TeeML = **1 image model** (`Qwen-Image-Edit`); no glm-5.2 |
| Galileo Router `sk-` | Separate from mainnet | Official: different UI, keys, balances |
| Storage | Addresses live | Flow + turbo indexer |
| ERC-8004 | Bytecode live | Identity + Reputation |
| Foundation AgenticID + verifier | LIVE on Galileo | config `chain_id:16602`; code 295/295 |
| iTransfer | Official path exists; **PIT tx UNVERIFIED** | Do not mark ✅ until `ownerOf` changes in a PIT tx |
| PitDeskID at gate address | **Not deployed** | code 0 |
| Hyperliquid testnet | LIVE | Chrome `app.hyperliquid-testnet.xyz` wallet connected, ~999 USDC fixture |
| Faucet | LIVE | `https://faucet.0g.ai` 0.1 0G/day |

**Kill:** claiming Galileo has the same 5 Private TeeML chat models as mainnet. Claiming iTransfer LIVE without a PIT transfer tx. Using Router Omni (`TeeTLS` on testnet Router) as if it were Direct TeeML.

---

## E. MAINNET ONLY (`PIT_NETWORK=mainnet`)

Purpose: **production** — real money, real 0G receipts, real Direct glm-5.2, real HL execution, production-safe Agentic ID subset (mint / authorize / revoke; **no** Foundation iTransfer).

| Capability | Status | Evidence |
|---|---|---|
| Aristotle chain | LIVE | `eth_chainId` 16661 |
| Direct TeeML glm-5.2 | LIVE (M17) | provider `0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D`, teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9` |
| Router Private TeeML | 5 SKUs | Chrome `pc.0g.ai/models` + Router GET N=29 |
| Independent 0GM providers | Catalog yes; **fixture unfunded** | Do not claim 3 independent sealed models |
| Storage `--proof` | LIVE (S6) | Official Go client |
| ERC-8004 | LIVE | Identity + Reputation |
| PitDeskID mint/authz | LIVE (M17) | tokenId 1 owner fixture |
| Foundation iTransfer | **IMPOSSIBLE** | attestor+verifier code 0; config still 16602 |
| Hyperliquid mainnet | LIVE | Chrome ETH mark ~2450; M15 dust oid |

**Kill:** showing Transfer as available. Router fallback. Shipping fixture payer as the product user.

---

## F. OPTIONAL

| Name | Required? | Who owns it? | How obtained? | What it grants | Where stored | Never expose? |
|---|---|---|---|---|---|---|
| `PIT_PC_ROUTER_API_KEY` (`sk-`) | **No for PIT** | Whoever created it on pc.0g.ai | Dashboard → API Keys | Router inference only. **Cannot** replace Direct. **Cannot** read a user’s book unless PIT sends it. | If ever used: server secret store, never browser | **Yes** |
| `PIT_PC_MGMT_KEY` (`mk-`) | **No for PIT** | Same | Dashboard → Management Keys | Account admin; not inference | server secret store | **Yes** |
| `PIT_X402_*` | No | Optional | Herald | Public evidence pack | `.env` | No |
| `PIT_0G_SWAP_ROUTER` / Zia / W0G / USDCE / WETH | No | 0G DeFi | RESOURCES | Labeled dust only | `.env` | No |
| Privy App ID | No until wallet lib chosen | PIT | Privy dashboard | Connect UX only — **not** HL keys | `.env` public id | App secret yes |
| `PIT_GAMMA_URL` / Dexscreener | No | Public | URLs | Research context | `.env` | No |

**Do not add** `PIT_PC_ROUTER_API_KEY` to the product book path. Creating one is **not required** for implementation.

---

## Forbidden

- `PIT_OPERATOR_KEY`
- HL session key in `.env`, browser, MCP, Vercel, Render, plaintext SQLite
- User seed / mnemonic / raw HL API secret fields in UI
- Silent mock if indexer/compute unset
- Using fixture `PIT_OG_PAYER_KEY` as a hosted trading operator
- Putting Router `sk-` in the client
- Falling back from Direct to Router (Hanami)
- Treating `PIT_MASTER_ADDRESS` as global user

---

## Hyperliquid public (no secrets)

Implementation agent must read before HL modules:

- https://hyperliquid.gitbook.io/hyperliquid-docs
- Exchange: `.../for-developers/api/exchange-endpoint.md`
- Nonces / API wallets: `.../for-developers/api/nonces-and-api-wallets.md`
- Account modes: `.../trading/account-abstraction-modes.md`
- Testnet faucet: `.../onboarding/testnet-faucet.md` (requires prior mainnet deposit)

Tick: `round(float(f"{px:.5g}"), 6-szDecimals)` for perps. Min notional **$10**. Unified: USDC in **spot** clearinghouse; perp `accountValue=0` is expected.

---

## Databases

| DB | Scope | Isolation |
|---|---|---|
| Desktop/CLI SQLite | Per workspace file | `pit/<network>/<workspaceId>/ledger.db` encrypted at rest |
| Testnet vs mainnet | Separate files | Never share cloids / memory / sessions across networks |
| Hosted server DB | **Not v1 for session keys** | If a read-only web cache exists, it holds **no** session material |

---

## `.env.example` (infra + empty user slots)

```dotenv
# A. GLOBAL — deploy machine. No user keys. No Router sk-.
PIT_NETWORK=mainnet
PIT_ALLOW_FALLBACKS=false

# E. MAINNET (Aristotle 16661)
PIT_CHAIN_RPC_URL=https://evmrpc.0g.ai
PIT_CHAIN_ID=16661
PIT_EXPLORER_BASE=https://chainscan.0g.ai
PIT_ROUTER_URL=https://router-api.0g.ai
PIT_SERVING_CONTRACT=0x47340d900bdFec2BD393c626E12ea0656F938d84
PIT_LEDGER_CONTRACT=0x2dE54c845Cd948B72D2e32e39586fe89607074E3
PIT_STORAGE_INDEXER=https://indexer-storage-turbo.0g.ai
PIT_STORAGE_FLOW=0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526
PIT_8004_IDENTITY=0x8004A169FB4a3325136EB29fA0ceB6D2e539a432
PIT_8004_REPUTATION=0x8004BAa17C55a88189AE136b182e5fdA19dE9b63
PIT_DESK_ID_CONTRACT=0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35
PIT_HL_INFO_URL=https://api.hyperliquid.xyz/info
PIT_HL_EXCHANGE_URL=https://api.hyperliquid.xyz/exchange

# D. TESTNET (Galileo 16602) — used only when PIT_NETWORK=testnet
PIT_TESTNET_RPC_URL=https://evmrpc-testnet.0g.ai
PIT_TESTNET_CHAIN_ID=16602
PIT_TESTNET_EXPLORER=https://chainscan-galileo.0g.ai
PIT_TESTNET_ROUTER_URL=https://router-api-testnet.integratenetwork.work
PIT_TESTNET_SERVING_CONTRACT=0xa79F4c8311FF93C06b8CfB403690cc987c93F91E
PIT_TESTNET_LEDGER_CONTRACT=0xE70830508dAc0A97e6c087c75f402f9Be669E406
PIT_TESTNET_STORAGE_INDEXER=https://indexer-storage-testnet-turbo.0g.ai
PIT_TESTNET_STORAGE_FLOW=0x22E03a6A89B950F1c82ec5e74F8eCa321a105296
PIT_TESTNET_8004_IDENTITY=0x8004A818BFB912233c491871b3d84c89A494BD9e
PIT_TESTNET_8004_REPUTATION=0x8004B663056A597Dffe9eCcC1965A193B7388713
PIT_TESTNET_DESK_ID_CONTRACT=
PIT_HL_TESTNET_INFO=https://api.hyperliquid-testnet.xyz/info
PIT_HL_TESTNET_EXCHANGE=https://api.hyperliquid-testnet.xyz/exchange
PIT_AGENTICID_CONFIG_URL=https://agenticid.0g.ai/config

PIT_STORAGE_CLI=
PIT_DEPLOYER_KEY=
PIT_POLICY_CONTRACT=
PIT_RECEIPTS_CONTRACT=
PIT_FORECASTS_CONTRACT=
PIT_X402_ENABLED=false
PIT_X402_FACILITATOR=https://facilitator.heraldprotocol.xyz
PIT_0G_SWAP_ROUTER=0x807F4E281B7A3B324825C64ca53c69F0b418dE40
PIT_0G_ZIA_ROUTER=0x18cCa38E51c4C339A6BD6e174025f08360FEEf30
PIT_0G_W0G=0x1Cd0690fF9a693f5EF2dD976660a8dAFc81A109c
PIT_0G_USDCE=0x1f3AA82227281cA364bFb3d253B0f1af1Da6473E
PIT_0G_WETH=0x564770837Ef8bbF077cFe54E5f6106538c815B22

# B. PRODUCT: NOT set in server env. Created per workspace.
# PIT_MASTER_ADDRESS=
# PIT_HL_SUBACCOUNT=
# PIT_MEMORY_KEY=
# PIT_8004_AGENT_ID=
# PIT_PC_ROUTER_API_KEY=   # FORBIDDEN on private book. Do not set.
```

---

## API key decision (M18)

**DEPENDS — Direct private-book path does not require a pc.0g.ai Router `sk-`/`mk-` key.**

- **Confidential book / committee / researcher / challenger / risk:** **NO** dashboard API key. Auth is the **user’s wallet** signing Direct `app-sk-` (ephemeral 255). Provider funding is Ledger `transferFund` from **that user**.
- **Router `sk-`:** optional only for unlabeled public catalog/debug. Created at https://pc.0g.ai/dashboard/api-keys (mainnet) or https://pc.testnet.0g.ai/dashboard/api-keys (testnet). Scope = OpenAI-compatible Router inference; trust-mode radios (Standard / Verified / Private) are Router-only. Server-side if ever used. **Never** per-user for the book. **Cannot** access the private book unless PIT sends it. Rotate: revoke in Console.
- **Persistent Direct `createApiKey`:** optional per-user per-provider in OS keychain — **not** the dashboard `sk-` key.

**Do not ask the user to paste an API key for PIT implementation.**
