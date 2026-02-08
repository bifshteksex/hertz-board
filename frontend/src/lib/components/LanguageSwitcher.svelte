<script lang="ts">
	import { i18n } from '$lib/stores/i18n.svelte';
	import PixelMenu from './PixelMenu.svelte';

	interface Props {
		position?: 'top-right' | 'inline';
	}

	let { position = 'top-right' }: Props = $props();

	const languages = [
		{ code: 'en', label: 'EN', name: 'English' },
		{ code: 'ru', label: 'RU', name: 'Русский' },
		{ code: 'zh', label: '中', name: '中文' }
	] as const;

	let showMenu = $state(false);
	const currentLocale = $derived(i18n.locale);
	const currentLanguage = $derived(languages.find((l) => l.code === currentLocale));

	function handleLanguageChange(code: 'en' | 'ru' | 'zh') {
		i18n.setLocale(code);
		showMenu = false;
	}

	function toggleMenu() {
		showMenu = !showMenu;
	}
</script>

<div
	class="z-50"
	class:fixed={position === 'top-right'}
	class:top-6={position === 'top-right'}
	class:right-6={position === 'top-right'}
>
	<button
		onclick={toggleMenu}
		class="flex h-12 w-12 items-center justify-center border-2 border-[#372d2e] bg-white text-sm font-semibold text-[#372d2e] transition-all duration-150 hover:-translate-y-0.5 hover:bg-gray-100 hover:shadow-[0_4px_0_#372d2e] active:translate-y-0 active:shadow-none"
		title="Change language"
		aria-label="Change language"
	>
		{currentLanguage?.label}
	</button>

	<PixelMenu show={showMenu}>
		{#each languages as lang}
			<button
				onclick={() => handleLanguageChange(lang.code)}
				class="pixel-menu-item"
				class:active={currentLocale === lang.code}
			>
				<span class="text-sm font-semibold" class:text-blue-600={currentLocale === lang.code}>
					{lang.label}
				</span>
				<span class="text-xs text-gray-500">{lang.name}</span>
			</button>
		{/each}
	</PixelMenu>
</div>
