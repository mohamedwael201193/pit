import { useEffect, useMemo, useRef, useState } from "react";
import { AuthorizeGate } from "./AuthorizeGate";
import { EmptyHome } from "./EmptyHome";
import { NAMED } from "./namedStates";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PermissionsCard } from "./Permissions";
import { PolicyLaw } from "./PolicyLaw";
import { PreviewNote } from "./PreviewNote";
import { SessionNote } from "./SessionNote";
import {
  bindWallet,
  createLocalSession,
  cancelResearch,
  describeBindError,
  doctor,
  localStatus,
  pairCode,
  pinLocalPolicy,
  prettyCode,
  researchStatus,
  revokeLocalSession,
  startResearch,
  wakeCompanion,
  type DoctorCheck,
  type LocalStatus,
} from "./companion";
import { explainStop } from "./explain";
import { probes, type Probe } from "./readiness";

type Net = "mainnet" | "testnet";
type View = "home" | "watch" | "research" | "activity" | "policy" | "security" | "account" | "settings";

type Coin = { coin: string; reason: string; mark: number; eligible?: boolean };

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
] as const;

function sleep(ms: number) {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

function ResearchProgress({
  stage,
  elapsedMs,
  coin,
  onCancel,
}: {
  stage: string;
  elapsedMs: number;
  coin: string;
  onCancel: () => void;
}) {
  const current = RESEARCH_STAGES.indexOf(stage as (typeof RESEARCH_STAGES)[number]);
  return (
    <article className="card" role="status">
      <p className="label">LIVE SEALED REQUEST</p>
      <h2>{stage.replaceAll("_", " ")}</h2>
      <p>
        {coin || "ETH"} · {(elapsedMs / 1000).toFixed(1)}s elapsed. This is a live Direct round-trip, not a timer.
      </p>
      <ol className="pipe stages">
        {RESEARCH_STAGES.map((name, i) => (
          <li key={name} className={i === current ? "lit" : i < current ? "done" : ""}>
            {name.replaceAll("_", " ")}
          </li>
        ))}
      </ol>
      <button type="button" className="linkish" onClick={onCancel}>
        Cancel
      </button>
    </article>
  );
}

const HEALTH = "https://pit-health.onrender.com";
const SETUP_KEY = "pit.desk.setup";
const PAIR_URL = "https://pit0g.vercel.app/pair";

const RAIL: { id: View; label: string }[] = [
  { id: "home", label: "Home" },
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
        Type this code at {PAIR_URL}. It expires in two minutes and works once. The website never receives a session
        key.
      </p>
      {expires ? <p className="fine">Expires {expires}</p> : null}
      <button type="button" className="linkish" onClick={() => void copy()} disabled={!code}>
        {copied ? "Copied" : "Copy code"}
      </button>
    </article>
  );
}

function ProbeList({ items }: { items: Probe[] }) {
  return (
    <ul className="probes">
      {items.map((p) => (
        <li key={p.id} className={markProbe(p)}>
          <strong>{p.state === "ok" ? "ready" : p.state === "fail" ? "fail" : "waiting"}</strong>
          <span>{p.label}</span>
          <em>{p.detail}</em>
        </li>
      ))}
    </ul>
  );
}

function Setup({
  step,
  setStep,
  net,
  setNet,
  items,
  code,
  expires,
  companionUp,
  walletDraft,
  setWalletDraft,
  bindError,
  bindBusy,
  boundWallet,
  pinned,
  agent,
  sessionAlive,
  onBind,
  onSession,
  onPolicy,
  onResearch,
  onCancelResearch,
  onDone,
  checks,
  researchBusy,
  researchVerified,
  researchStage,
  researchElapsed,
}: {
  step: number;
  setStep: (n: number) => void;
  net: Net;
  setNet: (n: Net) => void;
  items: Probe[];
  code: string;
  expires: string;
  companionUp: boolean;
  walletDraft: string;
  setWalletDraft: (v: string) => void;
  bindError: string | null;
  bindBusy: boolean;
  boundWallet: string;
  pinned: boolean;
  agent: string;
  sessionAlive: boolean;
  onBind: () => void;
  onSession: () => void;
  onPolicy: () => void;
  onResearch: () => void;
  onCancelResearch: () => void;
  onDone: () => void;
  checks: DoctorCheck[];
  researchBusy: boolean;
  researchVerified: boolean;
  researchStage: string;
  researchElapsed: number;
}) {
  const directOk = Boolean(checks.find((c) => c.name === "direct_auth" && c.ok));
  const directDetail = checks.find((c) => c.name === "direct_auth")?.detail;
  const hlAgent = checks.find((c) => c.name === "hl_agent");
  const last = 9;
  return (
    <section className="setup">
      <p className="eyebrow">FIRST RUN · {step + 1} / {last + 1}</p>
      {step === 0 ? (
        <>
          <h1>Your private trading desk.</h1>
          <p className="lead">PIT researches in a sealed enclave, then waits for you. Nothing leaves this machine without an exact authorize.</p>
        </>
      ) : null}
      {step === 1 ? (
        <>
          <h1>Your wallet stays yours.</h1>
          <p className="lead">{NAMED.SEED_FORBIDDEN}</p>
          <p className="lead">The browser never holds a Hyperliquid session key, a memory key, or a Direct token.</p>
        </>
      ) : null}
      {step === 2 ? (
        <>
          <h1>Pair the browser to this machine.</h1>
          <p className="lead">Launch is local. The one-time code never includes a secret.</p>
          <PairingBlock code={code} expires={expires} companionUp={companionUp} />
        </>
      ) : null}
      {step === 3 ? (
        <>
          <h1>Pick one world and stay there.</h1>
          <NetworkToggle net={net} onChange={setNet} />
          <NetworkBanner net={net} />
        </>
      ) : null}
      {step === 4 ? (
        <>
          <h1>Connect your wallet.</h1>
          <p className="lead">
            Paste the public 0x address on this computer, or pair the browser and connect with Privy. PIT never asks for
            a seed or a private key.
          </p>
          <form
            className="bind-form card"
            onSubmit={(e) => {
              e.preventDefault();
              onBind();
            }}
          >
            <label htmlFor="desk-wallet">Public wallet</label>
            <input
              id="desk-wallet"
              autoComplete="off"
              spellCheck={false}
              placeholder="0x…"
              value={walletDraft}
              onChange={(e) => setWalletDraft(e.target.value)}
            />
            <button type="submit" className="linkish" disabled={bindBusy || !companionUp}>
              {boundWallet ? "Wallet bound on this computer" : "Bind this computer"}
            </button>
            {boundWallet ? <p className="fine">Bound {boundWallet}</p> : null}
            {bindError ? (
              <p className="err" role="alert">
                {bindError}
              </p>
            ) : null}
          </form>
        </>
      ) : null}
      {step === 5 ? (
        <>
          <h1>Protect my strategy.</h1>
          <p className="lead">
            PIT sends your private strategy only through 0G’s verified sealed path. Sign in the paired browser. This
            computer stores the token. The website never receives it.
          </p>
          <p className="fine">This is not a withdraw. It cannot place a Hyperliquid order. It lasts 24 hours.</p>
          <p className="fine">{directOk ? "Sealed-path signature is on this computer." : directDetail || "Waiting for the paired-browser signature."}</p>
          <a className="linkish" href={PAIR_URL.replace("/pair", "/app")} target="_blank" rel="noreferrer">
            Open paired site
          </a>
        </>
      ) : null}
      {step === 6 ? (
        <>
          <h1>Hyperliquid is order and cancel only.</h1>
          <p className="lead">order ✓ cancel ✓ withdraw ✗ leverage ✗ transfer ✗</p>
          <PermissionsCard />
          <SessionNote />
          <button type="button" className="linkish" onClick={onSession} disabled={bindBusy || !companionUp || !boundWallet}>
            {sessionAlive ? "Local session is live" : "Create local session"}
          </button>
          {agent ? <p className="fine">Agent {agent}. Approve this agent on Hyperliquid. PIT cannot withdraw.</p> : null}
          {hlAgent ? <p className="fine">{hlAgent.ok ? "extraAgents lists this session." : hlAgent.detail}</p> : null}
          {bindError ? (
            <p className="err" role="alert">
              {bindError}
            </p>
          ) : null}
        </>
      ) : null}
      {step === 7 ? (
        <>
          <h1>Pin a policy before research.</h1>
          <p className="lead">The model cannot raise clip, leverage, or permissions. Pin writes a hash file on this computer.</p>
          <PolicyLaw pinned={pinned} onPin={onPolicy} busy={bindBusy || !boundWallet} />
        </>
      ) : null}
      {step === 8 ? (
        <>
          <h1>Security check, then a real research test.</h1>
          <p className="lead">Each row is a live probe. Waiting is honest. Green is never invented. Research is a live sealed Direct request.</p>
          <ProbeList items={items} />
          <button type="button" className="linkish" onClick={onResearch} disabled={researchBusy || !companionUp}>
            {researchBusy ? "Research running…" : "Run a real research test"}
          </button>
          {researchBusy ? (
            <ResearchProgress stage={researchStage} elapsedMs={researchElapsed} coin="ETH" onCancel={onCancelResearch} />
          ) : null}
          {researchVerified ? <p className="fine">RESEARCH VERIFIED. VerifyE2EE matched the on-chain teeSigner.</p> : null}
        </>
      ) : null}
      {step === 9 ? (
        <>
          <h1>Ready when the probes are real.</h1>
          <p className="lead">Watch is live public marks. Private research stays sealed. Authorize stays on this computer.</p>
          <ProbeList items={items} />
        </>
      ) : null}
      <div className="row">
        {step > 0 ? (
          <button type="button" className="off" onClick={() => setStep(step - 1)}>
            Back
          </button>
        ) : null}
        {step < last ? (
          <button type="button" className="on" onClick={() => setStep(step + 1)}>
            Continue
          </button>
        ) : (
          <button type="button" className="on" onClick={onDone}>
            Open the desk
          </button>
        )}
      </div>
    </section>
  );
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
  const [researchRoles, setResearchRoles] = useState<Array<{ role?: string; verify_e2ee?: string; pubkey_signer?: string }>>([]);
  const [researchBusy, setResearchBusy] = useState(false);
  const [researchStage, setResearchStage] = useState("READING_MARKET");
  const [researchElapsed, setResearchElapsed] = useState(0);
  const [researchCoin, setResearchCoin] = useState("ETH");
  const [researchEvidence, setResearchEvidence] = useState<string>("");
  const researchGen = useRef(0);
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
    void wakeCompanion();
  }, []);

  useEffect(() => {
    let gone = false;
    let timer = 0;
    const tick = () => {
      setTicks((n) => n + 1);
      localStatus()
        .then((s) => {
          if (gone) return;
          setCompanionUp(Boolean(s));
          setStatus(s);
          setSessionAlive(Boolean(s?.sessionAlive));
          setAgent(s?.agent || "");
          if (s?.network === "testnet" || s?.network === "mainnet") setNet(s.network);
          if (s?.wallet) setWalletDraft((cur) => cur || s.wallet || "");
        })
        .catch(() => {
          if (!gone) {
            setCompanionUp(false);
            setStatus(null);
          }
        });
      pairCode()
        .then((p) => {
          if (gone) return;
          setCode(p?.code || "");
          setExpires(p?.expires || "");
        })
        .catch(() => undefined);
      doctor()
        .then((c) => {
          if (!gone) {
            setChecks(c);
            setPinned(Boolean(c.find((x) => x.name === "policy" && x.ok)));
          }
        })
        .catch(() => undefined);
      fetch(`${HEALTH}/watch?network=${net}`)
        .then((r) => r.json() as Promise<{ coins?: Coin[]; sign?: boolean; trade?: boolean }>)
        .then((body) => {
          if (gone || body.sign || body.trade) return;
          setCoins(Array.isArray(body.coins) ? body.coins : []);
        })
        .catch(() => {
          if (!gone) setCoins([]);
        });
    };
    const loop = () => {
      tick();
      timer = window.setTimeout(loop, companionUp ? 4000 : 800);
    };
    loop();
    return () => {
      gone = true;
      window.clearTimeout(timer);
    };
  }, [net, companionUp]);

  useEffect(() => {
    if (!sessionAlive) return;
    try {
      window.localStorage.setItem(SETUP_KEY, "1");
    } catch {
      /* ignore */
    }
    setSetupDone(true);
  }, [sessionAlive]);

  const explained = explainStop(researchStop);
  const companionStuck = !companionUp && ticks >= 5;
  const items = useMemo(
    () => probes(checks, status, companionUp, researchRoles.some((r) => String(r.verify_e2ee).toUpperCase() === "OK")),
    [checks, status, companionUp, researchRoles],
  );
  const eligible = coins.filter((c) => c.eligible);
  const walletCheck = checks.find((c) => c.name === "wallet");

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
    setSessionAlive(true);
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

  async function researchThis(coin?: string) {
    const gen = ++researchGen.current;
    const want = (coin || "ETH").toUpperCase();
    setResearchNote(null);
    setResearchRoles([]);
    setResearchEvidence("");
    setResearchCoin(want);
    if (!companionUp) {
      setResearchStop("companion_down");
      setView("research");
      return;
    }
    const sealer = checks.find((c) => c.name === "direct_sealer");
    if (sealer && !sealer.ok) {
      setResearchStop("sealer_not_wired");
      setView("research");
      return;
    }
    const auth = checks.find((c) => c.name === "direct_auth");
    if (auth && !auth.ok) {
      setResearchStop("direct_token_required");
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
      const started = await startResearch(want);
      if (gen !== researchGen.current) return;
      if (started.error && !started.running) {
        setResearchStop(started.error);
        return;
      }
      if (started.stage) setResearchStage(started.stage);
      for (;;) {
        await sleep(400);
        if (gen !== researchGen.current) return;
        const st = await researchStatus();
        if (st.stage) setResearchStage(st.stage);
        if (typeof st.elapsed_ms === "number") setResearchElapsed(st.elapsed_ms);
        if (st.evidence) setResearchEvidence(JSON.stringify(st.evidence, null, 2));
        if (st.running) continue;
        if (st.error) {
          setResearchStop(st.error);
          return;
        }
        const roles = Array.isArray(st.roles) ? st.roles : [];
        const verified = roles.length > 0 && roles.every((x) => String(x.verify_e2ee).toUpperCase() === "OK");
        if (!verified) {
          setResearchStop("TEE_VERIFY_FAIL");
          return;
        }
        setResearchRoles(roles);
        setResearchNote(st.note || "Sealed committee verified on this computer.");
        return;
      }
    } catch (e) {
      if (gen !== researchGen.current) return;
      const msg = e instanceof Error ? e.message : "companion_http";
      setResearchStop(msg || "companion_http");
    } finally {
      window.clearInterval(tick);
      if (gen === researchGen.current) setResearchBusy(false);
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

  return (
    <div className="app">
      <aside className="rail">
        <div className="rail-brand">
          <div className="word">PIT.</div>
          <p className="kicker">0.1.5 · local execution</p>
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
          <p>PIT 0.1.5</p>
          <button type="button" className="ghost" onClick={() => setView("settings")}>
            Help / Diagnostics
          </button>
        </div>
      </aside>

      <div className="stage">
        <header className="bar">
          <div>
            <p className="eyebrow">WORKSPACE</p>
            <p>{status?.workspace || status?.wallet || walletCheck?.detail || "unbound"}</p>
          </div>
          <NetworkToggle net={net} onChange={setNet} />
          <div className="bar-meta">
            <p className="pair-chip">{code ? prettyCode(code) : companionUp ? "code rotating" : "starting companion"}</p>
            <p>Wallet {walletCheck?.ok ? "bound" : "waiting"}</p>
            <p>Session {sessionAlive ? "order/cancel live" : "none"}</p>
            {agent ? <p className="fine">Agent {agent}</p> : null}
          </div>
        </header>

        {!setupDone ? (
          <Setup
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
            sessionAlive={sessionAlive}
            onBind={() => void onBind()}
            onSession={() => void onSession()}
            onPolicy={() => void onPolicy()}
            onResearch={() => void researchThis()}
            onCancelResearch={() => void onCancelResearch()}
            onDone={finishSetup}
            checks={checks}
            researchBusy={researchBusy}
            researchVerified={Boolean(researchNote && !explained)}
            researchStage={researchStage}
            researchElapsed={researchElapsed}
          />
        ) : null}

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
                <ProbeList items={items} />
              </article>
            </div>
            <EmptyHome />
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
                    <th>Policy</th>
                    <th>Why</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  {coins.map((c) => (
                    <tr key={c.coin}>
                      <td>{c.coin}</td>
                      <td className="mark-num">{c.mark}</td>
                      <td>{c.eligible ? "PASS" : "BLOCKED"}</td>
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
              never places the order.
            </p>
            {researchBusy ? (
              <ResearchProgress
                stage={researchStage}
                elapsedMs={researchElapsed}
                coin={researchCoin}
                onCancel={() => void onCancelResearch()}
              />
            ) : null}
            {researchNote && !explained && !researchBusy ? (
              <article className="card">
                <p className="label">RESEARCH VERIFIED</p>
                <p>{researchNote}</p>
                <ul className="doctor">
                  {researchRoles.map((role) => (
                    <li key={role.role}>
                      <strong>{role.verify_e2ee}</strong> {role.role} {role.pubkey_signer}
                    </li>
                  ))}
                </ul>
              </article>
            ) : null}
            {explained && !researchBusy ? (
              <article className="card stop" role="alert">
                <p className="label">RESEARCH STOPPED</p>
                <h2>{explained.title}</h2>
                <p>{explained.body}</p>
                <button type="button" className="linkish" onClick={() => setTechOpen((v) => !v)}>
                  {techOpen ? "Hide technical evidence" : "View technical evidence"}
                </button>
                {techOpen ? (
                  <pre className="pipe evidence">
                    Code {researchStop}
                    {"\n"}
                    {researchEvidence || "Verification is fail-closed. Router fallback is impossible."}
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
            <article className="card">
              <p className="label">THIS MACHINE</p>
              <p>No order id is shown until Hyperliquid accepts one after you type AUTHORIZE.</p>
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
            <AuthorizeGate sessionAlive={sessionAlive} />
            <PermissionsCard />
            <SessionNote />
            <article className="card">
              <p className="label">REVOKE</p>
              <p>Kill the local session, then remove the PIT agent from your Hyperliquid account.</p>
              <button type="button" className="linkish" onClick={() => void onRevoke()} disabled={bindBusy || !sessionAlive}>
                Revoke local session
              </button>
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
              <p>Session {sessionAlive ? "alive" : "none"}</p>
            </article>
            <p className="fine">{NAMED.TRANSFER_NOT_LIVE}</p>
          </main>
        ) : null}

        {setupDone && view === "settings" ? (
          <main className="page dense">
            <p className="eyebrow">DIAGNOSTICS</p>
            <h1>Settings</h1>
            <NetworkToggle net={net} onChange={setNet} />
            <NetworkBanner net={net} />
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
            <p className="fine">{NAMED.TWO_WALLETS}</p>
            <p className="fine">{NAMED.TRANSFER_NOT_LIVE}</p>
          </main>
        ) : null}
      </div>
    </div>
  );
}
