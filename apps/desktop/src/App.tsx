import { FormEvent, useEffect, useMemo, useRef, useState } from "react";
import { NetworkToggle } from "./NetworkToggle";
import { committeeDeny } from "./committee";
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
  fetchAutomation,
  fetchCalibration,
  fetchChatThreads,
  fetchIdentity,
  fetchMission,
  fetchPositions,
  fetchSecurity,
  fetchUpdate,
  fetchWatch,
  forgetMemory,
  localStatus,
  mutateChatThread,
  pairCode,
  rotatePairCode,
  pinLocalPolicy,
  fetchPolicy,
  previewPolicy,
  fetchDemo,
  setDemoMode,
  postMission,
  researchEvidence,
  researchStatus,
  revokeLocalSession,
  saveAutomation,
  setKillSwitch,
  startResearch,
  wakeCompanion,
  type AccountSummary,
  type ActivityEvent,
  type BindResult,
  type ChatThread,
  type DoctorCheck,
  type LocalStatus,
  type MissionPublic,
  type AutoPrefs,
  type HostPolicy,
  type SecurityDomain,
  type VenuePosition,
} from "./companion";
import { LINKS, hyperliquidAPI, hyperliquidApp } from "./links";
import { openExternal, useNativeExternalLinks } from "./open";
import { nextFix } from "./nextFix";
import { probes } from "./readiness";
import { committeeVerified } from "./honesty";
import { AutomationCenter } from "./AutomationCenter";
import { CommandChat } from "./CommandChat";
import { ActivityTimeline } from "./ActivityTimeline";
import { ProofTimeline } from "./ProofTimeline";
import { PositionsPanel } from "./PositionsPanel";
import { SecurityCenter } from "./SecurityCenter";
import { SetupWizard } from "./SetupWizard";
import { TitleBar } from "./TitleBar";
import { BootGate } from "./BootGate";
import { CommandPalette } from "./CommandPalette";
import { ThreadRail } from "./ThreadRail";
import { DeskHome } from "./DeskHome";
import { WatchBook, type MarketCoin } from "./WatchBook";
import { ResearchBoard } from "./ResearchBoard";
import { StrategyHealth } from "./StrategyHealth";
import { askNotify, deskNotify } from "./notify";
import { DESKTOP_VERSION } from "./version";

type Net = "mainnet" | "testnet";
type View = "home" | "chat" | "markets" | "research" | "portfolio" | "activity" | "automation" | "health" | "security";

type Coin = MarketCoin;

function mapView(v: string): View {
  if (v === "watch" || v === "markets") return "markets";
  if (v === "positions" || v === "portfolio") return "portfolio";
  if (v === "policy" || v === "account" || v === "settings") return "security";
  if (v === "automation") return "automation";
  if (v === "health" || v === "calibration" || v === "skills") return "health";
  if (v === "preview") return "research";
  if (v === "home" || v === "chat" || v === "research" || v === "activity" || v === "security") return v;
  return "home";
}

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

function eventLine(ev: ActivityEvent) {
  const kind = String(ev.kind || "event");
  const label =
    kind === "opportunity"
      ? "Opportunity"
      : kind.startsWith("research")
        ? "Research"
        : kind.includes("preview")
          ? "Preview"
          : kind.includes("order") || kind.includes("fill")
            ? "Order"
            : kind.replaceAll(".", " ");
  return `${label} ${ev.market || ""} ${ev.status || ""}`.trim();
}

const SETUP_KEY = "pit.desk.setup";

const RAIL: { id: View; label: string }[] = [
  { id: "home", label: "Desk" },
  { id: "chat", label: "Chat" },
  { id: "markets", label: "Markets" },
  { id: "research", label: "Research" },
  { id: "portfolio", label: "Portfolio" },
  { id: "activity", label: "Activity" },
  { id: "automation", label: "Automation" },
  { id: "health", label: "Health" },
  { id: "security", label: "Security" },
];

export function App() {
  useNativeExternalLinks();
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
  const [calibN, setCalibN] = useState(0);
  const [calibNeed, setCalibNeed] = useState(0);
  const [calibEnough, setCalibEnough] = useState(false);
  const [calibSkills, setCalibSkills] = useState<Array<{ id?: string; title?: string; version?: string; n?: number; copy?: string }>>([]);
  const [identityNote, setIdentityNote] = useState("Transfer of Agentic ID is not live on mainnet.");
  const [updateNote, setUpdateNote] = useState("This build is checksum-verified, not OS-signed.");
  const [restartAllowed, setRestartAllowed] = useState(true);
  const [hypothesis, setHypothesis] = useState<"none" | "long" | "short">("none");
  const researchGen = useRef(0);
  const researchBusyRef = useRef(false);
  researchBusyRef.current = researchBusy;
  const walletBoundRef = useRef("");
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
  const [booted, setBooted] = useState(false);
  const [palette, setPalette] = useState(false);
  const [thread, setThread] = useState("desk");
  const [threads, setThreads] = useState<ChatThread[]>([{ id: "desk", title: "Desk" }]);
  const [memoryEpoch, setMemoryEpoch] = useState(0);
  const [summary, setSummary] = useState<AccountSummary>({});
  const [hostPolicy, setHostPolicy] = useState<HostPolicy | null>(null);
  const [policyHash, setPolicyHash] = useState("");
  const [policyConsequences, setPolicyConsequences] = useState<string[]>([]);
  const [policyAllowed, setPolicyAllowed] = useState<string[]>([]);
  const [policyRefused, setPolicyRefused] = useState<string[]>([]);
  const [capitalNote, setCapitalNote] = useState("");
  const [buyingPower, setBuyingPower] = useState<number | undefined>(undefined);
  const [powerSource, setPowerSource] = useState("");
  const [demoReplay, setDemoReplay] = useState(false);
  const [autoPrefs, setAutoPrefs] = useState<AutoPrefs>({ watch: true, notify: true, auto_research: false, cadence_minutes: 15, trigger: "policy_pass" });
  const [mission, setMission] = useState<MissionPublic>({ mode: "manual", mission: { mode: "manual" } });
  const [confirmHours, setConfirmHours] = useState(8);
  const [openConfirm, setOpenConfirm] = useState(false);
  const [bestWhy, setBestWhy] = useState("");
  const [scanned, setScanned] = useState(0);
  const fillKey = useRef("");
  const lastNotify = useRef(0);

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
    const tick = () => {
      void fetchActivity().then((ev) => {
        if (!gone) setActivity(ev);
      });
      void fetchMission().then((m) => {
        if (!gone) setMission(m);
      });
    };
    tick();
    const id = window.setInterval(tick, 3000);
    return () => {
      gone = true;
      window.clearInterval(id);
    };
  }, []);

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
            if (s?.wallet) {
              walletBoundRef.current = s.wallet;
              setWalletDraft((cur) => cur || s.wallet || "");
            }
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
            setCoins((Array.isArray(body.coins) ? body.coins : []) as Coin[]);
            if (body.bestWhy) setBestWhy(body.bestWhy);
            if (typeof body.scanned === "number") setScanned(body.scanned);
            if (body.capitalNote) setCapitalNote(body.capitalNote);
            if (typeof body.buyingPower === "number") setBuyingPower(body.buyingPower);
            if (body.powerSource) setPowerSource(body.powerSource);
          })
          .catch(() => {
            if (!gone) setCoins([]);
          })
          .finally(() => {
            watchBusy = false;
          });
        void fetchActivity().then((ev) => {
          if (!gone) {
            setActivity(ev);
            const last = ev[ev.length - 1];
            if (last?.kind === "opportunity" && last.ts && last.ts !== lastNotify.current) {
              lastNotify.current = last.ts;
              deskNotify("Watch", `${last.market || "A market"} is interesting under your policy.`);
            }
          }
        });
        void fetchPositions().then((p) => {
          if (gone) return;
          setPositions(p.positions);
          setPositionAccount(p.account || "");
          setPositionErr(p.error || "");
          setSummary(p.summary || {});
        });
        void fetchPolicy().then((p) => {
          if (gone || !p.policy) return;
          setHostPolicy(p.policy);
          if (typeof p.pinned === "boolean") setPinned(p.pinned);
          if (p.hash) setPolicyHash(p.hash);
          if (p.allowed) setPolicyAllowed(p.allowed);
          if (p.refused) setPolicyRefused(p.refused);
        });
        void fetchDemo().then((d) => {
          if (!gone) setDemoReplay(d.pref === "replay");
        });
        void fetchSecurity().then((d) => {
          if (!gone) setSecurityDomains(d);
        });
        void fetchCalibration().then((c) => {
          if (gone) return;
          setCalibCopy(c.copy || "NOT ENOUGH DATA");
          setCalibN(c.n || 0);
          setCalibNeed(c.need || 0);
          setCalibEnough(Boolean(c.enough));
          setCalibSkills(c.skills || []);
        });
        void fetchIdentity().then((id) => {
          if (!gone) setIdentityNote(id.note || "Transfer of Agentic ID is not live on mainnet.");
        });
        void fetchUpdate().then((u) => {
          if (gone) return;
          setRestartAllowed(u.restart_allowed !== false);
          setUpdateNote(u.note || "This build is checksum-verified, not OS-signed.");
        });
        void fetchMission().then((m) => {
          if (!gone) setMission(m);
        });
      }
    };
    const loop = () => {
      tick();
      if (!gone) timer = window.setTimeout(loop, 4000);
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

  useEffect(() => {
    askNotify();
  }, []);

  useEffect(() => {
    if (companionUp || ticks >= 8) setBooted(true);
  }, [companionUp, ticks]);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setPalette((v) => !v);
      }
      if (e.key === "Escape") setPalette(false);
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  useEffect(() => {
    if (!companionUp) return;
    let gone = false;
    void fetchChatThreads().then((rows) => {
      if (!gone && rows.length) setThreads(rows);
    });
    void fetchAutomation().then((p) => {
      if (!gone) setAutoPrefs(p);
    });
    void fetchMission().then((m) => {
      if (!gone) setMission(m);
    });
    return () => {
      gone = true;
    };
  }, [companionUp, thread]);

  useEffect(() => {
    if (status?.lastOrder?.status === "filled" && lastOid) {
      const key = `fill-${lastOid}`;
      if (fillKey.current === key) return;
      fillKey.current = key;
      deskNotify("Order filled", `OID ${lastOid}. Historical fills never appear inside a new preview.`);
    }
  }, [status?.lastOrder?.status, lastOid]);

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
  const attention = nextFix(companionUp, status, checks, items, sessionAlive, net, {
    buyingPower: buyingPower ?? Number(summary.buyingPower || 0),
    execGate: summary.execGate,
    execWhy: summary.execWhy,
  });
  const showChat = view === "chat";
  const protectedOk = Boolean(checks.find((c) => c.name === "direct_auth")?.ok);
  const computeReady = Boolean(checks.find((c) => c.name === "direct_credit")?.ok);
  const researchAllowed = protectedOk && computeReady && pinned;
  const awaitingAuth = Boolean(preview?.eligible && pinned && sessionAlive && !researchBusy);
  const doing = researchBusy
    ? `Researching ${researchCoin}`
    : awaitingAuth
      ? "Awaiting AUTHORIZE"
      : mission.status === "ACTIVE" || status?.missionRunning || mission.running
        ? mission.mode === "guarded"
          ? "Guarded Autonomy live"
          : `Mission ${status?.mode || mission.mode}`
        : attention.title === "Desk is ready"
          ? "Idle"
          : attention.title;

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

  async function onRotatePair() {
    const p = await rotatePairCode();
    if (p?.code) setCode(p.code);
    if (p?.expires) setExpires(p.expires);
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

  async function onPolicy(draft?: HostPolicy) {
    setBindBusy(true);
    setBindError(null);
    const r = await pinLocalPolicy(draft);
    setBindBusy(false);
    if (r.error) {
      setBindError(describeBindError(r.error));
      return;
    }
    setPinned(true);
    if (r.policy) setHostPolicy(r.policy);
    if (r.hash) setPolicyHash(r.hash);
    if (r.consequences) setPolicyConsequences(r.consequences);
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
    const want = (coin || coins.find((c) => c.eligible)?.coin || "").toUpperCase();
    setResearchNote(null);
    setResearchEvidence("");
    setAuthErr(null);
    setPollMiss(false);
    if (!want) {
      setResearchNote("No eligible market yet. Open Markets and wait for live books.");
      setView("markets");
      return;
    }
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
      <TitleBar status={doing} />
      <BootGate
        open={!booted}
        rows={[
          { id: "ws", label: "loading secure workspace", state: companionUp ? "ok" : "wait" },
          { id: "companion", label: "starting companion", state: companionUp ? "ok" : ticks >= 8 ? "fail" : "wait" },
          { id: "wallet", label: "checking wallet binding", state: items.find((p) => p.id === "wallet")?.state === "ok" ? "ok" : "wait" },
          { id: "protect", label: "checking Protect signature", state: checks.find((c) => c.name === "direct_auth")?.ok ? "ok" : "wait" },
          { id: "compute", label: "checking private compute", state: checks.find((c) => c.name === "direct_credit")?.ok ? "ok" : "wait" },
          { id: "policy", label: "checking policy", state: pinned ? "ok" : "wait" },
          { id: "session", label: "checking trading session", state: sessionAlive ? "ok" : "wait" },
          { id: "venue", label: "checking venue connection", state: checks.find((c) => c.name === "hl_agent")?.ok ? "ok" : "wait" },
        ]}
        stuck={companionStuck ? "Companion is slow on this computer. Close old PIT windows and launch again. A sealed job is not cancelled." : undefined}
      />
      <CommandPalette
        open={palette}
        onClose={() => setPalette(false)}
        actions={[
          { id: "desk", label: "Open Desk", run: () => setView("home") },
          { id: "chat", label: "Open Chat", run: () => setView("chat") },
          { id: "markets", label: "Open Markets", run: () => setView("markets") },
          { id: "research", label: "Open Research", run: () => setView("research") },
          { id: "portfolio", label: "Open Portfolio", run: () => setView("portfolio") },
          { id: "activity", label: "Open Activity", run: () => setView("activity") },
          { id: "automation", label: "Open Automation", run: () => setView("automation") },
          { id: "health", label: "Open Strategy Health", run: () => setView("health") },
          { id: "security", label: "Open Security", run: () => setView("security") },
          { id: "start", label: "Start research", run: () => void researchThis() },
          { id: "hl", label: "Open Hyperliquid", run: () => void openExternal(hyperliquidApp(net)) },
          { id: "hlapi", label: "Open Hyperliquid API", run: () => void openExternal(hyperliquidAPI(net)) },
          { id: "og", label: "Open 0G Private Compute", run: () => void openExternal(LINKS.pcAdvanced) },
          { id: "check", label: "Check system", run: () => void onCheck() },
          { id: "preview", label: "Show current preview", run: () => setView("research") },
          { id: "connecttest", label: "Prepare connection-test preview", run: () => void onConnectionPreview() },
          { id: "act", label: "Show latest activity", run: () => setView("activity") },
        ]}
      />
      <div className="work">
      <aside className="rail">
        <div className="rail-brand">
          <div className="word">PIT.</div>
          <p className="kicker">{status?.version || DESKTOP_VERSION} · local execution</p>
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
          <p>{status?.wallet ? `${status.wallet.slice(0, 6)}…${status.wallet.slice(-4)}` : "unbound"}</p>
          <p>{sessionAlive ? "session live" : "no session"}</p>
          <p>{checks.find((c) => c.name === "direct_credit")?.ok ? "compute ready" : "compute action"}</p>
          <p>{companionUp ? "companion live" : "starting companion"}</p>
          <button type="button" className="ghost" onClick={() => setView("security")}>
            Help / Diagnostics
          </button>
        </div>
      </aside>

      <div className="stage">
        <header className="bar">
          <div>
            <p className="eyebrow">
              {net === "mainnet" ? "MAINNET" : "TESTNET"} · {view === "home" ? "DESK" : view.toUpperCase()}
            </p>
            <p className="bar-doing">
              {doing}
              {researchBusy ? ` · ${Math.round(researchElapsed / 1000)}s` : ""}
            </p>
          </div>
          <p className="fine" style={{ margin: 0 }}>
            {attention.title}
          </p>
          <div className="bar-meta">
            <div className="bar-chips" aria-label="System">
              <span className={companionUp ? "chip ok" : "chip fail"}>Health {companionUp ? "live" : "down"}</span>
              <span className={mission.status === "ACTIVE" || mission.running ? "chip ok" : "chip"}>
                Autonomy {mission.mode === "guarded" && mission.running ? "live" : mission.mode === "research_only" ? "research" : "manual"}
              </span>
              <span className={sessionAlive ? "chip ok" : "chip fail"}>Session {sessionAlive ? "live" : "none"}</span>
              <span className={checks.find((c) => c.name === "direct_credit")?.ok ? "chip ok" : "chip fail"}>
                Compute {checks.find((c) => c.name === "direct_credit")?.ok ? "ready" : "action"}
              </span>
              <span className={status?.wallet ? "chip ok" : "chip"}>
                Account {status?.wallet ? `${status.wallet.slice(0, 6)}…${status.wallet.slice(-4)}` : "unbound"}
              </span>
              <span className={demoReplay ? "chip fail" : "chip ok"}>{demoReplay ? "REPLAY" : "LIVE"}</span>
            </div>
            <NetworkToggle net={net} onChange={setNet} />
            {mission.mode === "guarded" && mission.running ? (
              <button
                type="button"
                className="kill-switch compact"
                onClick={() => {
                  void postMission({ stop: true }).then(setMission);
                }}
              >
                Stop autonomy
              </button>
            ) : null}
          </div>
        </header>
        {companionUp && status?.version && String(status.version) !== DESKTOP_VERSION ? (
          <article className="card stop" role="status">
            <p className="label">COMPANION VERSION</p>
            <p>
              This window expects PIT {DESKTOP_VERSION}. The local companion is {status.version}. Close PIT, install the matching
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
          <div className={showChat ? "desk-body with-threads" : "desk-body solo"}>
            {showChat ? (
              <ThreadRail
                threads={threads}
                active={thread}
                onSelect={setThread}
                onNew={() => {
                  void mutateChatThread("new").then((r) => {
                    if (r.threads) setThreads(r.threads);
                    if (r.id) setThread(r.id);
                  });
                }}
                onRename={(id, title) => {
                  void mutateChatThread("rename", id, title).then((r) => {
                    if (r.threads) setThreads(r.threads);
                  });
                }}
                onDelete={(id) => {
                  void mutateChatThread("delete", id).then((r) => {
                    if (r.threads) setThreads(r.threads);
                    if (thread === id) setThread("desk");
                  });
                }}
              />
            ) : null}
            {showChat ? (
            <CommandChat
              key={`${thread}-${memoryEpoch}`}
              thread={thread}
              onNavigate={(v) => {
                if (v === "setup") {
                  setSetupDone(false);
                  setSetupStep(0);
                  return;
                }
                if (v === "preview") setView("research");
                else setView(mapView(v));
              }}
              onConfirmAutonomy={(hours) => {
                setConfirmHours(hours || 8);
                setOpenConfirm(true);
                setView("automation");
              }}
              onResearch={(c) => void researchThis(c)}
              onOpenPreview={() => setView("research")}
              onStop={() => void onCancelResearch()}
              live={{
                coin: mission.mission?.best_coin || researchCoin,
                stage: mission.stage || (researchBusy ? researchStage : ""),
                gate: mission.block_reason,
                awaiting: awaitingAuth,
                running: Boolean(mission.running || mission.mission?.running),
              }}
              island={{
                busy: researchBusy,
                coin: researchCoin,
                stage: researchStage,
                elapsedMs: researchElapsed,
                jobId: researchJobId,
                pollMiss,
                roles: researchRoles,
                kind: researchKind,
              }}
            />
            ) : null}
            {showChat ? null : (
            <div className="book">

        {setupDone && view === "home" ? (
            <DeskHome
              ready={attention.title === "Desk is ready"}
              doing={doing}
              items={items}
              attention={attention}
              code={code}
              companionUp={companionUp}
              sessionAlive={sessionAlive}
              computeReady={computeReady}
              protectedOk={protectedOk}
              policyPinned={pinned}
              hlApproved={Boolean(checks.find((c) => c.name === "hl_agent" && c.ok))}
              researchBusy={researchBusy}
              researchStage={researchStage}
              researchKind={researchKind}
              awaitingAuth={awaitingAuth}
              expires={expires}
              paired={Boolean(status?.paired)}
              pairingDevices={status?.pairingDevices}
              onRotatePair={() => void onRotatePair()}
              coins={coins}
              lastEvent={activity.length ? eventLine(activity[activity.length - 1]) : undefined}
              mode={status?.mode || mission.mode}
              exposure={summary.totalNtlPos || summary.accountValue}
              onResearch={(c) => void researchThis(c)}
              onGo={(v) => setView(v)}
            />
        ) : null}

        {setupDone && view === "markets" ? (
          <WatchBook
            coins={coins}
            bestWhy={bestWhy}
            scanned={scanned}
            execGate={summary.execGate}
            execWhy={summary.execWhy}
            computeReady={researchAllowed}
            researchBusy={researchBusy}
            capitalNote={capitalNote || summary.capitalNote}
            buyingPower={buyingPower ?? Number(summary.buyingPower || 0)}
            powerSource={powerSource || summary.powerSource}
            fundHref={hyperliquidApp(net)}
            pinned={pinned}
            onPin={() => setView("security")}
            onResearch={(c) => void researchThis(c)}
          />
        ) : null}

        {setupDone && view === "portfolio" ? (
          <main className="page dense">
            <p className="eyebrow">Portfolio</p>
            <h1>Live exposure</h1>
            <p className="lead">Venue truth for the connected trading account. PIT cannot withdraw.</p>
            <PositionsPanel
              account={positionAccount || status?.wallet}
              positions={positions}
              error={positionErr}
              lastOrder={status?.lastOrder}
              summary={summary}
              onReduceOnlyClose={(c) => void onReduceOnlyClose(c)}
              closeBusy={bindBusy}
              net={net}
            />
          </main>
        ) : null}

        {setupDone && view === "research" ? (
          <ResearchBoard
            coin={researchCoin}
            hypothesis={hypothesis}
            setHypothesis={setHypothesis}
            researchBusy={researchBusy}
            researchStage={researchStage}
            researchElapsed={researchElapsed}
            researchRoles={researchRoles}
            pollMiss={pollMiss}
            researchJobId={researchJobId}
            researchKind={researchKind}
            researchNote={researchNote}
            researchStop={researchStop}
            preview={preview}
            previewHash={previewHash}
            authTyped={authTyped}
            setAuthTyped={setAuthTyped}
            authBusy={authBusy}
            authErr={authErr}
            lastOid={lastOid}
            status={status}
            sessionAlive={sessionAlive}
            checks={checks}
            techOpen={techOpen}
            setTechOpen={setTechOpen}
            researchEvidenceText={researchEvidenceText}
            eligible={eligible}
            coins={coins}
            net={net}
            onResearch={(c) => void researchThis(c)}
            onCancel={() => void onCancelResearch()}
            onAuthorize={(e) => void onAuthorize(e)}
            onCancelBound={(e) => void onCancelBound(e)}
            onCheck={() => void onCheck()}
          />
        ) : null}

        {setupDone && view === "activity" ? (
          <main className="page dense">
            <p className="eyebrow">Activity</p>
            <h1>Evidence trail</h1>
            <p className="lead">Market scan, research, committee, policy, preview hash, execution, OID, fill, receipt. Historical fills never appear inside a new preview.</p>
            <p className={demoReplay ? "err" : "fine"} role="status">
              {demoReplay ? "DEMO REHEARSAL — recorded real evidence, not live." : "Live desk. Replay is labeled separately."}
            </p>
            <button
              type="button"
              className="linkish"
              onClick={() => {
                const next = demoReplay ? "live" : "replay";
                void setDemoMode(next).then((r) => setDemoReplay(r.mode === "replay"));
              }}
            >
              {demoReplay ? "Exit rehearsal" : "Replay recorded evidence"}
            </button>
            <ProofTimeline />
            <ActivityTimeline events={activity} lastOid={lastOid} lastOrder={status?.lastOrder} />
          </main>
        ) : null}

        {setupDone && view === "automation" ? (
          <AutomationCenter
            mission={mission}
            prefs={autoPrefs}
            busy={bindBusy}
            kill={Boolean(status?.kill)}
            openConfirm={openConfirm}
            confirmHours={confirmHours}
            onConfirmConsumed={() => setOpenConfirm(false)}
            onMode={(mode) => {
              void postMission({ mode }).then(setMission);
            }}
            onEnable={async (typed, hours) => {
              setBindBusy(true);
              const r = await postMission({ typed, hours, mode: "guarded" });
              setMission(r);
              setBindBusy(false);
              return r;
            }}
            onStop={async () => {
              setBindBusy(true);
              const r = await postMission({ stop: true });
              setMission(r);
              setBindBusy(false);
              return r;
            }}
            onSavePrefs={(p) => {
              setAutoPrefs(p);
              void saveAutomation(p).then(setAutoPrefs);
            }}
            onOpenHistory={() => setView("activity")}
            net={net}
            wallet={status?.wallet}
          />
        ) : null}

        {setupDone && view === "health" ? (
          <StrategyHealth copy={calibCopy} n={calibN} need={calibNeed} enough={calibEnough} skills={calibSkills} />
        ) : null}

        {setupDone && view === "security" ? (
          <SecurityCenter
            domains={securityDomains}
            checks={checks}
            items={items}
            net={net}
            status={status}
            agent={agent}
            sessionAlive={sessionAlive}
            approved={Boolean(checks.find((c) => c.name === "hl_agent" && c.ok))}
            tradingCapital={summary.accountValue}
            summary={summary}
            policy={hostPolicy}
            pinned={pinned}
            policyHash={policyHash}
            consequences={policyConsequences}
            allowed={policyAllowed}
            refused={policyRefused}
            identityNote={identityNote}
            calibCopy={calibCopy}
            updateNote={updateNote}
            restartAllowed={restartAllowed}
            busy={bindBusy}
            onSession={() => void onSession()}
            onPolicyPreview={(p) => {
              void previewPolicy(p).then((r) => {
                if (r.consequences) setPolicyConsequences(r.consequences);
                if (r.allowed) setPolicyAllowed(r.allowed);
                if (r.refused) setPolicyRefused(r.refused);
              });
            }}
            onPolicyPin={(p) => void onPolicy(p)}
            onCheck={() => void onCheck()}
            onRevoke={() => void onRevoke()}
            onKill={(on) => void onKill(on)}
            onForget={() => {
              void forgetMemory().then(() => setMemoryEpoch((n) => n + 1));
            }}
            code={code}
            expires={expires}
            companionUp={companionUp}
            paired={Boolean(status?.paired)}
            pairingDevices={status?.pairingDevices}
            onRotatePair={() => void onRotatePair()}
          />
        ) : null}
            </div>
            )}
          </div>
        )}
      </div>
      </div>
    </div>
  );
}
