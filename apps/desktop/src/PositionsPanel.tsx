export type VenuePosition = {
  coin?: string;
  sz?: string;
  entryPx?: string;
  markPx?: number;
  unrealizedPnl?: string;
  leverage?: string;
  marginUsed?: string;
  account?: string;
  policyClipUsd?: number;
  source?: string;
};

export function PositionsPanel({
  account,
  positions,
  error,
  lastOrder,
  onReduceOnlyClose,
  closeBusy,
}: {
  account?: string;
  positions: VenuePosition[];
  error?: string;
  lastOrder?: { oid?: string; status?: string; market?: string; side?: string; sz?: number };
  onReduceOnlyClose?: (coin: string) => void;
  closeBusy?: boolean;
}) {
  return (
    <article className="card">
      <p className="label">POSITIONS</p>
      <p>Trading account {account || "unbound"} (master, not the PIT agent).</p>
      {error === "HYPERLIQUID_OUTAGE" ? <p role="status">Hyperliquid did not answer. PIT will not invent a book.</p> : null}
      {error === "WRONG_ACCOUNT_QUERY" ? <p>No positions on this Hyperliquid account.</p> : null}
      {error && error !== "HYPERLIQUID_OUTAGE" && error !== "WRONG_ACCOUNT_QUERY" ? <p>{error}</p> : null}
      {positions.length === 0 && !error ? <p>No open positions returned by Hyperliquid for this account.</p> : null}
      {positions.length ? (
        <table className="desk-table">
          <thead>
            <tr>
              <th>Market</th>
              <th>Size</th>
              <th>Entry</th>
              <th>Mark</th>
              <th>uPnL</th>
              <th>Margin</th>
              <th>Leverage</th>
              <th>Policy clip</th>
            </tr>
          </thead>
          <tbody>
            {positions.map((p) => (
              <tr key={p.coin}>
                <td>{p.coin}</td>
                <td>{p.sz}</td>
                <td>{p.entryPx || "—"}</td>
                <td>{p.markPx ?? "—"}</td>
                <td>{p.unrealizedPnl || "—"}</td>
                <td>{p.marginUsed || "—"}</td>
                <td>{p.leverage || "—"}</td>
                <td>{p.policyClipUsd ?? "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : null}
      {lastOrder?.oid ? (
        <p className="fine">
          Last PIT order {lastOrder.market} {lastOrder.side} {lastOrder.sz} · {lastOrder.status} · OID {lastOrder.oid}. Flatten only with a
          reduce-only close that YOU authorize. PIT cannot withdraw. Take-profit and stop-loss are not available.
        </p>
      ) : (
        <p className="fine">No PIT order on this machine yet.</p>
      )}
      {onReduceOnlyClose && positions.length ? (
        <div className="cta-row">
          {positions.map((p) =>
            p.coin ? (
              <button
                key={`close-${p.coin}`}
                type="button"
                className="linkish"
                disabled={closeBusy}
                onClick={() => onReduceOnlyClose(p.coin!)}
              >
                Prepare reduce-only close {p.coin}
              </button>
            ) : null,
          )}
        </div>
      ) : null}
    </article>
  );
}
