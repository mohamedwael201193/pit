import { FormEvent } from "react";
import { NAMED } from "./namedStates";
import { NetworkBanner } from "./NetworkBanner";
import { NetworkToggle } from "./NetworkToggle";
import { PermissionsCard } from "./Permissions";
import { PolicyLaw } from "./PolicyLaw";
import { LINKS, hyperliquidAPI } from "./links";
import { ExternalLink } from "./ExternalLink";
import { type DoctorCheck } from "./companion";
import { PairingDock } from "./PairingDock";
import { HyperliquidCard } from "./HyperliquidCard";
import { OnboardRail } from "./OnboardRail";
import { computeOnboard, onboardInput } from "./onboard";
import type { Probe } from "./readiness";

type Net = "mainnet" | "testnet";

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
  paired,
  onRotate,
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
  paired?: boolean;
  onRotate?: () => void;
}) {
  const hlAgent = checks.find((c) => c.name === "hl_agent");
  const protectOk = Boolean(checks.find((c) => c.name === "direct_auth")?.ok);
  const board = computeOnboard(
    onboardInput(
      companionUp,
      { sign: false, paired: Boolean(paired), wallet: boundWallet },
      checks,
      sessionAlive,
    ),
  );
  const currentN = board.steps.find((s) => s.id === board.current)?.n ?? 1;
  const view = Math.min(Math.max(step, 0), currentN - 1);
  const viewId = board.steps[view]?.id || board.current;

  return (
    <section className="setup">
      <p className="eyebrow">
        FIRST RUN · STEP {currentN} OF 5 · {board.steps[currentN - 1]?.title}
      </p>
      <OnboardRail steps={board.steps} />
      {viewId === "pair" ? (
        <>
          <h1>Pair this browser with PIT Desktop.</h1>
          <p className="lead">
            Type the one-time code on the website. It expires in two minutes and works once. {NAMED.SEED_FORBIDDEN} The
            website never receives a session key.
          </p>
          <PairingDock
            code={code}
            expires={expires}
            companionUp={companionUp}
            paired={paired}
            onRotate={onRotate}
          />
        </>
      ) : null}
      {viewId === "protect" ? (
        <>
          <h1>Protect my strategy.</h1>
          <p className="lead">
            Sign in the bound wallet. This computer stores the sealed-path authorization for 24 hours. The website never
            receives it.
          </p>
          <p>
            {protectOk
              ? "This computer stored the authorization."
              : "Waiting for the wallet signature. Check again after you sign."}
          </p>
          <ExternalLink className="primary" href={LINKS.protect}>
            Protect my strategy
          </ExternalLink>
          <details className="sec-fold">
            <summary>Advanced recovery — bind a public address here</summary>
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
          </details>
        </>
      ) : null}
      {viewId === "hyperliquid" ? (
        <>
          <h1>Connect Hyperliquid.</h1>
          <p className="lead">
            PIT creates the agent on this computer. Approve it with your master wallet. Do not paste an API wallet into
            PIT.
          </p>
          <NetworkToggle net={net} onChange={setNet} />
          <NetworkBanner net={net} />
          <HyperliquidCard
            net={net}
            agent={agent}
            agentName={agentName}
            sessionAlive={sessionAlive}
            approved={Boolean(hlAgent?.ok)}
            approvedDetail={hlAgent?.detail}
            busy={bindBusy}
            account={boundWallet}
            onCreateSession={onSession}
            onCheck={onCheck}
          />
          <PermissionsCard />
        </>
      ) : null}
      {viewId === "policy" ? (
        <>
          <h1>Pin policy.</h1>
          <p className="lead">The model cannot raise clip, leverage, or permissions. After pin, Security keeps the editor so you can change and re-pin anytime.</p>
          <PolicyLaw pinned={pinned} onPin={onPolicy} busy={bindBusy || !boundWallet} />
        </>
      ) : null}
      {viewId === "ready" ? (
        <>
          <h1>Ready to trade.</h1>
          <p className="lead">Markets is live public marks. Private research stays sealed. Authorize stays on this computer.</p>
          <ul className="onboard-checks">
            <li>Browser paired ✓</li>
            <li>Strategy protected ✓</li>
            <li>Hyperliquid agent verified ✓</li>
            <li>Policy pinned — edit anytime on Security ✓</li>
          </ul>
        </>
      ) : null}
      <div className="row">
        {view > 0 ? (
          <button type="button" className="off" onClick={() => setStep(view - 1)}>
            Back
          </button>
        ) : null}
        {board.ready ? (
          <button type="button" className="on" onClick={onDone}>
            Open the desk
          </button>
        ) : (
          <button
            type="button"
            className="on"
            onClick={() => {
              onCheck();
              setStep(currentN - 1);
            }}
          >
            Check again
          </button>
        )}
      </div>
      <p className="fine">
        Official Hyperliquid API is {hyperliquidAPI(net).replace("https://", "")}. PIT still cannot withdraw.
      </p>
    </section>
  );
}
