export function Attention({ count }: { count: number }) {
  const copy =
    count <= 0
      ? "No opportunities match your policy."
      : `${count} opportunities match your policy.`;
  return (
    <div>
      <h2 className="text-[1.25rem] font-semibold tracking-[-0.03em]">Watch</h2>
      <p className="mt-3 max-w-[42ch] text-[1.0625rem] text-[rgb(240_231_212/0.85)]">{copy}</p>
      <p className="mt-2 max-w-[42ch] text-[0.875rem] text-[rgb(240_231_212/0.55)]">
        Watch does not invent cards and does not place orders. The full Hyperliquid universe and execution sizing stay on PIT Desktop.
      </p>
    </div>
  );
}
