import { useState, type FormEvent, type ReactNode } from "react";
import { Link } from "react-router-dom";
import {
  ARISTOTLE_EXPLORER,
  HISTORICAL_FILL,
  HISTORICAL_TEE_SIGNER,
  IDENTITY_8004,
  PIT_AGENT,
  RESEARCH_TXS,
  STORAGE_INDEXER,
  STORAGE_PROOFS,
  VERIFIED_FILL,
} from "./facts";
import { PageHead } from "../ui/PageHead";
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
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="What was verified, and how"
        lede={`This page will not say Verified unless this browser can show the check. Health process ${health?.version ?? "—"}. Sign is always false on this feed.`}
      />

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
          status="RECORDED ROOTS"
          proves="This desk filed research and order objects on Aristotle Flow. Roots are Flow event topics from those transactions. Open the indexer with a real root. This site will not badge a missing object."
          how={`Indexer ${STORAGE_INDEXER}. Confirm the root in the transaction log topic on ChainScan.`}
          extra={
            <ul className="mt-2 flex flex-col gap-2 text-[0.8125rem] leading-5">
              {STORAGE_PROOFS.map((p) => (
                <li key={p.root} className="break-all">
                  <span className="text-[rgb(240_231_212/0.55)]">{p.label}</span>
                  {p.tx ? (
                    <>
                      {" · "}
                      <a className="text-[#d82f2f]" href={`${ARISTOTLE_EXPLORER}/tx/${p.tx}`} target="_blank" rel="noreferrer">
                        tx {p.tx}
                      </a>
                    </>
                  ) : null}
                  <span className="block text-[rgb(240_231_212/0.75)]">root {p.root}</span>
                </li>
              ))}
            </ul>
          }
        />
        <ProofRow
          name="0G Chain"
          status="BROWSER READ"
          proves="A transaction exists on Aristotle if the RPC returns it. Paste a hash below. This browser calls evmrpc.0g.ai."
          how="This page does not run VerifyE2EE."
        />
        <ProofRow
          name="This-job 0G"
          status="RECORDED"
          proves={`Job ${VERIFIED_FILL.job}. Research and order evidence for that job. Not a historical fallback. This page did not re-run VerifyE2EE.`}
          how="Open the Aristotle explorer. Confirm the hashes match the job."
          extra={
            <p className="mt-2 flex flex-col gap-1 text-[0.8125rem] leading-5">
              <a className="text-[#d82f2f] break-all" href={`${ARISTOTLE_EXPLORER}/tx/${VERIFIED_FILL.researchTx}`} target="_blank" rel="noreferrer">
                Research {VERIFIED_FILL.researchTx}
              </a>
              <span className="break-all text-[rgb(240_231_212/0.7)]">root {VERIFIED_FILL.researchRoot}</span>
              <a className="text-[#d82f2f] break-all" href={`${ARISTOTLE_EXPLORER}/tx/${VERIFIED_FILL.orderTx}`} target="_blank" rel="noreferrer">
                Order {VERIFIED_FILL.orderTx}
              </a>
              <span className="break-all text-[rgb(240_231_212/0.7)]">root {VERIFIED_FILL.orderRoot}</span>
            </p>
          }
        />
        <ProofRow
          name="Recorded research txs"
          status="RECORDED"
          proves="Earlier sealed research jobs on this desk, each with its own Aristotle filing. Not a historical fallback for the HYPE fill."
          how="Open ChainScan. Each hash is a separate job receipt."
          extra={
            <ul className="mt-2 flex flex-col gap-1 text-[0.8125rem] leading-5">
              {RESEARCH_TXS.map((p) => (
                <li key={p.tx} className="break-all">
                  <a className="text-[#d82f2f]" href={`${ARISTOTLE_EXPLORER}/tx/${p.tx}`} target="_blank" rel="noreferrer">
                    {p.label} {p.tx}
                  </a>
                </li>
              ))}
            </ul>
          }
        />
        <ProofRow
          name="Hyperliquid"
          status="RECORDED OID"
          proves={`Latest recorded fill on this desk: OID ${VERIFIED_FILL.oid} / ${VERIFIED_FILL.market} ${VERIFIED_FILL.sz} @ ${VERIFIED_FILL.px}, job ${VERIFIED_FILL.job}. Older ETH OID ${HISTORICAL_FILL.oid} is historical and was not flattened.`}
          how="This site does not fetch another account’s fills and does not run a live mission stream. Verify the OID on the account that received it."
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
          proves="AUTHORIZE or ARM SLEEP MISSION on desktop."
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

function ProofRow({
  name,
  status,
  proves,
  how,
  extra,
}: {
  name: string;
  status: string;
  proves: string;
  how: string;
  extra?: ReactNode;
}) {
  return (
    <article className="grid gap-2 py-5 md:grid-cols-[9rem_8rem_1fr] md:gap-6">
      <h2 className="font-semibold">{name}</h2>
      <p className="text-[0.6875rem] tracking-[0.14em] text-[#d82f2f]">{status}</p>
      <div>
        <p className="text-[0.9375rem] leading-6">{proves}</p>
        <p className="mt-1 text-[0.8125rem] leading-5 text-[rgb(240_231_212/0.5)]">HOW: {how}</p>
        {extra}
      </div>
    </article>
  );
}
