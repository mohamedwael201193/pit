import { useEffect, useState } from "react";
import { Attention } from "./Attention";
import { Bezel } from "./ui/Surface";
import { DiagramEmptyWatch } from "./diagrams/pitGuide";
import { namedState } from "./namedStates";

type WatchPayload = {
  count?: number;
  copy?: string;
  trade?: boolean;
  sign?: boolean;
  coins?: { coin: string; reason: string; mark: number; eligible?: boolean; network?: string }[];
};

export function WatchHome({ network = "mainnet" }: { network?: "mainnet" | "testnet" }) {
  const base = import.meta.env.VITE_HEALTH_URL;
  const [count, setCount] = useState(0);
  const [coins, setCoins] = useState<NonNullable<WatchPayload["coins"]>>([]);
  const [fail, setFail] = useState<string | null>(null);

  useEffect(() => {
    if (!base) {
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
          setCount(0);
          setCoins([]);
          return;
        }
        setCount(typeof body.count === "number" ? body.count : 0);
        const rows = Array.isArray(body.coins) ? body.coins : [];
        const pass = rows.filter((c) => c.eligible).slice(0, 6);
        setCoins(pass);
        setFail(null);
      })
      .catch(() => {
        if (!gone) {
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
      <Attention count={count} />
      {fail ? (
        <p role="alert" className="mt-4 max-w-[48ch] text-[0.9375rem] text-[#ff7a7a]">
          {fail} Empty Watch is the honest state.
        </p>
      ) : null}
      {coins.length > 0 ? (
        <ul className="mt-6 grid gap-4 lg:grid-cols-2">
          {coins.map((c) => (
            <li key={c.coin} className="overflow-hidden border border-[rgb(240_231_212/0.22)] bg-[#141414]">
              <div className="p-6">
                <p className="font-mono text-[0.75rem] tracking-[0.12em] text-[#d82f2f]">
                  hyperliquid:perp:{c.coin}
                </p>
                <p className="mt-2 text-[1.35rem] font-semibold tracking-[-0.03em]">{c.coin}</p>
                <dl className="mt-4 grid gap-3 text-[0.875rem] sm:grid-cols-2">
                  <div>
                    <dt className="text-[rgb(240_231_212/0.5)]">Mark</dt>
                    <dd className="font-mono">{c.mark}</dd>
                  </div>
                  <div>
                    <dt className="text-[rgb(240_231_212/0.5)]">Side</dt>
                    <dd>not decided on this website</dd>
                  </div>
                  <div>
                    <dt className="text-[rgb(240_231_212/0.5)]">Entry / invalidation / edge</dt>
                    <dd>not invented. Open desktop to seal the book.</dd>
                  </div>
                  <div>
                    <dt className="text-[rgb(240_231_212/0.5)]">Policy</dt>
                    <dd>{c.eligible ? "PASS" : "BLOCKED"}</dd>
                  </div>
                  <div>
                    <dt className="text-[rgb(240_231_212/0.5)]">Research</dt>
                    <dd>not sealed on web</dd>
                  </div>
                  <div>
                    <dt className="text-[rgb(240_231_212/0.5)]">Confidence</dt>
                    <dd>NOT ENOUGH DATA</dd>
                  </div>
                  <div className="sm:col-span-2">
                    <dt className="text-[rgb(240_231_212/0.5)]">Why</dt>
                    <dd>{c.reason}</dd>
                  </div>
                  <div className="sm:col-span-2">
                    <dt className="text-[rgb(240_231_212/0.5)]">Action</dt>
                    <dd>
                      AUTHORIZE exists only on desktop or CLI after a verified preview. This card cannot trade.
                    </dd>
                  </div>
                </dl>
              </div>
            </li>
          ))}
        </ul>
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
