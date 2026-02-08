<script lang="ts">
	interface Props {
		x: number;
		y: number;
		size: number;
		visible: boolean;
		isEraser?: boolean;
	}

	let { x, y, size, visible, isEraser = false }: Props = $props();
</script>

{#if visible}
	<div
		class="pointer-events-none fixed z-10000 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 transition-[width,height] duration-100 ease-in-out {isEraser
			? 'border-red-500/70 bg-red-500/10 shadow-[0_0_0_1px_rgba(255,100,100,0.4)]'
			: 'border-black/50 shadow-[0_0_0_1px_rgba(255,255,255,0.8)] dark:border-white/60 dark:shadow-[0_0_0_1px_rgba(0,0,0,0.5)]'}"
		style:left="{x}px"
		style:top="{y}px"
		style:width="{size}px"
		style:height="{size}px"
		role="presentation"
		aria-hidden="true"
	></div>
{/if}
