import type { ComponentType } from "react";
import { motion, useReducedMotion } from "motion/react";
import { cn } from "../lib/cn";

export function ChoiceCard({
  title,
  body,
  onClick,
  active,
  badge,
  Diagram,
}: {
  title: string;
  body: string;
  onClick: () => void;
  active?: boolean;
  badge?: string;
  Diagram?: ComponentType<{ className?: string }>;
}) {
  const reduce = useReducedMotion();
  return (
    <motion.button
      type="button"
      onClick={onClick}
      whileHover={reduce ? undefined : { y: -2 }}
      whileTap={reduce ? undefined : { scale: 0.99 }}
      transition={{ type: "spring", stiffness: 280, damping: 22 }}
      className={cn(
        "overflow-hidden border bg-[#141414] text-left transition-colors",
        "hover:border-[#d82f2f]/55",
        active ? "border-[#d82f2f]" : "border-[rgb(240_231_212/0.25)]",
      )}
    >
      {Diagram ? (
        <div className="overflow-hidden border-b border-[rgb(240_231_212/0.22)]">
          <Diagram className="aspect-[16/10] w-full" />
        </div>
      ) : null}
      <div className="p-6">
        <div className="flex items-start justify-between gap-3">
          <h3 className="text-[1.25rem] font-semibold tracking-[-0.03em] text-[var(--guide-cream)]">{title}</h3>
          {badge ? (
            <span className="rounded-full bg-[#d82f2f] px-2 py-0.5 text-[0.6875rem] font-bold text-black">
              {badge}
            </span>
          ) : null}
        </div>
        <p className="mt-2 text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">{body}</p>
      </div>
    </motion.button>
  );
}
