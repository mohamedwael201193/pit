import { Link } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";
import { DirectSign } from "../DirectSign";
import { BindDesk } from "../BindDesk";
import { useWatch } from "./Watch";
import { classifyError } from "../namedStates";
import { useState } from "react";
import { OnboardRail } from "./OnboardRail";

export function ProtectPage() {
  const { ready, authenticated, login, user } = usePrivy();
  const { desktop } = useWatch();
  const addr = user?.wallet?.address;
  const paired = Boolean(typeof window !== "undefined" && sessionStorage.getItem("pit_device"));
  const [error, setError] = useState<string | null>(null);

  const connect = async () => {
    setError(null);
    try {
      await login();
    } catch (err) {
      const named = classifyError(err instanceof Error ? err.message : "connect failed");
      setError(`${named.title}. ${named.body} ${named.next}`);
    }
  };

  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Protect my strategy"
        lede="Desktop setup, step two: sign in the bound wallet to link this browser with PIT on this computer. The Direct token stays on the machine. This site never holds it and cannot place an order."
      />
      <OnboardRail current={2} />

      <ol className="mt-10 grid gap-6">
        <Bezel as="li">
          <p className="font-mono text-[0.7rem] tracking-[0.14em] text-[rgb(240_231_212/0.45)]">STEP 1 · PAIR</p>
          <p className="mt-2 text-[1.25rem] font-semibold">Pair this browser with PIT Desktop</p>
          <p className="mt-3 max-w-[52ch] text-[0.975rem] leading-6 text-[rgb(240_231_212/0.7)]">
            {desktop.present
              ? `PIT Desktop is live${desktop.version ? ` · ${desktop.version}` : ""}. Type the one-time code from the machine.`
              : "Launch PIT Desktop on this computer first. Pairing is loopback only."}
          </p>
          <p className="mt-2 text-[0.875rem] text-[rgb(240_231_212/0.55)]">
            {paired ? "This browser already has a pairing token for this computer." : "Not paired yet. Complete step 1 first."}
          </p>
          <Link
            className="mt-6 inline-flex rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4] no-underline"
            to="/pair"
          >
            Open pairing
          </Link>
        </Bezel>

        <Bezel as="li">
          <p className="font-mono text-[0.7rem] tracking-[0.14em] text-[rgb(240_231_212/0.45)]">STEP 2 · SIGN IN WALLET</p>
          <p className="mt-2 text-[1.25rem] font-semibold">Sign in to link the wallet with desktop</p>
          <p className="mt-3 max-w-[52ch] text-[0.975rem] leading-6 text-[rgb(240_231_212/0.7)]">
            Connect the wallet you already use. No seed field exists. After you sign Protect my strategy, PIT stores the
            24-hour sealed-path token on this computer only. Then return to PIT Desktop for Connect Hyperliquid.
          </p>
          {!paired ? (
            <p className="mt-6 text-[0.975rem] text-[#ff7a7a]" role="status">
              Complete pairing first. Connect your wallet stays locked until this browser is paired.
            </p>
          ) : !ready ? (
            <p className="mt-6">Loading wallet connect</p>
          ) : !authenticated || !addr ? (
            <button
              type="button"
              className="mt-6 rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4]"
              onClick={() => void connect()}
            >
              Connect your wallet
            </button>
          ) : (
            <div className="mt-6">
              <p className="text-[0.75rem] tracking-[0.16em] text-[rgb(240_231_212/0.55)]">BOUND WALLET</p>
              <p className="mt-2 font-mono break-all text-[0.9375rem]">{addr}</p>
              <BindDesk network="mainnet" />
              <DirectSign />
            </div>
          )}
          {error ? (
            <p role="alert" className="mt-4 text-[0.875rem] text-[#ff7a7a]">
              {error}
            </p>
          ) : null}
        </Bezel>
      </ol>
    </div>
  );
}
