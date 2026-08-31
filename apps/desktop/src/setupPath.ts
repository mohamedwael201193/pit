import type { DoctorCheck, LocalStatus } from "./companion";
import { checkNamed } from "./companion";
import { LINKS, hyperliquidAPI } from "./links";

type Tone = "READY" | "ACTION" | "BLOCKED";

export type PathStep = {
  id: string;
  title: string;
  why: string;
  tone: Tone;
  href?: string;
  hrefLabel?: string;
  go?: "security" | "research" | "watch" | "settings" | "policy" | "account";
  goLabel?: string;
};

export function setupPath(
  companionUp: boolean,
  status: LocalStatus | null,
  checks: DoctorCheck[],
  sessionAlive: boolean,
  net: string,
): PathStep[] {
  const paired = Boolean(status?.paired);
  const wallet = checkNamed(checks, "wallet");
  const direct = checkNamed(checks, "direct_auth");
  const policy = checkNamed(checks, "policy");
  const agent = checkNamed(checks, "hl_agent");
  return [
    {
      id: "pair",
      title: "Pair this browser",
      why: "A one-time code opens a read-only channel. Session keys never leave this computer.",
      tone: paired ? "READY" : companionUp ? "ACTION" : "BLOCKED",
      href: LINKS.pair,
      hrefLabel: "Open pairing",
      go: "security",
      goLabel: "Show code",
    },
    {
      id: "protect",
      title: "Protect my strategy",
      why: "One wallet signature stores the sealed-path token on this computer. The website never receives it.",
      tone: direct?.ok ? "READY" : paired ? "ACTION" : "BLOCKED",
      href: LINKS.protect,
      hrefLabel: "Protect my strategy",
    },
    {
      id: "hyperliquid",
      title: "Connect Hyperliquid",
      why: "PIT creates the agent on this computer. You approve order and cancel only with the master wallet.",
      tone: agent?.ok && sessionAlive ? "READY" : direct?.ok ? "ACTION" : "BLOCKED",
      href: hyperliquidAPI(net),
      hrefLabel: "Approve PIT on Hyperliquid",
      go: "security",
      goLabel: "Open Security",
    },
    {
      id: "policy",
      title: "Pin policy",
      why: "Clip, assets, and kill live on this computer. The model cannot raise them.",
      tone: policy?.ok ? "READY" : agent?.ok ? "ACTION" : "BLOCKED",
      go: "policy",
      goLabel: "Open Policy",
    },
    {
      id: "ready",
      title: "Ready to trade",
      why: "Research, preview, and AUTHORIZE stay on this computer. Chat cannot authorize.",
      tone: paired && wallet?.ok && direct?.ok && sessionAlive && agent?.ok && policy?.ok && !status?.kill ? "READY" : "BLOCKED",
      go: "research",
      goLabel: "Open Research",
    },
  ];
}
