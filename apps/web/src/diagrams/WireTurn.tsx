import type { MotionValue } from "motion/react";
import { motion } from "motion/react";

export function WireTurn({
  className,
  rotate,
}: {
  className?: string;
  rotate: MotionValue<number> | number;
}) {
  const teeth = 48;
  return (
    <motion.svg viewBox="0 0 100 100" className={className} style={{ rotate }} aria-hidden="true">
      <circle cx="50" cy="50" r="48.4" fill="none" stroke="#0a0a0a" strokeWidth="0.45" opacity="0.35" />
      {Array.from({ length: teeth }, (_, i) => {
        const a = (i / teeth) * Math.PI * 2 - Math.PI / 2;
        const long = i % 6 === 0;
        const mid = i % 2 === 0;
        const inner = long ? 22 : mid ? 28 : 34;
        return (
          <line
            key={i}
            x1={50 + inner * Math.cos(a)}
            y1={50 + (inner) * Math.sin(a)}
            x2={50 + 47.2 * Math.cos(a)}
            y2={50 + 47.2 * Math.sin(a)}
            stroke="#0a0a0a"
            strokeWidth={long ? 1.7 : mid ? 1.15 : 0.85}
            strokeLinecap="round"
          />
        );
      })}
      <circle cx="50" cy="50" r="20.5" fill="none" stroke="#0a0a0a" strokeWidth="1.15" />
      <circle cx="50" cy="50" r="14.2" fill="none" stroke="#0a0a0a" strokeWidth="0.9" strokeDasharray="2.2 2.8" />
      <circle cx="50" cy="50" r="8.4" fill="none" stroke="#0a0a0a" strokeWidth="1.05" />
      <circle cx="50" cy="50" r="4.1" fill="#0a0a0a" />
      <circle cx="50" cy="50" r="1.35" fill="#f0e7d4" />
    </motion.svg>
  );
}

export function WireIris({ className = "" }: { className?: string }) {
  return (
    <svg viewBox="0 0 100 100" className={className} aria-hidden="true">
      <circle cx="50" cy="50" r="16.5" fill="none" stroke="#f0e7d4" strokeWidth="0.7" opacity="0.85" />
      {Array.from({ length: 12 }, (_, i) => {
        const a = (i / 12) * Math.PI * 2;
        return (
          <line
            key={i}
            x1={50 + 7 * Math.cos(a)}
            y1={50 + 7 * Math.sin(a)}
            x2={50 + 15.2 * Math.cos(a)}
            y2={50 + 15.2 * Math.sin(a)}
            stroke="#f0e7d4"
            strokeWidth="0.7"
            strokeLinecap="round"
            opacity="0.7"
          />
        );
      })}
    </svg>
  );
}
