<script lang="ts">
	import { autosaveStore } from '$lib/stores/autosave.svelte';
	import { i18n } from '$lib/stores/i18n.svelte';
	import { Check, Cloud, CloudOff, Loader2, AlertCircle } from 'lucide-svelte';
	import { formatDistanceToNow } from 'date-fns';
	import { ru, zhCN, enUS } from 'date-fns/locale';

	// Reactive state from autosave store
	const status = $derived(autosaveStore.status);
	const pendingCount = $derived(autosaveStore.pendingCount);
	const lastSaveTime = $derived(autosaveStore.lastSaveTime);
	const lastError = $derived(autosaveStore.lastError);

	// Get date-fns locale based on current language
	const dateLocale = $derived.by(() => {
		const currentLocale = i18n.currentLocale;
		if (currentLocale === 'ru') return ru;
		if (currentLocale === 'zh') return zhCN;
		return enUS;
	});

	// Computed properties
	const statusText = $derived.by(() => {
		if (status === 'saving') return i18n.t('canvas.saveStatus.saving');
		if (status === 'saved') {
			if (lastSaveTime) {
				const time = formatDistanceToNow(lastSaveTime, { addSuffix: true, locale: dateLocale });
				return i18n.t('canvas.saveStatus.savedAgo', { time });
			}
			return i18n.t('canvas.saveStatus.saved');
		}
		if (status === 'error') return i18n.t('canvas.saveStatus.saveFailed');
		if (pendingCount > 0)
			return i18n.t('canvas.saveStatus.unsavedChanges', { count: pendingCount.toString() });
		return i18n.t('canvas.saveStatus.allSaved');
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
	{#if statusIcon}
		{@const StatusIcon = statusIcon}
		<StatusIcon size={16} class={shouldAnimate ? 'animate-spin' : ''} />
	{/if}
	<span>{statusText}</span>
	{#if lastError && status === 'error'}
		<button
			class="ml-2 text-xs underline hover:no-underline"
			onclick={() => alert(lastError)}
			title={lastError}
		>
			{i18n.t('canvas.saveStatus.details')}
		</button>
	{/if}
</div>
