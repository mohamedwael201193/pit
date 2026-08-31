import type { DoctorCheck, LocalStatus } from "./companion";
import { checkNamed } from "./companion";

export type OnboardId = "pair" | "protect" | "hyperliquid" | "policy" | "ready";

export type StageLabel =
  | "NOT CONNECTED"
  | "CONNECTING"
  | "WAITING FOR APPROVAL"
  | "VERIFIED"
  | "ACTION REQUIRED"
  | "FAILED"
  | "LOCKED";

export type OnboardStep = {
  id: OnboardId;
  n: number;
  title: string;
  state: StageLabel;
  done: boolean;
  current: boolean;
  locked: boolean;
  why: string;
};

export type OnboardInput = {
  companionUp: boolean;
  paired: boolean;
  walletOk: boolean;
  protectOk: boolean;
  sessionAlive: boolean;
  hlApproved: boolean;
  policyPinned: boolean;
  kill?: boolean;
};

export function onboardInput(
  companionUp: boolean,
  status: LocalStatus | null,
  checks: DoctorCheck[],
  sessionAlive: boolean,
): OnboardInput {
  return {
    companionUp,
    paired: Boolean(status?.paired),
    walletOk: Boolean(checkNamed(checks, "wallet")?.ok || status?.wallet),
    protectOk: Boolean(checkNamed(checks, "direct_auth")?.ok),
    sessionAlive,
    hlApproved: Boolean(checkNamed(checks, "hl_agent")?.ok),
    policyPinned: Boolean(checkNamed(checks, "policy")?.ok),
    kill: Boolean(status?.kill),
  };
}

export function computeOnboard(input: OnboardInput): { steps: OnboardStep[]; current: OnboardId; ready: boolean } {
  const pairDone = Boolean(input.companionUp && input.paired);
  const protectDone = Boolean(pairDone && input.protectOk);
  const hlDone = Boolean(protectDone && input.sessionAlive && input.hlApproved);
  const policyDone = Boolean(hlDone && input.policyPinned);
  const ready = Boolean(policyDone && !input.kill);

  let current: OnboardId = "pair";
  if (pairDone && !protectDone) current = "protect";
  else if (protectDone && !hlDone) current = "hyperliquid";
  else if (hlDone && !policyDone) current = "policy";
  else if (ready) current = "ready";
  else if (pairDone && protectDone && hlDone && policyDone) current = "ready";

  const pairState: StageLabel = !input.companionUp
    ? "NOT CONNECTED"
    : pairDone
      ? "VERIFIED"
      : "ACTION REQUIRED";
  const protectState: StageLabel = !pairDone
    ? "LOCKED"
    : protectDone
      ? "VERIFIED"
      : "ACTION REQUIRED";
  const hlState: StageLabel = !protectDone
    ? "LOCKED"
    : !input.sessionAlive
      ? "ACTION REQUIRED"
      : input.hlApproved
        ? "VERIFIED"
        : "WAITING FOR APPROVAL";
  const policyState: StageLabel = !hlDone ? "LOCKED" : policyDone ? "VERIFIED" : "ACTION REQUIRED";
  const readyState: StageLabel = input.kill ? "FAILED" : ready ? "VERIFIED" : "LOCKED";

  const steps: OnboardStep[] = [
    {
      id: "pair",
      n: 1,
      title: "Pair this browser",
      state: pairState,
      done: pairDone,
      current: current === "pair",
      locked: false,
      why: "A one-time code opens a read-only channel. Session keys never leave this computer.",
    },
    {
      id: "protect",
      n: 2,
      title: "Protect my strategy",
      state: protectState,
      done: protectDone,
      current: current === "protect",
      locked: !pairDone,
      why: "The bound wallet signs on this computer. The website never receives the Direct token.",
    },
    {
      id: "hyperliquid",
      n: 3,
      title: "Connect Hyperliquid",
      state: hlState,
      done: hlDone,
      current: current === "hyperliquid",
      locked: !protectDone,
      why: "PIT creates the agent on this computer. You approve order and cancel only. Withdraw stays impossible.",
    },
    {
      id: "policy",
      n: 4,
      title: "Pin policy",
      state: policyState,
      done: policyDone,
      current: current === "policy",
      locked: !hlDone,
      why: "You set clip and assets. Re-pin anytime. The model cannot raise them.",
    },
    {
      id: "ready",
      n: 5,
      title: "Ready to trade",
      state: readyState,
      done: ready,
      current: current === "ready",
      locked: !policyDone,
      why: "Research, preview, and AUTHORIZE stay on this computer.",
    },
  ];

  return { steps, current, ready };
}
