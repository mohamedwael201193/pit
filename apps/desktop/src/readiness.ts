import type { DoctorCheck, LocalStatus } from "./companion";
import { checkNamed } from "./companion";

export type ProbeState = "ok" | "waiting" | "fail" | "optional";

export type Probe = {
  id: string;
  label: string;
  state: ProbeState;
  detail: string;
};

function fromCheck(id: string, label: string, c: DoctorCheck | undefined, waitingDetail: string): Probe {
  if (!c) return { id, label, state: "waiting", detail: waitingDetail };
  if (c.ok) return { id, label, state: "ok", detail: c.detail };
  return { id, label, state: "waiting", detail: c.detail || waitingDetail };
}

export function probes(checks: DoctorCheck[], status: LocalStatus | null, companionUp: boolean, teeVerified = false): Probe[] {
  const wallet = fromCheck("wallet", "Wallet", checkNamed(checks, "wallet"), "Connect your wallet, then bind this machine.");
  const hl = fromCheck("hyperliquid", "Hyperliquid", checkNamed(checks, "hyperliquid"), "Public book not reached yet.");
  const sealer = checkNamed(checks, "direct_sealer");
  const auth = checkNamed(checks, "direct_auth");
  let direct: Probe = {
    id: "direct",
    label: "0G Direct",
    state: "waiting",
    detail: "Direct TeeML is not armed until you sign Protect my strategy on the paired browser.",
  };
  if (sealer?.ok && auth?.ok) {
    direct = { id: "direct", label: "0G Direct", state: "ok", detail: "Wallet-signed Direct token on this computer. Research still verifies each response." };
  } else if (sealer && !sealer.ok) {
    direct = { id: "direct", label: "0G Direct", state: "waiting", detail: sealer.detail };
  } else if (auth && !auth.ok) {
    direct = { id: "direct", label: "0G Direct", state: "waiting", detail: auth.detail };
  }
  const credit = checkNamed(checks, "direct_credit");
  const compute = fromCheck("direct_credit", "0G compute", credit, "Private compute credit not confirmed.");
  const teeCheck = checkNamed(checks, "tee");
  const tee: Probe =
    teeVerified || teeCheck?.ok
      ? {
          id: "tee",
          label: "TEE verification",
          state: "ok",
          detail: teeCheck?.detail || "VerifyE2EE matched the on-chain teeSigner on this machine.",
        }
      : {
          id: "tee",
          label: "TEE verification",
          state: "waiting",
          detail: teeCheck?.detail || "No sealed research has been verified on this machine yet.",
        };
  const storage = fromCheck("storage", "Storage", checkNamed(checks, "storage"), "Official Go storage client not found.");
  if (storage.state === "ok") {
    storage.detail = "Official client present. A proof is verified only after an upload/download.";
  }
  const policy = fromCheck("policy", "Policy", checkNamed(checks, "policy"), "Policy is not pinned for this workspace.");
  const agent = fromCheck(
    "hl_agent",
    "Hyperliquid agent",
    checkNamed(checks, "hl_agent"),
    "Create a local session, then approve it on Hyperliquid. PIT cannot withdraw.",
  );
  const session: Probe = status?.sessionAlive
    ? { id: "session", label: "Session", state: "ok", detail: "Live order/cancel session on this machine." }
    : fromCheck("session", "Session", checkNamed(checks, "session"), "Approve an order/cancel agent on this machine.");
  const local: Probe = companionUp
    ? { id: "local", label: "Local companion", state: "ok", detail: status?.version ? `PIT ${status.version}` : "Loopback companion is up." }
    : { id: "local", label: "Local companion", state: "waiting", detail: "Waiting for 127.0.0.1:17373." };
  const chain = fromCheck("0g_rpc", "0G Chain", checkNamed(checks, "0g_rpc"), "Aristotle RPC not reached.");
  const kill: Probe = status?.kill
    ? { id: "kill", label: "Kill switch", state: "fail", detail: "Kill switch is on. Discovery continues. AUTHORIZE stays off." }
    : { id: "kill", label: "Kill switch", state: "ok", detail: "Kill switch is off. Host law still sizes and gates every order." };
  const pairing: Probe = companionUp
    ? {
        id: "pairing",
        label: "Pairing",
        state: status?.paired ? "ok" : "waiting",
        detail: status?.paired
          ? `Browser paired (${status.pairingDevices || 1} device). Desktop is the signer.`
          : "Type the one-time code on the pairing page. The website never receives a session key.",
      }
    : { id: "pairing", label: "Pairing", state: "waiting", detail: "Launch PIT Desktop first." };
  const identity: Probe = {
    id: "identity",
    label: "Identity",
    state: "optional",
    detail: "Transfer of Agentic ID is not live on mainnet. Trading does not wait on it.",
  };
  const execution: Probe =
    policy.state === "ok" && session.state === "ok" && agent.state === "ok" && kill.state === "ok"
      ? { id: "execution", label: "Execution readiness", state: "ok", detail: "Pinned policy, live session, and Hyperliquid agent. AUTHORIZE still required." }
      : {
          id: "execution",
          label: "Execution readiness",
          state: "waiting",
          detail: "Pin policy, approve the PIT agent, and keep a live session. Chat cannot skip this.",
        };
  return [local, wallet, pairing, direct, compute, chain, storage, tee, hl, agent, session, policy, execution, kill, identity];
}
