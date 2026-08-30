import { useEffect, useId, useState, type ReactNode } from "react";
import { Link, NavLink, Outlet, useLocation } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { List, X } from "@phosphor-icons/react";
import { PitMark } from "../brand/PitMark";
import { cn } from "../lib/cn";
import { PUBLIC_NAV } from "./facts";
import { WatchProvider, useWatch } from "./Watch";
import { ChatPanel } from "./ChatPanel";

export function PublicShell({ children }: { children?: ReactNode }) {
  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
  }, []);

  return (
    <WatchProvider>
      <div className="intel-shell">
        <Skip />
        <TopNav />
        <DesktopBanner />
        <main id="intel-main" className="intel-main">
          {children ?? <Outlet />}
        </main>
        <ChatPanel />
        <MobileDock />
      </div>
    </WatchProvider>
  );
}

function Skip() {
  return (
    <a href="#intel-main" className="sr-only focus:not-sr-only focus:fixed focus:top-3 focus:left-3 focus:z-50 focus:bg-[#d82f2f] focus:px-3 focus:py-2 focus:text-black">
      Skip to content
    </a>
  );
}

function TopNav() {
  const [open, setOpen] = useState(false);
  const location = useLocation();
  const { authenticated } = usePrivy();
  const menuId = useId();

  useEffect(() => setOpen(false), [location.pathname]);

  return (
    <header className="intel-nav">
      <div className="intel-nav-inner">
        <Link to="/" className="flex items-center gap-2.5 no-underline" aria-label="PIT home">
          <PitMark />
          <span className="hidden font-semibold tracking-tight sm:inline">PIT</span>
        </Link>
        <nav aria-label="Public" className="hidden items-center gap-0.5 lg:flex">
          {PUBLIC_NAV.map((l) => (
            <NavLink
              key={l.to}
              to={l.to}
              end={l.end}
              className={({ isActive }) => cn("intel-link", isActive && "intel-link-on")}
            >
              {l.label}
            </NavLink>
          ))}
        </nav>
        <div className="flex items-center gap-2">
          <Link to="/signin" className="intel-ghost hidden sm:inline-flex">
            {authenticated ? "Overview" : "Connect"}
          </Link>
          <Link to="/download" className="intel-cta hidden sm:inline-flex">
            Download PIT
          </Link>
          <button
            type="button"
            className="intel-icon-btn lg:hidden"
            aria-expanded={open}
            aria-controls={menuId}
            aria-label={open ? "Close menu" : "Open menu"}
            onClick={() => setOpen((v) => !v)}
          >
            {open ? <X size={18} /> : <List size={18} />}
          </button>
        </div>
      </div>
      {open ? (
        <div id={menuId} className="intel-drawer lg:hidden">
          {PUBLIC_NAV.map((l) => (
            <NavLink key={l.to} to={l.to} end={l.end} className="intel-drawer-link">
              {l.label}
            </NavLink>
          ))}
          <Link to="/signin" className="intel-drawer-link">
            {authenticated ? "Overview" : "Connect wallet"}
          </Link>
        </div>
      ) : null}
    </header>
  );
}

function DesktopBanner() {
  const { desktop } = useWatch();
  if (!desktop.present) return null;
  return (
    <div className="intel-banner" role="status">
      <span className="text-[#7dffb3]">PIT DESKTOP DETECTED</span>
      <span>secure local companion{desktop.version ? ` · ${desktop.version}` : ""}</span>
      <span className="hidden sm:inline">private research available · execution stays on this computer</span>
      <span className="text-[rgb(240_231_212/0.5)]">This browser still cannot receive a session key.</span>
    </div>
  );
}

function MobileDock() {
  return (
    <div className="intel-dock sm:hidden">
      <Link to="/radar" className="intel-ghost">
        Radar
      </Link>
      <Link to="/download" className="intel-cta flex-1 justify-center">
        Download PIT Desktop
      </Link>
    </div>
  );
}
