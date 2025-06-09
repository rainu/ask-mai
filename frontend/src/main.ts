import './log.ts'
import './debug.ts'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { RouteMeta, Router } from 'vue-router'

// Vuetify
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { createVuetify, ThemeDefinition } from 'vuetify'
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

	const vuetifyOpts = {
		components,
		directives,
		theme: {
			themes: {} as Record<string, ThemeDefinition>,
		}
	}

	const themes = {} as Record<string, ThemeDefinition>
	if(ac.Config.Themes.dark) themes["dark"] = ac.Config.Themes.dark
	if(ac.Config.Themes.light) themes["light"] = ac.Config.Themes.light
	for (let name in ac.Config.Themes.custom) {
		themes[name] = ac.Config.Themes.custom[name] as ThemeDefinition
	}

	for (let name in themes) {
		vuetifyOpts.theme.themes[name] = {} as ThemeDefinition
		vuetifyOpts.theme.themes[name].dark = themes[name].dark
		vuetifyOpts.theme.themes[name].colors = themes[name].colors || {}
		vuetifyOpts.theme.themes[name].variables = themes[name].variables || {}
	}
	console.log(vuetifyOpts)

	app
		.use(pinia)
		.use(router)
		.use(i18n(langPrefix))
		.use(createVuetify(vuetifyOpts))
		.mount('#app')
})
