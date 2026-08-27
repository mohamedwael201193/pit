import { WatchHome } from "./WatchHome";

export function EmptyWatch({ network = "mainnet" }: { network?: "mainnet" | "testnet" }) {
  return <WatchHome network={network} />;
}
