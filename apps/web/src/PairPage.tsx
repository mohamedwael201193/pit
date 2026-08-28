import { useState } from "react";
import { Link } from "react-router-dom";
import { PitMark } from "./brand/PitMark";

const COMPANION = "http://127.0.0.1:17373";

export function PairPage() {
  const [code, setCode] = useState("");
  const [msg, setMsg] = useState<string | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function pair() {
    setErr(null);
    setMsg(null);
    try {
      const r = await fetch(`${COMPANION}/pair`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code: code.replace(/[^A-Za-z0-9]/g, "") }),
      });
      const text = await r.text();
      let body: { ok?: boolean; canSign?: boolean; sign?: boolean; device?: string } = {};
      try {
        body = JSON.parse(text) as typeof body;
      } catch {
        body = {};
      }
      if (body.sign || body.canSign) {
        setErr("Pairing refused. This site cannot hold a session key.");
        return;
      }
      if (!r.ok) {
        if (text.includes("pairing_expired")) {
          setErr("That code expired. Open PIT on this computer and type the new code shown there.");
          return;
        }
        if (text.includes("pairing_denied")) {
          setErr("Pairing refused. Open PIT on this computer and type the code shown there.");
          return;
        }
        setErr("Pairing refused. Open PIT on this computer and type the code shown there.");
        return;
      }
      if (body.device) {
        sessionStorage.setItem("pit_device", body.device);
      }
      setMsg("PAIRING COMPLETE. This browser can view. It cannot sign and cannot hold a session key. Next: Protect your private research.");
    } catch {
      setErr(
        "PIT is not reachable on this computer. Launch the Windows app first. If Chrome asks to access other apps on this device, choose Allow. PIT only uses 127.0.0.1.",
      );
    }
  }

  return (
    <div className="guide-shell min-h-[100dvh] bg-[#0a0a0a] px-6 py-16 text-[#f0e7d4]">
      <div className="mx-auto max-w-[40rem]">
        <PitMark />
        <p className="mt-10 font-mono text-[0.75rem] tracking-[0.14em] text-[#d82f2f]">LOCAL PAIRING</p>
        <h1 className="mt-3 text-[2.5rem] font-semibold tracking-[-0.04em]">Pair this browser with PIT on this machine.</h1>
        <p className="mt-4 max-w-[42ch] text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.7)]">
          PIT never asks for a seed phrase. The one-time code lives on your desktop. This site never receives your
          session key. After pairing, sign Protect my strategy from the bound wallet.
        </p>
        <ol className="mt-8 grid gap-3 text-[0.975rem] leading-6 text-[rgb(240_231_212/0.75)]">
          <li>1. Open PIT Desktop on this computer. A pairing code appears there.</li>
          <li>2. Type the code shown in PIT Desktop. It expires in two minutes and works once.</li>
          <li>3. If Chrome asks to access other apps on this device, choose Allow. That is loopback only.</li>
          <li>4. After pairing, this browser can view. It cannot sign orders or hold a session key.</li>
        </ol>
        <p className="mt-6 font-semibold text-[#f0e7d4]">Open PIT Desktop</p>
        <p className="mt-1 text-[0.975rem] leading-6 text-[rgb(240_231_212/0.75)]">Enter the code shown there.</p>
        <div className="mt-6 flex flex-wrap gap-3">
          <a
            className="rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4]"
            href="https://github.com/mohamedwael201193/pit/releases/latest"
          >
            Download PIT
          </a>
          <a
            className="rounded-full border border-[rgb(240_231_212/0.35)] px-6 py-3 font-semibold text-[#f0e7d4]"
            href="http://127.0.0.1:17373/health"
          >
            Open PIT Desktop
          </a>
        </div>
        <label className="mt-10 block">
          <span className="text-[0.75rem] tracking-[0.12em] text-[rgb(240_231_212/0.5)]">ONE-TIME CODE</span>
          <input
            className="mt-2 w-full border border-[rgb(240_231_212/0.25)] bg-[#141414] px-4 py-3 font-mono text-[1.25rem] tracking-[0.2em] outline-none focus:border-[#d82f2f]"
            autoComplete="off"
            spellCheck={false}
            value={code}
            onChange={(e) => setCode(e.target.value.toUpperCase())}
            placeholder="ABCD-EFGH"
            aria-label="pairing code"
          />
        </label>
        <button
          type="button"
          className="mt-6 rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4]"
          onClick={() => void pair()}
        >
          Pair this browser
        </button>
        {msg ? (
          <div className="mt-6">
            <p className="text-[0.975rem] text-[#f0e7d4]">{msg}</p>
            <p className="mt-2 text-[0.975rem] leading-6 text-[rgb(240_231_212/0.75)]">
              Your browser is read-only. The private authorization stays on this computer.
            </p>
            <Link className="mt-4 inline-block rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4]" to="/app">
              Protect my strategy
            </Link>
          </div>
        ) : null}
        {err ? (
          <p className="mt-4 text-[0.975rem] text-[#ff7a7a]" role="alert">
            {err}
          </p>
        ) : null}
        <p className="mt-10 text-[0.875rem] text-[rgb(240_231_212/0.5)]">
          macOS and Linux installers are not claimed until they are packaged and tested. Source build is documented in
          the README.
        </p>
        <Link className="mt-6 inline-block text-[#d82f2f]" to="/app">
          Protect private research
        </Link>
        <p className="mt-4 text-[0.875rem] text-[rgb(240_231_212/0.5)]">
          Next: Protect my strategy, then Connect Hyperliquid. If PIT Desktop is not open, launch it on this computer first.
        </p>
      </div>
    </div>
  );
}
