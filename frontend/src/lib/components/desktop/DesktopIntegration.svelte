<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { isTauri, parseWorkspaceElements } from '$lib/services/tauri';
	import UpdateChecker from './UpdateChecker.svelte';
	import { canvasStore } from '$lib/stores/canvas.svelte';

	let mounted = $state(false);

	async function handleOpenFile() {
		if (!isTauri()) return;

		try {
			const { openFileDialog, importWorkspace: importWs } = await import('$lib/services/tauri');
			const filePath = await openFileDialog();

			if (filePath) {
				const workspaceData = await importWs(filePath);
				const elements = parseWorkspaceElements(workspaceData);

				// Navigate to a new workspace or load into current
				// For now, we'll load into canvas store
				canvasStore.setElements(elements);

				console.log('Imported workspace:', workspaceData.name);
			}
		} catch (error) {
			console.error('Failed to import workspace:', error);
		}
	}

	async function handleSaveFile() {
		if (!isTauri()) return;

		try {
			const { saveFileDialog, exportWorkspace } = await import('$lib/services/tauri');
			const defaultName = 'workspace.hertzboard';
			const filePath = await saveFileDialog(defaultName);

			if (filePath) {
				const elements = canvasStore.elements;
				await exportWorkspace(filePath, 'workspace-id', 'My Workspace', elements);

				console.log('Exported workspace to:', filePath);
			}
		} catch (error) {
			console.error('Failed to export workspace:', error);
		}
	}

	function handleTauriEvents() {
		// Listen for custom events from Tauri global shortcuts
		window.addEventListener('tauri-open-file', handleOpenFile);
		window.addEventListener('tauri-save-file', handleSaveFile);
	}

	function cleanup() {
		window.removeEventListener('tauri-open-file', handleOpenFile);
		window.removeEventListener('tauri-save-file', handleSaveFile);
	}

	onMount(() => {
		mounted = true;

		if (isTauri()) {
			handleTauriEvents();

			// Log that we're running in Tauri
			console.log('🖥️ Running in Tauri desktop mode');
		}
	});

	onDestroy(() => {
		cleanup();
	});
</script>

{#if mounted && isTauri()}
	<UpdateChecker />
{/if}
