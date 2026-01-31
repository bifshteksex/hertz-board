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
	aria-label="Question icon"
>
	<g clip-path="url(#clip0_1128_12003)">
		<path d="M32 3.05H30.48V28.95H32V3.05Z" fill={color} />
		<path d="M30.4799 28.95H28.95V30.48H30.4799V28.95Z" fill={color} />
		<path d="M30.4799 1.52H28.95V3.05H30.4799V1.52Z" fill={color} />
		<path d="M28.9501 30.48H3.05005V32H28.9501V30.48Z" fill={color} />
		<path
			d="M27.43 24.38H25.91V25.9H24.38V27.43H25.91V28.95H27.43V27.43H28.95V25.9H27.43V24.38Z"
			fill={color}
		/>
		<path
			d="M27.43 3.05H25.91V4.57H24.38V6.09H25.91V7.62H27.43V6.09H28.95V4.57H27.43V3.05Z"
			fill={color}
		/>
		<path
			d="M21.34 6.09H19.81V4.57H12.19V6.09H10.67V7.62H9.15002V12.19H12.19V9.14H13.72V7.62H16.76V9.14H18.29V12.19H16.76V13.71H15.24V15.24H13.72V18.28H18.29V16.76H19.81V15.24H21.34V13.71H22.86V7.62H21.34V6.09Z"
			fill={color}
		/>
		<path
			d="M18.2899 19.81H13.7199V21.33H12.1899V25.9H13.7199V27.43H18.2899V25.9H19.8099V21.33H18.2899V19.81Z"
			fill={color}
		/>
		<path
			d="M6.10005 24.38H4.57005V25.9H3.05005V27.43H4.57005V28.95H6.10005V27.43H7.62005V25.9H6.10005V24.38Z"
			fill={color}
		/>
		<path
			d="M6.10005 3.05H4.57005V4.57H3.05005V6.09H4.57005V7.62H6.10005V6.09H7.62005V4.57H6.10005V3.05Z"
			fill={color}
		/>
		<path d="M28.9501 0H3.05005V1.52H28.9501V0Z" fill={color} />
		<path d="M3.05003 28.95H1.53003V30.48H3.05003V28.95Z" fill={color} />
		<path d="M3.05003 1.52H1.53003V3.05H3.05003V1.52Z" fill={color} />
		<path d="M1.53 3.05H0V28.95H1.53V3.05Z" fill={color} />
	</g>
	<defs>
		<clipPath id="clip0_1128_12003">
			<rect width="32" height="32" fill="white" />
		</clipPath>
	</defs>
</svg>
