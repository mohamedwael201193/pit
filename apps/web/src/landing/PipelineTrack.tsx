import { PIPELINE } from "../diagrams/PipelineCard";

export function PipelineTrack() {
  return (
    <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit">
        <h2 className="max-w-4xl text-4xl leading-[0.95] tracking-[-0.04em] md:text-5xl">
          Market in. Proof out. You in the middle.
        </h2>
        <div className="guide-track mt-10">
          {PIPELINE.map((label) => (
            <article
              key={label}
              className="border border-[rgb(240_231_212/0.28)] bg-[#141414] p-6"
            >
              <p className="font-mono text-[0.75rem] tracking-[0.14em] text-[#d82f2f]">{label}</p>
              <p className="mt-3 text-[0.95rem] leading-6 text-[rgb(240_231_212/0.75)]">{copyFor(label)}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}

function copyFor(label: string): string {
  switch (label) {
    case "MARKET":
      return "Live Hyperliquid public books. Empty Watch is real.";
    case "PRIVATE":
      return "Your thesis stays in the envelope. Not in a chat log.";
    case "SEALED":
      return "Direct TeeML HPKE. Router keys are refused.";
    case "RESEARCH":
      return "Researcher reads the book you attached.";
    case "CHALLENGE":
      return "Challenger attacks the thesis in a separate envelope.";
    case "RISK":
      return "Risk scores what remains after the challenge.";
    case "POLICY":
      return "Your clip, assets, and kill. The model cannot raise them.";
    case "AUTHORIZE":
      return "Desktop or CLI. This browser cannot sign the order.";
    case "EXECUTE":
      return "Order or cancel only. extraAgents must list the session.";
    case "PROVE":
      return "0G Storage with the official Go client --proof.";
    case "LEARN":
      return "Brier and ECE after outcomes. No fake 72 percent.";
    default:
      return "";
  }
}
