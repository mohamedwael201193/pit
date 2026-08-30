import { Suspense, useEffect } from "react";
import { Link, NavLink, Outlet } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import {
  ClockCounterClockwise,
  Coins,
  DownloadSimple,
  Gear,
  House,
  Plus,
  ShieldCheck,
  User,
} from "@phosphor-icons/react";
import type { Icon } from "@phosphor-icons/react";
import { PitMark } from "../brand/PitMark";
import { cn } from "../lib/cn";
import { WalletGate } from "./WalletGate";

interface NavItem {
  to: string;
  label: string;
  icon: Icon;
  end: boolean;
}

const ITEMS: readonly NavItem[] = [
  { to: "/app", label: "Overview", icon: House, end: true },
  { to: "/app/activity", label: "My Missions", icon: ClockCounterClockwise, end: false },
  { to: "/app/account", label: "My Agent", icon: User, end: false },
  { to: "/app/verify", label: "My Proof", icon: ShieldCheck, end: false },
  { to: "/capital", label: "My Capital", icon: Coins, end: false },
];

export function AppShell() {
  const { ready, authenticated, user, logout } = usePrivy();
  const addr = user?.wallet?.address;

  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
    document.querySelector('meta[name="theme-color"]')?.setAttribute("content", "#D82F2F");
  }, []);

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
      <a
        href="#app-main"
        className="sr-only focus:not-sr-only focus:fixed focus:top-4 focus:left-4 focus:z-50 focus:rounded-full focus:bg-[#d82f2f] focus:px-4 focus:py-2 focus:text-black"
      >
        Skip to content
      </a>
      <Rail />
      <div className="min-w-0">
        <TopBar address={addr} onSignOut={logout} />
        <main id="app-main" className="px-5 pt-8 pb-24 sm:px-6 lg:px-10 lg:pt-10">
          <Suspense fallback={<p className="text-[rgb(240_231_212/0.55)]">Loading…</p>}>
            <Outlet />
          </Suspense>
        </main>
      </div>
    </div>
  );
}

function Rail() {
  return (
    <div className="sticky top-0 hidden h-[100dvh] flex-col border-r border-[rgb(240_231_212/0.25)] bg-[#141414] px-4 py-6 lg:flex">
      <Link to="/" className="mb-8 inline-flex px-1" aria-label="PIT home">
        <PitMark />
      </Link>
      <Link
        to="/radar"
        className="mb-2 inline-flex items-center justify-center gap-2 rounded-full border border-[rgb(240_231_212/0.28)] px-3 py-2.5 text-[0.9375rem] font-medium text-[var(--guide-cream)] no-underline"
      >
        Radar
      </Link>
      <Link
        to="/download"
        className="mb-4 inline-flex items-center justify-center gap-2 rounded-full bg-[#d82f2f] px-3 py-2.5 text-[0.9375rem] font-semibold text-[#f0e7d4] no-underline"
      >
        <DownloadSimple size={16} weight="bold" aria-hidden="true" />
        Download
      </Link>
      <nav aria-label="Main" className="flex flex-col gap-1">
        {ITEMS.map((item) => (
          <RailLink key={item.to} item={item} />
        ))}
      </nav>
      <Link
        to="/app/settings"
        className="mt-auto inline-flex items-center gap-3 rounded-full px-3 py-2.5 text-[0.875rem] text-[rgb(240_231_212/0.45)] no-underline hover:bg-[rgb(240_231_212/0.08)] hover:text-[var(--guide-cream)]"
      >
        <Gear size={18} aria-hidden="true" />
        Settings
      </Link>
    </div>
  );
}

function RailLink({ item }: { item: NavItem }) {
  const Glyph = item.icon;
  return (
    <NavLink
      to={item.to}
      end={item.end}
      className={({ isActive }) =>
        cn(
          "group inline-flex items-center gap-3 rounded-full px-3 py-2.5 text-[0.9375rem] font-medium no-underline transition-colors duration-150",
          isActive
            ? "bg-[#d82f2f] text-black"
            : "text-[rgb(240_231_212/0.72)] hover:bg-[rgb(240_231_212/0.08)] hover:text-[var(--guide-cream)]",
        )
      }
    >
      {({ isActive }) => (
        <>
          <Glyph
            size={18}
            weight={isActive ? "fill" : "regular"}
            aria-hidden="true"
            className={isActive ? "text-black" : "text-[rgb(240_231_212/0.45)]"}
          />
          {item.label}
        </>
      )}
    </NavLink>
  );
}

function TopBar({ address, onSignOut }: { address?: string; onSignOut: () => void }) {
  return (
    <header className="sticky top-0 z-30 border-b border-[rgb(240_231_212/0.25)] bg-[#1a1a1a]/95">
      <div className="flex h-14 items-center justify-between gap-4 px-5 sm:px-6 lg:px-10">
        <div className="flex min-w-0 items-center gap-3">
          <Link to="/" className="lg:hidden" aria-label="PIT home">
            <PitMark />
          </Link>
          <span className="truncate text-[0.9375rem] font-bold text-[var(--guide-cream)]">
            {address ? "Orders stay on desktop" : "No session yet"}
          </span>
          <span className="hidden font-mono text-[0.8125rem] text-[rgb(240_231_212/0.45)] sm:inline">
            {address ? "this browser cannot authorize" : "start on desktop"}
          </span>
        </div>
        <div className="flex items-center gap-2">
          <Link
            to="/app/start"
            className="inline-flex items-center gap-1 rounded-full border border-[rgb(240_231_212/0.28)] px-3 py-1.5 text-[0.8125rem] font-medium text-[var(--guide-cream)] no-underline lg:hidden"
          >
            <Plus size={14} weight="bold" aria-hidden="true" />
            Setup
          </Link>
          {address ? (
            <span
              title={address}
              className="hidden items-center gap-2 rounded-full border border-[rgb(240_231_212/0.35)] px-3 py-1.5 font-mono text-[0.8125rem] text-[rgb(240_231_212/0.7)] sm:inline-flex"
            >
              {address.slice(0, 8)}...{address.slice(-4)}
            </span>
          ) : null}
          <button
            type="button"
            onClick={onSignOut}
            className="rounded-full border border-[rgb(240_231_212/0.35)] px-3 py-1.5 text-[0.8125rem] text-[var(--guide-cream)] transition-colors hover:bg-[var(--guide-cream)] hover:text-black"
          >
            Sign out
          </button>
        </div>
      </div>
      <nav
        aria-label="Main"
        className="-mx-px flex gap-1 overflow-x-auto border-t border-[rgb(240_231_212/0.25)] px-4 py-2 lg:hidden"
      >
        {ITEMS.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.end}
            className={({ isActive }) =>
              cn(
                "shrink-0 rounded-full px-3 py-1.5 text-[0.875rem] font-medium whitespace-nowrap no-underline transition-colors",
                isActive ? "bg-[#d82f2f] text-black" : "text-[rgb(240_231_212/0.65)]",
              )
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>
    </header>
  );
}
