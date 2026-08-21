const CARDS = [
  { title: "Max trade", value: "10 USD" },
  { title: "Max leverage", value: "1x" },
  { title: "Allowed assets", value: "ETH, BTC" },
  { title: "Kill switch", value: "off" },
];

export function PolicyPanel() {
  return (
    <div className="mt-8 border-t border-[#1d1f24] pt-6">
      <h2 className="text-xl tracking-tight">Your policy is the law</h2>
      <dl className="mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2">
        {CARDS.map((c) => (
          <div key={c.title}>
            <dt className="text-sm opacity-70">{c.title}</dt>
            <dd className="font-mono text-lg">{c.value}</dd>
          </div>
        ))}
      </dl>
      <p className="mt-3 max-w-[40ch] text-sm opacity-80">The model cannot raise clip, leverage, or permissions.</p>
    </div>
  );
}
