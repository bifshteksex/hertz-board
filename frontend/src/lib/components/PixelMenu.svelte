<script lang="ts">
	interface Props {
		show?: boolean;
		position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left';
		children?: import('svelte').Snippet;
	}

	let { show = false, position = 'bottom-right', children }: Props = $props();
</script>

{#if show}
	<div
		class="pixel-menu"
		class:top-right={position === 'top-right'}
		class:top-left={position === 'top-left'}
		class:bottom-right={position === 'bottom-right'}
		class:bottom-left={position === 'bottom-left'}
	>
		{@render children?.()}
	</div>
{/if}

<style>
	.pixel-menu {
		position: absolute;
		background: white;
		border: 2px solid #372d2e;
		min-width: 160px;
		box-shadow: 4px 4px 0 rgba(55, 45, 46, 0.2);
		z-index: 100;
	}

	.pixel-menu.top-right {
		bottom: calc(100% + 8px);
		right: 0;
	}

	.pixel-menu.top-left {
		bottom: calc(100% + 8px);
		left: 0;
	}

	.pixel-menu.bottom-right {
		top: calc(100% + 8px);
		right: 0;
	}

	.pixel-menu.bottom-left {
		top: calc(100% + 8px);
		left: 0;
	}

	:global(.pixel-menu-item) {
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
		padding: 12px 16px;
		background: white;
		border: none;
		border-bottom: 1px solid #e5e7eb;
		cursor: pointer;
		transition: background 0.15s;
		text-align: left;
		font-size: 14px;
		color: #374151;
	}

	:global(.pixel-menu-item:last-child) {
		border-bottom: none;
	}

	:global(.pixel-menu-item:hover) {
		background: #f3f4f6;
	}

	:global(.pixel-menu-item.active) {
		background: #dbeafe;
		font-weight: 600;
		color: #2563eb;
	}

	:global(.pixel-menu-item.danger) {
		color: #dc2626;
	}

	:global(.pixel-menu-item.danger:hover) {
		background: #fee2e2;
	}

	:global(.pixel-menu-separator) {
		height: 1px;
		background: #e5e7eb;
		margin: 4px 0;
	}
</style>
