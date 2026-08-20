# PIT — Research history (append-only)

---

## M1 — Judge + competitor file

- **Source:** `D:\route\0g\03_COMPETITORS_WINNERS_AND_CODE_REVERSE_ENGINEERING.md` entire file (v8).
- **Discovery:** Judge model = real value + real 0G + mainnet + verifiable + tested. Penalizes mocks, env gates, deploy-only, fake 7857, predates-window. Trading/LP agents with policy vaults are **high competition**. 4lpha 14 pts; Maneki is HL + 0G sticker; TRAIDE receipts are Galileo.
- **Candidate:** none yet.
- **Kill:** “be 4lpha but nicer.”
- **Next:** read clones + VANTA so we do not copy them.

---

## M2 — VANTA is the anti-pattern for this brief

- **Source:** `VANTA/FINAL_GLOBAL_WINNING_PRODUCT.md`.
- **Discovery:** VANTA = ERC-7857 iNFT trader + SIA as rule compiler + TEE-born keys + Flashbots + sellable agent. Stack does not support Flashbots on 16661; 7857 production IDs Galileo-only.
- **Kill:** VANTA and any rename (ownable overnight trader).
- **Next:** venue truth.

---

## M3 — HL scoped credential is real; 0G CLOB is not the book

- **Source:** `research/VEXPP_EXECUTION_2026-08-24.md`; live `api.hyperliquid.xyz` 2026-08-26; DeFiLlama.
- **Discovery:** `approveAgent` cannot withdraw **FACT**. Daily loss **not** protocol-enforced. Agent can still trade the subaccount to zero. HL ETH ~$1.78B day notional, 676k ETH OI. TradeGPT/Zia **$1.91M TVL**. JAINE **$2,907**. Chain 0G TVL **$2.88M** (mostly LST + TradeGPT, not JAINE).
- **Architectural decision:** HL = primary execution; 0G = sealed intelligence + receipts; TradeGPT = competitor and optional dust only.
- **Kill:** 0G-only DeFi terminal; JAINE as “ETH/USDC venue.”
- **Next:** demand map Area B.

---

## M4 — Demand is a terminal, not more alerts

- **Source:** `research/DEMAND_MAP_2026-08-24.md` Area B.
- **Discovery:** Users pay TradingView/Nansen for the **terminal**, mute alerts, pay Discords because bots were gamed. Almost nobody treats Nansen watchlists as MNPI. Privacy that sells: **don’t paste the book into ChatGPT** + **don’t give a bot withdraw**.
- **Kill:** alert-accountability product (WARD family); public forecast tape (DOCKET).
- **Next:** 50 internal shapes.

---

## M5 — Fifty candidates (internal; not separate files)

Kill codes: `4LPHA` `VANTA` `MANEKI` `ZIA` `PRED` `YLD` `WALLET` `JOURN` `MKT` `SAT` `LIQ` `FIRE` `NET`.

1. Autonomous HL bot overnight — `VANTA`/`MANEKI`  
2. Policy vault LP on Zia — `4LPHA`/`ZIA`  
3. iNFT sellable trader — `VANTA`  
4. SIA compiles English to orders — `FIRE`/`VANTA`  
5. TEE-born spending key — `VANTA`/`FIRE`  
6. Flashbots-protected 0G mempool — stack missing  
7. JAINE ETH/USDC desk — `LIQ`  
8. TradeGPT skin — `ZIA`  
9. Clark-style autonomous rebalancer on 0G — `ZIA`/`4LPHA`  
10. Polymarket CLOB agent — `PRED`; scoped creds weak  
11. Kalshi API bot — custodial  
12. Copy-trade Nansen Smart Money — late, `SAT`  
13. Yield optimizer — banned/weak  
14. Generic wallet spend limits — feature not product `WALLET`  
15. Safe Zodiac on 16661 — no modules  
16. 4337 session on 16661 — no bundler  
17. Prediction divergence terminal only — `PRED` no execute  
18. Funding-rate alert SaaS — mute `SAT`  
19. LLM sizes from prose — `FIRE`  
20. Multi-agent fake swarm — complexity  
21. ERC-7857 desk NFT — Galileo IDs `VANTA`  
22. ERC-8004 as trader score — unused registry  
23. x402 pay-per-thesis — 16661 not on Coinbase x402  
24. PaymentLayer trade escrow — no withdraw  
25. DA-attested fills — DA dead  
26. Telegram execute — injection  
27. Web trading UI with hot key — `FIRE`  
28. Operator backend signs HL — `FIRE`  
29. Maneki stocks-only — `MANEKI`  
30. Agon EdgeVault clone on Galileo — testnet trap  
31. TRAIDE keeper as the product — not a desk  
32. NeoTrade dmg + 0G sticker — `NET`  
33. Confidential MM desk — they stay in-house (intel ceiling F3)  
34. Packaged vault AUM — not user-authorize  
35. Cross-margin all coins — too much blast radius  
36. On-device LLM only — fails 0G necessity  
37. Router `verify_tee` as proof — Hanami-level, insufficient  
38. processResponse as caller-byte bind — false on plaintext  
39. Dust 0G swap as the hero fill — costume  
40. MEV rebate on 16661 — no builder market  
41. Options desk — no venue researched live  
42. Lending loop Aave on 0G — not in llama top  
43. Bridge-and-trade via Oku as core — extra failure  
44. Social trade sharing — leaks book  
45. Leaderboard of PnL on 16661 — tell/alpha leak  
46. Whisper voice orders — STT scheme UNVERIFIED  
47. z-image chart gen — no counterparty  
48. Fine-tune LoRA trader — 48h ack, not v1  
49. MCP authorize from Cursor — `FIRE`  
50. **PIT: human-authorize, sealed thesis, HL session cannot-withdraw, 0G receipt** — **KEEP**  
51. PIT + Base Uniswap as primary — worse demo, keep as v1.1  
52. PIT + PM execution — kill v1 `PRED`  
53. PIT overnight loop with caps — `VANTA` autonomy; v2 maybe  
54. PIT as Zia liquidity manager — `4LPHA`  
55. PIT with builder fee take-rate v1 — skip until master `approveBuilderFee`

**Winner:** #50 PIT.

---

## M6 — Why 84 not 98

- **Blocker:** matching engine is HL; 0G cannot host ETH-size CLOB today.
- **Solution:** do not fake 0G liquidity; make 0G load-bearing for **private book + bind + memory + receipts**.
- **0G Pay/x402:** unused honestly (0/1).
- **Still UNVERIFIED:** live sealed round-trip (S1); extraAgents on a funded account (S5).
- **Next:** Phase 0 S1–S4. Do not implement UI first.

---

## M7 — Planning package written

- **Files:** five markdown files in `D:\route\0g\PIT\`.
- **Decision:** skip ERC-7857; skip Telegram; skip web trading; MCP cannot authorize.
- **Next:** spike S1.

---

## M8 — Revision mandate: Private Alpha OS (not a new category)

- **Source:** user revision brief 2026-08-26; existing five PIT files; `03_COMPETITORS_WINNERS_AND_CODE_REVERSE_ENGINEERING.md` through EOF (Wave 2/3, APAC, Apollo, judge quotes).
- **Change:** keep PIT as a trading/DeFi desk; upgrade from “private HL terminal” to forecast → challenge → execute → calibrate.
- **Rejected:** replacing PIT; journals; EvoEvo clone; 4lpha-better; VANTA iNFT trader.

---

## M9 — NeoSoul reverse-engineering

- **Source:** https://github.com/NeoSoul-AI (10 public repos); clones under `research/clones/`; sites neosoul.ai, neotrade.pro; org README traction **REPORTED**.
- **Discovery:**
  - `rubric-prediction-skill@c3c3985`: deterministic engine v2.3.0; LLM must not hand-compute; `engine` field required.
  - `neotrade-wallet-sdk@4d92cfd`: signing gateway; **size must not live in the generic gate** (live 2026-08-16 sell-cap bug).
  - `neotrade-skills@445db07`: skills are prompt playbooks, mostly Polymarket; journal/debrief is the useful pattern.
  - `evoevo-contracts@9308aa3`: 8004 as identity; UUPS committee oracle — **not** PIT.
  - `0g-builders` README: EvoEvo Identity `0x8004Ae53…` — **not** official 0G 8004 Identity `0x8004A169…`.
  - neotrade-release: artifacts only; closed runtime.
- **Taken conceptually:** structured forecast, outcome journal, performance-weighted roles, gateway structure, 8004 identity tuple, evidence pack *shape*.
- **Deliberately not copied:** EvoEvo PM network, committee oracle, public leaderboard, mnemonic trading key, raw-key x402 package, Buffett-on-PM skills, UUPS, Brave-search product.
- **Twitter:** agent-reach doctor twitter CLI missing; used GitHub + site fetches.

---

## M10 — Stack revalidation 2026-08-26

- **Models:** 29 SKUs. TeeML still `glm-5.2`, `0gm-1.0-35b-a3b`, `0gm-1.0-35b-a3b-sia`, whisper, z-image. `glm-5.3` is **TeeTLS**, not TeeML. Claude/GPT have **no verifiability**.
- **HL:** meta 200, 232 perps.
- **0G TVL:** $2,871,774 llama `zero-gravity`.
- **Herald x402:** `eip155:16661` exact+upto **FACT**. Coinbase x402.org: no 16661.
- **8004:** official Identity + Reputation proxies have bytecode.
- **Attestor:** Galileo 16602 only.
- **DA:** eth_call out of gas — still unusable.
- **Storage:** official Go client AES-256-CTR / ECIES — supersedes custom AES note in M7-era PIT.
- **Decision:** include 7857 (production IDs, iTransfer reverts honestly); include 8004 reputation as Brier feedback; optional Herald x402 for public evidence; DA excluded.

---

## M11 — Score change 84 → 92

- **Why up:** forecast engine, committee weights, calibration pipeline, 8004 writes, 7857 production IDs, official storage encryption, Herald x402 mapped to a real (optional) purchase, denser 16661 txs, clearer “not HL+0G” story.
- **Why not 95:** DA 0; iTransfer cannot be live; fill remains HL; residual subaccount risk; thin business/moat.
- **Hard impossibilities documented.** Do not inflate.

---

## M12 — Canonical files rewritten

- **Files:** five PIT markdown files updated in place. No extra planning files. No application code.
- **Next:** Phase 0 S1 (funded LedgerManager sealed VerifyE2EE).

---

## M13 — Feasibility + competitor reverse-engineering gate (2026-08-26 evening)

- **Files:** five canonical PIT markdown files updated in place. No product UI. No `pit/` application tree started.
- **Env:** gitignored `PIT/.env`. Payer `0xBDfCeE82Bd42FEfA58ee850B3709636a8B6b0034` (funded 0G key). New PIT-only deployer `0xaAE3EAC0d6665832fe0E5036d61CE2DBC6ECAC2a`. `PIT_MEMORY_KEY` generated. `PIT_MASTER_ADDRESS` / `PIT_HL_SUBACCOUNT` empty. Session HL key never written to env.
- **Clones (depth 50) under `research/clones/`:** Axiom `6f0c496`, 4lpha `a271c79`, Knole `d0855a5`, Hanami `4ec4865`, Talos `34a7bf6`, OptimIEra `0c8d6cb`, diversifi `865d438`, cognivern `a13a81c`, NEXUS `6cd254c`, Lumen `47ee2a9`, TRAIDE `1a28411`, Maneki `c3847c2`.
- **DA re-test:** prior “eth_call OOG” is **false**. `epochNumber()=0`, `params()` succeed, `quorumCount(0)=0`, no mainnet DAEntrance (code length 0). NEXUS “0G DA” is docs-only (Solidity mapping, no DA call). Coal stamps epoch 0. Lineage DA sink throws. **Keep DA = 0. Do not claim DA.**
- **7857 re-test:** attestor `https://agenticid.0g.ai/config` still `chain_id: 16602`. Lumen honest `OracleNotLive()`. 4lpha omits iTransfer from prod ABI. Axiom software TEE PK. Talos always-true verifier. NEXUS trusted ECDSA. **Keep iTransfer deduction. Copy Lumen honesty.**
- **S1 GREEN:** Direct TeeML glm-5.2 `VerifyE2EE` scheme `zg-sig-v1/e2ee-ct`. Provider `0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D`. teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`. Ledger deposit `0xf98e5c4a…`; `transferFund` `0xdad09e14…` with serviceName **`inference-v1.0`** (`"inference"` reverts). `ZG-Res-Key` `fec6251d-4708-451a-8130-a9551680ded2`. Not Router. Not `processResponse` as bind. Direct `-onchain` still rejected.
- **S2 GREEN (gate unit):** `BindPreview` ignores model `sz=100`. Promote into `core/internal/preview` on implementation start.
- **S3 PARTIAL:** SDK classes confirmed (user-signed vs L1). Ephemeral-agent withdraw3 is treated as that EOA being the user — **never fund session key**. Code allowlist not yet in product. Funded `approveAgent` not done (HL equity $0).
- **S4 NOT READY:** no HL USDC / subaccount. Not a stack kill.
- **S6 GREEN:** official CLI encrypt+`--proof` roundtrip; wrong root fail; wrong key fail. Root `0x2f0419508df7dc546f9aa84d214723e72cc2547fa1f69d828c9187510cfbda5d`. `--encryption-key` needs `0x`.
- **8004 GREEN (with architecture constraint):** register agentId **3489333** tx `0xa2f67529…`. Owner `giveFeedback` reverts `Self-feedback not allowed`. Deployer client feedback tx `0xd4892d63…`. Do not use NeoSoul Identity `0x8004Ae53…`.
- **Catalog:** 29 models; TeeML 5 / TeeTLS 16 / none 8. glm-5.2 still TeeML. Claude/GPT still forbidden on book path. Herald still `eip155:16661`. HL 232 perps.
- **Score:** **92/100 unchanged** on the 35-dim product sheet. Execution-architecture rubric (M14) scores Option D at **86/100** — different sheet, do not add them. S1/S6/8004 were already at cap; proving them makes 92 executable, not 93/95. Hard ceilings remain DA / iTransfer / fill-off-0G for the **product book**.
- **Next implementation step (when user says go):** scaffold `D:\route\0g\pit\` contracts + Go core; promote S1 then S2; write HL allowlist fail-closed **before** any session key; Phase 1 contracts with honest iTransfer revert; **S4 univ3 dust on 16661 before HL USDC**. Frontend last.

---

## M14 — Execution architecture decision (2026-08-26)

- **Question:** Hyperliquid-first vs 0G-native-first vs multi-venue router vs execution-neutral vs other.
- **Live venues (treat Hub screenshots as leads, then probed):**
  - Hub Swap “Powered By JAINE” FACT UI. JAINE llama TVL **$2,914**. Hub “JAINE” page addr `0x7bbc…1404` is **st0G**, not a router.
  - Hub Bridge Ethereum→0G W0G, Chainlink lead. **Onramp, not a DEX.**
  - Hub Cross-Chain Swap = Khalani/TokenFlight. Sample history **5 0G → 5 0G** same-chain. **Intent rail, not a CLOB.**
  - Oku public app defaulted to **Ethereum** USDC/WETH; `has0g=false` in swap UI. Oku **contracts** on 16661 have bytecode (SwapRouter02 `0x807F…dE40`).
  - Morpho Blue on 0G: **$32.20** TVL / **$2.32** borrowed. Ignore.
  - Pool `balanceOf` block **42656196**: Zia USDC.e/WETH **334,519 + 132.48 WETH**; Oku W0G/USDC.e **502,586 W0G + 43,214 USDC.e**; Zia W0G/WETH **dead** (0.000426 WETH). Llama uniswap-v3 0G slice **timed out** UNVERIFIED.
  - HL: 232 perps; ETH mark **$2462.1**; day ntl **~$1.52B**; OI **~671k**; `szDecimals=4`; min notional **$10** gitbook FACT. Testnet faucet requires **prior mainnet deposit**. Payer HL equity **$0**.
- **Competitors (source, not README):**
  - 4lpha `curated-routes.ts` already owns Zia+Oku routers/pools; `package.json` has storage SDK **not** compute SDK; judge: finish V4, leave sandbox.
  - Axiom: MockUSDC `0x354CA53b…`, `InMemoryStorage` if indexer unset.
  - NEXUS `CompositeReceiptMinter.sol`: mapping + events, **no DA call**; Galileo-only deploy script.
  - Maneki public repo: HL `/info` only; cannot place orders.
  - TRAIDE README: “not live trading”; Galileo.
  - Lumen: honest `OracleNotLive` on iTransfer; not a desk.
  - Hanami: Direct env-gated, falls back to Router; `processResponse`.
  - Cannes Shawarma: Uniswap Trading API on **Arbitrum 42161**; 0G is compute.
- **PIT Router hypothesis:** **rejected.** Perp vs AMM vs bridge are not interchangeable. Judge does not reward integration count. 4lpha already is the 0G DEX router.
- **Decision:** **Option D — execution-neutral core, HL-default product book, labeled `0g:univ3` dust adapter, no smart router.** Overturn “HL is the product identity.” Overturn 0G-native-first. Overturn multi-venue best-ex.
- **Scores (user execution rubric /100):** A 78 · B 67 · C 63 · D **86** · E forecast-only 64. 35-dim product sheet stays **92**. Not 95. Path to 98–100 is **external** (DA, Aristotle attestor, 0G CLOB).
- **Kill:** JAINE as venue; Khalani as venue; Morpho; Hub Cross-Chain as trade; 4lpha PolicyVault clone; smart router.

---

## M15 — Feasibility + security + execution gate (2026-08-26)

- **Chrome (real):** mainnet portfolio `0xBDfC...0034`, Total Equity **$4.80**, no positions. Testnet trade: Unified, **999.00 USDC** available, Isolated 10x shown on HYPE book (not used).
- **API vs UI:** unified accounts keep USDC in **spot**. `clearinghouseState.accountValue=0` on both nets is expected (official account-abstraction docs). Not a zero-balance bug.
- **S3 code:** `PIT/_gate/hlcore` `go test ./...` PASS. `PitExecutor` `local_all_denied=true`. Session keys in `_gate/secrets/` (gitignored), never `.env`, never printed.
- **approveAgent:** master-signed. Testnet agent `0x71e9d258E3bC46A937C126890929CBccEDC7272D` name `PIT-M15` `validUntil` 1h. Mainnet agent `0x0b68BA7B3B79AB8357E6C34749459eC4b44d9F02` name `PIT-M15m`.
- **Protocol hole:** agent `updateLeverage` **succeeds** if PIT is bypassed (testnet ok). Restored via master. Allowlist is load-bearing.
- **S4 testnet:** first order rejected `Price must be divisible by tick size`. Retry with 5-sigfig round: limit **1597.2**, sz **0.0075**. oid **58551048563** resting, status open, cancel success, status canceled. USDC 999→999. Positions [].
- **S4 mainnet:** intended printed then sent. limit **1601.0**, sz **0.0063**, notional **$10.0863**. oid **526980359929** resting, cancel success. USDC 4.8→4.8. Chrome still $4.80 / no positions.
- **S1 re-run:** Direct TeeML `VerifyE2EE` OK. scheme `zg-sig-v1/e2ee-ct:c1c9a15b…:0a0f8903…`. `ZG-Res-Key` `af5ca525-bedb-4f60-80d2-ada008be667b`. Signer `0xA46EA4FC…` matches on-chain teeSigner. Prompt was `PONG` not a book committee.
- **S6 re-run:** root `0x3b4b3b772aae7195109ded219ca861da7eb3ca51776538e3486f7084b4ef193a`. Roundtrip match; wrong root fail; wrong key fail.
- **8004:** ownerOf `0xBDfC…0034`, URI `ipfs://pit-desk-v1`. Owner `giveFeedback` mined **status 0** (revert). Do not treat as success.
- **Chrome PIT app E2E:** **not run** — no product UI (correct; UI forbidden this pass).
- **Committee-over-book E2E:** **partial** — executor+sizer+preview proven; 3 sealed roles not wired yet.
- **Competitor upgrades taken as plan, not copy:** persist cloid/`usedActionHashes` (4lpha + NeoSoul gateway). Do not copy PolicyVault.
- **Score:** 92 unchanged. Not 95.
- **Verdict:** **NOT ready for product UI.** **Ready to start `pit/` contracts + Go core** (promote `_gate/hlcore`).

---

## M16 — Pre-implementation planning gate (2026-08-26)

**No product implementation.** Canon rewritten in place: `FINAL_FLAGSHIP_PRODUCT.md`, `IMPLEMENTATION_PLAN.md` (phases 0–29), `RESOURCES.md`, `ENVIRONMENT.md`. This report.

### Product reopen (first principles)

- Ask-only + fixture `PIT_MASTER_ADDRESS` would have shipped a **one-wallet demo**. That is now an explicit kill.
- Daily use requires **Opportunity Watch** (policy-eligible cards, user AUTHORIZE) + **Strategy Health** (calibration as a surface, not a hidden JSON).
- **Counterfactual Lab** deferred: theater without N≥30; PASSES already journal `would-have`; half-size is a policy card.
- KEPT (`D:\route\Flare\kept`, https://kept-ruby-five.vercel.app/) is the visual/onboarding bar. Coral `#D82F2F`, cream `#F0E7D4`, Host Grotesk, split WalletGate, Get started cards, named failure states. PIT ring is **not** a savings circle.
- Orchestra clone (`research/clones/Orchestra`): Groq + AUTO_EXECUTE — confirms Watch demand; **refuse** their execution model.

### Added

- Opportunity Watch (search/recommend/prepare; never sign).
- Strategy Health (surface existing stats).
- Multi-user workspace isolation; SIWE; per-user keychain; web cannot hold HL session.
- HL onboarding for non-experts including unified-spot USDC truth and protocol `updateLeverage` disclosure.
- Durable exactly-once ledger plan (4lpha/NeoSoul ideas, PIT-native SQLite).
- KEPT-grade 13-beat first run + welcome-back Watch home.
- ENVIRONMENT three-class secrets (GLOBAL / PER-USER / DEV FIXTURE). Fixture payer is not the product owner.

### Killed

- Overnight auto-trade; smart router; Telegram; copy-trade; social board; yield; meme universe; fake DA/iTransfer/hardware TEE; AI size/calldata; Counterfactual Lab v1; global operator key; session in `.env`/localStorage/Vercel; Orchestra AUTO_EXECUTE; Groq.

### Score

- **Proven:** **92**. Unchanged. Do not count planning as TEE/Storage points.
- **Planned after Chrome E2E of this spec:** **93** (UX row), still not 95. Ceilings: DA=0, iTransfer IMPOSSIBLE, fill≠0G CLOB.

### Revalidation this pass

- Agentic ID config re-fetch: `chain_id: 16602`, RPC `https://evmrpc-testnet.0g.ai` — **FACT**. iTransfer on Aristotle still **IMPOSSIBLE**.
- 03 file used through competitor/showcase sections including Orchestra / Shawarma / ETHGlobal Cannes list (line ~526). Judge pattern unchanged.
- KEPT source tokens read from `guide.css` / `WalletGate.tsx` / `StartFlow.tsx` / `namedStates.ts` / `keptGuide.tsx`. Live site not re-crawled this pass (screenshots + local source sufficient) — live CSS **UNVERIFIED** vs deploy, local is FACT.

---

### FINAL GATE REPORT (M16)

#### Green (FACT from prior gates + this planning)

| Item | Status |
|---|---|
| Architecture Option D execution-neutral | FACT M14 |
| Direct TeeML `VerifyE2EE` (PONG, not book committee) | FACT M15 |
| Storage encrypt + `--proof` | FACT M15 |
| ERC-8004 register agentId 3489333 | FACT (fixture owner) |
| HL allowlist order/cancel | FACT `_gate/hlcore` tests |
| HL testnet oid 58551048563 cancel | FACT |
| HL mainnet oid 526980359929 cancel | FACT |
| `PIT_ALLOW_FALLBACKS=false` required in plan | SPEC |
| Canon: multi-user, Watch, KEPT, phases 0–29 | SPEC this pass |
| iTransfer attestor Galileo-only | FACT re-fetch |

#### Red / not implemented (must not be called GREEN)

| Item | Status |
|---|---|
| Product `pit/` tree / Tauri / web | **NOT IMPLEMENTED** |
| Contracts issued on 16661 | **NOT IMPLEMENTED** |
| Durable cloid ledger | **NOT IMPLEMENTED** (in-memory gate) |
| 3-role committee **over book** | **NOT IMPLEMENTED** |
| Opportunity Watch runtime | **NOT IMPLEMENTED** |
| Multi-user E2E | **NOT IMPLEMENTED** |
| Playwright / Chrome product E2E | **IMPOSSIBLE** until UI |
| Product CI | **NOT IMPLEMENTED** |
| Per-user memory keys (not fixture) | **NOT IMPLEMENTED** |
| 8004 owner self-feedback | **RED** (revert); use reporter |

#### Unverified

| Item | Status |
|---|---|
| Llama uniswap-v3 0G TVL line | UNVERIFIED (timeout); pool balances FACT |
| KEPT production CSS vs local | UNVERIFIED (use local) |
| Orchestra live 0G vs Groq in their demo | UNVERIFIED; README admits Groq |
| Second independent user onboarding | UNVERIFIED (never run) |

#### External blockers

| Item | Status |
|---|---|
| 0G DA mainnet | **IMPOSSIBLE** (epoch 0, no DAEntrance) |
| ERC-7857 iTransfer on 16661 | **IMPOSSIBLE** (attestor 16602) |
| HL protocol: agent `updateLeverage` | **IMPOSSIBLE** to disable on-chain; mitigate TTL + allowlist + disclose |
| HL min order $10 | Constraint |
| Testnet faucet | Requires prior mainnet deposit |
| Fill on 0G CLOB | **IMPOSSIBLE** at ETH size (JAINE ~$2.9k) |

#### Optional

| Item | Status |
|---|---|
| Labeled 0G univ3 dust | OPTIONAL |
| Herald x402 public pack | OPTIONAL (`PIT_X402_ENABLED=false`) |
| SOL on allowlist | OPTIONAL after ETH/BTC calibration |
| Counterfactual Lab | OPTIONAL v1.1 |
| Builder fee | OPTIONAL, master-signed later |
| Subaccount isolation | OPTIONAL recommended |

#### User-required actions (implementation)

1. Do **not** paste user keys into server env.
2. Each real user: wallet with 16661 native for Ledger/pins; HL account with ≥$10 USDC (or research-only).
3. Testnet users: mainnet deposit first, then faucet.
4. Demo wallet `0xBDfC…0034` remains a **fixture** only.
5. Deployer `0xaAE3…C2a` for contracts + 8004 **reporter**, never Identity owner of user desks.

#### Env (classified)

See `ENVIRONMENT.md`. Infra: RPC, contract addresses, `PIT_ALLOW_FALLBACKS=false`, HL public URLs, deployer key on deploy machine. Per-user: never in `.env`. Fixture keys stay gitignored.

#### Chrome tests to run **after** UI (Phase 23)

1. Landing → Get started (coral/cream, no dashboard).
2. Connect wallet — decline, wrong network, success.
3. HL identify — spot USDC shown if unified.
4. Policy cards → pin tx on 16661.
5. Session create → approveAgent in wallet → permissions card + countdown.
6. Watch empty vs Watch ≥1 **real** candidate (not seeded).
7. VIEW THESIS: ring phases match network timing.
8. AUTHORIZE dust testnet order + cancel.
9. `/verify` recomputes TeeML signer.
10. Sign out / second profile: no leaked cards.
11. Reduced motion.
12. Web build: cannot AUTHORIZE HL.

#### Mainnet / testnet tests (Phase 9/26)

- Repeat M15 dust on **a second address**.
- Contract verify on chainscan.
- Storage `--proof` with **workspace B** key failing on A’s root.

#### Explorer evidence already in RESOURCES

Ledger, transferFund, 8004 register, 8004 client feedback, Storage roots, HL oids. New evidence only after new txs.

#### Competitor checks to repeat after implementation

Axiom still mocked? 4lpha still vault? Knole still no fill? Hanami still Router fallback? Lumen still honest iTransfer? Maneki still sticker? TRAIDE still Galileo? Orchestra still Groq? PIT still no DA/iTransfer claims? PIT `/verify` live? Two workspaces isolated?

### Planning checklist (user Part 28)

1–18: **GREEN** for planning artifacts. Product code **not** started (correct).

### Verdict

**PRE-IMPLEMENTATION GATE CLOSED — READY** to start `IMPLEMENTATION_PLAN.md` Phase 0 then Phase 1 (`pit/internal/identity`).

**NOT READY** to ship UI, to skip to desktop, or to call the repo complete.

**Start command:**

```powershell
cd D:\route\0g\PIT
cast chain-id --rpc-url https://evmrpc.0g.ai
# then Phase 1 Go workspace identity — not npm create, not Tauri first
```

---

## M17 — Hardening: Agentic ID + real sealed committee (2026-08-26)

**Source:** live Aristotle RPC; `https://agenticid.0g.ai/config` (Chrome); `pc.0g.ai/models` Private TeeML 5; Hub `hub.0g.ai` wallet `0xf76e…71a3`; HL Chrome `0xBDfC…0034` ETH mark **2445.9** (committee snapshot was **2438.5**); ChainScan deploy+mint; Foundry `PitDeskID`; Direct TeeML glm-5.2; official `0g-agentic-id` HEAD `b8f4845`; clones Lumen/Knole/4lpha/NEXUS/Talos.

### FACT

- Official Foundation attestor **Galileo 16602**. Chrome config JSON: `chain_id: 16602`, `chain_rpc: https://evmrpc-testnet.0g.ai`, AgenticID `0x3449…5648`, verifier `0x9D48…f980`, frameworks `openclaw` / `hermes` / `prime-agent`.
- `eth_getCode` Foundation AgenticID+verifier: Aristotle **0**, Galileo **295**. PitDeskID: Aristotle **7071**, Galileo **0**.
- SDK `CHAIN_ID = 16602`. `ZERO_G_MAINNET` exists; comment says no mainnet AgenticID deployment.
- `DEPLOYMENT.md` §2.4: seal-bound agents **re-enable** ERC-721 `transferFrom` and **`iTransferFrom` reverts**. Galileo sealed runtime is not Aristotle iTransfer.
- Official IERC7857Authorize is an **off-chain usage allowlist**, not an on-chain transfer gate.
- Docs ERC-7857 page is a sketch (`transfer`/`clone`) ≠ production `iTransferFrom` in source. PIT follows source IDs `0x2afbede9` / `0xdf597d99` / `0x74f8628b` / `0x80ac58cd`.
- Router Private TeeML **5**. glm-5.2 Direct: provider `0x7DCF…87D`, teeSigner `0xA46E…46B9`, url `compute-network-19.integratenetwork.work`.
- HL fixture: $4.80 USDC, no positions. Min notional $10.

### VERIFIED LIVE (explorer)

| Event | Tx / id |
|---|---|
| PitDeskID deploy | `0x2d5b688bf09bb72cb44b092da0c27cbe87a623141872e62b84cb95ecf7e90c24` contract `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` block 42705540 |
| mint tokenId 1 | `0x9494e3faec6d950942d1bfec53c4a13e6f28378da8e01ecd24b3e3de62e5c7d0` owner `0xBDfC…0034` |
| authorizeUsage | `0xdf13b262a875602386708f4c4a121e8ab039e28849d8148e8b991d26fff080d4` |
| unauthorized authorize | `0x74c8ead0…fa8dc5e` status 0 (expected) |
| revoke | `0xfd2e38b62f72fd8f23367aa6194991a91e833429f888b4c9c503ef8cb5440af9` |
| iTransfer/iClone | revert `0x0c9c36b4`; transferFrom revert `0x72e59941` |
| Foundry | 7/7 PASS |
| Ledger top-ups | `0x56eb0d63…3c3a03` (0.6 0G); `0xf3555b9a…218e09` (0.5 0G) |

### VERIFIED IN SOURCE (competitors — do not copy)

- **Lumen:** honest `OracleNotLive` + live `authorizeUsage`.
- **Knole:** custom id `0x4b396f04`; owner-move `iTransferFrom` without Foundation attestor.
- **4lpha:** omits iTransfer from prod ABI.
- **Talos:** always-valid verifier.
- **NEXUS:** `authorizeUsage(executor, permissions)` — different shape.

### Committee (real ETH book, not PONG)

Independence: **prompt_envelope_only_SAME_PROVIDER** (fixture Ledger only on glm-5.2).

| Scenario | Sealed VerifyE2EE | Engine |
|---|---|---|
| S1 $15/1x live tape | researcher / challenger / risk **OK** (`7a705ed4…` / `e9e81db2…` / `200d4933…`) | Researcher **sell**; Challenger **killed**; `eligible=false` `challenger_killed` |
| S2 forced 50x | retry **OK** (`101562ca…` / `aeca6159…` / `38810a5f…`) | Researcher `none`; `no_side` |
| S3 $5 policy | retry **OK** (`fc0f4743…` / `ec921a87…` / `8b7110d9…`) | sell + survives; `below_min_notional` |

Book sha256 `05de115772903f1e47751852adf3c31fcf3c1817124a0fbf5af58b8ebb3763cc`. Engine v2 (not v1: v1 matched `"thesis_killed"` substring even when false; v1 hardcoded side buy).

A sealed run with `eligible=true` was **not observed**. Engine unit: committee-agree + $15 → eligible **PASS**.

### Security matrix

| Attack | Result |
|---|---|
| 1 Router downgrade | **PASS** (committee.exe exit 10) |
| 2 Plaintext fallback | **PASS** (no Router path; `Allow-Fallbacks: false`) |
| 3 Wrong model | Expected fail-closed on Direct URL/model bind — **not live-swapped** this gate |
| 4 Wrong provider | Same — **not live-swapped** |
| 5 Wrong signer | **PASS** (recovered `0xA46E…` ≠ `0x000…001`) |
| 6–7 Response / ciphertext tamper | **PASS** (content-binding mismatch) |
| 8 Prompt mutation | VerifyE2EE binds request bytes — expected FAIL; **not separately hashed post-seal** this gate |
| 9 Replay | Serving generation/nonce — **specified**, not e2e this gate |
| 10 Wrong workspace | 7857 unauthorized authorize **PASS**; Storage cross-key prior S6 **PASS**; committee host check **product** |
| 11–12 Stale / mid-committee policy | Engine re-check rule **PASS**; do not apply S3 committee JSON under a $15 policy without a new sealed run |
| 13–14 Size mutation / LLM sz | **PASS** (sizer; LLM 9.99 ignored) |
| 15–18 Session expiry, kill, logout, crash | Kill switch engine **PASS**; session/crash/logout = Phase 10/25 **not this gate** |

Private prompts / sealed_input JSON **redacted** from evidence (sha256 stubs). Role outputs may quote public mark/funding and policy caps.

### IMPOSSIBLE

- ERC-7857 **iTransfer / iClone** on Aristotle 16661 (Foundation attestor absent).
- Claiming three independent TeeML models while only glm-5.2 is funded.
- 0G DA; fill on 0G CLOB; fake TEE.

### OPTIONAL

- Fund 0GM / 0GM-SIA Ledger subaccounts for real model diversity.
- Galileo sealed-agent lane (OpenClaw) as a **separate** product — not Aristotle iTransfer.
- Non-standard transferable wrapper, explicitly labeled, **not** claimed as ERC-7857 iTransfer.
- Labeled HL dust S4 (fixture $4.80 < $10 min; skipped).

### Verdict

**iTransfer = NOT LIVE.**

**Hardening primitives: proven.** Promote `_gate/desk7857` + `_gate/committee` in Phases 14 and 6.

**READY** to continue `IMPLEMENTATION_PLAN.md` Phase 0→1 (unchanged). **NOT READY** to ship UI, to claim S1-eligible on this ETH tape, or to call iTransfer live.

**Proven score: 92.** Do not inflate: no UI, envelope-only independence, S1 eligible not observed, DA/iTransfer ceilings unchanged.

---

## M18 — Final 0G / environment / testnet–mainnet revalidation (2026-08-26)

Research + verification only. Product UI **not** started. No Router API key created. No secrets recorded.

### Official sources (fresh)

- docs.0g.ai: Router overview/auth/comparison/privacy; Direct inference.md; testnet-overview; mainnet-overview; ai-context Agentic ID
- build.0g.ai (index)
- pc.0g.ai/models, /sdk, /dashboard/overview, /dashboard/api-keys (Chrome, read-only)
- pc.testnet.0g.ai/models (Chrome)
- agenticid.0g.ai/config (Chrome + curl)
- hub.0g.ai/discover?network=testnet (Chrome)
- chainscan.0g.ai PitDeskID + mint tx (already open)
- github.com/0gfoundation clones: `0g-agentic-id` `b8f4845`, `0g-pc-e2ee` `0ba9fb0`, `0g-compute-ts-sdk` `bff2913` v0.9.0, `0g-storage-client` `3d953af`
- Live RPC: Aristotle + Galileo `eth_chainId` + `eth_getCode`
- Live catalogs: Router mainnet N=29; Router testnet N=2; Galileo `getAllServices` total 6

### Clones rechecked

Hanami Direct still optional and **falls back to Router**. Lumen still OracleNotLive + dual 16661/16602 honesty. No winner showed a live Aristotle Foundation attestor.

### API key

Dashboard API Keys = **Router `sk-` only** (Chrome copy: OpenAI-compatible API). Trust modes Standard/Verified/Private are Router routing tiers. Management Keys = `mk-`. Direct uses wallet-signed `app-sk-` (tokenId 255 ephemeral). **PIT confidential book does not need a dashboard key.** Existing truncated `sk-` prefixes were visible in Console; values **not** copied into env/git/history.

### Catalog / roles

Mainnet Private TeeML still 5. Fixture Direct committee still **only glm-5.2 funded**. Galileo Private filter = **1 image model**; Direct ack’d TeeML **chat** = `qwen/qwen2.5-omni-7b` (Router labels Omni TeeTLS — ignore). No glm-5.2 on Galileo.

### Agentic ID / iTransfer

Attestor still `chain_id: 16602`. Aristotle Foundation AgenticID + verifier **code 0**. Galileo **295/295**. PitDeskID 16661 **7071**, 16602 **0**. iTransfer on mainnet **IMPOSSIBLE**. iTransfer on Galileo = official path, **PIT tx UNVERIFIED**. ERC-721 fallback still not official for DeskID. No new mainnet verifier. SDK still Galileo-default.

### Dual environment

Ship **Mainnet + Testnet**. Mainnet = production-safe subset. Testnet = protocol laboratory with a **different** compute catalog. Do not copy SKUs across networks. Fixture `.env` `PIT_DESK_ID_CONTRACT` pinned to live Aristotle address (public). `PIT_HL_NETWORK=Testnet` remains a gate leftover, not product default.

### Files updated

`FINAL_FLAGSHIP_PRODUCT.md` · `IMPLEMENTATION_PLAN.md` (capability matrix + per-phase TN/MN) · `RESOURCES.md` · `ENVIRONMENT.md` (A–F rewrite) · this append.

### Verdict

**READY** to implement Phase 0→1 with dual-env contracts. **NOT READY** to ship UI. **NOT READY** to claim Galileo iTransfer LIVE or Galileo committee LIVE. **Router API key not required.**

**Proven score: 92.** Unchanged. Dual-env honesty is a documentation win, not a score bump.

---

## M19 — Production implementation start (2026-08-26 21:56+03)

DATE/TIME: 2026-08-26
PHASE: 0 then 1+ core libraries
GOAL: Close planning. Implement identity, policy, session allowlist, engine, ledger, compute deny-router, storage namespaces, CLI/MCP/SDK skeletons, DeskID+Policy contracts, wallet-connect web shell.
FILES CHANGED: `pit/**`, `contracts/src`, `contracts/test`, `apps/web/**`, `.env.example`, `.gitignore`, `README.md`, `.github/workflows/ci.yml`, this file.
COMMANDS RUN: live catalog + chainId + attestor + HL meta; `go test ./...`; `forge test` (next).
EXTERNAL RESOURCES CONSULTED: 0G docs Direct vs Router; Hyperliquid info endpoints; Privy wallet connect (app id only in frontend).
RESULT: Phase 0 facts unchanged (Aristotle 16661, Galileo 16602, attestor 16602, glm-5.2 TeeML still listed). Product Go module created. Router URL compile-time denied for private book.
TESTS: `go test ./...` all packages ok (ledger, siwe, workspace, policy, engine, session, compute, mcp, sdk). Foundry: 10 passed (7 DeskID + 3 Policy including fuzz).
EVIDENCE: this entry; no secrets logged.
TX HASH / OID: none this step (library + contracts tests).
FAILURE: none yet.
FIX: n/a
RETRY: n/a
FINAL STATE: implementation in progress. Web can connect a wallet; it cannot hold a session key.
NEXT STEP: `go test ./...`, Foundry tests, then continue sealed committee wiring and live execution tests.

Git: 30 commits on `main` pushed to https://github.com/mohamedwael201193/pit (HEAD `484cd0d`). `.env` not committed. Continue remaining phases; subsequent work ships as 15 commits dated the work day.

---

## M20 — Library phases 2/4/6–19/21/25 (2026-08-26 21:14+03)

DATE/TIME: 2026-08-26 21:14+03
PHASE: 2 wallet states · 4 session keygen · 6 committee/verify · 7 forecast · 8 watch safety · 9 exec gateway · 10 ledger recover · 11 storage proof args · 12 memory kinds · 13 8004 host · 14 DeskID host gate · 15 calibration health · 16 MCP stdin · 17 SDK no-session · 18 CLI TTY authorize · 19 Receipts/Forecasts/Memory contracts · 21 web onboarding · 25 S3 red-team matrix
GOAL: Land fail-closed libraries and contract tests for the remaining core loop without claiming live TeeML/HL/Storage txs.
FILES CHANGED: `pit/internal/wallet/` `session/keygen*` `compute/committee.go` `compute/verify.go` `engine/forecast*` `exec/` `ledger/recover*` `storage/proof*` `memory/` `chain8004/` `deskid/` `calib/health*` `watch/safety*` `cli/` `verify/` `redteam/` `cmd/pit/main.go` `cmd/mcp/` `mcp/mcp.go` `sdk/` `contracts/src/PitReceipts.sol` `PitForecasts.sol` `PitMemory.sol` + Foundry tests · `apps/web/src/App.tsx` `styles.css` · `.env.example` `IMPLEMENTATION_PLAN.md` `README.md` this file.
COMMANDS RUN: `go test ./...` (all packages ok) · `forge test` (13 passed: 7 DeskID + 3 Policy + 2 Receipts/Forecasts + 1 Memory) · local `.env` public dual-env keys appended (testnet RPCs/addrs, empty receipts/forecasts/memory placeholders, `VITE_PRIVY_APP_ID`). Secrets not logged, `.env` not committed.
EXTERNAL RESOURCES CONSULTED: Direct vs Router (deny `router-api` on book path) · Hyperliquid extraAgents / order|cancel allowlist · 0G storage Go client `--encryption-key` `0x` + `--proof` · ERC-8004 reporter ≠ owner · ERC-7857 `isAuthorized` + Aristotle `AttestorNotOnAristotle`.
RESULT: Named connect states. Session export denied. Envelope+scheme `zg-sig-v1/e2ee-ct` + teeSigner match required. Host `p` ignores model probability. Gateway denies withdraw/leverage/sendAsset/approveAgent. Recover never blindly reposts signed/timeout. Storage TS client forbidden. Memory rejects global `PIT_MEMORY_KEY` in product mode. 8004 IDs not portable. Desk transfer refused on mainnet. Health card empty until N≥30. CLI piped authorize denied. Web shows YOU/YOUR steps and real progress labels; session keys still forbidden in browser.
TESTS: Go unit all packages pass. Foundry 13/13.
EVIDENCE: this entry. No private keys, tokens, or session material in git.
TX HASH / OID: none (libraries + local Foundry). Live Seal+VerifyE2EE, HL order, Storage upload, user 8004/7857 txs still open.
FAILURE: first `storage` proof test used Python string-repeat syntax; fixed to `strings.Repeat` (64 hex chars + `0x`).
FIX: `proof_test.go` key length 66.
RETRY: `go test ./internal/storage/...` pass.
FINAL STATE: core loop is coded fail-closed. Product binary still does not send a live Direct TeeML request or a live HL order.
NEXT STEP: vendor/promote Direct HPKE Seal+VerifyE2EE into `pit/internal/compute`, live testnet dust order, Storage `--proof` upload, desktop shell, then Chrome E2E.

`.env` (gitignored): added public Galileo registry lines and empty `PIT_*_CONTRACT` placeholders for receipts/forecasts/memory. Did not print or commit secrets.


