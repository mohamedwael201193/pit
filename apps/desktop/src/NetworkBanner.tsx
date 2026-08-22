export function NetworkBanner({ net }: { net: "mainnet" | "testnet" }) {
  if (net === "testnet") {
    return <p className="fine">TESTNET laboratory. Do not mix with production.</p>;
  }
  return <p className="fine">MAINNET production. Transfer of Agentic ID is not live.</p>;
}
