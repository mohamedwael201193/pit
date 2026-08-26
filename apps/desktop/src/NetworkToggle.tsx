type Net = "mainnet" | "testnet";

export function NetworkToggle({ net, onChange }: { net: Net; onChange: (n: Net) => void }) {
  return (
    <div className="row">
      {(["mainnet", "testnet"] as const).map((n) => (
        <button key={n} type="button" className={net === n ? "on" : "off"} onClick={() => onChange(n)}>
          {n === "mainnet" ? "MAINNET" : "TESTNET"}
        </button>
      ))}
    </div>
  );
}
