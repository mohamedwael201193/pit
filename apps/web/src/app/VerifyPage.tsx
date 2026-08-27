import { useState } from "react";
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
      <h1 className="text-4xl tracking-[-0.04em]">Verify a receipt on the matching explorer.</h1>
      <p className="mt-4 max-w-[48ch] text-[rgb(240_231_212/0.75)]">
        Recompute from chain and storage proof. Not from a screenshot.
      </p>
      <NetworkToggle net={net} onChange={setNet} />
      <VerifyForm hash={hash} root={root} explorer={explorer} net={net} onHash={setHash} onRoot={setRoot} />
    </div>
  );
}
