import { LINKS } from "./links";

export type AutoPrefs = {
  watch?: boolean;
  auto_research?: boolean;
  notify?: boolean;
  cadence_minutes?: number;
  trigger?: string;
  markets?: string[];
};

export function AutomationPane({
  prefs,
  busy,
  onSave,
}: {
  prefs: AutoPrefs;
  busy?: boolean;
  onSave: (p: AutoPrefs) => void;
}) {
  const cadence = prefs.cadence_minutes || 15;
  return (
    <article className="card">
      <p className="label">Automation</p>
      <p>Watch, discover, research, notify, and prepare. Automation cannot AUTHORIZE live money movement.</p>
      <label className="fine">
        <input
          type="checkbox"
          checked={Boolean(prefs.watch)}
          disabled={busy}
          onChange={(e) => onSave({ ...prefs, watch: e.target.checked })}
        />{" "}
        Watch live policy markets
      </label>
      <label className="fine">
        <input
          type="checkbox"
          checked={Boolean(prefs.notify)}
          disabled={busy}
          onChange={(e) => onSave({ ...prefs, notify: e.target.checked })}
        />{" "}
        Notify when something is interesting
      </label>
      <label className="fine">
        <input
          type="checkbox"
          checked={Boolean(prefs.auto_research)}
          disabled={busy}
          onChange={(e) => onSave({ ...prefs, auto_research: e.target.checked })}
        />{" "}
        Start sealed research automatically (still stops at human approval)
      </label>
      <label className="fine">
        Cadence
        <select
          value={String(cadence)}
          disabled={busy}
          onChange={(e) => onSave({ ...prefs, cadence_minutes: Number(e.target.value) })}
        >
          <option value="5">Every 5 minutes</option>
          <option value="15">Every 15 minutes</option>
          <option value="60">Every hour</option>
        </select>
      </label>
      <label className="fine">
        Trigger
        <select
          value={prefs.trigger || "policy_pass"}
          disabled={busy}
          onChange={(e) => onSave({ ...prefs, trigger: e.target.value })}
        >
          <option value="policy_pass">Any policy-eligible market</option>
          <option value="mark_below_oracle">Mark below oracle</option>
          <option value="funding">Nonzero funding</option>
        </select>
      </label>
      <p className="fine">
        Markets follow pinned policy unless you restrict them later. Compute money is not trading capital.{" "}
        <a href={LINKS.pcAdvanced} target="_blank" rel="noreferrer">
          Open 0G Private Compute
        </a>
      </p>
    </article>
  );
}
