import { committeeVerified, oidBelongsToPreview, researchCardTitle } from "../src/honesty";

export function assertHonesty() {
  const one = [{ role: "researcher", verify_e2ee: "OK" }];
  if (committeeVerified(one)) throw new Error("one role verified");
  if (researchCardTitle("COMMITTEE_INCOMPLETE") === "RESEARCH VERIFIED") throw new Error("incomplete title");
  const three = [
    { role: "researcher", verify_e2ee: "OK" },
    { role: "challenger", verify_e2ee: "OK" },
    { role: "risk", verify_e2ee: "OK" },
  ];
  if (!committeeVerified(three)) throw new Error("three roles");
  if (oidBelongsToPreview("0xold", "0xnew")) throw new Error("stale oid");
  if (!oidBelongsToPreview("0xabc", "0xabc")) throw new Error("matching hash");
}
