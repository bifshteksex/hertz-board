<script lang="ts">
	import { onMount } from 'svelte';

	interface Props {
		value?: string;
		onChange: (_color: string) => void;
		label?: string;
	}

	let { value = '#000000', onChange, label }: Props = $props();

	// Preset colors
	const presetColors = [
		'#000000', // Black
		'#ffffff', // White
		'#ef4444', // Red
		'#f97316', // Orange
		'#eab308', // Yellow
		'#22c55e', // Green
		'#3b82f6', // Blue
		'#a855f7', // Purple
		'#ec4899', // Pink
		'#64748b', // Gray
		'#fef3c7', // Light Yellow (sticky)
		'#dbeafe' // Light Blue
	];

	// Recently used colors (stored in localStorage)
	let recentColors = $state<string[]>([]);

	// Custom color picker state
	let showCustomPicker = $state(false);
	let customColor = $state(value);
	let hue = $state(0);
	let saturation = $state(100);
	let lightness = $state(50);

	onMount(() => {
		// Load recent colors from localStorage
		const stored = localStorage.getItem('hertz-board-recent-colors');
		if (stored) {
			try {
				recentColors = JSON.parse(stored);
			} catch {
				recentColors = [];
			}
		}

		// Parse current color to HSL
		parseColorToHSL(value);
	});

	function selectPresetColor(color: string) {
		onChange(color);
		addToRecentColors(color);
	}

	function handleCustomColorChange() {
		const color = hslToHex(hue, saturation, lightness);
		customColor = color;
		onChange(color);
		addToRecentColors(color);
	}

	function handleHexInput(e: Event) {
		const input = e.target as HTMLInputElement;
		let hex = input.value;

		// Add # if missing
		if (!hex.startsWith('#')) {
			hex = '#' + hex;
		}

		// Validate hex color
		if (/^#[0-9A-F]{6}$/i.test(hex)) {
			customColor = hex;
			onChange(hex);
			addToRecentColors(hex);
			parseColorToHSL(hex);
		}
	}

	function addToRecentColors(color: string) {
		// Add to recent colors (max 12)
		if (!recentColors.includes(color)) {
			recentColors = [color, ...recentColors.slice(0, 11)];
			localStorage.setItem('hertz-board-recent-colors', JSON.stringify(recentColors));
		}
	}

	function parseColorToHSL(hex: string) {
		// Convert hex to RGB
		const r = parseInt(hex.slice(1, 3), 16) / 255;
		const g = parseInt(hex.slice(3, 5), 16) / 255;
		const b = parseInt(hex.slice(5, 7), 16) / 255;

		const max = Math.max(r, g, b);
		const min = Math.min(r, g, b);
		const delta = max - min;

		// Calculate lightness
		lightness = ((max + min) / 2) * 100;

		if (delta === 0) {
			hue = 0;
			saturation = 0;
		} else {
			// Calculate saturation
			saturation = (delta / (1 - Math.abs(max + min - 1))) * 100;

			// Calculate hue
			if (max === r) {
				hue = ((g - b) / delta + (g < b ? 6 : 0)) * 60;
			} else if (max === g) {
				hue = ((b - r) / delta + 2) * 60;
			} else {
				hue = ((r - g) / delta + 4) * 60;
			}
		}
	}

	function hslToHex(h: number, s: number, l: number): string {
		h = h / 360;
		s = s / 100;
		l = l / 100;

		let r, g, b;

		if (s === 0) {
			r = g = b = l;
		} else {
			const hue2rgb = (p: number, q: number, t: number) => {
				if (t < 0) t += 1;
				if (t > 1) t -= 1;
				if (t < 1 / 6) return p + (q - p) * 6 * t;
				if (t < 1 / 2) return q;
				if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
				return p;
			};

			const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
			const p = 2 * l - q;

			r = hue2rgb(p, q, h + 1 / 3);
			g = hue2rgb(p, q, h);
			b = hue2rgb(p, q, h - 1 / 3);
		}

		const toHex = (x: number) => {
			const hex = Math.round(x * 255).toString(16);
			return hex.length === 1 ? '0' + hex : hex;
		};

		return `#${toHex(r)}${toHex(g)}${toHex(b)}`;
	}
</script>

<div class="flex flex-col gap-3 rounded-lg bg-white p-3">
	{#if label}
		<div class="text-[13px] font-semibold text-gray-700">{label}</div>
	{/if}

	<!-- Current color display -->
	<div class="flex items-center gap-2">
		<div
			class="h-8 w-8 shrink-0 rounded-md border-2 border-gray-200"
			style="background-color: {value}"
		></div>
		<input
			type="text"
			class="flex-1 rounded-md border border-gray-300 px-2.5 py-1.5 font-mono text-[13px] uppercase focus:border-blue-500 focus:outline-none"
			{value}
			onblur={handleHexInput}
			placeholder="#000000"
		/>
	</div>

	<!-- Preset colors -->
	<div class="flex flex-col gap-1.5">
		<div class="text-[11px] font-semibold tracking-wider text-gray-500 uppercase">Presets</div>
		<div class="grid grid-cols-6 gap-1.5">
			{#each presetColors as color}
				<button
					class="h-8 w-8 cursor-pointer rounded-md border-2 p-0 transition-all duration-150 hover:scale-110 hover:border-gray-400 {value ===
					color
						? '!border-[3px] border-blue-500 shadow-[0_0_0_2px_rgba(59,130,246,0.2)]'
						: 'border-gray-200'}"
					style="background-color: {color}"
					onclick={() => selectPresetColor(color)}
					title={color}
					aria-label={`Select color ${color}`}
				></button>
			{/each}
		</div>
	</div>

	<!-- Recent colors -->
	{#if recentColors.length > 0}
		<div class="flex flex-col gap-1.5">
			<div class="text-[11px] font-semibold tracking-wider text-gray-500 uppercase">Recent</div>
			<div class="grid grid-cols-6 gap-1.5">
				{#each recentColors as color}
					<button
						class="h-8 w-8 cursor-pointer rounded-md border-2 p-0 transition-all duration-150 hover:scale-110 hover:border-gray-400 {value ===
						color
							? '!border-[3px] border-blue-500 shadow-[0_0_0_2px_rgba(59,130,246,0.2)]'
							: 'border-gray-200'}"
						style="background-color: {color}"
						onclick={() => selectPresetColor(color)}
						title={color}
						aria-label={`Select color ${color}`}
					></button>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Custom color picker toggle -->
	<button
		class="cursor-pointer rounded-md border border-gray-200 bg-gray-100 px-3 py-2 text-[13px] font-medium text-gray-700 transition-all duration-150 hover:bg-gray-200"
		onclick={() => (showCustomPicker = !showCustomPicker)}
	>
		{showCustomPicker ? 'Hide' : 'Show'} Custom Picker
	</button>

	{#if showCustomPicker}
		<div class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3">
			<!-- Hue slider -->
			<div class="grid grid-cols-[80px_1fr_auto] items-center gap-2">
				<label for="hue-slider" class="text-xs font-medium text-gray-500">Hue</label>
				<input
					id="hue-slider"
					type="range"
					min="0"
					max="360"
					bind:value={hue}
					oninput={handleCustomColorChange}
					class="h-1.5 w-full cursor-pointer appearance-none rounded-sm outline-none [&::-moz-range-thumb]:h-4 [&::-moz-range-thumb]:w-4 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-blue-500 [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:shadow-sm [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-blue-500 [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow-sm"
					style="background: linear-gradient(to right, #ff0000, #ffff00, #00ff00, #00ffff, #0000ff, #ff00ff, #ff0000);"
				/>
				<span class="w-10 text-right font-mono text-xs text-gray-500">{Math.round(hue)}°</span>
			</div>

			<!-- Saturation slider -->
			<div class="grid grid-cols-[80px_1fr_auto] items-center gap-2">
				<label for="saturation-slider" class="text-xs font-medium text-gray-500">Saturation</label>
				<input
					id="saturation-slider"
					type="range"
					min="0"
					max="100"
					bind:value={saturation}
					oninput={handleCustomColorChange}
					class="h-1.5 w-full cursor-pointer appearance-none rounded-sm outline-none [&::-moz-range-thumb]:h-4 [&::-moz-range-thumb]:w-4 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-blue-500 [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:shadow-sm [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-blue-500 [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow-sm"
					style="background: linear-gradient(to right, #d1d5db, hsl({hue}, 100%, 50%));"
				/>
				<span class="w-10 text-right font-mono text-xs text-gray-500"
					>{Math.round(saturation)}%</span
				>
			</div>

			<!-- Lightness slider -->
			<div class="grid grid-cols-[80px_1fr_auto] items-center gap-2">
				<label for="lightness-slider" class="text-xs font-medium text-gray-500">Lightness</label>
				<input
					id="lightness-slider"
					type="range"
					min="0"
					max="100"
					bind:value={lightness}
					oninput={handleCustomColorChange}
					class="h-1.5 w-full cursor-pointer appearance-none rounded-sm outline-none [&::-moz-range-thumb]:h-4 [&::-moz-range-thumb]:w-4 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-blue-500 [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:shadow-sm [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-blue-500 [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow-sm"
					style="background: linear-gradient(to right, #000000, hsl({hue}, 100%, 50%), #ffffff);"
				/>
				<span class="w-10 text-right font-mono text-xs text-gray-500">{Math.round(lightness)}%</span
				>
			</div>

			<!-- Preview -->
			<div class="flex items-center gap-3 rounded-md border border-gray-200 bg-white p-2">
				<div
					class="h-10 w-10 rounded-md border-2 border-gray-200"
					style="background-color: {customColor}"
				></div>
				<span class="font-mono text-sm font-semibold text-gray-700">{customColor}</span>
			</div>
		</div>
	{/if}
</div>
