type Theme = 'light' | 'dark';

class ThemeStore {
	private currentTheme = $state<Theme>('light');

	constructor() {
		// Load saved theme from localStorage on initialization
		if (typeof window !== 'undefined') {
			const savedTheme = localStorage.getItem('theme') as Theme | null;
			if (savedTheme && ['light', 'dark'].includes(savedTheme)) {
				this.currentTheme = savedTheme;
			} else {
				// Check system preference
				const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
				this.currentTheme = prefersDark ? 'dark' : 'light';
			}
			this.applyTheme();
		}
	}

	get theme(): Theme {
		return this.currentTheme;
	}

	get isDark(): boolean {
		return this.currentTheme === 'dark';
	}

	setTheme(theme: Theme) {
		this.currentTheme = theme;
		if (typeof window !== 'undefined') {
			localStorage.setItem('theme', theme);
			this.applyTheme();
		}
	}

	toggleTheme() {
		this.setTheme(this.currentTheme === 'light' ? 'dark' : 'light');
	}

	private applyTheme() {
		if (typeof window !== 'undefined') {
			const root = document.documentElement;
			if (this.currentTheme === 'dark') {
				root.classList.add('dark');
			} else {
				root.classList.remove('dark');
			}
		}
	}
}

export const themeStore = new ThemeStore();
