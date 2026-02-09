<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { WifiOff, Wifi, Cloud, CloudOff } from 'lucide-svelte';
	import { isTauri, saveOfflineWorkspace } from '$lib/services/tauri';
	import { canvasStore } from '$lib/stores/canvas.svelte';

	let isOnline = $state(navigator.onLine);
	let lastSaved = $state<Date | null>(null);
	let autoSaveTimer: number | undefined;

	interface Props {
		workspaceId?: string;
		workspaceName?: string;
	}

	let { workspaceId, workspaceName }: Props = $props();

	function updateOnlineStatus() {
		isOnline = navigator.onLine;
	}

	async function saveToLocal() {
		if (!isTauri() || !workspaceId || !workspaceName) return;

		try {
			const elements = canvasStore.elements;
			await saveOfflineWorkspace(workspaceId, workspaceName, elements);
			lastSaved = new Date();
		} catch (error) {
			console.error('Failed to save offline:', error);
		}
	}

	function setupAutoSave() {
		if (!isTauri()) return;

		// Auto-save every 30 seconds when offline
		autoSaveTimer = window.setInterval(() => {
			if (!isOnline) {
				saveToLocal();
			}
		}, 30000);
	}

	function formatTimeSince(date: Date): string {
		const seconds = Math.floor((Date.now() - date.getTime()) / 1000);

		if (seconds < 60) return 'just now';
		if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
		if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
		return `${Math.floor(seconds / 86400)}d ago`;
	}

	onMount(() => {
		window.addEventListener('online', updateOnlineStatus);
		window.addEventListener('offline', updateOnlineStatus);
		setupAutoSave();

		// Save immediately when going offline
		if (!isOnline) {
			saveToLocal();
		}
	});

	onDestroy(() => {
		window.removeEventListener('online', updateOnlineStatus);
		window.removeEventListener('offline', updateOnlineStatus);
		if (autoSaveTimer) {
			clearInterval(autoSaveTimer);
		}
	});
</script>

{#if isTauri()}
	<div class="offline-indicator" class:offline={!isOnline}>
		<div class="status-icon">
			{#if isOnline}
				<Wifi size={16} />
			{:else}
				<WifiOff size={16} />
			{/if}
		</div>

		<div class="status-text">
			{#if isOnline}
				<span class="status-label">Online</span>
			{:else}
				<span class="status-label">Offline Mode</span>
			{/if}
		</div>

		{#if !isOnline && isTauri()}
			<div class="offline-badge">
				<CloudOff size={12} />
				{#if lastSaved}
					<span class="saved-time">Saved {formatTimeSince(lastSaved)}</span>
				{:else}
					<span class="saved-time">Local storage</span>
				{/if}
			</div>
		{/if}

		{#if isOnline && isTauri()}
			<div class="online-badge">
				<Cloud size={12} />
				<span>Synced</span>
			</div>
		{/if}
	</div>
{/if}

<style>
	.offline-indicator {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 6px 12px;
		background: white;
		border: 1px solid #d1d5db;
		border-radius: 8px;
		font-size: 13px;
		transition: all 0.3s;
	}

	.offline-indicator.offline {
		background: #fef3c7;
		border-color: #fbbf24;
	}

	.status-icon {
		display: flex;
		align-items: center;
		color: #10b981;
	}

	.offline .status-icon {
		color: #f59e0b;
	}

	.status-text {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	.status-label {
		font-weight: 500;
		color: #374151;
	}

	.offline-badge,
	.online-badge {
		display: flex;
		align-items: center;
		gap: 4px;
		padding: 2px 8px;
		border-radius: 4px;
		font-size: 11px;
		font-weight: 500;
	}

	.offline-badge {
		background: #fde68a;
		color: #92400e;
	}

	.online-badge {
		background: #d1fae5;
		color: #065f46;
	}

	.saved-time {
		color: inherit;
	}
</style>
