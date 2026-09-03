import { HISTORICAL_FILL, PIT_AGENT, VERIFIED_FILL } from "./facts";
import { compact, fundingLabel, markLabel, usd } from "./format";
import type { PublicCoin, WatchView } from "./types";
import { coinMin } from "./venue";

export type ChatTurn = { q: string; a: string };

export const STARTERS = [
  "What is happening?",
  "Why is SOL interesting?",
  "Why didn't PIT trade?",
  "How does PIT protect my strategy?",
  "What happened during the last mission?",
  "Show me proof.",
  "How does PIT use 0G?",
] as const;

function mentionedCoin(q: string, coins: PublicCoin[]) {
  const ordered = [...coins].sort((a, b) => b.coin.length - a.coin.length);
  return ordered.find((c) => {
    const token = c.coin.toLowerCase().replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    if (token.length < 2) return false;
    return new RegExp(`(?:^|[^a-z0-9])${token}(?:$|[^a-z0-9])`).test(q);
  });
}

export function answerChat(raw: string, watch: WatchView | null): string {
  const q = raw.trim().toLowerCase();
  if (!q) return "Ask what is happening, why a market is watched, or how to verify. This chat cannot authorize, pin policy, or enable autonomy.";

  if (/(authori[sz]e|enable autonomy|arm sleep|sleep mission|pin policy|sign this|session key|seed phrase|raise max trade|execute btc now)/.test(q)) {
    return "This chat is informational. It cannot authorize, pin policy, arm a Sleep Mission, or hold keys. Open PIT Desktop for those actions.";
  }

  if (q.includes("0g") || q.includes("tee") || q.includes("storage") || q.includes("compute")) {
    return "0G Compute (Direct TeeML) is the private research path — the website never receives the sealed prompt. 0G Storage holds durable public-safe evidence when a receipt is published. 0G Chain is where a judge can read a transaction from the public RPC. Agentic ID is mint, ownership, and authorizeUsage on Aristotle. ERC-8004 identity and reputation are chain facts — this site does not invent a score.";
  }

  if (q.includes("proof") || q.includes("verify") || q.includes("oid")) {
    return `Open /proof. Latest recorded desk fill: ${VERIFIED_FILL.market} size ${VERIFIED_FILL.sz} OID ${VERIFIED_FILL.oid}, job ${VERIFIED_FILL.job}. Research 0G ${VERIFIED_FILL.researchTx}. Order 0G ${VERIFIED_FILL.orderTx}. An older ETH fill OID ${HISTORICAL_FILL.oid} is HISTORICAL and was not flattened. TEE VerifyE2EE is not run in this browser, so this page will not badge Verified. Paste a 0G Chain transaction hash to read it from Aristotle RPC.`;
  }

  if (q.includes("last mission") || (q.includes("happened") && q.includes("mission"))) {
    return `There is no live public mission stream. Latest recorded fill: ${VERIFIED_FILL.market} ${VERIFIED_FILL.sz} @ ${VERIFIED_FILL.px}, OID ${VERIFIED_FILL.oid}. Older ETH OID ${HISTORICAL_FILL.oid} is HISTORICAL. Open /missions/${VERIFIED_FILL.id}/replay.`;
  }

  if (q.includes("protect") || q.includes("strategy") || q.includes("private") || q.includes("sealed")) {
    return "Private reasoning is sealed with 0G Direct. This browser never receives the private book, session key, Direct token plaintext, or memory. Policy, AUTHORIZE, and Sleep Missions stay on PIT Desktop. 0G lets PIT prove the machine without exposing the intelligence.";
  }

  if (q.includes("didn't trade") || q.includes("did not trade") || q.includes("no trade") || q.includes("stand-down") || q.includes("stood down")) {
    return "A public health feed does not include your account's buying power, so this site cannot call a book actionable. Policy-eligible marks are not orders. If desktop stood down, that is a successful stand-down, not a crash — the host refused size, risk, or policy. This chat will not invent a private thesis.";
  }

  if (q.includes("desktop") || q.includes("download") || q.includes("install")) {
    return "PIT Desktop is the private brain: policy, keys, session, authorization, autonomy, execution. This website discovers and proves. Use Download PIT Desktop on this site — the installer is served as a file, not a GitHub Releases page. Verify SHA256. It is not Authenticode-signed. macOS and Linux are not claimed until packaged and tested.";
  }

  const coins = watch?.coins ?? [];
  const sol = coins.find((c) => c.coin?.toUpperCase() === "SOL");

  if (q.includes("happening") || q.includes("now") || q.includes("radar") || (q.includes("market") && !mentionedCoin(q, coins))) {
    if (!watch) {
      return "Public health is not loaded yet. If it stays empty, the watch process is down — this site will not invent scanned counts.";
    }
    return `Live public watch: ${watch.scanned ?? 0} Hyperliquid perps scanned, ${watch.count ?? 0} policy-eligible. Account-actionable count on this public feed is ${coins.filter((c) => c.executionFeasible).length} — buying power is not attached for website origins. Best public-eligible: ${watch.best?.coin ?? "none"}. ${watch.bestWhy ?? ""} Agent ${PIT_AGENT.name} executes only from desktop.`;
  }

  const asked = mentionedCoin(q, coins);

  if (/\bsol\b/.test(q) || (asked && asked.coin === "SOL") || (sol && q.includes("interesting") && !asked)) {
    if (!sol) {
      return "SOL is not in the current public watch payload. PIT will not invent a mark.";
    }
    const min = coinMin(sol);
    return `SOL is on the live Hyperliquid book this health process scanned. Mark $${markLabel(sol.mark)}${sol.oracle ? `, oracle $${markLabel(sol.oracle)}` : ""}. Funding ${fundingLabel(sol.funding)}. Open interest ${compact(sol.openInterest)}. Venue min about ${usd(min)}. ${sol.eligible ? "It is policy-eligible for research on the public default policy." : "It is not policy-eligible on the public default policy."} Private thesis is sealed. This is not a trade recommendation.`;
  }

  if (asked) {
    const min = coinMin(asked);
    return `${asked.coin} mark $${markLabel(asked.mark)}, funding ${fundingLabel(asked.funding)}, OI ${compact(asked.openInterest)}, min ${usd(min)}. ${asked.eligible ? "Policy-eligible on the public default list." : "Not on the public eligible list."} No private thesis. Open /radar/${asked.coin} for the public market sheet.`;
  }

  return "This chat answers public intelligence only. Try: what is happening, why SOL is watched, how 0G is used, or show me proof. For private command, open PIT Desktop.";
}
