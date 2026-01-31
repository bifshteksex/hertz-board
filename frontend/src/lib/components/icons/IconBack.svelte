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
	aria-label="Back icon"
>
	<g clip-path="url(#clip0_1128_12658)">
		<path
			d="M30.48 9.14H28.95V6.09H27.43V4.57H25.9V3.05H22.86V1.52H19.81V0H12.19V1.52H9.14V3.05H6.09V4.57H4.57V6.09H3.05V9.14H1.52V12.19H0V19.81H1.52V22.86H3.05V25.9H4.57V27.43H6.09V28.95H9.14V30.48H12.19V32H19.81V30.48H22.86V28.95H25.9V27.43H27.43V25.9H28.95V22.86H30.48V19.81H32V12.19H30.48V9.14ZM25.9 21.33H24.38V22.86H22.86V24.38H18.29V22.86H19.81V21.33H21.33V19.81H19.81V18.29H13.71V21.33H12.19V19.81H10.67V18.29H9.14V16.76H7.62V15.24H6.09V13.71H7.62V12.19H9.14V10.67H10.67V9.14H12.19V7.62H13.71V10.67H22.86V12.19H24.38V13.71H25.9V21.33Z"
			fill={color}
		/>
	</g>
	<defs>
		<clipPath id="clip0_1128_12658">
			<rect width="32" height="32" fill="white" />
		</clipPath>
	</defs>
</svg>
