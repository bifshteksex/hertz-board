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
		backToLogin: 'Вернуться к входу'
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
			deleteWarning: 'После удаления аккаунта восстановление невозможно. Пожалуйста, будьте уверены.',
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
			languageZh: '中文'
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
