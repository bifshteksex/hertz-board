<script lang="ts">
	/**
	 * ShapeCreator - компонент для создания фигур drag-to-create
	 */
	import type { Tool } from '$lib/stores/canvas.svelte';

	interface Props {
		tool: Tool;
		startX: number;
		startY: number;
		currentX: number;
		currentY: number;
		shiftPressed?: boolean;
	}

	let { tool, startX, startY, currentX, currentY, shiftPressed = false }: Props = $props();

	// Вычисляем размеры и позицию
	const x = $derived(Math.min(startX, currentX));
	const y = $derived(Math.min(startY, currentY));
	const width = $derived(Math.abs(currentX - startX));
	const height = $derived(Math.abs(currentY - startY));

	// Для сохранения пропорций при shift
	const finalWidth = $derived(shiftPressed ? Math.max(width, height) : width);
	const finalHeight = $derived(shiftPressed ? Math.max(width, height) : height);

	// Стили preview
	const previewStyle = {
		fill: 'rgba(59, 130, 246, 0.1)',
		stroke: '#3b82f6',
		strokeWidth: 2,
		strokeDasharray: '5 5'
	};
</script>

{#if tool === 'rectangle'}
	<rect
		{x}
		{y}
		width={finalWidth}
		height={finalHeight}
		fill={previewStyle.fill}
		stroke={previewStyle.stroke}
		stroke-width={previewStyle.strokeWidth}
		stroke-dasharray={previewStyle.strokeDasharray}
		pointer-events="none"
	/>
{:else if tool === 'ellipse'}
	<ellipse
		cx={x + finalWidth / 2}
		cy={y + finalHeight / 2}
		rx={finalWidth / 2}
		ry={finalHeight / 2}
		fill={previewStyle.fill}
		stroke={previewStyle.stroke}
		stroke-width={previewStyle.strokeWidth}
		stroke-dasharray={previewStyle.strokeDasharray}
		pointer-events="none"
	/>
{:else if tool === 'triangle'}
	<polygon
		points="{x + finalWidth / 2},{y} {x + finalWidth},{y + finalHeight} {x},{y + finalHeight}"
		fill={previewStyle.fill}
		stroke={previewStyle.stroke}
		stroke-width={previewStyle.strokeWidth}
		stroke-dasharray={previewStyle.strokeDasharray}
		pointer-events="none"
	/>
{:else if tool === 'line'}
	<line
		x1={startX}
		y1={startY}
		x2={currentX}
		y2={currentY}
		stroke={previewStyle.stroke}
		stroke-width={previewStyle.strokeWidth}
		stroke-dasharray={previewStyle.strokeDasharray}
		pointer-events="none"
	/>
{:else if tool === 'arrow'}
	<defs>
		<marker
			id="preview-arrowhead"
			markerWidth="10"
			markerHeight="10"
			refX="9"
			refY="3"
			orient="auto"
			markerUnits="strokeWidth"
		>
			<path d="M0,0 L0,6 L9,3 z" fill={previewStyle.stroke} />
		</marker>
	</defs>
	<line
		x1={startX}
		y1={startY}
		x2={currentX}
		y2={currentY}
		stroke={previewStyle.stroke}
		stroke-width={previewStyle.strokeWidth}
		stroke-dasharray={previewStyle.strokeDasharray}
		marker-end="url(#preview-arrowhead)"
		pointer-events="none"
	/>
{/if}

<style>
	/* Styles are inline for this component */
</style>
