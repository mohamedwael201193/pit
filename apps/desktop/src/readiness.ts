import type { DoctorCheck, LocalStatus } from "./companion";
import { checkNamed } from "./companion";

export type ProbeState = "ok" | "waiting" | "fail";

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

export function probes(checks: DoctorCheck[], status: LocalStatus | null, companionUp: boolean): Probe[] {
  const wallet = fromCheck("wallet", "Wallet", checkNamed(checks, "wallet"), "Connect your wallet, then bind this machine.");
  const hl = fromCheck("hyperliquid", "Hyperliquid", checkNamed(checks, "hyperliquid"), "Public book not reached yet.");
  const sealer = checkNamed(checks, "direct_sealer");
  const auth = checkNamed(checks, "direct_auth");
  let direct: Probe = {
    id: "direct",
    label: "0G Direct",
    state: "waiting",
    detail: "Direct TeeML is not armed until the sealer binary and PIT_DIRECT_AUTH_FILE are present.",
  };
  if (sealer?.ok && auth?.ok) {
    direct = { id: "direct", label: "0G Direct", state: "ok", detail: "Sealer and auth file present. Research still verifies each response." };
  } else if (sealer && !sealer.ok) {
    direct = { id: "direct", label: "0G Direct", state: "waiting", detail: sealer.detail };
  } else if (auth && !auth.ok) {
    direct = { id: "direct", label: "0G Direct", state: "waiting", detail: auth.detail };
  }
  const tee: Probe = {
    id: "tee",
    label: "TEE verification",
    state: "waiting",
    detail: "No sealed research has been verified on this machine in this session.",
  };
  const storage = fromCheck("storage", "Storage", checkNamed(checks, "storage"), "Official Go storage client not found.");
  if (storage.state === "ok") {
    storage.detail = "Official client present. A proof is verified only after an upload/download.";
  }
  const policy = fromCheck("policy", "Policy", checkNamed(checks, "policy"), "Policy is not pinned for this workspace.");
  const session: Probe = status?.sessionAlive
    ? { id: "session", label: "Session", state: "ok", detail: "Live order/cancel session on this machine." }
    : fromCheck("session", "Session", checkNamed(checks, "session"), "Approve an order/cancel agent on this machine.");
  const local: Probe = companionUp
    ? { id: "local", label: "Local companion", state: "ok", detail: status?.version ? `PIT ${status.version}` : "Loopback companion is up." }
    : { id: "local", label: "Local companion", state: "waiting", detail: "Waiting for 127.0.0.1:17373." };
  return [local, wallet, hl, direct, tee, storage, policy, session];
}
