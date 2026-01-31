<script lang="ts">
	import { X, Search } from 'lucide-svelte';
	import { keyboardManager, formatShortcut, categoryNames } from '$lib/utils/keyboard';

	interface Props {
		onClose: () => void;
	}

	let { onClose }: Props = $props();

	let searchQuery = $state('');
	let panelElement: HTMLDivElement;

	// Get all shortcuts grouped by category
	const shortcuts = $derived(keyboardManager.getShortcuts());
	const categories = $derived(keyboardManager.getCategories());

	// Filter shortcuts by search query
	const filteredShortcuts = $derived(() => {
		if (!searchQuery.trim()) return shortcuts;

		const query = searchQuery.toLowerCase();
		return shortcuts.filter(
			(s) =>
				s.description.toLowerCase().includes(query) ||
				s.key.toLowerCase().includes(query) ||
				formatShortcut(s.key).toLowerCase().includes(query)
		);
	});

	// Group filtered shortcuts by category
	const groupedShortcuts = $derived(() => {
		const grouped = new Map<string, typeof shortcuts>();

		filteredShortcuts().forEach((shortcut) => {
			const category = shortcut.category;
			if (!grouped.has(category)) {
				grouped.set(category, []);
			}
			grouped.get(category)!.push(shortcut);
		});

		return grouped;
	});

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	function handleEscape(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			onClose();
		}
	}

	$effect(() => {
		document.addEventListener('keydown', handleEscape);
		return () => {
			document.removeEventListener('keydown', handleEscape);
		};
	});
</script>

<div
	class="fixed inset-0 z-[10000] flex items-center justify-center bg-black/50"
	onclick={handleBackdropClick}
>
	<div
		bind:this={panelElement}
		class="animate-fade-in flex max-h-[80vh] w-full max-w-3xl flex-col bg-white shadow-2xl"
	>
		<!-- Header -->
		<div class="flex items-center justify-between border-b border-gray-200 p-6">
			<div>
				<h2 class="text-2xl font-semibold text-gray-900">Keyboard Shortcuts</h2>
				<p class="mt-1 text-sm text-gray-500">Learn all the shortcuts to work faster</p>
			</div>
			<button
				onclick={onClose}
				class="rounded-lg p-2 transition-colors hover:bg-gray-100"
				aria-label="Close"
			>
				<X size={24} class="text-gray-500" />
			</button>
		</div>

		<!-- Search -->
		<div class="border-b border-gray-200 p-6">
			<div class="relative">
				<Search size={20} class="absolute top-1/2 left-3 -translate-y-1/2 text-gray-400" />
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search shortcuts..."
					class="w-full rounded-lg border border-gray-300 py-2 pr-4 pl-10 focus:ring-2 focus:ring-blue-500 focus:outline-none"
				/>
			</div>
		</div>

		<!-- Shortcuts List -->
		<div class="flex-1 overflow-y-auto p-6">
			{#if groupedShortcuts().size === 0}
				<div class="py-12 text-center text-gray-500">
					<p>No shortcuts found for "{searchQuery}"</p>
				</div>
			{:else}
				<div class="space-y-8">
					{#each categories as category}
						{#if groupedShortcuts().has(category)}
							<div>
								<h3 class="mb-3 text-sm font-semibold tracking-wider text-gray-400 uppercase">
									{categoryNames[category]}
								</h3>
								<div class="space-y-2">
									{#each groupedShortcuts().get(category) || [] as shortcut}
										<div
											class="flex items-center justify-between rounded-lg px-3 py-2 hover:bg-gray-50"
										>
											<span class="text-gray-700">{shortcut.description}</span>
											<kbd class="kbd">{formatShortcut(shortcut.key)}</kbd>
										</div>
									{/each}
								</div>
							</div>
						{/if}
					{/each}
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="border-t border-gray-200 bg-gray-50 p-6">
			<p class="text-center text-sm text-gray-600">
				Press <kbd class="kbd-small">?</kbd> anytime to show this panel
			</p>
		</div>
	</div>
</div>

<style>
	.kbd {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-width: 2rem;
		padding: 0.25rem 0.5rem;
		font-family: ui-monospace, monospace;
		font-size: 0.75rem;
		font-weight: 600;
		color: #374151;
		background: #f3f4f6;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
	}

	.kbd-small {
		padding: 0.125rem 0.375rem;
		font-size: 0.625rem;
		min-width: 1.5rem;
		font-family: ui-monospace, monospace;
		font-weight: 600;
		color: #374151;
		background: #f3f4f6;
		border: 1px solid #d1d5db;
		border-radius: 0.25rem;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
	}

	.animate-fade-in {
		animation: fadeIn 0.2s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}
</style>
