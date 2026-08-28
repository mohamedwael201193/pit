export function PolicyLaw({
  pinned,
  onPin,
  busy,
}: {
  pinned?: boolean;
  onPin?: () => void;
  busy?: boolean;
}) {
  const rows = [
    ["Max trade", pinned ? "10 USD (pinned)" : "example 10 USD until your file is pinned"],
    ["Max daily loss", pinned ? "50 USD (pinned)" : "example 50 USD until your file is pinned"],
    ["Max leverage", pinned ? "1x (pinned)" : "example 1x until your file is pinned"],
    ["Max slippage", pinned ? "80 bps (pinned)" : "example 80 bps until your file is pinned"],
    ["Cooldown", "0 s"],
    ["Session TTL", "3600 s"],
    ["Min calibration", "NOT ENOUGH DATA until N is large enough"],
    ["Allowed", "ETH, BTC, SOL, HYPE, DOGE, AVAX on Hyperliquid"],
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
      {onPin ? (
        <button type="button" className="linkish" onClick={onPin} disabled={busy}>
          {pinned ? "Policy pinned on this computer" : "Pin policy on this computer"}
        </button>
      ) : null}
    </div>
  );
}
