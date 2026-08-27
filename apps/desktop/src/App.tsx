import { useEffect, useMemo, useState } from "react";
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
  doctor,
  localStatus,
  pairCode,
  prettyCode,
  wakeCompanion,
  type DoctorCheck,
  type LocalStatus,
} from "./companion";
import { explainStop } from "./explain";
import { probes, type Probe } from "./readiness";

type Net = "mainnet" | "testnet";
type View = "home" | "watch" | "research" | "activity" | "policy" | "security" | "account" | "settings";

type Coin = { coin: string; reason: string; mark: number; eligible?: boolean };

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
  onDone,
}: {
  step: number;
  setStep: (n: number) => void;
  net: Net;
  setNet: (n: Net) => void;
  items: Probe[];
  code: string;
  expires: string;
  companionUp: boolean;
  onDone: () => void;
}) {
  const last = 8;
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
          <h1>Connect your wallet in the browser.</h1>
          <p className="lead">Open pit0g.vercel.app after pairing. Approve nothing that looks like a seed or a withdraw permission.</p>
        </>
      ) : null}
      {step === 4 ? (
        <>
          <h1>Pick one world and stay there.</h1>
          <NetworkToggle net={net} onChange={setNet} />
          <NetworkBanner net={net} />
        </>
      ) : null}
      {step === 5 ? (
        <>
          <h1>Hyperliquid is order and cancel only.</h1>
          <p className="lead">order ✓ cancel ✓ withdraw ✗ leverage ✗ transfer ✗</p>
          <PermissionsCard />
          <SessionNote />
        </>
      ) : null}
      {step === 6 ? (
        <>
          <h1>Pin a policy before research.</h1>
          <p className="lead">These are example defaults until your workspace pins a file. The model cannot raise them.</p>
          <PolicyLaw />
        </>
      ) : null}
      {step === 7 ? (
        <>
          <h1>Security check</h1>
          <p className="lead">Each row is a live probe. Waiting is honest. Green is never invented.</p>
          <ProbeList items={items} />
        </>
      ) : null}
      {step === 8 ? (
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
  const [techOpen, setTechOpen] = useState(false);
  const [ticks, setTicks] = useState(0);
  const [setupStep, setSetupStep] = useState(0);
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
          if (!gone) setChecks(c);
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
    tick();
    const id = window.setInterval(tick, 4000);
    return () => {
      gone = true;
      window.clearInterval(id);
    };
  }, [net]);

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
  const items = useMemo(() => probes(checks, status, companionUp), [checks, status, companionUp]);
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

  function researchThis() {
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
    setResearchStop("direct_token_required");
    setView("research");
  }

  return (
    <div className="app">
      <aside className="rail">
        <div className="rail-brand">
          <div className="word">PIT.</div>
          <p className="kicker">0.1.2 · local execution</p>
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
          <p>PIT 0.1.2</p>
          <button type="button" className="ghost" onClick={() => setView("settings")}>
            Help / Diagnostics
          </button>
        </div>
      </aside>

      <div className="stage">
        <header className="bar">
          <div>
            <p className="eyebrow">WORKSPACE</p>
            <p>{status?.workspace || walletCheck?.detail || "unbound"}</p>
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
            onDone={finishSetup}
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
                        <button type="button" className="linkish" onClick={researchThis}>
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
            {explained ? (
              <article className="card stop" role="alert">
                <p className="label">STOPPED</p>
                <h2>{explained.title}</h2>
                <p>{explained.body}</p>
                <button type="button" className="linkish" onClick={() => setTechOpen((v) => !v)}>
                  {techOpen ? "Hide technical evidence" : "View technical evidence"}
                </button>
                {techOpen ? (
                  <p className="fine">
                    Code {researchStop}. Verification is fail-closed. Router fallback is impossible.
                  </p>
                ) : null}
                <button type="button" onClick={() => setResearchStop(null)}>
                  Retry
                </button>
              </article>
            ) : (
              <p className="fine">Private research has not been run on this machine in this session.</p>
            )}
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
                    <button type="button" className="linkish" onClick={researchThis}>
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
            <PolicyLaw />
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
