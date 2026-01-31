/**
 * Autosave Store
 *
 * Reactive store that wraps the autosave service and provides
 * Svelte 5 reactive state for autosave status.
 */

import { autosaveService, type SaveStatus } from '$lib/services/autosave';
import type { CanvasElement } from '$lib/types/api';

class AutosaveStore {
	private _status = $state<SaveStatus>('idle');
	private _pendingCount = $state(0);
	private _lastSaveTime = $state<Date | null>(null);
	private _lastError = $state<string | null>(null);

	// Callback for element ID changes
	public onElementIdChanged: ((frontendId: string, backendId: string) => void) | null = null;

	// Callback for getting full element data
	public getElementData: ((elementId: string) => CanvasElement | undefined) | null = null;

	constructor() {
		console.log('[AutosaveStore] Initializing...');

		// Subscribe to autosave service events
		autosaveService.onStatusChange = (status: SaveStatus) => {
			console.log(`[AutosaveStore] Status changed: ${status}`);
			this._status = status;
			if (status === 'error') {
				// Keep the error status visible
			} else if (status === 'saved') {
				this._lastSaveTime = new Date();
				this._lastError = null;
			}
		};

		autosaveService.onSaveComplete = (count: number) => {
			console.log(`[AutosaveStore] ✅ Saved ${count} changes successfully`);
		};

		autosaveService.onSaveError = (error: Error) => {
			this._lastError = error.message;
			console.error('[AutosaveStore] ❌ Save failed:', error);
		};

		autosaveService.onElementIdChanged = (frontendId: string, backendId: string) => {
			console.log(`[AutosaveStore] 🔄 Element ID changed: ${frontendId} → ${backendId}`);
			if (this.onElementIdChanged) {
				this.onElementIdChanged(frontendId, backendId);
			}
		};

		autosaveService.getElementData = (elementId: string) => {
			if (this.getElementData) {
				return this.getElementData(elementId);
			}
			return undefined;
		};

		// Update pending count periodically
		if (typeof window !== 'undefined') {
			setInterval(() => {
				const newCount = autosaveService.getPendingCount();
				if (newCount !== this._pendingCount) {
					console.log(`[AutosaveStore] Pending count: ${this._pendingCount} → ${newCount}`);
					this._pendingCount = newCount;
				}
			}, 100);
		}
	}

	get status() {
		return this._status;
	}

	get pendingCount() {
		return this._pendingCount;
	}

	get lastSaveTime() {
		return this._lastSaveTime;
	}

	get lastError() {
		return this._lastError;
	}

	get isSaving() {
		return this._status === 'saving';
	}

	get hasError() {
		return this._status === 'error';
	}

	get isSaved() {
		return this._status === 'saved';
	}

	/**
	 * Set workspace ID for autosave
	 */
	setWorkspaceId(workspaceId: string | null) {
		console.log(`[AutosaveStore] Setting workspace ID: ${workspaceId}`);
		autosaveService.setWorkspaceId(workspaceId);
	}

	/**
	 * Track a change for autosave
	 */
	trackChange(elementId: string, updates: Partial<CanvasElement>, isNew = false) {
		console.log(
			`[AutosaveStore] 📝 Tracking change for element ${elementId} (isNew: ${isNew}):`,
			updates
		);
		autosaveService.trackChange(elementId, updates, isNew);
		this._pendingCount = autosaveService.getPendingCount();
	}

	/**
	 * Track multiple changes
	 */
	trackChanges(changes: Array<{ id: string; updates: Partial<CanvasElement> }>) {
		console.log(`[AutosaveStore] 📝 Tracking ${changes.length} changes`);
		autosaveService.trackChanges(changes);
		this._pendingCount = autosaveService.getPendingCount();
	}

	/**
	 * Force immediate save
	 */
	async flush() {
		await autosaveService.flush();
		this._pendingCount = autosaveService.getPendingCount();
	}

	/**
	 * Clear pending changes
	 */
	clear() {
		autosaveService.clear();
		this._pendingCount = 0;
		this._lastError = null;
	}

	/**
	 * Pause autosave (e.g., during undo/redo)
	 */
	pause() {
		autosaveService.pause();
	}

	/**
	 * Resume autosave
	 */
	resume() {
		autosaveService.resume();
	}
}

export const autosaveStore = new AutosaveStore();
