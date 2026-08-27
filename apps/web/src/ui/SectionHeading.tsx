import type { ReactNode } from "react";
import { cn } from "../lib/cn";
import { Reveal } from "./Reveal";

export function SectionHeading({
  title,
  body,
  className,
  id,
}: {
  title: ReactNode;
  body?: ReactNode;
  className?: string;
  id?: string;
}) {
  return (
    <Reveal className={cn("flex flex-col gap-4 items-start", className)}>
      <h2 id={id} className="guide-display max-w-[14ch]">
        {title}
      </h2>
      {body ? <p className="max-w-[58ch] text-[1.2rem] leading-8 text-[rgb(240_231_212/0.75)]">{body}</p> : null}
    </Reveal>
  );
}
