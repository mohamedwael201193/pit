# PIT — Resources

Only items used for architecture. Unverified claims are not listed as FACT.

Probed **2026-08-26** (M18 revalidation) unless noted.

**Implementation agent:** before each `IMPLEMENTATION_PLAN.md` phase, read the **External resources** row for that phase **and** the matching section here. Do not invent repos. Do not use Galileo addresses on Aristotle. Do not use Aristotle glm-5.2 SKUs on Galileo.

**Env classification:** `ENVIRONMENT.md` A–F. Fixture payer `0xBDfC…0034` is **not** the product owner.

---

## M18 live revalidation (do not skip)

| Check | Result | Evidence |
|---|---|---|
| Aristotle `eth_chainId` | `16661` | RPC `https://evmrpc.0g.ai` |
| Galileo `eth_chainId` | `16602` | RPC `https://evmrpc-testnet.0g.ai` |
| Router mainnet catalog | HTTP 200, **N=29** | `GET https://router-api.0g.ai/v1/models` |
| Router TeeML | **5**: `glm-5.2`, `0gm-1.0-35b-a3b`, `0gm-1.0-35b-a3b-sia`, `whisper-large-v3`, `z-image-turbo` | same |
| Router TeeTLS | 16 | same |
| Router no-verifiability | 8 Claude/GPT family | **forbidden** on book |
| Router **testnet** catalog | HTTP 200, **N=2**: `qwen-image-edit` TeeML, `qwen2.5-omni` **TeeTLS** | `GET https://router-api-testnet.integratenetwork.work/v1/models` |
| Chrome `pc.0g.ai/models` | Private TeeML **5** (logged-in, 5.09 0G Mainnet) | real Chrome |
| Chrome `pc.testnet.0g.ai/models` | Private TeeML **1** = **Qwen-Image-Edit**; Chat **0** under that filter | real Chrome |
| Galileo Direct `getAllServices` | total **6**; TeeML ack **true**: `qwen/qwen2.5-omni-7b` + `qwen/qwen-image-edit-2511` | serving `0xa79F…F91E` |
| Galileo unacked TeeML | `openai/gpt-oss-20b`, `google/gemma-3-27b-it` | **do not use** until teeAck |
| Attestor config | `chain_id: 16602`; AgenticID `0x34493302287308f565cf3409daadedf4c8895648`; verifier `0x9d48fcce51b4b39fcb6e4bd0840f75a987cef980` | Chrome + curl |
| Foundation bytecode | Aristotle **0 / 0**; Galileo **295 / 295** | eth_getCode |
| PitDeskID | Aristotle **7071**; Galileo **0** | eth_getCode |
| Galileo 8004 | Identity+Reputation code **130** | eth_getCode |
| Galileo serving/ledger | both code **502** | SDK `constants.ts` + RPC |
| Chrome Hub | `hub.0g.ai/discover?network=testnet` selector Galileo | real Chrome |
| Chrome HL mainnet | ETH ~2450 | real Chrome |
| Chrome HL testnet | connected `0xBDfC…0034`, ~999 USDC, HYPE mark ~50–60 | real Chrome |
| Chrome PC API Keys | Router `sk-` create UI; trust modes Standard/Verified/Private; **no key created this pass** | real Chrome — **values not recorded** |
| `0g-agentic-id` HEAD | `b8f4845` sdk 0.1.2 | git log |
| `0g-pc-e2ee` HEAD | `0ba9fb0` | git log |
| `0g-compute-ts-sdk` HEAD | `bff2913` v0.9.0 | git log |

**Roles (mainnet):** RESEARCHER / CHALLENGER / RISK = **glm-5.2** Direct (only funded TeeML chat provider in the fixture). 0GM / 0GM-SIA are catalog-Private but **unfunded** for this wallet — not independent Direct sealed providers.

**Roles (testnet):** only ack’d TeeML **chat** is `qwen/qwen2.5-omni-7b`. Image model is not a committee. PIT VerifyE2EE on Galileo = **UNVERIFIED**.

---

## 0G chain and docs

| Source | URL / address | Status | Verified | Why it matters |
|---|---|---|---|---|
| Aristotle mainnet | chainId `16661` `0x4115`, RPC `https://evmrpc.0g.ai`, explorer `https://chainscan.0g.ai` | FACT | eth_chainId, block `0x28ac737` | Policy, forecasts, receipts, 8004, 7857 |
| Docs | https://docs.0g.ai/ | — | M18 local 0g-doc | Router holds prompts — confidential path unsafe |
| ERC-8004 docs | https://docs.0g.ai/developer-hub/building-on-0g/agentic-id/erc8004 | FACT (local 0g-doc copy) | Identity + Reputation addresses | Official registries |
| Builder hub | https://build.0g.ai/ | — | — | SDK index |
| PC mainnet | https://pc.0g.ai · /models · /sdk · /dashboard/api-keys | FACT Chrome M18 | Advanced = Direct; API Keys = Router `sk-` only | Do not mix |
| PC testnet | https://pc.testnet.0g.ai | FACT Chrome M18 | Separate keys/balances (official overview.md) | Do not reuse mainnet `sk-` |
| Router vs Direct | https://docs.0g.ai/developer-hub/building-on-0g/compute-network/router/comparison.md | FACT docs | Separate Payment Layer vs Ledger subaccounts | Direct for book |
| Router auth | …/router/authentication.md | FACT docs | `sk-` inference; `mk-` admin; create at Dashboard → API Keys | **Not required for PIT Direct** |
| Router privacy | …/router/privacy.md | FACT docs | Prompts in memory to route/bill; privacy mode still trusts Router | Forbidden for book |
| Direct inference | …/compute-network/inference.md | FACT docs | Wallet-signed provider calls | PIT path |
| Direct token | `0g-serving-broker/doc/api-token.md` + SDK `getHeader` | FACT source | Ephemeral tokenId 255 (`app-sk-`); persistent 0–254 `createApiKey` | Not dashboard `sk-` |
| InferenceServing | `0x47340d900bdFec2BD393c626E12ea0656F938d84` | FACT bytecode | eth_getCode (prior) | getService teeSigner |
| LedgerManager | `0x2dE54c845Cd948B72D2e32e39586fe89607074E3` | FACT bytecode | | Direct compute credit |
| PaymentLayer | `0xA3b15Bd2aD18BFB6b5f92D8AA9F444Dd59d1cE32` | FACT bytecode | **no user withdraw** | Do not escrow trades |
| Storage Flow | `0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526` | FACT bytecode | | Uploads |
| Storage indexer | `https://indexer-storage-turbo.0g.ai` | — | | Mainnet |
| Galileo serving | `0xa79F4c8311FF93C06b8CfB403690cc987c93F91E` | FACT M18 bytecode 502 | SDK testnet.inference | Direct on 16602 |
| Galileo ledger | `0xE70830508dAc0A97e6c087c75f402f9Be669E406` | FACT M18 bytecode 502 | SDK testnet.ledger | |
| Galileo Flow | `0x22E03a6A89B950F1c82ec5e74F8eCa321a105296` | FACT M18 bytecode 295 | testnet-overview | |
| Galileo indexer | `https://indexer-storage-testnet-turbo.0g.ai` | SOURCE docs/clones | | |
| ERC-8004 Identity Galileo | `0x8004A818BFB912233c491871b3d84c89A494BD9e` | FACT M18 code 130 | | |
| ERC-8004 Reputation Galileo | `0x8004B663056A597Dffe9eCcC1965A193B7388713` | FACT M18 code 130 | | |
| ERC-8004 Identity | `0x8004A169FB4a3325136EB29fA0ceB6D2e539a432` | FACT bytecode 2026-08-26 | eth_getCode len 262 (proxy) | Register desk agent |
| ERC-8004 Reputation | `0x8004BAa17C55a88189AE136b182e5fdA19dE9b63` | FACT bytecode 2026-08-26 | eth_getCode len 262 (proxy) | Calibration feedback |
| Router | `https://router-api.0g.ai` | FACT | 29 models | Catalog only |
| DASigners `0x000…1000` | epochNumber `0x61ba5a0d`; params | FACT 2026-08-26 re-gate | `epochNumber()=0`; `params()` tokensPerVote=30, maxVotesPerSigner=102400, maxQuorums=10, epochBlocks=460800, encodedSlices=3072; `quorumCount(0)=0`. Prior “eth_call OOG” is stale. | Callable ≠ usable blob DA |
| DAEntrance | docs: Galileo `0xE75A073dA5bb7b0eC622170Fd268f35E675a957B` | FACT | Mainnet table **has no DAEntrance**; Aristotle code at that address **length 0** | Do not claim DA |
| Agentic ID attestor | https://agenticid.0g.ai/config | FACT M17 Chrome + curl | `chain_id: 16602`; RPC `https://evmrpc-testnet.0g.ai`; AgenticID `0x34493302287308f565cf3409daadedf4c8895648`; TEE verifier `0x9d48fcce51b4b39fcb6e4bd0840f75a987cef980`; frameworks openclaw / hermes / prime-agent | iTransfer not on Aristotle |
| Agentic ID bytecode | same addresses | FACT M17 `eth_getCode` | Aristotle len **0** / **0**; Galileo len **295** / **295** (proxy). PitDeskID Aristotle **7071**, Galileo **0** | Do not call Foundation iTransfer LIVE |
| Agentic ID SDK note | `0g-agentic-id` `sdk/typescript/src/constants.ts` | FACT source | `CHAIN_ID = 16602`; `ZERO_G_MAINNET` id 16661 with comment “AgenticID is testnet-only today — no mainnet contract deployment” | Mint/own/authorizeUsage on 16661 with production IDs; iTransfer reverts honestly |
| Foundation §2.4 | `contracts/DEPLOYMENT.md` | FACT source | Seal-bound agents: ERC-721 transferFrom **re-enabled**, `iTransferFrom` **reverts** `AgenticIDSealedAgentUseTransfer` | Galileo sealed runtime ≠ Aristotle iTransfer |

---

## Live model catalog (Router + Chrome 2026-08-26)

`GET https://router-api.0g.ai/v1/models` HTTP 200. Chrome `pc.0g.ai/models` filter **Private (TeeML) 5** (logged-in Hub wallet `5.09 0G Mainnet`): `0GM-1.0-35B-A3B` (262k, Private), `GLM-5.2` (1,048,576, Private), `0GM-1.0-35B-A3B-SIA` (32,768, Private), plus Image `z-image-turbo` and Speech `whisper-large-v3` under the same filter.

On-chain `getAllServices` may label additional SKUs `TeeML` (M17 saw glm-5.1, qwen, deepseek). **Do not treat that as Router Private-5.** Sealed book requires Router `verifiability=TeeML` **and** a funded Direct provider with matching `teeSigner`. Fixture committee used glm-5.2 only.

Direct TeeML glm-5.2 (M17 committee): provider `0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D`, url `https://compute-network-19.integratenetwork.work`, teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`.

`getAllServices` page size max **10**. `transferFund` type-2 txs need `gasPrice`. Unsettled fees + 1.0 0G reserve: locked 1.800 < min 1.844 → HTTP 400 until top-up.

New SKUs vs earlier snapshot (still forbidden on book path if no verifiability): `claude-fable-5`, `claude-opus-4-8`, `deepseek-v4-*`, `seedance-2.5`, `qwen3.6`/`3.7`/`3.8`, `gpt-5.5`, `gpt-5.6-luna`/`sol`/`terra`. Claude/GPT family remain **forbidden** on the sealed book path.

| SKU | verifiability | PIT role |
|---|---|---|
| `glm-5.2` | **TeeML** | M17 sealed researcher **and** challenger **and** risk (same provider; envelope independence) |
| `0gm-1.0-35b-a3b` | **TeeML** | Optional vision risk |
| `0gm-1.0-35b-a3b-sia` | **TeeML** | Risk scorer on structured thesis (**not** order compiler) |
| `whisper-large-v3` | TeeML | Not v1 desk |
| `z-image-turbo` | TeeML | Not v1 desk |
| `glm-5.3` | TeeTLS | Public unlabeled notes only; default off |
| `glm-5` / `glm-5.1` / Qwen / Kimi / DeepSeek / MiniMax / hy3 | TeeTLS | Same restriction |
| Claude / GPT family | **no verifiability field** | **Forbidden** |

---

## 0G SDKs / source (local)

| Item | Path / version | HEAD | Verified | Why |
|---|---|---|---|---|
| `0g-pc-e2ee` | `D:\route\0g\research\0g-pc-e2ee` | `0ba9fb0b` 2026-08-20 | protocol tests (prior) | Seal + VerifyE2EE |
| Binding | `protocol/proof/proof.go` | same | KAT | Caller-byte bind |
| Direct `-onchain` reject | `client/cmd/internal/proxycli/proxycli.go` | same | source | Wire getService |
| Storage TS SDK | `@0gfoundation/0g-storage-ts-sdk@1.2.11` | npm | TODO ignores proof | **Do not use for proofs** |
| Storage Go client | `D:\route\0g\research\0g-storage-client` HEAD `3d953af` | `--encryption-key` AES-256-CTR v1 **requires `0x` prefix**; `--encrypt` ECIES v2; `--proof` on download | FACT source + live S6 | **Official encryption** |
| Gate Storage CLI | `D:\route\0g\PIT\_gate\bin\0g-storage-client.exe` | encrypt+upload+download `--proof` SHA256 match; wrong root fail; wrong key fail | FACT 2026-08-26 | Root `0x2f0419508df7dc546f9aa84d214723e72cc2547fa1f69d828c9187510cfbda5d` |
| Compute TS SDK | `@0gfoundation/0g-compute-ts-sdk@0.9.0` HEAD `bff2913` | `processResponse` is EIP-191 on provider text, **not** caller-byte E2EE bind; `transferFund` serviceName **`inference-v1.0`** | FACT source + live revert on `"inference"` | Ledger only |
| `0g-agentic-id` | `D:\route\0g\research\0g-agentic-id` HEAD `b8f4845` | IERC7857 / Authorize / Cloneable; attestor `chain_id: 16602` | interface IDs `0x2afbede9` / `0xdf597d99` / `0x74f8628b` | Production 7857 |
| Agentic ID SDK | `@0gfoundation/0g-agenticid-sdk@0.1.2` | Galileo-default | not the mint path on 16661 | |
| Official mcp-0g | stale Galileo | not on npm | unused | |

---

## x402 / Pay

| Item | URL | Status | Verified | Why |
|---|---|---|---|---|
| Herald facilitator | https://facilitator.heraldprotocol.xyz/supported | FACT 200 | kinds include `eip155:16661` exact + upto; facilitator `0x686Ca1f3BAf7F7Df3334f2f1A65AE314ee9CDb29` | Optional paid **public** evidence pack |
| Coinbase x402.org | https://x402.org/facilitator/supported | FACT 200 | **no** `eip155:16661` (84532 etc.) | Cannot use as 0G rail |
| 0G PaymentLayer | see above | no withdraw | unused for trades | |

---

## 0G DeFi (live 2026-08-26 evening re-probe)

| Protocol | URL / address | TVL / reserves | Role |
|---|---|---|---|
| Chain 0G | llama `zero-gravity` | **$2,871,774** FACT earlier day; venue_probe Gimo **$3.02M** + TradeGPT **$1.91M** (LST + DEX, not JAINE) | Context |
| TradeGPT / Zia DEX | router `0x18cCa38E51c4C339A6BD6e174025f08360FEEf30` factory `0x6F3945Ab27296D1D66D8EEb042ff1B4fb2E0CE70` | llama **$1,907,740** FACT; USDC.e/WETH 0.3% pool `0x22b46C…E942eC` **334,519 USDC.e + 132.48 WETH** FACT block `42656196` | Competitor NL DEX. Same contracts as labeled PIT dust path. Do not become 4lpha |
| Oku Uniswap V3 on 0G | SwapRouter02 `0x807F4E281B7A3B324825C64ca53c69F0b418dE40` factory `0xcb2436774C3e191c85056d248EF4260ce5f27A9D` | W0G/USDC.e 1% `0xCE7737…C71` **502,586 W0G + 43,214 USDC.e** FACT | **PIT labeled dust adapter** (user-signed). Oku **frontend** defaulted to Ethereum; `has0g=false` in swap UI FACT — do not brand “Oku-first” |
| JAINE | https://jaine.app/ · Hub Swap “Powered By JAINE” | llama **$2,914** FACT | Official Hub UI. **Not** PIT’s adapter. Page addr `0x7bbc…1404` is **st0G**, not a router FACT |
| Bond DEX | llama slug `bond` | **$184,946** FACT | Other 0G DEX. **UNVERIFIED** swap API / agent exec. Not v1 |
| Morpho Blue on 0G | llama `morpho-blue` chain 0G | **$32.20** TVL, **$2.32** borrowed FACT | **Not a venue** |
| Gimo LST | llama `gimo` | **$3,019,596** FACT | Staking, not trading |
| Llama uniswap-v3 0G slice | api.llama.fi | **timed out** UNVERIFIED as a Llama line item | Use pool `balanceOf` instead |
| Khalani / TokenFlight | Hub `/khalani/transfer` | UI FACT; sample history **5 0G → 5 0G** same-chain | Intent/cross-chain **funding rail**. Not a CLOB |
| Hub Bridge | Ethereum ↔ 0G W0G | UI “Powered By Chainlink” lead | Onramp. Not a DEX |
| W0G | `0x1Cd0690fF9a693f5EF2dD976660a8dAFc81A109c` | symbol W0G FACT | Wrap native for univ3 |
| USDC.e | `0x1f3AA82227281cA364bFb3d253B0f1af1Da6473E` | — | Dust quote token |
| WETH on 0G | `0x564770837Ef8bbF077cFe54E5f6106538c815B22` | Zia pool 132.48 WETH FACT | Spot ETH on 0G. Not Jordan’s book |

4lpha `lib/contracts/curated-routes.ts` already allowlists these exact Zia/Oku routers and pools. Copying that adapter-as-product **is** competing in their slot.

---

## Hyperliquid

**Must-read before Phase 3/4/9:** https://hyperliquid.gitbook.io/hyperliquid-docs

| Item | URL | Status | Why |
|---|---|---|---|
| Docs home | https://hyperliquid.gitbook.io/hyperliquid-docs | — | Onboarding, tick, accounts |
| Info | `https://api.hyperliquid.xyz/info` | FACT 200, **232 perps**; ETH `szDecimals=4`, mark **$2462.1**, day ntl **~$1.52B**, OI **~671k ETH** 2026-08-26 evening | Product default book |
| Exchange | `https://api.hyperliquid.xyz/exchange` | documented; error **“Order must have minimum value of $10.”** FACT gitbook | Orders |
| Exchange endpoint | https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/exchange-endpoint.md | FACT docs | `order`, `cancel`, `approveAgent` |
| Nonces / API wallets | https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/nonces-and-api-wallets.md | FACT docs | Session agents, `validUntil` |
| Account modes | https://hyperliquid.gitbook.io/hyperliquid-docs/trading/account-abstraction-modes.md | FACT docs | Unified USDC in **spot**; perps `accountValue=0` expected |
| Testnet | `https://api.hyperliquid-testnet.xyz/info` · `/exchange` | 210 assets FACT | Rehearsal |
| Testnet faucet | https://hyperliquid.gitbook.io/hyperliquid-docs/onboarding/testnet-faucet.md | **1,000 mock USDC only if that address has deposited on mainnet** FACT docs | Second-user onboarding must warn |
| extraAgents | POST `/info` `{type:extraAgents,user}` | documented | Judge-visible |
| Tick / sigfigs | gitbook tick size + 5 significant figures | FACT M15 | `1597.31` rejected; `1597.2` accepted |
| `clientOrderId` | exchange docs | documented | Durable exactly-once key |
| HL signing classes (Python SDK) | `hyperliquid.utils.signing` | FACT source | `withdraw3`/`usdSend`/`spotSend` = **user-signed**; `updateLeverage`/`vaultTransfer`/`subAccountTransfer`/`order`/`cancel` = **L1**. PIT refuses withdraw/send **and** `updateLeverage`/`vaultTransfer`/`subAccountTransfer`. Protocol: stolen agent **can** `updateLeverage` — disclose. Never fund session EOA. |
| Fixture HL (DEV ONLY) | extraAgents + spot | FACT 2026-08-26 M15 | Testnet spot **999 USDC**. Mainnet spot **4.8 USDC**. Agents `PIT-M15` `0x71e9…272D` / `PIT-M15m` `0x0b68…9F02`. **Not** a global operator. |

---

## Other market APIs

| API | URL | Status | Why |
|---|---|---|---|
| CoinGecko ping | https://api.coingecko.com/api/v3/ping | 200 `gecko_says` | Spot ref |
| Polymarket Gamma | https://gamma-api.polymarket.com/markets?limit=1&closed=false | 200 prior | Research only |
| DexScreener | https://api.dexscreener.com/latest/dex/search?q=ETH%20USDC | 200 prior | Research only |

---

## NeoSoul research (cloned 2026-08-26)

Do **not** invent extra repos. Org https://github.com/NeoSoul-AI (10 public).

| Repo | URL | HEAD cloned | What was proven | What was not copied |
|---|---|---|---|---|
| rubric-prediction-skill | https://github.com/NeoSoul-AI/rubric-prediction-skill | `c3c3985` 2026-04-28 | Engine v2.3.0; LLM proposes, script computes; fail-closed `insufficient_spec`; 82★ | Python as prod runtime; EvoEvo submit scripts |
| neotrade-skills | https://github.com/NeoSoul-AI/neotrade-skills | `445db07` 2026-08-14 | Skill = SKILL.md; journal/debrief; PM playbooks; Binance kline skills | PM essays; Binance execution; “read and apply” as the only engine |
| neotrade-wallet-sdk | https://github.com/NeoSoul-AI/neotrade-wallet-sdk | `4d92cfd` 2026-08-25 | Signing gateway: no raw-hash, fail-closed, exactly-once, **size not in gateway** (2026-08-16 sell-cap bug); x402 pkg takes raw key | BIP-32 HL trading key; `@neotrade/x402` raw-key transfers as PIT’s x402 |
| neotrade-release | https://github.com/NeoSoul-AI/neotrade-release | `f2a18a6` (local shallow; remote newer) | Artifacts only; self-custodial claim; not notarized | Closed-source runtime |
| evoevo-contracts | https://github.com/NeoSoul-AI/evoevo-contracts | `9308aa3` 2026-08-24/25 | 8004 as identity layer; binding + reasoning commitments; UUPS committee oracle | Committee oracle; UUPS; PM settlement |
| evoevo-agent-kit | https://github.com/NeoSoul-AI/evoevo-agent-kit | `eb599a9` 2026-06-24 | Identity key `(chain, registry, agentId)`; 8004 feedback after settle | Hosted opinion polling |
| trusted-evidence-engine | https://github.com/NeoSoul-AI/trusted-evidence-engine | `ff60e87` 2026-06-16 | Evidence pack schema; CoinGecko/Polymarket providers; not a signer | Brave-key product; TEE server as PIT |
| 0g-builders | https://github.com/NeoSoul-AI/0g-builders | README | Lists EvoEvo contracts including Identity `0x8004Ae533a0301CbD7508373b663756D26DfB028` | **Do not use this registry** — not the official 0G 8004 Identity |
| erc-8004-contracts | fork of erc-8004/erc-8004-contracts | not re-cloned this pass | Standard | — |
| new-api | https://github.com/NeoSoul-AI/new-api | listed API 2026-08-24 | LLM hub | Unused |
| Sites | https://neosoul.ai/ · https://neotrade.pro/ · https://evoevo.ai/ · https://x.com/NeoSoulAI | Jina/local fetches | Desktop product + EvoEvo network | Closed core |

Twitter CLI was not installed (`agent-reach doctor` twitter warn). Web/GitHub used instead.

---

## Competitors (from `03_COMPETITORS_WINNERS_AND_CODE_REVERSE_ENGINEERING.md`)

Use **these** URLs only:

| Project | GitHub / product | Local clone HEAD | Notes |
|---|---|---|---|
| Axiom | https://github.com/symulacr/axiom-protocol | `6f0c496` 2026-07-22 | 15 pts; MockUSDC; InMemoryStorage; software `AXIOM_TEE_SIGNER_PK` |
| 4lpha | https://github.com/kann420/4lpha-0G · https://0g.4lpha.tech/ | `a271c79` 2026-08-23 | 14 pts; no compute SDK; prod ABI **omits** iTransfer/iClone |
| Knole | https://github.com/Pratiikpy/Knole | `d0855a5` 2026-08-24 | 13 pts; journal; custom 7857 IDs `0x4b396f04`; owner-move iTransfer without Foundation attestor |
| Hanami | https://github.com/ajanaku1/hanami | local clone M18 | Direct optional; **falls back to Router** (`og-compute-direct.ts`); default `OG_DIRECT_ENABLED` false; Router `sk-` in `og-compute.ts` | PIT kill to copy |
| Talos | https://github.com/enliven17/talos | `34a7bf6` 2026-07-18 | Groq fallback; `TalosTestnetVerifier.verifyTransferValidity` **always `valid: true`** |
| OptimIEra | https://github.com/shrixvxv-prog/OptimIEra | `0c8d6cb` 2026-07-19 | Agentic ID **PLANNED**, not minted |
| diversifi | https://github.com/thisyearnofear/diversify | `865d438` 2026-08-26 | mocks; evidence-buy; storage-as-DA (`persistence-service.ts`); Galileo indexer default |
| cognivern | https://github.com/thisyearnofear/cognivern | `a13a81c` 2026-08-26 | spend-control proofs; not a desk |
| NEXUS | https://github.com/harsh11067/NEXUS | `6cd254c` 2026-06-22 | **Galileo only**. `CompositeReceiptMinter.sol` stores mapping/events — **no DA call**. iTransfer = trusted ECDSA. Docs overclaim “0G DA” |
| Lumen | https://github.com/OoJae/lumen | `47ee2a9` 2026-08-25 | Honest `OracleNotLive()` on iTransfer/iClone. Gateway/Direct plaintext, not HPKE |
| TRAIDE | https://github.com/NoBanks/traide-keeper-0g | `1a28411` 2026-08-14 | Galileo receipts + Storage merkle; **not live trading** |
| Maneki | https://github.com/maraitona/ManekiAI-Terminal | `c3847c2` 2026-08-02 | Public HL charts only; **no trading/keys** |
| Aevum | **no public repo** | — | Do not invent |
| Orchestra | https://github.com/youssef-jeddi/Orchestra | `research/clones/Orchestra` | ETHGlobal Cannes: NL swap + Groq/0G switch + **AUTO_EXECUTE** small txs + Ledger for large. **Refuse** Groq, auto-exec, Uniswap-as-book. Watcher idea is weaker than PIT Watch (still authorize). |
| Judge file | `D:\route\0g\03_COMPETITORS_WINNERS_AND_CODE_REVERSE_ENGINEERING.md` | — | 2530 lines, v8 |

---

## Prior local research (kept)

`research/VEXPP_EXECUTION_2026-08-24.md` — HL agent cannot withdraw **FACT**; daily loss **not** protocol-enforced.
`research/DEMAND_MAP_2026-08-24.md` — Area B terminals.
`research/MARKET_INTEL_CEILING_2026-08-24.md` — intel-only ceiling.
`VANTA/` — **do not clone** (iNFT + SIA compiler + Flashbots).
`PIT/` prior 84/100 terminal — superseded by Private Alpha OS (this revision).
Gate evidence (gitignored): `PIT/_gate/evidence/`. Do not commit `.env`.

---

## Gate live proofs (2026-08-26 evening)

| Proof | Explorer / id | Result |
|---|---|---|
| Ledger deposit 3.2 0G | `https://chainscan.0g.ai/tx/0xf98e5c4af1674aa1de1564716f462d7203eddcc5710b3cb250c7dbf10caf1fdb` | SUCCESS |
| transferFund 1.2 0G → glm-5.2 TeeML (`inference-v1.0`) | `https://chainscan.0g.ai/tx/0xdad09e14a428d6c698db344972232edeedb3876dd0ad235319c360ebdb43899d` | SUCCESS (`"inference"` reverts) |
| S1 VerifyE2EE | scheme `zg-sig-v1/e2ee-ct`; `ZG-Res-Key` `fec6251d-4708-451a-8130-a9551680ded2`; teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9` | GREEN — recovered signer matched on-chain |
| Fund deployer 0.2 0G | payer `0xBDfC…0034` → deployer `0xaAE3…C2a` native transfer | SUCCESS (hash in gate run; do not treat truncated prefixes as explorer IDs) |
| ERC-8004 register | `https://chainscan.0g.ai/tx/0xa2f67529745a662163b84fe10f855a3aa25596f9bc4d4c604d2abefbc3f3ff7d` | agentId **3489333** owner=payer URI `ipfs://pit-desk-v1` |
| ERC-8004 owner giveFeedback | — | revert `Self-feedback not allowed` |
| ERC-8004 client giveFeedback | `https://chainscan.0g.ai/tx/0xd4892d6345e748080d4a377891e8f555776bb52c9762713dc48c6f815fdbfa25` | SUCCESS from deployer |
| Storage encrypt/`--proof` | root `0x2f0419508df7dc546f9aa84d214723e72cc2547fa1f69d828c9187510cfbda5d` | GREEN; wrong root fail; wrong key fail |
| PitDeskID deploy | `https://chainscan.0g.ai/tx/0x2d5b688bf09bb72cb44b092da0c27cbe87a623141872e62b84cb95ecf7e90c24` | SUCCESS; contract `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35`; block 42705540; Chrome Aristotle Success, unverified |
| PitDeskID mint tokenId 1 | `https://chainscan.0g.ai/tx/0x9494e3faec6d950942d1bfec53c4a13e6f28378da8e01ecd24b3e3de62e5c7d0` | SUCCESS; owner `0xBDfC…0034`; URI `0g://pit-desk-m17-cipher`; Chrome “Mint 1 of PITDESK” |
| PitDeskID authorizeUsage | `https://chainscan.0g.ai/tx/0xdf13b262a875602386708f4c4a121e8ab039e28849d8148e8b991d26fff080d4` | SUCCESS; deployer authorized then revoked `0xfd2e38b6…440af9` |
| PitDeskID unauthorized authorize | `0x74c8ead0747f2bb7e9f0b6f6b84295dcc9252e1c80e8b5e78fa54aa50fa8dc5e` | REVERT (expected) |
| PitDeskID iTransfer/iClone | selector `0x0c9c36b4` | REVERT `AttestorNotOnAristotle` (expected) |
| Ledger top-up glm-5.2 (M17 retry) | `https://chainscan.0g.ai/tx/0x56eb0d630f786343101a980df1ea0b676986cd9148401159925a0d91893c3a03` | SUCCESS 0.6 0G |
| Ledger top-up glm-5.2 (M17 finish) | `https://chainscan.0g.ai/tx/0xf3555b9a85b67997759614edadc126daad709a39bc1a4dbd16ec88099b218e09` | SUCCESS 0.5 0G |
| M17 committee S1 researcher | zg-res-key `7a705ed4-2341-41ef-acca-713afb34259e` | VerifyE2EE OK; glm-5.2; teeSigner `0xA46E…46B9` |
| M17 committee S1 challenger | `e9e81db2-bcdb-4e0f-aea5-45c1091988f1` | VerifyE2EE OK; thesis_killed=true |
| M17 committee S1 risk | `200d4933-278f-4c46-ada2-6c7d1f682ffa` | VerifyE2EE OK; risk_class medium |
| M17 S2 retry (3 roles) | `101562ca…` / `aeca6159…` / `38810a5f…` | all VerifyE2EE OK; engine `no_side` |
| M17 S3 retry (3 roles) | `fc0f4743…` / `ec921a87…` / `8b7110d9…` | all VerifyE2EE OK; engine `below_min_notional` |
| M17 tamper ciphertext | zg-res-key `d025880b-5258-41ac-ac1a-42c42dd2566c` | VerifyE2EE FAIL (content-binding mismatch) — PASS |
| M17 wrong signer | `c01ad1df-5c8b-4715-b1a9-47945e3f16d3` | VerifyE2EE FAIL (recovered `0xA46E…` ≠ `0x000…001`) — PASS |

Payer `0xBDfCeE82Bd42FEfA58ee850B3709636a8B6b0034`. Deployer `0xaAE3EAC0d6665832fe0E5036d61CE2DBC6ECAC2a`. **Do not print private keys.**

---

## Demand pricing

TradingView Premium **$59.95/mo** annual · Nansen Pro **$49–$69** · Discord potion-class **$69–$99 REPORTED**.

---

## KEPT (visual benchmark — not business logic)

| Item | Path / URL | What to copy | What not to copy |
|---|---|---|---|
| Local source | `D:\route\Flare\kept` | Tokens, motion, onboarding pacing, named states, SVG postcard language | Savings circles, XRPL, FXRP, Xaman, rail B |
| Live | https://kept-ruby-five.vercel.app/ | Landing coral hero, split WalletGate, Get started cards, AppShell | Circle join/create product |
| Guide CSS | `apps/web/src/styles/guide.css` | `#D82F2F` coral, `#F0E7D4` cream, `#1A1A1A` ink, Host Grotesk, grain, mega type, pills | — |
| App CSS | `apps/web/src/styles/app.css` | Dark terminal surfaces, Geist, radius 20/14/8, ease-out-quint | Jade accent — PIT uses coral as the one accent |
| Hero | `apps/web/src/landing/Hero.tsx` | Sticky pin, sunburst SVG, postcard | Turn seats |
| WalletGate | `apps/web/src/components/app/WalletGate.tsx` | Split dark/coral, stacked wallet CTAs, no seed field | XRPL pairing |
| StartFlow | `apps/web/src/routes/app/StartFlow.tsx` | Two-card choose path | Asset FXRP/USDT0 |
| namedStates | `apps/web/src/lib/namedStates.ts` | Failure → sentence + next action | KEPT-specific ids |
| Diagrams | `apps/web/src/components/diagrams/keptGuide.tsx` | Cream/coral postcard SVG | Eight-seat circle |
| App router | `apps/web/src/App.tsx` | `/` marketing vs `/app` shell | Circle routes |
| Motion | `motion/react` + `prefers-reduced-motion` | Real state only | Fake 10s spinners |

PIT ring nodes (custom): Market → Seal → Research → Challenge → Risk → Policy → Authorize → Execute → Prove → Learn.

---

## Clone file pointers (implementation consult)

Paths under `D:\route\0g\research\clones\` unless noted. Read README + the listed files; do not copy bytecode.

| Project | Must-read files | Trap |
|---|---|---|
| Axiom | storage mock, TEE signer env, MockUSDC | Software TEE PK; InMemoryStorage |
| 4lpha | `usedActionHashes`, `curated-routes.ts`, PolicyVault | Sandbox liquidity; omit iTransfer in prod ABI |
| Knole | E2EE / journal path | Custom 7857 IDs; owner-move iTransfer |
| Hanami | verify UI; Direct `/v1/proxy` | Router fallback; Direct off by default |
| Talos | Groq fallback; `verifyTransferValidity` | Always `valid: true` |
| OptimIEra | README claims | Agentic ID planned not minted |
| diversifi | `persistence-service.ts`; env gates | Storage-as-DA; Galileo indexer |
| cognivern | spend proofs | Not a desk |
| NEXUS | `CompositeReceiptMinter.sol` | Galileo; mapping ≠ DA |
| Lumen | iTransfer `OracleNotLive()` | Gateway plaintext |
| TRAIDE | keeper receipts / merkle | Galileo; not live trading |
| Maneki | HL charts | No trading |
| NeoSoul | `neotrade-wallet-sdk` gateway; `evoevo-contracts` 8004 bind; `0g-builders` **wrong** Identity `0x8004Ae53…` | Closed core; mnemonic HL key; do not use their 8004 registry |

---

## 0G docs / hub (phase index)

| Source | URL | Use |
|---|---|---|
| Docs | https://docs.0g.ai/ | Compute/Storage/Chain. Router holds prompts — unsafe for book |
| ERC-8004 | https://docs.0g.ai/developer-hub/building-on-0g/agentic-id/erc8004 | Official registries only |
| Builder hub | https://build.0g.ai/ | SDK index |
| Compute | Direct TeeML + vendored `0g-pc-e2ee` | Not TS `processResponse` as E2EE |
| Storage | Go client `--proof` | Not TS SDK 1.2.11 proofs |
| Agentic ID | https://agenticid.0g.ai/config | FACT `chain_id: 16602` |
| Hub | https://hub.0g.ai/ | Onramp/JAINE UI — not PIT venue |
| Explorer | https://chainscan.0g.ai | Evidence |

Vendored: `D:\route\0g\research\0g-pc-e2ee` commit `0ba9fb0b`; `D:\route\0g\research\0g-storage-client` `3d953af`; `D:\route\0g\research\0g-agentic-id` `b8f4845`.

---

## Gate HL execution (M15) — DEV FIXTURE only

| Proof | Id | Result |
|---|---|---|
| Testnet order+cancel | oid **58551048563** limit 1597.2 sz 0.0075 | canceled; USDC 999→999 |
| Mainnet order+cancel | oid **526980359929** limit 1601.0 sz 0.0063 notional $10.0863 | canceled; USDC 4.8→4.8 |
| Allowlist | `_gate/hlcore` `go test` | PASS; Python local_all_denied |
| Protocol hole | agent `updateLeverage` | succeeds if PIT bypassed |

---

## Implementation-agent reading order

1. `FINAL_FLAGSHIP_PRODUCT.md` (thesis, multi-user, Watch, KEPT, honesty)
2. `ENVIRONMENT.md` (A–F: global / user / fixtures / testnet / mainnet / optional)
3. This file (pinned addresses)
4. `IMPLEMENTATION_PLAN.md` **current phase only** + its resource row
5. 03 file sections for the competitor named in that phase
6. KEPT files listed above before any UI
7. HL gitbook pages listed above before any HL module
8. Vendored 0G repos before Compute/Storage/7857

