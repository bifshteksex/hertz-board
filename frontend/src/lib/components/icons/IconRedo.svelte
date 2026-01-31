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
	aria-label="Redo icon"
>
	<path
		d="M2.52602 19.5995L4.20117 13.734L5.66273 14.1514L3.98758 20.0169L2.52602 19.5995Z"
		fill={color}
	/>
	<path
		d="M5.66279 14.1514L6.50037 11.2187L7.97154 11.6388L7.13397 14.5716L5.66279 14.1514Z"
		fill={color}
	/>
	<path
		d="M7.9715 11.6388L8.38892 10.1773L9.85048 10.5947L9.43306 12.0562L7.9715 11.6388Z"
		fill={color}
	/>
	<path
		d="M9.85053 10.5947L10.2679 9.13312L11.7295 9.55053L11.3121 11.0121L9.85053 10.5947Z"
		fill={color}
	/>
	<path
		d="M26.8037 12.2645L28.2748 12.6847L28.6923 11.2231L30.1538 11.6405L28.0613 18.9676L20.7342 16.875L21.1544 15.4038L22.6159 15.8212L23.0334 14.3597L24.5045 14.7798L24.9219 13.3183L26.3835 13.7357L26.8037 12.2645Z"
		fill={color}
	/>
	<path
		d="M21.9893 12.4807L22.4094 11.0095L25.3422 11.8471L24.922 13.3183L21.9893 12.4807Z"
		fill={color}
	/>
	<path
		d="M11.7295 9.55054L12.1497 8.07936L15.0824 8.91694L14.6622 10.3881L11.7295 9.55054Z"
		fill={color}
	/>
	<path
		d="M15.0823 8.91694L15.4998 7.45538L22.8268 9.54794L22.4094 11.0095L15.0823 8.91694Z"
		fill={color}
	/>
</svg>
