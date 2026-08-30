import TokenAVAX from "@web3icons/react/icons/tokens/TokenAVAX";
import TokenBTC from "@web3icons/react/icons/tokens/TokenBTC";
import TokenDOGE from "@web3icons/react/icons/tokens/TokenDOGE";
import TokenETH from "@web3icons/react/icons/tokens/TokenETH";
import TokenHYPE from "@web3icons/react/icons/tokens/TokenHYPE";
import TokenSOL from "@web3icons/react/icons/tokens/TokenSOL";

const COINS = [
  { symbol: "ETH", Icon: TokenETH },
  { symbol: "BTC", Icon: TokenBTC },
  { symbol: "SOL", Icon: TokenSOL },
  { symbol: "HYPE", Icon: TokenHYPE },
  { symbol: "DOGE", Icon: TokenDOGE },
  { symbol: "AVAX", Icon: TokenAVAX },
] as const;

export function TokenRow() {
  return (
    <ul className="mt-10 flex flex-wrap justify-center gap-3">
      {COINS.map(({ symbol, Icon }) => (
        <li
          key={symbol}
          className="flex items-center gap-2.5 border border-[rgb(240_231_212/0.32)] bg-[#141414] px-3.5 py-2.5"
        >
          <Icon size={28} variant="mono" color="#f0e7d4" aria-hidden="true" />
          <span className="text-[0.95rem] font-semibold tracking-[-0.03em] text-[var(--guide-cream)]">{symbol}</span>
        </li>
      ))}
    </ul>
  );
}
