import { useState } from "react";
import { usePrivy } from "@privy-io/react-auth";
import { PageHead } from "../ui/PageHead";
import { NetworkToggle } from "../NetworkToggle";
import { NetworkBanner } from "../NetworkBanner";
import { IsolateNote } from "../IsolateNote";
import { KillNote } from "../KillNote";
import { TransferNote } from "../TransferNote";
import { NoSession } from "../NoSession";
import { RefreshNote } from "../RefreshNote";
import { SiweBind } from "../SiweBind";
import { VerifyForm } from "../VerifyForm";

type Net = "mainnet" | "testnet";

export function AccountPage() {
  const { authenticated, user } = usePrivy();
  const [net, setNet] = useState<Net>("mainnet");
  const [hash, setHash] = useState("");
  const [root, setRoot] = useState("");
  const explorer = net === "mainnet" ? "https://chainscan.0g.ai" : "https://chainscan-galileo.0g.ai";
  const addr = user?.wallet?.address;

  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead title="Account" lede="Wallet, network, kill, and verify. Session keys never live here." />
      {addr ? <p className="mt-6 font-mono text-[0.9375rem] break-all text-[rgb(240_231_212/0.75)]">{addr}</p> : null}
      <SiweBind connected={authenticated} />
      <NetworkToggle net={net} onChange={setNet} />
      <NetworkBanner net={net} />
      <NoSession />
      <RefreshNote />
      <IsolateNote />
      <KillNote />
      <TransferNote />
      <div className="mt-10 max-w-[36rem] border border-[rgb(240_231_212/0.25)] p-6">
        <VerifyForm hash={hash} root={root} explorer={explorer} net={net} onHash={setHash} onRoot={setRoot} />
      </div>
    </div>
  );
}
