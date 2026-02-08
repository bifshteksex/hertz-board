<script lang="ts">
	import IconCross from '$components/icons/IconCross.svelte';
	interface Props {
		show?: boolean;
		title?: string;
		onClose?: () => void;
		children?: import('svelte').Snippet;
	}

	let { show = false, title, onClose, children }: Props = $props();

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget && onClose) {
			onClose();
		}
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Escape' && onClose) {
			onClose();
		}
	}
</script>

{#if show}
	<div
		class="fixed inset-0 z-[1000] flex items-center justify-center bg-black/50 p-4 dark:bg-black/70"
		onclick={handleBackdropClick}
		onkeydown={handleKeyDown}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div
			class="flex max-h-[90vh] w-full max-w-[500px] flex-col overflow-hidden border-3 border-[#372d2e] bg-white shadow-[8px_8px_0_rgba(55,45,46,0.3)] dark:border-gray-700 dark:bg-gray-900 dark:shadow-[8px_8px_0_rgba(0,0,0,0.5)]"
		>
			{#if title}
				<div
					class="flex items-center justify-between border-b-2 border-[#372d2e] bg-gray-50 px-5 py-4 dark:border-gray-700 dark:bg-gray-800"
				>
					<h2 class="m-0 text-xl font-bold text-gray-800 dark:text-gray-100">{title}</h2>
					{#if onClose}
						<button
							onclick={onClose}
							class="flex h-8 w-8 items-center justify-center border-2 border-[#372d2e] bg-white text-gray-700 transition-all duration-150 hover:-translate-y-0.5 hover:bg-red-100 hover:text-red-600 hover:shadow-[0_2px_0_#372d2e] active:translate-y-0 active:shadow-none dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-red-900/30 dark:hover:text-red-400 dark:hover:shadow-[0_2px_0_rgba(75,85,99,1)]"
							aria-label="Close"
						>
							<IconCross size={22} />
						</button>
					{/if}
				</div>
			{/if}

			<div class="flex-1 overflow-y-auto p-5">
				{@render children?.()}
			</div>
		</div>
	</div>
{/if}

<style>
	/* Global styles for modal buttons */
	:global(.pixel-modal-buttons) {
		display: flex;
		gap: 12px;
		margin-top: 24px;
	}

	:global(.pixel-modal-buttons > *) {
		flex: 1;
	}
</style>
