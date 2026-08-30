import { Link } from "react-router-dom";
import { motion, useMotionValue, useReducedMotion, useTransform } from "motion/react";

export function HeroCtas() {
  const reduce = useReducedMotion();
  const mx = useMotionValue(0);
  const my = useMotionValue(0);
  const x = useTransform(mx, [-48, 48], [-6, 6]);
  const y = useTransform(my, [-48, 48], [-6, 6]);

  return (
    <motion.div
      className="hero-ctas mt-8 flex flex-wrap items-center gap-x-4 gap-y-3"
      initial={reduce ? false : { opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.7, delay: 0.12, ease: [0.16, 1, 0.3, 1] }}
      onPointerMove={(e) => {
        if (reduce) return;
        const r = e.currentTarget.getBoundingClientRect();
        mx.set(e.clientX - (r.left + r.width / 2));
        my.set(e.clientY - (r.top + r.height / 2));
      }}
      onPointerLeave={() => {
        mx.set(0);
        my.set(0);
      }}
    >
      <motion.div className="w-full sm:w-auto" style={reduce ? undefined : { x, y }}>
        <Link to="/radar" className="pill pill-ink">
          Explore PIT
        </Link>
      </motion.div>
      <Link to="/download" className="pill pill-ghost">
        Download PIT Desktop
      </Link>
      <Link to="/missions" className="pill pill-cream">
        See Sleep Missions
      </Link>
    </motion.div>
  );
}
