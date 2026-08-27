export const RING = [
  "PRIVATE_BOOK",
  "SEALING",
  "TEE",
  "TEE_SIGNATURE",
  "ONCHAIN_SIGNER",
  "STORAGE",
  "RECEIPT",
  "CALIBRATION",
] as const;

export function Ring() {
  return (
    <ol className="text-[26px] leading-[1.35] tracking-[-0.03em]">
      {RING.map((s) => (
        <li key={s}>{s}</li>
      ))}
    </ol>
  );
}
