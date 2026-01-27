type Locale = 'en' | 'ru' | 'zh';

interface Translations {
	[key: string]: string | Translations;
}

interface LocaleData {
	[key: string]: Translations;
}

class I18nStore {
	private currentLocale = $state<Locale>('en');
	private translations = $state<LocaleData>({});

	constructor() {
		// Load saved locale from localStorage on initialization
		if (typeof window !== 'undefined') {
			const savedLocale = localStorage.getItem('locale') as Locale | null;
			if (savedLocale && ['en', 'ru', 'zh'].includes(savedLocale)) {
				this.currentLocale = savedLocale;
			}
		}
	}

	get locale(): Locale {
		return this.currentLocale;
	}

	setLocale(locale: Locale) {
		this.currentLocale = locale;
		if (typeof window !== 'undefined') {
			localStorage.setItem('locale', locale);
		}
	}

	setTranslations(locale: Locale, data: Translations) {
		this.translations[locale] = data;
	}

	t(key: string, params?: Record<string, string | number>): string {
		const keys = key.split('.');
		let value: string | Translations | undefined = this.translations[this.currentLocale];

		for (const k of keys) {
			if (value && typeof value === 'object') {
				value = value[k];
			} else {
				value = undefined;
				break;
			}
		}

		let result = typeof value === 'string' ? value : key;

		// Replace parameters
		if (params) {
			Object.entries(params).forEach(([param, val]) => {
				result = result.replace(new RegExp(`\\{${param}\\}`, 'g'), String(val));
			});
		}

		return result;
	}
}

export const i18n = new I18nStore();
