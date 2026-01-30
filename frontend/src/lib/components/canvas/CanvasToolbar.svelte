<script lang="ts">
	import { canvasStore, type Tool } from '$lib/stores/canvas.svelte';
	import {
		MousePointer2,
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
		ArrowRightLeft,
		ChevronUp,
		ChevronDown,
		ChevronsUp,
		ChevronsDown,
		Group,
		Ungroup,
		Grid3x3,
		Magnet
	} from 'lucide-svelte';

	const activeTool = $derived(canvasStore.activeTool);
	const selectedCount = $derived(canvasStore.selectedIds.length);
	const showGrid = $derived(canvasStore.showGrid);
	const snapToGrid = $derived(canvasStore.snapToGrid);

	const tools: Array<{ id: Tool; icon: any; label: string }> = [
		{ id: 'select', icon: MousePointer2, label: 'Select (V)' },
		{ id: 'text', icon: Type, label: 'Text (T)' },
		{ id: 'rectangle', icon: Square, label: 'Rectangle (R)' },
		{ id: 'ellipse', icon: Circle, label: 'Circle (C)' },
		{ id: 'triangle', icon: Triangle, label: 'Triangle' },
		{ id: 'line', icon: Minus, label: 'Line (L)' },
		{ id: 'arrow', icon: ArrowRight, label: 'Arrow (A)' },
		{ id: 'freehand', icon: Pencil, label: 'Pen (P)' },
		{ id: 'sticky', icon: StickyNote, label: 'Sticky Note (S)' },
		{ id: 'list', icon: List, label: 'List' },
		{ id: 'connector', icon: ArrowRightLeft, label: 'Connector' },
		{ id: 'image', icon: ImageIcon, label: 'Image (I)' }
	];

	function selectTool(tool: Tool) {
		canvasStore.setTool(tool);
	}

	function handleBringToFront() {
		canvasStore.bringToFront();
	}

	function handleSendToBack() {
		canvasStore.sendToBack();
	}

	function handleBringForward() {
		canvasStore.bringForward();
	}

	function handleSendBackward() {
		canvasStore.sendBackward();
	}

	function handleGroup() {
		canvasStore.groupSelected();
	}

	function handleUngroup() {
		canvasStore.ungroupSelected();
	}

	function toggleGrid() {
		canvasStore.toggleGrid();
	}

	function toggleSnap() {
		canvasStore.toggleSnap();
	}
</script>

<div class="toolbar">
	<!-- Left section: Drawing tools -->
	<div class="toolbar-section">
		{#each tools as tool}
			{@const Icon = tool.icon}
			<button
				class="tool-btn"
				class:active={activeTool === tool.id}
				onclick={() => selectTool(tool.id)}
				title={tool.label}
				aria-label={tool.label}
			>
				<Icon size={18} />
			</button>
		{/each}
	</div>

	<!-- Middle section: Z-order and grouping (visible when selection exists) -->
	{#if selectedCount > 0}
		<div class="toolbar-section">
			<div class="separator"></div>

			<!-- Z-order controls -->
			<button
				class="tool-btn"
				onclick={handleBringToFront}
				title="Bring to Front (Ctrl+Shift+])"
				aria-label="Bring to front"
			>
				<ChevronsUp size={18} />
			</button>

			<button
				class="tool-btn"
				onclick={handleBringForward}
				title="Bring Forward (Ctrl+])"
				aria-label="Bring forward"
			>
				<ChevronUp size={18} />
			</button>

			<button
				class="tool-btn"
				onclick={handleSendBackward}
				title="Send Backward (Ctrl+[)"
				aria-label="Send backward"
			>
				<ChevronDown size={18} />
			</button>

			<button
				class="tool-btn"
				onclick={handleSendToBack}
				title="Send to Back (Ctrl+Shift+[)"
				aria-label="Send to back"
			>
				<ChevronsDown size={18} />
			</button>

			<div class="separator"></div>

			<!-- Grouping controls -->
			{#if selectedCount > 1}
				<button
					class="tool-btn"
					onclick={handleGroup}
					title="Group (Ctrl+G)"
					aria-label="Group elements"
				>
					<Group size={18} />
				</button>
			{/if}

			<button
				class="tool-btn"
				onclick={handleUngroup}
				title="Ungroup (Ctrl+Shift+G)"
				aria-label="Ungroup elements"
			>
				<Ungroup size={18} />
			</button>
		</div>
	{/if}

	<!-- Right section: View controls -->
	<div class="toolbar-section ml-auto">
		<button
			class="tool-btn"
			class:active={showGrid}
			onclick={toggleGrid}
			title="Toggle Grid (Ctrl+&apos;)"
			aria-label="Toggle grid"
		>
			<Grid3x3 size={18} />
		</button>

		<button
			class="tool-btn"
			class:active={snapToGrid}
			onclick={toggleSnap}
			title="Snap to Grid (Ctrl+Shift+&apos;)"
			aria-label="Toggle snap to grid"
		>
			<Magnet size={18} />
		</button>
	</div>
</div>

<style>
	.toolbar {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 8px 16px;
		background: white;
		border-bottom: 1px solid #e5e7eb;
	}

	.toolbar-section {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	.tool-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;
		border: none;
		background: transparent;
		border-radius: 6px;
		cursor: pointer;
		color: #374151;
		transition: all 0.15s;
	}

	.tool-btn:hover {
		background: #f3f4f6;
		color: #111827;
	}

	.tool-btn.active {
		background: #dbeafe;
		color: #2563eb;
	}

	.tool-btn:active {
		background: #bfdbfe;
	}

	.separator {
		width: 1px;
		height: 24px;
		background: #e5e7eb;
		margin: 0 4px;
	}

	.ml-auto {
		margin-left: auto;
	}
</style>
