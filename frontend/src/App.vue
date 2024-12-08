<template>
	<v-app :theme="theme">
		<v-main>
			<v-container class="pa-0 ma-0" fluid>
				<RouterView />
			</v-container>
		</v-main>
	</v-app>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { Quit } from '../wailsjs/runtime'

export default defineComponent({
	data() {
		let theme = ''

		if (this.$appConfig.UI.Theme === 'system') {
			theme = 'light'
			if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
				theme = 'dark'
			}
		} else {
			theme = this.$appConfig.UI.Theme
		}

		return {
			theme,
		}
	},
	methods: {
		handleGlobalKeydown(event: KeyboardEvent) {
			const code = event.code.toLowerCase() === this.$appConfig.UI.QuitShortcut.Code.toLowerCase()
			const ctrl = event.ctrlKey === this.$appConfig.UI.QuitShortcut.Ctrl
			const shift = event.shiftKey === this.$appConfig.UI.QuitShortcut.Shift
			const alt = event.altKey === this.$appConfig.UI.QuitShortcut.Alt
			const meta = event.metaKey === this.$appConfig.UI.QuitShortcut.Meta

			if (code && ctrl && shift && alt && meta) {
				Quit()
			}
		},
	},
	mounted() {
		window.addEventListener('keydown', this.handleGlobalKeydown)
	},
})
</script>

<style scoped></style>
