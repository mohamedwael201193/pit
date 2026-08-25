export function PolicyLaw() {
  const rows = [
    ["Max trade", "10 USD"],
    ["Max daily loss", "50 USD"],
    ["Max leverage", "1x"],
    ["Max slippage", "80 bps"],
    ["Cooldown", "0 s"],
    ["Session TTL", "3600 s"],
    ["Min calibration", "0.00"],
    ["Allowed", "ETH, BTC on Hyperliquid"],
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
