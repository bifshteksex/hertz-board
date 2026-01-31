<script lang="ts">
	import { canvasStore } from '$lib/stores/canvas.svelte';
	import { canvas } from '$lib/stores/canvasWithHistory.svelte';
	import type { CanvasElement } from '$lib/types/api';
	import ColorPicker from './ColorPicker.svelte';
	import {
		AlignLeft,
		AlignCenter,
		AlignRight,
		ChevronUp,
		ChevronDown,
		ChevronsUp,
		ChevronsDown,
		Lock,
		Unlock
	} from 'lucide-svelte';

	const selectedElements = $derived(canvasStore.selectedElements);
	const selectedCount = $derived(canvasStore.selectedIds.length);

	// Get first selected element for editing
	const element = $derived(selectedElements[0]);

	// Common properties
	let posX = $state(0);
	let posY = $state(0);
	let width = $state(0);
	let height = $state(0);
	let rotation = $state(0);
	let opacity = $state(100);
	let locked = $state(false);

	// Style properties
	let backgroundColor = $state('#3b82f6');
	let strokeColor = $state('#1e40af');
	let strokeWidth = $state(2);
	let borderRadius = $state(0);
	let color = $state('#000000');

	// Text properties
	let fontFamily = $state('Inter');
	let fontSize = $state(16);
	let fontWeight = $state('normal');
	let textAlign = $state<'left' | 'center' | 'right'>('left');

	// Update local state when element changes
	$effect(() => {
		if (element) {
			posX = element.pos_x;
			posY = element.pos_y;
			width = element.width || 0;
			height = element.height || 0;
			rotation = element.rotation || 0;
			opacity = (element.style?.opacity || 1) * 100;
			locked = element.locked || false;

			// Style properties
			backgroundColor = element.style?.backgroundColor || '#3b82f6';
			strokeColor = element.style?.strokeColor || '#1e40af';
			strokeWidth = element.style?.strokeWidth || 2;
			borderRadius = element.style?.borderRadius || 0;
			color = element.style?.color || '#000000';

			// Text properties
			fontFamily = element.style?.fontFamily || 'Inter';
			fontSize = element.style?.fontSize || 16;
			fontWeight = element.style?.fontWeight || 'normal';
			textAlign = element.style?.textAlign || 'left';
		}
	});

	function updateElement(updates: Partial<CanvasElement>) {
		if (!element) return;
		canvas.updateElement(element.id, updates);
	}

	function updateStyle(styleUpdates: Record<string, any>) {
		if (!element) return;
		canvas.updateElement(element.id, {
			style: { ...element.style, ...styleUpdates }
		});
	}

	function handlePositionChange(prop: 'pos_x' | 'pos_y', value: number) {
		updateElement({ [prop]: value });
	}

	function handleSizeChange(prop: 'width' | 'height', value: number) {
		if (value < 1) value = 1; // Minimum size
		updateElement({ [prop]: value });
	}

	function handleRotationChange(value: number) {
		updateElement({ rotation: value });
	}

	function handleOpacityChange(value: number) {
		updateStyle({ opacity: value / 100 });
	}

	function toggleLock() {
		updateElement({ locked: !locked });
	}

	function handleColorChange(prop: string, value: string) {
		updateStyle({ [prop]: value });
	}

	function handleStrokeWidthChange(value: number) {
		updateStyle({ strokeWidth: value });
	}

	function handleBorderRadiusChange(value: number) {
		updateStyle({ borderRadius: value });
	}

	function handleFontFamilyChange(value: string) {
		updateStyle({ fontFamily: value });
	}

	function handleFontSizeChange(value: number) {
		updateStyle({ fontSize: value });
	}

	function handleFontWeightChange(value: string) {
		updateStyle({ fontWeight: value });
	}

	function handleTextAlignChange(value: 'left' | 'center' | 'right') {
		updateStyle({ textAlign: value });
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

	// Font families
	const fontFamilies = [
		'Inter',
		'Arial',
		'Helvetica',
		'Times New Roman',
		'Georgia',
		'Courier New',
		'Comic Sans MS',
		'Verdana'
	];
</script>

{#if selectedCount === 0}
	<div class="flex h-full w-[280px] items-center justify-center border-l border-gray-200 bg-white">
		<div class="px-6 text-center text-sm text-gray-400">
			<p>Select an element to edit properties</p>
		</div>
	</div>
{:else if selectedCount === 1 && element}
	<div class="flex h-full w-[280px] flex-col overflow-hidden border-l border-gray-200 bg-white">
		<div class="flex items-center justify-between border-b border-gray-200 px-4 py-4">
			<h3 class="m-0 text-sm font-semibold text-gray-900">Properties</h3>
			<button
				class="cursor-pointer rounded-md border border-gray-200 bg-transparent px-1.5 py-1.5 text-gray-500 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900"
				onclick={toggleLock}
				title={locked ? 'Unlock' : 'Lock'}
			>
				{#if locked}
					<Lock size={16} />
				{:else}
					<Unlock size={16} />
				{/if}
			</button>
		</div>

		<div class="flex-1 overflow-y-auto px-4 py-4">
			<!-- Position & Size -->
			<div class="mb-6">
				<div class="mb-3 text-xs font-semibold tracking-wide text-gray-500 uppercase">
					Position & Size
				</div>

				<div class="mb-2 grid grid-cols-2 gap-2">
					<div class="flex flex-col gap-1.5">
						<label for="pos-x" class="text-xs font-medium text-gray-700">X</label>
						<input
							id="pos-x"
							type="number"
							class="rounded-md border border-gray-300 bg-white px-2 py-1.5 text-[13px] text-gray-900 focus:border-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
							bind:value={posX}
							onchange={() => handlePositionChange('pos_x', posX)}
							disabled={locked}
						/>
					</div>
					<div class="flex flex-col gap-1.5">
						<label for="pos-y" class="text-xs font-medium text-gray-700">Y</label>
						<input
							id="pos-y"
							type="number"
							class="rounded-md border border-gray-300 bg-white px-2 py-1.5 text-[13px] text-gray-900 focus:border-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
							bind:value={posY}
							onchange={() => handlePositionChange('pos_y', posY)}
							disabled={locked}
						/>
					</div>
				</div>

				<div class="mb-2 grid grid-cols-2 gap-2">
					<div class="flex flex-col gap-1.5">
						<label for="width" class="text-xs font-medium text-gray-700">W</label>
						<input
							id="width"
							type="number"
							class="rounded-md border border-gray-300 bg-white px-2 py-1.5 text-[13px] text-gray-900 focus:border-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
							bind:value={width}
							onchange={() => handleSizeChange('width', width)}
							disabled={locked}
						/>
					</div>
					<div class="flex flex-col gap-1.5">
						<label for="height" class="text-xs font-medium text-gray-700">H</label>
						<input
							id="height"
							type="number"
							class="rounded-md border border-gray-300 bg-white px-2 py-1.5 text-[13px] text-gray-900 focus:border-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
							bind:value={height}
							onchange={() => handleSizeChange('height', height)}
							disabled={locked}
						/>
					</div>
				</div>

				<div class="mb-2 flex flex-col gap-1.5">
					<label for="rotation" class="text-xs font-medium text-gray-700">Rotation</label>
					<div class="flex items-center gap-2">
						<input
							id="rotation"
							type="range"
							min="0"
							max="360"
							class="h-1.5 flex-1 appearance-none rounded-full bg-gray-200 outline-none [&::-moz-range-thumb]:h-3.5 [&::-moz-range-thumb]:w-3.5 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-blue-500 [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:shadow-sm [&::-webkit-slider-thumb]:h-3.5 [&::-webkit-slider-thumb]:w-3.5 [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-blue-500 [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow-sm"
							bind:value={rotation}
							oninput={() => handleRotationChange(rotation)}
							disabled={locked}
						/>
						<span class="min-w-[40px] text-right font-mono text-xs text-gray-500">{rotation}°</span>
					</div>
				</div>

				<div class="mb-2 flex flex-col gap-1.5">
					<label for="opacity" class="text-xs font-medium text-gray-700">Opacity</label>
					<div class="flex items-center gap-2">
						<input
							id="opacity"
							type="range"
							min="0"
							max="100"
							class="h-1.5 flex-1 appearance-none rounded-full bg-gray-200 outline-none [&::-moz-range-thumb]:h-3.5 [&::-moz-range-thumb]:w-3.5 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-blue-500 [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:shadow-sm [&::-webkit-slider-thumb]:h-3.5 [&::-webkit-slider-thumb]:w-3.5 [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-blue-500 [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow-sm"
							bind:value={opacity}
							oninput={() => handleOpacityChange(opacity)}
							disabled={locked}
						/>
						<span class="min-w-[40px] text-right font-mono text-xs text-gray-500"
							>{Math.round(opacity)}%</span
						>
					</div>
				</div>
			</div>

			<!-- Text Properties -->
			{#if element.type === 'text' || element.type === 'sticky' || element.type === 'list'}
				<div class="mb-6">
					<div class="mb-3 text-xs font-semibold tracking-wide text-gray-500 uppercase">Text</div>

					<div class="mb-2 flex flex-col gap-1.5">
						<label for="font-family" class="text-xs font-medium text-gray-700">Font</label>
						<select
							id="font-family"
							class="rounded-md border border-gray-300 bg-white px-2 py-1.5 text-[13px] text-gray-900 focus:border-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
							bind:value={fontFamily}
							onchange={() => handleFontFamilyChange(fontFamily)}
							disabled={locked}
						>
							{#each fontFamilies as font}
								<option value={font}>{font}</option>
							{/each}
						</select>
					</div>

					<div class="mb-2 grid grid-cols-2 gap-2">
						<div class="flex flex-col gap-1.5">
							<label for="font-size" class="text-xs font-medium text-gray-700">Size</label>
							<input
								id="font-size"
								type="number"
								min="8"
								max="128"
								class="rounded-md border border-gray-300 bg-white px-2 py-1.5 text-[13px] text-gray-900 focus:border-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
								bind:value={fontSize}
								onchange={() => handleFontSizeChange(fontSize)}
								disabled={locked}
							/>
						</div>
						<div class="flex flex-col gap-1.5">
							<label for="font-weight" class="text-xs font-medium text-gray-700">Weight</label>
							<select
								id="font-weight"
								class="rounded-md border border-gray-300 bg-white px-2 py-1.5 text-[13px] text-gray-900 focus:border-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
								bind:value={fontWeight}
								onchange={() => handleFontWeightChange(fontWeight)}
								disabled={locked}
							>
								<option value="normal">Normal</option>
								<option value="bold">Bold</option>
								<option value="lighter">Light</option>
							</select>
						</div>
					</div>

					<div class="mb-2 flex flex-col gap-1.5">
						<label class="text-xs font-medium text-gray-700">Alignment</label>
						<div class="grid grid-cols-3 gap-1">
							<button
								class="flex cursor-pointer items-center justify-center rounded-md border border-gray-200 bg-white px-2 py-2 text-gray-500 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-50 {textAlign ===
								'left'
									? 'border-blue-500 bg-blue-50 text-blue-600'
									: ''}"
								onclick={() => handleTextAlignChange('left')}
								disabled={locked}
							>
								<AlignLeft size={16} />
							</button>
							<button
								class="flex cursor-pointer items-center justify-center rounded-md border border-gray-200 bg-white px-2 py-2 text-gray-500 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-50 {textAlign ===
								'center'
									? 'border-blue-500 bg-blue-50 text-blue-600'
									: ''}"
								onclick={() => handleTextAlignChange('center')}
								disabled={locked}
							>
								<AlignCenter size={16} />
							</button>
							<button
								class="flex cursor-pointer items-center justify-center rounded-md border border-gray-200 bg-white px-2 py-2 text-gray-500 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-50 {textAlign ===
								'right'
									? 'border-blue-500 bg-blue-50 text-blue-600'
									: ''}"
								onclick={() => handleTextAlignChange('right')}
								disabled={locked}
							>
								<AlignRight size={16} />
							</button>
						</div>
					</div>

					<div class="mb-3">
						<ColorPicker
							value={color}
							onChange={(value) => handleColorChange('color', value)}
							label="Text Color"
						/>
					</div>
				</div>
			{/if}

			<!-- Shape Properties -->
			{#if ['rectangle', 'ellipse', 'triangle', 'line', 'arrow'].includes(element.type)}
				<div class="mb-6">
					<div class="mb-3 text-xs font-semibold tracking-wide text-gray-500 uppercase">
						Fill & Stroke
					</div>

					<div class="mb-3">
						<ColorPicker
							value={backgroundColor}
							onChange={(value) => handleColorChange('backgroundColor', value)}
							label="Fill Color"
						/>
					</div>

					<div class="mb-3">
						<ColorPicker
							value={strokeColor}
							onChange={(value) => handleColorChange('strokeColor', value)}
							label="Stroke Color"
						/>
					</div>

					<div class="mb-2 flex flex-col gap-1.5">
						<label for="stroke-width" class="text-xs font-medium text-gray-700">Stroke Width</label>
						<div class="flex items-center gap-2">
							<input
								id="stroke-width"
								type="range"
								min="0"
								max="20"
								step="0.5"
								class="h-1.5 flex-1 appearance-none rounded-full bg-gray-200 outline-none [&::-moz-range-thumb]:h-3.5 [&::-moz-range-thumb]:w-3.5 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-blue-500 [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:shadow-sm [&::-webkit-slider-thumb]:h-3.5 [&::-webkit-slider-thumb]:w-3.5 [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-blue-500 [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow-sm"
								bind:value={strokeWidth}
								oninput={() => handleStrokeWidthChange(strokeWidth)}
								disabled={locked}
							/>
							<span class="min-w-[40px] text-right font-mono text-xs text-gray-500"
								>{strokeWidth}px</span
							>
						</div>
					</div>

					{#if element.type === 'rectangle'}
						<div class="mb-2 flex flex-col gap-1.5">
							<label for="border-radius" class="text-xs font-medium text-gray-700"
								>Corner Radius</label
							>
							<div class="flex items-center gap-2">
								<input
									id="border-radius"
									type="range"
									min="0"
									max="50"
									class="h-1.5 flex-1 appearance-none rounded-full bg-gray-200 outline-none [&::-moz-range-thumb]:h-3.5 [&::-moz-range-thumb]:w-3.5 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-blue-500 [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:shadow-sm [&::-webkit-slider-thumb]:h-3.5 [&::-webkit-slider-thumb]:w-3.5 [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-blue-500 [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow-sm"
									bind:value={borderRadius}
									oninput={() => handleBorderRadiusChange(borderRadius)}
									disabled={locked}
								/>
								<span class="min-w-[40px] text-right font-mono text-xs text-gray-500"
									>{borderRadius}px</span
								>
							</div>
						</div>
					{/if}
				</div>
			{/if}

			<!-- Drawing Properties -->
			{#if element.type === 'freehand'}
				<div class="mb-6">
					<div class="mb-3 text-xs font-semibold tracking-wide text-gray-500 uppercase">Stroke</div>

					<div class="mb-3">
						<ColorPicker
							value={strokeColor}
							onChange={(value) => handleColorChange('strokeColor', value)}
							label="Stroke Color"
						/>
					</div>

					<div class="mb-2 flex flex-col gap-1.5">
						<label for="stroke-width-freehand" class="text-xs font-medium text-gray-700"
							>Stroke Width</label
						>
						<div class="flex items-center gap-2">
							<input
								id="stroke-width-freehand"
								type="range"
								min="1"
								max="20"
								step="0.5"
								class="h-1.5 flex-1 appearance-none rounded-full bg-gray-200 outline-none [&::-moz-range-thumb]:h-3.5 [&::-moz-range-thumb]:w-3.5 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-blue-500 [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:shadow-sm [&::-webkit-slider-thumb]:h-3.5 [&::-webkit-slider-thumb]:w-3.5 [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-blue-500 [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow-sm"
								bind:value={strokeWidth}
								oninput={() => handleStrokeWidthChange(strokeWidth)}
								disabled={locked}
							/>
							<span class="min-w-[40px] text-right font-mono text-xs text-gray-500"
								>{strokeWidth}px</span
							>
						</div>
					</div>
				</div>
			{/if}

			<!-- Layer Controls -->
			<div class="mb-0">
				<div class="mb-3 text-xs font-semibold tracking-wide text-gray-500 uppercase">Layers</div>

				<div class="grid grid-cols-2 gap-2">
					<button
						class="flex cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-50"
						onclick={handleBringToFront}
						disabled={locked}
						title="Bring to Front"
					>
						<ChevronsUp size={16} />
						<span>To Front</span>
					</button>
					<button
						class="flex cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-50"
						onclick={handleBringForward}
						disabled={locked}
						title="Bring Forward"
					>
						<ChevronUp size={16} />
						<span>Forward</span>
					</button>
					<button
						class="flex cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-50"
						onclick={handleSendBackward}
						disabled={locked}
						title="Send Backward"
					>
						<ChevronDown size={16} />
						<span>Backward</span>
					</button>
					<button
						class="flex cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-50"
						onclick={handleSendToBack}
						disabled={locked}
						title="Send to Back"
					>
						<ChevronsDown size={16} />
						<span>To Back</span>
					</button>
				</div>
			</div>
		</div>
	</div>
{:else}
	<div class="flex h-full w-[280px] flex-col overflow-hidden border-l border-gray-200 bg-white">
		<div class="flex items-center justify-between border-b border-gray-200 px-4 py-4">
			<h3 class="m-0 text-sm font-semibold text-gray-900">Multiple Selection</h3>
		</div>
		<div class="flex-1 overflow-y-auto px-4 py-4">
			<p class="mb-4 rounded-md bg-gray-100 px-3 py-3 text-center text-[13px] text-gray-500">
				{selectedCount} elements selected
			</p>

			<!-- Layer Controls for multiple selection -->
			<div class="mb-0">
				<div class="mb-3 text-xs font-semibold tracking-wide text-gray-500 uppercase">Layers</div>
				<div class="grid grid-cols-2 gap-2">
					<button
						class="flex cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900"
						onclick={handleBringToFront}
						title="Bring to Front"
					>
						<ChevronsUp size={16} />
						<span>To Front</span>
					</button>
					<button
						class="flex cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900"
						onclick={handleBringForward}
						title="Bring Forward"
					>
						<ChevronUp size={16} />
						<span>Forward</span>
					</button>
					<button
						class="flex cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900"
						onclick={handleSendBackward}
						title="Send Backward"
					>
						<ChevronDown size={16} />
						<span>Backward</span>
					</button>
					<button
						class="flex cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900"
						onclick={handleSendToBack}
						title="Send to Back"
					>
						<ChevronsDown size={16} />
						<span>To Back</span>
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
