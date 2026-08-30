// MOCK TEST HARNESS — chat agent copy. Never stub VerifyE2EE or a live order.

import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { CHAT_AGENT_COPY } from "../src/AgentRun";

export function assertChatAgentCopy() {
  if (CHAT_AGENT_COPY.cannotAuthorize !== "Chat cannot AUTHORIZE") throw new Error("cannot authorize copy");
  if (!CHAT_AGENT_COPY.acceptOnDesk.includes("AUTHORIZE")) throw new Error("accept on desk");
  const here = dirname(fileURLToPath(import.meta.url));
  const chat = readFileSync(join(here, "../src/CommandChat.tsx"), "utf8");
  if (chat.includes("authorizePreview")) throw new Error("chat must not call authorizePreview");
  if (chat.includes('aria-label="type AUTHORIZE"')) throw new Error("chat must not host the AUTHORIZE field");
  const dock = readFileSync(join(here, "../src/PairingDock.tsx"), "utf8");
  if (!dock.includes('aria-label="Browser pairing"')) throw new Error("pair dock");
  const app = readFileSync(join(here, "../src/App.tsx"), "utf8");
  if (!app.includes("compact")) throw new Error("pair strip on the desk");
}
