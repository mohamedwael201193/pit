import type { MouseEvent, ReactNode } from "react";
import { isAllowedHttps } from "./allowedUrl";
import { openExternal } from "./open";

export function ExternalLink({
  href,
  className,
  children,
  title,
}: {
  href: string;
  className?: string;
  children: ReactNode;
  title?: string;
}) {
  const allowed = isAllowedHttps(href);
  function onClick(e: MouseEvent<HTMLAnchorElement>) {
    e.preventDefault();
    if (allowed) void openExternal(href);
  }
  return (
    <a className={className} href={href} title={title} target="_blank" rel="noreferrer noopener" onClick={onClick}>
      {children}
    </a>
  );
}
