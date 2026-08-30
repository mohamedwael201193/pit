import { Link, useParams } from "react-router-dom";
import { HISTORICAL_FILL } from "./facts";
import type { EvidenceKind } from "./types";

type Beat = { title: string; kind: EvidenceKind; detail: string };

const HISTORICAL_BEATS: Beat[] = [
  { title: "Market discovered", kind: "HISTORICAL", detail: `${HISTORICAL_FILL.market} on Hyperliquid mainnet.` },
  { title: "Research started", kind: "HISTORICAL", detail: "Private Direct TeeML. Prompt not published." },
  { title: "Roles completed", kind: "HISTORICAL", detail: "Committee ran on desktop. Transcript sealed." },
  { title: "Host decision", kind: "HISTORICAL", detail: "Host allowed a sized preview. Model cannot raise clip." },
  { title: "Policy", kind: "HISTORICAL", detail: "Pinned policy on the machine. Not readable as a private blob here." },
  { title: "Execution", kind: "HISTORICAL", detail: "Desktop AUTHORIZE. This website cannot sign." },
  { title: "OID", kind: "HISTORICAL", detail: HISTORICAL_FILL.oid },
  { title: "Fill", kind: "HISTORICAL", detail: `${HISTORICAL_FILL.sz} @ ${HISTORICAL_FILL.px}` },
  { title: "Reconciliation", kind: "ABSENT", detail: "No public reconciliation object on this site." },
  { title: "0G proof", kind: "ABSENT", detail: "No public storage root or chain tx is attached to this replay." },
];

export function ReplayPage() {
  const { id = "" } = useParams();
  if (id !== HISTORICAL_FILL.id) {
    return (
      <div>
        <p className="intel-kicker">Mission replay</p>
        <h1 className="intel-title mt-2">{id || "Unknown"}</h1>
        <p className="intel-lede">
          No public-safe mission with that id. PIT will not invent a timeline. If a desktop run ended with no trade,
          that is a successful stand-down, not a system failure — the public reason lives on the desktop receipt.
        </p>
        <Link to="/missions" className="intel-ghost mt-6 inline-flex">
          Back to missions
        </Link>
      </div>
    );
  }

  return (
    <div>
      <Link to="/missions" className="intel-ghost">
        ← Missions
      </Link>
      <p className="intel-kicker mt-6">Replay · {HISTORICAL_FILL.kind}</p>
      <h1 className="intel-title mt-2">
        {HISTORICAL_FILL.market} {HISTORICAL_FILL.oid}
      </h1>
      <p className="intel-lede">{HISTORICAL_FILL.note}</p>
      <ol className="mt-10 divide-y divide-[rgb(240_231_212/0.12)] border-y border-[rgb(240_231_212/0.12)]">
        {HISTORICAL_BEATS.map((b) => (
          <li key={b.title} className="grid gap-1 py-4 sm:grid-cols-[1fr_7rem_2fr] sm:items-baseline">
            <p className="font-medium">{b.title}</p>
            <p className={`text-[0.6875rem] tracking-[0.14em] ${kindColor(b.kind)}`}>{b.kind}</p>
            <p className="text-[0.875rem] leading-6 text-[rgb(240_231_212/0.65)]">{b.detail}</p>
          </li>
        ))}
      </ol>
    </div>
  );
}

function kindColor(k: EvidenceKind): string {
  if (k === "LIVE") return "text-[#7dffb3]";
  if (k === "ABSENT") return "text-[rgb(240_231_212/0.4)]";
  return "text-[#d82f2f]";
}
