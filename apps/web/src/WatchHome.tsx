import { useEffect, useState } from "react";
import { Attention } from "./Attention";
import { Bezel } from "./ui/Surface";
import { DiagramEmptyWatch } from "./diagrams/pitGuide";
import { namedState } from "./namedStates";

type Coin = {
  coin: string;
  reason: string;
  mark: number;
  eligible?: boolean;
  executionFeasible?: boolean;
  why?: string;
};

type WatchPayload = {
  count?: number;
  scanned?: number;
  copy?: string;
  trade?: boolean;
  sign?: boolean;
  coins?: Coin[];
};

function markLabel(n: number) {
  if (!Number.isFinite(n) || n <= 0) return "—";
  if (n >= 1000) return n.toLocaleString(undefined, { maximumFractionDigits: 1 });
  if (n >= 1) return n.toLocaleString(undefined, { maximumFractionDigits: 4 });
  return n.toLocaleString(undefined, { maximumFractionDigits: 6 });
}

export function WatchHome({ network = "mainnet" }: { network?: "mainnet" | "testnet" }) {
  const base = import.meta.env.VITE_HEALTH_URL;
  const [scanned, setScanned] = useState(0);
  const [count, setCount] = useState(0);
  const [coins, setCoins] = useState<Coin[]>([]);
  const [fail, setFail] = useState<string | null>(null);

  useEffect(() => {
    if (!base) {
      setScanned(0);
      setCount(0);
      setCoins([]);
      return;
    }
    let gone = false;
    fetch(`${base.replace(/\/$/, "")}/watch?network=${network}`)
      .then((r) => r.json() as Promise<WatchPayload>)
      .then((body) => {
        if (gone) return;
        if (body.trade || body.sign) {
          setScanned(0);
          setCount(0);
          setCoins([]);
          return;
        }
        const rows = Array.isArray(body.coins) ? body.coins : [];
        const pass = rows.filter((c) => c.eligible).slice(0, 8);
        setScanned(typeof body.scanned === "number" ? body.scanned : rows.length);
        setCount(typeof body.count === "number" ? body.count : pass.length);
        setCoins(pass);
        setFail(null);
      })
      .catch(() => {
        if (!gone) {
          setScanned(0);
          setCount(0);
          setCoins([]);
          setFail(namedState("BACKEND_UNREACHABLE").body);
        }
      });
    return () => {
      gone = true;
    };
  }, [base, network]);

  return (
    <div>
      <Attention scanned={scanned} count={count} />
      {fail ? (
        <p role="alert" className="mt-4 max-w-[48ch] text-[0.9375rem] text-[#ff7a7a]">
          {fail} Empty Watch is the honest state.
        </p>
      ) : null}
      {coins.length > 0 ? (
        <div className="mt-6 overflow-hidden border-t border-[rgb(240_231_212/0.18)]">
          <table className="w-full text-left text-[0.9375rem]">
            <caption className="sr-only">Policy-eligible public books. This site cannot trade.</caption>
            <thead>
              <tr className="border-b border-[rgb(240_231_212/0.12)] font-mono text-[0.7rem] tracking-[0.12em] text-[rgb(240_231_212/0.45)]">
                <th className="py-3 pr-4 font-medium">Asset</th>
                <th className="py-3 pr-4 font-medium">Mark</th>
                <th className="py-3 pr-4 font-medium">Policy</th>
                <th className="hidden py-3 font-medium sm:table-cell">Why it is listed</th>
              </tr>
            </thead>
            <tbody>
              {coins.map((c) => (
                <tr key={c.coin} className="border-b border-[rgb(240_231_212/0.08)]">
                  <td className="py-3 pr-4 font-semibold tracking-[-0.02em]">{c.coin}</td>
                  <td className="py-3 pr-4 font-mono text-[rgb(240_231_212/0.85)]">{markLabel(c.mark)}</td>
                  <td className="py-3 pr-4 text-[rgb(240_231_212/0.7)]">Eligible · not executable here</td>
                  <td className="hidden py-3 text-[rgb(240_231_212/0.6)] sm:table-cell">{c.why || c.reason}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <p className="mt-4 max-w-[56ch] text-[0.875rem] text-[rgb(240_231_212/0.5)]">
            Side, size, and AUTHORIZE are not decided on this website. Open PIT Desktop to research against your policy and actual buying power.
          </p>
        </div>
      ) : (
        <Bezel className="mt-6">
          <DiagramEmptyWatch className="aspect-[21/9] w-full border border-[rgb(240_231_212/0.2)]" />
          <p className="mt-6 max-w-[42ch] text-[0.9375rem] text-[rgb(240_231_212/0.65)]">
            Live books only. PIT will not invent a card to fill the home.
          </p>
        </Bezel>
      )}
    </div>
  );
}
