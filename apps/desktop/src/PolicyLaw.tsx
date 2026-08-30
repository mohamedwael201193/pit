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
      v: pinned ? "Pinned host clip" : "Unpinned — pin on this computer",
      why: "Host sizes every clip to this ceiling. Venue minimum is per book after szDecimals — never a universal $10.",
      hit: "A larger idea is refused.",
    },
    {
      k: "Max position",
      v: "Same as max trade",
      why: "Same ceiling as the clip. Guarded Autonomy cannot grow past it. Venue min is per book after szDecimals.",
      hit: "A larger book is refused.",
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
      why: "Markets and research stay inside this universe unless you pin a change.",
      hit: "A coin outside the list is blocked before a sealed request starts.",
    },
    {
      k: "Allowed venues",
      v: "hyperliquid",
      why: "PIT only reads and posts this venue.",
      hit: "Another venue is refused.",
    },
    {
      k: "Kill switch",
      v: "Off until you flip it",
      why: "You halt new orders. The model cannot turn this off.",
      hit: "New AUTHORIZE posts are refused while it is on.",
    },
    {
      k: "Max open positions",
      v: "1",
      why: "Guarded Autonomy cannot add a market while this ceiling is full.",
      hit: "A new order is refused. Existing positions are not flattened.",
    },
    {
      k: "Consecutive losses",
      v: "3",
      why: "A losing streak stops the mission. The model cannot raise this.",
      hit: "Guarded Autonomy stands down until you enable it again.",
    },
    {
      k: "Daily loss",
      v: "$50",
      why: "Realized loss beyond this stops the mission.",
      hit: "Guarded Autonomy stops. Positions are not flattened.",
    },
    {
      k: "Withdraw / transfer",
      v: "Impossible",
      why: "The session can order and cancel only.",
      hit: "Those actions are denied at the host.",
    },
    {
      k: "Policy mutation",
      v: "Forbidden for the model",
      why: "Only you pin policy on this computer.",
      hit: "Chat and the model cannot raise clip or permissions.",
    },
    {
      k: "Max slippage",
      v: pinned ? "80 bps" : "80 bps until pinned",
      why: "Host rejects a book that would slip past this band.",
      hit: "The preview is not eligible.",
    },
    {
      k: "Liquidity / cooldown / uncertainty / session",
      v: "Host gated",
      why: "Thin books, cooldown, uncertainty, and session TTL are host law.",
      hit: "The preview or mission is refused.",
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
