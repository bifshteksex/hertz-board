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
		backToLogin: 'Back to Login'
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
			oauthWarning: 'You signed in with {provider}. Password changes are not available for OAuth accounts.',
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
			languageZh: '中文'
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
