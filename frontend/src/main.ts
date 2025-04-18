import './log.ts'
import './debug.ts'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { Router } from 'vue-router'

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
import { model } from '../wailsjs/go/models.ts'
import { GetApplicationConfig } from '../wailsjs/go/controller/Controller'

declare module '@vue/runtime-core' {
	interface ComponentCustomProperties {
		$router: Router
		$i18n: I18nOptions
		$t: (key: string, ...args: any[]) => string
		$appConfig: model.Config
	}
}

declare global {
  interface Window {
    $appConfig: model.Config
  }
}

const pinia = createPinia()
const app = createApp(App)

GetApplicationConfig().then((cfg: model.Config) => {
	app.config.globalProperties.$appConfig = cfg
	window.$appConfig = cfg
	const langPrefix = cfg.UI.Language.split('_')[0]
	app
		.use(pinia)
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
