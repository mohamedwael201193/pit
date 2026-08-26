// MOCK TEST HARNESS — public UI copy only. Never stub VerifyE2EE success.

export const RING = [
  "PRIVATE_BOOK",
  "SEALING",
  "TEE",
  "TEE_SIGNATURE",
  "ONCHAIN_SIGNER",
  "STORAGE",
  "RECEIPT",
  "CALIBRATION",
];

export function assertRing(labels: string[]) {
  if (labels.join(",") !== RING.join(",")) {
    throw new Error("ring");
  }
}
