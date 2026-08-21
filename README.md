# PIT — Private Alpha OS

Private research. Controlled execution. A desk that learns.

PIT seals your private trading book into 0G Direct TeeML, challenges its own thesis, sizes with a host engine (the model cannot set size), then waits for **you** to authorize an exact preview. A short-lived session can place or cancel. It cannot withdraw.

## 20 seconds

1. Connect **your wallet**.
2. Pick **MAINNET** or **TESTNET**.
3. Connect **your Hyperliquid account**.
4. Set **your policy**.
5. Create **your session** (order + cancel only).
6. Ask or wait for an opportunity.
7. Authorize the exact preview — or do nothing.

## Architecture

```
YOUR PRIVATE BOOK + PUBLIC MARKET
  → HPKE seal
  → Direct 0G TeeML (not the Router)
  → VerifyE2EE
  → on-chain teeSigner check
  → host engine + policy
  → YOU AUTHORIZE
  → session order|cancel
  → storage proof + receipt
  → calibration
```

Two environments, never mixed in one workspace:

| | TESTNET | MAINNET |
|---|---|---|
| 0G chain | Galileo 16602 | Aristotle 16661 |
| Book | Hyperliquid testnet | Hyperliquid mainnet |
| Sealed chat | Different catalog (no glm-5.2) | glm-5.2 Direct (proven path) |
| Agentic ID transfer | Official path exists; not claimed until a PIT transfer tx | Disabled |

## Prerequisites

- Go 1.22+
- Foundry (`forge`)
- Node 20+ (web)
- A wallet (Rabby / MetaMask). No seed phrase is ever collected.

## Environment

Copy `.env.example` to `.env` for local fixtures. **Do not commit `.env`.**

- Public RPCs and contract addresses are already filled.
- `PIT_PRIVY_APP_ID` is a public application id for wallet connect.
- `PIT_PRIVY_APP_SECRET` is server-only. Never put it in the web bundle.
- Hyperliquid session keys are generated locally and stored in the OS keychain (desktop/CLI). The web app cannot sign orders.

## Commands

```powershell
cd pit
go test ./...
go run ./cmd/pit init --network testnet --wallet 0xYourAddress

cd ..\contracts
forge install foundry-rs/forge-std OpenZeppelin/openzeppelin-contracts --no-commit
forge test -vvv
```

Live Hyperliquid book (optional):

```powershell
cd pit
go test -tags live ./internal/hl -count=1
```

Web (connect only):

```powershell
cd apps\web
copy .env.example .env
npm install
npm run dev
```

Open http://localhost:3000 — add `localhost:3000` in the Privy dashboard allowed origins.

## Contract addresses (mainnet)

- Desk ID: `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35`
- InferenceServing: `0x47340d900bdFec2BD393c626E12ea0656F938d84`
- LedgerManager: `0x2dE54c845Cd948B72D2e32e39586fe89607074E3`
- ERC-8004 Identity: `0x8004A169FB4a3325136EB29fA0ceB6D2e539a432`
- ERC-8004 Reputation: `0x8004BAa17C55a88189AE136b182e5fdA19dE9b63`

Testnet serving/ledger/8004: see `.env.example`. Galileo Desk ID is not deployed yet.

## How to reproduce 0G integration

Private book path is Direct TeeML + HPKE + VerifyE2EE. Router `sk-` keys are forbidden on that path. Catalog URL is used only to list models.

## CLI

`pit init|login|policy|ask|opportunities|forecast|preview|authorize|cancel|status|resolve|card|verify|kill`

Authorize is interactive. Piped `yes` is rejected; type `AUTHORIZE` on a TTY with `--i-understand`.

## MCP

Read-only tools: market, opportunities, forecast, status, card, verify. No authorize, order, cancel, or key export.

## SDK

`pit/sdk` is typed and **cannot sign**. Desktop/CLI own the session.

## Security model

- No global user identity.
- Session allowlist: `order`, `cancel` only.
- Preview hash binds the exact order.
- Durable SQLite ledger is per workspace + network.
- Memory keys are per workspace (`0x` + 32 bytes).
- Web never receives a Hyperliquid session key.

## Limitations (honest)

- Foundation Agentic ID transfer is not live on Aristotle.
- Galileo sealed committee is not the same model set as mainnet.
- Web onboarding copy is live; Hyperliquid session creation stays on desktop/CLI.
- Do not claim hardware quotes unless the verifier is wired.

## Desktop

`apps/desktop` is the local authorize surface. Session keys stay on the machine. Web remains connect-and-inspect.

```powershell
cd apps\desktop
npm install
npm run dev
```

Connect wallet → select network → read capability list → (desktop) connect trading account → set policy → create session → first private analysis → preview → authorize → verify on explorer.
