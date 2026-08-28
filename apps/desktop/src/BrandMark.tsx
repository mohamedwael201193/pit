/** Token and chain marks from the cryptocurrency-icons geometry (spothq, MIT). Unknown tickers render as a letter badge — never an invented logo. */

import type { ReactNode } from "react";

const SIZE = 16;

export function BrandMark({
  symbol,
  size = SIZE,
  title,
}: {
  symbol: string;
  size?: number;
  title?: string;
}) {
  const s = String(symbol || "").toUpperCase();
  if (s === "ETH") return <EthMark size={size} />;
  if (s === "BTC") return <BtcMark size={size} />;
  if (s === "SOL") return <SolMark size={size} />;
  if (s === "DOGE") return <DogeMark size={size} />;
  if (s === "AVAX") return <AvaxMark size={size} />;
  if (s === "HYPE" || s === "HYPERLIQUID" || s === "HL") return <HyperliquidMark size={size} title={title || "Hyperliquid"} />;
  if (s === "0G" || s === "OG") return <OgMark size={size} />;
  return <LetterMark letter={s.slice(0, 1) || "?"} size={size} title={title || s} />;
}

function wrap(size: number, fill: string, children: ReactNode, title: string) {
  return (
    <svg width={size} height={size} viewBox="0 0 32 32" aria-hidden={title ? undefined : true} aria-label={title} role="img">
      <title>{title}</title>
      <circle cx="16" cy="16" r="16" fill={fill} />
      {children}
    </svg>
  );
}

function EthMark({ size }: { size: number }) {
  return wrap(
    size,
    "#627eea",
    <>
      <path fill="#fff" fillOpacity="0.6" d="M16 4.5 9.5 16.2 16 19.9z" />
      <path fill="#fff" d="M16 4.5 22.5 16.2 16 19.9z" />
      <path fill="#fff" fillOpacity="0.6" d="M16 21.3 9.5 17.4 16 27.5z" />
      <path fill="#fff" d="M16 21.3 22.5 17.4 16 27.5z" />
    </>,
    "ETH",
  );
}

function BtcMark({ size }: { size: number }) {
  return wrap(
    size,
    "#f7931a",
    <path
      fill="#fff"
      d="M22.2 14.1c.3-2.1-1.3-3.2-3.5-4l.7-2.8-1.7-.4-.7 2.7c-.5-.1-.9-.2-1.4-.3l.7-2.7-1.7-.4-.7 2.8c-.4-.1-.8-.2-1.1-.2v-.1l-2.4-.6-.5 1.8s1.3.3 1.3.3c.7.2.8.6.9 1l-.9 3.5c.1 0 .1 0 .2.1h-.2l-1.2 4.8c-.1.2-.3.6-.8.4 0 0-1.3-.3-1.3-.3l-.9 2 2.2.5c.4.1.8.2 1.2.3l-.7 2.9 1.7.4.7-2.8c.5.1 1 .3 1.4.3l-.7 2.8 1.7.4.7-2.9c2.9.6 5.1.3 6-2.3.7-2.1 0-3.4-1.6-4.2 1.1-.3 2-1 2.2-2.6zm-3.3 5.2c-.5 2.1-4.1 1-5.3.7l.9-3.8c1.2.3 5 .9 4.4 3.1zm.5-5.3c-.5 1.9-3.5.9-4.5.7l.9-3.4c1 .2 4.1.7 3.6 2.7z"
    />,
    "BTC",
  );
}

function SolMark({ size }: { size: number }) {
  return wrap(
    size,
    "#000",
    <>
      <path fill="#14f195" d="M9.5 20.8h11.4l-2.2 2.4H7.3z" />
      <path fill="#9945ff" d="M9.5 14.8h11.4l-2.2 2.4H7.3z" />
      <path fill="#00ffa3" d="M20.9 8.8H9.5L7.3 11.2h13.6z" />
    </>,
    "SOL",
  );
}

function DogeMark({ size }: { size: number }) {
  return wrap(
    size,
    "#c2a633",
    <path fill="#fff" d="M13.2 8.2h4.2c3.6 0 6.1 2.4 6.1 5.8 0 3.5-2.5 5.9-6.2 5.9h-2.4V24h-1.7zm1.7 1.6v7.3h2.4c2.5 0 4.2-1.5 4.2-3.7 0-2.1-1.7-3.6-4.2-3.6z" />,
    "DOGE",
  );
}

function AvaxMark({ size }: { size: number }) {
  return wrap(size, "#e84142", <path fill="#fff" d="M16 7.5 22.8 21H19.6L16 14.2 12.4 21H9.2z" />, "AVAX");
}

function HyperliquidMark({ size, title }: { size: number; title: string }) {
  return wrap(size, "#072723", <path fill="#97fce4" d="M11 8h2.2v6.2H19V8H21.2v16H19v-7.2h-5.8V24H11z" />, title);
}

function OgMark({ size }: { size: number }) {
  return wrap(
    size,
    "#111",
    <text x="16" y="21" textAnchor="middle" fill="#f0e7d4" fontSize="11" fontFamily="ui-sans-serif, system-ui" fontWeight="700">
      0G
    </text>,
    "0G",
  );
}

function LetterMark({ letter, size, title }: { letter: string; size: number; title: string }) {
  return wrap(
    size,
    "#1c1e22",
    <text x="16" y="21" textAnchor="middle" fill="#f0e7d4" fontSize="14" fontFamily="ui-sans-serif, system-ui" fontWeight="700">
      {letter}
    </text>,
    title,
  );
}
