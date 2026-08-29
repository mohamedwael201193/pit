use serde::{Deserialize, Serialize};
use std::fs::{self, OpenOptions};
use std::io::{Read, Write};
use std::net::{SocketAddr, TcpStream};
use std::path::{Path, PathBuf};
use std::process::{Command, Stdio};
use std::sync::Mutex;
use std::time::Duration;
use tauri::{Manager, PhysicalPosition, PhysicalSize};

#[derive(Serialize, Deserialize)]
struct Bounds {
    x: i32,
    y: i32,
    w: u32,
    h: u32,
}

fn bounds_file() -> Option<PathBuf> {
    let base = std::env::var_os("APPDATA")?;
    Some(PathBuf::from(base).join("pit").join("window-bounds.json"))
}

fn restore_bounds(window: &tauri::WebviewWindow) {
    let Some(path) = bounds_file() else {
        return;
    };
    let Ok(raw) = fs::read_to_string(path) else {
        return;
    };
    let Ok(b) = serde_json::from_str::<Bounds>(&raw) else {
        return;
    };
    if b.w < 1100 || b.h < 720 {
        return;
    }
    let _ = window.set_position(PhysicalPosition::new(b.x, b.y));
    let _ = window.set_size(PhysicalSize::new(b.w, b.h));
}

fn save_bounds(window: &tauri::Window) {
    let Ok(pos) = window.outer_position() else {
        return;
    };
    let Ok(size) = window.outer_size() else {
        return;
    };
    if size.width < 1100 || size.height < 720 {
        return;
    }
    let Some(path) = bounds_file() else {
        return;
    };
    if let Some(dir) = path.parent() {
        let _ = fs::create_dir_all(dir);
    }
    let body = Bounds {
        x: pos.x,
        y: pos.y,
        w: size.width,
        h: size.height,
    };
    if let Ok(raw) = serde_json::to_string(&body) {
        let _ = fs::write(path, raw);
    }
}

#[tauri::command]
fn companion_url() -> String {
    "http://127.0.0.1:17373".into()
}

#[tauri::command]
fn export_session() -> Result<String, String> {
    Err("session_export_denied".into())
}

async fn json_get(path: String, secs: u64) -> Result<serde_json::Value, String> {
    tauri::async_runtime::spawn_blocking(move || {
        let raw = loopback_exchange_timeout("GET", &path, None, secs)?;
        serde_json::from_str(&raw).map_err(|e| e.to_string())
    })
    .await
    .map_err(|_| "companion_down".to_string())?
}

async fn json_post(path: String, body: serde_json::Value, secs: u64) -> Result<serde_json::Value, String> {
    tauri::async_runtime::spawn_blocking(move || loopback_json_post_timeout(&path, &body, secs))
        .await
        .map_err(|_| "companion_down".to_string())?
}

#[tauri::command]
async fn local_status() -> Result<serde_json::Value, String> {
    json_get("/local/status".into(), 2).await
}

#[tauri::command]
async fn local_code() -> Result<serde_json::Value, String> {
    json_get("/local/code".into(), 2).await
}

#[tauri::command]
async fn local_doctor() -> Result<serde_json::Value, String> {
    json_get("/local/doctor".into(), 12).await
}

#[tauri::command]
async fn local_init(wallet: String, network: String) -> Result<serde_json::Value, String> {
    let body = serde_json::json!({ "wallet": wallet, "network": network });
    json_post("/local/init".into(), body, 8).await
}

#[tauri::command]
async fn local_session() -> Result<serde_json::Value, String> {
    json_post("/local/session".into(), serde_json::json!({}), 8).await
}

#[tauri::command]
async fn local_connection_preview(coin: Option<String>, reduce_only: Option<bool>) -> Result<serde_json::Value, String> {
    let body = serde_json::json!({
        "coin": coin.unwrap_or_else(|| "ETH".into()),
        "reduceOnly": reduce_only.unwrap_or(false)
    });
    json_post("/local/connection-preview".into(), body, 20).await
}

#[tauri::command]
async fn local_policy(policy: Option<serde_json::Value>) -> Result<serde_json::Value, String> {
    json_post("/local/policy".into(), policy.unwrap_or_else(|| serde_json::json!({})), 8).await
}

#[tauri::command]
async fn local_policy_get() -> Result<serde_json::Value, String> {
    json_get("/local/policy".into(), 8).await
}

#[tauri::command]
async fn local_revoke_session() -> Result<serde_json::Value, String> {
    json_post("/local/revoke-session".into(), serde_json::json!({}), 8).await
}

#[tauri::command]
async fn local_direct_intent() -> Result<serde_json::Value, String> {
    json_get("/local/direct-intent".into(), 20).await
}

#[tauri::command]
async fn local_direct_status() -> Result<serde_json::Value, String> {
    json_get("/local/direct-status".into(), 2).await
}

#[tauri::command]
async fn local_research_start(coin: String, hypothesis: Option<String>) -> Result<serde_json::Value, String> {
    let mut body = serde_json::json!({ "coin": coin });
    if let Some(h) = hypothesis {
        if !h.trim().is_empty() {
            body["hypothesis"] = serde_json::Value::String(h);
        }
    }
    json_post("/local/research/start".into(), body, 8).await
}

#[tauri::command]
async fn local_research_status() -> Result<serde_json::Value, String> {
    json_get("/local/research/status".into(), 15).await
}

#[tauri::command]
async fn local_research_result() -> Result<serde_json::Value, String> {
    json_get("/local/research/result".into(), 8).await
}

#[tauri::command]
async fn local_research_cancel() -> Result<serde_json::Value, String> {
    json_post("/local/research/cancel".into(), serde_json::json!({}), 8).await
}

#[tauri::command]
async fn local_authorize(typed: String, hash: String) -> Result<serde_json::Value, String> {
    let body = serde_json::json!({ "typed": typed, "hash": hash });
    json_post("/local/authorize".into(), body, 30).await
}

#[tauri::command]
async fn local_cancel_order(typed: String) -> Result<serde_json::Value, String> {
    let body = serde_json::json!({ "typed": typed });
    json_post("/local/cancel".into(), body, 30).await
}

#[tauri::command]
async fn local_watch(network: String) -> Result<serde_json::Value, String> {
    let net = if network == "testnet" { "testnet" } else { "mainnet" };
    json_get(format!("/watch?network={net}"), 8).await
}

#[tauri::command]
async fn local_kill(on: bool) -> Result<serde_json::Value, String> {
    let body = serde_json::json!({ "on": on });
    json_post("/local/kill".into(), body, 8).await
}

#[tauri::command]
async fn local_activity() -> Result<serde_json::Value, String> {
    json_get("/local/activity".into(), 8).await
}

#[tauri::command]
async fn local_positions() -> Result<serde_json::Value, String> {
    json_get("/local/positions".into(), 12).await
}

#[tauri::command]
async fn local_chat(text: String, thread: Option<String>) -> Result<serde_json::Value, String> {
    json_post(
        "/local/chat".into(),
        serde_json::json!({ "text": text, "thread": thread.unwrap_or_else(|| "desk".into()) }),
        20,
    )
    .await
}

#[tauri::command]
async fn local_chat_log(thread: Option<String>) -> Result<serde_json::Value, String> {
    let t = thread.unwrap_or_else(|| "desk".into());
    json_get(format!("/local/chat/log?thread={t}"), 8).await
}

#[tauri::command]
async fn local_chat_threads() -> Result<serde_json::Value, String> {
    json_get("/local/chat/threads".into(), 8).await
}

#[tauri::command]
async fn local_chat_thread(op: String, id: Option<String>, title: Option<String>) -> Result<serde_json::Value, String> {
    json_post(
        "/local/chat/thread".into(),
        serde_json::json!({ "op": op, "id": id, "title": title }),
        8,
    )
    .await
}

#[tauri::command]
fn window_min(window: tauri::Window) {
    let _ = window.minimize();
}

#[tauri::command]
fn window_toggle_max(window: tauri::Window) {
    if window.is_maximized().unwrap_or(false) {
        let _ = window.unmaximize();
    } else {
        let _ = window.maximize();
    }
}

#[tauri::command]
fn open_url(url: String) -> Result<bool, String> {
    let url = allowed_https(&url)?;
    open_default_browser(&url)
}

fn allowed_https(url: &str) -> Result<String, String> {
    let u = url.trim();
    if u.len() > 2048 {
        return Err("url_too_long".into());
    }
    if !u.starts_with("https://") {
        return Err("https_required".into());
    }
    if u.bytes().any(|b| b < 0x20 || matches!(b, b'"' | b'\\' | b'<' | b'>' | b'`' | b'|')) {
        return Err("url_denied".into());
    }
    let rest = &u[8..];
    if rest.contains('@') {
        return Err("url_denied".into());
    }
    let hostport = rest.split(|c| c == '/' || c == '?' || c == '#').next().unwrap_or("");
    if hostport.contains(':') {
        return Err("url_denied".into());
    }
    let host = hostport.to_ascii_lowercase();
    const ALLOW: &[&str] = &[
        "pit0g.vercel.app",
        "app.hyperliquid.xyz",
        "app.hyperliquid-testnet.xyz",
        "pc.0g.ai",
        "0g.ai",
        "www.0g.ai",
        "docs.0g.ai",
        "chainscan.0g.ai",
        "chainscan-galileo.0g.ai",
        "hyperliquid.info",
        "github.com",
    ];
    if !ALLOW.iter().any(|h| host == *h) {
        return Err("host_denied".into());
    }
    if host == "github.com"
        && !rest
            .to_ascii_lowercase()
            .starts_with("github.com/mohamedwael201193/pit")
    {
        return Err("host_denied".into());
    }
    Ok(u.to_string())
}

#[cfg(windows)]
fn open_default_browser(url: &str) -> Result<bool, String> {
    use std::os::windows::process::CommandExt;
    const CREATE_NO_WINDOW: u32 = 0x0800_0000;
    let status = Command::new("rundll32")
        .args(["url.dll,FileProtocolHandler", url])
        .creation_flags(CREATE_NO_WINDOW)
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()
        .map_err(|e| e.to_string())?;
    if status.success() {
        Ok(true)
    } else {
        Err("open_failed".into())
    }
}

#[cfg(not(windows))]
fn open_default_browser(url: &str) -> Result<bool, String> {
    let status = Command::new("xdg-open")
        .arg(url)
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()
        .map_err(|e| e.to_string())?;
    if status.success() {
        Ok(true)
    } else {
        Err("open_failed".into())
    }
}

#[tauri::command]
fn window_close(window: tauri::Window) {
    save_bounds(&window);
    let _ = window.close();
}

#[tauri::command]
async fn local_memory_forget() -> Result<serde_json::Value, String> {
    json_post("/local/memory/forget".into(), serde_json::json!({}), 8).await
}

#[tauri::command]
async fn local_calibration() -> Result<serde_json::Value, String> {
    json_get("/local/calibration".into(), 8).await
}

#[tauri::command]
async fn local_security() -> Result<serde_json::Value, String> {
    json_get("/local/security".into(), 12).await
}

#[tauri::command]
async fn local_identity() -> Result<serde_json::Value, String> {
    json_get("/local/identity".into(), 8).await
}

#[tauri::command]
async fn local_update() -> Result<serde_json::Value, String> {
    json_get("/local/update".into(), 4).await
}

#[tauri::command]
async fn local_explain() -> Result<serde_json::Value, String> {
    json_get("/local/explain".into(), 8).await
}

#[tauri::command]
async fn local_models() -> Result<serde_json::Value, String> {
    json_get("/local/models".into(), 12).await
}

#[tauri::command]
async fn local_research(coin: String) -> Result<serde_json::Value, String> {
    local_research_start(coin, None).await
}

#[tauri::command]
async fn ensure_companion() -> Result<bool, String> {
    tauri::async_runtime::spawn_blocking(|| {
        start_companion();
        companion_listening()
    })
    .await
    .map_err(|_| "companion_down".to_string())
}

fn companion_addr() -> SocketAddr {
    "127.0.0.1:17373".parse().expect("loopback companion")
}

fn companion_listening() -> bool {
    TcpStream::connect_timeout(&companion_addr(), Duration::from_millis(400)).is_ok()
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

static LOOPBACK: Mutex<Option<TcpStream>> = Mutex::new(None);
static STARTING: Mutex<()> = Mutex::new(());

fn take_loopback() -> Result<TcpStream, String> {
    if let Ok(mut slot) = LOOPBACK.lock() {
        if let Some(stream) = slot.take() {
            return Ok(stream);
        }
    }
    let stream = TcpStream::connect_timeout(&companion_addr(), Duration::from_millis(400))
        .map_err(|_| "companion_down".to_string())?;
    let _ = stream.set_nodelay(true);
    Ok(stream)
}

fn drop_loopback() {
    if let Ok(mut slot) = LOOPBACK.lock() {
        *slot = None;
    }
}

fn header_end(buf: &[u8]) -> Option<usize> {
    buf.windows(4).position(|w| w == b"\r\n\r\n").map(|i| i + 4)
}

fn http_status(header: &str) -> u16 {
    header
        .split_whitespace()
        .nth(1)
        .and_then(|s| s.parse().ok())
        .unwrap_or(0)
}

fn content_length(header: &str) -> Option<usize> {
    for line in header.lines() {
        let (k, v) = line.split_once(':')?;
        if k.eq_ignore_ascii_case("content-length") {
            return v.trim().parse().ok();
        }
    }
    None
}

fn read_http(stream: &mut TcpStream) -> Result<(u16, String), String> {
    let mut buf = Vec::new();
    let mut tmp = [0u8; 2048];
    loop {
        let n = stream.read(&mut tmp).map_err(|_| "companion_down".to_string())?;
        if n == 0 {
            return Err("companion_down".into());
        }
        buf.extend_from_slice(&tmp[..n]);
        if let Some(end) = header_end(&buf) {
            let header = String::from_utf8_lossy(&buf[..end]).into_owned();
            let status = http_status(&header);
            let want = content_length(&header).unwrap_or(0);
            let mut body = buf[end..].to_vec();
            while body.len() < want {
                let n = stream.read(&mut tmp).map_err(|_| "companion_down".to_string())?;
                if n == 0 {
                    break;
                }
                body.extend_from_slice(&tmp[..n]);
            }
            if want > 0 {
                body.truncate(want);
            }
            return Ok((status, String::from_utf8_lossy(&body).trim().to_string()));
        }
        if buf.len() > 2_000_000 {
            return Err("companion_down".into());
        }
    }
}

fn loopback_exchange_timeout(method: &str, path: &str, json: Option<&str>, read_secs: u64) -> Result<String, String> {
    let mut stream = take_loopback()?;
    let _ = stream.set_read_timeout(Some(Duration::from_secs(read_secs)));
    let mut req = format!("{method} {path} HTTP/1.1\r\nHost: 127.0.0.1:17373\r\nConnection: close\r\n");
    if let Some(body) = json {
        req.push_str("Content-Type: application/json\r\n");
        req.push_str(&format!("Content-Length: {}\r\n\r\n{}", body.len(), body));
    } else {
        req.push_str("Content-Length: 0\r\n\r\n");
    }
    if stream.write_all(req.as_bytes()).is_err() {
        drop_loopback();
        let mut stream = TcpStream::connect_timeout(&companion_addr(), Duration::from_millis(400))
            .map_err(|_| "companion_down".to_string())?;
        let _ = stream.set_nodelay(true);
        let _ = stream.set_read_timeout(Some(Duration::from_secs(read_secs)));
        stream
            .write_all(req.as_bytes())
            .map_err(|_| "companion_down".to_string())?;
        let (status, body) = read_http(&mut stream)?;
        drop_loopback();
        return finish_http(status, body);
    }
    match read_http(&mut stream) {
        Ok((status, body)) => {
            drop_loopback();
            finish_http(status, body)
        }
        Err(e) => {
            drop_loopback();
            Err(e)
        }
    }
}

fn finish_http(status: u16, body: String) -> Result<String, String> {
    if status != 200 {
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

const SIDECAR_VERSION: &str = "0.5.0";

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

fn apply_sponsor_env(cmd: &mut Command, bin: &Path) {
    if std::env::var_os("PIT_DIRECT_SPONSOR_FILE").is_some() {
        return;
    }
    if let Some(dir) = bin.parent() {
        let p = dir.join("direct-sponsor.json");
        if p.is_file() {
            cmd.env("PIT_DIRECT_SPONSOR_FILE", p);
            return;
        }
    }
    if let Some(base) = std::env::var_os("LOCALAPPDATA") {
        let p = PathBuf::from(base).join("PIT").join("direct-sponsor.json");
        if p.is_file() {
            cmd.env("PIT_DIRECT_SPONSOR_FILE", p);
        }
    }
}

fn spawn_sidecar(bin: &Path) {
    let mut cmd = Command::new(bin);
    cmd.arg("companion");
    if let Some(dir) = bin.parent() {
        cmd.current_dir(dir);
    }
    cmd.env("PIT_ALLOW_FALLBACKS", "false");
    cmd.env("PIT_COMPANION", "1");
    apply_sponsor_env(&mut cmd, bin);
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

fn sealed_job_active() -> bool {
    let Ok(base) = std::env::var("APPDATA") else {
        return false;
    };
    let raw = fs::read_to_string(PathBuf::from(base).join("pit").join("research-job.json")).unwrap_or_default();
    raw.contains("\"running\":true")
}

fn start_companion() {
    let _gate = STARTING.lock().unwrap_or_else(|e| e.into_inner());
    if companion_listening() && sealed_job_active() {
        return;
    }
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
        .setup(|app| {
            if let Some(w) = app.get_webview_window("main") {
                restore_bounds(&w);
                let _ = w.show();
                let _ = w.set_focus();
            }
            Ok(())
        })
        .on_window_event(|window, event| {
            if matches!(
                event,
                tauri::WindowEvent::Moved(_) | tauri::WindowEvent::Resized(_)
            ) {
                save_bounds(window);
            }
        })
        .invoke_handler(tauri::generate_handler![
            companion_url,
            export_session,
            local_status,
            local_code,
            local_doctor,
            local_init,
            local_session,
            local_connection_preview,
            local_policy,
            local_policy_get,
            local_revoke_session,
            local_direct_intent,
            local_direct_status,
            local_research,
            local_research_start,
            local_research_status,
            local_research_result,
            local_research_cancel,
            local_authorize,
            local_cancel_order,
            local_watch,
            local_kill,
            local_activity,
            local_positions,
            local_chat,
            local_chat_log,
            local_chat_threads,
            local_chat_thread,
            local_memory_forget,
            local_calibration,
            local_security,
            local_identity,
            local_update,
            local_explain,
            local_models,
            ensure_companion,
            open_url,
            window_min,
            window_toggle_max,
            window_close
        ])
        .run(tauri::generate_context!())
        .expect("PIT desktop failed to start");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn allow_official_https() {
        assert!(allowed_https("https://app.hyperliquid.xyz/trade/BTC").is_ok());
        assert!(allowed_https("https://app.hyperliquid-testnet.xyz/API").is_ok());
        assert!(allowed_https("https://chainscan.0g.ai/tx/0xabc").is_ok());
        assert!(allowed_https("https://github.com/mohamedwael201193/pit/releases/latest").is_ok());
        assert!(allowed_https("https://pit0g.vercel.app/pair").is_ok());
        assert!(allowed_https("https://pc.0g.ai/sdk/dashboard/funds").is_ok());
        assert!(allowed_https("https://hyperliquid.info").is_ok());
        assert!(allowed_https("http://app.hyperliquid.xyz").is_err());
        assert!(allowed_https("https://evil.example").is_err());
        assert!(allowed_https("https://github.com/other/repo").is_err());
        assert!(allowed_https("https://user@app.hyperliquid.xyz").is_err());
        assert!(allowed_https("https://app.hyperliquid.xyz:443/x").is_err());
        assert!(allowed_https("javascript:alert(1)").is_err());
    }
}
