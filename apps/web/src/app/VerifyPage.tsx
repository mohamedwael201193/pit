import { useState } from "react";
import { PageHead } from "../ui/PageHead";
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
      <div className="mt-8 border border-[rgb(240_231_212/0.25)] p-6">
        <VerifyForm hash={hash} root={root} explorer={explorer} net={net} onHash={setHash} onRoot={setRoot} />
      </div>
    </div>
  );
}
