<script lang="ts">
	import type { CanvasElement } from '$lib/types/api';
	import type { Viewport } from '$lib/stores/canvas.svelte';

	interface Props {
		elements: CanvasElement[];
		viewport: Viewport;
	}

	let { elements, viewport }: Props = $props();

	// Вычисляем bounding box всех выделенных элементов
	const bounds = $derived(() => {
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

		// Преобразуем в screen координаты
		const topLeft = {
			x: minX * viewport.zoom + viewport.x,
			y: minY * viewport.zoom + viewport.y
		};
		const bottomRight = {
			x: maxX * viewport.zoom + viewport.x,
			y: maxY * viewport.zoom + viewport.y
		};

		return {
			x: topLeft.x,
			y: topLeft.y,
			width: bottomRight.x - topLeft.x,
			height: bottomRight.y - topLeft.y
		};
	});

	const handleSize = 8; // размер handle в пикселях
	const rotateHandleDistance = 30; // расстояние rotation handle от верха

	// Позиции handles
	const handles = $derived(
		bounds()
			? (() => {
					const b = bounds()!;
					const hw = handleSize / 2;

					return {
						// Corner handles
						topLeft: { x: b.x - hw, y: b.y - hw, cursor: 'nwse-resize' },
						topRight: { x: b.x + b.width - hw, y: b.y - hw, cursor: 'nesw-resize' },
						bottomLeft: { x: b.x - hw, y: b.y + b.height - hw, cursor: 'nesw-resize' },
						bottomRight: { x: b.x + b.width - hw, y: b.y + b.height - hw, cursor: 'nwse-resize' },

						// Edge handles
						top: { x: b.x + b.width / 2 - hw, y: b.y - hw, cursor: 'ns-resize' },
						right: { x: b.x + b.width - hw, y: b.y + b.height / 2 - hw, cursor: 'ew-resize' },
						bottom: { x: b.x + b.width / 2 - hw, y: b.y + b.height - hw, cursor: 'ns-resize' },
						left: { x: b.x - hw, y: b.y + b.height / 2 - hw, cursor: 'ew-resize' },

						// Rotate handle
						rotate: {
							x: b.x + b.width / 2 - hw,
							y: b.y - rotateHandleDistance - hw,
							cursor: 'grab'
						}
					};
				})()
			: null
	);
</script>

{#if bounds()}
	<g class="selection-handles" pointer-events="all">
		<!-- Selection border -->
		<rect
			x={bounds()!.x}
			y={bounds()!.y}
			width={bounds()!.width}
			height={bounds()!.height}
			fill="none"
			stroke="#3b82f6"
			stroke-width="2"
			pointer-events="none"
		/>

		<!-- Rotation handle line -->
		<line
			x1={bounds()!.x + bounds()!.width / 2}
			y1={bounds()!.y}
			x2={bounds()!.x + bounds()!.width / 2}
			y2={bounds()!.y - rotateHandleDistance}
			stroke="#3b82f6"
			stroke-width="1"
			pointer-events="none"
		/>

		<!-- Corner handles -->
		{#if handles}
			{#each Object.entries(handles) as [name, handle]}
				<rect
					class="handle handle-{name}"
					x={handle.x}
					y={handle.y}
					width={handleSize}
					height={handleSize}
					fill="white"
					stroke="#3b82f6"
					stroke-width="2"
					rx="2"
					style:cursor={handle.cursor}
					data-handle={name}
				/>
			{/each}

			<!-- Rotate handle (circle) -->
			<circle
				class="handle handle-rotate"
				cx={handles.rotate.x + handleSize / 2}
				cy={handles.rotate.y + handleSize / 2}
				r={handleSize / 2}
				fill="white"
				stroke="#3b82f6"
				stroke-width="2"
				style:cursor={handles.rotate.cursor}
				data-handle="rotate"
			/>
		{/if}
	</g>
{/if}

<style>
	.handle {
		pointer-events: all;
	}

	.handle:hover {
		fill: #3b82f6;
	}
</style>
