import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { PageHead } from "./ui/PageHead";
import { Bezel } from "./ui/Surface";
import { windowsInstallerUrl } from "./public/facts";

const COMPANION = "http://127.0.0.1:17373";

export function PairPage() {
  const navigate = useNavigate();
  const [code, setCode] = useState("");
  const [msg, setMsg] = useState<string | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [desk, setDesk] = useState("Checking this computer…");
  const [deskOk, setDeskOk] = useState(false);
  const [busy, setBusy] = useState(false);

  async function probeDesktop() {
    try {
      const r = await fetch(`${COMPANION}/health`);
      const body = (await r.json()) as { ok?: boolean; version?: string; pairing?: boolean };
      if (body.ok) {
        setDeskOk(true);
        setDesk(`PIT Desktop is live on this computer · ${body.version || "connected"}. Type the code shown there.`);
        return;
      }
    } catch {
      /* fall through */
    }
    setDeskOk(false);
    setDesk("PIT Desktop is not reachable on this computer. Launch the Windows app first. This page will not leave pairing.");
  }

  useEffect(() => {
    void probeDesktop();
    const t = window.setInterval(() => void probeDesktop(), 4000);
    return () => window.clearInterval(t);
  }, []);

  async function pair() {
    setErr(null);
    setMsg(null);
    setBusy(true);
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
          setErr("That code expired. Open PIT Desktop and regenerate, then type the new code.");
          return;
        }
        if (text.includes("pairing_denied")) {
          setErr("Pairing refused. Open PIT Desktop and type the code shown there.");
          return;
        }
        setErr("Pairing refused. Open PIT Desktop and type the code shown there.");
        return;
      }
      if (body.device) {
        sessionStorage.setItem("pit_device", body.device);
      }
      setMsg("PAIRING COMPLETE. This browser can view. It cannot sign and cannot hold a session key.");
      window.setTimeout(() => navigate("/protect"), 600);
    } catch {
      setErr(
        "PIT is not reachable on this computer. Launch the Windows app first. If Chrome asks to access other apps on this device, choose Allow. PIT only uses 127.0.0.1.",
      );
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="mx-auto grid max-w-[80rem] gap-10 lg:grid-cols-[minmax(0,1fr)_20rem]">
      <div>
        <PageHead
          title="Pair this browser with PIT on this machine."
          lede="PIT never asks for a seed phrase. The one-time code lives on your desktop. This site never receives your session key. After pairing, sign Protect my strategy from the bound wallet."
        />
        <p className="mt-6 text-[0.875rem] text-[rgb(240_231_212/0.5)]">
          Pairing is a late step. Explore radar and proof first.
        </p>
        <ol className="mt-8 grid gap-3 text-[0.975rem] leading-6 text-[rgb(240_231_212/0.75)]">
          <li>1. Open PIT Desktop on this computer. A pairing code appears there with an expiry.</li>
          <li>2. Type the code shown in PIT Desktop. It expires in two minutes and works once. Use Regenerate on desktop if it expired.</li>
          <li>3. If Chrome asks to access other apps on this device, choose Allow. That is loopback only.</li>
          <li>4. After pairing, this browser can view. It cannot sign orders or hold a session key.</li>
        </ol>
        <p className={`mt-6 text-[0.975rem] ${deskOk ? "text-[#7dffb3]" : "text-[#ff7a7a]"}`} role="status">
          {desk}
        </p>
        <div className="mt-6 flex flex-wrap gap-3">
          {deskOk ? (
            <button
              type="button"
              className="rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4]"
              onClick={() => document.querySelector<HTMLInputElement>('input[aria-label="pairing code"]')?.focus()}
            >
              Pair this browser
            </button>
          ) : (
            <a
              className="rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4]"
              href={windowsInstallerUrl()}
            >
              Download PIT
            </a>
          )}
          <button
            type="button"
            className="rounded-full border border-[rgb(240_231_212/0.35)] px-6 py-3 font-semibold text-[#f0e7d4]"
            onClick={() => void probeDesktop()}
          >
            Open PIT Desktop
          </button>
          {deskOk ? (
            <a
              className="rounded-full border border-[rgb(240_231_212/0.35)] px-6 py-3 font-semibold text-[#f0e7d4]"
              href={windowsInstallerUrl()}
            >
              Download PIT
            </a>
          ) : null}
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
          className="mt-6 rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4] disabled:opacity-50"
          disabled={busy}
          onClick={() => void pair()}
        >
          {busy ? "Pairing…" : "Pair this browser"}
        </button>
        {msg ? (
          <div className="mt-6">
            <p className="text-[0.975rem] text-[#f0e7d4]">{msg}</p>
            <p className="mt-2 text-[0.975rem] leading-6 text-[rgb(240_231_212/0.75)]">
              Your browser is read-only. The private authorization stays on this computer.
            </p>
            <Link className="mt-4 inline-block rounded-full bg-[#d82f2f] px-6 py-3 font-semibold text-[#f0e7d4]" to="/protect">
              Protect my strategy
            </Link>
          </div>
        ) : (
          <p className="mt-6 text-[0.875rem] text-[rgb(240_231_212/0.5)]">
            After pairing, open{" "}
            <Link className="text-[#d82f2f]" to="/protect">
              Protect my strategy
            </Link>
            .
          </p>
        )}
        {err ? (
          <p className="mt-4 text-[0.975rem] text-[#ff7a7a]" role="alert">
            {err}
          </p>
        ) : null}
      </div>
      <Bezel>
        <p className="font-mono text-[0.7rem] tracking-[0.14em] text-[rgb(240_231_212/0.45)]">WHAT THIS SITE CANNOT DO</p>
        <ul className="mt-4 grid gap-2 text-[0.95rem] leading-6 text-[rgb(240_231_212/0.7)]">
          <li>Sign orders</li>
          <li>Hold a session key</li>
          <li>Authorize a trade</li>
          <li>Change policy</li>
        </ul>
        <p className="mt-8 text-[0.875rem] text-[rgb(240_231_212/0.5)]">
          macOS and Linux installers are not claimed until they are packaged and tested. Source build is documented in
          the README.
        </p>
        <Link className="mt-6 inline-block text-[#d82f2f]" to="/protect">
          Protect my strategy
        </Link>
      </Bezel>
    </div>
  );
}
