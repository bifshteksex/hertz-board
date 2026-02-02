<script lang="ts">
	import { canvasStore, type Tool } from '$lib/stores/canvas.svelte';
	import { i18n } from '$lib/stores/i18n.svelte';
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
	const brushColor = $derived(canvasStore.brushColor);
	const brushWidth = $derived(canvasStore.brushWidth);

	// Submenu state
	let showShapeSubmenu = $state(false);
	let showListSubmenu = $state(false);
	let activeShapeType = $state<'rectangle' | 'ellipse' | 'triangle' | 'line' | 'arrow'>(
		'rectangle'
	);
	let activeListType = $state<'bullet' | 'numbered' | 'checkbox'>('bullet');

	const basicTools = $derived([
		{ id: 'select' as Tool, icon: MousePointer2, label: i18n.t('canvas.toolbar.tools.select') },
		{ id: 'text' as Tool, icon: Type, label: i18n.t('canvas.toolbar.tools.text') },
		{ id: 'freehand' as Tool, icon: Pencil, label: i18n.t('canvas.toolbar.tools.pen') },
		{ id: 'sticky' as Tool, icon: StickyNote, label: i18n.t('canvas.toolbar.tools.sticky') },
		{ id: 'image' as Tool, icon: ImageIcon, label: i18n.t('canvas.toolbar.tools.image') },
		{
			id: 'connector' as Tool,
			icon: ArrowRightLeft,
			label: i18n.t('canvas.toolbar.tools.connector')
		}
	]);

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

	function handleBrushColorChange(e: Event) {
		const input = e.target as HTMLInputElement;
		canvasStore.setBrushColor(input.value);
	}

	function handleBrushWidthChange(e: Event) {
		const input = e.target as HTMLInputElement;
		canvasStore.setBrushWidth(Number(input.value));
	}
</script>

<div
	class="flex items-center gap-2 border-b border-gray-200 bg-white px-4 py-2 dark:border-gray-700 dark:bg-gray-800"
>
	<!-- Left section: Undo/Redo and Help -->
	<div class="flex items-center gap-1">
		<button
			class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:pointer-events-none disabled:opacity-40 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
			disabled={!canUndo}
			onclick={onUndo}
			title={i18n.t('canvas.toolbar.undoShortcut')}
			aria-label={i18n.t('canvas.toolbar.undo')}
		>
			<IconUndo size={18} />
		</button>

		<button
			class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:pointer-events-none disabled:opacity-40 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
			disabled={!canRedo}
			onclick={onRedo}
			title={i18n.t('canvas.toolbar.redoShortcut')}
			aria-label={i18n.t('canvas.toolbar.redo')}
		>
			<IconRedo size={18} />
		</button>

		<div class="mx-1 h-6 w-px bg-gray-200 dark:bg-gray-700"></div>

		<button
			class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
			onclick={onShowHelp}
			title={i18n.t('canvas.toolbar.help')}
			aria-label={i18n.t('canvas.toolbar.helpAria')}
		>
			<IconQuestion size={24} />
		</button>

		<div class="mx-1 h-6 w-px bg-gray-200 dark:bg-gray-700"></div>
	</div>

	<!-- Middle section: Drawing tools -->
	<div class="flex items-center gap-1">
		{#each basicTools as tool}
			{@const Icon = tool.icon}
			<button
				class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 active:bg-blue-200 data-[active=true]:bg-blue-100 data-[active=true]:text-blue-600 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100 dark:data-[active=true]:bg-blue-900/30 dark:data-[active=true]:text-blue-400"
				data-active={activeTool === tool.id}
				onclick={() => selectTool(tool.id)}
				title={tool.label}
				aria-label={tool.label}
			>
				<Icon size={18} />
			</button>
		{/each}

		<!-- Shape tool with submenu -->
		<div class="relative">
			<button
				class="flex h-9 w-9 items-center justify-center gap-0.5 rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 data-[active=true]:bg-blue-100 data-[active=true]:text-blue-600 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100 dark:data-[active=true]:bg-blue-900/30 dark:data-[active=true]:text-blue-400"
				data-active={['rectangle', 'ellipse', 'triangle', 'line', 'arrow'].includes(activeTool)}
				onclick={toggleShapeSubmenu}
				title={i18n.t('canvas.toolbar.tools.shapes')}
				aria-label={i18n.t('canvas.toolbar.tools.shapes')}
			>
				{#if shapeIcons[activeShapeType]}
					{@const ShapeIcon = shapeIcons[activeShapeType]}
					<ShapeIcon size={18} />
				{/if}
				<SubmenuIcon size={12} class="-ml-1 opacity-50" />
			</button>
			{#if showShapeSubmenu}
				<ShapeSubmenu onSelect={handleShapeSelect} />
			{/if}
		</div>

		<!-- List tool with submenu -->
		<div class="relative">
			<button
				class="flex h-9 w-9 items-center justify-center gap-0.5 rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 data-[active=true]:bg-blue-100 data-[active=true]:text-blue-600 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100 dark:data-[active=true]:bg-blue-900/30 dark:data-[active=true]:text-blue-400"
				data-active={activeTool === 'list'}
				onclick={toggleListSubmenu}
				title={i18n.t('canvas.toolbar.tools.lists')}
				aria-label={i18n.t('canvas.toolbar.tools.lists')}
			>
				{#if listIcons[activeListType]}
					{@const ListIcon = listIcons[activeListType]}
					<ListIcon size={18} />
				{/if}
				<SubmenuIcon size={12} class="-ml-1 opacity-50" />
			</button>
			{#if showListSubmenu}
				<ListSubmenu onSelect={handleListSelect} />
			{/if}
		</div>
	</div>

	<!-- Brush settings (visible when freehand tool is active) -->
	{#if activeTool === 'freehand'}
		<div class="flex items-center gap-2">
			<div class="mx-1 h-6 w-px bg-gray-200 dark:bg-gray-700"></div>

			<!-- Color picker -->
			<div class="flex items-center gap-1.5">
				<label
					class="relative flex h-9 w-9 cursor-pointer items-center justify-center overflow-hidden rounded-md border-2 border-gray-300 transition-all hover:border-gray-400 dark:border-gray-600 dark:hover:border-gray-500"
					title="Brush color"
				>
					<input
						type="color"
						value={brushColor}
						oninput={handleBrushColorChange}
						class="absolute inset-0 h-full w-full cursor-pointer opacity-0"
					/>
					<div
						class="pointer-events-none h-6 w-6 rounded"
						style="background-color: {brushColor};"
					></div>
				</label>
			</div>

			<!-- Width slider -->
			<div class="flex items-center gap-2">
				<span class="text-xs text-gray-500 dark:text-gray-400">Width:</span>
				<input
					type="range"
					min="1"
					max="50"
					value={brushWidth}
					oninput={handleBrushWidthChange}
					class="h-2 w-24 cursor-pointer appearance-none rounded-lg bg-gray-200 dark:bg-gray-700"
					title="Brush width: {brushWidth}px"
				/>
				<span class="min-w-[2rem] text-xs text-gray-600 dark:text-gray-300">{brushWidth}px</span>
			</div>
		</div>
	{/if}

	<!-- Middle section: Z-order and grouping (visible when selection exists) -->
	{#if selectedCount > 0}
		<div class="flex items-center gap-1">
			<div class="mx-1 h-6 w-px bg-gray-200 dark:bg-gray-700"></div>

			<!-- Z-order controls -->
			<button
				class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
				onclick={handleBringToFront}
				title={i18n.t('canvas.toolbar.zorder.bringToFront')}
				aria-label={i18n.t('canvas.toolbar.zorder.bringToFront')}
			>
				<ChevronsUp size={18} />
			</button>

			<button
				class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
				onclick={handleBringForward}
				title={i18n.t('canvas.toolbar.zorder.bringForward')}
				aria-label={i18n.t('canvas.toolbar.zorder.bringForward')}
			>
				<ChevronUp size={18} />
			</button>

			<button
				class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
				onclick={handleSendBackward}
				title={i18n.t('canvas.toolbar.zorder.sendBackward')}
				aria-label={i18n.t('canvas.toolbar.zorder.sendBackward')}
			>
				<ChevronDown size={18} />
			</button>

			<button
				class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
				onclick={handleSendToBack}
				title={i18n.t('canvas.toolbar.zorder.sendToBack')}
				aria-label={i18n.t('canvas.toolbar.zorder.sendToBack')}
			>
				<ChevronsDown size={18} />
			</button>

			<div class="mx-1 h-6 w-px bg-gray-200 dark:bg-gray-700"></div>

			<!-- Grouping controls -->
			{#if selectedCount > 1}
				<button
					class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
					onclick={handleGroup}
					title={i18n.t('canvas.toolbar.grouping.group')}
					aria-label={i18n.t('canvas.toolbar.grouping.groupElements')}
				>
					<Group size={18} />
				</button>
			{/if}

			<button
				class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
				onclick={handleUngroup}
				title={i18n.t('canvas.toolbar.grouping.ungroup')}
				aria-label={i18n.t('canvas.toolbar.grouping.ungroupElements')}
			>
				<Ungroup size={18} />
			</button>
		</div>
	{/if}

	<!-- Right section: View controls -->
	<div class="ml-auto flex items-center gap-1">
		<button
			class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 data-[active=true]:bg-blue-100 data-[active=true]:text-blue-600 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100 dark:data-[active=true]:bg-blue-900/30 dark:data-[active=true]:text-blue-400"
			data-active={showGrid}
			onclick={toggleGrid}
			title={i18n.t('canvas.toolbar.view.toggleGrid')}
			aria-label={i18n.t('canvas.toolbar.view.gridAria')}
		>
			<Grid3x3 size={18} />
		</button>

		<button
			class="flex h-9 w-9 items-center justify-center rounded-md border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 data-[active=true]:bg-blue-100 data-[active=true]:text-blue-600 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100 dark:data-[active=true]:bg-blue-900/30 dark:data-[active=true]:text-blue-400"
			data-active={snapToGrid}
			onclick={toggleSnap}
			title={i18n.t('canvas.toolbar.view.toggleSnap')}
			aria-label={i18n.t('canvas.toolbar.view.snapAria')}
		>
			<Magnet size={18} />
		</button>
	</div>
</div>
