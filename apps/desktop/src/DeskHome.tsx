import { BrandMark } from "./BrandMark";
import { EvidenceStrip } from "./EvidenceStrip";
import { PairingDock } from "./PairingDock";
import { ExternalLink } from "./ExternalLink";
import { compactNum } from "./format";
import type { NextFix } from "./nextFix";
import type { Probe } from "./readiness";

type Coin = {
  coin: string;
  why?: string;
  trend?: string;
  mark: number;
  eligible?: boolean;
  executionFeasible?: boolean;
  previewReady?: boolean;
  rankGroup?: number;
};

function Chip({ ok, label, value }: { ok: boolean; label: string; value: string }) {
  return (
    <span className={ok ? "chip ok" : "chip fail"}>
      {label} {value}
    </span>
  );
}

export function DeskHome({
  ready,
  doing,
  items,
  attention,
  code,
  companionUp,
  sessionAlive,
  computeReady,
  protectedOk,
  policyPinned,
  hlApproved,
  researchBusy,
  researchStage,
  researchKind,
  awaitingAuth,
  expires,
  paired,
  pairingDevices,
  onRotatePair,
  coins,
  lastEvent,
  mode,
  exposure,
  onResearch,
  onGo,
}: {
  ready: boolean;
  doing: string;
  items: Probe[];
  attention: NextFix;
  code: string;
  companionUp: boolean;
  sessionAlive: boolean;
  computeReady: boolean;
  protectedOk: boolean;
  policyPinned: boolean;
  hlApproved: boolean;
  researchBusy?: boolean;
  researchStage?: string;
  researchKind?: string;
  awaitingAuth?: boolean;
  expires?: string;
  paired?: boolean;
  pairingDevices?: number;
  onRotatePair?: () => void;
  coins: Coin[];
  lastEvent?: string;
  mode?: string;
  exposure?: string;
  onResearch: (coin: string) => void;
  onGo: (view: "markets" | "research" | "security" | "chat" | "automation" | "portfolio" | "activity") => void;
}) {
  const showPair = true;
  const best =
    coins.find((c) => c.previewReady) ||
    coins.find((c) => c.executionFeasible) ||
    coins.find((c) => c.eligible);
  const liveBook = policyPinned && Boolean(coins.find((c) => c.executionFeasible) || coins.find((c) => c.previewReady));
  const sealedNow = policyPinned && Boolean(researchBusy);
  const modeLabel = mode === "guarded" ? "Guarded Autonomy" : mode === "research_only" ? "Research Only" : "Manual";
  return (
    <main className="page dense desk-home">
      <div className="page-head">
        <div>
          <p className="eyebrow">Desk</p>
          <h1>{researchBusy ? doing : awaitingAuth ? "Waiting for you" : ready ? "Ready to discover" : attention.title}</h1>
        </div>
      </div>
      <p className="lead">{doing}</p>
      <ol className="demo-path" aria-label="New user path">
        <li className={items.find((p) => p.id === "wallet")?.state === "ok" ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            1. Connect wallet
          </button>
        </li>
        <li className={protectedOk ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            2. Protect 0G strategy
          </button>
        </li>
        <li className={hlApproved ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            3. Connect Hyperliquid
          </button>
        </li>
        <li className={sessionAlive ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            4. Scoped session
          </button>
        </li>
        <li className={policyPinned ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            5. Pin policy{policyPinned ? "" : " — draft only until you pin"}
          </button>
        </li>
        <li className={liveBook ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("markets")}>
            6. Live opportunities
          </button>
        </li>
        <li className={sealedNow ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => (best ? onResearch(best.coin) : onGo("research"))}>
            7. Sealed research
          </button>
        </li>
        <li className={awaitingAuth ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("research")}>
            8. AUTHORIZE on this computer
          </button>
        </li>
      </ol>
      <div className="chip-row" aria-label="Readiness">
        <Chip ok={protectedOk} label="Research" value={protectedOk ? "protected" : "needs protect"} />
        <Chip ok={computeReady} label="Compute" value={computeReady ? "funded" : "needs funds"} />
        <Chip ok={sessionAlive} label="Session" value={sessionAlive ? "live" : "none"} />
        <Chip ok={hlApproved} label="Hyperliquid" value={hlApproved ? "approved" : "needs approval"} />
        <Chip ok={policyPinned} label="Policy" value={policyPinned ? "pinned" : "unpinned"} />
      </div>
      <ol className="workflow" aria-label="Desk workflow">
        <li className={best && policyPinned ? "on" : ""}>
          <span>Discover</span>
          <strong>{best && policyPinned ? best.coin : "none"}</strong>
        </li>
        <li className={researchBusy || (researchKind && policyPinned) ? "on" : ""}>
          <span>Research</span>
          <strong>{researchBusy ? (researchStage || "running").replaceAll("_", " ") : researchKind ? researchKind.replaceAll("_", " ") : "idle"}</strong>
        </li>
        <li className={awaitingAuth ? "on" : ""}>
          <span>Decision</span>
          <strong>{awaitingAuth ? "awaiting AUTHORIZE" : "none"}</strong>
        </li>
        <li>
          <span>Proof</span>
          <strong>0G trail</strong>
        </li>
      </ol>
      <dl className="metrics">
        <div>
          <dt>What PIT is doing</dt>
          <dd>{researchBusy ? doing : awaitingAuth ? "Preview waiting" : runningCopy(mode)}</dd>
        </div>
        <div>
          <dt>Best opportunity</dt>
          <dd>{best ? `${best.coin} ${compactNum(best.mark)}` : "none"}</dd>
        </div>
        <div>
          <dt>Current exposure</dt>
          <dd>{exposure || "—"}</dd>
        </div>
        <div>
          <dt>Autonomy</dt>
          <dd>{modeLabel}</dd>
        </div>
      </dl>
      <section className="next-row">
        <div>
          <p className="label">Needs you</p>
          <h2>{attention.title}</h2>
          <p className="fine" style={{ margin: 0 }}>
            {attention.why}
          </p>
        </div>
        <div className="cta-row">
          {ready && !researchBusy && !awaitingAuth ? (
            <button type="button" className="primary" onClick={() => onGo("markets")}>
              Open Markets
            </button>
          ) : null}
          {awaitingAuth ? (
            <button type="button" className="primary" onClick={() => onGo("research")}>
              Open preview
            </button>
          ) : null}
          {researchBusy ? (
            <button type="button" className="primary" onClick={() => onGo("research")}>
              Open Research
            </button>
          ) : null}
          <button type="button" className="linkish" onClick={() => onGo("chat")}>
            Ask PIT
          </button>
          <button type="button" className="linkish" onClick={() => onGo("automation")}>
            Automation
          </button>
          {attention.href ? (
            <ExternalLink className="linkish" href={attention.href}>
              {attention.hrefLabel || "Open official page"}
            </ExternalLink>
          ) : null}
        </div>
      </section>
      {best ? (
        <section>
          <p className="label">Best opportunity</p>
          <ul className="desk-ops">
            <li>
              <BrandMark symbol={best.coin} />
              <strong>{best.coin}</strong>
              <span className="mark-num">{compactNum(best.mark)}</span>
              <span className="fine" style={{ margin: 0 }}>
                {best.why || best.trend || "In policy universe."}
              </span>
              <button type="button" className="primary" disabled={researchBusy || !protectedOk || !computeReady || !policyPinned} onClick={() => onResearch(best.coin)}>
                Research privately
              </button>
            </li>
          </ul>
        </section>
      ) : (
        <p className="empty">No opportunities match your policy yet. Empty is honest.</p>
      )}
      <EvidenceStrip onOpen={() => onGo("activity")} />
      {lastEvent ? <p className="fine">Recently: {lastEvent}</p> : null}
      {showPair ? (
        <PairingDock
          code={code}
          expires={expires}
          companionUp={companionUp}
          paired={paired}
          devices={pairingDevices}
          onRotate={onRotatePair}
        />
      ) : null}
    </main>
  );
}

function runningCopy(mode?: string) {
  if (mode === "guarded") return "Guarded Autonomy is live inside your policy.";
  if (mode === "research_only") return "Research Only — scan and prepare, never execute.";
  return "Manual. Waiting for you.";
}
