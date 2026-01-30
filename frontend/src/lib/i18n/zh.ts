export default {
	common: {
		loading: '加载中...',
		save: '保存',
		cancel: '取消',
		delete: '删除',
		edit: '编辑',
		close: '关闭',
		search: '搜索',
		create: '创建',
		update: '更新',
		confirm: '确认',
		back: '返回'
	},
	nav: {
		dashboard: '仪表板',
		workspaces: '工作区',
		settings: '设置',
		logout: '登出'
	},
	auth: {
		appName: 'HertzBoard',
		login: '登录',
		register: '注册',
		email: '邮箱',
		password: '密码',
		name: '全名',
		confirmPassword: '确认密码',
		forgotPassword: '忘记密码？',
		noAccount: '没有账号？',
		hasAccount: '已有账号？',
		signIn: '登录',
		signUp: '注册',
		orContinueWith: '或继续使用',
		resetPassword: '重置密码',
		backToLogin: '返回登录',
		loginTitle: '登录您的账户',
		registerTitle: '创建您的账户',
		fullName: '全名',
		emailPlaceholder: '邮箱地址',
		passwordPlaceholder: '密码',
		passwordMinPlaceholder: '密码（最少8个字符）',
		confirmPasswordPlaceholder: '确认密码',
		signingIn: '登录中...',
		creatingAccount: '创建账户中...',
		googleButton: 'Google',
		githubButton: 'GitHub',
		resetTitle: '重置您的密码',
		resetDescription: '输入您的邮箱地址，我们将向您发送重置密码的链接。',
		resetSuccess:
			'请检查您的邮箱以获取重置密码的链接。如果几分钟内没有收到，请检查您的垃圾邮件文件夹。',
		sendResetLink: '发送重置链接',
		sending: '发送中...',
		errorPasswordMatch: '密码不匹配',
		errorPasswordLength: '密码必须至少8个字符'
	},
	dashboard: {
		title: '工作区',
		subtitle: '管理您的协作看板',
		newWorkspace: '新建工作区',
		searchPlaceholder: '搜索工作区...',
		loading: '加载工作区中...',
		noWorkspaces: '未找到工作区',
		createFirst: '创建您的第一个工作区',
		member: '成员',
		members: '成员',
		role: '角色',
		menu: {
			open: '打开',
			rename: '重命名',
			copyLink: '复制链接',
			share: '分享',
			duplicate: '复制',
			delete: '删除'
		},
		modal: {
			create: {
				title: '创建新工作区',
				name: '名称',
				description: '描述（可选）',
				namePlaceholder: '我的工作区',
				descriptionPlaceholder: '描述您的工作区...',
				creating: '创建中...',
				create: '创建'
			},
			duplicate: {
				title: '复制工作区',
				newName: '新工作区名称',
				copyOf: '这将创建"{name}"的副本',
				duplicating: '复制中...',
				duplicate: '复制'
			},
			rename: {
				title: '重命名工作区',
				saving: '保存中...',
				saveChanges: '保存更改'
			}
		},
		alerts: {
			deleteConfirm: '您确定要删除此工作区吗？',
			shareComingSoon: '分享功能将在未来版本中实现',
			linkCopied: '链接已复制到剪贴板！'
		},
		time: {
			today: '今天',
			yesterday: '昨天',
			daysAgo: '{count}天前',
			weeksAgo: '{count}周前',
			monthsAgo: '{count}月前'
		}
	},
	settings: {
		title: '设置',
		subtitle: '管理您的账户设置和偏好',
		tabs: {
			profile: '个人资料',
			password: '密码',
			account: '账户',
			preferences: '偏好设置'
		},
		profile: {
			title: '个人资料信息',
			fullName: '全名',
			avatarUrl: '头像网址',
			avatarPlaceholder: 'https://example.com/avatar.jpg',
			avatarHint: '可选：输入您的头像图片网址',
			provider: '提供商',
			saveChanges: '保存更改',
			saving: '保存中...',
			successMessage: '个人资料更新成功！'
		},
		password: {
			title: '更改密码',
			current: '当前密码',
			new: '新密码',
			confirm: '确认新密码',
			hint: '最少8个字符',
			change: '更改密码',
			changing: '更改中...',
			successMessage: '密码更改成功！',
			oauthWarning: '您使用 {provider} 登录。OAuth账户无法更改密码。',
			errorMatch: '密码不匹配',
			errorLength: '密码必须至少8个字符'
		},
		account: {
			title: '账户信息',
			email: '邮箱',
			accountType: '账户类型',
			emailVerified: '邮箱已验证',
			memberSince: '注册时间',
			verified: '是',
			notVerified: '否',
			dangerZone: '危险区域',
			deleteWarning: '一旦删除您的账户，将无法恢复。请确认。',
			deleteAccount: '删除账户'
		},
		preferences: {
			title: '偏好设置',
			language: '语言',
			theme: '主题',
			themeLight: '浅色',
			themeDark: '深色',
			languageEn: 'English',
			languageRu: 'Русский',
			languageZh: '中文',
			languageLabel: '语言',
			themeLabel: '主题',
			languageHint: '选择您喜欢的界面语言',
			themeHint: '选择您喜欢的配色方案'
		}
	},
	landing: {
		subtitle: '实时协作工作空间平台',
		login: '登录',
		signUp: '注册',
		features: {
			realtime: {
				title: '实时协作',
				description: '与您的团队实时协作，具有实时光标和在线状态'
			},
			canvas: {
				title: '强大的画布',
				description: '在无限画布上使用文本、形状、图像等进行创作'
			},
			tech: {
				title: '基于 Go 构建',
				description: '由 CloudWeGo Hertz 和 WebSockets 提供支持的高性能后端'
			}
		},
		footer: '用 ❤️ 制作'
	},
	workspaceDetail: {
		loading: '加载工作区中...',
		backToDashboard: '返回仪表板',
		share: '分享',
		canvasComingSoon: '画布即将推出',
		canvasDescription: '画布编辑器将在第6阶段实现',
		workspaceId: '工作区 ID:',
		role: '角色:',
		errors: {
			missingId: '工作区 ID 缺失',
			loadFailed: '加载工作区失败'
		}
	},
	workspace: {
		title: '工作区',
		createNew: '创建工作区',
		myWorkspaces: '我的工作区',
		recentWorkspaces: '最近的工作区',
		noWorkspaces: '还没有工作区',
		createFirst: '创建您的第一个工作区开始使用'
	},
	errors: {
		generic: '发生错误',
		network: '网络错误。请检查您的连接。',
		unauthorized: '未授权。请登录。',
		notFound: '未找到',
		serverError: '服务器错误。请稍后再试。'
	}
};
