export function NetworkBanner({ net }: { net: "mainnet" | "testnet" }) {
  if (net === "testnet") {
    return (
      <p className="mt-4 max-w-[48ch] text-sm opacity-80">
        TESTNET is the full integration lab. Galileo and Hyperliquid testnet stay on this workspace only.
      </p>
    );
  }
  return (
    <p className="mt-4 max-w-[48ch] text-sm opacity-80">
      MAINNET is production. Transfer of Agentic ID is not live. Direct TeeML stays on Aristotle.
    </p>
  );
}
