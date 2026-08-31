import { useState } from "react";
import { answerChat, STARTERS, type ChatTurn } from "./chat";
import { windowsInstallerUrl } from "./facts";
import { useWatchSafe } from "./Watch";

export function ChatPanel({
  floating = true,
  open: openProp,
  onOpenChange,
}: {
  floating?: boolean;
  open?: boolean;
  onOpenChange?: (v: boolean) => void;
}) {
  const watch = useWatchSafe();
  const [inner, setInner] = useState(false);
  const [input, setInput] = useState("");
  const [turns, setTurns] = useState<ChatTurn[]>([]);
  const open = openProp ?? inner;

  const setOpen = (v: boolean) => {
    onOpenChange?.(v);
    if (openProp === undefined) setInner(v);
  };

  const ask = (q: string) => {
    const text = q.trim();
    if (!text) return;
    setTurns((t) => [...t, { q: text, a: answerChat(text, watch?.watch ?? null) }]);
    setInput("");
    setOpen(true);
  };

  return (
    <>
      {floating ? (
        <button type="button" className="intel-ask" onClick={() => setOpen(!open)} aria-expanded={open}>
          {open ? "Close" : "Ask PIT"}
        </button>
      ) : null}
      {open ? (
        <aside className={floating ? "intel-chat" : "desk-chat"} aria-label="Public intelligence chat">
          <header className="mb-4 flex items-start justify-between gap-3">
            <div>
              <p className="text-[0.75rem] tracking-[0.14em] text-[rgb(240_231_212/0.45)] uppercase">Public intelligence</p>
              <p className="mt-1 text-[0.875rem] leading-5 text-[rgb(240_231_212/0.6)]">
                Informational. No wallet signing. No authorize. No policy pin. No autonomy.
              </p>
            </div>
            <button
              type="button"
              className="rounded-full border border-[rgb(240_231_212/0.28)] px-3 py-1 text-[0.75rem] text-[var(--guide-cream)]"
              onClick={() => setOpen(false)}
            >
              Close
            </button>
          </header>
          <div className="flex flex-wrap gap-1.5">
            {STARTERS.map((s) => (
              <button key={s} type="button" className="intel-chip" onClick={() => ask(s)}>
                {s}
              </button>
            ))}
          </div>
          <div className="mt-4 grid max-h-[40vh] gap-4 overflow-y-auto">
            {turns.map((t, i) => (
              <div key={`${t.q}-${i}`}>
                <p className="text-[0.8125rem] font-medium">{t.q}</p>
                <p className="mt-1 text-[0.875rem] leading-6 text-[rgb(240_231_212/0.75)]">{t.a}</p>
              </div>
            ))}
          </div>
          <form
            className="mt-4 flex gap-2"
            onSubmit={(e) => {
              e.preventDefault();
              ask(input);
            }}
          >
            <input
              className="intel-input"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="What is happening?"
              aria-label="Ask PIT"
            />
            <button type="submit" className="intel-cta">
              Ask
            </button>
          </form>
          <p className="mt-3 text-[0.75rem] text-[rgb(240_231_212/0.45)]">
            Private command lives on desktop.{" "}
            <a href={windowsInstallerUrl()} className="text-[#d82f2f]">
              Download PIT Desktop
            </a>
          </p>
        </aside>
      ) : null}
    </>
  );
}
