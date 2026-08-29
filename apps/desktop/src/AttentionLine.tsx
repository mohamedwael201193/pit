export function AttentionLine({ count }: { count: number }) {
  if (count <= 0) {
    return <p className="fine">No books pass default policy. Watch does not place orders.</p>;
  }
  return <p className="fine">{count} books pass default policy. Execution stays on this computer.</p>;
}
