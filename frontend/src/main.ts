import { createApp } from 'vue'

// Vuetify
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import App from './App.vue'
import router from './router'
import i18n from './i18n'
import { I18nOptions } from 'vue-i18n'
import { config } from '../wailsjs/go/models.ts'
import { GetApplicationConfig } from '../wailsjs/go/controller/Controller'

declare module '@vue/runtime-core' {
	interface ComponentCustomProperties {
		$i18n: I18nOptions
		$t: (key: string, ...args: any[]) => string
		$appConfig: config.Config
	}
}

const app = createApp(App)

GetApplicationConfig().then((cfg: config.Config) => {
	app.config.globalProperties.$appConfig = cfg
	const langPrefix = cfg.UI.Language.split('_')[0]
	app
		.use(router)
		.use(i18n(langPrefix))
		.use(
			createVuetify({
				components,
				directives,
			}),
		)
		.mount('#app')
})
