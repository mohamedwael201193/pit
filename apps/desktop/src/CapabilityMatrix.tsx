export type CapRow = {
  id: string;
  label: string;
  mainnet: "live" | "partial" | "off";
  testnet: "live" | "partial" | "off";
  note: string;
};

export const CAPABILITY: CapRow[] = [
  {
    id: "direct",
    label: "0G Direct TeeML",
    mainnet: "live",
    testnet: "off",
    note: "Aristotle glm-5.3. Galileo sealed ask stays off until VerifyE2EE is proven.",
  },
  {
    id: "tee",
    label: "TEE verification",
    mainnet: "live",
    testnet: "off",
    note: "Recovered signer must equal on-chain teeSigner. A 200 is not success.",
  },
  {
    id: "storage",
    label: "Storage proof",
    mainnet: "partial",
    testnet: "partial",
    note: "Official Go client upload/download --proof. Workspace key stays on this machine.",
  },
  {
    id: "hl",
    label: "Hyperliquid",
    mainnet: "partial",
    testnet: "partial",
    note: "Order and cancel after Hyperliquid lists the local PIT agent. Withdraw is impossible.",
  },
  {
    id: "agentic",
    label: "Agentic ID",
    mainnet: "partial",
    testnet: "partial",
    note: "TokenId 1 exists. Mint, ownerOf, authorizeUsage, and revoke are live on Aristotle. iTransfer and iClone UNAVAILABLE (AttestorNotOnAristotle). Trading does not wait on mint.",
  },
  {
    id: "erc8004",
    label: "ERC-8004",
    mainnet: "partial",
    testnet: "partial",
    note: "Identity ownerOf is live. Reputation giveFeedback is an on-chain write from a reporter that is not the owner. Self-feedback is rejected.",
  },
];

function mark(v: CapRow["mainnet"]) {
  if (v === "live") return "LIVE";
  if (v === "partial") return "PARTIAL";
  return "OFF";
}

export function CapabilityMatrix({ net }: { net: "mainnet" | "testnet" }) {
  return (
    <article className="card">
      <p className="label">CAPABILITY</p>
      <p className="fine">Honest per network. Green is never invented.</p>
      <table className="watch-table">
        <thead>
          <tr>
            <th>Feature</th>
            <th>{net === "testnet" ? "Testnet" : "Mainnet"}</th>
            <th>Why</th>
          </tr>
        </thead>
        <tbody>
          {CAPABILITY.map((row) => (
            <tr key={row.id}>
              <td>{row.label}</td>
              <td>{mark(net === "testnet" ? row.testnet : row.mainnet)}</td>
              <td>{row.note}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </article>
  );
}
