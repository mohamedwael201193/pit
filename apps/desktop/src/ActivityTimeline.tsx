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

export function ActivityTimeline({ events, lastOid }: { events: EventRow[]; lastOid?: string }) {
  const grouped = new Map<string, EventRow[]>();
  for (const ev of events) {
    const k = dayLabel(ev.ts);
    grouped.set(k, [...(grouped.get(k) || []), ev]);
  }
  return (
    <article className="card">
      <p className="label">TIMELINE</p>
      <p className="fine">Historical fills live here. They never appear inside a new exact preview.</p>
      {events.length === 0 && !lastOid ? (
        <p>Empty is honest until this machine records a research, preview, or order.</p>
      ) : null}
      {[...grouped.entries()].map(([day, rows]) => (
        <div key={day}>
          <p className="label">{day}</p>
          <ul className="timeline">
            {rows.map((ev, i) => (
              <li key={`${ev.ts}-${i}`}>
                <strong>{ev.kind || ev.action || "event"}</strong>
                {ev.market ? ` ${ev.market}` : ""} {ev.status || ""}
                {ev.job_id ? ` · job ${ev.job_id}` : ""}
                {ev.oid ? ` · OID ${ev.oid}` : ""}
                {ev.reason ? ` · ${ev.reason}` : ""}
              </li>
            ))}
          </ul>
        </div>
      ))}
      {lastOid ? <p className="fine">Last venue OID on this machine: {lastOid}. Not a new preview.</p> : null}
    </article>
  );
}
