export default {
	seo: {
		defaultTitle: 'HertzBoard - Real-time Collaborative Workspace',
		defaultDescription:
			'Create, collaborate, and visualize ideas in real-time with HertzBoard - a powerful collaborative whiteboard platform built with modern technologies.',
		keywords:
			'collaborative whiteboard, real-time collaboration, visual workspace, online whiteboard, team collaboration, canvas editor',
		ogSiteName: 'HertzBoard',
		twitterCard: 'summary_large_image',
		pages: {
			landing: {
				title: 'HertzBoard - Real-time Collaborative Workspace',
				description:
					'Create, collaborate, and visualize ideas in real-time with your team on an infinite canvas.'
			},
			login: {
				title: 'Sign In - HertzBoard',
				description: 'Sign in to your HertzBoard account to access your collaborative workspaces.'
			},
			register: {
				title: 'Create Account - HertzBoard',
				description: 'Create a new HertzBoard account and start collaborating with your team today.'
			},
			dashboard: {
				title: 'Dashboard - HertzBoard',
				description: 'Manage your collaborative workspaces and boards.'
			},
			settings: {
				title: 'Settings - HertzBoard',
				description: 'Manage your account settings and preferences.'
			},
			workspace: {
				title: '{name} - HertzBoard',
				description: 'Collaborate on {name} workspace in real-time.'
			}
		}
	},
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
		footer: 'Made with {heart} by'
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
	},
	canvas: {
		toolbar: {
			undo: 'Undo',
			redo: 'Redo',
			undoShortcut: 'Undo (Ctrl+Z)',
			redoShortcut: 'Redo (Ctrl+Y)',
			help: 'Keyboard Shortcuts (?)',
			helpAria: 'Show keyboard shortcuts',
			tools: {
				select: 'Select (V)',
				text: 'Text (T)',
				pen: 'Pen (P)',
				eraser: 'Eraser (E)',
				sticky: 'Sticky Note (S)',
				image: 'Image (I)',
				connector: 'Connector',
				shapes: 'Shapes',
				lists: 'Lists'
			},
			shapes: {
				rectangle: 'Rectangle',
				ellipse: 'Circle',
				triangle: 'Triangle',
				line: 'Line',
				arrow: 'Arrow'
			},
			lists: {
				bullet: 'Bullet List',
				numbered: 'Numbered List',
				checkbox: 'Checkbox List'
			},
			zorder: {
				bringToFront: 'Bring to Front (Ctrl+Shift+])',
				bringForward: 'Bring Forward (Ctrl+])',
				sendBackward: 'Send Backward (Ctrl+[)',
				sendToBack: 'Send to Back (Ctrl+Shift+[)'
			},
			grouping: {
				group: 'Group (Ctrl+G)',
				ungroup: 'Ungroup (Ctrl+Shift+G)',
				groupElements: 'Group elements',
				ungroupElements: 'Ungroup elements'
			},
			view: {
				toggleGrid: "Toggle Grid (Ctrl+')",
				toggleSnap: "Snap to Grid (Ctrl+Shift+')",
				gridAria: 'Toggle grid',
				snapAria: 'Toggle snap to grid'
			}
		},
		shortcuts: {
			title: 'Keyboard Shortcuts',
			subtitle: 'Learn all the shortcuts to work faster',
			search: 'Search shortcuts...',
			noResults: 'No shortcuts found for "{query}"',
			footer: 'Press',
			footerKey: '?',
			footerText: 'anytime to show this panel',
			categories: {
				tools: 'Tools',
				edit: 'Edit',
				selection: 'Selection',
				layers: 'Layers',
				grouping: 'Grouping',
				view: 'View',
				other: 'Other'
			}
		},
		workspace: {
			loading: 'Loading workspace...',
			backToDashboard: 'Back to Dashboard',
			toolbar: {
				layers: 'Layers',
				share: 'Share'
			}
		},
		contextMenu: {
			cut: 'Cut',
			copy: 'Copy',
			paste: 'Paste',
			duplicate: 'Duplicate',
			delete: 'Delete',
			bringToFront: 'Bring to Front',
			sendToBack: 'Send to Back',
			group: 'Group',
			ungroup: 'Ungroup',
			lock: 'Lock',
			unlock: 'Unlock',
			copyLink: 'Copy Link'
		},
		imageUploader: {
			uploading: 'Uploading image...',
			dismiss: 'Dismiss',
			errorFileType: 'Please select an image file',
			errorFileSize: 'Image size must be less than 10MB',
			errorUploadFailed: 'Failed to upload image'
		},
		saveStatus: {
			saving: 'Saving...',
			saved: 'Saved',
			savedAgo: 'Saved {time}',
			saveFailed: 'Save failed',
			unsavedChanges: '{count} unsaved changes',
			allSaved: 'All changes saved',
			details: 'Details'
		},
		activeUsers: {
			title: 'Active Users',
			noUsers: 'No active users',
			you: '(You)',
			active: 'Active',
			idle: 'Idle',
			justNow: 'just now',
			secondsAgo: '{count}s ago',
			minutesAgo: '{count}m ago',
			hoursAgo: '{count}h ago',
			closeButton: 'Close'
		},
		connectionStatus: {
			connected: 'Connected',
			disconnected: 'Disconnected',
			syncing: 'Syncing...',
			error: 'Connection Error',
			retry: 'Retry',
			reconnect: 'Reconnect'
		}
	},
	workspaceMembers: {
		title: 'Workspace Members',
		inviteMember: 'Invite Member',
		loading: 'Loading members...',
		noMembers: 'No members yet',
		owner: 'Owner',
		editor: 'Editor',
		viewer: 'Viewer',
		menu: {
			makeViewer: 'Make Viewer',
			makeEditor: 'Make Editor',
			removeMember: 'Remove Member'
		},
		modal: {
			invite: {
				title: 'Invite Member',
				emailLabel: 'Email Address',
				emailPlaceholder: 'colleague@example.com',
				roleLabel: 'Role',
				roleEditor: 'Editor - Can edit content',
				roleViewer: 'Viewer - Can only view',
				roleEditorDesc: 'Editors can create, edit, and delete content.',
				roleViewerDesc: 'Viewers can only view content, not edit.',
				cancel: 'Cancel',
				send: 'Send Invitation',
				sending: 'Sending...'
			},
			success: {
				title: 'Invitation Sent!',
				description: 'Share this link with the user:',
				copy: 'Copy',
				copied: 'Copied!',
				expiresIn: 'Link expires in 7 days',
				done: 'Done'
			}
		},
		alerts: {
			removeConfirm: 'Are you sure you want to remove this member?'
		}
	}
};
