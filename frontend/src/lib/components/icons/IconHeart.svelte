<script lang="ts">
	interface Props {
		size?: number;
		color?: string;
		class?: string;
	}

	let { size = 32, color: _color = 'currentColor', class: className = '' }: Props = $props();
</script>

<svg
	width={size}
	height={size}
	viewBox="0 0 20 17"
	class={className}
	fill="none"
	xmlns="http://www.w3.org/2000/svg"
>
	<path
		d="M18.29 1.53H16.76V0H12.19V1.53H10.67V3.05H9.14V1.53H7.62V0H3.05V1.53H1.53V3.05H0V7.62H1.53V9.15H3.05V10.67H4.57V12.19H6.1V13.72H7.62V15.24H9.14V16.76H10.67V15.24H12.19V13.72H13.72V12.19H15.24V10.67H16.76V9.15H18.29V7.62H19.81V3.05H18.29V1.53ZM16.76 7.62H15.24V4.57H12.19V3.05H15.24V4.57H16.76V7.62Z"
		fill="#F30C0C"
	/>
</svg>
