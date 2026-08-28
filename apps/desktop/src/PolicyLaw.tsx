export function PolicyLaw({
  pinned,
  onPin,
  busy,
}: {
  pinned?: boolean;
  onPin?: () => void;
  busy?: boolean;
}) {
  const cells = [
    {
      k: "Max trade",
      v: pinned ? "$10" : "$10 example until pinned",
      why: "Host sizes every clip to this ceiling. The model cannot raise it.",
      hit: "A larger idea is refused. No order is sent.",
    },
    {
      k: "Max leverage",
      v: pinned ? "1x" : "1x example until pinned",
      why: "Session cannot change leverage. PIT cannot withdraw.",
      hit: "Any attempt to raise leverage is denied on this computer.",
    },
    {
      k: "Allowed assets",
      v: "ETH BTC SOL HYPE DOGE AVAX",
      why: "Watch and research stay inside this universe.",
      hit: "A coin outside the list is blocked before a sealed request starts.",
    },
    {
      k: "Kill switch",
      v: "Off until you flip it",
      why: "You halt new orders. The model cannot turn this off.",
      hit: "New AUTHORIZE posts are refused while it is on.",
    },
    {
      k: "Max slippage",
      v: pinned ? "80 bps" : "80 bps example until pinned",
      why: "Host rejects a book that would slip past this band.",
      hit: "The preview is not eligible. No order is sent.",
    },
  ];
  return (
    <div className="card">
      <p className="label">Your policy</p>
      <p className="fine">Host enforced. Chat can explain this. Chat cannot mutate it.</p>
      <div className="policy-grid">
        {cells.map((c) => (
          <div key={c.k} className="policy-cell">
            <p className="label">{c.k}</p>
            <strong>{c.v}</strong>
            <p className="host-enforced">Host enforced</p>
            <p className="fine">{c.why}</p>
            <p className="fine">If violated: {c.hit}</p>
          </div>
        ))}
      </div>
      {onPin ? (
        <button type="button" className="linkish" onClick={onPin} disabled={busy}>
          {pinned ? "Policy pinned on this computer" : "Pin policy on this computer"}
        </button>
      ) : null}
    </div>
  );
}
