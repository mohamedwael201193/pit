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
  if (k.includes("mission.enabled") || k === "mission.enabled") return "Mission started";
  if (k.includes("mission.stopped") || k === "mission.stopped") return "Mission stopped";
  if (k.includes("autonomous") || k.includes("guarded")) return "Autonomous action";
  if (k.includes("opportunity")) return "Opportunity found";
  if (k.includes("research.start") || k === "research_started") return "Research started";
  if (k.includes("research") && (k.includes("done") || k.includes("complete"))) return "Research completed";
  if (k.startsWith("research")) return "Research";
  if (k.includes("committee")) return "Committee";
  if (k.includes("preview")) return "Preview";
  if (k.includes("author")) return "Approval";
  if (k.includes("fill")) return "Fill";
  if (k.includes("cancel")) return "Cancel";
  if (k.includes("order")) return "Order";
  if (k.includes("policy")) return "Policy";
  if (k.includes("session")) return "Session";
  if (k.includes("receipt") || k.includes("proof")) return "Proof";
  if (k.includes("position")) return "Position";
  if (k.includes("mission")) return "Mission";
  if (k.includes("automation")) return "Prepared";
  return k.replaceAll(".", " ");
}

function humanStatus(status?: string) {
  const s = String(status || "");
  if (s === "READY_ELIGIBLE") return "eligible preview";
  if (s === "READY_STOOD_DOWN") return "stood down";
  if (s === "filled") return "filled (historical, not a new preview)";
  if (s === "COMMITTEE_INCOMPLETE") return "incomplete";
  return s.replaceAll("_", " ");
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
  const grouped = new Map<string, EventRow[]>();
  for (const ev of events) {
    const k = dayLabel(ev.ts);
    grouped.set(k, [...(grouped.get(k) || []), ev]);
  }
  return (
    <section>
      <p className="fine">Historical fills live here. They never appear inside a new exact preview.</p>
      {events.length === 0 && !lastOid ? (
        <p className="empty">Empty is honest until this machine records a research, preview, or order.</p>
      ) : null}
      {[...grouped.entries()].map(([day, rows]) => (
        <div key={day}>
          <p className="label">{day}</p>
          <ul className="timeline">
            {rows.map((ev, i) => (
              <li key={`${ev.ts}-${i}`}>
                <time>{ev.ts ? new Date(ev.ts).toLocaleTimeString() : ""}</time>
                <strong>{humanKind(ev.kind)}</strong>
                {ev.market ? ` ${ev.market}` : ""} {humanStatus(ev.status)}
                {ev.job_id ? ` · job ${ev.job_id.slice(0, 8)}` : ""}
                {ev.oid ? ` · OID ${ev.oid}` : ""}
                {ev.reason ? ` · ${ev.reason}` : ""}
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
