<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/services/api';
	import { Mail, CheckCircle, XCircle, Loader } from 'lucide-svelte';

	let token = $state('');
	let isLoading = $state(true);
	let isAccepting = $state(false);
	let error = $state('');
	let success = $state(false);
	let workspaceName = $state('');

	onMount(() => {
		// Get token from URL query params
		const params = new URLSearchParams(window.location.search);
		const tokenParam = params.get('token');

		if (!tokenParam) {
			error = 'Invalid invitation link';
			isLoading = false;
			return;
		}

		token = tokenParam;
		isLoading = false;
	});

	async function acceptInvitation() {
		if (!token) return;

		isAccepting = true;
		error = '';

		try {
			const response = await api.acceptInvitation(token);
			success = true;
			workspaceName = response.workspace?.name || 'Workspace';

			// Redirect to workspace after 2 seconds
			setTimeout(() => {
				if (response.workspace?.id) {
					goto(`/workspace/${response.workspace.id}`);
				} else {
					goto('/dashboard');
				}
			}, 2000);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to accept invitation';
		} finally {
			isAccepting = false;
		}
	}

	function goToDashboard() {
		goto('/dashboard');
	}
</script>

<svelte:head>
	<title>Accept Workspace Invitation</title>
</svelte:head>

<div class="invite-page">
	<div class="invite-container">
		{#if isLoading}
			<div class="loading-state">
				<Loader size={48} class="spinner" />
				<p>Loading invitation...</p>
			</div>
		{:else if success}
			<div class="success-state">
				<div class="success-icon">
					<CheckCircle size={64} />
				</div>
				<h1>Welcome aboard!</h1>
				<p>You've successfully joined <strong>{workspaceName}</strong></p>
				<p class="redirect-note">Redirecting to workspace...</p>
			</div>
		{:else if error}
			<div class="error-state">
				<div class="error-icon">
					<XCircle size={64} />
				</div>
				<h1>Invitation Invalid</h1>
				<p class="error-message">{error}</p>
				<p class="help-text">
					This invitation may have expired, been revoked, or already been accepted.
				</p>
				<button class="button primary" onclick={goToDashboard}>Go to Dashboard</button>
			</div>
		{:else}
			<div class="invite-prompt">
				<div class="invite-icon">
					<Mail size={64} />
				</div>
				<h1>You've been invited!</h1>
				<p>You've been invited to join a workspace on Hertz Board.</p>

				<div class="actions">
					<button class="button secondary" onclick={goToDashboard} disabled={isAccepting}>
						Decline
					</button>
					<button class="button primary" onclick={acceptInvitation} disabled={isAccepting}>
						{isAccepting ? 'Accepting...' : 'Accept Invitation'}
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.invite-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		padding: 20px;
	}

	.invite-container {
		background: white;
		border-radius: 16px;
		padding: 48px;
		max-width: 500px;
		width: 100%;
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
		text-align: center;
	}

	.loading-state,
	.success-state,
	.error-state,
	.invite-prompt {
		animation: fadeIn 0.3s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	:global(.spinner) {
		color: #667eea;
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

	.loading-state p {
		margin-top: 16px;
		color: #6b7280;
		font-size: 16px;
	}

	.success-icon,
	.error-icon,
	.invite-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0 auto 24px;
	}

	.success-icon {
		color: #10b981;
	}

	.error-icon {
		color: #ef4444;
	}

	.invite-icon {
		color: #667eea;
	}

	h1 {
		font-size: 28px;
		font-weight: 700;
		color: #111827;
		margin: 0 0 12px;
	}

	p {
		font-size: 16px;
		color: #6b7280;
		margin: 0 0 8px;
		line-height: 1.6;
	}

	p strong {
		color: #111827;
		font-weight: 600;
	}

	.error-message {
		color: #ef4444;
		font-weight: 500;
		margin-bottom: 16px;
	}

	.help-text {
		font-size: 14px;
		color: #9ca3af;
		margin-bottom: 24px;
	}

	.redirect-note {
		font-size: 14px;
		color: #9ca3af;
		margin-top: 16px;
	}

	.actions {
		display: flex;
		gap: 12px;
		margin-top: 32px;
		justify-content: center;
	}

	.button {
		padding: 12px 24px;
		border-radius: 8px;
		font-size: 16px;
		font-weight: 500;
		transition: all 0.2s;
		border: none;
		flex: 1;
		max-width: 200px;
	}

	.button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.button.secondary {
		background: white;
		color: #374151;
		border: 2px solid #e5e7eb;
	}

	.button.secondary:hover:not(:disabled) {
		background: #f9fafb;
		border-color: #d1d5db;
	}

	.button.primary {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
	}

	.button.primary:hover:not(:disabled) {
		box-shadow: 0 6px 16px rgba(102, 126, 234, 0.5);
		transform: translateY(-1px);
	}

	@media (max-width: 640px) {
		.invite-container {
			padding: 32px 24px;
		}

		h1 {
			font-size: 24px;
		}

		p {
			font-size: 14px;
		}

		.actions {
			flex-direction: column;
		}

		.button {
			max-width: 100%;
		}
	}
</style>
