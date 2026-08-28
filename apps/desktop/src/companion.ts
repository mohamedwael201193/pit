export const COMPANION = "http://127.0.0.1:17373";

export type LocalStatus = {
  sessionAlive?: boolean;
  agent?: string;
  agentName?: string;
  sessionExpires?: number;
  hypothesis?: string;
  network?: string;
  wallet?: string;
  workspace?: string;
  kill?: boolean;
  sign?: boolean;
  trade?: boolean;
  version?: string;
  lastOrder?: { oid?: string; cloid?: string; hash?: string; market?: string; side?: string; sz?: number; status?: string; posted?: boolean; cancelled?: boolean; venue?: string };
};

export type PairCode = {
  code?: string;
  expires?: string;
  sign?: boolean;
  trade?: boolean;
};

export type DoctorCheck = {
  name: string;
  ok: boolean;
  detail: string;
};

export type BindResult = {
  ok?: boolean;
  wallet?: string;
  workspace?: string;
  network?: string;
  agent?: string;
  error?: string;
  sign?: boolean;
  trade?: boolean;
  provider?: string;
  model?: string;
  expiresAt?: number;
  source?: string;
  digest?: string;
  message?: string;
  explain?: string;
  note?: string;
  roles?: Array<{
    role?: string;
    verify_e2ee?: string;
    pubkey_signer?: string;
    teeSigner?: string;
    proposed_side?: string;
    survives?: boolean;
    kill?: boolean;
  }>;
  verify?: boolean;
  stage?: string;
  elapsed_ms?: number;
  running?: boolean;
  done?: boolean;
  job_id?: string;
  workspace_id?: string;
  coin?: string;
  transient?: boolean;
  poll?: string;
  seq?: number;
  heartbeat_unix_ms?: number;
  terminal?: boolean;
  terminal_kind?: string;
  card_title?: string;
  retryable?: boolean;
  current_stage?: string;
  evidence?: unknown;
  preview?: {
    eligible?: boolean;
    deny?: string;
    market?: string;
    side?: string;
    sz?: number;
    limitPx?: string;
    hash?: string;
    cloid?: string;
    expiryUnixMs?: number;
    notionalUsd?: number;
    reasons?: string[];
    kind?: string;
    note?: string;
  };
  preview_hash?: string;
  deny?: string;
  eligible?: boolean;
  oid?: string;
  posted?: boolean;
  kind?: string;
  market?: string;
  side?: string;
  sz?: number;
  limitPx?: string;
  hash?: string;
  cloid?: string;
};

type TauriWindow = Window & {
  __TAURI_INTERNALS__?: { invoke: (cmd: string, args?: unknown) => Promise<unknown> };
  __TAURI__?: { core?: { invoke: (cmd: string, args?: unknown) => Promise<unknown> } };
};

export function nativeInvoke<T = unknown>(cmd: string, args?: unknown): Promise<T> | null {
  const w = window as TauriWindow;
  const inv = w.__TAURI_INTERNALS__?.invoke ?? w.__TAURI__?.core?.invoke;
  if (typeof inv !== "function") return null;
  return (args === undefined ? inv(cmd) : inv(cmd, args)) as Promise<T>;
}

async function nativeJson<T extends object>(cmd: string, args?: unknown): Promise<T | null> {
  const p = nativeInvoke<T>(cmd, args);
  if (!p) return null;
  try {
    return await p;
  } catch {
    return null;
  }
}

async function nativeJsonOrError<T extends object>(cmd: string, args?: unknown): Promise<T> {
  const p = nativeInvoke<T>(cmd, args);
  if (!p) {
    throw new Error("companion_down");
  }
  return p;
}

async function fetchJson<T extends object>(path: string, init?: RequestInit): Promise<T | null> {
  try {
    const r = await fetch(`${COMPANION}${path}`, init);
    if (!r.ok) return null;
    return (await r.json()) as T;
  } catch {
    return null;
  }
}

function rejectSecrets<T extends { sign?: boolean; trade?: boolean }>(body: T | null): T | null {
  if (!body || body.sign || body.trade) return null;
  return body;
}

export async function wakeCompanion(): Promise<boolean> {
  const p = nativeInvoke<boolean>("ensure_companion");
  if (p) {
    try {
      return Boolean(await p);
    } catch {
      return false;
    }
  }
  const health = await fetchJson<{ ok?: boolean }>("/health");
  return Boolean(health);
}

export async function localStatus(): Promise<LocalStatus | null> {
  const native = rejectSecrets(await nativeJson<LocalStatus>("local_status"));
  if (native) return native;
  return rejectSecrets(await fetchJson<LocalStatus>("/local/status"));
}

export async function pairCode(): Promise<PairCode | null> {
  const native = rejectSecrets(await nativeJson<PairCode>("local_code"));
  if (native) return native;
  return rejectSecrets(await fetchJson<PairCode>("/local/code"));
}

export async function doctor(): Promise<DoctorCheck[]> {
  const native = await nativeJson<{ checks?: DoctorCheck[]; sign?: boolean }>("local_doctor");
  if (native && !native.sign) {
    return Array.isArray(native.checks) ? native.checks : [];
  }
  const fetched = await fetchJson<{ checks?: DoctorCheck[]; sign?: boolean }>("/local/doctor");
  if (!fetched || fetched.sign) return [];
  return Array.isArray(fetched.checks) ? fetched.checks : [];
}

export async function bindWallet(wallet: string, network: string): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_init", { wallet, network });
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function createLocalSession(): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_session");
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function pinLocalPolicy(): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_policy");
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function fetchWatch(network: string): Promise<{ coins?: Array<{ coin: string; reason: string; why?: string; trend?: string; rank?: number; freshness?: string; mark: number; eligible?: boolean; oracle?: number; funding?: number; openInterest?: number; timestamp?: string }>; sign?: boolean; trade?: boolean }> {
  const native = rejectSecrets(await nativeJson<{ coins?: Array<{ coin: string; reason: string; why?: string; trend?: string; rank?: number; freshness?: string; mark: number; eligible?: boolean; oracle?: number; funding?: number; openInterest?: number; timestamp?: string }>; sign?: boolean; trade?: boolean }>("local_watch", { network }));
  if (native) return native;
  const fetched = await fetchJson<{ coins?: Array<{ coin: string; reason: string; mark: number; eligible?: boolean; oracle?: number; funding?: number; openInterest?: number }>; sign?: boolean; trade?: boolean }>(`/watch?network=${network}`);
  if (fetched) return fetched;
  try {
    const r = await fetch(`https://pit-health.onrender.com/watch?network=${network}`);
    if (!r.ok) return {};
    return (await r.json()) as { coins?: Array<{ coin: string; reason: string; mark: number; eligible?: boolean; oracle?: number; funding?: number; openInterest?: number }>; sign?: boolean; trade?: boolean };
  } catch {
    return {};
  }
}

export async function setKillSwitch(on: boolean): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_kill", { on });
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function revokeLocalSession(): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_revoke_session");
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function directIntent(): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_direct_intent");
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function directStatus(): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_direct_status");
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function startResearch(coin: string, hypothesis?: string): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_research_start", { coin, hypothesis: hypothesis || "none" });
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function researchStatus(): Promise<BindResult> {
  const native = rejectSecrets(await nativeJson<BindResult>("local_research_status"));
  if (native) return native;
  const fetched = rejectSecrets(await fetchJson<BindResult>("/local/research/status"));
  if (fetched) return fetched;
  return { error: "POLL_FAILED", poll: "POLL_FAILED", transient: true, running: true };
}

export async function researchEvidence(): Promise<BindResult> {
  const native = rejectSecrets(await nativeJson<BindResult>("local_research_result"));
  if (native) return native;
  const fetched = rejectSecrets(await fetchJson<BindResult>("/local/research/result"));
  if (fetched) return fetched;
  return { error: "POLL_FAILED", transient: true };
}

export async function cancelResearch(): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_research_cancel");
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function connectionPreview(coin?: string, reduceOnly?: boolean): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_connection_preview", {
      coin: coin || "ETH",
      reduceOnly: Boolean(reduceOnly),
    });
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function authorizePreview(typed: string, hash: string): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_authorize", { typed, hash });
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function cancelBoundOrder(typed: string): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_cancel_order", { typed });
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export type ChatMessage = { ts?: number; role?: string; text?: string; tool?: string; thread?: string; coin?: string };
export type ChatThread = { id: string; title: string; created?: number; updated?: number; preview?: string };
export type DirectModel = {
  model?: string;
  label?: string;
  path?: string;
  provider?: string;
  verifiability?: string;
  proven_e2ee?: boolean;
  private_book?: boolean;
  capability?: string;
  note?: string;
  latency?: string;
  cost?: string;
};
export type ChatReply = BindResult & {
  reply?: string;
  tool?: string;
  execute?: boolean;
  start_research?: boolean;
  coin?: string;
  navigate?: string;
  open_url?: string;
  thread?: string;
};
export type ActivityEvent = {
  ts?: number;
  kind?: string;
  market?: string;
  action?: string;
  status?: string;
  job_id?: string;
  preview_hash?: string;
  oid?: string;
  reason?: string;
};
export type SecurityDomain = {
  id?: string;
  state?: string;
  why?: string;
  means?: string;
  do?: string;
  href?: string;
  hrefLabel?: string;
};
export type VenuePosition = {
  coin?: string;
  sz?: string;
  entryPx?: string;
  markPx?: number;
  unrealizedPnl?: string;
  leverage?: string;
  marginUsed?: string;
  account?: string;
  policyClipUsd?: number;
};

async function localGet<T extends object>(cmd: string, path: string): Promise<T | null> {
  const native = rejectSecrets(await nativeJson<T>(cmd));
  if (native) return native;
  return rejectSecrets(await fetchJson<T>(path));
}

export async function sendDeskCommand(text: string, thread = "desk", signal?: AbortSignal): Promise<ChatReply> {
  try {
    const native = await nativeJsonOrError<ChatReply>("local_chat", { text, thread });
    if (signal?.aborted) throw new DOMException("Aborted", "AbortError");
    if (native.sign || native.trade) return { error: "companion_denied", execute: false };
    return { ...native, execute: false };
  } catch (e) {
    if (e instanceof DOMException && e.name === "AbortError") throw e;
    const body = await fetchJson<ChatReply>("/local/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text, thread }),
      signal,
    });
    if (!body || body.sign || body.trade) {
      return { error: "companion_http", execute: false, reply: "PIT could not reach the local companion." };
    }
    return { ...body, execute: false };
  }
}

export async function fetchChatLog(thread = "desk"): Promise<ChatMessage[]> {
  const qs = `/local/chat/log?thread=${encodeURIComponent(thread)}`;
  const native = rejectSecrets(await nativeJson<{ messages?: ChatMessage[]; sign?: boolean }>("local_chat_log", { thread }));
  if (native && Array.isArray(native.messages)) return native.messages;
  const body = rejectSecrets(await fetchJson<{ messages?: ChatMessage[]; sign?: boolean }>(qs));
  return Array.isArray(body?.messages) ? body.messages : [];
}

export async function fetchChatThreads(): Promise<ChatThread[]> {
  const body = await localGet<{ threads?: ChatThread[]; sign?: boolean }>("local_chat_threads", "/local/chat/threads");
  return Array.isArray(body?.threads) ? body.threads : [];
}

export async function mutateChatThread(op: "new" | "rename" | "delete", id?: string, title?: string): Promise<{ threads?: ChatThread[]; id?: string; error?: string }> {
  try {
    const native = await nativeJsonOrError<{ threads?: ChatThread[]; id?: string; error?: string; sign?: boolean }>(
      "local_chat_thread",
      { op, id, title },
    );
    if (native.sign) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg };
  }
}

export async function fetchModels(): Promise<DirectModel[]> {
  const cat = await fetchModelCatalog();
  return cat.models;
}

export async function fetchModelCatalog(): Promise<{
  models: DirectModel[];
  private_verified: DirectModel[];
  other_chat: DirectModel[];
  unsupported: DirectModel[];
}> {
  const body = await localGet<{
    models?: DirectModel[];
    groups?: { private_verified?: DirectModel[]; other_chat?: DirectModel[]; unsupported?: DirectModel[] };
    sign?: boolean;
  }>("local_models", "/local/models");
  const models = Array.isArray(body?.models) ? body.models : [];
  return {
    models,
    private_verified: body?.groups?.private_verified || models.filter((m) => m.private_book),
    other_chat: body?.groups?.other_chat || [],
    unsupported: body?.groups?.unsupported || [],
  };
}

export type AutoPrefs = {
  watch?: boolean;
  auto_research?: boolean;
  notify?: boolean;
  cadence_minutes?: number;
  trigger?: string;
  markets?: string[];
  execute?: boolean;
};

export async function fetchAutomation(): Promise<AutoPrefs> {
  const body = await localGet<{ prefs?: AutoPrefs; sign?: boolean }>("local_automation", "/local/automation");
  return body?.prefs || { watch: true, notify: true, auto_research: false, cadence_minutes: 15, trigger: "policy_pass" };
}

export async function saveAutomation(prefs: AutoPrefs): Promise<AutoPrefs> {
  const body = await fetchJson<{ prefs?: AutoPrefs; error?: string; execute?: boolean }>("/local/automation", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ ...prefs, execute: false }),
  });
  if (body?.execute) return prefs;
  return body?.prefs || prefs;
}

export async function fetchActivity(): Promise<ActivityEvent[]> {
  const body = await localGet<{ events?: ActivityEvent[]; sign?: boolean }>("local_activity", "/local/activity");
  return Array.isArray(body?.events) ? body.events : [];
}

export type AccountSummary = {
  accountValue?: string;
  totalMarginUsed?: string;
  totalNtlPos?: string;
  withdrawable?: string;
};

export async function fetchPositions(): Promise<{
  account?: string;
  positions: VenuePosition[];
  error?: string;
  lastOrder?: LocalStatus["lastOrder"];
  summary?: AccountSummary;
}> {
  const body = await localGet<{
    account?: string;
    positions?: VenuePosition[];
    error?: string;
    lastOrder?: LocalStatus["lastOrder"];
    summary?: AccountSummary;
    sign?: boolean;
  }>("local_positions", "/local/positions");
  return {
    account: body?.account,
    positions: Array.isArray(body?.positions) ? body.positions : [],
    error: body?.error,
    lastOrder: body?.lastOrder,
    summary: body?.summary,
  };
}

export async function fetchSecurity(): Promise<SecurityDomain[]> {
  const body = await localGet<{ domains?: SecurityDomain[]; sign?: boolean }>("local_security", "/local/security");
  return Array.isArray(body?.domains) ? body.domains : [];
}

export async function fetchCalibration(): Promise<{ n?: number; need?: number; copy?: string; enough?: boolean }> {
  return (
    (await localGet<{ n?: number; need?: number; copy?: string; enough?: boolean; sign?: boolean }>(
      "local_calibration",
      "/local/calibration",
    )) || { copy: "NOT ENOUGH DATA" }
  );
}

export async function fetchIdentity(): Promise<{ itransfer?: string; iclone?: string; note?: string }> {
  return (
    (await localGet<{ itransfer?: string; iclone?: string; note?: string; sign?: boolean }>("local_identity", "/local/identity")) || {
      itransfer: "UNAVAILABLE",
      iclone: "UNAVAILABLE",
    }
  );
}

export async function fetchUpdate(): Promise<{
  version?: string;
  research_running?: boolean;
  restart_allowed?: boolean;
  authenticode?: boolean;
  note?: string;
}> {
  return (
    (await localGet<{
      version?: string;
      research_running?: boolean;
      restart_allowed?: boolean;
      authenticode?: boolean;
      note?: string;
      sign?: boolean;
    }>("local_update", "/local/update")) || { authenticode: false, restart_allowed: true }
  );
}

export async function forgetMemory(): Promise<BindResult> {
  try {
    const native = await nativeJsonOrError<BindResult>("local_memory_forget");
    if (native.sign || native.trade) return { error: "companion_denied" };
    return native;
  } catch (e) {
    const msg = e instanceof Error ? e.message : "companion_http";
    return { error: msg || "companion_http" };
  }
}

export async function fetchExplain(): Promise<BindResult> {
  return (await localGet<BindResult>("local_explain", "/local/explain")) || {};
}

export async function runResearch(coin: string): Promise<BindResult> {
  return startResearch(coin);
}

export function prettyCode(code: string) {
  const raw = code.replace(/[^A-Za-z0-9]/g, "").toUpperCase();
  if (raw.length !== 8) return raw;
  return `${raw.slice(0, 4)}-${raw.slice(4)}`;
}

export function checkNamed(checks: DoctorCheck[], name: string): DoctorCheck | undefined {
  return checks.find((c) => c.name === name);
}

export function describeBindError(code: string) {
  if (code === "workspace_owned") {
    return "This computer already belongs to another wallet. Revoke and forget on this machine first.";
  }
  if (code === "network_switch_denied") {
    return "This workspace is already bound to the other network. Forget the workspace to switch.";
  }
  if (code === "wallet_required") {
    return "Paste your public 0x address. PIT never asks for a seed or a private key.";
  }
  if (code === "unbound") {
    return "Bind your public wallet on this computer first.";
  }
  if (code === "direct_token_required" || code === "direct_challenge_required") {
    return "Pair the browser, then sign Protect my strategy. PIT never asks you to edit an env file.";
  }
  if (code === "direct_token_expired") {
    return "The sealed-path signature expired. Sign again in the paired browser.";
  }
  if (code === "galileo_e2ee_unproven") {
    return "TESTNET sealed research is off until VerifyE2EE is proven. Switch to MAINNET.";
  }
  if (code === "direct_ledger") {
    return "The sealed-path signature is on this computer. Fund Direct at pc.0g.ai Advanced with the same wallet. PIT does not ask for a private key.";
  }
  if (code === "direct_provider_http") {
    return "The sealed provider rejected the request. PIT did not fall back to Router.";
  }
  if (code === "companion_down") {
    return "The local companion is not running yet.";
  }
  return code;
}
