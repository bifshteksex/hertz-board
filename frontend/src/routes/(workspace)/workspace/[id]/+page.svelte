<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';
	import { api } from '$lib/services/api';
	import { canvasStore } from '$lib/stores/canvas.svelte';
	import { canvas } from '$lib/stores/canvasWithHistory.svelte';
	import { historyStore } from '$lib/stores/history.svelte';
	import type { Workspace } from '$lib/types/api';
	import { ArrowLeft, Users, Layers as LayersIcon } from 'lucide-svelte';
	import Canvas from '$lib/components/canvas/Canvas.svelte';
	import CanvasToolbar from '$lib/components/canvas/CanvasToolbar.svelte';
	import PropertiesPanel from '$lib/components/canvas/PropertiesPanel.svelte';
	import LayersPanel from '$lib/components/canvas/LayersPanel.svelte';
	import KeyboardShortcuts from '$lib/components/canvas/KeyboardShortcuts.svelte';
	import ShortcutsPanel from '$lib/components/canvas/ShortcutsPanel.svelte';
	import SaveStatus from '$lib/components/workspace/SaveStatus.svelte';

	let showLayersPanel = $state(false);
	let showShortcutsHelp = $state(false);

	// Derived state для undo/redo
	const canUndo = $derived(historyStore.canUndo);
	const canRedo = $derived(historyStore.canRedo);

	const workspaceId = $derived($page.params.id);
	let workspace = $state<Workspace | null>(null);
	let isLoading = $state(true);
	let error = $state('');

	onMount(async () => {
		if (!workspaceId) {
			error = 'Workspace ID is missing';
			isLoading = false;
			return;
		}

		try {
			// Загружаем workspace
			workspace = await api.getWorkspace(workspaceId);

			// Устанавливаем workspace в canvas store
			canvasStore.setWorkspaceId(workspaceId);

			// Инициализируем автосохранение
			canvas.initAutosave(workspaceId);

			// Загружаем элементы canvas
			try {
				console.log('[Workspace] Loading elements for workspace:', workspaceId);
				const response = await api.listElements(workspaceId);
				console.log('[Workspace] Loaded elements:', response);
				const backendElements = response?.elements || [];

				// Преобразуем из формата backend в формат frontend
				const elements = backendElements.map((el: any) => {
					const elementData = el.element_data || {};
					const position = elementData.position || {};
					const size = elementData.size || {};

					// Определяем тип элемента
					let elementType = el.element_type;
					if (el.element_type === 'shape') {
						// Для shape берем shape_type из element_data, если есть
						elementType = elementData.shape_type || 'rectangle'; // дефолт - rectangle
					}

					return {
						id: el.id,
						workspace_id: el.workspace_id,
						type: elementType,
						pos_x: position.x || 0,
						pos_y: position.y || 0,
						width: size.width || 0,
						height: size.height || 0,
						rotation: elementData.rotation || 0,
						z_index: el.z_index || 0,
						content: elementData.content,
						html_content: elementData.html_content,
						style: elementData.style || {},
						parent_id: el.parent_id,
						created_at: el.created_at,
						updated_at: el.updated_at,
						created_by: el.created_by,
						updated_by: el.updated_by
					};
				});

				canvasStore.setElements(elements);
				console.log('[Workspace] Elements set in store, count:', elements.length);
			} catch (err) {
				console.error('[Workspace] Failed to load elements:', err);
				// Начинаем с пустого workspace
				canvasStore.setElements([]);
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load workspace';
		} finally {
			isLoading = false;
		}
	});

	onDestroy(() => {
		// Сохраняем все pending изменения перед выходом
		canvas.saveNow();

		// Очищаем canvas store, history и autosave при выходе
		canvas.stopAutosave();
		canvasStore.reset();
		historyStore.clear();
	});

	// Keyboard shortcuts handlers
	function handleUndo() {
		canvas.undo();
	}

	function handleRedo() {
		canvas.redo();
	}

	function handleCopy() {
		// TODO: implement copy
		console.log('Copy not implemented yet');
	}

	function handleCut() {
		// TODO: implement cut
		console.log('Cut not implemented yet');
	}

	function handlePaste() {
		// TODO: implement paste
		console.log('Paste not implemented yet');
	}

	function handleDuplicate() {
		// TODO: implement duplicate
		console.log('Duplicate not implemented yet');
	}

	function handleDelete() {
		const selectedIds = canvasStore.selectedIds;
		if (selectedIds.length > 0) {
			canvas.deleteElements(selectedIds);
		}
	}

	function handleSave() {
		// TODO: implement auto-save or force save
		console.log('Save triggered');
	}

	function handleShowHelp() {
		showShortcutsHelp = true;
	}
</script>

{#if isLoading}
	<div class="flex h-screen items-center justify-center">
		<div class="text-center">
			<div
				class="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent"
			></div>
			<p class="text-gray-600">Loading workspace...</p>
		</div>
	</div>
{:else if error}
	<div class="flex h-screen items-center justify-center">
		<div class="text-center">
			<p class="mb-4 text-red-600">{error}</p>
			<button
				onclick={() => goto('/dashboard')}
				class="rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700"
			>
				Back to Dashboard
			</button>
		</div>
	</div>
{:else if workspace}
	<div class="flex h-screen flex-col">
		<!-- Workspace Header -->
		<header
			class="flex shrink-0 items-center justify-between border-b border-gray-200 bg-white px-6 py-3"
		>
			<div class="flex items-center gap-4">
				<button
					onclick={() => goto('/dashboard')}
					class="rounded-lg p-2 text-gray-600 transition hover:bg-gray-100"
					title="Back to dashboard"
				>
					<ArrowLeft size={20} />
				</button>
				<div>
					<h1 class="text-lg font-semibold text-gray-900">{workspace.name}</h1>
					{#if workspace.description}
						<p class="text-sm text-gray-600">{workspace.description}</p>
					{/if}
				</div>
			</div>

			<div class="flex items-center gap-4">
				<!-- Save Status Indicator -->
				<SaveStatus />

				<div class="h-6 w-px bg-gray-300"></div>

				<button
					class="flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-2 text-sm transition hover:bg-gray-50"
					class:bg-blue-50={showLayersPanel}
					class:border-blue-500={showLayersPanel}
					onclick={() => (showLayersPanel = !showLayersPanel)}
					title="Toggle layers panel"
				>
					<LayersIcon size={16} />
					Layers
				</button>
				<button
					class="flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-2 text-sm transition hover:bg-gray-50"
					title="Share workspace"
				>
					<Users size={16} />
					Share
				</button>
			</div>
		</header>

		<!-- Toolbar -->
		<CanvasToolbar
			{canUndo}
			{canRedo}
			onUndo={handleUndo}
			onRedo={handleRedo}
			onShowHelp={handleShowHelp}
		/>

		<!-- Keyboard Shortcuts Component (headless) -->
		<KeyboardShortcuts
			onCut={handleCut}
			onCopy={handleCopy}
			onPaste={handlePaste}
			onDuplicate={handleDuplicate}
			onDelete={handleDelete}
			onUndo={handleUndo}
			onRedo={handleRedo}
			onSave={handleSave}
			onShowHelp={handleShowHelp}
		/>

		<!-- Shortcuts Help Panel (modal) -->
		{#if showShortcutsHelp}
			<ShortcutsPanel onClose={() => (showShortcutsHelp = false)} />
		{/if}

		<!-- Canvas Area with panels -->
		<div class="flex min-h-0 flex-1">
			<!-- Layers Panel (optional, left side) -->
			{#if showLayersPanel}
				<LayersPanel />
			{/if}

			<!-- Canvas (center, takes remaining space) -->
			<div class="min-h-0 flex-1">
				<Canvas />
			</div>

			<!-- Properties Panel (right side, always visible) -->
			<PropertiesPanel />
		</div>
	</div>
{/if}
