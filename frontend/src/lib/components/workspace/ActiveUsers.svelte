<script lang="ts">
	import { presenceStore } from '$lib/stores/presence.svelte';
	import { authStore } from '$lib/stores/auth.svelte';
	import type { UserPresence } from '$lib/types/websocket';
	import { User } from 'lucide-svelte';

	import IconUsers from '$components/icons/IconUsers.svelte';

	interface Props {
		onUserClick?: (_user: UserPresence) => void;
	}

	let { onUserClick }: Props = $props();

	const users = $derived(presenceStore.users);
	const userCount = $derived(presenceStore.userCount);
	const currentUserId = $derived(authStore.user?.id);

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

<div class="relative">
	<!-- Trigger button -->
	<button
		class="flex items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition-all hover:border-gray-300 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:border-gray-600 dark:hover:bg-gray-700"
		onclick={toggleExpanded}
		title="Active users"
	>
		<IconUsers size={24} />
		{#if userCount > 0}
			<span
				class="flex h-5 min-w-5 items-center justify-center rounded-full bg-blue-500 px-1.5 text-[11px] font-semibold text-white"
				>{userCount}</span
			>
		{/if}
	</button>

	<!-- Dropdown panel -->
	{#if isExpanded}
		<div
			class="animate-slideDown absolute top-[calc(100%+8px)] right-0 z-1000 flex max-h-100 w-70 flex-col border border-gray-200 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-800"
		>
			<div
				class="flex items-center justify-between border-b border-gray-200 px-4 py-3 dark:border-gray-700"
			>
				<span class="text-sm font-semibold text-gray-900 dark:text-gray-50"
					>Active Users ({userCount})</span
				>
				<button
					class="flex h-6 w-6 items-center justify-center rounded border-none bg-transparent text-xl text-gray-500 transition-all hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-50"
					onclick={toggleExpanded}
					aria-label="Close">×</button
				>
			</div>

			<div
				class="scrollbar-thin scrollbar-thumb-gray-300 dark:scrollbar-thumb-gray-600 scrollbar-track-transparent hover:scrollbar-thumb-gray-400 dark:hover:scrollbar-thumb-gray-500 flex-1 overflow-y-auto p-2"
			>
				{#if users.length === 0}
					<div
						class="flex flex-col items-center justify-center p-8 px-4 text-center text-gray-400 dark:text-gray-600"
					>
						<User size={32} strokeWidth={1.5} />
						<p class="mt-3 text-sm">No active users</p>
					</div>
				{:else}
					{#each users as user (user.user_id)}
						<button
							class="flex w-full items-center gap-3 rounded-md border-none bg-transparent px-3 py-2.5 text-left transition-colors hover:bg-gray-50 dark:hover:bg-gray-700"
							onclick={() => handleUserClick(user)}
						>
							<!-- Avatar -->
							<div
								class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-[13px] font-semibold text-white"
								style="background-color: {user.user_color}"
							>
								{getInitials(user.user_name)}
							</div>

							<!-- User info -->
							<div class="min-w-0 flex-1">
								<div
									class="flex items-center gap-1.5 overflow-hidden text-sm font-medium text-ellipsis whitespace-nowrap text-gray-900 dark:text-gray-50"
								>
									{user.user_name}
									{#if currentUserId && user.user_id === currentUserId}
										<span class="shrink-0 text-xs font-medium text-blue-500 dark:text-blue-400"
											>(You)</span
										>
									{/if}
								</div>
								<div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
									{user.cursor ? 'Active' : 'Idle'} · {getTimeSinceLastSeen(user.last_seen)}
								</div>
							</div>

							<!-- Indicator -->
							<div
								class="h-2 w-2 shrink-0 rounded-full"
								style="background-color: {user.user_color}"
							></div>
						</button>
					{/each}
				{/if}
			</div>
		</div>

		<!-- Backdrop -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="fixed inset-0 z-999 bg-transparent" onclick={toggleExpanded}></div>
	{/if}
</div>

<style>
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

	.animate-slideDown {
		animation: slideDown 0.2s ease-out;
	}

	/* Custom scrollbar */
	:global(.scrollbar-thin)::-webkit-scrollbar {
		width: 6px;
	}

	:global(.scrollbar-thin)::-webkit-scrollbar-track {
		background: transparent;
	}

	:global(.scrollbar-thin)::-webkit-scrollbar-thumb {
		background: #d1d5db;
		border-radius: 3px;
	}

	:global(.scrollbar-thin)::-webkit-scrollbar-thumb:hover {
		background: #9ca3af;
	}

	:global(.dark .scrollbar-thin)::-webkit-scrollbar-thumb {
		background: #4b5563;
	}

	:global(.dark .scrollbar-thin)::-webkit-scrollbar-thumb:hover {
		background: #6b7280;
	}
</style>
