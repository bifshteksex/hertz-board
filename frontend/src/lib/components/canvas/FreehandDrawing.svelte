<script lang="ts">
	/**
	 * FreehandDrawing - компонент для рисования от руки
	 */
	import { getStroke } from 'perfect-freehand';

	interface Point {
		x: number;
		y: number;
		pressure?: number;
	}

	interface Props {
		points: Point[];
		color?: string;
		width?: number;
		opacity?: number;
	}

	let { points, color = '#000000', width = 2, opacity = 1 }: Props = $props();

	// Генерируем SVG path из точек используя perfect-freehand
	const pathData = $derived.by(() => {
		if (points.length < 2) return '';

		const stroke = getStroke(points, {
			size: width * 2,
			thinning: 0.5,
			smoothing: 0.5,
			streamline: 0.5,
			simulatePressure: true
		});

		if (stroke.length === 0) return '';

		const d = stroke.reduce((acc, [x0, y0], i, arr) => {
			const [x1, y1] = arr[(i + 1) % arr.length];
			if (i === 0) {
				return `M ${x0},${y0} Q ${x1},${y1}`;
			}
			if (i === arr.length - 1) {
				return `${acc} ${x0},${y0} Z`;
			}
			return `${acc} ${x0},${y0} T ${x1},${y1}`;
		}, '');

		return d;
	});
</script>

{#if pathData}
	<path d={pathData} fill={color} {opacity} pointer-events="none" />
{/if}

<style>
	/* Styles are inline for this component */
</style>
