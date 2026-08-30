"""Deploy web (Vercel) and health (Render). Never print secrets."""
from __future__ import annotations

import json
import subprocess
import urllib.error
import urllib.request
from pathlib import Path

ROOT = Path(r"d:\route\0g\PIT")
WEB = ROOT / "apps" / "web"


def load_env() -> dict[str, str]:
    out: dict[str, str] = {}
    for line in (ROOT / ".env").read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, v = line.split("=", 1)
        out[k.strip()] = v.strip().strip('"').strip("'")
    return out


def redact(text: str, secrets: list[str]) -> str:
    for s in secrets:
        if s:
            text = text.replace(s, "[redacted]")
    return text


def api(method: str, url: str, token: str, data: dict | None = None):
    body = None if data is None else json.dumps(data).encode()
    r = urllib.request.Request(
        url,
        data=body,
        headers={
            "Authorization": f"Bearer {token}",
            "Accept": "application/json",
            "Content-Type": "application/json",
        },
        method=method,
    )
    with urllib.request.urlopen(r, timeout=45) as resp:
        raw = resp.read().decode()
        return json.loads(raw) if raw else {}


def main() -> None:
    env = load_env()
    vercel = env.get("vercel_token") or env.get("VERCEL_TOKEN") or ""
    render = env.get("render_api_key") or env.get("RENDER_API_KEY") or ""
    privy = env.get("VITE_PRIVY_APP_ID") or "cmtafcijw02av0cl1ay81om7m"
    secrets = [vercel, render, env.get("github_token", ""), env.get("GITHUB_TOKEN", "")]

    if not vercel:
        raise SystemExit("missing vercel token")

    pid = "prj_1fUFZk2BM9x9OyVoXKCvChNdqlAF"
    print("vercel: clear rootDirectory for CLI upload from apps/web")
    try:
        api("PATCH", f"https://api.vercel.com/v9/projects/{pid}", vercel, {"rootDirectory": None})
    except urllib.error.HTTPError:
        api("PATCH", f"https://api.vercel.com/v9/projects/{pid}", vercel, {"rootDirectory": ""})

    cmd = [
        "npx",
        "--yes",
        "vercel@latest",
        "--prod",
        "--yes",
        "--token",
        vercel,
        "--env",
        f"VITE_PRIVY_APP_ID={privy}",
        "--build-env",
        f"VITE_PRIVY_APP_ID={privy}",
        "--env",
        "VITE_HEALTH_URL=https://pit-health.onrender.com",
        "--build-env",
        "VITE_HEALTH_URL=https://pit-health.onrender.com",
    ]
    print("vercel: deploying apps/web")
    try:
        r = subprocess.run(cmd, cwd=WEB, capture_output=True, text=True, shell=True)
        print(redact((r.stdout or "") + (r.stderr or ""), secrets)[-4000:])
        print("vercel exit", r.returncode)
    finally:
        print("vercel: restore git rootDirectory apps/web")
        api("PATCH", f"https://api.vercel.com/v9/projects/{pid}", vercel, {"rootDirectory": "apps/web"})

    if render:
        req = urllib.request.Request(
            "https://api.render.com/v1/services?limit=50",
            headers={"Authorization": f"Bearer {render}", "Accept": "application/json"},
        )
        try:
            with urllib.request.urlopen(req, timeout=30) as resp:
                body = json.loads(resp.read().decode())
            names = []
            for item in body:
                svc = item.get("service") or item
                names.append({"id": svc.get("id"), "name": svc.get("name"), "type": svc.get("type")})
            print("render services", json.dumps(names))
            pit = next((n for n in names if n.get("name") == "pit-health"), None)
            if pit and pit.get("id"):
                deploy_req = urllib.request.Request(
                    f"https://api.render.com/v1/services/{pit['id']}/deploys",
                    data=b"{}",
                    headers={
                        "Authorization": f"Bearer {render}",
                        "Accept": "application/json",
                        "Content-Type": "application/json",
                    },
                    method="POST",
                )
                with urllib.request.urlopen(deploy_req, timeout=30) as resp:
                    d = json.loads(resp.read().decode())
                print("render deploy", d.get("id") or d.get("status") or "ok")
            else:
                print("render: pit-health not found; create from render.yaml in dashboard")
        except urllib.error.HTTPError as e:
            print("render http", e.code, redact(e.read().decode()[:400], secrets))
    else:
        print("render: no api key")


if __name__ == "__main__":
    main()
