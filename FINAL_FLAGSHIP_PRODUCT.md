# PIT — Final Flagship Product

**Wave 3 · 0G Aristotle `16661` · submit 2026-08-30T15:00:00Z**

**Architecture lock:** execution-neutral core (M14 Option D). Hyperliquid is the default **live book**. 0G is the **private OS**. No smart router. No overnight auto-trade in v1.

**Pre-implementation gate (M16):** planning rewritten for **arbitrary users**, KEPT-grade onboarding, policy-gated **Opportunity Watch**, and **Strategy Health**. Product code is **not started**.

**Hardening gate (M17):** Agentic ID mint/own/authorizeUsage/revoke is **live on Aristotle**; Foundation `iTransfer` remains **NOT LIVE** on `16661`. Sealed 3-role committee over a **real Hyperliquid ETH book** is **live in `_gate`** (Direct TeeML + HPKE + VerifyE2EE). Product UI is still **not started**. Proven score remains **92**. Planned product (if this spec is implemented without mocks) is **93** — +1 Frequency/UX from daily Watch, not from DA/iTransfer/CLOB (those remain hard ceilings).

**Revalidation gate (M18, 2026-08-26):** Dual environment is **required**, not optional. **Mainnet** = production-safe subset (proven Direct glm-5.2 + PitDeskID mint/authz). **Testnet** = protocol laboratory (Galileo attestor live; Direct TeeML **chat** = `qwen/qwen2.5-omni-7b`, **not** glm-5.2; iTransfer official path **UNVERIFIED** until a PIT transfer tx). Router `sk-`/`mk-` keys are **not required** for the private book. Product UI still **not started**.

---

## Thesis (one paragraph)

PIT is a **Private Alpha OS** for a serious trader who already has capital on Hyperliquid (or is willing to fund a tiny labeled 0G swap). The user owns the money. The user writes a policy a human can read. PIT **searches** the public book, **seals** the private strategy into 0G Compute, **challenges** itself, **scores** with a deterministic engine, and **notifies** when something is eligible. The user **authorizes the exact preview**. A one-hour session key that **cannot withdraw** places or cancels. 0G files the research and the receipt. Outcomes update **that user’s** skill weights. Tomorrow the desk is still theirs.

```
USER OWNS THE MONEY
  → USER DEFINES THE POLICY
  → PIT DISCOVERS OPPORTUNITIES
  → PIT RESEARCHES PRIVATELY IN 0G
  → PIT CHALLENGES ITSELF
  → PIT EXPLAINS THE OPPORTUNITY
  → POLICY CHECK
  → USER AUTHORIZES EXACT ACTION
  → SCOPED SESSION EXECUTES
  → 0G PROVES RESEARCH/RECEIPT
  → OUTCOME RESOLVES
  → PIT LEARNS WHICH SKILLS WORKED
  → NEXT OPPORTUNITY IS BETTER
```

**The agent may search, recommend, and prepare. The agent cannot authorize, cannot withdraw, cannot change policy, cannot compile calldata.**

If 0G Compute, Storage `--proof`, or on-chain pins are down, PIT **stops**. It does not become a ChatGPT wrapper.

---

## Two environments (user-visible, never hidden)

The UI always shows which world the user is in. Testnet-only capabilities **must not** appear as mainnet-enabled.

### TESTNET — “Full test environment”

0G Galileo `16602` + Hyperliquid testnet. Faucet, fake fills, CI, Chrome E2E, Foundation Agentic ID attestor.

| Capability | Status shown to user |
|---|---|
| Chain / explorer | ✅ Galileo |
| Direct TeeML chat | ⚠️ `qwen2.5-omni-7b` on-chain TeeML ack — **not glm-5.2**; PIT VerifyE2EE **not yet run** on Galileo |
| E2EE / on-chain signer | ⚠️ same protocol as mainnet; **unproven on this chain** |
| Storage + proof | ✅ addresses live; product upload not yet |
| Agentic ID mint (PitDeskID) | ❌ gate contract **not** on Galileo (code 0) — deploy testnet DeskID first |
| Foundation Agentic ID | ✅ attestor + verifier bytecode live |
| authorizeUsage / revoke | ✅ after testnet DeskID or Foundation mint |
| iTransfer / iClone | ⚠️ official attestor path — **UNVERIFIED** (no PIT `ownerOf` change tx) |
| ERC-8004 | ✅ Galileo registries live |
| Hyperliquid | ✅ testnet book + faucet |
| Researcher / Challenger / Risk | ⚠️ one funded TeeML chat SKU until proven otherwise — envelope independence only |
| Router catalog | 2 models (image TeeML + Omni **TeeTLS** on Router) — **forbidden** for the book |

### MAINNET — “Production”

0G Aristotle `16661` + Hyperliquid mainnet. Real money. Real receipts.

| Capability | Status shown to user |
|---|---|
| Chain / explorer | ✅ Aristotle |
| Direct TeeML chat | ✅ glm-5.2 provider funded (M17) |
| E2EE + VerifyE2EE + teeSigner | ✅ M17 committee |
| Storage + proof | ✅ S6 |
| Agentic ID mint / own / encrypted hash | ✅ PitDeskID |
| authorizeUsage / revoke / isAuthorized | ✅ |
| iTransfer / iClone | ❌ **not live** — Foundation attestor + verifier **absent** (code 0; config still 16602) |
| ERC-8004 | ✅ |
| Hyperliquid | ✅ real book; min $10 |
| Researcher / Challenger / Risk | ✅ same glm-5.2 provider (envelope independence). 0GM SKUs exist on Router Private-5 but **fixture unfunded** — do not claim independent Direct providers |
| Router `sk-` path | ❌ never used for YOUR private book |

Default after onboarding: **Mainnet / Production**. Testnet is a deliberate lab switch, not a silent fallback.

---

## Why a user needs this

A capable trader already has a terminal (HL, Axiom-trade, etc.). They do **not** need another chart. They need:

1. **A private book** that does not sit in OpenAI logs (Knole’s job, but Knole does not fill).
2. **A policy fence** they can explain to themselves at 2am (max clip, daily loss, leverage 1x, ETH/BTC only).
3. **A desk that hunts** while they sleep **without** spending (Watch → notify → they authorize).
4. **A memory that is statistics**, not “the model learned.”
5. **A receipt** a judge or a partner can verify on Aristotle without seeing the book.

PIT is daily software because Watch + calibration make **re-entry** valuable. Ask-once is a demo. **Welcome back. Five opportunities match your rules** is the product.

---

## Why 0G is load-bearing

| Layer | What the user sees | What is true |
|---|---|---|
| Compute | “Sealing strategy + private book to 0G Compute” → “TEE response verified” | Direct TeeML HPKE `VerifyE2EE`. Signer = on-chain `teeSigner`. Not Router. Not “hardware verified” unless that exact quote exists. |
| Storage | “Encrypted memory anchored on 0G Storage” | Official AES-256-CTR / ECIES + Go `download --proof`. TS SDK ignores proof — **forbidden for proof claims**. |
| Aristotle | “Decision committed on 0G Aristotle” | Policy pin, forecast commit, receipt, 8004, 7857. |
| Calibration | “Outcome recorded” | Brier/ECE on-chain as **hashes + scores**, never positions. |
| Verify | “Signer matches on-chain TeeML identity” | `/verify` recomputes from explorer, not a screenshot. |

**Honesty (product copy):**

- Do **not** claim DA. Epoch 0, no mainnet DAEntrance.
- Do **not** claim iTransfer on **mainnet**. Attestor `chain_id: 16602` Galileo. Re-verified M18 (Chrome config + `eth_getCode` length **0** on 16661, **295** on 16602). Galileo official iTransfer path exists; PIT has **not** executed a transfer tx — UI shows ⚠️ not ✅.
- Do **not** claim a 0G CLOB. JAINE ~$2.9k TVL. Fill is Hyperliquid.
- Do **not** say “99.9% secure.”
- Do **not** fake progress. If TeeML takes 15s, the ring is in `researching` for 15s.

---

## M17 — strongest truthful Agentic ID + sealed committee (gate)

Not product UI. Gate proofs live under gitignored `PIT/_gate/evidence/`.

### Agentic ID lifecycle (ship this, nothing stronger)

```
USER WALLET
  → MINT PitDeskID (encrypted intelligence hash + 0g:// URI)
  → VERIFY ownerOf
  → AUTHORIZE USAGE (off-chain allowlist; official IERC7857Authorize)
  → USE in private trading workflow (host checks isAuthorized)
  → REVOKE / RENEW
  → RECORD RECEIPT
  → PERFORMANCE / CALIBRATION
  → VERIFY
```

If Foundation attestor ever appears on 16661 (config `chain_id: 16661` **and** `eth_getCode(AgenticID)` nonzero **and** a real `iTransferFrom` that changes `ownerOf`): **STOP**, re-verify, then add:

```
TRANSFER → NEW OWNER → NEW AUTHORIZATION → INTELLIGENCE ACCESS CONTROL
```

Until then: `iTransferFrom` / `iCloneFrom` revert `AttestorNotOnAristotle`. ERC-721 `transferFrom` reverts `ERC7857UseITransferFrom`. `verifier() = address(0)`.

**Rejected (not honest ERC-7857 transfer):**
- Self-deployed always-true TEE oracle (Talos anti-pattern).
- Knole-style owner-move without Foundation attestor / custom interface id `0x4b396f04`.
- 4lpha omitting iTransfer from the prod ABI while implying transfer.
- A “transferable wrapper” labeled as native iTransfer.
- Two-network “Aristotle iTransfer” via Galileo. Official Foundation path on Galileo is **seal-bound runtime** (openclaw / hermes / prime-agent): `DEPLOYMENT.md` §2.4 **re-enables ERC-721 transferFrom and makes iTransferFrom revert** for sealed agents. That is not PIT desk transfer on 16661.

**authorizeUsage is the real delegated lifecycle today.** Official interface: the allowlist does **not** gate on-chain transfer; it is what off-chain inference/use checks.

Gate contract `PitDeskID` `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` (Aristotle). Production IDs `0x2afbede9` / `0xdf597d99` / `0x74f8628b` / `0x80ac58cd`. Foundry 7/7. Live mint tokenId 1 to fixture `0xBDfC…0034`.

### Sealed committee (replace S1 PONG)

```
PUBLIC MARKET (HL l2 + mark/oracle/funding/OI)
  + PRIVATE USER BOOK / POLICY  → ENCRYPT (HPKE)
  → Direct TeeML (not Router)
  → RESPONSE
  → VerifyE2EE (scheme zg-sig-v1/e2ee-ct)
  → teeSigner == on-chain identity
```

Repeat for Researcher, Challenger, Risk. Fail closed if verification fails.

**Authoritative size is the Go/host engine**, never the LLM. Engine denies: `challenger_killed`, `no_side`, `below_min_notional` ($10), `leverage_above_policy`, `kill_switch`, `coin_not_allowed`.

**Independence (honest):** Router Private TeeML = **5** (Chrome 2026-08-26: GLM-5.2, 0GM-1.0-35B-A3B, 0GM-1.0-35B-A3B-SIA, whisper, z-image). Text committee can use three SKUs. Gate fixture has a Ledger subaccount only on glm-5.2 provider `0x7DCF…87D` / teeSigner `0xA46E…46B9`. M17 roles = **prompt-envelope independence, same provider**. Do not claim three independent models until each role’s provider is funded and acknowledged.

**Scenario facts (live ETH book mark 2438.5, funding +1.11319e-05, sha256 `05de1157…`):**
- S1 $15/1x: Researcher proposed **sell**; Challenger **killed** (mark-below-oracle is a short disadvantage); engine `eligible=false` / `challenger_killed`. All three VerifyE2EE **OK**.
- S2 forced 50x long: Researcher `proposed_side=none`; engine `no_side`.
- S3 $5 clip: Researcher **sell**, Challenger **survives**, Risk medium; engine `below_min_notional`.

A sealed `$15` run that reaches `eligible=true` was **not observed** on this tape. The engine path **is** proven: committee-agree + clip ≥ $10 + 1x → eligible. Do not fake S1-eligible.

---

## What we added (M16) and what we killed

### Added (survived adversarial eval)

**A. Opportunity Watch** — PIT continuously (or on a user-set poll, default 5 min while desktop is open; OS notification if backgrounded) scans the **allowlisted** universe under **pinned policy**. Pipeline is the same as ask. Output is a card: VIEW THESIS / PASS / AUTHORIZE. **No auto-place.** This is the frequency engine. Without it PIT is a chatbot with extra steps.

**B. Strategy Health** — Home and `/health` show which skills are overconfident, which regimes are toxic, calibration drift. This is the existing calibration store **surfaced**, not a new ML trainer. Users will not “trust the agent” unless they can see it rotting.

### Killed / deferred

| Idea | Verdict | Why |
|---|---|---|
| Counterfactual Lab (“what if we passed / half size”) | **DEFER v1.1** | Easy to fake; needs a filled history; Watch PASSES already journal as first-class `would-have`. Half-size is a **policy** (`max clip`), not a simulator. Ship after N≥30 resolved clips or it is theater. |
| Overnight autonomous trading | **KILL v1** | Judge + user: agent must not bypass AUTHORIZE. |
| Smart multi-venue router | **KILL** | Perp vs AMM vs bridge are different instruments. 4lpha collision. |
| Telegram / copy trade / social board / yield farm / meme universe | **KILL** | Leak or toy. |
| Fake DA / storage-as-DA / fake hardware TEE / fake 7857 transfer | **KILL** | Judge pattern. |
| AI probability as truth / AI calldata | **KILL** | Engine + allowlist. |
| Huge token universe | **KILL** | ETH/BTC (+ optional SOL) until calibration exists. |
| Generic chatbot | **KILL** | Ask is a **desk command**, not a companion. |

---

## Multi-user architecture

PIT is **not** a single-wallet demo. Development keys in `PIT/.env` are **fixtures**. Production onboarding is the same code for user 1 and user 1000.

```
USER A                         USER B
wallet A                       wallet B
HL account/sub A               HL account/sub B
session A (keychain)           session B (keychain)
policy A                       policy B
skills A                       skills B
memory A                       memory B
forecasts A                    forecasts B
calibration A                  calibration B
```

No cross-user reads/writes. No global trading operator. No shared session key. Hosted **web** never holds an HL session key.

### Identity objects

| Object | Definition | Isolation |
|---|---|---|
| **User** | Human; proven by EVM signature (SIWE). Never a private key in PIT. | — |
| **Workspace** | UUID. One active workspace per user in v1 (future: multiple desks). | All rows keyed `workspaceId` |
| **Wallet connection** | Address + chain 16661 for 0G txs; same address typically is HL master | Public to that user |
| **HL master** | Address that `approveAgent`s. Proven by that signature, not by `.env` | Per workspace |
| **HL session** | Named agent, `validUntil`, order/cancel only | Key in OS keychain `pit/{ws}/session/{id}` |
| **Storage encryption key** | Per workspace; `0x` + 32 bytes for official client | Keychain `pit/{ws}/memory` |
| **Memory namespace** | `ws/{workspaceId}/…` on 0G Storage | Decrypt fail across ws |
| **Agent identity** | ERC-8004 token owned by user wallet | On-chain owner |
| **Policy namespace** | Pinned hash + local JSON | Private |
| **Forecast / receipt namespace** | DB + chain keyed by ws | IDOR tests |

**Web** may: connect wallet, pin policy (user-signed txs), view public receipts, `/verify`. **Web may not:** create HL session keys, sign HL orders, export memory keys.

**Desktop + CLI** may authorize. **MCP is read-only.**

---

## User account model (every object)

Every persisted object has: `workspaceId`, `createdAt`, `version`, `provenance` (who/what wrote it), `authz` (who may mutate), `visibility` (`private` | `hash-public` | `public-proof`).

| Object | Visibility | Auth to mutate |
|---|---|---|
| User / Workspace | private | SIWE |
| Wallet connection | address is visible to user | reconnect |
| Trading account (HL) | private positions | user + session |
| Session | private; agent address may be shown | user create/revoke |
| Policy | private; **hash** on-chain | user pin |
| Skill definition | private | user pin |
| Skill stats | private; **Brier hash** on 8004 | engine |
| Opportunity | private | watch |
| Forecast | private; **commit hash** on-chain | engine |
| Preview | private; binds order | engine |
| Execution | private; HL oid public-enough on HL | session + AUTHORIZE |
| Receipt | **public proof** (hashes, not book) | contracts |
| Outcome | private numbers; public scores | resolver |
| Memory blob | private ciphertext | storage + memory key |
| Calibration card | private UI; public 8004 | engine |

---

## Wallet connection (never keys)

**Copy in product:** “PIT never asks for your private key.”

EVM (Rabby / MetaMask / WalletConnect):

1. Connect.
2. Sign SIWE / workspace bind (EIP-191).
3. Sign **only** the on-chain txs the user initiates (policy pin, 8004 register, 7857 mint, Ledger top-up, `approveAgent` to HL).

Hyperliquid:

1. User identifies the HL account (usually the same address).
2. PIT generates a **local** session key in OS keychain (Tauri `stronghold` / Windows Credential Manager / macOS Keychain).
3. User signs **`approveAgent`** in their wallet (extraAgents).
4. Session expires. Cannot withdraw. Cannot be exported to browser JS.

```
WALLET  ≠  PIT SESSION  ≠  PIT OPERATOR
```

There is **no operator**. Fixture payer `0xBDfC…0034` is a **test wallet**, not architecture.

**Forbidden storage for session keys:** `localStorage`, Vercel, Render, plaintext DB, `.env`, frontend bundle, MCP.

---

## Hyperliquid onboarding (non-expert)

Users are not assumed to know HL.

| Step | Screen | Truth |
|---|---|---|
| 1 | Connect wallet | SIWE |
| 2 | Testnet / Mainnet | Explicit. Testnet faucet needs a prior mainnet deposit (official). |
| 3 | Identify HL account | `clearinghouseState` + **spot** USDC (unified). Do not tell the user they have $0 when spot is funded. |
| 4 | Capital source | Spot / perp unified. Optional isolated **subaccount**. |
| 5 | Subaccount | Create/select if they want isolation. User-signed. |
| 6 | Max PIT trading amount | Policy `maxClipUsd` + optional subaccount cap. |
| 7 | Create session | Keygen in keychain. |
| 8 | Approve scoped agent | Wallet `approveAgent`, name `PIT-{shortWs}`, `validUntil` ≤ 1h. |
| 9 | Permissions card | See below. |
| 10 | Expiry countdown | Real `validUntil`. |
| 11 | First action | “Tiny test order” **or** “Research-only.” Never skip to size. |

**PIT CAN:** place orders, cancel orders, read market data, read permitted trading state.

**PIT CANNOT:** withdraw, send funds, change leverage (PIT refuses; **protocol still allows** `updateLeverage` if the agent key is stolen — disclose this), transfer subaccounts, approve another agent, change policy.

HL protocol hole is **documented in onboarding**, not hidden. Mitigation: short TTL, keychain, kill switch, never export key.

---

## Policy model (human → executable)

Policies are cards a normal user can set. Every row maps to a **host check** (and a pin hash). LLM cannot override.

| User card | Code check | Fail behavior |
|---|---|---|
| Max trade | Host sizer USD notional | Refuse preview |
| Daily loss halt | Local realized PnL from fills | Kill signing |
| Max leverage | Preview `leverage` field; never send `updateLeverage` | Refuse |
| Allowed assets | Coin allowlist | Refuse |
| Allowed market types | `perp` / `spot` / `none` | Refuse |
| Allowed venues | `hyperliquid` / `0g-univ3-labeled` / `none` | Refuse |
| Min skill calibration | Stats store vs pin | Refuse |
| Cooldown | Clock since last fill | Refuse |
| Session expiry | `validUntil` | Refuse |
| Emergency kill | Flag + revokeAgent | Block all signs |
| Never trade when X | Named predicates (news halt = **not v1**; funding extreme = skill gate) | Refuse |
| Min liquidity | `impactPxs` / book depth vs size | Refuse |
| Max uncertainty | Engine `uncertainty` vs pin | Refuse |

**Re-check policy after every async step** (committee return, network delay). Preview hash binds the exact order. User authorizes **that** hash.

---

## Opportunity Watch

Not “click Generate Forecast.”

**Home (returning user):**

> Good evening.  
> 3 opportunities match your policy.

Each card: market, skill tag, confidence, **calibration of that skill**, risk, expected edge (bps), liquidity, invalidation one-liner. Actions: **VIEW THESIS** / **PASS** / **AUTHORIZE**.

Opening a card is a **real state machine** (same as ask):

```
LIVE DATA → SEALED RESEARCH → CHALLENGER → RISK → ENGINE → POLICY → PREVIEW
```

Watch **candidates** are cheap (public info + deterministic prefilter). **Sealed committee** runs only when the user opens VIEW or when Watch is set to “prepare next” (optional, still no sign).

Polling: desktop loop while unlocked. No cloud worker holding session keys. If the app is closed, Watch is closed. Optional OS notification: “PIT: ETH funding-basis eligible — open to review.”

---

## Skills + memory (no magical learning)

Memory kinds (tagged, queryable, never mixed):

`observation` · `forecast` · `execution` · `outcome` · `error` · `calibration` · `role_performance` · `skill_performance` · `risk_event` · `policy_decision`

The LLM does not “get smarter.” The **stats store** updates Brier, ECE, win/loss by regime, challenger catch rate. Weights are explicit numbers on the Strategy Health card.

---

## Personalization

Two users never share book, policies, strategies, session keys, memory, skills, calibration, forecasts, private evidence. ERC-8004 public page shows **hashes and scores**, not the book.

---

## Execution pipeline (authoritative)

```
MARKET DATA
  → CANDIDATE (Watch or Ask)
  → SEALED RESEARCH (0G TeeML)
  → CHALLENGER
  → RISK
  → DETERMINISTIC ENGINE
  → FORECAST
  → POLICY CHECK
  → PREVIEW HASH
  → USER AUTHORIZATION
  → SESSION SIGN (order|cancel only)
  → RECEIPT (0G)
  → STORAGE MEMORY
  → OUTCOME
  → CALIBRATION
```

**Invariants**

- LLM cannot set authoritative size, cannot change policy, cannot emit calldata.
- Session cannot withdraw, transfer capital, or change leverage **via PIT**.
- Session expires; kill switch blocks signing.
- `clientOrderId` on every action; durable exactly-once (not in-memory).
- Preview binds the final order; policy re-check after async.
- Wrong workspace cannot read receipt or memory.

---

## Durable exactly-once (upgrade from in-memory gate)

Ledger: local encrypted SQLite (desktop/CLI) keyed by workspace. Semantics from 4lpha `usedActionHashes` + NeoSoul `insertIfAbsent` — **PIT-native**, not a copy of their contracts.

| Case | Expected |
|---|---|
| Same action twice | Second `already-applied`, no second HL post |
| Concurrent same cloid | One winner; loser `conflict` |
| Restart between preview and sign | Preview still valid if TTL/policy hold; else regenerate |
| Crash after sign, before receipt | Recovery: query HL by cloid; file receipt; no duplicate |
| Crash after receipt | Idempotent file |
| Network timeout | Retry **read** status; do not blindly re-post |
| Duplicate click | UI disabled + ledger |
| Browser refresh | Web cannot sign; desktop recovers from ledger |

---

## Surfaces (one core)

| Surface | Can authorize | Session key | Notes |
|---|---|---|---|
| **Desktop** | Yes | OS keychain | Flagship. KEPT-grade UI. |
| **CLI** | Yes | keychain | `pit init/login/policy/ask/opportunities/forecast/preview/authorize/cancel/status/resolve/card/verify/kill` |
| **Web** | No HL sign | None | Wallet connect, research view, cards, policy **txs**, `/verify` |
| **MCP** | **No** | None | Read: markets, opportunities, forecast, status, card, verify |
| **SDK** | No secrets in browser builds | — | Typed API; no duplicated executor |

MCP **must not** authorize, place, transfer, export sessions.

---

## Frontend — KEPT visual system, PIT grammar

**Do not copy KEPT business logic** (savings circles, XRPL, FXRP). **Do copy** visual and interaction philosophy from `D:\route\Flare\kept` and https://kept-ruby-five.vercel.app/

### Tokens (adapt; do not invent a third palette)

From `apps/web/src/styles/guide.css` + `app.css`:

| Token | Value | PIT use |
|---|---|---|
| Coral | `#D82F2F` | Primary CTA, landing hero field, accent (replaces jade in app shell) |
| Ink | `#1A1A1A` / `#0A0A0A` | Type on coral; marketing |
| Cream | `#F0E7D4` / paper `#F9F5E3` | Type on dark; postcard fill |
| App ink | `#08090B` surfaces `#0E1013`… | Logged-in terminal |
| Type | Host Grotesk Variable (marketing); Geist Variable (app/mono) | Oversized editorial heads |
| Radius | Pills `999px` for CTA; marketing hairlines; app `20/14/8` | KEPT split: landing = Swiss guide; app = dark terminal |
| Motion | `cubic-bezier(0.16, 1, 0.3, 1)` | Real state, `prefers-reduced-motion` |
| Grain | SVG turbulence overlay on coral | Landing only |

**Landing:** coral sticky hero, mega wordmark **PIT.**, kicker, pill nav (How it works / Help), black **Get started**. Sunburst SVG becomes the **desk ring** (not a savings circle).

**Sign-in:** split — dark functional left / coral diagram right (WalletGate pattern). Buttons: Rabby/MetaMask recommended; WalletConnect; **no** seed field. Named errors (port KEPT `namedStates.ts` to PIT ids: `WRONG_NETWORK`, `SIGNATURE_DECLINED`, `SESSION_EXPIRED`, `HL_UNFUNDED`, `POLICY_BLOCK`, `TEE_VERIFY_FAIL`, …).

**Get started:** AppShell + two cards (Research-only vs Connect trading) — StartFlow pattern, not a dashboard grid.

### PIT ring (custom SVG — not KEPT seats)

Ten nodes, one lit = current **real** phase:

`Market → Seal → Research → Challenge → Risk → Policy → Authorize → Execute → Prove → Learn`

Center copy is live: `WEEK` is wrong. Center = `ETH · 31 bps` or `SEALING` or `oid 52698…`. If a step has not started, it is dim. No CSS animation pretending TeeML finished.

**No fake charts. No placeholder cards. No “99.9% secure.”**

---

## First-time journey (14 beats)

Copy uses **YOU / YOUR WALLET / YOUR TRADING ACCOUNT / YOUR SESSION / YOUR POLICY / YOUR MONEY**. Never “owner,” “master,” or “operator” for a normal user.

**Landing**

> PIT  
> Private Alpha OS.  
> Private research. Controlled execution. A desk that learns.  
> **Get started**

Then:

1. **Connect YOUR WALLET** — “Your keys stay yours. PIT never asks for a seed phrase.”
2. **Select Testnet or Mainnet** — show the capability list for that world (see Two environments). Confirm: Testnet money is not Mainnet money.
3. **Connect YOUR TRADING ACCOUNT** — Hyperliquid. “PIT never gets withdrawal authority.”
4. **Confirm the account** — show **YOUR** spot USDC (unified), not a global PIT balance.
5. **Set YOUR POLICY** — cards: max trade, daily loss, leverage, assets, venues, calibration floor.
6. **Create YOUR SESSION** — scoped, ~1h, `order` + `cancel` only. Show the exact permission list before the signature.
7. **First private analysis** — e.g. “Find the lowest-risk ETH opportunity today.”
8. **Live 0G pipeline** (every state is a real backend event — no fake animation):

```
YOUR PRIVATE BOOK
  → SEALING (HPKE)
  → 0G TEE (Direct TeeML)
  → TEE SIGNATURE
  → ON-CHAIN SIGNER (teeSigner match)
  → STORAGE
  → RECEIPT
  → CALIBRATION
```

9. **Opportunity discovery** (Watch) — only cards that pass YOUR POLICY.
10. **Exact preview** — size from the engine, not the model.
11. **YOU AUTHORIZE** — the only step that can spend YOUR MONEY.
12. **Execute** via YOUR SESSION (or research-only skip).
13. **Verify** — explorer + `/verify` + Storage `--proof`.
14. **Learn** — YOUR calibration card (empty until enough samples; never fake 72%).

**No private keys. No seed phrase. No raw API secret. No Router `sk-` in the client.**

---

## Re-entry (the daily product)

Do **not** restart setup.

> Welcome back.  
> Session expired → create a new session.  
> Wallet, policies, memory, skills, performance, agent identity remain.  
> **5 new opportunities match your rules.**

---

## Demo spine (every arrow = backend event)

Wallet connected → **network chosen** → HL connected → Policy pinned → Session created → Live snapshot → 0G seal → Researcher → Challenger → Risk → Engine → Forecast → Policy pass → Preview → AUTHORIZE → Real HL order → Order status → 0G receipt → Storage proof → Outcome → Calibration card.

Judge path: Chrome + explorer + `/verify` + `go test` + Foundry.

---

## Real execution venues

| Venue | Role |
|---|---|
| Hyperliquid perps | Default book. Testnet + mainnet dust **FACT** (M15). |
| 0G univ3 (Oku/Zia) | **Labeled dust only**, `forecast.market = 0g:univ3:*`. Not a router. |
| `none` | Research-only. |

Min HL order **$10**. Tick/5-sigfig rounding required.

---

## Contracts (to implement — not yet issued)

`PitPolicy` · `PitForecasts` (commit-reveal + resolve) · `PitReceipts` · `PitMemory` (root registry) · official 8004 (already live).

`PitDeskID` **gate-issued** on 16661: `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` (unverified source on ChainScan). Phase 19 may redeploy from `pit/contracts` with the same semantics; do not treat the gate address as the only forever product address until verify+CREATE2 is decided.

Replay: **off-chain durable ledger** is source of truth for HL cloids; chain is proof, not the HL mutex (HL is not 0G). Optional on-chain `usedPreviewHash` to prevent double-file.

---

## Security red team (must be tests, not prose)

User A cannot: read B memory/book/policies; execute B forecast; cancel B order; use B session; replay B receipt; decrypt B storage.

Attacks: bearer mixup, tenant injection, object-id guessing, storage path collision, receipt replay, session replay, websocket leak, cached private data, refresh, logout/login confusion.

HL: allowlist `order`/`cancel` only; sizer; preview bind; kill; expiry.

0G: no Router on book; `VerifyE2EE` required; Storage `--proof`; no DA claim.

---

## Competitor red team (M16)

Judge pattern (03 file, notmartin / @Damiclone): real txs, real 0G primitives, mainnet, browser attestations, reproducible proof, tests, CI, Solidity tests, security, in-window git. Penalize: mocks, fallbacks, env-gating, deploy-only, docs-only, fake 7857/DA, external fallback as 0G, broken concurrency, missing tests, pre-wave work.

| Dimension | PIT now (FACT) | Competitor best | Gap | Improve? | Path | Test | Kill |
|---|---|---|---|---|---|---|---|
| Real txs | HL testnet oid 58551048563 + mainnet 526980359929 + cancel; 8004 register | Zia $1.91M; 4lpha vault | Density vs Zia | Yes, user flows not more size | Dust + Watch demo | Chrome+API | Fake fills |
| 0G primitives | Direct TeeML VerifyE2EE; Storage --proof; official 8004; **M17 3-role committee over live ETH book** | Knole E2EE; Hanami browser attest; Lumen honesty | Product not wired; envelope-only independence | Yes | promote `_gate/committee` | M17 evidence | Router fallback |
| Mainnet | Aristotle + HL mainnet dust | 4lpha mainnet vault | 4lpha is a vault we refuse | No (wrong product) | Stay desk | — | Clone PolicyVault |
| Attestation | VerifyE2EE + /verify planned | Hanami recompute | UI not built | Yes | KEPT verify page | Playwright | Screenshot-only |
| Tests/CI | `_gate` tests; product CI **not** | 4lpha/Knole CI | Missing product CI | **Must** | Foundry+go+Playwright | CI green | Ship without |
| Solidity tests | Gate `PitDeskID` 7/7; product suite **not** | 4lpha usedActionHashes tests | Product contracts | **Must** | forge in `contracts/` | fuzz replay | Deploy-only |
| Security | Allowlist FACT; in-memory cloid | NeoSoul insertIfAbsent | Durable ledger | **Must** | SQLite+tests | crash cases | In-memory prod |
| Concurrency | Gate sequential | 4lpha mutex | Multi-user | **Must** | ws isolation tests | A/B E2E | Shared session |
| Mocks | `ALLOW_FALLBACKS=false` | Axiom MockUSDC | Don’t regress | Guard | startup exit | — | Env-gate |
| 7857 | Production IDs **live**; mint/own/authorize/revoke **live**; iTransfer IMPOSSIBLE | Lumen OracleNotLive | Honest | Do not copy Knole owner-move | Foundry 7/7 + live txs | ChainScan mint | Fake iTransfer |
| DA | Unused FACT | — | Ceiling 0 | No | — | — | Storage-as-DA |
| Product UX | None | KEPT editorial; NeoTrade desktop | **Largest product gap** | **Must** | KEPT tokens | Chrome E2E | Tailwind dashboard |
| Multi-user | Fixture-shaped env | NeoSoul tenants | **Architecture gap** | **Must** | workspace | IDOR tests | Global operator |
| Opportunity loop | Ask-only in docs | Terminals (Axiom-trade) hunt | Frequency | Watch | Poll+policy | No ghost cards | Fake opps |
| Privacy | Sealed path FACT | Knole | Committee on book | Wire | — | Log leak tests | Public book |
| Dual env | M18 Galileo ≠ Aristotle catalogs | Lumen 16661/16602 honesty | Galileo Omni unproven | Yes | separate DeskID | Chrome both PCs | Copy mainnet SKUs to testnet |

**NeoSoul:** take mechanics (structured forecast, gateway, skills). Refuse public agent economy, mnemonic trading key, closed-core-only proof.

**4lpha:** take `usedActionHashes` idea. Refuse sandbox vault + mock USDC as the product.

**Axiom:** take nothing from MockUSDC / software TEE PK.

**Knole:** take privacy journal tone. Add fill + Brier.

**Hanami:** take /verify UX. **Refuse** `OG_DIRECT_ENABLED=false` + Router fallback (`og-compute-direct.ts` falls back on any Direct miss). PIT exits.

**Orchestra (Cannes clone, not a 0G wave winner):** NL portfolio agent, Groq fallback, Safe + Ledger, **AUTO_EXECUTE** under a USD limit. Confirms Watch is real demand. **Refuse** auto-place, Groq, Uniswap-as-primary-book.

**NEXUS / Lumen / Maneki / TRAIDE:** honesty, W0G proofs, don’t become a sticker on HL.

---

## Business model / moat

Desktop subscription later. v1 take rate 0. Moat = first honest sealed desk with Watch + cannot-withdraw session + private calibration, not a 4lpha rename.

---

## Score

**Proven (gates M14–M18):** **92 / 100**. Hard ceilings: DA=0, iTransfer on 16661, fill≠0G CLOB. M18 locked dual-env honesty (Galileo ≠ Aristotle catalogs) — **do not add points**.

**Planned (this spec, if implemented without mocks):** **93**. The +1 is Frequency/UX (Watch + KEPT onboarding + re-entry). Do **not** score 95. Do not double-count TEE/Storage already at cap.

35-dimension sheet: keep M15 numbers except **Frequency 3→3** (already cap) — the extra point is taken from **UX 3→4** on the 4-point UX row (was 3). If UX was already 3/4, planned UX = 4/4 **after Chrome E2E**, not before. **Until UI exists, report 92.**

---

## Limitations (honest)

- HL agent can `updateLeverage` if key is extracted (protocol). Disclose + short TTL.
- Unified USDC lives in spot; UI must not show $0.
- Committee-over-book is **gate-proven** (M17); **not** product-wired. Independence is prompt-envelope / same glm-5.2 provider until other TeeML providers are funded.
- Sealed `$15` committee did **not** reach `eligible=true` on the live ETH tape (Challenger killed). Do not invent a passing thesis.
- 8004 owner self-feedback reverts — use separate reporter.
- Chrome product E2E **impossible** until UI exists.
- Llama 0G TVL slice timeout = UNVERIFIED (use Hyperbeat/DefiLlama).
- Product contracts **NOT ISSUED** except gate `PitDeskID` (see Contracts).
- Durable ledger **NOT IMPLEMENTED** (in-memory gate only).
- Galileo Direct TeeML chat exists (`qwen2.5-omni-7b` ack) but **PIT has not VerifyE2EE’d it**. Testnet committee is not a copy of the M17 glm-5.2 run.
- Galileo Router Private TeeML chat count = **0** (Chrome: Private filter = image only). Do not use testnet Router for the book.
- Fixture `PIT/.env` `PIT_DESK_ID_CONTRACT` is now pinned to the live Aristotle address (public, not a key).

---

## 0G Private Computer API key

**DEPENDS.** The private book uses **Direct / Advanced**, which authenticates with a **wallet-signed** ephemeral `app-sk-` (tokenId 255). That is **not** the dashboard key.

Dashboard keys (`sk-` inference, `mk-` management) are created at https://pc.0g.ai/dashboard/api-keys (mainnet) and the testnet twin. They authenticate the **Router**. Official docs: Router still handles prompts in memory; `verify_tee` is Router-asserted. **Forbidden** on YOUR PRIVATE BOOK.

PIT does **not** need the user to create or paste a Router API key to implement or run the confidential path.

---

## Implementation-ready?

**Planning gate (M16): READY to start `pit/` per IMPLEMENTATION_PLAN phases 0→1.**  
**Hardening (M17): Agentic ID lifecycle + sealed committee primitives proven in `_gate` on Aristotle. iTransfer on 16661 = NOT LIVE.**  
**Revalidation (M18): Dual Mainnet+Testnet architecture locked. Router API key not required for Direct. Galileo iTransfer = official path, PIT-unverified. Product UI still not started.**  
**Product UI/contracts: not shipped.**  
**Do not call the repo “complete.”**
