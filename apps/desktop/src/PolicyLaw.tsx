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
    {
      k: "Max trade",
      v: pinned ? "$10" : "$10 until pinned",
      why: "Host sizes every clip to this ceiling.",
      hit: "A larger idea is refused.",
    },
    {
      k: "Max leverage",
      v: pinned ? "1x" : "1x until pinned",
      why: "Session cannot change leverage. PIT cannot withdraw.",
      hit: "Raising leverage is denied here.",
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
      v: pinned ? "80 bps" : "80 bps until pinned",
      why: "Host rejects a book that would slip past this band.",
      hit: "The preview is not eligible.",
    },
  ];
  return (
    <section>
      <p className="fine">Host enforced. Chat can explain this. Chat cannot mutate it.</p>
      <table className="desk-table">
        <thead>
          <tr>
            <th>Control</th>
            <th>Value</th>
            <th>Why</th>
            <th>If violated</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((c) => (
            <tr key={c.k}>
              <td>{c.k}</td>
              <td>{c.v}</td>
              <td>{c.why}</td>
              <td>{c.hit}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {onPin ? (
        <button type="button" className="linkish" onClick={onPin} disabled={busy}>
          {pinned ? "Policy pinned on this computer" : "Pin policy on this computer"}
        </button>
      ) : null}
    </section>
  );
}
