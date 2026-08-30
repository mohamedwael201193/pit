import { PairingDock } from "./PairingDock";
import { ComputeCard } from "./ComputeCard";
import { ExternalLink } from "./ExternalLink";
import { HyperliquidCard } from "./HyperliquidCard";
import { PolicyEditor, emptyPolicy } from "./PolicyEditor";
import { explorerAddress, hyperliquidAPI, LINKS } from "./links";
import type { AccountSummary, DoctorCheck, HostPolicy, LocalStatus, SecurityDomain } from "./companion";
import { checkNamed } from "./companion";
import type { NextFix } from "./nextFix";
import type { Probe } from "./readiness";

export type { SecurityDomain };

function tone(state?: string) {
  const s = String(state || "").toUpperCase();
  if (s === "READY") return "ok";
  if (s === "BLOCKED") return "fail";
  return "wait";
}

function probeOf(items: Probe[], id: string) {
  return items.find((p) => p.id === id);
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
  attention,
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
  attention: NextFix;
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
  const wallet = status?.wallet || "";
  const hl = checkNamed(checks, "hl_agent");
  const ready =
    attention.title === "Desk is ready" || attention.title.startsWith("Watching");
  const pairNow = /pair|wallet/i.test(attention.title);
  const spine = [
    probeOf(items, "local"),
    probeOf(items, "session"),
    probeOf(items, "hl_agent"),
    probeOf(items, "policy"),
  ].filter(Boolean) as Probe[];

  return (
    <main className="page security-page">
      <header className="sec-head">
        <div>
          <p className="eyebrow">Security</p>
          <h1>{ready ? "Workspace ready" : "Do this next"}</h1>
          <p className="sec-line">{attention.why}</p>
        </div>
        <button type="button" className="linkish" onClick={onCheck} disabled={busy}>
          Check again
        </button>
      </header>

      <section className={ready ? "sec-next ok" : "sec-next"} aria-label="Next action">
        <p className="sec-step">{ready ? "Ready" : "1"}</p>
        <div>
          <h2>{attention.title}</h2>
          <p>{attention.fix}</p>
        </div>
        <NextControl
          attention={attention}
          busy={busy}
          onSession={onSession}
          onPin={() => onPolicyPin(policy || emptyPolicy())}
          onCheck={onCheck}
        />
      </section>

      {spine.length ? (
        <ol className="sec-spine" aria-label="Setup spine">
          {spine.map((p) => (
            <li key={p.id} className={p.state === "ok" ? "ok" : p.state === "fail" ? "fail" : "wait"}>
              <strong>{p.label}</strong>
              <span>{p.state === "ok" ? "ready" : p.state === "fail" ? "blocked" : "needed"}</span>
            </li>
          ))}
        </ol>
      ) : null}

      <section className="sec-block" id="session">
        <h2>Session</h2>
        <p className="sec-line">Order and cancel only. Withdraw stays impossible. Chat cannot authorize.</p>
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
      </section>

      {summary?.execWhy ? (
        <p className="fine" role="status">
          {summary.execGate ? `Execution blocked: ${summary.execGate.replaceAll("_", " ")}. ` : ""}
          {summary.execWhy}
        </p>
      ) : null}

      <section className="sec-block" id="policy">
        <h2>Policy</h2>
        <p className="sec-line">You edit. You pin. The model cannot raise clip, leverage, or permissions.</p>
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
      </section>

      {code !== undefined ? (
        <section className="sec-block" id="pairing">
          <h2>Browser</h2>
          <p className="sec-line">Optional. Orders still sign on this computer if the site is unpaired.</p>
          {pairNow ? (
            <PairingDock
              code={code}
              expires={expires}
              companionUp={Boolean(companionUp)}
              paired={paired}
              devices={pairingDevices}
              onRotate={onRotatePair}
            />
          ) : (
            <details className="sec-fold">
              <summary>{paired ? "Browser paired" : "Pair a browser"}</summary>
              <PairingDock
                code={code}
                expires={expires}
                companionUp={Boolean(companionUp)}
                paired={paired}
                devices={pairingDevices}
                onRotate={onRotatePair}
              />
            </details>
          )}
        </section>
      ) : null}

      <section className="sec-block" id="compute">
        <h2>Private compute</h2>
        <p className="sec-line">0G Direct credit is not Hyperliquid buying power.</p>
        <ComputeCard checks={checks} onCheck={onCheck} />
      </section>

      <section className="sec-danger" aria-label="Halt and revoke">
        <div>
          <h2>Halt</h2>
          <p>You flip this. The model cannot. New orders stop. Positions are not flattened.</p>
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
        </div>
        <div>
          <h2>Revoke</h2>
          <p>Delete the local session, then remove the PIT agent on Hyperliquid if you want the listing gone.</p>
          <div className="cta-row">
            <button type="button" className="linkish" onClick={onRevoke} disabled={busy || !sessionAlive}>
              Revoke local session
            </button>
            <ExternalLink className="linkish" href={hyperliquidAPI(net)}>
              Open Hyperliquid API
            </ExternalLink>
          </div>
        </div>
      </section>

      <details className="sec-fold">
        <summary>Workspace identity</summary>
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
      </details>

      <details className="sec-fold">
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
                <td>{d.id}</td>
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

      <details className="sec-fold">
        <summary>Diagnostics</summary>
        {checks.length === 0 ? (
          <p>Waiting for the local companion on 127.0.0.1:17373.</p>
        ) : (
          <ul className="doctor">
            {checks.map((c) => (
              <li key={c.name}>
                <strong>{c.ok ? "ok" : "fail"}</strong> {c.name} - {c.detail}
              </li>
            ))}
          </ul>
        )}
        <p className="fine">{updateNote}</p>
        <p className="fine">Restart {restartAllowed ? "allowed" : "refused because research is running."}</p>
        <button type="button" className="linkish" onClick={onForget}>
          Forget this workspace memory
        </button>
      </details>
    </main>
  );
}

function NextControl({
  attention,
  busy,
  onSession,
  onPin,
  onCheck,
}: {
  attention: NextFix;
  busy?: boolean;
  onSession: () => void;
  onPin: () => void;
  onCheck: () => void;
}) {
  const t = attention.title;
  if (/Create a local session/i.test(t)) {
    return (
      <button type="button" className="primary" onClick={onSession} disabled={busy}>
        Create session
      </button>
    );
  }
  if (/Pin a trading policy|Policy cap/i.test(t)) {
    return (
      <a className="primary" href="#policy">
        Edit policy
      </a>
    );
  }
  if (/Approve PIT/i.test(t) && attention.href) {
    return (
      <ExternalLink className="primary" href={attention.href}>
        {attention.hrefLabel || "Approve PIT"}
      </ExternalLink>
    );
  }
  if (/Protect my strategy/i.test(t) && attention.href) {
    return (
      <ExternalLink className="primary" href={attention.href}>
        {attention.hrefLabel || "Protect my strategy"}
      </ExternalLink>
    );
  }
  if (/pair|wallet/i.test(t) && attention.href) {
    return (
      <ExternalLink className="primary" href={attention.href}>
        {attention.hrefLabel || "Open pairing"}
      </ExternalLink>
    );
  }
  if (/Fund private research|Direct funds/i.test(t) && attention.href) {
    return (
      <ExternalLink className="primary" href={attention.href}>
        {attention.hrefLabel || "Open 0G Direct funds"}
      </ExternalLink>
    );
  }
  if (/Check private compute|Check again/i.test(t)) {
    return (
      <button type="button" className="primary" onClick={onCheck} disabled={busy}>
        Check again
      </button>
    );
  }
  if (attention.href && attention.hrefLabel) {
    return (
      <ExternalLink className="primary" href={attention.href}>
        {attention.hrefLabel}
      </ExternalLink>
    );
  }
  if (/Desk is ready/i.test(t)) {
    return (
      <a className="primary" href="#policy">
        Review policy
      </a>
    );
  }
  return null;
}
