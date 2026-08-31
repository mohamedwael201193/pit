import { BrandMark } from "./BrandMark";
import { deskHeadline } from "./deskCopy";
import { EvidenceStrip } from "./EvidenceStrip";
import { ExternalLink } from "./ExternalLink";
import { accountSizeGate, compactNum, compactUsd, marketSizeGate, nearestPolicyClip, nearestVenueMin } from "./format";
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
  minNotional?: number;
  execGate?: string;
  policyClip?: number;
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
  code: _code,
  companionUp: _companionUp,
  sessionAlive,
  computeReady,
  protectedOk,
  policyPinned,
  hlApproved,
  researchBusy,
  researchStage,
  researchKind,
  awaitingAuth,
  expires: _expires,
  paired,
  pairingDevices: _pairingDevices,
  onRotatePair: _onRotatePair,
  coins,
  lastEvent,
  mode,
  exposure,
  buyingPower,
  execWhy,
  execGate,
  capitalNote,
  powerSource,
  fundHref,
  experienceWhy,
  routes,
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
  buyingPower?: number;
  execWhy?: string;
  execGate?: string;
  capitalNote?: string;
  powerSource?: string;
  fundHref?: string;
  experienceWhy?: string;
  routes?: Array<{ action: string; coin?: string; reason: string; execution: string }>;
  onResearch: (coin: string) => void;
  onGo: (view: "markets" | "research" | "security" | "chat" | "automation" | "portfolio" | "activity") => void;
}) {
  const best =
    coins.find((c) => c.previewReady) ||
    coins.find((c) => c.executionFeasible) ||
    coins.find((c) => c.eligible);
  const ranked = [...coins]
    .filter((c) => c.executionFeasible || c.previewReady || c.eligible)
    .sort((a, b) => Number(Boolean(b.executionFeasible)) - Number(Boolean(a.executionFeasible)) || Number(Boolean(b.previewReady)) - Number(Boolean(a.previewReady)))
    .slice(0, 6);
  const liveBook = policyPinned && Boolean(coins.find((c) => c.executionFeasible) || coins.find((c) => c.previewReady));
  const sealedNow = policyPinned && Boolean(researchBusy);
  const modeLabel = mode === "guarded" ? "Sleep Mission" : mode === "research_only" ? "Research" : "Manual";
  const venueMin = nearestVenueMin(coins);
  const execN = coins.filter((c) => c.executionFeasible).length;
  const gate = accountSizeGate({
    have: buyingPower,
    venueMin,
    executable: execN,
    execGate,
    execWhy,
    policyClip: nearestPolicyClip(coins),
  });
  const nearest = [...coins]
    .filter((c) => (c.executionFeasible || c.eligible) && (c.minNotional || 0) > 0)
    .sort((a, b) => (a.minNotional || 0) - (b.minNotional || 0))[0];
  const heroTitle = deskHeadline({
    researchBusy,
    doing,
    awaitingAuth,
    ready,
    canOpen: gate.canOpen,
    execN,
    researchKind,
    attentionTitle: attention.title,
    policyTight: gate.cta === "policy",
  });
  const path = (
      <ol className="demo-path" aria-label="New user path">
        <li className={paired ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            1. Pair this browser
          </button>
        </li>
        <li className={protectedOk ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            2. Protect my strategy
          </button>
        </li>
        <li className={hlApproved && sessionAlive ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            3. Connect Hyperliquid
          </button>
        </li>
        <li className={policyPinned ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            4. Pin policy{policyPinned ? "" : " — draft only until you pin"}
          </button>
        </li>
        <li className={paired && protectedOk && hlApproved && sessionAlive && policyPinned ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => onGo("security")}>
            5. Verify ready
          </button>
        </li>
        <li className={liveBook || sealedNow ? "on" : ""}>
          <button type="button" className="linkish" onClick={() => (best ? onResearch(best.coin) : onGo("research"))}>
            6. Use agent
          </button>
        </li>
      </ol>
  );

  return (
    <main className="page dense desk-home">
      <section className="desk-hero">
        <div>
          <p className="eyebrow">Desk</p>
          <h1>{heroTitle}</h1>
          {doing && doing !== heroTitle ? <p className="lead">{doing}</p> : null}
          <p className="capital-line" role="status">
            This account {compactUsd(gate.have)}
            {nearest ? ` · nearest floor ${nearest.coin} ${compactUsd(nearest.minNotional)}` : ` · this market min ${compactUsd(gate.min)}`}
            {gate.canOpen ? "" : gate.cta === "policy" ? ` · policy gap ${compactUsd(gate.policyGap)}` : ` · ${compactUsd(gate.shortfall)} short`}
            {powerSource ? ` · ${powerSource.replaceAll("_", " ")}` : ""}
            {execGate ? ` · ${execGate.replaceAll("_", " ")}` : ""}
          </p>
          {gate.canOpen ? <p className="fine">{gate.detail}</p> : <p>{execWhy || capitalNote || gate.detail}</p>}
          {experienceWhy ? <p className="fine">{experienceWhy}</p> : null}
          {routes && routes.length ? (
            <p className="route-rail" aria-label="Capital router">
              {routes.map((r) => (
                <span key={r.action} className={`route-pill ${r.execution}`} title={r.reason}>
                  {r.action.toUpperCase()} {r.execution}
                </span>
              ))}
            </p>
          ) : null}
          {gate.canOpen && best && !researchBusy && !awaitingAuth ? (
            <p className="fine">
              <button type="button" className="primary" onClick={() => onResearch(best.coin)}>
                Research {best.coin} privately
              </button>
            </p>
          ) : gate.cta === "policy" ? (
            <p className="fine">
              <button type="button" className="linkish" onClick={() => onGo("security")}>
                Open Policy
              </button>
            </p>
          ) : !gate.canOpen && fundHref ? (
            <p className="fine">
              <ExternalLink className="linkish" href={fundHref}>
                Fund this Hyperliquid account
              </ExternalLink>
            </p>
          ) : null}
        </div>
        <section className="next-row hero-next">
          <div>
            <p className="label">{attention.title === "Desk is ready" || gate.canOpen ? "Next" : "Needs you"}</p>
            <h2>{attention.title}</h2>
            <p className="fine" style={{ margin: 0 }}>
              {attention.why}
            </p>
          </div>
          <div className="cta-row">
            {ready && !researchBusy && !awaitingAuth ? (
              <button type="button" className={gate.canOpen ? "linkish" : "primary"} onClick={() => onGo("markets")}>
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
      </section>
      {ready ? (
        <details className="card">
          <summary>Setup path</summary>
          {path}
        </details>
      ) : (
        path
      )}
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
        <li className={researchBusy || (researchKind === "READY_ELIGIBLE" && policyPinned) ? "on" : researchKind === "READY_STOOD_DOWN" ? "calm" : ""}>
          <span>Research</span>
          <strong>
            {researchBusy
              ? (researchStage || "running").replaceAll("_", " ")
              : researchKind === "READY_STOOD_DOWN"
                ? "no trade survived challenge"
                : researchKind
                  ? researchKind.replaceAll("_", " ")
                  : "idle"}
          </strong>
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
          <dt>Ranked books</dt>
          <dd>{ranked.length ? `${execN} executable · ${ranked.length} shown` : "none"}</dd>
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
      {ranked.length ? (
        <section>
          <p className="label">{execN ? "Executable now" : "Policy-eligible books"}</p>
          <ul className="book-list desk-books" aria-label="Ranked books">
            {ranked.map((c) => {
              const chip = marketSizeGate(c.coin, buyingPower, c.minNotional, c.executionFeasible, c.execGate, c.policyClip);
              return (
                <li key={c.coin}>
                  <button type="button" className="book-row" onClick={() => onGo("markets")}>
                    <span className="asset">
                      <BrandMark symbol={c.coin} size={14} />
                      <strong>{c.coin}</strong>
                    </span>
                    <span className="tile-mark">{compactNum(c.mark)}</span>
                    <span className={`layer-chip ${c.executionFeasible ? "ok" : "pass"}`}>{chip.chip}</span>
                    <span className="fine" style={{ margin: 0 }}>
                      {compactUsd(c.minNotional)} · clip {compactUsd(c.policyClip)}
                    </span>
                  </button>
                </li>
              );
            })}
          </ul>
        </section>
      ) : (
        <p className="empty">No opportunities match your policy yet. Empty is honest.</p>
      )}
      <EvidenceStrip onOpen={() => onGo("activity")} />
      {lastEvent ? <p className="fine">Recently: {lastEvent}</p> : null}
      {!paired ? (
        <p className="fine">Pair this browser first. Session keys stay on this computer.</p>
      ) : null}
    </main>
  );
}

function runningCopy(mode?: string) {
  if (mode === "guarded") return "Sleep Mission is live inside your policy.";
  if (mode === "research_only") return "Research Only — scan and prepare, never execute.";
  return "Manual. Waiting for you.";
}
