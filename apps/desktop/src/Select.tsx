import { useEffect, useRef, useState } from "react";

export function Select({
  value,
  onChange,
  options,
  disabled,
  label,
  id,
}: {
  value: string;
  onChange: (v: string) => void;
  options: { value: string; label: string }[];
  disabled?: boolean;
  label?: string;
  id?: string;
}) {
  const [open, setOpen] = useState(false);
  const root = useRef<HTMLDivElement>(null);
  const current = options.find((o) => o.value === value);
  useEffect(() => {
    function onDoc(e: MouseEvent) {
      if (!root.current?.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", onDoc);
    return () => document.removeEventListener("mousedown", onDoc);
  }, []);
  return (
    <div className="pit-select" ref={root}>
      {label ? (
        <label className="pit-select-label" htmlFor={id}>
          {label}
        </label>
      ) : null}
      <button
        id={id}
        type="button"
        className="pit-select-btn"
        disabled={disabled}
        aria-haspopup="listbox"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
      >
        <span>{current?.label || value}</span>
        <span aria-hidden>▾</span>
      </button>
      {open ? (
        <ul className="pit-select-menu" role="listbox">
          {options.map((o) => (
            <li key={o.value}>
              <button
                type="button"
                role="option"
                aria-selected={o.value === value}
                className={o.value === value ? "on" : ""}
                onClick={() => {
                  onChange(o.value);
                  setOpen(false);
                }}
              >
                {o.label}
              </button>
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  );
}
