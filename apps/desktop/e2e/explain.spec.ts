import { explainStop } from "../src/explain";

export function assertNamedResearchWhy() {
  const poll = explainStop("POLL_FAILED");
  if (!poll || poll.title === "Research stopped" || /crash/i.test(poll.body)) {
    throw new Error("poll is not a research failure");
  }
  const stood = explainStop("READY_STOOD_DOWN");
  if (!stood || !/no trade survived|stood down/i.test(stood.title) || !/not a crash|checking the next/i.test(stood.body)) {
    throw new Error("stand-down must not be a crash");
  }
  const credit = explainStop("DIRECT_CREDIT_INSUFFICIENT");
  if (!credit || !/compute/i.test(credit.body)) {
    throw new Error("compute why");
  }
  const tee = explainStop("TEE_VERIFY_FAIL");
  if (!tee || !/TEE/i.test(tee.title)) {
    throw new Error("tee why");
  }
}
