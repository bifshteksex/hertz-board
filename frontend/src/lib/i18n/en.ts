export default {
	common: {
		loading: 'Loading...',
		save: 'Save',
		cancel: 'Cancel',
		delete: 'Delete',
		edit: 'Edit',
		close: 'Close',
		search: 'Search',
		create: 'Create',
		update: 'Update',
		confirm: 'Confirm',
		back: 'Back'
	},
	nav: {
		dashboard: 'Dashboard',
		workspaces: 'Workspaces',
		settings: 'Settings',
		logout: 'Logout'
	},
	auth: {
		appName: 'HertzBoard',
		login: 'Log In',
		register: 'Register',
		email: 'Email',
		password: 'Password',
		name: 'Full Name',
		confirmPassword: 'Confirm Password',
		forgotPassword: 'Forgot Password?',
		noAccount: "Don't have an account?",
		hasAccount: 'Already have an account?',
		signIn: 'Sign In',
		signUp: 'Sign Up',
		orContinueWith: 'Or continue with',
		resetPassword: 'Reset Password',
		backToLogin: 'Back to Login',
		loginTitle: 'Sign in to your account',
		registerTitle: 'Create your account',
		fullName: 'Full name',
		emailPlaceholder: 'Email address',
		passwordPlaceholder: 'Password',
		passwordMinPlaceholder: 'Password (min 8 characters)',
		confirmPasswordPlaceholder: 'Confirm password',
		signingIn: 'Signing in...',
		creatingAccount: 'Creating account...',
		googleButton: 'Google',
		githubButton: 'GitHub',
		resetTitle: 'Reset your password',
		resetDescription: "Enter your email address and we'll send you a link to reset your password.",
		resetSuccess:
			"Check your email for a link to reset your password. If it doesn't appear within a few minutes, check your spam folder.",
		sendResetLink: 'Send reset link',
		sending: 'Sending...',
		errorPasswordMatch: 'Passwords do not match',
		errorPasswordLength: 'Password must be at least 8 characters'
	},
	dashboard: {
		title: 'Workspaces',
		subtitle: 'Manage your collaborative boards',
		newWorkspace: 'New Workspace',
		searchPlaceholder: 'Search workspaces...',
		loading: 'Loading workspaces...',
		noWorkspaces: 'No workspaces found',
		createFirst: 'Create your first workspace',
		member: 'member',
		members: 'members',
		role: 'Role',
		menu: {
			open: 'Open',
			rename: 'Rename',
			copyLink: 'Copy link',
			share: 'Share',
			duplicate: 'Duplicate',
			delete: 'Delete'
		},
		modal: {
			create: {
				title: 'Create New Workspace',
				name: 'Name',
				description: 'Description (optional)',
				namePlaceholder: 'My Workspace',
				descriptionPlaceholder: 'Describe your workspace...',
				creating: 'Creating...',
				create: 'Create'
			},
			duplicate: {
				title: 'Duplicate Workspace',
				newName: 'New workspace name',
				copyOf: 'This will create a copy of "{name}"',
				duplicating: 'Duplicating...',
				duplicate: 'Duplicate'
			},
			rename: {
				title: 'Rename Workspace',
				saving: 'Saving...',
				saveChanges: 'Save Changes'
			}
		},
		alerts: {
			deleteConfirm: 'Are you sure you want to delete this workspace?',
			shareComingSoon: 'Share functionality will be implemented in a future phase',
			linkCopied: 'Link copied to clipboard!'
		},
		time: {
			today: 'Today',
			yesterday: 'Yesterday',
			daysAgo: '{count} days ago',
			weeksAgo: '{count} weeks ago',
			monthsAgo: '{count} months ago'
		}
	},
	settings: {
		title: 'Settings',
		subtitle: 'Manage your account settings and preferences',
		tabs: {
			profile: 'Profile',
			password: 'Password',
			account: 'Account',
			preferences: 'Preferences'
		},
		profile: {
			title: 'Profile Information',
			fullName: 'Full Name',
			avatarUrl: 'Avatar URL',
			avatarPlaceholder: 'https://example.com/avatar.jpg',
			avatarHint: 'Optional: Enter a URL to your profile picture',
			provider: 'Provider',
			saveChanges: 'Save Changes',
			saving: 'Saving...',
			successMessage: 'Profile updated successfully!'
		},
		password: {
			title: 'Change Password',
			current: 'Current Password',
			new: 'New Password',
			confirm: 'Confirm New Password',
			hint: 'Minimum 8 characters',
			change: 'Change Password',
			changing: 'Changing...',
			successMessage: 'Password changed successfully!',
			oauthWarning:
				'You signed in with {provider}. Password changes are not available for OAuth accounts.',
			errorMatch: 'Passwords do not match',
			errorLength: 'Password must be at least 8 characters'
		},
		account: {
			title: 'Account Information',
			email: 'Email',
			accountType: 'Account Type',
			emailVerified: 'Email Verified',
			memberSince: 'Member Since',
			verified: 'Yes',
			notVerified: 'No',
			dangerZone: 'Danger Zone',
			deleteWarning: 'Once you delete your account, there is no going back. Please be certain.',
			deleteAccount: 'Delete Account'
		},
		preferences: {
			title: 'Preferences',
			language: 'Language',
			theme: 'Theme',
			themeLight: 'Light',
			themeDark: 'Dark',
			languageEn: 'English',
			languageRu: 'Русский',
			languageZh: '中文',
			languageLabel: 'Language',
			themeLabel: 'Theme',
			languageHint: 'Select your preferred language for the interface',
			themeHint: 'Choose your preferred color scheme'
		}
	},
	landing: {
		subtitle: 'Real-time collaborative workspace platform',
		login: 'Login',
		signUp: 'Sign Up',
		features: {
			realtime: {
				title: 'Real-time Collaboration',
				description: 'Work together with your team in real-time with live cursors and presence'
			},
			canvas: {
				title: 'Powerful Canvas',
				description: 'Create with text, shapes, images, and more on an infinite canvas'
			},
			tech: {
				title: 'Built with Go',
				description: 'High-performance backend powered by CloudWeGo Hertz and WebSockets'
			}
		},
		footer: 'Made with ❤️ by'
	},
	workspaceDetail: {
		loading: 'Loading workspace...',
		backToDashboard: 'Back to Dashboard',
		share: 'Share',
		canvasComingSoon: 'Canvas Coming Soon',
		canvasDescription: 'The canvas editor will be implemented in Phase 6',
		workspaceId: 'Workspace ID:',
		role: 'Role:',
		errors: {
			missingId: 'Workspace ID is missing',
			loadFailed: 'Failed to load workspace'
		}
	},
	workspace: {
		title: 'Workspaces',
		createNew: 'Create Workspace',
		myWorkspaces: 'My Workspaces',
		recentWorkspaces: 'Recent Workspaces',
		noWorkspaces: 'No workspaces yet',
		createFirst: 'Create your first workspace to get started'
	},
	errors: {
		generic: 'An error occurred',
		network: 'Network error. Please check your connection.',
		unauthorized: 'Unauthorized. Please log in.',
		notFound: 'Not found',
		serverError: 'Server error. Please try again later.'
	}
};
