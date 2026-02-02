<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte';
	import { i18n } from '$lib/stores/i18n.svelte';
	import PixelButton from '$lib/components/PixelButton.svelte';
	import IconHeart from '$components/icons/IconHeart.svelte';
	import IconGitHub from '$components/icons/IconGitHub.svelte';
	import IconStar from '$components/icons/IconStar.svelte';

	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
	import tv from '$lib/tv-new.webp';

	const currentYear = new Date().getFullYear();

	// SEO metadata
	const title = $derived(i18n.t('seo.pages.landing.title'));
	const description = $derived(i18n.t('seo.pages.landing.description'));
	const keywords = $derived(i18n.t('seo.keywords'));

	// GitHub stars state
	let starCount = $state<number | null>(null);

	async function fetchGitHubStars() {
		try {
			const response = await fetch('https://api.github.com/repos/bifshteksex/hertz-board');
			if (response.ok) {
				const data = await response.json();
				starCount = data.stargazers_count;
			}
		} catch (error) {
			console.error('Failed to fetch GitHub stars:', error);
		}
	}

	onMount(async () => {
		// Check for OAuth cookies and redirect immediately
		const cookies = document.cookie.split(';');
		const hasAccessToken = cookies.some((cookie) => cookie.trim().startsWith('access_token='));

		if (hasAccessToken) {
			// Give the API client time to process cookies
			await new Promise((resolve) => setTimeout(resolve, 100));
		}

		// Initialize auth and redirect if already logged in
		await authStore.initialize();
		if (authStore.isAuthenticated) {
			goto('/dashboard');
			return;
		}

		// Fetch GitHub stars
		fetchGitHubStars();
	});
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content={keywords} />

	<!-- Open Graph / Facebook -->
	<meta property="og:type" content="website" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:site_name" content={i18n.t('seo.ogSiteName')} />

	<!-- Twitter -->
	<meta name="twitter:card" content={i18n.t('seo.twitterCard')} />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />

	<!-- Canonical URL -->
	<link rel="canonical" href="https://hertz.rshang.in/" />
</svelte:head>

<div class="flex h-screen flex-col justify-between">
	<LanguageSwitcher />

	<div class="flex h-full flex-col items-center justify-center">
		<div class="flex items-center">
			<h1 class="mb-4 flex flex-col font-[CoralPixels] text-[10rem] leading-30 text-gray-900">
				<span>Hertz</span>
				<span class="ml-16">Board</span>
			</h1>
			<a href="https://www.cloudwego.io/" target="_blank" class="">
				<img class="w-72 transition-transform duration-300 hover:scale-105" src={tv} alt="TV" /></a
			>
		</div>
		<p class="mb-8 text-xl text-gray-600">{i18n.t('landing.subtitle')}</p>

		<div class="mb-12 space-x-4 text-2xl">
			<PixelButton variant="main" href="/auth/login">{i18n.t('landing.login')}</PixelButton>
			<PixelButton variant="second" href="/auth/register">{i18n.t('landing.signUp')}</PixelButton>
		</div>
	</div>

	<div class="grid border-x-2 border-t-2 border-[#372d2e] md:grid-cols-3">
		<div class=" border-r-2 border-[#372d2e] bg-white p-6">
			<h3 class="mb-2 text-lg font-semibold text-gray-900">
				{i18n.t('landing.features.realtime.title')}
			</h3>
			<p class="text-sm text-gray-600">
				{i18n.t('landing.features.realtime.description')}
			</p>
		</div>
		<div class=" border-r-2 border-[#372d2e] bg-white p-6">
			<h3 class="mb-2 text-lg font-semibold text-gray-900">
				{i18n.t('landing.features.canvas.title')}
			</h3>
			<p class="text-sm text-gray-600">
				{i18n.t('landing.features.canvas.description')}
			</p>
		</div>
		<div class=" bg-white p-6">
			<h3 class="mb-2 text-lg font-semibold text-gray-900">
				{i18n.t('landing.features.tech.title')}
			</h3>
			<p class="text-sm text-gray-600">
				{i18n.t('landing.features.tech.description')}
			</p>
		</div>
	</div>
	<div class="flex items-center justify-center border-2 border-[#372d2e] p-6">
		<span class="border-r-2 border-[#372d2e] px-3">© {currentYear} HertzBoard</span>
		<span class="flex items-center justify-center gap-1 border-r-2 border-[#372d2e] px-3">
			{i18n.t('landing.footer').split('{heart}')[0]}
			<span class="pixel-heartbeat">
				<IconHeart size={20} />
			</span>
			{i18n.t('landing.footer').split('{heart}')[1]}
			<a class="ml-1 underline" target="_blank" href="https://rshang.in/">Roman Shangin</a></span
		>
		<a
			target="_blank"
			href="https://github.com/bifshteksex/hertz-board"
			class="flex items-center gap-1 px-3"
		>
			<IconGitHub />
			<span class="underline">GitHub</span>
			{#if starCount !== null}
				<span
					class="ml-1 flex items-center gap-x-1 border border-[#372d2e] bg-gray-100 px-2 py-1 text-sm font-semibold"
				>
					<IconStar size={16} />
					<p>{starCount.toLocaleString()}</p>
				</span>
			{/if}
		</a>
	</div>
</div>

<style>
	@keyframes pixelHeartbeat {
		0% {
			transform: scale(1);
		}
		10% {
			transform: scale(1.2);
		}
		20% {
			transform: scale(1);
		}
		30% {
			transform: scale(1.3);
		}
		40% {
			transform: scale(1);
		}
		100% {
			transform: scale(1);
		}
	}

	.pixel-heartbeat {
		display: inline-flex;
		animation: pixelHeartbeat 2s steps(1) infinite;
	}
</style>
