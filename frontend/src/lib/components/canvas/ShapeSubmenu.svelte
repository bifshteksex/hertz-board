<script lang="ts">
	import { Square, Circle, Triangle, Minus, ArrowRight } from 'lucide-svelte';
	import { i18n } from '$lib/stores/i18n.svelte';

	interface Props {
		onSelect: (_shape: 'rectangle' | 'ellipse' | 'triangle' | 'line' | 'arrow') => void;
	}

	let { onSelect }: Props = $props();

	const shapes = $derived([
		{ id: 'rectangle' as const, icon: Square, label: i18n.t('canvas.toolbar.shapes.rectangle') },
		{ id: 'ellipse' as const, icon: Circle, label: i18n.t('canvas.toolbar.shapes.ellipse') },
		{ id: 'triangle' as const, icon: Triangle, label: i18n.t('canvas.toolbar.shapes.triangle') },
		{ id: 'line' as const, icon: Minus, label: i18n.t('canvas.toolbar.shapes.line') },
		{ id: 'arrow' as const, icon: ArrowRight, label: i18n.t('canvas.toolbar.shapes.arrow') }
	]);

	function handleSelect(shape: (typeof shapes)[number]['id']) {
		onSelect(shape);
	}
</script>

<div
	class="absolute top-[calc(100%+4px)] left-0 z-[1000] flex min-w-[140px] flex-col gap-0.5 rounded-lg border border-gray-200 bg-white p-1 shadow-md"
>
	{#each shapes as shape}
		{@const Icon = shape.icon}
		<button
			class="flex w-full items-center gap-2 rounded-md border-none bg-transparent px-3 py-2 text-left text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 active:bg-gray-200"
			onclick={() => handleSelect(shape.id)}
			title={shape.label}
			aria-label={shape.label}
		>
			<Icon size={18} />
			<span class="text-sm font-medium">{shape.label}</span>
		</button>
	{/each}
</div>
