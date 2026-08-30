import { FormEvent, useEffect, useRef, useState } from "react";
import {
  streamDeskCommand,
  pickChatModel,
  fetchChatLog,
  fetchModelCatalog,
  type ActivityEvent,
  type BindResult,
  type ChatMessage,
  type ChatReply,
  type DirectModel,
  type LocalStatus,
} from "./companion";
import { openExternal } from "./open";
import { AgentRun, CHAT_AGENT_COPY } from "./AgentRun";
import type { MarketCoin } from "./WatchBook";

const PRIMARY = [
  { label: "Find best opportunity", q: "Find the best opportunity available right now." },
  { label: "Find best long", q: "Find me the best long" },
  { label: "Scan all markets", q: "Scan all markets" },
  { label: "What can I trade now?", q: "What can I trade with my current capital?" },
  { label: "Why didn't PIT trade?", q: "Why didn't you trade?" },
  { label: "While I sleep", q: "Trade this while I sleep" },
];

const MORE = [
  { label: "Find best short", q: "Find me the best short" },
  { label: "Research BTC", q: "Research BTC" },
  { label: "Research ETH", q: "Research ETH" },
  { label: "Compare top", q: "Compare top opportunities" },
  { label: "Explain my policy", q: "Explain my policy" },
  { label: "What is executable?", q: "What is executable?" },
  { label: "Show current exposure", q: "Show current exposure" },
  { label: "Show proof", q: "Show proof" },
  { label: "Review Sleep Mission", q: "Review Sleep Mission" },
];

function huntLike(text: string) {
  return /find (me )?(the )?(best|strongest)|scan all markets|what can i trade|research (btc|eth|avax|sol|the strongest)|research the strongest/i.test(
    text,
  );
}

export function CommandChat({
  thread,
  onNavigate,
  onResearch,
  onOpenPreview,
  onStop,
  onConfirmAutonomy,
  onTradeNow,
  island,
  coins,
  scanned,
  buyingPower,
  powerSource,
  watchAgeMs,
  bestWhy,
  activity,
  preview,
  previewHash,
  lastOrder,
  lastOid,
  evidence,
  researchNote,
  researchStop,
  huntRejected,
  huntSurvived,
  pinned,
  sessionAlive,
  autonomy,
  authBusy,
  authErr,
}: {
  thread: string;
  onNavigate: (view: string) => void;
  onResearch: (coin: string, hypothesis?: "none" | "long" | "short") => void;
  onOpenPreview: () => void;
  onConfirmAutonomy?: (hours: number) => void;
  onStop?: () => void;
  onTradeNow: () => void;
  island?: {
    busy: boolean;
    coin: string;
    stage: string;
    elapsedMs: number;
    jobId: string;
    pollMiss?: boolean;
    updatedAt?: number;
    roles: Array<{ role?: string; verify_e2ee?: string; proposed_side?: string; survives?: boolean; kill?: boolean }>;
    kind?: string;
  };
  coins?: MarketCoin[];
  scanned?: number;
  buyingPower?: number;
  powerSource?: string;
  watchAgeMs?: number;
  bestWhy?: string;
  activity?: ActivityEvent[];
  preview?: BindResult["preview"] | null;
  previewHash?: string;
  lastOrder?: LocalStatus["lastOrder"] | null;
  lastOid?: string;
  evidence?: unknown;
  researchNote?: string | null;
  researchStop?: string | null;
  huntRejected?: string[];
  huntSurvived?: string;
  pinned?: boolean;
  sessionAlive?: boolean;
  autonomy?: string;
  authBusy?: boolean;
  authErr?: string | null;
}) {
  const [draft, setDraft] = useState("");
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [lines, setLines] = useState<ChatMessage[]>([]);
  const [privateModels, setPrivate] = useState<DirectModel[]>([]);
  const [otherModels, setOther] = useState<DirectModel[]>([]);
  const [unsupported, setUnsupported] = useState<DirectModel[]>([]);
  const [catalog, setCatalog] = useState<DirectModel[]>([]);
  const [catalogNote, setCatalogNote] = useState("");
  const [modelOpen, setModelOpen] = useState(false);
  const [moreOpen, setMoreOpen] = useState(false);
  const [picked, setPicked] = useState<DirectModel | null>(null);
  const end = useRef<HTMLDivElement>(null);
  const log = useRef<HTMLDivElement>(null);
  const abort = useRef<AbortController | null>(null);

  useEffect(() => {
    let gone = false;
    void (async () => {
      const got = await fetchChatLog(thread);
      if (!gone) setLines(got);
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
      setCatalog(c.official_catalog || []);
      setCatalogNote(c.catalog_note || "");
      const host = c.other_chat.find((m) => m.model === "host-parsed");
      const remembered = [...c.other_chat, ...c.private_verified].find((m) => m.model === c.picked);
      setPicked(remembered || host || c.private_verified[0] || c.models[0] || null);
    });
    return () => {
      gone = true;
    };
  }, []);

  useEffect(() => {
    const box = log.current;
    if (box) box.scrollTop = box.scrollHeight;
    else end.current?.scrollIntoView({ block: "end" });
  }, [lines, busy]);

  async function ask(text: string) {
    if (!text || busy) return;
    if (island?.busy && huntLike(text) && !/why|stop|proof|reject/i.test(text)) {
      setLines((cur) => [
        ...cur,
        { role: "user", text, ts: Date.now(), thread },
        {
          role: "pit",
          text: `Still researching ${island.coin}. Live stages stay on this screen.`,
          ts: Date.now(),
          thread,
        },
      ]);
      return;
    }
    setBusy(true);
    setErr(null);
    setMoreOpen(false);
    setLines((cur) => [...cur, { role: "user", text, ts: Date.now(), thread }]);
    abort.current?.abort();
    const ac = new AbortController();
    abort.current = ac;
    const streamIdx = { n: -1 };
    try {
      setLines((cur) => {
        streamIdx.n = cur.length;
        return [...cur, { role: "pit", text: "", ts: Date.now(), thread }];
      });
      const r = await streamDeskCommand(text, thread, picked?.model || "host-parsed", ac.signal, (delta) => {
        setLines((cur) => {
          const next = [...cur];
          const i = streamIdx.n >= 0 ? streamIdx.n : next.length - 1;
          if (next[i]?.role === "pit") next[i] = { ...next[i], text: (next[i].text || "") + delta };
          return next;
        });
      });
      if (ac.signal.aborted) return;
      let reply = r.reply || r.error || "PIT could not complete that command.";
      if (r.tool === "research.status") {
        reply = `Still researching ${r.coin || island?.coin || "the live book"}. Watch the stages above.`;
      }
      setLines((cur) => {
        const next = [...cur];
        const i = streamIdx.n >= 0 ? streamIdx.n : next.length - 1;
        if (next[i]?.role === "pit") next[i] = { ...next[i], text: reply, tool: r.tool, coin: r.coin };
        return next;
      });
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
    if (r.execute || r.sign || r.trade) return;
    if (r.open_url) void openExternal(r.open_url);
    if (r.tool === "mission.enable_required") {
      onConfirmAutonomy?.(r.hours || 8);
      return;
    }
    if (r.start_research) {
      const hyp = r.hypothesis === "long" || r.hypothesis === "short" ? r.hypothesis : undefined;
      onResearch(r.coin || "", hyp);
      return;
    }
    if (r.navigate === "preview") onOpenPreview();
    if (r.navigate === "security") onNavigate("security");
    if (r.navigate === "automation") onNavigate("automation");
    if (r.navigate === "setup") onNavigate("setup");
  }

  const researchSku = privateModels[0];

  return (
    <section className="command cockpit-shell" aria-label="PIT agent">
      <div className="command-head">
        <div>
          <p className="label">PIT AGENT</p>
          <div className="honesty-row">
            <span className="honesty-chip">Private research</span>
            <span className="honesty-chip">{CHAT_AGENT_COPY.cannotAuthorize}</span>
            <span className="honesty-chip">Cannot pin</span>
          </div>
        </div>
        <div className="model-pick">
          <button type="button" aria-haspopup="listbox" aria-expanded={modelOpen} onClick={() => setModelOpen((v) => !v)}>
            PIT AGENT · {picked?.private_book ? "Private + Verified" : "Desk command"}
          </button>
          {modelOpen ? (
            <div className="model-menu wide" role="listbox">
              <p className="label">PIT PRIVATE · Recommended</p>
              {privateModels.length === 0 ? (
                <p className="fine">No verified Direct SKU on this network.</p>
              ) : (
                privateModels.map((m) => (
                  <ModelRow key={m.model} m={m} picked={picked} onPick={() => { setPicked(m); setModelOpen(false); void pickChatModel(m.model || "host-parsed"); }} />
                ))
              )}
              <p className="label">FAST · desk command</p>
              {otherModels.map((m) => (
                <ModelRow key={m.model} m={m} picked={picked} onPick={() => { setPicked(m); setModelOpen(false); void pickChatModel(m.model || "host-parsed"); }} />
              ))}
              <p className="label">AVAILABLE PROVIDERS · listing only</p>
              <details className="catalog-disclosure">
                <summary>Show listings — not used for chat or private research</summary>
                <p className="fine">{catalogNote || "Listed SKUs are not inference paths. Private book stays Direct TeeML."}</p>
                {catalog.slice(0, 31).map((m) => (
                  <ModelRow key={m.model} m={m} picked={null} onPick={() => undefined} disabled />
                ))}
              </details>
              {unsupported.map((m) => (
                <ModelRow key={m.model} m={m} picked={picked} onPick={() => undefined} disabled />
              ))}
              <p className="fine">Desk commands stay host-parsed. Direct TeeML is used only for sealed research. Catalog presence is not privacy.</p>
            </div>
          ) : null}
        </div>
      </div>

      <div className="cockpit-live">
        <AgentRun
          busy={Boolean(island?.busy)}
          coin={island?.coin || ""}
          stage={island?.stage || ""}
          elapsedMs={island?.elapsedMs || 0}
          jobId={island?.jobId || ""}
          pollMiss={Boolean(island?.pollMiss)}
          updatedAt={island?.updatedAt}
          roles={island?.roles || []}
          kind={island?.kind || ""}
          coins={coins || []}
          scanned={scanned || 0}
          buyingPower={buyingPower}
          powerSource={powerSource}
          watchAgeMs={watchAgeMs}
          bestWhy={bestWhy}
          preview={preview || null}
          previewHash={previewHash}
          lastOrder={lastOrder}
          lastOid={lastOid}
          activity={activity || []}
          evidence={evidence}
          researchNote={researchNote}
          researchStop={researchStop}
          huntRejected={huntRejected || []}
          huntSurvived={huntSurvived || ""}
          pinned={Boolean(pinned)}
          sessionAlive={Boolean(sessionAlive)}
          autonomy={autonomy || "manual"}
          researchSku={researchSku}
          authBusy={authBusy}
          authErr={authErr}
          onAsk={ask}
          onOpenPreview={onOpenPreview}
          onOpenPolicy={() => onNavigate("security")}
          onOpenAutomation={() => onNavigate("automation")}
          onOpenActivity={() => onNavigate("activity")}
          onStop={() => onStop?.()}
          onTradeNow={onTradeNow}
        />
      </div>

      <div className="transcript" role="log" ref={log}>
        {lines.length === 0 ? (
          <p className="chat-empty">Ask PIT what to trade. It will scan live books, research privately, and wait for TRADE NOW on this computer.</p>
        ) : null}
        {lines.map((m, i) => (
          <article key={`${m.ts}-${i}`} className={`turn ${m.role === "user" ? "user" : "pit"}`}>
            <div className="turn-meta">
              <span className="who">{m.role === "user" ? "You" : "PIT"}</span>
              {m.ts ? <time>{new Date(m.ts).toLocaleTimeString()}</time> : null}
            </div>
            <TurnBody text={displayTurn(m)} />
          </article>
        ))}
        <div ref={end} />
      </div>

      <form className="composer" onSubmit={onSubmit}>
        <div className="prompt-chips">
          {PRIMARY.map((p) => (
            <button
              key={p.label}
              type="button"
              className="chip-btn"
              disabled={Boolean(island?.busy && huntLike(p.q))}
              onClick={() => void ask(p.q)}
            >
              {p.label}
            </button>
          ))}
          <button type="button" className="chip-btn" onClick={() => setMoreOpen((v) => !v)}>More</button>
        </div>
        {moreOpen ? (
          <div className="prompt-chips more">
            {MORE.map((p) => (
              <button key={p.label} type="button" className="chip-btn" onClick={() => void ask(p.q)}>{p.label}</button>
            ))}
          </div>
        ) : null}
        <textarea
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          placeholder="Ask PIT what to trade, research, explain, or watch…"
          rows={3}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              const text = draft.trim();
              if (!text) return;
              setDraft("");
              void ask(text);
            }
          }}
        />
        <div className="composer-row">
          <p className="fine">Enter send · Shift+Enter newline · Esc stop research · Ctrl+K command</p>
          <button type="submit" className="primary" disabled={busy || !draft.trim()}>Send</button>
        </div>
        {err ? <p className="fine">{err}</p> : null}
      </form>
    </section>
  );
}

function displayTurn(m: ChatMessage) {
  const text = String(m.text || "").trim();
  if (!text) return m.role === "pit" ? "Working…" : "";
  return text;
}

function TurnBody({ text }: { text: string }) {
  const parts = String(text || "").split(/\n+/).filter(Boolean);
  if (!parts.length) return <p className="turn-body">Working…</p>;
  return (
    <div className="turn-body">
      {parts.map((ln, i) => (
        <p key={i}>{ln}</p>
      ))}
    </div>
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
  const on = picked?.model === m.model;
  return (
    <button type="button" role="option" aria-selected={on} disabled={disabled} onClick={onPick}>
      <strong>{m.label || m.model}</strong>
      <span>{m.path || m.capability || ""}</span>
    </button>
  );
}
