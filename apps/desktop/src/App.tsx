import { FormEvent, useEffect, useMemo, useRef, useState } from "react";
import { AuthorizeGate } from "./AuthorizeGate";
import { CapabilityMatrix } from "./CapabilityMatrix";
import { EmptyHome } from "./EmptyHome";
import { KillNote } from "./KillNote";
import { NAMED } from "./namedStates";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PermissionsCard } from "./Permissions";
import { PolicyLaw } from "./PolicyLaw";
import { PreviewNote } from "./PreviewNote";
import { SessionNote } from "./SessionNote";
import { HyperliquidCard } from "./HyperliquidCard";
import { committeeDeny, explainCommittee } from "./committee";
import {
  authorizePreview,
  bindWallet,
  createLocalSession,
  cancelBoundOrder,
  cancelResearch,
  connectionPreview,
  describeBindError,
  doctor,
  fetchActivity,
  fetchCalibration,
  fetchIdentity,
  fetchPositions,
  fetchSecurity,
  fetchUpdate,
  fetchWatch,
  forgetMemory,
  localStatus,
  pairCode,
  pinLocalPolicy,
  prettyCode,
  researchEvidence,
  researchStatus,
  revokeLocalSession,
  setKillSwitch,
  startResearch,
  wakeCompanion,
  type ActivityEvent,
  type BindResult,
  type DoctorCheck,
  type LocalStatus,
  type SecurityDomain,
  type VenuePosition,
} from "./companion";
import { explainStop, explainStopHref } from "./explain";
import { LINKS, hyperliquidAPI, hyperliquidApp } from "./links";
import { nextFix } from "./nextFix";
import { probes, type Probe } from "./readiness";
import { setupPath } from "./setupPath";
import { WelcomePath } from "./WelcomePath";
import { committeeVerified, oidBelongsToPreview, researchCardTitle } from "./honesty";
import { CommandChat } from "./CommandChat";
import { EvidenceDrawer } from "./EvidenceDrawer";
import { ActivityTimeline } from "./ActivityTimeline";
import { PositionsPanel } from "./PositionsPanel";
import { SecurityCenter } from "./SecurityCenter";
import { ComputeCard } from "./ComputeCard";
import { SetupWizard } from "./SetupWizard";

type Net = "mainnet" | "testnet";
type View = "home" | "watch" | "research" | "activity" | "policy" | "security" | "account" | "settings";

type Coin = { coin: string; reason: string; mark: number; eligible?: boolean; oracle?: number; funding?: number; openInterest?: number };

const RESEARCH_STAGES = [
  "READING_MARKET",
  "SEALING_PRIVATE_BOOK",
  "CONTACTING_PRIVATE_PROVIDER",
  "RECEIVING_SEALED_RESPONSE",
  "VERIFYING_TEE_SIGNATURE",
  "RESEARCHER",
  "CHALLENGER",
  "RISK",
  "DETERMINISTIC_ENGINE",
  "POLICY",
  "PREVIEW",
  "READY",
] as const;

type ResearchRole = {
  role?: string;
  verify_e2ee?: string;
  pubkey_signer?: string;
  proposed_side?: string;
  survives?: boolean;
  kill?: boolean;
};

function sleep(ms: number) {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

function roleVerified(roles: ResearchRole[], name: string) {
  return roles.some(
    (r) => String(r.role || "").toLowerCase() === name && String(r.verify_e2ee || "").toUpperCase() === "OK",
  );
}

function canonicalResearchStage(stage: string, roles: ResearchRole[]) {
  const s = (stage || "").toUpperCase();
  if (s === "RISK_START" || s === "RISK_HTTP_REQUEST" || s === "RISK_HTTP_RESPONSE" || s === "RISK_E2EE_VERIFY") {
    return "RISK";
  }
  if (s.endsWith("_VERIFIED")) {
    const base = s.slice(0, -"_VERIFIED".length);
    if (base === "RESEARCHER") return "CHALLENGER";
    if (base === "CHALLENGER") return "RISK";
    if (base === "RISK") return "DETERMINISTIC_ENGINE";
  }
  if (s.endsWith("_FAILED")) return s.slice(0, -"_FAILED".length);
  if (s === "CONTACTING_PRIVATE_PROVIDER" || s === "RECEIVING_SEALED_RESPONSE" || s === "VERIFYING_TEE_SIGNATURE") {
    if (roleVerified(roles, "challenger")) return "RISK";
    if (roleVerified(roles, "researcher")) return "CHALLENGER";
    return s;
  }
  if ((RESEARCH_STAGES as readonly string[]).includes(s)) return s;
  return s;
}

function stageMark(name: string, stage: string, roles: ResearchRole[]) {
  const s = canonicalResearchStage(stage, roles);
  if (name === "RESEARCHER") {
    if (roleVerified(roles, "researcher") || ["CHALLENGER", "RISK", "DETERMINISTIC_ENGINE", "POLICY", "PREVIEW", "READY"].includes(s)) {
      return "done";
    }
    if (s === "RESEARCHER") return "lit";
    return "";
  }
  if (name === "CHALLENGER") {
    if (roleVerified(roles, "challenger") || ["RISK", "DETERMINISTIC_ENGINE", "POLICY", "PREVIEW", "READY"].includes(s)) {
      return "done";
    }
    if (s === "CHALLENGER") return "lit";
    return "";
  }
  if (name === "RISK") {
    if (roleVerified(roles, "risk") || ["DETERMINISTIC_ENGINE", "POLICY", "PREVIEW", "READY"].includes(s)) return "done";
    if (s === "RISK") return "lit";
    return "";
  }
  const current = RESEARCH_STAGES.indexOf(s as (typeof RESEARCH_STAGES)[number]);
  const i = RESEARCH_STAGES.indexOf(name as (typeof RESEARCH_STAGES)[number]);
  if (current < 0 || i < 0) return "";
  if (i < current) return "done";
  if (i === current) return "lit";
  return "";
}

function ResearchProgress({
  stage,
  elapsedMs,
  coin,
  roles,
  pollMiss,
  onCancel,
}: {
  stage: string;
  elapsedMs: number;
  coin: string;
  roles: ResearchRole[];
  pollMiss?: boolean;
  onCancel: () => void;
}) {
  const shown = canonicalResearchStage(stage, roles);
  const riskLive = shown === "RISK" && !roleVerified(roles, "risk");
  return (
    <article className="card" role="status">
      <p className="label">LIVE SEALED REQUEST</p>
      <h2>{shown.replaceAll("_", " ")}</h2>
      <p>
        {coin || "ETH"} · {(elapsedMs / 1000).toFixed(1)}s elapsed. This is a live Direct round-trip, not a timer.
      </p>
      {riskLive ? (
        <p>
          Risk is running · {(elapsedMs / 1000).toFixed(1)}s. The provider is still working. PIT is not spinning a fake
          timer.
        </p>
      ) : null}
      {pollMiss ? (
        <p role="status">Connection check missed — research is still running.</p>
      ) : null}
      <ol className="pipe stages">
        {RESEARCH_STAGES.map((name) => {
          const mark = stageMark(name, stage, roles);
          const prefix = mark === "done" ? "✓ " : mark === "lit" ? "● " : "○ ";
          return (
            <li key={name} className={mark}>
              {prefix}
              {name.replaceAll("_", " ")}
            </li>
          );
        })}
      </ol>
      <button type="button" className="linkish" onClick={onCancel}>
        Cancel
      </button>
    </article>
  );
}

const SETUP_KEY = "pit.desk.setup";

const RAIL: { id: View; label: string }[] = [
  { id: "home", label: "Desk" },
  { id: "watch", label: "Watch" },
  { id: "research", label: "Research" },
  { id: "activity", label: "Activity" },
  { id: "policy", label: "Policy" },
  { id: "security", label: "Security" },
  { id: "account", label: "Account" },
  { id: "settings", label: "Settings" },
];

function markProbe(p: Probe) {
  if (p.state === "ok") return "pass";
  if (p.state === "fail") return "fail";
  return "wait";
}

function PairingBlock({
  code,
  expires,
  companionUp,
}: {
  code: string;
  expires: string;
  companionUp: boolean;
}) {
  const display = code ? prettyCode(code) : companionUp ? "rotating…" : "waiting for local PIT";
  const [copied, setCopied] = useState(false);
  async function copy() {
    if (!code) return;
    try {
      await navigator.clipboard.writeText(prettyCode(code));
      setCopied(true);
      window.setTimeout(() => setCopied(false), 1600);
    } catch {
      setCopied(false);
    }
  }
  return (
    <article className="card pair-card">
      <p className="label">PAIR THIS COMPUTER</p>
      <p className="pair-code" aria-label="pairing code">
        {display}
      </p>
      <p className="fine">
        Type this code at {LINKS.pair}. It expires in two minutes and works once. The website never receives a session
        key.
      </p>
      {expires ? <p className="fine">Expires {expires}</p> : null}
      <button type="button" className="linkish" onClick={() => void copy()} disabled={!code}>
        {copied ? "Copied" : "Copy code"}
      </button>
    </article>
  );
}

function ProbeList({
  items,
  net,
  onGo,
}: {
  items: Probe[];
  net: string;
  onGo?: (view: View) => void;
}) {
  return (
    <ul className="probes">
      {items.map((p) => (
        <li key={p.id} className={markProbe(p)}>
          <strong>{p.state === "ok" ? "ready" : p.state === "fail" ? "fail" : "waiting"}</strong>
          <span>{p.label}</span>
          <em>{p.detail}</em>
          {p.state !== "ok" ? <ProbeAction id={p.id} net={net} onGo={onGo} /> : null}
        </li>
      ))}
    </ul>
  );
}

function ProbeAction({ id, net, onGo }: { id: string; net: string; onGo?: (view: View) => void }) {
  if (id === "hyperliquid" || id === "hl_agent" || id === "session") {
    return (
      <span className="cta-row">
        {id === "session" ? (
          <button type="button" className="linkish" onClick={() => onGo?.("security")}>
            Create session
          </button>
        ) : null}
        <a className="linkish" href={hyperliquidApp(net)} target="_blank" rel="noreferrer">
          Open Hyperliquid
        </a>
        <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
          Open Hyperliquid API
        </a>
      </span>
    );
  }
  if (id === "direct") {
    return (
      <a className="linkish" href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
        Open 0G Private Compute
      </a>
    );
  }
  if (id === "tee") {
    return (
      <button type="button" className="linkish" onClick={() => onGo?.("research")}>
        Open Research
      </button>
    );
  }
  if (id === "wallet" || id === "local") {
    return (
      <a className="linkish" href={LINKS.pair} target="_blank" rel="noreferrer">
        Open pairing
      </a>
    );
  }
  return null;
}

export function App() {
  const [view, setView] = useState<View>("home");
  const [net, setNet] = useState<Net>("mainnet");
  const [sessionAlive, setSessionAlive] = useState(false);
  const [agent, setAgent] = useState("");
  const [code, setCode] = useState("");
  const [expires, setExpires] = useState("");
  const [companionUp, setCompanionUp] = useState(false);
  const [status, setStatus] = useState<LocalStatus | null>(null);
  const [checks, setChecks] = useState<DoctorCheck[]>([]);
  const [coins, setCoins] = useState<Coin[]>([]);
  const [researchStop, setResearchStop] = useState<string | null>(null);
  const [researchNote, setResearchNote] = useState<string | null>(null);
  const [researchRoles, setResearchRoles] = useState<ResearchRole[]>([]);
  const [researchBusy, setResearchBusy] = useState(false);
  const [researchStage, setResearchStage] = useState("READING_MARKET");
  const [researchElapsed, setResearchElapsed] = useState(0);
  const [researchCoin, setResearchCoin] = useState("ETH");
  const [researchEvidenceText, setResearchEvidence] = useState<string>("");
  const [pollMiss, setPollMiss] = useState(false);
  const [researchJobId, setResearchJobId] = useState("");
  const [researchKind, setResearchKind] = useState("");
  const [preview, setPreview] = useState<NonNullable<BindResult["preview"]> | null>(null);
  const [previewHash, setPreviewHash] = useState("");
  const [authTyped, setAuthTyped] = useState("");
  const [authBusy, setAuthBusy] = useState(false);
  const [authErr, setAuthErr] = useState<string | null>(null);
  const [lastOid, setLastOid] = useState("");
  const [activity, setActivity] = useState<ActivityEvent[]>([]);
  const [positions, setPositions] = useState<VenuePosition[]>([]);
  const [positionAccount, setPositionAccount] = useState("");
  const [positionErr, setPositionErr] = useState("");
  const [securityDomains, setSecurityDomains] = useState<SecurityDomain[]>([]);
  const [calibCopy, setCalibCopy] = useState("NOT ENOUGH DATA");
  const [identityNote, setIdentityNote] = useState("Transfer of Agentic ID is not live on mainnet.");
  const [updateNote, setUpdateNote] = useState("This build is checksum-verified, not OS-signed.");
  const [restartAllowed, setRestartAllowed] = useState(true);
  const [hypothesis, setHypothesis] = useState<"none" | "long" | "short">("none");
  const researchGen = useRef(0);
  const researchBusyRef = useRef(false);
  researchBusyRef.current = researchBusy;
  const [techOpen, setTechOpen] = useState(false);
  const [ticks, setTicks] = useState(0);
  const [setupStep, setSetupStep] = useState(0);
  const [walletDraft, setWalletDraft] = useState("");
  const [bindError, setBindError] = useState<string | null>(null);
  const [bindBusy, setBindBusy] = useState(false);
  const [pinned, setPinned] = useState(false);
  const [setupDone, setSetupDone] = useState(() => {
    try {
      return window.localStorage.getItem(SETUP_KEY) === "1";
    } catch {
      return false;
    }
  });

  useEffect(() => {
    let gone = false;
    void (async () => {
      await wakeCompanion();
      if (gone) return;
      const st = await researchStatus();
      if (gone || st.transient) return;
      if (st.job_id) setResearchJobId(st.job_id);
      if (st.stage) setResearchStage(st.stage);
      if (typeof st.elapsed_ms === "number") setResearchElapsed(st.elapsed_ms);
      if (st.coin) setResearchCoin(String(st.coin).toUpperCase());
      const roles = Array.isArray(st.roles) ? st.roles : [];
      if (roles.length) setResearchRoles(roles);
      if (st.terminal_kind) setResearchKind(String(st.terminal_kind));
      if (st.preview) {
        setPreview(st.preview);
        setPreviewHash(st.preview_hash || st.preview.hash || "");
      }
      const verified = committeeVerified(roles, st.verify);
      const deny = String(st.deny || st.preview?.deny || "");
      if ((verified || st.verify || committeeDeny(deny)) && !st.running) {
        setResearchStop(null);
        setResearchNote(st.note || "Sealed committee verified on this computer.");
        return;
      }
      if (!st.running && st.error && !verified && !committeeDeny(deny)) {
        setResearchStop(st.error);
        return;
      }
      if (!st.running || !st.job_id) return;
      const gen = ++researchGen.current;
      setResearchBusy(true);
      setResearchStop(null);
      setView("research");
      const wall = Date.now() - (Number(st.elapsed_ms) || 0);
      const tick = window.setInterval(() => {
        if (gen === researchGen.current) setResearchElapsed(Date.now() - wall);
      }, 250);
      try {
        for (;;) {
          await sleep(1000);
          if (gone || gen !== researchGen.current) return;
          const next = await researchStatus();
          if (next.transient) {
            setPollMiss(true);
            continue;
          }
          setPollMiss(false);
          if (next.job_id) setResearchJobId(next.job_id);
          if (next.stage) setResearchStage(next.stage);
          if (typeof next.elapsed_ms === "number") setResearchElapsed(next.elapsed_ms);
          const nextRoles = Array.isArray(next.roles) ? next.roles : [];
          if (nextRoles.length) setResearchRoles(nextRoles);
          if (next.preview) {
            setPreview(next.preview);
            setPreviewHash(next.preview_hash || next.preview.hash || "");
          }
          if (next.terminal_kind) setResearchKind(String(next.terminal_kind));
          const ok = committeeVerified(nextRoles, next.verify);
          const nextDeny = String(next.deny || next.preview?.deny || "");
          if ((ok || next.verify || committeeDeny(nextDeny)) && !next.running) {
            setResearchStop(null);
            setResearchNote(next.note || "Sealed committee verified on this computer.");
            return;
          }
          if (next.running) continue;
          if (next.error && !ok && !committeeDeny(nextDeny)) {
            setResearchStop(next.error);
            return;
          }
          setResearchStop(ok ? null : next.error || "COMMITTEE_INCOMPLETE");
          if (ok) setResearchNote(next.note || "Sealed committee verified on this computer.");
          return;
        }
      } finally {
        window.clearInterval(tick);
        if (gen === researchGen.current) {
          setResearchBusy(false);
          setPollMiss(false);
        }
      }
    })();
    return () => {
      gone = true;
    };
  }, []);

  useEffect(() => {
    if (!techOpen) return;
    let gone = false;
    void researchEvidence().then((st) => {
      if (gone || !st.evidence) return;
      setResearchEvidence(JSON.stringify(st.evidence, null, 2));
    });
    return () => {
      gone = true;
    };
  }, [techOpen]);

  useEffect(() => {
    let gone = false;
    let timer = 0;
    let statusBusy = false;
    let codeBusy = false;
    let doctorBusy = false;
    let watchBusy = false;
    let lastDoctor = 0;
    let lastWatch = 0;
    const tick = () => {
      if (gone) return;
      setTicks((n) => n + 1);
      const now = Date.now();
      if (!statusBusy && !researchBusyRef.current) {
        statusBusy = true;
        localStatus()
          .then((s) => {
            if (gone) return;
            setCompanionUp(Boolean(s));
            setStatus(s);
            setSessionAlive(Boolean(s?.sessionAlive));
            setAgent(s?.agent || "");
            if (s?.hypothesis === "long" || s?.hypothesis === "short" || s?.hypothesis === "none") {
              setHypothesis(s.hypothesis);
            }
            if (s?.lastOrder?.oid) setLastOid(String(s.lastOrder.oid));
            if (s?.network === "testnet" || s?.network === "mainnet") setNet(s.network);
            if (s?.wallet) setWalletDraft((cur) => cur || s.wallet || "");
          })
          .catch(() => {
            if (!gone && !researchBusyRef.current) {
              setCompanionUp(false);
              setStatus(null);
            }
          })
          .finally(() => {
            statusBusy = false;
          });
      }
      if (!researchBusyRef.current && !codeBusy) {
        codeBusy = true;
        pairCode()
          .then((p) => {
            if (gone) return;
            setCode(p?.code || "");
            setExpires(p?.expires || "");
          })
          .catch(() => undefined)
          .finally(() => {
            codeBusy = false;
          });
      }
      if (!researchBusyRef.current && !doctorBusy && now - lastDoctor >= 15000) {
        doctorBusy = true;
        lastDoctor = now;
        doctor()
          .then((c) => {
            if (!gone) {
              setChecks(c);
              setPinned(Boolean(c.find((x) => x.name === "policy" && x.ok)));
            }
          })
          .catch(() => undefined)
          .finally(() => {
            doctorBusy = false;
          });
      }
      if (!researchBusyRef.current && !watchBusy && now - lastWatch >= 8000) {
        watchBusy = true;
        lastWatch = now;
        fetchWatch(net)
          .then((body) => {
            if (gone || body.sign || body.trade) return;
            setCoins(Array.isArray(body.coins) ? body.coins : []);
          })
          .catch(() => {
            if (!gone) setCoins([]);
          })
          .finally(() => {
            watchBusy = false;
          });
        void fetchActivity().then((ev) => {
          if (!gone) setActivity(ev);
        });
        void fetchPositions().then((p) => {
          if (gone) return;
          setPositions(p.positions);
          setPositionAccount(p.account || "");
          setPositionErr(p.error || "");
        });
        void fetchSecurity().then((d) => {
          if (!gone) setSecurityDomains(d);
        });
        void fetchCalibration().then((c) => {
          if (!gone) setCalibCopy(c.copy || "NOT ENOUGH DATA");
        });
        void fetchIdentity().then((id) => {
          if (!gone) setIdentityNote(id.note || "Transfer of Agentic ID is not live on mainnet.");
        });
        void fetchUpdate().then((u) => {
          if (gone) return;
          setRestartAllowed(u.restart_allowed !== false);
          setUpdateNote(u.note || "This build is checksum-verified, not OS-signed.");
        });
      }
    };
    const loop = () => {
      tick();
      if (!gone) timer = window.setTimeout(loop, 2000);
    };
    loop();
    return () => {
      gone = true;
      window.clearTimeout(timer);
    };
  }, [net]);

  useEffect(() => {
    let gone = false;
    researchStatus()
      .then((st) => {
        if (gone || st.sign || st.trade) return;
        const roles = Array.isArray(st.roles) ? st.roles : [];
        const verified = committeeVerified(roles, st.verify);
        if (st.stage) setResearchStage(st.stage);
        if (typeof st.elapsed_ms === "number") setResearchElapsed(st.elapsed_ms);
        if (roles.length) setResearchRoles(roles);
      if (st.terminal_kind) setResearchKind(String(st.terminal_kind));
        if (st.coin) setResearchCoin(st.coin);
        if (st.preview) {
          setPreview(st.preview);
          setPreviewHash(st.preview_hash || st.preview.hash || "");
        }
        if (verified && !st.running) {
          setResearchStop(null);
          setResearchNote(st.note || "Sealed committee verified on this computer.");
          return;
        }
        if (committeeDeny(String(st.deny || st.preview?.deny || "")) && !st.running) {
          setResearchStop(null);
          setResearchNote(st.note || "Sealed committee verified on this computer.");
          return;
        }
        if (st.error && !st.running && !verified) setResearchStop(st.error);
      })
      .catch(() => undefined);
    return () => {
      gone = true;
    };
  }, []);

  useEffect(() => {
    if (!sessionAlive) return;
    try {
      window.localStorage.setItem(SETUP_KEY, "1");
    } catch {
      /* ignore */
    }
    setSetupDone(true);
  }, [sessionAlive]);

  const denyCode = String(preview?.deny || "");
  const committee = committeeDeny(denyCode) || committeeDeny(researchStop);
  const explained = committee ? null : explainStop(researchStop);
  const committeeCopy = committeeDeny(denyCode)
    ? explainCommittee(denyCode)
    : committeeDeny(researchStop)
      ? explainCommittee(researchStop || "")
      : null;
  const companionStuck = !companionUp && ticks >= 5;
  const items = useMemo(
    () =>
      probes(
        checks,
        status,
        companionUp,
        researchRoles.some((r) => String(r.verify_e2ee).toUpperCase() === "OK") ||
          Boolean(checks.find((c) => c.name === "tee" && c.ok)),
      ),
    [checks, status, companionUp, researchRoles],
  );
  const eligible = coins.filter((c) => c.eligible);
  const walletCheck = checks.find((c) => c.name === "wallet");
  const attention = nextFix(companionUp, status, checks, items, sessionAlive, net);
  const path = setupPath(companionUp, status, checks, sessionAlive, net);
  const doing = researchBusy
    ? `Researching ${researchCoin}`
    : preview?.eligible
      ? "Awaiting AUTHORIZE"
      : "Idle";
  const cost = researchBusy
    ? "Private compute (0G Direct)"
    : preview?.eligible
      ? "Trading capital only after you type AUTHORIZE"
      : "None";

  function finishSetup() {
    try {
      window.localStorage.setItem(SETUP_KEY, "1");
    } catch {
      /* ignore */
    }
    setSetupDone(true);
    setView("home");
  }

  async function onBind() {
    setBindBusy(true);
    setBindError(null);
    const r = await bindWallet(walletDraft, net);
    setBindBusy(false);
    if (r.error) {
      setBindError(describeBindError(r.error));
      return;
    }
    const s = await localStatus();
    if (s) {
      setStatus(s);
      if (s.wallet) setWalletDraft(s.wallet);
    }
    setChecks(await doctor());
  }

  async function onCheck() {
    setChecks(await doctor());
  }

  async function onSession() {
    setBindBusy(true);
    setBindError(null);
    const r = await createLocalSession();
    setBindBusy(false);
    if (r.error) {
      setBindError(describeBindError(r.error));
      return;
    }
    if (r.agent) setAgent(r.agent);
    const s = await localStatus();
    if (s) {
      setStatus(s);
      setSessionAlive(Boolean(s.sessionAlive));
      if (s.agent) setAgent(s.agent);
    }
    setChecks(await doctor());
  }

  async function onPolicy() {
    setBindBusy(true);
    setBindError(null);
    const r = await pinLocalPolicy();
    setBindBusy(false);
    if (r.error) {
      setBindError(describeBindError(r.error));
      return;
    }
    setPinned(true);
    setChecks(await doctor());
  }

  async function onRevoke() {
    setBindBusy(true);
    const r = await revokeLocalSession();
    setBindBusy(false);
    if (r.error) {
      setBindError(describeBindError(r.error));
      return;
    }
    setSessionAlive(false);
    setAgent("");
    setChecks(await doctor());
  }

  async function onKill(on: boolean) {
    setBindBusy(true);
    setBindError(null);
    const r = await setKillSwitch(on);
    setBindBusy(false);
    if (r.error) {
      setBindError(describeBindError(r.error));
      return;
    }
    setStatus((s) => (s ? { ...s, kill: on } : s));
    setChecks(await doctor());
  }

  async function researchThis(coin?: string) {
    const gen = ++researchGen.current;
    const want = (coin || "ETH").toUpperCase();
    setResearchNote(null);
    setResearchEvidence("");
    setAuthErr(null);
    setPollMiss(false);
    setResearchCoin(want);
    if (!companionUp) {
      setResearchStop("COMPANION_NOT_RUNNING");
      setView("research");
      return;
    }
    const sealer = checks.find((c) => c.name === "direct_sealer");
    if (sealer && !sealer.ok) {
      setResearchStop("DIRECT_PROVIDER_UNAVAILABLE");
      setView("research");
      return;
    }
    const auth = checks.find((c) => c.name === "direct_auth");
    if (auth && !auth.ok) {
      setResearchStop("DIRECT_NOT_AUTHORIZED");
      setView("research");
      return;
    }
    const credit = checks.find((c) => c.name === "direct_credit");
    if (credit && !credit.ok) {
      setResearchStop("DIRECT_CREDIT_INSUFFICIENT");
      setView("research");
      return;
    }
    setResearchBusy(true);
    setResearchStop(null);
    setResearchStage("READING_MARKET");
    setResearchElapsed(0);
    setView("research");
    const wall = Date.now();
    const tick = window.setInterval(() => {
      if (gen === researchGen.current) setResearchElapsed(Date.now() - wall);
    }, 250);
    try {
      const started = await startResearch(want, hypothesis);
      if (gen !== researchGen.current) return;
      if (started.error && !started.running) {
        setResearchStop(started.error);
        return;
      }
      applyStatus(started);
      await followJob(gen, wall);
    } catch (e) {
      if (gen !== researchGen.current) return;
      setPollMiss(true);
      await followJob(gen, wall);
    } finally {
      window.clearInterval(tick);
      if (gen === researchGen.current) {
        setResearchBusy(false);
        setPollMiss(false);
      }
    }

    function applyStatus(st: BindResult) {
      if (st.job_id) setResearchJobId(st.job_id);
      if (st.stage) setResearchStage(st.stage);
      if (typeof st.elapsed_ms === "number") setResearchElapsed(st.elapsed_ms);
      if (st.coin) setResearchCoin(String(st.coin).toUpperCase());
      const roles = Array.isArray(st.roles) ? st.roles : [];
      if (roles.length) setResearchRoles(roles);
      if (st.terminal_kind) setResearchKind(String(st.terminal_kind));
      if (st.preview) {
        setPreview(st.preview);
        setPreviewHash(st.preview_hash || st.preview.hash || "");
      }
      return roles;
    }

    async function followJob(genNow: number, startedAt: number) {
      for (;;) {
        await sleep(1000);
        if (genNow !== researchGen.current) return;
        const st = await researchStatus();
        if (st.transient) {
          setPollMiss(true);
          continue;
        }
        setPollMiss(false);
        const roles = applyStatus(st);
        const verified = committeeVerified(roles, st.verify);
        const deny = String(st.deny || st.preview?.deny || "");
        if ((verified || st.verify || committeeDeny(deny)) && !st.running) {
          setResearchStop(null);
          setResearchNote(st.note || "Sealed committee verified on this computer.");
          return;
        }
        if (st.running) continue;
        if (st.error && !verified && !committeeDeny(deny)) {
          setResearchStop(st.error);
          return;
        }
        if (roles.length === 0 && Date.now() - startedAt < 180000) {
          continue;
        }
        if (!verified) {
          setResearchStop("COMMITTEE_INCOMPLETE");
          return;
        }
        setResearchStop(null);
        setResearchNote(st.note || "Sealed committee verified on this computer.");
        return;
      }
    }
  }

  async function onCancelResearch() {
    researchGen.current += 1;
    setResearchBusy(false);
    const r = await cancelResearch();
    setResearchStop(r.error || "research_cancelled");
    if (r.stage) setResearchStage(r.stage);
    if (typeof r.elapsed_ms === "number") setResearchElapsed(r.elapsed_ms);
  }

  async function onAuthorize(e: FormEvent) {
    e.preventDefault();
    setAuthBusy(true);
    setAuthErr(null);
    const r = await authorizePreview(authTyped, previewHash);
    setAuthBusy(false);
    if (r.error) {
      setAuthErr(r.error);
      return;
    }
    if (r.oid) setLastOid(String(r.oid));
    setAuthTyped("");
    setChecks(await doctor());
  }

  async function onCancelBound(e: FormEvent) {
    e.preventDefault();
    setAuthBusy(true);
    setAuthErr(null);
    const r = await cancelBoundOrder(authTyped);
    setAuthBusy(false);
    if (r.error) {
      setAuthErr(r.error);
      return;
    }
    setAuthTyped("");
  }

  async function onConnectionPreview() {
    setBindBusy(true);
    setAuthErr(null);
    const r = await connectionPreview("ETH");
    setBindBusy(false);
    if (r.error) {
      setAuthErr(r.error);
      setView("research");
      return;
    }
    setPreview({
      eligible: true,
      kind: r.kind || "connection_test",
      market: r.market,
      side: r.side,
      sz: r.sz,
      limitPx: r.limitPx,
      hash: r.hash,
      cloid: r.cloid,
      note: r.note,
    });
    setPreviewHash(r.hash || "");
    setResearchStop(null);
    setResearchNote(r.note || "Connection test preview. This is not a research recommendation.");
    setView("research");
  }

  async function onReduceOnlyClose(coin: string) {
    setBindBusy(true);
    setAuthErr(null);
    const r = await connectionPreview(coin, true);
    setBindBusy(false);
    if (r.error) {
      setAuthErr(r.error);
      setView("research");
      return;
    }
    setPreview({
      eligible: true,
      kind: r.kind || "reduce_only_close",
      market: r.market,
      side: r.side,
      sz: r.sz,
      limitPx: r.limitPx,
      hash: r.hash,
      cloid: r.cloid,
      note: r.note,
    });
    setPreviewHash(r.hash || "");
    setResearchStop(null);
    setResearchNote(r.note || "Reduce-only close preview. Type AUTHORIZE on this computer. PIT cannot withdraw.");
    setView("research");
  }

  return (
    <div className="app">
      <aside className="rail">
        <div className="rail-brand">
          <div className="word">PIT.</div>
          <p className="kicker">{status?.version || "0.2.0"} · local execution</p>
        </div>
        <nav className="rail-nav" aria-label="Desk">
          {RAIL.map((item) => (
            <button
              key={item.id}
              type="button"
              className={view === item.id ? "on" : ""}
              onClick={() => {
                setSetupDone(true);
                setView(item.id);
              }}
            >
              {item.label}
            </button>
          ))}
        </nav>
        <div className="rail-foot">
          <p>{net === "mainnet" ? "MAINNET" : "TESTNET"}</p>
          <p>{companionUp ? "companion live" : "starting companion"}</p>
          <p>{status?.version || "PIT 0.2.0"}</p>
          <button type="button" className="ghost" onClick={() => setView("settings")}>
            Help / Diagnostics
          </button>
        </div>
      </aside>

      <div className="stage">
        <header className="bar">
          <div>
            <p className="eyebrow">
              WHERE · {net === "mainnet" ? "MAINNET" : "TESTNET"} · {view.toUpperCase()}
            </p>
            <p>{status?.workspace || status?.wallet || walletCheck?.detail || "unbound"}</p>
          </div>
          <div>
            <p className="eyebrow">WHAT PIT IS DOING</p>
            <p>{doing}</p>
            <p className="fine">Needs: {attention.title}</p>
            <p className="fine">Cost: {cost}</p>
          </div>
          <div className="bar-meta">
            <NetworkToggle net={net} onChange={setNet} />
            <p className="pair-chip">{code ? prettyCode(code) : companionUp ? "code rotating" : "starting companion"}</p>
            <p>Session {sessionAlive ? "order/cancel live" : "none"} · Compute {checks.find((c) => c.name === "direct_credit")?.ok ? "ready" : "action"}</p>
            {agent ? <p className="fine">PIT Agent {agent}</p> : null}
            <p className="island" role="status">
              {researchBusy
                ? `Job ${researchJobId || "…"} · ${researchStage} · ${Math.round(researchElapsed / 1000)}s`
                : activity.length
                  ? `${activity[activity.length - 1].kind || "event"} ${activity[activity.length - 1].market || ""} ${activity[activity.length - 1].status || ""}`.trim()
                  : "No new desk events"}
              {restartAllowed ? "" : " · Restart blocked while research runs"}
            </p>
          </div>
        </header>
        {companionUp && status?.version && !String(status.version).includes("0.2.0") ? (
          <article className="card stop" role="status">
            <p className="label">COMPANION VERSION</p>
            <p>
              This window expects PIT 0.2.0. The local companion is {status.version}. Close PIT, install the matching
              desktop, then launch again. A running sealed job is not cancelled by this warning.
            </p>
          </article>
        ) : null}

        {!setupDone ? (
          <SetupWizard
            step={setupStep}
            setStep={setSetupStep}
            net={net}
            setNet={setNet}
            items={items}
            code={code}
            expires={expires}
            companionUp={companionUp}
            walletDraft={walletDraft}
            setWalletDraft={setWalletDraft}
            bindError={bindError}
            bindBusy={bindBusy}
            boundWallet={status?.wallet || ""}
            pinned={pinned}
            agent={agent}
            agentName={status?.agentName}
            sessionAlive={sessionAlive}
            onBind={() => void onBind()}
            onSession={() => void onSession()}
            onPolicy={() => void onPolicy()}
            onResearch={() => void researchThis()}
            onCheck={() => void onCheck()}
            onDone={finishSetup}
            checks={checks}
            researchBusy={researchBusy}
            researchVerified={committeeVerified(researchRoles)}
          />
        ) : (
          <div className="desk-body">
            <CommandChat
              onNavigate={(v) => setView(v as View)}
              onResearch={(c) => void researchThis(c)}
              onOpenPreview={() => setView("research")}
            />
            <div className="book">

        {setupDone && view === "home" ? (
          <main className="page dense">
            <p className="eyebrow">HOME</p>
            <h1>Authorize on this computer. Never in the browser.</h1>
            <p className="lead">{NAMED.SEED_FORBIDDEN}</p>
            <NetworkBanner net={net} />
            <div className="desk-grid">
              <PairingBlock code={code} expires={expires} companionUp={companionUp} />
              <article className="card">
                <p className="label">READINESS</p>
                <ProbeList items={items} net={net} onGo={(v) => setView(v)} />
              </article>
            </div>
            <HyperliquidCard
              net={net}
              agent={agent}
              agentName={status?.agentName}
              sessionAlive={sessionAlive}
              sessionExpires={status?.sessionExpires}
              approved={Boolean(checks.find((c) => c.name === "hl_agent" && c.ok))}
              approvedDetail={checks.find((c) => c.name === "hl_agent")?.detail}
              busy={bindBusy}
              onCreateSession={() => void onSession()}
              onConnectionPreview={() => void onConnectionPreview()}
              onCheck={() => void onCheck()}
              onRevoke={() => void onRevoke()}
            />
            <PositionsPanel
              account={positionAccount || status?.wallet}
              positions={positions}
              error={positionErr}
              lastOrder={status?.lastOrder}
              onReduceOnlyClose={(c) => void onReduceOnlyClose(c)}
              closeBusy={bindBusy}
            />
            <EmptyHome count={eligible.length} next={attention} onGo={(v) => setView(v)} />
            <WelcomePath steps={path} onGo={(v) => setView(v)} />
            {companionStuck && !explained ? (
              <article className="card stop" role="status">
                <p className="label">LOCAL COMPANION</p>
                <h2>PIT is waiting for the process on this computer</h2>
                <p>
                  No order was placed. Close any old PIT window, reinstall so the local companion is included, then
                  launch PIT again. A terminal is not required.
                </p>
              </article>
            ) : null}
          </main>
        ) : null}

        {setupDone && view === "watch" ? (
          <main className="page dense">
            <p className="eyebrow">LIVE BOOKS</p>
            <h1>Watch</h1>
            <p className="lead">Public Hyperliquid marks only. This window cannot place an order.</p>
            {coins.length === 0 ? (
              <p className="fine">Empty Watch is the honest state until live books arrive.</p>
            ) : (
              <table className="desk-table">
                <thead>
                  <tr>
                    <th>Market</th>
                    <th>Mark</th>
                    <th>Oracle</th>
                    <th>Funding</th>
                    <th>Open interest</th>
                    <th>Policy</th>
                    <th>Research</th>
                    <th>Why</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  {coins.map((c) => (
                    <tr key={c.coin}>
                      <td>{c.coin}</td>
                      <td className="mark-num">{c.mark}</td>
                      <td>{c.oracle || "—"}</td>
                      <td>{c.funding ?? "—"}</td>
                      <td>{c.openInterest ? Math.round(c.openInterest) : "—"}</td>
                      <td>{c.eligible ? "PASS" : "BLOCKED"}</td>
                      <td>{checks.find((x) => x.name === "direct_credit")?.ok ? "Ready" : "Needs compute"}</td>
                      <td>{c.reason}</td>
                      <td>
                        <button type="button" className="linkish" onClick={() => void researchThis(c.coin)}>
                          Research this
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
            <p className="fine">Confidence NOT ENOUGH DATA. Side is not decided on this surface.</p>
          </main>
        ) : null}

        {setupDone && view === "research" ? (
          <main className="page dense">
            <p className="eyebrow">SEALED PATH</p>
            <h1>Research</h1>
            <p className="lead">
              Private book → Direct TeeML → researcher / challenger / risk → host size → policy → exact preview. Watch
              never places the order. Same provider is labeled as role separation, not three independent models.
            </p>
            <article className="card">
              <p className="label">SEALED HYPOTHESIS</p>
              <p>This is your intent, sealed into the private book. The committee may still stand down. Host sizes.</p>
              <div className="cta-row">
                {(["none", "long", "short"] as const).map((h) => (
                  <button
                    key={h}
                    type="button"
                    className={hypothesis === h ? "linkish on" : "linkish off"}
                    onClick={() => setHypothesis(h)}
                    disabled={researchBusy}
                  >
                    {h === "none" ? "No bias" : h === "long" ? "Consider long" : "Consider short"}
                  </button>
                ))}
              </div>
            </article>
            <ComputeCard checks={checks} onCheck={() => void onCheck()} />
            <article className="card">
              <p className="label">PRIVATE RESEARCH</p>
              <p>Provider Direct · glm-5.2 · Role separation, not three independent models. Estimated ~3 0G locked.</p>
              <button
                type="button"
                className="on"
                onClick={() => void researchThis(researchCoin)}
                disabled={researchBusy || !checks.find((c) => c.name === "direct_credit")?.ok}
              >
                Start Research
              </button>
            </article>
            {researchBusy ? (
              <ResearchProgress
                stage={researchStage}
                elapsedMs={researchElapsed}
                coin={researchCoin}
                roles={researchRoles}
                pollMiss={pollMiss}
                onCancel={() => void onCancelResearch()}
              />
            ) : null}
            {researchNote && !explained && !researchBusy ? (
              <article className="card">
                <p className="label">{researchCardTitle(researchKind, committeeVerified(researchRoles))}</p>
                <p>{researchNote}</p>
                {researchJobId ? <p className="fine">Job {researchJobId}</p> : null}
                <ul className="doctor">
                  {researchRoles.map((role) => (
                    <li key={role.role}>
                      <strong>{role.verify_e2ee}</strong> {role.role}
                      {role.proposed_side ? ` side ${role.proposed_side}` : ""}
                      {role.survives === false ? " stood down" : ""}
                      {role.kill ? " kill" : ""}
                      {role.pubkey_signer ? ` ${role.pubkey_signer}` : ""}
                    </li>
                  ))}
                </ul>
              </article>
            ) : null}
            {committeeCopy && !researchBusy ? (
              <article className="card" role="status">
                <p className="label">COMMITTEE DECISION</p>
                <h2>{committeeCopy.title}</h2>
                <p>{committeeCopy.body}</p>
                <p className="fine">Change the sealed hypothesis or wait for a different book. PIT will not invent a side.</p>
              </article>
            ) : null}
            {!researchBusy && preview ? (
              <article className="card">
                <p className="label">EXACT PREVIEW</p>
                {preview.eligible ? (
                  <>
                    {preview.kind === "connection_test" ? (
                      <p className="fine">Connection test. This is not a research recommendation. Host sized a policy clip.</p>
                    ) : preview.kind === "reduce_only_close" ? (
                      <p className="fine">Reduce-only close. This is not a research recommendation. Type AUTHORIZE to send it. PIT cannot withdraw.</p>
                    ) : (
                      <p className="label">OPPORTUNITY FOUND</p>
                    )}
                    <p>
                      {preview.market} {preview.side} {preview.sz} @ {preview.limitPx}
                    </p>
                    <p className="fine">Hash {preview.hash || previewHash}. Model size was ignored. Host sized this clip.</p>
                    {sessionAlive ? (
                      <form onSubmit={(e) => void onAuthorize(e)}>
                        <input
                          aria-label="type AUTHORIZE"
                          autoComplete="off"
                          value={authTyped}
                          onChange={(ev) => setAuthTyped(ev.target.value)}
                          placeholder="Type AUTHORIZE"
                        />
                        <button type="submit" disabled={authBusy || !previewHash}>
                          Authorize
                        </button>
                      </form>
                    ) : (
                      <p>Create a local session, then type AUTHORIZE here.</p>
                    )}
                    {oidBelongsToPreview(status?.lastOrder?.hash, previewHash, preview.hash) && lastOid ? (
                      status?.lastOrder?.status !== "filled" && !status?.lastOrder?.cancelled ? (
                        <form onSubmit={(e) => void onCancelBound(e)}>
                          <p className="fine">OID {lastOid} is resting for this preview. Type AUTHORIZE again to cancel. PIT cannot withdraw.</p>
                          <button type="submit" disabled={authBusy}>
                            Cancel this order
                          </button>
                        </form>
                      ) : status?.lastOrder?.status === "filled" ? (
                        <p className="fine">
                          OID {lastOid} FILLED for this preview. Flatten only with a reduce-only close that YOU authorize. PIT cannot withdraw.
                        </p>
                      ) : null
                    ) : null}
                    {authErr ? (
                      <p className="err" role="alert">
                        {authErr === "approveAgent_required"
                          ? "Approve this agent on Hyperliquid before PIT will send an order."
                          : authErr}
                      </p>
                    ) : null}
                    {authErr === "approveAgent_required" ? (
                      <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
                        Open Hyperliquid
                      </a>
                    ) : null}
                  </>
                ) : (
                  <p>
                    {preview.deny === "no_side"
                      ? "Stand down. The committee did not propose a side. This is a verified result, not a crash. No order was placed."
                      : `Host did not size a trade (${preview.deny || "no_side"}). Committee envelopes stay independent. The model cannot raise clip. No order was placed.`}
                  </p>
                )}
              </article>
            ) : null}
            {explained && !researchBusy ? (
              <article className="card stop" role="alert">
                <p className="label">RESEARCH STOPPED</p>
                <h2>{explained.title}</h2>
                <p>{explained.body}</p>
                {researchRoles.length ? (
                  <ul className="doctor">
                    {researchRoles.map((role) => (
                      <li key={role.role}>
                        <strong>{role.verify_e2ee || "STOP"}</strong> {role.role}
                      {role.proposed_side ? ` side ${role.proposed_side}` : ""} {role.pubkey_signer}
                      </li>
                    ))}
                  </ul>
                ) : null}
                {explainStopHref(researchStop) ? (
                  <a className="linkish" href={explainStopHref(researchStop)?.href} target="_blank" rel="noreferrer">
                    {explainStopHref(researchStop)?.label}
                  </a>
                ) : null}
                <button type="button" className="linkish" onClick={() => setTechOpen((v) => !v)}>
                  {techOpen ? "Hide technical evidence" : "View technical evidence"}
                </button>
                {techOpen ? (
                  <pre className="pipe evidence">
                    Code {researchStop}
                    {"\n"}
                    {researchEvidenceText || "Verification is fail-closed. Router fallback is impossible."}
                  </pre>
                ) : null}
                <button type="button" onClick={() => void researchThis(researchCoin)} disabled={researchBusy}>
                  Retry
                </button>
              </article>
            ) : null}
            {!researchBusy && !researchNote && !explained ? (
              <p className="fine">Private research has not been run on this machine in this session.</p>
            ) : null}
            <PreviewNote />
            <EvidenceDrawer
              jobId={researchJobId}
              roles={researchRoles}
              preview={preview}
              previewHash={previewHash}
              kind={researchKind}
              deny={preview?.deny}
              evidence={researchEvidenceText}
            />
            {eligible.length === 0 ? (
              <p className="fine">No policy-eligible market is waiting. Watch does not invent cards.</p>
            ) : (
              <ul className="watch-grid">
                {eligible.map((c) => (
                  <li key={c.coin} className="card">
                    <p className="label">{c.coin}</p>
                    <p className="mark-num">{c.mark}</p>
                    <p>Policy PASS</p>
                    <p className="fine">{c.reason}</p>
                    <button type="button" className="linkish" onClick={() => void researchThis(c.coin)}>
                      Research this
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </main>
        ) : null}

        {setupDone && view === "activity" ? (
          <main className="page dense">
            <p className="eyebrow">LEDGER</p>
            <h1>Activity</h1>
            <p className="lead">Exact-once orders, cancels, receipts, and stops. Empty is honest until this machine records one.</p>
            <ActivityTimeline events={activity} lastOid={lastOid} />
            <article className="card">
              <p className="label">LAST VENUE ORDER</p>
              {lastOid ? (
                <>
                  <p>
                    {status?.lastOrder?.status === "filled"
                      ? "ORDER FILLED"
                      : status?.lastOrder?.cancelled
                        ? "ORDER CANCELED"
                        : "ORDER SUBMITTED"}
                  </p>
                  <p>OID {lastOid}</p>
                  {status?.lastOrder?.market ? (
                    <p>
                      {status.lastOrder.market} {status.lastOrder.side} {status.lastOrder.sz}
                    </p>
                  ) : null}
                  {status?.lastOrder?.status === "filled" ? (
                    <p className="fine">This size is a position. Cancel does not apply to a filled order.</p>
                  ) : null}
                </>
              ) : (
                <p>No order id is shown until Hyperliquid accepts one after you type AUTHORIZE.</p>
              )}
            </article>
          </main>
        ) : null}

        {setupDone && view === "policy" ? (
          <main className="page dense">
            <p className="eyebrow">CONSTRAINTS</p>
            <h1>Policy</h1>
            <p className="lead">Host engine enforces this. The model cannot raise clip, leverage, or permissions.</p>
            <PolicyLaw pinned={pinned} onPin={() => void onPolicy()} busy={bindBusy} />
          </main>
        ) : null}

        {setupDone && view === "security" ? (
          <main className="page dense">
            <p className="eyebrow">PERMISSIONS</p>
            <h1>Security</h1>
            <p className="lead">Order and cancel only. Withdraw is impossible through PIT.</p>
            <SecurityCenter
              domains={securityDomains}
              net={net}
              onSession={() => void onSession()}
              onPolicy={() => void onPolicy()}
              onCheck={() => void onCheck()}
              busy={bindBusy}
            />
            <HyperliquidCard
              net={net}
              agent={agent}
              agentName={status?.agentName}
              sessionAlive={sessionAlive}
              sessionExpires={status?.sessionExpires}
              approved={Boolean(checks.find((c) => c.name === "hl_agent" && c.ok))}
              approvedDetail={checks.find((c) => c.name === "hl_agent")?.detail}
              busy={bindBusy}
              onCreateSession={() => void onSession()}
              onConnectionPreview={() => void onConnectionPreview()}
              onCheck={() => void onCheck()}
              onRevoke={() => void onRevoke()}
            />
            <AuthorizeGate
              sessionAlive={sessionAlive}
              agent={agent}
              agentName={status?.agentName}
              net={net}
              busy={bindBusy}
              onCreateSession={() => void onSession()}
            />
            <PermissionsCard />
            <SessionNote />
            <KillNote />
            <article className="card">
              <p className="label">KILL SWITCH</p>
              <p>You flip this. The model cannot. New orders stop on this workspace until you turn it off.</p>
              <button type="button" className="linkish" onClick={() => void onKill(true)} disabled={bindBusy || Boolean(status?.kill)}>
                {status?.kill ? "Kill switch is on" : "Halt new orders"}
              </button>
              {status?.kill ? (
                <button type="button" className="linkish" onClick={() => void onKill(false)} disabled={bindBusy}>
                  Resume this workspace
                </button>
              ) : null}
            </article>
            <article className="card">
              <p className="label">REVOKE</p>
              <p>Delete the local session, then remove the PIT agent from your Hyperliquid account.</p>
              <button type="button" className="linkish" onClick={() => void onRevoke()} disabled={bindBusy || !sessionAlive}>
                Revoke local session
              </button>
              <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
                Open Hyperliquid
              </a>
            </article>
          </main>
        ) : null}

        {setupDone && view === "account" ? (
          <main className="page dense">
            <p className="eyebrow">IDENTITY</p>
            <h1>Account</h1>
            <p className="lead">{NAMED.TWO_WALLETS}</p>
            <article className="card">
              <p className="label">THIS WORKSPACE</p>
              <p>Wallet {walletCheck?.ok ? walletCheck.detail : "unbound"}</p>
              <p>Network {net}</p>
              <p>Agent {agent || "none"}</p>
              {status?.agentName ? <p>Agent name {status.agentName} (under 17 characters on Hyperliquid)</p> : null}
              <p>Session {sessionAlive ? "order/cancel live" : "none"}</p>
              <p>Kill {status?.kill ? "on" : "off"}</p>
            </article>
            <article className="card">
              <p className="label">DESK IDENTITY</p>
              <p>{identityNote}</p>
              <p className="fine">Identity is optional. Trading does not wait on mint. Transfer of Agentic ID is unavailable on mainnet.</p>
            </article>
            <article className="card">
              <p className="label">CALIBRATION</p>
              <p>{calibCopy}</p>
            </article>
            {!sessionAlive ? (
              <article className="card">
                <p className="label">CONNECT TRADING</p>
                <p>Create an order/cancel session on this computer, then approve that agent on Hyperliquid.</p>
                <button type="button" className="linkish" onClick={() => void onSession()} disabled={bindBusy || !companionUp}>
                  Create local session
                </button>
                {agent ? <p className="fine">Agent {agent}</p> : null}
                <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
                  Open Hyperliquid API
                </a>
              </article>
            ) : null}
            <p className="fine">{NAMED.TRANSFER_NOT_LIVE}</p>
          </main>
        ) : null}

        {setupDone && view === "settings" ? (
          <main className="page dense">
            <p className="eyebrow">DIAGNOSTICS</p>
            <h1>Settings</h1>
            <NetworkToggle net={net} onChange={setNet} />
            <NetworkBanner net={net} />
            <CapabilityMatrix net={net} />
            <article className="card">
              <p className="label">UPDATES</p>
              <p>{updateNote}</p>
              <p>Restart {restartAllowed ? "allowed" : "refused — research is running. PIT will not replace pit.exe under a live job."}</p>
              <a className="linkish" href={LINKS.releases} target="_blank" rel="noreferrer">
                Open release
              </a>
            </article>
            <article className="card">
              <p className="label">MEMORY</p>
              <p>Forget wipes working memory and chat on this workspace. Receipts and venue positions stay.</p>
              <button
                type="button"
                className="linkish"
                onClick={() => {
                  void forgetMemory();
                }}
              >
                Forget this workspace memory
              </button>
            </article>
            <article className="card">
              <p className="label">DOCTOR</p>
              {checks.length === 0 ? (
                <p>Waiting for the local companion on 127.0.0.1:17373.</p>
              ) : (
                <ul className="doctor">
                  {checks.map((c) => (
                    <li key={c.name}>
                      <strong>{c.ok ? "ok" : "fail"}</strong> {c.name} — {c.detail}
                    </li>
                  ))}
                </ul>
              )}
            </article>
            <article className="card">
              <p className="label">OFFICIAL LINKS</p>
              <a className="linkish" href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
                Open 0G Private Compute
              </a>
              <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
                Open Hyperliquid API
              </a>
              <a className="linkish" href={LINKS.releases} target="_blank" rel="noreferrer">
                Open release
              </a>
              <a className="linkish" href={LINKS.pair} target="_blank" rel="noreferrer">
                Open pairing
              </a>
            </article>
            <p className="fine">{NAMED.TWO_WALLETS}</p>
            <p className="fine">{NAMED.TRANSFER_NOT_LIVE}</p>
          </main>
        ) : null}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
