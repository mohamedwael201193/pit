import { FormEvent } from "react";
import { NAMED } from "./namedStates";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PermissionsCard } from "./Permissions";
import { PolicyLaw } from "./PolicyLaw";
import { SessionNote } from "./SessionNote";
import { LINKS, hyperliquidAPI, hyperliquidApp } from "./links";
import { prettyCode, type DoctorCheck } from "./companion";
import { ComputeCard } from "./ComputeCard";
import type { Probe } from "./readiness";

type Net = "mainnet" | "testnet";
type View = "home" | "watch" | "research" | "activity" | "policy" | "security" | "account" | "settings";
type ResearchRole = { role?: string; verify_e2ee?: string; proposed_side?: string; survives?: boolean; kill?: boolean; pubkey_signer?: string };

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
  return (
    <article className="card pair-card">
      <p className="label">PAIR THIS COMPUTER</p>
      <p className="pair-code" aria-label="pairing code">
        {display}
      </p>
      <p className="fine">
        Type this code on the pairing page. It expires in two minutes and works once. The website never receives a session key.
      </p>
      {expires ? <p className="fine">Expires {expires}</p> : null}
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
          {p.state !== "ok" && (p.id === "hyperliquid" || p.id === "hl_agent" || p.id === "session") ? (
            <span className="cta-row">
              <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
                Open Hyperliquid API
              </a>
            </span>
          ) : null}
          {p.state !== "ok" && p.id === "direct" ? (
            <a className="linkish" href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
              Open 0G Private Compute
            </a>
          ) : null}
          {p.state !== "ok" && (p.id === "wallet" || p.id === "local") ? (
            <a className="linkish" href={LINKS.pair} target="_blank" rel="noreferrer">
              Open pairing
            </a>
          ) : null}
          {p.state !== "ok" && p.id === "tee" && onGo ? (
            <button type="button" className="linkish" onClick={() => onGo("research")}>
              Open Research
            </button>
          ) : null}
        </li>
      ))}
    </ul>
  );
}

export function SetupWizard({
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
  agentName,
  sessionAlive,
  onBind,
  onSession,
  onPolicy,
  onResearch,
  onCheck,
  onDone,
  checks,
  researchBusy,
  researchVerified,
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
  agentName?: string;
  sessionAlive: boolean;
  onBind: () => void;
  onSession: () => void;
  onPolicy: () => void;
  onResearch: () => void;
  onCheck: () => void;
  onDone: () => void;
  checks: DoctorCheck[];
  researchBusy: boolean;
  researchVerified: boolean;
}) {
  const hlAgent = checks.find((c) => c.name === "hl_agent");
  const last = 12;
  const titles = [
    "Welcome",
    "Choose network",
    "Pair this computer",
    "Connect wallet",
    "Protect my strategy",
    "Private compute",
    "Policy",
    "Connect Hyperliquid",
    "Secure session",
    "Approve PIT",
    "System check",
    "Optional first research",
    "Ready",
  ];
  return (
    <section className="setup">
      <p className="eyebrow">FIRST RUN · {step + 1} / {last + 1} · {titles[step]}</p>
      {step === 0 ? (
        <>
          <h1>Your private trading desk.</h1>
          <p className="lead">
            You own the money. You write the law. PIT researches in private compute, then waits. It is not a chatbot and it
            will not trade by itself.
          </p>
        </>
      ) : null}
      {step === 1 ? (
        <>
          <h1>Pick one world and stay there.</h1>
          <NetworkToggle net={net} onChange={setNet} />
          <NetworkBanner net={net} />
          <p className="fine">MAINNET is production. TESTNET is the lab. Mixing compute and venue across worlds is refused.</p>
        </>
      ) : null}
      {step === 2 ? (
        <>
          <h1>Pair the browser to this machine.</h1>
          <p className="lead">Launch is local. The one-time code never includes a secret.</p>
          <PairingBlock code={code} expires={expires} companionUp={companionUp} />
          <a className="linkish" href={LINKS.pair} target="_blank" rel="noreferrer">
            Open pairing
          </a>
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 3 ? (
        <>
          <h1>Connect your wallet.</h1>
          <p className="lead">{NAMED.SEED_FORBIDDEN} Paste the public 0x address, or pair and connect in the browser.</p>
          <form
            className="bind-form card"
            onSubmit={(e: FormEvent) => {
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
      {step === 4 ? (
        <>
          <h1>Protect my strategy.</h1>
          <p className="lead">
            Sign in the paired browser. This computer stores the sealed-path authorization for 24 hours. The website never
            receives it. This cannot withdraw and cannot place a Hyperliquid order.
          </p>
          <p>{checks.find((c) => c.name === "direct_auth")?.ok ? "Protected on this computer." : "Waiting for the wallet signature."}</p>
          <a className="linkish" href={LINKS.app} target="_blank" rel="noreferrer">
            Protect my strategy
          </a>
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 5 ? (
        <>
          <h1>Private compute is not trading capital.</h1>
          <ComputeCard checks={checks} onCheck={onCheck} />
        </>
      ) : null}
      {step === 6 ? (
        <>
          <h1>Pin a policy before research.</h1>
          <p className="lead">The model cannot raise clip, leverage, or permissions.</p>
          <PolicyLaw pinned={pinned} onPin={onPolicy} busy={bindBusy || !boundWallet} />
        </>
      ) : null}
      {step === 7 ? (
        <>
          <h1>Connect Hyperliquid.</h1>
          <p className="lead">PIT needs your trading account. Open the official app. PIT still cannot withdraw.</p>
          <a className="linkish" href={hyperliquidApp(net)} target="_blank" rel="noreferrer">
            Open Hyperliquid
          </a>
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 8 ? (
        <>
          <h1>Create a secure session.</h1>
          <PermissionsCard />
          <SessionNote />
          <button type="button" className="linkish" onClick={onSession} disabled={bindBusy || !companionUp || !boundWallet}>
            {sessionAlive ? "Refresh secure session" : "Create secure PIT session"}
          </button>
          {agent ? <p className="fine">PIT Agent {agentName || ""} {agent}. If Hyperliquid still lists this agent, PIT reuses it.</p> : null}
          {bindError ? (
            <p className="err" role="alert">
              {bindError}
            </p>
          ) : null}
        </>
      ) : null}
      {step === 9 ? (
        <>
          <h1>Approve PIT on Hyperliquid.</h1>
          <p className="lead">
            Open Hyperliquid API. Authorize API Wallet with the name and address below. Then Check again. PIT cannot withdraw.
          </p>
          {agentName ? <p>Name {agentName}</p> : null}
          {agent ? <p>PIT Agent {agent}</p> : null}
          <p>{hlAgent?.ok ? "Your trading account is ready." : hlAgent?.detail || "Waiting for Hyperliquid approval."}</p>
          <a className="linkish" href={hyperliquidAPI(net)} target="_blank" rel="noreferrer">
            Approve PIT
          </a>
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 10 ? (
        <>
          <h1>System check.</h1>
          <ProbeList items={items} net={net} />
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 11 ? (
        <>
          <h1>Optional first research.</h1>
          <p className="lead">Not required to finish. This spends private compute, not trading capital. It will not place an order.</p>
          <button type="button" className="linkish" onClick={onResearch} disabled={researchBusy || !companionUp}>
            {researchBusy ? "Research running…" : "Run a real research test"}
          </button>
          {researchVerified ? <p className="fine">RESEARCH VERIFIED. Three roles matched the on-chain TEE signer.</p> : null}
          <button type="button" className="linkish" onClick={() => setStep(step + 1)}>
            Skip
          </button>
        </>
      ) : null}
      {step === 12 ? (
        <>
          <h1>Ready when the probes are real.</h1>
          <p className="lead">Watch is live public marks. Private research stays sealed. Authorize stays on this computer.</p>
          <ProbeList items={items} net={net} />
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
