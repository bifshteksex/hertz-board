<script lang="ts">
	import { List, ListOrdered, CheckSquare } from 'lucide-svelte';
	import { i18n } from '$lib/stores/i18n.svelte';

	interface Props {
		onSelect: (_listType: 'bullet' | 'numbered' | 'checkbox') => void;
	}

	let { onSelect }: Props = $props();

	const listTypes = $derived([
		{ id: 'bullet' as const, icon: List, label: i18n.t('canvas.toolbar.lists.bullet') },
		{ id: 'numbered' as const, icon: ListOrdered, label: i18n.t('canvas.toolbar.lists.numbered') },
		{ id: 'checkbox' as const, icon: CheckSquare, label: i18n.t('canvas.toolbar.lists.checkbox') }
	]);

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
			class="flex w-full items-center gap-2 rounded-md border-none bg-transparent px-3 py-2 text-left text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 active:bg-gray-200"
			onclick={() => handleSelect(listType.id)}
			title={listType.label}
			aria-label={listType.label}
		>
			<Icon size={18} />
			<span class="text-sm font-medium">{listType.label}</span>
		</button>
	{/each}
</div>
