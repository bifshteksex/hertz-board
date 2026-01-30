<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte';
	import { onMount } from 'svelte';

	let { children } = $props();
	let isChecking = $state(true);

	onMount(async () => {
		// Initialize auth store
		await authStore.initialize();

		isChecking = false;

		// Redirect to login if not authenticated
		if (!authStore.isAuthenticated) {
			goto('/auth/login');
		}
	});
</script>

{#if isChecking || !authStore.isInitialized}
	<div class="flex min-h-screen items-center justify-center">
		<div class="text-center">
			<div
				class="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent"
			></div>
			<p class="text-gray-600">Loading...</p>
		</div>
	</div>
{:else if authStore.isAuthenticated}
	<!-- Workspace layout: full screen, no sidebar -->
	{@render children()}
{/if}
