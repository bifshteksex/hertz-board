<script lang="ts">
	import { List, ListOrdered, CheckSquare } from 'lucide-svelte';

	interface Props {
		onSelect: (_listType: 'bullet' | 'numbered' | 'checkbox') => void;
	}

	let { onSelect }: Props = $props();

	const listTypes = [
		{ id: 'bullet' as const, icon: List, label: 'Bullet List' },
		{ id: 'numbered' as const, icon: ListOrdered, label: 'Numbered List' },
		{ id: 'checkbox' as const, icon: CheckSquare, label: 'Checkbox List' }
	];

	function handleSelect(listType: (typeof listTypes)[number]['id']) {
		onSelect(listType);
	}
</script>

<div
	class="absolute top-[calc(100%+4px)] left-0 z-[1000] flex min-w-[160px] flex-col gap-0.5 rounded-lg border border-gray-200 bg-white p-1 shadow-md"
>
	{#each listTypes as listType}
		{@const Icon = listType.icon}
		<button
			class="flex w-full cursor-pointer items-center gap-2 rounded-md border-none bg-transparent px-3 py-2 text-left text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 active:bg-gray-200"
			onclick={() => handleSelect(listType.id)}
			title={listType.label}
			aria-label={listType.label}
		>
			<Icon size={18} />
			<span class="text-sm font-medium">{listType.label}</span>
		</button>
	{/each}
</div>
