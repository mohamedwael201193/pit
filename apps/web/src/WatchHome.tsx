import { useEffect, useState } from "react";
import { Attention } from "./Attention";

type WatchPayload = {
  count?: number;
  copy?: string;
  trade?: boolean;
  sign?: boolean;
  coins?: { coin: string; reason: string; mark: number }[];
};

export function WatchHome() {
  const base = import.meta.env.VITE_HEALTH_URL;
  const [count, setCount] = useState(0);
  const [coins, setCoins] = useState<{ coin: string; reason: string; mark: number }[]>([]);

  useEffect(() => {
    if (!base) {
      return;
    }
    let gone = false;
    fetch(`${base.replace(/\/$/, "")}/watch`)
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
      })
      .catch(() => {
        if (!gone) {
          setCount(0);
          setCoins([]);
        }
      });
    return () => {
      gone = true;
    };
  }, [base]);

  return (
    <div>
      <Attention count={count} />
      {coins.length > 0 ? (
        <ul className="mt-4 grid gap-2">
          {coins.map((c) => (
            <li key={c.coin} className="font-mono text-sm">
              {c.coin} {c.reason} mark={c.mark}
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  );
}
