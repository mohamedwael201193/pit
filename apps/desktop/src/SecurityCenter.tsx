import { BrandMark } from "./BrandMark";
import { PairingDock } from "./PairingDock";
import { ComputeCard } from "./ComputeCard";
import { ExternalLink } from "./ExternalLink";
import { HyperliquidCard } from "./HyperliquidCard";
import { PolicyEditor, emptyPolicy } from "./PolicyEditor";
import { explorerAddress, hyperliquidAPI, LINKS } from "./links";
import type { AccountSummary, DoctorCheck, HostPolicy, LocalStatus, SecurityDomain } from "./companion";
import { checkNamed } from "./companion";
import type { Probe } from "./readiness";

export type { SecurityDomain };

function tone(state?: string) {
  const s = String(state || "").toUpperCase();
  if (s === "READY") return "ok";
  if (s === "BLOCKED") return "fail";
  return "wait";
}

export function SecurityCenter({
  domains,
  checks,
  items,
  net,
  status,
  agent,
  sessionAlive,
  approved,
  tradingCapital,
  summary,
  policy,
  pinned,
  policyHash,
  consequences,
  allowed,
  refused,
  identityNote,
  calibCopy,
  updateNote,
  restartAllowed,
  busy,
  onSession,
  onPolicyPreview,
  onPolicyPin,
  onCheck,
  onRevoke,
  onKill,
  onForget,
  code,
  expires,
  companionUp,
  paired,
  pairingDevices,
  onRotatePair,
}: {
  domains: SecurityDomain[];
  checks: DoctorCheck[];
  items: Probe[];
  net: string;
  status: LocalStatus | null;
  agent: string;
  sessionAlive: boolean;
  approved: boolean;
  tradingCapital?: string;
  summary?: AccountSummary;
  policy?: HostPolicy | null;
  pinned?: boolean;
  policyHash?: string;
  consequences?: string[];
  allowed?: string[];
  refused?: string[];
  identityNote?: string;
  calibCopy?: string;
  updateNote?: string;
  restartAllowed?: boolean;
  busy?: boolean;
  onSession: () => void;
  onPolicyPreview: (p: HostPolicy) => void;
  onPolicyPin: (p: HostPolicy) => void;
  onCheck: () => void;
  onRevoke: () => void;
  onKill: (on: boolean) => void;
  onForget: () => void;
  code?: string;
  expires?: string;
  companionUp?: boolean;
  paired?: boolean;
  pairingDevices?: number;
  onRotatePair?: () => void;
}) {
  const missing = items.filter((p) => p.state !== "ok");
  const wallet = status?.wallet || "";
  const hl = checkNamed(checks, "hl_agent");
  return (
    <main className="page dense">
      <div className="page-head">
        <div>
          <p className="eyebrow">Security</p>
          <h1>{missing.length ? "Action required" : "Workspace ready"}</h1>
        </div>
        <button type="button" className="linkish" onClick={onCheck} disabled={busy}>
          Check again
        </button>
      </div>
      <p className="lead">What is ready, what is missing, and the one next action. Order and cancel only. Chat cannot authorize or pin policy.</p>

      {code !== undefined ? (
        <PairingDock
          code={code}
          expires={expires}
          companionUp={Boolean(companionUp)}
          paired={paired}
          devices={pairingDevices}
          onRotate={onRotatePair}
        />
      ) : null}

      <section className="ready-list" aria-label="Setup readiness">
        {items.map((p) => (
          <article key={p.id} className={`ready-item ${p.state}`}>
            <p className="label">{p.label}</p>
            <p className={p.state === "ok" ? "ok" : p.state === "optional" ? "mute" : "bad"}>
              {p.state === "ok" ? "READY" : p.state === "optional" ? "OPTIONAL" : p.state === "fail" ? "BLOCKED" : "ACTION REQUIRED"}
            </p>
            <p>{p.detail}</p>
            {p.id === "wallet" && p.state !== "ok" ? <p className="fine">Open Desk first-run or bind this computer.</p> : null}
            {p.id === "direct" && p.state !== "ok" ? (
              <ExternalLink className="primary" href={LINKS.app}>
                Protect my strategy
              </ExternalLink>
            ) : null}
            {p.id === "direct_credit" && p.state !== "ok" && !(p.detail || "").toLowerCase().includes("unread") && !(p.detail || "").toLowerCase().includes("sponsor") ? (
              <ExternalLink className="primary" href={LINKS.pcAdvanced}>
                Open 0G Direct funds
              </ExternalLink>
            ) : null}
            {p.id === "pairing" && p.state !== "ok" ? (
              <ExternalLink className="primary" href={LINKS.pair}>
                Open pairing
              </ExternalLink>
            ) : null}
            {p.id === "session" && p.state !== "ok" ? (
              <button type="button" className="primary" onClick={onSession} disabled={busy}>
                Create session
              </button>
            ) : null}
            {p.id === "hl_agent" && p.state !== "ok" ? (
              <ExternalLink className="primary" href={hyperliquidAPI(net)}>
                Approve PIT on Hyperliquid
              </ExternalLink>
            ) : null}
            {p.id === "policy" && p.state !== "ok" ? (
              <button type="button" className="primary" onClick={() => onPolicyPin(policy || emptyPolicy())} disabled={busy}>
                Pin default policy
              </button>
            ) : null}
            {p.id === "storage" && p.state !== "ok" ? (
              <ExternalLink className="linkish" href="https://docs.0g.ai">
                0G docs
              </ExternalLink>
            ) : null}
          </article>
        ))}
      </section>

      <HyperliquidCard
        net={net}
        agent={agent}
        agentName={status?.agentName}
        sessionAlive={sessionAlive}
        sessionExpires={status?.sessionExpires}
        approved={approved}
        approvedDetail={hl?.detail}
        busy={busy}
        tradingCapital={tradingCapital || summary?.accountValue}
        account={wallet}
        onCreateSession={onSession}
        onCheck={onCheck}
        onRevoke={onRevoke}
      />

      {summary?.execWhy ? (
        <p className="fine" role="status">
          {summary.execGate ? `Execution blocked: ${summary.execGate.replaceAll("_", " ")}. ` : ""}
          {summary.execWhy}
        </p>
      ) : null}

      <ComputeCard checks={checks} onCheck={onCheck} />

      <PolicyEditor
        current={policy}
        consequences={consequences}
        allowed={allowed}
        refused={refused}
        pinned={pinned}
        policyHash={policyHash}
        busy={busy}
        onPreview={onPolicyPreview}
        onPin={onPolicyPin}
      />

      <article className="card">
        <p className="label">This workspace</p>
        <p>Wallet {wallet || "unbound"}</p>
        <p>Network {net === "mainnet" ? "MAINNET" : "TESTNET"}</p>
        <p>PIT Agent {status?.agentName || "none"}</p>
        <p>Agent address {agent || "none"}</p>
        <p>
          Session {sessionAlive ? "Active" : "none"}
          {status?.sessionExpires ? ` until ${new Date(status.sessionExpires).toISOString().replace(".000Z", "Z")}` : ""}
        </p>
        {status?.workspace ? <p>Desk ID {status.workspace}</p> : null}
        {wallet ? (
          <ExternalLink className="linkish" href={explorerAddress(wallet)}>
            View on explorer
          </ExternalLink>
        ) : null}
        <p className="fine">{identityNote}</p>
        <p className="fine">{calibCopy || "NOT ENOUGH DATA"}</p>
      </article>

      <details className="card">
        <summary>Official domains</summary>
        <table className="desk-table">
          <thead>
            <tr>
              <th>Domain</th>
              <th>State</th>
              <th>Do</th>
            </tr>
          </thead>
          <tbody>
            {domains.map((d) => (
              <tr key={d.id}>
                <td>
                  <span className="asset">
                    {d.id === "hyperliquid" ? <BrandMark symbol="HL" /> : null}
                    {d.id === "compute" || d.id === "tee" || d.id === "storage" ? <BrandMark symbol="0G" /> : null}
                    {d.id}
                  </span>
                </td>
                <td className={tone(d.state)}>{d.state || "NEEDS ACTION"}</td>
                <td>
                  <p>{d.why}</p>
                  {d.href ? (
                    <ExternalLink className="linkish" href={d.href}>
                      {d.hrefLabel || "Open official page"}
                    </ExternalLink>
                  ) : null}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </details>

      <article className="card">
        <p className="label">Kill switch</p>
        <p>You flip this. The model cannot. New orders stop until you turn it off. Positions are not flattened.</p>
        <div className="cta-row">
          <button type="button" className="primary" onClick={() => onKill(true)} disabled={busy || Boolean(status?.kill)}>
            {status?.kill ? "Kill switch is on" : "Halt new orders"}
          </button>
          {status?.kill ? (
            <button type="button" className="linkish" onClick={() => onKill(false)} disabled={busy}>
              Resume this workspace
            </button>
          ) : null}
        </div>
      </article>

      <article className="card">
        <p className="label">Revoke</p>
        <p>Delete the local session, then remove the PIT agent from Hyperliquid if you want the venue listing gone.</p>
        <div className="cta-row">
          <button type="button" className="linkish" onClick={onRevoke} disabled={busy || !sessionAlive}>
            Revoke local session
          </button>
          <ExternalLink className="linkish" href={hyperliquidAPI(net)}>
            Open Hyperliquid API
          </ExternalLink>
        </div>
      </article>

      <details className="card">
        <summary>Diagnostics</summary>
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
        <p className="fine">{updateNote}</p>
        <p className="fine">Restart {restartAllowed ? "allowed" : "refused — research is running."}</p>
        <button type="button" className="linkish" onClick={onForget}>
          Forget this workspace memory
        </button>
      </details>
    </main>
  );
}