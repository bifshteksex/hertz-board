<script lang="ts">
	import { api } from '$lib/services/api';
	import { i18n } from '$lib/stores/i18n.svelte';
	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';

	let email = $state('');
	let error = $state('');
	let success = $state(false);
	let isLoading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		isLoading = true;

		try {
			await api.forgotPassword({ email });
			success = true;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to send reset email';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center px-4 py-12 sm:px-6 lg:px-8">
	<LanguageSwitcher />

	<div class="w-full max-w-md space-y-8">
		<div>
			<h1 class="text-center text-4xl font-bold text-gray-900">{i18n.t('auth.appName')}</h1>
			<h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">
				{i18n.t('auth.resetTitle')}
			</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				{i18n.t('auth.resetDescription')}
			</p>
		</div>

		{#if success}
			<div class="rounded-md bg-green-50 p-4">
				<p class="text-sm text-green-800">
					{i18n.t('auth.resetSuccess')}
				</p>
			</div>
			<div class="text-center">
				<a href="/auth/login" class="font-medium text-blue-600 hover:text-blue-500">
					{i18n.t('auth.backToLogin')}
				</a>
			</div>
		{:else}
			<form class="mt-8 space-y-6" onsubmit={handleSubmit}>
				{#if error}
					<div class="rounded-md bg-red-50 p-4">
						<p class="text-sm text-red-800">{error}</p>
					</div>
				{/if}

				<div>
					<label for="email" class="sr-only">{i18n.t('auth.email')}</label>
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						bind:value={email}
						class="relative block w-full border bg-white px-3 py-2 text-gray-900 ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:z-10 focus:ring-2 focus:ring-blue-600 focus:ring-inset sm:text-sm sm:leading-6"
						placeholder={i18n.t('auth.emailPlaceholder')}
					/>
				</div>

				<div>
					<button
						type="submit"
						disabled={isLoading}
						class="flex w-full items-center justify-center border-2 border-[#372d2e] bg-white px-3 py-2 text-sm font-semibold text-[#372d2e] transition-all duration-150 hover:-translate-y-0.5 hover:bg-gray-100 hover:shadow-[0_4px_0_#372d2e] active:translate-y-0 active:shadow-none disabled:cursor-not-allowed disabled:opacity-50"
					>
						{isLoading ? i18n.t('auth.sending') : i18n.t('auth.sendResetLink')}
					</button>
				</div>

				<div class="text-center text-sm">
					<a href="/auth/login" class="font-medium text-blue-600 hover:text-blue-500">
						{i18n.t('auth.backToLogin')}
					</a>
				</div>
			</form>
		{/if}
	</div>
</div>
