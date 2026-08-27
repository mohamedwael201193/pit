import { cn } from "../lib/cn";

export function PitMark({
  className,
  variant = "pill",
}: {
  className?: string;
  variant?: "pill" | "bare";
}) {
  return (
    <span
      className={cn(
        "inline-grid place-items-center overflow-hidden",
        variant === "pill" && "h-9 w-9 rounded-full bg-[#f0e7d4] p-0.5 sm:h-10 sm:w-10",
        variant === "bare" && "h-9 w-9",
        className,
      )}
    >
      <img src="/mark.svg" alt="" width={36} height={36} className="h-full w-full object-cover" decoding="async" />
    </span>
  );
}
