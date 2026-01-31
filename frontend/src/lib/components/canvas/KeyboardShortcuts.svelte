<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { canvasStore } from '$lib/stores/canvas.svelte';
	import { keyboardManager } from '$lib/utils/keyboard';

	// Props for callbacks
	interface Props {
		onCut?: () => void;
		onCopy?: () => void;
		onPaste?: () => void;
		onDuplicate?: () => void;
		onDelete?: () => void;
		onUndo?: () => void;
		onRedo?: () => void;
		onSave?: () => void;
		onShowHelp?: () => void;
	}

	let { onCut, onCopy, onPaste, onDuplicate, onDelete, onUndo, onRedo, onSave, onShowHelp }: Props =
		$props();

	function registerShortcuts() {
		// Clear existing shortcuts
		keyboardManager.clear();

		// Tools
		keyboardManager.register('v', 'Select Tool', 'tools', () => {
			canvasStore.setTool('select');
		});

		keyboardManager.register('t', 'Text Tool', 'tools', () => {
			canvasStore.setTool('text');
		});

		keyboardManager.register('r', 'Rectangle Tool', 'tools', () => {
			canvasStore.setTool('rectangle');
		});

		keyboardManager.register('c', 'Circle Tool', 'tools', () => {
			canvasStore.setTool('ellipse');
		});

		keyboardManager.register('p', 'Pen Tool', 'tools', () => {
			canvasStore.setTool('freehand');
		});

		keyboardManager.register('s', 'Sticky Note', 'tools', () => {
			canvasStore.setTool('sticky');
		});

		keyboardManager.register('i', 'Image Tool', 'tools', () => {
			canvasStore.setTool('image');
		});

		keyboardManager.register('l', 'List Tool', 'tools', () => {
			canvasStore.setTool('list');
		});

		// Edit operations
		keyboardManager.register('ctrl+z', 'Undo', 'edit', () => {
			onUndo?.();
		});

		keyboardManager.register('ctrl+shift+z', 'Redo', 'edit', () => {
			onRedo?.();
		});

		keyboardManager.register('ctrl+y', 'Redo (alternative)', 'edit', () => {
			onRedo?.();
		});

		keyboardManager.register('ctrl+c', 'Copy', 'edit', () => {
			onCopy?.();
		});

		keyboardManager.register('ctrl+x', 'Cut', 'edit', () => {
			onCut?.();
		});

		keyboardManager.register('ctrl+v', 'Paste', 'edit', () => {
			onPaste?.();
		});

		keyboardManager.register('ctrl+d', 'Duplicate', 'edit', () => {
			onDuplicate?.();
		});

		keyboardManager.register('del', 'Delete', 'edit', () => {
			onDelete?.();
		});

		keyboardManager.register('backspace', 'Delete (alternative)', 'edit', () => {
			onDelete?.();
		});

		// Selection
		keyboardManager.register('ctrl+a', 'Select All', 'selection', () => {
			canvasStore.selectAll();
		});

		keyboardManager.register('esc', 'Clear Selection', 'selection', () => {
			canvasStore.clearSelection();
		});

		// Layers
		keyboardManager.register('ctrl+]', 'Bring Forward', 'layers', () => {
			canvasStore.bringForward();
		});

		keyboardManager.register('ctrl+[', 'Send Backward', 'layers', () => {
			canvasStore.sendBackward();
		});

		keyboardManager.register('ctrl+shift+]', 'Bring to Front', 'layers', () => {
			canvasStore.bringToFront();
		});

		keyboardManager.register('ctrl+shift+[', 'Send to Back', 'layers', () => {
			canvasStore.sendToBack();
		});

		// Grouping
		keyboardManager.register('ctrl+g', 'Group', 'grouping', () => {
			canvasStore.groupSelected();
		});

		keyboardManager.register('ctrl+shift+g', 'Ungroup', 'grouping', () => {
			canvasStore.ungroupSelected();
		});

		// View
		keyboardManager.register('ctrl+0', 'Reset Zoom', 'view', () => {
			canvasStore.resetZoom();
		});

		keyboardManager.register('ctrl+=', 'Zoom In', 'view', () => {
			canvasStore.zoom(0.1);
		});

		keyboardManager.register('ctrl+-', 'Zoom Out', 'view', () => {
			canvasStore.zoom(-0.1);
		});

		keyboardManager.register("ctrl+'", 'Toggle Grid', 'view', () => {
			canvasStore.toggleGrid();
		});

		keyboardManager.register("ctrl+shift+'", 'Toggle Snap to Grid', 'view', () => {
			canvasStore.toggleSnap();
		});

		// Other
		keyboardManager.register('ctrl+s', 'Save', 'other', () => {
			onSave?.();
		});

		keyboardManager.register('?', 'Show Keyboard Shortcuts', 'other', () => {
			onShowHelp?.();
		});

		keyboardManager.register('shift+?', 'Show Keyboard Shortcuts (alternative)', 'other', () => {
			onShowHelp?.();
		});
	}

	function handleKeyDown(event: KeyboardEvent) {
		keyboardManager.handleKeyDown(event);
	}

	onMount(() => {
		registerShortcuts();
		window.addEventListener('keydown', handleKeyDown);
	});

	onDestroy(() => {
		window.removeEventListener('keydown', handleKeyDown);
		keyboardManager.clear();
	});
</script>

<!-- This component is headless - it only registers keyboard shortcuts -->
