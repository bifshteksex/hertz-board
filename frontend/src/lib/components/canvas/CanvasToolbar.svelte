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
		Magnet,
		ChevronDown as SubmenuIcon
	} from 'lucide-svelte';
	import IconUndo from '$components/icons/IconUndo.svelte';
	import IconRedo from '$components/icons/IconRedo.svelte';
	import ShapeSubmenu from './ShapeSubmenu.svelte';
	import ListSubmenu from './ListSubmenu.svelte';

	import IconQuestion from '$components/icons/IconQuestion.svelte';

	interface Props {
		canUndo?: boolean;
		canRedo?: boolean;
		onUndo?: () => void;
		onRedo?: () => void;
		onShowHelp?: () => void;
	}

	let { canUndo = false, canRedo = false, onUndo, onRedo, onShowHelp }: Props = $props();

	const activeTool = $derived(canvasStore.activeTool);
	const selectedCount = $derived(canvasStore.selectedIds.length);
	const showGrid = $derived(canvasStore.showGrid);
	const snapToGrid = $derived(canvasStore.snapToGrid);

	// Submenu state
	let showShapeSubmenu = $state(false);
	let showListSubmenu = $state(false);
	let activeShapeType = $state<'rectangle' | 'ellipse' | 'triangle' | 'line' | 'arrow'>(
		'rectangle'
	);
	let activeListType = $state<'bullet' | 'numbered' | 'checkbox'>('bullet');

	const basicTools: Array<{ id: Tool; icon: any; label: string }> = [
		{ id: 'select', icon: MousePointer2, label: 'Select (V)' },
		{ id: 'text', icon: Type, label: 'Text (T)' },
		{ id: 'freehand', icon: Pencil, label: 'Pen (P)' },
		{ id: 'sticky', icon: StickyNote, label: 'Sticky Note (S)' },
		{ id: 'image', icon: ImageIcon, label: 'Image (I)' },
		{ id: 'connector', icon: ArrowRightLeft, label: 'Connector' }
	];

	// Shape icons based on active type
	const shapeIcons = {
		rectangle: Square,
		ellipse: Circle,
		triangle: Triangle,
		line: Minus,
		arrow: ArrowRight
	};

	const listIcons = {
		bullet: List,
		numbered: List,
		checkbox: List
	};

	function selectTool(tool: Tool) {
		canvasStore.setTool(tool);
		// Close any open submenus
		showShapeSubmenu = false;
		showListSubmenu = false;
	}

	function toggleShapeSubmenu() {
		showShapeSubmenu = !showShapeSubmenu;
		showListSubmenu = false;
	}

	function toggleListSubmenu() {
		showListSubmenu = !showListSubmenu;
		showShapeSubmenu = false;
	}

	function handleShapeSelect(shape: typeof activeShapeType) {
		activeShapeType = shape;
		canvasStore.setTool(shape);
		showShapeSubmenu = false;
	}

	function handleListSelect(listType: typeof activeListType) {
		activeListType = listType;
		// Store list type for later use
		(canvasStore as any).activeListType = listType;
		canvasStore.setTool('list');
		showListSubmenu = false;
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
	<!-- Left section: Undo/Redo and Help -->
	<div class="toolbar-section">
		<button
			class="tool-btn"
			class:disabled={!canUndo}
			disabled={!canUndo}
			onclick={onUndo}
			title="Undo (Ctrl+Z)"
			aria-label="Undo"
		>
			<IconUndo size={18} />
		</button>

		<button
			class="tool-btn"
			class:disabled={!canRedo}
			disabled={!canRedo}
			onclick={onRedo}
			title="Redo (Ctrl+Y)"
			aria-label="Redo"
		>
			<IconRedo size={18} />
		</button>

		<div class="separator"></div>

		<button
			class="tool-btn"
			onclick={onShowHelp}
			title="Keyboard Shortcuts (?)"
			aria-label="Show keyboard shortcuts"
		>
			<IconQuestion size={24} />
		</button>

		<div class="separator"></div>
	</div>

	<!-- Middle section: Drawing tools -->
	<div class="toolbar-section">
		{#each basicTools as tool}
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

		<!-- Shape tool with submenu -->
		<div class="tool-with-submenu">
			<button
				class="tool-btn"
				class:active={['rectangle', 'ellipse', 'triangle', 'line', 'arrow'].includes(activeTool)}
				onclick={toggleShapeSubmenu}
				title="Shapes"
				aria-label="Shapes"
			>
				<svelte:component this={shapeIcons[activeShapeType]} size={18} />
				<SubmenuIcon size={12} class="submenu-indicator" />
			</button>
			{#if showShapeSubmenu}
				<ShapeSubmenu onSelect={handleShapeSelect} />
			{/if}
		</div>

		<!-- List tool with submenu -->
		<div class="tool-with-submenu">
			<button
				class="tool-btn"
				class:active={activeTool === 'list'}
				onclick={toggleListSubmenu}
				title="Lists"
				aria-label="Lists"
			>
				<svelte:component this={listIcons[activeListType]} size={18} />
				<SubmenuIcon size={12} class="submenu-indicator" />
			</button>
			{#if showListSubmenu}
				<ListSubmenu onSelect={handleListSelect} />
			{/if}
		</div>
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

	.tool-btn:disabled,
	.tool-btn.disabled {
		opacity: 0.4;
		cursor: not-allowed;
		pointer-events: none;
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

	.tool-with-submenu {
		position: relative;
	}

	.tool-with-submenu .tool-btn {
		display: flex;
		align-items: center;
		gap: 2px;
	}

	:global(.submenu-indicator) {
		opacity: 0.5;
		margin-left: -4px;
	}
</style>
