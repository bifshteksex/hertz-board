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

<div
	class="absolute right-6 bottom-6 z-10 flex items-center gap-1 border border-gray-200 bg-white p-1 shadow-md dark:border-gray-700 dark:bg-gray-800"
>
	<button
		class="flex size-8 items-center justify-center border-none bg-transparent text-gray-700 transition-all duration-150 hover:enabled:bg-gray-100 hover:enabled:text-gray-900 active:enabled:bg-gray-200 disabled:cursor-not-allowed disabled:opacity-40 dark:text-gray-300 dark:hover:enabled:bg-gray-700 dark:hover:enabled:text-gray-100"
		onclick={handleZoomOut}
		disabled={!canZoomOut}
		title="Zoom Out (Ctrl + -)"
		aria-label="Zoom out"
	>
		<IconZoomOut size={18} />
	</button>

	<button
		class="h-8 min-w-15 border-none bg-transparent px-2 text-[13px] font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 active:bg-gray-200 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
		onclick={handleResetZoom}
		title="Reset Zoom (Ctrl + 0)"
		aria-label="Reset zoom to 100%"
	>
		{zoomPercent}%
	</button>

	<button
		class="flex size-8 items-center justify-center border-none bg-transparent text-gray-700 transition-all duration-150 hover:enabled:bg-gray-100 hover:enabled:text-gray-900 active:enabled:bg-gray-200 disabled:cursor-not-allowed disabled:opacity-40 dark:text-gray-300 dark:hover:enabled:bg-gray-700 dark:hover:enabled:text-gray-100"
		onclick={handleZoomIn}
		disabled={!canZoomIn}
		title="Zoom In (Ctrl + +)"
		aria-label="Zoom in"
	>
		<IconZoomIn size={18} />
	</button>

	<div class="mx-1 h-6 w-px bg-gray-200 dark:bg-gray-700"></div>

	<button
		class="flex size-8 items-center justify-center border-none bg-transparent text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 active:bg-gray-200 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100"
		onclick={handleFitToScreen}
		title="Fit to Screen"
		aria-label="Fit all elements to screen"
	>
		<Maximize2 size={18} />
	</button>
</div>
