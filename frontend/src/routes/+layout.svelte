<script lang="ts">
	import '../app.css';
	import '$lib/i18n';
	import { i18n } from '$lib/stores/i18n.svelte';
	import { browser } from '$app/environment';

	// Update HTML lang attribute when locale changes
	$effect(() => {
		if (browser) {
			document.documentElement.lang = i18n.locale;
		}
	});

	let { children } = $props();
</script>

{@render children()}
