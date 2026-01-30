<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { workspaceStore } from '$lib/stores/workspace.svelte';
	import { i18n } from '$lib/stores/i18n.svelte';
	import {
		Plus,
		MoreVertical,
		Users,
		Calendar,
		Search,
		Edit2,
		Copy,
		Share2,
		Trash2
	} from 'lucide-svelte';
	import type { Workspace } from '$lib/types/api';

	let showCreateModal = $state(false);
	let createName = $state('');
	let createDescription = $state('');
	let createError = $state('');
	let isCreating = $state(false);

	let showRenameModal = $state(false);
	let renameWorkspace = $state<Workspace | null>(null);
	let renameName = $state('');
	let renameDescription = $state('');
	let renameError = $state('');
	let isRenaming = $state(false);

	let showDuplicateModal = $state(false);
	let duplicateWorkspace = $state<Workspace | null>(null);
	let duplicateName = $state('');
	let duplicateError = $state('');
	let isDuplicating = $state(false);

	let activeMenuId = $state<string | null>(null);
	let searchQuery = $state('');

	// Helper function to get workspace role
	function getWorkspaceRole(workspace: Workspace): string {
		return workspace.user_role || workspace.role || 'viewer';
	}

	onMount(async () => {
		try {
			await workspaceStore.loadWorkspaces();
		} catch (error) {
			console.error('Failed to load workspaces:', error);
			// If unauthorized, the layout will handle redirect
		}
	});

	async function handleCreateWorkspace(e: Event) {
		e.preventDefault();
		createError = '';
		isCreating = true;

		try {
			const workspace = await workspaceStore.createWorkspace({
				name: createName,
				description: createDescription || undefined
			});
			showCreateModal = false;
			createName = '';
			createDescription = '';
			// Navigate to the workspace
			goto(`/workspace/${workspace.id}`);
		} catch (err) {
			createError = err instanceof Error ? err.message : 'Failed to create workspace';
		} finally {
			isCreating = false;
		}
	}

	async function handleDeleteWorkspace(id: string) {
		if (!confirm(i18n.t('dashboard.alerts.deleteConfirm'))) {
			return;
		}

		try {
			await workspaceStore.deleteWorkspace(id);
			activeMenuId = null;
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete workspace');
		}
	}

	function handleDuplicateWorkspace(workspace: Workspace) {
		duplicateWorkspace = workspace;
		duplicateName = `${workspace.name} (Copy)`;
		showDuplicateModal = true;
		activeMenuId = null;
	}

	async function submitDuplicate(e: Event) {
		e.preventDefault();
		if (!duplicateWorkspace) return;

		duplicateError = '';
		isDuplicating = true;

		try {
			await workspaceStore.duplicateWorkspace(duplicateWorkspace.id, duplicateName);
			showDuplicateModal = false;
			duplicateWorkspace = null;
			duplicateName = '';
		} catch (err) {
			duplicateError = err instanceof Error ? err.message : 'Failed to duplicate workspace';
		} finally {
			isDuplicating = false;
		}
	}

	function handleRenameWorkspace(workspace: Workspace) {
		renameWorkspace = workspace;
		renameName = workspace.name;
		renameDescription = workspace.description || '';
		showRenameModal = true;
		activeMenuId = null;
	}

	async function submitRename(e: Event) {
		e.preventDefault();
		if (!renameWorkspace) return;

		renameError = '';
		isRenaming = true;

		try {
			await workspaceStore.updateWorkspace(renameWorkspace.id, {
				name: renameName,
				description: renameDescription || undefined
			});
			showRenameModal = false;
			renameWorkspace = null;
			renameName = '';
			renameDescription = '';
		} catch (err) {
			renameError = err instanceof Error ? err.message : 'Failed to rename workspace';
		} finally {
			isRenaming = false;
		}
	}

	function handleShareWorkspace(_workspace: Workspace) {
		// TODO: Implement share modal in future phase
		alert(i18n.t('dashboard.alerts.shareComingSoon'));
		activeMenuId = null;
	}

	function handleCopyLink(workspace: Workspace) {
		const url = `${window.location.origin}/workspace/${workspace.id}`;
		navigator.clipboard.writeText(url).then(() => {
			alert(i18n.t('dashboard.alerts.linkCopied'));
			activeMenuId = null;
		});
	}

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		const days = Math.floor(diff / (1000 * 60 * 60 * 24));

		if (days === 0) return i18n.t('dashboard.time.today');
		if (days === 1) return i18n.t('dashboard.time.yesterday');
		if (days < 7) return i18n.t('dashboard.time.daysAgo', { count: days.toString() });
		if (days < 30)
			return i18n.t('dashboard.time.weeksAgo', {
				count: Math.floor(days / 7).toString()
			});
		if (days < 365)
			return i18n.t('dashboard.time.monthsAgo', {
				count: Math.floor(days / 30).toString()
			});
		return date.toLocaleDateString();
	}

	function handleSearch() {
		workspaceStore.loadWorkspaces({ query: searchQuery || undefined, offset: 0 });
	}

	$effect(() => {
		// Debounced search
		const timer = setTimeout(handleSearch, 300);
		return () => clearTimeout(timer);
	});
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">{i18n.t('dashboard.title')}</h1>
			<p class="mt-1 text-gray-600">{i18n.t('dashboard.subtitle')}</p>
		</div>
		<button
			onclick={() => (showCreateModal = true)}
			class="flex cursor-pointer items-center gap-2 bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700"
		>
			<Plus size={20} />
			{i18n.t('dashboard.newWorkspace')}
		</button>
	</div>

	<!-- Search -->
	<div class="relative">
		<Search class="absolute top-1/2 left-3 -translate-y-1/2 text-gray-400" size={20} />
		<input
			type="text"
			bind:value={searchQuery}
			placeholder={i18n.t('dashboard.searchPlaceholder')}
			class="w-full border border-gray-300 py-2 pr-4 pl-10 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
		/>
	</div>

	<!-- Workspaces Grid -->
	{#if workspaceStore.isLoading && workspaceStore.workspaces.length === 0}
		<div class="flex items-center justify-center py-12">
			<div class="text-center">
				<div
					class="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent"
				></div>
				<p class="text-gray-600">{i18n.t('dashboard.loading')}</p>
			</div>
		</div>
	{:else if workspaceStore.workspaces.length === 0}
		<div class="flex flex-col items-center justify-center py-12">
			<p class="mb-4 text-gray-600">{i18n.t('dashboard.noWorkspaces')}</p>
			<button
				onclick={() => (showCreateModal = true)}
				class="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700"
			>
				<Plus size={20} />
				{i18n.t('dashboard.createFirst')}
			</button>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each workspaceStore.workspaces as workspace (workspace.id)}
				<div
					class="group relative flex flex-col border border-gray-200 bg-white shadow-sm transition hover:shadow-md"
				>
					<!-- Thumbnail -->
					<button
						onclick={() => goto(`/workspace/${workspace.id}`)}
						class="aspect-video w-full overflow-hidden bg-gradient-to-br from-blue-50 to-indigo-100"
					>
						{#if workspace.thumbnail_url}
							<img
								src={workspace.thumbnail_url}
								alt={workspace.name}
								class="h-full w-full object-cover"
							/>
						{:else}
							<div class="flex h-full items-center justify-center text-4xl font-bold text-gray-300">
								{workspace.name?.charAt(0)?.toUpperCase() || 'W'}
							</div>
						{/if}
					</button>

					<!-- Content -->
					<div class="flex flex-1 flex-col p-4">
						<button onclick={() => goto(`/workspace/${workspace.id}`)} class="mb-2 text-left">
							<h3 class="line-clamp-1 font-semibold text-gray-900">{workspace.name}</h3>
							{#if workspace.description}
								<p class="mt-1 line-clamp-2 text-sm text-gray-600">{workspace.description}</p>
							{/if}
						</button>

						<div class="mt-auto space-y-2 pt-4">
							<div class="flex items-center gap-2 text-xs text-gray-500">
								<Calendar size={14} />
								{formatDate(workspace.updated_at)}
							</div>
							{#if workspace.member_count}
								<div class="flex items-center gap-2 text-xs text-gray-500">
									<Users size={14} />
									{workspace.member_count}
									{workspace.member_count === 1
										? i18n.t('dashboard.member')
										: i18n.t('dashboard.members')}
								</div>
							{/if}
							<div class="text-xs text-gray-500">
								{i18n.t('dashboard.role')}:
								<span class="font-medium capitalize">{getWorkspaceRole(workspace)}</span>
							</div>
						</div>
					</div>

					<!-- Menu -->
					<div class="absolute top-2 right-2">
						<button
							onclick={(e) => {
								e.stopPropagation();
								activeMenuId = activeMenuId === workspace.id ? null : workspace.id;
							}}
							class="bg-white p-1 opacity-0 shadow-sm transition group-hover:opacity-100 hover:bg-gray-50"
						>
							<MoreVertical size={16} />
						</button>

						{#if activeMenuId === workspace.id}
							<div
								class="ring-opacity-5 absolute top-8 right-0 z-10 w-56 bg-white shadow-lg ring-1 ring-black"
							>
								<div class="py-1">
									<button
										onclick={() => {
											activeMenuId = null;
											goto(`/workspace/${workspace.id}`);
										}}
										class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
									>
										<Search size={16} />
										{i18n.t('dashboard.menu.open')}
									</button>

									{#if getWorkspaceRole(workspace) === 'owner' || getWorkspaceRole(workspace) === 'editor'}
										<button
											onclick={() => handleRenameWorkspace(workspace)}
											class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
										>
											<Edit2 size={16} />
											{i18n.t('dashboard.menu.rename')}
										</button>
									{/if}

									<button
										onclick={() => handleCopyLink(workspace)}
										class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
									>
										<Copy size={16} />
										{i18n.t('dashboard.menu.copyLink')}
									</button>

									{#if getWorkspaceRole(workspace) === 'owner' || getWorkspaceRole(workspace) === 'editor'}
										<button
											onclick={() => handleShareWorkspace(workspace)}
											class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
										>
											<Share2 size={16} />
											{i18n.t('dashboard.menu.share')}
										</button>
									{/if}

									{#if getWorkspaceRole(workspace) === 'owner'}
										<div class="my-1 border-t border-gray-200"></div>

										<button
											onclick={() => handleDuplicateWorkspace(workspace)}
											class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
										>
											<Copy size={16} />
											{i18n.t('dashboard.menu.duplicate')}
										</button>

										<button
											onclick={() => handleDeleteWorkspace(workspace.id)}
											class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-red-600 hover:bg-gray-100"
										>
											<Trash2 size={16} />
											{i18n.t('dashboard.menu.delete')}
										</button>
									{/if}
								</div>
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Create Workspace Modal -->
{#if showCreateModal}
	<div class="bg-opacity-50 fixed inset-0 z-50 flex items-center justify-center bg-black">
		<div class="w-full max-w-md rounded-lg bg-white p-6 shadow-xl">
			<h2 class="mb-4 text-xl font-bold text-gray-900">{i18n.t('dashboard.modal.create.title')}</h2>

			<form onsubmit={handleCreateWorkspace}>
				{#if createError}
					<div class="mb-4 rounded-md bg-red-50 p-3">
						<p class="text-sm text-red-800">{createError}</p>
					</div>
				{/if}

				<div class="space-y-4">
					<div>
						<label for="name" class="block text-sm font-medium text-gray-700"
							>{i18n.t('dashboard.modal.create.name')}</label
						>
						<input
							id="name"
							type="text"
							required
							bind:value={createName}
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
							placeholder={i18n.t('dashboard.modal.create.namePlaceholder')}
						/>
					</div>

					<div>
						<label for="description" class="block text-sm font-medium text-gray-700">
							{i18n.t('dashboard.modal.create.description')}
						</label>
						<textarea
							id="description"
							bind:value={createDescription}
							rows="3"
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
							placeholder={i18n.t('dashboard.modal.create.descriptionPlaceholder')}
						></textarea>
					</div>
				</div>

				<div class="mt-6 flex gap-3">
					<button
						type="button"
						onclick={() => {
							showCreateModal = false;
							createName = '';
							createDescription = '';
							createError = '';
						}}
						class="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-gray-700 transition hover:bg-gray-50"
					>
						{i18n.t('common.cancel')}
					</button>
					<button
						type="submit"
						disabled={isCreating}
						class="flex-1 rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
					>
						{isCreating
							? i18n.t('dashboard.modal.create.creating')
							: i18n.t('dashboard.modal.create.create')}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Duplicate Workspace Modal -->
{#if showDuplicateModal && duplicateWorkspace}
	<div class="bg-opacity-50 fixed inset-0 z-50 flex items-center justify-center bg-black">
		<div class="w-full max-w-md rounded-lg bg-white p-6 shadow-xl">
			<h2 class="mb-4 text-xl font-bold text-gray-900">
				{i18n.t('dashboard.modal.duplicate.title')}
			</h2>

			<form onsubmit={submitDuplicate}>
				{#if duplicateError}
					<div class="mb-4 rounded-md bg-red-50 p-3">
						<p class="text-sm text-red-800">{duplicateError}</p>
					</div>
				{/if}

				<div class="mb-4">
					<label for="duplicate-name" class="block text-sm font-medium text-gray-700">
						{i18n.t('dashboard.modal.duplicate.newName')}
					</label>
					<input
						id="duplicate-name"
						type="text"
						required
						bind:value={duplicateName}
						class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
						placeholder="My Workspace (Copy)"
					/>
					<p class="mt-1 text-xs text-gray-500">
						{i18n.t('dashboard.modal.duplicate.copyOf', { name: duplicateWorkspace.name })}
					</p>
				</div>

				<div class="mt-6 flex gap-3">
					<button
						type="button"
						onclick={() => {
							showDuplicateModal = false;
							duplicateWorkspace = null;
							duplicateName = '';
							duplicateError = '';
						}}
						class="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-gray-700 transition hover:bg-gray-50"
					>
						{i18n.t('common.cancel')}
					</button>
					<button
						type="submit"
						disabled={isDuplicating}
						class="flex-1 rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
					>
						{isDuplicating
							? i18n.t('dashboard.modal.duplicate.duplicating')
							: i18n.t('dashboard.modal.duplicate.duplicate')}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Rename Workspace Modal -->
{#if showRenameModal && renameWorkspace}
	<div class="bg-opacity-50 fixed inset-0 z-50 flex items-center justify-center bg-black">
		<div class="w-full max-w-md rounded-lg bg-white p-6 shadow-xl">
			<h2 class="mb-4 text-xl font-bold text-gray-900">{i18n.t('dashboard.modal.rename.title')}</h2>

			<form onsubmit={submitRename}>
				{#if renameError}
					<div class="mb-4 rounded-md bg-red-50 p-3">
						<p class="text-sm text-red-800">{renameError}</p>
					</div>
				{/if}

				<div class="space-y-4">
					<div>
						<label for="rename-name" class="block text-sm font-medium text-gray-700"
							>{i18n.t('dashboard.modal.create.name')}</label
						>
						<input
							id="rename-name"
							type="text"
							required
							bind:value={renameName}
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
							placeholder={i18n.t('dashboard.modal.create.namePlaceholder')}
						/>
					</div>

					<div>
						<label for="rename-description" class="block text-sm font-medium text-gray-700">
							{i18n.t('dashboard.modal.create.description')}
						</label>
						<textarea
							id="rename-description"
							bind:value={renameDescription}
							rows="3"
							class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
							placeholder={i18n.t('dashboard.modal.create.descriptionPlaceholder')}
						></textarea>
					</div>
				</div>

				<div class="mt-6 flex gap-3">
					<button
						type="button"
						onclick={() => {
							showRenameModal = false;
							renameWorkspace = null;
							renameName = '';
							renameDescription = '';
							renameError = '';
						}}
						class="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-gray-700 transition hover:bg-gray-50"
					>
						{i18n.t('common.cancel')}
					</button>
					<button
						type="submit"
						disabled={isRenaming}
						class="flex-1 rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
					>
						{isRenaming
							? i18n.t('dashboard.modal.rename.saving')
							: i18n.t('dashboard.modal.rename.saveChanges')}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
