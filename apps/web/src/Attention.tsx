export function Attention({ count }: { count: number }) {
  const copy =
    count <= 0
      ? "No opportunities match your policy."
      : `${count} opportunities match your policy.`;
  return (
    <div className="mt-8 border-t border-[#1d1f24] pt-6">
      <h2 className="text-xl tracking-tight">What needs your attention</h2>
      <p className="mt-3 max-w-[42ch] text-base opacity-85">{copy}</p>
      <p className="mt-2 max-w-[42ch] text-sm opacity-70">
        Watch does not invent cards and does not place orders.
      </p>
    </div>
  );
}
