<script lang="ts">
	import { api } from '$lib/services/api';

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

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md space-y-8">
		<div>
			<h1 class="text-center text-4xl font-bold text-gray-900">HertzBoard</h1>
			<h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">
				Reset your password
			</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				Enter your email address and we'll send you a link to reset your password.
			</p>
		</div>

		{#if success}
			<div class="rounded-md bg-green-50 p-4">
				<p class="text-sm text-green-800">
					Check your email for a link to reset your password. If it doesn't appear within a few
					minutes, check your spam folder.
				</p>
			</div>
			<div class="text-center">
				<a href="/auth/login" class="font-medium text-blue-600 hover:text-blue-500">
					Back to login
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
					<label for="email" class="sr-only">Email address</label>
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						bind:value={email}
						class="relative block w-full rounded-md border-0 px-3 py-2 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:z-10 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
						placeholder="Email address"
					/>
				</div>

				<div>
					<button
						type="submit"
						disabled={isLoading}
						class="group relative flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{isLoading ? 'Sending...' : 'Send reset link'}
					</button>
				</div>

				<div class="text-center text-sm">
					<a href="/auth/login" class="font-medium text-blue-600 hover:text-blue-500">
						Back to login
					</a>
				</div>
			</form>
		{/if}
	</div>
</div>
