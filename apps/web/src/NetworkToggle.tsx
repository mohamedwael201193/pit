type Net = "mainnet" | "testnet";

export function NetworkToggle({ net, onChange }: { net: Net; onChange: (n: Net) => void }) {
  return (
    <div className="mt-8 flex gap-2">
      {(["mainnet", "testnet"] as const).map((n) => (
        <button
          key={n}
          type="button"
          onClick={() => onChange(n)}
          className={
            net === n
              ? "rounded-full bg-coral px-6 py-3 font-semibold text-white"
              : "rounded-full bg-[#15171b] px-6 py-3 font-semibold text-cream"
          }
        >
          {n === "mainnet" ? "MAINNET" : "TESTNET"}
        </button>
      ))}
    </div>
  );
}
