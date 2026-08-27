import { useState } from "react";
import { usePrivy } from "@privy-io/react-auth";
import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";
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
    <div className="mx-auto flex max-w-[80rem] flex-col gap-8">
      <PageHead title="Account" lede="Wallet, network, kill, and verify. Session keys never live here." />

      <Bezel>
        {addr ? <p className="font-mono text-[0.9375rem] break-all text-[var(--guide-cream)]">{addr}</p> : null}
        <div className="mt-4">
          <SiweBind connected={authenticated} />
        </div>
        <NoSession />
        <RefreshNote />
      </Bezel>

      <Bezel>
        <NetworkToggle net={net} onChange={setNet} />
        <NetworkBanner net={net} />
      </Bezel>

      <Bezel>
        <IsolateNote />
        <KillNote />
        <TransferNote />
      </Bezel>

      <Bezel className="max-w-[36rem]">
        <VerifyForm hash={hash} root={root} explorer={explorer} net={net} onHash={setHash} onRoot={setRoot} />
      </Bezel>
    </div>
  );
}
