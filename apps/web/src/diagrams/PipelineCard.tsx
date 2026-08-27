export const PIPELINE = [
  "MARKET",
  "PRIVATE",
  "SEALED",
  "RESEARCH",
  "CHALLENGE",
  "RISK",
  "POLICY",
  "AUTHORIZE",
  "EXECUTE",
  "PROVE",
  "LEARN",
] as const;

export function PipelineCard({ className = "" }: { className?: string }) {
  return (
    <figure className={`border border-black bg-black/5 ${className}`}>
      <svg viewBox="0 0 320 240" className="aspect-[4/3] w-full" role="img" aria-label="PIT desk pipeline">
        <rect width="320" height="168" fill="#d82f2f" />
        <text x="16" y="28" fill="#0a0a0a" fontSize="13" fontWeight="700" letterSpacing="0.08em">
          PIT desk
        </text>
        <text x="16" y="48" fill="#0a0a0a" fontSize="11">
          private book never leaves the seal
        </text>
        {PIPELINE.map((label, i) => {
          const col = i < 6 ? 0 : 1;
          const row = i < 6 ? i : i - 6;
          const x = 24 + col * 148;
          const y = 72 + row * 16;
          return (
            <g key={label}>
              <circle cx={x} cy={y - 4} r={label === "AUTHORIZE" ? 5 : 3} fill={label === "AUTHORIZE" ? "#fff" : "#0a0a0a"} />
              <text x={x + 12} y={y} fill="#0a0a0a" fontSize="10" fontFamily="ui-monospace, monospace">
                {label}
              </text>
            </g>
          );
        })}
        <rect y="168" width="320" height="72" fill="#f0e7d4" />
        <text x="16" y="194" fill="#0a0a0a" fontSize="11">
          You authorize the exact preview.
        </text>
        <text x="16" y="214" fill="#0a0a0a" fontSize="11">
          The session cannot withdraw.
        </text>
      </svg>
    </figure>
  );
}
