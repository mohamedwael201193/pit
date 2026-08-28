import { committeeVerified, oidBelongsToPreview, researchCardTitle } from "../src/honesty";
import { researchWhyCopy } from "../src/researchWhy";

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
  if (researchCardTitle("READY_STOOD_DOWN") !== "COMMITTEE STOOD DOWN") throw new Error("stood down title");
  if (researchCardTitle("COMMITTEE_INCOMPLETE") === "RESEARCH VERIFIED") throw new Error("incomplete title");
  if (oidBelongsToPreview("0xold", "0xnew")) throw new Error("stale oid");
  if (!oidBelongsToPreview("0xabc", "0xabc")) throw new Error("matching hash");
  const why = researchWhyCopy({ coin: "ETH", kind: "READY_STOOD_DOWN", roles: three, deny: "no_side" });
  if (why.length !== 8) throw new Error("why rows");
  if (!why[6].a.toLowerCase().includes("stand-down")) throw new Error("stand-down why");
}
