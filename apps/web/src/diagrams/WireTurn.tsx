import type { MotionValue } from "motion/react";
import { motion } from "motion/react";

export function WireTurn({
  className,
  rotate,
}: {
  className?: string;
  rotate: MotionValue<number> | number;
}) {
  const teeth = 32;
  return (
    <motion.svg viewBox="0 0 100 100" className={className} style={{ rotate }} aria-hidden="true">
      {Array.from({ length: teeth }, (_, i) => {
        const a = (i / teeth) * Math.PI * 2 - Math.PI / 2;
        return (
          <line
            key={i}
            x1={50 + 30 * Math.cos(a)}
            y1={50 + 30 * Math.sin(a)}
            x2={50 + 47 * Math.cos(a)}
            y2={50 + 47 * Math.sin(a)}
            stroke="#0a0a0a"
            strokeWidth="1.2"
            strokeLinecap="round"
          />
        );
      })}
      <circle cx="50" cy="50" r="11" fill="none" stroke="#0a0a0a" strokeWidth="1.15" />
      <circle cx="50" cy="50" r="3.5" fill="#0a0a0a" />
    </motion.svg>
  );
}
