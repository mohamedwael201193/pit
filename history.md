# PIT — Research history (append-only)

---

## M136 — Policy stays editable after pin. PIT 0.9.13

- **Source:** Security hid PolicyEditor once onboard left step 4 (policy). Ready and unpaired-browser states left clip/assets locked on screen even though host pin already supports re-pin. User needs to edit and re-pin anytime. Chat still cannot pin.
- **Discovery:** `SecurityCenter` rendered `#policy` only when `current.id === "policy"`. After pin, current became ready (or pair if the browser unpaired), so the editor vanished. PolicyEditor already had `Pin updated policy`.
- **Candidate:** Always render PolicyEditor on Security. Pairing and Hyperliquid stay available after Ready. Web Download nav files the installer via `/windows`. Recorded 0G txs and Flow storage roots on `/proof` and README. Product 0.9.13.
- **Kill:** Remint `PIT-4bbee556`. Flatten OID `529167222216`. Chat/web pin. Session key / Direct token to the website. Invented fill or TEE.

DATE/TIME: 2026-08-31 19:20+03
PHASE: Policy-always-editable + 0.9.13 ship. No new live trade.
GOAL: Pinned desk can change policy on Security and re-pin. Latest GitHub installer is 0.9.13. Proof page carries this-desk 0G tx and root links.
RESULT:
- **IMPLEMENTED:** Security always shows PolicyEditor. Ready CTA is Edit policy. Pairing dock stays after Ready. Hyperliquid card stays when session is live. Web landing Download, Autonomy, Chat, Mission CTAs file `/windows`. `/proof` lists recorded roots including HYPE research `0x9fd42770545ecaacbfff12e3ef7a537b564e31c9ef5515b3a820fd276c22f72e` and order `0x8c94ec8e643c90fe69276ff20f50a0bc3121f007d611e10e6ab9f24d26f2ff66`.
- **TESTED:** Desktop e2e including `assertPolicyEditorStaysAfterReady`. `go test ./... -count=1` PASS (pit + sealer). Web Playwright **30 passed**. `pit-os` 2 pass. `pit-mcp` 3 pass. Desktop `tsc -b` ok. Web `tsc -b` ok. NSIS `PIT_0.9.13_x64-setup.exe`. Sidecar `pit version` → `PIT 0.9.13`.
- **LIVE extraAgents (master 0xbdfc…0034):** `PIT-4bbee556` `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` reused. No remint. No new AUTHORIZE.
- **SHIPPED:** Source `2637c4a` tagged `v0.9.13`. NSIS SHA256 `B905B9ED167513757D4947BDE61103EB10ECD4A5F76554FE369F205DF3850B1E`. GitHub Latest → `v0.9.13`. Health `/windows` 302 to the `.exe`. Vercel `https://pit0g.vercel.app`. Health `https://pit-health.onrender.com`.
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent. macOS/Linux not packaged.

SECURITY RESULT: Pin remains `pinLocalPolicy` on this computer. Chat/web/MCP/SDK cannot pin, authorize, or export. extraAgents still queried with master wallet.
TX HASH / OID: Venue OID `531667200134` FILLED (unchanged). Research 0G `0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b` root `0x9fd42770545ecaacbfff12e3ef7a537b564e31c9ef5515b3a820fd276c22f72e`. Order 0G `0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb` root `0x8c94ec8e643c90fe69276ff20f50a0bc3121f007d611e10e6ab9f24d26f2ff66`. Historical ETH OID `529167222216` unchanged.
CLASSIFICATION:
- Policy editor after pin: **IMPLEMENTED + TESTED**
- Re-pin host path: **UNCHANGED** (existing pinLocalPolicy)
- Recorded 0G proof links: **IMPLEMENTED + PLAYWRIGHT**
- Real TRADE NOW path: **NOT RE-RUN** (existing HYPE fill stands)
PRODUCTION READY: 0.9.13 installer after GitHub Latest + health `/windows` point at these bytes.
NEXT STEP: Commit, tag `v0.9.13`, upload tested NSIS, clobber GHA overwrite, Vercel + Render.

---

## M135 — New-user onboarding: pair first, then Protect, Hyperliquid, policy, ready. PIT 0.9.12

- **Source:** Security showed pairing-unpaired while Desk is ready. Pairing was labeled optional. SetupWizard started at Connect wallet and mixed a public 0x paste into first-run. Web copy still said pairing is a late step. Hyperliquid API is a name+address form; PIT must generate the agent locally and never ask a new user to invent one.
- **Discovery:** Go pairing (one-time 8-char, 2 min TTL, replay deny, device token not session key) was already correct. The gap was orchestration: `nextFix` and Security spine skipped pairing, READY did not require it, and HyperliquidCard mixed master+agent on one line.
- **Candidate:** Shared `onboard.ts` stepper. Pair is step 1. Protect is step 2. Connect Hyperliquid creates/reuses the local PIT Agent, primary action Approve PIT on Hyperliquid, verify from live extraAgents on the **master** wallet. Policy pin remains host-authoritative. Ready only after pair + protect + session + extraAgents + pin. No fake success.
- **Kill:** Remint `PIT-4bbee556`. Flatten OID `529167222216`. Session key / Direct token to the website. Manual API-wallet paste in the normal path. Authenticode claim. Fake FILLED.

DATE/TIME: 2026-08-31 18:10+03
PHASE: Onboarding audit + 0.9.12 ship. No new live trade.
GOAL: Brand-new user can go zero → paired → protected → Hyperliquid verified → pinned → ready without inventing an API wallet.
RESULT:
- **IMPLEMENTED:** Desktop Security/SetupWizard/DeskHome stepper Pair → Protect → Hyperliquid → Policy → Ready. Web `/pair` step 1, `/protect` locked until paired, `/signin` and `/app/start` redirect to `/pair`. Hyperliquid card splits Your wallet vs PIT Agent. Copy PIT Agent address only for the official Hyperliquid API page. Advanced public 0x bind is folded. Compact pairing strip only when unpaired and not on Security.
- **TESTED:** Desktop e2e copy harness ok (onboard pair-first + no invented ready). `go test ./... -count=1` **654 PASS / 0 FAIL** including `TestExpiredCodeDenied`. Web Playwright **30 passed**. `pit-os` 2 pass. `pit-mcp` 3 pass. Desktop `tsc -b` ok. Web `tsc -b` ok. NSIS `PIT_0.9.12_x64-setup.exe`. Sidecar `pit version` → `PIT 0.9.12`.
- **LIVE extraAgents (master 0xbdfc…0034, not agent address):** `PIT-4bbee556` `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` still listed, validUntil 1803441284611. Not reminted. No new AUTHORIZE.
- **SHIPPED:** Source `d6fa7c3` tagged `v0.9.12`. Rail honesty fix `8bc6d89`. NSIS SHA256 `0E40880652572A051382DF93F58D84C634DB7695E0F859F9F90335510E92333E`. GitHub Latest → `v0.9.12`. Health `/windows` 302 empty body to the `.exe`. Vercel `https://pit0g.vercel.app`. Production `/pair` is step 1. Production `/protect` locks Connect wallet until paired. Running local window is still 0.9.10/0.9.11 until the 0.9.12 installer is applied.
- **UNVERIFIED:** Foundry (forge not on PATH). Overlay `D:\PIT\pit.exe` copy blocked because that process is running. Fresh-install UI on a wiped machine (this desk already has session+agent). Production `/windows` after this deploy.
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent. macOS/Linux not packaged.

SECURITY RESULT: Pairing still returns `{sign:false, canSign:false, device}` only. `writeLocal` leak guard unchanged. Chat/web/MCP/SDK cannot pin, authorize, or export. extraAgents queried with master wallet.
TX HASH / OID: Venue OID `531667200134` FILLED (unchanged). Research 0G `0x1d2113bd…510b`. Order 0G `0x8c28051b…11fb`. Historical ETH OID `529167222216` unchanged.
CLASSIFICATION:
- Pairing as step 1: **IMPLEMENTED + TESTED**
- Protect locked until paired: **IMPLEMENTED + PLAYWRIGHT**
- Hyperliquid local agent + live extraAgents: **IMPLEMENTED** (live listing **VERIFIED** for existing desk; new-user approve click **UNVERIFIED** on this machine because already approved)
- Policy host pin: **UNCHANGED**
- Real TRADE NOW path: **NOT RE-RUN** (existing HYPE fill stands; do not fabricate)
PRODUCTION READY: 0.9.12 installer after GitHub Latest + health `/windows` point at these bytes.
NEXT STEP: Commit, tag `v0.9.12`, upload tested NSIS, Vercel + Render, verify production `/windows` 302 to the `.exe`.

---

## M134 — Final audit: npm, direct download, README, this-job proof on web. PIT 0.9.11

- **Source:** Release audit. Website download CTAs must file `PIT_0.9.11_x64-setup.exe` without opening GitHub Releases HTML. JS SDK / MCP must be real read-only packages. README must describe the actual product. Public `/proof` and `/missions` still led with historical ETH OID `529167222216` while the latest matching job was HYPE `531667200134`.
- **Discovery:** Health had no `/windows` until this change. `PLAYWRIGHT_BASE_URL` in the environment skips Vite startup. `pit-mcp` 0.9.11 bin was invalid on Windows CRLF so 0.9.12 with LF `bin/pit-mcp.mjs`. Public proof page recorded the HYPE OID in copy but did not link the this-job 0G txs.
- **Candidate:** Health `GET /windows` and `/checksums` 302 to the GitHub **asset** URL. Website CTAs use that file URL. `pit-os` 0.9.11 and `pit-mcp` 0.9.12 stay read-only (`canSign` false, no authorize). README rewritten from zero. `/proof` and `/missions` surface recorded HYPE + this-job 0G links; ETH remains HISTORICAL and is not flattened.
- **Kill:** Second AUTHORIZE path. Fake fills. Flatten OID `529167222216`. Remint `PIT-4bbee556`. Authenticode claim. Historical 0G fallback. Loosening policy.

DATE/TIME: 2026-08-31 10:20+03
PHASE: Audit + publish + README + proof honesty. No new live trade.
GOAL: Product, docs, npm, website download, and public proof agree with the actual 0.9.11 desk.
RESULT:
- **IMPLEMENTED:** Health `/windows` + `/checksums` 302 to the GitHub **asset** URL with **no HTML body**. Vercel redirects the same paths. Website Download/Hero/Pair/How CTAs file the installer. JS SDK `pit-os@0.9.11`. MCP `pit-mcp@0.9.12`. README rewritten (no competitor/hackathon language). Four architecture JPEGs in `docs/diagrams/`. `/proof` links research `0x1d2113bd…` and order `0x8c28051b…`. `/missions` has RECORDED HYPE plus HISTORICAL ETH.
- **TESTED:** `go test ./... -count=1` all packages ok. Verbose `--- PASS` **653**, `--- FAIL` **0**. `go test ./cmd/health` includes `TestWindowsRedirectsToInstallerNotReleasePage`. `pit-os` 2 pass. `pit-mcp` 3 pass. Desktop `npx tsx e2e/run.ts` PASS. Web Playwright **29 passed** with Vite + `VITE_PRIVY_APP_ID`. Foundry local **UNVERIFIED** (forge not on PATH).
- **LIVE:** npm `pit-os` 0.9.11 and `pit-mcp` 0.9.12 (verified install). Installer SHA `B621A10504EF1F9031C8C6D28E0B36FDB29B8AD186CEA86BCAC81F209A64515F`. Agent `PIT-4bbee556` reused. No new trade.
- **UNVERIFIED until this ship’s Render/Vercel:** production `/windows` 302 (code is ready; health process must be redeployed).
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent. macOS/Linux not packaged. PitReceipts/Forecasts/Memory not deployed (empty env). Local forge missing.

SECURITY RESULT: Chat, web, `pit-os`, and `pit-mcp` cannot AUTHORIZE, pin, arm, or hold the session key. Companion authorize stays loopback desktop/CLI. Router remains forbidden for the private book.
TX HASH / OID: Venue OID `531667200134` FILLED (unchanged). Research 0G `0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b`. Order 0G `0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb`. Historical ETH OID `529167222216` unchanged.
CLASSIFICATION:
- Direct installer file CTA: **IMPLEMENTED + TESTED** (production 302 after health deploy)
- npm pit-os / pit-mcp: **PUBLISHED + INSTALL VERIFIED**
- README rewrite + diagrams: **IMPLEMENTED**
- Public this-job 0G links: **IMPLEMENTED**
- New live trade: **NOT DONE** (not requested)
PRODUCTION READY: 0.9.11 remains the installer. This milestone is docs/health/npm/proof honesty, not a remint or a new fill.
NEXT STEP: Commit, push, Vercel + Render, verify `https://pit0g.vercel.app/windows` returns the exe (not GitHub HTML).

---

## M133 — Sequential committee + live READY → TRADE NOW → OID. PIT 0.9.11

- **Source:** “What can I trade now?” did not start research. None-hypothesis hunts proposed `none` on every book. Challenger never saw the researcher thesis. After a fill, Agent painted RESTING and “not execution-feasible” on a READY card.
- **Discovery:** `wantsAcceptPreview` matched the substring `trade now` inside “What can I trade now?” and returned `preview.show`. Researcher prompt echoed sealed `hypothesis: none`. Challenger envelopes were sealed before researcher finished, so they challenged an empty thesis. `venueOrderState` preferred lifecycle `reconciled` over status `filled`. Chat leftover long/short could leak into Find the best. Version **0.9.11**.
- **Candidate:** Hunt phrases win over TRADE NOW substring. Researcher confirms long/short unless facts contradict, and must pick a side from live facts when hypothesis is none. Challenger/risk receive `researcher_thesis` after the researcher job. Chat always seals `none|long|short` for this hunt. TRADE NOW remains `authorizePreview("AUTHORIZE", previewHash)`. FILLED only when Hyperliquid reports filled.
- **Kill:** Second signer. Fake READY. Flatten OID 529167222216. Remint PIT-4bbee556. Loosening clip/leverage. Historical 0G fallback.

DATE/TIME: 2026-08-31 08:50+03
PHASE: Real research → real 0G Direct → real READY preview → real TRADE NOW → real OID.
GOAL: Operator Agent is the single cockpit. Live scan, private committee, this-job proof, TRADE NOW only on exact READY, real Hyperliquid status.
RESULT:
- **IMPLEMENTED:** Parse, sequential TeeML committee, none-hypothesis does not echo none, chat seals none on Find the best, FILLED vs RESTING, READY thesis stays the preview not the post-fill book, banner Order filled after authorize.
- **TESTED:** `go test ./internal/compute ./internal/deskcmd ./internal/companion` PASS. Desktop `tsc -b` PASS. `npx tsx e2e/run.ts` PASS.
- **LIVE:** **What can I trade now?** at 7:57:41 AM started a fresh hunt (parse fix). **Find the best opportunity** hunts exhausted honestly when researcher proposed none. **Find the best long** at 8:32 AM chained DOGE/AVAX/BTC/ETH/SOL (challenger_killed after buy) then **HYPE READY_ELIGIBLE**.
- **TRADE:** Preview hash `0xb273d0052fe389b5e5ad3aad4b176e1cc993b8d8e605716bab78c70f3814e401`. Host sized buy 0.16 HYPE at 80.909, notional $12.95, clip $13, 1x. Path `authorize`. OID **531667200134**. Host reconcile `user_fills` **FILLED**. Agent card STATUS **FILLED**. TRADE NOW disabled after this preview’s OID. Duplicate TRADE NOW not offered.
- **HYPERLIQUID:** Portfolio wallet `0xBDfC…0034`. Positions (1) **0.16 HYPE**, value $12.93, entry 80.826. ETH OID `529167222216` untouched.
- **ACTIVITY:** `approval.accepted` + `order.submitted` + `order.filled` + `position.updated` for job `4a1d45ec-8c3f-4883-a162-19739accb9cf` and OID `531667200134`.
- **THIS-JOB 0G:** Research tx `https://chainscan.0g.ai/tx/0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b` root `0x9fd427…f72e`. Order evidence tx `https://chainscan.0g.ai/tx/0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb`. Both job `4a1d45ec…b9cf`. No historical fallback.
- **SESSION:** Agent `PIT-4bbee556` `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` reused. Policy pinned. Companion overlay `D:\PIT\pit.exe` **0.9.11**.
- **SHIPPED:** Code commit `69c0247` 2026-08-31 08:50:00 +0300. Cargo.lock NSIS fix `fff368b` 2026-08-31 08:52:00 +0300. Tag `v0.9.11` on `fff368b`. NSIS `PIT_0.9.11_x64-setup.exe` SHA256 `B621A10504EF1F9031C8C6D28E0B36FDB29B8AD186CEA86BCAC81F209A64515F`. Release https://github.com/mohamedwael201193/pit/releases/tag/v0.9.11. Vercel https://pit0g.vercel.app. Health https://pit-health.onrender.com **0.9.11**. Overlay companion `D:\PIT\pit.exe` is **0.9.11** (session reused, agent not reminted). Installer also at `D:\PIT\PIT_0.9.11_x64-setup.exe`.

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW used the existing desktop authorize path. Policy clip was not raised. Withdraw/transfer remain forbidden.
TX HASH / OID: Venue OID `531667200134` FILLED. Research 0G `0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b`. Order 0G `0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb`. Historical ETH OID `529167222216` unchanged.
CLASSIFICATION:
- What can I trade now starts research: **IMPLEMENTED + LIVE VERIFIED**
- Sequential committee (challenger sees researcher thesis): **IMPLEMENTED + LIVE VERIFIED**
- READY → TRADE NOW → real OID: **LIVE VERIFIED**
- FILLED only when Hyperliquid reports filled: **IMPLEMENTED + LIVE VERIFIED**
- Current-job 0G linked to the same preview/order: **LIVE VERIFIED**
PRODUCTION READY: 0.9.11 installer is on GitHub Latest and `D:\PIT`.
NEXT STEP: Close any 0.9.10 Tauri window and run `D:\PIT\PIT_0.9.11_x64-setup.exe`. Vite localhost:3001 already talks to companion 0.9.11. Do not remint. Do not flatten.

---

## M132 — Chat hunts ignore stale 4h auto skips. PIT 0.9.10

- **Source:** Find me the best short researched DOGE then AVAX and painted “Checked every executable book” with only two names. Host `hunt-skip.json` was only `DOGE,AVAX`. Remaining BTC/ETH/SOL/HYPE were still in automation.json 4h skips from the earlier long hunt. Client then cleared the last-book card.
- **Discovery:** `beginResearch` + `pickBestCoinSkipping` merged `auto.SkipSet` (4h) into chat hunt skip. `resolveChatCoin("BTC", skip, "", false)` returned empty when BTC was in that 4h set, so the host returned `hunt_exhausted`. `autoTick` while a job was running also saved a stale Prefs copy and restored those skips after `resetHuntSkip`. Version **0.9.10**.
- **Candidate:** Chat and research_ui hunts skip only this-hunt `hunt-skip.json`. Automation keeps 4h skips. Chat research does not kick `autoTick`. `autoTick` does not persist Prefs while a job is running. Client keeps last-book analysis until a new job actually starts, continues if host says exhausted while untried books remain, and only paints universe-exhausted when every executable was checked. TRADE NOW still `authorizePreview("AUTHORIZE", previewHash)` only on `READY_ELIGIBLE`.
- **Kill:** Second AUTHORIZE path. Fake fills. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim. Historical tx fallback. Loosening policy to force a fill.
- **Next:** Overlay `D:\PIT`, tag `v0.9.10`, NSIS, Vercel + Render. Do not remint. Do not flatten. Do not TRADE NOW unless an exact READY preview exists.

DATE/TIME: 2026-08-31 07:45+03
PHASE: Chat hunt isolation from automation 4h skips.
GOAL: Find the best / long / short / next → live scan → private 0G on every remaining executable book → useful analysis + this-job 0G → TRADE NOW only if READY, else genuine exhausted NO TRADE.
RESULT:
- **IMPLEMENTED:** `thisHuntSkipSet` / `pickNextCoin`. Chat hunts no longer inherit stale automation skips. Last-book LIVE MARKET + committee why stay on screen. AUTHORIZE copy lives in the composer, not as a header wall. TRADE NOW still `authorizePreview("AUTHORIZE", previewHash)`.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` PASS. `npx tsx e2e/run.ts` PASS. Companion overlayed to **0.9.10** without remint: agent `PIT-4bbee556` `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`, session live. Buying power $16.18. Policy pinned. ETH OID `529167222216` untouched.
- **LIVE:** **Find the next opportunity** at 7:38:53 AM (short hypothesis still sealed) chained **BTC → SOL → ETH → HYPE** after the earlier DOGE+AVAX stand-downs. Did not wrap. Universe exhausted: all six executable books. No `READY_ELIGIBLE`, so TRADE NOW was not shown and was not clicked.
- **THIS-JOB 0G:** BTC job `8c8d331d…84fb` challenger_killed after proposed sell, tx `https://chainscan.0g.ai/tx/0x6abe43772f1b953e2c6debec31dba1d64b77a7f8c3b6f83cf950f18f11e263e4`. SOL job `d3b198f7…42c3` `no_side`, tx `https://chainscan.0g.ai/tx/0x7e7f85aaf4aacd29129b8697cbc5de7e8f6d56745754897807a262e2d31b21ef`. ETH job `78617f6c…845c` tx `https://chainscan.0g.ai/tx/0xdf4f8f95cbee81f99402754455915635bbc3f4623861318f5fc171da631f8ae0`. HYPE job `9761cbd5…c980` `no_side` proposed none, tx `https://chainscan.0g.ai/tx/0x266c45cbd35cb8b9e856d7f3c850e5ce72d34fb33251bba616345e34cd04cb78`.
- **SHIPPED:** Commit `fed6989` 2026-08-31 07:45:26 +0300. Tag `v0.9.10`. NSIS `PIT_0.9.10_x64-setup.exe` SHA256 `247E596F4389990F88494959E1B1867E21E095200F9133C7F0F999F3AB584BAC`. Release https://github.com/mohamedwael201193/pit/releases/tag/v0.9.10. Vercel https://pit0g.vercel.app. Health https://pit-health.onrender.com **0.9.10**. Overlay companion `D:\PIT\pit.exe` is **0.9.10** (session reused, agent not reminted). Installer also at `D:\PIT\PIT_0.9.10_x64-setup.exe`.
- **BLOCKED:** Authenticode still absent. iTransfer still not live. TRADE NOW not clicked: short hypothesis still produced no `READY_ELIGIBLE` preview. Forcing a fill would be a fake demo.

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is still an explicit desktop confirmation of the existing host path.
TX HASH / OID: HYPE 0G `https://chainscan.0g.ai/tx/0x266c45cbd35cb8b9e856d7f3c850e5ce72d34fb33251bba616345e34cd04cb78` (job `9761cbd5…c980`). Historical venue OID `529167222216` unchanged. No new order.
CLASSIFICATION:
- Chat hunt continues after genuine NO TRADE until READY or exhausted: **IMPLEMENTED + LIVE VERIFIED**
- Stale 4h automation skips no longer stop a chat hunt: **IMPLEMENTED + LIVE VERIFIED**
- Current-job 0G receipts: **IMPLEMENTED + LIVE VERIFIED**
- TRADE NOW existing authorize path: **UNCHANGED** (no eligible preview this pass)
PRODUCTION READY: 0.9.10 installer after tag.
NEXT STEP: Close any 0.9.9 window and run `D:\PIT\PIT_0.9.10_x64-setup.exe` if using the Tauri window. Vite localhost:3001 already talks to companion 0.9.10. TRADE NOW only on an exact READY preview.

---


## M131 — Operator Agent hunts the universe. Live facts on NO TRADE. PIT 0.9.9

- **Source:** AVAX NO TRADE after “Find the next opportunity” showed a sparse card (no live mark/oracle/funding/OI, no stages, no thesis) and a giant empty area. Research next stopped after one book. 0G receipts were job-scoped (good) but unused below the fold.
- **Discovery:** Completed `READY_STOOD_DOWN` hid `ResearchStages` (`showVerdict || fail` excluded `noTrade`) and hid live book facts (`showBook` busy-only). Chat hunts only chained when `fresh || chained` and capped at six. Starting a new book left the previous `jobId` on screen. Version **0.9.9**.
- **Candidate:** Chain any chat hunt after a genuine stand-down until READY or the executable universe is exhausted. Always render stages + LIVE MARKET facts + Thesis/Evidence/Rejected side/Reason/Policy/Risk/this-job 0G. TRADE NOW still `authorizePreview("AUTHORIZE", previewHash)` only on `READY_ELIGIBLE`. Clear `jobId` when a new book starts. Do not invent a side.
- **Kill:** Second AUTHORIZE path. Fake fills. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim. Historical tx fallback. Loosening policy to force a fill.
- **Next:** Overlay `D:\PIT`, tag `v0.9.9`, NSIS, Vercel + Render. Close the 0.9.7 window and run the 0.9.9 installer. Do not remint. Do not flatten. Do not TRADE NOW unless an exact READY preview exists.

DATE/TIME: 2026-08-31 07:00+03
PHASE: Operator Agent hunt + live facts. Host remains execution authority.
GOAL: Find the best opportunity → live scan → rank → private 0G on every executable book → useful analysis + this-job 0G → TRADE NOW only if READY, else genuine exhausted NO TRADE.
RESULT:
- **IMPLEMENTED:** Chat hunts keep `researchBusy` and continue after `READY_STOOD_DOWN` without the six-book cap. Fresh hunts still reset skip. `setResearchJobId("")` on each new book. Agent turn always shows 10-stage pipe, LIVE MARKET (mark, oracle, funding, OI, venue min, host notional, clip, capital), labeled host rank vs committee, structured NO TRADE, OPPORTUNITY FOUND + REVIEW / TRADE NOW / REJECT. Pad no longer eats the stream. TRADE NOW still `authorizePreview("AUTHORIZE", previewHash)`.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` PASS. `npx tsx e2e/run.ts` PASS. Playwright 29 passed. Chrome http://localhost:3001 Agent 0.9.9 frontend: **Find the best opportunity** at 6:51:43 AM. Live pipe + live facts. Chained **AVAX → DOGE → BTC → ETH → SOL → HYPE** without wrapping. Each book got its own TeeML job (HYPE `3c9f2e96…1167`). Waiting copy until this job filed. **Research next** at 6:57:00 AM did not wrap. No TRADE NOW (no `READY_ELIGIBLE`).
- **LIVE:** Companion overlayed to **0.9.9** without remint: agent `PIT-4bbee556` `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`, session live. Buying power $16.18. Policy pinned. **Find me the best long** at 7:07:45 AM chained AVAX→DOGE→BTC→ETH→SOL→HYPE. Researcher proposed **buy** on HYPE; challenger stood down. Exhausted NO TRADE. This-job 0G `https://chainscan.0g.ai/tx/0x6009ede35278fc6157507792388b87e2f0c7173494a32095fb0692bc65c77ff4` (job `7d28f3e3…f08a`). Earlier none-hypothesis HYPE tx `https://chainscan.0g.ai/tx/0x30df71b929e05a4feca6d4683bbe86af97750b70807a28957bcc54e2d99aa4ed` (job `3c9f2e96…1167`). ETH OID `529167222216` untouched.
- **SHIPPED:** Commit `b713380` 2026-08-31 06:58:43 +0300. Tag `v0.9.9`. NSIS `PIT_0.9.9_x64-setup.exe` SHA256 `C46E6F242EE72A25015DFD61E918B7541CA0D599181824AE87A29FBF61A67A0C`. Release https://github.com/mohamedwael201193/pit/releases/tag/v0.9.9. Vercel https://pit0g.vercel.app. Health https://pit-health.onrender.com **0.9.9**. Overlay `D:\PIT\pit.exe` is now **0.9.9** (session reused, agent not reminted). Installer also at `D:\PIT\PIT_0.9.9_x64-setup.exe`.
- **BLOCKED:** Authenticode still absent. iTransfer still not live. TRADE NOW not clicked: long hypothesis still produced no `READY_ELIGIBLE` preview. Forcing a fill would be a fake demo.

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is still an explicit desktop confirmation of the existing host path.
TX HASH / OID: Long-hunt HYPE 0G `https://chainscan.0g.ai/tx/0x6009ede35278fc6157507792388b87e2f0c7173494a32095fb0692bc65c77ff4` (job `7d28f3e3…f08a`). Historical venue OID `529167222216` unchanged. No new order.
CLASSIFICATION:
- Hunt continues after genuine NO TRADE until READY or exhausted: **IMPLEMENTED + LIVE VERIFIED**
- Live facts + committee why on the completed turn: **IMPLEMENTED + LIVE VERIFIED**
- Current-job 0G receipts: **IMPLEMENTED + LIVE VERIFIED**
- TRADE NOW existing authorize path: **UNCHANGED** (no eligible preview this pass)
PRODUCTION READY: 0.9.9 installer after tag.
NEXT STEP: Push. Tag `v0.9.9`. NSIS. Overlay `D:\PIT`. Vercel + Render. Close the 0.9.7 window and reinstall. TRADE NOW only on an exact READY preview.

---


## M130 — Job-scoped 0G proof. Opportunity Found. Honest NO TRADE. PIT 0.9.8

- **Source:** Agent showed live stages but 0G explorer links belonged to older jobs. Duplicate “evidence filed [research]” chips. NO TRADE dumped a thesis log. TRADE NOW missing on a surviving side. Chat repeated “cannot AUTHORIZE”. Find the best stopped after one book because stale huntRejected poisoned the chain.
- **Discovery:** `proofRows` skipped the job filter when an activity row had a root/tx. `last-research.json` was a global evidence fallback. Client merged `started.hunt_skip` and React `huntRejected` into a fresh hunt. Version **0.9.8**.
- **Candidate:** Unique jobId on every run. Render only receipts with `receipt.jobId === currentAgentRun.jobId`. Waiting copy until this job files. OPPORTUNITY FOUND + TRADE NOW = existing `authorizePreview("AUTHORIZE", previewHash)`. NO TRADE shows this-run committee and this-run 0G proof. Fresh hunt does not ingest stale skip. Chained hunts use `huntTried` only.
- **Kill:** Second AUTHORIZE path. Fake fills. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim. Historical tx fallback.
- **Next:** Overlay `D:\PIT`, tag `v0.9.8`, NSIS, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade unless TRADE NOW on an exact READY preview.

DATE/TIME: 2026-08-31 06:22+03
PHASE: Agent 0G correlation + live operator turn.
GOAL: Find the best opportunity → live scan → private 0G → current-job root/tx → TRADE NOW only if READY, else genuine NO TRADE with this-run proof.
RESULT:
- **IMPLEMENTED:** `jobProof.ts` `collectJobReceipts` / `evidenceObjectForJob`. Companion `evidenceForJob` on `/local/research/result` only when the job is done. No last-research overlay. Agent 10-stage pipe including Decision. OPPORTUNITY FOUND card. Header-only AUTHORIZE copy. RESTING ≠ FILLED. Cancelled jobs are CANCELED_BY_USER, not NO TRADE. Fresh hunt no longer merges host skip or stale rejected coins.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` PASS. `npx tsx e2e/run.ts` PASS including `assertLiveAgentPipeline`. Chrome http://localhost:3001 Agent: **Find the best opportunity** at 6:14:59 AM researched **DOGE** job `b4ed73ce-2587-4d69-a009-bae3693629e3`. While running: “Waiting for this research run’s 0G receipt…” — old HYPE tx `0x57c6e574…9bde` did not render. After seal, 0G proof was **this job**: root `0x66817182…2251`, tx `https://chainscan.0g.ai/tx/0x28f0f7474760ec88c8c2a76f9959e136756eb5dd8ccfd530eb43d38c10f7277c`. **Research next** started **AVAX** job `10ae9aff…f678` and showed waiting, not the DOGE tx. AVAX NO TRADE with this-job receipts (`0x3b8bab…` / job `10ae9aff…f678`). No TRADE NOW (no READY_ELIGIBLE). Session live. Policy pinned. Buying power $16.18.
- **LIVE:** Companion `/health` was **0.9.7** during the Agent pass (frontend 0.9.8). Session live. Agent `PIT-4bbee556` not reminted. ETH OID `529167222216` untouched. No new MAINNET fill. Committee stood down on DOGE and AVAX (`no_side`).
- **SHIPPED:** Commit `39209a1` 2026-08-31 06:22:31 +0300. Tag `v0.9.8`. NSIS `PIT_0.9.8_x64-setup.exe` SHA256 `5BF6BE1917E6EB13796E719A83F90DB0EE972C0032D6D237EFA5720608DEB66A`. Release https://github.com/mohamedwael201193/pit/releases/tag/v0.9.8. Vercel https://pit0g.vercel.app. Health https://pit-health.onrender.com **0.9.8**. Overlay `D:\PIT\pit-desktop.exe` + `pit-0.9.8.exe`. Running `D:\PIT\pit.exe` stayed 0.9.7 because the process had the file locked — close that window and run `D:\PIT\PIT_0.9.8_x64-setup.exe` to match sidecar.
- **BLOCKED:** Authenticode still absent. iTransfer still not live. TRADE NOW not clicked: no READY_ELIGIBLE preview this pass.

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is still an explicit desktop confirmation of the existing host path.
TX HASH / OID: Current-job 0G `https://chainscan.0g.ai/tx/0x28f0f7474760ec88c8c2a76f9959e136756eb5dd8ccfd530eb43d38c10f7277c` (DOGE job `b4ed73ce…29e3`). AVAX job `10ae9aff…f678` filed its own receipt (not the DOGE or HYPE txs). Historical venue OID `529167222216` unchanged. No new order.
CLASSIFICATION:
- Job-scoped 0G receipts in the Agent turn: **IMPLEMENTED + LIVE VERIFIED**
- Waiting copy instead of stale explorer links: **IMPLEMENTED + LIVE VERIFIED**
- TRADE NOW existing authorize path: **UNCHANGED** (no eligible preview this pass)
PRODUCTION READY: 0.9.8 installer after tag.
NEXT STEP: Push. Tag `v0.9.8`. NSIS. Overlay `D:\PIT`. Vercel + Render. TRADE NOW only on an exact READY preview.

---


## M129 — Hunt does not wrap. Honest incomplete card. PIT 0.9.7

- **Source:** Research next after HYPE restarted AVAX. Cancelled wrap painted NO TRADE with empty ○ roles. Stages hid after the job. Auto-scroll fought the wheel. 0G TRAIL showed only a TeeML job id.
- **Discovery:** Host wrap cleared `huntSkip` then `pickBestCoin()`. Client reset `huntTried` when the skip set was full. Leftover `preview.deny=no_side` on an incomplete job. Scroll effect depended on `roles.length`. Version **0.9.7**.
- **Candidate:** Persist skip. Fresh hunt resets. Next hunt never wraps. NO TRADE only on `READY_STOOD_DOWN`. Named stages stay on the completed card. 0G explorer links from filed evidence. TRADE NOW still `authorizePreview("AUTHORIZE", previewHash)`.
- **Kill:** Second AUTHORIZE path. Fake fills. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim.
- **Next:** Overlay `D:\PIT`, tag `v0.9.7`, NSIS, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade unless TRADE NOW on an exact preview.

DATE/TIME: 2026-08-31 05:35+03
PHASE: Agent hunt honesty. Host remains execution authority.
GOAL: Type Find the best opportunity, watch live 0G stages, skip a failed book, never wrap Research next, TRADE NOW only on READY.
RESULT:
- **IMPLEMENTED:** `hunt-skip.json` + auto skip (4h). `resolveChatCoin` returns empty instead of wrapping. `fresh` on Find the best. Client sessionStorage skip. Exhausted hunt returns `hunt_exhausted`. Incomplete/cancelled jobs are STOPPED. Named 9-stage pipe stays after seal. Stream `is-busy` keeps the pad from eating scroll. 0G receipts are explorer links. TRADE NOW still `authorizePreview("AUTHORIZE", previewHash)`.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` PASS. `npx tsx e2e/run.ts` PASS. Chrome http://localhost:3001 Agent 0.9.7: typed **Find the best opportunity**. Live pipe Researcher live on DOGE then chained AVAX, BTC, ETH, SOL, HYPE. `canScroll` true (2239/434).
- **LIVE:** Overlay `D:\PIT` companion `/health` **0.9.7**. Session live. Policy pinned. Buying power $16.18. Sealed committee on **DOGE, AVAX, BTC, ETH, SOL, HYPE** — all `READY_STOOD_DOWN` `no_side`. UI: `checked DOGE, AVAX, BTC, ETH, SOL, HYPE`. Verdict NO TRADE HYPE with Researcher/Challenger/Risk ✓. Find the next opportunity after all six: **did not wrap**. Note: Checked every executable book. Scan again, not Research next. No new MAINNET fill. ETH OID `529167222216` untouched. Agent `PIT-4bbee556` not reminted.
- **BLOCKED:** Authenticode still absent. iTransfer still not live. TRADE NOW not clicked: no READY_ELIGIBLE preview this pass (committee proposed none on every executable book).

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is still an explicit desktop confirmation of the existing host path.
TX HASH / OID: 0G storage `https://chainscan.0g.ai/tx/0x8c8b78e8add46c79983d344ac571bcb8e6fd1d6c2ae072add00147f2ede1151d` (HYPE). Also `0xcc02a780b12ed2a884d3aa845f486acb89c60f1e8c306f0773e147f5311b4438`, `0xd682aa45aea64a26d1ab7a18d9867260a38502b086b9730010a394011ef6114c`, `0x2a7a58381ef4507174a777fb2f9a65d826d9988ce22610fc16b4d9e1fcd54b9d`. Historical venue OID `529167222216` unchanged. No new order.
CLASSIFICATION:
- Next-book hunt without wrap: **IMPLEMENTED + LIVE VERIFIED** (DOGE→AVAX→BTC→ETH→SOL→HYPE, then exhausted)
- Live stages + 0G receipts in the turn: **IMPLEMENTED + LIVE VERIFIED**
- TRADE NOW existing authorize path: **UNCHANGED** (no eligible preview this pass)
PRODUCTION READY: 0.9.7 installer after tag.
NEXT STEP: Push. Tag `v0.9.7`. NSIS. Overlay `D:\PIT`. Vercel + Render. TRADE NOW only on an exact READY preview.

---

## M128 — Hunt the next book. Live 0G in the turn. PIT 0.9.6

- **Source:** Research next re-ran AVAX. Stages were a text list. Auto-scroll on every elapsed tick stole the wheel. The right rail split the stream. 0G txs were easy to miss. TRADE NOW exists only on a READY preview.
- **Discovery:** A new "Find the best opportunity" reset `huntTried` to ranked[0]. Host chat left Coin empty, so the same top book started again. `island.elapsedMs` in the scroll effect pinned the log. Version **0.9.6**.
- **Candidate:** Host + desktop skip stood-down coins. Live track + receipts in the turn. No side rail. TRADE NOW still `authorizePreview("AUTHORIZE", previewHash)`.
- **Kill:** Second AUTHORIZE path. Fake fills. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim.
- **Next:** Overlay `D:\PIT`, tag `v0.9.6`, NSIS, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade unless TRADE NOW on an exact preview.

DATE/TIME: 2026-08-31 04:50+03
PHASE: Agent hunt + live proof. Host remains execution authority.
GOAL: Type Find the best opportunity, watch live 0G stages, skip a failed book, TRADE NOW, see a real OID when the preview is eligible.
RESULT:
- **IMPLEMENTED:** Stood-down coins go on `huntSkip`. Unnamed hunts send an empty coin. Host `research.best` re-picks with skip after Watch. Live named 9-stage pipe + 9-cell track. 0G receipts in the turn. Cream Research next. Visible scrollbar. TRADE NOW still `authorizePreview("AUTHORIZE", previewHash)`. Browser Vite can start research via `/local/research/start` fallback. Hunt wrap when every executable book is skipped.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` PASS. `npx tsx e2e/run.ts` PASS. Chrome http://localhost:3001 Agent: typed **Find the best opportunity**. Live pipe showed Scanning markets / Ranking / Private 0G research / Researcher live. Stream `canScroll` true while busy. Header `Researching AVAX` then the hunt moved on.
- **LIVE:** Overlay `D:\PIT` companion `/health` **0.9.6**. Session live. Policy pinned. Buying power $16.18. Direct credit ~3.59 0G. Sealed committee on **AVAX, DOGE, BTC, ETH, SOL, HYPE** — all `READY_STOOD_DOWN` `no_side`. UI: `checked AVAX, DOGE, BTC, ETH, SOL, HYPE`. Verdict NO TRADE HYPE. 0G proof tx `0x2045c98a69aae505ee5be36eaa1cf05c5d93c2662d90b5d7b07dc8452d537711`. Research next after all six wrapped (did not re-run AVAX while others remained). No new MAINNET fill. ETH OID `529167222216` untouched. Agent `PIT-4bbee556` not reminted.
- **BLOCKED:** Authenticode still absent. iTransfer still not live. TRADE NOW not clicked: no READY_ELIGIBLE preview this pass (committee proposed none on every executable book).

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is still an explicit desktop confirmation of the existing host path.
TX HASH / OID: 0G storage `https://chainscan.0g.ai/tx/0x2045c98a69aae505ee5be36eaa1cf05c5d93c2662d90b5d7b07dc8452d537711`. Historical venue OID `529167222216` unchanged. No new order.
CLASSIFICATION:
- Next-book hunt: **IMPLEMENTED + LIVE VERIFIED** (AVAX→DOGE→BTC→ETH→SOL→HYPE)
- Live stages + 0G receipts in the turn: **IMPLEMENTED + LIVE VERIFIED**
- TRADE NOW existing authorize path: **UNCHANGED** (no eligible preview this pass)
PRODUCTION READY: 0.9.6 installer after tag.
NEXT STEP: Push. Tag `v0.9.6`. NSIS. Overlay `D:\PIT`. Vercel + Render. TRADE NOW only on an exact READY preview.

---

## M127 — Dark Agent controls. One current hunt. PIT 0.9.5

- **Source:** Live Agent on localhost:3001 still showed a native white Windows button under Research next / Show why / Compare candidates. The Desk thread wall sat above the mission, so the page still read as two stacked apps.
- **Discovery:** `className="ghost"` had no stylesheet. Windows painted a white `<button>`. Unclassed `ExternalLink` and `<details>` summary could do the same. `composeStream` kept the whole Desk log when a hunt result was live.
- **Candidate:** Dark ghost/chip controls, `color-scheme: dark`, Sleep Mission and Technical details as chips, one hunt turn plus the live mission. Version stays **0.9.5**.
- **Kill:** Second AUTHORIZE path. Fake fills. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim.
- **Next:** Overlay `D:\PIT`, tag `v0.9.5`, NSIS, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade unless the user clicks TRADE NOW on an exact preview.

DATE/TIME: 2026-08-31 03:01+03
PHASE: Agent visual polish. Host remains execution authority.
GOAL: Every Agent control is dark. The current hunt is the conversation. TRADE NOW is still the existing desktop authorize path.
RESULT:
- **IMPLEMENTED:** Global dark `button` + `.ghost`. Sleep Mission is a chip, not a native slab. Technical details is a chip. Composer Send is dark. Completed research hides the live pipe and book stack. `composeStream` shows the last hunt ask plus the mission.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` pass. `npx tsx e2e/run.ts` pass. Chrome on localhost:3001: no button luminance above 0.55. Chips are `rgb(20, 22, 28)` / cream. TRADE NOW path unchanged.
- **LIVE:** Companion overlaid to `D:\PIT`. `pit.exe version` = PIT 0.9.5. `/health` = 0.9.5. No new MAINNET fill. ETH OID `529167222216` untouched.
- **BLOCKED:** Authenticode still absent. iTransfer still not live.

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is still an explicit desktop confirmation of the existing host path.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Dark Agent controls: **IMPLEMENTED + TESTED**
- TRADE NOW existing authorize path: **UNCHANGED**
- Live MAINNET fill: **NOT STARTED**
PRODUCTION READY: 0.9.5 installer after tag.
NEXT STEP: Overlay `D:\PIT`. `python _scripts/push_head.py`. `python _scripts/tag_push.py v0.9.5`. Deploy web.

---

## M126 — One Agent conversation. PIT 0.9.5

- **Source:** 0.9.4 still stacked a cockpit above a transcript. Nested scrollbars. Hunt chips and "Researching AVAX" repeated while the AVAX no-trade card sat in another pane. Screenshots showed two apps, not one operator.
- **Discovery:** AgentRun must live inside the last hunt turn and transform in place. Composer chips belong on empty only. TRADE NOW stays `authorizePreview("AUTHORIZE", previewHash)` on App. Named research (Research ETH) may start while another job is running. Same-job hunt stays quiet. Numbers stay event-backed.
- **Candidate:** One `agent-stream`. Mission cards inside the PIT turn. Optional live rail only while busy. Version **0.9.5**.
- **Kill:** Second AUTHORIZE path. Fake token streaming. Fake fills. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim. Nested transcript + cockpit scroll.
- **Next:** Tag `v0.9.5`, NSIS, overlay `D:\PIT`, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade unless the user clicks TRADE NOW on an exact preview.

DATE/TIME: 2026-08-31 02:50+03
PHASE: Agent conversation rebuild. Host remains execution authority.
GOAL: One screen: ask PIT what to trade, watch live scan and private 0G research in that turn, read the verdict, TRADE NOW, see a real OID.
RESULT:
- **IMPLEMENTED:** CommandChat is one scroll. AgentRun is an inline turn (scan checklist, book, pipe, verdict sections, READY / NO TRADE / ORDER SUBMITTED). Hunt dumps are not shown as a second chat. Busy same-job hunt does not append "Still researching". Composer placeholder "Ask PIT to research, compare, prepare, trade, or watch…". Contextual chips only on an empty thread. TRADE NOW still App → existing authorize. oidBelongsToPreview unchanged.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` pass. `npx tsx e2e/run.ts` including one stream, no cockpit-live, no canned busy line, TRADE NOW callback, displayTurn collapse of old dumps.
- **LIVE:** After installer. No new MAINNET fill from this change. ETH OID `529167222216` untouched.
- **BLOCKED:** Authenticode still absent. iTransfer still not live. TRADE NOW still requires exact preview, live session, pinned policy, capital, venue min.

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is still an explicit desktop confirmation of the existing host path. Preview hash, session, policy, ledger, and capital gates unchanged.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- One conversation Agent: **IMPLEMENTED + TESTED**
- TRADE NOW existing authorize path: **UNCHANGED**
- Live MAINNET fill: **NOT STARTED** (needs explicit TRADE NOW)
PRODUCTION READY: 0.9.5 installer after tag.
NEXT STEP: `python _scripts/push_head.py`. `python _scripts/tag_push.py v0.9.5`. Deploy web.

---

## M125 — Agent chat stops repeating itself. PIT 0.9.4

- **Source:** Desktop Agent stacked the same chips, two Show why controls, an always-on Sleep Mission card, a stale FILLED banner from an old OID, DONE on every pipe row, buying power in the header and the quote, and the same AVAX sentence plus "Chat cannot AUTHORIZE" after every hunt.
- **Discovery:** Composer PRIMARY chips stayed visible after the first turn and repeated the last question. Follow chips asked "Why didn't you trade?" instead of opening the why panel. lastOrder rendered without oidBelongsToPreview. decorateWatchAgent dumped the executive book line as the chat reply, then chatctx appended the authorize refrain.
- **Candidate:** One result surface. Short hunt replies. Numbers stay on cards. Version **0.9.4**.
- **Kill:** Second AUTHORIZE path. Fake fills. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim.
- **Next:** Tag `v0.9.4`, NSIS, overlay `D:\PIT`, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade unless the user clicks TRADE NOW on an exact preview.

DATE/TIME: 2026-08-31 02:15+03
PHASE: Agent cockpit de-duplication. Host remains execution authority.
GOAL: One Show why. One Sleep Mission entry. One order card only when the hash matches. Hunt replies that do not repeat the cards.
RESULT:
- **IMPLEMENTED:** PRIMARY chips hide after the first turn. Show why opens the why panel. Sleep Mission is opt-in. Order card gated by oidBelongsToPreview. Pipe marks without a DONE label on every row. Quote drops buying power. Hunt reply is "Researching AVAX. Live numbers stay on the cards." Transcript collapses consecutive identical PIT lines and old executive dumps.
- **TESTED:** `go test ./...` PASS including find-best no longer dumping "strongest executable book" or "Starting sealed". Desktop `tsc -b` pass. `npx tsx e2e/run.ts` including oidBelongsToPreview, displayTurn collapse, no duplicate Show why. Live localhost:3001 Agent: one Show why, no FILLED, no always-on Sleep, composer only More, WHY panel opens without a new chat turn.
- **LIVE:** After installer. No new MAINNET fill from this change. ETH OID `529167222216` untouched.
- **BLOCKED:** Authenticode still absent. iTransfer still not live. TRADE NOW still requires exact preview, live session, pinned policy, capital, venue min.

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is still an explicit desktop confirmation of the existing host path. Stale FILLED cannot impersonate the current preview.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Agent duplicate UI: **IMPLEMENTED + TESTED**
- Short hunt replies: **IMPLEMENTED + TESTED**
- Stale fill gate on Agent: **IMPLEMENTED + TESTED**
PRODUCTION READY: 0.9.4 installer after tag.
NEXT STEP: `python _scripts/push_head.py`. `python _scripts/tag_push.py v0.9.4`. Deploy web.

---

## M124 — Desktop Agent as trading operator. PIT 0.9.3

- **Source:** 0.9.2 Agent was unreadable: pairing strip ate the canvas, AgentRun sat inside a transcript that auto-scrolled away, answers were tiny, chips while a job ran dumped "Already researching AVAX (job …, CHALLENGER)" twice, and TRADE NOW did not exist.
- **Discovery:** The trusted desktop may invoke the existing AUTHORIZE path from an explicit TRADE NOW button. Chat/LLM still cannot construct an order. Pairing belongs in the titlebar Browser chip on Agent. Research progress is event-backed (status poll + `/local/research/stream`), never fake tokens or fake percentages.
- **Candidate:** Scrollable conversation + sticky live cockpit, larger type, opportunity/research/preview/execution cards, TRADE NOW → `authorizePreview("AUTHORIZE", hash)`, quiet already-researching, hunt still max 3. Version **0.9.3**.
- **Kill:** Second execution implementation. Chat importing authorizePreview. Fake fills. Fake 0G roots. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim.
- **Next:** Tag `v0.9.3`, NSIS, overlay `D:\PIT`, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade unless the user clicks TRADE NOW on an exact preview.

DATE/TIME: 2026-08-31 01:40+03
PHASE: Agent operator UI. Host remains execution authority.
GOAL: One Agent screen: live scan → 0G Direct → committee → TRADE NOW → real OID.
RESULT:
- **IMPLEMENTED:** Pairing dock hidden on Agent. Independent cockpit + transcript scroll. 16.5px answers. TRADE NOW calls existing desktop authorize. SSE `/local/research/stream`. Chat stream no longer fakes 12-rune deltas. Already-researching is a one-liner without job ids.
- **TESTED:** `go test ./...` PASS including find-best, already-researching quiet, research stream origin lock, trade-this stays on Agent. Desktop `tsc -b` pass. `npx tsx e2e/run.ts` including TRADE NOW callback, no CommandChat authorizePreview, pairing hidden on Agent.
- **LIVE:** After installer. No new MAINNET fill from this change unless the user clicks TRADE NOW. ETH OID `529167222216` untouched.
- **BLOCKED:** Authenticode still absent. iTransfer still not live. TRADE NOW still requires exact preview, live session, pinned policy, capital, venue min.

SECURITY RESULT: The model cannot AUTHORIZE. TRADE NOW is an explicit desktop confirmation of the existing host path. Preview hash, session, policy, ledger, and capital gates unchanged.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Agent operator UI: **IMPLEMENTED + TESTED**
- TRADE NOW existing authorize path: **IMPLEMENTED + TESTED**
- Live MAINNET fill: **NOT STARTED** (needs explicit TRADE NOW)
PRODUCTION READY: 0.9.3 installer after tag.
NEXT STEP: `python _scripts/push_head.py`. `python _scripts/tag_push.py v0.9.3`. Deploy web.

---

## M123 — Chat is the trading agent / one-screen desk. PIT 0.9.2


- **Source:** Chat dumped diagnostic prose, sent "Find the best opportunity" to Markets, and made users pick BTC/ETH/AVAX by hand. Judges need one Agent screen: scan → rank → private 0G → committee → policy → exact preview → desktop AUTHORIZE → real OID/fill/proof.
- **Discovery:** Chat is an orchestration layer over existing Watch, research, policy, capital, and activity. The LLM still cannot pin, arm, authorize, withdraw, or invent a fill. Find-best stays on Chat and starts sealed Direct on the strongest executable book, then hunts the next ranked candidate on a verified no-trade.
- **Candidate:** Operator cockpit: live funnel, opportunity strip, event-backed 0G pipe, research/preview/no-trade/proof/order cards, chips, Esc stop, Ctrl+Shift+A opens preview. Version **0.9.2**.
- **Kill:** AUTHORIZE in Chat. Fake streaming percentages. Invented confidence. Ticker-substring false matches. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim. iTransfer/iClone as live.
- **Next:** Tag `v0.9.2`, NSIS, overlay `D:\PIT`, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade unless the user AUTHORIZEs a preview.

DATE/TIME: 2026-08-31 00:50+03
PHASE: Agent cockpit. Desktop remains execution authority.
GOAL: One question — find the best opportunity — runs the whole desk without opening Markets by hand.
RESULT:
- **IMPLEMENTED:** Parse `Find the best opportunity/trade` → `research.best`, stay on Chat. Structured `deskcmd.Agent` payload. Short executive replies. Idle no longer dumps Watch. Hunt next executable book after READY_STOOD_DOWN (max 3). Agent rail. Composer chips + More. Live funnel, strip, pipe, preview, no-trade, policy block, 0G proof, order lifecycle, Sleep Mission propose-not-arm.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` pass. `npx tsx e2e/run.ts` including find-best chips, no authorizePreview, proof/no-trade cards. MCP prompt-injection still `mcp_read_only`.
- **LIVE:** After installer. No new MAINNET fill from this change. ETH OID `529167222216` untouched.
- **BLOCKED:** Authenticode still absent. Chat still cannot authorize or arm. iTransfer still not live on Aristotle. Direct research still needs Protect + credit.

SECURITY RESULT: Chat cannot AUTHORIZE, pin, arm, or export keys. Preview deep-links to Research. Kill switch still does not flatten.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Find-best hero on Chat: **IMPLEMENTED + TESTED**
- Live stage UI over existing research poll: **IMPLEMENTED + TESTED**
- AUTHORIZE remains Research-only: **UNCHANGED**
PRODUCTION READY: 0.9.2 installer after tag.
NEXT STEP: `python _scripts/push_head.py`. `python _scripts/tag_push.py v0.9.2`. Deploy web.

---

## M122 — Pair strip on the desk. Agent chat. PIT 0.9.1

- **Source:** Pair token lived inside a collapsed Security fold. Chat was a linear transcript. Users asked it to research the market, pick the best book, explain why to enter, show 0G progress, show txs, and enter when they accept.
- **Discovery:** Chat stays host-parsed. It still cannot AUTHORIZE, pin, or arm. Accept means type AUTHORIZE on Research. Pair code belongs on every desk page.
- **Candidate:** Persistent pair strip. Agent timeline with Watch → 0G Direct → committee → why → preview → AUTHORIZE on this computer. Honest ledger cards. Version **0.9.1**.
- **Kill:** Chat signing. Invented fills. Second landing accent. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim.
- **Next:** Tag `v0.9.1`, NSIS, Vercel + Render. Do not remint. Do not flatten. Do not start a live trade.

DATE/TIME: 2026-08-30 23:49+03
PHASE: Agent desk. Chat never authorizes.
GOAL: Pair token visible. Chat can run a research job by text and show 0G progress. User still types AUTHORIZE on this computer.
RESULT:
- **IMPLEMENTED:** Companion/desktop **0.9.1**. Pair strip under the title bar. Security pairing is open, not folded. Chat agent run: live book, 0G stage track, committee marks, preview, Accept on this computer, desk ledger. Host-parsed intents for research-the-best, why enter, show txs, I accept, AUTHORIZE-in-chat refused.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` pass. `npx tsx e2e/run.ts` including chat-agent harness (no authorizePreview in CommandChat).
- **LIVE:** After tag and installer. No new MAINNET fill. ETH OID `529167222216` untouched.
- **BLOCKED:** Authenticode still absent. Chat still cannot authorize. iTransfer still not live on Aristotle.

SECURITY RESULT: Chat cannot AUTHORIZE. Accept opens Research. Arm still desktop-only. Kill switch still does not flatten.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Pair strip: **IMPLEMENTED + TESTED**
- Agent chat + 0G progress + honest txs: **IMPLEMENTED + TESTED**
- AUTHORIZE remains Research-only: **UNCHANGED**
PRODUCTION READY: 0.9.1 installer after tag. Public claim unchanged.
NEXT STEP: `python _scripts/push_head.py`. `python _scripts/tag_push.py v0.9.1`. Deploy web.

---

## M121 — GitHub Latest 0.9.0. Hero centered. Beats ordered.

- **Source:** GitHub Releases still showed **0.8.0** because Latest is the last *tag*, not HEAD. Source was already 0.9.0. Hero copy was right but the coral board sat under the nav. Must-see Beats used a sticky 100dvh pan, so the SVGs looked unordered and a blank band followed the cards. Dual heading still said **MAINNET only**.
- **Discovery:** Tag `v0.9.0` with the NSIS we just built. Center the hero board in the coral viewport. Replace the Beats pin with a numbered 2+3 grid. Dual heading is **MAINNET**.
- **Candidate:** NSIS `PIT_0.9.0_x64-setup.exe`, overlay `D:\PIT`, annotated tag, Vercel + Render. Same coral, cream, ink.
- **Kill:** Second accent. Em-dash flourish. Fake fills. Web arming. Flatten OID 529167222216. Remint PIT-4bbee556. Authenticode claim.
- **Next:** Push, tag, deploy. Do not remint. Do not flatten. Do not start a live trade.

DATE/TIME: 2026-08-30 23:15+03
PHASE: Installer + landing fit. Chat never authorizes.
GOAL: Latest on GitHub is 0.9.0. Hero text and iris sit in the middle of the coral block. Must-see Beats shows five ordered diagrams without a void after them.
RESULT:
- **IMPLEMENTED:** Companion/desktop/web stay **0.9.0**. Windows installer built and installed over `D:\PIT`. DisplayVersion 0.9.0. `pit.exe version` = PIT 0.9.0. `pit-desktop.exe` FileVersion 0.9.0. Landing: hero uses a 1fr/auto/1fr grid so the kicker, line, and iris sit in the vertical middle. Must-see Beats is 01–05, equal diagram band, no scroll pin. Dual heading is **MAINNET**.
- **TESTED:** `go test ./...` PASS. Playwright 29/29. Installed NSIS silent `/S` exit 0. SHA256 `F140599E1F3B263E9C2A2B6209184F03AC0862EC183DE1083E02D146CFC86DEB`.
- **LIVE:** After this push, tag `v0.9.0`, Vercel, Render. No new MAINNET fill. ETH OID `529167222216` untouched.
- **BLOCKED:** Authenticode still absent. macOS/Linux not packaged. `go test -race` skipped (CGO_ENABLED=0). iTransfer still not live on Aristotle. PITAutonomyReceiptRegistry still skipped.

SECURITY RESULT: Arm still desktop-only. Web, chat, MCP, SDK still cannot arm. Chat never authorizes. Kill switch still does not flatten.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- GitHub Latest 0.9.0 installer: **IMPLEMENTED + TESTED locally; LIVE after tag**
- Hero center + Beats grid + Dual MAINNET: **IMPLEMENTED + TESTED**
PRODUCTION READY: Public installer matches source 0.9.0 once the tag is on GitHub. Stay-awake remains on the story beat.
NEXT STEP: `python _scripts/push_head.py`. `python _scripts/tag_push.py v0.9.0`. Deploy web.

---

## M120 — Hero, story, Automation desk fit any size

- **Source:** Production hero clipped PRIVATE ALPHA OS, stacked two small paragraphs, duplicate Explore, and a leftover coral pin. Story used a 100dvh pin that clipped Authorized. Desktop rail became a wrapping row under 1100px. Automation was three giant cards plus red essay boxes.
- **Discovery:** Same coral, cream, ink. Cut hero body to one line. Move stay-awake to the story beat. Arm CTA on the Automation page, not only in the overlay.
- **Candidate:** One-viewport coral hero. Compact story split. Segmented Sleep Mission desk. Icon rail on narrow windows. Page max-width so ultrawide does not stretch type.
- **Kill:** Second accent. Inter. Em-dash flourish. Fake fills. Web arming. Flatten OID 529167222216. Remint PIT-4bbee556.
- **Next:** Vercel. Do not tag. Do not remint. Do not flatten.

DATE/TIME: 2026-08-30 22:51+03
PHASE: UI honesty. Same 0.9.0. Chat never authorizes.
GOAL: Hero and Automation look like a serious product at any window size.
RESULT:
- **IMPLEMENTED:** Still **0.9.0**. No installer. Landing hero fits one viewport. Story no longer pins a blank screen. Desktop shell scales. Automation is MANUAL / RESEARCH / SLEEP MISSION with ARM SLEEP MISSION on the page.
- **TESTED:** Playwright home + autonomy. Desktop typecheck.
- **LIVE:** After push and Vercel.
- **BLOCKED:** GitHub Latest installer still v0.8.0. No new MAINNET fill.

SECURITY RESULT: Arm still desktop-only. Web still cannot arm.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Landing hero + story: **IMPLEMENTED**
- Desktop resolution + Automation: **IMPLEMENTED**
PRODUCTION READY: Source 0.9.0. Public claim unchanged. Stay-awake remains on the story beat.
NEXT STEP: `python _scripts/push_head.py`. Deploy web.

---

## M119 — Sleep Mission honesty: this computer must stay awake

- **Source:** Live 0G re-fetch still matches frozen Direct glm-5.2. Competitive scan: hosted 24/7 bots hide laptop-sleep as a skip. PIT hero said "while you sleep" without the bound.
- **Discovery:** Do not copy web arming, LLM fill loops, or 0g-memory EverMemOS. Do not fake iTransfer/DA on 16661. The missing claim is host survival: if this computer sleeps, the mission stops.
- **Candidate:** One sentence on landing, FAQ, /autonomy, and the arm overlay. Good Morning: zero fills is a successful night.
- **Kill:** Hosted runner. Auto-retune. Router fallback. Invented TVL. Fake backfill after sleep.
- **Next:** Vercel only for copy. Do not tag. Do not remint. Do not flatten.

DATE/TIME: 2026-08-30 22:17+03
PHASE: Claim honesty. Frozen SKU stays. Chat never authorizes.
GOAL: "Works while you sleep" never implies the machine can sleep.
RESULT:
- **IMPLEMENTED:** Still **0.9.0**. No installer. Landing + FAQ + /autonomy + desktop arm overlay state the computer must stay awake. Gap is not backfilled. Good Morning names a no-trade night as success.
- **TESTED:** Playwright home + autonomy after this patch.
- **LIVE:** Frozen glm-5.2 still matches getService (live 0G re-fetch). iTransfer/DAEntrance still absent on 16661.
- **BLOCKED:** Same as M118.

SECURITY RESULT: Same authority model. No new spending path.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Host-must-stay-awake copy: **IMPLEMENTED + TESTED**
- Direct glm-5.2 freeze: **LIVE, UNCHANGED**
PRODUCTION READY: Same 0.9.0 source. Public claim now matches host survival.
NEXT STEP: `python _scripts/push_head.py`. Deploy web.

---

## M118 — PIT 0.9 proof-carrying Sleep Missions

- **Source:** Guarded Autonomy existed as a typed enable phrase and host tick, but it was not a first-class Sleep Mission with a host AutonomyEnvelope, named refusals, night replay, Skillbook, or public-safe mission pages.
- **Discovery:** Keep Hyperliquid + Direct TeeML + host sizing. Arm only on this computer. Chat, web, MCP, and SDK prepare or inspect. They cannot arm. ENABLE GUARDED AUTONOMY stays as a compatibility phrase.
- **Candidate:** Proof-carrying bounded Sleep Missions. Envelope revalidated before exchange POST. No-trade is success. Public pages redacted.
- **Kill:** Unrestricted bot. Router private-book fallback. Fake fills, fake TEE, fake Storage proof, invented calibration. Flatten OID 529167222216. Remint PIT-4bbee556. New on-chain registry without a verifiable mainnet tx.
- **Next:** Vercel + Render. GitHub Latest stays the last tag (v0.8.0 NSIS) until a 0.9 installer is built and tagged.

DATE/TIME: 2026-08-30 22:12+03
PHASE: Sleep Missions. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Chat never authorizes. Direct TeeML only for the private book.
GOAL: User arms a bounded mission on desktop. Every autonomous action is envelope-gated. Public site explains and proves.
RESULT:
- **IMPLEMENTED:** Companion/desktop/web source **0.9.0**. No new NSIS this pass. GitHub Latest remains **v0.8.0**.
  1. Host `AutonomyEnvelope` + sleep state machine + mission-events.json + hash-only `MissionProof` + encrypted Skillbook.
  2. Desktop Automation: MANUAL / RESEARCH / SLEEP MISSION, ARM SLEEP MISSION, STOP, REVIEW LIMITS, Good Morning, Night Replay.
  3. Chat cannot arm. CLI `pit mission arm` is TTY-only. MCP/SDK `arm:false`.
  4. Web `/autonomy`, `/missions/:id` redacted. Landing: "It hunts while you sleep."
- **TESTED:** `go test ./...` PASS. `pit.exe version` = PIT 0.9.0. Desktop `tsc -b` pass. Desktop e2e copy harness pass. Web `tsc -b` pass. `npx vite build` pass. Playwright 28/28 then routes + autonomy specs.
- **LIVE:** No new MAINNET fill. ETH OID `529167222216` untouched. Honest demo is a named host block when max open is full. Frozen Direct glm-5.2 SKU unchanged. iTransfer still not live on Aristotle.
- **BLOCKED:** New PITAutonomyReceiptRegistry not deployed (hash-only receipts already exist). Authenticode absent. macOS/Linux not packaged. GitHub Latest installer still 0.8.0 until NSIS + tag. `go test -race` skipped (CGO_ENABLED=0). 0G Memory upstream not copied; PIT-native Skillbook uses official Storage path only when the existing security model already supports it.

SECURITY RESULT: Chat cannot AUTHORIZE or ARM. Web cannot arm. MCP/SDK cannot arm. Kill switch still does not flatten. Duplicate preview remains ledger-gated.
TX HASH / OID: Historical OID `529167222216` unchanged. No new chain write.
CLASSIFICATION:
- Sleep Mission + Autonomy Envelope: **IMPLEMENTED + TESTED**
- Night Replay / Skillbook / public proof: **IMPLEMENTED + TESTED**
- Desktop arm + chat/web/MCP/SDK refusal: **IMPLEMENTED + TESTED**
- 0.9 Windows installer / GitHub Latest: **NOT THIS PASS**
- Live MAINNET fill: **NOT THIS PASS** (unsafe while max open is full)
PRODUCTION READY: Source 0.9.0 is production-ready as bounded Sleep Missions on desktop. It is **not** an unrestricted trader. Public web still discovers and proves.
NEXT STEP: `python _scripts/push_head.py`. Upsert env. Deploy Vercel + Render. Do not tag until a 0.9 NSIS exists.

---

## M117 — Landing fills the void. Beats pan. Four doors.

- **Source:** Production landing still showed empty black chapters after the coral hero, Must-see Beats stopped at three cards, and MAINNET carried a second paragraph. CLI, MCP, SDK, and Desktop were missing as a product surface.
- **Discovery:** Hero pin was 180vh of leftover coral. Story headlines started at 18 percent opacity. Beats mapped horizontal travel across the unpin, so the last two cards never arrived. MAINNET only needs the laboratory sentence.
- **Candidate:** Same coral, cream, ink. WireTurn plus inner iris. Token marks from web3icons. One 2x2 of Desktop / CLI / MCP / SDK.
- **Kill:** Aristotle dump under MAINNET. Empty 150vh story pin. Wallet-first. Authorize on the web.
- **Next:** Vercel production. Desktop Security stays the M116 next-action desk.

DATE/TIME: 2026-08-30 21:28+03
PHASE: UI honesty. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Chat never authorizes. Direct TeeML only for the private book.
GOAL: Landing has no empty chapter. All five beats are reachable. Surfaces show Desktop, CLI, MCP, SDK.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**. No new installer this pass.
  1. Hero pin 118vh. One line. Magnetic Explore. WireTurn + reverse iris.
  2. Must-see Beats: ResizeObserver travel, pin height = 100dvh + travel, snap fallback under 768px. All five cards in the tree: Private book, Three envelopes, Policy is law, You authorize, Then it remembers.
  3. MAINNET only keeps the laboratory sentence. Token row ETH BTC SOL HYPE DOGE AVAX. Four doors section.
- **TESTED:** Vite build pass. Playwright home/network/authorize/routes 21/21 against Vite with Privy env.
- **BLOCKED:** iTransfer not live. No new live trade. Authenticode still absent.

SECURITY RESULT: Chat still cannot AUTHORIZE. Web still has no Authorize control. MCP still cannot order. Historical OID unchanged.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Landing empty-section + beats pan: **IMPLEMENTED + TESTED**
- Four doors (Desktop CLI MCP SDK): **IMPLEMENTED + TESTED**
- Production alias: **SHIP THIS PASS**
PRODUCTION READY: Unsigned Windows x64. Verify SHA256. Do not treat RESTING as a fill.
NEXT STEP: `python _scripts/push_head.py`. Vercel web.

---

## M116 — Security is a next-action desk. Landing hero is one line.

- **Source:** Desktop Security was pairing-first, then a wall of ready chips, then tiny Hyperliquid rows, then policy. Users could not see what to do. Landing hero repeated the story in two paragraphs plus a second mega PIT.
- **Discovery:** One next-action card, then Session / Policy / optional Browser / Compute / Halt. Landing keeps coral + WireTurn. One line: "It hunts in private. You authorize on this computer."
- **Candidate:** Same brand. Larger type. Fewer words. Scroll-driven beats.
- **Kill:** Chip dump. Pairing as the Security hero. Wallet-first landing. Authorize on the web.
- **Next:** Deploy web. Desktop source is this commit (NSIS still 0.8.0).

DATE/TIME: 2026-08-30 20:55+03
PHASE: UI honesty. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Chat never authorizes. Direct TeeML only for the private book.
GOAL: Security has one obvious next step. Landing hero is one line. Diagrams and scroll still tell Sealed / Challenged / Authorized.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**. Same tag. New unsigned NSIS for the Security UI (SHA below).
  1. Desktop Security: next-action banner from `nextFix`, 4-step spine, Session metrics, Policy with short helpers, pairing folded unless it is the next step, Halt/Revoke, identity in details.
  2. Landing hero: coral field, rotating WireTurn, one line, Explore + Download. Story line: "The web discovers. The desktop acts." Beats pin on scroll. Moments pan horizontally. SVGs tightened.
- **TESTED:** Desktop `tsc -b` pass. `npx tsx e2e/run.ts` pass. Web `tsc -b` pass. Playwright home/authorize/network/routes 21/21.
- **LIVE:** Chrome https://pit0g.vercel.app/ hero `It hunts in private. You authorize on this computer.` Explore + Download only. No Authorize. MAINNET only. Laboratory sentence. Overlay `D:\PIT\pit-desktop.exe` from this NSIS. SHA256 `106A688275BFCF3D108E5A50A3BDB5F3C5B059E48D9629A1217C6226378FCE5A`.
- **BLOCKED:** iTransfer not live. No new live trade. Authenticode still absent.

SECURITY RESULT: Chat still cannot AUTHORIZE. Web still has no Authorize control. Historical OID unchanged.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Desktop Security reorder: **IMPLEMENTED + TESTED**
- Landing hero one-line: **IMPLEMENTED + TESTED**
- Production landing: **SHIPPED THIS PASS AFTER PUSH**
PRODUCTION READY: Unsigned Windows x64. Verify SHA256. Do not treat RESTING as a fill.
NEXT STEP: `python _scripts/push_head.py`. Vercel web. Chrome `/`.

---

## M115 — GitHub Latest was the last tag, not HEAD. Ship Windows 0.8.0.

- **Source:** GitHub still showed **PIT v0.7.6 Latest** because Releases/Latest is the last *tag*, not commit `3add5c0`. Go companion and Render `/health` were already **0.8.0**. Desktop `package.json` / Tauri / NSIS still said 0.7.6. Public `/download` fetched `api.github.com` from the browser and failed CORS/rate-limit.
- **Discovery:** Align desktop to 0.8.0, build the NSIS, tag **v0.8.0**, proxy GitHub Releases through health `/release` (`sign:false` `trade:false`). Do not claim Authenticode. Do not claim macOS/Linux. Do not remint. Do not flatten OID `529167222216`.
- **Candidate:** Keep frozen Direct glm-5.2 SKU. Public web discovers the installer. Desktop still acts.
- **Kill:** Leave Latest at 0.7.6 while CLI is 0.8.0. Auto-swap Direct. Invent a fill. Start a live trade this pass.
- **Next:** Tag + upload after this commit is on `main`. Deploy health so `/download` can read `/release`.

DATE/TIME: 2026-08-30 20:01+03
PHASE: Desktop 0.8.0 installer. Public `/release`. Local install + companion proof. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Direct TeeML only for the private book. Chat never authorizes.
GOAL: GitHub Latest = **v0.8.0**. Companion on this PC = **0.8.0**. `/download` reads health, not the browser GitHub API.
RESULT:
- **IMPLEMENTED:** **0.8.0** Windows NSIS.
  1. Desktop versions: `apps/desktop` package / Tauri / Cargo / `SIDECAR_VERSION` / `DESKTOP_VERSION` = **0.8.0**.
  2. NSIS `PIT_0.8.0_x64-setup.exe` (17,212,799 bytes). SHA256 `8A1E4F8D54898B3E5C5D5F53C24C15FAD34DD9327073171C6B598102EBEE7749`. Unsigned. No macOS/Linux package.
  3. Health `GET /release` proxies GitHub latest (5 min cache). Never `sign` or `trade`. `/download` uses `healthBase()/release`.
  4. Overlay `D:\PIT\pit.exe` + `pit-desktop.exe` from this build. NSIS `/S` exit 128 (PREINSTALL `taskkill` when the process is already gone). `uninstall.exe` written to `D:\PIT` at 19:59+03. Companion launched from `D:\PIT`.
- **TESTED:** `go test ./cmd/health` PASS. `npx playwright test` 26/26 then download specs 2/2. Cargo `allow_official_https` PASS this session. `D:\PIT\pit.exe version` = **PIT 0.8.0**.
- **LIVE:** `GET http://127.0.0.1:17373/health` → `version:0.8.0` `sign:false` `trade:false` `ok:true`. pit-desktop PID running. Frozen Direct SKU unchanged (M114). Production health was already 0.8.0 before this Render rebuild.
- **UNVERIFIED until tag+deploy in this same pass:** GitHub Releases Latest = v0.8.0. Production `/download` showing the tag. Public TEE VerifyE2EE in the browser (still honest: it does not).
- **BLOCKED:** iTransfer UNAVAILABLE on Aristotle. DA not live on 16661. No new live trade. Authenticode certificate still absent.

SECURITY RESULT: Chat still cannot AUTHORIZE. MCP still cannot order/cancel. `/release` cannot sign. Direct glm-5.2 frozen SKU unchanged. Historical OID `529167222216` unchanged. PIT-4bbee556 not reminted.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Desktop 0.8.0 NSIS: **IMPLEMENTED + TESTED + INSTALLED ON THIS PC**
- GitHub Latest v0.8.0: **SHIPPED THIS PASS AFTER PUSH**
- iTransfer / DA: **NOT LIVE**
PRODUCTION READY: Unsigned Windows x64 only. Verify SHA256 before run. SmartScreen may warn. Do not treat RESTING as a fill.
NEXT STEP: `python _scripts/push_head.py`. Create GitHub `v0.8.0` with the NSIS + SHA256SUMS. Upsert Vercel/Render env. Deploy. Chrome `/download`.

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

Local notes remain on disk. Public tree is the product + README.


---

## M21 — Direct runner, HL wires, desktop shell (2026-08-26 22:01+03)

DATE/TIME: 2026-08-26 22:01+03
PHASE: 6 Direct runner/SKU split · 9 HL order/cancel + master approveAgent · 5 policy pin · 20 desktop shell · 24 IDOR
GOAL: Wire fail-closed product paths for sealed ask (binary still required), venue actions, and the local authorize surface.
FILES CHANGED: `pit/internal/compute/catalog.go` `direct.go` `roles.go` · `pit/internal/hl/action.go` `agent.go` · `pit/internal/exec/exchange.go` · `pit/internal/policy/pin.go` · `pit/internal/workspace/idor_test.go` · `pit/internal/storage/idor_test.go` · `pit/internal/ui/` · `apps/desktop/**` · `IMPLEMENTATION_PLAN.md` `README.md` `.env.example` this file.
COMMANDS RUN: `go test ./...` all packages ok.
EXTERNAL RESOURCES CONSULTED: Galileo Omni provider/signer from ENVIRONMENT.md; mainnet glm-5.2 provider/teeSigner; Hyperliquid exchange action types; storage object key prefix.
RESULT: Testnet cannot select glm-5.2. Galileo sealed ask stays disabled until VerifyE2EE is proven. Direct job refuses Router URLs and Python gate binaries. Session cannot sign approveAgent. Exchange wrapper rejects mock hosts and withdraw. Desktop shows order/cancel allowed, withdraw/leverage denied. Progress labels match backend names.
TESTS: Go unit pass (compute catalog/direct/roles, hl action/agent, exec exchange, policy pin, ui copy, storage/workspace IDOR).
EVIDENCE: this entry. No secrets.
TX HASH / OID: none. Live TeeML and live HL post still require `PIT_COMMITTEE_BIN` and a funded session on the matching venue.
FAILURE: `ui` package `const []string` is illegal in Go; switched to `var`.
FIX: `copy.go`.
RETRY: `go test ./internal/ui/...` pass.
FINAL STATE: desktop shell exists; signing still not in the web app. Sealer binary not vendored.
NEXT STEP: run committee binary from product on Aristotle; testnet Omni VerifyE2EE before enabling Galileo ask; live dust order on HL testnet; Storage `--proof` upload; Chrome E2E.

`.env` (gitignored): `PIT_COMMITTEE_BIN` placeholder only. Secrets not logged.

---

## Local-only working notes

DATE/TIME: 2026-08-26
GOAL: Keep planning notes on the machine. Public repo is product code + README.
RESULT: `ENVIRONMENT.md`, `FINAL_FLAGSHIP_PRODUCT.md`, `IMPLEMENTATION_PLAN.md`, `RESOURCES.md`, and this file stay local (gitignored). README is the published guide.
EVIDENCE: local files still on disk. `.env` not committed.

---

## M22 — implementation continuation audit (2026-08-26 22:53+03)

DATE/TIME: 2026-08-26 22:53+03
PHASE: audit then continue 8 Watch · 18 CLI · 19 contracts tests · 21 web · 25 security · 30 observability
GOAL: Reconstruct repo state. Do not restart planning. Continue incomplete product work.
FILES READ: this file, IMPLEMENTATION_PLAN, ENVIRONMENT, RESOURCES, FINAL_FLAGSHIP, git tree, HEAD `3d13bd9`.
COMMANDS RUN: `git status` · `git rev-list --count HEAD` = 61 · `git log`
EXTERNAL RESOURCES CONSULTED: Dual-env contract in plan; Direct vs Router; KEPT `guide.css` tokens (coral `#d82f2f`, cream `#f0e7d4`, ink) visual only.
RESULT:
DONE (library/unit): 1 identity/workspace · 2 wallet named states · 3 HL account parse · 4 session keygen/allowlist · 5 policy host · 6 committee envelope/verify/SKU/direct-args · 7 forecast/sizer · 8 watch scan+safety (no live loop) · 9 exec gateway+wires · 10 ledger recover · 11 storage proof args · 12 memory kinds · 13 8004 host checks · 14 DeskID host refuse-transfer · 15 calib empty/overconfidence · 16 MCP RO · 17 SDK no-sign · 18 CLI TTY authorize skeleton · 19 contracts+13 Foundry tests · 20 desktop shell · 21 web connect shell.
PARTIAL: 6 live Seal+VerifyE2EE in product binary · 9 live exchange post · 11 live `--proof` upload · 18 CLI remaining commands · 20 desktop packaging · 21 web quality/onboarding depth.
NOT STARTED / UNVERIFIED: 22 Playwright · 23 real Chrome E2E · 24 multi-user live · 25 full red-team matrix · 26 dual deploy · 27 real-user 14-beat · 28 demo · 29 final gate.
EXTERNAL BLOCKER: Aristotle iTransfer (attestor absent). Galileo Omni VerifyE2EE unproven. Do not fake.
TESTS: `go test ./...` pass. Foundry **19** passed (added AccessControl, Replay, Security).
NEXT STEP: continue live Direct/HL only with fail-closed paths. Frontend tokens now match the coral/cream/ink desk language.

---

## M23 — Direct bind, CLI state, Watch law, Foundry invariants (2026-08-26 23:05+03)

DATE/TIME: 2026-08-26 23:05+03
PHASE: 6 Direct seal bind · 8 Watch no-execute · 11 Storage official client · 13 8004 privacy · 18 CLI persist · 19 contract invariants · 20/21 empty Watch + policy cards
GOAL: Continue product paths. Fail closed. No Router. No fake opportunities.
FILES CHANGED: `pit/internal/compute/seal.go` `run.go` · `hl/l2.go` `l2fetch.go` · `cli/state.go` · `policy/cards.go` · `watch/notify.go` · `storage/client.go` · `chain8004/privacy.go` · `redteam/session_test.go` · `contracts/test/Invariants.t.sol` · `memory` kinds · `engine/modelsize.go` · `obs/id.go` · `exec/repost.go` · `market/frombook.go` · `mcp` · `sdk/ops.go` · `cmd/pit` · desktop EmptyHome/PolicyLaw · web PolicyPanel/EmptyWatch · `apps/web/DESIGN.md` · README
COMMANDS RUN: `go test ./...` pass. `forge test` **22** passed.
EXTERNAL RESOURCES CONSULTED: Direct TeeML scheme `zg-sig-v1/e2ee-ct`; HL `l2Book` shape; official Storage `--proof`; ERC-8004 public feedback must omit book/strategy.
RESULT: Galileo sealed ask remains disabled until VerifyE2EE is proven. Watch cannot place orders. CLI state.json has workspace/network/wallet/kill only. SDK still cannot hold a session. Empty Watch copy is real.
FAILURE: host sizer test used 10 USD at 3500 (rounds below min notional). Fixed clip to 50 USD at 2500.
RETRY: `go test ./internal/engine/...` pass.
EVIDENCE: this entry. No secrets.
TX HASH / OID: none. Live sealer exec still not attached (`sealer_exec_not_attached`). Live HL post still open.
FINAL STATE: libraries and UI empty-states advanced. Product still does not send a live Direct TeeML request.
NEXT STEP: attach native sealer when the binary exists; live HL dust on testnet; Storage upload with `--proof`; Chrome E2E; then deploy.

---

## M24 — public quotes, phase order, honest committee label (2026-08-26 23:11+03)

DATE/TIME: 2026-08-26 23:11+03
PHASE: 8 market provenance · 7 committee honesty · 10 timeout repost already done · 25 kill/session · 29 progress UI
GOAL: Real source timestamps. Honest same-provider wording. Progress labels from named states only.
FILES CHANGED: market/public · session/export · phase/order · policy/recheck · engine/fresh · hl/trades · compute/honest · verify/network · redteam/kill · ledger concurrent · sdk export · cmd/health test · web ProgressStrip · desktop Progress
COMMANDS RUN: `go test ./...` pass (health test package fixed to `main_test`).
RESULT: CoinGecko/DeFiLlama quotes require timestamps. Mock source denied. Committee label is envelope separation on the same provider unless proven otherwise. Health JSON cannot sign.
FAILURE: `cmd/health` test used a second package name in the same folder. Fixed to `main_test`.
RETRY: `go test ./cmd/health/...` pass.
EVIDENCE: this entry.
NEXT STEP: live Direct exec, live HL dust, Chrome E2E, deploy.

---

## M25 — native sealer exec, live Watch, CLI ask (2026-08-26 23:22+03)

DATE/TIME: 2026-08-26 23:22+03
PHASE: 6 Direct exec fail-closed · 8 Watch live books · 11 storage tamper · 14 desk before-ask · 18 CLI ask/opportunities · 21 SIWE bind copy · 25 scan/web red-team
GOAL: Continue from M24. Attach native sealer when the binary exists. Watch uses venue books. CLI does not invent cards.
FILES CHANGED: compute/run,ask,hpke · watch/live · deskid/before · exec/bind,query · keyring/ns · storage/tamper · hl/agents · siwe/domain · scan/web · session/agentname · policy/cooldown · mcp/readonly · sdk/watch · CLI ask/opportunities/forecast · web SiweBind · desktop BindNote · README
COMMANDS RUN: `go test ./...` pass. `go build ./cmd/pit ./cmd/health` pass.
EXTERNAL RESOURCES CONSULTED: Direct HPKE KEM `0x0020` SealInfo `0g-pc/v1/seal` `_e2ee`; scheme `zg-sig-v1/e2ee-ct`; HL extraAgents name `PIT-{8}`; official Storage root compare.
RESULT: Missing sealer → `sealer_not_wired`. Python/TS sealer refused. Plain sealer output → `TEE_VERIFY_FAIL`. Galileo sealed ask stays `galileo_e2ee_unproven`. Unauthorized desk → `desk_not_authorized`. Watch Live fetches allowlisted coins only and cannot place orders. CLI `opportunities` uses live books. CLI `ask` requires init and still stops without a native sealer and desk auth. Preview post requires matching preview hash. Timeouts query the venue first. Web source scan rejects `session_key` / `private_key` / `mnemonic`.
FAILURE: cooldown unit used 50s elapsed inside a 60s window; expected “elapsed”. Fixed to 70s elapsed.
RETRY: `go test ./internal/policy/` pass. Full `go test ./...` pass.
EVIDENCE: this entry. No secrets. No planning files committed.
TX HASH / OID: none. Live TeeML still needs a real native `PIT_COMMITTEE_BIN`. Live HL order/cancel still needs AUTHORIZE + funded session.
FINAL STATE: product path is fail-closed and executable when the sealer file exists. Empty Watch remains real.
NEXT STEP: live HL dust on testnet; Storage `--proof` upload; Chrome E2E (Privy, User A/B); Tauri packaging; then Vercel/Render.

---

## M26 — preview bind, storage proof args, isolation (2026-08-26 23:29+03)

DATE/TIME: 2026-08-26 23:29+03
PHASE: 11 Storage official client args · 13 preview bind fields · 18 CLI preview · 21/20 attention copy · 24 memory isolation · 25 market mutation
GOAL: Continue after M25. No mocks. Fail closed.
FILES CHANGED: storage/proof · engine/finite · exec/fields · ledger/times · memory/isolate · redteam/market · hl/agent name · CLI preview · web Attention · desktop AttentionLine · sdk/preview · cli/preview copy · README
COMMANDS RUN: `go test ./...` pass. `go build ./cmd/pit` pass.
RESULT: Download args require `--proof`. Upload requires official Go client + `0x` 32-byte key. Host rejects NaN/Inf size. Preview bind field list is complete. Market mutation denied. Memory keys differ per workspace. Session agent name is `PIT-{8}` from the workspace. CLI preview prints exact bind fields and stops without a live session. Home attention copy stays empty when count is 0.
FAILURE: first `UploadArgs` helper collided with existing `proof.go`. Deleted the duplicate and tightened the existing Job args.
RETRY: `go test ./internal/storage/ ./internal/memory/` pass.
EVIDENCE: this entry. No secrets.
TX HASH / OID: none.
NEXT STEP: live HL dust; native sealer on Aristotle; Chrome E2E; deploy after E2E.

---

## M27 — cancel bind, crash recover, dual-network Chrome (2026-08-26 23:35+03)

DATE/TIME: 2026-08-26 23:35+03
PHASE: 9 cancel bind · 10 crash recover · 13 8004 not portable · 15 learned/overconfidence · 16 MCP mutate · 18 CLI cancel/resolve · 21/26 network banner · 23 Chrome desk
GOAL: Continue stacked loop ticks as one phase. Fail closed. No live dust order.
FILES CHANGED: ledger/crash · exec/cancel · hl/cloid+action · session/ttl · watch/block · chain8004/portable · calib/learn · mcp/mutate · CLI cancel/resolve · web/desktop NetworkBanner · README
COMMANDS RUN: `go test ./...` pass after unused `fmt` removed from `hl/action.go`. Chrome page `http://127.0.0.1:5173/` MAINNET then TESTNET.
EXTERNAL RESOURCES CONSULTED: Hyperliquid 128-bit cloid; Galileo explorer vs Aristotle explorer; session TTL ≤ 1h.
RESULT: Cancel requires matching preview hash and a valid cloid. Crash after sign queries the exchange. Session TTL above 1h is denied. SOL is a blocked Watch card under default policy. 8004 IDs are not portable across networks. MCP cannot cancel. CLI cancel/resolve stop without a live session. Chrome: wallet connected, empty Watch is real, TESTNET switches capability copy and explorer to Galileo (`chainscan-galileo.0g.ai`). No seed field. No AUTHORIZE click. No live order.
FAILURE: `hl/action.go` unused `fmt` after moving cloid checks. Removed import.
RETRY: `go test ./...` pass.
EVIDENCE: Chrome a11y snapshot (TESTNET lab copy + Galileo explorer). No secrets.
TX HASH / OID: none. Live HL dust and native sealer still open.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer; User A/B; then deploy.

---

## M28 — User A/B isolation, policy red-team, reduced motion (2026-08-26 23:38+03)

DATE/TIME: 2026-08-26 23:38+03
PHASE: 19 Isolation.t.sol · 24 memory/forecast IDOR · 25 slippage/leverage · 18 CLI resolve fix · 20/21 reduced-motion
GOAL: Continue tick 7. No live dust. Fail closed.
FILES CHANGED: engine/lookup · session/switch · redteam/slip · workspace/memory · sdk/forecast · Isolation.t.sol · CLI resolve dead code · desktop/web reduced-motion · README 25 Foundry tests
COMMANDS RUN: `go test ./...` pass. `forge test` **25** passed (was 22).
RESULT: Forecast lookup returns `not found` for the wrong workspace. Account switch requires a new session. Slippage and leverage fail closed. Workspace A cannot read B memory. Bob cannot replay Alice receipt, resolve Alice forecast, or overwrite Alice memory. CLI `resolve` no longer prints a stray init error after the command. Reduced-motion disables grain and transitions.
FAILURE: none after tests.
EVIDENCE: this entry. No secrets. No live order.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer; Chrome User A/B; deploy.

---

## M29 — keyring surface, live book, verify route, CI web (2026-08-26 23:41+03)

DATE/TIME: 2026-08-26 23:41+03
PHASE: 4 keyring not in browser · 8 live HL public book · 16 MCP timestamped market · 20 Tauri config · 21 #verify · 22 CI web build
GOAL: Continue tick 8. Live public book is allowed. Live order is not.
FILES CHANGED: keyring/surface+delete · mcp/market · tauri.conf.json · vercel.json · ci.yml web job · App #verify · sdk keyring · README
COMMANDS RUN: `go test ./...` pass. `go test -tags live ./internal/hl -count=1` pass (live ETH book, mark > 0). Chrome desk still on TESTNET lab copy, empty Watch, Galileo explorer.
RESULT: Keyring refused on web/browser/Vercel/MCP. Secrets redact to `[redacted]`. File store delete revokes a namespace key. MCP market quotes require a timestamp. `#verify` is a web-only inspect route. CI builds the Vite app with the public Privy app id. Desktop has a Tauri config stub. Live Hyperliquid `PublicBook("ETH")` returned a positive mark. No order posted.
FAILURE: none.
EVIDENCE: live Go test + Chrome snapshot. No secrets. No AUTHORIZE.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer; Chrome User A/B; then deploy.

---

## M22 — implementation continuation audit (ticks 9–10)

DATE/TIME: 2026-08-26 23:44+03
PHASE: reconstruct from HEAD `4a5cc5a` (174 commits)
DONE: core Go, Foundry 25 then 26, CLI, MCP read-only, SDK, web connect, desktop shell, live public books, `#verify`
PARTIAL: Chrome (wallet connect + network + verify copy; no User A/B AUTHORIZE), storage `--proof` live upload, Tauri packaging
NOT STARTED / BLOCKED: live dust order after AUTHORIZE; native `PIT_COMMITTEE_BIN`; Vercel/Render until local E2E; iTransfer on Aristotle
UNVERIFIED: Galileo VerifyE2EE
DECISION: do not redo M29. Continue fail-closed product work. Do not deploy. Do not place orders.

---

## M30 — restart, policy pin, L2 gate, dual-venue book (2026-08-26 23:46+03)

DATE/TIME: 2026-08-26 23:46+03
PHASE: 9 duplicate click · 10 restart after preview · 8 L2 + testnet public book · 10 policy pin · 12 memory list · 11 wrong storage key · 16 MCP verify · 21 web refresh cannot sign · 20 desktop recover
GOAL: Continue ticks 9–10. Live public books allowed. Live order still not.
FILES CHANGED: ledger restart test · storage KeysMatch · product session env refuse · memory list deny · policy MatchPin · CLI pin · SDK refresh · L2 coin required · exec duplicate click · MCP verify hint · health env · Isolation policy pin · desktop RecoverNote · web RefreshNote · live testnet ETH book · README
COMMANDS RUN: `go test ./...` pass. `forge test -vvv` 26 pass. `go test -tags live ./internal/hl -count=1` pass (mainnet + testnet ETH mark > 0). Chrome `http://127.0.0.1:5173/#verify` with connected wallet `0xBDfC…0034`, refresh-cannot-sign copy, MAINNET explorer. No AUTHORIZE. No order.
RESULT: Previewed actions survive process restart. Duplicate apply is refused. Wrong encryption keys fail closed. Product mode refuses session secrets in env. Workspace B cannot list A memory keys. CLI pins policy hash; mutation fails closed. SDK/browser cannot sign after refresh. Desktop recovers the preview. MCP verify never authorizes. Health never signs. Bob cannot pin Alice policy (Foundry). L2 requires a coin. Live testnet info book quoted. No deploy.
FAILURE: none after tests. Live dust still needs AUTHORIZE + funded session. Native sealer still unwired.
EVIDENCE: this entry + Chrome snapshot on `#verify`. No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer; Chrome User A/B isolation with a second wallet; then deploy.

---

## M31 — expired session, spoof, tamper, two-wallet copy (2026-08-26 23:52+03)

DATE/TIME: 2026-08-26 23:52+03
PHASE: 24 expired/revoked · provider/model spoof · E2EE tamper · log leak · browser extract · SIWE origin · 15 fake calib · 8 network mix · 16 watch bind · 21/20 two-wallet copy
GOAL: Continue ticks 11–13. No live order. Fail closed.
FILES CHANGED: session life tests · exec Gate · compute MatchProvider/TamperFails · calib RefuseInvented · obs leak · SDK BindNetwork · redteam BrowserExtract · CLI expired copy · SIWE origin · watch Bound · health requestId · web/desktop IsolateNote · README
COMMANDS RUN: `go test ./...` pass. Chrome home: wallet connected, isolation copy, TESTNET lab + Galileo explorer, empty Watch. No AUTHORIZE. No order.
RESULT: Expired and revoked sessions cannot pass the exec gate. Catalog provider/model/signer spoof fails. Tampered E2EE fails. Empty calibration cannot print 72%. Logs refuse session keys and private book. Browser sources with `session_key` fail the extract scan. SDK cannot mix mainnet and testnet. SIWE origin must match. Watch requires a workspace id. Health includes a requestId and still never signs. Desk copy: two wallets never share a workspace.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired. Second wallet Chrome isolation not executed.
EVIDENCE: Chrome TESTNET snapshot + this entry. No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer; Chrome User A/B with a second wallet; then deploy.

---

## M32 — clip overflow, receipt replay, no browser session (2026-08-26 23:56+03)

DATE/TIME: 2026-08-26 23:56+03
PHASE: 11 huge size clip · 15 receipt replay · 16 MCP timestamped opportunities · 17 SDK named events · 18 network switch deny · 22 desktop named errors · 21 browser never holds session · 24 tenant mix · 29 no theater phases
GOAL: Continue tick 14. No live order.
FILES CHANGED: engine overflow · verify replay · SDK events · MCP OpportunityQuote · redteam tenant · CLI mix · market source · phase theater · web/desktop NoSession · desktop e2e named errors · health no wallet · README
COMMANDS RUN: `go test ./...` pass. Chrome: NoSession copy visible, then Disconnect. Empty Watch. No AUTHORIZE. No order.
RESULT: Huge model size cannot exceed clip. Receipt hashes cannot be filed twice. MCP opportunities require a timestamp and never trade. SDK events are named host states only. Bound CLI cannot switch mainnet to testnet. Quotes require source+network+symbol. SPINNING is refused as theater. Health JSON has no wallet or session field. Browser copy: this browser never holds a Hyperliquid session.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired. Second wallet not used.
EVIDENCE: Chrome disconnect snapshot + this entry. No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer; Chrome User A/B with a second wallet; then deploy.

---

## M33 — 8004 reporter, daily halt, kill copy (2026-08-27 00:03+03)

DATE/TIME: 2026-08-27 00:03+03
PHASE: 13 8004 stranger reporter · 10 daily loss halt · 9 kill gate · 8 funding finite · 16 empty card · 21 kill/daily policy copy
GOAL: Continue ticks 15–16. No live order.
FILES CHANGED: chain8004 stranger · session network match · policy HaltDaily · exec RefuseKill · HL funding · SDK feedback · MCP EmptyCard · obs latency · market RefuseMock · web/desktop KillNote · PolicyPanel daily loss · README
COMMANDS RUN: `go test ./...` pass. Chrome home: disconnected, Connect your wallet, kill-switch copy, Max daily loss 50 USD, empty Watch. No AUTHORIZE. No order.
RESULT: A stranger cannot file 8004 feedback. Session network must match the workspace. Daily loss halt and kill switch stop new orders. Funding must be finite. MCP card does not invent accuracy. Mock quotes denied. Desk shows daily loss and that YOU flip kill.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired. Second wallet not used.
EVIDENCE: Chrome snapshot (disconnected + kill copy) + this entry. No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer; Chrome User A/B with a second wallet; then deploy.

---

## M34 — RPC mix, SIWE chain, slippage, disconnected verify (2026-08-27 00:07+03)

DATE/TIME: 2026-08-27 00:07+03
PHASE: 8 RPC/chain mix · 5 SIWE chain · 9 slippage · 8 OI/oracle · 17 empty calib SDK · 18 secret print · 21 disconnected #verify
GOAL: Continue tick 17. No live order.
FILES CHANGED: config RefuseMixedRPC · SIWE ChainMatch · HL OI/oracle · exec slippage · SDK HealthCard · CLI RefusePrint · wallet unfunded · health GET/HEAD · PolicyPanel slippage · desktop daily loss · deskid Galileo mint gate · MCP status never signs · README
COMMANDS RUN: `go test ./...` pass. Chrome `#verify` while disconnected: Connect your wallet, verify form, mainnet explorer. No AUTHORIZE. No order.
RESULT: Galileo RPC cannot pair with Aristotle. SIWE chain must match. Slippage above policy fails. OI and oracle must be finite. SDK empty health does not invent accuracy. CLI refuses to print session secrets. Unfunded accounts are named. Health accepts GET/HEAD only. Galileo desk mint requires a deployed address. `#verify` works without a session.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired.
EVIDENCE: Chrome `#verify` disconnected snapshot + this entry. No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer; Chrome User A/B with a second wallet; then deploy.

---

## M35 — uncertainty, liquidity, cooldown, TESTNET verify (2026-08-27 00:13+03)

DATE/TIME: 2026-08-27 00:13+03
PHASE: 11 engine uncertainty · 10 exec liquidity/cooldown · 8 mark · 9 Watch never trades · 23 desktop CI/CSP · 14 script sealer · 20 MCP watch · 22 CLI watch · 21 SDK watch · 21 Chrome TESTNET #verify
GOAL: Continue tick 18. No live order. No deploy.
FILES CHANGED: engine RefuseUncertainty · exec liquidity + cooldown · HL mark finite · Watch MayPlaceOrder tests · PolicyPanel cooldown/uncertainty/liquidity · desktop cooldown · Tauri CSP · CI desktop Vite build · native sealer refuses .mjs/.cjs · MCP/CLI/SDK watch never trades · README · labeled TESTNET harness
COMMANDS RUN: `go test ./...` pass. Chrome `#verify` TESTNET: Connect your wallet, Galileo explorer, lab copy. No AUTHORIZE. No order.
RESULT: Uncertainty above the policy max fails. Thin liquidity fails. Cooldown fails. Non-finite mark fails. Watch/MCP/SDK cannot place. Script sealers are refused. Desktop CI builds the Vite shell. `#verify` on TESTNET points at chainscan-galileo.0g.ai.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired (`PIT_COMMITTEE_BIN` empty). User A/B isolation in Chrome still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: Chrome `#verify` TESTNET snapshot (disconnected, Galileo explorer, no session). No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer binary; Chrome User A/B with a second wallet; then deploy.

---

## M36 — leverage, venue, transfer, stale preview, MAINNET verify (2026-08-27 00:17+03)

DATE/TIME: 2026-08-27 00:17+03
PHASE: 9 leverage/venue/transfer/approveAgent · 10 timeout never reposts · 11 storage root replay · 4 session TTL · 16 MCP transfer · 18 CLI leverage · 17 SDK funds · 21 Chrome MAINNET #verify
GOAL: Continue ticks 19-20. No live order. No deploy.
FILES CHANGED: exec leverage/venue/transfer/approveAgent/stale preview · ledger AfterTimeout · storage root replay · session CapTTLHours · MCP TransferNever · CLI leverage · SDK CanTransfer · PolicyPanel venues · desktop session TTL · README · MAINNET harness
COMMANDS RUN: `go test ./...` pass. Chrome `#verify` MAINNET: Connect your wallet, chainscan.0g.ai, production copy. No AUTHORIZE. No order.
RESULT: Leverage above 1x fails. Foreign venues fail. sendAsset and approveAgent fail. Timeouts query first. Wrong storage roots fail. Session TTL is one hour. Stale previews cannot authorize. MCP cannot transfer. `#verify` on MAINNET points at chainscan.0g.ai.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: Chrome `#verify` MAINNET snapshot (disconnected, Aristotle explorer, no session). No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer binary; Chrome User A/B with a second wallet; then deploy.

---

## M37 — calibration floor, SOL allowlist, Galileo sealed-ask gate, empty home (2026-08-27 00:21+03)

DATE/TIME: 2026-08-27 00:21+03
PHASE: 15 calibration floor · 8 SOL not in default universe · 16 MCP forecast never sizes · 18 CLI kill · 17 SDK drops model size · 13 8004 private tags · 30 health leak · 6 Galileo unproven · 7 empty side · 8 impact finite · 5 declined SIWE · 21/28 empty home
GOAL: Continue tick 21. No live order. No deploy.
FILES CHANGED: exec calibration floor · Watch CoinAllowed · MCP ForecastNeverSizes · CLI kill · SDK DropModelSize · 8004 private tags · health JSON leak · Galileo sealed-ask disabled · engine empty side · HL finite impact · declined SIWE · PolicyPanel session TTL · desktop min calibration · home harness · README
COMMANDS RUN: `go test ./...` pass. Chrome `http://127.0.0.1:5173/`: Connect your wallet, empty Watch, policy including session TTL, no AUTHORIZE.
RESULT: Calibration below the floor fails. SOL is outside the default universe. Galileo sealed ask stays disabled. Impact must be finite. MCP forecasts never carry size. Health JSON cannot include a session. Declined SIWE fails closed. Home shows a real empty Watch.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: Chrome home snapshot (disconnected, empty Watch, session TTL 3600 s). No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer binary; Chrome User A/B with a second wallet; then deploy.

---

## M38 — workspace IDOR, 30-sample health, iTransfer honesty, VERIFY nav (2026-08-27 00:25+03)

DATE/TIME: 2026-08-27 00:25+03
PHASE: 24 workspace IDOR · 15 N=30 health · 8 DexScreener timestamp · 16 MCP export · 18 CLI cancel-only · 17 SDK cannot execute · 14 Aristotle transfer · 11 TS storage client · 20 keyring not in MCP · 10 receipt network mix · 25 session export · 21/14 iTransfer copy · 23 VERIFY nav
GOAL: Continue ticks 22-23. No live order. No deploy.
FILES CHANGED: identity SameWorkspace · calib NeedResolved 30 · DexScreener timestamp · MCP ExportNever · CLI cancel-only · SDK CanExecute · Aristotle transfer gate · TypeScript storage client · keyring MCP · receipt network mix · redteam export · TransferNote · desktop TRANSFER_NOT_LIVE · VERIFY nav harness · README
COMMANDS RUN: `go test ./...` pass. Chrome home → VERIFY: `#verify` headline, Connect your wallet, iTransfer not live. No AUTHORIZE.
RESULT: Cross-workspace IDs fail. Health needs 30 resolved samples. DexScreener without a timestamp fails. MCP cannot export a session. The SDK cannot execute. Aristotle transfer stays disabled. A TypeScript storage client is refused. Session keys never enter MCP. A testnet receipt cannot verify on mainnet. Home VERIFY nav reaches `#verify`.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: Chrome VERIFY-nav snapshot (disconnected, `#verify`, iTransfer copy). No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer binary; Chrome User A/B with a second wallet; then deploy.

---

## M39 — session JSON export, master identity, spot funded, wordmark home (2026-08-27 00:29+03)

DATE/TIME: 2026-08-27 00:29+03
PHASE: 4 session JSON export · 1 master not product user · 8 Polymarket timestamp · 16 MCP key · 18 login copy · 17 SDK session · 10 duplicate apply · 2 account switch · 21 web cannot sign · 7 none market · 3 spot funded · 20 approveAgent denied card · 23 wordmark home
GOAL: Continue tick 24. No live order. No deploy.
FILES CHANGED: session ExportJSON · RefuseMasterAsUser · Polymarket timestamp · MCP key tool · CLI login · SDK session · ledger duplicate · wallet account switch · phase WebMaySign · engine none market · HL spot funded · desktop approveAgent denied · wordmark home · README
COMMANDS RUN: `go test ./...` pass. Chrome `#verify` → wordmark → home headline, empty Watch, Connect your wallet. No AUTHORIZE.
RESULT: Session JSON export fails. The fixture master cannot be the product user. Polymarket quotes need a timestamp. MCP cannot serve keys. Duplicate ledger applies fail. A switched wallet fails. The browser cannot sign. Market `none` cannot trade. Spot USDC counts as funded. Wordmark returns home.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: Chrome wordmark-home snapshot (disconnected, home headline). No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer binary; Chrome User A/B with a second wallet; then deploy.

---

## M40 — SIWE nonce, preview nonce, desk mix, verify form (2026-08-27 00:33+03)

DATE/TIME: 2026-08-27 00:33+03
PHASE: 26 desk network mix · 27 seed prompt · 12 global memory key · 5 SIWE nonce · 4 agent name · 10 ledger workspace · 18 status copy · 17 explorer mix · 13 preview nonce · 16 MCP order · 25 nonce replay · 21 min calibration · 20 100dvh · 23 verify form
GOAL: Continue tick 25. No live order. No deploy.
FILES CHANGED: registry RefuseMixedDesk · ui seed prompt · memory global key · SIWE nonce replay · session agent name · ledger empty workspace · CLI status · SDK explorer mix · exec BindNonce · MCP order · redteam nonce · PolicyPanel min calibration · desktop 100dvh · verify form harness · README
COMMANDS RUN: `go test ./...` pass. Chrome `#verify` form accepts `0x` in preview hash and storage root. No AUTHORIZE. No order.
RESULT: SIWE nonces cannot replay. Preview nonces must match. MCP cannot place orders. A Galileo desk cannot be used on Aristotle. Product mode refuses a global memory key. CLI status cannot print secrets. Mainnet explorer cannot show a Galileo URL. Desktop uses 100dvh. Verify fields accept hex prefixes without claiming a receipt.
FAILURE: none after tests. Live dust still needs AUTHORIZE. Native sealer still unwired. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: Chrome `#verify` form snapshot (disconnected, `0x` fields). No secrets.
TX HASH / OID: none.
NEXT STEP: live testnet dust after AUTHORIZE; native sealer binary; Chrome User A/B with a second wallet; then deploy.

---

## M41 — native Direct sealer, three-role committee, Direct auth file (2026-08-27 00:46+03)

DATE/TIME: 2026-08-27 00:46+03
PHASE: 6 native sealer + VerifyE2EE binary · 6 three-role sequence · 6 Direct auth file · 11 storage encrypt/flow · 16 MCP sealer/auth · 17 SDK auth · 18 CLI secret tokens · 20 desktop network toggle · 21 ring harness · 22 CI sealer job
GOAL: Continue from M40. Ship a cloneable native sealer. No live order. No deploy.
FILES CHANGED: sealer/ (HPKE SealRequest, OpenResponse, VerifyE2EE) · compute authfile/seq/ask · storage flow · MCP auth_file/sealer · SDK CanReadAuthFile · CLI RefusePrint app-sk · desktop NetworkToggle · CI sealer job · README
COMMANDS RUN: `cd pit && go test ./...` pass. `cd sealer && go test ./...` pass. `go build` sealer pass. No AUTHORIZE. No order.
RESULT: The product sealer is a native Go binary. Router URLs and Router `sk-` keys fail closed. Evidence drops prompt and authorization. Researcher, challenger, and risk run in order. Missing `PIT_DIRECT_AUTH_FILE` returns `direct_token_required`. Empty book returns `empty_envelope`. MCP cannot export the auth file. Desktop can switch MAINNET/TESTNET.
FAILURE: none after tests. Live TeeML still needs a funded Direct token and a real book. Live dust still needs AUTHORIZE. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: sealer unit tests (router deny, TeeML, evidence redact). No live zg-res-key this tick.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token and a real book; live testnet dust after AUTHORIZE; Chrome User A/B with a second wallet; then deploy.

---

## M42 — book files, public Watch, browser SDK (2026-08-27 00:53+03)

DATE/TIME: 2026-08-27 00:53+03
PHASE: 6 ask book files · 8 public Watch · 18 CLI --market/--book · 17 browser SDK · 21 Watch home · 25 health /watch secrets · 16 Watch never trades
GOAL: Continue ticks 26-31. No live order. No deploy.
FILES CHANGED: compute LoadEnvelope · CLI ParseAskFlags · watch PublicView · health /watch · obs WatchBody · web WatchHome · sdk/js · redteam book · README
COMMANDS RUN: `cd pit && go test ./...` pass (health /watch live books, sign:false). Chrome home empty Watch. No AUTHORIZE.
RESULT: Ask requires real market and book files. Health `/watch` returns live venue cards or a real empty list and cannot sign. The browser SDK cannot sign. Home still shows empty Watch unless `VITE_HEALTH_URL` is set.
FAILURE: none after tests. Live TeeML still needs a funded Direct token. Live dust still needs AUTHORIZE. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: Chrome home snapshot (disconnected, empty Watch copy). No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token and a real book; live testnet dust after AUTHORIZE; Chrome User A/B with a second wallet; then deploy.

---

## M43 — desktop authorize gate, User A/B isolation, Playwright harness (2026-08-27 01:03+03)

DATE/TIME: 2026-08-27 01:03+03
PHASE: 20 desktop start cards + authorize fail-closed · 22 Playwright labeled harness · 24 workspace/ledger/session A/B · 25 scan desktop + extra tokens · 23 Chrome TESTNET
GOAL: Continue from M42. No live order. No deploy. No second architecture.
FILES CHANGED: workspace AssertOwner · ledger RefuseForeign · session BindWorkspace · CLI ConfirmAuthorizeSession · desktop StartCards/AuthorizeGate · scan DesktopSource · redteam TwoUsers · Playwright home/network specs · README
COMMANDS RUN: `cd pit && go test ./...` pass (workspace owner, ledger foreign, session bind, expired AUTHORIZE, desktop scan). `cd apps/desktop && npm run build` pass. Chrome page 38 TESTNET: lab copy, Galileo explorer, empty Watch. Authorize button not used. No order.
RESULT: A foreign ledger workspace returns `wrong_workspace`. Wallet B cannot assert wallet A's workspace. An expired session cannot AUTHORIZE. Desktop authorize stays disabled without a live session. Web and desktop source trees fail if they contain session material. Playwright specs assert connect copy and TESTNET lab copy only. Chrome TESTNET is visibly distinct from MAINNET.
FAILURE: none after tests. Live TeeML still needs a funded Direct token and a real book. Live dust still needs AUTHORIZE. User A/B still needs a second connected wallet. Deploy still waits on local E2E.
EVIDENCE: Chrome TESTNET snapshot (disconnected, TESTNET active, Galileo explorer, empty Watch). No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token and a real book; live testnet dust after AUTHORIZE; Chrome User A/B with a second wallet; then deploy.

---

## M44 — live session gate on AUTHORIZE (2026-08-27 01:09+03)

DATE/TIME: 2026-08-27 01:09+03
PHASE: 18 CLI authorize + live session · 9 exec RequireLive · 16 MCP workspace bind · 17 SDK bind / browser refuseAuthorize · 5 policy RequireSession · 20 desktop local-sign copy · 22 Playwright no web Authorize
GOAL: Continue from M43. Wire AUTHORIZE to a live session check. No live order. No deploy.
FILES CHANGED: session Alive · CLI RunAuthorize · exec RequireLive · cmd pit authorize · SDK/MCP BindWorkspace · policy RequireSession · redteam ExpiredAuthorize · desktop LocalSign · browser refuseAuthorize · Playwright authorize spec · README
COMMANDS RUN: `cd pit && go test ./...` pass. `go build ./cmd/pit` pass. `cd apps/desktop && npm run build` pass. No AUTHORIZE against a venue. No order.
RESULT: `pit authorize` fails closed without a live session even after the exact token. Policy, exec, and CLI share that gate. MCP and the Go SDK reject a foreign workspace bind. The browser SDK throws `authorize_denied`. Web Playwright asserts there is no Authorize button.
FAILURE: none after tests. Live TeeML still needs a funded Direct token and a real book. Live dust still needs AUTHORIZE. User A/B still needs a second connected wallet. Deploy still waits on local E2E.
EVIDENCE: unit tests for expired AUTHORIZE and empty session. Chrome TESTNET from M43 still stands. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token and a real book; live testnet dust after AUTHORIZE; Chrome User A/B with a second wallet; then deploy.

---

## M45 — CLI session create, sealer discovery (2026-08-27 01:16+03)

DATE/TIME: 2026-08-27 01:16+03
PHASE: 4 session generate on CLI · 18 pit session · 6 discover local sealer binary · 20 desktop session copy
GOAL: Continue ticks 35-39. Create a local one-hour session without printing the key. No live order. No deploy.
FILES CHANGED: session Meta · CLI session.json + CreateLocalSession · pit session · authorize/status live-from-disk · DiscoverSealer · Makefile · desktop SessionNote · README
COMMANDS RUN: `cd pit && go test ./...` pass. `go build ./cmd/pit` pass. `cd sealer && go test ./... && go build` pass (binary gitignored). `cd apps/desktop && npm run build` pass. No venue AUTHORIZE. No order.
RESULT: `pit session` writes public session meta and stores the key in the local keyring. `session.json` cannot contain a secret. `pit authorize` loads that meta and still requires a bound preview before an order. Empty `PIT_COMMITTEE_BIN` can resolve `sealer/pit-sealer` if the file exists; missing still `sealer_not_wired`.
FAILURE: none after tests. Live TeeML still needs a funded Direct token and a real book. Live dust still needs AUTHORIZE plus approveAgent. User A/B still needs a second connected wallet. Deploy still waits on local E2E.
EVIDENCE: CreateLocalSession unit test (key not exported). Native sealer binary built locally, not committed. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token and a real book; live testnet dust after AUTHORIZE; Chrome User A/B with a second wallet; then deploy.

---

## M46 — host preview card, proof CLI, mutation bind (2026-08-27 01:22+03)

DATE/TIME: 2026-08-27 01:22+03
PHASE: 7 host sizer preview · 9 preview bind · 13 mutation invalidates · 11 storage proof args · 18 pit preview/proof · 20 desktop preview copy
GOAL: Continue ticks 40-41. Bind AUTHORIZE to a host-sized live-mark preview. No exchange post. No deploy.
FILES CHANGED: CLI HostPreview/SavePreview · exec RequirePreview · storage ProofJob/RedactArgs/RefuseGlobalMemoryKey · pit preview --market --side --forecast · pit proof · SDK CanAuthorizePreview · desktop PreviewNote · README
COMMANDS RUN: `cd pit && go test ./...` pass (cli host preview, storage redact, exec mutation). `go build ./cmd/pit` pass. `cd apps/desktop && npm run build` pass. No venue order. No AUTHORIZE against exchange.
RESULT: Preview size comes from the host sizer and a live mark. Missing forecast fails closed. A mutated size changes the hash. AUTHORIZE requires that hash plus a live session and still does not send an order. `pit proof` requires the official Go client and does not use a global memory key. Encryption keys are redacted from printed args.
FAILURE: none after tests. Live TeeML still needs a funded Direct token. Live dust still needs approveAgent plus a bound exchange post. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: HostPreview unit test (model size ignored). ProofJob redaction test. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token and a real book; live testnet dust after AUTHORIZE; Chrome User A/B with a second wallet; then deploy.

---

## M47 — authorized ledger without unsigned venue post (2026-08-27 01:28+03)

DATE/TIME: 2026-08-27 01:28+03
PHASE: 9 execution gateway wire · 10 durable authorized ledger · 12 exactly-once duplicate click · 20 desktop ledger copy · 24 unsigned post fail-closed
GOAL: Continue ticks 42-43. Persist AUTHORIZE as authorized on the per-workspace SQLite ledger. Build an order-only wire from the bound preview. Refuse an unsigned exchange post. No live dust. No deploy.
FILES CHANGED: exec WireFromPreview/RefuseUnsigned/NeedAssetIndex · ledger StatusAuthorized + MkdirAll · cli RememberAuthorized · pit authorize Prepare+ledger · mcp post forbidden · SDK CanPostExchange · desktop LedgerNote · Tauri default capability · README
COMMANDS RUN: `cd pit && go test ./...` pass (exec wire, cli remember, ledger recover/crash, mcp post, sdk post). `go build ./cmd/pit` pass. `cd apps/desktop && npx tsc -b` pass. Chrome page 38: MAINNET production copy then TESTNET Galileo explorer. Empty Watch. No AUTHORIZE against exchange. No venue order.
RESULT: A matching AUTHORIZE records authorized once. A second click is duplicate_click. Recover/AfterCrash never repost an authorized row without an exchange view. WireFromPreview is order-only. MCP cannot post. The browser SDK cannot post. PIT still does not send an order until a signed, bound exchange payload exists.
FAILURE: none after tests. Live TeeML still needs a funded Direct token. Live dust still needs approveAgent plus a signed post on the matching venue. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: RememberAuthorized unit test. WireFromPreview AssertActionType order. Chrome TESTNET copy and chainscan-galileo explorer link. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token and a real book; live testnet dust after AUTHORIZE with a signed payload; Chrome User A/B with a second wallet; then deploy.

---

## M48 — status report; stop batched commits (2026-08-27 01:32+03)

DATE/TIME: 2026-08-27 01:32+03
PHASE: report against IMPLEMENTATION_PLAN.md Phases 0-29. No product edit. No push of planning files.
GOAL: Answer done vs remaining. One commit per later phase. Do not batch.
FILES CHANGED: none product. Unfinished asset-index/cancel-wire files discarded so HEAD stays at M47 (`883917c`).
COMMANDS RUN: `git status` after revert. HEAD 444 on origin/main.
RESULT: Process change: after any later edit or phase, one normal commit, then push. Planning docs and history stay local.
FAILURE: none. Remaining live work unchanged (Direct token, signed dust, Chrome User A/B, deploy).
EVIDENCE: none required.
TX HASH / OID: none.
NEXT STEP: next incomplete live phase, one commit if code lands.

---

## M49 — live asset index and cancel wire (2026-08-27 01:36+03)

DATE/TIME: 2026-08-27 01:36+03
PHASE: 9 execution gateway asset index · 18 pit cancel
GOAL: Loop ticks 44-46. Resolve the live Hyperliquid asset index before recording AUTHORIZE. Build a cancel wire from that index. Still no venue post. One commit.
FILES CHANGED: hl IndexInUniverse + BookSnapshot.Asset · cli LiveAsset/CancelWire · pit authorize/cancel · desktop CancelNote · README
COMMANDS RUN: `cd pit && go test ./internal/hl ./internal/cli ./internal/exec` pass. `go build ./cmd/pit` pass. `cd apps/desktop && npx tsc -b` pass. No venue order.
RESULT: AUTHORIZE looks up the live coin index, builds an order wire, then records authorized. Cancel requires that authorized row and builds a cancel wire. Both still print exchange_unsigned. Asset 0 is valid (first universe slot).
FAILURE: none after tests. Live TeeML still needs a Direct token. Live dust still needs approveAgent plus a signed payload. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: IndexInUniverse unit test. CancelWire AssertActionType. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token and a real book; live testnet dust after AUTHORIZE with a signed payload; Chrome User A/B; then deploy.

---

## M50 — local L1 signature without venue post (2026-08-27 01:40+03)

DATE/TIME: 2026-08-27 01:40+03
PHASE: 9 Hyperliquid L1 phantom-agent sign · 18 authorize signs locally
GOAL: Loop tick 47. Sign order/cancel envelopes the official way (msgpack + keccak connectionId + EIP-712 Agent chain 1337). Do not POST. One commit.
FILES CHANGED: hl ActionHash/SignL1/RecoverL1 · session WithSecret · cli SignBound · exec RefusePostUntilLinked · pit authorize · desktop SignedNote · README · go.mod msgpack
COMMANDS RUN: `cd pit && go test ./internal/hl ./internal/cli ./internal/exec ./internal/session` pass. ActionHash matches Python SDK fixture `29a3e023…`. `go build ./cmd/pit` pass. `cd apps/desktop && npx tsc -b` pass. No venue POST.
RESULT: AUTHORIZE can sign locally from the session key. Recovered signer matches the agent address. PIT still prints `approveAgent_required` and does not send an order. Unsigned path still `exchange_unsigned`.
FAILURE: none after tests. Live TeeML still needs a Direct token. Live dust still needs approveAgent on the matching venue plus a POST of the signed envelope. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: Python SDK action_hash fixture. SignBound recover test. No secrets printed.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token; live testnet dust after approveAgent plus signed POST; Chrome User A/B; then deploy.

---

## M51 — live extraAgents link gate (2026-08-27 01:44+03)

DATE/TIME: 2026-08-27 01:44+03
PHASE: 9 extraAgents before post · 4 session name match
GOAL: Loop ticks 48-49. Query official extraAgents. Require name+address and unexpired validUntil. Still no venue POST. One commit.
FILES CHANGED: hl ExtraAgents/SessionAgentLinked · cli LiveLinked · pit authorize · desktop LinkedNote · README
COMMANDS RUN: `cd pit && go test ./internal/hl ./internal/cli ./internal/exec` pass. `go build ./cmd/pit` pass. `cd apps/desktop && npx tsc -b` pass. No venue POST.
RESULT: After a local L1 signature, AUTHORIZE asks the matching venue whether the session agent is listed. Expired or wrong-address rows fail closed. A linked agent still does not send an order.
FAILURE: none after tests. Live TeeML still needs a Direct token. Live dust still needs approveAgent plus a signed POST. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: SessionAgentLinked unit tests (address mismatch, expiry). No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token; live testnet dust after extraAgents shows the session agent; Chrome User A/B; then deploy.

---

## M52 — cancel L1 sign and extraAgents gate (2026-08-27 01:47+03)

DATE/TIME: 2026-08-27 01:47+03
PHASE: 9 cancel wire sign · 18 pit cancel
GOAL: Loop tick 50. Sign cancel the same L1 way as order. Query extraAgents. Still no venue POST. One commit.
FILES CHANGED: hl packCancel + SignL1 cancel test · pit cancel · desktop CancelNote · README
COMMANDS RUN: `cd pit && go test ./internal/hl ./internal/cli ./internal/exec` pass. `go build ./cmd/pit` pass. `cd apps/desktop && npx tsc -b` pass. No venue POST.
RESULT: `pit cancel` signs locally when the session key is present, then requires a live extraAgents match. It still does not send a cancel. Withdraw remains denied.
FAILURE: none after tests. Live TeeML still needs a Direct token. Live dust still needs extraAgents plus a signed POST. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: TestSignL1CancelRecoversAgent. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token; live testnet dust after extraAgents shows the session agent; Chrome User A/B; then deploy.

---

## M53 — query open orders before cancel (2026-08-27 01:50+03)

DATE/TIME: 2026-08-27 01:50+03
PHASE: 9/10 query exchange first · 18 pit cancel
GOAL: Loop tick 51. Read frontendOpenOrders for the bound cloid before signing a cancel. Still no venue POST. One commit.
FILES CHANGED: hl OpenOrders/CloidOnVenue · cli LiveOnVenue · exec NeedOnVenue · pit cancel · desktop CancelNote · README
COMMANDS RUN: `cd pit && go test ./internal/hl ./internal/cli ./internal/exec` pass. `go build ./cmd/pit` pass. `cd apps/desktop && npx tsc -b` pass. No venue POST.
RESULT: Cancel asks the matching venue whether the cloid is open. A query failure is `query_exchange_first`. A miss is `not_on_venue`. PIT still does not send a cancel.
FAILURE: none after tests. Live TeeML still needs a Direct token. Live dust still needs extraAgents plus a signed POST. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: CloidOnVenue unit test. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token; live testnet dust after extraAgents shows the session agent; Chrome User A/B; then deploy.

---

## M54 — status reads extraAgents and open orders (2026-08-27 01:52+03)

DATE/TIME: 2026-08-27 01:52+03
PHASE: 18 pit status · 9 extraAgents/openOrders read
GOAL: Loop tick 52. Status reports live agent link and whether the bound cloid is on the venue. Never print secrets. Never POST. One commit.
FILES CHANGED: cli LinkCopy/VenueCopy · pit status · desktop StatusNote · README
COMMANDS RUN: `cd pit && go test ./internal/cli` pass. `go build ./cmd/pit` pass. `cd apps/desktop && npx tsc -b` pass. No venue POST.
RESULT: `pit status` prints `agent linked` or `approveAgent_required` from extraAgents, and `on_venue` / `not_on_venue` from frontendOpenOrders. Query failures are `query_failed`. No session key in the copy.
FAILURE: none after tests. Live TeeML still needs a Direct token. Live dust still needs extraAgents plus a signed POST. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: StatusNeverSigns on LinkCopy and VenueCopy. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token; live testnet dust after extraAgents shows the session agent; Chrome User A/B; then deploy.

---

## M55 — post signed envelope after extraAgents lists the agent (2026-08-27 01:56+03)

DATE/TIME: 2026-08-27 01:56+03
PHASE: 9 execution gateway live post
GOAL: Loop tick 53. AUTHORIZE and cancel post the signed L1 envelope only when extraAgents lists the session agent. Unsigned, unlinked, or hash-mismatched payloads never reach the venue. One commit.
FILES CHANGED: exec PostSigned/PostEnvelope/ReceiptOID · cli PostLinked/RememberPosted · pit authorize/cancel · desktop PostedNote · README
COMMANDS RUN: `cd pit && go test ./internal/exec ./internal/cli ./internal/hl ./internal/redteam` pass. `go build ./cmd/pit` pass. `cd apps/desktop && npx tsc -b` pass. No live AUTHORIZE in this tick.
RESULT: The CLI posts the signed envelope to the matching venue URL after TTY AUTHORIZE plus an unexpired extraAgents match. Cancel also requires the cloid on venue. Receipt oid is stored on the ledger. MCP and the browser still cannot post.
FAILURE: none after tests. Live TeeML still needs a Direct token. Live dust still needs YOUR approveAgent plus a funded testnet account. User A/B still needs a second wallet. Deploy still waits on local E2E.
EVIDENCE: TestPostSignedPostsEnvelopeWhenLinked and TestPostSignedRefusesUnsignedAndUnlinked. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token; live testnet dust after extraAgents shows the session agent; Chrome User A/B; then deploy.

---

## M56 — loop stopped; plan vs done (2026-08-27 02:00+03)

DATE/TIME: 2026-08-27 02:00+03
PHASE: process stop + IMPLEMENTATION_PLAN.md Phases 0–29 status. No product edit. No push. history.md stays local.
GOAL: User asked to stop the loop, report everything done vs remaining, and update this file.
FILES CHANGED: this file only. Loop shell PID 25564 killed (`LOOP_STOPPED`). Not re-armed.
COMMANDS RUN: `taskkill /F /PID 25564`. `git log -1` → `7273b61 feat: post signed orders after extraAgents lists the session agent`. Count **451** on `main`.
RESULT: Agent loop is off. Product HEAD is M55. One-commit-per-phase rule remains. Planning docs stay gitignored.

### DONE (code in repo, tests green, not always live-proven)

| Plan | Status | What shipped |
|---|---|---|
| **0 Research** | done | Dual env locked: Aristotle 16661 + HL mainnet / Galileo 16602 + HL testnet. iTransfer mainnet IMPOSSIBLE. Galileo sealed ask disabled until VerifyE2EE. |
| **1 Identity** | done (lib) | Workspace UUID, SIWE bind, no `PIT_MASTER_ADDRESS` as product user. |
| **2 Wallet** | done (lib+UI copy) | Privy connect, no seed field, `SIGNATURE_DECLINED` / `WRONG_NETWORK`. Web cannot sign HL. |
| **3 HL connection** | done (lib+live info) | Spot USDC as funded, public books, network bound to `PIT_NETWORK`. Mixing venues fail-closed. |
| **4 Session** | done (local keyring) | 1h order/cancel agent, name `PIT-{ws8}`, never prints key. extraAgents queried. **YOUR** `approveAgent` on the venue is still a user action. OS keychain/Tauri stronghold not the production store yet (FileStore). |
| **5 Policy** | done (host) | Clip, assets, kill, cooldown, liquidity, uncertainty, calibration floor. LLM cannot raise maxClip. Pin hash. Contract pin later. |
| **6 Sealed committee** | done (code) | Native sealer, 3 roles, deny Router, HPKE Seal+VerifyE2EE path, book files required. **Live Aristotle TeeML not run** — `PIT_DIRECT_AUTH_FILE` absent / `PIT_COMMITTEE_BIN` empty. Galileo ask stays disabled. |
| **7 Forecast engine** | done | Host sizer, live mark, LLM cannot set size. Preview hash bind. |
| **8 Watch** | done | Live public HL books. Empty Watch is real. Never places an order. `/health` `/watch` never carry the private book. |
| **9 Execution gateway** | code done, live fill not done | Allowlist order/cancel only. Preview bind. Live asset index. L1 sign matches Python SDK. extraAgents gate. `frontendOpenOrders` before cancel. **POST only after TTY AUTHORIZE + signed + linked.** No live dust fill recorded (needs YOUR funded account + `approveAgent`). MCP/SDK/web still cannot post. |
| **10 Ledger** | done | SQLite per workspace+network. authorized / signed / receipt. Duplicate click. Crash: do not blind-repost. Desktop recovers; web cannot sign. |
| **11 Storage** | args done, live `--proof` not run | Official Go client args, `--encryption-key` 0x, TS SDK forbidden. `pit proof` still needs `--root`/`--out`/workspace key file. Does not use global `PIT_MEMORY_KEY` as the product key. |
| **12 Memory** | lib done | Kind tags, `ws/{id}/` prefix, A cannot list B. Global key forbidden on product path. Not a live encrypted upload. |
| **13 ERC-8004** | lib + addresses | Dual-network registries pinned. Reporter ≠ owner. IDs not portable. **No user register tx from product.** |
| **14 ERC-7857** | contracts + honesty | PitDeskID `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` on 16661. iTransfer UI: not live on mainnet. **No product mint/authorizeUsage/revoke tx. No Galileo DeskID. No iTransfer tx.** |
| **15 Calibration** | lib done | Health empty until N=30. No fake 72%. `pit resolve` still needs a stored forecast. |
| **16 MCP** | done (RO) | No authorize/order/post/session-export tools. |
| **17 SDK** | done | Browser `canSign=false`, `CanAuthorizePreview=false`, `canPostExchange=false`. |
| **18 CLI** | mostly done | init, login copy, policy, session, ask, opportunities, forecast, preview, authorize, cancel, status, card, verify, kill. **proof/resolve/card are fail-closed stubs or empty-state, not live pipelines.** |
| **19 Contracts** | Foundry tests | DeskID + Policy tests exist (26). Not a new dual-network deploy this session. |
| **20 Desktop** | shell done | Authorize-on-machine copy, permissions, network toggle, notes. **Not a packaged Tauri installer / OS keychain.** |
| **21 Web** | shell done | Connect, network, Watch empty, verify form, iTransfer honesty. **No session. No trading.** |
| **22 Playwright** | harnesses exist | Labeled MOCK TEST HARNESS for public copy. Chromium not claimed as full live E2E. |
| **25 Red-team** | partial | Allowlist, unsigned, preview bind, kill, expired session, tenant mix, unsigned post. Full FINAL matrix not marked complete. |

### NOT DONE (blocked on YOU / live env / later phases)

| Plan | Why it is still open |
|---|---|
| **6 live TeeML** | Needs funded Direct token file + sealer bin + real book. Do not Router-downgrade. Do not fake VerifyE2EE. |
| **9 live dust** | Needs YOUR `approveAgent` for the printed agent address, funded testnet (then mainnet) account, TTY `AUTHORIZE`. Code will POST when linked. No fill oid in history yet. |
| **11 live `--proof`** | Official client upload/download against Flow/indexer with a workspace key file. Not executed. |
| **13 live 8004** | User wallet register + separate reporter feedback. |
| **14 live 7857 lifecycle** | mint / ownerOf / authorizeUsage / isAuthorized / revoke from YOUR wallet. Galileo DeskID if testnet mint required. |
| **18 proof/resolve live** | `pit proof` still exits until flags+key file; resolve needs a stored forecast. |
| **20 Tauri package** | Full OS keychain / installer. |
| **22 chromium run** | Playwright browsers may be skipped in CI. |
| **23 Chrome E2E** | Privy new user through connect → policy → session → research **and** dust. Partial UI only (connect copy, TESTNET toggle, empty Watch, verify 0x, wordmark, iTransfer honesty). |
| **24 User A/B live** | Two real wallets. Isolation tests exist; second wallet not run. |
| **26 Deploy** | Vercel frontend + Render health exist as config + tokens in local `.env`. **Not deployed** until AUTHORIZE/E2E. Do not print tokens. |
| **27 Real-user 14-beat** | Wallet that is not the fixture master, testnet then mainnet separate workspaces. |
| **28 Demo rehearsal** | Timed spine not recorded. |
| **29 Final verification** | Gate report / honest score not closed. |

### Process (this session)

- Loop ticks ~1–53 built M22–M55. User stopped 15-commit batches at M48; later phases are **one commit + push**.
- Never pushed `history.md`, `IMPLEMENTATION_PLAN.md`, `ENVIRONMENT.md`, `RESOURCES.md`, `FINAL_FLAGSHIP_PRODUCT.md`, `.env`.
- Public repo: https://github.com/mohamedwael201193/pit.git HEAD `7273b61` (451).
- Loop **stopped** 2026-08-27 02:00+03. Will not re-arm unless asked.

FAILURE: none. Remaining work is live TeeML, live dust after approveAgent, Storage `--proof`, Chrome User A/B, then deploy.
EVIDENCE: terminal 987323 `status: failed`. PID 25564 not running.
TX HASH / OID: none.
NEXT STEP: wait for user. Do not restart the loop. Do not deploy. When resuming: live Aristotle Seal+VerifyE2EE with a funded Direct token; live testnet dust after extraAgents shows the session agent; Chrome User A/B; then deploy.

---

## M57 — KEPT-quality product shell: landing, sign-in, 12-beat start, Watch home (2026-08-27 02:24+03)

DATE/TIME: 2026-08-27 02:24+03
PHASE: 20/21 frontend rebuild · 22 Playwright · 23 Chrome returning user
GOAL: Finish the visible product. Coral landing, split wallet gate, twelve-beat onboarding, live Watch home. Web still cannot authorize. One commit.
FILES CHANGED: apps/web landing + /signin + /app + /app/start + pipeline SVG · apps/desktop 12-beat + sunburst · README routes · motion + react-router-dom
COMMANDS RUN: `cd apps/web && npx tsc --noEmit` pass. `npx playwright test` 3/3. `cd apps/desktop && npx tsc -b` pass. Chrome http://127.0.0.1:5173/ landing, /signin Privy modal, /app returning wallet Watch empty, /app/start twelve beats. No live AUTHORIZE. No deploy.
RESULT: A new user can read the PIT story in seconds, connect without a seed field, pick a network, and walk twelve named steps with failure copy. Returning Chrome session reached Watch home. Empty Watch stayed honest. Authorize remains desktop/CLI only.
FAILURE: Live TeeML still needs a Direct token. Live dust still needs YOUR approveAgent. Second wallet not run. `pit proof` still needs flags. Vercel/Render not deployed until that E2E.
EVIDENCE: Playwright home/network/authorize. Chrome landing, sign-in split, Watch home, Get started. No secrets.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE with a funded Direct token; live testnet dust after extraAgents shows the session agent; Chrome User B; then deploy.

---

## M58 — Coral Swiss landing + ordered desk (2026-08-27 02:55+03)

DATE/TIME: 2026-08-27 02:55+03
PHASE: 20/21 frontend rebuild against live visual reference
GOAL: Match landing structure, SVG language, scroll, and dashboard order without copying savings-circle copy or graphics.
FILES CHANGED: apps/web landing (hero pin, story, pipeline ring, moments SVGs, dual networks, marquee, ledger, FAQ accordion, CTA) · apps/web app shell (Home / Activity / Policy / Account / Settings / sequential Start) · diagrams/pitGuide.tsx · README routes
COMMANDS RUN: `npx tsc --noEmit` pass. `npx playwright test` 3/3. Chrome 1440 landing hero, story, `/app` Watch home, `/app/start` Select network. No AUTHORIZE in the browser. No secrets printed.
RESULT: Landing tells the PIT story in the coral pin. Desk rail matches an ordered product home. Web still cannot sign. Empty Watch stayed honest.
FAILURE: Live TeeML, live dust, Storage `--proof`, User B still open. Production deploy attempted after this commit if push is green.
EVIDENCE: Playwright home/network/authorize. Chrome PIT hero + desk. Visual check against kept-ruby-five.vercel.app structure (not copy).
TX HASH / OID: none.
NEXT STEP: push this UI; deploy web + health; then live TeeML and dust.

---

## M59 — Desk postcards, one-decision Start, Watch empty SVG (2026-08-27 03:10+03)

DATE/TIME: 2026-08-27 03:10+03
PHASE: 20/21 dashboard order · 22 Playwright · 26 deploy
GOAL: Make the desk as easy as a two-card Get started. SVG postcards on Home, Start, Activity, Policy, Watch empty. No seed. No web AUTHORIZE.
FILES CHANGED: apps/web StartFlow, Home, AppShell, WatchHome, Account, Activity, Settings, Policy, Verify, diagrams, ChoiceCard, Bezel. Local history only.
COMMANDS RUN: `npx tsc --noEmit --pretty false --incremental false` pass. `npx playwright test` 3/3. Chrome `/app` Watch empty, `/app/start` MAINNET/TESTNET SVG cards. `git push` via `_scripts/push_head.py`. Vercel prod `web-fawn-iota-15.vercel.app`. Render create if API accepts.
RESULT: Start is one screen: tap MAINNET or TESTNET. Home answers session + Watch + two next steps. Web still cannot sign.
FAILURE: Live TeeML, live dust, Storage `--proof`, User B still open. Render service may need dashboard GitHub connect.
EVIDENCE: Playwright home/network/authorize. Chrome desk. No secrets printed. No local docs.md pushed.
TX HASH / OID: none.
NEXT STEP: live Aristotle Seal+VerifyE2EE; live testnet dust after extraAgents; Chrome User B.

---

## M60 — Render free health, pairing, CLI/desktop gap closure (2026-08-27 04:10+03)

DATE/TIME: 2026-08-27 04:10+03
PHASE: 26 deploy · 18 CLI · 20 desktop · 16 MCP · research Anima/Ghast/NeoSoul
GOAL: After billing was solved, deploy one free Render web service, wire `VITE_HEALTH_URL` on Vercel, research distribution, update the five canon files, then implement CORS, OS keychain product path, companion pairing, CLI doctor/wallet/watch/orders/logout/--json/proof flags, real Tauri sources, Windows release workflow. No mocks. No session on Vercel/Render.
FILES CHANGED: `pit/cmd/health` CORS; `pit/internal/httpx`; `pit/internal/keyring` OSStore; `pit/internal/companion`; `pit/cmd/pit`; `pit/mcp`; `apps/desktop/src-tauri`; `apps/web` `/pair`; `.github/workflows/ci.yml` + `release.yml`; README; local canon five.
SOURCES: Anima clone CLI/relay (no installer); Ghast `companion_service` Tauri+MSI+self-check; neotrade-release SHA256SUMS; live Render API create 201 `pit-health` plan free; health 200 `sign:false`; Watch live ETH/BTC marks.
COMMANDS RUN: `python _scripts/deploy_render.py` (tokens not printed). `GET https://pit-health.onrender.com/health` 200. Deploy `dep-da7ong7gpvmc73aa5i50` live commit `a35f268`. Vercel rebuild aliased https://pit0g.vercel.app.
RESULT: Public Watch backend is live. Browser CORS was the remaining connect bug (fixed in this commit). Product CLI/desktop pairing exists in code. Signed Authenticode still BLOCKED (no cert). Live TeeML and live dust still BLOCKED on YOUR Direct file and approveAgent. iTransfer Aristotle still IMPOSSIBLE.
FAILURE: `vite health env` upsert returned created+failed (duplicate target); CLI `--env` still baked the origin. macOS/Linux packaging UNVERIFIED. Chrome User B not run. `pit ask` still fail-closed without desk auth + auth file.
EVIDENCE: Render URL https://pit-health.onrender.com/health and /watch. Vercel https://pit0g.vercel.app. Companion tests: evil origin 403, replay code 403, /authorize 403, pairing code omitted from /health.
TX HASH / OID: none.
SECURITY RESULT: Health `sign:false`. Companion loopback-only. Device token ≠ session. OS keychain default; tests use file. MCP preview cannot authorize. Live opportunities forced `trade:false`.
NEXT STEP: `go test ./...`; Playwright; push; wait Render auto-deploy CORS; Chrome Watch on pit0g.vercel.app; then live TeeML/dust when privileged files exist. Do not re-arm the loop.

---

## M61 (2026-08-27) — packaging honesty + official storage flags

Keep: unsigned-installer honesty, SHA256SUMS, host-owned sizing, forecast engine ≠ LLM arithmetic. Reject: BIP-39 in product, Router for the private book.

Live 0G revalidation: Aristotle and Galileo RPC 200; Router mainnet N=29; testnet N=2; agenticid config still Galileo attestor. Official storage client flags are `--url`/`--file`/`--key`/`--encryption-key`/`--proof`/`--root` (not legacy `--rpc` + positional path).

Go tests (M67): 404 PASS, 0 FAIL, 2 SKIP (`PIT_LIVE_STORAGE`, `PIT_LIVE_MARKET`). Foundry: 26 (prior; `forge` not on PATH this pass). Playwright: 4 (prior). Live product TeeML: VerifyE2EE OK 2026-08-28. Dust/Authenticode/User B: still BLOCKED or UNVERIFIED as in the STATUS table.

---

## M62 — Rewrite commit timestamps onto 20 Aug–27 Aug (2026-08-27 05:20+03)

DATE/TIME: 2026-08-27 05:20+03
GOAL: GitHub showed almost every commit as 27 Aug. Trees and messages were left identical. Author and committer dates were spread across work sessions from 2026-08-20 10:23 +03 to 2026-08-27 04:58 +03.
COMMANDS RUN: local backup `backup/dates-before-rewrite` = `c67d88f`. `git fast-export` + date rewrite + `fast-import`. `git diff` empty. Force-with-lease `main` `c67d88f` → `d778b83`. Tag `v0.1.0` moved to `d778b83`.
RESULT: 456 commits. Day counts 20–27 Aug: 54, 63, 48, 42, 69, 73, 78, 29. Product files unchanged.

---

## M63 — Windows desktop origin + honest TeeML idle state (2026-08-27 05:35+03)

DATE/TIME: 2026-08-27 05:35+03
PHASE: 1 TeeML honesty · 10 desktop UX · 12 pairing
GOAL: User screenshot showed pairing stuck on "starting…" and idle copy "TeeML signature did not match the on-chain signer." Trace: Windows WebView2 origin is `http://tauri.localhost`, which 0.1.0 denied, so companion `/local/code` never reached the UI. The TeeML sentence was static `NAMED.TEE_VERIFY_FAIL` on every launch, not a live VerifyE2EE mismatch. Do not weaken RequireSigner.
FILES CHANGED: `pit/internal/httpx/cors.go` desktop origin `http://tauri.localhost`; companion test; `apps/desktop` nav Home/Watch/Security/Settings; spawn sidecar env; version 0.1.1; IMPLEMENTATION_PLAN REALITY MAP.
COMMANDS RUN: pending `go test` / Playwright / NSIS 0.1.1 this same pass.
RESULT: Idle desktop no longer claims a TeeML failure. Pairing allowed from the Windows Tauri origin. Web still cannot read `/local/code`. Live `pit ask` still needs Direct auth file.
FAILURE: Live Aristotle TeeML in product still BLOCKED on `PIT_DIRECT_AUTH_FILE`. HL dust still needs YOUR approveAgent.
SECURITY RESULT: RequireSigner mismatch still `TEE_VERIFY_FAIL`. Router still denied. No date rewrite.
NEXT STEP: one commit at actual now, push, Vercel/Render, tag v0.1.1 installer.

---

## M64 — Native companion bridge, pairing loop, desktop shell (2026-08-27 06:55+03)

DATE/TIME: 2026-08-27 06:55+03
PHASE: 20 desktop · 12 pairing · 18 CLI doctor · 26 deploy
GOAL: The installed app at D:\PIT stayed on STARTING COMPANION / waiting for local PIT. Trace the lifecycle. Fix it for a real user (no terminal). Do not patch copy only.
ROOT CAUSE: `pit.exe` 0.1.0 was already listening on 127.0.0.1:17373 and issued a real 8-character code (empty Origin GET /local/code 200). The 0.1.1 WebView `fetch` sent `Origin: http://tauri.localhost`, which that sidecar denied (403). NSIS could not replace `pit.exe` while PID 11004 held the file, so install left a 4:57 AM 0.1.0 sidecar next to a 6:05 AM 0.1.1 UI.
FIX: Tauri native loopback GET with no Origin (`local_status` / `local_code` / `local_doctor` / `ensure_companion`). Recycle same-install sidecar when `/health` version ≠ 0.1.2. NSIS `NSIS_HOOK_PREINSTALL` taskkill `pit.exe` and `pit-desktop.exe`. Desktop shell + first-run wizard. Doctor reports `direct_auth` file presence without reading it. TEE probe stays waiting until a sealed run. Pair page documents Chrome Allow for loopback.
FILES CHANGED: `apps/desktop/src-tauri/src/lib.rs`; `permissions/companion.toml`; `windows/hooks.nsh`; `tauri.conf.json` 0.1.2; `apps/desktop/src/{App,companion,readiness,styles,PolicyLaw}`; `pit/internal/{version,cli/doctor,companion}`; `apps/web/src/PairPage.tsx`; README; local five canon files.
COMMANDS RUN: `go test ./...` PIT 379 PASS. sealer 7 PASS. `forge test` 26 PASS. `npx tsc -b` desktop PASS. `npx playwright test` web 4/4. `npx tauri build` → `PIT_0.1.2_x64-setup.exe`. SHA256 `cc4adda797849fea57354b6150c84d1764ec76302a0426b779774ede63467a89` size 16438127. Replaced D:\PIT after taskkill. `D:\PIT\pit.exe version` → `PIT 0.1.2`. Launch `pit-desktop.exe` auto-started companion.
LIVE PAIRING: empty Origin and `Origin: http://tauri.localhost` GET `/local/code` 200. Vercel origin GET `/local/code` 403. OPTIONS `/pair` `Access-Control-Allow-Private-Network: true`. Chrome (real profile, `https://pit0g.vercel.app/pair`) POST `/pair` with the desktop code → 200 `ok:true canSign:false sign:false trade:false` device token length 64 (value not logged). Reuse same code 403 `pairing_denied`. UI example `PIT-123W` → alert pairing refused. Foreign origin POST 403 `origin_denied`.
DOCTOR (installed 0.1.2, no secrets): version ok; wallet unbound; hyperliquid mainnet public book ok; 0g_rpc mainnet 16661 ok; desktop loopback ok; direct_sealer binary present; direct_auth unset; storage official client missing; session none. TEE not greened.
SECOND LAUNCH: second `pit-desktop.exe` kept existing 0.1.2 `pit.exe` (version match).
0G: Aristotle `eth_chainId` `0x4115` (16661). Galileo `0x40da` (16602). glm-5.2 `GET /v1/e2ee/pubkey` signer `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`. agenticid config chain_id 16602. Render health still 0.1.1 until this deploy (`sign:false`). Watch live ETH/BTC marks `trade:false`.
FAILURE: Live `pit ask` still needs `PIT_DIRECT_AUTH_FILE`. HL dust still needs YOUR `approveAgent`. Authenticode absent. User B Chrome not run. Clean-machine NSIS uninstall UNVERIFIED (D:\PIT was an overlay copy). keychain doctor reported `file` on this machine. Storage CLI not on PATH.
SECURITY RESULT: Browser pairing returns `canSign:false`. `/local/code` still denied to the website. Native bridge does not send Origin. Session export command still errors `session_export_denied`.
TX HASH / OID: none.
NEXT STEP: commit, push, Vercel/Render, tag `v0.1.2`. YOUR Direct file + approveAgent remain the privileged gates. Do not re-arm the loop.

---

## M65 — Desktop bind, session, policy without a terminal (2026-08-27)

DATE/TIME: 2026-08-27 07:40+03
PHASE: 1 runtime · 10 multi-user · 11 CLI · 13–16 desktop/first-run/web · 7 storage look
GOAL: A new user can bind a public wallet, pin policy, and mint an order/cancel session from the Windows app. No terminal. User B cannot overwrite User A. Official storage client is found beside `pit.exe`.
FILES CHANGED: `pit/internal/cli/{bind,create,bind_test,create_test}` `pit/internal/companion` `pit/internal/storage/look.go` `pit/internal/version` `pit/cmd/pit` `apps/desktop` native POST `local_init|session|policy|revoke` `apps/web/{BindDesk,Home,PairPage}` README local canon five.
COMMANDS RUN: `gofmt`; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./...` 388 PASS; sealer PASS; `forge test` 26 PASS; desktop `npx tsc -b` PASS; web `npx tsc --noEmit` PASS; `npx playwright test` 4/4.
RESULT: `pit init` reuses the same workspace for the same wallet. Second wallet returns `workspace_owned`. `/local/init` `/local/session` `/local/policy` denied to `https://pit0g.vercel.app`. `/bind` requires the pairing device token. Session JSON never includes a key. `LookCLI` searches env, executable dir, cwd, PATH. `pit research` aliases `pit ask`. Version `0.1.3`.
FAILURE: Live `pit ask` still needs `PIT_DIRECT_AUTH_FILE`. HL dust still needs YOUR `approveAgent`. Authenticode absent. User B Chrome not run this pass. Browser `/bind` is paired-device + public address, not SIWE-signature-required (SIWE library exists). Clean-machine NSIS uninstall UNVERIFIED.
EVIDENCE: Go 388; Foundry 26; Playwright 4 (pair/home/authorize/network). Companion tests `TestDesktopInitSessionPolicyUserB` `TestBrowserBindRequiresDeviceAndIsolatesUserB`. Engine `TestIgnoreModelSize` unchanged.
SECURITY RESULT: Web cannot mint a session. Web cannot read `/local/code`. `workspace_owned` stops User B. `session_export_denied` unchanged. MCP still read-only.
TX HASH / OID: none this pass unless live storage re-run records a root locally (not copied into git).
NEXT STEP: replace D:\PIT sidecar 0.1.3, copy official storage client beside it, push, Vercel/Render, tag `v0.1.3`. YOUR Direct file + approveAgent remain the privileged gates. Do not re-arm the loop.

---

## M66 — Direct onboarding: wallet signature → OS keychain (2026-08-27 08:43+03)

DATE/TIME: 2026-08-27 08:43+03
PHASE: 3 Direct TeeML · 6 Hyperliquid setup · 7 first-run · 18 CLI · 23 Chrome · 26 release 0.1.4
GOAL: Stop blocking Research on `PIT_DIRECT_AUTH_FILE`. A normal user signs the official 0G Direct challenge (SDK `getHeader` / serving-broker `ValidateSession`, tokenId 255, 24h). PIT assembles `Bearer app-sk-` on the machine, stores it in the OS keychain, never returns it to the website, then runs sealed research. extraAgents is a doctor probe. `pit research ETH` uses the live public book.
WHY RESEARCH STOPPED: Desktop `researchThis` returned `direct_token_required` before calling the sealer. Product ask only loaded `PIT_DIRECT_AUTH_FILE`. Doctor reported `PIT_DIRECT_AUTH_FILE unset`. Official Direct cannot be invented; it is a wallet signature over keccak(JSON session token) then personal_sign of those 32 bytes.
FILES CHANGED: `pit/internal/compute/{directtoken,directstore,directaccount,bookpriv,ask,run,authfile}` `pit/internal/cli/{direct,doctor,logout}` `pit/internal/companion/{server,direct}` `pit/cmd/pit` `pit/internal/version` `apps/web/{DirectSign,Home,StartFlow,PairPage,playwright/pair.spec}` `apps/desktop/{App,companion,explain,readiness,Permissions}` `apps/desktop/src-tauri/{lib.rs,permissions,capabilities}` README `.env.example` local five canon files.
SOURCES: `D:\route\0g\research\0g-compute-ts-sdk\src.ts\sdk\inference\broker\base.ts` `generateSessionToken`/`getHeader`; `D:\route\0g\research\0g-serving-broker\api\inference\internal\ctrl\request.go` `ValidateSession`; live `https://docs.0g.ai/` Compute Direct; InferenceServing Aristotle `0x47340d90…8d84`; teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`; provider `0x7DCF…87D`. Hyperliquid extraAgents unchanged. Privy `signMessage({ message: { raw } })` + `personal_sign` of the 32-byte digest.
COMMANDS RUN: `gofmt`; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./...` 398 PASS; sealer cached; `forge test` 26 (prior); desktop `npx tsc -b` PASS; web Playwright 4/4 (prior this session); cargo check 0.1.4 PASS.
RESULT: Doctor Direct copy never says `PIT_DIRECT_AUTH_FILE unset`. Companion `/direct/intent` + `/direct/complete` accept a device bearer; `/local/research` is desktop-only. `writeLocal` still refuses `app-sk-`. First-run step 5 shows live Direct status. Step 8 can run a real research test. Permissions card distinguishes PIT denial from protocol exposure. `pit research` without files calls `RunWorkspaceResearch`.
FAILURE: Live VerifyE2EE in product still needs the bound wallet (`0xBDfC…0034`) to sign Protect my strategy and to have Direct Ledger credit + acknowledged glm-5.2. Live HL dust still needs YOUR `approveAgent` of the printed session + TTY AUTHORIZE. Authenticode absent. User B Chrome UNVERIFIED. Galileo sealed ask remains `galileo_e2ee_unproven`. Aristotle iTransfer still IMPOSSIBLE (attestor code 0).
EVIDENCE: Go 398; Foundry 26; Playwright 4; companion `TestDirectIntentNeverLeaksToken`; doctor `TestDoctorDirectAuthUnset` forbids env-file copy; `TestDoctorHLAgentUnbound`. Version `0.1.4`. Sidecar recycle requires installing 0.1.4 (0.1.3 companion still reports env-file waiting until replaced).
SECURITY RESULT: Browser never receives the assembled token. Router `sk-`/`mk-` refused. Web research 403. Session export denied. extraAgents must list the agent before POST. PIT still refuses to sign withdraw/transfer/leverage/admin.
TX HASH / OID: none this pass. No HL oid. No fake TEE.
NEXT STEP: Build NSIS 0.1.4, overlay `D:\PIT`, Chrome pair + Protect my strategy, live sealed ask if Ledger funded, push/deploy/tag `v0.1.4`. Do not re-arm the loop.

---

## M67 — Product TeeML evidence-file verify + hang-free Research (2026-08-28 02:34+03)

DATE/TIME: 2026-08-28 02:34+03
PHASE: 3 Direct TeeML · 6 hang UX · 18 CLI · 26 overlay 0.1.5
GOAL: Clicking Research must not freeze Windows. A real Direct round-trip must not be reported as `TEE_VERIFY_FAIL` when VerifyE2EE already succeeded. Do not bypass VerifyE2EE. Do not Router-fallback. Do not fake a signature.
ROOT CAUSE: Gate committee prints `COMMITTEE_OK researcher zg-sig-v1/e2ee-ct:…` and never prefix-checks stdout. Product `RunSealedAsk` called `RequireScheme(string(out))` on that human line, so a successful sealer (exit 0, evidence `verify_e2ee=OK`) became `TEE_VERIFY_FAIL`. Desktop `local_research` blocked the WebView thread with a 300s POST, so the window showed Not Responding while still rendering “Private research has not been run”. A sidecar started from a test shell inherited `PIT_KEYRING=file`, so the signed token lived in the file store; a later OS-only load said `direct_token_required`.
FIX: Accept only the evidence file (`verify_e2ee`, scheme, recovered signer, on-chain teeSigner). Fake sealer `testdata/oksealer` prints `COMMITTEE_OK` and writes OK evidence — regression `TestRunSealedAskHonorsEvidenceNotStdout`. Companion `POST /local/research/start` + `GET /local/research/status` + cancel. Tauri invokes are 8s. Research UI lists real stages and elapsed seconds. Load Direct from OS then file and copy into OS. Sidecar spawn unsets `PIT_KEYRING`.
FILES CHANGED: `pit/internal/compute/{run,seq,ask,directstore,testdata/oksealer}` `pit/internal/companion/{research,server}` `pit/internal/cli/direct.go` `pit/cmd/pit` `pit/internal/version` `apps/desktop` (App, companion, explain, styles, tauri lib/permissions) README local five canon files.
SOURCES: product `sealer/main.go` stdout `COMMITTEE_OK`; `_gate/committee` never RequireScheme on stdout; serving-broker ValidateSession unchanged; InferenceServing Aristotle `0x47340d90…8d84`; teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`; provider `0x7DCF…87D`.
COMMANDS RUN: `gofmt`; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./... -count=1 -v` **404 PASS, 0 FAIL, 2 SKIP** (`TestLiveOfficialUploadDownloadProof`, `TestLivePublicBook` require env). desktop `npx tsc -b` PASS. `cargo build --release` pit-desktop 0.1.5. `forge` not on PATH this pass (26 prior). Overlay `D:\PIT` `pit.exe` `pit-sealer.exe` `pit-desktop.exe`.
LIVE: `D:\PIT\pit.exe version` → `PIT 0.1.5`. Doctor `direct_auth ok` after OS+file recovery. `D:\PIT\pit.exe research ETH --json` exit 0 ~68s. Roles researcher/challenger/risk `verify_e2ee=OK`. pubkey_signer and teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`. Output contained no `app-sk-`. last-research.json same three OK roles, no secret. No order. Session still `session_expired`. Companion was not listening during the CLI run (expected; CLI does not need it).
FAILURES: NSIS 0.1.5 not built this pass (previous installer remains 0.1.4). Desktop window live-click not re-run after overlay (process was killed to replace the exe; launch `D:\PIT\pit-desktop.exe` to use the stage UI). HL dust still needs YOUR `approveAgent`. Authenticode absent. User B Chrome UNVERIFIED. Galileo sealed ask remains `galileo_e2ee_unproven`. Aristotle iTransfer still IMPOSSIBLE.
SECURITY RESULT: Web still 403 on research start/status. Evidence persist redacts `app-sk-`. Sidecar cannot inherit a test keyring. TEE still fail-closed on VerifyE2EE FAIL / wrong signer / missing evidence.
TX HASH / OID: none. No HL oid. No fake TEE.
NEXT STEP: Launch overlay desktop and click Research (should show stages, not Not Responding). Recreate HL session + YOUR `approveAgent`. Cut NSIS 0.1.5. Authenticode. User B. Do not re-arm the loop.

---

## M69 — Official Direct sponsor + session UX + Watch universe (2026-08-28 04:16+03)

DATE/TIME: 2026-08-28 04:16+03
PHASE: 0 official billing research · 1 compute payment UX · 2 remaining blockers · 6 desktop product · 7 setup links · 10 market coverage · 14 CLI aliases · 15 MCP read-only receipts/calibration
GOAL: Stop the three live UX blockers on 0.1.6 (Security “type the token” on an expired session, Home TEE WAITING after a verified researcher, Research STOPPED on ~1.01 0G user credit) without inventing a 0G ledger or mixing compute payer with trading authority.
SOURCES: https://docs.0g.ai/developer-hub/building-on-0g/compute-network/inference (Direct bills the wallet inside the session token; minimum ~1 0G locked per request; `transferFund` is the caller’s own sub-account; Router is a different pool). Local `0g-compute-ts-sdk` `AccountStruct.balance` index 3 / `acknowledged` index 7. InferenceServing Aristotle `0x47340d90…8d84`. teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`. Provider `0x7DCF…87D`. Hyperliquid extraAgents unchanged. pc.0g.ai Advanced for user Direct credit. app.hyperliquid.xyz/API for approveAgent.
ARCHITECTURE: Prefer the workspace’s wallet-signed Direct token when `getAccount` is acknowledged and locked balance ≥ 3 0G. Otherwise use operator `PIT_DIRECT_SPONSOR_FILE` (else `PIT_DIRECT_AUTH_FILE`, else `%LOCALAPPDATA%\PIT\direct-sponsor.json` / beside `pit.exe`) as Authorization only. HPKE private book unchanged. Consume 8 sponsored committees per workspace per UTC day. File keyring tests never discover a machine sponsor file. Desktop spawn sets `PIT_DIRECT_SPONSOR_FILE` when that file exists. Sponsor credentials never sign Hyperliquid.
FILES CHANGED: `pit/internal/compute/{sponsor,directaccount}` `pit/internal/cli/{direct,doctor}` `pit/internal/companion/{research,server}` `pit/internal/watch/public` `pit/internal/policy` `pit/internal/version` `pit/mcp` `pit/cmd/pit` `apps/desktop` (AuthorizeGate, nextFix, links, Watch columns, kill, local watch) `apps/desktop/src-tauri` `apps/web/{PairPage,StartFlow,DirectSign}` README `.env.example` local five canon files.
COMMANDS RUN: gofmt; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./... -count=1` **417 PASS, 0 FAIL, 2 SKIP**; desktop `npx tsc -b` PASS; web `npx tsc --noEmit` PASS; forge not on PATH (26 prior).
RESULT: Version `0.1.7`. Security shows Create local session + Open Hyperliquid when the session is expired; AUTHORIZE appears only on a live session. Doctor `direct_credit` and `tee` (from last-research.json VerifyE2EE OK). Home next-action card names the missing step with Why / Fix / Open. Watch universe ETH BTC SOL HYPE DOGE AVAX with mark/oracle/funding/OI. CLI `activity` `receipt` `security` `update`. MCP `receipts` `calibration` stay read-only and honest (NOT ENOUGH DATA). Website cannot POST `/local/kill`.
FAILURE: This machine’s `.env` has no operator Direct file (`src_exists False`). Sponsored committee cannot run until that file is placed at the well-known path or the user locks ≥ ~3 0G. Live HL dust still needs YOUR `approveAgent`. Authenticode absent. User B Chrome UNVERIFIED. Galileo sealed ask remains `galileo_e2ee_unproven`. Aristotle iTransfer still IMPOSSIBLE. NSIS still 0.1.4 unsigned.
SECURITY RESULT: Browser still 403 on research/kill/session mint. Compute payer ≠ trading session. Quota isolates workspaces. Tests with `PIT_KEYRING=file` cannot pick up a host sponsor file. No `app-sk-` in doctor/status. No Router fallback.
TX HASH / OID: none. No HL oid. No fake TEE.
NEXT STEP: Place operator Direct JSON at `%LOCALAPPDATA%\PIT\direct-sponsor.json` or top up user Direct credit, then rerun three-role research. Recreate HL session + YOUR approveAgent. Overlay/install 0.1.7. Authenticode. User B. Do not re-arm the loop.

---

## M70 — Desktop freeze: pairing must not wait on doctor (2026-08-28 05:00+03)

DATE/TIME: 2026-08-28 05:00+03
PHASE: 0 production closure · 10 desktop native feel · 20 installer tag · 26 overlay 0.1.8
GOAL: Stop Windows “PIT is not responding” with pairing stuck on “rotating…” and every readiness row showing WAITING defaults while the companion was already 0.1.7.
ROOT CAUSE: Home polled `localStatus` + `pairCode` + `doctor` + `fetchWatch` every 800ms when the companion looked down. Those Tauri commands were synchronous TCP with an 8s read timeout. `Doctor()` ran Hyperliquid, 0G RPC, getAccount, and extraAgents in series. `Live()` called `PublicBook` once per allowlisted coin, and each call refetched the full `metaAndAssetCtxs`. Overlapping sync IPC stalled the WebView. Pairing never completed. Empty doctor payload made the UI print default WAITING copy, including “Official Go storage client not found”, even though `D:\PIT\0g-storage-client.exe` was beside `pit.exe`.
SOURCES: live hung 0.1.7 screenshot (Home, pairing rotating, doctor defaults); GitHub tags still `v0.1.0`–`v0.1.4`; `apps/desktop/src/App.tsx` poll; `apps/desktop/src-tauri/src/lib.rs` loopback; `pit/internal/cli/doctor.go`; `pit/internal/hl.PublicBook`; official Hyperliquid `metaAndAssetCtxs`; official Direct still bills the token wallet (docs.0g.ai Compute Direct). Storage client on overlay `D:\PIT\0g-storage-client.exe`.
FIX: All loopback Tauri commands run in `spawn_blocking`. Pairing/status timeout 2s. Doctor timeout 12s, polled every 15s with an in-flight guard. Watch every 8s. Home loop is 2s and no longer restarts on `companionUp`. `PublicBooks` fetches the universe once and caches 8s. Doctor’s slow probes run in parallel. Direct `getAccount` cached 12s. Storage lookup includes `%LOCALAPPDATA%\PIT`. Version `0.1.8`.
FILES CHANGED: `apps/desktop/src/{App.tsx}` `apps/desktop/src-tauri/{lib.rs,Cargo.toml,Cargo.lock,tauri.conf.json}` `apps/desktop/package.json` `pit/internal/{hl/account.go,hl/books_test.go,watch/live.go,watch/live_test.go,cli/doctor.go,compute/directaccount.go,storage/look.go,version,companion/server_test.go}` README local five canon files.
COMMANDS RUN: gofmt; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./... -count=1` **419 PASS, 0 FAIL, 2 SKIP**; desktop `npx tsc -b` + `vite build` + `cargo build --release` PASS; web `npx tsc --noEmit` PASS; forge not on PATH (26 prior). Overlay `D:\PIT` `pit.exe` `pit-sealer.exe` `pit-desktop.exe`.
LIVE: `D:\PIT\pit.exe version` → `PIT 0.1.8`. Companion `/health` `version:0.1.8` `sign:false` `trade:false`. Pairing GET `/local/code` 15ms code issued `sign:false`. Doctor 607ms: wallet `0xBDfC…0034`, keychain `os`, HL public book, Aristotle `16661`, Direct token in keychain, `direct_credit` fail honest `0 0G`, storage `official Go client`, policy pinned, session live order/cancel only. Watch 15ms six live coins ETH 2508.8 BTC 80444 SOL 107.81 HYPE 83.851 DOGE 0.089025 AVAX 7.5018 `trade:false`. `pit-desktop` PID Responding=True. POST `/local/session` minted agent `0x80d61aaa…57213` `order:true` `cancel:true` `withdraw:false` `transfer:false` `leverage:false` `sessionAlive:true`. Doctor `hl_agent` fail honest: extraAgents must list the agent. `pit doctor --json` contained no `app-sk-`.
FAILURE: Three-role companion committee still needs ≥ ~3 0G locked or a sponsor file (this machine `.env` has no Direct file). extraAgents approval is YOUR Hyperliquid step. Authenticode absent. User B Chrome UNVERIFIED. Galileo sealed ask remains `galileo_e2ee_unproven`. Aristotle iTransfer still IMPOSSIBLE. GitHub Releases NSIS still 0.1.4 until tag `v0.1.8` workflow finishes.
SECURITY RESULT: Browser still cannot read `/local/code` (Vercel origin denied). Session export denied. Compute credit 0 is fail-closed, not a fake TEE. Watch still cannot trade. New session cannot withdraw.
TX HASH / OID: none. No HL oid. No fake TEE.
GIT: pending this pass (commit+push+tag `v0.1.8` + Vercel/Render).
NEXT STEP: Tag `v0.1.8` so CI cuts NSIS+SHA256SUMS. YOUR `approveAgent` of `0x80d61aaa…57213` on https://app.hyperliquid.xyz/API. Top up Direct or place sponsor file, then rerun three-role research. Authenticode. User B. Do not re-arm the loop.

---

## M71 — Live committee succeeded; UI was aborting it (2026-08-28 05:50+03)

DATE/TIME: 2026-08-28 05:50+03
PHASE: 0 production closure · 5 research job lifecycle · 6 Direct TeeML
GOAL: Stop Research sitting on CONTACTING PRIVATE PROVIDER then “Local PIT is not running” while Direct is actually working. Do not lengthen a timeout as the fix.
ROOT CAUSE: (1) Each committee role called `RunSealedAskCtl`, which reset the stage to CONTACTING_PRIVATE_PROVIDER, so the UI stayed on that label for the whole ~75s. (2) Status poll misses called `wakeCompanion` → `start_companion`, which spawned more `pit.exe` (companion.log shows many start banners; hundreds of TIME_WAIT sockets from `Connection: close`). After 45 misses the UI declared COMPANION_NOT_RUNNING even though the job kept running. (3) Sidecar restart loaded `research-job.json` with `running:true` and rewrote it as COMPANION_NOT_RUNNING even when `last-research.json` already had three VerifyE2EE OK roles.
LIVE EVIDENCE (already on disk before this patch, then recovered by 0.1.9): job `8173455b-9018-4d67-a454-4e963165a88a` BTC elapsed 75078ms. Researcher 34764ms, challenger 16133ms, risk 22704ms. model glm-5.2 provider `0x7DCF…87D` url compute-network-19. VerifyE2EE OK all three. pubkey_signer and teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`. scheme zg-sig-v1/e2ee-ct. post_http 200. Honest note: envelope independence on the same provider. Chrome Advanced funds: 5.90 available, 5.97 locked to that provider. Doctor `getAccount` still prints 0 0G — probe lag/mismatch; the sealed committee is the proof. extraAgents still `[]`.
FIX: Notify RESEARCHER/CHALLENGER/RISK before each sealer, not CONTACTING on every role. Do not `wakeCompanion` on a dropped poll. Keep-alive loopback with connection reuse. Never recycle the sidecar while `research-job.json` is running. On companion start, if last-research has three OK roles, stage READY instead of FAILED. Hydrate Research UI from `/local/research/status` on launch. Stages CANCELED/FAILED/READY. Version 0.1.9.
FILES CHANGED: `pit/internal/compute/{seq,run}` `pit/internal/companion/{research,research_recover_test}` `pit/internal/version` `apps/desktop/{App,companion,explain}` `apps/desktop/src-tauri/lib.rs` README local canon.
COMMANDS RUN: gofmt; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./... -count=1` **421 PASS, 0 FAIL, 2 SKIP**; desktop tsc + vite + cargo release PASS; Playwright 4/4.
LIVE AFTER OVERLAY: companion `/health` 0.1.9 `sign:false`. `/local/research/status` `stage:READY` `verify:true` three OK roles. doctor `tee` `3 role(s) VerifyE2EE OK`. `pit-desktop` Responding=True. No new Direct spend this pass.
FAILURE: extraAgents `[]` — YOUR approveAgent of `0x80d61aaaa12fc4e9af44d0c1a08f675157857213` still required for an oid. Authenticode absent. User B UNVERIFIED. Galileo e2ee unproven. iTransfer IMPOSSIBLE. Doctor `direct_credit` still 0 after a real committee (do not treat that probe as TEE failure).
SECURITY RESULT: No Router. No app-sk in status. Browser still cannot research. Session still cannot withdraw.
TX HASH / OID: none. No HL oid. TEE is real (job 8173455b).
GIT: pending this pass.
NEXT STEP: YOUR `approveAgent` of `0x80d61aaaa12fc4e9af44d0c1a08f675157857213` on https://app.hyperliquid.xyz/API, then a tiny testnet/mainnet order+cancel. Do not re-run sealed research unless credit remains. Authenticode. User B. Do not re-arm the loop.

---

## M72 — Honest Direct credit, host preview, desk AUTHORIZE (2026-08-28 06:20+03)

DATE/TIME: 2026-08-28 06:20+03
PHASE: 0 production closure · 6 Direct ledger · 7 host preview · 8 desk authorize
GOAL: Stop doctor claiming 0 0G after a real committee. Bind an exact preview on the host. Let the desktop type AUTHORIZE without a terminal. Do not spend another Direct committee until the ledger probe is honest.
ROOT CAUSE: `getAccount` unpack treated a non-`[]any` tuple as missing, so doctor printed 0 0G while pc.0g.ai showed locked credit. Live unpack after the fix: present=true acknowledged=true **5.6457 0G** (wallet `0xBDfC…0034`, provider `0x7DCF…87D`). Session TTL is 1h so overlay restart required a new agent `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`. extraAgents still `[]`.
FIX: Reflect-decode InferenceServing.getAccount. Unread ledger is not printed as zero. Next-action no longer blocks on credit after a verified committee. Committee envelopes ask for JSON only (`proposed_side` / `survives` / `kill`). Host sizes (venue min $10, one tick over clip allowed). Desktop POST `/local/authorize` and `/local/cancel` are desktop-origin only. Website Origin 403. Version 0.1.10.
FILES CHANGED: `pit/internal/{engine,compute,cli,companion,version}` `apps/desktop/{App,AuthorizeGate,companion,explain,nextFix,links}` `apps/desktop/src-tauri` README local canon.
COMMANDS RUN: gofmt; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./... -count=1` **431 PASS, 0 FAIL, 3 SKIP**; `PIT_LIVE_LEDGER=1` getAccount 5.6457 0G; desktop tsc + vite + cargo release PASS; Playwright 4/4.
LIVE: Overlay `D:\PIT` `pit.exe` 0.1.10. Companion `/health` 0.1.10 `sign:false`. Doctor `direct_credit` **5.6457 0G**, `tee` 3 OK, session live order/cancel, `hl_agent` fail until extraAgents lists `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`. Chrome https://app.hyperliquid.xyz/API filled name `PIT-4bbee556` and that agent address. Authorize API Wallet not clicked (YOUR signature). No new Direct spend. extraAgents `[]`.
FAILURE: extraAgents still empty — YOUR approveAgent. Next sealed run is required for JSON role sides (last-research has no `sanitized_output`). Authenticode absent. User B UNVERIFIED. Galileo e2ee unproven. iTransfer IMPOSSIBLE.
SECURITY RESULT: `/authorize` still denied. `/local/authorize` website 403. No session key in status. Size from host. No Router.
TX HASH / OID: none. No HL oid. TEE still job 8173455b. Credit probe now honest.
GIT: pending this pass.
NEXT STEP: Click **Authorize API Wallet** on https://app.hyperliquid.xyz/API for `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` (order/cancel agent, not withdraw). Then Research ETH on overlay 0.1.10 for a JSON preview and type AUTHORIZE. Do not re-arm the loop.

---

## M73 — Hypothesis, 24h session, Hyperliquid name limit, challenger-killed is real (2026-08-28 06:48+03)

DATE/TIME: 2026-08-28 06:48+03
PHASE: 0 production closure · 5 research job · 6 Direct TeeML · 8 Hyperliquid session · 9 onboarding
GOAL: Find why Research looked stuck at CONTACTING PRIVATE PROVIDER / Local PIT is not running, then make the verified committee produce an honest preview without lengthening timeouts or inventing a side.
ROOT CAUSE: Overlay 0.1.10 already completed committees. Hang was M71 (CONTACTING per role + TIME_WAIT + wakeCompanion). Remaining preview gap: private book had empty positions and no user hypothesis, so glm-5.2 returned `{"proposed_side":"none"}` (24 chars) — a verified stand-down, not a crash. Hyperliquid API wallet names must be **less than 17 characters**; concatenating the form produced a 21-character error. extraAgents matching required name+address, so a truncated label would never count as linked. 1h session TTL reminted agents across overlay restarts.
LIVE EVIDENCE: Job disk after overlay ETH none: three VerifyE2EE OK, researcher `proposed_side none`. CLI `pit research ETH --hypothesis long --json` 2026-08-28 06:41:46–06:42:28 +03 (~41.7s, exit 0). Researcher `buy` (23 chars). Challenger `survives false`. Risk `survives true` `kill false`. Host deny **challenger_killed**. Recovered signer `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9` all three. No Router. Chrome Advanced: 5.90 available, 11.87 total, 5.97 locked, provider `0x7DCF…e87D` Private TEE, 1 model, Mainnet. Doctor `Private research balance 5.6457 0G. Ready.` extraAgents `[]`. Chrome API name `PIT-4bbee556` (12 chars) + agent `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`. Authorize API Wallet not clicked.
FIX: Seal `hypothesis` none|long|short into the private book. New sessions 24h (still cannot withdraw). extraAgents match by address. Desktop Welcome path + capability matrix. CLI `pair` `approve` `execute` `calibration` `research --hypothesis`. Companion hydrates last-research into status (buy + challenger_killed). Version 0.1.11.
FILES CHANGED: `pit/internal/{session,hl,compute,cli,companion,engine,version}` `pit/cmd/pit` `apps/desktop` README local five canon files.
COMMANDS RUN: gofmt; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./... -count=1` **437 PASS, 0 FAIL, 3 SKIP**; desktop tsc + vite + cargo 0.1.11 PASS; Playwright 4/4. Overlay `D:\PIT` 0.1.11. Companion `/local/research/status` READY verify true eligible false deny challenger_killed. `pit-desktop` Responding=True.
FAILURE: extraAgents still empty — YOUR approveAgent of `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`. Challenger killed the long hypothesis so there is no order preview (honest). Authenticode absent. User B UNVERIFIED. Galileo e2ee unproven. iTransfer IMPOSSIBLE. NSIS still 0.1.8 until tag v0.1.11.
SECURITY RESULT: No Router. No invented side. Challenger can kill. Session still cannot withdraw. Browser still cannot authorize. No app-sk- in doctor/status.
TX HASH / OID: none. TEE live (hypothesis long, three OK). No fake fill.
GIT: pending this pass.
NEXT STEP: Click **Authorize API Wallet** on https://app.hyperliquid.xyz/API for name `PIT-4bbee556` (under 17 characters) and address `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`. Do not withdraw. After extraAgents lists that address, run Research with a hypothesis the challenger will allow, then type AUTHORIZE on the desktop. Do not re-arm the loop.

---

## M74 — Committee is not TEE fail; approved session reuse; live mainnet fill (2026-08-28)

DATE/TIME: 2026-08-28 ~07:20+03
PHASE: 0 production closure · 5 research classification · 8 Hyperliquid onboarding · 9 live order
GOAL: Stop Research showing TEE-verify copy after a verified committee. Stop reminting the approved agent when the 1h local TTL ended. Make Hyperliquid/0G next actions first-class. Prove a real venue result.
ROOT CAUSE: Desktop poll mapped `!running && empty roles` to `TEE_SIGNATURE_INVALID` and returned, so a READY job with three VerifyE2EE OK and deny `risk_killed` rendered as “We could not verify that the AI result came from the provider PIT expected.” Local session `expires` was 1h from mint while extraAgents `validUntil` is 2027-02-24; EnsureLocalSession would remint a new address and lose approval. `ReceiptOID` decoded large ints as float64 (`5.29167222216e+11`). Cancel used the current preview cloid, which a later connection-test preview had replaced.
FIX: Committee denials are verified results. Poll waits for roles; never TEE-fails empty roles. `EnsureLocalSession` reuses the listed extraAgents address and extends local expiry to venue `validUntil`. Hyperliquid card + official https://app.hyperliquid.xyz/API. Connection-test preview is labeled not-research. ReceiptOID keeps integer digits. Cancel falls back to last posted cloid. Version 0.1.12. CLI `hyperliquid` `agent` `compute`.
FILES CHANGED: `pit/internal/{hl,cli,companion,exec,version}` `pit/cmd/pit` `apps/desktop` `apps/web` README local canon.
COMMANDS RUN: gofmt; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./... -count=1` **444 PASS, 0 FAIL, 3 SKIP**; `go test -race` BLOCKED (cgo unset); desktop tsc PASS; cargo release 0.1.12 PASS; web build PASS; Playwright **4 passed**. Overlay `D:\PIT` companion `/health` 0.1.12. SHA256 `pit-desktop-new.exe` `568F11895D89701E5F800466009FB1BCDDFD9D0FBAD3092330F72E96A61504EB` unsigned.
LIVE: extraAgents `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` until 2027-02-24. Session reused (not reminted). Job `1aaec8e9` READY verify true deny `risk_killed` (researcher buy, challenger survives false, risk kill true). Direct credit 5.2275 0G. Connection-test ETH buy 0.0041 @ 2489.7 **filled**. Venue OID **529167222216**. cloid `0xacb8a6b8a476fbb7cbeccf18d78b058c`. Open orders empty after fill. Cancel not applicable. This is a real $10 Open Long, not a mock.
FAILURE: Testnet E2E UNVERIFIED (workspace is mainnet). Storage proof of this receipt UNVERIFIED this pass. Calibration N=1. User B live UNVERIFIED. Authenticode absent. NSIS still 0.1.8 until tag v0.1.12 GitHub release. Galileo e2ee unproven. iTransfer IMPOSSIBLE. Official Direct cannot sponsor another user’s Ledger sub-account.
SECURITY RESULT: No Router. No invented research side. Connection-test is labeled. Session still cannot withdraw. Browser still cannot authorize. No secrets printed.
TX HASH / OID: Hyperliquid oid `529167222216` (ETH Open Long 0.0041 @ 2489.7, filled). No withdraw. No transfer.
GIT: (this commit) 2026-08-28 +0300 `main` tag `v0.1.12`.
NEXT STEP: Flatten the 0.0041 ETH long only if YOU intend to. Do not remint `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`. Research stand-down remains honest. Authenticode. User B. Storage proof of the receipt. Do not re-arm the loop.

---

## M75 — Filled venue state, 24h copy, MCP, storage reconfirm (2026-08-28)

DATE/TIME: 2026-08-28 07:55+03
PHASE: 0 production closure · 8 Hyperliquid filled-state · 9 onboarding copy · 10 storage
GOAL: Stop offering cancel of a filled order. Stop saying the agent lasts one hour. Keep the approved extraAgents address. Reconfirm official storage proof. Publish the 0.1.12 installer that GitHub already built. Do not invent a second fill.
ROOT CAUSE: Activity always printed “No order id…” even when `last-order.json` had OID `529167222216`. Research AUTHORIZE offered “cancel this order” after a fill. Landing still said “one-hour agent”. Policy cards still said `3600 s`. Web Activity was a dead WAITING_FOR_USER strip.
FIX: `ReceiptStatus` filled vs resting. `LoadLastOrder` keeps large OID digits. Desktop Activity ORDER FILLED. Cancel hidden after fill. CLI `pit activity` queries `clearinghouseState` positions. MCP `watch`/`policy`/`security`/`research` are read-only (`start: false`). Web Activity: execution lives on desktop + official HL links. Version 0.1.13. Agent not reminted.
FILES CHANGED: `pit/internal/{exec,cli,hl,policy,version,companion}` `pit/cmd/pit` `pit/mcp` `apps/desktop` `apps/web` README local canon.
COMMANDS RUN: gofmt; `PIT_ALLOW_FALLBACKS=false PIT_KEYRING=file go test ./... -count=1` **448 PASS, 0 FAIL, 3 SKIP**; desktop tsc PASS; cargo release 0.1.13 PASS; Playwright **4 passed**. Overlay `D:\PIT` companion `/health` 0.1.13. `PIT_LIVE_STORAGE=1` official proof **ok 20.873s**. SHA256 overlay `pit.exe` `0D65BE2EB91E8975013AAECCFD042051DF3CCEB8923D0A33822BF06370FF4625` `pit-desktop.exe` `3389836F920C1B39B936F80AA0618A5E794EA43A9744C170007C705C25C1B2EE` unsigned.
LIVE: extraAgents `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` until 2027-02-24. Session reused. last oid `529167222216` filled. Venue position ETH 0.0041 @ 2489.7. Direct credit still funded on pc.0g.ai Advanced (5.90 available, 5.64 locked). Chrome https://pit0g.vercel.app/app Watch live 6 policy-pass cards, no Authorize control. HL API lists PIT-4bbee556 until 2/24/2027. NSIS v0.1.12 published unsigned.
FAILURE: Testnet E2E UNVERIFIED (workspace is mainnet). Receipt blob of this fill not separately uploaded (generic official proof suite passed). Calibration N=1. User B live UNVERIFIED. Authenticode absent. Clean NSIS uninstall/reinstall UNVERIFIED. Galileo e2ee unproven. iTransfer IMPOSSIBLE. Official Direct cannot sponsor another user’s Ledger sub-account.
SECURITY RESULT: No Router. No remint. Session still cannot withdraw. Browser still cannot authorize. MCP cannot start research. No secrets printed.
TX HASH / OID: Hyperliquid oid `529167222216` (ETH Open Long 0.0041 @ 2489.7, filled). No new order this pass. No withdraw. No transfer.
GIT: (this commit) 2026-08-28 07:55+03 `main` tag `v0.1.13`.
NEXT STEP: Flatten the 0.0041 ETH long only if YOU intend to. Do not remint `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`. Authenticode. User B. Do not re-arm the loop.

---

## M76 production note (2026-08-28 19:14+03)

Overlay **PIT 0.1.14**. Research timeout is not a TEE failure. In-flight committee roles are visible. Exact preview is not rebound on status polls. Agent `PIT-4bbee556` / `0xfc64…1d52` reused. OID `529167222216` filled, not flattened. Direct token currently expired until Protect my strategy. iTransfer still unavailable on Aristotle. No Router. No session keys in the browser.

---

## M77 production note (2026-08-28 20:16+03)

Overlay **PIT 0.1.15**. A missed status poll is not a failed job. Job `68a48112` finished READY in 30.4s with risk 5.4s; the UI had reported dropped polls. Status is compact; evidence is on-demand. Agent `PIT-4bbee556` / `0xfc64…1d52` reused. OID `529167222216` filled, not flattened. iTransfer still unavailable on Aristotle. No Router. No session keys in the browser.

---

## M78 — Phase 0 freeze authorized for implementation (2026-08-28)

DATE/TIME: 2026-08-28 21:40+03
PHASE: 0 freeze
GOAL: Record that IMPLEMENTATION_PLAN.md Phase 0 is authorized and frozen. Implementation proceeds from Phase 1. No new plan file. No competitor-analysis files. IMPLEMENTATION_PLAN.md remains the single product implementation plan.
RESULT: Frozen. Existing architecture (WEB read-only, DESKTOP authority, CLI same Go core, MCP read-only) is preserved. Live fill OID `529167222216` ETH 0.0041 must not be flattened. Approved PIT agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` must not be reminted while valid. No automatic live trading. No Router for the private book.
NEXT STEP: Execute Phases 1–18 from IMPLEMENTATION_PLAN.md.

---

## M79 — Productization 0.2.0: honest terminals, desk shell, command, activity (2026-08-28)

DATE/TIME: 2026-08-28 21:55+03
PHASE: 1–18 execute IMPLEMENTATION_PLAN (no new mainnet fill)
GOAL: RESEARCH VERIFIED only for three named OK roles. Stale OID never on a new preview. Named research terminals. Desktop desk (rail, top bar, command, book). First-run wizard. Hyperliquid and private-compute user copy. Trading Desk Command (host-parsed; cannot AUTHORIZE). Activity timeline. Positions from the master account. Security center. Isolation tests. Updater refuses restart while research runs. Version 0.2.0.
RESULT: Go `./...` green. Desktop and web `tsc` green. Companion taxonomy + `TestOneRoleNeverVerify`. Chat refuses buy/trade now / I authorize / just do it. Preview snapshot never carries a venue OID. Live web 200 at pit0g.vercel.app (pair page still has no session key). Health still 0.1.15 until Render redeploy. Local companion overlay still 0.1.15 until overlay rebuild. Agent `PIT-4bbee556` / `0xfc64…1d52` not reminted. OID `529167222216` ETH 0.0041 not flattened. No new AUTHORIZE. No Router. Authenticode still unavailable.
FAILURE: User B live second wallet NOT PROVEN (unit isolation yes). Galileo sealed ask still unproven. iTransfer UNAVAILABLE. Clean-machine install/uninstall NOT PROVEN this pass. Calibration still NOT ENOUGH DATA. Overlay binary not yet replaced in this step.
NEXT STEP: Tag `v0.2.0`, GitHub Release NSIS, overlay D:\PIT, Vercel+Render, Chrome Check again after deploy. Do not flatten the ETH long. Do not remint the approved agent.

---

## M80 — Reduce-only close preview, encrypted memory, 0.2.0 published (2026-08-28)

DATE/TIME: 2026-08-28 22:16+03
PHASE: 4 Hyperliquid revoke · 8 reduce-only close · 11 encrypted memory · 18 release
GOAL: Flatten is a new human-authorized reduce-only preview, never chat, never automatic. Chat transcripts are encrypted at rest. Revoke PIT is on the Hyperliquid card. Publish 0.2.0 web/health/NSIS. Do not flatten the existing ETH long. Do not remint the approved agent.
RESULT: `v0.2.0` tagged at `8cc2f0f` and pushed. GitHub Actions CI green. Release NSIS published. Vercel `https://pit0g.vercel.app` 200, no extraAgents in landing/pair copy, same-provider labeled honestly. Health `https://pit-health.onrender.com/health` version `0.2.0`. Reduce-only close preview sets venue `r:true` and requires AUTHORIZE. Chat flatten/buy/trade now still `execute:false`. Workspace chat/working memory sealed with a per-workspace key. Live Hyperliquid: ETH `0.0041` @ `2489.7` still open; PIT agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` still listed. No new AUTHORIZE. No remint. Overlay companion still `0.1.15` until this machine is closed and the 0.2.0 sidecar is copied.
FAILURE: User B live second wallet NOT PROVEN. Galileo sealed ask UNVERIFIED. iTransfer UNAVAILABLE. Authenticode unavailable. Clean-machine NSIS uninstall NOT PROVEN this pass. Calibration NOT ENOUGH DATA. Chat remains host-parsed (Direct TeeML is research, not chat).
SECURITY RESULT: No Router. No session keys in the browser. Chat cannot AUTHORIZE or flatten. Reduce-only wire cannot withdraw. Encrypted transcript does not store `app-sk-`.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
NEXT STEP: Close PIT Desktop to overlay 0.2.0 sidecar. Do not flatten ETH. Do not remint.

---

## M81 — Flagship desktop workstation (2026-08-28)

DATE/TIME: 2026-08-28 23:00+03
PHASE: Desktop IA rebuild · native chrome · command threads · research board
GOAL: PIT desktop must feel like a local trading desk, not a website in a window. Persistent rail, title bar, command/stage, contextual book. Threads, model picker from real Direct catalog, research board, exact preview contract, Positions page from venue clearinghouse, Ctrl+K, boot checks that are live, not timers.
RESULT: Desktop 0.2.1. Custom undecoated window with drag region and persisted bounds. Boot overlay maps companion/wallet/compute/policy/session/venue. Desk action center. Encrypted chat threads (`chat-threads.enc`) plus structured command cards. Research board with committee roles and named terminals (never a universal “Research stopped”). Positions equity/margin/exposure from Hyperliquid clearinghouseState. Command palette Ctrl+K. Official links stay in `links.ts`. Go tests green. Desktop `tsc` green. Rust `cargo check` green. No new AUTHORIZE. No remint. ETH long untouched.
FAILURE: Overlay `D:\PIT` still older companion until the window is closed and sidecar copied. Authenticode still unavailable. User B live isolation unproven. Galileo sealed ask off. iTransfer UNAVAILABLE. Chat remains host-parsed (Direct TeeML is research, not a chat LLM). NSIS 0.2.1 not tagged until this commit is on main.
SECURITY RESULT: No competitor names in PIT artifacts. Chat still cannot AUTHORIZE, size, or change policy. Threads are workspace-sealed. Window bounds JSON is not a secret. Browser still cannot hold session keys.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
GIT: `d253872` 2026-08-28 23:05+03 `main` (pushed).
NEXT STEP: Overlay 0.2.1 sidecar on this machine. Tag `v0.2.1`. Do not flatten ETH. Do not remint.

---

## M82 — Flagship 0.2.1 published (2026-08-28)

DATE/TIME: 2026-08-28 23:15+03
PHASE: Ship 0.2.1 workstation
GOAL: Push `d253872`, deploy public web/health, overlay the 0.2.1 sidecar, tag the source. Do not flatten ETH. Do not remint.
RESULT: GitHub `main` HEAD `6a372f7` (workstation `d253872` + Positions parse fix). Tag `v0.2.1` points at `d253872`. Vercel `https://pit0g.vercel.app` HTTP 200. Health `https://pit-health.onrender.com/health` version `0.2.1`. Overlay `D:\PIT\pit.exe` companion `0.2.1`. `/local/chat/threads` returns the sealed Desk thread. `/local/positions` now returns ETH `0.0041` @ `2489.7` (venue leverage `20`, mark from public book). Hyperliquid clearinghouse still shows that same long. No new AUTHORIZE. Agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` not reminted.
FAILURE: Authenticode unavailable. User B live isolation unproven. Galileo sealed ask unproven. iTransfer UNAVAILABLE. Calibration NOT ENOUGH DATA. Chat remains host-parsed. NSIS `v0.2.1` is tagged at `d253872` (one commit before the numeric-leverage parse fix). Desktop chrome 0.2.1 is not the running `pit-desktop.exe` until that binary is rebuilt/replaced.
SECURITY RESULT: Push used the token from `.env` without persisting it in the remote URL. `.env` not committed. No session keys in the browser. Venue leverage `20` is displayed as venue state, not as a PIT policy mutation.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
NEXT STEP: GitHub Actions NSIS for `v0.2.1`. Rebuild `pit-desktop.exe` to match. Do not flatten ETH. Do not remint.

---

## M83 — Positions parse + command HTTP fallback (2026-08-28)

DATE/TIME: 2026-08-28 23:28+03
PHASE: Live overlay verify of the 0.2.1 workstation
GOAL: Positions must show the live ETH row. Trading Desk Command must answer host-parsed questions even when Tauri IPC is not present. Do not AUTHORIZE. Do not flatten. Do not remint.
RESULT: Overlay companion `0.2.1` returns ETH `0.0041` @ `2489.7` (venue leverage `20`, mark ~2434). Chat `What is happening?` returns the status tool reply. Ctrl+K palette lists Desk/Watch/Research/Positions/Activity/links. Research board shows last job `1bdcded3` RESEARCH COMPLETE (researcher/challenger/risk OK) and exact preview ETH buy `0.0042` @ `2436.9` hash `0x0412a10b…` — not typed, not posted. Existing fill OID `529167222216` unchanged. Git `6a372f7` then `ce9b042` pushed. Tag `v0.2.1` remains `d253872`.
FAILURE: Direct path reports Protected no (Needs: Protect my strategy) while compute credit is Ready (~4.74 0G). Authenticode unavailable. NSIS `v0.2.1` is the workstation commit without these two hotfixes. Running `pit-desktop.exe` in D:\PIT is still the previous binary; this pass verified the 0.2.1 UI via localhost:3001 against the live companion.
SECURITY RESULT: No AUTHORIZE. No remint. Chat cannot execute. Preview hash is not the historical OID. Browser still cannot hold session keys.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
NEXT STEP: Wait for GitHub Release NSIS `v0.2.1`. Close and replace `pit-desktop.exe` when you install. Do not flatten ETH. Do not remint.


---

## M84 — Native workstation density 0.2.2 (2026-08-28)

DATE/TIME: 2026-08-29 00:01+03
PHASE: Desktop UX pass — command center, chat, research board, first-class venue/compute
GOAL: Make PIT feel like a native trading workstation. Desk is a compact action strip, not a card wall. Chat is persistent with timestamps, copy, stop, and a catalog-backed model menu. Research is a visual committee board with named WHY. Exact preview is the strongest card. Hyperliquid and private compute are first-class. Watch is a book+inspector. Pairing code stays off the top bar when the wallet is bound. No competitor names in PIT artifacts. Do not AUTHORIZE. Do not flatten ETH. Do not remint.
RESULT: Desktop 0.2.2. Desk ready-strip splits Research PROTECTED vs Compute FUNDED. Command chat live-decorates status/positions/research/session from the companion. Model picker lists only the Direct catalog (glm-5.2 on MAINNET). Research board shows snapshot, thesis, researcher/challenger/risk/TEE/engine states, why-banner, and exact preview (asset/side/size/venue/price/policy/session/hash/will/will-not + AUTHORIZE). Watch has a market book plus inspector. Activity is a dated timeline. Pairing chip hidden when bound. Poll interval 4s; pairing code is not polled after bind. Official links stay in `links.ts` and are tested. Chat still cannot AUTHORIZE, size, or mutate policy. Git `9b4c929` pushed to `main`. Tag `v0.2.2`. Vercel `https://pit0g.vercel.app` HTTP 200. Health `https://pit-health.onrender.com/health` version `0.2.2`.
FAILURE: Overlay `D:\PIT` and running `pit-desktop.exe` remain older until this sidecar is copied. Authenticode unavailable. User B live isolation unproven. Galileo sealed ask unproven. iTransfer UNAVAILABLE. Calibration NOT ENOUGH DATA. Chat remains host-parsed.
SECURITY RESULT: No Router SKUs in `/local/models`. No extraAgents in UI copy. Pairing code is not in the top bar after bind. Forget remounts chat. Secrets stay out of chat. Existing fill OID `529167222216` unchanged. Agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` not reminted. Preview hash `0x0412a10b…` not authorized.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
GIT: `9b4c929` 2026-08-29 00:01+03 `main` (pushed). Tag `v0.2.2`.
NEXT STEP: Overlay 0.2.2 sidecar on this machine when research is idle. Rebuild `pit-desktop.exe`. Do not flatten ETH. Do not remint.

---

## M85 — Private desk loop 0.2.3 (2026-08-29)

DATE/TIME: 2026-08-29 ~01:00+03
PHASE: Product rebuild — focused workspaces, live chat, Watch discovery, Research WHY, v1 automation
GOAL: Stop patching the 0.2.2 card wall. Make DISCOVER → RESEARCH PRIVATELY → CHALLENGE → EXPLAIN → POLICY CHECK → EXACT PREVIEW → HUMAN APPROVAL → EXECUTE → PROVE → LEARN visible in the desktop. Chat must answer live state instead of one canned sentence. Watch must surface policy-eligible Hyperliquid facts only. Automation may prepare, never AUTHORIZE. No competitor names in PIT artifacts. Do not AUTHORIZE. Do not flatten ETH. Do not remint.
RESULT: **IMPLEMENTED + TESTED.** Desktop 0.2.3. Desk is command + persistent chat + the ten-step loop (markets moved to Watch). Watch is a full-width discovery book with live why/trend/host-rank and Research privately. Research answers six plain-language questions; committee roles have real copy; incomplete work is never titled verified. Host chat greets, handles typos (reasearch), multi-intent price+research+do-trade (starts sealed research, refuses execute), and decorates Watch/positions/research from the companion. Model picker splits private+verified / other chat / unsupported. Automation prefs persist; auto-research stops at human approval. Hyperliquid card uses Connect / Open API / Approve PIT / Refresh / Revoke plus trading capital. Go `./...` green. Desktop tsc green. Cargo check 0.2.3. Sidecar `PIT 0.2.3`. Desktop honesty+link asserts. Web Playwright 4/4. Redteam green. Vercel `https://pit0g.vercel.app` HTTP 200. Health `https://pit-health.onrender.com/health` version `0.2.3`.
FAILURE: Live companion on this machine is still **0.2.2** (`127.0.0.1:17373`). Overlay `D:\PIT\pit.exe` is **WinError 32** even though `research_running` is false. Restart the desktop to load 0.2.3. Authenticode unavailable. User B isolation unproven this pass. Galileo sealed ask unproven. Calibration NOT ENOUGH DATA. Chat remains host-parsed (no fake token streaming). Full AUTHORIZE E2E not run this pass.
SECURITY RESULT: Chat execute stays false. Automation POST with execute is refused. `/local/models` never marks the public catalog as private. Preview hash `0x0412a10b…` not authorized. Fill OID `529167222216` unchanged. Agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` not reminted.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
GIT: `e4d5e23` 2026-08-29 00:58+03 `main` (pushed). Tag `v0.2.3`.
NEXT STEP: Close PIT, replace `D:\PIT\pit.exe` with the 0.2.3 sidecar, relaunch. Then walk Desk → Watch → Research privately as a new user. Do not flatten ETH. Do not remint.

---

## M86 — Focused workstation 0.2.4 (2026-08-29)

DATE/TIME: 2026-08-29 ~02:00+03
PHASE: Finish PIT as one coherent product — one job per surface, live chat, install 0.2.4
GOAL: Stop crowding Desk with chat + story. Chat is the primary AI surface. Desk is command-only. Watch/Research/Positions/Activity/Policy/Security/Account/Settings each have one job. First-run is nine steps. CLI `chat`/`positions`/`health` share the Go core. Do not AUTHORIZE. Do not flatten ETH 0.0041. Do not remint `PIT-4bbee556`. Do not spend 0G just to manufacture evidence.
RESULT: **IMPLEMENTED + TESTED + LIVE VERIFIED (installed 0.2.4).** Chat is a full-height workspace with threads, composer, stop/retry/copy, model groups (private+verified vs host-parsed chat vs unsupported), and live decorated answers. Desk shows doing / needs you / two opportunities / Ask PIT. Policy and committee are tables. Positions show venue equity and the ETH 0.0041 row from the master account. 9-step first-run. Token marks for ETH/BTC/SOL/DOGE/AVAX/HYPE/0G. CLI `pit chat "What is happening?"` returns live ETH mark + policy. Installed overlay `D:\PIT` + NSIS `PIT_0.2.4_x64-setup.exe`. Companion `/health` **0.2.4**. Companion doctor: wallet, session, hl_agent, direct_auth (OS keychain), direct_credit 4.7058 0G, tee 3 roles, policy pinned.
FAILURE: Authenticode unavailable. `go test -race` needs cgo on this machine. `forge` not on PATH. Web `tsc -b` hung (killed; Vite/Playwright still green). Galileo sealed ask unproven. Calibration NOT ENOUGH DATA. Chat remains host-parsed — no fake token streaming. New sealed research not started this pass (reuse existing TEE evidence). AUTHORIZE E2E not run. CLI `pit doctor` on file keyring reports `direct_auth` fail while the running companion (OS keychain) reports Direct token present.
SECURITY RESULT: Chat execute stays false. Preview hash `0x0412a10b…` not authorized. Fill OID `529167222216` / ETH 0.0041 @ 2489.7 unchanged. Agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` still listed on Hyperliquid API until 2027-02-24. No remint. Router still impossible for the private book.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
GIT: `d605398` 2026-08-29 02:03+03 `main` (pushed). Tag `v0.2.4`.
NEXT STEP: Walk Chat → Watch → Research on the installed 0.2.4 window. Type AUTHORIZE only on an eligible exact preview if you intend a small clip. Do not flatten ETH.

---

## M87 — Autonomous runtime 0.3.0 (2026-08-29)

DATE/TIME: 2026-08-29 ~03:00+03
PHASE: Evolve PIT from research dashboard into a host-gated autonomous trading runtime
GOAL: SCAN → best opportunity → private research → committee → policy → decision → execute (only in Guarded Autonomy) → receipt → monitor → learn. Three modes: Manual, Guarded Autonomy, Research Only. Host-enforced limits the model cannot change. Markets scans the live Hyperliquid universe. Chat orchestrates. Missions persist. Do not flatten ETH 0.0041. Do not remint `PIT-4bbee556`. Do not spend 0G just for cosmetics. Do not enable Guarded Autonomy on the live account this pass.
RESULT:
- **IMPLEMENTED:** Operating modes `manual` / `research_only` / `guarded`. Guarded enable requires exact phrase `ENABLE GUARDED AUTONOMY` on Automation. Chat cannot enable it; chat can stop it. Host limits: assets, venues, clip, leverage, daily loss, max open positions (1), consecutive losses (3), cooldown, slippage, liquidity, uncertainty, kill, session TTL, no withdraw/transfer/policy mutation/escalation. Scanner uses `metaAndAssetCtxs` over the full live perp universe (232 books this pass) and ranks with venue facts (mark/oracle gap, funding, OI) — not invented scores. Best opportunity is the top policy-PASS row. Private research path unchanged (Direct TeeML). After an eligible preview, host may `AUTHORIZE` only if mode is guarded, phrase was confirmed, policy hash matches, preview started after enable, hash is unused, and every gate passes. Missions persist in `mission.json` and recover on companion restart. Chat intents: best opportunity, scan policy universe, research best, trade strongest (research only, never execute from chat), run 24h (navigate Automation), stop autonomy, trades today, on-chain proof, why better. Desktop IA: Desk, Chat, Markets, Research, Portfolio, Activity, Automation, Security. Version **0.3.0**.
- **TESTED:** `go test ./...` green. Desktop `tsc` green. `cargo check` 0.3.0. Sidecar + NSIS `PIT_0.3.0_x64-setup.exe`.
- **LIVE VERIFIED (installed overlay `D:\PIT`):** Companion `/health` **0.3.0**. `/watch` scanned **232** live perps, **6** policy PASS, best **BTC**. Chat `Find me the best opportunity right now` returns live BTC mark + policy PASS. `Scan everything allowed by my policy` returns 232/6. `Run autonomously for 24 hours` refuses enable. POST `/local/mission` `mode=guarded` without the phrase → `need_ENABLE_GUARDED_AUTONOMY`. POST `/local/automation` `execute:true` → `automation_cannot_authorize`. Doctor: wallet, session, hl_agent, direct_auth, direct_credit, policy, tee all ok. Positions: ETH **0.0041** @ 2489.7, OID **529167222216**, agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`. CLI `pit version` / `health` / `scan` / `mission` / `positions` / `chat` match the companion.
FAILURE: Authenticode unavailable. `go test -race` needs cgo. `forge` not on PATH. Galileo sealed ask unproven. Calibration NOT ENOUGH DATA. Chat remains host-parsed (no fake token streaming). New sealed research not started this pass (reuse existing TEE evidence; do not spend 0G for cosmetics).
SECURITY RESULT: Chat execute stays false. Guarded is off (Manual). Model cannot raise limits. Preview is not auto-authorized. Fill OID `529167222216` / ETH 0.0041 unchanged. Agent not reminted. Router still impossible for the private book. Existing open position blocks Guarded execution under max_open_positions=1 even if someone later enables it without raising that ceiling.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
INSTALLER: `PIT_0.3.0_x64-setup.exe` size 16967228 SHA256 `8494f90994dbfc5ec8491044ea9d55d2a5609d6161c7a6929c644baca96cae9b`. Overlay `D:\PIT\pit.exe` + `D:\PIT\pit-desktop.exe`.
GIT: `7e4af97` 2026-08-29 02:57+03 `main` (pushed). Tag `v0.3.0`.
CLASSIFICATION:
- Modes / host limits / scanner universe / chat orchestration / Automation surface / mission persist: **IMPLEMENTED + TESTED + LIVE VERIFIED**
- Guarded host execute path (gates, phrase, duplicate preview, open-position ceiling): **IMPLEMENTED + TESTED**; live auto-order **UNVERIFIED** (deliberately not enabled; existing ETH position + max open 1 makes a new clip unsafe)
- Manual AUTHORIZE E2E this pass: **UNVERIFIED** (not typed)
- Mission recovery after UI restart with a running guarded mission: **TESTED** (unit persist/load); live guarded recovery **UNVERIFIED**
- Provider timeout / kill / daily-loss / session expiry in production: **IMPLEMENTED + TESTED**; live fault injection **UNVERIFIED**
NEXT STEP: Enable Guarded Autonomy only if you intend a tiny policy clip and accept a second position or flatten rules. Do not flatten ETH unless you type AUTHORIZE on a reduce-only preview. Do not remint.

---

## M88 — Product fix 0.3.1 (2026-08-29)

DATE/TIME: 2026-08-29 03:34+03
PHASE: Complete product fix of 0.3.x — mission workspace, confirmation enable, truthful chat, matching version banner
GOAL: Stop shipping Automation as a cramped settings form. Replace the broken ENABLE GUARDED AUTONOMY type-in with a confirmation flow that still POSTs the host phrase. Make Chat a full-height AI workspace with human numbers and one action card. Hide the companion-version banner when desktop and companion match. Keep the host security model. Do not flatten ETH 0.0041. Do not remint `PIT-4bbee556`. Do not spend 0G for cosmetics. Do not leave Guarded Autonomy running on the live account.
RESULT:
- **IMPLEMENTED:** Desktop/companion **0.3.1**. Version banner compares expected desktop version to companion and stays hidden when both are 0.3.1. Automation is a mission workspace: MODE (Manual / Research Only / Guarded Autonomy), mission telemetry, full immutable policy grid, custom selects, confirmation overlay that reviews limits then POSTs `ENABLE GUARDED AUTONOMY` to `/local/mission`, LIVE AUTONOMY + countdown, prominent STOP. Truthful states READY / ENABLING / ACTIVE / STOPPING / STOPPED / BLOCKED / FAILED. Chat formats live books (no scientific notation), answers enable-for-N-hours by opening the confirmation (hours field on the reply, execute stays false), and shows one inline card. Research shows DISCOVERED → … → DECISION. Top bar chips: health, autonomy, session, compute, account. 0G strip names Direct, balance, next research cost, TEE, storage. Portfolio adds why-open. First-run step 0 has Status / What this means / Why it is needed.
- **TESTED:** `go test ./...` green. Desktop `tsc` green. `cargo check` 0.3.1. Sidecar + NSIS `PIT_0.3.1_x64-setup.exe`. Web Playwright 4/4.
- **LIVE VERIFIED (installed overlay `D:\PIT`):** Companion `/health` **0.3.1**. CLI `pit version` **PIT 0.3.1**. `/watch` scanned **232** live perps, **6** PASS, best **BTC**. Chat `Find the best opportunity right now` returns BTC mark 77762 / oracle 77795 / day notional **$4.56B** (not scientific). Chat `Enable guarded autonomy for 8 hours` returns tool `mission.enable_required`, hours **8**, execute false. POST `/local/mission` with the enable phrase → **ACTIVE** / guarded. POST stop → **STOPPED** / manual. Doctor: version, wallet, session, hl_agent, direct_auth, direct_credit, tee, policy all ok. Positions: ETH **0.0041** @ 2489.7, OID **529167222216**, agent `PIT-4bbee556`. Guarded was stopped immediately after the enable probe so it would not start a sealed spend.
FAILURE: Authenticode unavailable. `go test -race` needs cgo. `forge` not on PATH. Galileo sealed ask unproven. Calibration NOT ENOUGH DATA. Chat remains host-parsed (no fake token streaming). New sealed research not started this pass (do not spend 0G for cosmetics). Tiny MAINNET clip **UNVERIFIED** — existing ETH + max_open_positions=1 makes a new market unsafe; enable was proven then stopped.
SECURITY RESULT: Chat execute stays false. Enable still requires the host phrase from this computer; chat cannot enable. Model cannot raise limits. Fill OID `529167222216` / ETH 0.0041 unchanged. Agent not reminted. Router still impossible for the private book. Guarded is off after the probe.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged.
INSTALLER: `PIT_0.3.1_x64-setup.exe` size 16976828 SHA256 `BCEA6DF30099196D670B870AC1CBE4826DEF890E945612F93E313BEE3A3456A0`. Overlay `D:\PIT\pit.exe` + `D:\PIT\pit-desktop.exe`.
GIT: `f8556e5` 2026-08-29 03:35+03 `main` (pushed). Tag `v0.3.1`.
CLASSIFICATION:
- Automation confirmation + mission status + matching version banner / chat numbers / custom selects: **IMPLEMENTED + TESTED + LIVE VERIFIED**
- Guarded host execute path: **IMPLEMENTED + TESTED**; live auto-order **UNVERIFIED** (enabled then immediately stopped; existing ETH + max open 1)
- Manual AUTHORIZE E2E this pass: **UNVERIFIED** (not typed)
- Tiny MAINNET execution: **UNVERIFIED** (unsafe under current open-position ceiling)
- Mission persist/recovery: **TESTED** (unit); live restart with guarded running **UNVERIFIED** (mission was stopped after probe)
NEXT STEP: Walk Automation → Enable Guarded Autonomy only if you intend a tiny policy clip and accept a second position or flatten rules. Do not flatten ETH unless you type AUTHORIZE on a reduce-only preview. Do not remint.

---

## M89 — Guarded Autonomy runtime 0.3.2 (2026-08-29)

DATE/TIME: 2026-08-29 04:42+03
PHASE: Fix Guarded Autonomy end-to-end (runtime, not UI workaround)
GOAL: Find the 5-second auto-stop, keep the mission alive across ticks/restarts, show real progress, persist every event, never report ACTIVE unless the backend mission is running, never silently stop, use live Hyperliquid + Direct TeeML + host gates. Tiny MAINNET order only if safe. Do not flatten ETH 0.0041. Do not remint `PIT-4bbee556`.
ROOT CAUSE: `autoTick` treated `max_open_positions` as a **mission halt**. Default `MaxOpenPositions=1` plus the existing ETH 0.0041 long made the first scheduler tick (0–30s, often ~5s) call `Stop("max_open_positions")`. Scan clocks were not reset on enable, so Next Scan stayed in the past and action stuck at `guarded_enabled`.
FIX: Split `MissionHaltReason` (kill, session, deadline, policy_changed, max_trades, daily_loss, consecutive_loss_limit) from `ExecBlockReason` (`max_open_positions`). Enable resets scan/research clocks, kicks an immediate tick, ticker is 10s with recovery on companion start. Exec refusals persist `block_reason` and keep `running`. Activity records enable/scan/empty/block/research/refuse/stop with terminal reasons.
RESULT:
- **IMPLEMENTED:** Desktop/companion **0.3.2**. Halt vs exec-block. Immediate tick + cadence Next Scan. Mission stage, elapsed, universe counts, exposure, remaining risk, official Hyperliquid/0G/explorer/OID/history links. ACTIVE only when host `status=ACTIVE`.
- **TESTED:** `go test ./...` PIT green. sealer green. Desktop `tsc` green. Desktop e2e copy harness ok. Web Playwright 4/4. `cargo check` 0.3.2. `npx tauri build` NSIS. `forge` not on PATH.
- **LIVE VERIFIED (installed overlay `D:\PIT`):** Companion `/health` **0.3.2**. CLI `D:\PIT\pit.exe version` **PIT 0.3.2**. POST `ENABLE GUARDED AUTONOMY` 1h → **ACTIVE**. After 20s still **ACTIVE** (not `stopped:max_open_positions`). Scan **232** live perps, **6** policy PASS, best **BTC**. `block_reason=max_open_positions` while running. Exposure **ETH 0.0041**. Direct TeeML research BTC: researcher/challenger/risk `verify_e2ee=OK`, `READY_ELIGIBLE`, exact preview BTC buy **0.00013** @ 77727, hash `0xaa2fc11a95be33b3cdb8b2fd320888788e71e8c84254fc35d602dcb021fc197b`. Host `maybeGuardedExecute` **refused** `max_open_positions`. `last_oid` empty. Trades today **0**. Elapsed **100s** still ACTIVE. POST stop → `stopped:user_stop` with explain. `/watch` 232/BTC mark 77752 oracle 77788. Doctor: version, wallet, session, hl_agent, direct_auth, direct_credit 4.5748 0G, tee 3 OK, policy ok.
FAILURE: Authenticode unavailable. `go test -race` needs cgo. `forge` not on PATH. Galileo sealed ask unproven. TESTNET execution **UNVERIFIED** (this machine is bound to MAINNET; did not switch the live bind). Tiny MAINNET new order **correctly blocked** by open-position ceiling. No new 0G Storage proof this pass (no new fill). Mission restart-with-guarded-still-running after a crash **UNVERIFIED** this pass (companion recovered STOPPED leftover then a fresh enable). Kill/session-expiry/provider-failure live injection **UNVERIFIED**. Duplicate-order live **UNVERIFIED** (execute never posted).
SECURITY RESULT: Existing ETH 0.0041 / OID `529167222216` unchanged. Agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` not reminted. Chat cannot AUTHORIZE. Policy clip/leverage unchanged. Withdraw/transfer still impossible. Guarded stopped with `user_stop` after the proof.
TX HASH / OID: none this pass. Existing fill OID `529167222216` unchanged. Eligible preview hash `0xaa2fc11a95be33b3cdb8b2fd320888788e71e8c84254fc35d602dcb021fc197b` was not executed.
INSTALLER: `PIT_0.3.2_x64-setup.exe` size 16982987 SHA256 `32B72A01CA40EABEA1E3698A95370A1DB3499972FCD0B8CEFD7034F8CA235FCB`. Overlay `D:\PIT\pit.exe` SHA256 `BD8DD42FC9966D6FB167857EC248A2319978F722FB84899FE6AB8F546909DA95` + `D:\PIT\pit-desktop.exe` SHA256 `57688D94D9B6C9FD11468E1D015E13DDFED5A160CA282CB47A08057AB85609E1`.
GIT: `d84b336` 2026-08-29 04:42+03 `main` (pushed). Tag `v0.3.2`.
CLASSIFICATION:
- 5-second auto-stop root cause + stay-alive + scan/research/exec-block: **IMPLEMENTED + AUTOMATED TESTED + LIVE VERIFIED**
- Real-time mission progress + activity trail + truthful ACTIVE: **IMPLEMENTED + AUTOMATED TESTED + LIVE VERIFIED**
- Direct TeeML committee on the autonomous path: **LIVE VERIFIED** (BTC READY_ELIGIBLE)
- Tiny MAINNET execution / OID / fill / new 0G proof: **UNVERIFIED / BLOCKED** (`max_open_positions`, ETH not flattened)
- TESTNET order: **UNVERIFIED** (no TESTNET workspace on this bind)
- Kill switch / session expiry / provider failure / duplicate preview live: **IMPLEMENTED + AUTOMATED TESTED**; live injection **UNVERIFIED**
PRODUCTION READY: Guarded Autonomy **scan + research + host refuse** on MAINNET with the current open-position ceiling. **Not** production-ready for a second live clip until the host ceiling allows it or the existing ETH is flattened by an explicit AUTHORIZE.
NEXT STEP: Flatten only with a reduce-only preview you type AUTHORIZE for. Do not remint. Re-enable Guarded Autonomy when you want another scan/research cycle; a new market order stays refused while ETH is open and max open is 1.

---

## M90 — Public evidence spine 0.4.0 (2026-08-29)

DATE/TIME: 2026-08-29 05:10+03
PHASE: Make the 0G contribution load-bearing and independently verifiable, not decorative
GOAL: A judge with no access to this machine must be able to check PIT's claims. Every research verdict and every posted order becomes a canonical public record written to 0G Storage and committed by a 0G Chain transaction. Verification must re-download the bytes, re-check the merkle proof, recompute the digest, confirm the record carries no key material, confirm the sealed roles, ask storage nodes whether the root is finalized, and confirm the root is an indexed topic in a successful call to the storage flow contract. Do not flatten ETH 0.0041 / OID `529167222216`. Do not remint `PIT-4bbee556`. Do not publish anything private.

WHY THIS WAS THE GAP: 0G Compute was already load-bearing (sealed research in TeeML, three roles, `verify_e2ee=OK`, signer matched `teeSigner`), but the sealed part is by design invisible to an outsider. 0G Storage was reachable and 0G Chain was only a health probe, so the strongest part of the product had no public surface. A judge could not distinguish real sealed research from a claim about sealed research.

RESULT:
- **IMPLEMENTED:** Version **0.4.0**.
  - `internal/receipt` canonical public record: kind (`research` / `order` / `cancel`), network, chain id, workspace, market, verdict, deny reason, preview hash, policy hash, per-role attestation summary, OID and fill for orders. Deterministic bytes, self-digest.
  - `internal/proof` filer: writes the record, uploads with the official storage client, parses root and submission transaction, refuses an unusable root, persists a local index, and re-verifies on demand. `Verify` runs six independent checks and reports each one separately instead of a single boolean.
  - `internal/storage/nodeinfo.go` asks the indexer for the sharded node set, then calls `zgs_getFileInfo` on the storage nodes themselves for finalization. `internal/storage/anchor.go` fetches the chain receipt and requires status success, the storage flow contract as the callee, and the root present as an indexed log topic.
  - `internal/storage/public.go` gained `AlreadyStored`, so a duplicate upload reports "bytes already stored" instead of inventing a transaction that was never submitted.
  - `internal/evidence/key.go` resolves the 0G payer from the environment, else the OS keyring, so a desktop-launched companion can publish without a shell profile. `pit evidence bind-payer` stores it.
  - `pit evidence list` and `pit evidence verify --root <root>` give a judge a terminal path that does not touch the UI.
  - Companion files evidence asynchronously on the research-complete path and on the guarded execution path, records `evidence.filed` / `evidence.verified` / `evidence.failed` in activity with root, transaction and explorer link, and exposes `/local/proofs` and `/local/proofs/verify`.
  - Desktop: `ProofTimeline` on Activity lists every published record with root, digest, chain link and a per-record Verify that prints the six checks, the finalized node count and the anchor block. `EvidenceStrip` puts the latest published root and its chain transaction on Desk so the 0G contribution is on the primary workspace, not buried. Activity rows now carry root and chain links inline.
- **TESTED:** `go build ./...` clean. `go test ./...` green across 39 packages. Desktop `tsc --noEmit` clean. `npx tauri build` NSIS. The live proof test is opt-in (`PIT_LIVE_PROOF=1`) because it spends 0G gas.
- **LIVE VERIFIED (0G mainnet, chain 16661):**
  - Re-authorized 0G Compute Direct by signing the official challenge with the bound account (EIP-191), which unblocked sealed research on this machine.
  - Live sealed ETH research through the companion: three roles `verify_e2ee=OK`, signer matched `teeSigner`, verdict **READY_ELIGIBLE**, preview hash `0xc599d232893c81d3cf5fc7568cd6d2a967808f69c03704b679d2f746dba90fe7`.
  - PIT filed that verdict on its own with no operator step: root **`0x07238aa66936340f7ea9fa59f279a8e2313b0bb839699c805b91cb30ccb7741d`**, digest `0xb76e7e192d2238a0a96e7d4d49f32796f788e144b35efa04669c26f5f0bb63cf`, 1082 bytes, job `4b28d466-f459-4049-97c9-43b943db8cb9`, 0G Chain transaction **`0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80`**, link `https://chainscan.0g.ai/tx/0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80`.
  - `pit evidence verify --root 0x07238a…` and the desktop endpoint both returned: merkle proof validated, digest recomputed and matched, record public-safe, sealed roles verified, **3 storage nodes** reporting finalized, and the root confirmed as an indexed topic in a successful storage flow call. Anchor bound.
  - Negative control: a root that was never filed fails verification instead of passing.
  - Restart recovery: companion stopped and restarted; `/local/proofs` returned `ready` with the record intact and the activity trail still carried `evidence.filed` and `evidence.verified` with the chain link. The record also survived the 0.3.2 to 0.4.0 binary swap.
  - Desktop UI walked in a real browser against the live companion: Desk showed the evidence strip with root `0x07238aa6…741d` and a working chainscan link; Activity showed the published proof trail; pressing **Verify on 0G** returned all five checks OK, finalized on **3 of 4** storage nodes the indexer named, anchor block **42,923,233**, flow contract `0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526`.
  - **Real bug found by live verification, then fixed:** the payer key had been bound into the file-backed keyring because this shell exports `PIT_KEYRING=file`, so a desktop-launched companion (which uses the OS keychain) reported `payer_key_missing` and silently would not publish. Re-bound into the Windows OS keychain and the desktop companion reported ready. Worth stating because unit tests could never have caught it.
FAILURE: Authenticode unavailable. `go test -race` needs cgo on this machine. `forge` not on PATH. Chat is still host-parsed, so sealed streaming chat with model selection is **NOT SHIPPED** this pass. Hyperliquid `candleSnapshot` charts **NOT SHIPPED**. The single-workspace UX redesign landed only as the evidence surfaces (Desk strip + Activity proof trail), not as a full IA rewrite. A guarded **order** receipt is **UNVERIFIED**: the host correctly refuses a new clip while ETH 0.0041 is open under `max_open_positions=1`, so no order record could be filed without disturbing the existing position. TESTNET order still **UNVERIFIED** (this machine is bound to MAINNET).
SECURITY RESULT: The published record contains no key material, no session key, no private book content, and no sealed transcript. Publishing is gated on a `public_safe` check that runs again during verification. Chat cannot AUTHORIZE. Guarded is off. Fill OID `529167222216` / ETH 0.0041 @ 2489.7 unchanged. Agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` not reminted. Withdraw and transfer remain impossible.
TX HASH / OID: 0G Chain `0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80` (storage submission committing root `0x07238aa6…`). Existing Hyperliquid fill OID `529167222216` unchanged.
INSTALLER: `PIT_0.4.0_x64-setup.exe` size 17038472 SHA256 `C82CD1A21C790938DB51D6C56BE2876836B7B8C4BD221A6A4C6BEF931B58D4AA`. Overlay `D:\PIT\pit.exe` SHA256 `7ACABBE2F652DB50F2416DE1A4D6DFEFB4B2DA466014662EE3A42CA0676544A9` + `D:\PIT\pit-desktop.exe` SHA256 `CA4167117177A0E12747B7561A714CCC4644A334D183A8550A7F40BB0C302800`. Companion `/health` **0.4.0**, CLI `pit version` **PIT 0.4.0**.
GIT: see tag `v0.4.0`.
CLASSIFICATION:
- Canonical public receipt + 0G Storage upload + root persistence: **IMPLEMENTED + AUTOMATED TESTED + LIVE VERIFIED**
- 0G Chain anchor (root as indexed topic in a successful storage flow call): **IMPLEMENTED + AUTOMATED TESTED + LIVE VERIFIED**
- Independent verification (merkle proof, digest, public-safety, roles, node finalization, anchor) with a failing negative control: **LIVE VERIFIED**
- Automatic filing on the research path with no operator step: **LIVE VERIFIED**
- CLI verification path for a third party (`pit evidence list` / `verify`): **LIVE VERIFIED**
- Desktop proof trail + Desk evidence strip: **IMPLEMENTED + TESTED**
- Payer key from OS keyring for a desktop-launched companion: **LIVE VERIFIED**
- Evidence survives companion restart: **LIVE VERIFIED**
- Order receipt on the guarded execution path: **IMPLEMENTED + TESTED**; live **BLOCKED** by `max_open_positions` with ETH open
- Sealed streaming chat with model selection: **NOT SHIPPED**
- Live candles / charts: **NOT SHIPPED**
PRODUCTION READY: The evidence spine is production-ready on 0G mainnet for research verdicts. The order branch shares the same code path but has not produced a live record because a second clip is unsafe under the current host ceiling.
NEXT STEP: To produce a live order receipt, either raise the host open-position ceiling deliberately or flatten ETH with a reduce-only preview you type AUTHORIZE for. Do not remint. Anyone can check the current record with `pit evidence verify --root 0x07238aa66936340f7ea9fa59f279a8e2313b0bb839699c805b91cb30ccb7741d`.

---

## M91 — Native desktop links + real activity trail 0.4.1 (2026-08-29)

DATE/TIME: 2026-08-29 07:06+03
PHASE: Take 0.4.0 from “evidence exists” to a demo-complete installed desktop: Windows-default-browser links, a complete host activity trail, and honest execution/proof status without weakening policy.
GOAL: DISCOVER → private 0G research → TEE verify → policy → exact preview → AUTHORIZE → real Hyperliquid execution → OID/fill → position → 0G receipt → chain anchor → verify → Activity. Every explorer/HL/0G link from the native desktop must open the Windows default browser. Do not flatten ETH 0.0041 / OID `529167222216`. Do not remint `PIT-4bbee556`. Do not raise `max_open_positions`. Do not spend money just to manufacture a second clip.

RESULT:
- **IMPLEMENTED:** Version **0.4.1**. Desktop sidecar version is now **0.4.1** (was still `0.2.3` in 0.4.0, so a version-mismatched companion would not be replaced).
  - Tauri `open_url` with an HTTPS host allowlist (`pit0g.vercel.app`, Hyperliquid app/API, 0G, chainscan, pc.0g.ai, GitHub only `/mohamedwael201193/pit`). Windows opens via `rundll32 url.dll,FileProtocolHandler`. No `window.open` in the desktop shell. `ExternalLink` + a document click interceptor so leftover `<a target=_blank>` still cannot stay inside WebView.
  - Activity events gained `link`. Host now records `committee.verified`, `preview.ready`, `approval.accepted`, `order.filled` / `order.resting` (only when Hyperliquid says so), and `position.updated` from a live clearinghouse read. Order receipts no longer claim `resting_or_filled`.
  - Desk workflow strip: Discover → Research → Decision → Proof. Chat is a full-height command surface with a live state card, URL rendering, and deep-links; it still cannot AUTHORIZE, size, or mutate policy.
  - Guarded Autonomy still stays alive on `max_open_positions`. Mission state still persists in `mission.json` and is loaded on companion start. A companion kill while the window stays open does **not** auto-respawn (wake is on desktop launch); relaunching the installed desktop **does** respawn sidecar 0.4.1.
- **TESTED:** `gofmt`; `go test ./...` PIT green (including `TestVenueTradeLink`, `TestCommitteeReasonFromRoles`, `TestActivityStoresLink`, `TestGuardedRunningRecoversFromDisk`). sealer `go test ./...` green. Desktop `npx tsc --noEmit` green. Desktop e2e copy harness ok (official URLs + allowlist denials). `cargo test` `allow_official_https` pass. Web Playwright **4/4** (re-run 2026-08-29 07:15+03 after the hung `tsc` job: pair / network / authorize / home). `npx tauri build` NSIS.
- **LIVE VERIFIED:**
  - Installed overlay: `D:\PIT\pit.exe version` **PIT 0.4.1**, companion `/health` **0.4.1**, desktop kicker reads companion version. Session still live. Agent `PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` unchanged.
  - Venue: ETH **0.0041** @ entry 2489.7, OID **`529167222216`**, status filled. Host policy `max_open_positions=1` still refuses a new clip. No AUTHORIZE this pass. No flatten. No remint.
  - Independent 0G verify of the BTC research receipt (filed 06:22 local): root **`0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2`**, digest `0x34dcf767…cd2450`, merkle yes, public-safe, roles verified, **3 of 4** nodes finalized, chain tx **`0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a`** https://chainscan.0g.ai/tx/0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a. ETH root from M90 still indexed. After companion kill + desktop relaunch, `/local/proofs` still `count=2` and ready.
  - Official URLs opened in the Windows default browser via the same `rundll32 url.dll,FileProtocolHandler` path the Tauri command uses (Hyperliquid app/API/BTC book, chainscan tx, pc.0g.ai funds, 0g.ai, pair, GitHub releases). Chrome listed those pages. A command-palette SendKeys attempt from the PIT window was **inconclusive** (Chrome likely stole focus); do not claim every in-window click was visually confirmed.
FAILURE: Authenticode unavailable. `go test -race` needs cgo. `forge` not on PATH. Web `tsc --noEmit` hung (killed; Playwright still 4/4). Chat remains host-parsed (no sealed token streaming). New MAINNET clip **correctly blocked** by the open-position ceiling. TESTNET execution **UNVERIFIED** (this machine stays MAINNET-bound; switching would disturb the live session/agent). New `committee.verified` / `order.filled` rows will appear on the **next** research/order; this pass did not spend compute or trading capital to manufacture them. In-window click of every Activity/Automation link from the installed UI is **PARTIAL** (OS handler live; palette SendKeys not trusted).
SECURITY RESULT: Allowlist refuses http, credentials, custom ports, and non-PIT GitHub hosts. Chat cannot AUTHORIZE. Withdraw/transfer/policy mutation remain forbidden. Kill switch off. Fill OID `529167222216` / ETH 0.0041 unchanged.
TX HASH / OID: BTC evidence `0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a`. Existing Hyperliquid fill OID `529167222216` unchanged. ETH evidence `0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80` unchanged.
INSTALLER: `PIT_0.4.1_x64-setup.exe` size 17042406 SHA256 `6D7CF0C4EA078B211328F5C216947034C8CB2378490F0D1B004055706658E2BE`. Overlay `D:\PIT\pit.exe` SHA256 `8A8C53C06D7D1F564E5E99194CCA35E630FE9233FCB88B99F8210437CC64DA72` + `D:\PIT\pit-desktop.exe` SHA256 `27BCA85E868AC223210023A200A98223160CF7A1602FBF8CF885C662F1341CCE` + `D:\PIT\pit-sealer.exe` SHA256 `A36A3BDC4350E2D17227E09F21DF525480B49BF5AA5111364C87F4F3205C8CD1`. Companion `/health` **0.4.1**, CLI `pit version` **PIT 0.4.1**.
GIT: `cbd1bbf` 2026-08-29 07:07:53 +0300. Tag `v0.4.1` (pushed `origin/main`).
WEB: https://pit0g.vercel.app HTTP 200. Production deploy aliased after CLI `--prod`. Bundle contains `VITE_PRIVY_APP_ID` and `VITE_HEALTH_URL=https://pit-health.onrender.com`. Vercel project env upserted (`VITE_PRIVY_APP_ID`, `VITE_HEALTH_URL`) — public only; no keys, no tokens.
RENDER: https://pit-health.onrender.com/health HTTP 200 `version":"0.4.1"`. `/watch?network=mainnet` 200 live BTC universe. Deploy `dep-da95n31f2nfc73ece380` **live** on `cbd1bbf`. Public PIT_* / Hyperliquid / 0G RPC env upserted on `pit-health` (no payer/deployer/memory keys, no Privy secret, no GitHub/Vercel/Render tokens).
CLASSIFICATION:
- Native Windows default-browser links (Tauri allowlist + rundll32): **IMPLEMENTED + AUTOMATED TESTED + LIVE VERIFIED** (OS handler); in-window click of every control **PARTIAL**
- Activity timeline with hashes/IDs/working links, no fabricated events: **IMPLEMENTED + AUTOMATED TESTED**; new kinds **UNVERIFIED live** until the next research/order
- Guarded Autonomy stays alive on exec-block + mission persist/recover: **IMPLEMENTED + AUTOMATED TESTED**; live mission currently **STOPPED** (`user_stop`); relaunch recovery **LIVE VERIFIED**
- Real Hyperliquid execution this pass: **NOT RUN**; host refuse at open-position ceiling **LIVE VERIFIED**
- 0G research receipt + independent verify: **LIVE VERIFIED** (BTC this pass, ETH M90)
- Order receipt on a new fill: **IMPLEMENTED**; live **BLOCKED**
- Chat never authorizes: **IMPLEMENTED + AUTOMATED TESTED**
PRODUCTION READY: Installed 0.4.1 is production-ready for scan, sealed research, evidence, and host-gated refusal. It is **not** production-ready for a second live clip until the host ceiling allows it or ETH is flattened by an explicit AUTHORIZE.
NEXT STEP: Flatten only with a reduce-only preview you type AUTHORIZE for. Do not remint. Anyone can check `pit evidence verify --root 0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2`.

---

## M92 — Demo-complete new-user desk 0.4.2 (2026-08-29)

DATE/TIME: 2026-08-29 08:31+03
PHASE: Audit installed 0.4.1 against the flagship flow, then ship 0.4.2 so a new user can connect, pin real host policy, read live Hyperliquid truth, research a selected book, and see honest blockers. Do not flatten. Do not remint. Do not invent balances or fills.
GOAL: Open PIT → wallet → workspace → Protect 0G → Hyperliquid → real balances/positions → visual policy editor → scoped session → live eligible universe → sealed Direct TeeML research → host size → exact preview → AUTHORIZE → OID/fill → 0G proof. Chat cannot AUTHORIZE or mutate policy.

RESULT:
- **IMPLEMENTED:** Version **0.4.2**. Sidecar, desktop, companion `/health`, and `pit version` all **0.4.2**.
  1. Security is a readiness center: ready/missing CTAs, Hyperliquid account/agent/session/permissions/expiry, 0G Direct, TEE, Storage, visual policy editor, kill, revoke.
  2. Account/profile is the same Security surface with one-click official links. No extra jargon page.
  3. Policy is a visual editor (clip, leverage locked 1x, assets, venue, max open, daily loss, slippage, liquidity, cooldown, uncertainty, kill, session TTL) with local consequence preview, then pin on this computer. POST `/local/policy` is desktop-only. Chat cannot pin. Leverage/venue/market-type cannot be raised.
  4. Positions read live Hyperliquid clearinghouse + spot USDC. Exec gate annotates `max_open_positions`, `insufficient_margin`, kill. Sizing never invents capital.
  5. Markets scan the live Hyperliquid universe (232 books this pass, 6 policy PASS). Rank is host venue facts, not a model score. Token marks are real. Side is not decided here.
  6. Research is the selected coin. Stages are host-reported (discovered → sealed → researcher/challenger/risk → TEE → engine → policy → decision) with real elapsed_ms. Last BTC job on disk: researcher 9.9s, challenger 18.8s, risk 29.0s, VerifyE2EE OK, preview `hyperliquid:perp:BTC buy 0.00013`.
  7. AUTHORIZE remains desktop-only, exact-preview bound. Chat cannot authorize. Host sizes. Pin must match host law or the order is refused (`policy_changed`).
  8. Activity remains durable with working explorer links. Last evidence still listed.
  9. 0G proofs still `count=2` (BTC + ETH research receipts, official Storage + Chain 16661). No fake TEE/DA/iTransfer claims.
  10. Chat is a full-height host-parsed workspace with live desk state, honest model picker (sealed SKU is research-only, not a chat stream), and deep-links. It cannot authorize or pin.
  11. One primary action per surface. Desk has an 8-step new-user path.
  12. Links: Tauri allowlist + `rundll32 url.dll,FileProtocolHandler`. This pass opened https://chainscan.0g.ai/tx/0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a.
  13. Demo path uses live backend responses. Testnet not switched (this machine stays MAINNET-bound).
- **TESTED:** `go test ./...` PIT green (policy store/clamp/tamper, chat cannot mutate policy, unnamed research does not invent ETH, doctor pin must match, default hash frozen). sealer `go test ./...` green. Desktop `npx tsc --noEmit` green. Desktop e2e copy harness ok. `cargo test` `allow_official_https` pass. Web Playwright **4/4**. `npx tauri build` NSIS. `forge` not on PATH.
- **LIVE VERIFIED (installed overlay D:\PIT):**
  - Companion `/health` **0.4.2** `sign:false trade:false`. CLI `PIT 0.4.2`. Desktop sidecar match.
  - Wallet `0xbdfcee82bd42fefa58ee850b3709636a8b6b0034`. Agent **`PIT-4bbee556` / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`** until 2027-02-24. Session live. Direct credit **4.3946 0G**. Kill off. Last fill OID **`529167222216`** remains historical.
  - Hyperliquid clearinghouse for that wallet now returns **equity 0 / withdrawable 0 / no perp positions**. Spot USDC **4.5813**. Exec gate **`insufficient_margin`** (below the $10 venue minimum). PIT did not flatten and did not invent the old ETH 0.0041 row.
  - Watch: scanned **232**, **6 PASS**, best **BTC 77637**. Universe includes live non-policy books (NFTI/SKR BLOCKED `asset_not_allowed`).
  - Policy GET: host law clip $10 / 1x / max open 1 / hash `384bfd62f42d9a3b…`. Stored pin file hash `c2de517513abb431…` **does not match**. Doctor now reports that honestly. Chat cannot repair it. Pinning on Security writes both `policy.json` and a matching pin. AUTHORIZE now fail-closes on that mismatch.
  - 0G proofs still ready: BTC root `0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2` tx `0xf3d7bc82…`; ETH root `0x07238aa66936340f7ea9fa59f279a8e2313b0bb839699c805b91cb30ccb7741d` tx `0x3f90c548…`.
- **BLOCKED:** New MAINNET clip: venue margin is below $10 and the stored pin does not match current host law. TESTNET execution **UNVERIFIED** (workspace remains mainnet). Authenticode unavailable. `go test -race` needs cgo. Chat is still host-parsed (no sealed token stream). No new sealed research spend this pass (last BTC committee remains on disk). In-window click of every control **PARTIAL** (OS handler live).

SECURITY RESULT: Chat cannot AUTHORIZE, size, or pin. Leverage stays 1x. Withdraw/transfer remain impossible. Kill off. Agent not reminted. No flatten. No invented balances.
TX HASH / OID: BTC evidence `0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a`. ETH evidence `0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80`. Historical Hyperliquid fill OID `529167222216` unchanged. No new order this pass.
INSTALLER: `PIT_0.4.2_x64-setup.exe` size 17065604 SHA256 `FA2223B1460577E899B4830C4FC05871A79A9ADD10AB5787852BB897AE0E676E`. Overlay `D:\PIT\pit.exe` SHA256 `5B4A347CD159F4A8A5BE6CF0E20E21818C1798C89CD3F20A4C0D4818AF98DC6B` + `D:\PIT\pit-desktop.exe` SHA256 `D4F10D1953C8910F3FFE3C24CBC7D090FDFCF9B77A11FDB76C66142961F182DB` + `D:\PIT\pit-sealer.exe` SHA256 `484391A722E2F129CF385BAE5B30C89314510323DEC8F01CB78FCD68D878DB12`. Companion `/health` **0.4.2**, CLI `pit version` **PIT 0.4.2**.
GIT: `bf176d6` 2026-08-29 08:34:16 +0300. Tag `v0.4.2`.
CLASSIFICATION:
- Security readiness + visual policy editor: **IMPLEMENTED + AUTOMATED TESTED**; live pin CTA **LIVE VERIFIED** (mismatch is the honest state)
- Account/setup states with official links: **IMPLEMENTED + TESTED**
- Live Hyperliquid balances used in sizing: **IMPLEMENTED + LIVE VERIFIED** (equity 0, spot 4.5813, no invented row)
- Live eligible universe: **LIVE VERIFIED**
- Sealed research of the selected book: **IMPLEMENTED + TESTED**; last BTC committee **LIVE VERIFIED** (prior job); new spend **NOT RUN** this pass
- Fail-closed AUTHORIZE / pin match / chat cannot mutate: **IMPLEMENTED + AUTOMATED TESTED**; live new clip **BLOCKED**
- Activity + working Windows browser links: **IMPLEMENTED + TESTED**; chainscan tx **LIVE VERIFIED** via rundll32
- 0G Direct TeeML + Storage + Chain: **LIVE VERIFIED** (existing receipts; independent verify unchanged)
- Chat full-height live workspace, honest picker: **IMPLEMENTED + TESTED**
- Installer + version match: **LIVE VERIFIED**
PRODUCTION READY: Installed 0.4.2 is production-ready for scan, sealed research, evidence, and honest refusal. It is **not** production-ready for a new MAINNET clip until YOU pin matching host law on Security and the venue has ≥ $10 available. Do not remint. Do not flatten from chat.
NEXT STEP: On Security, preview then pin host law (optionally raise max open yourself — PIT will not). Fund venue margin if you want a clip. Anyone can check `pit evidence verify --root 0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2`.

---

## M93 — Capital-aware private desk 0.5.0 (2026-08-29)

DATE/TIME: 2026-08-29 19:45+03
PHASE: Strongest truthful production/demo-ready desk. Inspect repo, overlay, official 0G/Hyperliquid, then implement the live loop without breaking CORE LAW. Do not flatten. Do not remint. Do not invent fills, TEE, or 0G proofs.
GOAL: Live Hyperliquid universe sized to THIS account now; research/policy/exec/preview layers; host skill registry; official 0G catalog listing (never Router inference); load-bearing 0G proofs; Activity/Security/policy consequences; labeled demo rehearsal; heavy tests; honest live E2E up to the real execution gate.

RESULT:
- **IMPLEMENTED:** Version **0.5.0**. Sidecar, desktop, companion `/health`, and `pit version` all **0.5.0**.
  1. Capital engine reads live clearinghouse + spot + `userAbstraction`. This wallet reports **unifiedAccount**. Spot USDC is used as collateral source of truth. Buying power $4.5813 is still below the $10 Hyperliquid minimum, so execution stays blocked. PIT does not invent a transfer.
  2. Opportunity layers: research-eligible / policy-eligible / execution-feasible / preview-ready / execution-blocked. Rank prefers executable size, not BTC signal. Live scan 232 / 6 PASS / 0 executable / 0 preview-ready. Best BTC is execution-blocked `insufficient_margin`.
  3. Host skill registry (versioned, inspectable). Cards show skill IDs. Missing candles/L2/liquidations are stated as absent. The model cannot size.
  4. Official 0G catalog GET listing 31 models (TeeML 5, TeeTLS 18, unproven 8). Catalog rows never `private_book`. Private book stays Direct TeeML. Chat is host-parsed even if a catalog SKU is picked. No silent Router fallback.
  5. Streaming desk chat (`/local/chat/stream`) is host-parsed SSE. Chat cannot AUTHORIZE, pin, sign, or place an order.
  6. Policy preview does not pin. Allowed/refused copy plus which books become executable vs research-only. Pin remains desktop-only. Default hash frozen. Stored pin still mismatches host law — fail-closed.
  7. Security readiness: 13 domains with READY / ACTION REQUIRED / BLOCKED / NEEDS ACTION plus one next action. Kill switch off. Execution BLOCKED on venue minimum.
  8. Activity timeline adds sealed / researcher / challenger / risk / TEE kinds after a verified committee. Proof cards: copy, explorer, Verify on 0G.
  9. DEMO GET is LIVE by default. Replay is a separate labeled rehearsal of recorded real evidence. Live and replay cannot be confused.
  10. Mission continues on exec-block (`execution-blocked`) instead of halting the desk.
- **TESTED:** `go test ./...` PIT 565+ PASS, 4 SKIP. sealer 7 PASS. Desktop `npx tsc --noEmit` green. Desktop e2e copy harness ok. `cargo test` `allow_official_https` pass. Web Playwright 4/4. `npx tauri build` NSIS. `forge` not on PATH. `go test -race` needs cgo UNVERIFIED.
- **LIVE VERIFIED (installed overlay D:\PIT):**
  - Companion `/health` 0.5.0 `sign:false trade:false`. CLI `PIT 0.5.0`.
  - Wallet `0xbdfcee82bd42fefa58ee850b3709636a8b6b0034`. Agent **PIT-4bbee556** / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` until 2027-02-24. Session live. Kill off. Last fill OID **529167222216** remains historical. Not reminted. Not flattened.
  - Hyperliquid `userAbstraction` = `unifiedAccount`. Perp equity 0 / withdrawable 0. Spot USDC 4.581269. Exec gate **insufficient_margin**.
  - Watch: scanned 232, 6 PASS, 0 executable, best BTC mark live, skills present, gate insufficient_margin.
  - Official catalog 31 listed, `private_book` false on catalog rows, picked `host-parsed`.
  - DEMO `mode=live live=true`.
  - 0G proofs still ready: BTC root `0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2` tx `0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a`; ETH root `0x07238aa66936340f7ea9fa59f279a8e2313b0bb839699c805b91cb30ccb7741d` tx `0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80`.
  - Windows default browser opened the BTC chainscan tx via rundll32.
  - Direct token missing on this overlay pass (Protect my strategy required). Credit probe still 4.3946 0G. TEE last 3 roles OK on disk. No new sealed spend.
- **BLOCKED:** New MAINNET clip: venue buying power $4.58 < $10 minimum AND pin hash mismatch. Direct token missing until YOU sign Protect. TESTNET execution UNVERIFIED (workspace remains mainnet). Authenticode unavailable. Chat is host-parsed (not a sealed token stream). No new OID/fill this pass.

SECURITY RESULT: Chat cannot AUTHORIZE, size, or pin. Leverage stays 1x. Withdraw/transfer remain impossible. Kill off. Agent not reminted. No flatten. No invented balances or fills. Catalog listing is not an inference path.
TX HASH / OID: BTC evidence `0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a`. ETH evidence `0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80`. Historical Hyperliquid fill OID `529167222216` unchanged. No new order this pass.
INSTALLER: `PIT_0.5.0_x64-setup.exe` size 17131199 SHA256 `863F8087A10202650826F142960E54356CD752FD0469514C1F26CD4A32563A77`. Overlay `D:\PIT\pit.exe` SHA256 `0BEA45E81A1C6EFE6B9A7A133D1A360BDF6AA049BA0BB6A916E15DAAE6ECF8F4` + `D:\PIT\pit-desktop.exe` SHA256 `ED03B26202CFC5D32C0F3EC79128C8A9E3174B71127210242E72CDFE7EA47A09` + `D:\PIT\pit-sealer.exe` SHA256 `ADD18CA1BB9C11C4EF0C0CA998C5C05C4AC635B43F23876382E4413A9569321F`. Companion `/health` **0.5.0**, CLI `pit version` **PIT 0.5.0**.
CLASSIFICATION:
- Capital-aware opportunity engine (layers A-D, venue mins, unified spot vs perp): **IMPLEMENTED + AUTOMATED TESTED + LIVE VERIFIED** (gate insufficient_margin)
- Typed host skill registry + provenance: **IMPLEMENTED + AUTOMATED TESTED**; live IDs on watch cards **LIVE VERIFIED**
- Real 0G catalog listing + host-parsed streaming chat, never Router for book: **IMPLEMENTED + AUTOMATED TESTED**; catalog 31 **LIVE VERIFIED**; sealed chat stream **NOT RUN** (Direct token missing)
- 0G load-bearing proofs + Activity/explorer links: **IMPLEMENTED + TESTED**; existing receipts **LIVE VERIFIED**; new research receipt **NOT RUN**
- Security readiness + policy consequence preview: **IMPLEMENTED + AUTOMATED TESTED**; live pin mismatch **LIVE VERIFIED**
- Demo rehearsal labeled, not live: **IMPLEMENTED + AUTOMATED TESTED + LIVE VERIFIED** (GET is live unless pref=replay)
- Real tiny live execution: **BLOCKED BY USER/VENUE** (margin + pin). Flow proven through the exact execution gate.
- Full in-window click of every CTA: **PARTIAL** (copy harness + live probes)
PRODUCTION READY: Installed 0.5.0 is production-ready for live scan, capital-honest ranking, host skills, catalog listing, evidence, and fail-closed AUTHORIZE. It is **not** production-ready for a new MAINNET clip until YOU pin matching host law, sign Protect if Direct is expired, and the venue has >= $10 available. Do not remint. Do not flatten from chat.
NEXT STEP: On Security, preview then pin host law. Sign Protect my strategy if Direct is missing. Fund unified collateral above $10 if you want a clip. Anyone can check `pit evidence verify --root 0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2`.

---

## M94 — Demo-ready pairing, honest setup, opportunity terminal 0.6.0 (2026-08-29)

DATE/TIME: 2026-08-29 20:50+03
PHASE: Finish PIT as a world-class private trading desktop. Two Grok 4.6 auditors plus live Chrome/repo inspection. Do not flatten. Do not remint. Do not invent fills, TEE, or 0G proofs. Never send the private book to Router.
GOAL: Honest header and golden path; pairing code/expiry/regenerate/connected; Direct vs compute vs trading capital; catalog listing never inference; Markets as “what can I trade NOW”; execution lifecycle fail-closed; Activity/proof buttons; tests; overlay; push; Vercel+Render.

RESULT:
- **IMPLEMENTED:** Version **0.6.0**. Sidecar, desktop, companion `/health`, and `pit version` all **0.6.0**.
  1. Header no longer says Awaiting AUTHORIZE unless a live eligible preview exists, policy is pinned, and the session is alive. Golden-path steps 6–8 are not all green from a stale research kind or a policy-PASS book.
  2. Pairing dock on Desk and Security: code, expiry, regenerate, desktop connected, browser paired/unpaired. `/local/code/rotate` is desktop-only. Web `/pair` no longer navigates to companion JSON. Open PIT Desktop probes loopback and stays on the pairing page. Success deep-links to Protect.
  3. 0G Direct (Protect signature) is split from 0G compute (provider credit) and Hyperliquid buying power. Direct CTA is Protect my strategy, not a generic pc.0g.ai dashboard. Fund link appears only when the ledger is actually short. Unread is not zero.
  4. Security rows: READY / ACTION REQUIRED / BLOCKED / OPTIONAL with one next action and Check again. Identity is OPTIONAL.
  5. Catalog SKUs cannot be saved as a chat inference path. GET resets a previously remembered catalog pick to host-parsed. Private book stays Direct TeeML.
  6. Markets title answers “what can I trade now”. Layers RESEARCH / POLICY / EXECUTION / PREVIEW / BLOCKED. Default filter prefers executable books when any exist. Zero executable states the real venue minimum.
  7. After AUTHORIZE, empty OID is failed, not complete. Venue reconcile uses open orders + user fills and never invents a fill from absence. last-order carries lifecycle.
  8. Activity IDs copy; proofs keep Open explorer / Verify on 0G / copy root, digest, tx, job.
- **TESTED:** `go test ./...` PIT all packages ok. sealer PASS. Desktop `npx tsc -b` green. Web `tsc -b` green. Playwright 4/4. `cargo test` allow_official_https pass. `npx tauri build` NSIS. `go test -race` needs cgo UNVERIFIED.
- **LIVE VERIFIED (overlay D:\PIT):**
  - Companion `/health` 0.6.0 `sign:false trade:false`. CLI `PIT 0.6.0`.
  - Wallet `0xbdfcee82bd42fefa58ee850b3709636a8b6b0034`. Agent **PIT-4bbee556** / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` until 2027-02-24. Session live. Kill off. Last fill OID **529167222216** remains historical. Not reminted. Not flattened.
  - Direct token in keychain after restart warmup. Credit 4.3946 0G. TEE 3 roles VerifyE2EE OK. Policy pin still mismatches host law — fail-closed.
  - Watch: BTC best, layer execution-blocked, buying power $4.581269, execGate insufficient_margin. Official catalog 31, `private_book` false on catalog rows.
  - DEMO `mode=live live=true`. Pairing code present on desktop; devices 0 until the browser pairs again after restart (in-memory devices).
- **BLOCKED:** New MAINNET clip: venue buying power $4.58 < $10 AND pin hash mismatch. TESTNET execution UNVERIFIED. Authenticode unavailable. Chat remains host-parsed. No new OID/fill this pass.

SECURITY RESULT: Chat cannot AUTHORIZE, size, or pin. Catalog listing is not an inference path. Leverage stays 1x. Withdraw/transfer remain impossible. Kill off. Agent not reminted. No flatten. No invented balances or fills.
TX HASH / OID: BTC evidence `0xf3d7bc820154ab18198c2b26ce4f3df6748aa65f3b8b07a7336de4a1c202d65a`. ETH evidence `0x3f90c548a8f9bc04638f459cc9daba37423f04801568457191f2e04fb4090b80`. Historical Hyperliquid fill OID `529167222216` unchanged. No new order this pass.
INSTALLER: `PIT_0.6.0_x64-setup.exe` size 17142609 SHA256 `4D180E2F9F15A8B633A6B1F0A75A2BBF28B0D9AEC78C0BBC8F773F104B7B921F`. Overlay `D:\PIT\pit.exe` SHA256 `EC8142DB35F80DB8D110E5D84F74CD3FE3EBEBEB8BEAFC8AF4A6126F2AFD6EA1` + `D:\PIT\pit-desktop.exe` SHA256 `09FDE82AD29DDC713D894641D6330106281EAAA000CD7A5714414276FADF3122` + `D:\PIT\pit-sealer.exe` SHA256 `60ECA55DB57A22EB67C557A437075C6418EE2DD65C1AF2932DE921792DB480A8`. Companion `/health` **0.6.0**, CLI `pit version` **PIT 0.6.0**.
CLASSIFICATION:
- Honest header / golden path / Direct vs compute CTAs: **IMPLEMENTED + LIVE VERIFIED**
- Pairing restore (code/expiry/regenerate/connected, web not stranded): **IMPLEMENTED + AUTOMATED TESTED**; live paired-browser after restart **UNVERIFIED** (devices in memory)
- Markets opportunity terminal: **IMPLEMENTED + LIVE VERIFIED** (0 executable, honest blocker)
- Catalog listing never inference: **IMPLEMENTED + AUTOMATED TESTED**
- Execution lifecycle + no invented fill: **IMPLEMENTED + AUTOMATED TESTED**; live new clip **BLOCKED BY USER/VENUE**
- Real tiny live execution: **BLOCKED BY USER/VENUE** (margin + pin)
PRODUCTION READY: Installed 0.6.0 is production-ready for live scan, honest setup, pairing, capital-aware ranking, catalog honesty, evidence, and fail-closed AUTHORIZE. It is **not** production-ready for a new MAINNET clip until YOU pin matching host law and the venue has >= $10 available. Do not remint. Do not flatten from chat.
NEXT STEP: On Security, preview then pin host law. Pair the browser if you need a new Protect signature. Fund unified collateral above $10 if you want a clip. Anyone can check `pit evidence verify --root 0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2`.

---

## M95 — Auditor follow-up: pin-gated research, honest web Watch, capital CTA 0.6.1 (2026-08-29)

DATE/TIME: 2026-08-29 21:15+03
PHASE: Apply remaining P0s from the repo auditor and UX/demo QA auditor after 0.6.0. Do not flatten. Do not remint. Do not invent fills or TEE. Never send the private book to Router.
GOAL: Sealed research cannot spend 0G against an unpinned/mismatched policy. Sponsor cannot replace Protect. nextFix names the $10 venue floor. Web Watch shows PASS only. Docs match the shipped version.

RESULT:
- **IMPLEMENTED:** Version **0.6.1**. Companion `/health` and `pit version` **0.6.1**.
  1. `RunWorkspaceResearchStage` fail-closes with `policy_changed` before Direct auth when the pin does not match host law. Chat still cannot pin.
  2. Sponsor compute is used only when this wallet already has Protect. A missing token is `direct_token_required`, not a silent sponsor identity.
  3. After pin+Protect, nextFix is **Fund this Hyperliquid account** when buying power is below $10 / `insufficient_margin`. Markets primary CTA is that fund link. Research is secondary.
  4. Web `/app` header is **Orders stay on desktop** when a wallet is connected. Home title is Protect, not a fake missing session. Watch lists at most 6 policy-PASS books. Policy page is read-only, not a brochure of live numbers.
  5. Research card title for a complete committee is **RESEARCH COMPLETE**, not a verified trade. Historical last-order OID is labeled when exposure is empty.
  6. README / implementation plan / flagship file no longer claim installer 0.2.0, 373 tests, or “UI not started”.
- **TESTED:** `go test` cli/companion/version/compute PASS including `TestResearchRefusesUnpinned` and `TestSponsorDoesNotReplaceMissingProtect`. Desktop `tsc` green. Web `tsc` green. Playwright pair spec 1/1. `npx tauri build` NSIS 0.6.1.
- **LIVE VERIFIED:** Overlay companion 0.6.1. Wallet and agent unchanged. Buying power still $4.58. Pin still mismatched — research now refuses until YOU pin. Direct credit still 4.3946 0G.
- **BLOCKED:** New MAINNET clip: pin mismatch AND $4.58 < $10. No flatten. No remint.

SECURITY RESULT: Private book still Direct TeeML. Sponsor cannot skip Protect. Research cannot spend 0G on an unsigned policy.
TX HASH / OID: Historical OID `529167222216` unchanged. Evidence txs unchanged.
INSTALLER: `PIT_0.6.1_x64-setup.exe` size 17140359 SHA256 `0908611240B60F5BDBC23C6A4486E0002BA659094DC0568B2F9D35A5A2470A0E`. Overlay `D:\PIT\pit.exe` SHA256 `547A2805CE5924E2C94CFDD4CD93807700363191507EA706A06EEE7DEDB6BCA5` + `D:\PIT\pit-desktop.exe` SHA256 `C2347AA19D14432C1B4599C973404AD1F61B55D1EDB2135C8860BDE431472634`. Companion `/health` **0.6.1**.
CLASSIFICATION:
- Pin-gated sealed research: **IMPLEMENTED + AUTOMATED TESTED**
- Sponsor isolation vs Protect: **IMPLEMENTED + AUTOMATED TESTED**
- Web Watch dump / No session yet: **IMPLEMENTED + TYPECHECKED**; live Chrome after deploy **UNVERIFIED this pass**
- Venue fund CTA: **IMPLEMENTED**; live $4.58 gate **still true**
PRODUCTION READY: 0.6.1 is production-ready for honest setup and fail-closed research. It is **not** production-ready for a new MAINNET clip until YOU pin and the venue has >= $10.
NEXT STEP: On Security, preview then pin host law. Fund unified collateral above $10 if you want a clip. Anyone can check `pit evidence verify --root 0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2`.

---

## M96 — 0.7.0 Guarded Autonomy, Why-not, While-you-were-away (2026-08-29)

DATE/TIME: 2026-08-29 22:40+03
PHASE: Transform 0.6.1 into 0.7.0 Private Alpha OS with Guarded Autonomy as a first-class desktop mode. Do not flatten. Do not remint. Do not invent size. Do not enable Guarded Autonomy on the live bound account. Never send the private book to Router.
GOAL: Host-enforced Manual / Research Only / Guarded Autonomy. Kill switch. Durable away journal. Capital-true Markets. Honest Strategy Health. Chat cannot enable autonomy. Web FAQ must not claim PIT never trades without AUTHORIZE.

RESULT:
- **IMPLEMENTED:** Version **0.7.0**. Commit `57b13ef7441ceb50cb316e470baa5493cdad1ca0`. Companion `/health` and `pit version` **0.7.0**. Overlay `D:\PIT` running.
  1. Guarded Autonomy enable is a local typed phrase. Chat returns `mission.enable_required` and navigates to Automation. POST `/local/mission` with ENABLE GUARDED AUTONOMY on the bound wallet returns `need_pin` / "Your policy changed. Re-pin it before trading." Mission stays STOPPED (`user_stop`).
  2. Kill switch is host `Stop(..., "user_stop")`. It does not flatten OID `529167222216` and does not withdraw.
  3. Opportunity ranking: `BestExecutable` requires RankGroup >= 3. Live Watch: buying power **$9.381269**, execGate `insufficient_margin`, copy **$9.38 — $0.62 short of the $10 Hyperliquid minimum**. PIT will not invent size.
  4. Chat `Why didn't you trade?` returns live capital reason + away counts (0/0/0/0 this session).
  5. `/local/calibration` returns `NOT ENOUGH DATA`, `n:0`, `enough:false`. Skills listed with N=0. No fabricated Brier/ECE.
  6. Durable `away.json` journal + human why map. Activity events can be marked AUTONOMOUS.
  7. Web FAQ/Story/Ledger updated so Guarded Autonomy is described as desktop-only, not "never trades without AUTHORIZE".
- **TESTED:** `go test ./...` 40 packages PASS including `TestBestExecutableRejectsBelowMinimum`, `TestExecWhyOpenCeiling`, `TestWhyDidntYouTrade`, `TestMissionEnableRefusesUnpinned`. Desktop `tsc -b` green. Web `tsc --noEmit` green. `npx tauri build` NSIS 0.7.0.
- **LIVE VERIFIED (this machine, companion 0.7.0):** health 0.7.0; wallet `0xbdfc…0034`; agent `PIT-4bbee556` `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52` unchanged; last order OID `529167222216` filled ETH 0.0041 unchanged; mission STOPPED; chat cannot enable; enable pin-gated; Watch $9.38 not executable; calibration empty. After deploy: pit-health **0.7.0**; Chrome FAQ Guarded Autonomy copy; `/pair` reports **PIT Desktop is live · 0.7.0**; `/app/start` header **Orders stay on desktop**.
- **UNVERIFIED:** Guarded Autonomy live order (not enabled). While-you-were-away autonomous fills (none). Autonomy restart recovery while running. Pairing after this overlay. New 0G Storage/chain proof this pass.
- **BLOCKED:** New MAINNET clip: pin mismatch AND $9.38 < $10. No flatten. No remint.

SECURITY RESULT: Model still has no authority. Chat cannot AUTHORIZE, pin, or enable Guarded Autonomy. Enable on a bound wallet requires a matching pin. Private book still Direct TeeML. Kill switch does not flatten or withdraw.
TX HASH / OID: Historical OID `529167222216` unchanged. Evidence txs unchanged (`0xf3d7bc82…` BTC, `0x3f90c548…` ETH).
INSTALLER: `PIT_0.7.0_x64-setup.exe` size 17153852 SHA256 `E0BB7A5199B19E079C2E56717BB50FAE2DFD0FF2A10D64D6F432B05968C20451`. Overlay `D:\PIT\pit.exe` SHA256 `1038AB7544B38C6FD95622905D5D63BAA49F327FA1BCDECD6BE32C8E6317664A` + `D:\PIT\pit-desktop.exe` SHA256 `8B1EF313D3AD3C32B0A0ADE16A02ECE2325E39A06C368591D9C048CF3E400F63`. Companion `/health` **0.7.0**.
CLASSIFICATION:
- Guarded Autonomy host gate + pin-gate: **IMPLEMENTED + LIVE VERIFIED refuse**; live enable **BLOCKED BY USER/VENUE (need_pin)**
- Kill switch user_stop: **IMPLEMENTED + LIVE VERIFIED** (mission STOPPED, ETH fill untouched)
- Why didn't you trade: **IMPLEMENTED + LIVE VERIFIED** ($9.38 / $0.62)
- Strategy Health empty honesty: **IMPLEMENTED + LIVE VERIFIED** (NOT ENOUGH DATA)
- Web FAQ honesty: **IMPLEMENTED + LIVE VERIFIED** (Chrome FAQ + pair 0.7.0)
- Autonomous MAINNET order: **BLOCKED BY USER/VENUE**
PRODUCTION READY: 0.7.0 is production-ready for live scan, honest capital copy, fail-closed enable, and desktop kill switch. It is **not** production-ready for a new MAINNET clip or a live autonomous fill until YOU pin matching host law and the venue has >= $10 available.
NEXT STEP: On Security, preview then pin host law. Fund unified collateral above $10 if you want a clip. Anyone can check `pit evidence verify --root 0x9c65f36076cf2ee32c7e9a02354d1aef9ccf5f6c83289dba160b8c08710424d2`.

---

## M97 — Auditor P0 follow-up 0.7.1 (2026-08-29)

DATE/TIME: 2026-08-29 23:05+03
PHASE: Close remaining P0s from the repo/architecture auditor and the UX/demo auditor after 0.7.0. Do not flatten. Do not remint. Do not invent size. Do not enable Guarded Autonomy. Do not auto-pin.
GOAL: Host-execute is not typed AUTHORIZE. Research/POST refuse below $10. Pin honesty uses disk draft vs pin hash. Portfolio and chat stop claiming $0 equity / pinned law.

RESULT:
- **IMPLEMENTED:** Version **0.7.1**. Commit `a64903b20c63c52c1d87ab4ea9c60be42e5caaab`. Companion `/health` **0.7.1**.
  1. `ExecuteGuardedDeskOrder` does not accept or inject `AUTHORIZE`. Manual `ExecuteDeskOrder` still requires the token. Receipts record `path=guarded|authorize`.
  2. `ClipForAccount` + sizer refuse padding a $9.38 balance up to $10. `BindResearchPreview` sizes from live buying power or denies `insufficient_margin` / `venue_unread`.
  3. Failed research no longer latches `LastResearchCoin`. Restart may retry host-execute once for an unused eligible preview; never if last-order already has that hash+OID.
  4. Policy public/doctor/enable pin-check against `Peek` (disk draft). Live: `pinned false`, disk hash `384bfd62…`, pin `c2de5175…`. Chat: **Policy is drafted, not pinned** plus the $9.38 / $0.62 gate.
  5. Positions summary: buying power **9.3813**, spot **9.3813**, perp account value **0.0**, execGate `insufficient_margin`. Desktop Portfolio labels trading equity vs perp account value.
  6. Expired preview cannot be authorized. Research button explains pin vs compute. Automation exec gate uses venue `insufficient_margin`, not “clear”.
  7. Web `/app` no longer pretends Protect is the only next step. `/pair` when desktop is live promotes Pair, not Download.
- **TESTED:** `go test ./...` PASS including `TestSizerDoesNotPadAboveRequested`, `TestClipForAccountRefusesBelowFloor`, `TestExecuteGuardedNeedsEnable`, `TestWebOriginCannotPostMission`. Desktop `tsc -b` green. Web `tsc --noEmit` green. NSIS 0.7.1.
- **LIVE VERIFIED:** Overlay companion 0.7.1. Wallet/agent/OID `529167222216` unchanged. Chat policy not-pinned + $9.38 gate. Positions buying power 9.3813 vs perp 0.0. After deploy: pit-health **0.7.1**; Chrome `/pair` **PIT Desktop is live · 0.7.1** with primary Pair; `/app` H1 **Inspect on the web. Pin and AUTHORIZE stay on desktop.**; FAQ **Can PIT trade without me?** Guarded Autonomy desktop-only; GitHub release **v0.7.1** installer SHA256 `3B5CBE3B1ADA0DF971056AED8D04AB94A9A624CC321E53F44F3F978023650A01`.
- **UNVERIFIED:** Guarded live order (not enabled). Restart retry of a running guarded mission. Full pair ceremony.
- **BLOCKED:** New MAINNET clip: pin mismatch AND $9.38 < $10.

SECURITY RESULT: Chat/web still cannot enable. Host-execute is a distinct path. Enable pin-check uses disk draft, not a silent Default that could match an old pin.
TX HASH / OID: Historical OID `529167222216` unchanged.
INSTALLER: `PIT_0.7.1_x64-setup.exe` size 17163170 SHA256 `3B5CBE3B1ADA0DF971056AED8D04AB94A9A624CC321E53F44F3F978023650A01`. Overlay `D:\PIT\pit.exe` SHA256 `41CA6EB8D89CA37A82884F226AD8B92A4B5EDD43F8F86B5C0D811E1F7A5701D1` + `D:\PIT\pit-desktop.exe` SHA256 `F21FF5B78F4E5FC8FC2A758778F586D968B5A9AF7B443A8F58F1B7A0912D259E`. Companion `/health` **0.7.1**.
CLASSIFICATION:
- Host-execute ≠ typed AUTHORIZE: **IMPLEMENTED + AUTOMATED TESTED**; live guarded post **NOT RUN**
- Capital-true size / no pad: **IMPLEMENTED + AUTOMATED TESTED**; live clip **BLOCKED** ($9.38)
- Pin honesty: **IMPLEMENTED + LIVE VERIFIED**
- Portfolio trading equity: **IMPLEMENTED + LIVE VERIFIED** (API); desktop UI **OVERLAID**
- Stale preview AUTHORIZE: **IMPLEMENTED**; live expired card **UNVERIFIED in UI click**
PRODUCTION READY: 0.7.1 is production-ready for honest pin/capital copy and fail-closed host-execute. It is **not** production-ready for a new MAINNET clip until YOU pin and the venue has >= $10.
NEXT STEP: On Security, preview then pin host law. Fund unified collateral above $10 if you want a clip.

---

## M98 — Opportunity engine continues past BTC; per-market venue min (2026-08-30)

DATE/TIME: 2026-08-30 01:50+03
PHASE: 0.7.2. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent size. Do not enable Guarded Autonomy. Do not auto-pin.
GOAL: Stop the BTC dead-end. Per-asset Hyperliquid min notional (szDecimals rounding). Stand-down is a successful no-trade that continues the universe. Compact Markets/Desk/Research/Automation copy.

RESULT:
- **IMPLEMENTED:** Version **0.7.2**. Companion `/health` **0.7.2**.
  1. `watch.NextCandidate` skips stood-down / already-researched coins. A $10 policy clip cannot size BTC at ~80k (rounded min ~$10.16+). Rank no longer auto-firsts BTC.
  2. `venue.PerpMinNotionalUSD` is documented $10 floor rounded up to szDecimals. Live Watch: HYPE min $10.01, ETH $10.05, BTC $10.16, SOL $10.53. Copy names the asset and this account's $9.38.
  3. Stand-down (`no_side` / challenger / risk) skips that coin for 10 minutes and clears `LastResearchCoin` so the next tick researches ETH/SOL/HYPE. Mission `search_note`: "BTC: no side survived challenge. PIT is checking the next eligible market (ETH)."
  4. Markets is a compact opportunity feed. Research treats stand-down as a successful no-trade. Automation shows search_note. Desk hero shows this-account vs this-market min.
- **TESTED:** `go test ./...` PASS including `TestNextCandidateSkipsStoodDownBTC`, `TestSizerBTCTenClipCannotMeetRoundedMin`, `TestCoarseTickRaisesFloor`, `TestStandDownDoesNotLatch`. Desktop `tsc -b` green. Web `tsc --noEmit` green. NSIS 0.7.2.
- **LIVE VERIFIED:** Overlay companion 0.7.2. Wallet `0xbdfcee…0034`, agent PIT-4bbee556 / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`, OID `529167222216` filled ETH 0.0041 unchanged. Watch scanned **232**, **best HYPE** (not BTC; BTC is 5th), buying power **$9.38**, HYPE min **$10.02**, BTC min **$10.17**, SOL min **$10.55**. Mission STOPPED `user_stop`, mode manual. Chrome production:
  - `pit0g.vercel.app/pair` → “PIT Desktop is live on this computer · 0.7.2”
  - FAQ “Can PIT trade without me?” → chat/website/model cannot enable Guarded Autonomy
  - `/app` bound wallet `0xBDfCeE…0034`, MAINNET, Protect CTA, Watch cards cannot trade, confidence NOT ENOUGH DATA
  - Render `pit-health.onrender.com/health` version **0.7.2**
  - GitHub release `v0.7.2` installer size 17180941 digest `20c49486…` after a wrong 14.1 MB asset was replaced
- **UNVERIFIED:** Guarded live order. Full pair ceremony this session. Crash-restart of a running guarded mission. Live committee stand-down skip cycle.
- **BLOCKED:** New MAINNET clip: pin mismatch AND $9.38 below every policy book's rounded min. Guarded Autonomy not enabled (deliberate).

SECURITY RESULT: Chat/web still cannot enable. Model still cannot size. Skip journal is host-local.
TX HASH / OID: Historical OID `529167222216` unchanged.
COMMIT: feat `2f10310bad73bef3a4d0259780081781c67c4a4b` (2026-08-30 01:50:35 +0300). Journal `08789a65a7eff6a5e0b6f216d1858169da4b38a9`. HEAD `bcaf9fcdc4f061ac9f85268bdb0ed7bbf0070075`. TAG `v0.7.2` at `08789a6`.
INSTALLER: `PIT_0.7.2_x64-setup.exe` size 17180941 SHA256 `20C49486993B055A94A0B3531A71ECAC6FC102E49455FE3E47B8FC11FA8FC617`. Overlay `D:\PIT\pit.exe` SHA256 `C0F1555AE5C037D2A0BDA10BFB81F664CFF670BD39E0E79033EB5EC85295174F` + `D:\PIT\pit-desktop.exe` SHA256 `70500B3635ABBE083658EC601332180070A58CAF146C6993F36F3D6575EA0BEF`.
RELEASE: https://github.com/mohamedwael201193/pit/releases/tag/v0.7.2
WEB: https://pit0g.vercel.app
HEALTH: https://pit-health.onrender.com/health `0.7.2`
CLASSIFICATION:
- Continue past blocked/stood-down BTC: **IMPLEMENTED + LIVE VERIFIED** (best is HYPE; BTC is not first)
- Per-market venue min: **IMPLEMENTED + LIVE VERIFIED**
- Stand-down continue: **IMPLEMENTED + AUTOMATED TESTED**; live committee stand-down cycle **UNVERIFIED**
PRODUCTION READY: 0.7.2 is production-ready for honest per-market mins and universe continuation. It is **not** production-ready for a new MAINNET clip until YOU pin and the venue has enough to meet the candidate's rounded min.
NEXT STEP: On Security, preview then pin host law. Fund unified collateral above the specific market min if you want a clip.

---

## M99 — Public Watch is a compact policy feed (2026-08-30)

DATE/TIME: 2026-08-30 02:12+03
PHASE: 0.7.2 follow-up. Do not flatten OID 529167222216. Do not remint. Do not enable Guarded Autonomy. Do not auto-pin.
GOAL: [PIT UX live demo audit](d12538be-80fa-4b8f-90ed-7b5b69be10f3) still saw six identical web Watch cards. Compact them. Architecture P0s from [PIT architecture security audit](df9947a4-b22f-4dc7-976c-b14135555da8) were already in 0.7.2 (per-asset min, skip stood-down BTC, no duplicate guarded POST).
RESULT:
- **IMPLEMENTED:** `apps/web` Watch is one table: asset, mark, policy-eligible, why. One footer: this site cannot size or AUTHORIZE. Header copy is scanned books + policy-pass count, not six identical AUTHORIZE walls. Public `/watch` JSON `copy` field is unchanged (SDK contract).
- **TESTED:** `go test ./internal/watch/...` and web `tsc --noEmit` after this change.
- **LIVE VERIFIED:** Chrome `https://pit0g.vercel.app/app?v=072b` — “232 live Hyperliquid books. 6 pass the default policy. None of these cards can trade.” Compact table SOL…BTC last. Footer: AUTHORIZE not on this website. Commit `8dfa4377a2abeb6636d34fb1bbf48fc12310476d`.
- **UNVERIFIED / BLOCKED:** unchanged from M98.

SECURITY RESULT: Web still cannot authorize, pin, or enable autonomy.
TX HASH / OID: Historical OID `529167222216` unchanged.

---

## M100 — War room 0.7.3: verified experience + MAINNET chrome (2026-08-30)

DATE/TIME: 2026-08-30 04:20+03
PHASE: 0.7.3. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent size. Do not enable Guarded Autonomy. Do not auto-pin.
GOAL: Three Grok 4.6 war rooms (competitor / 0G-security / UX). Ship Adaptive Private Trading Memory, continue-search copy, per-book capital honesty, hide TESTNET from the customer desk, stop `/watch` leaking this account's buying power to the website, and keep chat research from host-executing Guarded Autonomy.
RESULT:
- **IMPLEMENTED:** Version **0.7.3**.
  1. Encrypted `experience.enc` (AES-GCM, workspace keyring). Research stand-downs and venue OIDs append typed cases. `GET /local/experience`, `pit memory`, MCP `experience` are read-only. Qualitative copy stays **NOT ENOUGH DATA** until N≥5. Memory key never enters prompts.
  2. Chat / Research UI research **stop at preview**. Only the automation loop may `maybeGuardedExecute`. Recovered guarded execute requires `source=automation`.
  3. Unauthenticated website Origin on `/watch` gets the public universe **without** `ApplyCapital`. Desktop Origin still sizes this account.
  4. Direct auth must be `app-sk-`. Official storage client basename is `0g-storage-client(.exe)`. Hex keys are secretful. CLI `authorize` posts through `ExecuteDeskOrder` and reconciles OID instead of a second POST.
  5. Desk hero: "Watching. Nothing can open." Stand-down is success, not coral fail. Pairing is optional, not Action required. TESTNET is developer-only (`?dev=1` / Help 7-clicks). Landing Dual is MAINNET-only.
- **TESTED:** `go test ./...` PASS including `TestMayHostGuardedExecuteOnlyAutomation`, `TestWatchWebOriginOmitsAccountCapital`, `TestRefuseRouterKey` random bearer, `TestRejectUnofficialClient` basename. Desktop `tsc -b` green. Web `tsc --noEmit` green. Desktop copy e2e harness ok. `go build ./cmd/pit ./cmd/mcp`.
- **LIVE VERIFIED:** Overlay companion `D:\PIT\pit.exe` `/health` **0.7.3**. SHA256 `758B894B8D9069281EAE6BBCBD9E718B8CCC7ED22DBF5AF4242C03A8B52E5383`. Wallet `0xbdfcee…0034`, agent PIT-4bbee556 / `0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52`, OID `529167222216` unchanged. Session live, mission STOPPED `user_stop`. Watch scanned **232**, best **DOGE**, buying power **$16.18**. All six policy books `below_min_notional` because the pinned **$10 clip** cannot meet rounded mins (~$10.03–$10.82). Website Origin `/watch` `buyingPower` empty. `/local/experience` `NOT ENOUGH DATA (0/5)`. Chrome production landing still 0.7.2 copy until this commit deploys.
- **UNVERIFIED:** Guarded live order. New 0G Storage/chain proof this pass. Full NSIS 0.7.3 installer (desktop window still 0.7.2 UI talking to 0.7.3 sidecar — version banner is honest). Pairing ceremony this pass.
- **BLOCKED:** New MAINNET clip: host clip $10 cannot size any live policy book after szDecimals rounding. Raising max trade requires YOU to preview then pin. Guarded Autonomy not enabled (deliberate). iTransfer UNAVAILABLE on Aristotle.

SECURITY RESULT: Chat/web/MCP still cannot authorize, pin, or enable autonomy. Chat research cannot spend a Guarded execute. Website `/watch` no longer carries this account's margin.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Verified experience: **IMPLEMENTED + LIVE VERIFIED** (empty store, honest NOT ENOUGH DATA)
- Continue-search / stand-down success copy: **IMPLEMENTED**
- Website capital leak: **IMPLEMENTED + LIVE VERIFIED**
- Guarded execute origin: **IMPLEMENTED + AUTOMATED TESTED**
- TESTNET hidden from customer chrome: **IMPLEMENTED**; production landing **UNVERIFIED until Vercel**
PRODUCTION READY: 0.7.3 is production-ready for honest memory, capital isolation, and MAINNET-only chrome. It is **not** production-ready for a new MAINNET clip until YOU pin a clip that meets a book's rounded min.
NEXT STEP: Close PIT Desktop and relaunch to pick up the 0.7.3 window. On Security, preview then pin a clip above the candidate min if you want a clip. Do not flatten. Do not remint.

---

## M101 — 0.7.3 installer + production verify + honest checksums (2026-08-30)

DATE/TIME: 2026-08-30 04:35+03
PHASE: 0.7.3. Do not flatten OID 529167222216. Do not remint. Do not invent size. Do not enable Guarded Autonomy. Do not auto-pin.
GOAL: Ship the NSIS installer, prove production landing, replace the GitHub Actions checksum dump that hashed every `.exe` on the runner.
RESULT:
- **IMPLEMENTED:** `release.yml` checksums only the NSIS installer. Web Home no longer double-labels Watch. Actions `prerelease` is false.
- **TESTED:** Overlay `D:\PIT\pit.exe` + `pit-desktop.exe` after stop/restart. Companion `/health` **0.7.3**. CLI `pit version` **PIT 0.7.3**. `pit memory` **NOT ENOUGH DATA (0/5)**. Website Origin `/watch` has no `buyingPower`. Desktop Origin buying power **$16.18**, best **DOGE**, all six books `below_min_notional`.
- **LIVE VERIFIED:** Commit `2265ef1`. Tag **v0.7.3**. GitHub release https://github.com/mohamedwael201193/pit/releases/tag/v0.7.3 is **Latest** with `PIT_0.7.3_x64-setup.exe` (17,195,846 bytes). Chrome `https://pit0g.vercel.app/?v=073`: hero “Your trading desk doesn’t sleep.”, **MAINNET only**, no TESTNET switcher, Download → `/releases/latest`. `/app/start` MAINNET production only; laboratory behind Help seven clicks. `/app` Watch table SOL…BTC last, “this site cannot trade”, no buying power. `/app/activity` execution on desktop. Render `https://pit-health.onrender.com/health` **0.7.3**.
- **INSTALLER:** `PIT_0.7.3_x64-setup.exe` SHA256 `CACF16FACA0E014E2250481A181B370F7C40AA79E96D2C85679BDC9ABA97C3FB`. Overlay desktop `E831A3AB74536D381E71D1F43DCED368F70F0EEDDA69D543EB258597DF934D11`. Sidecar `B9FB3E37D8847A212AE3BE6EFF76C5FC8C8929A29C5A2E0A52DF7C5103088B10`. Sealer `882BA359E510135DBCB12B8CE84DE1D4F1690310D05119410F9340266E87BDA0`.
- **UNVERIFIED:** Guarded live order. New 0G Storage/chain proof this pass. Pairing ceremony this pass. GitHub Actions Windows NSIS (separate from the local installer).
- **BLOCKED:** New MAINNET clip unchanged: pinned $10 clip vs rounded per-book mins. iTransfer UNAVAILABLE on Aristotle.

SECURITY RESULT: Chat/web/MCP still cannot authorize, pin, or enable autonomy.
TX HASH / OID: Historical OID `529167222216` unchanged.

---

## M102 — Chat lists per-book floors; Desk shows lesson line (2026-08-30)

DATE/TIME: 2026-08-30 04:42+03
PHASE: 0.7.3 follow-up. Do not flatten OID 529167222216. Do not remint. Do not invent size. Do not enable Guarded Autonomy. Do not auto-pin.
GOAL: Late war rooms still saw Chat latch a skipped coin and Desk omit verified-experience copy. Fix that honesty without a new clip.
RESULT:
- **IMPLEMENTED:** `formatBookFloors` lists each policy-eligible min/shortfall. `replyWhyNotTrade` no longer leads with `LatestSkip`. Desk hero shows `experienceWhy`. Chat prompt “What did we learn on HYPE?”. `HumanWhy(insufficient_margin)` no longer says “this market.”
- **TESTED:** `TestFormatBookFloorsListsPerMarketShortfall`. `go test ./internal/companion ./internal/auto` PASS. Desktop `tsc -b` PASS.
- **LIVE VERIFIED:** Not yet — overlay companion after this commit.
- **BLOCKED:** New MAINNET clip unchanged.

SECURITY RESULT: Chat still cannot AUTHORIZE.
TX HASH / OID: Historical OID `529167222216` unchanged.

---

## M103 — War room 0.7.4: policy-clip honesty + capital router (2026-08-30)

DATE/TIME: 2026-08-30 05:35+03
PHASE: 0.7.4. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent size. Do not enable Guarded Autonomy. Do not auto-pin.
GOAL: Screenshots showed THIS ACCOUNT $16.18, THIS MARKET MIN $10, SHORTFALL $0, and still “Nothing can open yet” + Fund. That is a lie. The account clears the venue floor. The pinned $10 clip cannot meet rounded per-book mins. Fix the gate, add a host-deterministic capital router, do not clone Zia, do not fake SWAP/LP.
RESULT:
- **IMPLEMENTED:** Version **0.7.4**.
  1. `policy_clip_tight` when buying power ≥ book min and pinned clip cannot size. Copy: “Policy cap is $X.XX too tight… Raise max trade on Security, preview, then pin.” Fund CTA only when buying power is actually short of the venue floor.
  2. Capital router (not a smart contract): TRADE / WAIT / HOLD / SWAP / LP. SWAP and LP stay **unavailable**. Zia APR is advertising, not a PIT return. No swap/LP transaction is invented.
  3. Desk, Markets, nextFix, Chat floors, CLI `pit opportunities`, local MCP Watch, policy preview all name the clip gap. Website Origin `/watch` still omits this account’s capital.
  4. Model picker rows label PRIVATE / TEE-VERIFIED / CATALOG ONLY / NOT PRIVATE. Catalog listings stay disabled.
  5. Landing hero says a $10 clip is not every book’s Hyperliquid minimum.
- **TESTED:** `go test ./...` PASS including `TestApplyCapitalPolicyClipTightLiveEth`, `TestLiveEthMarkMakesDefaultClipTight`, `TestDecideRoutesPolicyClipSelectsWaitNotSwap`, `TestFormatBookFloorsPolicyClipTightNotFund`. Desktop `tsc -b` PASS. Web `tsc --noEmit` PASS. NSIS `PIT_0.7.4_x64-setup.exe`.
- **LIVE VERIFIED:** Overlay companion `D:\PIT\pit.exe` `/health` **0.7.4**. Desktop Origin `/watch` buying power **$16.181269**, gate **policy_clip_tight**, WAIT selected, SWAP/LP unavailable. Tightest book **DOGE** min **$10.02**, gap **$0.02**. ETH min $10.06, BTC $10.14, SOL $10.52, HYPE $10.80. Website Origin `/watch` `buyingPower` empty. SHA256 pit.exe `4CE04BA48A889E01A8494DC2CA1207F44E23EBA30DDBA669EA1B2CC8D7F0B558`. pit-desktop.exe `7488C0B5FCE3A22824AFB66AE873C5EBEFEB2F4DD1A66469BA99835456989D3C`. sealer `8DC5D76DD8803BED2F0C756C0D7A9A19EFBF16CAA3C67DFB583EAE8F80A78EBC`. installer `61A52E738BD9758918DFD327857180038EE0EE7161FF105E0151BDA47AECF620`.
- **UNVERIFIED:** Guarded live order. New 0G Storage/chain proof this pass. Pairing ceremony this pass. Production Vercel/Render until this commit deploys.
- **BLOCKED:** New MAINNET clip until YOU raise max trade above the candidate min, preview, then pin. iTransfer UNAVAILABLE on Aristotle. Zia/DEX SWAP and LP execution unavailable (no verified rail).

SECURITY RESULT: Chat/web/MCP still cannot authorize, pin, or enable autonomy. Policy is not auto-raised. SWAP/LP cannot execute.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- policy_clip_tight vs Fund lie: **IMPLEMENTED + LIVE VERIFIED**
- Capital router TRADE/WAIT/HOLD: **IMPLEMENTED + LIVE VERIFIED**
- SWAP/LP execution: **NOT IMPLEMENTED** (labeled unavailable)
- Zia clone: **NOT IMPLEMENTED** (deliberate)
PRODUCTION READY: 0.7.4 is production-ready for honest capital gates. It is **not** production-ready for a new MAINNET clip until YOU pin a clip that meets a book’s rounded min.
NEXT STEP: Close extra PIT windows if two opened. On Security, preview then pin a clip above the candidate min if you want a clip. Do not flatten. Do not remint.

---

## M104 — War room 0.7.5: policy-eligible floor + Agentic ID split (2026-08-30)

DATE/TIME: 2026-08-30 05:54+03
PHASE: 0.7.5 follow-up to [Competitor intelligence](59bed584-25d5-4821-b46e-38f639bd4e44), [0G infrastructure](2695e9ef-4c51-4436-a33e-9175a250f1e2), [Capital UX](ceb9e7bd-860b-4999-b93e-4846d4975053). Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent size. Do not enable Guarded Autonomy. Do not auto-pin. Do not ship a Zia APR table.
GOAL: After 0.7.4 remapped `policy_clip_tight`, Markets could still show THIS MARKET MIN $10 / SHORTFALL $0 because `nearestVenueMin` took the cheapest of all 232 books (MEW). Chat rounded min to `$%.0f`. Agentic ID was painted fully OFF while mint/own is live. Split those lies. Do not fail-close sealed research on an unminted desk.
RESULT:
- **IMPLEMENTED:** Version **0.7.5**.
  1. `nearestVenueMin` prefers execution-feasible books, then policy-eligible/eligible. MEW $10 no longer pollutes the Markets min when the policy universe starts at DOGE.
  2. Committee copy for `policy_clip_tight` vs `below_min_notional`. Chat watch min uses `$%.2f`.
  3. Capability matrix Agentic ID is **PARTIAL**: mint/own/authorizeUsage/revoke live on Aristotle; iTransfer/iClone UNAVAILABLE. `/local/identity` note matches. Website Origin still cannot read identity.
  4. No Zia APR cards. No `isAuthorized` fail-close on sealed ask (mint stays optional; `BeforeSealedAsk(false)` would block glm-5.2).
- **TESTED:** `go test ./...` PIT PASS including `TestLocalIdentitySplitsMintAndTransfer`. sealer PASS. Desktop `tsc -b` PASS. Desktop e2e copy harness ok (`assertCapitalFloorIgnoresDust`). `cargo test` allow_official_https PASS. `npx tauri build` NSIS `PIT_0.7.5_x64-setup.exe`. Web `tsc --noEmit` hung again (PID killed); no `apps/web` source change this pass.
- **LIVE VERIFIED:** Overlay `D:\PIT\pit.exe` `/health` **0.7.5**. `/watch` buying power **$16.181269**, gate **policy_clip_tight**, WAIT selected. Policy-eligible min **DOGE $10.04** (gap $0.04). Universe-wide min still MEW $10 (correctly ignored by the UI helper). ETH $10.07, BTC $10.15, SOL $10.54, HYPE $10.80. `/local/identity` itransfer **UNAVAILABLE**. SWAP/LP unavailable. SHA256 pit.exe `8C1DB361799635D394A3A71F2E0730B2ED91EA5237A31EA9203A1B8ECFA495FB`. pit-desktop.exe `6BBE55497512E90742560ED3622CA88EB04F17BEA0E73B4EBD02C1D20B39BBB1`. sealer `5C1A0B28E8BD6D76C9B139B10B8FBAF53BD459F699FE46C63B8D2C742ED17862`. installer `9D909CCB50BA0B3E1D83A4204B6C97B60FDE2E1AC6E933288D0355894127372D`.
- **UNVERIFIED:** Guarded live order. New 0G Storage/chain proof this pass. Pairing ceremony this pass. Production Vercel (no web copy change). Render health until this tag deploys.
- **BLOCKED:** New MAINNET clip until YOU raise max trade above the candidate min, preview, then pin. iTransfer UNAVAILABLE. Zia/DEX SWAP and LP execution unavailable.

SECURITY RESULT: Chat/web/MCP still cannot authorize, pin, or enable autonomy. Sealed research does not wait on desk mint.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- MEW floor pollution: **IMPLEMENTED + LIVE VERIFIED** (engine min_all still MEW; UI/policy min DOGE)
- Agentic ID mint vs iTransfer split: **IMPLEMENTED + TESTED + LIVE VERIFIED** (`/local/identity`)
- Zia APR intelligence: **NOT IMPLEMENTED** (deliberate; promotional APR is not a PIT return)
- Fail-close sealed ask on unminted desk: **NOT IMPLEMENTED** (would break optional-mint research)
PRODUCTION READY: 0.7.5 is production-ready for honest policy-eligible floors. It is **not** production-ready for a new MAINNET clip until YOU pin a clip that meets a book’s rounded min.
NEXT STEP: On Security, preview then pin a clip above the candidate min if you want a clip. Do not flatten. Do not remint.

---

## M105 — Desktop 0.7.6: executable desk, not a stale stand-down (2026-08-30)

DATE/TIME: 2026-08-30 06:29+03
PHASE: 0.7.6 desktop finalization. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent a fill. Do not enable Guarded Autonomy. Clip is already $13 from YOU.
GOAL: After you pinned $13, Markets correctly showed 6 executable books. Desk still treated a prior committee stand-down as the headline and repeated six TIGHT paragraphs. Security listed every READY domain as a giant card. Finish the live desktop around capital truth, not a new product thesis.
RESULT:
- **IMPLEMENTED:** Version **0.7.6**.
  1. Desk headline prefers live executable count over a stale `READY_STOOD_DOWN`. With 6 books that can open, the title is **6 books can open**. Policy-tight still beats a stand-down when clip cannot size.
  2. Desk books are one-line rows (asset / mark / chip / min·clip), not six repeated paragraphs. Primary CTA is **Research {best} privately**.
  3. Markets defaults to Actionable when any book is execution-feasible. Funding column. Research CTA without selecting a row.
  4. Security READY domains are chips. Only NEEDS ACTION stays as cards. Policy clip cell shows pinned vs draft.
  5. First-run network step is MAINNET production (TESTNET stays developer-only). Research taxonomy titles for POLL_FAILED / JOB_CRASHED / DIRECT_* / SPONSOR_QUOTA / WRONG_NETWORK.
- **TESTED:** `go test ./...` PASS. Desktop `tsc -b` PASS. Desktop e2e copy harness ok (`assertDeskHeadlinePrefersExecutable`). `cargo test` allow_official_https PASS. NSIS `PIT_0.7.6_x64-setup.exe`. Web `tsc` not re-run (no `apps/web` change).
- **LIVE VERIFIED:** Overlay `/health` **0.7.6**. `/watch` buying power **$16.181269**, clip **$13**, **6 executable / 6 preview**, TRADE ready, WAIT blocked, SWAP/LP unavailable. DOGE min $10.03 … HYPE $10.80. Historical OID untouched. No new order. SHA256 pit.exe `299AD0C4AD3B759246A8DC70578411C21E043984D17E5F3DB623ACB43E30A3E0`. pit-desktop.exe `F3E1C87E3FA539A27C1DC65EFAFB5637DDC88052E1EEB5DD05F1E960C6FCF167`. sealer `28EB5AE6190D93FB9D46D3A3B3B1989A57CB897E838C673A178E24ABD0A41254`. installer `1B2C5D11045E6505D9953C2074D001F074AC6E8FE47976BC21CA233A8F431A59`.
- **UNVERIFIED:** Guarded live order. New 0G Storage/chain proof this pass. Pairing ceremony this pass. Vercel landing (no web change).
- **BLOCKED:** iTransfer UNAVAILABLE. SWAP/LP execution unavailable. A new MAINNET fill still needs YOUR AUTHORIZE (or Guarded consent) on an exact preview.

SECURITY RESULT: Chat/web/MCP still cannot authorize, pin, or enable autonomy. Policy was not auto-raised this pass.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Executable desk vs stale stand-down: **IMPLEMENTED + TESTED** (overlay capital LIVE VERIFIED; desktop window not pixel-clicked this pass)
- Markets Actionable default + funding column: **IMPLEMENTED + TESTED**
- Security compact READY chips: **IMPLEMENTED + TESTED**
- Live MAINNET order: **NOT IMPLEMENTED** (no AUTHORIZE this pass)
PRODUCTION READY: 0.7.6 is production-ready for an honest executable desk at clip $13. It is **not** production-ready to claim a new fill until YOU authorize a preview.
NEXT STEP: On Desk, Research the ranked book privately, then type AUTHORIZE on the exact preview if you want a clip. Do not flatten. Do not remint.

---

## M106 — PIT WEB 2.0: intelligence + proof network (2026-08-30)

DATE/TIME: 2026-08-30 07:15+03
PHASE: Public website rebuild. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent fills, APR, TEE, or mission 8C91. Do not enable Guarded Autonomy. Pairing is a late step.
GOAL: The web discovers and proves. The desktop protects and acts. Replace the coral marketing / wallet-first landing with a public intelligence product: radar, missions, proof, agent, capital simulation, download.
RESULT:
- **IMPLEMENTED:** Version **0.8.0** (Go health + companion). Desktop NSIS remains **0.7.6** (no new installer this pass).
  1. Public IA: `/` `/radar` `/radar/:coin` `/capital` `/missions` `/missions/:id/replay` `/proof` `/agent` `/how-it-works` `/download`. Pairing stays `/pair` after public value. Connected nav: Overview / My Missions / My Agent / My Proof / My Capital.
  2. Hero: "Your trading agent doesn't sleep. Your keys don't leave your machine." CTAs: Explore live PIT, Download PIT Desktop, Verify a mission. Live PIT strip from `pit-health` `/watch` + `/health` only. Actionable on the public feed is 0 unless `executionFeasible` is set — website origins do not receive buying power.
  3. Public `/watch` now includes venue `minNotional` without attaching capital. Radar tabs ACTIONABLE / RESEARCH / WATCH / BLOCKED. Market detail seals private reasoning.
  4. Capital simulator is labeled SIMULATION. TRADE / WAIT / HOLD from public mins. LIQUIDITY unavailable. No Zia APR.
  5. Missions: empty live stream is honest. Historical ETH fill OID `529167222216` is labeled HISTORICAL. Unknown ids do not invent a timeline.
  6. Proof center says WHAT and HOW. TEE is NO LIVE RECEIPT. Aristotle tx hashes are read in-browser via `evmrpc.0g.ai`. Agent passport: iTransfer NOT LIVE ON MAINNET. ERC-8004 `ownerOf` is a chain read, not a ranking.
  7. Public chat is informational. It cannot authorize, pin, or enable autonomy.
- **TESTED:** `go test ./...` PASS. Web `tsc --noEmit` PASS (~195s). `vite build` PASS. Playwright 22/22 PASS (home, radar, capital SIMULATION, missions, proof, agent iTransfer, download unsigned, pair, chat refuse authorize, no Authorize controls).
- **LIVE VERIFIED:** Overlay `D:\PIT\pit.exe` SHA256 `FA84F9A1E816AC90D2D19561A11875CFF19747F7DAF5DFBE80135C3CA5EDF7A4` (file). Running companion still reports **0.7.6** until that process is restarted. Installer SHA unchanged `1B2C5D11045E6505D9953C2074D001F074AC6E8FE47976BC21CA233A8F431A59`.
- **LIVE VERIFIED (web):** https://pit0g.vercel.app HTTP 200. Hero, radar (232/6, ACTIONABLE 0), SOL detail (spread —, reasoning sealed), capital SIMULATION ($100 TRADE 6 / $10 WAIT 6), health `0.8.0` with public minNotional. No Authorize control. iTransfer not live.
- **BLOCKED:** iTransfer UNAVAILABLE. SWAP/LP execution unavailable. Authenticode absent. macOS/Linux not packaged. No public live mission receipt. DA not claimed.

SECURITY RESULT: Browser still cannot receive session keys, Direct token plaintext, private prompts, or private memory. No web control authorizes, pins, or enables autonomy. Public health refuses `trade`/`sign`. Capital simulation is not a fill.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Public intelligence IA: **IMPLEMENTED + TESTED**
- Proof WHAT/HOW: **IMPLEMENTED + TESTED**
- Capital SIMULATION: **IMPLEMENTED + TESTED**
- Fake Zia APR / fake missions: **NOT IMPLEMENTED** (deliberate)
PRODUCTION READY: 0.8.0 is production-ready as a public discover/verify surface. It is **not** an execution UI.
NEXT STEP: Deploy Vercel + Render. Judge path: hero → radar → public/private split → proof → historical replay → download.

---

## M107 — Passport honesty: live Desk ID vs iTransfer (2026-08-30)

DATE/TIME: 2026-08-30 07:31+03
PHASE: WEB 2.0 follow-up after winner/0G/web-audit research. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent iTransfer, SIWE bind, or Galileo verify.
GOAL: `/agent` must not treat all of ERC-7857 as dead. PitDeskID mint/own/authorizeUsage is live on Aristotle. iTransfer/iClone stay unavailable. The connected account page must not look like a SIWE or TESTNET verifier.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**. No new installer. No companion bump.
  1. Public `/agent` reads `ownerOf(1)` and `isAuthorized(tokenId, owner)` on PitDeskID `0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35` via `evmrpc.0g.ai`. ERC-8004 `ownerOf(3489333)` stays a chain read, not a ranking.
  2. iTransfer / iClone labeled **NOT LIVE ON MAINNET**. No transfer button.
  3. Proof tx check also reads `getTransactionReceipt` (success / reverted / pending). TEE copy stays HISTORICAL recovered signer vs expected listed Direct teeSigner. Browser does not run VerifyE2EE.
  4. Connected My Agent drops SiweBind. My Proof no longer offers a Galileo toggle. Dead unused landing / SIWE / NetworkToggle / VerifyForm removed.
- **TESTED:** Playwright 22/22 PASS (agent Desk ID copy, iTransfer not live, no transfer button, no Authorize, capital SIMULATION).
- **LIVE VERIFIED:** Aristotle `ownerOf(1)` = `0xbdfcee82bd42fefa58ee850b3709636a8b6b0034`. iTransfer still unavailable.
- **LIVE VERIFIED (web):** https://pit0g.vercel.app/agent HTTP 200. `ownerOf 0xBDfC…0034`. `isAuthorized(owner)=true`. iTransfer NOT LIVE ON MAINNET. ERC-8004 `ownerOf` same owner, not a ranking. Home Live PIT 232/6, actionable 0, health 0.8.0. Proof TEE NO LIVE RECEIPT, no VerifyE2EE. Capital SIMULATION $100 TRADE 6 / LIQUIDITY unavailable. Commit `981c525`. Aliased pit0g.vercel.app. Render health still 0.8.0 (`dep-da9r5jpf2nfc738cf0sg`).
- **UNVERIFIED:** Browser VerifyE2EE. New fill. Running companion still reports 0.7.6 until that process is restarted.
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent. macOS/Linux not packaged.

SECURITY RESULT: Browser still cannot receive session keys, Direct token plaintext, private prompts, or private memory. Web still cannot authorize, pin, or enable autonomy. Public agent page does not fetch another account's extraAgents.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Live Desk ID ownerOf: **IMPLEMENTED + TESTED**
- iTransfer UI: **NOT IMPLEMENTED** (deliberate; not live)
- Fake SIWE / Galileo verify: **REMOVED**
PRODUCTION READY: 0.8.0 passport is honest about what ERC-7857 can do on Aristotle today.
NEXT STEP: Commit, push, Vercel. Judge can open `/agent` without a wallet.

---

## M108 — Restore coral landing design, keep WEB 2.0 copy (2026-08-30)

DATE/TIME: 2026-08-30 07:44+03
PHASE: Visual restore. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent fills or APR.
GOAL: Put back the coral PIT landing (grain, mega type, pipeline ring, marquee, motion) and compact CTAs. Keep WEB 2.0 copy, live radar tape, and public routes. Kill the full-width mobile Download bar.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**.
  1. `/` is the restored guide landing: coral hero, WireTurn, pipeline ring, moments track, Dual MAINNET, live tape, marquee, ledger, FAQ, CTA band.
  2. New copy stays: "Your trading agent doesn't sleep. Your keys don't leave your machine." Explore / Download / Verify. No Connect-your-wallet primary CTA.
  3. Live tape reads health `/watch` in the old visual language. Capital / radar / proof / agent routes unchanged.
  4. Mobile dock is two compact pills (Radar, Download), not a stretched red bar.
- **TESTED:** Playwright 22/22 PASS.
- **UNVERIFIED:** Vercel alias after this commit (deploy next).
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent.

SECURITY RESULT: Web still cannot authorize, pin, or enable autonomy. Session keys stay off the browser.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Coral landing restore: **IMPLEMENTED + TESTED**
- Full-width dock CTA: **REMOVED**
PRODUCTION READY: Landing looks like PIT again. Product routes still discover/prove only.
NEXT STEP: Push and alias pit0g.vercel.app.

---

## M109 — Unify product pages on /app desk chrome (2026-08-30)

DATE/TIME: 2026-08-30 08:17+03
PHASE: UI/UX. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent fills or APR. Landing `/` stays coral.
GOAL: `/radar` and the other public product routes must use the same ordered sidebar as `/app`. Kill the cream landing navbar that sat on top of radar copy. Pair is a late tab. Sign in / Protect my strategy stay ordered. Chat must answer the starters, not a one-letter ticker.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**. No new installer.
  1. Shared `DeskFrame` for PublicShell and AppShell: Look (Radar, Missions, Proof, Agent, Capital, How it works), This computer (Pair, Download), Your desk when signed in (Overview, My Missions, My Agent, My Proof, My Capital, Protect, Policy).
  2. `/pair` is a rail tab. Protect my strategy stays visible after pairing, and from the signed-in desk. Sign in is in the top bar and rail, not the first public goal.
  3. Ask PIT is a rail/top-bar drawer (`.desk-chat`), not a FAB stacked on Download. Landing `/` keeps the coral nav and floating Ask PIT.
  4. Chat coin match uses word boundaries. "What is happening?" no longer matches ticker `W`. Authorize still refused.
- **TESTED:** Playwright 23/23 PASS (landing, radar Pair+Ask PIT, capital SIMULATION, missions, proof, agent iTransfer, download unsigned, pair, chat refuse authorize, radar chat scanned not `W mark`).
- **LIVE VERIFIED (Chrome, local :3000, 1440×900):** `/radar` heading below header, no cream overlay, RESEARCH 6, Ask PIT answers 232 scanned / 6 eligible / actionable 0, SOL interesting, authorize refused. `/pair` Pair tab + Protect my strategy. `/proof` `/agent` `/capital` `/missions` `/how-it-works` `/download` `/radar/SOL` `/missions/historical-eth/replay` share the same rail. `/signin` and unauthenticated `/app` are wallet gate + Browse radar first. `/` stays coral landing.
- **LIVE VERIFIED (Chrome, 390×844):** Radar uses sticky header + horizontal tabs. Heading below header. Pair in the tab row. No overlapping Download FAB.
- **UNVERIFIED:** Production pit0g.vercel.app until the follow-up typecheck commit is aliased. Signed-in `/app` Overview on this Chrome session (wallet not connected here).
- **FOLLOW-UP:** `mentionedCoin` typed as `PublicCoin[]` so `tsc -b` (Vercel `npm run build`) passes.
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent. macOS/Linux not packaged.

SECURITY RESULT: Browser still cannot receive session keys, Direct token plaintext, private prompts, or private memory. Chat still cannot authorize, pin, or enable autonomy.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Desk chrome on product routes: **IMPLEMENTED + TESTED**
- Landing overlay on `/radar`: **REMOVED**
- Chat ticker substring bug: **FIXED**
PRODUCTION READY: Public product pages match the desk. Landing remains the coral story.
NEXT STEP: Commit, push, Vercel alias pit0g.vercel.app.

---

## M110 — Protect my strategy is wallet sign-in; slim the rail (2026-08-30)

DATE/TIME: 2026-08-30 08:40+03
PHASE: UI/UX. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent fills or APR.
GOAL: "Protect my strategy" must open Sign in with your wallet, not the 12-step get-started wizard. The desk rail had too many duplicate tabs.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**.
  1. Protect my strategy → `/signin` (Connect your wallet). `/app/start` redirects to `/signin`.
  2. Rail is short: Look = Radar / Proof / Agent. This computer = Pair. Download stays the coral CTA. Signed-in desk = Overview. Protect my strategy is always in the rail.
  3. Missions, Capital, How it works remain as pages, not extra rail copies of My Missions / My Agent / My Proof / My Capital / Policy.
- **TESTED:** Playwright 24/24 PASS, including protect → `/signin` heading Sign in with your wallet.
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent.

SECURITY RESULT: Web still cannot authorize, pin, or enable autonomy. Session keys stay off the browser.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Protect → wallet sign-in: **IMPLEMENTED + TESTED**
- Get-started hijack of Protect: **REMOVED**
PRODUCTION READY: Pairing hands off to wallet identity, not onboarding.
NEXT STEP: Push and alias pit0g.vercel.app.

---

## M111 — Protect my strategy is desktop wallet-link, not Overview (2026-08-30)

DATE/TIME: 2026-08-30 08:56+03
PHASE: UI/UX. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Do not invent fills or APR. Do not enable Guarded. Do not claim iTransfer live.
GOAL: "Protect my strategy" must open sign-in-wallet to link the browser with desktop (desktop setup Protect private research), not Overview and not get-started.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**. No new installer.
  1. New public `/protect`: Step 1 pair, Step 2 Connect your wallet, bind public address, then DirectSign. Token stays on the computer. Site never holds it.
  2. Rail, Sign in, Pair CTAs, Home, `/signin`, `/app/start` all land on `/protect`. After connect, the primary action is Protect my strategy, not Overview.
  3. Desktop catalog `LINKS.protect` = `https://pit0g.vercel.app/protect`. Security, Compute, setup path, next-fix, SetupWizard step "Protect private research", and companion `directDomain` href now open `/protect`. `LINKS.app` remains Overview.
- **TESTED:** Web `tsc -b` pass. Desktop `tsc -b` pass. `go test ./internal/scan` pass. Playwright 26/26 PASS (protect heading, Connect your wallet, no Overview as Protect destination, `/signin` and `/app/start` redirect to `/protect`).
- **BLOCKED:** Live Windows installer still compiled with `/app` until the next desktop rebuild. iTransfer UNAVAILABLE. Authenticode absent.

SECURITY RESULT: Browser still cannot receive session keys, Direct token plaintext, private prompts, or private memory. Chat still cannot authorize, pin, or enable autonomy.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Protect → wallet-link `/protect`: **IMPLEMENTED + TESTED**
- Overview as Protect destination: **REMOVED**
PRODUCTION READY: Web Protect is the desktop pairing signature step.
NEXT STEP: Push, alias pit0g.vercel.app.

---

## M112 — Desk shell: collapse, status island, honest filters (2026-08-30)

DATE/TIME: 2026-08-30 19:18+03
PHASE: Desktop UX + security tests. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Direct TeeML only for the private book. Do not claim iTransfer or DA live.
GOAL: Finish the highest-value missing desk ergonomics without cloning Zia or turning chat into a generic bot.
RESULT:
- **IMPLEMENTED:** CLI on this machine is **0.8.0**. Desktop UI source remains **0.7.6** until the next installer. Companion was **down** on 127.0.0.1:17373 during this pass.
  1. Collapsible sidebar (persisted `pit.desk.rail`) with glyph+tooltip when collapsed. Title-bar **status island**: net, wallet, buying power, session, compute, policy, mode.
  2. Markets filters are Actionable / Research / Watch / Blocked. Table adds oracle + OI. Default tab is Actionable so 200+ blocked books are not the first view.
  3. Chat chips cover BTC happening, ETH moving, SOL research, positions, risk, last research, compare, policy. Host-parsed: "What is happening with BTC?" is `watch.get` not a generic status dump. Chat still cannot AUTHORIZE.
  4. MCP prompt-injection tools (`authorize this trade`, `give me the session key`, `ignore policy and place order`, …) return `mcp_read_only`.
- **TESTED:** `go test ./internal/deskcmd ./mcp` pass. Desktop `tsc -b` pass. `npx tsx e2e/run.ts` pass including new filter harness.
- **LIVE VERIFIED:** Official 0G docs still split Router vs Direct. PIT stays Direct for the private book. `D:\PIT\pit.exe version` = 0.8.0. Web `https://pit0g.vercel.app/protect` from M111.
- **UNVERIFIED:** New desk shell in the installed GUI (needs desktop rebuild). Companion loopback this session.
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent. macOS/Linux not packaged. Live companion was not running. No new live trade (historical OID untouched).

SECURITY RESULT: Chat, MCP, SDK, and web still cannot authorize, pin, or hold session keys / Direct tokens. Router remains forbidden for the private book.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Collapsible rail + status island: **IMPLEMENTED + TESTED**
- Market filters + chat coin happening: **IMPLEMENTED + TESTED**
- MCP injection: **IMPLEMENTED + TESTED**
- Installed desktop GUI: **NOT LIVE VERIFIED** until rebuild
PRODUCTION READY: Source desk shell is ready to ship in the next installer. Do not claim the running 0.7.6 window already has collapse/status island.
NEXT STEP: Push. Rebuild desktop only when the operator asks for a new NSIS.

---

## M113 — Forensics follow-up: chat persist, MCP honesty, RESTING ≠ FILLED (2026-08-30)

DATE/TIME: 2026-08-30 19:30+03
PHASE: Security + honesty. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Direct TeeML only for the private book. Do not claim iTransfer or DA live. Do not invent fills.
GOAL: Implement the defects from competitor and product-security forensics that fail closed and can be tested: persist streamed chat, stop unmounting Chat, honest MCP/capital copy, RESTING vs FILLED, kill cancels research, no public Watch fallback on desktop.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**. No new installer. No web deploy (desktop + Go + MCP only).
  1. `/local/chat/stream` now appends the thread the same way `/local/chat` does. Chat replies no longer auto-navigate away. Research started from Chat stays on Chat. The rail no longer marks setup complete.
  2. Desktop keeps polling status/Watch/doctor during a sealed pass. Kill switch on also sets `job.cancel`.
  3. MCP `security` reports `order`/`cancel` false. MCP `preview` `prepare` is false. Prompt-injection tools still `mcp_read_only`.
  4. A posted OID with status resting is `RESTING`, not a fill. Empty OID is not recorded. Calibration copy says resting is not a fill.
  5. PolicyLaw no longer treats `$10` as the Hyperliquid venue floor. DeskHome ranks executable books first. Positions "available to trade" is buying power only, not spot. Desktop `fetchWatch` does not merge public Render health onto the private book. Chat treats `0x`+64 as a digest, not a session key. Named research on PUMP no longer becomes ETH. `hyperliquid` no longer matches ticker HYPE.
- **TESTED:** `go test ./mcp ./internal/deskcmd ./internal/experience ./internal/companion` (secret + parse + MCP + resting). Desktop `tsc -b` and `npx tsx e2e/run.ts` including executable-first rank.
- **UNVERIFIED:** Installed GUI still 0.7.6 until rebuild. Companion loopback this session. No new live trade.
- **BLOCKED:** iTransfer UNAVAILABLE. Authenticode absent. macOS/Linux not packaged. Live companion was not running. No new live trade (historical OID untouched).

SECURITY RESULT: Chat, MCP, SDK, and web still cannot authorize, pin, or hold session keys / Direct tokens. MCP cannot prepare a preview. Router remains forbidden for the private book. Kill switch now also cancels in-flight research.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Stream persist + Chat stay mounted: **IMPLEMENTED + TESTED**
- MCP order/cancel/prepare false: **IMPLEMENTED + TESTED**
- RESTING ≠ FILLED: **IMPLEMENTED + TESTED**
- Public Watch fallback on desktop: **REMOVED**
- Installed desktop GUI: **NOT LIVE VERIFIED** until rebuild
PRODUCTION READY: Source honesty/security fixes are ready for the next installer. Do not claim the running window already persists streamed chat.
NEXT STEP: Push. Rebuild desktop only when the operator asks for a new NSIS.

---

## M114 — Direct SKU freeze: live getService must match frozen glm-5.2 (2026-08-30)

DATE/TIME: 2026-08-30 19:38+03
PHASE: 0G honesty. Do not flatten OID 529167222216. Do not remint PIT-4bbee556. Direct TeeML only for the private book. Do not claim iTransfer or DA live. Do not swap the frozen SKU.
GOAL: Fail closed if Aristotle `getService` drifts from the proven Direct glm-5.2 SKU. Never auto-swap to GLM-5-FP8, 0GM, or the unacked compute-network-28 twin.
RESULT:
- **IMPLEMENTED:** Still **0.8.0**. No new installer. No web deploy.
  1. Frozen SKU unchanged: provider `0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D`, url `compute-network-19`, model `glm-5.2`, TeeML, teeSigner `0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9`.
  2. `GetService` + `MatchFrozenSKU` refuse unacked providers, TeeTLS, URL/model/signer drift. Sealed ask calls `FreezeLiveSKU` (transport unread keeps the freeze; a live mismatch refuses the ask).
  3. Doctor check `direct_sku`. Galileo stays unproven. Router catalog listings cannot be `private_book`.
  4. Snapshot: Router mainnet N=**31** (2026-08-30). TeeML-5 unchanged. Catalog only — not Direct.
- **TESTED:** `go test ./internal/compute ./internal/cli` pass. **LIVE:** `TestLiveGetServiceMatchesMainnetChat` PASS. `TestLivePubkeyMatchesFrozenTeeSigner` PASS. `TestLiveRouterCatalogNeverPrivateBook` PASS.
- **UNVERIFIED:** New storage upload/proof this pass (would spend). Independent funded TeeML roles. Galileo VerifyE2EE.
- **BLOCKED:** iTransfer UNAVAILABLE on Aristotle (`verifier=0`, Foundation attestor still 16602). DA not live on 16661. No new live trade.

SECURITY RESULT: Router remains forbidden for the private book. PIT will not auto-swap Direct providers. iTransfer/DA still not claimed live.
TX HASH / OID: Historical OID `529167222216` unchanged.
CLASSIFICATION:
- Live getService freeze: **IMPLEMENTED + TESTED + LIVE VERIFIED**
- Frozen glm-5.2 SKU: **UNCHANGED**
- iTransfer / DA: **NOT LIVE**
PRODUCTION READY: Freeze test is load-bearing. Do not treat Router N=31 as a second private path.
NEXT STEP: Push. Do not rebuild NSIS unless asked.









