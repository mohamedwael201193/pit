use std::process::Command;

#[tauri::command]
fn companion_url() -> String {
    "http://127.0.0.1:17373".into()
}

#[tauri::command]
fn export_session() -> Result<String, String> {
    Err("session_export_denied".into())
}

fn spawn_companion() {
    if let Ok(exe) = std::env::current_exe() {
        if let Some(dir) = exe.parent() {
            for name in ["pit.exe", "pit"] {
                let bin = dir.join(name);
                if bin.exists() {
                    let mut cmd = Command::new(bin);
                    cmd.arg("companion");
                    #[cfg(windows)]
                    {
                        use std::os::windows::process::CommandExt;
                        const CREATE_NO_WINDOW: u32 = 0x0800_0000;
                        cmd.creation_flags(CREATE_NO_WINDOW);
                    }
                    let _ = cmd.spawn();
                    return;
                }
            }
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
