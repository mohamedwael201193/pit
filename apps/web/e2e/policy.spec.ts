// MOCK TEST HARNESS — public UI copy only. Never stub AUTHORIZE or a live order.

export function assertPolicyCards(titles: string[]) {
  const need = ["Max trade", "Cooldown", "Max uncertainty", "Min liquidity"];
  for (const t of need) {
    if (!titles.includes(t)) {
      throw new Error("policy card " + t);
    }
  }
}
