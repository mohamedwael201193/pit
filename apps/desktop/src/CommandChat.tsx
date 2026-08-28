import { FormEvent, useEffect, useRef, useState } from "react";
import { sendDeskCommand, fetchChatLog, type ChatMessage } from "./companion";

export function CommandChat({
  onNavigate,
  onResearch,
  onOpenPreview,
}: {
  onNavigate: (view: string) => void;
  onResearch: (coin: string) => void;
  onOpenPreview: () => void;
}) {
  const [draft, setDraft] = useState("");
  const [busy, setBusy] = useState(false);
  const [lines, setLines] = useState<ChatMessage[]>([]);
  const end = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let gone = false;
    void (async () => {
      const log = await fetchChatLog();
      if (!gone && log.length) setLines(log);
    })();
    return () => {
      gone = true;
    };
  }, []);

  useEffect(() => {
    end.current?.scrollIntoView({ block: "end" });
  }, [lines]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const text = draft.trim();
    if (!text || busy) return;
    setDraft("");
    setBusy(true);
    setLines((cur) => [...cur, { role: "user", text, ts: Date.now() }]);
    try {
      const r = await sendDeskCommand(text);
      const reply = r.reply || "PIT did not execute.";
      setLines((cur) => [...cur, { role: "pit", text: reply, tool: r.tool, ts: Date.now() }]);
      if (r.open_url) window.open(r.open_url, "_blank", "noopener,noreferrer");
      if (r.navigate === "preview") onOpenPreview();
      else if (r.navigate) onNavigate(r.navigate);
      if (r.start_research && r.coin) onResearch(r.coin);
    } catch {
      setLines((cur) => [
        ...cur,
        { role: "pit", text: "Local PIT did not answer. The sealed job is not cancelled by a missed chat poll.", ts: Date.now() },
      ]);
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="command" aria-label="Trading desk command">
      <p className="label">TRADING DESK COMMAND</p>
      <p className="fine">
        Host-parsed. Chat never sends your private book to a model. Research still uses Direct TeeML. Chat cannot
        AUTHORIZE, size, or change policy.
      </p>
      <div className="transcript" role="log">
        {lines.length === 0 ? (
          <p className="fine">Try: Research ETH privately. Why is ETH interesting? Show me the evidence. Open Hyperliquid.</p>
        ) : (
          lines.map((m, i) => (
            <p key={`${m.ts}-${i}`} className={m.role === "user" ? "turn user" : "turn pit"}>
              <strong>{m.role === "user" ? "You" : "PIT"}</strong> {m.text}
            </p>
          ))
        )}
        <div ref={end} />
      </div>
      <form onSubmit={(e) => void onSubmit(e)}>
        <label htmlFor="desk-cmd" className="label">
          Command
        </label>
        <textarea
          id="desk-cmd"
          rows={2}
          value={draft}
          disabled={busy}
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              e.currentTarget.form?.requestSubmit();
            }
          }}
          placeholder="Research ETH privately."
        />
        <button type="submit" className="linkish" disabled={busy || !draft.trim()}>
          {busy ? "Working…" : "Send"}
        </button>
      </form>
    </section>
  );
}
