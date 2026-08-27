import { Link } from "react-router-dom";

export function LandingNav() {
  return (
    <header className="absolute inset-x-0 top-0 z-10">
      <div className="container-pit flex h-16 items-center justify-between">
        <Link to="/" className="text-[1.35rem] font-bold tracking-[-0.06em] text-white no-underline" aria-label="PIT home">
          PIT.
        </Link>
        <nav className="hidden items-center gap-1 rounded-full bg-[#f0e7d4] px-1 py-1 sm:flex">
          <a className="rounded-full px-4 py-2 text-[0.875rem] font-medium text-black no-underline" href="#story">
            How it works
          </a>
          <a className="rounded-full px-4 py-2 text-[0.875rem] font-medium text-black no-underline" href="#verify">
            Verify
          </a>
        </nav>
        <Link to="/signin" className="pill pill-ink h-10 text-[0.875rem]">
          Get started
        </Link>
      </div>
    </header>
  );
}
