#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

fn main() {
    #[cfg(target_os = "linux")]
    {
        // Fixes WebKitGTK DMA-BUF Wayland protocol crash (Error 71) on Linux
        if std::env::var_os("__NV_DISABLE_EXPLICIT_SYNC").is_none() {
            unsafe {
                std::env::set_var("__NV_DISABLE_EXPLICIT_SYNC", "1");
            }
        }
    }

    tauri::Builder::default()
        .setup(|app| {
            use tauri::Manager;
            if let Some(window) = app.get_webview_window("main") {
                let _ = window.set_fullscreen(true);
            }
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
