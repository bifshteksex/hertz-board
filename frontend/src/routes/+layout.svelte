<script lang="ts">
	import '../app.css';
	import '$lib/i18n';
	import { i18n } from '$lib/stores/i18n.svelte';
	import { themeStore } from '$lib/stores/theme.svelte';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';

	// Update HTML lang attribute when locale changes
	$effect(() => {
		if (browser) {
			document.documentElement.lang = i18n.locale;
		}
	});

	// Handle theme updates on route changes
	$effect(() => {
		if (browser) {
			// Access $page.url to make this effect reactive to route changes
			$page.url.pathname;
			themeStore.handleRouteChange();
		}
	});

	let { children } = $props();
</script>

{@render children()}
