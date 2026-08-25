export function PermissionsCard() {
  return (
    <div className="card">
      <p className="label">YOUR SESSION</p>
      <ul className="perms">
        <li>order — allowed</li>
        <li>cancel — allowed</li>
        <li>withdraw — denied</li>
        <li>leverage — denied</li>
        <li>approveAgent — denied after the session exists</li>
      </ul>
      <p className="fine">
        The session key is created on this machine. It never goes to the web app, Vercel, or a server env.
      </p>
    </div>
  );
}
