<script lang="ts">
	import type { UserPresence } from '$lib/types/websocket';

	interface Props {
		user: UserPresence;
		viewportZoom?: number;
		viewportOffsetX?: number;
		viewportOffsetY?: number;
	}

	let { user, viewportZoom = 1, viewportOffsetX = 0, viewportOffsetY = 0 }: Props = $props();

	// Transform cursor position based on viewport
	const screenX = $derived(user.cursor ? user.cursor.x * viewportZoom + viewportOffsetX : -9999);
	const screenY = $derived(user.cursor ? user.cursor.y * viewportZoom + viewportOffsetY : -9999);

	// Show cursor only if position is valid
	const isVisible = $derived(user.cursor !== undefined && user.cursor !== null);
</script>

{#if isVisible}
	<div
		class="user-cursor"
		style="
			--user-color: {user.user_color};
			transform: translate({screenX}px, {screenY}px);
		"
	>
		<!-- Cursor SVG (pointer arrow) -->
		<svg
			width="20"
			height="20"
			viewBox="0 0 20 20"
			fill="none"
			xmlns="http://www.w3.org/2000/svg"
			class="cursor-icon"
		>
			<path
				d="M2 2L18 9L9 10L7 18L2 2Z"
				fill="var(--user-color)"
				stroke="white"
				stroke-width="1.5"
				stroke-linejoin="round"
			/>
		</svg>

		<!-- User name label -->
		<div class="user-label" style="background-color: var(--user-color)">
			{user.user_name}
		</div>
	</div>
{/if}

<style>
	.user-cursor {
		position: absolute;
		top: 0;
		left: 0;
		pointer-events: none;
		z-index: 9999;
		transition: transform 0.1s ease-out;
		will-change: transform;
	}

	.cursor-icon {
		display: block;
		filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
	}

	.user-label {
		position: absolute;
		top: 16px;
		left: 12px;
		padding: 3px 6px;
		border-radius: 3px;
		font-size: 11px;
		font-weight: 500;
		color: white;
		white-space: nowrap;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
		opacity: 1;
		pointer-events: none;
	}
</style>
