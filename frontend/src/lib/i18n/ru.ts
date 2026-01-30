export default {
	common: {
		loading: 'Загрузка...',
		save: 'Сохранить',
		cancel: 'Отмена',
		delete: 'Удалить',
		edit: 'Редактировать',
		close: 'Закрыть',
		search: 'Поиск',
		create: 'Создать',
		update: 'Обновить',
		confirm: 'Подтвердить',
		back: 'Назад'
	},
	nav: {
		dashboard: 'Панель',
		workspaces: 'Рабочие области',
		settings: 'Настройки',
		logout: 'Выйти'
	},
	auth: {
		appName: 'HertzBoard',
		login: 'Войти',
		register: 'Регистрация',
		email: 'Email',
		password: 'Пароль',
		name: 'Полное имя',
		confirmPassword: 'Подтвердите пароль',
		forgotPassword: 'Забыли пароль?',
		noAccount: 'Нет аккаунта?',
		hasAccount: 'Уже есть аккаунт?',
		signIn: 'Войти',
		signUp: 'Зарегистрироваться',
		orContinueWith: 'Или продолжить с',
		resetPassword: 'Сбросить пароль',
		backToLogin: 'Вернуться к входу',
		loginTitle: 'Войдите в свой аккаунт',
		registerTitle: 'Создайте свой аккаунт',
		fullName: 'Полное имя',
		emailPlaceholder: 'Адрес электронной почты',
		passwordPlaceholder: 'Пароль',
		passwordMinPlaceholder: 'Пароль (минимум 8 символов)',
		confirmPasswordPlaceholder: 'Подтвердите пароль',
		signingIn: 'Вход...',
		creatingAccount: 'Создание аккаунта...',
		googleButton: 'Google',
		githubButton: 'GitHub',
		resetTitle: 'Сброс пароля',
		resetDescription: 'Введите ваш email и мы отправим вам ссылку для сброса пароля.',
		resetSuccess:
			'Проверьте вашу почту на наличие ссылки для сброса пароля. Если письмо не появится в течение нескольких минут, проверьте папку со спамом.',
		sendResetLink: 'Отправить ссылку для сброса',
		sending: 'Отправка...',
		errorPasswordMatch: 'Пароли не совпадают',
		errorPasswordLength: 'Пароль должен содержать минимум 8 символов'
	},
	dashboard: {
		title: 'Рабочие области',
		subtitle: 'Управление вашими совместными досками',
		newWorkspace: 'Новая рабочая область',
		searchPlaceholder: 'Поиск рабочих областей...',
		loading: 'Загрузка рабочих областей...',
		noWorkspaces: 'Рабочие области не найдены',
		createFirst: 'Создайте свою первую рабочую область',
		member: 'участник',
		members: 'участников',
		role: 'Роль',
		menu: {
			open: 'Открыть',
			rename: 'Переименовать',
			copyLink: 'Копировать ссылку',
			share: 'Поделиться',
			duplicate: 'Дублировать',
			delete: 'Удалить'
		},
		modal: {
			create: {
				title: 'Создать новую рабочую область',
				name: 'Название',
				description: 'Описание (опционально)',
				namePlaceholder: 'Моя рабочая область',
				descriptionPlaceholder: 'Опишите вашу рабочую область...',
				creating: 'Создание...',
				create: 'Создать'
			},
			duplicate: {
				title: 'Дублировать рабочую область',
				newName: 'Название новой рабочей области',
				copyOf: 'Это создаст копию "{name}"',
				duplicating: 'Дублирование...',
				duplicate: 'Дублировать'
			},
			rename: {
				title: 'Переименовать рабочую область',
				saving: 'Сохранение...',
				saveChanges: 'Сохранить изменения'
			}
		},
		alerts: {
			deleteConfirm: 'Вы уверены, что хотите удалить эту рабочую область?',
			shareComingSoon: 'Функция общего доступа будет реализована в будущей версии',
			linkCopied: 'Ссылка скопирована в буфер обмена!'
		},
		time: {
			today: 'Сегодня',
			yesterday: 'Вчера',
			daysAgo: '{count} дней назад',
			weeksAgo: '{count} недель назад',
			monthsAgo: '{count} месяцев назад'
		}
	},
	settings: {
		title: 'Настройки',
		subtitle: 'Управление настройками аккаунта и предпочтениями',
		tabs: {
			profile: 'Профиль',
			password: 'Пароль',
			account: 'Аккаунт',
			preferences: 'Предпочтения'
		},
		profile: {
			title: 'Информация профиля',
			fullName: 'Полное имя',
			avatarUrl: 'URL аватара',
			avatarPlaceholder: 'https://example.com/avatar.jpg',
			avatarHint: 'Опционально: Введите URL вашей фотографии профиля',
			provider: 'Провайдер',
			saveChanges: 'Сохранить изменения',
			saving: 'Сохранение...',
			successMessage: 'Профиль успешно обновлен!'
		},
		password: {
			title: 'Изменить пароль',
			current: 'Текущий пароль',
			new: 'Новый пароль',
			confirm: 'Подтвердите новый пароль',
			hint: 'Минимум 8 символов',
			change: 'Изменить пароль',
			changing: 'Изменение...',
			successMessage: 'Пароль успешно изменен!',
			oauthWarning: 'Вы вошли через {provider}. Изменение пароля недоступно для OAuth аккаунтов.',
			errorMatch: 'Пароли не совпадают',
			errorLength: 'Пароль должен содержать минимум 8 символов'
		},
		account: {
			title: 'Информация об аккаунте',
			email: 'Email',
			accountType: 'Тип аккаунта',
			emailVerified: 'Email подтвержден',
			memberSince: 'Участник с',
			verified: 'Да',
			notVerified: 'Нет',
			dangerZone: 'Опасная зона',
			deleteWarning:
				'После удаления аккаунта восстановление невозможно. Пожалуйста, будьте уверены.',
			deleteAccount: 'Удалить аккаунт'
		},
		preferences: {
			title: 'Предпочтения',
			language: 'Язык',
			theme: 'Тема',
			themeLight: 'Светлая',
			themeDark: 'Темная',
			languageEn: 'English',
			languageRu: 'Русский',
			languageZh: '中文',
			languageLabel: 'Язык',
			themeLabel: 'Тема',
			languageHint: 'Выберите предпочитаемый язык интерфейса',
			themeHint: 'Выберите предпочитаемую цветовую схему'
		}
	},
	landing: {
		subtitle: 'Платформа для совместной работы в реальном времени',
		login: 'Войти',
		signUp: 'Регистрация',
		features: {
			realtime: {
				title: 'Совместная работа в реальном времени',
				description:
					'Работайте вместе с вашей командой в реальном времени с живыми курсорами и присутствием'
			},
			canvas: {
				title: 'Мощный холст',
				description:
					'Создавайте с помощью текста, фигур, изображений и многого другого на бесконечном холсте'
			},
			tech: {
				title: 'Построено на Go',
				description: 'Высокопроизводительный бэкенд на базе CloudWeGo Hertz и WebSockets'
			}
		},
		footer: 'Сделано с ❤️'
	},
	workspaceDetail: {
		loading: 'Загрузка рабочей области...',
		backToDashboard: 'Назад к панели',
		share: 'Поделиться',
		canvasComingSoon: 'Холст скоро появится',
		canvasDescription: 'Редактор холста будет реализован в фазе 6',
		workspaceId: 'ID рабочей области:',
		role: 'Роль:',
		errors: {
			missingId: 'Отсутствует ID рабочей области',
			loadFailed: 'Не удалось загрузить рабочую область'
		}
	},
	workspace: {
		title: 'Рабочие области',
		createNew: 'Создать рабочую область',
		myWorkspaces: 'Мои рабочие области',
		recentWorkspaces: 'Недавние рабочие области',
		noWorkspaces: 'Пока нет рабочих областей',
		createFirst: 'Создайте свою первую рабочую область для начала работы'
	},
	errors: {
		generic: 'Произошла ошибка',
		network: 'Ошибка сети. Проверьте подключение.',
		unauthorized: 'Не авторизован. Пожалуйста, войдите.',
		notFound: 'Не найдено',
		serverError: 'Ошибка сервера. Попробуйте позже.'
	}
};
