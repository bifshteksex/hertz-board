<script lang="ts">
	import type { UserPresence } from '$lib/types/websocket';
	import type { CanvasElement } from '$lib/types/api';

	interface Props {
		user: UserPresence;
		elements: CanvasElement[];
		viewportZoom?: number;
		viewportOffsetX?: number;
		viewportOffsetY?: number;
	}

	let {
		user,
		elements,
		viewportZoom = 1,
		viewportOffsetX = 0,
		viewportOffsetY = 0
	}: Props = $props();

	// Get selected elements for this user
	const selectedElements = $derived(
		elements.filter((el) => user.selected_elements.includes(el.id))
	);
</script>

{#each selectedElements as element (element.id)}
	<div
		class="user-selection"
		style="
			--user-color: {user.user_color};
			left: {element.pos_x * viewportZoom + viewportOffsetX}px;
			top: {element.pos_y * viewportZoom + viewportOffsetY}px;
			width: {(element.width || 100) * viewportZoom}px;
			height: {(element.height || 100) * viewportZoom}px;
			transform: rotate({element.rotation || 0}deg);
		"
	>
		<!-- Selection outline -->
		<div class="selection-outline"></div>

		<!-- User name badge -->
		<div class="user-badge" style="background-color: var(--user-color)">
			{user.user_name}
		</div>
	</div>
{/each}

<style>
	.user-selection {
		position: absolute;
		pointer-events: none;
		z-index: 9998;
		transform-origin: center center;
	}

	.selection-outline {
		position: absolute;
		inset: -2px;
		border: 2px solid var(--user-color);
		border-radius: 4px;
		opacity: 0.8;
		animation: pulse 2s ease-in-out infinite;
	}

	.user-badge {
		position: absolute;
		top: -24px;
		left: -2px;
		padding: 2px 6px;
		border-radius: 3px;
		font-size: 11px;
		font-weight: 500;
		color: white;
		white-space: nowrap;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 0.8;
		}
		50% {
			opacity: 0.4;
		}
	}
</style>
