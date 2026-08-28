import { useState } from "react";
import { usePrivy } from "@privy-io/react-auth";

type PrivySign = {
  user?: { wallet?: { address?: string } };
  signMessage?: (message: unknown, opts?: unknown) => Promise<unknown>;
};

const COMPANION = "http://127.0.0.1:17373";

type Challenge = {
  ok?: boolean;
  message?: string;
  digest?: string;
  provider?: string;
  model?: string;
  explain?: string;
  expiresAt?: number;
  error?: string;
  sign?: boolean;
  trade?: boolean;
};

function deviceToken() {
  return sessionStorage.getItem("pit_device") || "";
}

async function signRawDigest(digest: string, address: string): Promise<string> {
  const eth = (window as unknown as { ethereum?: { request: (args: { method: string; params: string[] }) => Promise<string> } }).ethereum;
  if (eth?.request) {
    return eth.request({ method: "personal_sign", params: [digest, address] });
  }
  throw new Error("SIGNATURE_DECLINED");
}

export function DirectSign() {
  const privy = usePrivy() as unknown as PrivySign;
  const addr = privy.user?.wallet?.address || "";
  const signMessage = privy.signMessage;
  const [msg, setMsg] = useState<string | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function run() {
    setErr(null);
    setMsg(null);
    const device = deviceToken();
    if (!device) {
      setErr("Pair this browser first. Open /pair and type the code from the desktop.");
      return;
    }
    if (!addr) {
      setErr("Connect your wallet first. PIT never asks for a seed phrase.");
      return;
    }
    setBusy(true);
    try {
      const intentRes = await fetch(`${COMPANION}/direct/intent`, {
        method: "POST",
        headers: { Authorization: `Bearer ${device}` },
      });
      const intent = (await intentRes.json()) as Challenge;
      if (intent.sign || intent.trade) {
        setErr("This site cannot hold a Direct token.");
        return;
      }
      if (intentRes.status === 401) {
        setErr("Pair this browser first. Open /pair and type the code from the desktop.");
        return;
      }
      if (!intentRes.ok || !intent.digest || !intent.message) {
        setErr(intent.error === "unbound" ? "Bind your wallet on this computer first." : intent.error || "Direct challenge refused.");
        return;
      }
      let signature = "";
      try {
        signature = await signRawDigest(intent.digest, addr);
      } catch {
        signature = "";
      }
      if (!signature && signMessage) {
        try {
          const out = await signMessage({ message: { raw: intent.digest } });
          if (typeof out === "string") signature = out;
          else if (out && typeof out === "object" && "signature" in out) signature = String((out as { signature: string }).signature);
        } catch {
          signature = "";
        }
      }
      const doneRes = await fetch(`${COMPANION}/direct/complete`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${device}`,
        },
        body: JSON.stringify({ signature }),
      });
      const done = (await doneRes.json()) as { ok?: boolean; error?: string; sign?: boolean; expiresAt?: number };
      if (done.sign) {
        setErr("This site cannot hold a Direct token.");
        return;
      }
      if (!doneRes.ok) {
        setErr(done.error || "Signature was not accepted. No token was stored in the browser.");
        return;
      }
      setMsg("Sealed-path authorization is on this computer. This site never received the token.");
    } catch {
      setErr("PIT is not reachable on this computer, or the signature was declined. No token left the wallet into this page.");
    } finally {
      setBusy(false);
    }
  }

  if (!addr) return null;

  return (
    <div className="mt-6 max-w-[46ch]">
      <p className="text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.78)]">
        PIT sends your private strategy only through 0G’s verified sealed path. This signature lasts 24 hours. It cannot
        withdraw. It cannot place a Hyperliquid order. Direct credit lives at pc.0g.ai — switch to Advanced. That is
        provider credit, not a Hyperliquid balance.
      </p>
      <button
        type="button"
        className="mt-4 rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4] disabled:opacity-50"
        disabled={busy}
        onClick={() => void run()}
      >
        {busy ? "Waiting for signature…" : "Protect my strategy"}
      </button>
      {msg ? <p className="mt-4 text-[0.975rem] text-[#f0e7d4]">{msg}</p> : null}
      {err ? (
        <p className="mt-4 text-[0.975rem] text-[#ff7a7a]" role="alert">
          {err}
        </p>
      ) : null}
    </div>
  );
}
