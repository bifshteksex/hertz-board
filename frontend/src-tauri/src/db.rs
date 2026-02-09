use rusqlite::{Connection, params};
use tauri::{AppHandle, Manager};
use std::sync::Mutex;
use crate::commands::WorkspaceData;

pub struct DbState {
    pub conn: Mutex<Connection>,
}

/// Initialize SQLite database for offline storage
pub fn init_database(app: &AppHandle) -> Result<(), Box<dyn std::error::Error>> {
    let app_dir = app
        .path()
        .app_data_dir()
        .expect("Failed to get app data directory");

    std::fs::create_dir_all(&app_dir)?;

    let db_path = app_dir.join("hertzboard.db");
    let conn = Connection::open(db_path)?;

    // Create workspaces table
    conn.execute(
        "CREATE TABLE IF NOT EXISTS workspaces (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT,
            elements TEXT NOT NULL,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL,
            needs_sync INTEGER NOT NULL DEFAULT 0,
            last_synced_at TEXT
        )",
        [],
    )?;

    // Create index for faster queries
    conn.execute(
        "CREATE INDEX IF NOT EXISTS idx_workspaces_updated
         ON workspaces(updated_at DESC)",
        [],
    )?;

    app.manage(DbState {
        conn: Mutex::new(conn),
    });

    Ok(())
}

/// Get all offline workspaces
pub fn get_offline_workspaces(app: &AppHandle) -> Result<Vec<WorkspaceData>, String> {
    let state = app.state::<DbState>();
    let conn = state.conn.lock().unwrap();

    let mut stmt = conn
        .prepare(
            "SELECT id, name, description, elements, created_at, updated_at
             FROM workspaces
             ORDER BY updated_at DESC"
        )
        .map_err(|e| format!("Failed to prepare statement: {}", e))?;

    let workspaces = stmt
        .query_map([], |row| {
            Ok(WorkspaceData {
                id: row.get(0)?,
                name: row.get(1)?,
                description: row.get(2)?,
                elements: row.get(3)?,
                created_at: row.get(4)?,
                updated_at: row.get(5)?,
            })
        })
        .map_err(|e| format!("Failed to query workspaces: {}", e))?
        .collect::<Result<Vec<_>, _>>()
        .map_err(|e| format!("Failed to collect workspaces: {}", e))?;

    Ok(workspaces)
}

/// Save workspace to offline database
pub fn save_offline_workspace(
    app: &AppHandle,
    workspace: WorkspaceData,
) -> Result<(), String> {
    let state = app.state::<DbState>();
    let conn = state.conn.lock().unwrap();

    conn.execute(
        "INSERT OR REPLACE INTO workspaces
         (id, name, description, elements, created_at, updated_at, needs_sync)
         VALUES (?1, ?2, ?3, ?4, ?5, ?6, 1)",
        params![
            workspace.id,
            workspace.name,
            workspace.description,
            workspace.elements,
            workspace.created_at,
            workspace.updated_at,
        ],
    )
    .map_err(|e| format!("Failed to save workspace: {}", e))?;

    Ok(())
}

/// Check if workspace needs sync
pub fn workspace_needs_sync(app: &AppHandle, workspace_id: &str) -> Result<bool, String> {
    let state = app.state::<DbState>();
    let conn = state.conn.lock().unwrap();

    let needs_sync: i32 = conn
        .query_row(
            "SELECT needs_sync FROM workspaces WHERE id = ?1",
            params![workspace_id],
            |row| row.get(0),
        )
        .map_err(|e| format!("Failed to check sync status: {}", e))?;

    Ok(needs_sync == 1)
}

/// Mark workspace as synced
pub fn mark_workspace_synced(app: &AppHandle, workspace_id: &str) -> Result<(), String> {
    let state = app.state::<DbState>();
    let conn = state.conn.lock().unwrap();

    let now = chrono::Utc::now().to_rfc3339();

    conn.execute(
        "UPDATE workspaces SET needs_sync = 0, last_synced_at = ?1 WHERE id = ?2",
        params![now, workspace_id],
    )
    .map_err(|e| format!("Failed to mark workspace as synced: {}", e))?;

    Ok(())
}

/// Delete offline workspace
pub fn delete_offline_workspace(app: &AppHandle, workspace_id: &str) -> Result<(), String> {
    let state = app.state::<DbState>();
    let conn = state.conn.lock().unwrap();

    conn.execute(
        "DELETE FROM workspaces WHERE id = ?1",
        params![workspace_id],
    )
    .map_err(|e| format!("Failed to delete workspace: {}", e))?;

    Ok(())
}
