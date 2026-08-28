import { FormEvent, useEffect, useRef, useState } from "react";
import {
  sendDeskCommand,
  fetchChatLog,
  fetchModels,
  type ChatMessage,
  type ChatReply,
  type DirectModel,
} from "./companion";

export function CommandChat({
  thread,
  onNavigate,
  onResearch,
  onOpenPreview,
  onStop,
  island,
}: {
  thread: string;
  onNavigate: (view: string) => void;
  onResearch: (coin: string) => void;
  onOpenPreview: () => void;
  onStop?: () => void;
  island?: {
    busy: boolean;
    coin: string;
    stage: string;
    elapsedMs: number;
    jobId: string;
    pollMiss?: boolean;
    roles: Array<{ role?: string; verify_e2ee?: string }>;
    kind?: string;
  };
}) {
  const [draft, setDraft] = useState("");
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [lines, setLines] = useState<ChatMessage[]>([]);
  const [models, setModels] = useState<DirectModel[]>([]);
  const [modelOpen, setModelOpen] = useState(false);
  const [picked, setPicked] = useState<DirectModel | null>(null);
  const end = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let gone = false;
    void (async () => {
      const log = await fetchChatLog(thread);
      if (!gone) setLines(log);
    })();
    return () => {
      gone = true;
    };
  }, [thread]);

  useEffect(() => {
    let gone = false;
    void fetchModels().then((m) => {
      if (gone) return;
      setModels(m);
      setPicked(m[0] || null);
    });
    return () => {
      gone = true;
    };
  }, []);

  useEffect(() => {
    end.current?.scrollIntoView({ block: "end" });
  }, [lines, island?.stage, island?.elapsedMs]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const text = draft.trim();
    if (!text || busy) return;
    setDraft("");
    setBusy(true);
    setErr(null);
    setLines((cur) => [...cur, { role: "user", text, ts: Date.now(), thread }]);
    try {
      const r = await sendDeskCommand(text, thread);
      const reply = r.reply || r.error || "PIT could not complete that command.";
      setLines((cur) => [...cur, { role: "pit", text: reply, tool: r.tool, ts: Date.now(), thread, coin: r.coin }]);
      applyReply(r);
    } catch {
      setErr("Local PIT did not answer. A sealed job is not cancelled by a missed chat poll.");
      setLines((cur) => [
        ...cur,
        { role: "pit", text: "Local PIT did not answer. The sealed job is not cancelled by a missed chat poll.", ts: Date.now() },
      ]);
    } finally {
      setBusy(false);
    }
  }

  function applyReply(r: ChatReply) {
    if (r.open_url) window.open(r.open_url, "_blank", "noopener,noreferrer");
    if (r.navigate === "preview") onOpenPreview();
    else if (r.navigate) onNavigate(r.navigate);
    if (r.start_research && r.coin) onResearch(r.coin);
  }

  const modelLabel = picked
    ? `${picked.private_book ? "Private + Verified" : picked.label || "Direct"} · ${picked.model}`
    : "Direct";

  return (
    <section className="command" aria-label="Trading desk command">
      <div className="command-head">
        <div>
          <p className="label">Command</p>
          <p className="fine">Host-parsed. Private book never leaves this computer. Chat cannot AUTHORIZE.</p>
        </div>
        <div className="model-pick">
          <button type="button" aria-haspopup="listbox" aria-expanded={modelOpen} onClick={() => setModelOpen((v) => !v)}>
            {modelLabel}
          </button>
          {modelOpen ? (
            <div className="model-menu" role="listbox">
              <p className="label">Private research models</p>
              {models.length === 0 ? (
                <p className="fine">No Direct-compatible model is listed for this network.</p>
              ) : (
                models.map((m) => (
                  <button
                    key={m.model}
                    type="button"
                    className={picked?.model === m.model ? "on" : "linkish"}
                    onClick={() => {
                      setPicked(m);
                      setModelOpen(false);
                    }}
                  >
                    {m.model}
                    <span className="fine" style={{ display: "block", marginTop: 2 }}>
                      {m.path || "Direct"} · {m.private_book ? "Private" : "not private"} ·{" "}
                      {m.proven_e2ee ? "Verified" : "unproven"} · {m.verifiability || "TeeML"}
                      {m.provider
                        ? ` · ${m.provider.length > 12 ? `${m.provider.slice(0, 6)}…${m.provider.slice(-4)}` : m.provider}`
                        : ""}
                    </span>
                    {m.latency ? (
                      <span className="fine" style={{ display: "block" }}>
                        {m.latency}
                      </span>
                    ) : null}
                    {m.cost ? (
                      <span className="fine" style={{ display: "block" }}>
                        {m.cost}
                      </span>
                    ) : null}
                  </button>
                ))
              )}
              <p className="fine">Only verified Direct models appear. Router SKUs are not listed as private.</p>
            </div>
          ) : null}
        </div>
      </div>
      <div className="transcript" role="log">
        {lines.length === 0 ? (
          <p className="fine">Try: What is happening? Research ETH privately. Show my positions. Explain my policy.</p>
        ) : (
          lines.map((m, i) => (
            <div key={`${m.ts}-${i}`} className={m.role === "user" ? "turn user" : "turn pit"}>
              <div className="turn-meta">
                <span className="who">{m.role === "user" ? "You" : "PIT"}</span>
                {m.ts ? <time dateTime={new Date(m.ts).toISOString()}>{new Date(m.ts).toLocaleTimeString()}</time> : null}
                {m.role === "pit" ? (
                  <button
                    type="button"
                    className="copy-turn"
                    onClick={() => void navigator.clipboard.writeText(m.text || "")}
                  >
                    Copy
                  </button>
                ) : null}
              </div>
              <p style={{ margin: 0 }}>{m.text}</p>
              {m.tool ? <ChatCard tool={m.tool} coin={m.coin} onNavigate={onNavigate} onOpenPreview={onOpenPreview} /> : null}
            </div>
          ))
        )}
        {island?.busy ? (
          <article className="chat-card" role="status">
            <p className="label">Researching {island.coin}</p>
            <p>
              {island.stage.replaceAll("_", " ")} · {(island.elapsedMs / 1000).toFixed(1)}s elapsed
            </p>
            {island.pollMiss ? <p role="status">Live view delayed — research is still running.</p> : null}
            <p className="fine">
              Researcher {roleMark(island.roles, "researcher")} · Challenger {roleMark(island.roles, "challenger")} · Risk{" "}
              {roleMark(island.roles, "risk")}
            </p>
            {island.jobId ? <p className="fine">Job {island.jobId}</p> : null}
            <button type="button" className="linkish" onClick={onOpenPreview}>
              Open research
            </button>
          </article>
        ) : null}
        {err ? (
          <p className="err" role="alert">
            {err}
          </p>
        ) : null}
        <div ref={end} />
      </div>
      <form className="composer" onSubmit={(e) => void onSubmit(e)}>
        <textarea
          id="desk-cmd"
          rows={2}
          value={draft}
          disabled={busy}
          aria-label="Ask PIT"
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              e.currentTarget.form?.requestSubmit();
            }
          }}
          placeholder="What is happening?"
        />
        <div className="composer-row">
          <p className="fine" style={{ margin: 0 }}>
            Enter send · Shift+Enter newline
          </p>
          {island?.busy && onStop ? (
            <button type="button" className="linkish" onClick={onStop}>
              Stop research
            </button>
          ) : null}
          <button type="submit" className="primary" disabled={busy || !draft.trim()}>
            {busy ? "Working…" : "Send"}
          </button>
        </div>
      </form>
    </section>
  );
}

function roleMark(roles: Array<{ role?: string; verify_e2ee?: string }>, name: string) {
  return roles.some((r) => String(r.role).toLowerCase() === name && String(r.verify_e2ee).toUpperCase() === "OK")
    ? "verified"
    : "pending";
}

function ChatCard({
  tool,
  coin,
  onNavigate,
  onOpenPreview,
}: {
  tool?: string;
  coin?: string;
  onNavigate: (view: string) => void;
  onOpenPreview: () => void;
}) {
  if (!tool || tool === "help" || tool === "status") return null;
  if (tool === "research.start") {
    return (
      <article className="chat-card">
        <p className="label">Research</p>
        <p>Sealed Direct pass for {coin || "ETH"}. Compute money, not trading capital.</p>
        <button type="button" className="linkish" onClick={onOpenPreview}>
          Open research
        </button>
      </article>
    );
  }
  if (tool === "preview.show" || tool === "refuse_execute") {
    return (
      <article className="chat-card">
        <p className="label">Exact preview</p>
        <p>Review happens on this computer. Chat cannot AUTHORIZE.</p>
        <button type="button" className="linkish" onClick={onOpenPreview}>
          Review preview
        </button>
      </article>
    );
  }
  if (tool === "watch.get") {
    return (
      <article className="chat-card">
        <p className="label">Market</p>
        <button type="button" className="linkish" onClick={() => onNavigate("watch")}>
          Open Watch
        </button>
      </article>
    );
  }
  if (tool === "positions.get") {
    return (
      <article className="chat-card">
        <p className="label">Positions</p>
        <button type="button" className="linkish" onClick={() => onNavigate("positions")}>
          Open Positions
        </button>
      </article>
    );
  }
  if (tool === "activity.list") {
    return (
      <article className="chat-card">
        <p className="label">Activity</p>
        <button type="button" className="linkish" onClick={() => onNavigate("activity")}>
          Open Activity
        </button>
      </article>
    );
  }
  if (tool === "policy.get") {
    return (
      <article className="chat-card">
        <p className="label">Policy</p>
        <p>Host enforced. Chat cannot mutate it.</p>
        <button type="button" className="linkish" onClick={() => onNavigate("policy")}>
          Open Policy
        </button>
      </article>
    );
  }
  if (tool === "session.status") {
    return (
      <article className="chat-card">
        <p className="label">Hyperliquid</p>
        <button type="button" className="linkish" onClick={() => onNavigate("security")}>
          Open Security
        </button>
      </article>
    );
  }
  if (tool === "research.result") {
    return (
      <article className="chat-card">
        <p className="label">Evidence</p>
        <button type="button" className="linkish" onClick={onOpenPreview}>
          Open research
        </button>
      </article>
    );
  }
  return null;
}
