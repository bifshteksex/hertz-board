/**
 * Canvas Store - управление состоянием Canvas
 *
 * Ответственность:
 * - Viewport (позиция, zoom)
 * - Элементы canvas
 * - Выделение
 * - Активный инструмент
 * - Состояние взаимодействия (drag, pan)
 */

import type { CanvasElement } from '$lib/types/api';

// Типы для Canvas

export type Tool =
	| 'select'
	| 'text'
	| 'rectangle'
	| 'ellipse'
	| 'triangle'
	| 'line'
	| 'arrow'
	| 'freehand'
	| 'sticky'
	| 'list'
	| 'image'
	| 'connector';

export interface Viewport {
	x: number;
	y: number;
	zoom: number; // 0.1 (10%) to 5.0 (500%)
}

export interface Point {
	x: number;
	y: number;
}

export interface Rect {
	x: number;
	y: number;
	width: number;
	height: number;
}

export interface SelectionBox {
	start: Point;
	current: Point;
}

// Константы
export const MIN_ZOOM = 0.1;
export const MAX_ZOOM = 5.0;
export const ZOOM_STEP = 0.1;
export const GRID_SIZE = 20; // базовый размер сетки в пикселях

class CanvasStore {
	// Viewport
	private _viewport = $state<Viewport>({ x: 0, y: 0, zoom: 1 });

	// Элементы (ключ - element.id)
	private _elements = $state<Map<string, CanvasElement>>(new Map());

	// Выделение
	private _selectedIds = $state<Set<string>>(new Set());

	// Активный инструмент
	private _activeTool = $state<Tool>('select');

	// Состояние взаимодействия
	private _isDragging = $state(false);
	private _isPanning = $state(false);
	private _isSelecting = $state(false);
	private _selectionBox = $state<SelectionBox | null>(null);

	// Настройки
	private _showGrid = $state(true);
	private _snapToGrid = $state(false);
	private _gridSize = $state(GRID_SIZE);

	// Workspace ID (текущий открытый workspace)
	private _workspaceId = $state<string | null>(null);

	// Getters (derived state)
	get viewport() {
		return this._viewport;
	}

	get elements() {
		return Array.from(this._elements.values());
	}

	get selectedIds() {
		return Array.from(this._selectedIds);
	}

	get selectedElements() {
		return this.elements.filter((el) => this._selectedIds.has(el.id));
	}

	get activeTool() {
		return this._activeTool;
	}

	get isDragging() {
		return this._isDragging;
	}

	get isPanning() {
		return this._isPanning;
	}

	get isSelecting() {
		return this._isSelecting;
	}

	get selectionBox() {
		return this._selectionBox;
	}

	get showGrid() {
		return this._showGrid;
	}

	get snapToGrid() {
		return this._snapToGrid;
	}

	get gridSize() {
		return this._gridSize;
	}

	get workspaceId() {
		return this._workspaceId;
	}

	// Viewport управление

	/**
	 * Установить viewport
	 */
	setViewport(viewport: Partial<Viewport>) {
		this._viewport = {
			...this._viewport,
			...viewport,
			zoom: Math.max(MIN_ZOOM, Math.min(MAX_ZOOM, viewport.zoom ?? this._viewport.zoom))
		};
	}

	/**
	 * Pan viewport (переместить)
	 */
	pan(dx: number, dy: number) {
		this._viewport.x += dx;
		this._viewport.y += dy;
	}

	/**
	 * Zoom viewport
	 * @param delta - изменение zoom (положительное = zoom in, отрицательное = zoom out)
	 * @param centerX - X координата центра zoom (в screen coordinates)
	 * @param centerY - Y координата центра zoom (в screen coordinates)
	 */
	zoom(delta: number, centerX?: number, centerY?: number) {
		const oldZoom = this._viewport.zoom;
		const newZoom = Math.max(MIN_ZOOM, Math.min(MAX_ZOOM, oldZoom + delta));

		if (newZoom === oldZoom) return;

		// Если указана точка центра, zoom к этой точке
		if (centerX !== undefined && centerY !== undefined) {
			// Преобразуем screen координаты в canvas координаты при старом zoom
			const canvasX = (centerX - this._viewport.x) / oldZoom;
			const canvasY = (centerY - this._viewport.y) / oldZoom;

			// После zoom нужно скорректировать viewport так,
			// чтобы точка canvas осталась под курсором
			this._viewport.x = centerX - canvasX * newZoom;
			this._viewport.y = centerY - canvasY * newZoom;
		}

		this._viewport.zoom = newZoom;
	}

	/**
	 * Zoom к указанному уровню
	 */
	zoomTo(zoom: number, centerX?: number, centerY?: number) {
		const delta = zoom - this._viewport.zoom;
		this.zoom(delta, centerX, centerY);
	}

	/**
	 * Сбросить zoom к 100%
	 */
	resetZoom() {
		this._viewport.zoom = 1;
	}

	/**
	 * Fit to screen - показать все элементы
	 */
	fitToScreen(containerWidth: number, containerHeight: number, padding = 50) {
		if (this._elements.size === 0) {
			// Если нет элементов, центрируем viewport
			this._viewport = { x: containerWidth / 2, y: containerHeight / 2, zoom: 1 };
			return;
		}

		// Находим bounding box всех элементов
		const bounds = this.getElementsBounds();
		if (!bounds) return;

		// Вычисляем zoom чтобы все элементы поместились
		const zoomX = (containerWidth - padding * 2) / bounds.width;
		const zoomY = (containerHeight - padding * 2) / bounds.height;
		const zoom = Math.min(zoomX, zoomY, MAX_ZOOM);

		// Центрируем
		const centerX = bounds.x + bounds.width / 2;
		const centerY = bounds.y + bounds.height / 2;

		this._viewport = {
			x: containerWidth / 2 - centerX * zoom,
			y: containerHeight / 2 - centerY * zoom,
			zoom: Math.max(MIN_ZOOM, zoom)
		};
	}

	// Координатные преобразования

	/**
	 * Преобразовать screen координаты в canvas координаты
	 */
	screenToCanvas(screenX: number, screenY: number): Point {
		return {
			x: (screenX - this._viewport.x) / this._viewport.zoom,
			y: (screenY - this._viewport.y) / this._viewport.zoom
		};
	}

	/**
	 * Преобразовать canvas координаты в screen координаты
	 */
	canvasToScreen(canvasX: number, canvasY: number): Point {
		return {
			x: canvasX * this._viewport.zoom + this._viewport.x,
			y: canvasY * this._viewport.zoom + this._viewport.y
		};
	}

	/**
	 * Snap точку к сетке (если включено)
	 */
	snapPoint(point: Point): Point {
		if (!this._snapToGrid) return point;

		const grid = this._gridSize;
		return {
			x: Math.round(point.x / grid) * grid,
			y: Math.round(point.y / grid) * grid
		};
	}

	// Управление элементами

	/**
	 * Установить элементы (обычно при загрузке workspace)
	 */
	setElements(elements: CanvasElement[]) {
		this._elements.clear();
		// Сортируем по z_index при загрузке
		const sorted = [...elements].sort((a, b) => (a.z_index || 0) - (b.z_index || 0));
		sorted.forEach((el) => this._elements.set(el.id, el));
	}

	/**
	 * Добавить элемент
	 */
	addElement(element: CanvasElement) {
		this._elements.set(element.id, element);
		// Trigger reactivity in Svelte 5
		this._elements = new Map(this._elements);
	}

	/**
	 * Добавить несколько элементов
	 */
	addElements(elements: CanvasElement[]) {
		elements.forEach((el) => this._elements.set(el.id, el));
		// Trigger reactivity in Svelte 5
		this._elements = new Map(this._elements);
	}

	/**
	 * Обновить элемент
	 */
	updateElement(id: string, updates: Partial<CanvasElement>) {
		const element = this._elements.get(id);
		if (!element) return;

		this._elements.set(id, { ...element, ...updates });
		// Trigger reactivity in Svelte 5
		this._elements = new Map(this._elements);
	}

	/**
	 * Batch update нескольких элементов (более эффективно)
	 */
	updateElements(updates: Array<{ id: string; updates: Partial<CanvasElement> }>) {
		updates.forEach(({ id, updates: elementUpdates }) => {
			const element = this._elements.get(id);
			if (element) {
				this._elements.set(id, { ...element, ...elementUpdates });
			}
		});
		// Trigger reactivity once for all updates
		this._elements = new Map(this._elements);
	}

	/**
	 * Удалить элемент
	 */
	deleteElement(id: string) {
		this._elements.delete(id);
		this._selectedIds.delete(id);
		// Trigger reactivity in Svelte 5
		this._elements = new Map(this._elements);
		this._selectedIds = new Set(this._selectedIds);
	}

	/**
	 * Удалить несколько элементов
	 */
	deleteElements(ids: string[]) {
		ids.forEach((id) => {
			this._elements.delete(id);
			this._selectedIds.delete(id);
		});
		// Trigger reactivity in Svelte 5
		this._elements = new Map(this._elements);
		this._selectedIds = new Set(this._selectedIds);
	}

	/**
	 * Получить элемент по ID
	 */
	getElement(id: string): CanvasElement | undefined {
		return this._elements.get(id);
	}

	/**
	 * Получить bounding box всех элементов или выбранных
	 */
	getElementsBounds(elementIds?: string[]): Rect | null {
		const elements = elementIds
			? (elementIds.map((id) => this._elements.get(id)).filter(Boolean) as CanvasElement[])
			: this.elements;

		if (elements.length === 0) return null;

		let minX = Infinity;
		let minY = Infinity;
		let maxX = -Infinity;
		let maxY = -Infinity;

		elements.forEach((el) => {
			minX = Math.min(minX, el.pos_x);
			minY = Math.min(minY, el.pos_y);
			maxX = Math.max(maxX, el.pos_x + (el.width || 0));
			maxY = Math.max(maxY, el.pos_y + (el.height || 0));
		});

		return {
			x: minX,
			y: minY,
			width: maxX - minX,
			height: maxY - minY
		};
	}

	// Управление выделением

	/**
	 * Выделить элемент
	 */
	select(id: string, addToSelection = false) {
		if (!addToSelection) {
			this._selectedIds.clear();
		}
		this._selectedIds.add(id);
		// Trigger reactivity in Svelte 5
		this._selectedIds = new Set(this._selectedIds);
	}

	/**
	 * Выделить несколько элементов
	 */
	selectMultiple(ids: string[], addToSelection = false) {
		if (!addToSelection) {
			this._selectedIds.clear();
		}
		ids.forEach((id) => this._selectedIds.add(id));
		// Trigger reactivity in Svelte 5
		this._selectedIds = new Set(this._selectedIds);
	}

	/**
	 * Убрать выделение с элемента
	 */
	deselect(id: string) {
		this._selectedIds.delete(id);
		// Trigger reactivity in Svelte 5
		this._selectedIds = new Set(this._selectedIds);
	}

	/**
	 * Очистить выделение
	 */
	clearSelection() {
		this._selectedIds.clear();
		// Trigger reactivity in Svelte 5
		this._selectedIds = new Set(this._selectedIds);
	}

	/**
	 * Переключить выделение элемента
	 */
	toggleSelection(id: string) {
		if (this._selectedIds.has(id)) {
			this._selectedIds.delete(id);
		} else {
			this._selectedIds.add(id);
		}
		// Trigger reactivity in Svelte 5
		this._selectedIds = new Set(this._selectedIds);
	}

	/**
	 * Выделить все элементы
	 */
	selectAll() {
		this._elements.forEach((_, id) => this._selectedIds.add(id));
		// Trigger reactivity in Svelte 5
		this._selectedIds = new Set(this._selectedIds);
	}

	/**
	 * Проверить, выделен ли элемент
	 */
	isSelected(id: string): boolean {
		return this._selectedIds.has(id);
	}

	// Selection Box (прямоугольное выделение)

	/**
	 * Начать selection box
	 */
	startSelectionBox(canvasX: number, canvasY: number) {
		const point = { x: canvasX, y: canvasY };
		this._selectionBox = { start: point, current: point };
		this._isSelecting = true;
	}

	/**
	 * Обновить selection box
	 */
	updateSelectionBox(canvasX: number, canvasY: number) {
		if (!this._selectionBox) return;

		const point = { x: canvasX, y: canvasY };
		this._selectionBox = { ...this._selectionBox, current: point };
	}

	/**
	 * Завершить selection box и выделить элементы
	 */
	endSelectionBox(addToSelection = false) {
		if (!this._selectionBox) return;

		const box = this.getSelectionBoxRect();
		if (!box) {
			this._selectionBox = null;
			this._isSelecting = false;
			return;
		}

		// Найти элементы которые пересекаются с selection box
		const intersecting = this.elements.filter((el) =>
			this.rectIntersects(box, {
				x: el.pos_x,
				y: el.pos_y,
				width: el.width || 0,
				height: el.height || 0
			})
		);

		if (!addToSelection) {
			this._selectedIds.clear();
		}

		intersecting.forEach((el) => this._selectedIds.add(el.id));

		this._selectionBox = null;
		this._isSelecting = false;
	}

	/**
	 * Отменить selection box
	 */
	cancelSelectionBox() {
		this._selectionBox = null;
		this._isSelecting = false;
	}

	/**
	 * Получить Rect для selection box
	 */
	getSelectionBoxRect(): Rect | null {
		if (!this._selectionBox) return null;

		const { start, current } = this._selectionBox;
		const x = Math.min(start.x, current.x);
		const y = Math.min(start.y, current.y);
		const width = Math.abs(current.x - start.x);
		const height = Math.abs(current.y - start.y);

		return { x, y, width, height };
	}

	/**
	 * Проверить пересечение двух прямоугольников
	 */
	private rectIntersects(a: Rect, b: Rect): boolean {
		return !(
			a.x + a.width < b.x ||
			b.x + b.width < a.x ||
			a.y + a.height < b.y ||
			b.y + b.height < a.y
		);
	}

	// Управление инструментами

	/**
	 * Установить активный инструмент
	 */
	setTool(tool: Tool) {
		this._activeTool = tool;
		// При смене инструмента сбрасываем выделение
		if (tool !== 'select') {
			this.clearSelection();
		}
	}

	// Управление состоянием взаимодействия

	setIsDragging(value: boolean) {
		this._isDragging = value;
	}

	setIsPanning(value: boolean) {
		this._isPanning = value;
	}

	// Настройки

	setShowGrid(value: boolean) {
		this._showGrid = value;
	}

	setSnapToGrid(value: boolean) {
		this._snapToGrid = value;
	}

	setGridSize(size: number) {
		this._gridSize = size;
	}

	toggleGrid() {
		this._showGrid = !this._showGrid;
	}

	toggleSnap() {
		this._snapToGrid = !this._snapToGrid;
	}

	// Workspace management

	setWorkspaceId(id: string | null) {
		this._workspaceId = id;
	}

	// Z-order управление

	/**
	 * Переместить элемент(ы) на передний план
	 */
	bringToFront(ids?: string[]) {
		const targetIds = ids || this.selectedIds;
		if (targetIds.length === 0) return;

		const maxZ = Math.max(...this.elements.map((el) => el.z_index || 0), 0);

		const updates = targetIds.map((id, index) => ({
			id,
			updates: { z_index: maxZ + index + 1 }
		}));

		this.updateElements(updates);
	}

	/**
	 * Переместить элемент(ы) на задний план
	 */
	sendToBack(ids?: string[]) {
		const targetIds = ids || this.selectedIds;
		if (targetIds.length === 0) return;

		const minZ = Math.min(...this.elements.map((el) => el.z_index || 0), 0);

		const updates = targetIds.map((id, index) => ({
			id,
			updates: { z_index: minZ - targetIds.length + index }
		}));

		this.updateElements(updates);
	}

	/**
	 * Переместить элемент(ы) на один слой вперед
	 */
	bringForward(ids?: string[]) {
		const targetIds = ids || this.selectedIds;
		if (targetIds.length === 0) return;

		const updates = targetIds
			.map((id) => {
				const element = this.getElement(id);
				if (!element) return null;
				return {
					id,
					updates: { z_index: (element.z_index || 0) + 1 }
				};
			})
			.filter((u) => u !== null);

		this.updateElements(updates as Array<{ id: string; updates: Partial<CanvasElement> }>);
	}

	/**
	 * Переместить элемент(ы) на один слой назад
	 */
	sendBackward(ids?: string[]) {
		const targetIds = ids || this.selectedIds;
		if (targetIds.length === 0) return;

		const updates = targetIds
			.map((id) => {
				const element = this.getElement(id);
				if (!element) return null;
				return {
					id,
					updates: { z_index: (element.z_index || 0) - 1 }
				};
			})
			.filter((u) => u !== null);

		this.updateElements(updates as Array<{ id: string; updates: Partial<CanvasElement> }>);
	}

	// Группировка

	/**
	 * Сгруппировать выделенные элементы
	 */
	groupSelected() {
		if (this.selectedIds.length < 2) return;

		const groupId = crypto.randomUUID();

		this.selectedIds.forEach((id) => {
			this.updateElement(id, { group_id: groupId });
		});

		return groupId;
	}

	/**
	 * Разгруппировать выделенные элементы
	 */
	ungroupSelected() {
		if (this.selectedIds.length === 0) return;

		this.selectedIds.forEach((id) => {
			const element = this.getElement(id);
			if (element?.group_id) {
				this.updateElement(id, { group_id: undefined });
			}
		});
	}

	/**
	 * Получить все элементы группы
	 */
	getGroupElements(groupId: string): CanvasElement[] {
		return this.elements.filter((el) => el.group_id === groupId);
	}

	/**
	 * Проверить, является ли элемент частью группы
	 */
	isGrouped(id: string): boolean {
		const element = this.getElement(id);
		return !!element?.group_id;
	}

	/**
	 * Очистить все состояние (при выходе из workspace)
	 */
	reset() {
		this._viewport = { x: 0, y: 0, zoom: 1 };
		this._elements.clear();
		this._selectedIds.clear();
		this._activeTool = 'select';
		this._isDragging = false;
		this._isPanning = false;
		this._isSelecting = false;
		this._selectionBox = null;
		this._workspaceId = null;
	}
}

// Singleton instance
export const canvasStore = new CanvasStore();
