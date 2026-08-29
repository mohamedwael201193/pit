const CARDS = [
  { title: "Where policy lives", value: "PIT Desktop" },
  { title: "What this page is", value: "Read-only" },
  { title: "Clip / leverage", value: "Host law on the machine" },
  { title: "Allowed venues", value: "hyperliquid" },
  { title: "Session", value: "order + cancel only" },
  { title: "Withdraw / transfer", value: "impossible through PIT" },
];

export function PolicyPanel() {
  return (
    <div className="mt-8">
      <h2 className="text-[1.25rem] font-semibold tracking-[-0.03em]">Your policy is the law</h2>
      <p className="mt-3 max-w-[48ch] text-[0.975rem] leading-7 text-[rgb(240_231_212/0.75)]">
        Numbers on this website are not the pinned host law. Preview, confirm, and pin on PIT Desktop. Chat cannot pin.
      </p>
      <dl className="mt-5 grid grid-cols-1 gap-3 sm:grid-cols-2">
        {CARDS.map((c) => (
          <div key={c.title} className="border border-[rgb(240_231_212/0.22)] bg-[#141414] p-5">
            <dt className="text-[0.8125rem] text-[rgb(240_231_212/0.55)]">{c.title}</dt>
            <dd className="mt-2 font-mono text-[1.125rem]">{c.value}</dd>
          </div>
        ))}
      </dl>
      <p className="mt-4 max-w-[40ch] text-[0.875rem] text-[rgb(240_231_212/0.7)]">
        The model cannot raise clip, leverage, or permissions.
      </p>
    </div>
  );
}
