// MOCK TEST HARNESS — chat agent copy. Never stub VerifyE2EE or a live order.

import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { CHAT_AGENT_COPY } from "../src/AgentRun";

export function assertChatAgentCopy() {
  if (CHAT_AGENT_COPY.cannotAuthorize !== "The model cannot AUTHORIZE") throw new Error("cannot authorize copy");
  if (!CHAT_AGENT_COPY.acceptOnDesk.includes("TRADE NOW")) throw new Error("accept on desk");
  if (!CHAT_AGENT_COPY.cannotAuthorize.includes("cannot AUTHORIZE")) throw new Error("cannot authorize copy");
  const here = dirname(fileURLToPath(import.meta.url));
  const chat = readFileSync(join(here, "../src/CommandChat.tsx"), "utf8");
  if (chat.includes("authorizePreview")) throw new Error("chat must not call authorizePreview");
  if (chat.includes('aria-label="type AUTHORIZE"')) throw new Error("chat must not host the AUTHORIZE field");
  if (!chat.includes("Ask PIT what to trade")) throw new Error("composer placeholder");
  if (!chat.includes("Find best opportunity")) throw new Error("hero chip");
  if (!chat.includes("cockpit-live")) throw new Error("live cockpit must sit outside the transcript");
  const run = readFileSync(join(here, "../src/AgentRun.tsx"), "utf8");
  if (!run.includes("READY TO TRADE")) throw new Error("preview card");
  if (!run.includes("NO TRADE")) throw new Error("no-trade card");
  if (!run.includes("TRADE NOW")) throw new Error("trade now");
  if (!run.includes("onTradeNow")) throw new Error("trade now callback");
  if (run.includes("authorizePreview")) throw new Error("agent run must not import authorizePreview");
  if (!run.includes("PRIVATE 0G RESEARCH")) throw new Error("live pipe");
  if (!run.includes("ORDER SUBMITTED")) throw new Error("execution card");
  const dock = readFileSync(join(here, "../src/PairingDock.tsx"), "utf8");
  if (!dock.includes('aria-label="Browser pairing"')) throw new Error("pair dock");
  const app = readFileSync(join(here, "../src/App.tsx"), "utf8");
  if (!app.includes("setupDone && !showChat")) throw new Error("pair strip hidden on Agent");
  if (!app.includes('label: "Agent"')) throw new Error("rail Agent");
  if (!app.includes("onAgentTrade")) throw new Error("desktop TRADE NOW handoff");
  if (!app.includes('authorizePreview("AUTHORIZE"')) throw new Error("TRADE NOW must use existing authorize path");
  if (!app.includes("openResearchStream")) throw new Error("research event stream");
  const companion = readFileSync(join(here, "../src/companion.ts"), "utf8");
  if (!companion.includes("/local/research/stream")) throw new Error("research SSE");
}
