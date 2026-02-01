<script lang="ts">
	import { presenceStore } from '$lib/stores/presence.svelte';
	import type { UserPresence } from '$lib/types/websocket';
	import { Users, User } from 'lucide-svelte';

	interface Props {
		onUserClick?: (_user: UserPresence) => void;
	}

	let { onUserClick }: Props = $props();

	const users = $derived(presenceStore.users);
	const userCount = $derived(presenceStore.userCount);

	let isExpanded = $state(false);

	function toggleExpanded() {
		isExpanded = !isExpanded;
	}

	function handleUserClick(user: UserPresence) {
		onUserClick?.(user);
		// Auto-collapse after click
		isExpanded = false;
	}

	function getInitials(name: string): string {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase()
			.slice(0, 2);
	}

	function getTimeSinceLastSeen(lastSeen: string): string {
		const now = Date.now();
		const then = new Date(lastSeen).getTime();
		const diffMs = now - then;
		const diffSec = Math.floor(diffMs / 1000);

		if (diffSec < 5) return 'just now';
		if (diffSec < 60) return `${diffSec}s ago`;
		const diffMin = Math.floor(diffSec / 60);
		if (diffMin < 60) return `${diffMin}m ago`;
		const diffHour = Math.floor(diffMin / 60);
		return `${diffHour}h ago`;
	}
</script>

<div class="active-users">
	<!-- Trigger button -->
	<button class="users-button" onclick={toggleExpanded} title="Active users">
		<Users size={20} />
		{#if userCount > 0}
			<span class="user-count">{userCount}</span>
		{/if}
	</button>

	<!-- Dropdown panel -->
	{#if isExpanded}
		<div class="users-dropdown">
			<div class="dropdown-header">
				<span class="dropdown-title">Active Users ({userCount})</span>
				<button class="close-button" onclick={toggleExpanded} aria-label="Close">×</button>
			</div>

			<div class="users-list">
				{#if users.length === 0}
					<div class="empty-state">
						<User size={32} strokeWidth={1.5} />
						<p>No other users online</p>
					</div>
				{:else}
					{#each users as user (user.user_id)}
						<button class="user-item" onclick={() => handleUserClick(user)}>
							<!-- Avatar -->
							<div class="user-avatar" style="background-color: {user.user_color}">
								{getInitials(user.user_name)}
							</div>

							<!-- User info -->
							<div class="user-info">
								<div class="user-name">{user.user_name}</div>
								<div class="user-status">
									{user.cursor ? 'Active' : 'Idle'} · {getTimeSinceLastSeen(user.last_seen)}
								</div>
							</div>

							<!-- Indicator -->
							<div class="user-indicator" style="background-color: {user.user_color}"></div>
						</button>
					{/each}
				{/if}
			</div>
		</div>

		<!-- Backdrop -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="backdrop" onclick={toggleExpanded}></div>
	{/if}
</div>

<style>
	.active-users {
		position: relative;
	}

	.users-button {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 8px 12px;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s;
		color: #374151;
		font-size: 14px;
		font-weight: 500;
	}

	.users-button:hover {
		background: #f9fafb;
		border-color: #d1d5db;
	}

	.user-count {
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 20px;
		height: 20px;
		padding: 0 6px;
		background: #3b82f6;
		color: white;
		border-radius: 10px;
		font-size: 11px;
		font-weight: 600;
	}

	.users-dropdown {
		position: absolute;
		top: calc(100% + 8px);
		right: 0;
		width: 280px;
		max-height: 400px;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
		z-index: 1000;
		display: flex;
		flex-direction: column;
		animation: slideDown 0.2s ease-out;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.dropdown-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 12px 16px;
		border-bottom: 1px solid #e5e7eb;
	}

	.dropdown-title {
		font-size: 14px;
		font-weight: 600;
		color: #111827;
	}

	.close-button {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		background: transparent;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 20px;
		color: #6b7280;
		transition: all 0.2s;
	}

	.close-button:hover {
		background: #f3f4f6;
		color: #111827;
	}

	.users-list {
		flex: 1;
		overflow-y: auto;
		padding: 8px;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 32px 16px;
		color: #9ca3af;
		text-align: center;
	}

	.empty-state p {
		margin-top: 12px;
		font-size: 14px;
	}

	.user-item {
		display: flex;
		align-items: center;
		gap: 12px;
		width: 100%;
		padding: 10px 12px;
		background: transparent;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		transition: background 0.2s;
		text-align: left;
	}

	.user-item:hover {
		background: #f9fafb;
	}

	.user-avatar {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;
		border-radius: 50%;
		color: white;
		font-size: 13px;
		font-weight: 600;
		flex-shrink: 0;
	}

	.user-info {
		flex: 1;
		min-width: 0;
	}

	.user-name {
		font-size: 14px;
		font-weight: 500;
		color: #111827;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.user-status {
		font-size: 12px;
		color: #6b7280;
		margin-top: 2px;
	}

	.user-indicator {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		flex-shrink: 0;
	}

	.backdrop {
		position: fixed;
		inset: 0;
		z-index: 999;
		background: transparent;
	}

	/* Scrollbar styling */
	.users-list::-webkit-scrollbar {
		width: 6px;
	}

	.users-list::-webkit-scrollbar-track {
		background: transparent;
	}

	.users-list::-webkit-scrollbar-thumb {
		background: #d1d5db;
		border-radius: 3px;
	}

	.users-list::-webkit-scrollbar-thumb:hover {
		background: #9ca3af;
	}
</style>
