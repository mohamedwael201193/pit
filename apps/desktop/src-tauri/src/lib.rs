use std::fs::{self, OpenOptions};
use std::io::{Read, Write};
use std::net::{SocketAddr, TcpStream};
use std::path::{Path, PathBuf};
use std::process::{Command, Stdio};
use std::time::Duration;

#[tauri::command]
fn companion_url() -> String {
    "http://127.0.0.1:17373".into()
}

#[tauri::command]
fn export_session() -> Result<String, String> {
    Err("session_export_denied".into())
}

#[tauri::command]
fn local_status() -> Result<serde_json::Value, String> {
    loopback_json("/local/status")
}

#[tauri::command]
fn local_code() -> Result<serde_json::Value, String> {
    loopback_json("/local/code")
}

#[tauri::command]
fn local_doctor() -> Result<serde_json::Value, String> {
    loopback_json("/local/doctor")
}

#[tauri::command]
fn local_init(wallet: String, network: String) -> Result<serde_json::Value, String> {
    let body = serde_json::json!({ "wallet": wallet, "network": network });
    loopback_json_post("/local/init", &body)
}

#[tauri::command]
fn local_session() -> Result<serde_json::Value, String> {
    loopback_json_post("/local/session", &serde_json::json!({}))
}

#[tauri::command]
fn local_policy() -> Result<serde_json::Value, String> {
    loopback_json_post("/local/policy", &serde_json::json!({}))
}

#[tauri::command]
fn local_revoke_session() -> Result<serde_json::Value, String> {
    loopback_json_post("/local/revoke-session", &serde_json::json!({}))
}

#[tauri::command]
fn local_direct_intent() -> Result<serde_json::Value, String> {
    let raw = loopback_exchange_timeout("GET", "/local/direct-intent", None, 20)?;
    serde_json::from_str(&raw).map_err(|e| e.to_string())
}

#[tauri::command]
fn local_direct_status() -> Result<serde_json::Value, String> {
    loopback_json("/local/direct-status")
}

#[tauri::command]
fn local_research_start(coin: String) -> Result<serde_json::Value, String> {
    let body = serde_json::json!({ "coin": coin });
    loopback_json_post_timeout("/local/research/start", &body, 8)
}

#[tauri::command]
fn local_research_status() -> Result<serde_json::Value, String> {
    loopback_json("/local/research/status")
}

#[tauri::command]
fn local_research_cancel() -> Result<serde_json::Value, String> {
    loopback_json_post("/local/research/cancel", &serde_json::json!({}))
}

#[tauri::command]
fn local_research(coin: String) -> Result<serde_json::Value, String> {
    local_research_start(coin)
}

#[tauri::command]
fn ensure_companion() -> Result<bool, String> {
    start_companion();
    Ok(companion_listening())
}

fn companion_addr() -> SocketAddr {
    "127.0.0.1:17373".parse().expect("loopback companion")
}

fn companion_listening() -> bool {
    TcpStream::connect_timeout(&companion_addr(), Duration::from_millis(250)).is_ok()
}

fn loopback_json(path: &str) -> Result<serde_json::Value, String> {
    let raw = loopback_get(path)?;
    serde_json::from_str(&raw).map_err(|e| e.to_string())
}

fn loopback_json_post(path: &str, body: &serde_json::Value) -> Result<serde_json::Value, String> {
    loopback_json_post_timeout(path, body, 8)
}

fn loopback_json_post_timeout(path: &str, body: &serde_json::Value, read_secs: u64) -> Result<serde_json::Value, String> {
    let payload = body.to_string();
    let raw = loopback_exchange_timeout("POST", path, Some(&payload), read_secs)?;
    serde_json::from_str(&raw).map_err(|e| e.to_string())
}

fn loopback_get(path: &str) -> Result<String, String> {
    loopback_exchange("GET", path, None)
}

fn loopback_exchange(method: &str, path: &str, json: Option<&str>) -> Result<String, String> {
    loopback_exchange_timeout(method, path, json, 8)
}

fn loopback_exchange_timeout(method: &str, path: &str, json: Option<&str>, read_secs: u64) -> Result<String, String> {
    let mut stream = TcpStream::connect_timeout(&companion_addr(), Duration::from_secs(2))
        .map_err(|_| "companion_down".to_string())?;
    let _ = stream.set_read_timeout(Some(Duration::from_secs(read_secs)));
    let mut req = format!("{method} {path} HTTP/1.1\r\nHost: 127.0.0.1:17373\r\nConnection: close\r\n");
    if let Some(body) = json {
        req.push_str("Content-Type: application/json\r\n");
        req.push_str(&format!("Content-Length: {}\r\n\r\n{}", body.len(), body));
    } else {
        req.push_str("\r\n");
    }
    stream
        .write_all(req.as_bytes())
        .map_err(|_| "companion_down".to_string())?;
    let mut buf = Vec::new();
    stream
        .read_to_end(&mut buf)
        .map_err(|_| "companion_down".to_string())?;
    let text = String::from_utf8_lossy(&buf);
    let status_line = text.lines().next().unwrap_or("");
    let ok = status_line.contains(" 200 ") || status_line.ends_with(" 200");
    let body = text
        .split_once("\r\n\r\n")
        .or_else(|| text.split_once("\n\n"))
        .map(|(_, b)| b.trim().to_string())
        .unwrap_or_default();
    if !ok {
        if let Ok(v) = serde_json::from_str::<serde_json::Value>(&body) {
            if let Some(err) = v.get("error").and_then(|x| x.as_str()) {
                return Err(err.to_string());
            }
        }
        return Err("companion_http".into());
    }
    Ok(body)
}

fn install_dir() -> Option<PathBuf> {
    std::env::current_exe()
        .ok()
        .and_then(|p| p.parent().map(|d| d.to_path_buf()))
}

fn sidecar_paths() -> Vec<PathBuf> {
    let mut out = Vec::new();
    if let Some(dir) = install_dir() {
        for name in [
            "pit.exe",
            "pit",
            "pit-x86_64-pc-windows-msvc.exe",
            "pit-aarch64-pc-windows-msvc.exe",
        ] {
            out.push(dir.join(name));
        }
    }
    out
}

fn first_sidecar() -> Option<PathBuf> {
    sidecar_paths().into_iter().find(|p| p.exists())
}

fn log_file() -> Option<std::fs::File> {
    let base = std::env::var_os("LOCALAPPDATA").map(PathBuf::from)?;
    let dir = base.join("pit");
    fs::create_dir_all(&dir).ok()?;
    OpenOptions::new()
        .create(true)
        .append(true)
        .open(dir.join("companion.log"))
        .ok()
}

fn listening_pid() -> Option<u32> {
    let out = Command::new("netstat").args(["-ano", "-p", "tcp"]).output().ok()?;
    let text = String::from_utf8_lossy(&out.stdout);
    for line in text.lines() {
        if !line.contains("127.0.0.1:17373") {
            continue;
        }
        if !line.to_ascii_uppercase().contains("LISTEN") {
            continue;
        }
        let pid = line.split_whitespace().last()?.parse().ok()?;
        if pid > 0 {
            return Some(pid);
        }
    }
    None
}

fn process_path(pid: u32) -> Option<PathBuf> {
    let cmd = format!("(Get-CimInstance Win32_Process -Filter 'ProcessId={pid}').ExecutablePath");
    let out = Command::new("powershell")
        .args(["-NoProfile", "-NonInteractive", "-Command", &cmd])
        .output()
        .ok()?;
    let s = String::from_utf8_lossy(&out.stdout).trim().to_string();
    if s.is_empty() {
        None
    } else {
        Some(PathBuf::from(s))
    }
}

fn same_install(path: &Path) -> bool {
    let Some(dir) = install_dir() else {
        return false;
    };
    path.parent().map(|p| p == dir).unwrap_or(false)
}

const SIDECAR_VERSION: &str = "0.1.5";

fn companion_version() -> Option<String> {
    let raw = loopback_get("/health").ok()?;
    let v: serde_json::Value = serde_json::from_str(&raw).ok()?;
    v.get("version")
        .and_then(|x| x.as_str())
        .map(|s| s.trim().to_string())
}

fn stop_our_listener() {
    let Some(pid) = listening_pid() else {
        return;
    };
    match process_path(pid) {
        Some(path) if same_install(&path) => {}
        _ => return,
    }
    let _ = Command::new("taskkill")
        .args(["/F", "/PID", &pid.to_string()])
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status();
    for _ in 0..20 {
        if !companion_listening() {
            break;
        }
        std::thread::sleep(Duration::from_millis(100));
    }
}

fn spawn_sidecar(bin: &Path) {
    let mut cmd = Command::new(bin);
    cmd.arg("companion");
    if let Some(dir) = bin.parent() {
        cmd.current_dir(dir);
    }
    cmd.env("PIT_ALLOW_FALLBACKS", "false");
	cmd.env_remove("PIT_SESSION_KEY");
    cmd.env_remove("HL_SECRET");
    cmd.env_remove("PIT_KEYRING");
    cmd.env_remove("GIT_AUTHOR_DATE");
    cmd.env_remove("GIT_COMMITTER_DATE");
    if let Some(log) = log_file() {
        if let Ok(errlog) = log.try_clone() {
            cmd.stdout(Stdio::from(log));
            cmd.stderr(Stdio::from(errlog));
        }
    }
    #[cfg(windows)]
    {
        use std::os::windows::process::CommandExt;
        const CREATE_NO_WINDOW: u32 = 0x0800_0000;
        cmd.creation_flags(CREATE_NO_WINDOW);
    }
    let _ = cmd.spawn();
}

fn start_companion() {
    let bin = first_sidecar();
    if companion_listening() {
        let ours = listening_pid()
            .and_then(process_path)
            .map(|p| same_install(&p))
            .unwrap_or(false);
        let ver = companion_version();
        let stale = ours && ver.as_deref().is_some_and(|v| v != SIDECAR_VERSION);
        if stale {
            stop_our_listener();
        } else if companion_listening() {
            return;
        }
    }
    if let Some(bin) = bin {
        if !companion_listening() {
            spawn_sidecar(&bin);
        }
    } else {
        return;
    }
    for _ in 0..50 {
        if companion_listening() {
            return;
        }
        std::thread::sleep(Duration::from_millis(100));
    }
}

pub fn run() {
    start_companion();
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            companion_url,
            export_session,
            local_status,
            local_code,
            local_doctor,
            local_init,
            local_session,
            local_policy,
            local_revoke_session,
            local_direct_intent,
            local_direct_status,
            local_research,
            local_research_start,
            local_research_status,
            local_research_cancel,
            ensure_companion
        ])
        .run(tauri::generate_context!())
        .expect("PIT desktop failed to start");
}
