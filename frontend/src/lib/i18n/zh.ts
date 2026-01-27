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
		backToLogin: '返回登录'
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
			languageZh: '中文'
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
