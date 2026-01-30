<script lang="ts">
	import { onMount } from 'svelte';
	import { canvasStore } from '$lib/stores/canvas.svelte';
	import type { CanvasElement as CanvasElementType } from '$lib/types/api';
	import CanvasGrid from './CanvasGrid.svelte';
	import CanvasElement from './CanvasElement.svelte';
	import SelectionBox from './SelectionBox.svelte';
	import SelectionHandles from './SelectionHandles.svelte';
	import ZoomControls from './ZoomControls.svelte';
	import ShapeCreator from './ShapeCreator.svelte';
	import FreehandDrawing from './FreehandDrawing.svelte';
	import ImageUploader from './ImageUploader.svelte';
	import ConnectorCreator from './ConnectorCreator.svelte';

	let canvasContainer: HTMLDivElement;
	let svgCanvas: SVGSVGElement;

	// Mouse state
	let isPanningActive = $state(false);
	let panStart = $state({ x: 0, y: 0 });
	let isSpacePressed = $state(false);

	// Element manipulation state
	let isDraggingElements = $state(false);
	let dragStart = $state<{ x: number; y: number } | null>(null);
	let elementStartPositions = $state<Map<string, { x: number; y: number }>>(new Map());

	let isResizing = $state(false);
	let resizeHandle = $state<string | null>(null);
	let resizeStartData = $state<{
		mouseX: number;
		mouseY: number;
		elements: Map<string, { x: number; y: number; width: number; height: number }>;
	} | null>(null);

	// Shape creation state
	let isCreatingShape = $state(false);
	let shapeCreateStart = $state<{ x: number; y: number } | null>(null);
	let shapeCreateCurrent = $state<{ x: number; y: number } | null>(null);
	let isShiftPressed = $state(false);

	// Freehand drawing state
	let isDrawing = $state(false);
	let drawingPoints = $state<{ x: number; y: number; pressure?: number }[]>([]);

	// Connector creation state
	let isCreatingConnector = $state(false);
	let connectorStart = $state<{ x: number; y: number } | null>(null);
	let connectorCurrent = $state<{ x: number; y: number } | null>(null);
	let connectorType = $state<'straight' | 'curved' | 'elbow'>('straight');

	// Image uploader
	let imageUploader = $state<ImageUploader>();

	// Track if any element is being edited
	let isEditingText = $state(false);

	onMount(() => {
		// Fit to screen при загрузке
		if (canvasContainer) {
			const rect = canvasContainer.getBoundingClientRect();
			canvasStore.fitToScreen(rect.width, rect.height);
		}

		// Keyboard shortcuts
		const handleKeyDown = (e: KeyboardEvent) => {
			// Skip keyboard shortcuts when editing text
			if (isEditingText) return;

			// Shift для сохранения пропорций
			if (e.code === 'ShiftLeft' || e.code === 'ShiftRight') {
				isShiftPressed = true;
			}

			// Space для panning
			if (e.code === 'Space' && !isSpacePressed && !isDraggingElements && !isResizing) {
				e.preventDefault();
				isSpacePressed = true;
				if (canvasContainer) {
					canvasContainer.style.cursor = 'grab';
				}
			}

			// Delete для удаления выделенных элементов
			if ((e.code === 'Delete' || e.code === 'Backspace') && canvasStore.selectedIds.length > 0) {
				e.preventDefault();
				canvasStore.deleteElements(canvasStore.selectedIds);
			}

			// Ctrl/Cmd + A для выделения всех
			if ((e.ctrlKey || e.metaKey) && e.code === 'KeyA') {
				e.preventDefault();
				canvasStore.selectAll();
			}

			// Escape для отмены выделения
			if (e.code === 'Escape') {
				canvasStore.clearSelection();
			}

			// Ctrl/Cmd + 0 для reset zoom
			if ((e.ctrlKey || e.metaKey) && e.code === 'Digit0') {
				e.preventDefault();
				canvasStore.resetZoom();
			}
		};

		const handleKeyUp = (e: KeyboardEvent) => {
			if (e.code === 'Space') {
				isSpacePressed = false;
				if (canvasContainer && !isPanningActive) {
					canvasContainer.style.cursor = 'default';
				}
			}
			if (e.code === 'ShiftLeft' || e.code === 'ShiftRight') {
				isShiftPressed = false;
			}
		};

		window.addEventListener('keydown', handleKeyDown);
		window.addEventListener('keyup', handleKeyUp);

		return () => {
			window.removeEventListener('keydown', handleKeyDown);
			window.removeEventListener('keyup', handleKeyUp);

			// Cleanup RAF
			if (rafId !== null) {
				cancelAnimationFrame(rafId);
				rafId = null;
			}
		};
	});

	// Helper functions for creating elements

	function handleImageUploaded(imageUrl: string, width: number, height: number) {
		if (!canvasStore.workspaceId) return;

		// Create image element at center of viewport
		const centerX = canvasStore.viewport.x + window.innerWidth / (2 * canvasStore.viewport.zoom);
		const centerY = canvasStore.viewport.y + window.innerHeight / (2 * canvasStore.viewport.zoom);

		// Scale down large images
		const maxSize = 400;
		let finalWidth = width;
		let finalHeight = height;
		if (width > maxSize || height > maxSize) {
			const scale = Math.min(maxSize / width, maxSize / height);
			finalWidth = width * scale;
			finalHeight = height * scale;
		}

		const newElement: CanvasElementType = {
			id: crypto.randomUUID(),
			workspace_id: canvasStore.workspaceId,
			type: 'image',
			pos_x: centerX - finalWidth / 2,
			pos_y: centerY - finalHeight / 2,
			width: finalWidth,
			height: finalHeight,
			z_index: canvasStore.elements.length,
			image_url: imageUrl
		};

		canvasStore.addElement(newElement);
	}

	function createShapeElement(
		type: 'rectangle' | 'ellipse' | 'triangle' | 'line' | 'arrow',
		startX: number,
		startY: number,
		endX: number,
		endY: number,
		maintainAspect: boolean = false
	) {
		if (!canvasStore.workspaceId) return;

		const x = Math.min(startX, endX);
		const y = Math.min(startY, endY);
		let width = Math.abs(endX - startX);
		let height = Math.abs(endY - startY);

		// Maintain aspect ratio if shift pressed
		if (maintainAspect) {
			const size = Math.max(width, height);
			width = size;
			height = size;
		}

		// Minimum size
		if (width < 10 || height < 10) return;

		const newElement: CanvasElementType = {
			id: crypto.randomUUID(),
			workspace_id: canvasStore.workspaceId,
			type: type,
			pos_x: type === 'line' || type === 'arrow' ? startX : x,
			pos_y: type === 'line' || type === 'arrow' ? startY : y,
			width: type === 'line' || type === 'arrow' ? endX - startX : width,
			height: type === 'line' || type === 'arrow' ? endY - startY : height,
			z_index: canvasStore.elements.length,
			style: {
				backgroundColor: type === 'line' || type === 'arrow' ? undefined : '#3b82f6',
				strokeColor: '#1e40af',
				strokeWidth: 2,
				opacity: 1
			}
		};

		canvasStore.addElement(newElement);
		canvasStore.select(newElement.id);
	}

	async function createFreehandElement(points: { x: number; y: number; pressure?: number }[]) {
		if (!canvasStore.workspaceId) return;
		if (points.length < 2) return;

		// Convert points to SVG path using perfect-freehand
		const { getStroke } = await import('perfect-freehand');
		const stroke = getStroke(points, {
			size: 2 * 2,
			thinning: 0.5,
			smoothing: 0.5,
			streamline: 0.5,
			simulatePressure: true
		});

		if (stroke.length === 0) return;

		const pathData = stroke.reduce((acc, [x0, y0], i, arr) => {
			const [x1, y1] = arr[(i + 1) % arr.length];
			if (i === 0) return `M ${x0},${y0} Q ${x1},${y1}`;
			if (i === arr.length - 1) return `${acc} ${x0},${y0} Z`;
			return `${acc} ${x0},${y0} T ${x1},${y1}`;
		}, '');

		const newElement: CanvasElementType = {
			id: crypto.randomUUID(),
			workspace_id: canvasStore.workspaceId,
			type: 'freehand',
			pos_x: points[0].x,
			pos_y: points[0].y,
			z_index: canvasStore.elements.length,
			path_data: pathData,
			style: {
				strokeColor: '#000000',
				strokeWidth: 2,
				opacity: 1
			}
		};

		canvasStore.addElement(newElement);
	}

	function createTextElement(x: number, y: number) {
		if (!canvasStore.workspaceId) return;

		const newElement: CanvasElementType = {
			id: crypto.randomUUID(),
			workspace_id: canvasStore.workspaceId,
			type: 'text',
			pos_x: x,
			pos_y: y,
			width: 200,
			height: 100,
			z_index: canvasStore.elements.length,
			content: '',
			html_content: '<p>Double-click to edit</p>',
			style: {
				fontSize: 16,
				fontFamily: 'Inter',
				color: '#000000',
				textAlign: 'left'
			}
		};

		canvasStore.addElement(newElement);
		canvasStore.select(newElement.id);
	}

	function createStickyNote(x: number, y: number) {
		if (!canvasStore.workspaceId) return;

		const colors = ['#fef3c7', '#d1fae5', '#dbeafe', '#fce7f3', '#e0e7ff'];
		const randomColor = colors[Math.floor(Math.random() * colors.length)];

		const newElement: CanvasElementType = {
			id: crypto.randomUUID(),
			workspace_id: canvasStore.workspaceId,
			type: 'sticky',
			pos_x: x,
			pos_y: y,
			width: 200,
			height: 200,
			z_index: canvasStore.elements.length,
			content: '',
			html_content: '<p>Double-click to edit</p>',
			style: {
				backgroundColor: randomColor,
				fontSize: 14,
				color: '#000000'
			}
		};

		canvasStore.addElement(newElement);
		canvasStore.select(newElement.id);
	}

	function createListElement(x: number, y: number, listType: 'bullet' | 'numbered' | 'checkbox') {
		if (!canvasStore.workspaceId) return;

		const newElement: CanvasElementType = {
			id: crypto.randomUUID(),
			workspace_id: canvasStore.workspaceId,
			type: 'list',
			pos_x: x,
			pos_y: y,
			width: 250,
			height: 200,
			z_index: canvasStore.elements.length,
			content: '',
			html_content: '<ul><li>Item 1</li><li>Item 2</li></ul>',
			style: {
				fontSize: 14,
				fontFamily: 'Inter',
				color: '#000000',
				listType: listType
			}
		};

		canvasStore.addElement(newElement);
		canvasStore.select(newElement.id);
	}

	function createConnectorElement(
		startX: number,
		startY: number,
		endX: number,
		endY: number,
		type: 'straight' | 'curved' | 'elbow' = 'straight'
	) {
		if (!canvasStore.workspaceId) return;

		const newElement: CanvasElementType = {
			id: crypto.randomUUID(),
			workspace_id: canvasStore.workspaceId,
			type: 'connector',
			pos_x: Math.min(startX, endX),
			pos_y: Math.min(startY, endY),
			z_index: canvasStore.elements.length,
			connector_data: {
				startX: startX,
				startY: startY,
				endX: endX,
				endY: endY,
				endArrow: true,
				startArrow: false
			},
			style: {
				strokeColor: '#1e40af',
				strokeWidth: 2,
				connectorType: type,
				opacity: 1
			}
		};

		canvasStore.addElement(newElement);
		canvasStore.select(newElement.id);
	}

	// Element manipulation functions

	function startDragging(e: MouseEvent) {
		if (canvasStore.selectedIds.length === 0) return;

		isDraggingElements = true;
		canvasStore.setIsDragging(true);

		const canvasPoint = canvasStore.screenToCanvas(e.clientX, e.clientY);
		dragStart = canvasPoint;

		// Сохраняем начальные позиции всех выделенных элементов
		elementStartPositions.clear();
		canvasStore.selectedElements.forEach((el) => {
			elementStartPositions.set(el.id, { x: el.pos_x, y: el.pos_y });
		});
	}

	function updateDragging(e: MouseEvent) {
		if (!isDraggingElements || !dragStart) return;

		const canvasPoint = canvasStore.screenToCanvas(e.clientX, e.clientY);
		let dx = canvasPoint.x - dragStart.x;
		let dy = canvasPoint.y - dragStart.y;

		// Shift для ограничения по осям
		if (e.shiftKey) {
			if (Math.abs(dx) > Math.abs(dy)) {
				dy = 0;
			} else {
				dx = 0;
			}
		}

		// Обновляем позиции всех выделенных элементов
		canvasStore.selectedElements.forEach((el) => {
			const startPos = elementStartPositions.get(el.id);
			if (startPos) {
				let newX = startPos.x + dx;
				let newY = startPos.y + dy;

				// Snap to grid если включено
				if (canvasStore.snapToGrid) {
					const snapped = canvasStore.snapPoint({ x: newX, y: newY });
					newX = snapped.x;
					newY = snapped.y;
				}

				canvasStore.updateElement(el.id, {
					pos_x: newX,
					pos_y: newY
				});
			}
		});
	}

	function endDragging() {
		isDraggingElements = false;
		canvasStore.setIsDragging(false);
		dragStart = null;
		elementStartPositions.clear();

		// TODO: Отправить обновления на сервер
	}

	function startResize(e: MouseEvent, handle: string) {
		if (canvasStore.selectedIds.length === 0) return;

		isResizing = true;
		resizeHandle = handle;

		const canvasPoint = canvasStore.screenToCanvas(e.clientX, e.clientY);

		// Сохраняем начальные размеры и позиции
		const elementsData = new Map();
		canvasStore.selectedElements.forEach((el) => {
			elementsData.set(el.id, {
				x: el.pos_x,
				y: el.pos_y,
				width: el.width || 0,
				height: el.height || 0
			});
		});

		resizeStartData = {
			mouseX: canvasPoint.x,
			mouseY: canvasPoint.y,
			elements: elementsData
		};
	}

	function updateResize(e: MouseEvent) {
		if (!isResizing || !resizeHandle || !resizeStartData) return;

		const canvasPoint = canvasStore.screenToCanvas(e.clientX, e.clientY);
		const dx = canvasPoint.x - resizeStartData.mouseX;
		const dy = canvasPoint.y - resizeStartData.mouseY;

		const maintainAspect = e.shiftKey;

		// Для каждого выделенного элемента
		canvasStore.selectedElements.forEach((el) => {
			const startData = resizeStartData!.elements.get(el.id);
			if (!startData) return;

			let newX = startData.x;
			let newY = startData.y;
			let newWidth = startData.width;
			let newHeight = startData.height;

			// В зависимости от handle применяем изменения
			switch (resizeHandle) {
				case 'topLeft':
					newX = startData.x + dx;
					newY = startData.y + dy;
					newWidth = startData.width - dx;
					newHeight = startData.height - dy;
					break;
				case 'topRight':
					newY = startData.y + dy;
					newWidth = startData.width + dx;
					newHeight = startData.height - dy;
					break;
				case 'bottomLeft':
					newX = startData.x + dx;
					newWidth = startData.width - dx;
					newHeight = startData.height + dy;
					break;
				case 'bottomRight':
					newWidth = startData.width + dx;
					newHeight = startData.height + dy;
					break;
				case 'top':
					newY = startData.y + dy;
					newHeight = startData.height - dy;
					break;
				case 'right':
					newWidth = startData.width + dx;
					break;
				case 'bottom':
					newHeight = startData.height + dy;
					break;
				case 'left':
					newX = startData.x + dx;
					newWidth = startData.width - dx;
					break;
			}

			// Сохранение пропорций
			if (
				maintainAspect &&
				(resizeHandle === 'topLeft' ||
					resizeHandle === 'topRight' ||
					resizeHandle === 'bottomLeft' ||
					resizeHandle === 'bottomRight')
			) {
				const aspectRatio = startData.width / startData.height;
				newHeight = newWidth / aspectRatio;

				// Корректируем позицию если нужно
				if (resizeHandle === 'topLeft' || resizeHandle === 'topRight') {
					newY = startData.y + startData.height - newHeight;
				}
			}

			// Минимальные размеры
			newWidth = Math.max(20, newWidth);
			newHeight = Math.max(20, newHeight);

			canvasStore.updateElement(el.id, {
				pos_x: newX,
				pos_y: newY,
				width: newWidth,
				height: newHeight
			});
		});
	}

	function endResize() {
		isResizing = false;
		resizeHandle = null;
		resizeStartData = null;

		// TODO: Отправить обновления на сервер
	}

	// Mouse handlers

	// Получить координаты в canvas space с учетом viewport
	function getCanvasPoint(e: MouseEvent): { x: number; y: number } {
		if (!svgCanvas) return { x: 0, y: 0 };

		// Получаем координаты относительно SVG элемента
		const rect = svgCanvas.getBoundingClientRect();
		const x = e.clientX - rect.left;
		const y = e.clientY - rect.top;

		// Преобразуем screen координаты в canvas координаты с учетом viewport
		return canvasStore.screenToCanvas(x, y);
	}

	function handleMouseDown(e: MouseEvent) {
		// Средняя кнопка мыши или Space + левая кнопка = panning
		if (e.button === 1 || (e.button === 0 && isSpacePressed)) {
			e.preventDefault();
			isPanningActive = true;
			canvasStore.setIsPanning(true);
			panStart = { x: e.clientX, y: e.clientY };
			if (canvasContainer) {
				canvasContainer.style.cursor = 'grabbing';
			}
			return;
		}

		// Левая кнопка мыши - начало выделения/манипуляций (если инструмент = select)
		if (e.button === 0 && canvasStore.activeTool === 'select' && !isSpacePressed) {
			const target = e.target as SVGElement;

			// Проверяем клик по resize handle
			const handleEl = target.closest('[data-handle]');
			if (handleEl) {
				const handle = handleEl.getAttribute('data-handle');
				if (handle && handle !== 'rotate') {
					startResize(e, handle);
					return;
				}
				if (handle === 'rotate') {
					// TODO: implement rotation
					return; // Prevent deselection
				}
			}

			// Проверяем, кликнули ли по элементу
			const elementId = target.closest('[data-element-id]')?.getAttribute('data-element-id');

			if (elementId) {
				// Клик по элементу
				const addToSelection = e.shiftKey || e.ctrlKey || e.metaKey;
				if (addToSelection) {
					canvasStore.toggleSelection(elementId);
				} else {
					if (!canvasStore.isSelected(elementId)) {
						canvasStore.select(elementId);
					}
				}
				// Начать dragging выделенных элементов
				startDragging(e);
			} else {
				// Клик по пустому месту - начать selection box
				if (!e.shiftKey && !e.ctrlKey && !e.metaKey) {
					canvasStore.clearSelection();
				}
				const point = getCanvasPoint(e);
				canvasStore.startSelectionBox(point.x, point.y);
			}
			return;
		}

		// Shape creation tools
		const shapeTools = ['rectangle', 'ellipse', 'triangle', 'line', 'arrow'];
		if (e.button === 0 && shapeTools.includes(canvasStore.activeTool)) {
			const point = getCanvasPoint(e);
			isCreatingShape = true;
			shapeCreateStart = point;
			shapeCreateCurrent = point;
			return;
		}

		// Freehand drawing tool
		if (e.button === 0 && canvasStore.activeTool === 'freehand') {
			const point = getCanvasPoint(e);
			isDrawing = true;
			drawingPoints = [point];
			return;
		}

		// Text tool - single click creates text
		if (e.button === 0 && canvasStore.activeTool === 'text') {
			const point = getCanvasPoint(e);
			createTextElement(point.x, point.y);
			canvasStore.setTool('select');
			return;
		}

		// Sticky note tool - single click creates sticky
		if (e.button === 0 && canvasStore.activeTool === 'sticky') {
			const point = getCanvasPoint(e);
			createStickyNote(point.x, point.y);
			canvasStore.setTool('select');
			return;
		}

		// Image tool - trigger file upload
		if (e.button === 0 && canvasStore.activeTool === 'image') {
			imageUploader?.openFileDialog();
			canvasStore.setTool('select');
			return;
		}

		// List tool - single click creates list (TODO: add submenu for list type)
		if (e.button === 0 && canvasStore.activeTool === 'list') {
			const point = getCanvasPoint(e);
			createListElement(point.x, point.y, 'bullet');
			canvasStore.setTool('select');
			return;
		}

		// Connector tool - drag to create connector
		if (e.button === 0 && canvasStore.activeTool === 'connector') {
			const point = getCanvasPoint(e);
			isCreatingConnector = true;
			connectorStart = point;
			connectorCurrent = point;
			return;
		}
	}

	let rafId: number | null = null;
	let pendingMouseEvent: MouseEvent | null = null;

	function handleMouseMove(e: MouseEvent) {
		// Используем requestAnimationFrame для throttling
		pendingMouseEvent = e;

		if (rafId === null) {
			rafId = requestAnimationFrame(() => {
				if (!pendingMouseEvent) return;

				const event = pendingMouseEvent;
				pendingMouseEvent = null;
				rafId = null;

				// Panning
				if (isPanningActive) {
					const dx = event.clientX - panStart.x;
					const dy = event.clientY - panStart.y;
					canvasStore.pan(dx, dy);
					panStart = { x: event.clientX, y: event.clientY };
					return;
				}

				// Element dragging
				if (isDraggingElements) {
					updateDragging(event);
					return;
				}

				// Element resizing
				if (isResizing) {
					updateResize(event);
					return;
				}

				// Update selection box
				if (canvasStore.isSelecting) {
					const point = getCanvasPoint(event);
					canvasStore.updateSelectionBox(point.x, point.y);
					return;
				}

				// Shape creation
				if (isCreatingShape && shapeCreateStart) {
					const point = getCanvasPoint(event);
					shapeCreateCurrent = point;
					return;
				}

				// Freehand drawing
				if (isDrawing) {
					const point = getCanvasPoint(event);
					drawingPoints = [...drawingPoints, point];
					return;
				}

				// Connector creation
				if (isCreatingConnector && connectorStart) {
					const point = getCanvasPoint(event);
					connectorCurrent = point;
					return;
				}
			});
		}
	}

	function handleMouseUp(e: MouseEvent) {
		if (isPanningActive) {
			isPanningActive = false;
			canvasStore.setIsPanning(false);
			if (canvasContainer) {
				canvasContainer.style.cursor = isSpacePressed ? 'grab' : 'default';
			}
			return;
		}

		if (isDraggingElements) {
			endDragging();
			return;
		}

		if (isResizing) {
			endResize();
			return;
		}

		// End selection box
		if (canvasStore.isSelecting) {
			const addToSelection = e.shiftKey || e.ctrlKey || e.metaKey;
			canvasStore.endSelectionBox(addToSelection);
			return;
		}

		// End shape creation
		if (isCreatingShape && shapeCreateStart && shapeCreateCurrent) {
			const tool = canvasStore.activeTool as
				| 'rectangle'
				| 'ellipse'
				| 'triangle'
				| 'line'
				| 'arrow';
			createShapeElement(
				tool,
				shapeCreateStart.x,
				shapeCreateStart.y,
				shapeCreateCurrent.x,
				shapeCreateCurrent.y,
				isShiftPressed
			);
			isCreatingShape = false;
			shapeCreateStart = null;
			shapeCreateCurrent = null;
			canvasStore.setTool('select');
			return;
		}

		// End freehand drawing
		if (isDrawing) {
			createFreehandElement(drawingPoints);
			isDrawing = false;
			drawingPoints = [];
			canvasStore.setTool('select');
			return;
		}

		// End connector creation
		if (isCreatingConnector && connectorStart && connectorCurrent) {
			createConnectorElement(
				connectorStart.x,
				connectorStart.y,
				connectorCurrent.x,
				connectorCurrent.y,
				connectorType
			);
			isCreatingConnector = false;
			connectorStart = null;
			connectorCurrent = null;
			canvasStore.setTool('select');
			return;
		}
	}

	function handleWheel(e: WheelEvent) {
		e.preventDefault();

		// Ctrl/Cmd + Wheel = zoom
		if (e.ctrlKey || e.metaKey) {
			const delta = -e.deltaY * 0.001; // конвертируем wheel delta в zoom delta
			canvasStore.zoom(delta, e.clientX, e.clientY);
		} else {
			// Обычный scroll = pan
			canvasStore.pan(-e.deltaX, -e.deltaY);
		}
	}

	// Derived values для SVG transform
	const viewport = $derived(canvasStore.viewport);
	const transform = $derived(`translate(${viewport.x}, ${viewport.y}) scale(${viewport.zoom})`);

	const elements = $derived(canvasStore.elements);
	const selectedIds = $derived(canvasStore.selectedIds);
	const showGrid = $derived(canvasStore.showGrid);
	const selectionBox = $derived(canvasStore.selectionBox);

	// Виртуализация: рендерим только видимые элементы + небольшой margin
	const visibleElements = $derived(() => {
		if (!canvasContainer) return elements;

		const rect = canvasContainer.getBoundingClientRect();
		const margin = 200; // margin в пикселях

		return elements.filter((el) => {
			// Преобразуем позицию элемента в screen координаты
			const screenPos = {
				x: el.pos_x * viewport.zoom + viewport.x,
				y: el.pos_y * viewport.zoom + viewport.y
			};
			const screenWidth = (el.width || 0) * viewport.zoom;
			const screenHeight = (el.height || 0) * viewport.zoom;

			// Проверяем пересечение с видимой областью (с margin)
			return !(
				screenPos.x + screenWidth < -margin ||
				screenPos.x > rect.width + margin ||
				screenPos.y + screenHeight < -margin ||
				screenPos.y > rect.height + margin
			);
		});
	});
</script>

<div
	class="canvas-container"
	bind:this={canvasContainer}
	onmousedown={handleMouseDown}
	onmousemove={handleMouseMove}
	onmouseup={handleMouseUp}
	onmouseleave={handleMouseUp}
	onwheel={handleWheel}
	role="img"
	aria-label="Canvas workspace"
>
	<svg bind:this={svgCanvas} class="canvas-svg" width="100%" height="100%">
		<!-- Background -->
		<rect width="100%" height="100%" fill="#f8f9fa" />

		<!-- Main canvas group с viewport transform -->
		<g {transform}>
			<!-- Grid (внутри transform чтобы двигалась с viewport) -->
			{#if showGrid}
				<CanvasGrid {viewport} />
			{/if}

			<!-- Render elements (только видимые для производительности) -->
			{#each visibleElements() as element (element.id)}
				<CanvasElement
					{element}
					isSelected={selectedIds.includes(element.id)}
					onEditStart={() => (isEditingText = true)}
					onEditEnd={() => (isEditingText = false)}
				/>
			{/each}

			<!-- Selection box (внутри transform, в canvas coordinates) -->
			{#if selectionBox}
				<SelectionBox box={selectionBox} {viewport} />
			{/if}

			<!-- Shape creation preview -->
			{#if isCreatingShape && shapeCreateStart && shapeCreateCurrent}
				<ShapeCreator
					tool={canvasStore.activeTool}
					startX={shapeCreateStart.x}
					startY={shapeCreateStart.y}
					currentX={shapeCreateCurrent.x}
					currentY={shapeCreateCurrent.y}
					shiftPressed={isShiftPressed}
				/>
			{/if}

			<!-- Freehand drawing preview -->
			{#if isDrawing && drawingPoints.length > 0}
				<FreehandDrawing points={drawingPoints} color="#000000" width={2} />
			{/if}

			<!-- Connector creation preview -->
			{#if isCreatingConnector && connectorStart && connectorCurrent}
				<ConnectorCreator
					startX={connectorStart.x}
					startY={connectorStart.y}
					currentX={connectorCurrent.x}
					currentY={connectorCurrent.y}
					{connectorType}
				/>
			{/if}
		</g>

		<!-- Selection handles (в screen coordinates) -->
		{#if selectedIds.length > 0}
			<SelectionHandles elements={canvasStore.selectedElements} {viewport} />
		{/if}
	</svg>

	<!-- Zoom controls -->
	<ZoomControls />

	<!-- Image uploader (hidden, triggered programmatically) -->
	{#if canvasStore.workspaceId}
		<ImageUploader
			bind:this={imageUploader}
			workspaceId={canvasStore.workspaceId}
			onImageUploaded={handleImageUploaded}
		/>
	{/if}
</div>

<style>
	.canvas-container {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: hidden;
		background: #e5e7eb;
		cursor: default;
		user-select: none;
	}

	.canvas-container:focus {
		outline: none;
	}

	.canvas-container:focus-visible {
		outline: 2px solid #3b82f6;
		outline-offset: -2px;
	}

	.canvas-svg {
		display: block;
		touch-action: none; /* Disable browser touch gestures */
	}
</style>
