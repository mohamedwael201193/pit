import type { ReactNode, SVGProps } from "react";
import { PIPELINE } from "./pipeline";

const C = {
  coral: "#D82F2F",
  cream: "#F0E7D4",
  ink: "#1A1A1A",
  black: "#0A0A0A",
  muted: "rgb(26 26 26 / 0.55)",
};

function Frame({
  children,
  className = "",
  label = "",
  fill = C.cream,
}: {
  children: ReactNode;
  className?: string;
  label?: string;
  fill?: string;
}) {
  return (
    <svg
      viewBox="0 0 640 480"
      className={className}
      role="img"
      aria-label={label || undefined}
      xmlns="http://www.w3.org/2000/svg"
    >
      <rect width="640" height="480" fill={fill} />
      {children}
    </svg>
  );
}

function Face({
  children,
  size = 14,
  weight = "700",
  fill,
  ...rest
}: SVGProps<SVGTextElement> & { size?: number; weight?: string }) {
  return (
    <text
      fill={fill}
      fontFamily="Host Grotesk, system-ui, sans-serif"
      fontSize={size}
      fontWeight={weight}
      {...rest}
    >
      {children}
    </text>
  );
}

/** Hero postcard: eleven pipeline seats, AUTHORIZE is the lit seat. */
export function DiagramHeroPostcard({ className = "" }: { className?: string }) {
  const cx = 320;
  const cy = 228;
  const r = 118;
  const seats = PIPELINE.length;
  return (
    <Frame className={className} label="PIT desk. You authorize the exact preview." fill={C.coral}>
      {Array.from({ length: 36 }, (_, i) => (
        <circle key={i} cx={(i * 97) % 640} cy={(i * 53) % 400} r="0.8" fill={C.black} opacity="0.12" />
      ))}
      <Face x="36" y="44" fill={C.cream} size={22} letterSpacing="-0.04em">
        PIT desk
      </Face>
      <Face x="36" y="70" fill={C.black} size={14} weight="500" opacity={0.7}>
        private book never leaves the seal
      </Face>

      <circle cx={cx} cy={cy} r={r} fill="none" stroke={C.black} strokeWidth="1.5" />
      <circle
        cx={cx}
        cy={cy}
        r={r - 28}
        fill="none"
        stroke={C.black}
        strokeWidth="1"
        strokeDasharray="4 6"
        opacity="0.45"
      />

      {PIPELINE.map((label, i) => {
        const a = (i / seats) * Math.PI * 2 - Math.PI / 2;
        const x = cx + r * Math.cos(a);
        const y = cy + r * Math.sin(a);
        const lit = label === "AUTHORIZE";
        return (
          <g key={label}>
            <circle
              cx={x}
              cy={y}
              r={lit ? 18 : 14}
              fill={lit ? C.cream : C.black}
              stroke={C.black}
              strokeWidth="1.5"
            />
            <Face x={x} y={y + 4} textAnchor="middle" fill={lit ? C.black : C.cream} size={11} weight="700">
              {i + 1}
            </Face>
          </g>
        );
      })}

      <Face x={cx} y={cy - 6} textAnchor="middle" fill={C.cream} size={13} weight="600" letterSpacing="0.14em">
        AUTHORIZE
      </Face>
      <Face x={cx} y={cy + 28} textAnchor="middle" fill={C.cream} size={28} weight="800" letterSpacing="-0.05em">
        exact preview
      </Face>
    </Frame>
  );
}

/** Private book sealed into an envelope. */
export function DiagramPrivate({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="Your private book is sealed">
      <rect x="40" y="36" width="28" height="8" fill={C.coral} />
      <Face x="76" y="48" fill={C.coral} size={13} letterSpacing="0.2em">
        PRIVATE
      </Face>
      <path d="M88 148 H552 L560 168 V348 H80 V168 Z" fill={C.ink} />
      <path d="M80 168 L320 286 L560 168" fill={C.coral} />
      <path d="M80 168 L320 286 L560 168" fill="none" stroke={C.black} strokeWidth="1.5" />
      <circle cx="320" cy="252" r="70" fill="none" stroke={C.ink} strokeWidth="1.2" opacity="0.35" />
      <circle cx="320" cy="252" r="58" fill={C.cream} />
      <circle cx="320" cy="252" r="58" fill="none" stroke={C.black} strokeWidth="1.5" />
      <circle cx="320" cy="252" r="46" fill="none" stroke={C.ink} strokeWidth="0.8" strokeDasharray="3 4" />
      <Face x="320" y="244" textAnchor="middle" fill={C.ink} size={13} weight="700" letterSpacing="0.14em">
        BOOK
      </Face>
      <Face x="320" y="268" textAnchor="middle" fill={C.ink} size={18} weight="800">
        sealed
      </Face>
      <Face x="40" y="430" fill={C.ink} size={16} weight="500">
        Router keys are refused. Direct TeeML only.
      </Face>
    </Frame>
  );
}

/** Three envelopes: Researcher, Challenger, Risk. */
export function DiagramSealed({ className = "" }: { className?: string }) {
  const roles = ["RESEARCH", "CHALLENGE", "RISK"];
  return (
    <Frame className={className} label="Three sealed envelopes" fill={C.ink}>
      <rect x="40" y="32" width="28" height="8" fill={C.coral} />
      <Face x="76" y="44" fill={C.cream} size={13} letterSpacing="0.2em">
        COMMITTEE
      </Face>
      {roles.map((role, i) => (
        <g key={role} transform={`translate(${48 + i * 190} 96)`}>
          <rect width="168" height="248" fill={i === 1 ? C.coral : C.cream} />
          <rect x="8" y="8" width="152" height="232" fill="none" stroke={i === 1 ? C.ink : C.ink} strokeWidth="1" opacity="0.18" />
          <rect x="16" y="24" width="136" height="108" fill={i === 1 ? C.cream : C.ink} />
          <Face
            x="84"
            y="86"
            textAnchor="middle"
            fill={i === 1 ? C.ink : C.cream}
            size={13}
            weight="800"
            letterSpacing="0.1em"
          >
            {role}
          </Face>
          <circle cx="84" cy="188" r="18" fill={i === 1 ? C.ink : C.coral} />
          <circle cx="84" cy="188" r="6" fill={i === 1 ? C.cream : C.ink} />
        </g>
      ))}
      <Face x="40" y="430" fill="rgb(240 231 212 / 0.72)" size={15}>
        Same provider is labeled as role separation. Not three TEEs.
      </Face>
    </Frame>
  );
}

/** Policy clip blocking an oversized order. */
export function DiagramPolicy({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="Policy clips size before you see a preview">
      <Face x="40" y="56" fill={C.coral} size={14} letterSpacing="0.18em">
        POLICY
      </Face>
      <rect x="48" y="100" width="240" height="220" fill={C.ink} />
      <Face x="168" y="160" textAnchor="middle" fill={C.cream} size={15} weight="700">
        Model wants
      </Face>
      <Face x="168" y="210" textAnchor="middle" fill={C.coral} size={36} weight="800">
        400
      </Face>
      <Face x="168" y="244" textAnchor="middle" fill="rgb(240 231 212 / 0.55)" size={13}>
        USD size
      </Face>

      <rect x="352" y="100" width="240" height="220" fill={C.coral} />
      <Face x="472" y="160" textAnchor="middle" fill={C.cream} size={15} weight="700">
        Your clip
      </Face>
      <Face x="472" y="210" textAnchor="middle" fill={C.cream} size={36} weight="800">
        10
      </Face>
      <Face x="472" y="244" textAnchor="middle" fill={C.black} size={13} opacity={0.7}>
        USD max
      </Face>

      <rect x="220" y="300" width="200" height="88" fill="none" stroke={C.ink} strokeWidth="2" />
      <Face x="320" y="352" textAnchor="middle" fill={C.ink} size={18} weight="800">
        host sizes
      </Face>
      <Face x="40" y="440" fill={C.ink} size={15}>
        The model cannot raise clip, leverage, or permissions.
      </Face>
    </Frame>
  );
}

/** Exact preview, user types AUTHORIZE. */
export function DiagramAuthorize({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="You authorize one exact preview" fill={C.ink}>
      <Face x="40" y="52" fill={C.cream} size={14} letterSpacing="0.18em">
        PREVIEW
      </Face>
      <Face x="40" y="110" fill={C.cream} size={28} weight="800" letterSpacing="-0.04em">
        ETH  buy  10 USD
      </Face>
      <Face x="40" y="150" fill="rgb(240 231 212 / 0.55)" size={14}>
        hash binds market, side, size, session
      </Face>
      <rect x="40" y="200" width="360" height="160" fill={C.cream} />
      <Face x="60" y="250" fill={C.ink} size={13} letterSpacing="0.14em">
        TYPE ON DESKTOP
      </Face>
      <Face x="60" y="310" fill={C.coral} size={36} weight="800">
        AUTHORIZE
      </Face>
      <rect x="420" y="200" width="180" height="160" fill={C.coral} />
      <Face x="510" y="270" textAnchor="middle" fill={C.black} size={14} weight="700">
        this browser
      </Face>
      <Face x="510" y="304" textAnchor="middle" fill={C.cream} size={22} weight="800">
        cannot
      </Face>
      <Face x="40" y="430" fill="rgb(240 231 212 / 0.7)" size={15}>
        Session lives in the OS keychain. Not here.
      </Face>
    </Frame>
  );
}

/** Receipt plus storage proof. */
export function DiagramProve({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="Prove the order from chain and storage">
      <Face x="40" y="56" fill={C.coral} size={14} letterSpacing="0.18em">
        PROVE
      </Face>
      {[0, 1, 2, 3, 4].map((i) => (
        <rect
          key={i}
          x={40 + i * 18}
          y={110 + i * 28}
          width="240"
          height="48"
          fill={i === 4 ? C.coral : C.ink}
          opacity={0.35 + i * 0.15}
        />
      ))}
      <Face x="56" y="142" fill={C.cream} size={13} weight="700">
        0G receipt
      </Face>
      <Face x="360" y="140" fill={C.ink} size={14} weight="700" letterSpacing="0.1em">
        STORAGE PROOF
      </Face>
      {[0, 1, 2].map((i) => (
        <g key={i}>
          <circle cx={400 + i * 70} cy={220} r="28" fill="none" stroke={C.ink} strokeWidth="2" />
          <circle cx={400 + i * 70} cy={220} r="8" fill={C.coral} />
        </g>
      ))}
      <path d="M428 220 H442 M498 220 H512" stroke={C.ink} strokeWidth="2" opacity="0.5" />
      <rect x="40" y="340" width="560" height="88" fill={C.ink} />
      <Face x="60" y="392" fill={C.cream} size={18} weight="700">
        Verify reconstructs facts. Not a screenshot.
      </Face>
    </Frame>
  );
}

export function DiagramLearn({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="Calibration waits for resolved forecasts" fill={C.ink}>
      <Face x="40" y="56" fill={C.cream} size={14} letterSpacing="0.18em">
        LEARN
      </Face>
      <Face x="40" y="140" fill={C.cream} size={42} weight="800" letterSpacing="-0.05em">
        NOT ENOUGH
      </Face>
      <Face x="40" y="192" fill={C.coral} size={42} weight="800" letterSpacing="-0.05em">
        DATA
      </Face>
      <rect x="40" y="240" width="560" height="140" fill={C.cream} />
      <Face x="60" y="300" fill={C.ink} size={16} weight="700">
        Brier and ECE after outcomes.
      </Face>
      <Face x="60" y="336" fill={C.muted} size={15}>
        Need 30 resolved forecasts before skill scores.
      </Face>
    </Frame>
  );
}

export function DiagramMainnet({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="MAINNET production desk" fill={C.coral}>
      <Face x="40" y="56" fill={C.cream} size={14} letterSpacing="0.16em">
        PRODUCTION
      </Face>
      <Face x="40" y="140" fill={C.cream} size={36} weight="800" letterSpacing="-0.05em">
        Aristotle
      </Face>
      <Face x="40" y="188" fill={C.black} size={22} weight="700">
        16661
      </Face>
      <rect x="40" y="230" width="560" height="160" fill={C.ink} />
      <Face x="60" y="290" fill={C.cream} size={16} weight="700">
        Hyperliquid mainnet
      </Face>
      <Face x="60" y="328" fill="rgb(240 231 212 / 0.65)" size={14}>
        Transfer of Agentic ID is not live.
      </Face>
    </Frame>
  );
}

export function DiagramTestnet({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="TESTNET laboratory desk">
      <Face x="40" y="56" fill={C.coral} size={14} letterSpacing="0.16em">
        LABORATORY
      </Face>
      <Face x="40" y="140" fill={C.ink} size={36} weight="800" letterSpacing="-0.05em">
        Galileo
      </Face>
      <Face x="40" y="188" fill={C.coral} size={22} weight="700">
        16602
      </Face>
      <rect x="40" y="230" width="560" height="160" fill={C.ink} />
      <Face x="60" y="290" fill={C.cream} size={16} weight="700">
        Hyperliquid testnet
      </Face>
      <Face x="60" y="328" fill="rgb(240 231 212 / 0.65)" size={14}>
        Sealed ask stays off until proven.
      </Face>
    </Frame>
  );
}

export function DiagramSession({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="Session lives in the OS keychain" fill={C.ink}>
      <Face x="40" y="56" fill={C.cream} size={14} letterSpacing="0.16em">
        KEYCHAIN
      </Face>
      <rect x="80" y="110" width="480" height="220" fill="none" stroke={C.cream} strokeWidth="2" />
      <rect x="140" y="160" width="360" height="120" fill={C.coral} />
      <Face x="320" y="214" textAnchor="middle" fill={C.black} size={14} weight="700" letterSpacing="0.12em">
        ONE HOUR
      </Face>
      <Face x="320" y="250" textAnchor="middle" fill={C.cream} size={28} weight="800">
        local session
      </Face>
      <Face x="40" y="400" fill="rgb(240 231 212 / 0.7)" size={15}>
        This browser cannot hold the agent. Desktop or CLI.
      </Face>
    </Frame>
  );
}

export function DiagramHyperliquid({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="Connect your Hyperliquid trading account">
      <Face x="40" y="56" fill={C.coral} size={14} letterSpacing="0.16em">
        VENUE
      </Face>
      <circle cx="200" cy="240" r="90" fill={C.ink} />
      <Face x="200" y="236" textAnchor="middle" fill={C.cream} size={18} weight="800">
        YOU
      </Face>
      <path d="M300 240 H360" stroke={C.coral} strokeWidth="3" />
      <path d="M348 228 L364 240 L348 252" fill="none" stroke={C.coral} strokeWidth="3" />
      <rect x="380" y="150" width="200" height="180" fill={C.coral} />
      <Face x="480" y="230" textAnchor="middle" fill={C.cream} size={16} weight="800">
        HL
      </Face>
      <Face x="480" y="262" textAnchor="middle" fill={C.black} size={13}>
        public state
      </Face>
      <Face x="40" y="420" fill={C.ink} size={15}>
        Spot USDC counts as funded. Web reads. It does not trade.
      </Face>
    </Frame>
  );
}

export function DiagramEmptyWatch({ className = "" }: { className?: string }) {
  return (
    <Frame className={className} label="Empty Watch is an honest state" fill={C.ink}>
      <Face x="40" y="56" fill={C.cream} size={14} letterSpacing="0.16em">
        WATCH
      </Face>
      <circle cx="320" cy="210" r="70" fill="none" stroke={C.cream} strokeWidth="2" strokeDasharray="6 8" />
      <Face x="320" y="216" textAnchor="middle" fill={C.coral} size={18} weight="800">
        0
      </Face>
      <Face x="40" y="360" fill={C.cream} size={28} weight="800" letterSpacing="-0.04em">
        No invented cards.
      </Face>
      <Face x="40" y="410" fill="rgb(240 231 212 / 0.6)" size={15}>
        Live books only. PIT will not fill the home to look busy.
      </Face>
    </Frame>
  );
}

export function DiagramWideBanner({ className = "" }: { className?: string }) {
  return (
    <svg viewBox="0 0 1280 360" className={className} role="img" aria-label="Market to proof pipeline">
      <rect width="1280" height="360" fill={C.ink} />
      <path d="M76 140 H1204" stroke={C.coral} strokeWidth="2" opacity="0.55" />
      {PIPELINE.map((label, i) => {
        const x = 40 + i * 112;
        const lit = label === "AUTHORIZE";
        return (
          <g key={label}>
            <circle cx={x + 36} cy={140} r={lit ? 30 : 20} fill={lit ? C.coral : C.ink} stroke={C.cream} strokeWidth="2" />
            {lit ? (
              <circle cx={x + 36} cy={140} r="12" fill={C.cream} />
            ) : (
              <Face x={x + 36} y={145} textAnchor="middle" fill={C.cream} size={11} weight="700">
                {i + 1}
              </Face>
            )}
            <Face x={x + 36} y={204} textAnchor="middle" fill={C.cream} size={12} weight="700" letterSpacing="0.08em">
              {label}
            </Face>
          </g>
        );
      })}
      <Face x="40" y="300" fill="rgb(240 231 212 / 0.72)" size={18} weight="600">
        Live books in. You in the middle. Proof out.
      </Face>
    </svg>
  );
}
