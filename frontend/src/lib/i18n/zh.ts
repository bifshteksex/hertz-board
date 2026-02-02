export default {
	seo: {
		defaultTitle: 'HertzBoard - 实时协作工作空间',
		defaultDescription:
			'使用HertzBoard实时创建、协作和可视化创意 - 基于现代技术构建的强大协作白板平台。',
		keywords: '协作白板, 实时协作, 可视化工作空间, 在线白板, 团队协作, 画布编辑器',
		ogSiteName: 'HertzBoard',
		twitterCard: 'summary_large_image',
		pages: {
			landing: {
				title: 'HertzBoard - 实时协作工作空间',
				description: '在无限画布上与您的团队实时创建、协作和可视化创意。'
			},
			login: {
				title: '登录 - HertzBoard',
				description: '登录您的HertzBoard账户以访问您的协作工作空间。'
			},
			register: {
				title: '创建账户 - HertzBoard',
				description: '创建新的HertzBoard账户，立即开始与您的团队协作。'
			},
			dashboard: {
				title: '仪表板 - HertzBoard',
				description: '管理您的协作工作空间和看板。'
			},
			settings: {
				title: '设置 - HertzBoard',
				description: '管理您的账户设置和偏好。'
			},
			workspace: {
				title: '{name} - HertzBoard',
				description: '在{name}工作空间中实时协作。'
			}
		}
	},
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
		footer: '用 {heart} 制作'
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
	},
	canvas: {
		toolbar: {
			undo: '撤销',
			redo: '重做',
			undoShortcut: '撤销 (Ctrl+Z)',
			redoShortcut: '重做 (Ctrl+Y)',
			help: '键盘快捷键 (?)',
			helpAria: '显示键盘快捷键',
			tools: {
				select: '选择 (V)',
				text: '文本 (T)',
				pen: '画笔 (P)',
				sticky: '便签 (S)',
				image: '图片 (I)',
				connector: '连接线',
				shapes: '形状',
				lists: '列表'
			},
			shapes: {
				rectangle: '矩形',
				ellipse: '圆形',
				triangle: '三角形',
				line: '直线',
				arrow: '箭头'
			},
			lists: {
				bullet: '项目符号列表',
				numbered: '编号列表',
				checkbox: '复选框列表'
			},
			zorder: {
				bringToFront: '置于顶层 (Ctrl+Shift+])',
				bringForward: '上移一层 (Ctrl+])',
				sendBackward: '下移一层 (Ctrl+[)',
				sendToBack: '置于底层 (Ctrl+Shift+[)'
			},
			grouping: {
				group: '编组 (Ctrl+G)',
				ungroup: '取消编组 (Ctrl+Shift+G)',
				groupElements: '编组元素',
				ungroupElements: '取消编组元素'
			},
			view: {
				toggleGrid: "切换网格 (Ctrl+')",
				toggleSnap: "吸附到网格 (Ctrl+Shift+')",
				gridAria: '切换网格',
				snapAria: '切换网格吸附'
			}
		},
		shortcuts: {
			title: '键盘快捷键',
			subtitle: '学习所有快捷键以提高工作效率',
			search: '搜索快捷键...',
			noResults: '未找到"{query}"的快捷键',
			footer: '按',
			footerKey: '?',
			footerText: '随时显示此面板',
			categories: {
				tools: '工具',
				edit: '编辑',
				selection: '选择',
				layers: '图层',
				grouping: '编组',
				view: '视图',
				other: '其他'
			}
		},
		workspace: {
			loading: '加载工作区中...',
			backToDashboard: '返回仪表板',
			toolbar: {
				layers: '图层',
				share: '分享'
			}
		},
		contextMenu: {
			cut: '剪切',
			copy: '复制',
			paste: '粘贴',
			duplicate: '复制',
			delete: '删除',
			bringToFront: '置于顶层',
			sendToBack: '置于底层',
			group: '编组',
			ungroup: '取消编组',
			lock: '锁定',
			unlock: '解锁',
			copyLink: '复制链接'
		},
		imageUploader: {
			uploading: '上传图片中...',
			dismiss: '关闭',
			errorFileType: '请选择图片文件',
			errorFileSize: '图片大小必须小于10MB',
			errorUploadFailed: '上传图片失败'
		},
		saveStatus: {
			saving: '保存中...',
			saved: '已保存',
			savedAgo: '{time}前保存',
			saveFailed: '保存失败',
			unsavedChanges: '{count}个未保存的更改',
			allSaved: '所有更改已保存',
			details: '详情'
		}
	}
};
