const CARDS = [
  { title: "Max trade", value: "10 USD" },
  { title: "Max daily loss", value: "50 USD" },
  { title: "Max leverage", value: "1x" },
  { title: "Allowed assets", value: "ETH, BTC" },
  { title: "Allowed venues", value: "hyperliquid" },
  { title: "Max slippage", value: "80 bps" },
  { title: "Cooldown", value: "0 s" },
  { title: "Max uncertainty", value: "1.0" },
  { title: "Min liquidity", value: "0 USD" },
  { title: "Session TTL", value: "3600 s" },
  { title: "Min calibration", value: "0.00" },
  { title: "Kill switch", value: "off" },
];

export function PolicyPanel() {
  return (
    <div className="mt-8">
      <h2 className="text-[1.25rem] font-semibold tracking-[-0.03em]">Your policy is the law</h2>
      <dl className="mt-5 grid grid-cols-1 gap-3 sm:grid-cols-2">
        {CARDS.map((c) => (
          <div key={c.title} className="border border-[rgb(240_231_212/0.22)] bg-[#141414] p-5">
            <dt className="text-[0.8125rem] text-[rgb(240_231_212/0.55)]">{c.title}</dt>
            <dd className="mt-2 font-mono text-[1.125rem]">{c.value}</dd>
          </div>
        ))}
      </dl>
      <p className="mt-4 max-w-[40ch] text-[0.875rem] text-[rgb(240_231_212/0.7)]">
        The model cannot raise clip, leverage, or permissions.
      </p>
    </div>
  );
}
