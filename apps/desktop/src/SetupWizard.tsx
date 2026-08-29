import { FormEvent } from "react";
import { NAMED } from "./namedStates";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PermissionsCard } from "./Permissions";
import { PolicyLaw } from "./PolicyLaw";
import { SessionNote } from "./SessionNote";
import { LINKS, hyperliquidAPI, hyperliquidApp } from "./links";
import { ExternalLink } from "./ExternalLink";
import { prettyCode, type DoctorCheck } from "./companion";
import { ComputeCard } from "./ComputeCard";
import type { Probe } from "./readiness";

type Net = "mainnet" | "testnet";

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
    <section>
      <p className="pair-chip" aria-label="pairing code">
        {display}
      </p>
      <p className="fine">
        Type this code on the pairing page. It expires in two minutes and works once. The website never receives a session key.
      </p>
      {expires ? <p className="fine">Expires {expires}</p> : null}
    </section>
  );
}

function StepMeta({ status, meaning, why }: { status: string; meaning: string; why: string }) {
  return (
    <dl className="setup-meta">
      <div>
        <dt>Status</dt>
        <dd>{status}</dd>
      </div>
      <div>
        <dt>What this means</dt>
        <dd>{meaning}</dd>
      </div>
      <div>
        <dt>Why it is needed</dt>
        <dd>{why}</dd>
      </div>
    </dl>
  );
}

export function SetupWizard({
  step,
  setStep,
  net,
  setNet,
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
  onCheck,
  onDone,
  checks,
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
  const last = 8;
  const titles = [
    "Connect wallet",
    "Choose network",
    "Connect Hyperliquid",
    "Create PIT session",
    "Approve PIT",
    "Protect private research",
    "Private compute",
    "Set policy",
    "Ready to discover",
  ];
  return (
    <section className="setup">
      <p className="eyebrow">
        FIRST RUN · {step + 1} / {last + 1} · {titles[step]}
      </p>
      {step === 0 ? (
        <>
          <h1>Connect your wallet.</h1>
          <StepMeta
            status={boundWallet ? "Bound" : companionUp ? "Waiting for bind" : "Companion starting"}
            meaning="This computer is paired to a public 0x address. The seed never enters PIT."
            why="Every later step — session, policy, research — is bound to this wallet."
          />
          <p className="lead">
            Pair the browser to this machine, then bind the public 0x address. {NAMED.SEED_FORBIDDEN}
          </p>
          <PairingBlock code={code} expires={expires} companionUp={companionUp} />
          <ExternalLink className="primary" href={LINKS.pair}>
            Open official page
          </ExternalLink>
          <form
            className="bind-form"
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
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 1 ? (
        <>
          <h1>Choose Mainnet or Testnet.</h1>
          <NetworkToggle net={net} onChange={setNet} />
          <NetworkBanner net={net} />
          <p className="fine">MAINNET is production. TESTNET is the lab. Mixing compute and venue across worlds is refused.</p>
        </>
      ) : null}
      {step === 2 ? (
        <>
          <h1>Connect Hyperliquid.</h1>
          <p className="lead">Open the official app for this network. PIT still cannot withdraw.</p>
          <ExternalLink className="primary" href={hyperliquidApp(net)}>
            Open official page
          </ExternalLink>
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 3 ? (
        <>
          <h1>Create or verify the PIT session.</h1>
          <PermissionsCard />
          <SessionNote />
          <button type="button" className="primary" onClick={onSession} disabled={bindBusy || !companionUp || !boundWallet}>
            {sessionAlive ? "Session live on this computer" : "Create PIT session"}
          </button>
          {agent ? (
            <p className="fine">
              PIT Agent {agentName || ""} {agent}. If Hyperliquid still lists this agent, PIT reuses it.
            </p>
          ) : null}
          {bindError ? (
            <p className="err" role="alert">
              {bindError}
            </p>
          ) : null}
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 4 ? (
        <>
          <h1>Approve PIT on Hyperliquid.</h1>
          <p className="lead">
            Open the official API page. Authorize API Wallet with the name and address below. PIT cannot withdraw.
          </p>
          {agentName ? <p>Name {agentName}</p> : null}
          {agent ? <p>PIT Agent {agent}</p> : null}
          <p>{hlAgent?.ok ? "Your trading account is ready." : hlAgent?.detail || "Waiting for Hyperliquid approval."}</p>
          <ExternalLink className="primary" href={hyperliquidAPI(net)}>
            Open official page
          </ExternalLink>
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 5 ? (
        <>
          <h1>Protect private research.</h1>
          <p className="lead">
            Sign in the paired browser. This computer stores the sealed-path authorization for 24 hours. The website never
            receives it.
          </p>
          <p>{checks.find((c) => c.name === "direct_auth")?.ok ? "Protected on this computer." : "Waiting for the wallet signature."}</p>
          <ExternalLink className="primary" href={LINKS.app}>
            Open official page
          </ExternalLink>
          <button type="button" className="linkish" onClick={onCheck}>
            Check again
          </button>
        </>
      ) : null}
      {step === 6 ? (
        <>
          <h1>Verify private compute is ready.</h1>
          <p className="lead">This is compute money. It is not Hyperliquid trading capital.</p>
          <ComputeCard checks={checks} onCheck={onCheck} />
        </>
      ) : null}
      {step === 7 ? (
        <>
          <h1>Set policy.</h1>
          <p className="lead">The model cannot raise clip, leverage, or permissions.</p>
          <PolicyLaw pinned={pinned} onPin={onPolicy} busy={bindBusy || !boundWallet} />
        </>
      ) : null}
      {step === 8 ? (
        <>
          <h1>Ready to discover.</h1>
          <p className="lead">Markets is live public marks. Private research stays sealed. Authorize stays on this computer.</p>
          <p>
            {boundWallet ? "Wallet connected. " : "Wallet still unbound. "}
            {sessionAlive ? "Session live. " : "Session still needed. "}
            {hlAgent?.ok ? "Hyperliquid approved. " : "Hyperliquid still needs approval. "}
            {pinned ? "Policy pinned." : "Policy still unpinned."}
          </p>
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
