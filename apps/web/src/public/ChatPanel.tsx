import { useState } from "react";
import { Link } from "react-router-dom";
import { answerChat, STARTERS, type ChatTurn } from "./chat";
import { useWatch } from "./Watch";

export function ChatPanel() {
  const { watch } = useWatch();
  const [open, setOpen] = useState(false);
  const [input, setInput] = useState("");
  const [turns, setTurns] = useState<ChatTurn[]>([]);

  const ask = (q: string) => {
    const text = q.trim();
    if (!text) return;
    setTurns((t) => [...t, { q: text, a: answerChat(text, watch) }]);
    setInput("");
    setOpen(true);
  };

  return (
    <>
      <button type="button" className="intel-ask" onClick={() => setOpen((v) => !v)} aria-expanded={open}>
        {open ? "Close" : "Ask PIT"}
      </button>
      {open ? (
        <aside className="intel-chat" aria-label="Public intelligence chat">
          <header className="mb-4">
            <p className="intel-kicker">Public intelligence</p>
            <p className="mt-1 text-[0.875rem] leading-5 text-[rgb(240_231_212/0.6)]">
              Informational. No wallet signing. No authorize. No policy pin. No autonomy.
            </p>
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
            <Link to="/download" className="text-[#d82f2f]">
              Open PIT Desktop
            </Link>
          </p>
        </aside>
      ) : null}
    </>
  );
}
