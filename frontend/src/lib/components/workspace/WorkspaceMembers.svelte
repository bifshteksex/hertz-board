<script lang="ts">
	import { api } from '$lib/services/api';
	import type { WorkspaceMember, WorkspaceRole } from '$lib/types/api';
	import { Users, Mail, MoreVertical, UserX, Eye, Edit3, Crown, Copy, Check } from 'lucide-svelte';
	import { i18n } from '$lib/stores/i18n.svelte';

	interface Props {
		workspaceId: string;
		currentUserRole?: WorkspaceRole;
	}

	let { workspaceId, currentUserRole = 'viewer' }: Props = $props();

	let members = $state<WorkspaceMember[]>([]);
	let isLoading = $state(true);
	let error = $state('');

	// Invite modal
	let showInviteModal = $state(false);
	let inviteEmail = $state('');
	let inviteRole = $state<WorkspaceRole>('editor');
	let inviteError = $state('');
	let isInviting = $state(false);
	let inviteSuccess = $state(false);
	let inviteUrl = $state('');

	// Menu state
	let activeMenuId = $state<string | null>(null);

	// Копирование ссылки
	let copiedLink = $state(false);

	const canManageMembers = $derived(currentUserRole === 'owner' || currentUserRole === 'editor');
	const canInvite = $derived(canManageMembers);

	$effect(() => {
		loadMembers();
	});

	async function loadMembers() {
		try {
			isLoading = true;
			error = '';
			members = await api.listMembers(workspaceId);
			console.log('[WorkspaceMembers] Loaded members:', members);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load members';
			console.error('[WorkspaceMembers] Failed to load members:', err);
		} finally {
			isLoading = false;
		}
	}

	async function handleInvite(e: Event) {
		e.preventDefault();
		inviteError = '';
		inviteSuccess = false;
		isInviting = true;

		try {
			const response = await api.inviteMember(workspaceId, {
				email: inviteEmail,
				role: inviteRole
			});

			// Успех!
			inviteSuccess = true;

			// Формируем полный URL
			const baseUrl = window.location.origin;
			inviteUrl = `${baseUrl}/workspace/invite?token=${response.token}`;

			// Очищаем форму
			inviteEmail = '';
			inviteRole = 'editor';
		} catch (err) {
			inviteError = err instanceof Error ? err.message : 'Failed to send invitation';
		} finally {
			isInviting = false;
		}
	}

	async function handleUpdateRole(memberId: string, newRole: WorkspaceRole) {
		try {
			await api.updateMemberRole(workspaceId, memberId, { role: newRole });
			await loadMembers();
			activeMenuId = null;
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update role');
		}
	}

	async function handleRemoveMember(memberId: string) {
		if (!confirm(i18n.t('workspaceMembers.alerts.removeConfirm'))) {
			return;
		}

		try {
			await api.removeMember(workspaceId, memberId);
			await loadMembers();
			activeMenuId = null;
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to remove member');
		}
	}

	function getRoleColor(role: WorkspaceRole) {
		switch (role) {
			case 'owner':
				return '#f59e0b';
			case 'editor':
				return '#3b82f6';
			case 'viewer':
				return '#6b7280';
		}
	}

	function getRoleLabel(role: WorkspaceRole) {
		return i18n.t(`workspaceMembers.${role}`);
	}

	function getInitials(name?: string, email?: string): string {
		const displayName = name || email || '?';
		return displayName
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase()
			.slice(0, 2);
	}

	async function copyInviteLink() {
		try {
			await navigator.clipboard.writeText(inviteUrl);
			copiedLink = true;
			setTimeout(() => {
				copiedLink = false;
			}, 2000);
		} catch (err) {
			console.error('Failed to copy:', err);
		}
	}

	function closeInviteModal() {
		showInviteModal = false;
		inviteSuccess = false;
		inviteUrl = '';
		inviteEmail = '';
		inviteRole = 'editor';
		inviteError = '';
	}
</script>

<div class="overflow-hidden border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
	<div class="flex items-center justify-between border-b border-gray-200 p-4 dark:border-gray-700">
		<div class="flex items-center gap-2 text-gray-900 dark:text-gray-50">
			<Users size={20} />
			<h3 class="m-0 text-base font-semibold">{i18n.t('workspaceMembers.title')}</h3>
			<span
				class="flex h-6 min-w-6 items-center justify-center bg-gray-100 px-2 text-[13px] font-semibold text-gray-500 dark:bg-gray-700 dark:text-gray-400"
				>{members.length}</span
			>
		</div>

		{#if canInvite}
			<button
				class="flex items-center gap-1.5 border-2 border-[#372d2e] bg-[#fcbd80] px-4 py-2 text-sm font-medium text-[#372d2e] transition-all duration-150 hover:-translate-y-0.5 hover:bg-[#dda877] hover:shadow-[0_4px_0_#372d2e] active:translate-y-0 active:shadow-none"
				onclick={() => (showInviteModal = true)}
			>
				<Mail size={16} />
				{i18n.t('workspaceMembers.inviteMember')}
			</button>
		{/if}
	</div>

	{#if isLoading}
		<div class="p-8 text-center text-gray-500 dark:text-gray-400">
			{i18n.t('workspaceMembers.loading')}
		</div>
	{:else if error}
		<div class="p-8 text-center text-red-600 dark:text-red-500">{error}</div>
	{:else if members.length === 0}
		<div
			class="flex flex-col items-center justify-center p-12 px-4 text-gray-400 dark:text-gray-600"
		>
			<Users size={48} strokeWidth={1} />
			<p class="mt-3 text-sm">{i18n.t('workspaceMembers.noMembers')}</p>
		</div>
	{:else}
		<div class="p-2">
			{#each members as member (member.id)}
				<div
					class="relative flex items-center gap-3 rounded-md p-3 transition-colors hover:bg-gray-50 dark:hover:bg-gray-700"
				>
					<!-- Avatar -->
					<div
						class="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-full bg-linear-to-br from-indigo-500 to-purple-600 text-sm font-semibold text-white"
					>
						{#if member.user_avatar}
							<img
								src={member.user_avatar}
								alt={member.user_name}
								class="h-full w-full object-cover"
							/>
						{:else}
							{getInitials(member.user_name, member.user_email)}
						{/if}
					</div>

					<!-- Info -->
					<div class="min-w-0 flex-1">
						<div
							class="overflow-hidden text-sm font-medium text-ellipsis whitespace-nowrap text-gray-900 dark:text-gray-50"
						>
							{member.user_name || 'Unknown User'}
						</div>
						<div
							class="overflow-hidden text-[13px] text-ellipsis whitespace-nowrap text-gray-500 dark:text-gray-400"
						>
							{member.user_email}
						</div>
					</div>

					<!-- Role badge -->
					{#if member.role === 'owner'}
						<div
							class="flex shrink-0 items-center gap-1 bg-gray-50 px-3 py-1 text-[13px] font-medium dark:bg-gray-700"
							style="color: {getRoleColor(member.role)}"
						>
							<Crown size={14} />
							{getRoleLabel(member.role)}
						</div>
					{:else if member.role === 'editor'}
						<div
							class="flex shrink-0 items-center gap-1 bg-gray-50 px-3 py-1 text-[13px] font-medium dark:bg-gray-700"
							style="color: {getRoleColor(member.role)}"
						>
							<Edit3 size={14} />
							{getRoleLabel(member.role)}
						</div>
					{:else}
						<div
							class="flex shrink-0 items-center gap-1 bg-gray-50 px-3 py-1 text-[13px] font-medium dark:bg-gray-700"
							style="color: {getRoleColor(member.role)}"
						>
							<Eye size={14} />
							{getRoleLabel(member.role)}
						</div>
					{/if}

					<!-- Actions menu -->
					{#if canManageMembers && member.role !== 'owner'}
						<div class="relative">
							<button
								class="flex h-7 w-7 items-center justify-center rounded border-none bg-transparent text-gray-500 transition-all hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-50"
								onclick={() => (activeMenuId = activeMenuId === member.id ? null : member.id)}
							>
								<MoreVertical size={16} />
							</button>

							{#if activeMenuId === member.id}
								<div
									class="animate-slideDown absolute top-[calc(100%+4px)] right-0 z-100 min-w-40 rounded-md border border-gray-200 bg-white p-1 shadow-lg dark:border-gray-700 dark:bg-gray-800"
								>
									{#if currentUserRole === 'owner'}
										<button
											class="flex w-full items-center gap-2 rounded border-none bg-transparent px-3 py-2 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-50 dark:text-gray-300 dark:hover:bg-gray-700"
											onclick={() => handleUpdateRole(member.id, 'viewer')}
											disabled={member.role === 'viewer'}
										>
											<Eye size={14} />
											{i18n.t('workspaceMembers.menu.makeViewer')}
										</button>
										<button
											class="flex w-full items-center gap-2 rounded border-none bg-transparent px-3 py-2 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-50 dark:text-gray-300 dark:hover:bg-gray-700"
											onclick={() => handleUpdateRole(member.id, 'editor')}
											disabled={member.role === 'editor'}
										>
											<Edit3 size={14} />
											{i18n.t('workspaceMembers.menu.makeEditor')}
										</button>
										<div class="my-1 h-px bg-gray-200 dark:bg-gray-700"></div>
									{/if}
									<button
										class="flex w-full items-center gap-2 rounded border-none bg-transparent px-3 py-2 text-left text-sm text-red-600 transition-colors hover:bg-red-50 dark:text-red-500 dark:hover:bg-red-950"
										onclick={() => handleRemoveMember(member.id)}
									>
										<UserX size={14} />
										{i18n.t('workspaceMembers.menu.removeMember')}
									</button>
								</div>

								<!-- Backdrop -->
								<!-- svelte-ignore a11y_click_events_have_key_events -->
								<!-- svelte-ignore a11y_no_static_element_interactions -->
								<div class="fixed inset-0 z-99" onclick={() => (activeMenuId = null)}></div>
							{/if}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Invite Modal -->
{#if showInviteModal}
	<div
		class="animate-fadeIn fixed inset-0 z-1000 flex items-center justify-center bg-black/50 p-4"
		onclick={closeInviteModal}
		role="presentation"
	>
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div
			class="animate-slideUp max-h-[90vh] w-full max-w-125 overflow-auto bg-white shadow-xl dark:bg-gray-800"
			onclick={(e) => e.stopPropagation()}
		>
			<div
				class="flex items-center justify-between border-b border-gray-200 p-5 px-6 dark:border-gray-700"
			>
				<h3 class="m-0 text-lg font-semibold text-gray-900 dark:text-gray-50">
					{i18n.t('workspaceMembers.modal.invite.title')}
				</h3>
				<button
					class="flex h-8 w-8 items-center justify-center rounded-md border-none bg-transparent text-2xl text-gray-500 transition-all hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-50"
					onclick={closeInviteModal}>&times;</button
				>
			</div>

			{#if inviteSuccess}
				<div class="p-6">
					<div class="text-center">
						<div
							class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-blue-100 text-blue-500"
						>
							<Mail size={32} />
						</div>
						<h4 class="m-0 mb-2 text-lg font-semibold text-gray-900 dark:text-gray-50">
							{i18n.t('workspaceMembers.modal.success.title')}
						</h4>
						<p class="m-0 mb-4 text-sm text-gray-500 dark:text-gray-400">
							{i18n.t('workspaceMembers.modal.success.description')}
						</p>

						<div class="mb-3 flex gap-2">
							<input
								type="text"
								readonly
								value={inviteUrl}
								class="flex-1 border border-gray-300 bg-gray-50 px-3 py-2.5 text-[13px] text-gray-700 dark:border-gray-600 dark:bg-gray-900 dark:text-gray-300"
							/>
							<button
								class="flex items-center gap-1.5 border-2 border-[#372d2e] bg-[#fcbd80] px-4 py-2.5 text-sm font-medium whitespace-nowrap text-[#372d2e] transition-all duration-150 hover:-translate-y-0.5 hover:bg-[#dda877] hover:shadow-[0_4px_0_#372d2e] active:translate-y-0 active:shadow-none"
								onclick={copyInviteLink}
							>
								{#if copiedLink}
									<Check size={16} />
									{i18n.t('workspaceMembers.modal.success.copied')}
								{:else}
									<Copy size={16} />
									{i18n.t('workspaceMembers.modal.success.copy')}
								{/if}
							</button>
						</div>

						<p class="text-[13px] text-gray-400 dark:text-gray-600">
							{i18n.t('workspaceMembers.modal.success.expiresIn')}
						</p>
					</div>
				</div>
				<div
					class="flex items-center justify-end gap-3 border-t border-gray-200 p-4 px-6 dark:border-gray-700"
				>
					<button
						class="border-2 border-[#372d2e] bg-white px-5 py-2.5 text-sm font-medium text-[#372d2e] transition-all duration-150 hover:-translate-y-0.5 hover:bg-gray-50 hover:shadow-[0_4px_0_#372d2e] active:translate-y-0 active:shadow-none dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-gray-700"
						onclick={closeInviteModal}>{i18n.t('workspaceMembers.modal.success.done')}</button
					>
				</div>
			{:else}
				<form onsubmit={handleInvite}>
					<div class="p-6">
						{#if inviteError}
							<div
								class="mb-4 border border-red-200 bg-red-50 px-3 py-3 text-sm text-red-600 dark:border-red-900 dark:bg-red-950 dark:text-red-400"
							>
								{inviteError}
							</div>
						{/if}

						<div class="mb-5">
							<label
								for="invite-email"
								class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300"
								>{i18n.t('workspaceMembers.modal.invite.emailLabel')}</label
							>
							<input
								id="invite-email"
								type="email"
								bind:value={inviteEmail}
								placeholder={i18n.t('workspaceMembers.modal.invite.emailPlaceholder')}
								required
								disabled={isInviting}
								class="w-full border border-gray-300 px-3 py-2.5 text-sm text-gray-900 transition-all focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-50 dark:border-gray-600 dark:bg-gray-900 dark:text-gray-50 dark:focus:ring-blue-500/20 dark:disabled:bg-slate-950"
							/>
						</div>

						<div class="mb-0">
							<label
								for="invite-role"
								class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300"
								>{i18n.t('workspaceMembers.modal.invite.roleLabel')}</label
							>
							<select
								id="invite-role"
								bind:value={inviteRole}
								disabled={isInviting}
								class="w-full border border-gray-300 px-3 py-2.5 text-sm text-gray-900 transition-all focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-50 dark:border-gray-600 dark:bg-gray-900 dark:text-gray-50 dark:focus:ring-blue-500/20 dark:disabled:bg-slate-950"
							>
								<option value="editor">{i18n.t('workspaceMembers.modal.invite.roleEditor')}</option>
								<option value="viewer">{i18n.t('workspaceMembers.modal.invite.roleViewer')}</option>
							</select>
							<p class="mt-1.5 text-[13px] text-gray-500 dark:text-gray-400">
								{#if inviteRole === 'editor'}
									{i18n.t('workspaceMembers.modal.invite.roleEditorDesc')}
								{:else}
									{i18n.t('workspaceMembers.modal.invite.roleViewerDesc')}
								{/if}
							</p>
						</div>
					</div>

					<div
						class="flex items-center justify-end gap-3 border-t border-gray-200 p-4 px-6 dark:border-gray-700"
					>
						<button
							type="button"
							class="border-2 border-[#372d2e] bg-white px-5 py-2.5 text-sm font-medium text-[#372d2e] transition-all duration-150 hover:-translate-y-0.5 hover:bg-gray-50 hover:shadow-[0_4px_0_#372d2e] active:translate-y-0 active:shadow-none disabled:cursor-not-allowed disabled:opacity-50 dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-gray-700"
							onclick={closeInviteModal}
							disabled={isInviting}
						>
							{i18n.t('workspaceMembers.modal.invite.cancel')}
						</button>
						<button
							type="submit"
							class="border-2 border-[#372d2e] bg-[#fcbd80] px-5 py-2.5 text-sm font-medium text-[#372d2e] transition-all duration-150 hover:-translate-y-0.5 hover:bg-[#dda877] hover:shadow-[0_4px_0_#372d2e] active:translate-y-0 active:shadow-none disabled:cursor-not-allowed disabled:opacity-50"
							disabled={isInviting}
						>
							{isInviting
								? i18n.t('workspaceMembers.modal.invite.sending')
								: i18n.t('workspaceMembers.modal.invite.send')}
						</button>
					</div>
				</form>
			{/if}
		</div>
	</div>
{/if}

<style>
	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-4px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.animate-slideDown {
		animation: slideDown 0.15s ease-out;
	}

	.animate-fadeIn {
		animation: fadeIn 0.2s ease-out;
	}

	.animate-slideUp {
		animation: slideUp 0.3s ease-out;
	}
</style>
