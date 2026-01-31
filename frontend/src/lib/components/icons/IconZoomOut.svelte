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
	aria-label="Zoom out icon"
>
	<g clip-path="url(#clip0_1128_11234)">
		<path d="M32 9.14001H30.48V13.71H32V9.14001Z" fill={color} />
		<path d="M32 1.51999H30.48V6.08999H32V1.51999Z" fill={color} />
		<path
			d="M30.48 25.9H28.95V24.38H27.43V22.86H25.91V25.9H22.86V27.43H24.38V28.95H25.91V30.48H27.43V32H30.48V30.48H32V27.43H30.48V25.9Z"
			fill={color}
		/>
		<path d="M30.48 0H25.91V1.52H30.48V0Z" fill={color} />
		<path d="M28.9501 18.28H27.4301V22.86H28.9501V18.28Z" fill={color} />
		<path d="M27.43 15.24H25.91V18.28H27.43V15.24Z" fill={color} />
		<path d="M24.38 19.81H16.76V21.33H24.38V19.81Z" fill={color} />
		<path d="M25.91 13.71H22.86V15.24H25.91V13.71Z" fill={color} />
		<path d="M22.86 27.43H18.29V28.95H22.86V27.43Z" fill={color} />
		<path d="M22.86 12.19H18.29V13.71H22.86V12.19Z" fill={color} />
		<path d="M22.86 0H18.29V1.52H22.86V0Z" fill={color} />
		<path d="M18.29 25.9H15.24V27.43H18.29V25.9Z" fill={color} />
		<path d="M18.29 13.71H15.24V15.24H18.29V13.71Z" fill={color} />
		<path d="M15.24 22.86H13.72V25.9H15.24V22.86Z" fill={color} />
		<path d="M15.24 15.24H13.72V18.28H15.24V15.24Z" fill={color} />
		<path d="M13.7199 18.28H12.1899V22.86H13.7199V18.28Z" fill={color} />
		<path d="M13.72 30.48H9.15002V32H13.72V30.48Z" fill={color} />
		<path d="M13.72 0H9.15002V1.52H13.72V0Z" fill={color} />
		<path d="M6.10003 30.48H1.53003V32H6.10003V30.48Z" fill={color} />
		<path d="M6.10003 0H1.53003V1.52H6.10003V0Z" fill={color} />
		<path d="M1.53 25.9H0V30.48H1.53V25.9Z" fill={color} />
		<path d="M1.53 18.28H0V22.86H1.53V18.28Z" fill={color} />
		<path d="M1.53 9.14001H0V13.71H1.53V9.14001Z" fill={color} />
		<path d="M1.53 1.51999H0V6.08999H1.53V1.51999Z" fill={color} />
	</g>
	<defs>
		<clipPath id="clip0_1128_11234">
			<rect width="32" height="32" fill="white" />
		</clipPath>
	</defs>
</svg>
