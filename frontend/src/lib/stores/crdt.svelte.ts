import type { OperationPayload, OperationType } from '$lib/types/websocket';
import type { CanvasElement } from '$lib/types/api';

/**
 * CRDT Store for handling distributed operations
 * Implements Last-Write-Wins (LWW) strategy with Lamport clock
 */
class CRDTStore {
	// State vector: tracks last seen timestamp per user
	private _stateVector = $state<Record<string, number>>({});

	// Lamport clock for this client
	private _clock = $state<number>(0);

	// Current user ID
	private _userId = $state<string | null>(null);

	// Pending operations queue (for when offline)
	private _pendingOperations = $state<OperationPayload[]>([]);

	/**
	 * Initialize CRDT with user ID
	 */
	initialize(userId: string): void {
		this._userId = userId;
		this._clock = 0;
		this._stateVector = {};
		this._pendingOperations = [];
	}

	/**
	 * Get current state vector
	 */
	get stateVector(): Record<string, number> {
		return { ...this._stateVector };
	}

	/**
	 * Get current Lamport clock value
	 */
	get clock(): number {
		return this._clock;
	}

	/**
	 * Get pending operations
	 */
	get pendingOperations(): OperationPayload[] {
		return [...this._pendingOperations];
	}

	/**
	 * Tick the Lamport clock (increment)
	 */
	tick(): number {
		this._clock++;
		return this._clock;
	}

	/**
	 * Update Lamport clock based on received timestamp
	 */
	updateClock(receivedTimestamp: number): number {
		if (receivedTimestamp > this._clock) {
			this._clock = receivedTimestamp;
		}
		this._clock++;
		return this._clock;
	}

	/**
	 * Update state vector for a user
	 */
	updateStateVector(userId: string, timestamp: number): void {
		if (!this._stateVector[userId] || timestamp > this._stateVector[userId]) {
			this._stateVector[userId] = timestamp;
		}
	}

	/**
	 * Set state vector (from sync response)
	 */
	setStateVector(stateVector: Record<string, number>): void {
		this._stateVector = { ...stateVector };

		// Update clock to max timestamp in state vector
		const maxTimestamp = Math.max(...Object.values(stateVector), 0);
		if (maxTimestamp > this._clock) {
			this._clock = maxTimestamp;
		}
	}

	/**
	 * Create a new operation with Lamport timestamp
	 */
	createOperation(
		elementId: string,
		workspaceId: string,
		opType: OperationType,
		data?: Record<string, unknown>
	): OperationPayload {
		if (!this._userId) {
			throw new Error('CRDT not initialized with user ID');
		}

		const timestamp = this.tick();

		const operation: OperationPayload = {
			element_id: elementId,
			workspace_id: workspaceId,
			user_id: this._userId,
			op_type: opType,
			data,
			timestamp
		};

		// Update our state vector
		this.updateStateVector(this._userId, timestamp);

		return operation;
	}

	/**
	 * Check if an operation should be applied (based on LWW)
	 */
	shouldApplyOperation(operation: OperationPayload, currentVersion?: number): boolean {
		// Update our clock
		this.updateClock(operation.timestamp);

		// If no current version, always apply
		if (currentVersion === undefined || currentVersion === null) {
			return true;
		}

		// Apply if operation timestamp is greater (Last-Write-Wins)
		if (operation.timestamp > currentVersion) {
			return true;
		}

		// If timestamps are equal, use user ID for deterministic ordering
		if (operation.timestamp === currentVersion) {
			// This is a rare case, but we need deterministic resolution
			// Use lexicographic ordering of user IDs
			return operation.user_id > (this._userId || '');
		}

		return false;
	}

	/**
	 * Apply an operation to an element
	 * Returns updated element or null if operation should be ignored
	 */
	applyOperation(
		operation: OperationPayload,
		currentElement?: CanvasElement
	): CanvasElement | null {
		// Update clock and state vector
		this.updateClock(operation.timestamp);
		this.updateStateVector(operation.user_id, operation.timestamp);

		switch (operation.op_type) {
			case 'create':
				return this.applyCreate(operation);

			case 'update':
				return this.applyUpdate(operation, currentElement);

			case 'delete':
				return this.applyDelete(operation, currentElement);

			case 'move':
				return this.applyMove(operation, currentElement);

			default:
				console.warn('[CRDT] Unknown operation type:', operation.op_type);
				return null;
		}
	}

	/**
	 * Add operation to pending queue
	 */
	addPendingOperation(operation: OperationPayload): void {
		this._pendingOperations.push(operation);
	}

	/**
	 * Clear pending operations
	 */
	clearPendingOperations(): void {
		this._pendingOperations = [];
	}

	/**
	 * Reset CRDT state
	 */
	reset(): void {
		this._clock = 0;
		this._stateVector = {};
		this._pendingOperations = [];
		this._userId = null;
	}

	// Private methods for applying operations

	private applyCreate(operation: OperationPayload): CanvasElement {
		const element: CanvasElement = {
			id: operation.element_id,
			workspace_id: operation.workspace_id,
			...(operation.data as Partial<CanvasElement>),
			version: operation.timestamp
		} as CanvasElement;

		return element;
	}

	private applyUpdate(
		operation: OperationPayload,
		currentElement?: CanvasElement
	): CanvasElement | null {
		if (!currentElement) {
			console.warn('[CRDT] Cannot update non-existent element:', operation.element_id);
			return null;
		}

		// Check if we should apply (LWW)
		if (!this.shouldApplyOperation(operation, currentElement.version)) {
			console.log('[CRDT] Ignoring older operation');
			return null;
		}

		// Merge operation data with current element
		const updatedElement: CanvasElement = {
			...currentElement,
			...(operation.data as Partial<CanvasElement>),
			version: operation.timestamp
		};

		return updatedElement;
	}

	private applyDelete(
		operation: OperationPayload,
		currentElement?: CanvasElement
	): CanvasElement | null {
		if (!currentElement) {
			// Already deleted
			return null;
		}

		// Check if we should apply (LWW)
		if (!this.shouldApplyOperation(operation, currentElement.version)) {
			console.log('[CRDT] Ignoring older delete operation');
			return currentElement; // Keep current element
		}

		// Return null to indicate deletion
		return null;
	}

	private applyMove(
		operation: OperationPayload,
		currentElement?: CanvasElement
	): CanvasElement | null {
		if (!currentElement) {
			console.warn('[CRDT] Cannot move non-existent element:', operation.element_id);
			return null;
		}

		// Check if we should apply (LWW)
		if (!this.shouldApplyOperation(operation, currentElement.version)) {
			console.log('[CRDT] Ignoring older move operation');
			return null;
		}

		// Update position
		const updatedElement: CanvasElement = {
			...currentElement,
			pos_x: (operation.data?.pos_x as number) ?? currentElement.pos_x,
			pos_y: (operation.data?.pos_y as number) ?? currentElement.pos_y,
			version: operation.timestamp
		};

		return updatedElement;
	}
}

export const crdtStore = new CRDTStore();
