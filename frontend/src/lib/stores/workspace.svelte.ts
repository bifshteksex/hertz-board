import { api } from '$lib/services/api';
import type { Workspace, WorkspaceFilters, CreateWorkspaceRequest } from '$lib/types/api';

class WorkspaceStore {
	private _workspaces = $state<Workspace[]>([]);
	private _total = $state(0);
	private _isLoading = $state(false);
	private _filters = $state<WorkspaceFilters>({
		limit: 50,
		offset: 0,
		sort_by: 'updated_at',
		sort_order: 'desc'
	});

	get workspaces() {
		return this._workspaces;
	}

	get total() {
		return this._total;
	}

	get isLoading() {
		return this._isLoading;
	}

	get filters() {
		return this._filters;
	}

	async loadWorkspaces(filters?: WorkspaceFilters) {
		this._isLoading = true;
		try {
			const mergedFilters = { ...this._filters, ...filters };
			this._filters = mergedFilters;

			const response = await api.listWorkspaces(mergedFilters);
			this._workspaces = response.workspaces;
			this._total = response.total;
		} catch (error) {
			console.error('Failed to load workspaces:', error);
			throw error;
		} finally {
			this._isLoading = false;
		}
	}

	async createWorkspace(data: CreateWorkspaceRequest) {
		this._isLoading = true;
		try {
			const workspace = await api.createWorkspace(data);
			this._workspaces = [workspace, ...this._workspaces];
			this._total += 1;
			return workspace;
		} finally {
			this._isLoading = false;
		}
	}

	async updateWorkspace(id: string, data: { name?: string; description?: string }) {
		this._isLoading = true;
		try {
			const updated = await api.updateWorkspace(id, data);
			this._workspaces = this._workspaces.map((w) => (w.id === id ? updated : w));
			return updated;
		} finally {
			this._isLoading = false;
		}
	}

	async deleteWorkspace(id: string) {
		this._isLoading = true;
		try {
			await api.deleteWorkspace(id);
			this._workspaces = this._workspaces.filter((w) => w.id !== id);
			this._total -= 1;
		} finally {
			this._isLoading = false;
		}
	}

	async duplicateWorkspace(id: string, name: string) {
		this._isLoading = true;
		try {
			const duplicated = await api.duplicateWorkspace(id, name);
			this._workspaces = [duplicated, ...this._workspaces];
			this._total += 1;
			return duplicated;
		} finally {
			this._isLoading = false;
		}
	}

	setFilters(filters: WorkspaceFilters) {
		this._filters = { ...this._filters, ...filters };
	}

	reset() {
		this._workspaces = [];
		this._total = 0;
		this._filters = {
			limit: 50,
			offset: 0,
			sort_by: 'updated_at',
			sort_order: 'desc'
		};
	}
}

export const workspaceStore = new WorkspaceStore();
