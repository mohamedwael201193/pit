import { useEffect, type ReactNode } from "react";
import { Link, Outlet } from "react-router-dom";
import { LandingNav } from "../landing/Nav";
import { WatchProvider, useWatch } from "./Watch";
import { ChatPanel } from "./ChatPanel";

export function PublicShell({ children }: { children?: ReactNode }) {
  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
  }, []);

  return (
    <WatchProvider>
      <div className="guide-shell intel-shell">
        <a
          href="#intel-main"
          className="sr-only focus:not-sr-only focus:fixed focus:top-3 focus:left-3 focus:z-50 focus:bg-[#d82f2f] focus:px-3 focus:py-2 focus:text-black"
        >
          Skip to content
        </a>
        <LandingNav />
        <main id="intel-main" className="intel-main pt-24">
          <DesktopBanner />
          {children ?? <Outlet />}
        </main>
        <ChatPanel />
        <MobileDock />
      </div>
    </WatchProvider>
  );
}

function DesktopBanner() {
  const { desktop } = useWatch();
  if (!desktop.present) return null;
  return (
    <div className="intel-banner mb-8 border border-[rgb(240_231_212/0.2)]" role="status">
      <span className="text-[#f0e7d4]">PIT Desktop detected</span>
      <span>secure local companion{desktop.version ? ` · ${desktop.version}` : ""}</span>
      <span className="hidden sm:inline">private research available. execution stays on this computer</span>
      <span className="text-[rgb(240_231_212/0.5)]">This browser still cannot receive a session key.</span>
    </div>
  );
}

function MobileDock() {
  return (
    <div className="intel-dock sm:hidden">
      <Link to="/radar" className="pill pill-line h-11 px-4 text-[0.875rem]">
        Radar
      </Link>
      <Link to="/download" className="pill pill-coral h-11 px-4 text-[0.875rem]">
        Download
      </Link>
    </div>
  );
}
