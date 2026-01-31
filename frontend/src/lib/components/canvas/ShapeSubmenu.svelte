<script lang="ts">
	import { Square, Circle, Triangle, Minus, ArrowRight } from 'lucide-svelte';

	interface Props {
		onSelect: (_shape: 'rectangle' | 'ellipse' | 'triangle' | 'line' | 'arrow') => void;
	}

	let { onSelect }: Props = $props();

	const shapes = [
		{ id: 'rectangle' as const, icon: Square, label: 'Rectangle' },
		{ id: 'ellipse' as const, icon: Circle, label: 'Circle' },
		{ id: 'triangle' as const, icon: Triangle, label: 'Triangle' },
		{ id: 'line' as const, icon: Minus, label: 'Line' },
		{ id: 'arrow' as const, icon: ArrowRight, label: 'Arrow' }
	];

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
			class="flex w-full cursor-pointer items-center gap-2 rounded-md border-none bg-transparent px-3 py-2 text-left text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 active:bg-gray-200"
			onclick={() => handleSelect(shape.id)}
			title={shape.label}
			aria-label={shape.label}
		>
			<Icon size={18} />
			<span class="text-sm font-medium">{shape.label}</span>
		</button>
	{/each}
</div>
