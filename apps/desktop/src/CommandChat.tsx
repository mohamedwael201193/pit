import { FormEvent, useEffect, useRef, useState } from "react";
import {
  sendDeskCommand,
  fetchChatLog,
  fetchModelCatalog,
  type ChatMessage,
  type ChatReply,
  type DirectModel,
} from "./companion";

const PROMPTS = [
  "What is happening?",
  "What is interesting right now?",
  "What is the ETH setup?",
  "Research BTC privately.",
  "Why did the committee reject this?",
  "What is my current exposure?",
  "Show me today's opportunities.",
  "Why can't PIT execute this?",
];

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
  const [privateModels, setPrivate] = useState<DirectModel[]>([]);
  const [otherModels, setOther] = useState<DirectModel[]>([]);
  const [unsupported, setUnsupported] = useState<DirectModel[]>([]);
  const [modelOpen, setModelOpen] = useState(false);
  const [picked, setPicked] = useState<DirectModel | null>(null);
  const [lastUser, setLastUser] = useState("");
  const end = useRef<HTMLDivElement>(null);
  const abort = useRef<AbortController | null>(null);

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
    void fetchModelCatalog().then((c) => {
      if (gone) return;
      setPrivate(c.private_verified);
      setOther(c.other_chat);
      setUnsupported(c.unsupported);
      const host = c.other_chat.find((m) => m.model === "host-parsed");
      setPicked(host || c.private_verified[0] || c.models[0] || null);
    });
    return () => {
      gone = true;
    };
  }, []);

  useEffect(() => {
    end.current?.scrollIntoView({ block: "end" });
  }, [lines, island?.stage, island?.elapsedMs, busy]);

  async function ask(text: string) {
    if (!text || busy) return;
    setBusy(true);
    setErr(null);
    setLastUser(text);
    setLines((cur) => [...cur, { role: "user", text, ts: Date.now(), thread }]);
    abort.current?.abort();
    const ac = new AbortController();
    abort.current = ac;
    try {
      const r = await sendDeskCommand(text, thread, ac.signal);
      if (ac.signal.aborted) return;
      const reply = r.reply || r.error || "PIT could not complete that command.";
      setLines((cur) => [...cur, { role: "pit", text: reply, tool: r.tool, ts: Date.now(), thread, coin: r.coin }]);
      applyReply(r);
    } catch (e) {
      if (ac.signal.aborted) return;
      const msg = e instanceof Error && e.name === "AbortError" ? "Stopped." : "Local PIT did not answer. A sealed job is not cancelled by a missed chat poll.";
      setErr(msg);
      setLines((cur) => [...cur, { role: "pit", text: msg, ts: Date.now() }]);
    } finally {
      setBusy(false);
    }
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const text = draft.trim();
    if (!text) return;
    setDraft("");
    await ask(text);
  }

  function applyReply(r: ChatReply) {
    if (r.open_url) window.open(r.open_url, "_blank", "noopener,noreferrer");
    if (r.navigate === "preview") onOpenPreview();
    else if (r.navigate) onNavigate(r.navigate);
    if (r.start_research && r.coin) onResearch(r.coin);
  }

  const researchSku = privateModels[0];
  const modelLabel = picked?.model === "host-parsed"
    ? `Chat · host-parsed${researchSku ? ` · Research ${researchSku.model}` : ""}`
    : picked?.private_book
      ? `Research · Private + Verified · ${picked.model}`
      : `${picked?.label || "Desk"} · ${picked?.model || "host-parsed"}`;

  return (
    <section className="command" aria-label="Trading desk command">
      <div className="command-head">
        <div>
          <p className="label">Chat</p>
          <p className="fine">Host-parsed on this computer. Private book never leaves. Chat cannot AUTHORIZE.</p>
        </div>
        <div className="model-pick">
          <button type="button" aria-haspopup="listbox" aria-expanded={modelOpen} onClick={() => setModelOpen((v) => !v)}>
            {modelLabel}
          </button>
          {modelOpen ? (
            <div className="model-menu" role="listbox">
              <p className="label">Private + verified</p>
              {privateModels.length === 0 ? (
                <p className="fine">No verified Direct SKU on this network.</p>
              ) : (
                privateModels.map((m) => (
                  <ModelRow key={m.model} m={m} picked={picked} onPick={() => { setPicked(m); setModelOpen(false); }} />
                ))
              )}
              <p className="label">General chat</p>
              {otherModels.map((m) => (
                <ModelRow key={m.model} m={m} picked={picked} onPick={() => { setPicked(m); setModelOpen(false); }} />
              ))}
              <p className="label">Unsupported for private research</p>
              {unsupported.map((m) => (
                <ModelRow key={m.model} m={m} picked={picked} onPick={() => { setPicked(m); setModelOpen(false); }} disabled />
              ))}
              <p className="fine">
                Desk chat is host-parsed. glm-5.2 is used only on the proven Direct path for the sealed book. Catalog presence is not privacy.
              </p>
            </div>
          ) : null}
        </div>
      </div>
      <div className="transcript" role="log">
        {lines.length === 0 ? (
          <div className="chat-empty">
            <p>Ask the live desk. Answers come from this computer, not a canned script.</p>
            <div className="prompt-chips">
              {PROMPTS.map((p) => (
                <button key={p} type="button" className="chip-btn" onClick={() => void ask(p)}>
                  {p}
                </button>
              ))}
            </div>
          </div>
        ) : (
          lines.map((m, i) => (
            <div key={`${m.ts}-${i}`} className={m.role === "user" ? "turn user" : "turn pit"}>
              <div className="turn-meta">
                <span className="who">{m.role === "user" ? "You" : "PIT"}</span>
                {m.ts ? <time dateTime={new Date(m.ts).toISOString()}>{new Date(m.ts).toLocaleTimeString()}</time> : null}
                {m.role === "pit" ? (
                  <button type="button" className="copy-turn" onClick={() => void navigator.clipboard.writeText(m.text || "")}>
                    Copy
                  </button>
                ) : null}
              </div>
              <p style={{ margin: 0 }}>{m.text}</p>
              {m.tool ? <ChatCard tool={m.tool} coin={m.coin} onNavigate={onNavigate} onOpenPreview={onOpenPreview} /> : null}
            </div>
          ))
        )}
        {busy ? (
          <article className="chat-card" role="status">
            <p className="label">Working</p>
            <p>Host is reading live desk state. This is not a sealed model stream.</p>
          </article>
        ) : null}
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
          rows={3}
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
          {busy ? (
            <button
              type="button"
              className="linkish"
              onClick={() => {
                abort.current?.abort();
                setBusy(false);
              }}
            >
              Stop
            </button>
          ) : null}
          {lastUser && !busy ? (
            <button type="button" className="linkish" onClick={() => void ask(lastUser)}>
              Retry
            </button>
          ) : null}
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

function ModelRow({
  m,
  picked,
  onPick,
  disabled,
}: {
  m: DirectModel;
  picked: DirectModel | null;
  onPick: () => void;
  disabled?: boolean;
}) {
  return (
    <button type="button" className={picked?.model === m.model ? "on" : "linkish"} onClick={onPick} disabled={disabled}>
      {m.model}
      <span className="fine" style={{ display: "block", marginTop: 2 }}>
        {m.path || "Direct"} · {m.private_book ? "Private" : "not private"} · {m.proven_e2ee ? "Verified" : "unproven"}
        {m.capability ? ` · ${m.capability}` : ""}
      </span>
      {m.note ? (
        <span className="fine" style={{ display: "block" }}>
          {m.note}
        </span>
      ) : null}
    </button>
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
  if (!tool || tool === "help") return null;
  if (tool === "status" || tool === "greet") {
    return (
      <article className="chat-card">
        <p className="label">Desk</p>
        <button type="button" className="linkish" onClick={() => onNavigate("home")}>
          Open Desk
        </button>
      </article>
    );
  }
  if (tool === "setup.guide") {
    return (
      <article className="chat-card">
        <p className="label">Setup</p>
        <button type="button" className="linkish" onClick={() => onNavigate("setup")}>
          Open setup
        </button>
      </article>
    );
  }
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
        <p>Live Hyperliquid marks. Side is not decided on Watch.</p>
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
