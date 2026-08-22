export function AttentionLine({ count }: { count: number }) {
  if (count <= 0) {
    return <p className="fine">No opportunities match your policy. Watch does not place orders.</p>;
  }
  return <p className="fine">{count} opportunities match your policy.</p>;
}
