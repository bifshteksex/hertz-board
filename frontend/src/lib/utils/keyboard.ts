/**
 * Keyboard Shortcuts System
 *
 * Управление глобальными горячими клавишами для Canvas
 */

export type ShortcutKey = string;
export type ShortcutHandler = (event: KeyboardEvent) => void | boolean;

export interface Shortcut {
	key: string; // e.g., 'v', 'ctrl+z', 'ctrl+shift+g'
	description: string;
	category: ShortcutCategory;
	handler: ShortcutHandler;
}

export type ShortcutCategory =
	| 'tools'
	| 'edit'
	| 'selection'
	| 'layers'
	| 'grouping'
	| 'view'
	| 'other';

/**
 * Parse keyboard event to shortcut key string
 */
export function eventToShortcutKey(event: KeyboardEvent): string {
	const parts: string[] = [];

	if (event.ctrlKey || event.metaKey) parts.push('ctrl');
	if (event.shiftKey) parts.push('shift');
	if (event.altKey) parts.push('alt');

	// Normalize key
	let key = event.key.toLowerCase();

	// Special keys mapping
	const specialKeys: Record<string, string> = {
		escape: 'esc',
		delete: 'del',
		' ': 'space',
		arrowup: 'up',
		arrowdown: 'down',
		arrowleft: 'left',
		arrowright: 'right'
	};

	key = specialKeys[key] || key;

	// Ignore modifier keys alone
	if (['control', 'shift', 'alt', 'meta'].includes(key)) {
		return '';
	}

	parts.push(key);

	return parts.join('+');
}

/**
 * Format shortcut key for display
 */
export function formatShortcut(shortcut: string): string {
	return shortcut
		.split('+')
		.map((part) => {
			// Capitalize first letter
			const formatted = part.charAt(0).toUpperCase() + part.slice(1);

			// Platform-specific formatting
			if (part === 'ctrl') {
				// Use Cmd symbol on Mac
				return navigator.platform.includes('Mac') ? '⌘' : 'Ctrl';
			}
			if (part === 'shift') return '⇧';
			if (part === 'alt') return navigator.platform.includes('Mac') ? '⌥' : 'Alt';

			return formatted;
		})
		.join('+');
}

/**
 * Keyboard Shortcuts Manager
 */
class KeyboardManager {
	private shortcuts = new Map<ShortcutKey, Shortcut>();
	private enabled = true;

	/**
	 * Register a keyboard shortcut
	 */
	register(key: string, description: string, category: ShortcutCategory, handler: ShortcutHandler) {
		this.shortcuts.set(key, { key, description, category, handler });
	}

	/**
	 * Unregister a keyboard shortcut
	 */
	unregister(key: string) {
		this.shortcuts.delete(key);
	}

	/**
	 * Handle keyboard event
	 * Returns true if event was handled
	 */
	handleKeyDown(event: KeyboardEvent): boolean {
		if (!this.enabled) return false;

		// Don't handle shortcuts when typing in inputs
		const target = event.target as HTMLElement;
		if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
			// Allow Escape to blur input
			if (event.key === 'Escape') {
				target.blur();
				return true;
			}
			return false;
		}

		const shortcutKey = eventToShortcutKey(event);
		if (!shortcutKey) return false;

		const shortcut = this.shortcuts.get(shortcutKey);
		if (!shortcut) return false;

		// Call handler
		const result = shortcut.handler(event);

		// Prevent default if handler didn't return false
		if (result !== false) {
			event.preventDefault();
			event.stopPropagation();
			return true;
		}

		return false;
	}

	/**
	 * Get all registered shortcuts
	 */
	getShortcuts(): Shortcut[] {
		return Array.from(this.shortcuts.values());
	}

	/**
	 * Get shortcuts by category
	 */
	getShortcutsByCategory(category: ShortcutCategory): Shortcut[] {
		return this.getShortcuts().filter((s) => s.category === category);
	}

	/**
	 * Get all categories
	 */
	getCategories(): ShortcutCategory[] {
		const categories = new Set<ShortcutCategory>();
		this.shortcuts.forEach((s) => categories.add(s.category));
		return Array.from(categories);
	}

	/**
	 * Enable/disable shortcuts
	 */
	setEnabled(enabled: boolean) {
		this.enabled = enabled;
	}

	/**
	 * Clear all shortcuts
	 */
	clear() {
		this.shortcuts.clear();
	}
}

// Singleton instance
export const keyboardManager = new KeyboardManager();

/**
 * Category display names
 */
export const categoryNames: Record<ShortcutCategory, string> = {
	tools: 'Tools',
	edit: 'Edit',
	selection: 'Selection',
	layers: 'Layers',
	grouping: 'Grouping',
	view: 'View',
	other: 'Other'
};
