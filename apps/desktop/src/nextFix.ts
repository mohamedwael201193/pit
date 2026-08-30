import { LINKS, hyperliquidAPI, hyperliquidApp } from "./links";
import type { DoctorCheck, LocalStatus } from "./companion";
import { checkNamed } from "./companion";
import type { Probe } from "./readiness";

export type NextFix = {
  title: string;
  why: string;
  fix: string;
  href?: string;
  hrefLabel?: string;
  go?: "security" | "research" | "markets" | "automation";
  goLabel?: string;
};

export function nextFix(
  companionUp: boolean,
  status: LocalStatus | null,
  checks: DoctorCheck[],
  items: Probe[],
  sessionAlive: boolean,
  net: string,
  capital?: { buyingPower?: number; execGate?: string; execWhy?: string },
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
  if (!sessionAlive) {
    return {
      title: "Create a local session",
      why: "Orders and cancels are signed on this computer. Withdraw is impossible through PIT.",
      fix: "Create the session here, then approve that agent on Hyperliquid.",
      href: hyperliquidAPI(net),
      hrefLabel: "Open Hyperliquid API",
      go: "security",
      goLabel: "Open Security",
    };
  }
  const agent = checkNamed(checks, "hl_agent");
  if (agent && !agent.ok) {
    return {
      title: "Approve PIT on Hyperliquid",
      why: "Hyperliquid must list this session before AUTHORIZE can send an order.",
      fix: "Open Hyperliquid API. Authorize API Wallet using the PIT agent name and address on Security. PIT still cannot withdraw.",
      href: hyperliquidAPI(net),
      hrefLabel: "Open Hyperliquid API",
      go: "security",
      goLabel: "Open Security",
    };
  }
  const policy = checkNamed(checks, "policy");
  if (!policy?.ok) {
    return {
      title: "Pin a trading policy",
      why: "The model cannot raise clip, leverage, or permissions. Pin writes a hash on this computer.",
      fix: "Open Security and pin the default until you change it on purpose.",
      go: "security",
      goLabel: "Open Security",
    };
  }
  const credit = checkNamed(checks, "direct_credit");
  if (credit && !credit.ok) {
    const unread = credit.detail.toLowerCase().includes("unread");
    const sponsored = credit.detail.toLowerCase().includes("sponsor");
    if (!unread && !sponsored) {
      return {
        title: "Fund private research",
        why: "Sealed research is billed as 0G provider credit on this wallet. That is not Hyperliquid buying power and not PIT infrastructure.",
        fix: "Open 0G Direct funds in Advanced mode with the same wallet. PIT only asks because this ledger is short.",
        href: LINKS.pcAdvanced,
        hrefLabel: "Open 0G Direct funds",
      };
    }
    return {
      title: "Check private compute again",
      why: unread
        ? "The provider ledger was unread. That is not a zero balance and not a reason to visit a generic dashboard."
        : "PIT can sponsor sealed research within the daily workspace cap.",
      fix: "Use Check again on Security. Protect stays on this computer.",
      go: "security",
      goLabel: "Check again",
    };
  }
  const tee = items.find((p) => p.id === "tee");
  if (tee && tee.state === "fail") {
    return {
      title: "TEE verification failed",
      why: "VerifyE2EE did not match the on-chain signer. Idle is not a failure. A failed verify is.",
      fix: "Open Research and retry sealed research. PIT did not place an order.",
      go: "research",
      goLabel: "Open Research",
    };
  }
  const power = capital?.buyingPower;
  const gate = String(capital?.execGate || "");
  if (gate === "insufficient_margin" || (typeof power === "number" && power > 0 && power < 10)) {
    return {
      title: "Watching. Nothing can open.",
      why: `This account has ${typeof power === "number" ? `$${power.toFixed(2)}` : "less than the venue floor"}. PIT will not invent size. 0G compute credit is a different pile of money.`,
      fix: capital?.execWhy || "Available venue margin is below this market's Hyperliquid minimum.",
      href: hyperliquidApp(net),
      hrefLabel: "Open Hyperliquid",
      go: "markets",
      goLabel: "Open Markets",
    };
  }
  return {
    title: "Desk is ready",
    why: "Markets is public marks. Research stays sealed. Authorize stays on this computer.",
    fix: "Pick a market, run research, then type AUTHORIZE on the exact preview — or enable Guarded Autonomy on Automation.",
    go: "markets",
    goLabel: "Open Markets",
  };
}
