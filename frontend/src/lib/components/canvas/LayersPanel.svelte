<script lang="ts">
	import { canvasStore } from '$lib/stores/canvas.svelte';
	import { canvas } from '$lib/stores/canvasWithHistory.svelte';
	import type { CanvasElement } from '$lib/types/api';
	import {
		Eye,
		EyeOff,
		Lock,
		Unlock,
		Type,
		Square,
		Circle,
		Triangle,
		Minus,
		ArrowRight,
		Pencil,
		StickyNote,
		Image as ImageIcon,
		List,
		ArrowRightLeft
	} from 'lucide-svelte';

	const elements = $derived(canvasStore.elements);
	const selectedIds = $derived(canvasStore.selectedIds);

	// Sort elements by z-index (top to bottom in UI)
	const sortedElements = $derived(
		[...elements].sort((a, b) => (b.z_index || 0) - (a.z_index || 0))
	);

	// Element type icons
	const typeIcons: Record<string, any> = {
		text: Type,
		rectangle: Square,
		ellipse: Circle,
		triangle: Triangle,
		line: Minus,
		arrow: ArrowRight,
		freehand: Pencil,
		sticky: StickyNote,
		image: ImageIcon,
		list: List,
		connector: ArrowRightLeft
	};

	// Element visibility (stored locally, not persisted)
	let visibilityMap = $state<Record<string, boolean>>({});

	// Initialize visibility for all elements
	$effect(() => {
		elements.forEach((el) => {
			if (visibilityMap[el.id] === undefined) {
				visibilityMap[el.id] = true;
			}
		});
	});

	function getElementLabel(element: CanvasElement): string {
		// Try to get a meaningful label from content
		if (element.content) {
			const text = element.content.substring(0, 30);
			return text.length < element.content.length ? `${text}...` : text;
		}

		// Fallback to type
		return element.type.charAt(0).toUpperCase() + element.type.slice(1);
	}

	function handleElementClick(elementId: string, event: MouseEvent) {
		if (event.ctrlKey || event.metaKey) {
			// Multi-select
			canvasStore.toggleSelection(elementId);
		} else {
			// Single select
			canvasStore.select(elementId);
		}
	}

	function toggleVisibility(elementId: string, event: MouseEvent) {
		event.stopPropagation();
		visibilityMap[elementId] = !visibilityMap[elementId];
		// TODO: Implement actual visibility hiding in Canvas
	}

	function toggleLock(elementId: string, event: MouseEvent) {
		event.stopPropagation();
		const element = elements.find((el) => el.id === elementId);
		if (element) {
			canvas.updateElement(elementId, { locked: !element.locked });
		}
	}

	// Drag & Drop for reordering
	let draggedElement: string | null = $state(null);
	let dropTarget: string | null = $state(null);

	function handleDragStart(elementId: string, event: DragEvent) {
		draggedElement = elementId;
		if (event.dataTransfer) {
			event.dataTransfer.effectAllowed = 'move';
		}
	}

	function handleDragOver(elementId: string, event: DragEvent) {
		event.preventDefault();
		if (draggedElement && draggedElement !== elementId) {
			dropTarget = elementId;
		}
	}

	function handleDragLeave() {
		dropTarget = null;
	}

	function handleDrop(targetElementId: string, event: DragEvent) {
		event.preventDefault();

		if (!draggedElement || draggedElement === targetElementId) {
			dropTarget = null;
			draggedElement = null;
			return;
		}

		const draggedEl = elements.find((el) => el.id === draggedElement);
		const targetEl = elements.find((el) => el.id === targetElementId);

		if (!draggedEl || !targetEl) {
			dropTarget = null;
			draggedElement = null;
			return;
		}

		// Reorder z-indices: move dragged element to target position
		// Note: sortedElements is sorted by z-index descending (high to low)
		// So we need to work with this reversed order
		const currentOrder = [...sortedElements]; // Already sorted high to low
		const draggedIndex = currentOrder.findIndex((el) => el.id === draggedElement);
		const targetIndex = currentOrder.findIndex((el) => el.id === targetElementId);

		if (draggedIndex === -1 || targetIndex === -1) {
			dropTarget = null;
			draggedElement = null;
			return;
		}

		// Remove dragged element from array
		const [removed] = currentOrder.splice(draggedIndex, 1);
		// Insert at target position
		currentOrder.splice(targetIndex, 0, removed);

		// Reassign z-indices based on new order (reversed because display is high to low)
		// Highest z-index at index 0, lowest at the end
		const maxZ = currentOrder.length - 1;
		const batchUpdates = currentOrder.map((el, index) => ({
			id: el.id,
			updates: { z_index: maxZ - index }
		}));

		// Use batch update for better performance and proper reactivity
		canvas.updateElements(batchUpdates);

		dropTarget = null;
		draggedElement = null;
	}

	function handleDragEnd() {
		dropTarget = null;
		draggedElement = null;
	}
</script>

<div class="flex h-full w-[250px] flex-col overflow-hidden border-r border-gray-200 bg-white">
	<div class="flex items-center justify-between border-b border-gray-200 px-4 py-4">
		<h3 class="m-0 text-sm font-semibold text-gray-900">Layers</h3>
		<div class="text-xs text-gray-400">{elements.length} elements</div>
	</div>

	<div
		class="flex-1 overflow-y-auto [&::-webkit-scrollbar]:w-1.5 [&::-webkit-scrollbar-thumb]:rounded [&::-webkit-scrollbar-thumb]:bg-gray-300 [&::-webkit-scrollbar-thumb:hover]:bg-gray-400 [&::-webkit-scrollbar-track]:bg-transparent"
	>
		{#if sortedElements.length === 0}
			<div class="flex items-center justify-center px-6 py-12 text-center text-sm text-gray-400">
				<p>No elements yet</p>
			</div>
		{:else}
			<div class="p-2">
				{#each sortedElements as element (element.id)}
					{@const Icon = typeIcons[element.type] || Square}
					{@const isSelected = selectedIds.includes(element.id)}
					{@const isVisible = visibilityMap[element.id] !== false}
					{@const isLocked = element.locked || false}
					{@const isDragging = draggedElement === element.id}
					{@const isDropTarget = dropTarget === element.id}

					<div
						class="mb-0.5 flex cursor-pointer items-center gap-2 rounded-md px-2 py-2 transition-all duration-150 select-none hover:bg-gray-100"
						class:bg-blue-100={isSelected}
						class:text-blue-600={isSelected}
						class:opacity-40={isDragging || !isVisible}
						class:border-t-2={isDropTarget}
						class:border-blue-500={isDropTarget}
						onclick={(e) => handleElementClick(element.id, e)}
						draggable="true"
						ondragstart={(e) => handleDragStart(element.id, e)}
						ondragover={(e) => handleDragOver(element.id, e)}
						ondragleave={handleDragLeave}
						ondrop={(e) => handleDrop(element.id, e)}
						ondragend={handleDragEnd}
					>
						<!-- Type icon -->
						<div
							class="flex flex-shrink-0 items-center justify-center"
							class:text-gray-500={!isSelected}
							class:text-blue-600={isSelected}
						>
							<Icon size={16} />
						</div>

						<!-- Label -->
						<div
							class="flex-1 overflow-hidden text-[13px] font-medium text-ellipsis whitespace-nowrap"
							class:text-gray-700={!isSelected}
							class:text-blue-600={isSelected}
							title={getElementLabel(element)}
						>
							{getElementLabel(element)}
						</div>

						<!-- Controls -->
						<div
							class="flex gap-1 opacity-0 transition-opacity duration-150 group-hover:opacity-100 has-[.active]:opacity-100"
						>
							<!-- Visibility toggle -->
							<button
								class="flex cursor-pointer items-center justify-center rounded border-none bg-transparent p-1 text-gray-400 transition-all duration-150 hover:bg-white hover:text-gray-700"
								class:text-red-500={!isVisible}
								class:active={!isVisible}
								onclick={(e) => toggleVisibility(element.id, e)}
								title={isVisible ? 'Hide' : 'Show'}
							>
								{#if isVisible}
									<Eye size={14} />
								{:else}
									<EyeOff size={14} />
								{/if}
							</button>

							<!-- Lock toggle -->
							<button
								class="flex cursor-pointer items-center justify-center rounded border-none bg-transparent p-1 text-gray-400 transition-all duration-150 hover:bg-white hover:text-gray-700"
								class:text-red-500={isLocked}
								class:active={isLocked}
								onclick={(e) => toggleLock(element.id, e)}
								title={isLocked ? 'Unlock' : 'Lock'}
							>
								{#if isLocked}
									<Lock size={14} />
								{:else}
									<Unlock size={14} />
								{/if}
							</button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
	.flex:hover .opacity-0 {
		opacity: 1;
	}
</style>
