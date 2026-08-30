export function deskHeadline(opts: {
  researchBusy?: boolean;
  doing: string;
  awaitingAuth?: boolean;
  ready: boolean;
  canOpen: boolean;
  execN: number;
  researchKind?: string;
  attentionTitle: string;
  policyTight?: boolean;
}): string {
  if (opts.researchBusy) return opts.doing;
  if (opts.awaitingAuth) return "Waiting for you";
  if (opts.canOpen) {
    return opts.execN === 1 ? "1 book can open" : `${opts.execN} books can open`;
  }
  if (opts.policyTight) return "Policy cap is too tight";
  if (opts.researchKind === "READY_STOOD_DOWN") return "Committee stood down. Checking next.";
  if (opts.ready) return "Watching. Nothing can open.";
  return opts.attentionTitle;
}
