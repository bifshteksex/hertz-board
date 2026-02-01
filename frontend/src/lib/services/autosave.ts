/**
 * Autosave Service
 *
 * Manages automatic saving of canvas elements with debouncing and batching.
 * - Tracks changes to canvas elements
 * - Debounces save operations to reduce API calls
 * - Batches multiple changes into single API request
 * - Handles error recovery and retry logic
 */

import type { CanvasElement } from '$lib/types/api';
import { api } from './api';

export type SaveStatus = 'idle' | 'saving' | 'saved' | 'error';

export interface AutosaveConfig {
	/** Debounce delay in milliseconds (default: 1000ms) */
	debounceMs?: number;
	/** Maximum number of items to batch (default: 50) */
	maxBatchSize?: number;
	/** Enable debug logging */
	debug?: boolean;
	/** Callback when save status changes */
	onStatusChange?: (status: SaveStatus) => void;
	/** Callback when save completes */
	onSaveComplete?: (count: number) => void;
	/** Callback when save fails */
	onSaveError?: (error: Error) => void;
	/** Callback when element IDs change after creation (frontendId → backendId) */
	onElementIdChanged?: (frontendId: string, backendId: string) => void;
}

interface PendingChange {
	elementId: string;
	updates: Partial<CanvasElement>;
	timestamp: number;
	isNew?: boolean; // Флаг для новых элементов (CREATE)
}

interface ElementData {
	position: { x: number | undefined; y: number | undefined };
	size: { width: number; height: number };
	style?: Record<string, unknown>;
	rotation: number;
	shape_type?: string;
	content?: string;
	html_content?: string;
}

interface BatchUpdateItem {
	id: string;
	element_data?: Partial<ElementData>;
	z_index?: number;
	parent_id?: string | null;
}

interface BatchCreateResponse {
	elements: Array<{ id: string }>;
}

export class AutosaveService {
	private workspaceId: string | null = null;
	private pendingChanges = new Map<string, PendingChange>();
	private debounceTimer: number | null = null;
	private status: SaveStatus = 'idle';
	private saveInProgress = false;

	// Public callbacks that can be set externally
	public onStatusChange: (status: SaveStatus) => void = () => {};
	public onSaveComplete: (count: number) => void = () => {};
	public onSaveError: (error: Error) => void = () => {};
	public onElementIdChanged: (frontendId: string, backendId: string) => void = () => {};
	public getElementData: (elementId: string) => CanvasElement | undefined = () => undefined;

	// Config
	private debounceMs: number;
	private maxBatchSize: number;
	private debug: boolean;

	constructor(config: AutosaveConfig = {}) {
		this.debounceMs = config.debounceMs ?? 1000;
		this.maxBatchSize = config.maxBatchSize ?? 50;
		this.debug = config.debug ?? false;

		if (config.onStatusChange) this.onStatusChange = config.onStatusChange;
		if (config.onSaveComplete) this.onSaveComplete = config.onSaveComplete;
		if (config.onSaveError) this.onSaveError = config.onSaveError;
		if (config.onElementIdChanged) this.onElementIdChanged = config.onElementIdChanged;
	}

	/**
	 * Set the current workspace ID
	 */
	setWorkspaceId(workspaceId: string | null) {
		if (this.workspaceId !== workspaceId) {
			// Workspace changed, flush pending changes for old workspace
			if (this.pendingChanges.size > 0 && this.workspaceId) {
				this.log('Workspace changed, flushing pending changes');
				this.flush();
			}
			this.workspaceId = workspaceId;
		}
	}

	/**
	 * Track a change to an element
	 */
	trackChange(elementId: string, updates: Partial<CanvasElement>, isNew = false) {
		console.log(`[Autosave] trackChange called for ${elementId} (isNew: ${isNew}):`, updates);

		if (!this.workspaceId) {
			console.warn('[Autosave] ⚠️ No workspace ID set, skipping autosave');
			this.log('No workspace ID set, skipping autosave');
			return;
		}

		// Merge with existing pending changes for this element
		const existing = this.pendingChanges.get(elementId);
		const merged = existing ? { ...existing.updates, ...updates } : updates;

		this.pendingChanges.set(elementId, {
			elementId,
			updates: merged,
			timestamp: Date.now(),
			isNew: isNew || existing?.isNew // Если был new, остается new
		});

		console.log(
			`[Autosave] Tracked change for element ${elementId}, total pending: ${this.pendingChanges.size}`
		);
		this.log(`Tracked change for element ${elementId}, total pending: ${this.pendingChanges.size}`);

		// Update status
		this.setStatus('idle');

		// Reset debounce timer
		this.scheduleSave();
	}

	/**
	 * Track multiple changes at once
	 */
	trackChanges(changes: Array<{ id: string; updates: Partial<CanvasElement> }>) {
		changes.forEach(({ id, updates }) => this.trackChange(id, updates));
	}

	/**
	 * Schedule a save with debouncing
	 */
	private scheduleSave() {
		// Clear existing timer
		if (this.debounceTimer !== null) {
			clearTimeout(this.debounceTimer);
		}

		// Schedule new save
		this.debounceTimer = window.setTimeout(() => {
			this.save();
		}, this.debounceMs);
	}

	/**
	 * Force immediate save without debouncing
	 */
	async flush(): Promise<void> {
		if (this.debounceTimer !== null) {
			clearTimeout(this.debounceTimer);
			this.debounceTimer = null;
		}
		return this.save();
	}

	/**
	 * Perform the actual save operation
	 */
	private async save(): Promise<void> {
		console.log('[Autosave] 🔵 save() called');

		if (!this.workspaceId) {
			console.warn('[Autosave] ⚠️ No workspace ID, cannot save');
			this.log('No workspace ID, cannot save');
			return;
		}

		if (this.pendingChanges.size === 0) {
			console.log('[Autosave] ℹ️ No pending changes, skipping save');
			this.log('No pending changes, skipping save');
			return;
		}

		if (this.saveInProgress) {
			console.log('[Autosave] ⏳ Save already in progress, rescheduling');
			this.log('Save already in progress, rescheduling');
			this.scheduleSave();
			return;
		}

		// Get changes to save
		const changesToSave = Array.from(this.pendingChanges.values());
		const workspaceId = this.workspaceId;

		console.log(`[Autosave] 💾 Saving ${changesToSave.length} changes to workspace ${workspaceId}`);
		console.log('[Autosave] Changes to save:', changesToSave);
		this.log(`Saving ${changesToSave.length} changes to workspace ${workspaceId}`);
		this.setStatus('saving');
		this.saveInProgress = true;

		try {
			// Разделяем на CREATE и UPDATE
			const newElements = changesToSave.filter((c) => c.isNew);
			const updatedElements = changesToSave.filter((c) => !c.isNew);

			console.log(
				`[Autosave] Changes: ${newElements.length} new, ${updatedElements.length} updates`
			);

			// Сохраняем новые элементы (CREATE)
			if (newElements.length > 0) {
				const batches = this.createBatches(newElements);
				console.log(
					`[Autosave] Creating ${newElements.length} new elements in ${batches.length} batch(es)`
				);

				for (let i = 0; i < batches.length; i++) {
					const batch = batches[i];
					console.log(
						`[Autosave] Creating batch ${i + 1}/${batches.length} (${batch.length} items)`
					);
					const idMapping = await this.saveBatchCreate(workspaceId, batch);
					console.log(`[Autosave] ✅ Create batch ${i + 1}/${batches.length} saved successfully`);

					// Обновляем ID элементов в pendingChanges и уведомляем canvas store
					if (idMapping) {
						console.log('[Autosave] Updating element IDs after creation:', idMapping);
						idMapping.forEach((backendId, frontendId) => {
							// Уведомляем canvas store об изменении ID
							this.onElementIdChanged(frontendId, backendId);

							// Если есть pending changes для старого ID, переносим их на новый ID
							const pendingChange = this.pendingChanges.get(frontendId);
							if (pendingChange) {
								this.pendingChanges.delete(frontendId);
								this.pendingChanges.set(backendId, {
									...pendingChange,
									elementId: backendId,
									isNew: false // Теперь элемент уже создан
								});
								console.log(`[Autosave] Remapped pending changes: ${frontendId} → ${backendId}`);
							}
						});
					}
				}
			}

			// Обновляем существующие элементы (UPDATE)
			if (updatedElements.length > 0) {
				const batches = this.createBatches(updatedElements);
				console.log(
					`[Autosave] Updating ${updatedElements.length} elements in ${batches.length} batch(es)`
				);

				for (let i = 0; i < batches.length; i++) {
					const batch = batches[i];
					console.log(
						`[Autosave] Updating batch ${i + 1}/${batches.length} (${batch.length} items)`
					);
					await this.saveBatchUpdate(workspaceId, batch);
					console.log(`[Autosave] ✅ Update batch ${i + 1}/${batches.length} saved successfully`);
				}
			}

			// Clear saved changes
			changesToSave.forEach((change) => {
				this.pendingChanges.delete(change.elementId);
			});

			console.log(`[Autosave] ✅ Successfully saved ${changesToSave.length} changes`);
			this.log(`Successfully saved ${changesToSave.length} changes`);
			this.setStatus('saved');
			this.onSaveComplete(changesToSave.length);

			// Reset to idle after a short delay
			setTimeout(() => {
				if (this.status === 'saved') {
					this.setStatus('idle');
				}
			}, 2000);
		} catch (error) {
			console.error('[Autosave] ❌ Save failed:', error);
			this.log(`Save failed: ${error}`);
			this.setStatus('error');
			this.onSaveError(error as Error);

			// Retry after a delay
			setTimeout(() => {
				if (this.pendingChanges.size > 0) {
					console.log('[Autosave] 🔄 Retrying failed save in 5 seconds...');
					this.log('Retrying failed save');
					this.scheduleSave();
				}
			}, 5000);
		} finally {
			this.saveInProgress = false;
		}
	}

	/**
	 * Split changes into batches
	 */
	private createBatches(changes: PendingChange[]): PendingChange[][] {
		const batches: PendingChange[][] = [];
		for (let i = 0; i < changes.length; i += this.maxBatchSize) {
			batches.push(changes.slice(i, i + this.maxBatchSize));
		}
		return batches;
	}

	/**
	 * Преобразовать frontend element type в backend element type
	 */
	private mapElementType(frontendType: string): string {
		// Frontend использует специфичные типы: rectangle, ellipse, triangle, arrow
		// Backend использует обобщенный тип: shape
		const shapeTypes = ['rectangle', 'ellipse', 'triangle', 'arrow', 'line'];
		if (shapeTypes.includes(frontendType)) {
			return 'shape';
		}
		// Остальные типы совпадают: text, image, drawing, sticky, list, connector, group
		return frontendType;
	}

	/**
	 * Save a single batch of NEW elements (CREATE)
	 * Returns a Map of frontend ID → backend ID
	 */
	private async saveBatchCreate(
		workspaceId: string,
		batch: PendingChange[]
	): Promise<Map<string, string>> {
		// Сохраняем порядок элементов для сопоставления ID
		const frontendIds = batch.map((change) => change.elementId);

		// Преобразуем CanvasElement в CreateElementRequest для backend
		const elements = batch.map((change) => {
			const element = change.updates;

			// Специальная обработка для фигур
			const elementData: ElementData = {
				position: { x: element.pos_x, y: element.pos_y },
				size: { width: element.width || 0, height: element.height || 0 },
				style: element.style || {},
				rotation: element.rotation || 0
			};

			// Для shape типов добавляем shape_type
			const backendType = this.mapElementType(element.type || '');
			if (backendType === 'shape') {
				elementData.shape_type = element.type; // rectangle, ellipse, etc
			}

			// Добавляем content если есть
			if (element.content) {
				elementData.content = element.content;
			}
			if (element.html_content) {
				elementData.html_content = element.html_content;
			}

			return {
				element_type: backendType,
				element_data: elementData,
				z_index: element.z_index || 0,
				parent_id: element.parent_id
			};
		});

		const requestBody = {
			elements
		};

		const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
		const url = `${apiUrl}/workspaces/${workspaceId}/elements/batch`;
		const token = api.getAccessToken();

		console.log(`[Autosave] 📤 CREATE request to: ${url}`);
		console.log('[Autosave] Request body:', JSON.stringify(requestBody, null, 2));
		console.log(`[Autosave] Auth token present: ${!!token}`);

		// Call API
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(requestBody)
		});

		console.log(`[Autosave] Response status: ${response.status} ${response.statusText}`);

		if (!response.ok) {
			const errorText = await response.text().catch(() => 'Unknown error');
			console.error(`[Autosave] Response error body: ${errorText}`);
			throw new Error(`Save failed: ${response.status} ${errorText}`);
		}

		const responseData = await response.json().catch(() => null);
		console.log('[Autosave] Response data:', responseData);

		// Создаем маппинг frontend ID → backend ID
		const idMapping = new Map<string, string>();
		const typedResponse = responseData as BatchCreateResponse | null;
		if (typedResponse && typedResponse.elements && Array.isArray(typedResponse.elements)) {
			typedResponse.elements.forEach((element, index) => {
				if (index < frontendIds.length && element.id) {
					const frontendId = frontendIds[index];
					const backendId = element.id;
					idMapping.set(frontendId, backendId);
					console.log(`[Autosave] ID mapping: ${frontendId} → ${backendId}`);
				}
			});
		}

		return idMapping;
	}

	/**
	 * Save a single batch of EXISTING elements (UPDATE)
	 */
	/**
	 * Build partial update from change data (fallback when full element not available)
	 */
	private buildPartialUpdate(change: PendingChange): BatchUpdateItem {
		const update: BatchUpdateItem = {
			id: change.elementId
		};

		// Если есть изменения позиции или размера, создаем element_data
		const hasPositionOrSize =
			change.updates.pos_x !== undefined ||
			change.updates.pos_y !== undefined ||
			change.updates.width !== undefined ||
			change.updates.height !== undefined ||
			change.updates.rotation !== undefined ||
			change.updates.style !== undefined ||
			change.updates.content !== undefined ||
			change.updates.html_content !== undefined;

		if (hasPositionOrSize) {
			update.element_data = {};

			// Позиция
			if (change.updates.pos_x !== undefined || change.updates.pos_y !== undefined) {
				update.element_data.position = {
					x: change.updates.pos_x,
					y: change.updates.pos_y
				};
			}

			// Размер
			if (change.updates.width !== undefined || change.updates.height !== undefined) {
				update.element_data.size = {
					width: change.updates.width ?? 0,
					height: change.updates.height ?? 0
				};
			}

			// Стиль
			if (change.updates.style !== undefined) {
				update.element_data.style = change.updates.style;
			}

			// Rotation
			if (change.updates.rotation !== undefined) {
				update.element_data.rotation = change.updates.rotation;
			}

			// Content
			if (change.updates.content !== undefined) {
				update.element_data.content = change.updates.content;
			}

			// HTML content
			if (change.updates.html_content !== undefined) {
				update.element_data.html_content = change.updates.html_content;
			}

			// Shape type
			if (change.updates.type) {
				const shapeTypes = ['rectangle', 'ellipse', 'triangle', 'arrow', 'line'];
				if (shapeTypes.includes(change.updates.type)) {
					update.element_data.shape_type = change.updates.type;
				}
			}
		}

		// Z-index
		if (change.updates.z_index !== undefined) {
			update.z_index = change.updates.z_index;
		}

		// Parent ID
		if (change.updates.parent_id !== undefined) {
			update.parent_id = change.updates.parent_id;
		}

		return update;
	}

	private async saveBatchUpdate(workspaceId: string, batch: PendingChange[]): Promise<void> {
		// Преобразуем обновления в формат backend BatchUpdateItem
		const updates = batch.map((change) => {
			const update: BatchUpdateItem = {
				id: change.elementId
			};

			// Получаем полные данные элемента из canvas store
			const fullElement = this.getElementData(change.elementId);

			if (!fullElement) {
				console.warn(`[Autosave] ⚠️ Element ${change.elementId} not found in canvas store`);
				// Fallback to partial updates if element not found
				return this.buildPartialUpdate(change);
			}

			// Объединяем полные данные элемента с обновлениями
			const mergedElement = { ...fullElement, ...change.updates };

			// Создаем element_data со всеми данными
			update.element_data = {};

			// Позиция (всегда включаем)
			update.element_data.position = {
				x: mergedElement.pos_x ?? 0,
				y: mergedElement.pos_y ?? 0
			};

			// Размер (всегда включаем)
			update.element_data.size = {
				width: mergedElement.width ?? 0,
				height: mergedElement.height ?? 0
			};

			// Стиль (всегда включаем)
			if (mergedElement.style) {
				update.element_data.style = mergedElement.style;
			}

			// Rotation (всегда включаем)
			update.element_data.rotation = mergedElement.rotation ?? 0;

			// Content (для текстовых элементов)
			if (mergedElement.content !== undefined) {
				update.element_data.content = mergedElement.content;
			}

			// HTML content
			if (mergedElement.html_content !== undefined) {
				update.element_data.html_content = mergedElement.html_content;
			}

			// Shape type (если это shape)
			if (mergedElement.type) {
				const shapeTypes = ['rectangle', 'ellipse', 'triangle', 'arrow', 'line'];
				if (shapeTypes.includes(mergedElement.type)) {
					update.element_data.shape_type = mergedElement.type;
				}
			}

			// Z-index
			if (mergedElement.z_index !== undefined) {
				update.z_index = mergedElement.z_index;
			}

			// Parent ID
			if (mergedElement.parent_id !== undefined) {
				update.parent_id = mergedElement.parent_id;
			}

			return update;
		});

		// Use batch update API endpoint
		const requestBody = {
			updates
		};

		const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
		const url = `${apiUrl}/workspaces/${workspaceId}/elements/batch`;
		const token = api.getAccessToken();

		console.log(`[Autosave] 📤 UPDATE request to: ${url}`);
		console.log('[Autosave] Request body:', JSON.stringify(requestBody, null, 2));
		console.log(`[Autosave] Auth token present: ${!!token}`);

		// Call API
		const response = await fetch(url, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(requestBody)
		});

		console.log(`[Autosave] Response status: ${response.status} ${response.statusText}`);

		if (!response.ok) {
			const errorText = await response.text().catch(() => 'Unknown error');
			console.error(`[Autosave] Response error body: ${errorText}`);
			throw new Error(`Save failed: ${response.status} ${errorText}`);
		}

		const responseData = await response.json().catch(() => null);
		console.log('[Autosave] Response data:', responseData);
	}

	/**
	 * Update save status
	 */
	private setStatus(status: SaveStatus) {
		if (this.status !== status) {
			this.status = status;
			this.onStatusChange(status);
			this.log(`Status changed to: ${status}`);
		}
	}

	/**
	 * Get current save status
	 */
	getStatus(): SaveStatus {
		return this.status;
	}

	/**
	 * Get number of pending changes
	 */
	getPendingCount(): number {
		return this.pendingChanges.size;
	}

	/**
	 * Clear all pending changes without saving
	 */
	clear() {
		if (this.debounceTimer !== null) {
			clearTimeout(this.debounceTimer);
			this.debounceTimer = null;
		}
		this.pendingChanges.clear();
		this.setStatus('idle');
		this.log('Cleared all pending changes');
	}

	/**
	 * Pause autosave (useful during undo/redo operations)
	 */
	pause() {
		if (this.debounceTimer !== null) {
			clearTimeout(this.debounceTimer);
			this.debounceTimer = null;
		}
		this.log('Autosave paused');
	}

	/**
	 * Resume autosave
	 */
	resume() {
		if (this.pendingChanges.size > 0) {
			this.scheduleSave();
		}
		this.log('Autosave resumed');
	}

	/**
	 * Debug logging
	 */
	private log(message: string) {
		if (this.debug) {
			console.log(`[Autosave] ${message}`);
		}
	}
}

// Singleton instance
export const autosaveService = new AutosaveService({
	debug: true, // Always enable debug for now
	debounceMs: 1000,
	maxBatchSize: 50
});
