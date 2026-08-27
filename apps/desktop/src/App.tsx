import { useEffect, useState } from "react";
import { AuthorizeGate } from "./AuthorizeGate";
import { IsolateNote } from "./IsolateNote";
import { KillNote } from "./KillNote";
import { NAMED } from "./namedStates";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PermissionsCard } from "./Permissions";
import { RecoverNote } from "./RecoverNote";
import { doctor, localStatus, pairCode, prettyCode, type DoctorCheck } from "./companion";

type Net = "mainnet" | "testnet";

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
];

export function App() {
  const [net, setNet] = useState<Net>("mainnet");
  const [sessionAlive, setSessionAlive] = useState(false);
  const [agent, setAgent] = useState("");
  const [code, setCode] = useState("");
  const [expires, setExpires] = useState("");
  const [companionUp, setCompanionUp] = useState(false);
  const [checks, setChecks] = useState<DoctorCheck[]>([]);

  useEffect(() => {
    let gone = false;
    const tick = () => {
      localStatus()
        .then((s) => {
          if (gone) return;
          setCompanionUp(Boolean(s));
          setSessionAlive(Boolean(s?.sessionAlive));
          setAgent(s?.agent || "");
          if (s?.network === "testnet" || s?.network === "mainnet") {
            setNet(s.network);
          }
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
    };
    tick();
    const id = window.setInterval(tick, 4000);
    return () => {
      gone = true;
      window.clearInterval(id);
    };
  }, []);

  return (
    <div className="shell">
      <header className="top">
        <div>
          <div className="word">PIT.</div>
          <p className="kicker">0.1.0 · local execution authority</p>
        </div>
        <nav>
          <span className="kicker">{companionUp ? "companion live" : "starting companion"}</span>
        </nav>
      </header>
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
              {code ? prettyCode(code) : "starting…"}
            </p>
            <p className="fine">
              Type this code at https://pit0g.vercel.app/pair. It expires in two minutes and works once. The
              website never receives a session key.
            </p>
            {expires ? <p className="fine">Expires {expires}</p> : null}
          </article>

          <article className="card">
            <p className="label">SESSION</p>
            <p>{sessionAlive ? "Live order/cancel session on this machine." : "No live session. Run pit session after you approve the printed agent."}</p>
            {agent ? <p className="fine">Agent {agent}</p> : null}
          </article>

          <AuthorizeGate sessionAlive={sessionAlive} />
          <PermissionsCard />

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

          <RecoverNote />
          <IsolateNote />
          <KillNote />
          <p className="fine">{NAMED.TWO_WALLETS}</p>
        </section>
        <section className="right">
          <p className="eyebrow">PIPELINE</p>
          <ol className="pipe">
            {PIPELINE.map((s) => (
              <li key={s}>{s}</li>
            ))}
          </ol>
          <p className="fine">{NAMED.TEE_VERIFY_FAIL}</p>
          <p className="fine">Revoke the Hyperliquid agent from your account if this machine is lost.</p>
        </section>
      </main>
    </div>
  );
}
