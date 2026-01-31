<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte';
	import { themeStore } from '$lib/stores/theme.svelte';
	import { i18n } from '$lib/i18n';
	import { onMount } from 'svelte';
	import { Menu, Bell, Home, Sun, Moon } from 'lucide-svelte';
	import IconSettings from '$components/icons/IconSettings.svelte';
	import IconSignout from '$components/icons/IconSignout.svelte';
	import PixelMenu from '$lib/components/PixelMenu.svelte';

	let { children } = $props();
	let showUserMenu = $state(false);
	let showMobileMenu = $state(false);

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

	async function handleLogout() {
		await authStore.logout();
		goto('/auth/login');
	}

	// Close dropdowns when clicking outside
	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-menu-container')) {
			showUserMenu = false;
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

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
	<div class="flex h-screen bg-gray-50">
		<!-- Sidebar -->
		<aside class="hidden w-64 flex-col border-r border-gray-200 bg-white md:flex">
			<div class="flex h-16 items-center border-b border-gray-200 px-6">
				<a href="/dashboard" class="text-2xl font-bold text-gray-900">HertzBoard</a>
			</div>

			<nav class="flex-1 space-y-1 px-3 py-4">
				<a
					href="/dashboard"
					class="flex items-center gap-3 rounded-lg px-3 py-2 text-gray-700 transition hover:bg-gray-100"
				>
					<Home size={20} />
					<span>{i18n.t('nav.workspaces')}</span>
				</a>
			</nav>

			<div class="border-t border-gray-200 p-4">
				<div class="user-menu-container relative">
					<button
						onclick={() => (showUserMenu = !showUserMenu)}
						class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left transition hover:bg-gray-100"
					>
						<div
							class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-600 text-white"
						>
							{authStore.user?.name?.charAt(0).toUpperCase() || 'U'}
						</div>
						<div class="flex-1 overflow-hidden">
							<p class="truncate text-sm font-medium text-gray-900">{authStore.user?.name}</p>
							<p class="truncate text-xs text-gray-500">{authStore.user?.email}</p>
						</div>
					</button>

					<PixelMenu show={showUserMenu} position="top-right">
						<a href="/settings" class="pixel-menu-item">
							<IconSettings size={16} />
							<span>{i18n.t('nav.settings')}</span>
						</a>
						<button onclick={handleLogout} class="pixel-menu-item">
							<IconSignout size={16} />
							<span>{i18n.t('nav.logout')}</span>
						</button>
					</PixelMenu>
				</div>
			</div>
		</aside>

		<!-- Main content -->
		<div class="flex flex-1 flex-col overflow-hidden">
			<!-- Header -->
			<header class="flex h-16 items-center justify-between border-b border-gray-200 bg-white px-6">
				<button onclick={() => (showMobileMenu = !showMobileMenu)} class="md:hidden">
					<Menu size={24} />
				</button>

				<div class="flex items-center gap-4">
					<button
						onclick={() => themeStore.toggleTheme()}
						class="rounded-lg p-2 text-gray-600 transition hover:bg-gray-100"
						title="Toggle theme"
					>
						{#if themeStore.isDark}
							<Sun size={20} />
						{:else}
							<Moon size={20} />
						{/if}
					</button>
					<button class="rounded-lg p-2 text-gray-600 transition hover:bg-gray-100">
						<Bell size={20} />
					</button>
				</div>
			</header>

			<!-- Page content -->
			<main class="flex-1 overflow-y-auto p-6">
				{@render children()}
			</main>
		</div>

		<!-- Mobile menu -->
		{#if showMobileMenu}
			<div class="fixed inset-0 z-50 md:hidden">
				<div
					role="button"
					tabindex="0"
					class="bg-opacity-50 absolute inset-0 bg-black"
					onclick={() => (showMobileMenu = false)}
					onkeydown={(e) => {
						if (e.key === 'Enter' || e.key === ' ') {
							showMobileMenu = false;
						}
					}}
				></div>
				<aside class="absolute top-0 left-0 h-full w-64 bg-white">
					<div class="flex h-16 items-center border-b border-gray-200 px-6">
						<a href="/dashboard" class="text-2xl font-bold text-gray-900">HertzBoard</a>
					</div>

					<nav class="flex-1 space-y-1 px-3 py-4">
						<a
							href="/dashboard"
							onclick={() => (showMobileMenu = false)}
							class="flex items-center gap-3 rounded-lg px-3 py-2 text-gray-700 transition hover:bg-gray-100"
						>
							<Home size={20} />
							<span>{i18n.t('nav.workspaces')}</span>
						</a>
					</nav>

					<div class="border-t border-gray-200 p-4">
						<div class="mb-2 flex items-center gap-3 px-3">
							<div
								class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-600 text-white"
							>
								{authStore.user?.name?.charAt(0).toUpperCase() || 'U'}
							</div>
							<div class="flex-1">
								<p class="truncate text-sm font-medium text-gray-900">{authStore.user?.name}</p>
								<p class="truncate text-xs text-gray-500">{authStore.user?.email}</p>
							</div>
						</div>
						<a
							href="/settings"
							onclick={() => (showMobileMenu = false)}
							class="flex items-center gap-2 rounded-lg px-3 py-2 text-gray-700 hover:bg-gray-100"
						>
							<IconSettings size={16} />
							{i18n.t('nav.settings')}
						</a>
						<button
							onclick={handleLogout}
							class="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-gray-700 hover:bg-gray-100"
						>
							<IconSignout size={16} />
							{i18n.t('nav.logout')}
						</button>
					</div>
				</aside>
			</div>
		{/if}
	</div>
{/if}
