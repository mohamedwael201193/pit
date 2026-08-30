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
      className="mt-10 flex flex-wrap items-center gap-3"
      initial={reduce ? false : { opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.7, delay: 0.16, ease: [0.16, 1, 0.3, 1] }}
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
      <motion.div style={reduce ? undefined : { x, y }}>
        <Link to="/radar" className="pill pill-ink">
          Explore PIT
        </Link>
      </motion.div>
      <Link to="/download" className="pill pill-ghost">
        Download PIT Desktop
      </Link>
      <Link to="/missions" className="text-[0.9375rem] font-medium text-black underline-offset-4 hover:underline">
        See Sleep Missions
      </Link>
    </motion.div>
  );
}
