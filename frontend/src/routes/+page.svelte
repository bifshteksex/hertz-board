<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte';
	import PixelButton from '$lib/components/PixelButton.svelte';
	import tv from '$lib/tv-new.webp';

	onMount(async () => {
		// Initialize auth and redirect if already logged in
		await authStore.initialize();
		if (authStore.isAuthenticated) {
			goto('/dashboard');
		}
	});
</script>

<div class="flex h-screen flex-col justify-between">
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
		<p class="mb-8 text-xl text-gray-600">Real-time collaborative workspace platform</p>

		<div class="mb-12 space-x-4 text-2xl">
			<PixelButton variant="main" href="/auth/login">Login</PixelButton>
			<PixelButton variant="second" href="/auth/register">Sign Up</PixelButton>
		</div>
	</div>

	<div class="grid border-x-2 border-t-2 border-[#372d2e] md:grid-cols-3">
		<div class=" border-r-2 border-[#372d2e] bg-white p-6">
			<h3 class="mb-2 text-lg font-semibold text-gray-900">Real-time Collaboration</h3>
			<p class="text-sm text-gray-600">
				Work together with your team in real-time with live cursors and presence
			</p>
		</div>
		<div class=" border-r-2 border-[#372d2e] bg-white p-6">
			<h3 class="mb-2 text-lg font-semibold text-gray-900">Powerful Canvas</h3>
			<p class="text-sm text-gray-600">
				Create with text, shapes, images, and more on an infinite canvas
			</p>
		</div>
		<div class=" bg-white p-6">
			<h3 class="mb-2 text-lg font-semibold text-gray-900">Built with Go</h3>
			<p class="text-sm text-gray-600">
				High-performance backend powered by CloudWeGo Hertz and WebSockets
			</p>
		</div>
	</div>
	<span class="flex justify-center border-2 border-[#372d2e] p-6"
		>Made with ❤️ by <a class="underline" target="_blank" href="https://rshang.in/">Roman Shangin</a
		></span
	>
</div>
