import { mkdirSync } from "node:fs";
import { spawnSync } from "node:child_process";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const desktop = resolve(here, "..");
const repo = resolve(desktop, "../..");
const binaries = join(desktop, "src-tauri", "binaries");
mkdirSync(binaries, { recursive: true });

const win = process.platform === "win32";
const pitOut = join(binaries, win ? "pit-x86_64-pc-windows-msvc.exe" : "pit-x86_64-pc-windows-msvc");
const sealerOut = join(binaries, win ? "pit-sealer-x86_64-pc-windows-msvc.exe" : "pit-sealer-x86_64-pc-windows-msvc");

function goBuild(cwd, outfile, pkg) {
  const r = spawnSync("go", ["build", "-o", outfile, pkg], { cwd, stdio: "inherit", shell: win });
  if (r.status !== 0) {
    process.exit(r.status ?? 1);
  }
}

goBuild(join(repo, "pit"), pitOut, "./cmd/pit");
goBuild(join(repo, "sealer"), sealerOut, ".");
