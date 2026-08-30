import { ExternalLink } from "./ExternalLink";
import type { ActivityEvent } from "./companion";

export const OG_STAGES = [
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
  "READY",
] as const;

export const CHAT_AGENT_COPY = {
  cannotAuthorize: "Chat cannot AUTHORIZE",
  acceptOnDesk: "Type AUTHORIZE on Research",
  pair: "Pair this browser",
} as const;

const STEPS = [
  { id: "watch", label: "Watch" },
  { id: "og", label: "0G Direct" },
  { id: "committee", label: "Committee" },
  { id: "why", label: "Why enter" },
  { id: "preview", label: "Preview" },
  { id: "auth", label: "AUTHORIZE here" },
] as const;

type Coin = {
  coin: string;
  why?: string;
  mark?: number;
  eligible?: boolean;
  executionFeasible?: boolean;
  previewReady?: boolean;
};

export function AgentRun({
  coins,
  island,
  awaiting,
  preview,
  activity,
  onOpenPreview,
  onResearch,
}: {
  coins: Coin[];
  island?: {
    busy: boolean;
    coin: string;
    stage: string;
    elapsedMs: number;
    jobId: string;
    pollMiss?: boolean;
    roles: Array<{ role?: string; verify_e2ee?: string }>;
  };
  awaiting?: boolean;
  preview?: {
    market?: string;
    side?: string;
    sz?: number;
    hash?: string;
    notionalUsd?: number;
    reasons?: string[];
    eligible?: boolean;
  } | null;
  activity: ActivityEvent[];
  onOpenPreview: () => void;
  onResearch: (coin: string) => void;
}) {
  const best =
    coins.find((c) => c.previewReady) ||
    coins.find((c) => c.executionFeasible) ||
    coins.find((c) => c.eligible);
  const stage = island?.stage || "";
  const stageI = Math.max(0, OG_STAGES.indexOf(stage as (typeof OG_STAGES)[number]));
  const active =
    awaiting || preview?.eligible
      ? "auth"
      : island?.busy
        ? stageI >= 5
          ? "committee"
          : "og"
        : best
          ? "watch"
          : "watch";
  const txs = activity.filter((e) => e.kind && /order|approval|evidence|preview|research/.test(e.kind)).slice(-6).reverse();

  return (
    <aside className="agent-run" aria-label="Agent run">
      <ol className="agent-steps">
        {STEPS.map((s) => (
          <li key={s.id} className={s.id === active ? "on" : ""}>
            {s.label}
          </li>
        ))}
      </ol>

      {best ? (
        <article className="agent-card">
          <p className="label">Best live book</p>
          <p className="agent-coin">{best.coin}</p>
          <p>{best.why || "Highest host rank among policy-eligible Hyperliquid books. Not a model score."}</p>
          <p className="fine">
            {best.executionFeasible ? "Executable for this account." : "Policy-visible. Execution still gated on this computer."}
          </p>
          <button type="button" className="primary" onClick={() => onResearch(best.coin)} disabled={Boolean(island?.busy)}>
            Research {best.coin} on 0G
          </button>
        </article>
      ) : (
        <article className="agent-card">
          <p className="label">Watch</p>
          <p>No policy-eligible book yet. Empty is honest. PIT will not invent an opportunity.</p>
        </article>
      )}

      {island?.busy ? (
        <article className="agent-card og" role="status">
          <p className="label">0G Direct TeeML · {island.coin}</p>
          <p className="agent-stage">{island.stage.replaceAll("_", " ")}</p>
          <div className="og-track" aria-hidden="true">
            {OG_STAGES.map((s, i) => (
              <span key={s} className={i <= stageI ? "lit" : ""} />
            ))}
          </div>
          <p>
            {(island.elapsedMs / 1000).toFixed(1)}s elapsed · Researcher {mark(island.roles, "researcher")} · Challenger{" "}
            {mark(island.roles, "challenger")} · Risk {mark(island.roles, "risk")}
          </p>
          {island.pollMiss ? <p role="status">Live view delayed. The sealed job is still running.</p> : null}
          {island.jobId ? <p className="fine">Job {island.jobId}</p> : null}
          <p className="fine">Compute money, not trading capital. {CHAT_AGENT_COPY.cannotAuthorize}.</p>
        </article>
      ) : null}

      {preview?.market ? (
        <article className="agent-card preview">
          <p className="label">Exact preview</p>
          <p className="agent-coin">
            {preview.market} {preview.side} {preview.sz}
          </p>
          {preview.notionalUsd ? <p>${preview.notionalUsd} host-sized</p> : null}
          {preview.hash ? <p className="fine">hash {preview.hash.slice(0, 18)}…</p> : null}
          {preview.reasons?.length ? <p className="fine">{preview.reasons[0]}</p> : null}
          <p className="fine">{CHAT_AGENT_COPY.cannotAuthorize}. {CHAT_AGENT_COPY.acceptOnDesk}.</p>
          <button type="button" className="primary" onClick={onOpenPreview}>
            Accept on this computer
          </button>
        </article>
      ) : awaiting ? (
        <article className="agent-card preview">
          <p className="label">Waiting for you</p>
          <p>Exact preview is ready. {CHAT_AGENT_COPY.acceptOnDesk}.</p>
          <button type="button" className="primary" onClick={onOpenPreview}>
            Accept on this computer
          </button>
        </article>
      ) : null}

      {txs.length > 0 ? (
        <article className="agent-card ledger">
          <p className="label">Desk ledger</p>
          <ul className="tx-rail">
            {txs.map((e, i) => (
              <li key={`${e.ts}-${e.kind}-${i}`}>
                <strong>{labelKind(e.kind)}</strong>
                <span>
                  {e.market || ""} {e.status || ""}
                </span>
                {e.oid ? <span className="fine">OID {e.oid}</span> : null}
                {e.tx_link ? (
                  <ExternalLink href={e.tx_link}>tx</ExternalLink>
                ) : e.tx ? (
                  <span className="fine">{e.tx.slice(0, 12)}…</span>
                ) : null}
              </li>
            ))}
          </ul>
          <p className="fine">Historical fills never appear inside a new preview.</p>
        </article>
      ) : null}
    </aside>
  );
}

function mark(roles: Array<{ role?: string; verify_e2ee?: string }>, name: string) {
  return roles.some((r) => String(r.role).toLowerCase() === name && String(r.verify_e2ee).toUpperCase() === "OK")
    ? "verified"
    : "pending";
}

function labelKind(kind?: string) {
  const k = String(kind || "event");
  if (k === "order.submitted") return "Order submitted";
  if (k === "order.filled") return "Fill";
  if (k === "approval.accepted") return "Authorized";
  if (k === "preview.ready") return "Preview ready";
  if (k === "research.started") return "0G research";
  if (k === "evidence.filed") return "Evidence";
  return k.replaceAll(".", " ");
}
