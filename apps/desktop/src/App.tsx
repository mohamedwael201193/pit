import { useEffect, useState } from "react";
import { AuthorizeGate } from "./AuthorizeGate";
import { NAMED } from "./namedStates";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PermissionsCard } from "./Permissions";
import { doctor, localStatus, pairCode, prettyCode, type DoctorCheck } from "./companion";
import { explainStop } from "./explain";

type Net = "mainnet" | "testnet";
type View = "home" | "watch" | "opportunities" | "security" | "settings";

type Coin = { coin: string; reason: string; mark: number; eligible?: boolean };

const HEALTH = "https://pit-health.onrender.com";

const PIPELINE = [
  "MARKET",
  "PRIVATE",
  "SEALED",
  "RESEARCH",
  "CHALLENGE",
  "RISK",
  "POLICY",
  "AUTHORIZE",
  "EXECUTE",
  "PROVE",
  "LEARN",
] as const;

export function App() {
  const [view, setView] = useState<View>("home");
  const [net, setNet] = useState<Net>("mainnet");
  const [sessionAlive, setSessionAlive] = useState(false);
  const [agent, setAgent] = useState("");
  const [code, setCode] = useState("");
  const [expires, setExpires] = useState("");
  const [companionUp, setCompanionUp] = useState(false);
  const [checks, setChecks] = useState<DoctorCheck[]>([]);
  const [coins, setCoins] = useState<Coin[]>([]);
  const [researchStop, setResearchStop] = useState<string | null>(null);
  const [techOpen, setTechOpen] = useState(false);
  const [ticks, setTicks] = useState(0);

  useEffect(() => {
    let gone = false;
    const tick = () => {
      setTicks((n) => n + 1);
      localStatus()
        .then((s) => {
          if (gone) return;
          setCompanionUp(Boolean(s));
          setSessionAlive(Boolean(s?.sessionAlive));
          setAgent(s?.agent || "");
          if (s?.network === "testnet" || s?.network === "mainnet") setNet(s.network);
        })
        .catch(() => {
          if (!gone) setCompanionUp(false);
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

  const lit = !companionUp ? "" : !sessionAlive ? "AUTHORIZE" : "MARKET";
  const stage = !companionUp ? "LOCAL" : !sessionAlive ? "SESSION" : "WATCH";
  const explained = explainStop(researchStop);
  const companionStuck = !companionUp && ticks >= 5;
  const doctorOk = checks.filter((c) => c.ok).length;
  const pairDisplay = code ? prettyCode(code) : companionUp ? "rotating…" : "waiting for local PIT";
  function researchThis() {
    if (!companionUp) {
      setResearchStop("companion_down");
      setView("home");
      return;
    }
    const sealer = checks.find((c) => c.name === "direct_sealer");
    if (sealer && !sealer.ok) {
      setResearchStop("sealer_not_wired");
      setView("home");
      return;
    }
    setResearchStop("direct_token_required");
    setView("home");
  }

  return (
    <div className="shell">
      <header className="top">
        <div>
          <div className="word">PIT.</div>
          <p className="kicker">0.1.1 · local execution authority</p>
        </div>
        <nav className="nav">
          {(["home", "watch", "opportunities", "security", "settings"] as const).map((id) => (
            <button key={id} type="button" className={view === id ? "on" : ""} onClick={() => setView(id)}>
              {id}
            </button>
          ))}
        </nav>
        <span className="kicker">{companionUp ? "companion live" : "starting companion"}</span>
      </header>

      {view === "home" ? (
        <main className="split">
          <section className="left">
            <p className="eyebrow">YOUR MACHINE</p>
            <h1>Authorize on this computer. Never in the browser.</h1>
            <p className="lead">{NAMED.SEED_FORBIDDEN}</p>
            <p className="lead">Your wallet stays yours. Your trading session cannot withdraw.</p>

            <NetworkToggle net={net} onChange={setNet} />
            <NetworkBanner net={net} />

            <article className="card pair-card">
              <p className="label">PAIR THIS COMPUTER</p>
              <p className="pair-code" aria-label="pairing code">
                {pairDisplay}
              </p>
              <p className="fine">
                Type this code at https://pit0g.vercel.app/pair. It expires in two minutes and works once. The website
                never receives a session key.
              </p>
              {expires ? <p className="fine">Expires {expires}</p> : null}
            </article>

            <article className="card">
              <p className="label">STATUS</p>
              <p>
                {companionUp ? "Local PIT is running." : "PIT is starting the local companion on this computer."}{" "}
                {sessionAlive ? "A live order/cancel session exists." : "No live session yet."}
              </p>
              {agent ? <p className="fine">Agent {agent}</p> : null}
              {checks.length > 0 ? (
                <p className="fine">
                  Doctor {doctorOk}/{checks.length} checks ok.
                </p>
              ) : null}
            </article>

            {companionStuck && !explained ? (
              <article className="card stop" role="status">
                <p className="label">LOCAL COMPANION</p>
                <h2>PIT is waiting for the process on this computer</h2>
                <p>
                  No order was placed. Reinstall so the local companion is included, or run <code>pit companion</code>{" "}
                  from a terminal.
                </p>
              </article>
            ) : null}

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
          </section>
          <section className="right">
            <p className="eyebrow">PIPELINE · {stage}</p>
            <ol className="pipe">
              {PIPELINE.map((s) => (
                <li key={s} className={s === lit ? "lit" : ""}>
                  {s}
                </li>
              ))}
            </ol>
            <p className="fine">Revoke the Hyperliquid agent from your account if this machine is lost.</p>
          </section>
        </main>
      ) : null}

      {view === "watch" ? (
        <main className="page">
          <p className="eyebrow">LIVE BOOKS</p>
          <h1>Watch</h1>
          <p className="lead">Public Hyperliquid marks only. This window cannot place an order.</p>
          {coins.length === 0 ? (
            <p className="fine">Empty Watch is the honest state until live books arrive.</p>
          ) : (
            <ul className="watch-grid">
              {coins.map((c) => (
                <li key={c.coin} className="card">
                  <p className="label">{c.coin}</p>
                  <p className="mark-num">{c.mark}</p>
                  <p>Policy {c.eligible ? "PASS" : "BLOCKED"}</p>
                  <p className="fine">{c.reason}</p>
                  <p className="fine">Confidence NOT ENOUGH DATA. Side is not decided on this surface.</p>
                </li>
              ))}
            </ul>
          )}
        </main>
      ) : null}

      {view === "opportunities" ? (
        <main className="page">
          <p className="eyebrow">POLICY FILTER</p>
          <h1>Opportunities</h1>
          <p className="lead">Live public marks. Research stays sealed on this machine. Nothing here places an order.</p>
          {coins.length === 0 ? (
            <p className="fine">Empty is honest until Hyperliquid books arrive.</p>
          ) : (
            <ul className="watch-grid">
              {coins.map((c) => (
                <li key={c.coin} className="card">
                  <p className="label">{c.coin}</p>
                  <p className="mark-num">{c.mark}</p>
                  <p>Policy {c.eligible ? "PASS" : "BLOCKED"}</p>
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

      {view === "security" ? (
        <main className="page">
          <p className="eyebrow">PERMISSIONS</p>
          <h1>Security</h1>
          <p className="lead">Order and cancel only. Withdraw is impossible through PIT.</p>
          <AuthorizeGate sessionAlive={sessionAlive} />
          <PermissionsCard />
          <article className="card">
            <p className="label">REVOKE</p>
            <p>Kill the local session, then remove the PIT agent from your Hyperliquid account.</p>
          </article>
        </main>
      ) : null}

      {view === "settings" ? (
        <main className="page">
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
  );
}
