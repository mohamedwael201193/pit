import * as Accordion from "@radix-ui/react-accordion";
import { Plus } from "@phosphor-icons/react";
import { SectionHeading } from "../ui/SectionHeading";

const QUESTIONS: readonly { q: string; a: string }[] = [
  {
    q: "Who holds my money?",
    a: "You do. PIT never custody. The session agent can order or cancel on Hyperliquid. It cannot withdraw.",
  },
  {
    q: "Where does signing happen?",
    a: "On desktop or CLI. This website can connect, inspect, and verify. It cannot create a session key.",
  },
  {
    q: "What if TeeML fails?",
    a: "PIT stops. There is no Router fallback for the private book.",
  },
  {
    q: "Can PIT trade without me?",
    a: "Not unless you enable Guarded Autonomy on PIT Desktop. Chat, this website, and the model cannot turn it on. Manual still requires AUTHORIZE on one exact preview. Research-only never executes. Duplicate clicks do not send a second order.",
  },
  {
    q: "What is the laboratory?",
    a: "TESTNET is the protocol laboratory. Galileo and Hyperliquid testnet stay on that workspace. Capabilities are not copied from production.",
  },
  {
    q: "How does it learn?",
    a: "After a real outcome, PIT scores Brier and ECE. Until thirty resolved forecasts it says not enough data. It never prints a fake 72 percent.",
  },
  {
    q: "How do I install PIT?",
    a: "Download the Windows installer from GitHub Releases. Verify SHA256 against SHA256SUMS. The installer is not Authenticode-signed until a certificate exists. Pair at /pair with the code shown on this machine.",
  },
  {
    q: "How do I revoke PIT?",
    a: "Kill the session on desktop, then revoke the Hyperliquid agent from your account. Logout deletes the local session. PIT never holds your seed, so there is nothing to export.",
  },
];

export function Faq() {
  return (
    <section id="faq" className="section-kept border-t border-[rgb(240_231_212/0.25)]">
      <div className="container-pit">
        <div className="grid gap-12 lg:grid-cols-[minmax(0,0.72fr)_minmax(0,1.28fr)] lg:gap-20">
          <SectionHeading title="Questions people ask before connecting" className="lg:sticky lg:top-28 lg:self-start" />
          <Accordion.Root type="single" collapsible className="border-t border-[rgb(240_231_212/0.25)]">
            {QUESTIONS.map((item) => (
              <Accordion.Item key={item.q} value={item.q} className="border-b border-[rgb(240_231_212/0.25)]">
                <Accordion.Header>
                  <Accordion.Trigger className="group flex w-full items-start justify-between gap-6 py-5 text-left">
                    <span className="text-[1.125rem] leading-7 font-bold tracking-[-0.02em] text-[var(--guide-cream)] group-hover:underline">
                      {item.q}
                    </span>
                    <Plus
                      size={18}
                      aria-hidden="true"
                      className="mt-1 shrink-0 text-[var(--guide-cream)] transition-transform duration-200 group-data-[state=open]:rotate-45"
                    />
                  </Accordion.Trigger>
                </Accordion.Header>
                <Accordion.Content className="overflow-hidden data-[state=closed]:animate-collapse-up data-[state=open]:animate-collapse-down">
                  <p className="max-w-[56ch] pt-1 pb-6 text-[1rem] leading-7 text-[rgb(240_231_212/0.72)]">{item.a}</p>
                </Accordion.Content>
              </Accordion.Item>
            ))}
          </Accordion.Root>
        </div>
      </div>
    </section>
  );
}
