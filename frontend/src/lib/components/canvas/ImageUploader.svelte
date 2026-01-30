<script lang="ts">
	/**
	 * ImageUploader - компонент для загрузки изображений на canvas
	 */
	import { api } from '$lib/services/api';

	interface Props {
		workspaceId: string;
		onImageUploaded?: (_imageUrl: string, _width: number, _height: number) => void;
	}

	let { workspaceId, onImageUploaded }: Props = $props();

	let fileInput: HTMLInputElement;
	let isUploading = $state(false);
	let uploadError = $state<string | null>(null);

	async function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		const files = target.files;
		if (!files || files.length === 0) return;

		await uploadFile(files[0]);
		target.value = ''; // Reset input
	}

	async function uploadFile(file: File) {
		// Validate file type
		if (!file.type.startsWith('image/')) {
			uploadError = 'Please select an image file';
			return;
		}

		// Validate file size (max 10MB)
		const maxSize = 10 * 1024 * 1024;
		if (file.size > maxSize) {
			uploadError = 'Image size must be less than 10MB';
			return;
		}

		isUploading = true;
		uploadError = null;

		try {
			const response = await api.uploadAsset(workspaceId, file);

			// Get image dimensions
			const img = new Image();
			img.src = response.url;
			await new Promise((resolve) => {
				img.onload = resolve;
			});

			if (onImageUploaded) {
				onImageUploaded(response.url, img.width, img.height);
			}
		} catch (error) {
			console.error('Failed to upload image:', error);
			uploadError = error instanceof Error ? error.message : 'Failed to upload image';
		} finally {
			isUploading = false;
		}
	}

	async function handlePaste(e: ClipboardEvent) {
		const items = e.clipboardData?.items;
		if (!items) return;

		for (let i = 0; i < items.length; i++) {
			const item = items[i];
			if (item.type.startsWith('image/')) {
				e.preventDefault();
				const file = item.getAsFile();
				if (file) {
					await uploadFile(file);
				}
				break;
			}
		}
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		const files = e.dataTransfer?.files;
		if (!files || files.length === 0) return;

		const file = files[0];
		if (file.type.startsWith('image/')) {
			await uploadFile(file);
		}
	}

	function triggerFileInput() {
		fileInput.click();
	}

	// Export functions for external use
	export function openFileDialog() {
		triggerFileInput();
	}
</script>

<svelte:window onpaste={handlePaste} ondrop={handleDrop} ondragover={(e) => e.preventDefault()} />

<input
	bind:this={fileInput}
	type="file"
	accept="image/*"
	onchange={handleFileSelect}
	style="display: none;"
/>

{#if isUploading}
	<div class="upload-overlay">
		<div class="upload-spinner"></div>
		<p>Uploading image...</p>
	</div>
{/if}

{#if uploadError}
	<div class="upload-error">
		<p>{uploadError}</p>
		<button onclick={() => (uploadError = null)}>Dismiss</button>
	</div>
{/if}

<style>
	.upload-overlay {
		position: fixed;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		z-index: 1000;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}

	.upload-spinner {
		width: 40px;
		height: 40px;
		border: 4px solid #f3f4f6;
		border-top-color: #3b82f6;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.upload-error {
		position: fixed;
		top: 1rem;
		right: 1rem;
		background: #fee2e2;
		border: 1px solid #ef4444;
		color: #991b1b;
		padding: 1rem;
		border-radius: 6px;
		z-index: 1000;
		max-width: 300px;
	}

	.upload-error button {
		margin-top: 0.5rem;
		padding: 0.25rem 0.5rem;
		background: #ef4444;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	.upload-error button:hover {
		background: #dc2626;
	}
</style>
