import { LINKS, hyperliquidAPI } from "./links";
import type { DoctorCheck, LocalStatus } from "./companion";
import { checkNamed } from "./companion";
import type { Probe } from "./readiness";

export type NextFix = {
  title: string;
  why: string;
  fix: string;
  href?: string;
  hrefLabel?: string;
  go?: "security" | "research" | "watch" | "settings";
  goLabel?: string;
};

export function nextFix(
  companionUp: boolean,
  status: LocalStatus | null,
  checks: DoctorCheck[],
  items: Probe[],
  sessionAlive: boolean,
  net: string,
): NextFix {
  if (!companionUp) {
    return {
      title: "Launch PIT on this computer",
      why: "The local companion is the execution authority. The browser cannot research or authorize.",
      fix: "Close old PIT windows and open this app again.",
    };
  }
  const wallet = checkNamed(checks, "wallet");
  if (!wallet?.ok) {
    return {
      title: "Connect your wallet",
      why: "A public 0x address binds this machine. PIT never asks for a seed phrase.",
      fix: "Pair the browser or paste the public address on Home setup.",
      href: LINKS.pair,
      hrefLabel: "Open pairing",
    };
  }
  const direct = checkNamed(checks, "direct_auth");
  if (!direct?.ok) {
    return {
      title: "Protect my strategy",
      why: "Sealed research needs a wallet signature on this computer. The website never receives the token.",
      fix: "Pair this browser, then sign Protect my strategy from the bound wallet.",
      href: LINKS.app,
      hrefLabel: "Open paired site",
    };
  }
  if (status?.kill) {
    return {
      title: "Kill switch is on",
      why: "You halted new orders on this workspace. The model cannot turn this off.",
      fix: "Open Security and resume only if you intend to trade again.",
      go: "security",
      goLabel: "Open Security",
    };
  }
  const teeOk =
    Boolean(checks.find((c) => c.name === "tee" && c.ok)) || items.find((p) => p.id === "tee")?.state === "ok";
  if (!sessionAlive) {
    return {
      title: "Create a local session",
      why: "Orders and cancels are signed on this computer. Withdraw is impossible through PIT.",
      fix: "Create the session here, then approve that agent on Hyperliquid.",
      href: hyperliquidAPI(net),
      hrefLabel: "Open Hyperliquid",
      go: "security",
      goLabel: "Open Security",
    };
  }
  const agent = checkNamed(checks, "hl_agent");
  if (agent && !agent.ok) {
    return {
      title: "Approve PIT on Hyperliquid",
      why: "extraAgents must list this session before AUTHORIZE can send an order.",
      fix: "On Hyperliquid API, approve the agent printed on Security. PIT still cannot withdraw.",
      href: hyperliquidAPI(net),
      hrefLabel: "Open Hyperliquid",
      go: "security",
      goLabel: "Open Security",
    };
  }
  const policy = checkNamed(checks, "policy");
  if (!policy?.ok) {
    return {
      title: "Pin a trading policy",
      why: "The model cannot raise clip, leverage, or permissions. Pin writes a hash on this computer.",
      fix: "Open Policy and pin the default until you change it on purpose.",
      go: "watch",
      goLabel: "Open Watch after pin",
    };
  }
  const credit = checkNamed(checks, "direct_credit");
  if (credit && !credit.ok && !teeOk) {
    return {
      title: "Fund private research",
      why: "0G Direct bills the wallet inside the sealed-path token. That is provider credit, not a Hyperliquid balance. Three sealed roles need about 3 0G locked.",
      fix: "Open pc.0g.ai Advanced funds with the same wallet.",
      href: LINKS.pcAdvanced,
      hrefLabel: "Open 0G compute",
    };
  }
  const tee = items.find((p) => p.id === "tee");
  if (tee && tee.state !== "ok") {
    return {
      title: "Run sealed research",
      why: "TEE is only proven after VerifyE2EE matches the on-chain signer on this machine.",
      fix: "Open Research and run a live sealed request. Waiting is honest until that happens.",
      go: "research",
      goLabel: "Open Research",
    };
  }
  return {
    title: "Desk is ready",
    why: "Watch is public marks. Research stays sealed. Authorize stays on this computer.",
    fix: "Pick a market, run research, then type AUTHORIZE on the exact preview.",
    go: "watch",
    goLabel: "Open Watch",
  };
}
