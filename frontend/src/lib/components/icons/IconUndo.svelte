<script lang="ts">
	interface Props {
		size?: number;
		color?: string;
		class?: string;
	}

	let { size = 32, color = 'currentColor', class: className = '' }: Props = $props();
</script>

<svg
	width={size}
	height={size}
	viewBox="0 0 32 32"
	fill="none"
	xmlns="http://www.w3.org/2000/svg"
	class={className}
	role="img"
	aria-label="Undo icon"
>
	<path
		d="M30.1538 19.5995L28.4786 13.734L27.0171 14.1514L28.6922 20.0169L30.1538 19.5995Z"
		fill={color}
	/>
	<path
		d="M27.017 14.1514L26.1794 11.2187L24.7083 11.6388L25.5458 14.5716L27.017 14.1514Z"
		fill={color}
	/>
	<path
		d="M24.7083 11.6388L24.2909 10.1773L22.8293 10.5947L23.2467 12.0562L24.7083 11.6388Z"
		fill={color}
	/>
	<path
		d="M22.8293 10.5947L22.4119 9.13312L20.9503 9.55053L21.3677 11.0121L22.8293 10.5947Z"
		fill={color}
	/>
	<path
		d="M5.87614 12.2645L4.40496 12.6847L3.98755 11.2231L2.52598 11.6405L4.61855 18.9676L11.9456 16.875L11.5254 15.4038L10.0639 15.8212L9.64646 14.3597L8.17528 14.7798L7.75786 13.3183L6.2963 13.7357L5.87614 12.2645Z"
		fill={color}
	/>
	<path
		d="M10.6905 12.4807L10.2704 11.0095L7.33764 11.8471L7.75781 13.3183L10.6905 12.4807Z"
		fill={color}
	/>
	<path
		d="M20.9503 9.55054L20.5302 8.07936L17.5974 8.91694L18.0176 10.3881L20.9503 9.55054Z"
		fill={color}
	/>
	<path
		d="M17.5975 8.91694L17.1801 7.45538L9.85301 9.54794L10.2704 11.0095L17.5975 8.91694Z"
		fill={color}
	/>
</svg>
