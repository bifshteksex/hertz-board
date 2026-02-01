<script lang="ts">
	import type { Viewport } from '$lib/stores/canvas.svelte';
	import { canvasStore } from '$lib/stores/canvas.svelte';

	interface Props {
		viewport: Viewport;
	}

	let { viewport }: Props = $props();

	// Вычисляем динамический размер сетки в зависимости от zoom
	const gridSize = $derived(() => {
		const baseSize = canvasStore.gridSize;
		const zoom = viewport.zoom;

		// При малом zoom увеличиваем шаг сетки
		if (zoom < 0.25) return baseSize * 4;
		if (zoom < 0.5) return baseSize * 2;
		if (zoom > 2) return baseSize / 2;

		return baseSize;
	});

	// Вычисляем размер точек - фиксированный размер в screen space
	// Делим на zoom чтобы точки оставались одинакового размера на экране
	const dotSize = $derived(() => {
		const zoom = viewport.zoom;
		const baseSize = 1.5; // Базовый размер в screen pixels

		// Делим на zoom чтобы компенсировать масштабирование
		return baseSize / zoom;
	});

	// ID для pattern (уникальный)
	const patternId = 'grid-pattern';
</script>

<defs>
	<!-- Определяем pattern для сетки -->
	<pattern
		id={patternId}
		x="0"
		y="0"
		width={gridSize()}
		height={gridSize()}
		patternUnits="userSpaceOnUse"
	>
		<circle cx={gridSize() / 2} cy={gridSize() / 2} r={dotSize()} fill="#cbd5e1" opacity="0.5" />
	</pattern>
</defs>

<!-- Применяем pattern к большому прямоугольнику -->
<rect
	class="canvas-grid"
	x="-10000"
	y="-10000"
	width="20000"
	height="20000"
	fill="url(#{patternId})"
/>

<style>
	.canvas-grid {
		pointer-events: none;
	}
</style>
