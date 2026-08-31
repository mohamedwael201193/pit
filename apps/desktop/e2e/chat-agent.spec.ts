// MOCK TEST HARNESS — chat agent copy. Never stub VerifyE2EE or a live order.

import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { CHAT_AGENT_COPY } from "../src/AgentRun";
import { displayTurn } from "../src/CommandChat";

export function assertChatAgentCopy() {
  if (CHAT_AGENT_COPY.cannotAuthorize !== "The model cannot AUTHORIZE") throw new Error("cannot authorize copy");
  if (!CHAT_AGENT_COPY.acceptOnDesk.includes("TRADE NOW")) throw new Error("accept on desk");
  if (!CHAT_AGENT_COPY.cannotAuthorize.includes("cannot AUTHORIZE")) throw new Error("cannot authorize copy");
  const here = dirname(fileURLToPath(import.meta.url));
  const chat = readFileSync(join(here, "../src/CommandChat.tsx"), "utf8");
  if (chat.includes("authorizePreview")) throw new Error("chat must not call authorizePreview");
  if (chat.includes('aria-label="type AUTHORIZE"')) throw new Error("chat must not host the AUTHORIZE field");
  if (!chat.includes("Ask PIT to research, compare, prepare, trade, or watch")) throw new Error("composer placeholder");
  if (!chat.includes("Find best opportunity")) throw new Error("hero chip");
  if (chat.includes("cockpit-live")) throw new Error("cockpit must not sit above the transcript");
  if (!chat.includes("agent-stream")) throw new Error("one conversation stream");
  if (!chat.includes("agent-workspace")) throw new Error("one agent workspace");
  if (!chat.includes("unnamed ? \"\"")) throw new Error("next/best hunt must not reuse the last coin from the host");
  if (!chat.includes("[lines, island?.stage, island?.kind, island?.coin, island?.busy")) throw new Error("elapsed ticks must not steal scroll");
  const run = readFileSync(join(here, "../src/AgentRun.tsx"), "utf8");
  if (!run.includes("READY")) throw new Error("preview card");
  if (!run.includes("NO TRADE")) throw new Error("no-trade card");
  if (!run.includes("TRADE NOW")) throw new Error("trade now");
  if (!run.includes("onTradeNow")) throw new Error("trade now callback");
  if (run.includes("authorizePreview")) throw new Error("agent run must not import authorizePreview");
  if (!run.includes("agent-track")) throw new Error("live research track");
  if (!run.includes("agent-pipe")) throw new Error("named research stages");
  if (!run.includes("Find the next opportunity")) throw new Error("research next must skip the last book");
  if (!run.includes("0G PROOF") && !run.includes("0G TRAIL")) throw new Error("0G receipts in the turn");
  if (!run.includes("ORDER SUBMITTED")) throw new Error("execution card");
  if (!run.includes("oidBelongsToPreview")) throw new Error("stale fill must be gated");
  if (run.includes('markState === "done" ? "done"')) throw new Error("pipe must not label every row done");
  if (run.includes('className="linkish"') && run.includes("Show why")) throw new Error("duplicate show why");
  if (!chat.includes("lines.length === 0")) throw new Error("hero chips hide after the first turn");
  if (!chat.includes("visibleTurns")) throw new Error("duplicate pit lines must collapse");
  if (!chat.includes("Live numbers stay on the cards")) throw new Error("old hunt dump must collapse");
  if (!chat.includes("composeStream")) throw new Error("mission must live inside the conversation");
  if (!chat.includes("if (huntUser) return")) throw new Error("mission view must drop the old desk wall");
  if (chat.includes("Still researching ${island.coin}")) throw new Error("busy hunt must not append a second canned line");
  const css = readFileSync(join(here, "../src/styles.css"), "utf8");
  if (!css.includes("button.ghost")) throw new Error("ghost buttons must be dark-styled");
  if (!css.includes("color-scheme: dark")) throw new Error("native widgets must follow the dark desk");
  if (run.includes('["Research next", "Find the best opportunity"]')) throw new Error("research next must skip the last book");
  if (!run.includes("agent-receipts")) throw new Error("0G receipts in the turn");
  if (!run.includes("hyperliquidTrade")) throw new Error("Hyperliquid trade link");
  const collapsed = displayTurn({
    role: "pit",
    text: "AVAX is the strongest executable book among 6 of 232 live Hyperliquid perps. Mark 7.32. Venue min $10.02. Host clip $12.95. Buying power $16.18. Starting sealed 0G Direct on this computer. Chat cannot AUTHORIZE.",
    ts: 1,
  });
  if (collapsed !== "Researching AVAX. Live numbers stay on the cards.") throw new Error(collapsed);
  if (displayTurn({ role: "pit", text: "Still researching AVAX. Watch the stages above.", ts: 2 }) !== "Still researching AVAX. Watch the stages above.") {
    throw new Error("status line must stay");
  }
  const dock = readFileSync(join(here, "../src/PairingDock.tsx"), "utf8");
  if (!dock.includes('aria-label="Browser pairing"')) throw new Error("pair dock");
  const app = readFileSync(join(here, "../src/App.tsx"), "utf8");
  if (!app.includes("setupDone && !showChat")) throw new Error("pair strip hidden on Agent");
  if (!app.includes('label: "Agent"')) throw new Error("rail Agent");
  if (!app.includes("onAgentTrade")) throw new Error("desktop TRADE NOW handoff");
  if (!app.includes('authorizePreview("AUTHORIZE"')) throw new Error("TRADE NOW must use existing authorize path");
  if (!app.includes("openResearchStream")) throw new Error("research event stream");
  if (!app.includes("unnamed hunts") && !chat.includes("unnamed ? \"\"")) throw new Error("next hunt must ignore host coin");
  const companion = readFileSync(join(here, "../src/companion.ts"), "utf8");
  if (!companion.includes("/local/research/stream")) throw new Error("research SSE");
  if (!companion.includes('"/local/research/start"')) throw new Error("browser research start fallback");
}
