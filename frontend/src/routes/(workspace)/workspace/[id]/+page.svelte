<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';
	import { api } from '$lib/services/api';
	import { canvasStore } from '$lib/stores/canvas.svelte';
	import type { Workspace } from '$lib/types/api';
	import { ArrowLeft, Users } from 'lucide-svelte';
	import Canvas from '$lib/components/canvas/Canvas.svelte';
	import CanvasToolbar from '$lib/components/canvas/CanvasToolbar.svelte';

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

			// Загружаем элементы canvas
			try {
				const elements = await api.listElements(workspaceId);
				canvasStore.setElements(elements || []);
			} catch {
				// Если нет элементов или ошибка загрузки, создаем демо-контент
				createDemoElements();
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load workspace';
		} finally {
			isLoading = false;
		}
	});

	onDestroy(() => {
		// Очищаем canvas store при выходе
		canvasStore.reset();
	});

	function createDemoElements() {
		// Создаем несколько демо-элементов для тестирования
		const demoElements = [
			{
				id: crypto.randomUUID(),
				workspace_id: workspaceId!,
				type: 'text' as const,
				pos_x: 100,
				pos_y: 100,
				width: 300,
				height: 100,
				content: 'Welcome to HertzBoard! 🎨',
				z_index: 0,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString(),
				style: {
					fontSize: 32,
					fontFamily: 'Inter, sans-serif',
					color: '#1f2937',
					textAlign: 'center' as const
				}
			},
			{
				id: crypto.randomUUID(),
				workspace_id: workspaceId!,
				type: 'rectangle' as const,
				pos_x: 150,
				pos_y: 250,
				width: 200,
				height: 150,
				z_index: 1,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString(),
				style: {
					backgroundColor: '#3b82f6',
					strokeColor: '#1e40af',
					strokeWidth: 2,
					borderRadius: 8
				}
			},
			{
				id: crypto.randomUUID(),
				workspace_id: workspaceId!,
				type: 'ellipse' as const,
				pos_x: 400,
				pos_y: 250,
				width: 180,
				height: 180,
				z_index: 2,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString(),
				style: {
					backgroundColor: '#10b981',
					strokeColor: '#059669',
					strokeWidth: 3
				}
			},
			{
				id: crypto.randomUUID(),
				workspace_id: workspaceId!,
				type: 'sticky' as const,
				pos_x: 100,
				pos_y: 450,
				width: 200,
				height: 200,
				content: 'Try dragging me!\n\nShift+Drag to lock axis\nShift+Resize to keep aspect ratio',
				z_index: 3,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString(),
				style: {
					backgroundColor: '#fef3c7',
					fontSize: 14,
					color: '#92400e'
				}
			},
			{
				id: crypto.randomUUID(),
				workspace_id: workspaceId!,
				type: 'text' as const,
				pos_x: 350,
				pos_y: 480,
				width: 250,
				height: 80,
				content: 'Phase 6: Canvas Engine ✅',
				z_index: 4,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString(),
				style: {
					fontSize: 20,
					fontFamily: 'Inter, sans-serif',
					color: '#7c3aed',
					textAlign: 'left'
				}
			},
			{
				id: crypto.randomUUID(),
				workspace_id: workspaceId!,
				type: 'arrow' as const,
				pos_x: 350,
				pos_y: 350,
				width: 150,
				height: 100,
				z_index: 5,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString(),
				style: {
					strokeColor: '#ef4444',
					strokeWidth: 3
				}
			}
		];

		canvasStore.setElements(demoElements as any[]);
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

			<div class="flex items-center gap-2">
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
		<CanvasToolbar />

		<!-- Canvas Area (takes remaining space) -->
		<div class="min-h-0 flex-1">
			<Canvas />
		</div>
	</div>
{/if}
