import { useEffect } from "react";
import { Hero } from "./Hero";
import { Story } from "./Story";
import { PipelineRing } from "./PipelineRing";
import { Moments } from "./Moments";
import { Dual } from "./Dual";
import { LiveTape } from "./LiveTape";
import { Marquee } from "./Marquee";
import { Ledger } from "./Ledger";
import { Faq } from "./Faq";
import { CtaBand } from "./CtaBand";
import { LandingNav } from "./Nav";
import { WatchProvider } from "../public/Watch";
import { ChatPanel } from "../public/ChatPanel";

export function Landing() {
  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
  }, []);

  return (
    <WatchProvider>
      <div className="guide-shell relative min-h-[100dvh] w-full max-w-full overflow-x-hidden">
        <LandingNav />
        <Hero />
        <Story />
        <PipelineRing />
        <Moments />
        <Dual />
        <LiveTape />
        <Marquee />
        <Ledger />
        <Faq />
        <CtaBand />
        <ChatPanel />
      </div>
    </WatchProvider>
  );
}
