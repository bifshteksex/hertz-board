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
	aria-label="Sign out icon"
>
	<g clip-path="url(#clip0_1128_11069)">
		<path
			d="M25.145 9.14001V13.71H19.045V16.76H25.145V21.33H26.665V19.81H28.195V18.29H29.715V16.76H31.235V13.71H29.715V12.19H28.195V10.67H26.665V9.14001H25.145Z"
			fill={color}
		/>
		<path d="M22.0951 21.33H20.5751V27.43H22.0951V21.33Z" fill={color} />
		<path d="M22.0951 1.51999H20.5751V9.13999H22.0951V1.51999Z" fill={color} />
		<path d="M14.475 30.48V28.95H20.575V27.43H14.475V7.62H12.955V30.48H14.475Z" fill={color} />
		<path d="M12.955 30.48H8.38501V32H12.955V30.48Z" fill={color} />
		<path d="M11.425 16.76H9.90503V19.81H11.425V16.76Z" fill={color} />
		<path d="M12.955 6.10001H9.90503V7.62001H12.955V6.10001Z" fill={color} />
		<path d="M8.38496 28.95H5.33496V30.48H8.38496V28.95Z" fill={color} />
		<path d="M9.90498 4.57001H6.85498V6.10001H9.90498V4.57001Z" fill={color} />
		<path d="M5.33503 27.43H2.28503V28.95H5.33503V27.43Z" fill={color} />
		<path d="M6.85506 3.04999H3.81506V4.56999H6.85506V3.04999Z" fill={color} />
		<path
			d="M2.28501 0V1.52H0.765015V27.43H2.28501V3.05H3.81501V1.52H20.575V0H2.28501Z"
			fill={color}
		/>
	</g>
	<defs>
		<clipPath id="clip0_1128_11069">
			<rect width="32" height="32" fill="white" />
		</clipPath>
	</defs>
</svg>
