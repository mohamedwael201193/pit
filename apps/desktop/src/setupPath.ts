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
  const wallet = checkNamed(checks, "wallet");
  const direct = checkNamed(checks, "direct_auth");
  const policy = checkNamed(checks, "policy");
  const agent = checkNamed(checks, "hl_agent");
  const tee = checkNamed(checks, "tee");
  return [
    {
      id: "wallet",
      title: "Connect wallet",
      why: "A public 0x address binds this machine. PIT never asks for a seed phrase.",
      tone: wallet?.ok ? "READY" : companionUp ? "ACTION" : "BLOCKED",
      href: LINKS.pair,
      hrefLabel: "Open pairing",
    },
    {
      id: "pair",
      title: "Pair this computer",
      why: "The browser watches. This desktop is the execution authority.",
      tone: companionUp ? "READY" : "ACTION",
      href: LINKS.pair,
      hrefLabel: "Open pairing",
    },
    {
      id: "protect",
      title: "Protect private research",
      why: "One wallet signature stores the sealed-path token on this computer. The website never receives it.",
      tone: direct?.ok ? "READY" : wallet?.ok ? "ACTION" : "BLOCKED",
      href: LINKS.app,
      hrefLabel: "Open paired site",
    },
    {
      id: "hl",
      title: "Connect Hyperliquid",
      why: "Open the official API page. PIT needs your account, a scoped agent, and a local session.",
      tone: agent?.ok && sessionAlive ? "READY" : wallet?.ok ? "ACTION" : "BLOCKED",
      href: hyperliquidAPI(net),
      hrefLabel: "Open Hyperliquid API",
      go: "security",
      goLabel: "Open Security",
    },
    {
      id: "approve",
      title: "Approve PIT agent",
      why: "Authorize API Wallet with the printed name and address. extraAgents must list that address. PIT cannot withdraw.",
      tone: agent?.ok ? "READY" : sessionAlive ? "ACTION" : "BLOCKED",
      href: hyperliquidAPI(net),
      hrefLabel: "Open Hyperliquid API",
      go: "security",
      goLabel: "Open Security",
    },
    {
      id: "session",
      title: "Create secure session",
      why: "Order and cancel only, on this computer. If extraAgents still lists this agent, PIT reuses it and does not mint a new address.",
      tone: sessionAlive ? "READY" : wallet?.ok ? "ACTION" : "BLOCKED",
      go: "security",
      goLabel: "Create session",
    },
    {
      id: "policy",
      title: "Set your policy",
      why: "Clip, assets, and kill live on this computer. The model cannot raise them.",
      tone: policy?.ok ? "READY" : wallet?.ok ? "ACTION" : "BLOCKED",
      go: "policy",
      goLabel: "Open Policy",
    },
    {
      id: "research",
      title: "Run your first research",
      why: "Sealed Direct TeeML. Stand-down is a verified result. Authorize stays on this computer.",
      tone: status?.kill ? "BLOCKED" : tee?.ok ? "READY" : direct?.ok ? "ACTION" : "BLOCKED",
      go: "research",
      goLabel: "Open Research",
    },
  ];
}
