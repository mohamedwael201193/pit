import { useEffect } from "react";
import { Hero } from "./Hero";
import { Story } from "./Story";
import { PipelineRing } from "./PipelineRing";
import { Moments } from "./Moments";
import { Dual } from "./Dual";
import { Marquee } from "./Marquee";
import { Ledger } from "./Ledger";
import { Faq } from "./Faq";
import { CtaBand } from "./CtaBand";
import { LandingNav } from "./Nav";

export function Landing() {
  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
  }, []);

  return (
    <div className="guide-shell relative min-h-[100dvh] w-full max-w-full overflow-x-hidden">
      <LandingNav />
      <Hero />
      <Story />
      <PipelineRing />
      <Moments />
      <Dual />
      <Marquee />
      <Ledger />
      <Faq />
      <CtaBand />
    </div>
  );
}
