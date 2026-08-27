use std::net::{SocketAddr, TcpStream};
use std::path::PathBuf;
use std::process::Command;
use std::time::Duration;

#[tauri::command]
fn companion_url() -> String {
    "http://127.0.0.1:17373".into()
}

#[tauri::command]
fn export_session() -> Result<String, String> {
    Err("session_export_denied".into())
}

fn companion_listening() -> bool {
    let addr: SocketAddr = "127.0.0.1:17373".parse().expect("loopback companion");
    TcpStream::connect_timeout(&addr, Duration::from_millis(200)).is_ok()
}

fn sidecar_paths() -> Vec<PathBuf> {
    let mut out = Vec::new();
    if let Ok(exe) = std::env::current_exe() {
        if let Some(dir) = exe.parent() {
            for name in [
                "pit.exe",
                "pit",
                "pit-x86_64-pc-windows-msvc.exe",
                "pit-aarch64-pc-windows-msvc.exe",
            ] {
                out.push(dir.join(name));
            }
        }
    }
    out
}

fn spawn_companion() {
    if companion_listening() {
        return;
    }
    for bin in sidecar_paths() {
        if !bin.exists() {
            continue;
        }
        let mut cmd = Command::new(&bin);
        cmd.arg("companion");
        cmd.env("PIT_ALLOW_FALLBACKS", "false");
        cmd.env_remove("PIT_SESSION_KEY");
        cmd.env_remove("HL_SECRET");
        #[cfg(windows)]
        {
            use std::os::windows::process::CommandExt;
            const CREATE_NO_WINDOW: u32 = 0x0800_0000;
            cmd.creation_flags(CREATE_NO_WINDOW);
        }
        if cmd.spawn().is_ok() {
            return;
        }
    }
}

pub fn run() {
    spawn_companion();
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![companion_url, export_session])
        .run(tauri::generate_context!())
        .expect("PIT desktop failed to start");
}
