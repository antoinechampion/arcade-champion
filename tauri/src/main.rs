#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::process::Command;
use std::env;
use std::path::PathBuf;

fn sidecar_path() -> PathBuf {
    let mut path = env::current_exe().expect("failed to get exe path");
    path.pop();
    // In dev: binary is in target/debug/, sidecar is in tauri/binaries/
    // In prod: both are in the same bundle directory
    let dev_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .join("binaries")
        .join(format!("arcade-champion-backend-{}", env!("TARGET")));
    if dev_path.exists() {
        return dev_path;
    }
    path.join("arcade-champion-backend")
}

fn main() {
    let backend_path = sidecar_path();
    println!("Starting backend: {:?}", backend_path);

    let _child = Command::new(&backend_path)
        .spawn()
        .unwrap_or_else(|e| panic!("failed to start backend at {:?}: {}", backend_path, e));

    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
