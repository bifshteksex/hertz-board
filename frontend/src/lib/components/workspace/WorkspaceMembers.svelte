<script lang="ts">
	import { api } from '$lib/services/api';
	import type { WorkspaceMember, WorkspaceRole } from '$lib/types/api';
	import { Users, Mail, MoreVertical, UserX, Eye, Edit3, Crown, Copy, Check } from 'lucide-svelte';

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
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load members';
			console.error('Failed to load members:', err);
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
		if (!confirm('Are you sure you want to remove this member?')) {
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

	function getRoleIcon(role: WorkspaceRole) {
		switch (role) {
			case 'owner':
				return Crown;
			case 'editor':
				return Edit3;
			case 'viewer':
				return Eye;
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
		return role.charAt(0).toUpperCase() + role.slice(1);
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

<div class="workspace-members">
	<div class="header">
		<div class="title">
			<Users size={20} />
			<h3>Workspace Members</h3>
			<span class="count">{members.length}</span>
		</div>

		{#if canInvite}
			<button class="invite-button" onclick={() => (showInviteModal = true)}>
				<Mail size={16} />
				Invite Member
			</button>
		{/if}
	</div>

	{#if isLoading}
		<div class="loading">Loading members...</div>
	{:else if error}
		<div class="error">{error}</div>
	{:else if members.length === 0}
		<div class="empty-state">
			<Users size={48} strokeWidth={1} />
			<p>No members yet</p>
		</div>
	{:else}
		<div class="members-list">
			{#each members as member (member.id)}
				{@const RoleIcon = getRoleIcon(member.role)}
				<div class="member-item">
					<!-- Avatar -->
					<div class="member-avatar">
						{#if member.user_avatar}
							<img src={member.user_avatar} alt={member.user_name} />
						{:else}
							{getInitials(member.user_name, member.user_email)}
						{/if}
					</div>

					<!-- Info -->
					<div class="member-info">
						<div class="member-name">{member.user_name || 'Unknown User'}</div>
						<div class="member-email">{member.user_email}</div>
					</div>

					<!-- Role badge -->
					{#if member.role === 'owner'}
						<div class="role-badge" style="color: {getRoleColor(member.role)}">
							<Crown size={14} />
							{getRoleLabel(member.role)}
						</div>
					{:else if member.role === 'editor'}
						<div class="role-badge" style="color: {getRoleColor(member.role)}">
							<Edit3 size={14} />
							{getRoleLabel(member.role)}
						</div>
					{:else}
						<div class="role-badge" style="color: {getRoleColor(member.role)}">
							<Eye size={14} />
							{getRoleLabel(member.role)}
						</div>
					{/if}

					<!-- Actions menu -->
					{#if canManageMembers && member.role !== 'owner'}
						<div class="member-actions">
							<button
								class="menu-button"
								onclick={() => (activeMenuId = activeMenuId === member.id ? null : member.id)}
							>
								<MoreVertical size={16} />
							</button>

							{#if activeMenuId === member.id}
								<div class="actions-menu">
									{#if currentUserRole === 'owner'}
										<button
											class="menu-item"
											onclick={() => handleUpdateRole(member.id, 'viewer')}
											disabled={member.role === 'viewer'}
										>
											<Eye size={14} />
											Make Viewer
										</button>
										<button
											class="menu-item"
											onclick={() => handleUpdateRole(member.id, 'editor')}
											disabled={member.role === 'editor'}
										>
											<Edit3 size={14} />
											Make Editor
										</button>
										<div class="menu-divider"></div>
									{/if}
									<button class="menu-item danger" onclick={() => handleRemoveMember(member.id)}>
										<UserX size={14} />
										Remove Member
									</button>
								</div>

								<!-- Backdrop -->
								<!-- svelte-ignore a11y_click_events_have_key_events -->
								<!-- svelte-ignore a11y_no_static_element_interactions -->
								<div class="menu-backdrop" onclick={() => (activeMenuId = null)}></div>
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
	<div class="modal-overlay" onclick={closeInviteModal} role="presentation">
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h3>Invite Member</h3>
				<button class="close-button" onclick={closeInviteModal}>&times;</button>
			</div>

			{#if inviteSuccess}
				<div class="modal-body">
					<div class="success-message">
						<div class="success-icon">
							<Mail size={32} />
						</div>
						<h4>Invitation Sent!</h4>
						<p>Share this link with the user:</p>

						<div class="invite-link-box">
							<input type="text" readonly value={inviteUrl} class="invite-link-input" />
							<button class="copy-button" onclick={copyInviteLink}>
								{#if copiedLink}
									<Check size={16} />
									Copied!
								{:else}
									<Copy size={16} />
									Copy
								{/if}
							</button>
						</div>

						<p class="expiry-note">Link expires in 7 days</p>
					</div>
				</div>
				<div class="modal-footer">
					<button class="button secondary" onclick={closeInviteModal}>Done</button>
				</div>
			{:else}
				<form onsubmit={handleInvite}>
					<div class="modal-body">
						{#if inviteError}
							<div class="error-message">{inviteError}</div>
						{/if}

						<div class="form-group">
							<label for="invite-email">Email Address</label>
							<input
								id="invite-email"
								type="email"
								bind:value={inviteEmail}
								placeholder="colleague@example.com"
								required
								disabled={isInviting}
							/>
						</div>

						<div class="form-group">
							<label for="invite-role">Role</label>
							<select id="invite-role" bind:value={inviteRole} disabled={isInviting}>
								<option value="editor">Editor - Can edit content</option>
								<option value="viewer">Viewer - Can only view</option>
							</select>
							<p class="help-text">
								{#if inviteRole === 'editor'}
									Editors can create, edit, and delete content.
								{:else}
									Viewers can only view content, not edit.
								{/if}
							</p>
						</div>
					</div>

					<div class="modal-footer">
						<button
							type="button"
							class="button secondary"
							onclick={closeInviteModal}
							disabled={isInviting}
						>
							Cancel
						</button>
						<button type="submit" class="button primary" disabled={isInviting}>
							{isInviting ? 'Sending...' : 'Send Invitation'}
						</button>
					</div>
				</form>
			{/if}
		</div>
	</div>
{/if}

<style>
	.workspace-members {
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		overflow: hidden;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 16px;
		border-bottom: 1px solid #e5e7eb;
	}

	.title {
		display: flex;
		align-items: center;
		gap: 8px;
		color: #111827;
	}

	.title h3 {
		font-size: 16px;
		font-weight: 600;
		margin: 0;
	}

	.count {
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 24px;
		height: 24px;
		padding: 0 8px;
		background: #f3f4f6;
		color: #6b7280;
		border-radius: 12px;
		font-size: 13px;
		font-weight: 600;
	}

	.invite-button {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 8px 16px;
		background: #3b82f6;
		color: white;
		border: none;
		border-radius: 6px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.invite-button:hover {
		background: #2563eb;
	}

	.loading,
	.error {
		padding: 32px;
		text-align: center;
		color: #6b7280;
	}

	.error {
		color: #dc2626;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 48px 16px;
		color: #9ca3af;
	}

	.empty-state p {
		margin-top: 12px;
		font-size: 14px;
	}

	.members-list {
		padding: 8px;
	}

	.member-item {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px;
		border-radius: 6px;
		transition: background 0.2s;
		position: relative;
	}

	.member-item:hover {
		background: #f9fafb;
	}

	.member-avatar {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		border-radius: 50%;
		font-size: 14px;
		font-weight: 600;
		flex-shrink: 0;
		overflow: hidden;
	}

	.member-avatar img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.member-info {
		flex: 1;
		min-width: 0;
	}

	.member-name {
		font-size: 14px;
		font-weight: 500;
		color: #111827;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.member-email {
		font-size: 13px;
		color: #6b7280;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.role-badge {
		display: flex;
		align-items: center;
		gap: 4px;
		padding: 4px 12px;
		background: #f9fafb;
		border-radius: 12px;
		font-size: 13px;
		font-weight: 500;
		flex-shrink: 0;
	}

	.member-actions {
		position: relative;
	}

	.menu-button {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		background: transparent;
		border: none;
		border-radius: 4px;
		color: #6b7280;
		cursor: pointer;
		transition: all 0.2s;
	}

	.menu-button:hover {
		background: #f3f4f6;
		color: #111827;
	}

	.actions-menu {
		position: absolute;
		top: calc(100% + 4px);
		right: 0;
		min-width: 160px;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
		padding: 4px;
		z-index: 100;
		animation: slideDown 0.15s ease-out;
	}

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

	.menu-item {
		display: flex;
		align-items: center;
		gap: 8px;
		width: 100%;
		padding: 8px 12px;
		background: transparent;
		border: none;
		border-radius: 4px;
		font-size: 14px;
		color: #374151;
		cursor: pointer;
		text-align: left;
		transition: background 0.2s;
	}

	.menu-item:hover:not(:disabled) {
		background: #f3f4f6;
	}

	.menu-item:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.menu-item.danger {
		color: #dc2626;
	}

	.menu-item.danger:hover {
		background: #fee2e2;
	}

	.menu-divider {
		height: 1px;
		background: #e5e7eb;
		margin: 4px 0;
	}

	.menu-backdrop {
		position: fixed;
		inset: 0;
		z-index: 99;
	}

	/* Modal styles */
	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 16px;
		animation: fadeIn 0.2s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.modal {
		background: white;
		border-radius: 12px;
		width: 100%;
		max-width: 500px;
		max-height: 90vh;
		overflow: auto;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
		animation: slideUp 0.3s ease-out;
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

	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 20px 24px;
		border-bottom: 1px solid #e5e7eb;
	}

	.modal-header h3 {
		font-size: 18px;
		font-weight: 600;
		color: #111827;
		margin: 0;
	}

	.close-button {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		border-radius: 6px;
		font-size: 24px;
		color: #6b7280;
		cursor: pointer;
		transition: all 0.2s;
	}

	.close-button:hover {
		background: #f3f4f6;
		color: #111827;
	}

	.modal-body {
		padding: 24px;
	}

	.form-group {
		margin-bottom: 20px;
	}

	.form-group:last-child {
		margin-bottom: 0;
	}

	.form-group label {
		display: block;
		font-size: 14px;
		font-weight: 500;
		color: #374151;
		margin-bottom: 6px;
	}

	.form-group input,
	.form-group select {
		width: 100%;
		padding: 10px 12px;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 14px;
		color: #111827;
		transition: all 0.2s;
	}

	.form-group input:focus,
	.form-group select:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.form-group input:disabled,
	.form-group select:disabled {
		background: #f9fafb;
		cursor: not-allowed;
	}

	.help-text {
		margin-top: 6px;
		font-size: 13px;
		color: #6b7280;
	}

	.error-message {
		padding: 12px;
		background: #fee2e2;
		border: 1px solid #fecaca;
		border-radius: 6px;
		color: #dc2626;
		font-size: 14px;
		margin-bottom: 16px;
	}

	.success-message {
		text-align: center;
	}

	.success-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 64px;
		height: 64px;
		margin: 0 auto 16px;
		background: #dbeafe;
		color: #3b82f6;
		border-radius: 50%;
	}

	.success-message h4 {
		font-size: 18px;
		font-weight: 600;
		color: #111827;
		margin: 0 0 8px;
	}

	.success-message p {
		font-size: 14px;
		color: #6b7280;
		margin: 0 0 16px;
	}

	.invite-link-box {
		display: flex;
		gap: 8px;
		margin-bottom: 12px;
	}

	.invite-link-input {
		flex: 1;
		padding: 10px 12px;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 13px;
		color: #374151;
		background: #f9fafb;
	}

	.copy-button {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 10px 16px;
		background: #3b82f6;
		color: white;
		border: none;
		border-radius: 6px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
		white-space: nowrap;
	}

	.copy-button:hover {
		background: #2563eb;
	}

	.expiry-note {
		font-size: 13px;
		color: #9ca3af;
	}

	.modal-footer {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: 12px;
		padding: 16px 24px;
		border-top: 1px solid #e5e7eb;
	}

	.button {
		padding: 10px 20px;
		border-radius: 6px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	.button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.button.secondary {
		background: white;
		color: #374151;
		border: 1px solid #d1d5db;
	}

	.button.secondary:hover:not(:disabled) {
		background: #f9fafb;
		border-color: #9ca3af;
	}

	.button.primary {
		background: #3b82f6;
		color: white;
	}

	.button.primary:hover:not(:disabled) {
		background: #2563eb;
	}
</style>
