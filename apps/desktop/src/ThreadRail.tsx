import { useMemo, useState } from "react";
import type { ChatThread } from "./companion";

export function ThreadRail({
  threads,
  active,
  onSelect,
  onNew,
  onRename,
  onDelete,
}: {
  threads: ChatThread[];
  active: string;
  onSelect: (id: string) => void;
  onNew: () => void;
  onRename: (id: string, title: string) => void;
  onDelete: (id: string) => void;
}) {
  const [q, setQ] = useState("");
  const rows = useMemo(() => {
    const n = q.trim().toLowerCase();
    if (!n) return threads;
    return threads.filter((t) => `${t.title} ${t.preview}`.toLowerCase().includes(n));
  }, [threads, q]);

  return (
    <aside className="threads" aria-label="Desk threads">
      <div className="threads-head">
        <p className="label">Threads</p>
        <button type="button" className="primary" onClick={onNew}>
          New
        </button>
      </div>
      <input
        className="threads-search"
        aria-label="Search threads"
        placeholder="Search"
        value={q}
        onChange={(e) => setQ(e.target.value)}
      />
      <ul className="thread-list">
        {rows.length === 0 ? (
          <li>
            <p className="fine" style={{ padding: "8px 12px" }}>
              No threads yet. New starts an empty command log on this workspace.
            </p>
          </li>
        ) : (
          rows.map((t) => (
            <li key={t.id} className={t.id === active ? "on" : ""}>
              <button type="button" className="row" onClick={() => onSelect(t.id)}>
                <span className="title">{t.title || "Desk"}</span>
                <span className="preview">{t.preview || "Empty thread"}</span>
                <span className="meta">
                  {t.updated ? new Date(t.updated).toLocaleString() : ""}
                  {threads.length > 1 ? (
                    <>
                      {" · "}
                      <span
                        role="button"
                        tabIndex={0}
                        onClick={(e) => {
                          e.stopPropagation();
                          const next = window.prompt("Rename thread", t.title || "");
                          if (next) onRename(t.id, next);
                        }}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            const next = window.prompt("Rename thread", t.title || "");
                            if (next) onRename(t.id, next);
                          }
                        }}
                      >
                        Rename
                      </span>
                      {t.id !== "desk" ? (
                        <>
                          {" · "}
                          <span
                            role="button"
                            tabIndex={0}
                            onClick={(e) => {
                              e.stopPropagation();
                              onDelete(t.id);
                            }}
                            onKeyDown={(e) => {
                              if (e.key === "Enter") onDelete(t.id);
                            }}
                          >
                            Delete
                          </span>
                        </>
                      ) : null}
                    </>
                  ) : null}
                </span>
              </button>
            </li>
          ))
        )}
      </ul>
    </aside>
  );
}
