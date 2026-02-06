<script lang="ts">
	/**
	 * FreehandDrawing - компонент для рисования от руки
	 *
	 * NOTE: Eraser tool can be implemented by:
	 * 1. Detecting intersection between eraser path and existing strokes
	 * 2. Splitting strokes at intersection points
	 * 3. Creating new segments from non-erased portions
	 * For MVP, users can delete entire freehand elements via Delete key
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

		// Конвертируем points в формат для getStroke: [x, y, pressure]
		const strokePoints = points.map((p) => [p.x, p.y, p.pressure ?? 0.5]);

		const stroke = getStroke(strokePoints, {
			size: width * 2,
			thinning: 0.5,
			smoothing: 0.5,
			streamline: 0.5,
			simulatePressure: true
		});

		if (stroke.length === 0) return '';

		// Создаем SVG path из stroke точек
		let d = '';
		for (let i = 0; i < stroke.length; i++) {
			const [x, y] = stroke[i];
			if (i === 0) {
				d = `M ${x},${y}`;
			} else {
				d += ` L ${x},${y}`;
			}
		}

		// Закрываем path
		if (stroke.length > 0) {
			d += ' Z';
		}

		return d;
	});
</script>

{#if pathData}
	<path d={pathData} fill={color} {opacity} />
{/if}
