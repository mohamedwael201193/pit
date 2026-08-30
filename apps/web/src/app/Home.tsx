import { Link } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { PageHead } from "../ui/PageHead";
import { ButtonLink } from "../ui/Button";
import { WatchHome } from "../WatchHome";

export function Home() {
  const { user } = usePrivy();
  const addr = user?.wallet?.address;

  return (
    <div className="mx-auto flex max-w-[80rem] flex-col gap-10">
      <PageHead
        title="Observe here. Act on the machine."
        lede="This browser can inspect public books and verify receipts. It cannot pin policy, arm a Sleep Mission, or hold a session key."
      />

      <p className="font-mono text-[0.875rem] text-[rgb(240_231_212/0.6)]">
        Connected as {addr ? `${addr.slice(0, 8)}…${addr.slice(-4)}` : "wallet"} · identity only
      </p>

      <div className="flex flex-wrap gap-2">
        <ButtonLink as={Link} to="/radar" size="lg">
          Live radar
        </ButtonLink>
        <ButtonLink as={Link} to="/proof" variant="secondary" size="lg">
          Proof
        </ButtonLink>
        <ButtonLink as={Link} to="/download" variant="secondary" size="lg">
          Download desktop
        </ButtonLink>
      </div>

      <section>
        <WatchHome network="mainnet" />
      </section>

      <p className="text-[0.875rem] text-[rgb(240_231_212/0.5)]">
        Pairing is a late step, after you have seen live intelligence.{" "}
        <Link to="/pair" className="text-[#d82f2f]">
          Pair this computer
        </Link>
        {" · "}
        <Link to="/protect" className="text-[#d82f2f]">
          Protect my strategy
        </Link>
        {" · "}
        <Link to="/app/policy" className="text-[#d82f2f]">
          Read the law
        </Link>
      </p>
    </div>
  );
}
