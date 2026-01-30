<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import TextAlign from '@tiptap/extension-text-align';
	import { TextStyle } from '@tiptap/extension-text-style';
	import { Color } from '@tiptap/extension-color';
	import { FontFamily } from '@tiptap/extension-font-family';

	interface Props {
		content: string;
		onUpdate: (_content: string, _html: string) => void;
		onBlur?: () => void;
		fontSize?: number;
		fontFamily?: string;
		color?: string;
		textAlign?: 'left' | 'center' | 'right';
	}

	let {
		content,
		onUpdate,
		onBlur,
		fontSize = 16,
		fontFamily = 'Inter',
		color = '#000000',
		textAlign = 'left'
	}: Props = $props();

	let element: HTMLDivElement;
	let editor: Editor | null = null;

	onMount(() => {
		editor = new Editor({
			element: element,
			extensions: [
				StarterKit.configure({
					heading: {
						levels: [1, 2, 3]
					}
				}),
				TextAlign.configure({
					types: ['heading', 'paragraph']
				}),
				TextStyle,
				Color,
				FontFamily
			],
			content: content,
			editorProps: {
				attributes: {
					class: 'tiptap-editor'
				}
			},
			onUpdate: ({ editor }) => {
				const text = editor.getText();
				const html = editor.getHTML();
				onUpdate(text, html);
			},
			onBlur: () => {
				if (onBlur) onBlur();
			},
			autofocus: 'end'
		});

		// Применяем начальные стили
		if (editor) {
			editor
				.chain()
				.focus()
				.setFontFamily(fontFamily)
				.setColor(color)
				.setTextAlign(textAlign)
				.run();
		}
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});

	// Reactive updates для стилей
	$effect(() => {
		if (editor && !editor.isDestroyed) {
			editor.chain().setFontFamily(fontFamily).run();
		}
	});

	$effect(() => {
		if (editor && !editor.isDestroyed) {
			editor.chain().setColor(color).run();
		}
	});

	$effect(() => {
		if (editor && !editor.isDestroyed) {
			editor.chain().setTextAlign(textAlign).run();
		}
	});
</script>

<div
	bind:this={element}
	class="text-editor-container"
	style:font-size="{fontSize}px"
	style:font-family={fontFamily}
></div>

<style>
	.text-editor-container {
		width: 100%;
		height: 100%;
		overflow: auto;
		padding: 8px;
	}

	.text-editor-container :global(.tiptap-editor) {
		width: 100%;
		height: 100%;
		outline: none;
	}

	.text-editor-container :global(.tiptap-editor p) {
		margin: 0 0 0.5em 0;
	}

	.text-editor-container :global(.tiptap-editor p:last-child) {
		margin-bottom: 0;
	}

	.text-editor-container :global(.tiptap-editor h1),
	.text-editor-container :global(.tiptap-editor h2),
	.text-editor-container :global(.tiptap-editor h3) {
		margin: 0.5em 0;
		font-weight: 600;
	}

	.text-editor-container :global(.tiptap-editor h1) {
		font-size: 2em;
	}

	.text-editor-container :global(.tiptap-editor h2) {
		font-size: 1.5em;
	}

	.text-editor-container :global(.tiptap-editor h3) {
		font-size: 1.25em;
	}

	.text-editor-container :global(.tiptap-editor strong) {
		font-weight: 700;
	}

	.text-editor-container :global(.tiptap-editor em) {
		font-style: italic;
	}

	.text-editor-container :global(.tiptap-editor u) {
		text-decoration: underline;
	}

	.text-editor-container :global(.tiptap-editor s) {
		text-decoration: line-through;
	}

	.text-editor-container :global(.tiptap-editor ul),
	.text-editor-container :global(.tiptap-editor ol) {
		padding-left: 1.5em;
		margin: 0.5em 0;
	}

	.text-editor-container :global(.tiptap-editor li) {
		margin: 0.25em 0;
	}

	.text-editor-container :global(.tiptap-editor code) {
		background-color: #f3f4f6;
		padding: 0.2em 0.4em;
		border-radius: 3px;
		font-family: 'Courier New', monospace;
		font-size: 0.9em;
	}

	.text-editor-container :global(.tiptap-editor pre) {
		background-color: #1f2937;
		color: #f3f4f6;
		padding: 1em;
		border-radius: 6px;
		overflow-x: auto;
		margin: 0.5em 0;
	}

	.text-editor-container :global(.tiptap-editor pre code) {
		background-color: transparent;
		padding: 0;
		color: inherit;
	}

	.text-editor-container :global(.tiptap-editor blockquote) {
		border-left: 3px solid #d1d5db;
		padding-left: 1em;
		margin: 0.5em 0;
		font-style: italic;
	}
</style>
