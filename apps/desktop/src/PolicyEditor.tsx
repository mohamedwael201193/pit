import { useEffect, useState } from "react";
import { type HostPolicy } from "./companion";

const ASSETS = ["ETH", "BTC", "SOL", "HYPE", "DOGE", "AVAX"];

export function emptyPolicy(p?: Partial<HostPolicy> | null): HostPolicy {
  return {
    version: p?.version || "v1",
    maxClipUsd: p?.maxClipUsd ?? 10,
    dailyLossUsd: p?.dailyLossUsd ?? 50,
    maxLeverage: 1,
    allowedAssets: p?.allowedAssets?.length ? p.allowedAssets : [...ASSETS],
    allowedMarketTypes: ["perp"],
    allowedVenues: ["hyperliquid"],
    minSkillCalibration: p?.minSkillCalibration ?? 0,
    cooldownSeconds: p?.cooldownSeconds ?? 0,
    sessionTtlSeconds: p?.sessionTtlSeconds ?? 3600,
    killSwitch: Boolean(p?.killSwitch),
    maxUncertainty: p?.maxUncertainty ?? 1,
    maxSlippageBps: p?.maxSlippageBps ?? 80,
    minLiquidityUsd: p?.minLiquidityUsd ?? 0,
    maxOpenPositions: p?.maxOpenPositions ?? 1,
    maxConsecutiveLosses: p?.maxConsecutiveLosses ?? 3,
  };
}

export function PolicyEditor({
  current,
  consequences,
  pinned,
  busy,
  clipFloor = 10,
  clipCeil = 50,
  onPreview,
  onPin,
}: {
  current?: HostPolicy | null;
  consequences?: string[];
  pinned?: boolean;
  busy?: boolean;
  clipFloor?: number;
  clipCeil?: number;
  onPreview: (draft: HostPolicy) => void;
  onPin: (draft: HostPolicy) => void;
}) {
  const [draft, setDraft] = useState<HostPolicy>(() => emptyPolicy(current));
  const [localLines, setLocal] = useState<string[]>([
    "Edit a field, then preview. Pinning writes host law on this computer. Chat cannot do this.",
  ]);
  const pinKey = [
    current?.maxClipUsd,
    current?.maxOpenPositions,
    current?.dailyLossUsd,
    current?.maxSlippageBps,
    current?.cooldownSeconds,
    current?.sessionTtlSeconds,
    current?.maxUncertainty,
    current?.maxConsecutiveLosses,
    current?.minLiquidityUsd,
    current?.killSwitch,
    (current?.allowedAssets || []).join(","),
  ].join("|");
  useEffect(() => {
    setDraft(emptyPolicy(current));
  }, [pinKey]);

  const note = consequences && consequences.length ? consequences : localLines;

  function set<K extends keyof HostPolicy>(k: K, v: HostPolicy[K]) {
    setDraft((d) => ({ ...d, [k]: v, maxLeverage: 1, allowedVenues: ["hyperliquid"], allowedMarketTypes: ["perp"] }));
  }

  function toggleAsset(a: string) {
    setDraft((d) => {
      const has = d.allowedAssets.includes(a);
      const next = has ? d.allowedAssets.filter((x) => x !== a) : [...d.allowedAssets, a];
      return { ...d, allowedAssets: next.length ? next : [a], maxLeverage: 1 };
    });
  }

  return (
    <section className="policy-editor">
      <p className="label">Policy</p>
      <p className="fine">You edit. You pin. The model cannot. Leverage is locked at 1x. Withdraw stays impossible.</p>
      <div className="policy-grid">
        <label className="policy-cell">
          Max trade / position (USD)
          <input type="number" min={clipFloor} max={clipCeil} step={1} value={draft.maxClipUsd} onChange={(e) => set("maxClipUsd", Number(e.target.value))} />
          <span className="fine">${clipFloor}–${clipCeil}. Venue minimum is $10.</span>
        </label>
        <label className="policy-cell">
          Daily loss halt (USD)
          <input type="number" min={1} max={500} step={1} value={draft.dailyLossUsd} onChange={(e) => set("dailyLossUsd", Number(e.target.value))} />
        </label>
        <label className="policy-cell">
          Max open positions
          <input type="number" min={1} max={5} step={1} value={draft.maxOpenPositions} onChange={(e) => set("maxOpenPositions", Number(e.target.value))} />
          <span className="fine">Raising this does not flatten existing positions.</span>
        </label>
        <label className="policy-cell">
          Consecutive loss limit
          <input type="number" min={1} max={10} step={1} value={draft.maxConsecutiveLosses} onChange={(e) => set("maxConsecutiveLosses", Number(e.target.value))} />
        </label>
        <label className="policy-cell">
          Max slippage (bps)
          <input type="number" min={10} max={500} step={1} value={draft.maxSlippageBps} onChange={(e) => set("maxSlippageBps", Number(e.target.value))} />
        </label>
        <label className="policy-cell">
          Liquidity floor (USD)
          <input type="number" min={0} max={1000000} step={1} value={draft.minLiquidityUsd} onChange={(e) => set("minLiquidityUsd", Number(e.target.value))} />
        </label>
        <label className="policy-cell">
          Cooldown (seconds)
          <input type="number" min={0} max={86400} step={1} value={draft.cooldownSeconds} onChange={(e) => set("cooldownSeconds", Number(e.target.value))} />
        </label>
        <label className="policy-cell">
          Uncertainty ceiling
          <input type="number" min={0} max={1} step={0.05} value={draft.maxUncertainty} onChange={(e) => set("maxUncertainty", Number(e.target.value))} />
        </label>
        <label className="policy-cell">
          Session TTL (seconds)
          <input type="number" min={300} max={86400} step={60} value={draft.sessionTtlSeconds} onChange={(e) => set("sessionTtlSeconds", Number(e.target.value))} />
        </label>
        <div className="policy-cell">
          Max leverage
          <strong>1x</strong>
          <span className="fine">Locked. Session cannot change venue leverage.</span>
        </div>
        <div className="policy-cell">
          Venue
          <strong>hyperliquid</strong>
          <span className="fine">Perps only. Spot is not a PIT market type.</span>
        </div>
        <label className="policy-cell">
          <input type="checkbox" checked={draft.killSwitch} onChange={(e) => set("killSwitch", e.target.checked)} /> Kill switch
          <span className="fine">Halts new orders. Does not flatten.</span>
        </label>
      </div>
      <p className="fine" style={{ marginTop: 12 }}>
        Allowed assets
      </p>
      <div className="cta-row">
        {ASSETS.map((a) => (
          <button key={a} type="button" className={draft.allowedAssets.includes(a) ? "on" : "linkish"} onClick={() => toggleAsset(a)}>
            {a}
          </button>
        ))}
      </div>
      <article className="card" style={{ marginTop: 12 }}>
        <p className="label">If you pin this</p>
        {note.map((line) => (
          <p key={line}>{line}</p>
        ))}
      </article>
      <div className="cta-row">
        <button
          type="button"
          className="linkish"
          disabled={busy}
          onClick={() => {
            const lines = previewLines(current, draft);
            setLocal(lines);
            onPreview(draft);
          }}
        >
          Preview consequences
        </button>
        <button type="button" className="primary" disabled={busy} onClick={() => onPin(draft)}>
          {pinned ? "Pin updated policy" : "Pin policy on this computer"}
        </button>
      </div>
    </section>
  );
}

function previewLines(current: Partial<HostPolicy> | null | undefined, draft: HostPolicy): string[] {
  const from = emptyPolicy(current);
  const out: string[] = [];
  if (from.maxClipUsd !== draft.maxClipUsd) {
    out.push(`Max trade moves from $${from.maxClipUsd} to $${draft.maxClipUsd}. Host sizes every clip to this ceiling.`);
  }
  if (from.maxOpenPositions !== draft.maxOpenPositions) {
    out.push(`Open position ceiling moves from ${from.maxOpenPositions} to ${draft.maxOpenPositions}. Existing positions are not flattened.`);
  }
  if (from.dailyLossUsd !== draft.dailyLossUsd) {
    out.push(`Daily loss halt moves from $${from.dailyLossUsd} to $${draft.dailyLossUsd}.`);
  }
  if (from.maxSlippageBps !== draft.maxSlippageBps) {
    out.push(`Slippage band moves from ${from.maxSlippageBps} bps to ${draft.maxSlippageBps} bps.`);
  }
  if (from.cooldownSeconds !== draft.cooldownSeconds) {
    out.push(`Cooldown moves from ${from.cooldownSeconds}s to ${draft.cooldownSeconds}s.`);
  }
  if (from.sessionTtlSeconds !== draft.sessionTtlSeconds) {
    out.push(`Session TTL moves from ${from.sessionTtlSeconds}s to ${draft.sessionTtlSeconds}s.`);
  }
  if (from.killSwitch !== draft.killSwitch) {
    out.push(draft.killSwitch ? "Kill switch will be ON." : "Kill switch will be OFF.");
  }
  if (from.allowedAssets.slice().sort().join() !== draft.allowedAssets.slice().sort().join()) {
    out.push(`Allowed assets become ${draft.allowedAssets.join(", ")}.`);
  }
  if (!out.length) out.push("No field change. Pinning writes the same host law again.");
  out.push("Leverage stays 1x. Venue, market type, withdraw, and transfer cannot be changed. The model cannot pin this.");
  return out;
}