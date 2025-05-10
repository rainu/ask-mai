import './log.ts'
import './debug.ts'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { RouteMeta, Router } from 'vue-router'

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
import { controller, model } from '../wailsjs/go/models.ts'
import { GetApplicationConfig } from '../wailsjs/go/controller/Controller'

declare module '@vue/runtime-core' {
	interface ComponentCustomProperties {
		$route: RouteMeta
		$router: Router
		$i18n: I18nOptions
		$t: (key: string, ...args: any[]) => string
		$appConfig: model.Config
		$appProfile: model.Profile
	}
}

declare global {
  interface Window {
    $appConfig: model.Config
    $appProfile: model.Profile
  }
}

const pinia = createPinia()
const app = createApp(App)

GetApplicationConfig().then((ac: controller.ApplicationConfig) => {
	app.config.globalProperties.$appConfig = ac.Config
	app.config.globalProperties.$appProfile = ac.ActiveProfile
	window.$appConfig = ac.Config
	window.$appProfile = ac.ActiveProfile
	const langPrefix = ac.ActiveProfile.UI.Language.split('_')[0]
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
