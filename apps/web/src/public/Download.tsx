import { useEffect, useState } from "react";
import { PageHead } from "../ui/PageHead";
import { fetchRelease } from "./api";
import { REPO, windowsChecksumsUrl, windowsInstallerUrl } from "./facts";

type Release = {
  tag?: string;
  name?: string;
  sha?: string | null;
  filename?: string | null;
  unsigned: boolean;
};

export function DownloadPage() {
  const [rel, setRel] = useState<Release | null>(null);
  const [fail, setFail] = useState<string | null>(null);
  const installer = windowsInstallerUrl();
  const checksums = windowsChecksumsUrl();

  useEffect(() => {
    const ac = new AbortController();
    fetchRelease(ac.signal)
      .then((body) => {
        if (!body.tag && !body.sha) throw new Error("release");
        setRel({
          tag: body.tag,
          name: body.name,
          sha: body.sha || null,
          filename: body.filename || "PIT_x64-setup.exe",
          unsigned: true,
        });
        setFail(null);
      })
      .catch(() => {
        if (!ac.signal.aborted) {
          setRel(null);
          setFail("Release metadata could not be read. The download button still fetches the latest installer from PIT health.");
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
          <dd>{rel?.tag ?? "latest"}</dd>
          <p>{rel?.name ?? rel?.filename ?? "PIT_x64-setup.exe"}</p>
        </div>
        <div className="intel-metric">
          <dt>SHA256</dt>
          <dd className="!text-[0.85rem] break-all">{rel?.sha ?? "open SHA256SUMS"}</dd>
          <p>Hash the file. Compare with SHA256SUMS from the same release.</p>
        </div>
        <div className="intel-metric">
          <dt>Signature</dt>
          <dd>unsigned</dd>
          <p>Authenticode is not claimed.</p>
        </div>
      </dl>

      {fail ? <p className="mt-4 text-[0.875rem] text-[#ff8a8a]">{fail}</p> : null}

      <ol className="intel-steps mt-10">
        <li>Click Download Windows installer. That is a file download (HTTP 302 to the release asset), not a GitHub Releases HTML page.</li>
        <li>Hash the file. Compare to SHA256SUMS.</li>
        <li>Install. Launch PIT Desktop. Pairing comes after you have seen radar and proof.</li>
      </ol>

      <div className="mt-10 flex flex-wrap gap-2.5">
        <a className="intel-cta" href={installer} download>
          Download Windows installer
        </a>
        <a className="intel-secondary" href={checksums}>
          SHA256SUMS
        </a>
        <a className="intel-ghost" href={REPO} target="_blank" rel="noreferrer">
          Source
        </a>
      </div>
    </div>
  );
}
