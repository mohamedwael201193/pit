<div align="center">

# PIT

**Private intelligence on 0G. A real Hyperliquid order on this computer.**

Product **0.9.13** · companion `127.0.0.1:17373`

[Website](https://pit0g.vercel.app) · [Health](https://pit-health.onrender.com/health) (`sign: false`) · [Windows installer](https://pit-health.onrender.com/windows) · [Launch film](https://youtu.be/zYgxDTI7jIk)

</div>

```
private book
    → 0G Direct TeeML glm-5.3
    → Researcher / Challenger / Risk
    → host VerifyE2EE
    → policy size
    → AUTHORIZE on this computer
    → Hyperliquid OID / fill
    → 0G Storage --proof
```

The website discovers and proves. The desktop protects and acts. The model cannot set size, pin policy, or AUTHORIZE. If Direct TeeML is gone, PIT **stops**. It does not fall back to Router inference.

---

## Recorded Mainnet fill

Job `4a1d45ec-8c3f-4883-a162-19739accb9cf` · agent `PIT-4bbee556` · FILLED only from Hyperliquid `user_fills`.

| | |
|---|---|
| OID | **531667200134** buy 0.16 HYPE (docs @ 80.909 · live `px` 80.826) |
| Research | [`0x1d2113bd…510b`](https://chainscan.0g.ai/tx/0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b) → Flow `0x62D4…7526` |
| Order evidence | [`0x8c28051b…11fb`](https://chainscan.0g.ai/tx/0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb) |
| Storage | [submission 211566](https://storagescan.0g.ai/submission/211566) |
| Wallet | `0xbdfcee82bd42fefa58ee850b3709636a8b6b0034` |
| Agent | `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` |

Historical ETH OID `529167222216` is older fill on the same account. Not flattened. Not reminted.

---

## Architecture

```mermaid
flowchart LR
  subgraph WEB["pit0g.vercel.app"]
    Radar[Radar / Proof / Pair]
  end
  subgraph HOST["PIT Desktop · 127.0.0.1:17373"]
    Policy[Policy pin]
    Keyring[OS keychain]
    Sealer[pit-sealer]
    TradeNow["TRADE NOW"]
  end
  subgraph OG["Aristotle 16661"]
    Direct["Direct TeeML glm-5.3"]
    Storage["Storage --proof"]
    Identity["Desk ID / ERC-8004 read"]
  end
  subgraph HL["Hyperliquid"]
    Order["order / cancel"]
    Fills["user_fills → FILLED"]
  end
  Radar -->|"8-char pair / 2 min"| HOST
  Policy --> Sealer
  Keyring -->|"app-sk-"| Direct
  Sealer -->|"HPKE + VerifyE2EE"| Direct
  TradeNow -->|"AUTHORIZE + preview hash"| Order
  Order --> Fills
  Sealer --> Storage
  HOST --> Identity
```

---

## 0G live integration

Judged from live `getService`, provider `/v1/models`, `/v1/e2ee/pubkey`, Aristotle RPC, and product code. **Not 100%.** Temporary provider Unhealthy during 0G recovery is not a SKU failure.

| Integration | Status | Proof | Blocker |
|---|---|---|---|
| Direct TeeML glm-5.3 | **LIVE** | `getService` model `glm-5.3`, verifiability `TeeML`, teeAck `true`, URL `https://compute-network-19.integratenetwork.work`, provider `0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D`. `/v1/models` lists only glm-5.3 TeeML. `getAllServices` has **1** glm-5.3 TeeML row on this provider. | Router may list glm-5.3 as TeeML/`is_healthy:false` for this provider, or TeeTLS on other providers. **Router is refused as inference.** |
| HPKE + VerifyE2EE | **LIVE** | Native `pit-sealer`. Scheme `zg-sig-v1/e2ee-ct`. Pubkey signer **equals** on-chain teeSigner `0x041a09E5bEF30fd776D66Bb892d18B97637C7C7c`. Live glm-5.3 PONG `ZG-Res-Key` `a797da3f-bcd3-4a3f-ba4d-ada54e6952c0`. Historical glm-5.2 recovered `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`. | Temporary Unhealthy on the Router catalog is not TeeTLS. |
| Sequential committee | **LIVE** | Researcher → Challenger → Risk, same Direct provider, independent HPKE envelopes, each `VerifyE2EE` **OK** (keys `a797da3f…` / `b86930b3…` / `ea506ef1…`). | Not a new Hyperliquid trade. Product hunt with a private-book thesis was not re-run. Same-provider honesty: `prompt_envelope_only_SAME_PROVIDER`. |
| Storage `--proof` | **LIVE** | Official Go `0g-storage-client` only. Flow [`0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526`](https://chainscan.0g.ai/address/0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526). Recorded roots `0x9fd42770…f72e` / `0x8c94ec8e…ff66`. | TypeScript SDK is not used for proofs. |
| Aristotle 16661 | **LIVE** | RPC `https://evmrpc.0g.ai`. Serving [`0x47340d90…8d84`](https://chainscan.0g.ai/address/0x47340d900bdFec2BD393c626E12ea0656F938d84). Ledger [`0x2dE54c84…74E3`](https://chainscan.0g.ai/address/0x2dE54c845Cd948B72D2e32e39586fe89607074E3). `PIT_NETWORK` binds 0G + Hyperliquid together. | Galileo sealed ask stays `galileo_e2ee_unproven`. Mixed RPC refused. |
| PitDeskID / ERC-7857 | **PARTIAL** | [`0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35`](https://chainscan.0g.ai/address/0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35). `ownerOf(1)` = bound wallet. IDs `0x2afbede9` / `0xdf597d99` / `0x74f8628b` / `0x80ac58cd`. | Product **reads**. It does not mint/authorize as a trading gate. `iTransferFrom` / `iCloneFrom` revert `AttestorNotOnAristotle`. |
| ERC-8004 | **PARTIAL** | Identity [`0x8004A169…a432`](https://chainscan.0g.ai/address/0x8004A169FB4a3325136EB29fA0ceB6D2e539a432) agent **3489333**. Live `ownerOf` = bound wallet. Register [`0xa2f67529…`](https://chainscan.0g.ai/tx/0xa2f67529745a662163b84fe10f855a3aa25596f9bc4d4c604d2abefbc3f3ff7d). | Product reads `ownerOf`. It does **not** submit register or `giveFeedback`. Reputation [`0x8004BAa1…9b63`](https://chainscan.0g.ai/address/0x8004BAa17C55a88189AE136b182e5fdA19dE9b63) is pinned, not a product write path. |
| Skillbook / Experience | **LIVE local** | AES-GCM `skillbook.enc` / `experience.enc` on this computer. | **0g-memory is not used.** `PIT_MEMORY_KEY` is refused in product mode. |
| 0G Pay / x402 / DA / PitMemory / PitForecasts / PitReceipts | **NOT USED** | Empty contract env. No product callers. | Do not call Foundry-only contracts production. |

Direct pin (identity is provider + URL + TeeML + teeAck; live `teeSigner` is copied from `getService` because TEE keys rotate):

```
provider     0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D
url          https://compute-network-19.integratenetwork.work
model        glm-5.3
verifiability TeeML
teeSigner    0x041a09E5bEF30fd776D66Bb892d18B97637C7C7c
```

Locked Direct credit is `getAccount.balance − pendingRefund`. Reclaiming funds are not spendable. PIT will not start a three-role committee on a ledger that is mostly in reclaim.

---

## Sealed committee

```mermaid
sequenceDiagram
  autonumber
  participant Host as PIT host
  participant Sealer as pit-sealer
  participant Direct as Direct TeeML glm-5.3
  participant TEE as on-chain teeSigner
  Host->>Sealer: Researcher envelope
  Sealer->>Direct: HPKE seal POST /v1/proxy/chat/completions
  Direct-->>Sealer: sealed response
  Sealer->>TEE: VerifyE2EE recovered == teeSigner
  Host->>Sealer: Challenger plus researcher_thesis
  Sealer->>Direct: independent sealed envelope
  Sealer->>TEE: VerifyE2EE
  Host->>Sealer: Risk envelope
  Sealer->>Direct: independent sealed envelope
  Sealer->>TEE: VerifyE2EE
  Host->>Host: policy clip, host size, exact preview
  Note over Host: AUTHORIZE only on this computer
```

Router `GET /v1/providers?model=glm-5.3` may list TeeTLS providers as healthy. Those URLs are compile-time denied for the private book.

---

## Authority

```mermaid
flowchart TB
  subgraph CAN["Can AUTHORIZE"]
    Desk["Desktop TRADE NOW"]
    CLI["CLI TTY + AUTHORIZE + --i-understand"]
  end
  subgraph CANNOT["Cannot AUTHORIZE"]
    Web["Website / chat"]
    MCP["pit-mcp / pit-os"]
    Model["The model"]
    Health["Health service"]
  end
  Desk --> Preview["preview hash + cloid once"]
  CLI --> Preview
  Preview --> HL["Hyperliquid order / cancel"]
  Web -.->|denied| Preview
  MCP -.->|no method| Preview
  Model -.->|cannot size or pin| Preview
```

Session allowlist: `order`, `cancel`.

Denied on the session key: `withdraw3`, `usdSend`, `spotSend`, `usdClassTransfer`, `sendAsset`, `updateLeverage`, `approveAgent`, and the rest of the withdraw/transfer/leverage/TWAP surface.

Pairing: 8-character code, two-minute TTL. `POST /pair` from `https://pit0g.vercel.app`. Response `{sign:false, canSign:false, device}`. Protect stays locked until paired.

---

## Desk path

1. Download PIT Desktop. Health `GET /windows` 302s to GitHub asset `PIT_0.9.13_x64-setup.exe`.
2. Companion listens on `127.0.0.1:17373` only.
3. Pair the browser. Protect: wallet signs Direct `app-sk-` (24h) into the OS keychain `pit/{network}/{workspace}/direct`.
4. Connect Hyperliquid. Local agent `PIT-{workspace8}`. Master wallet signs `approveAgent`. Ready requires live `extraAgents`. Recorded: `PIT-4bbee556` reused, not reminted.
5. Pin policy on Security. Chat cannot pin. Host sizes. The model cannot raise `maxClip`.
6. Agent hunts live books (ETH, BTC, SOL, HYPE, DOGE, AVAX under the recorded policy).
7. Private book enters Direct TeeML. Sequential sealed committee. Host `VerifyE2EE`.
8. TRADE NOW only on `READY_ELIGIBLE` with an exact preview. You authorize here, or you walk away.
9. FILLED is painted only from `user_fills`. This-job Storage `--proof` files research and order.

Recorded preview (HYPE job `4a1d45ec`): hash `0xb273d0052fe389b5e5ad3aad4b176e1cc993b8d8e605716bab78c70f3814e401` · host-sized buy 0.16 HYPE at 80.909 · notional $12.95 · clip $13 · 1x.

---

## Evidence

```mermaid
flowchart LR
  subgraph JOB["This job only"]
    Sealed["Committee evidence · prompt stripped"]
    Client["0g-storage-client --proof"]
    Flow["Flow 0x62D4...7526"]
    Scan["storagescan.0g.ai"]
  end
  subgraph VENUE["Hyperliquid"]
    OID["OID"]
    Fill["user_fills"]
  end
  Sealed --> Client --> Flow --> Scan
  OID --> Fill
```

| What | Evidence |
|---|---|
| HYPE research · `4a1d45ec` | [0x1d2113bd…](https://chainscan.0g.ai/tx/0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b) · root `0x9fd42770545ecaacbfff12e3ef7a537b564e31c9ef5515b3a820fd276c22f72e` |
| HYPE order evidence | [0x8c28051b…](https://chainscan.0g.ai/tx/0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb) · root `0x8c94ec8e643c90fe69276ff20f50a0bc3121f007d611e10e6ab9f24d26f2ff66` |
| StorageScan | [211566](https://storagescan.0g.ai/submission/211566) |
| Preview hash | `0xb273d0052fe389b5e5ad3aad4b176e1cc993b8d8e605716bab78c70f3814e401` |
| PitDeskID deploy / mint | [`0x2d5b688b…`](https://chainscan.0g.ai/tx/0x2d5b688bf09bb72cb44b092da0c27cbe87a623141872e62b84cb95ecf7e90c24) / [`0x9494e3fa…`](https://chainscan.0g.ai/tx/0x9494e3faec6d950942d1bfec53c4a13e6f28378da8e01ecd24b3e3de62e5c7d0) |
| ERC-8004 register | [`0xa2f67529…`](https://chainscan.0g.ai/tx/0xa2f67529745a662163b84fe10f855a3aa25596f9bc4d4c604d2abefbc3f3ff7d) |

Public `/proof` lists recorded this-desk filings. Each job has its own root. Historical receipts are not reused as a later trade's evidence.

---

## Surfaces

| Surface | Does | Does not |
|---|---|---|
| [pit0g.vercel.app](https://pit0g.vercel.app) | Radar, proof, pair, Protect, download | Hold a session key, AUTHORIZE, pin, arm |
| PIT Desktop | Policy, session, research, TRADE NOW, Activity, Portfolio, Sleep Mission | Withdraw or transfer |
| Health | `/health`, `/watch`, `/release`, `/windows` | Sign or trade |
| `pit-os` / `pit-mcp` | Read health, watch, loopback status | Sign, order, arm, export |

**Web** — Vite + Privy. Routes: `/`, `/radar`, `/capital`, `/autonomy`, `/missions`, `/proof`, `/agent`, `/how-it-works`, `/download`, `/pair`, `/protect`. Redirects: `/watch` → `/radar`, `/signin` → `/pair`, `/verify` → `/proof`, `/app/start` → `/pair`.

**Desktop** — Tauri `os.pit.desktop`. Sidecars `pit` + `pit-sealer`. Loopback `/local/*`. Denied: `/authorize`, `/export`, `/session` as companion aliases.

**CLI**

```powershell
cd pit
go run ./cmd/pit version
go run ./cmd/pit doctor
go run ./cmd/pit status
```

`pit version` prints `PIT 0.9.13`. `pit companion` binds `127.0.0.1:17373` only.

**SDK** — npm `pit-os@0.9.11`. `canSign` is always `false`. No authorize method.

**MCP** — npm `pit-mcp@0.9.12`. Tools: health, watch, release, companion, status, activity, positions, research_status, proofs, security. Not tools: authorize, order, cancel, export_session, arm.

---

## Memory and Sleep Mission

Skillbook (`skillbook.enc`) and Experience (`experience.enc`) are AES-sealed on this computer. `PublicSkillbook` stays **NOT ENOUGH DATA** until five verified rows. Forget is desktop Security. Two wallets never share a workspace.

Sleep Mission arms only on this host (`ARM SLEEP MISSION` / `pit mission arm` TTY). Chat, web, MCP, and the model cannot arm it.

---

## Security

- The model cannot AUTHORIZE, size, or raise clip.
- The browser never receives the session key or Direct token.
- Preview hash binding. Exactly-once `cloid`.
- FILLED only from Hyperliquid fill-ish status.
- Private book has no Router fallback and no TeeTLS SKU.
- Sealer VerifyE2EE runs on the host against on-chain `teeSigner`.
- `PIT_ALLOW_FALLBACKS` must be `false` or the process exits.
- No global memory key in `PIT_PRODUCT_MODE`.
- Copy `.env.example` to `.env`. **Do not commit `.env`.**

---

## Tests and CI

| Suite | Command | Last local |
|---|---|---|
| Go `pit` | `go test ./... -count=1` | **659 PASS / 0 FAIL / 4 SKIP** |
| Sealer | `go test ./...` | **7 PASS** |
| Foundry | `forge test -vvv` | **26 PASS** |
| Playwright | `npx playwright test` | **31 PASS** |
| pit-os / pit-mcp | `npm test` | **2 / 3 PASS** |
| Desktop e2e | `npx tsx e2e/run.ts` | **PASS** |

CI: `.github/workflows/ci.yml` (go, sealer, contracts, web, desktop, sdk, mcp). Release: Windows NSIS + `SHA256SUMS.txt`.

Installer SHA256 `B905B9ED167513757D4947BDE61103EB10ECD4A5F76554FE369F205DF3850B1E`. Tag [`v0.9.13`](https://github.com/mohamedwael201193/pit/releases/tag/v0.9.13).

---

## Developer setup

```powershell
copy .env.example .env
cd pit
go test ./...
go run ./cmd/pit companion
```

Desktop: `apps/desktop` — `npm install`, `npm run sidecar`, `npm run tauri`.
Web: `apps/web` — `VITE_PRIVY_APP_ID`, `VITE_HEALTH_URL`.
Sealer: Go 1.24 in `sealer/`.

```
pit/          CLI, companion, engine
sealer/       HPKE + VerifyE2EE
contracts/    Foundry (PitDeskID live; PitMemory/Forecasts/Receipts undeployed)
apps/web      Vite + Privy
apps/desktop  Local authorize surface
sdk/js        pit-os
sdk/mcp       pit-mcp
docs/diagrams SVG companions to the mermaid above
```

---

## License

MIT (`LICENSE`).
