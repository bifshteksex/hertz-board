use tauri::AppHandle;
use serde::{Deserialize, Serialize};
use std::path::PathBuf;

#[derive(Debug, Serialize, Deserialize)]
pub struct WorkspaceData {
    pub id: String,
    pub name: String,
    pub description: Option<String>,
    pub elements: String, // JSON string of canvas elements
    pub created_at: String,
    pub updated_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateInfo {
    pub available: bool,
    pub current_version: String,
    pub latest_version: Option<String>,
}

/// Open file dialog for importing workspace
#[tauri::command]
pub async fn open_file_dialog(app: AppHandle) -> Result<Option<PathBuf>, String> {
    use tauri_plugin_dialog::DialogExt;

    let file_path = app
        .dialog()
        .file()
        .add_filter("HertzBoard Workspace", &["hertzboard", "json"])
        .blocking_pick_file();

    Ok(file_path.and_then(|f| f.into_path().ok()))
}

/// Save file dialog for exporting workspace
#[tauri::command]
pub async fn save_file_dialog(
    app: AppHandle,
    default_name: String,
) -> Result<Option<PathBuf>, String> {
    use tauri_plugin_dialog::DialogExt;

    let file_path = app
        .dialog()
        .file()
        .add_filter("HertzBoard Workspace", &["hertzboard"])
        .set_file_name(&default_name)
        .blocking_save_file();

    Ok(file_path.and_then(|f| f.into_path().ok()))
}

/// Export workspace to file
#[tauri::command]
pub async fn export_workspace(
    file_path: String,
    workspace_data: WorkspaceData,
) -> Result<(), String> {
    use std::fs::File;
    use std::io::Write;

    let json_data = serde_json::to_string_pretty(&workspace_data)
        .map_err(|e| format!("Failed to serialize workspace: {}", e))?;

    let mut file = File::create(&file_path)
        .map_err(|e| format!("Failed to create file: {}", e))?;

    file.write_all(json_data.as_bytes())
        .map_err(|e| format!("Failed to write file: {}", e))?;

    Ok(())
}

/// Import workspace from file
#[tauri::command]
pub async fn import_workspace(file_path: String) -> Result<WorkspaceData, String> {
    use std::fs::File;
    use std::io::Read;

    let mut file = File::open(&file_path)
        .map_err(|e| format!("Failed to open file: {}", e))?;

    let mut contents = String::new();
    file.read_to_string(&mut contents)
        .map_err(|e| format!("Failed to read file: {}", e))?;

    let workspace_data: WorkspaceData = serde_json::from_str(&contents)
        .map_err(|e| format!("Failed to parse workspace data: {}", e))?;

    Ok(workspace_data)
}

/// Get offline workspaces from local database
#[tauri::command]
pub async fn get_offline_workspaces(app: AppHandle) -> Result<Vec<WorkspaceData>, String> {
    crate::db::get_offline_workspaces(&app)
}

/// Save workspace to local database for offline access
#[tauri::command]
pub async fn save_offline_workspace(
    app: AppHandle,
    workspace: WorkspaceData,
) -> Result<(), String> {
    crate::db::save_offline_workspace(&app, workspace)
}

/// Sync offline workspaces with server
#[tauri::command]
pub async fn sync_offline_workspaces(app: AppHandle) -> Result<Vec<String>, String> {
    // Get all offline workspaces
    let workspaces = crate::db::get_offline_workspaces(&app)?;

    let mut synced_ids = Vec::new();

    for workspace in workspaces {
        // Check if workspace needs sync (has local changes)
        if crate::db::workspace_needs_sync(&app, &workspace.id)? {
            synced_ids.push(workspace.id.clone());

            // Mark as synced
            crate::db::mark_workspace_synced(&app, &workspace.id)?;
        }
    }

    Ok(synced_ids)
}

/// Check for application updates
#[tauri::command]
pub async fn check_for_updates(app: AppHandle) -> Result<UpdateInfo, String> {
    use tauri_plugin_updater::UpdaterExt;

    let current_version = app.package_info().version.to_string();

    // Check for updates
    match app.updater() {
        Ok(updater) => {
            match updater.check().await {
                Ok(Some(update)) => {
                    Ok(UpdateInfo {
                        available: true,
                        current_version,
                        latest_version: Some(update.version),
                    })
                }
                Ok(None) => {
                    Ok(UpdateInfo {
                        available: false,
                        current_version,
                        latest_version: None,
                    })
                }
                Err(e) => {
                    log::error!("Failed to check for updates: {}", e);
                    Err(format!("Failed to check for updates: {}", e))
                }
            }
        }
        Err(e) => {
            log::error!("Updater not available: {}", e);
            Ok(UpdateInfo {
                available: false,
                current_version,
                latest_version: None,
            })
        }
    }
}
