import { useEffect, useState } from "react";
import { fetchProofs, type ProofsView } from "./companion";
import { ExternalLink } from "./ExternalLink";

function shortHash(v?: string) {
  const s = String(v || "");
  if (s.length <= 16) return s;
  return `${s.slice(0, 10)}...${s.slice(-4)}`;
}

function kindLabel(kind?: string) {
  switch (String(kind || "")) {
    case "research":
      return "research verdict";
    case "order":
      return "venue order";
    case "cancel":
      return "venue cancel";
    default:
      return "record";
  }
}

export function EvidenceStrip({ onOpen }: { onOpen: () => void }) {
  const [view, setView] = useState<ProofsView | null>(null);

  useEffect(() => {
    let alive = true;
    const load = async () => {
      const got = await fetchProofs();
      if (alive) setView(got);
    };
    void load();
    const t = window.setInterval(() => void load(), 20000);
    return () => {
      alive = false;
      window.clearInterval(t);
    };
  }, []);

  if (!view) return null;
  const latest = (view.receipts || [])[0];

  return (
    <section className="evidence-strip">
      <div className="evidence-lede">
        <p className="label">0G evidence</p>
        {latest ? (
          <p>
            Last {kindLabel(latest.kind)} published to 0G Storage as{" "}
            <span className="hash" title={latest.root}>
              {shortHash(latest.root)}
            </span>
            {latest.market ? ` for ${latest.market}` : ""}.
          </p>
        ) : view.ready ? (
          <p>Nothing published yet. The next research verdict lands on 0G Storage with a chain transaction.</p>
        ) : (
          <p>
            Publishing is off: {String(view.blocked || "unknown").replaceAll("_", " ")}. Bind the payer with{" "}
            <span className="hash">pit evidence bind-payer</span>.
          </p>
        )}
      </div>
      <div className="evidence-act">
        {latest?.tx_link ? (
          <ExternalLink href={latest.tx_link}>Chain transaction</ExternalLink>
        ) : null}
        <button type="button" onClick={onOpen}>
          {view.count ? `Proof trail (${view.count})` : "Proof trail"}
        </button>
      </div>
    </section>
  );
}
