import { createI18n } from 'vue-i18n'

import de from './locales/de.json'
import en from './locales/en.json'

const i18n = (locale: string) => createI18n({
	locale: locale,
	fallbackLocale: 'en',
	legacy: false,
	messages: {
		de: de,
		en: en,
	},
})

export default i18n
