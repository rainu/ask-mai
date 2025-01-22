<template>
	<v-hover>
		<template v-slot:default="{ props: hoverProps, isHovering }">
			<v-app :theme="theme" v-bind="hoverProps" :style="{ opacity: opacityValue(isHovering) }">
				<v-main>
					<v-container class="pa-0 ma-0" fluid style="overflow-x: auto;">
						<RouterView />
					</v-container>
				</v-main>
			</v-app>
		</template>
	</v-hover>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { TriggerRestart, Shutdown } from '../wailsjs/go/controller/Controller'

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
			opacity: this.$appConfig.UI.Window.BackgroundColor.A / 255,
		}
	},
	methods: {
		handleGlobalKeydown(event: KeyboardEvent) {
			let code = event.code.toLowerCase() === this.$appConfig.UI.QuitShortcut.Code.toLowerCase()
			let ctrl = event.ctrlKey === this.$appConfig.UI.QuitShortcut.Ctrl
			let shift = event.shiftKey === this.$appConfig.UI.QuitShortcut.Shift
			let alt = event.altKey === this.$appConfig.UI.QuitShortcut.Alt
			let meta = event.metaKey === this.$appConfig.UI.QuitShortcut.Meta

			if (code && ctrl && shift && alt && meta) {
				Shutdown()
			}

			code = event.code.toLowerCase() === this.$appConfig.Debug.RestartShortcut.Code.toLowerCase()
			ctrl = event.ctrlKey === this.$appConfig.Debug.RestartShortcut.Ctrl
			shift = event.shiftKey === this.$appConfig.Debug.RestartShortcut.Shift
			alt = event.altKey === this.$appConfig.Debug.RestartShortcut.Alt
			meta = event.metaKey === this.$appConfig.Debug.RestartShortcut.Meta
			if (code && ctrl && shift && alt && meta) {
				TriggerRestart()
			}
		},
		opacityValue(isHovering: boolean | null): number {
			switch (this.$appConfig.UI.Window.Translucent) {
				case 'ever':
					return this.opacity
				case 'hover':
					return isHovering ? 1 : this.opacity
			}
			return 1
		},
	},
	mounted() {
		window.addEventListener('keydown', this.handleGlobalKeydown)
	},
})
</script>

<style scoped></style>
