<script lang="ts">
	import { canvasStore, MIN_ZOOM, MAX_ZOOM, ZOOM_STEP } from '$lib/stores/canvas.svelte';
	import { Maximize2 } from 'lucide-svelte';
	import IconZoomIn from '$components/icons/IconZoomIn.svelte';
	import IconZoomOut from '$components/icons/IconZoomOut.svelte';

	const viewport = $derived(canvasStore.viewport);
	const zoomPercent = $derived(Math.round(viewport.zoom * 100));
	const canZoomIn = $derived(viewport.zoom < MAX_ZOOM);
	const canZoomOut = $derived(viewport.zoom > MIN_ZOOM);

	function handleZoomIn() {
		canvasStore.zoom(ZOOM_STEP);
	}

	function handleZoomOut() {
		canvasStore.zoom(-ZOOM_STEP);
	}

	function handleFitToScreen() {
		// Получаем размеры canvas container
		const container = document.querySelector('.canvas-container');
		if (container) {
			const rect = container.getBoundingClientRect();
			canvasStore.fitToScreen(rect.width, rect.height);
		}
	}

	function handleResetZoom() {
		canvasStore.resetZoom();
	}
</script>

<div class="zoom-controls">
	<button
		class="zoom-btn"
		onclick={handleZoomOut}
		disabled={!canZoomOut}
		title="Zoom Out (Ctrl + -)"
		aria-label="Zoom out"
	>
		<IconZoomOut size={18} />
	</button>

	<button
		class="zoom-percent"
		onclick={handleResetZoom}
		title="Reset Zoom (Ctrl + 0)"
		aria-label="Reset zoom to 100%"
	>
		{zoomPercent}%
	</button>

	<button
		class="zoom-btn"
		onclick={handleZoomIn}
		disabled={!canZoomIn}
		title="Zoom In (Ctrl + +)"
		aria-label="Zoom in"
	>
		<IconZoomIn size={18} />
	</button>

	<div class="separator"></div>

	<button
		class="zoom-btn"
		onclick={handleFitToScreen}
		title="Fit to Screen"
		aria-label="Fit all elements to screen"
	>
		<Maximize2 size={18} />
	</button>
</div>

<style>
	.zoom-controls {
		position: absolute;
		bottom: 24px;
		right: 24px;
		display: flex;
		align-items: center;
		gap: 4px;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		padding: 4px;
		box-shadow:
			0 4px 6px -1px rgb(0 0 0 / 0.1),
			0 2px 4px -2px rgb(0 0 0 / 0.1);
		z-index: 10;
	}

	.zoom-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		border: none;
		background: transparent;
		border-radius: 6px;
		cursor: pointer;
		color: #374151;
		transition: all 0.15s;
	}

	.zoom-btn:hover:not(:disabled) {
		background: #f3f4f6;
		color: #111827;
	}

	.zoom-btn:active:not(:disabled) {
		background: #e5e7eb;
	}

	.zoom-btn:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.zoom-percent {
		min-width: 60px;
		height: 32px;
		padding: 0 8px;
		border: none;
		background: transparent;
		border-radius: 6px;
		cursor: pointer;
		font-size: 13px;
		font-weight: 500;
		color: #374151;
		transition: all 0.15s;
	}

	.zoom-percent:hover {
		background: #f3f4f6;
		color: #111827;
	}

	.zoom-percent:active {
		background: #e5e7eb;
	}

	.separator {
		width: 1px;
		height: 24px;
		background: #e5e7eb;
		margin: 0 4px;
	}
</style>
