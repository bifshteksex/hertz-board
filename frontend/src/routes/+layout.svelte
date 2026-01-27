<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { themeStore } from '$lib/stores/theme.svelte';
	import '$lib/i18n';

	// Initialize theme on mount
	onMount(() => {
		// Theme is already initialized in the store constructor
		// This ensures it runs on client side
	});
</script>

<slot />
