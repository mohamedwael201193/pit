import { useEffect, useState } from "react";
import { PageHead } from "../ui/PageHead";
import { RELEASES, REPO } from "./facts";

type Release = {
  tag?: string;
  name?: string;
  html?: string;
  sha?: string | null;
  unsigned: boolean;
};

export function DownloadPage() {
  const [rel, setRel] = useState<Release | null>(null);
  const [fail, setFail] = useState<string | null>(null);

  useEffect(() => {
    const ac = new AbortController();
    fetch("https://api.github.com/repos/mohamedwael201193/pit/releases/latest", { signal: ac.signal })
      .then((r) => {
        if (!r.ok) throw new Error("github");
        return r.json() as Promise<{
          tag_name?: string;
          name?: string;
          html_url?: string;
          assets?: { name: string; browser_download_url: string }[];
        }>;
      })
      .then(async (body) => {
        const sums = (body.assets ?? []).find((a) => a.name.toLowerCase().includes("sha256"));
        let sha: string | null = null;
        if (sums) {
          const text = await fetch(sums.browser_download_url, { signal: ac.signal }).then((r) => r.text());
          const line = text.split(/\r?\n/).find((l) => /setup\.exe/i.test(l) || /PIT_.*x64/i.test(l));
          sha = line ? line.trim().split(/\s+/)[0] : text.slice(0, 64);
        }
        setRel({
          tag: body.tag_name,
          name: body.name,
          html: body.html_url,
          sha,
          unsigned: true,
        });
        setFail(null);
      })
      .catch(() => {
        if (!ac.signal.aborted) {
          setRel(null);
          setFail("GitHub Releases could not be read in this browser. Open the latest release and verify SHA256SUMS there.");
        }
      });
    return () => ac.abort();
  }, []);

  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Windows x64. Verify the bytes."
        lede="PIT Desktop is the private brain. This website will not claim a Windows signature that does not exist. macOS and Linux are not claimed until they are packaged and tested."
      />

      <dl className="intel-metrics mt-8">
        <div className="intel-metric">
          <dt>Platform</dt>
          <dd>Windows x64</dd>
        </div>
        <div className="intel-metric">
          <dt>Latest</dt>
          <dd>{rel?.tag ?? "see GitHub"}</dd>
          <p>{rel?.name ?? ""}</p>
        </div>
        <div className="intel-metric">
          <dt>SHA256</dt>
          <dd className="!text-[0.85rem] break-all">{rel?.sha ?? "open SHA256SUMS"}</dd>
          <p>Verify against the release asset. Do not trust a screenshot.</p>
        </div>
        <div className="intel-metric">
          <dt>Signature</dt>
          <dd>unsigned</dd>
          <p>Authenticode is not claimed.</p>
        </div>
      </dl>

      {fail ? <p className="mt-4 text-[0.875rem] text-[#ff8a8a]">{fail}</p> : null}

      <ol className="intel-steps mt-10">
        <li>Download the NSIS installer from GitHub Releases.</li>
        <li>Hash the file. Compare to SHA256SUMS on the same release.</li>
        <li>Install. Launch PIT Desktop. Pairing comes after you have seen radar and proof.</li>
      </ol>

      <div className="mt-10 flex flex-wrap gap-2.5">
        <a className="intel-cta" href={rel?.html ?? RELEASES} target="_blank" rel="noreferrer">
          GitHub release
        </a>
        <a className="intel-secondary" href={RELEASES} target="_blank" rel="noreferrer">
          Verify download
        </a>
        <a className="intel-ghost" href={REPO} target="_blank" rel="noreferrer">
          Source
        </a>
      </div>
    </div>
  );
}
