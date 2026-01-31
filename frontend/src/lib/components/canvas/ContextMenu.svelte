<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Scissors,
		Copy,
		Clipboard,
		Files,
		Trash2,
		Lock,
		Unlock,
		Group,
		Ungroup,
		ChevronsUp,
		ChevronsDown,
		Link
	} from 'lucide-svelte';

	interface Props {
		x: number;
		y: number;
		selectedCount: number;
		hasGroupedElements: boolean;
		isLocked: boolean;
		onCut: () => void;
		onCopy: () => void;
		onPaste: () => void;
		onDuplicate: () => void;
		onDelete: () => void;
		onLock: () => void;
		onUnlock: () => void;
		onGroup: () => void;
		onUngroup: () => void;
		onBringToFront: () => void;
		onSendToBack: () => void;
		onCopyLink?: () => void;
		onClose: () => void;
	}

	let {
		x,
		y,
		selectedCount,
		hasGroupedElements,
		isLocked,
		onCut,
		onCopy,
		onPaste,
		onDuplicate,
		onDelete,
		onLock,
		onUnlock,
		onGroup,
		onUngroup,
		onBringToFront,
		onSendToBack,
		onCopyLink,
		onClose
	}: Props = $props();

	let menuElement: HTMLDivElement;

	onMount(() => {
		// Close on outside click
		const handleClickOutside = (e: MouseEvent) => {
			if (menuElement && !menuElement.contains(e.target as Node)) {
				onClose();
			}
		};

		// Close on escape
		const handleEscape = (e: KeyboardEvent) => {
			if (e.code === 'Escape') {
				onClose();
			}
		};

		setTimeout(() => {
			document.addEventListener('click', handleClickOutside);
			document.addEventListener('keydown', handleEscape);
		}, 0);

		return () => {
			document.removeEventListener('click', handleClickOutside);
			document.removeEventListener('keydown', handleEscape);
		};
	});

	// Adjust position if menu goes off screen
	$effect(() => {
		if (menuElement) {
			const rect = menuElement.getBoundingClientRect();
			const viewportWidth = window.innerWidth;
			const viewportHeight = window.innerHeight;

			if (rect.right > viewportWidth) {
				menuElement.style.left = `${x - rect.width}px`;
			}

			if (rect.bottom > viewportHeight) {
				menuElement.style.top = `${y - rect.height}px`;
			}
		}
	});

	interface MenuItem {
		id: string;
		label: string;
		icon: any;
		shortcut?: string;
		action: () => void;
		disabled?: boolean;
		separator?: boolean;
	}

	const menuItems = $derived.by(() => {
		const items: MenuItem[] = [
			{
				id: 'cut',
				label: 'Cut',
				icon: Scissors,
				shortcut: 'Ctrl+X',
				action: () => {
					onCut();
					onClose();
				},
				disabled: isLocked
			},
			{
				id: 'copy',
				label: 'Copy',
				icon: Copy,
				shortcut: 'Ctrl+C',
				action: () => {
					onCopy();
					onClose();
				}
			},
			{
				id: 'paste',
				label: 'Paste',
				icon: Clipboard,
				shortcut: 'Ctrl+V',
				action: () => {
					onPaste();
					onClose();
				}
			},
			{
				id: 'duplicate',
				label: 'Duplicate',
				icon: Files,
				shortcut: 'Ctrl+D',
				action: () => {
					onDuplicate();
					onClose();
				},
				disabled: isLocked
			},
			{
				id: 'sep1',
				label: '',
				icon: null,
				action: () => {},
				separator: true
			},
			{
				id: 'delete',
				label: 'Delete',
				icon: Trash2,
				shortcut: 'Del',
				action: () => {
					onDelete();
					onClose();
				},
				disabled: isLocked
			}
		];

		// Add separator before layer controls
		items.push({
			id: 'sep2',
			label: '',
			icon: null,
			action: () => {},
			separator: true
		});

		// Layer controls
		items.push(
			{
				id: 'bring-to-front',
				label: 'Bring to Front',
				icon: ChevronsUp,
				shortcut: 'Ctrl+]',
				action: () => {
					onBringToFront();
					onClose();
				},
				disabled: isLocked
			},
			{
				id: 'send-to-back',
				label: 'Send to Back',
				icon: ChevronsDown,
				shortcut: 'Ctrl+[',
				action: () => {
					onSendToBack();
					onClose();
				},
				disabled: isLocked
			}
		);

		// Group/Ungroup (only for multiple selection)
		if (selectedCount > 1 || hasGroupedElements) {
			items.push({
				id: 'sep3',
				label: '',
				icon: null,
				action: () => {},
				separator: true
			});

			if (selectedCount > 1) {
				items.push({
					id: 'group',
					label: 'Group',
					icon: Group,
					shortcut: 'Ctrl+G',
					action: () => {
						onGroup();
						onClose();
					},
					disabled: isLocked
				});
			}

			if (hasGroupedElements) {
				items.push({
					id: 'ungroup',
					label: 'Ungroup',
					icon: Ungroup,
					shortcut: 'Ctrl+Shift+G',
					action: () => {
						onUngroup();
						onClose();
					},
					disabled: isLocked
				});
			}
		}

		// Lock/Unlock
		items.push(
			{
				id: 'sep4',
				label: '',
				icon: null,
				action: () => {},
				separator: true
			},
			{
				id: 'lock',
				label: isLocked ? 'Unlock' : 'Lock',
				icon: isLocked ? Unlock : Lock,
				action: () => {
					if (isLocked) {
						onUnlock();
					} else {
						onLock();
					}
					onClose();
				}
			}
		);

		// Copy Link (if available)
		if (onCopyLink) {
			items.push(
				{
					id: 'sep5',
					label: '',
					icon: null,
					action: () => {},
					separator: true
				},
				{
					id: 'copy-link',
					label: 'Copy Link',
					icon: Link,
					action: () => {
						onCopyLink();
						onClose();
					}
				}
			);
		}

		return items;
	});
</script>

<div
	bind:this={menuElement}
	class="fixed z-[10000] min-w-[200px] animate-[menuFadeIn_0.1s_ease-out] rounded-lg border border-gray-200 bg-white p-1 shadow-lg"
	style="left: {x}px; top: {y}px;"
>
	{#each menuItems as item}
		{#if item.separator}
			<div class="my-1 h-px bg-gray-200"></div>
		{:else}
			{@const Icon = item.icon}
			<button
				class="flex w-full cursor-pointer items-center gap-3 rounded-md border-none bg-transparent px-3 py-2 text-left text-sm font-medium text-gray-700 transition-all duration-150 hover:bg-gray-100 hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-40"
				onclick={item.action}
				disabled={item.disabled}
			>
				<Icon size={16} class="flex-shrink-0" />
				<span class="flex-1">{item.label}</span>
				{#if item.shortcut}
					<span class="font-mono text-xs text-gray-400">{item.shortcut}</span>
				{/if}
			</button>
		{/if}
	{/each}
</div>

<style>
	@keyframes menuFadeIn {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}
</style>
