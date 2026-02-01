<script lang="ts">
	import { collaborationStore } from '$lib/stores/collaboration.svelte';
	import { Wifi, WifiOff, Loader2, AlertCircle } from 'lucide-svelte';

	const isConnected = $derived(collaborationStore.isConnected);
	const isSyncing = $derived(collaborationStore.isSyncing);
	const error = $derived(collaborationStore.error);

	const statusText = $derived.by(() => {
		if (error) return 'Connection Error';
		if (isSyncing) return 'Syncing...';
		if (isConnected) return 'Connected';
		return 'Disconnected';
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

<div class="connection-status" style="--status-color: {statusColor}">
	<!-- Icon -->
	<div class="status-icon">
		{#if error}
			<AlertCircle size={16} />
		{:else if isSyncing}
			<Loader2 size={16} class="spinning" />
		{:else if isConnected}
			<Wifi size={16} />
		{:else}
			<WifiOff size={16} />
		{/if}
	</div>

	<!-- Status text -->
	<div class="status-text">{statusText}</div>

	<!-- Retry button (only on error) -->
	{#if error}
		<button class="retry-button" onclick={handleRetry} title="Reconnect">Retry</button>
	{/if}

	<!-- Status indicator dot -->
	<div class="status-dot"></div>
</div>

<style>
	.connection-status {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 6px 12px;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		font-size: 13px;
		color: #374151;
	}

	.status-icon {
		display: flex;
		align-items: center;
		color: var(--status-color);
	}

	.status-text {
		font-weight: 500;
		color: #111827;
	}

	.status-dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background-color: var(--status-color);
		animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
	}

	.retry-button {
		padding: 4px 8px;
		background: #ef4444;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 12px;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.retry-button:hover {
		background: #dc2626;
	}

	:global(.spinning) {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.5;
		}
	}
</style>
