/**
 * Canvas with History - обертка для canvas store с поддержкой undo/redo
 *
 * Предоставляет те же методы что и canvasStore, но с автоматической записью в history
 */

import { canvasStore } from './canvas.svelte';
import {
	historyStore,
	createElementAction,
	deleteElementAction,
	updateElementAction,
	batchUpdateAction
} from './history.svelte';
import { autosaveStore } from './autosave.svelte';
import type { CanvasElement } from '$lib/types/api';

class CanvasWithHistory {
	// Проксируем все геттеры из canvasStore
	get viewport() {
		return canvasStore.viewport;
	}
	get elements() {
		return canvasStore.elements;
	}
	get selectedIds() {
		return canvasStore.selectedIds;
	}
	get selectedElements() {
		return canvasStore.selectedElements;
	}
	get activeTool() {
		return canvasStore.activeTool;
	}
	get isDragging() {
		return canvasStore.isDragging;
	}
	get isPanning() {
		return canvasStore.isPanning;
	}
	get isSelecting() {
		return canvasStore.isSelecting;
	}
	get selectionBox() {
		return canvasStore.selectionBox;
	}
	get showGrid() {
		return canvasStore.showGrid;
	}
	get snapToGrid() {
		return canvasStore.snapToGrid;
	}
	get gridSize() {
		return canvasStore.gridSize;
	}
	get workspaceId() {
		return canvasStore.workspaceId;
	}

	// Viewport методы (без history - они не persistent)
	setViewport = canvasStore.setViewport.bind(canvasStore);
	pan = canvasStore.pan.bind(canvasStore);
	zoom = canvasStore.zoom.bind(canvasStore);
	zoomTo = canvasStore.zoomTo.bind(canvasStore);
	resetZoom = canvasStore.resetZoom.bind(canvasStore);
	fitToScreen = canvasStore.fitToScreen.bind(canvasStore);
	screenToCanvas = canvasStore.screenToCanvas.bind(canvasStore);
	canvasToScreen = canvasStore.canvasToScreen.bind(canvasStore);
	snapPoint = canvasStore.snapPoint.bind(canvasStore);

	// Selection методы (без history - transient state)
	select = canvasStore.select.bind(canvasStore);
	selectMultiple = canvasStore.selectMultiple.bind(canvasStore);
	deselect = canvasStore.deselect.bind(canvasStore);
	clearSelection = canvasStore.clearSelection.bind(canvasStore);
	toggleSelection = canvasStore.toggleSelection.bind(canvasStore);
	selectAll = canvasStore.selectAll.bind(canvasStore);
	isSelected = canvasStore.isSelected.bind(canvasStore);

	// Selection box методы (без history)
	startSelectionBox = canvasStore.startSelectionBox.bind(canvasStore);
	updateSelectionBox = canvasStore.updateSelectionBox.bind(canvasStore);
	endSelectionBox = canvasStore.endSelectionBox.bind(canvasStore);
	cancelSelectionBox = canvasStore.cancelSelectionBox.bind(canvasStore);
	getSelectionBoxRect = canvasStore.getSelectionBoxRect.bind(canvasStore);

	// Tool методы (без history)
	setTool = canvasStore.setTool.bind(canvasStore);

	// State методы (без history)
	setIsDragging = canvasStore.setIsDragging.bind(canvasStore);
	setIsPanning = canvasStore.setIsPanning.bind(canvasStore);

	// Settings методы (без history - UI state)
	setShowGrid = canvasStore.setShowGrid.bind(canvasStore);
	setSnapToGrid = canvasStore.setSnapToGrid.bind(canvasStore);
	setGridSize = canvasStore.setGridSize.bind(canvasStore);
	toggleGrid = canvasStore.toggleGrid.bind(canvasStore);
	toggleSnap = canvasStore.toggleSnap.bind(canvasStore);

	// Utility методы (без history)
	getElement = canvasStore.getElement.bind(canvasStore);
	getElementsBounds = canvasStore.getElementsBounds.bind(canvasStore);
	getGroupElements = canvasStore.getGroupElements.bind(canvasStore);
	isGrouped = canvasStore.isGrouped.bind(canvasStore);
	setWorkspaceId = canvasStore.setWorkspaceId.bind(canvasStore);
	reset = canvasStore.reset.bind(canvasStore);

	// Методы БЕЗ history (используются для прямых операций из API)
	setElements = canvasStore.setElements.bind(canvasStore);
	addElements = canvasStore.addElements.bind(canvasStore);

	/**
	 * Добавить элемент с записью в history
	 */
	addElement(element: CanvasElement) {
		// Создаем копию для history
		const elementCopy = { ...element };

		// Записываем action
		const action = createElementAction(
			elementCopy,
			(data) => canvasStore.addElement(data),
			(id) => canvasStore.deleteElement(id)
		);

		historyStore.record(action);

		// Выполняем операцию
		canvasStore.addElement(element);

		// Автосохранение - отслеживаем создание элемента (isNew = true)
		autosaveStore.trackChange(element.id, element, true);
	}

	/**
	 * Обновить элемент с записью в history
	 */
	updateElement(id: string, updates: Partial<CanvasElement>) {
		const element = canvasStore.getElement(id);
		if (!element) return;

		// Создаем snapshot старых значений
		const oldData: Partial<CanvasElement> = {};
		const newData: Partial<CanvasElement> = {};

		for (const key in updates) {
			const k = key as keyof CanvasElement;
			oldData[k] = element[k];
			newData[k] = updates[k];
		}

		// Определяем тип операции для description
		let description = 'Update element';
		if ('pos_x' in updates || 'pos_y' in updates) {
			description = 'Move element';
		} else if ('width' in updates || 'height' in updates) {
			description = 'Resize element';
		} else if ('rotation' in updates) {
			description = 'Rotate element';
		} else if ('z_index' in updates) {
			description = 'Change layer';
		}

		// Записываем action
		const action = updateElementAction(
			id,
			oldData,
			newData,
			(elementId, data) => canvasStore.updateElement(elementId, data),
			description
		);

		historyStore.record(action);

		// Выполняем операцию
		canvasStore.updateElement(id, updates);

		// Автосохранение - отслеживаем изменение элемента
		autosaveStore.trackChange(id, updates);
	}

	/**
	 * Batch update с записью в history
	 */
	updateElements(updates: Array<{ id: string; updates: Partial<CanvasElement> }>) {
		// Создаем snapshots
		const snapshots = updates.map(({ id, updates: elementUpdates }) => {
			const element = canvasStore.getElement(id);
			if (!element) return null;

			const oldData: Partial<CanvasElement> = {};
			const newData: Partial<CanvasElement> = {};

			for (const key in elementUpdates) {
				const k = key as keyof CanvasElement;
				oldData[k] = element[k];
				newData[k] = elementUpdates[k];
			}

			return { id, oldData, newData };
		});

		const validSnapshots = snapshots.filter((s) => s !== null) as Array<{
			id: string;
			oldData: Partial<CanvasElement>;
			newData: Partial<CanvasElement>;
		}>;

		if (validSnapshots.length === 0) return;

		// Записываем batch action
		const action = batchUpdateAction(
			validSnapshots,
			(id, data) => canvasStore.updateElement(id, data),
			`Update ${validSnapshots.length} elements`
		);

		historyStore.record(action);

		// Выполняем операцию
		canvasStore.updateElements(updates);

		// Автосохранение - отслеживаем batch изменения
		autosaveStore.trackChanges(updates);
	}

	/**
	 * Удалить элемент с записью в history
	 */
	deleteElement(id: string) {
		const element = canvasStore.getElement(id);
		if (!element) return;

		// Создаем копию для history
		const elementCopy = { ...element };

		// Записываем action
		const action = deleteElementAction(
			elementCopy,
			(data) => canvasStore.addElement(data),
			(elementId) => canvasStore.deleteElement(elementId)
		);

		historyStore.record(action);

		// Выполняем операцию
		canvasStore.deleteElement(id);

		// Автосохранение - отслеживаем удаление (отправим deleted_at)
		autosaveStore.trackChange(id, {
			deleted_at: new Date().toISOString()
		} as Partial<CanvasElement>);
	}

	/**
	 * Удалить несколько элементов с записью в history
	 */
	deleteElements(ids: string[]) {
		if (ids.length === 0) return;

		// Начинаем batch
		historyStore.startBatch(`Delete ${ids.length} elements`);

		// Удаляем каждый элемент
		ids.forEach((id) => this.deleteElement(id));

		// Завершаем batch
		historyStore.endBatch();
	}

	/**
	 * Z-order методы с записью в history
	 */

	bringToFront(ids?: string[]) {
		const targetIds = ids || canvasStore.selectedIds;
		if (targetIds.length === 0) return;

		// Создаем snapshots текущих z-индексов
		const snapshots = targetIds.map((id) => {
			const element = canvasStore.getElement(id);
			return element ? { id, oldZ: element.z_index || 0 } : null;
		});

		const validSnapshots = snapshots.filter((s) => s !== null) as Array<{
			id: string;
			oldZ: number;
		}>;

		// Выполняем операцию
		canvasStore.bringToFront(targetIds);

		// Создаем snapshots новых z-индексов
		const updates = validSnapshots.map(({ id, oldZ }) => {
			const element = canvasStore.getElement(id);
			return {
				id,
				oldData: { z_index: oldZ },
				newData: { z_index: element?.z_index || 0 }
			};
		});

		// Записываем action
		const action = batchUpdateAction(
			updates,
			(id, data) => canvasStore.updateElement(id, data),
			'Bring to front'
		);

		historyStore.record(action);
	}

	sendToBack(ids?: string[]) {
		const targetIds = ids || canvasStore.selectedIds;
		if (targetIds.length === 0) return;

		// Создаем snapshots
		const snapshots = targetIds.map((id) => {
			const element = canvasStore.getElement(id);
			return element ? { id, oldZ: element.z_index || 0 } : null;
		});

		const validSnapshots = snapshots.filter((s) => s !== null) as Array<{
			id: string;
			oldZ: number;
		}>;

		// Выполняем операцию
		canvasStore.sendToBack(targetIds);

		// Создаем updates
		const updates = validSnapshots.map(({ id, oldZ }) => {
			const element = canvasStore.getElement(id);
			return {
				id,
				oldData: { z_index: oldZ },
				newData: { z_index: element?.z_index || 0 }
			};
		});

		// Записываем action
		const action = batchUpdateAction(
			updates,
			(id, data) => canvasStore.updateElement(id, data),
			'Send to back'
		);

		historyStore.record(action);
	}

	bringForward(ids?: string[]) {
		const targetIds = ids || canvasStore.selectedIds;
		if (targetIds.length === 0) return;

		const snapshots = targetIds.map((id) => {
			const element = canvasStore.getElement(id);
			return element ? { id, oldZ: element.z_index || 0 } : null;
		});

		const validSnapshots = snapshots.filter((s) => s !== null) as Array<{
			id: string;
			oldZ: number;
		}>;

		canvasStore.bringForward(targetIds);

		const updates = validSnapshots.map(({ id, oldZ }) => {
			const element = canvasStore.getElement(id);
			return {
				id,
				oldData: { z_index: oldZ },
				newData: { z_index: element?.z_index || 0 }
			};
		});

		const action = batchUpdateAction(
			updates,
			(id, data) => canvasStore.updateElement(id, data),
			'Bring forward'
		);

		historyStore.record(action);
	}

	sendBackward(ids?: string[]) {
		const targetIds = ids || canvasStore.selectedIds;
		if (targetIds.length === 0) return;

		const snapshots = targetIds.map((id) => {
			const element = canvasStore.getElement(id);
			return element ? { id, oldZ: element.z_index || 0 } : null;
		});

		const validSnapshots = snapshots.filter((s) => s !== null) as Array<{
			id: string;
			oldZ: number;
		}>;

		canvasStore.sendBackward(targetIds);

		const updates = validSnapshots.map(({ id, oldZ }) => {
			const element = canvasStore.getElement(id);
			return {
				id,
				oldData: { z_index: oldZ },
				newData: { z_index: element?.z_index || 0 }
			};
		});

		const action = batchUpdateAction(
			updates,
			(id, data) => canvasStore.updateElement(id, data),
			'Send backward'
		);

		historyStore.record(action);
	}

	/**
	 * Группировка с записью в history
	 */

	groupSelected() {
		if (canvasStore.selectedIds.length < 2) return;

		const groupId = crypto.randomUUID();
		const targetIds = [...canvasStore.selectedIds];

		// Создаем updates
		const updates = targetIds.map((id) => ({
			id,
			oldData: { group_id: undefined as string | undefined },
			newData: { group_id: groupId }
		}));

		// Записываем action
		const action = batchUpdateAction(
			updates,
			(id, data) => canvasStore.updateElement(id, data),
			'Group elements'
		);

		historyStore.record(action);

		// Выполняем операцию
		canvasStore.groupSelected();
	}

	ungroupSelected() {
		if (canvasStore.selectedIds.length === 0) return;

		const targetIds = [...canvasStore.selectedIds];

		// Создаем snapshots
		const updates = targetIds
			.map((id) => {
				const element = canvasStore.getElement(id);
				if (!element?.group_id) return null;

				return {
					id,
					oldData: { group_id: element.group_id },
					newData: { group_id: undefined as string | undefined }
				};
			})
			.filter((u) => u !== null) as Array<{
			id: string;
			oldData: Partial<CanvasElement>;
			newData: Partial<CanvasElement>;
		}>;

		if (updates.length === 0) return;

		// Записываем action
		const action = batchUpdateAction(
			updates,
			(id, data) => canvasStore.updateElement(id, data),
			'Ungroup elements'
		);

		historyStore.record(action);

		// Выполняем операцию
		canvasStore.ungroupSelected();
	}

	/**
	 * Undo/Redo
	 */

	undo() {
		historyStore.undo();
	}

	redo() {
		historyStore.redo();
	}

	get canUndo() {
		return historyStore.canUndo;
	}

	get canRedo() {
		return historyStore.canRedo;
	}

	/**
	 * Batch operations
	 */

	startBatch(description: string) {
		historyStore.startBatch(description);
	}

	endBatch() {
		historyStore.endBatch();
	}

	cancelBatch() {
		historyStore.cancelBatch();
	}

	/**
	 * Autosave operations
	 */

	/**
	 * Инициализировать автосохранение для workspace
	 */
	initAutosave(workspaceId: string) {
		autosaveStore.setWorkspaceId(workspaceId);

		// Установить callback для обновления ID элементов после создания
		autosaveStore.onElementIdChanged = (frontendId: string, backendId: string) => {
			console.log(`[CanvasWithHistory] Updating element ID: ${frontendId} → ${backendId}`);
			// Обновляем ID в canvas store
			const element = canvasStore.elements.find((el) => el.id === frontendId);
			if (element) {
				// Удаляем элемент со старым ID и добавляем с новым
				canvasStore.deleteElement(frontendId);
				canvasStore.addElement({ ...element, id: backendId });
				console.log(`[CanvasWithHistory] ✅ Updated element ID in canvas store`);
			} else {
				console.warn(`[CanvasWithHistory] ⚠️ Element ${frontendId} not found in canvas store`);
			}
		};

		// Установить callback для получения полных данных элемента
		autosaveStore.getElementData = (elementId: string) => {
			const element = canvasStore.elements.find((el) => el.id === elementId);
			return element;
		};
	}

	/**
	 * Остановить автосохранение
	 */
	stopAutosave() {
		autosaveStore.clear();
		autosaveStore.setWorkspaceId(null);
	}

	/**
	 * Принудительно сохранить все изменения
	 */
	async saveNow() {
		await autosaveStore.flush();
	}

	/**
	 * Получить статус автосохранения
	 */
	get autosaveStatus() {
		return autosaveStore.status;
	}

	get autosavePendingCount() {
		return autosaveStore.pendingCount;
	}
}

// Export singleton
export const canvas = new CanvasWithHistory();
