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
	aria-label="Bin icon"
>
	<path
		d="M25.905 8.38003V25.14H27.435V8.38003H30.475V6.86003H28.955V5.33003H22.855V2.28003H21.335V5.33003H10.665V2.28003H9.14502V5.33003H3.05502V6.86003H1.52502V8.38003H4.57502V25.14H6.09502V8.38003H25.905Z"
		fill={color}
	/>
	<path d="M25.905 25.14H24.385V29.71H25.905V25.14Z" fill={color} />
	<path d="M24.385 29.71H7.625V31.24H24.385V29.71Z" fill={color} />
	<path d="M22.8551 11.43H21.3351V23.62H22.8551V11.43Z" fill={color} />
	<path d="M21.3351 23.62H19.8151V26.66H21.3351V23.62Z" fill={color} />
	<path d="M16.765 11.43H15.245V26.66H16.765V11.43Z" fill={color} />
	<path d="M21.335 0.76001H10.665V2.28001H21.335V0.76001Z" fill={color} />
	<path d="M12.195 23.62H10.665V26.66H12.195V23.62Z" fill={color} />
	<path d="M10.665 11.43H9.14502V23.62H10.665V11.43Z" fill={color} />
	<path d="M7.62497 25.14H6.09497V29.71H7.62497V25.14Z" fill={color} />
</svg>
