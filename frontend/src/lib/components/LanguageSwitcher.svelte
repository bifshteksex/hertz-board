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

<div class="language-switcher" class:fixed={position === 'top-right'}>
	<button
		onclick={toggleMenu}
		class="language-btn"
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
				<span class="lang-code">{lang.label}</span>
				<span class="lang-name">{lang.name}</span>
			</button>
		{/each}
	</PixelMenu>
</div>

<style>
	.language-switcher {
		position: relative;
		z-index: 50;
	}

	.language-switcher.fixed {
		position: fixed;
		top: 1.5rem;
		right: 1.5rem;
	}

	.language-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 48px;
		height: 48px;
		background: white;
		border: 2px solid #372d2e;
		font-size: 14px;
		font-weight: 600;
		color: #372d2e;
		transition: all 0.15s;
		cursor: pointer;
	}

	.language-btn:hover {
		background: #f3f4f6;
		transform: translateY(-2px);
		box-shadow: 0 4px 0 #372d2e;
	}

	.language-btn:active {
		transform: translateY(0);
		box-shadow: none;
	}

	.lang-code {
		font-size: 14px;
		font-weight: 600;
	}

	.lang-name {
		font-size: 12px;
		color: #6b7280;
	}

	:global(.pixel-menu-item.active .lang-code) {
		color: #2563eb;
	}
</style>
