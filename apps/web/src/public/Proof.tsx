import { useState, type FormEvent } from "react";
import { Link } from "react-router-dom";
import {
  ARISTOTLE_EXPLORER,
  HISTORICAL_FILL,
  HISTORICAL_TEE_SIGNER,
  IDENTITY_8004,
  PIT_AGENT,
  STORAGE_INDEXER,
} from "./facts";
import { shortAddr } from "./format";
import { readAristotleTx } from "./chain";
import { useWatch } from "./Watch";

export function ProofPage() {
  const { health } = useWatch();
  const [hash, setHash] = useState("");
  const [txMsg, setTxMsg] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const onTx = async (e: FormEvent) => {
    e.preventDefault();
    setBusy(true);
    const got = await readAristotleTx(hash);
    setBusy(false);
    if (!got.ok) {
      setTxMsg(got.reason);
      return;
    }
    setTxMsg(
      `${got.status === "success" ? "MATCH" : got.status.toUpperCase()} on Aristotle. From ${shortAddr(got.from)} to ${got.to ? shortAddr(got.to) : "contract-create"} in block ${got.blockNumber}. Status ${got.status}. This browser read evmrpc.0g.ai getTransaction + getTransactionReceipt. It did not run VerifyE2EE.`,
    );
  };

  return (
    <div>
      <p className="intel-kicker">Proof center</p>
      <h1 className="intel-title mt-2">What was verified, and how</h1>
      <p className="intel-lede">
        This page will not say Verified unless this browser can show the check. Health process {health?.version ?? "—"}.
        Sign is always false on this feed.
      </p>

      <div className="mt-10 divide-y divide-[rgb(240_231_212/0.12)] border-y border-[rgb(240_231_212/0.12)]">
        <ProofRow
          name="PIT agent"
          status="IDENTITY"
          proves={`${PIT_AGENT.name} is the Hyperliquid API wallet name for ${shortAddr(PIT_AGENT.address)}. iTransfer is not live on mainnet.`}
          how="Compare the printed agent on Hyperliquid API. This site does not hold the session key."
        />
        <ProofRow
          name="0G Compute"
          status="PRIVATE PATH"
          proves="Direct TeeML is the only private research path. Router is impossible for the private book."
          how="A live job receipt is required to recover a signer. None is published on this page."
        />
        <ProofRow
          name="TEE"
          status="NO LIVE RECEIPT"
          proves="Recover signer from Direct evidence, compare to the registered signer."
          how={`Historical recovered signer ${shortAddr(HISTORICAL_TEE_SIGNER)} from a prior sealed job is HISTORICAL. Expected listed Direct teeSigner is ${shortAddr(HISTORICAL_TEE_SIGNER)}. This page does not run VerifyE2EE.`}
        />
        <ProofRow
          name="0G Storage"
          status="NO PUBLIC ROOT"
          proves="Resolve the proof object and recompute the hash where the client supports it."
          how={`Open the indexer ${STORAGE_INDEXER} with a real root. This site will not badge a missing object.`}
        />
        <ProofRow
          name="0G Chain"
          status="BROWSER READ"
          proves="A transaction exists on Aristotle if the RPC returns it."
          how="Paste a hash. This browser calls evmrpc.0g.ai."
        />
        <ProofRow
          name="Hyperliquid"
          status="HISTORICAL OID"
          proves={`OID ${HISTORICAL_FILL.oid} / ${HISTORICAL_FILL.market} ${HISTORICAL_FILL.sz} is a recorded fill.`}
          how="This site does not fetch another account’s fills. Verify the OID on the account that received it."
        />
        <ProofRow
          name="Policy"
          status="DESKTOP"
          proves="Pinned clip, assets, and kill. The model cannot raise them."
          how="Open PIT Desktop. This website cannot pin policy."
        />
        <ProofRow
          name="Execution"
          status="DESKTOP"
          proves="AUTHORIZE or host-gated Guarded Autonomy."
          how="No web control can place an order."
        />
      </div>

      <form className="mt-10 max-w-[36rem]" onSubmit={(e) => void onTx(e)}>
        <label className="block">
          <span className="intel-kicker">Read an Aristotle transaction</span>
          <input
            className="intel-input mt-2"
            value={hash}
            onChange={(e) => setHash(e.target.value)}
            placeholder="0x…"
            aria-label="0G Chain transaction hash"
          />
        </label>
        <button type="submit" className="intel-cta mt-4" disabled={busy}>
          {busy ? "Reading RPC…" : "Read transaction"}
        </button>
        {txMsg ? <p className="mt-3 text-[0.875rem] leading-6 text-[rgb(240_231_212/0.75)]">{txMsg}</p> : null}
        <a className="intel-ghost mt-4 inline-flex" href={ARISTOTLE_EXPLORER} target="_blank" rel="noreferrer">
          Open ChainScan
        </a>
      </form>

      <p className="mt-8 text-[0.8125rem] text-[rgb(240_231_212/0.45)]">
        ERC-8004 identity {shortAddr(IDENTITY_8004)}. See /agent for a live ownerOf read. No ranking is invented.
      </p>
      <Link to="/agent" className="intel-ghost mt-4 inline-flex">
        Open agent passport
      </Link>
    </div>
  );
}

function ProofRow({ name, status, proves, how }: { name: string; status: string; proves: string; how: string }) {
  return (
    <article className="grid gap-2 py-5 md:grid-cols-[9rem_8rem_1fr] md:gap-6">
      <h2 className="font-semibold">{name}</h2>
      <p className="text-[0.6875rem] tracking-[0.14em] text-[#d82f2f]">{status}</p>
      <div>
        <p className="text-[0.9375rem] leading-6">{proves}</p>
        <p className="mt-1 text-[0.8125rem] leading-5 text-[rgb(240_231_212/0.5)]">HOW: {how}</p>
      </div>
    </article>
  );
}
