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
	aria-label="Search icon"
>
	<g clip-path="url(#clip0_1128_11711)">
		<path
			d="M24.38 19.81H22.85V22.86H19.81V24.38H16.76V25.9H21.33V27.43H22.85V28.95H24.38V30.48H25.9V32H27.43V30.48H28.95V28.95H30.47V27.43H32V25.9H30.47V24.38H28.95V22.86H27.43V21.33H25.9V16.76H24.38V19.81Z"
			fill={color}
		/>
		<path d="M27.43 10.67H25.9V16.76H27.43V10.67Z" fill={color} />
		<path d="M25.9 7.62H24.38V10.67H25.9V7.62Z" fill={color} />
		<path d="M24.38 4.57001H22.85V7.62001H24.38V4.57001Z" fill={color} />
		<path d="M22.85 12.19H21.33V15.24H22.85V12.19Z" fill={color} />
		<path d="M21.3301 9.14H19.8101V12.19H21.3301V9.14Z" fill={color} />
		<path d="M22.8501 3.05H19.8101V4.57H22.8501V3.05Z" fill={color} />
		<path d="M19.81 7.62H16.76V9.14H19.81V7.62Z" fill={color} />
		<path d="M19.81 1.52H16.76V3.05H19.81V1.52Z" fill={color} />
		<path d="M16.76 6.10001H13.71V7.62001H16.76V6.10001Z" fill={color} />
		<path d="M16.76 25.9H10.66V27.43H16.76V25.9Z" fill={color} />
		<path d="M16.76 0H10.66V1.52H16.76V0Z" fill={color} />
		<path d="M10.66 24.38H7.62V25.9H10.66V24.38Z" fill={color} />
		<path d="M10.66 1.52H7.62V3.05H10.66V1.52Z" fill={color} />
		<path d="M7.61995 22.86H4.56995V24.38H7.61995V22.86Z" fill={color} />
		<path d="M7.61995 3.05H4.56995V4.57H7.61995V3.05Z" fill={color} />
		<path d="M4.57005 19.81H3.05005V22.86H4.57005V19.81Z" fill={color} />
		<path d="M4.57005 4.57001H3.05005V7.62001H4.57005V4.57001Z" fill={color} />
		<path d="M3.05002 16.76H1.52002V19.81H3.05002V16.76Z" fill={color} />
		<path d="M3.05002 7.62H1.52002V10.67H3.05002V7.62Z" fill={color} />
		<path d="M1.52 10.67H0V16.76H1.52V10.67Z" fill={color} />
	</g>
	<defs>
		<clipPath id="clip0_1128_11711">
			<rect width="32" height="32" fill="white" />
		</clipPath>
	</defs>
</svg>
