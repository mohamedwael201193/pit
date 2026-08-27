import { ArrowRight } from "@phosphor-icons/react";
import type { ComponentPropsWithoutRef, ElementType, ReactNode } from "react";
import { cn } from "../lib/cn";

type Variant = "primary" | "secondary" | "ghost";
type Size = "md" | "lg";

const VARIANTS: Record<Variant, string> = {
  primary: "bg-[#d82f2f] text-black hover:opacity-90",
  secondary: "border border-[rgb(240_231_212/0.35)] text-[var(--guide-cream)] hover:border-[var(--guide-cream)]",
  ghost: "bg-transparent text-[rgb(240_231_212/0.72)] hover:text-[var(--guide-cream)]",
};

const SIZES: Record<Size, string> = {
  md: "min-h-11 px-5 py-2.5 text-[0.9375rem]",
  lg: "min-h-12 px-6 py-3 text-base",
};

export function Button({
  variant = "primary",
  size = "md",
  trailingArrow,
  className,
  children,
  ...rest
}: ComponentPropsWithoutRef<"button"> & {
  variant?: Variant;
  size?: Size;
  trailingArrow?: boolean;
}) {
  return (
    <button className={cn(base(variant, size, Boolean(trailingArrow)), className)} {...rest}>
      <span className="whitespace-nowrap">{children}</span>
      {trailingArrow ? <ArrowWell /> : null}
    </button>
  );
}

export function ButtonLink({
  as,
  variant = "primary",
  size = "md",
  trailingArrow,
  className,
  children,
  ...rest
}: {
  as?: ElementType;
  variant?: Variant;
  size?: Size;
  trailingArrow?: boolean;
  className?: string;
  children: ReactNode;
  href?: string;
  to?: string;
} & Record<string, unknown>) {
  const Tag = (as ?? "a") as ElementType;
  return (
    <Tag className={cn(base(variant, size, Boolean(trailingArrow)), className)} {...rest}>
      <span className="whitespace-nowrap">{children}</span>
      {trailingArrow ? <ArrowWell /> : null}
    </Tag>
  );
}

function base(variant: Variant, size: Size, hasArrow: boolean): string {
  return cn(
    "group inline-flex items-center justify-center gap-2 rounded-full font-medium no-underline",
    "transition-[opacity,border-color,color,transform] duration-150",
    "active:scale-[0.98] disabled:pointer-events-none disabled:opacity-45",
    SIZES[size],
    hasArrow && "pr-1.5",
    VARIANTS[variant],
  );
}

function ArrowWell() {
  return (
    <span
      aria-hidden="true"
      className="grid size-7 place-items-center rounded-full bg-black/12 transition-transform duration-200 group-hover:translate-x-0.5"
    >
      <ArrowRight size={14} weight="bold" />
    </span>
  );
}
