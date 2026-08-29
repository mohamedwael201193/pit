import { ExternalLink } from "./ExternalLink";

type EventRow = {
  ts?: number;
  kind?: string;
  market?: string;
  action?: string;
  status?: string;
  job_id?: string;
  preview_hash?: string;
  oid?: string;
  reason?: string;
  root?: string;
  tx?: string;
  tx_link?: string;
  digest?: string;
  link?: string;
};

function dayLabel(ts?: number) {
  if (!ts) return "Earlier";
  const d = new Date(ts);
  const now = new Date();
  const same = d.toDateString() === now.toDateString();
  const y = new Date(now);
  y.setDate(now.getDate() - 1);
  if (d.toDateString() === y.toDateString()) return "Yesterday";
  const week = Date.now() - ts < 7 * 86400000;
  if (same) return "Today";
  if (week) return "This week";
  return d.toISOString().slice(0, 10);
}

function humanKind(kind?: string) {
  const k = String(kind || "event");
  if (k === "mission.enabled") return "Mission started";
  if (k === "mission.stopped") return "Mission stopped";
  if (k.includes("mission.scanned")) return "Universe scanned";
  if (k.includes("mission.exec_blocked") || k.includes("mission.refused")) return "Execution blocked";
  if (k.includes("mission.empty")) return "No opportunity";
  if (k.includes("mission.scan_failed")) return "Scan failed";
  if (k.includes("opportunity")) return "Opportunity found";
  if (k === "research.sealed") return "Sealed";
  if (k === "researcher.verified") return "Researcher verified";
  if (k === "challenger.verified") return "Challenger verified";
  if (k === "risk.verified") return "Risk verified";
  if (k === "tee.verified") return "TEE verified";
  if (k === "policy.pinned") return "Policy pinned";
  if (k === "policy.failed" || k === "policy.fail") return "Policy fail";
  if (k === "candidate" || k.includes("candidate")) return "Candidate";
  if (k.includes("calibration") || k.includes("calib")) return "Calibration";
  if (k === "research.started") return "Research started";
  if (k === "research.verified") return "Research verified";
  if (k === "research.stood_down") return "Research stood down";
  if (k === "research.canceled") return "Research canceled";
  if (k === "research.failed") return "Research failed";
  if (k === "committee.verified") return "Committee verified";
  if (k === "preview.ready") return "Exact preview ready";
  if (k === "approval.accepted") return "Authorization accepted";
  if (k === "approval.rejected") return "Authorization refused";
  if (k === "order.submitted") return "Order submitted";
  if (k === "order.filled") return "Fill";
  if (k === "order.resting") return "Order resting";
  if (k === "order.canceled") return "Order canceled";
  if (k === "order.rejected") return "Order rejected";
  if (k === "position.updated") return "Position";
  if (k === "evidence.filed") return "Evidence published";
  if (k === "evidence.verified") return "Evidence verified";
  if (k === "evidence.failed" || k === "evidence.unavailable") return "Evidence not published";
  if (k.startsWith("research")) return "Research";
  if (k.includes("automation")) return "Prepared";
  if (k.includes("mission")) return "Mission";
  return k.replaceAll(".", " ");
}

function humanStatus(status?: string) {
  const s = String(status || "");
  if (s === "READY_ELIGIBLE") return "eligible preview";
  if (s === "READY_STOOD_DOWN") return "stood down";
  if (s === "awaiting_AUTHORIZE") return "awaiting AUTHORIZE";
  if (s === "filled") return "filled";
  if (s === "resting") return "resting";
  if (s === "COMMITTEE_INCOMPLETE") return "incomplete";
  return s.replaceAll("_", " ");
}

function short(v?: string, n = 10) {
  const s = String(v || "");
  if (s.length <= n + 4) return s;
  return `${s.slice(0, n)}…`;
}

function EventIds({ ev }: { ev: EventRow }) {
  const href = ev.tx_link || ev.link;
  return (
    <span className="timeline-ids">
      {ev.job_id ? <span title={ev.job_id}>job {short(ev.job_id, 8)}</span> : null}
      {ev.oid ? <span title={ev.oid}>OID {ev.oid}</span> : null}
      {ev.preview_hash ? <span title={ev.preview_hash}>preview {short(ev.preview_hash)}</span> : null}
      {ev.root ? <span title={ev.root}>root {short(ev.root)}</span> : null}
      {ev.digest ? <span title={ev.digest}>digest {short(ev.digest)}</span> : null}
      {ev.tx && !ev.tx_link ? <span title={ev.tx}>tx {short(ev.tx)}</span> : null}
      {href ? (
        <ExternalLink href={href}>{ev.tx_link ? "chain tx" : "open"}</ExternalLink>
      ) : null}
    </span>
  );
}

export function ActivityTimeline({
  events,
  lastOid,
  lastOrder,
}: {
  events: EventRow[];
  lastOid?: string;
  lastOrder?: { oid?: string; status?: string; market?: string; side?: string; sz?: number };
}) {
  const newestFirst = [...events].reverse();
  const grouped = new Map<string, EventRow[]>();
  for (const ev of newestFirst) {
    const k = dayLabel(ev.ts);
    grouped.set(k, [...(grouped.get(k) || []), ev]);
  }
  return (
    <section>
      <p className="fine">Every row is a host record from this computer. Historical fills never appear inside a new exact preview.</p>
      {events.length === 0 && !lastOid ? (
        <p className="empty">Empty is honest until this machine records a research, preview, or order.</p>
      ) : null}
      {[...grouped.entries()].map(([day, rows]) => (
        <div key={day}>
          <p className="label">{day}</p>
          <ul className="timeline">
            {rows.map((ev, i) => (
              <li key={`${ev.ts}-${ev.kind}-${i}`}>
                <time dateTime={ev.ts ? new Date(ev.ts).toISOString() : undefined}>
                  {ev.ts ? new Date(ev.ts).toLocaleTimeString() : ""}
                </time>
                <strong>{humanKind(ev.kind)}</strong>
                {ev.market ? ` ${ev.market}` : ""} {humanStatus(ev.status)}
                {ev.reason ? ` · ${ev.reason}` : ""}
                <EventIds ev={ev} />
              </li>
            ))}
          </ul>
        </div>
      ))}
      {lastOid ? (
        <p className="fine">
          Last venue OID on this machine: {lastOid}
          {lastOrder?.status === "filled" ? " · filled position, not a new preview" : ""}
          {lastOrder?.market ? ` · ${lastOrder.market} ${lastOrder.side} ${lastOrder.sz}` : ""}.
        </p>
      ) : null}
    </section>
  );
}
