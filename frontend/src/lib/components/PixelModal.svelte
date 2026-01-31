<script lang="ts">
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
		class="pixel-modal-backdrop"
		onclick={handleBackdropClick}
		onkeydown={handleKeyDown}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div class="pixel-modal">
			{#if title}
				<div class="pixel-modal-header">
					<h2 class="pixel-modal-title">{title}</h2>
					{#if onClose}
						<button onclick={onClose} class="pixel-modal-close" aria-label="Close">
							<svg
								width="16"
								height="16"
								viewBox="0 0 16 16"
								fill="none"
								xmlns="http://www.w3.org/2000/svg"
							>
								<path d="M2 2H4V4H2V2Z" fill="currentColor" />
								<path d="M4 4H6V6H4V4Z" fill="currentColor" />
								<path d="M6 6H8V8H6V6Z" fill="currentColor" />
								<path d="M8 8H10V10H8V8Z" fill="currentColor" />
								<path d="M10 6H12V8H10V6Z" fill="currentColor" />
								<path d="M12 4H14V6H12V4Z" fill="currentColor" />
								<path d="M14 2H16V4H14V2Z" fill="currentColor" />
								<path d="M14 10H16V12H14V10Z" fill="currentColor" />
								<path d="M12 12H14V14H12V12Z" fill="currentColor" />
								<path d="M10 10H12V12H10V10Z" fill="currentColor" />
								<path d="M6 10H8V12H6V10Z" fill="currentColor" />
								<path d="M4 12H6V14H4V12Z" fill="currentColor" />
								<path d="M2 14H4V16H2V14Z" fill="currentColor" />
							</svg>
						</button>
					{/if}
				</div>
			{/if}

			<div class="pixel-modal-content">
				{@render children?.()}
			</div>
		</div>
	</div>
{/if}

<style>
	.pixel-modal-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 16px;
	}

	.pixel-modal {
		background: white;
		border: 3px solid #372d2e;
		box-shadow: 8px 8px 0 rgba(55, 45, 46, 0.3);
		max-width: 500px;
		width: 100%;
		max-height: 90vh;
		overflow: hidden;
		display: flex;
		flex-direction: column;
	}

	.pixel-modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 16px 20px;
		border-bottom: 2px solid #372d2e;
		background: #f9fafb;
	}

	.pixel-modal-title {
		font-size: 20px;
		font-weight: 700;
		color: #1f2937;
		margin: 0;
	}

	.pixel-modal-close {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		border: 2px solid #372d2e;
		background: white;
		color: #374151;
		cursor: pointer;
		transition: all 0.15s;
	}

	.pixel-modal-close:hover {
		background: #fee2e2;
		color: #dc2626;
		transform: translateY(-2px);
		box-shadow: 0 2px 0 #372d2e;
	}

	.pixel-modal-close:active {
		transform: translateY(0);
		box-shadow: none;
	}

	.pixel-modal-content {
		padding: 20px;
		overflow-y: auto;
		flex: 1;
	}

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
