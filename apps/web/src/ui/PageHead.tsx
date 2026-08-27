import type { ReactNode } from "react";
import { Reveal } from "./Reveal";

export function PageHead({
  title,
  lede,
  action,
}: {
  title: string;
  lede: string;
  action?: ReactNode;
}) {
  return (
    <Reveal className="flex flex-col gap-6 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 className="display-l text-[var(--guide-cream)]">{title}</h1>
        <p className="mt-3 max-w-[56ch] text-[1.0625rem] leading-7 text-[rgb(240_231_212/0.72)]">{lede}</p>
      </div>
      {action ? <div className="shrink-0">{action}</div> : null}
    </Reveal>
  );
}
