import { AttentionLine } from "./AttentionLine";
import { ExternalLink } from "./ExternalLink";
import type { NextFix } from "./nextFix";

export function EmptyHome({
  count,
  next,
  onGo,
}: {
  count: number;
  next: NextFix;
  onGo?: (view: NonNullable<NextFix["go"]>) => void;
}) {
  return (
    <div className="card next-fix">
      <p className="label">WHAT NEEDS YOUR ATTENTION</p>
      <h2>{next.title}</h2>
      <p>{next.why}</p>
      <p className="fine">{next.fix}</p>
      <div className="cta-row">
        {next.href ? (
          <ExternalLink className="linkish" href={next.href}>
            {next.hrefLabel || "Open"}
          </ExternalLink>
        ) : null}
        {next.go && onGo ? (
          <button type="button" className="linkish" onClick={() => onGo(next.go!)}>
            {next.goLabel || "Open"}
          </button>
        ) : null}
      </div>
      <AttentionLine count={count} />
    </div>
  );
}
