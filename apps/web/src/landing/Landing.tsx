import { useEffect, useState } from "react";
import { usePrivy } from "@privy-io/react-auth";
import { Hero } from "./Hero";
import { Story } from "./Story";
import { PipelineTrack } from "./PipelineTrack";
import { Capabilities } from "./Capabilities";
import { Faq } from "./Faq";
import { CtaBand } from "./CtaBand";
import { NetworkToggle } from "../NetworkToggle";
import { NetworkBanner } from "../NetworkBanner";
import { IsolateNote } from "../IsolateNote";
import { KillNote } from "../KillNote";
import { TransferNote } from "../TransferNote";
import { NoSession } from "../NoSession";
import { RefreshNote } from "../RefreshNote";
import { SiweBind } from "../SiweBind";
import { PolicyPanel } from "../PolicyPanel";
import { EmptyWatch } from "../EmptyWatch";
import { ProgressStrip } from "../ProgressStrip";
import { VerifyForm } from "../VerifyForm";
import { Ring } from "../Ring";

type Net = "mainnet" | "testnet";

const CAP: Record<Net, string[]> = {
  mainnet: [
    "Production",
    "Direct TeeML on Aristotle",
    "Desk mint / authorize / revoke",
    "Transfer of Agentic ID is not live",
    "Hyperliquid mainnet",
  ],
  testnet: [
    "Full test environment",
    "Galileo chain + Hyperliquid testnet",
    "Experimental transfer only if proven",
    "Different model catalog than production",
  ],
};

export function Landing() {
  const { authenticated } = usePrivy();
  const [net, setNet] = useState<Net>("mainnet");
  const [hash, setHash] = useState("");
  const [root, setRoot] = useState("");
  const explorer = net === "mainnet" ? "https://chainscan.0g.ai" : "https://chainscan-galileo.0g.ai";

  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
  }, []);

  return (
    <div className="guide-shell relative min-h-[100dvh] overflow-x-hidden w-full max-w-full">
      <Hero />
      <Story />
      <PipelineTrack />
      <Capabilities />
      <section className="border-t border-[rgb(240_231_212/0.25)] py-20">
        <div className="container-pit grid gap-12 lg:grid-cols-[1.1fr_0.9fr]">
          <div>
            <p className="text-[11px] tracking-[0.22em] text-[#d82f2f]">YOUR DESK</p>
            <h2 className="mt-3 max-w-3xl text-4xl leading-[0.95] tracking-[-0.04em]">
              Inspect here. Authorize on the machine.
            </h2>
            <SiweBind connected={authenticated} />
            <NoSession />
            <RefreshNote />
            <NetworkToggle net={net} onChange={setNet} />
            <NetworkBanner net={net} />
            <IsolateNote />
            <KillNote />
            <TransferNote />
            <PolicyPanel />
            <EmptyWatch network={net} />
            <ProgressStrip current={authenticated ? "AUTHENTICATING" : "CONNECTING"} />
            <ul className="mt-6 list-disc pl-5">
              {CAP[net].map((line) => (
                <li key={line}>{line}</li>
              ))}
            </ul>
          </div>
          <div id="verify" className="guide-coral px-8 py-12 text-[#0a0a0a]">
            <Ring />
            <p className="mt-10 max-w-[28rem]">
              Web can connect and inspect. Signing Hyperliquid orders happens on desktop or CLI. Session keys never enter
              this browser.
            </p>
            <VerifyForm hash={hash} root={root} explorer={explorer} net={net} onHash={setHash} onRoot={setRoot} />
          </div>
        </div>
      </section>
      <Faq />
      <CtaBand />
    </div>
  );
}
