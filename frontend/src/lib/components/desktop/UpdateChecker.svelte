<script lang="ts">
	import { onMount } from 'svelte';
	import { checkForUpdates, isTauri, type UpdateInfo } from '$lib/services/tauri';
	import { X, Download, RefreshCw } from 'lucide-svelte';

	let updateInfo = $state<UpdateInfo | null>(null);
	let showUpdateDialog = $state(false);
	let error = $state<string | null>(null);

	async function checkUpdates() {
		if (!isTauri()) return;

		error = null;

		try {
			const info = await checkForUpdates();
			updateInfo = info;

			if (info.available) {
				showUpdateDialog = true;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to check for updates';
			console.error('Update check failed:', err);
		}
	}

	function closeDialog() {
		showUpdateDialog = false;
	}

	function downloadUpdate() {
		if (updateInfo?.latest_version) {
			// Open releases page
			window.open(
				`https://github.com/YOUR_USERNAME/hertz-board/releases/tag/v${updateInfo.latest_version}`,
				'_blank'
			);
		}
		closeDialog();
	}

	onMount(() => {
		// Check for updates on mount (after 5 seconds delay)
		if (isTauri()) {
			setTimeout(() => {
				checkUpdates();
			}, 5000);
		}
	});
</script>

{#if isTauri()}
	<div class="update-checker">
		{#if showUpdateDialog && updateInfo?.available}
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div class="update-dialog-backdrop" onclick={closeDialog}></div>
			<div class="update-dialog">
				<div class="dialog-header">
					<div class="header-icon">
						<Download size={24} />
					</div>
					<h2>Update Available</h2>
					<button class="close-btn" onclick={closeDialog} aria-label="Close">
						<X size={20} />
					</button>
				</div>

				<div class="dialog-content">
					<p class="version-info">A new version of HertzBoard is available!</p>
					<div class="version-details">
						<div class="version-item">
							<span class="version-label">Current Version:</span>
							<span class="version-value">{updateInfo.current_version}</span>
						</div>
						<div class="version-item">
							<span class="version-label">Latest Version:</span>
							<span class="version-value new">{updateInfo.latest_version}</span>
						</div>
					</div>
					<p class="update-message">
						Download the latest version to get new features, improvements, and bug fixes.
					</p>
				</div>

				<div class="dialog-actions">
					<button class="btn-secondary" onclick={closeDialog}>Later</button>
					<button class="btn-primary" onclick={downloadUpdate}>
						<Download size={16} />
						Download Update
					</button>
				</div>
			</div>
		{/if}

		{#if error}
			<div class="update-error">
				<span class="error-message">{error}</span>
				<button class="retry-btn" onclick={checkUpdates}>
					<RefreshCw size={14} />
					Retry
				</button>
			</div>
		{/if}
	</div>
{/if}

<style>
	.update-checker {
		position: relative;
	}

	.update-dialog-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 9998;
		backdrop-filter: blur(4px);
	}

	.update-dialog {
		position: fixed;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		background: white;
		border-radius: 12px;
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
		z-index: 9999;
		width: 90%;
		max-width: 500px;
		animation: slideIn 0.3s ease-out;
	}

	@keyframes slideIn {
		from {
			opacity: 0;
			transform: translate(-50%, -45%);
		}
		to {
			opacity: 1;
			transform: translate(-50%, -50%);
		}
	}

	.dialog-header {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 24px 24px 16px;
		border-bottom: 1px solid #e5e7eb;
	}

	.header-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
		border-radius: 10px;
		color: white;
	}

	.dialog-header h2 {
		flex: 1;
		margin: 0;
		font-size: 20px;
		font-weight: 600;
		color: #111827;
	}

	.close-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		border: none;
		background: transparent;
		border-radius: 6px;
		color: #6b7280;
		cursor: pointer;
		transition: all 0.2s;
	}

	.close-btn:hover {
		background: #f3f4f6;
		color: #111827;
	}

	.dialog-content {
		padding: 24px;
	}

	.version-info {
		margin: 0 0 16px;
		font-size: 15px;
		color: #374151;
	}

	.version-details {
		display: flex;
		flex-direction: column;
		gap: 12px;
		padding: 16px;
		background: #f9fafb;
		border-radius: 8px;
		margin-bottom: 16px;
	}

	.version-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.version-label {
		font-size: 14px;
		color: #6b7280;
	}

	.version-value {
		font-size: 14px;
		font-weight: 600;
		color: #374151;
		font-family: 'Courier New', monospace;
	}

	.version-value.new {
		color: #3b82f6;
	}

	.update-message {
		margin: 0;
		font-size: 14px;
		color: #6b7280;
		line-height: 1.6;
	}

	.dialog-actions {
		display: flex;
		gap: 12px;
		padding: 16px 24px 24px;
		border-top: 1px solid #e5e7eb;
	}

	.btn-secondary,
	.btn-primary {
		flex: 1;
		padding: 10px 16px;
		border: none;
		border-radius: 8px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
	}

	.btn-secondary {
		background: #f3f4f6;
		color: #374151;
	}

	.btn-secondary:hover {
		background: #e5e7eb;
	}

	.btn-primary {
		background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
		color: white;
	}

	.btn-primary:hover {
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
	}

	.update-error {
		position: fixed;
		bottom: 20px;
		right: 20px;
		background: #fef2f2;
		border: 1px solid #fecaca;
		border-radius: 8px;
		padding: 12px 16px;
		display: flex;
		align-items: center;
		gap: 12px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		z-index: 1000;
	}

	.error-message {
		font-size: 14px;
		color: #991b1b;
	}

	.retry-btn {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 6px 12px;
		background: white;
		border: 1px solid #fca5a5;
		border-radius: 6px;
		color: #dc2626;
		font-size: 13px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.retry-btn:hover {
		background: #fee2e2;
	}
</style>
