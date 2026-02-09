mod commands;
mod db;

use tauri::Manager;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_updater::Builder::new().build())
        .plugin(tauri_plugin_window_state::Builder::default().build())
        .plugin(tauri_plugin_store::Builder::default().build())
        .setup(|app| {
            if cfg!(debug_assertions) {
                app.handle().plugin(
                    tauri_plugin_log::Builder::default()
                        .level(log::LevelFilter::Info)
                        .build(),
                )?;
            }

            // Initialize database
            db::init_database(app.handle())?;

            // Setup global shortcuts (desktop only)
            #[cfg(not(any(target_os = "android", target_os = "ios")))]
            {
                use tauri_plugin_global_shortcut::{GlobalShortcutExt, Shortcut};

                // Ctrl+N - New workspace
                let new_workspace_shortcut: Shortcut = "Ctrl+N".parse().unwrap();
                app.global_shortcut().on_shortcut(new_workspace_shortcut, |app, _event, _shortcut| {
                    if let Some(window) = app.get_webview_window("main") {
                        let _ = window.eval("window.location.href = '/dashboard/new'");
                    }
                })?;

                // Ctrl+O - Open file
                let open_shortcut: Shortcut = "Ctrl+O".parse().unwrap();
                app.global_shortcut().on_shortcut(open_shortcut, |app, _event, _shortcut| {
                    if let Some(window) = app.get_webview_window("main") {
                        let _ = window.eval("window.dispatchEvent(new CustomEvent('tauri-open-file'))");
                    }
                })?;

                // Ctrl+S - Save/Export
                let save_shortcut: Shortcut = "Ctrl+S".parse().unwrap();
                app.global_shortcut().on_shortcut(save_shortcut, |app, _event, _shortcut| {
                    if let Some(window) = app.get_webview_window("main") {
                        let _ = window.eval("window.dispatchEvent(new CustomEvent('tauri-save-file'))");
                    }
                })?;

                app.handle().plugin(tauri_plugin_global_shortcut::Builder::new().build())?;
            }

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            commands::open_file_dialog,
            commands::save_file_dialog,
            commands::export_workspace,
            commands::import_workspace,
            commands::get_offline_workspaces,
            commands::save_offline_workspace,
            commands::sync_offline_workspaces,
            commands::check_for_updates,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
