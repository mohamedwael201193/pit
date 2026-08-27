import { useEffect, useState } from "react";
import { usePrivy } from "@privy-io/react-auth";

const COMPANION = "http://127.0.0.1:17373";

export function BindDesk({ network }: { network: "mainnet" | "testnet" }) {
  const { user } = usePrivy();
  const addr = user?.wallet?.address || "";
  const [msg, setMsg] = useState<string | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    const device = sessionStorage.getItem("pit_device");
    if (!device || !addr) {
      return;
    }
    let gone = false;
    fetch(`${COMPANION}/bind`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${device}`,
      },
      body: JSON.stringify({ wallet: addr, network }),
    })
      .then(async (r) => {
        const text = await r.text();
        let body: { ok?: boolean; sign?: boolean; trade?: boolean; error?: string; workspace?: string } = {};
        try {
          body = JSON.parse(text) as typeof body;
        } catch {
          body = {};
        }
        if (gone) return;
        if (body.sign || body.trade) {
          setErr("Bind refused. This site cannot hold a session key.");
          return;
        }
        if (r.status === 401) {
          setErr("Pair this browser first. Open /pair and type the code from the desktop.");
          return;
        }
        if (body.error === "workspace_owned") {
          setErr("This computer already belongs to another wallet. User B cannot overwrite User A.");
          return;
        }
        if (body.error === "network_switch_denied") {
          setErr("This workspace is already bound to the other network.");
          return;
        }
        if (!r.ok) {
          setErr("Bind refused. Launch PIT on this computer.");
          return;
        }
        setErr(null);
        setMsg(
          body.workspace
            ? `This browser sent your public address to the machine. Workspace bound. The site still cannot authorize.`
            : "Address sent to the local machine. The site still cannot authorize.",
        );
      })
      .catch(() => {
        if (!gone) {
          setErr("PIT is not reachable on this computer. Launch the Windows app, then pair.");
        }
      });
    return () => {
      gone = true;
    };
  }, [addr, network]);

  if (!addr) return null;

  return (
    <p className="mt-4 max-w-[46ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.72)]" role="status">
      {err ? <span className="text-[#ff7a7a]">{err}</span> : null}
      {!err && msg ? msg : null}
      {!err && !msg
        ? "After you pair, this page sends only your public address to the machine. No session key leaves the desktop."
        : null}
    </p>
  );
}
