import { invoke } from '@tauri-apps/api/core';
import type { CanvasElement } from '$lib/types/api';

export interface WorkspaceData {
	id: string;
	name: string;
	description?: string;
	elements: string; // JSON string of canvas elements
	created_at: string;
	updated_at: string;
}

export interface UpdateInfo {
	available: boolean;
	current_version: string;
	latest_version?: string;
}

/**
 * Check if running in Tauri environment
 */
export function isTauri(): boolean {
	return '__TAURI__' in window;
}

/**
 * Open file dialog for importing workspace
 */
export async function openFileDialog(): Promise<string | null> {
	if (!isTauri()) return null;
	return await invoke<string | null>('open_file_dialog');
}

/**
 * Save file dialog for exporting workspace
 */
export async function saveFileDialog(defaultName: string): Promise<string | null> {
	if (!isTauri()) return null;
	return await invoke<string | null>('save_file_dialog', { defaultName });
}

/**
 * Export workspace to file
 */
export async function exportWorkspace(
	filePath: string,
	workspaceId: string,
	workspaceName: string,
	elements: CanvasElement[]
): Promise<void> {
	if (!isTauri()) {
		throw new Error('Tauri is not available');
	}

	const workspaceData: WorkspaceData = {
		id: workspaceId,
		name: workspaceName,
		description: '',
		elements: JSON.stringify(elements),
		created_at: new Date().toISOString(),
		updated_at: new Date().toISOString()
	};

	await invoke('export_workspace', { filePath, workspaceData });
}

/**
 * Import workspace from file
 */
export async function importWorkspace(filePath: string): Promise<WorkspaceData> {
	if (!isTauri()) {
		throw new Error('Tauri is not available');
	}

	return await invoke<WorkspaceData>('import_workspace', { filePath });
}

/**
 * Get offline workspaces from local database
 */
export async function getOfflineWorkspaces(): Promise<WorkspaceData[]> {
	if (!isTauri()) return [];
	return await invoke<WorkspaceData[]>('get_offline_workspaces');
}

/**
 * Save workspace to local database for offline access
 */
export async function saveOfflineWorkspace(
	workspaceId: string,
	workspaceName: string,
	elements: CanvasElement[]
): Promise<void> {
	if (!isTauri()) {
		console.warn('Tauri is not available, skipping offline save');
		return;
	}

	const workspaceData: WorkspaceData = {
		id: workspaceId,
		name: workspaceName,
		description: '',
		elements: JSON.stringify(elements),
		created_at: new Date().toISOString(),
		updated_at: new Date().toISOString()
	};

	await invoke('save_offline_workspace', { workspace: workspaceData });
}

/**
 * Sync offline workspaces with server
 */
export async function syncOfflineWorkspaces(): Promise<string[]> {
	if (!isTauri()) return [];
	return await invoke<string[]>('sync_offline_workspaces');
}

/**
 * Check for application updates
 */
export async function checkForUpdates(): Promise<UpdateInfo> {
	if (!isTauri()) {
		return {
			available: false,
			current_version: '0.1.0'
		};
	}

	return await invoke<UpdateInfo>('check_for_updates');
}

/**
 * Parse elements from workspace data
 */
export function parseWorkspaceElements(workspaceData: WorkspaceData): CanvasElement[] {
	try {
		return JSON.parse(workspaceData.elements) as CanvasElement[];
	} catch (error) {
		console.error('Failed to parse workspace elements:', error);
		return [];
	}
}
