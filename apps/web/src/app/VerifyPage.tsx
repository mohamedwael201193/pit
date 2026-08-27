import { useState } from "react";
import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";
import { NetworkToggle } from "../NetworkToggle";
import { VerifyForm } from "../VerifyForm";

type Net = "mainnet" | "testnet";

export function VerifyPage() {
  const [net, setNet] = useState<Net>("mainnet");
  const [hash, setHash] = useState("");
  const [root, setRoot] = useState("");
  const explorer = net === "mainnet" ? "https://chainscan.0g.ai" : "https://chainscan-galileo.0g.ai";
  return (
    <div className="mx-auto max-w-[40rem]">
      <PageHead
        title="Verify a receipt"
        lede="Recompute from chain and storage proof. Not from a screenshot."
      />
      <NetworkToggle net={net} onChange={setNet} />
      <Bezel className="mt-8">
        <VerifyForm hash={hash} root={root} explorer={explorer} net={net} onHash={setHash} onRoot={setRoot} />
      </Bezel>
    </div>
  );
}
