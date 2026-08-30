type Net = "mainnet" | "testnet";

function labOn() {
  try {
    return localStorage.getItem("pit.developer") === "1" || new URLSearchParams(window.location.search).get("dev") === "1";
  } catch {
    return false;
  }
}

export function NetworkToggle({ net, onChange }: { net: Net; onChange: (n: Net) => void }) {
  if (!labOn()) {
    return (
      <p className="mt-4 text-[0.875rem] tracking-[0.12em] text-[rgb(240_231_212/0.55)]">MAINNET · Aristotle 16661</p>
    );
  }
  return (
    <div className="mt-4 flex gap-2">
      {(["mainnet", "testnet"] as const).map((n) => (
        <button
          key={n}
          type="button"
          onClick={() => onChange(n)}
          className={
            net === n
              ? "rounded-full border border-[#f0e7d4] bg-[#f0e7d4] px-5 py-2.5 text-[0.9375rem] font-medium text-black"
              : "rounded-full border border-[rgb(240_231_212/0.35)] px-5 py-2.5 text-[0.9375rem] font-medium text-[#f0e7d4] hover:border-[#f0e7d4]"
          }
        >
          {n === "mainnet" ? "MAINNET" : "TESTNET"}
        </button>
      ))}
    </div>
  );
}
