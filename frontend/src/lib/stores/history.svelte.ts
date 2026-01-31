/**
 * History Store - Undo/Redo система
 *
 * Ответственность:
 * - Управление стеком undo/redo операций
 * - Запись операций над элементами
 * - Откат и повтор операций
 */

import type { CanvasElement } from '$lib/types/api';

/**
 * Типы операций, поддерживающих undo/redo
 */
export type HistoryActionType =
	| 'create'
	| 'delete'
	| 'update'
	| 'move'
	| 'resize'
	| 'rotate'
	| 'style'
	| 'group'
	| 'ungroup'
	| 'z-order'
	| 'batch'; // группа операций

/**
 * Действие в истории
 */
export interface HistoryAction {
	type: HistoryActionType;
	timestamp: number;
	description: string;

	// Функции для выполнения undo/redo
	undo: () => void;
	redo: () => void;
}

/**
 * Batch action - группа связанных операций
 */
export interface BatchAction extends HistoryAction {
	type: 'batch';
	actions: HistoryAction[];
}

/**
 * История изменений элемента (для создания snapshot)
 */
export interface ElementSnapshot {
	id: string;
	data: Partial<CanvasElement>;
}

const MAX_HISTORY_SIZE = 100; // максимальное количество операций в истории

class HistoryStore {
	// Стеки undo/redo
	private _undoStack = $state<HistoryAction[]>([]);
	private _redoStack = $state<HistoryAction[]>([]);

	// Флаг для отключения записи истории (например, во время undo/redo)
	private _isRecording = $state(true);

	// Batch mode - группировка операций
	private _isBatching = $state(false);
	private _batchActions: HistoryAction[] = [];
	private _batchDescription = '';

	// Getters
	get canUndo() {
		return this._undoStack.length > 0;
	}

	get canRedo() {
		return this._redoStack.length > 0;
	}

	get undoStack() {
		return this._undoStack;
	}

	get redoStack() {
		return this._redoStack;
	}

	get lastAction() {
		return this._undoStack[this._undoStack.length - 1];
	}

	/**
	 * Записать действие в историю
	 */
	record(action: HistoryAction) {
		if (!this._isRecording) return;

		// Если в batch режиме, добавляем в batch
		if (this._isBatching) {
			this._batchActions.push(action);
			return;
		}

		// Добавляем в стек undo
		this._undoStack.push(action);

		// Ограничиваем размер стека
		if (this._undoStack.length > MAX_HISTORY_SIZE) {
			this._undoStack.shift();
		}

		// Очищаем стек redo при новой операции
		this._redoStack = [];
	}

	/**
	 * Отменить последнее действие
	 */
	undo() {
		const action = this._undoStack.pop();
		if (!action) return;

		// Отключаем запись истории во время undo
		this._isRecording = false;

		try {
			action.undo();
			this._redoStack.push(action);
		} finally {
			this._isRecording = true;
		}
	}

	/**
	 * Повторить отмененное действие
	 */
	redo() {
		const action = this._redoStack.pop();
		if (!action) return;

		// Отключаем запись истории во время redo
		this._isRecording = false;

		try {
			action.redo();
			this._undoStack.push(action);
		} finally {
			this._isRecording = true;
		}
	}

	/**
	 * Начать batch операцию (группировка связанных действий)
	 */
	startBatch(description: string) {
		this._isBatching = true;
		this._batchActions = [];
		this._batchDescription = description;
	}

	/**
	 * Завершить batch операцию
	 */
	endBatch() {
		if (!this._isBatching) return;

		this._isBatching = false;

		// Если есть действия в batch, создаем batch action
		if (this._batchActions.length > 0) {
			const batchAction: BatchAction = {
				type: 'batch',
				timestamp: Date.now(),
				description: this._batchDescription,
				actions: [...this._batchActions],
				undo: () => {
					// Откатываем все действия в обратном порядке
					for (let i = batchAction.actions.length - 1; i >= 0; i--) {
						batchAction.actions[i].undo();
					}
				},
				redo: () => {
					// Повторяем все действия в прямом порядке
					for (const action of batchAction.actions) {
						action.redo();
					}
				}
			};

			// Записываем batch как одно действие
			this._isRecording = true;
			this.record(batchAction);
		}

		this._batchActions = [];
		this._batchDescription = '';
	}

	/**
	 * Отменить batch операцию
	 */
	cancelBatch() {
		this._isBatching = false;
		this._batchActions = [];
		this._batchDescription = '';
	}

	/**
	 * Очистить историю
	 */
	clear() {
		this._undoStack = [];
		this._redoStack = [];
		this._batchActions = [];
		this._isBatching = false;
		this._isRecording = true;
	}

	/**
	 * Включить/выключить запись истории
	 */
	setRecording(enabled: boolean) {
		this._isRecording = enabled;
	}
}

// Singleton instance
export const historyStore = new HistoryStore();

/**
 * Helper функции для создания actions
 */

/**
 * Создать action для создания элемента
 */
export function createElementAction(
	elementData: CanvasElement,
	onCreate: (data: CanvasElement) => void,
	onDelete: (id: string) => void
): HistoryAction {
	return {
		type: 'create',
		timestamp: Date.now(),
		description: `Create ${elementData.type}`,
		undo: () => onDelete(elementData.id),
		redo: () => onCreate(elementData)
	};
}

/**
 * Создать action для удаления элемента
 */
export function deleteElementAction(
	elementData: CanvasElement,
	onCreate: (data: CanvasElement) => void,
	onDelete: (id: string) => void
): HistoryAction {
	return {
		type: 'delete',
		timestamp: Date.now(),
		description: `Delete ${elementData.type}`,
		undo: () => onCreate(elementData),
		redo: () => onDelete(elementData.id)
	};
}

/**
 * Создать action для обновления элемента
 */
export function updateElementAction(
	elementId: string,
	oldData: Partial<CanvasElement>,
	newData: Partial<CanvasElement>,
	onUpdate: (id: string, data: Partial<CanvasElement>) => void,
	description?: string
): HistoryAction {
	return {
		type: 'update',
		timestamp: Date.now(),
		description: description || 'Update element',
		undo: () => onUpdate(elementId, oldData),
		redo: () => onUpdate(elementId, newData)
	};
}

/**
 * Создать action для batch обновления элементов
 */
export function batchUpdateAction(
	updates: Array<{
		id: string;
		oldData: Partial<CanvasElement>;
		newData: Partial<CanvasElement>;
	}>,
	onUpdate: (id: string, data: Partial<CanvasElement>) => void,
	description: string
): HistoryAction {
	return {
		type: 'batch',
		timestamp: Date.now(),
		description,
		undo: () => {
			updates.forEach(({ id, oldData }) => onUpdate(id, oldData));
		},
		redo: () => {
			updates.forEach(({ id, newData }) => onUpdate(id, newData));
		}
	};
}
