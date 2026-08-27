export function PermissionsCard() {
  return (
    <div className="card">
      <p className="label">YOUR SESSION</p>
      <ul className="perms">
        <li>ORDER — allowed after extraAgents lists this agent</li>
        <li>CANCEL — allowed after extraAgents lists this agent</li>
        <li>WITHDRAW — denied by PIT. PIT never signs it.</li>
        <li>FUND TRANSFER — denied by PIT. PIT never signs it.</li>
        <li>LEVERAGE — denied by PIT policy. PIT never signs it.</li>
        <li>ACCOUNT ADMIN — denied by PIT. PIT never signs it.</li>
      </ul>
      <p className="fine">
        Hyperliquid may still expose those actions to a raw agent key. PIT refuses to sign them. The session key is
        created on this machine. It never goes to the web app, Vercel, or a server env.
      </p>
    </div>
  );
}
