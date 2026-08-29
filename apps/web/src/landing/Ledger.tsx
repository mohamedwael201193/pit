import { Reveal } from "../ui/Reveal";
import { SectionHeading } from "../ui/SectionHeading";

const ROWS = [
  { human: "Direct TeeML on Aristotle", call: "glm-5.2", status: "live" },
  { human: "Galileo sealed ask", call: "VerifyE2EE", status: "unverified" },
  { human: "Transfer of Agentic ID", call: "iTransfer", status: "unavailable" },
  { human: "Hyperliquid order", call: "Desktop AUTHORIZE or host-gated Guarded Autonomy", status: "desktop" },
  { human: "Storage proof", call: "0g-storage-client upload/download --proof", status: "workspace key" },
] as const;

export function Ledger() {
  return (
    <section id="verify" className="section-kept border-t border-[rgb(240_231_212/0.25)]">
      <div className="container-pit">
        <SectionHeading
          title="Every claim has a real limit under it"
          body="Open any of these in the matching explorer after a receipt exists. Nothing here is a screenshot."
        />
        <Reveal className="mt-12 divide-y divide-[rgb(240_231_212/0.25)] border-y border-[rgb(240_231_212/0.25)]">
          {ROWS.map((row) => (
            <div key={row.human} className="grid gap-2 py-5 md:grid-cols-[1.4fr_1fr_auto] md:items-baseline">
              <p className="text-[1.125rem] font-bold tracking-[-0.02em] text-[var(--guide-cream)]">{row.human}</p>
              <p className="font-mono text-[0.875rem] text-[rgb(240_231_212/0.55)]">{row.call}</p>
              <p className="text-[0.8125rem] font-medium tracking-[0.08em] text-[#d82f2f] uppercase">{row.status}</p>
            </div>
          ))}
        </Reveal>
      </div>
    </section>
  );
}
