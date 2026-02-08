<script lang="ts">
	import { collaborationStore } from '$lib/stores/collaboration.svelte';
	import { Wifi, WifiOff, Loader2, AlertCircle } from 'lucide-svelte';
	import { i18n } from '$lib/stores/i18n.svelte';

	const t = $derived(i18n.messages);
	const isConnected = $derived(collaborationStore.isConnected);
	const isSyncing = $derived(collaborationStore.isSyncing);
	const error = $derived(collaborationStore.error);

	const statusText = $derived.by(() => {
		if (error) return t.canvas.connectionStatus.error;
		if (isSyncing) return t.canvas.connectionStatus.syncing;
		if (isConnected) return t.canvas.connectionStatus.connected;
		return t.canvas.connectionStatus.disconnected;
	});

	const statusColor = $derived.by(() => {
		if (error) return '#EF4444'; // Red
		if (isSyncing) return '#F59E0B'; // Amber
		if (isConnected) return '#10B981'; // Green
		return '#6B7280'; // Gray
	});

	function handleRetry() {
		// The workspace page should handle reconnection
		window.location.reload();
	}
</script>

<div
	class="flex items-center gap-2 border border-gray-200 bg-white px-3 py-1.5 text-[13px] text-gray-700 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
	style="--status-color: {statusColor}"
>
	<!-- Icon -->
	<div class="flex items-center" style="color: var(--status-color)">
		{#if error}
			<AlertCircle size={30} />
		{:else if isSyncing}
			<Loader2 size={30} class="animate-spin" />
		{:else if isConnected}
			<Wifi size={30} />
		{:else}
			<WifiOff size={30} />
		{/if}
	</div>

	<!-- Status text -->
	<div class="font-medium text-gray-900 dark:text-gray-100">{statusText}</div>

	<!-- Retry button (only on error) -->
	{#if error}
		<button
			class="rounded border-none bg-red-600 px-2 py-1 text-xs font-medium text-white transition-colors hover:bg-red-700"
			onclick={handleRetry}
			title={t.canvas.connectionStatus.reconnect}>{t.canvas.connectionStatus.retry}</button
		>
	{/if}

	<!-- Status indicator dot -->
	<div
		class="h-2 w-2 animate-pulse rounded-full"
		style="background-color: var(--status-color)"
	></div>
</div>
