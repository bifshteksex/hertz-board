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
	aria-label="Settings icon"
>
	<g clip-path="url(#clip0_1128_12479)">
		<path d="M32 20.57H30.48V28.19H32V20.57Z" fill={color} />
		<path
			d="M30.48 2.28H28.95V0.759998H3.05V2.28H1.53V3.81H0V11.43H1.53V12.95H3.05V14.47H28.95V12.95H30.48V11.43H32V3.81H30.48V2.28ZM12.19 9.9H10.67V11.43H9.15V12.95H4.57V11.43H3.05V9.9H1.53V5.33H3.05V3.81H4.57V2.28H9.15V3.81H10.67V5.33H12.19V9.9Z"
			fill={color}
		/>
		<path d="M30.4799 28.19H28.95V29.71H30.4799V28.19Z" fill={color} />
		<path d="M30.4799 19.04H28.95V20.57H30.4799V19.04Z" fill={color} />
		<path d="M28.9501 22.09H27.4301V25.14H28.9501V22.09Z" fill={color} />
		<path d="M21.34 29.71V28.19H19.81V29.71H3.05005V31.24H28.9501V29.71H21.34Z" fill={color} />
		<path d="M27.43 20.57H25.91V22.09H27.43V20.57Z" fill={color} />
		<path d="M19.81 20.57H18.29V28.19H19.81V20.57Z" fill={color} />
		<path d="M10.67 5.33H9.15002V8.38H10.67V5.33Z" fill={color} />
		<path d="M9.14999 3.81H7.62V5.33H9.14999V3.81Z" fill={color} />
		<path d="M19.81 19.04V20.57H21.34V19.04H28.9501V17.52H3.05005V19.04H19.81Z" fill={color} />
		<path d="M3.05003 28.19H1.53003V29.71H3.05003V28.19Z" fill={color} />
		<path d="M3.05003 19.04H1.53003V20.57H3.05003V19.04Z" fill={color} />
		<path d="M1.53 20.57H0V28.19H1.53V20.57Z" fill={color} />
	</g>
	<defs>
		<clipPath id="clip0_1128_12479">
			<rect width="32" height="32" fill="white" />
		</clipPath>
	</defs>
</svg>
