import { Suspense, useEffect, useState } from "react";
import { Link, NavLink, Outlet } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import {
  Broadcast,
  DownloadSimple,
  Flag,
  Gear,
  House,
  Info,
  LinkSimple,
  LockKey,
  ShieldCheck,
  User,
} from "@phosphor-icons/react";
import type { Icon } from "@phosphor-icons/react";
import { PitMark } from "../brand/PitMark";
import { cn } from "../lib/cn";
import { ChatPanel } from "../public/ChatPanel";
import { useWatch } from "../public/Watch";
import { windowsInstallerUrl } from "../public/facts";

interface NavItem {
  to: string;
  label: string;
  icon: Icon;
  end: boolean;
}

const LOOK: readonly NavItem[] = [
  { to: "/radar", label: "Radar", icon: Broadcast, end: false },
  { to: "/missions", label: "Missions", icon: Flag, end: false },
  { to: "/autonomy", label: "Autonomy", icon: ShieldCheck, end: true },
  { to: "/proof", label: "Proof", icon: ShieldCheck, end: true },
  { to: "/agent", label: "Agent", icon: User, end: true },
];

const MACHINE: readonly NavItem[] = [{ to: "/pair", label: "Pair", icon: LinkSimple, end: true }];

const DESK: readonly NavItem[] = [{ to: "/app", label: "Overview", icon: House, end: true }];

const PROTECT: NavItem = { to: "/protect", label: "Protect my strategy", icon: LockKey, end: true };

export function DeskFrame() {
  const [chat, setChat] = useState(false);
  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
    document.querySelector('meta[name="theme-color"]')?.setAttribute("content", "#D82F2F");
  }, []);

  return (
    <div className="guide-shell guide-app min-h-[100dvh] lg:grid lg:grid-cols-[240px_minmax(0,1fr)]">
      <a
        href="#desk-main"
        className="sr-only focus:not-sr-only focus:fixed focus:top-4 focus:left-4 focus:z-50 focus:rounded-full focus:bg-[#d82f2f] focus:px-4 focus:py-2 focus:text-black"
      >
        Skip to content
      </a>
      <Rail chatOpen={chat} onAsk={() => setChat((v) => !v)} />
      <div className="min-w-0">
        <TopBar chatOpen={chat} onAsk={() => setChat((v) => !v)} />
        <main id="desk-main" className="px-5 pt-8 pb-16 sm:px-6 lg:px-10 lg:pt-10 lg:pb-20">
          <Suspense fallback={<p className="text-[rgb(240_231_212/0.55)]">Loading…</p>}>
            <Outlet />
          </Suspense>
        </main>
      </div>
      <ChatPanel floating={false} open={chat} onOpenChange={setChat} />
    </div>
  );
}

function Rail({ chatOpen, onAsk }: { chatOpen: boolean; onAsk: () => void }) {
  const { authenticated } = usePrivy();
  return (
    <div className="sticky top-0 hidden h-[100dvh] flex-col overflow-y-auto border-r border-[rgb(240_231_212/0.25)] bg-[#141414] px-4 py-6 lg:flex">
      <Link to="/" className="mb-6 inline-flex px-1" aria-label="PIT home">
        <PitMark />
      </Link>
      <a
        href={windowsInstallerUrl()}
        className="mb-6 inline-flex items-center justify-center gap-2 rounded-full bg-[#d82f2f] px-3 py-2.5 text-[0.9375rem] font-semibold text-[#f0e7d4] no-underline"
      >
        <DownloadSimple size={16} weight="bold" aria-hidden="true" />
        Download
      </a>
      <p className="mb-2 px-3 text-[0.6875rem] tracking-[0.14em] text-[rgb(240_231_212/0.4)] uppercase">Look</p>
      <nav aria-label="Look" className="mb-6 flex flex-col gap-1">
        {LOOK.map((item) => (
          <RailLink key={item.to} item={item} />
        ))}
      </nav>
      <p className="mb-2 px-3 text-[0.6875rem] tracking-[0.14em] text-[rgb(240_231_212/0.4)] uppercase">This computer</p>
      <nav aria-label="This computer" className="mb-6 flex flex-col gap-1">
        {MACHINE.map((item) => (
          <RailLink key={item.to} item={item} />
        ))}
      </nav>
      {authenticated ? (
        <>
          <p className="mb-2 px-3 text-[0.6875rem] tracking-[0.14em] text-[rgb(240_231_212/0.4)] uppercase">Your desk</p>
          <nav aria-label="Your desk" className="mb-4 flex flex-col gap-1">
            {DESK.map((item) => (
              <RailLink key={item.to} item={item} />
            ))}
          </nav>
        </>
      ) : null}
      <RailLink item={PROTECT} />
      <button
        type="button"
        onClick={onAsk}
        aria-expanded={chatOpen}
        className="mt-auto mb-2 inline-flex items-center justify-center rounded-full border border-[rgb(240_231_212/0.28)] px-3 py-2.5 text-[0.9375rem] font-medium text-[var(--guide-cream)]"
      >
        {chatOpen ? "Close chat" : "Ask PIT"}
      </button>
      <Link
        to={authenticated ? "/app/settings" : "/how-it-works"}
        className="inline-flex items-center gap-3 rounded-full px-3 py-2.5 text-[0.875rem] text-[rgb(240_231_212/0.45)] no-underline hover:bg-[rgb(240_231_212/0.08)] hover:text-[var(--guide-cream)]"
      >
        {authenticated ? <Gear size={18} aria-hidden="true" /> : <Info size={18} aria-hidden="true" />}
        {authenticated ? "Settings" : "How it works"}
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

function TopBar({ chatOpen, onAsk }: { chatOpen: boolean; onAsk: () => void }) {
  const { ready, authenticated, user, login, logout } = usePrivy();
  const { desktop } = useWatch();
  const addr = user?.wallet?.address;
  const mobile: readonly NavItem[] = authenticated ? [...LOOK, ...MACHINE, ...DESK, PROTECT] : [...LOOK, ...MACHINE, PROTECT];

  return (
    <header className="sticky top-0 z-30 border-b border-[rgb(240_231_212/0.25)] bg-[#1a1a1a]/95 backdrop-blur">
      <div className="flex min-h-14 items-center justify-between gap-3 px-5 py-2 sm:px-6 lg:px-10">
        <div className="flex min-w-0 items-center gap-3">
          <Link to="/" className="lg:hidden" aria-label="PIT home">
            <PitMark />
          </Link>
          <p className="truncate text-[0.8125rem] leading-5 text-[rgb(240_231_212/0.7)] sm:text-[0.9375rem]">
            Orders stay on desktop
            <span className="hidden sm:inline"> · this browser cannot authorize</span>
          </p>
          {desktop.present ? (
            <span className="hidden shrink-0 rounded-full border border-[rgb(240_231_212/0.28)] px-2.5 py-1 text-[0.75rem] text-[rgb(240_231_212/0.7)] md:inline">
              Desktop{desktop.version ? ` ${desktop.version}` : ""}
            </span>
          ) : null}
        </div>
        <div className="flex shrink-0 items-center gap-2">
          <button
            type="button"
            onClick={onAsk}
            aria-expanded={chatOpen}
            className="rounded-full border border-[rgb(240_231_212/0.35)] px-3 py-1.5 text-[0.8125rem] text-[var(--guide-cream)] lg:hidden"
          >
            {chatOpen ? "Close" : "Ask PIT"}
          </button>
          {!ready ? (
            <span className="text-[0.8125rem] text-[rgb(240_231_212/0.45)]">…</span>
          ) : authenticated && addr ? (
            <>
              <span
                title={addr}
                className="hidden items-center rounded-full border border-[rgb(240_231_212/0.35)] px-3 py-1.5 font-mono text-[0.8125rem] text-[rgb(240_231_212/0.7)] sm:inline-flex"
              >
                {addr.slice(0, 8)}…{addr.slice(-4)}
              </span>
              <button
                type="button"
                onClick={() => void logout()}
                className="rounded-full border border-[rgb(240_231_212/0.35)] px-3 py-1.5 text-[0.8125rem] text-[var(--guide-cream)] hover:bg-[var(--guide-cream)] hover:text-black"
              >
                Sign out
              </button>
            </>
          ) : (
            <button
              type="button"
              onClick={() => void login()}
              className="rounded-full bg-[#d82f2f] px-3 py-1.5 text-[0.8125rem] font-medium text-[#f0e7d4]"
            >
              Connect wallet
            </button>
          )}
        </div>
      </div>
      <nav
        aria-label="Main"
        className="-mx-px flex gap-1 overflow-x-auto border-t border-[rgb(240_231_212/0.25)] px-4 py-2 lg:hidden"
      >
        {mobile.map((item) => (
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
