<script lang="ts">
	import type { SelectionBox, Viewport } from '$lib/stores/canvas.svelte';

	interface Props {
		box: SelectionBox;
		viewport: Viewport;
	}

	let { box }: Props = $props();

	// box.start и box.current уже в canvas координатах
	// Преобразуем их обратно в screen координаты только для отображения
	const screenRect = $derived(() => {
		const startX = box.start.x;
		const startY = box.start.y;
		const currentX = box.current.x;
		const currentY = box.current.y;

		const x = Math.min(startX, currentX);
		const y = Math.min(startY, currentY);
		const width = Math.abs(currentX - startX);
		const height = Math.abs(currentY - startY);

		return { x, y, width, height };
	});
</script>

<rect
	class="selection-box"
	x={screenRect().x}
	y={screenRect().y}
	width={screenRect().width}
	height={screenRect().height}
	fill="rgba(59, 130, 246, 0.1)"
	stroke="#3b82f6"
	stroke-width="1"
	stroke-dasharray="4 4"
	pointer-events="none"
/>

<style>
	.selection-box {
		animation: dash 0.5s linear infinite;
	}

	@keyframes dash {
		to {
			stroke-dashoffset: -8;
		}
	}
</style>
