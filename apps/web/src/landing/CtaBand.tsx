import { Link } from "react-router-dom";

export function CtaBand() {
  return (
    <section className="guide-coral">
      <div className="container-pit flex min-h-[40vh] flex-col justify-end gap-6 py-16">
        <h2 className="guide-mega">PIT.</h2>
        <div className="flex flex-wrap gap-3">
          <Link to="/signin" className="pill pill-ink">
            Get started
          </Link>
          <a href="#verify" className="pill pill-ghost">
            Verify a receipt
          </a>
        </div>
      </div>
    </section>
  );
}
