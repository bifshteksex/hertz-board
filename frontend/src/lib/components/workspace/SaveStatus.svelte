<script lang="ts">
	import { autosaveStore } from '$lib/stores/autosave.svelte';
	import { Check, Cloud, CloudOff, Loader2, AlertCircle } from 'lucide-svelte';
	import { formatDistanceToNow } from 'date-fns';
	import { ru } from 'date-fns/locale';

	// Reactive state from autosave store
	const status = $derived(autosaveStore.status);
	const pendingCount = $derived(autosaveStore.pendingCount);
	const lastSaveTime = $derived(autosaveStore.lastSaveTime);
	const lastError = $derived(autosaveStore.lastError);

	// Computed properties
	const statusText = $derived.by(() => {
		if (status === 'saving') return 'Saving...';
		if (status === 'saved') {
			if (lastSaveTime) {
				return `Saved ${formatDistanceToNow(lastSaveTime, { addSuffix: true, locale: ru })}`;
			}
			return 'Saved';
		}
		if (status === 'error') return 'Save failed';
		if (pendingCount > 0) return `${pendingCount} unsaved changes`;
		return 'All changes saved';
	});

	const statusIcon = $derived.by(() => {
		switch (status) {
			case 'saving':
				return Loader2;
			case 'saved':
				return Check;
			case 'error':
				return AlertCircle;
			default:
				return pendingCount > 0 ? Cloud : CloudOff;
		}
	});

	const statusColor = $derived.by(() => {
		switch (status) {
			case 'saving':
				return 'text-blue-600';
			case 'saved':
				return 'text-green-600';
			case 'error':
				return 'text-red-600';
			default:
				return pendingCount > 0 ? 'text-yellow-600' : 'text-gray-500';
		}
	});

	const shouldAnimate = $derived(status === 'saving');
</script>

<div class="flex items-center gap-2 text-sm {statusColor}">
	<svelte:component this={statusIcon} size={16} class={shouldAnimate ? 'animate-spin' : ''} />
	<span>{statusText}</span>
	{#if lastError && status === 'error'}
		<button
			class="ml-2 text-xs underline hover:no-underline"
			onclick={() => alert(lastError)}
			title={lastError}
		>
			Details
		</button>
	{/if}
</div>
