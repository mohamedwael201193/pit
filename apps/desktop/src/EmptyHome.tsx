import { AttentionLine } from "./AttentionLine";

export function EmptyHome() {
  return (
    <div className="card">
      <p className="label">WHAT NEEDS YOUR ATTENTION</p>
      <AttentionLine count={0} />
      <p className="fine">Watch does not invent cards. It also does not place orders.</p>
    </div>
  );
}
