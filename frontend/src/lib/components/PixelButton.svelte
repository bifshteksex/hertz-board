<script lang="ts">
	export let variant: 'main' | 'second' = 'main';
	export let href: string | null = null;
	export let type: 'button' | 'submit' | 'reset' = 'button';
	export let disabled: boolean = false;
</script>

{#if href}
	<a
		{href}
		class="pixel-button cursor-pointer px-8 py-4 transition-transform duration-300 hover:scale-105"
		class:main={variant === 'main'}
		class:second={variant === 'second'}
		aria-disabled={disabled}
		style={disabled ? 'pointer-events: none; opacity: 0.5;' : ''}
	>
		<span
			class="text-4xl font-black uppercase"
			class:text-main={variant === 'main'}
			class:text-second={variant === 'second'}
		>
			<slot />
		</span>
	</a>
{:else}
	<button
		{type}
		class="pixel-button cursor-pointer px-8 py-4 transition-transform duration-300 hover:scale-105"
		class:main={variant === 'main'}
		class:second={variant === 'second'}
		{disabled}
	>
		<span
			class="text-4xl font-black uppercase"
			class:text-main={variant === 'main'}
			class:text-second={variant === 'second'}
		>
			<slot />
		</span>
	</button>
{/if}

<style>
	.pixel-button {
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}

	.pixel-button.main {
		border-image: url('../main.svg') 21 21 21 21 fill / 21px 21px 21px 21px;
	}

	.pixel-button.second {
		border-image: url('../second.svg') 21 21 21 21 fill / 21px 21px 21px 21px;
	}

	.text-main {
		color: #451b0b;
	}

	.text-second {
		color: #848484;
	}

	button:disabled,
	button[disabled] {
		cursor: not-allowed;
		opacity: 0.5;
		pointer-events: none;
	}

	button:disabled:hover {
		transform: none;
	}

	button:hover:disabled {
		/* prevent hover scale if disabled */
		transform: none;
	}
</style>
