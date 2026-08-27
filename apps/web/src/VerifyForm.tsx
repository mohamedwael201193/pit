export function VerifyForm({
  hash,
  root,
  explorer,
  net,
  onHash,
  onRoot,
}: {
  hash: string;
  root: string;
  explorer: string;
  net: string;
  onHash: (v: string) => void;
  onRoot: (v: string) => void;
}) {
  return (
    <div className="flex flex-col gap-3">
      <label className="flex flex-col gap-2">
        <span className="text-[0.875rem] font-medium">Preview hash</span>
        <input
          aria-label="preview hash"
          className="rounded-xl border border-[rgb(240_231_212/0.3)] bg-transparent px-3.5 py-3 text-[var(--guide-cream)]"
          placeholder="0x"
          value={hash}
          onChange={(e) => onHash(e.target.value)}
        />
      </label>
      <label className="flex flex-col gap-2">
        <span className="text-[0.875rem] font-medium">Storage root</span>
        <input
          aria-label="storage root"
          className="rounded-xl border border-[rgb(240_231_212/0.3)] bg-transparent px-3.5 py-3 text-[var(--guide-cream)]"
          placeholder="0x"
          value={root}
          onChange={(e) => onRoot(e.target.value)}
        />
      </label>
      <a
        className="mt-2 inline-flex h-12 items-center justify-center rounded-full border border-[var(--guide-cream)] px-6 text-center font-medium text-[var(--guide-cream)] no-underline hover:bg-[var(--guide-cream)] hover:text-black"
        href={explorer}
        target="_blank"
        rel="noreferrer"
      >
        Open {net} explorer
      </a>
    </div>
  );
}
