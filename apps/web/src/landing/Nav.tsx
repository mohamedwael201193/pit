import { useEffect, useId, useState } from "react";
import { Link, NavLink, useLocation } from "react-router-dom";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { List, X } from "@phosphor-icons/react";
import { PitMark } from "../brand/PitMark";
import { windowsInstallerUrl } from "../public/facts";

const LINKS = [
  { to: "/radar", label: "Radar" },
  { to: "/missions", label: "Missions" },
  { to: "/proof", label: "Proof" },
  { to: "/how-it-works", label: "How it works" },
] as const;

export function LandingNav() {
  const [open, setOpen] = useState(false);
  const location = useLocation();
  const reduce = useReducedMotion();
  const menuId = useId();

  useEffect(() => setOpen(false), [location.pathname]);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("keydown", onKey);
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKey);
      document.body.style.overflow = "";
    };
  }, [open]);

  return (
    <header className="pointer-events-none fixed inset-x-0 top-0 z-50 pt-4 pb-2">
      <div className="container-pit pointer-events-auto">
        <nav aria-label="Main" className="grid grid-cols-[auto_1fr_auto] items-center gap-3">
          <Link to="/" className="shrink-0" aria-label="PIT home">
            <PitMark />
          </Link>

          <ul className="hidden items-center justify-self-center gap-0.5 rounded-full bg-[#f0e7d4] px-2 py-1.5 lg:flex">
            {LINKS.map((l) => (
              <li key={l.to}>
                <NavLink
                  to={l.to}
                  className={({ isActive }) =>
                    `inline-flex h-9 items-center rounded-full px-3.5 text-[0.9375rem] font-medium text-black no-underline hover:bg-black/5 ${isActive ? "bg-black/10" : ""}`
                  }
                >
                  {l.label}
                </NavLink>
              </li>
            ))}
            <li>
              <a
                href={windowsInstallerUrl()}
                className="inline-flex h-9 items-center rounded-full px-3.5 text-[0.9375rem] font-medium text-black no-underline hover:bg-black/5"
              >
                Download
              </a>
            </li>
          </ul>

          <div className="flex h-10 w-10 shrink-0 items-center justify-end">
            <button
              type="button"
              aria-expanded={open}
              aria-controls={menuId}
              aria-label={open ? "Close menu" : "Open menu"}
              onClick={() => setOpen((v) => !v)}
              className="grid size-11 place-items-center rounded-full bg-[#f0e7d4] text-black lg:hidden"
            >
              {open ? <X size={20} /> : <List size={20} />}
            </button>
          </div>
        </nav>
      </div>

      <AnimatePresence>
        {open ? (
          <motion.div
            id={menuId}
            className="fixed inset-0 z-30 bg-[#1a1a1a] px-5 pt-24 lg:hidden"
            initial={reduce ? { opacity: 0 } : { opacity: 0, y: -8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={reduce ? { opacity: 0 } : { opacity: 0, y: -8 }}
            transition={{ duration: 0.24, ease: [0.32, 0.72, 0, 1] }}
          >
            <ul className="container-pit flex flex-col">
              {LINKS.map((l) => (
                <li key={l.to} className="border-b border-white/10">
                  <Link to={l.to} className="block py-4 text-[1.25rem] font-medium text-[#f0e7d4] no-underline">
                    {l.label}
                  </Link>
                </li>
              ))}
              <li className="border-b border-white/10">
                <a href={windowsInstallerUrl()} className="block py-4 text-[1.25rem] font-medium text-[#f0e7d4] no-underline">
                  Download
                </a>
              </li>
              <li className="pt-6">
                <Link
                  to="/radar"
                  className="inline-flex h-12 items-center rounded-full bg-[#d82f2f] px-6 text-base font-semibold text-black no-underline"
                >
                  Explore live PIT
                </Link>
              </li>
            </ul>
          </motion.div>
        ) : null}
      </AnimatePresence>
    </header>
  );
}
