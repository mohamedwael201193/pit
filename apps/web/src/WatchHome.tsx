import { useEffect, useState } from "react";
import { Attention } from "./Attention";
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
        setCoins(Array.isArray(body.coins) ? body.coins : []);
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
            <li
              key={c.coin}
              className="rounded-2xl border border-[rgb(240_231_212/0.22)] bg-[#141414] p-6"
            >
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
                  <dt className="text-[rgb(240_231_212/0.5)]">Policy</dt>
                  <dd>{c.eligible ? "eligible to research" : "blocked"}</dd>
                </div>
                <div className="sm:col-span-2">
                  <dt className="text-[rgb(240_231_212/0.5)]">Why it is on Watch</dt>
                  <dd>{c.reason}</dd>
                </div>
                <div className="sm:col-span-2">
                  <dt className="text-[rgb(240_231_212/0.5)]">Thesis / size / calibration</dt>
                  <dd>
                    Open desktop to seal the private book. Web does not invent a thesis. NOT ENOUGH DATA until 30
                    resolved forecasts.
                  </dd>
                </div>
              </dl>
            </li>
          ))}
        </ul>
      ) : (
        <p className="mt-6 max-w-[42ch] text-[0.9375rem] text-[rgb(240_231_212/0.65)]">
          Live books only. PIT will not invent a card to fill the home.
        </p>
      )}
    </div>
  );
}
