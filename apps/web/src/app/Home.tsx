import { useState } from "react";
import { Link } from "react-router-dom";
import { EmptyWatch } from "../EmptyWatch";
import { NetworkToggle } from "../NetworkToggle";
import { NetworkBanner } from "../NetworkBanner";
import { IsolateNote } from "../IsolateNote";
import { KillNote } from "../KillNote";
import { TransferNote } from "../TransferNote";
import { NoSession } from "../NoSession";
import { PolicyPanel } from "../PolicyPanel";
import { ProgressStrip } from "../ProgressStrip";

type Net = "mainnet" | "testnet";

export function Home() {
  const [net, setNet] = useState<Net>("mainnet");
  return (
    <div className="mx-auto max-w-[80rem]">
      <p className="text-[11px] tracking-[0.18em] text-[#d82f2f]">WHAT NEEDS YOUR ATTENTION</p>
      <h1 className="mt-2 max-w-4xl text-4xl tracking-[-0.04em] md:text-5xl">Your Watch. Live books only.</h1>
      <p className="mt-4 max-w-[48ch] text-[1.05rem] leading-7 text-[rgb(240_231_212/0.75)]">
        PIT may research and notify. PIT may not sign automatically. Authorize lives on desktop.
      </p>
      <NetworkToggle net={net} onChange={setNet} />
      <NetworkBanner net={net} />
      <EmptyWatch network={net} />
      <PolicyPanel />
      <ProgressStrip current="WAITING_FOR_USER" />
      <NoSession />
      <IsolateNote />
      <KillNote />
      <TransferNote />
      <p className="mt-8">
        <Link to="/app/start" className="pill pill-coral">
          Resume onboarding
        </Link>
      </p>
    </div>
  );
}
