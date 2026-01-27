<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte';

	onMount(async () => {
		// Initialize auth and redirect if already logged in
		await authStore.initialize();
		if (authStore.isAuthenticated) {
			goto('/dashboard');
		}
	});
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100"
>
	<div class="mx-auto max-w-4xl px-4 text-center">
		<h1 class="mb-4 text-6xl font-bold text-gray-900">HertzBoard</h1>
		<p class="mb-8 text-xl text-gray-600">
			Real-time collaborative workspace platform built with CloudWeGo Hertz
		</p>

		<div class="mb-12 space-x-4">
			<a
				href="/auth/login"
				class="inline-block rounded-lg bg-blue-600 px-6 py-3 text-white transition hover:bg-blue-700"
			>
				Login
			</a>
			<a
				href="/auth/register"
				class="inline-block rounded-lg bg-gray-200 px-6 py-3 text-gray-900 transition hover:bg-gray-300"
			>
				Sign Up
			</a>
		</div>

		<div class="grid gap-8 md:grid-cols-3">
			<div class="rounded-lg bg-white p-6 shadow-sm">
				<h3 class="mb-2 text-lg font-semibold text-gray-900">Real-time Collaboration</h3>
				<p class="text-sm text-gray-600">
					Work together with your team in real-time with live cursors and presence
				</p>
			</div>
			<div class="rounded-lg bg-white p-6 shadow-sm">
				<h3 class="mb-2 text-lg font-semibold text-gray-900">Powerful Canvas</h3>
				<p class="text-sm text-gray-600">
					Create with text, shapes, images, and more on an infinite canvas
				</p>
			</div>
			<div class="rounded-lg bg-white p-6 shadow-sm">
				<h3 class="mb-2 text-lg font-semibold text-gray-900">Built with Go</h3>
				<p class="text-sm text-gray-600">
					High-performance backend powered by CloudWeGo Hertz and WebSockets
				</p>
			</div>
		</div>
	</div>
</div>
