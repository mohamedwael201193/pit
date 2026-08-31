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
  { label: "Find best short", q: "Find me the best short" },
  { label: "What can I trade?", q: "What can I trade?" },
  { label: "Compare top opportunities", q: "Compare top opportunities" },
];

function huntLike(text: string) {
  return /find (me )?(the )?(best|strongest|next)|scan all markets|what can i trade|research (next|btc|eth|avax|sol)|research the strongest|compare top|next opportunity/i.test(
    text,
  );
}

function huntDump(text: string) {
  const t = String(text || "").trim();
  if (!t) return false;
  if (/^(Researching |Still researching |Scanning live Hyperliquid)/i.test(t)) return true;
  if (/Live numbers stay on the cards/i.test(t)) return true;
  if (/Live stages stay on this screen/i.test(t)) return true;
  if (/Watch the (live )?stages/i.test(t)) return true;
  if (/is the strongest executable book among/i.test(t)) return true;
  if (t === "Working…") return true;
  return false;
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
  net,
}: {
  thread: string;
  onNavigate: (view: string) => void;
  onResearch: (coin: string, hypothesis?: "none" | "long" | "short", opts?: { fresh?: boolean }) => void;
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
  net?: string;
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
  const [picked, setPicked] = useState<DirectModel | null>(null);
  const end = useRef<HTMLDivElement>(null);
  const log = useRef<HTMLDivElement>(null);
  const stick = useRef(true);
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
    if (!box || !stick.current) return;
    box.scrollTop = box.scrollHeight;
  }, [lines, island?.busy]);

  async function ask(text: string) {
    if (!text || busy) return;
    setBusy(true);
    setErr(null);
    setLines((cur) => [...cur, { role: "user", text, ts: Date.now(), thread }]);
    abort.current?.abort();
    const ac = new AbortController();
    abort.current = ac;
    const streamIdx = { n: -1 };
    const hunt = huntLike(text);
    try {
      if (!hunt) {
        setLines((cur) => {
          streamIdx.n = cur.length;
          return [...cur, { role: "pit", text: "", ts: Date.now(), thread }];
        });
      }
      const r = await streamDeskCommand(text, thread, picked?.model || "host-parsed", ac.signal, (delta) => {
        if (hunt) return;
        setLines((cur) => {
          const next = [...cur];
          const i = streamIdx.n >= 0 ? streamIdx.n : next.length - 1;
          if (next[i]?.role === "pit") next[i] = { ...next[i], text: (next[i].text || "") + delta };
          return next;
        });
      });
      if (ac.signal.aborted) return;
      if (r.tool === "research.status" || r.start_research) {
        if (streamIdx.n >= 0) {
          setLines((cur) => {
            const next = [...cur];
            const row = next[streamIdx.n];
            if (row?.role === "pit" && (huntDump(row.text || "") || !row.text)) next.splice(streamIdx.n, 1);
            return next;
          });
        }
      applyReply(r, text);
      return;
    }
    const reply = r.reply || r.error || "PIT could not complete that command.";
    setLines((cur) => {
      const next = [...cur];
      if (streamIdx.n >= 0 && next[streamIdx.n]?.role === "pit") {
        next[streamIdx.n] = { ...next[streamIdx.n], text: reply, tool: r.tool, coin: r.coin };
        return next;
      }
      return [...next, { role: "pit", text: reply, tool: r.tool, coin: r.coin, ts: Date.now(), thread }];
    });
    applyReply(r, text);
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

  function applyReply(r: ChatReply, asked = "") {
    if (r.execute || r.sign || r.trade) return;
    if (r.open_url) void openExternal(r.open_url);
    if (r.tool === "mission.enable_required") {
      onConfirmAutonomy?.(r.hours || 8);
      return;
    }
    if (r.start_research) {
      const hyp = r.hypothesis === "long" || r.hypothesis === "short" ? r.hypothesis : "none";
      const unnamed = /find (me )?(the )?(best|strongest|next)|what can i trade|research next|next opportunity|next eligible|next market/i.test(asked);
      const nextHunt = /research next|next opportunity|next eligible|next market/i.test(asked);
      const fresh = unnamed && !nextHunt;
      onResearch(unnamed ? "" : r.coin || "", hyp, { fresh: fresh });
      return;
    }
    if (r.navigate === "preview") onOpenPreview();
    if (r.navigate === "security") onNavigate("security");
    if (r.navigate === "automation") onNavigate("automation");
    if (r.navigate === "setup") onNavigate("setup");
  }

  const researchSku = privateModels[0];
  const showMission = Boolean(
    island?.busy || island?.kind || island?.jobId || preview?.eligible || (researchNote && island?.coin),
  );
  const stream = composeStream(lines, showMission);
  const chips = lines.length === 0 ? PRIMARY : [];

  return (
    <section className="command agent-workspace" aria-label="PIT agent">
      <div className="command-head">
        <div className="model-pick">
          <button type="button" aria-haspopup="listbox" aria-expanded={modelOpen} onClick={() => setModelOpen((v) => !v)}>
            {picked?.private_book ? "Private + Verified" : "Desk command"}
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
                <summary>Show listings - not used for chat or private research</summary>
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

      <div className="agent-body">
        <div
          className={`agent-stream${island?.busy ? " is-busy" : ""}`}
          role="log"
          ref={log}
          onWheel={(e) => {
            if (e.deltaY < 0) stick.current = false;
          }}
          onScroll={() => {
            const box = log.current;
            if (!box) return;
            stick.current = box.scrollHeight - box.scrollTop - box.clientHeight < 96;
          }}
        >
          {lines.length === 0 && !showMission ? (
            <p className="chat-empty">Ask PIT what to trade. It will scan live books, research privately, and wait for TRADE NOW on this computer.</p>
          ) : null}
          {stream.map((row, i) =>
            row.mission ? (
              <article key={`mission-${i}`} className="turn pit mission">
                <div className="turn-meta">
                  <span className="who">PIT</span>
                </div>
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
                  net={net || "mainnet"}
                />
              </article>
            ) : (
              <article key={`${row.m?.ts}-${i}`} className={`turn ${row.m?.role === "user" ? "user" : "pit"}`}>
                <div className="turn-meta">
                  <span className="who">{row.m?.role === "user" ? "You" : "PIT"}</span>
                  {row.m?.ts ? <time>{new Date(row.m.ts).toLocaleTimeString()}</time> : null}
                </div>
                <TurnBody text={row.text || ""} />
              </article>
            ),
          )}
          <div ref={end} />
        </div>
      </div>

      <form className="composer" onSubmit={onSubmit}>
        {chips.length ? (
          <div className="prompt-chips">
            {chips.map((p) => (
              <button key={p.label} type="button" className="chip-btn" disabled={Boolean(island?.busy && huntLike(p.q))} onClick={() => void ask(p.q)}>
                {p.label}
              </button>
            ))}
          </div>
        ) : null}
        <textarea
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          placeholder="Ask PIT to research, compare, prepare, trade, or watch…"
          rows={3}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              const text = draft.trim();
              if (!text) return;
              setDraft("");
              void ask(text);
            }
            if (e.key === "Escape") onStop?.();
          }}
        />
        <div className="composer-row">
          <p className="fine">{CHAT_AGENT_COPY.cannotAuthorize} · Enter send · Shift+Enter newline · Esc stop research · Ctrl+K command</p>
          <button type="submit" className="primary" disabled={busy || !draft.trim()}>
            Send
          </button>
        </div>
        {err ? <p className="fine">{err}</p> : null}
      </form>
    </section>
  );
}

export function displayTurn(m: ChatMessage) {
  let text = String(m.text || "").trim();
  if (!text) return m.role === "pit" ? "Working…" : "";
  text = text
    .replace(/\s*Starting sealed 0G Direct on this computer\.?/gi, "")
    .replace(/\s*Chat cannot AUTHORIZE\.?/gi, "")
    .replace(/\s*The model cannot AUTHORIZE\.?/gi, "")
    .replace(/\s+/g, " ")
    .trim();
  const dump = text.match(/^([A-Z0-9]+) is the strongest executable book among .+/i);
  if (dump) return `Researching ${dump[1]}. Live numbers stay on the cards.`;
  return text || (m.role === "pit" ? "Working…" : "");
}

function visibleTurns(lines: ChatMessage[]) {
  const out: { m: ChatMessage; text: string }[] = [];
  for (const m of lines) {
    const text = displayTurn(m);
    if (!text) continue;
    if (huntDump(text) || huntDump(String(m.text || ""))) continue;
    const prev = out[out.length - 1];
    if (prev && prev.m.role === m.role && prev.text === text) continue;
    out.push({ m, text });
  }
  return out;
}

function lastHuntUserIndex(lines: ChatMessage[]) {
  for (let i = lines.length - 1; i >= 0; i--) {
    if (lines[i].role === "user" && huntLike(String(lines[i].text || ""))) return i;
  }
  return -1;
}

function composeStream(lines: ChatMessage[], showMission: boolean) {
  const vis = visibleTurns(lines);
  if (!showMission) return vis.map((row) => ({ ...row, mission: false as const }));
  if (!vis.length) return [{ mission: true as const }];
  const huntAt = lastHuntUserIndex(lines);
  const huntUser = huntAt >= 0 ? vis.find((row) => row.m === lines[huntAt]) : undefined;
  if (huntUser) return [{ ...huntUser, mission: false as const }, { mission: true as const }];
  const lastUser = [...vis].reverse().find((r) => r.m.role === "user");
  if (lastUser && huntLike(String(lastUser.m.text || ""))) {
    return [{ ...lastUser, mission: false as const }, { mission: true as const }];
  }
  return [{ mission: true as const }];
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
