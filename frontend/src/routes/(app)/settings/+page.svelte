<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';
	import { themeStore } from '$lib/stores/theme.svelte';
	import { i18n } from '$lib/i18n';
	import { User, Lock, Mail, Settings as SettingsIcon } from 'lucide-svelte';

	let activeTab = $state<'profile' | 'password' | 'account' | 'preferences'>('profile');

	// Profile form
	let profileName = $state(authStore.user?.name || '');
	let profileAvatarUrl = $state(authStore.user?.avatar_url || '');
	let profileError = $state('');
	let profileSuccess = $state(false);
	let isUpdatingProfile = $state(false);

	// Password form
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let passwordError = $state('');
	let passwordSuccess = $state(false);
	let isChangingPassword = $state(false);

	$effect(() => {
		// Update form when user changes
		if (authStore.user) {
			profileName = authStore.user.name;
			profileAvatarUrl = authStore.user.avatar_url || '';
		}
	});

	async function handleUpdateProfile(e: Event) {
		e.preventDefault();
		profileError = '';
		profileSuccess = false;
		isUpdatingProfile = true;

		try {
			await authStore.updateProfile({
				name: profileName,
				avatar_url: profileAvatarUrl || undefined
			});
			profileSuccess = true;
			setTimeout(() => (profileSuccess = false), 3000);
		} catch (err) {
			profileError = err instanceof Error ? err.message : 'Failed to update profile';
		} finally {
			isUpdatingProfile = false;
		}
	}

	async function handleChangePassword(e: Event) {
		e.preventDefault();
		passwordError = '';
		passwordSuccess = false;

		// Validation
		if (newPassword !== confirmPassword) {
			passwordError = i18n.t('settings.password.errorMatch');
			return;
		}

		if (newPassword.length < 8) {
			passwordError = i18n.t('settings.password.errorLength');
			return;
		}

		isChangingPassword = true;

		try {
			await authStore.changePassword(currentPassword, newPassword);
			passwordSuccess = true;
			currentPassword = '';
			newPassword = '';
			confirmPassword = '';
			setTimeout(() => (passwordSuccess = false), 3000);
		} catch (err) {
			passwordError = err instanceof Error ? err.message : 'Failed to change password';
		} finally {
			isChangingPassword = false;
		}
	}
</script>

<div class="mx-auto max-w-4xl space-y-6">
	<div>
		<h1 class="text-3xl font-bold text-gray-900">{i18n.t('settings.title')}</h1>
		<p class="mt-1 text-gray-600">{i18n.t('settings.subtitle')}</p>
	</div>

	<!-- Tabs -->
	<div class="border-b border-gray-200">
		<nav class="-mb-px flex space-x-8">
			<button
				onclick={() => (activeTab = 'profile')}
				class={`flex items-center gap-2 border-b-2 px-1 py-4 text-sm font-medium transition ${
					activeTab === 'profile'
						? 'border-blue-600 text-blue-600'
						: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
				}`}
			>
				<User size={16} />
				{i18n.t('settings.tabs.profile')}
			</button>
			<button
				onclick={() => (activeTab = 'password')}
				class={`flex items-center gap-2 border-b-2 px-1 py-4 text-sm font-medium transition ${
					activeTab === 'password'
						? 'border-blue-600 text-blue-600'
						: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
				}`}
			>
				<Lock size={16} />
				{i18n.t('settings.tabs.password')}
			</button>
			<button
				onclick={() => (activeTab = 'account')}
				class={`flex items-center gap-2 border-b-2 px-1 py-4 text-sm font-medium transition ${
					activeTab === 'account'
						? 'border-blue-600 text-blue-600'
						: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
				}`}
			>
				<Mail size={16} />
				{i18n.t('settings.tabs.account')}
			</button>
			<button
				onclick={() => (activeTab = 'preferences')}
				class={`flex items-center gap-2 border-b-2 px-1 py-4 text-sm font-medium transition ${
					activeTab === 'preferences'
						? 'border-blue-600 text-blue-600'
						: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
				}`}
			>
				<SettingsIcon size={16} />
				{i18n.t('settings.tabs.preferences')}
			</button>
		</nav>
	</div>

	<!-- Profile Tab -->
	{#if activeTab === 'profile'}
		<div class="rounded-lg border border-gray-200 bg-white p-6">
			<h2 class="mb-6 text-lg font-semibold text-gray-900">{i18n.t('settings.profile.title')}</h2>

			<form onsubmit={handleUpdateProfile} class="space-y-6">
				{#if profileError}
					<div class="rounded-md bg-red-50 p-4">
						<p class="text-sm text-red-800">{profileError}</p>
					</div>
				{/if}

				{#if profileSuccess}
					<div class="rounded-md bg-green-50 p-4">
						<p class="text-sm text-green-800">{i18n.t('settings.profile.successMessage')}</p>
					</div>
				{/if}

				<!-- Avatar Preview -->
				<div class="flex items-center gap-4">
					<div
						class="flex h-20 w-20 items-center justify-center rounded-full bg-blue-600 text-3xl font-bold text-white"
					>
						{#if profileAvatarUrl}
							<img
								src={profileAvatarUrl}
								alt="Avatar"
								class="h-full w-full rounded-full object-cover"
							/>
						{:else}
							{profileName.charAt(0).toUpperCase()}
						{/if}
					</div>
					<div>
						<p class="text-sm font-medium text-gray-900">{authStore.user?.email}</p>
						<p class="text-xs text-gray-500">
							{i18n.t('settings.profile.provider')}:
							<span class="capitalize">{authStore.user?.provider}</span>
						</p>
					</div>
				</div>

				<div>
					<label for="name" class="block text-sm font-medium text-gray-700"
						>{i18n.t('settings.profile.fullName')}</label
					>
					<input
						id="name"
						type="text"
						required
						bind:value={profileName}
						class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
					/>
				</div>

				<div>
					<label for="avatar" class="block text-sm font-medium text-gray-700"
						>{i18n.t('settings.profile.avatarUrl')}</label
					>
					<input
						id="avatar"
						type="url"
						bind:value={profileAvatarUrl}
						class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
						placeholder={i18n.t('settings.profile.avatarPlaceholder')}
					/>
					<p class="mt-1 text-xs text-gray-500">{i18n.t('settings.profile.avatarHint')}</p>
				</div>

				<div class="flex justify-end">
					<button
						type="submit"
						disabled={isUpdatingProfile}
						class="rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
					>
						{isUpdatingProfile
							? i18n.t('settings.profile.saving')
							: i18n.t('settings.profile.saveChanges')}
					</button>
				</div>
			</form>
		</div>
	{/if}

	<!-- Password Tab -->
	{#if activeTab === 'password'}
		<div class="rounded-lg border border-gray-200 bg-white p-6">
			<h2 class="mb-6 text-lg font-semibold text-gray-900">{i18n.t('settings.password.title')}</h2>

			{#if authStore.user?.provider !== 'email'}
				<div class="rounded-md bg-yellow-50 p-4">
					<p class="text-sm text-yellow-800">
						{i18n.t('settings.password.oauthWarning', { provider: authStore.user?.provider || '' })}
					</p>
				</div>
			{:else}
				<form onsubmit={handleChangePassword} class="space-y-6">
					{#if passwordError}
						<div class="rounded-md bg-red-50 p-4">
							<p class="text-sm text-red-800">{passwordError}</p>
						</div>
					{/if}

					{#if passwordSuccess}
						<div class="rounded-md bg-green-50 p-4">
							<p class="text-sm text-green-800">{i18n.t('settings.password.successMessage')}</p>
						</div>
					{/if}

					<div>
						<label for="current-password" class="block text-sm font-medium text-gray-700">
							{i18n.t('settings.password.current')}
						</label>
						<input
							id="current-password"
							type="password"
							required
							bind:value={currentPassword}
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
						/>
					</div>

					<div>
						<label for="new-password" class="block text-sm font-medium text-gray-700">
							{i18n.t('settings.password.new')}
						</label>
						<input
							id="new-password"
							type="password"
							required
							bind:value={newPassword}
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
						/>
						<p class="mt-1 text-xs text-gray-500">{i18n.t('settings.password.hint')}</p>
					</div>

					<div>
						<label for="confirm-password" class="block text-sm font-medium text-gray-700">
							{i18n.t('settings.password.confirm')}
						</label>
						<input
							id="confirm-password"
							type="password"
							required
							bind:value={confirmPassword}
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
						/>
					</div>

					<div class="flex justify-end">
						<button
							type="submit"
							disabled={isChangingPassword}
							class="rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
						>
							{isChangingPassword
								? i18n.t('settings.password.changing')
								: i18n.t('settings.password.change')}
						</button>
					</div>
				</form>
			{/if}
		</div>
	{/if}

	<!-- Account Tab -->
	{#if activeTab === 'account'}
		<div class="space-y-6">
			<div class="rounded-lg border border-gray-200 bg-white p-6">
				<h2 class="mb-4 text-lg font-semibold text-gray-900">{i18n.t('settings.account.title')}</h2>
				<div class="space-y-4">
					<div class="flex justify-between">
						<span class="text-sm text-gray-600">{i18n.t('settings.account.email')}</span>
						<span class="text-sm font-medium text-gray-900">{authStore.user?.email}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-sm text-gray-600">{i18n.t('settings.account.accountType')}</span>
						<span class="text-sm font-medium text-gray-900 capitalize"
							>{authStore.user?.provider}</span
						>
					</div>
					<div class="flex justify-between">
						<span class="text-sm text-gray-600">{i18n.t('settings.account.emailVerified')}</span>
						<span
							class={`text-sm font-medium ${authStore.user?.email_verified ? 'text-green-600' : 'text-yellow-600'}`}
						>
							{authStore.user?.email_verified
								? i18n.t('settings.account.verified')
								: i18n.t('settings.account.notVerified')}
						</span>
					</div>
					<div class="flex justify-between">
						<span class="text-sm text-gray-600">{i18n.t('settings.account.memberSince')}</span>
						<span class="text-sm font-medium text-gray-900">
							{new Date(authStore.user?.created_at || '').toLocaleDateString()}
						</span>
					</div>
				</div>
			</div>

			<div class="rounded-lg border border-red-200 bg-red-50 p-6">
				<h2 class="mb-2 text-lg font-semibold text-red-900">
					{i18n.t('settings.account.dangerZone')}
				</h2>
				<p class="mb-4 text-sm text-red-700">
					{i18n.t('settings.account.deleteWarning')}
				</p>
				<button
					onclick={() => alert('Account deletion is not yet implemented')}
					class="rounded-lg border border-red-600 bg-white px-4 py-2 text-red-600 transition hover:bg-red-50"
				>
					{i18n.t('settings.account.deleteAccount')}
				</button>
			</div>
		</div>
	{/if}

	<!-- Preferences Tab -->
	{#if activeTab === 'preferences'}
		<div class="rounded-lg border border-gray-200 bg-white p-6">
			<h2 class="mb-6 text-lg font-semibold text-gray-900">Preferences</h2>

			<div class="space-y-6">
				<!-- Language Selection -->
				<div>
					<label for="language" class="mb-2 block text-sm font-medium text-gray-700">
						Language
					</label>
					<select
						id="language"
						value={i18n.locale}
						onchange={(e) => i18n.setLocale(e.currentTarget.value as 'en' | 'ru' | 'zh')}
						class="block w-full max-w-xs rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
					>
						<option value="en">English</option>
						<option value="ru">Русский</option>
						<option value="zh">中文</option>
					</select>
					<p class="mt-1 text-xs text-gray-500">Select your preferred language for the interface</p>
				</div>

				<!-- Theme Selection -->
				<div>
					<label class="mb-2 block text-sm font-medium text-gray-700">Theme</label>
					<div class="flex gap-4">
						<button
							onclick={() => themeStore.setTheme('light')}
							class={`flex items-center gap-2 rounded-lg border-2 px-4 py-3 transition ${
								themeStore.theme === 'light'
									? 'border-blue-600 bg-blue-50 text-blue-700'
									: 'border-gray-300 bg-white text-gray-700 hover:border-gray-400'
							}`}
						>
							<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
								/>
							</svg>
							Light
						</button>

						<button
							onclick={() => themeStore.setTheme('dark')}
							class={`flex items-center gap-2 rounded-lg border-2 px-4 py-3 transition ${
								themeStore.theme === 'dark'
									? 'border-blue-600 bg-blue-50 text-blue-700'
									: 'border-gray-300 bg-white text-gray-700 hover:border-gray-400'
							}`}
						>
							<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
								/>
							</svg>
							Dark
						</button>
					</div>
					<p class="mt-1 text-xs text-gray-500">Choose your preferred color scheme</p>
				</div>
			</div>
		</div>
	{/if}
</div>
