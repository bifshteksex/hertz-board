<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';
	import { User, Lock, Mail } from 'lucide-svelte';

	let activeTab = $state<'profile' | 'password' | 'account'>('profile');

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
			passwordError = 'Passwords do not match';
			return;
		}

		if (newPassword.length < 8) {
			passwordError = 'Password must be at least 8 characters';
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
		<h1 class="text-3xl font-bold text-gray-900">Settings</h1>
		<p class="mt-1 text-gray-600">Manage your account settings and preferences</p>
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
				Profile
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
				Password
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
				Account
			</button>
		</nav>
	</div>

	<!-- Profile Tab -->
	{#if activeTab === 'profile'}
		<div class="rounded-lg border border-gray-200 bg-white p-6">
			<h2 class="mb-6 text-lg font-semibold text-gray-900">Profile Information</h2>

			<form onsubmit={handleUpdateProfile} class="space-y-6">
				{#if profileError}
					<div class="rounded-md bg-red-50 p-4">
						<p class="text-sm text-red-800">{profileError}</p>
					</div>
				{/if}

				{#if profileSuccess}
					<div class="rounded-md bg-green-50 p-4">
						<p class="text-sm text-green-800">Profile updated successfully!</p>
					</div>
				{/if}

				<!-- Avatar Preview -->
				<div class="flex items-center gap-4">
					<div class="flex h-20 w-20 items-center justify-center rounded-full bg-blue-600 text-3xl font-bold text-white">
						{#if profileAvatarUrl}
							<img src={profileAvatarUrl} alt="Avatar" class="h-full w-full rounded-full object-cover" />
						{:else}
							{profileName.charAt(0).toUpperCase()}
						{/if}
					</div>
					<div>
						<p class="text-sm font-medium text-gray-900">{authStore.user?.email}</p>
						<p class="text-xs text-gray-500">
							Provider: <span class="capitalize">{authStore.user?.provider}</span>
						</p>
					</div>
				</div>

				<div>
					<label for="name" class="block text-sm font-medium text-gray-700">Full Name</label>
					<input
						id="name"
						type="text"
						required
						bind:value={profileName}
						class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					/>
				</div>

				<div>
					<label for="avatar" class="block text-sm font-medium text-gray-700">Avatar URL</label>
					<input
						id="avatar"
						type="url"
						bind:value={profileAvatarUrl}
						class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
						placeholder="https://example.com/avatar.jpg"
					/>
					<p class="mt-1 text-xs text-gray-500">Optional: Enter a URL to your profile picture</p>
				</div>

				<div class="flex justify-end">
					<button
						type="submit"
						disabled={isUpdatingProfile}
						class="rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{isUpdatingProfile ? 'Saving...' : 'Save Changes'}
					</button>
				</div>
			</form>
		</div>
	{/if}

	<!-- Password Tab -->
	{#if activeTab === 'password'}
		<div class="rounded-lg border border-gray-200 bg-white p-6">
			<h2 class="mb-6 text-lg font-semibold text-gray-900">Change Password</h2>

			{#if authStore.user?.provider !== 'email'}
				<div class="rounded-md bg-yellow-50 p-4">
					<p class="text-sm text-yellow-800">
						You signed in with {authStore.user?.provider}. Password changes are not available for OAuth accounts.
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
							<p class="text-sm text-green-800">Password changed successfully!</p>
						</div>
					{/if}

					<div>
						<label for="current-password" class="block text-sm font-medium text-gray-700">
							Current Password
						</label>
						<input
							id="current-password"
							type="password"
							required
							bind:value={currentPassword}
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
						/>
					</div>

					<div>
						<label for="new-password" class="block text-sm font-medium text-gray-700">
							New Password
						</label>
						<input
							id="new-password"
							type="password"
							required
							bind:value={newPassword}
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
						/>
						<p class="mt-1 text-xs text-gray-500">Minimum 8 characters</p>
					</div>

					<div>
						<label for="confirm-password" class="block text-sm font-medium text-gray-700">
							Confirm New Password
						</label>
						<input
							id="confirm-password"
							type="password"
							required
							bind:value={confirmPassword}
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
						/>
					</div>

					<div class="flex justify-end">
						<button
							type="submit"
							disabled={isChangingPassword}
							class="rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							{isChangingPassword ? 'Changing...' : 'Change Password'}
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
				<h2 class="mb-4 text-lg font-semibold text-gray-900">Account Information</h2>
				<div class="space-y-4">
					<div class="flex justify-between">
						<span class="text-sm text-gray-600">Email</span>
						<span class="text-sm font-medium text-gray-900">{authStore.user?.email}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-sm text-gray-600">Account Type</span>
						<span class="text-sm font-medium capitalize text-gray-900">{authStore.user?.provider}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-sm text-gray-600">Email Verified</span>
						<span class={`text-sm font-medium ${authStore.user?.email_verified ? 'text-green-600' : 'text-yellow-600'}`}>
							{authStore.user?.email_verified ? 'Yes' : 'No'}
						</span>
					</div>
					<div class="flex justify-between">
						<span class="text-sm text-gray-600">Member Since</span>
						<span class="text-sm font-medium text-gray-900">
							{new Date(authStore.user?.created_at || '').toLocaleDateString()}
						</span>
					</div>
				</div>
			</div>

			<div class="rounded-lg border border-red-200 bg-red-50 p-6">
				<h2 class="mb-2 text-lg font-semibold text-red-900">Danger Zone</h2>
				<p class="mb-4 text-sm text-red-700">
					Once you delete your account, there is no going back. Please be certain.
				</p>
				<button
					onclick={() => alert('Account deletion is not yet implemented')}
					class="rounded-lg border border-red-600 bg-white px-4 py-2 text-red-600 transition hover:bg-red-50"
				>
					Delete Account
				</button>
			</div>
		</div>
	{/if}
</div>
