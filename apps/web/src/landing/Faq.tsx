const ITEMS = [
  {
    q: "Who holds my money?",
    a: "You do. PIT never custody. The session agent can order or cancel on Hyperliquid. It cannot withdraw.",
  },
  {
    q: "Where does signing happen?",
    a: "On desktop or CLI. This website can connect, inspect, and verify. It cannot create a session key.",
  },
  {
    q: "What if TeeML fails?",
    a: "PIT stops. There is no Router fallback for the private book.",
  },
] as const;

export function Faq() {
  return (
    <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit grid gap-10 lg:grid-cols-[0.8fr_1.2fr]">
        <h2 className="text-4xl tracking-[-0.04em]">Questions worth answering once</h2>
        <dl className="grid gap-8">
          {ITEMS.map((item) => (
            <div key={item.q}>
              <dt className="text-[1.15rem] font-semibold">{item.q}</dt>
              <dd className="mt-2 max-w-[52ch] text-[1.05rem] leading-7 text-[rgb(240_231_212/0.75)]">{item.a}</dd>
            </div>
          ))}
        </dl>
      </div>
    </section>
  );
}
