export function Attention({ scanned, count }: { scanned: number; count: number }) {
  const copy =
    scanned <= 0
      ? "Public Watch is empty. PIT will not invent a book."
      : `${scanned} live Hyperliquid books. ${count} pass the default policy. None of these cards can trade.`;
  return (
    <div>
      <h2 id="watch-heading" className="text-[1.25rem] font-semibold tracking-[-0.03em]">
        Watch
      </h2>
      <p className="mt-3 max-w-[52ch] text-[1.0625rem] text-[rgb(240_231_212/0.85)]">{copy}</p>
      <p className="mt-2 max-w-[52ch] text-[0.875rem] text-[rgb(240_231_212/0.55)]">
        This website has no session and no buying power. Execution sizing, AUTHORIZE, and Guarded Autonomy stay on PIT Desktop.
      </p>
    </div>
  );
}
