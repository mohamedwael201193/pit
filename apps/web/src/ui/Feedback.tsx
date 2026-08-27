import type { ReactNode } from "react";
import { Warning } from "@phosphor-icons/react";
import { cn } from "../lib/cn";
import { Button } from "./Button";
import { PitMark } from "../brand/PitMark";

export function EmptyState({
  title,
  body,
  action,
  className,
}: {
  title: string;
  body: string;
  action?: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "flex flex-col items-center gap-4 border border-dashed border-[rgb(240_231_212/0.28)] px-6 py-16 text-center",
        className,
      )}
    >
      <PitMark variant="bare" className="opacity-90" />
      <h3 className="text-[1.0625rem] font-medium text-[var(--guide-cream)]">{title}</h3>
      <p className="max-w-[42ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.72)]">{body}</p>
      {action ? <div className="mt-2">{action}</div> : null}
    </div>
  );
}

export function ErrorState({
  title = "That did not load",
  body = "The request failed. Nothing was sent.",
  onRetry,
  className,
}: {
  title?: string;
  body?: string;
  onRetry?: () => void;
  className?: string;
}) {
  return (
    <div
      role="alert"
      className={cn("flex flex-col items-start gap-3 border border-[#ff7a7a]/25 p-6", className)}
    >
      <span className="inline-flex items-center gap-2 text-[#ff7a7a]">
        <Warning size={17} aria-hidden="true" />
        <span className="text-[0.9375rem] font-medium">{title}</span>
      </span>
      <p className="max-w-[52ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.72)]">{body}</p>
      {onRetry ? (
        <Button variant="secondary" onClick={onRetry} className="mt-1">
          Try again
        </Button>
      ) : null}
    </div>
  );
}

export function Skeleton({ className }: { className?: string }) {
  return <span aria-hidden="true" className={cn("block h-4 bg-[#2a2a2a]", className)} />;
}
