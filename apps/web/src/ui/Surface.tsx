import type { ReactNode } from "react";
import { cn } from "../lib/cn";

export function Bezel({
  children,
  className,
  coreClassName,
  as: Tag = "div",
}: {
  children: ReactNode;
  className?: string;
  coreClassName?: string;
  as?: "div" | "article" | "section" | "li";
}) {
  return (
    <Tag className={cn("border border-[rgb(240_231_212/0.28)] bg-[#222] p-1.5", className)}>
      <div className={cn("bg-[#2a2a2a] p-6 sm:p-8", coreClassName)}>{children}</div>
    </Tag>
  );
}
