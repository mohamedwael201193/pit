import { Link } from "react-router-dom";

export function HowPage() {
  return (
    <div>
      <p className="intel-kicker">How it works</p>
      <h1 className="intel-title mt-2">The web discovers and proves. The desktop protects and acts.</h1>
      <p className="intel-lede">
        MAINNET only. Aristotle 16661 and Hyperliquid mainnet. The laboratory exists for CI and developers, not for
        the public desk.
      </p>

      <ol className="intel-steps mt-10">
        <li>
          <strong>Explore.</strong> Live radar from public Hyperliquid marks.
        </li>
        <li>
          <strong>Understand 0G.</strong> Compute seals intelligence. Storage holds evidence. Chain is the public
          record.
        </li>
        <li>
          <strong>Verify proof.</strong> What was checked, and how.
        </li>
        <li>
          <strong>Meet PIT.</strong> Agent passport, not a wallet homepage.
        </li>
        <li>
          <strong>Simulate capital.</strong> Venue floors only. Labeled simulation.
        </li>
        <li>
          <strong>Download desktop.</strong> Private brain, policy, keys, session.
        </li>
        <li>
          <strong>Connect, then pair.</strong> Pairing is a late step after public value.
        </li>
      </ol>

      <section className="intel-section" id="og">
        <h2 className="intel-title">0G is the private OS</h2>
        <dl className="intel-grid-2 mt-8">
          <div className="intel-pair">
            <dt>0G Compute</dt>
            <dd>Private intelligence. Direct TeeML. No Router for the private book.</dd>
          </div>
          <div className="intel-pair">
            <dt>0G Storage</dt>
            <dd>Durable evidence / memory when a public-safe object exists.</dd>
          </div>
          <div className="intel-pair">
            <dt>0G Chain</dt>
            <dd>Verifiable proof. Read the transaction from the public RPC.</dd>
          </div>
          <div className="intel-pair">
            <dt>Agentic ID</dt>
            <dd>Identity. iTransfer is not live on mainnet.</dd>
          </div>
          <div className="intel-pair">
            <dt>ERC-8004</dt>
            <dd>Reputation / feedback. Never an invented leaderboard.</dd>
          </div>
        </dl>
      </section>

      <section className="intel-section">
        <h2 className="intel-title">Safer than giving a bot a wallet</h2>
        <p className="intel-lede">
          A session agent can order or cancel. It cannot withdraw. Chat, this website, and the model cannot authorize, pin
          policy, or enable Guarded Autonomy. Duplicate clicks do not send a second order.
        </p>
      </section>

      <div className="mt-10 flex flex-wrap gap-2.5">
        <Link to="/radar" className="intel-cta">
          Explore live PIT
        </Link>
        <Link to="/download" className="intel-secondary">
          Download PIT Desktop
        </Link>
        <Link to="/pair" className="intel-ghost">
          Pair later
        </Link>
      </div>
    </div>
  );
}
