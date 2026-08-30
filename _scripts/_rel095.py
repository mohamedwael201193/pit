"""Create or update GitHub release v0.9.5. Never print the token."""
from __future__ import annotations

import os
import subprocess
from pathlib import Path

ROOT = Path(r"d:\route\0g\PIT")
NSIS = ROOT / "apps/desktop/src-tauri/target/release/bundle/nsis/PIT_0.9.5_x64-setup.exe"
SUMS = ROOT / "apps/desktop/src-tauri/target/release/bundle/nsis/SHA256SUMS.txt"
NOTES = """Windows installer (NSIS). Not Authenticode-signed.

PIT 0.9.5 rebuilds the desktop Agent as one conversation. Research, progress, verdict, TRADE NOW, and OID live inside a single PIT turn. No stacked cockpit above a transcript. TRADE NOW still calls the existing desktop AUTHORIZE path. The model cannot AUTHORIZE.

Verify checksums. Pair at https://pit0g.vercel.app/pair after launch.
"""


def token() -> str:
    for line in (ROOT / ".env").read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, v = line.split("=", 1)
        if k.strip() in ("github_token", "GITHUB_TOKEN"):
            return v.strip().strip('"').strip("'")
    raise SystemExit("missing github token")


def main() -> None:
    tok = token()
    env = os.environ.copy()
    env["GH_TOKEN"] = tok
    if not NSIS.exists():
        raise SystemExit(f"missing {NSIS}")
    view = subprocess.run(
        ["gh", "release", "view", "v0.9.5", "--json", "tagName"],
        cwd=ROOT,
        capture_output=True,
        text=True,
        env=env,
    )
    assets = [str(NSIS), str(SUMS)] if SUMS.exists() else [str(NSIS)]
    if view.returncode != 0:
        r = subprocess.run(
            ["gh", "release", "create", "v0.9.5", *assets, "--title", "PIT 0.9.5", "--notes", NOTES],
            cwd=ROOT,
            capture_output=True,
            text=True,
            env=env,
        )
    else:
        r = subprocess.run(
            ["gh", "release", "upload", "v0.9.5", *assets, "--clobber"],
            cwd=ROOT,
            capture_output=True,
            text=True,
            env=env,
        )
    out = (r.stdout or "") + (r.stderr or "")
    print(out.replace(tok, "[redacted]"))
    print("exit", r.returncode)


if __name__ == "__main__":
    main()
