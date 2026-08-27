export function PolicyLaw() {
  const rows = [
    ["Max trade", "example 10 USD until your file is pinned"],
    ["Max daily loss", "example 50 USD until your file is pinned"],
    ["Max leverage", "example 1x until your file is pinned"],
    ["Max slippage", "example 80 bps until your file is pinned"],
    ["Cooldown", "example 0 s"],
    ["Session TTL", "example 3600 s"],
    ["Min calibration", "NOT ENOUGH DATA until N is large enough"],
    ["Allowed", "ETH, BTC on Hyperliquid after you pin them"],
    ["Kill", "off until you flip it"],
  ];
  return (
    <div className="card">
      <p className="label">YOUR POLICY</p>
      <ul className="perms">
        {rows.map(([k, v]) => (
          <li key={k}>
            {k}: {v}
          </li>
        ))}
      </ul>
      <p className="fine">The model cannot raise clip, leverage, or permissions.</p>
    </div>
  );
}
