import { AttentionLine } from "./AttentionLine";
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
          <a className="linkish" href={next.href} target="_blank" rel="noreferrer">
            {next.hrefLabel || "Open"}
          </a>
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
