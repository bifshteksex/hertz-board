mod commands;
mod db;

use std::panic;
use tauri::Manager;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    // Setup panic hook BEFORE anything else to catch early crashes
    panic::set_hook(Box::new(|panic_info| {
        let panic_message = if let Some(s) = panic_info.payload().downcast_ref::<&str>() {
            s.to_string()
        } else if let Some(s) = panic_info.payload().downcast_ref::<String>() {
            s.clone()
        } else {
            "Unknown panic".to_string()
        };

        let location = if let Some(location) = panic_info.location() {
            format!(
                "{}:{}:{}",
                location.file(),
                location.line(),
                location.column()
            )
        } else {
            "Unknown location".to_string()
        };

        let error_msg = format!("PANIC: {} at {}", panic_message, location);

        // Try to log to stderr (will be visible in console if run from terminal)
        eprintln!("{}", error_msg);

        // Try to write to a crash log file
        if let Ok(app_data) = std::env::var("APPDATA") {
            let log_dir = std::path::Path::new(&app_data)
                .join("com.hertzboard.app")
                .join("logs");
            let _ = std::fs::create_dir_all(&log_dir);
            let crash_log = log_dir.join("crash.log");

            let timestamp = chrono::Utc::now().to_rfc3339();
            let crash_entry = format!("\n[{}] {}\n", timestamp, error_msg);

            let _ = std::fs::OpenOptions::new()
                .create(true)
                .append(true)
                .open(crash_log)
                .and_then(|mut f| {
                    use std::io::Write;
                    f.write_all(crash_entry.as_bytes())
                });
        }
    }));

    log::info!("Starting HertzBoard application...");

    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_updater::Builder::new().build())
        .plugin(tauri_plugin_window_state::Builder::default().build())
        .plugin(tauri_plugin_store::Builder::default().build())
        .plugin(tauri_plugin_deep_link::init())
        .plugin(
            // Enable logging in BOTH debug AND release builds
            tauri_plugin_log::Builder::default()
                .targets([
                    tauri_plugin_log::Target::new(tauri_plugin_log::TargetKind::Stdout),
                    tauri_plugin_log::Target::new(tauri_plugin_log::TargetKind::LogDir {
                        file_name: Some("hertzboard".to_string()),
                    }),
                    tauri_plugin_log::Target::new(tauri_plugin_log::TargetKind::Webview),
                ])
                .level(if cfg!(debug_assertions) {
                    log::LevelFilter::Debug
                } else {
                    log::LevelFilter::Info
                })
                .max_file_size(5_000_000) // 5MB max file size
                .rotation_strategy(tauri_plugin_log::RotationStrategy::KeepAll)
                .build(),
        )
        .setup(|app| {
            log::info!("=== HertzBoard Initialization Started ===");
            log::info!("Version: {}", app.package_info().version);
            log::info!("Platform: {}", std::env::consts::OS);
            log::info!("Architecture: {}", std::env::consts::ARCH);

            // Setup deep link handler for OAuth callbacks
            let app_handle = app.handle().clone();
            tauri_plugin_deep_link::register("hertzboard", move |request| {
                log::info!("Deep link received: {}", request);

                // Parse the URL to extract OAuth tokens
                // Expected format: hertzboard://oauth/callback?access_token=xxx&refresh_token=yyy
                if request.starts_with("hertzboard://oauth/callback") {
                    if let Some(window) = app_handle.get_webview_window("main") {
                        // Send the deep link URL to the frontend
                        let script = format!(
                            "window.dispatchEvent(new CustomEvent('oauth-callback', {{ detail: '{}' }}))",
                            request.replace("'", "\\'")
                        );
                        let _ = window.eval(&script);
                    }
                }
            })
            .map_err(|e| {
                log::error!("Failed to register deep link handler: {}", e);
                e
            })?;

            // Log where logs are stored
            if let Ok(log_dir) = app.path().app_log_dir() {
                log::info!("Log directory: {:?}", log_dir);
            }

            // Initialize database
            log::info!("Initializing database...");
            match db::init_database(app.handle()) {
                Ok(_) => log::info!("Database initialized successfully"),
                Err(e) => {
                    log::error!("Failed to initialize database: {}", e);
                    return Err(e);
                }
            }

            // Setup global shortcuts (desktop only)
            #[cfg(not(any(target_os = "android", target_os = "ios")))]
            {
                use tauri_plugin_global_shortcut::{GlobalShortcutExt, Shortcut};

                log::info!("Setting up global shortcuts...");

                // CRITICAL FIX: Register plugin BEFORE using it!
                match app
                    .handle()
                    .plugin(tauri_plugin_global_shortcut::Builder::new().build())
                {
                    Ok(_) => log::info!("Global shortcut plugin registered successfully"),
                    Err(e) => {
                        log::error!("Failed to register global shortcut plugin: {}", e);
                        return Err(e.into());
                    }
                }

                // Now we can safely use the global shortcut API
                // Ctrl+N - New workspace
                let new_workspace_shortcut: Shortcut = match "Ctrl+N".parse() {
                    Ok(s) => s,
                    Err(e) => {
                        log::error!("Failed to parse Ctrl+N shortcut: {}", e);
                        return Err(e.into());
                    }
                };

                if let Err(e) = app.global_shortcut().on_shortcut(
                    new_workspace_shortcut,
                    |app, _event, _shortcut| {
                        log::debug!("Ctrl+N shortcut triggered");
                        if let Some(window) = app.get_webview_window("main") {
                            let _ = window.eval("window.location.href = '/dashboard/new'");
                        }
                    },
                ) {
                    log::error!("Failed to register Ctrl+N shortcut: {}", e);
                    return Err(e.into());
                }
                log::info!("Registered shortcut: Ctrl+N");

                // Ctrl+O - Open file
                let open_shortcut: Shortcut = match "Ctrl+O".parse() {
                    Ok(s) => s,
                    Err(e) => {
                        log::error!("Failed to parse Ctrl+O shortcut: {}", e);
                        return Err(e.into());
                    }
                };

                if let Err(e) =
                    app.global_shortcut()
                        .on_shortcut(open_shortcut, |app, _event, _shortcut| {
                            log::debug!("Ctrl+O shortcut triggered");
                            if let Some(window) = app.get_webview_window("main") {
                                let _ = window.eval(
                                    "window.dispatchEvent(new CustomEvent('tauri-open-file'))",
                                );
                            }
                        })
                {
                    log::error!("Failed to register Ctrl+O shortcut: {}", e);
                    return Err(e.into());
                }
                log::info!("Registered shortcut: Ctrl+O");

                // Ctrl+S - Save/Export
                let save_shortcut: Shortcut = match "Ctrl+S".parse() {
                    Ok(s) => s,
                    Err(e) => {
                        log::error!("Failed to parse Ctrl+S shortcut: {}", e);
                        return Err(e.into());
                    }
                };

                if let Err(e) =
                    app.global_shortcut()
                        .on_shortcut(save_shortcut, |app, _event, _shortcut| {
                            log::debug!("Ctrl+S shortcut triggered");
                            if let Some(window) = app.get_webview_window("main") {
                                let _ = window.eval(
                                    "window.dispatchEvent(new CustomEvent('tauri-save-file'))",
                                );
                            }
                        })
                {
                    log::error!("Failed to register Ctrl+S shortcut: {}", e);
                    return Err(e.into());
                }
                log::info!("Registered shortcut: Ctrl+S");

                log::info!("All global shortcuts registered successfully");
            }

            log::info!("=== HertzBoard Initialization Complete ===");
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
