import { i18n } from '$lib/stores/i18n.svelte';
import en from './en';
import ru from './ru';
import zh from './zh';

// Initialize translations
i18n.setTranslations('en', en);
i18n.setTranslations('ru', ru);
i18n.setTranslations('zh', zh);

export { i18n };
export type { default as Translations } from './en';
