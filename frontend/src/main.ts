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

declare module '@vue/runtime-core' {
	interface ComponentCustomProperties {
		$i18n: I18nOptions
		$t: typeof i18n.global.t
	}
}

createApp(App)
	.use(router)
	.use(i18n)
	.use(
		createVuetify({
			components,
			directives,
		}),
	)
	.mount('#app')
