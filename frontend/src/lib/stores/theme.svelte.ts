type Theme = 'light' | 'dark';

class ThemeStore {
	private currentTheme = $state<Theme>('light');
	private allowedRoutes: string[] = [];

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
			// Apply theme based on current route
			this.updateThemeForRoute();
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
			this.updateThemeForRoute();
		}
	}

	toggleTheme() {
		this.setTheme(this.currentTheme === 'light' ? 'dark' : 'light');
	}

	// Check if current route should have dark mode enabled
	private isThemeAllowedRoute(): boolean {
		if (typeof window === 'undefined') return false;

		const path = window.location.pathname;

		// Exclude landing page, auth pages (login, register, forgot-password)
		const excludedRoutes = ['/', '/auth/login', '/auth/register', '/auth/forgot-password'];

		return !excludedRoutes.includes(path);
	}

	// Update theme application based on route
	private updateThemeForRoute() {
		if (typeof window !== 'undefined') {
			const root = document.documentElement;

			// Only apply dark theme if on allowed route and theme is dark
			if (this.isThemeAllowedRoute() && this.currentTheme === 'dark') {
				root.classList.add('dark');
			} else {
				root.classList.remove('dark');
			}
		}
	}

	// Call this method when route changes
	handleRouteChange() {
		this.updateThemeForRoute();
	}

	// Apply theme without checking route (for internal use in app pages)
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
