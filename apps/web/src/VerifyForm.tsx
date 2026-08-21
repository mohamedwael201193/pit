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
    <div className="mt-8 flex flex-col gap-2">
      <p className="text-[11px] tracking-[0.16em]">VERIFY A RECEIPT</p>
      <input
        aria-label="preview hash"
        className="rounded-xl bg-paper px-3.5 py-3 text-ink"
        placeholder="preview hash 0x"
        value={hash}
        onChange={(e) => onHash(e.target.value)}
      />
      <input
        aria-label="storage root"
        className="rounded-xl bg-paper px-3.5 py-3 text-ink"
        placeholder="storage root 0x"
        value={root}
        onChange={(e) => onRoot(e.target.value)}
      />
      <a
        className="mt-2 inline-block rounded-full border border-black/20 px-6 py-3 text-center font-semibold no-underline"
        href={explorer}
        target="_blank"
        rel="noreferrer"
      >
        Open {net} explorer
      </a>
    </div>
  );
}
