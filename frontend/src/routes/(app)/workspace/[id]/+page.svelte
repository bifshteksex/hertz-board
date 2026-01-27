<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { api } from '$lib/services/api';
	import type { Workspace } from '$lib/types/api';
	import { ArrowLeft, Users, Share2 } from 'lucide-svelte';

	const workspaceId = $derived($page.params.id);
	let workspace = $state<Workspace | null>(null);
	let isLoading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			workspace = await api.getWorkspace(workspaceId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load workspace';
		} finally {
			isLoading = false;
		}
	});
</script>

{#if isLoading}
	<div class="flex h-screen items-center justify-center">
		<div class="text-center">
			<div
				class="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent"
			></div>
			<p class="text-gray-600">Loading workspace...</p>
		</div>
	</div>
{:else if error}
	<div class="flex h-screen items-center justify-center">
		<div class="text-center">
			<p class="mb-4 text-red-600">{error}</p>
			<button
				onclick={() => goto('/dashboard')}
				class="rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700"
			>
				Back to Dashboard
			</button>
		</div>
	</div>
{:else if workspace}
	<div class="flex h-screen flex-col">
		<!-- Workspace Header -->
		<header class="flex items-center justify-between border-b border-gray-200 bg-white px-6 py-3">
			<div class="flex items-center gap-4">
				<button
					onclick={() => goto('/dashboard')}
					class="rounded-lg p-2 text-gray-600 transition hover:bg-gray-100"
				>
					<ArrowLeft size={20} />
				</button>
				<div>
					<h1 class="text-lg font-semibold text-gray-900">{workspace.name}</h1>
					{#if workspace.description}
						<p class="text-sm text-gray-600">{workspace.description}</p>
					{/if}
				</div>
			</div>

			<div class="flex items-center gap-2">
				<button
					class="flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-2 text-sm transition hover:bg-gray-50"
				>
					<Users size={16} />
					Share
				</button>
			</div>
		</header>

		<!-- Canvas Area (Placeholder) -->
		<div class="flex flex-1 items-center justify-center bg-gray-50">
			<div class="text-center">
				<h2 class="mb-2 text-2xl font-bold text-gray-900">Canvas Coming Soon</h2>
				<p class="text-gray-600">The canvas editor will be implemented in Phase 6</p>
				<div class="mt-8 rounded-lg bg-white p-6 shadow-sm">
					<p class="mb-4 text-sm text-gray-700">
						Workspace ID: <code class="rounded bg-gray-100 px-2 py-1 font-mono text-xs"
							>{workspaceId}</code
						>
					</p>
					<p class="text-sm text-gray-700">
						Role: <span class="font-medium capitalize"
							>{workspace.user_role || workspace.role || 'viewer'}</span
						>
					</p>
				</div>
			</div>
		</div>
	</div>
{/if}
