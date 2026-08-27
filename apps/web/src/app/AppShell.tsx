import { Link, NavLink, Outlet } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { WalletGate } from "./WalletGate";

const ITEMS = [
  { to: "/app", label: "Home", end: true },
  { to: "/app/start", label: "Start", end: false },
  { to: "/verify", label: "Verify", end: false },
] as const;

export function AppShell() {
  const { ready, authenticated, user, logout } = usePrivy();
  const addr = user?.wallet?.address;

  if (!ready) {
    return (
      <div className="guide-shell grid min-h-[100dvh] place-items-center">
        <p>Loading wallet connect</p>
      </div>
    );
  }
  if (!authenticated) return <WalletGate />;

  return (
    <div className="guide-shell guide-app min-h-[100dvh] lg:grid lg:grid-cols-[240px_minmax(0,1fr)]">
      <aside className="sticky top-0 hidden h-[100dvh] flex-col border-r border-[rgb(240_231_212/0.25)] bg-[#141414] px-4 py-6 lg:flex">
        <Link to="/" className="mb-8 text-[1.35rem] font-bold tracking-[-0.06em] text-[#d82f2f] no-underline">
          PIT.
        </Link>
        <Link
          to="/app/start"
          className="mb-4 inline-flex items-center justify-center rounded-full bg-[#d82f2f] px-3 py-2.5 text-[0.9375rem] font-semibold text-black no-underline"
        >
          Start
        </Link>
        <nav aria-label="Main" className="flex flex-col gap-1">
          {ITEMS.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                `rounded-full px-3 py-2.5 text-[0.9375rem] font-medium no-underline ${
                  isActive ? "bg-[#d82f2f] text-black" : "text-[rgb(240_231_212/0.72)]"
                }`
              }
            >
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <div className="min-w-0">
        <header className="flex items-center justify-between gap-4 border-b border-[rgb(240_231_212/0.2)] px-5 py-4 sm:px-8">
          <p className="text-[0.875rem] text-[rgb(240_231_212/0.7)]">No session in this browser</p>
          <div className="flex items-center gap-3">
            <span className="rounded-full border border-[rgb(240_231_212/0.28)] px-3 py-1 font-mono text-[0.75rem]">
              {addr ? `${addr.slice(0, 8)}...${addr.slice(-4)}` : ""}
            </span>
            <button className="pill pill-line h-9 text-[0.8125rem]" type="button" onClick={logout}>
              Sign out
            </button>
          </div>
        </header>
        <main className="px-5 pt-8 pb-24 sm:px-6 lg:px-10 lg:pt-10">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
