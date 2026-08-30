import type { MotionValue } from "motion/react";
import { motion } from "motion/react";

export function WireTurn({
  className,
  rotate,
}: {
  className?: string;
  rotate: MotionValue<number> | number;
}) {
  const teeth = 40;
  return (
    <motion.svg viewBox="0 0 100 100" className={className} style={{ rotate }} aria-hidden="true">
      {Array.from({ length: teeth }, (_, i) => {
        const a = (i / teeth) * Math.PI * 2 - Math.PI / 2;
        const long = i % 4 === 0;
        return (
          <line
            key={i}
            x1={50 + (long ? 26 : 32) * Math.cos(a)}
            y1={50 + (long ? 26 : 32) * Math.sin(a)}
            x2={50 + 48 * Math.cos(a)}
            y2={50 + 48 * Math.sin(a)}
            stroke="#0a0a0a"
            strokeWidth={long ? 1.5 : 1.05}
            strokeLinecap="round"
          />
        );
      })}
      <circle cx="50" cy="50" r="18" fill="none" stroke="#0a0a0a" strokeWidth="1.1" />
      <circle cx="50" cy="50" r="10" fill="none" stroke="#0a0a0a" strokeWidth="1.05" />
      <circle cx="50" cy="50" r="3.4" fill="#0a0a0a" />
    </motion.svg>
  );
}
