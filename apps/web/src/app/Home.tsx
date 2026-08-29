import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";
import { ChoiceCard } from "../ui/ChoiceCard";
import { EmptyWatch } from "../EmptyWatch";
import { NetworkToggle } from "../NetworkToggle";
import { NetworkBanner } from "../NetworkBanner";
import { BindDesk } from "../BindDesk";
import { DirectSign } from "../DirectSign";
import { DiagramPolicy, DiagramSession } from "../diagrams/pitGuide";

type Net = "mainnet" | "testnet";

export function Home() {
  const { user } = usePrivy();
  const navigate = useNavigate();
  const [net, setNet] = useState<Net>("mainnet");
  const addr = user?.wallet?.address;

  return (
    <div className="mx-auto flex max-w-[80rem] flex-col gap-10">
      <PageHead
        title="Protect private research"
        lede="This browser can pair, sign Protect, and show public books. Authorize and session keys stay on PIT Desktop."
      />

      <Bezel>
        <p className="text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.78)]">
          Connected as{" "}
          <span className="font-mono text-[var(--guide-cream)]">
            {addr ? `${addr.slice(0, 8)}...${addr.slice(-4)}` : "wallet"}
          </span>
        </p>
        <NetworkToggle net={net} onChange={setNet} />
        <NetworkBanner net={net} />
        <BindDesk network={net} />
        <DirectSign />
      </Bezel>

      <section aria-labelledby="watch-heading">
        <h2 id="watch-heading" className="sr-only">
          Watch
        </h2>
        <EmptyWatch network={net} />
      </section>

      <div className="grid gap-4 sm:grid-cols-2">
        <ChoiceCard
          title="Finish setup"
          body="Twelve beats. Wallet, network, policy, then a session on the machine."
          Diagram={DiagramSession}
          onClick={() => navigate("/app/start")}
        />
        <ChoiceCard
          title="Pair this computer"
          body="One-time code. The browser never holds a session key."
          Diagram={DiagramSession}
          onClick={() => navigate("/pair")}
        />
        <ChoiceCard
          title="Read the law"
          body="Clip, assets, kill. The model cannot raise them."
          Diagram={DiagramPolicy}
          onClick={() => navigate("/app/policy")}
        />
      </div>
    </div>
  );
}
