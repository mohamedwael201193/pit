import { BrandMark } from "./BrandMark";
import { ExternalLink } from "./ExternalLink";
import { hyperliquidTrade } from "./links";

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
  summary,
  onReduceOnlyClose,
  closeBusy,
  net,
}: {
  account?: string;
  positions: VenuePosition[];
  error?: string;
  lastOrder?: { oid?: string; status?: string; market?: string; side?: string; sz?: number };
  summary?: {
    accountValue?: string;
    withdrawable?: string;
    totalNtlPos?: string;
    spotUsdc?: string;
    buyingPower?: string;
    powerSource?: string;
    execGate?: string;
    execWhy?: string;
  };
  onReduceOnlyClose?: (coin: string) => void;
  closeBusy?: boolean;
  net?: string;
}) {
  const tradePile = summary?.buyingPower || summary?.spotUsdc || summary?.accountValue;
  return (
    <section>
      <dl className="metrics">
        <div>
          <dt>Trading equity</dt>
          <dd>{tradePile ? `$${tradePile}` : "—"}</dd>
        </div>
        <div>
          <dt>Available to trade</dt>
          <dd>{summary?.buyingPower ? `$${summary.buyingPower}` : summary?.spotUsdc ? `$${summary.spotUsdc}` : "—"}</dd>
        </div>
        <div>
          <dt>Spot USDC</dt>
          <dd>{summary?.spotUsdc || "—"}</dd>
        </div>
        <div>
          <dt>Perp account value</dt>
          <dd>{summary?.accountValue || "—"}</dd>
        </div>
      </dl>
      {summary?.execWhy ? (
        <p className="fine" role="status">
          {summary.execGate ? `Execution blocked: ${String(summary.execGate).replaceAll("_", " ")}. ` : ""}
          {summary.execWhy} Venue leverage on an existing position is not PIT policy leverage.
        </p>
      ) : null}
      <p className="fine">
        Trading account {account || "unbound"} (master, not the PIT agent). Unified accounts show trading equity from
        spot USDC, not perp account value. Venue withdrawable is not a PIT withdraw
        {summary?.powerSource ? ` · ${summary.powerSource.replaceAll("_", " ")}` : ""}.
      </p>
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
              <th>Leverage</th>
              <th>Policy clip</th>
              <th>Why open</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {positions.map((p) => (
              <tr key={p.coin}>
                <td>
                  <span className="asset">
                    <BrandMark symbol={p.coin || ""} />
                    {p.coin}
                  </span>
                </td>
                <td>{p.sz}</td>
                <td>{p.entryPx || "—"}</td>
                <td>{p.markPx ?? "—"}</td>
                <td>{p.unrealizedPnl || "—"}</td>
                <td>{p.leverage || "—"}</td>
                <td>{p.policyClipUsd ?? "—"}</td>
                <td>
                  {lastOrder?.market === p.coin && lastOrder?.oid
                    ? `PIT fill OID ${lastOrder.oid}`
                    : "Venue position. Close needs a fresh reduce-only preview."}
                </td>
                <td>
                  {onReduceOnlyClose && p.coin ? (
                    <button type="button" className="linkish" disabled={closeBusy} onClick={() => onReduceOnlyClose(p.coin!)}>
                      Prepare close
                    </button>
                  ) : null}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : null}
      {lastOrder?.oid ? (
        <p className="fine">
          Last PIT order {lastOrder.market} {lastOrder.side} {lastOrder.sz} · {lastOrder.status} ·{" "}
          {net ? (
            <ExternalLink href={hyperliquidTrade(net, String(lastOrder.market || "").split(":").pop())}>
              OID {lastOrder.oid}
            </ExternalLink>
          ) : (
            <>OID {lastOrder.oid}</>
          )}
          {positions.length === 0 ? " · Historical fill — not an open position on this account." : ""}. Flatten only with a reduce-only close that YOU authorize.
        </p>
      ) : (
        <p className="fine">No PIT order on this machine yet.</p>
      )}
    </section>
  );
}
