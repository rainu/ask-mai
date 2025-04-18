<template>
	<v-hover>
		<template v-slot:default="{ props: hoverProps, isHovering }">
			<v-app :theme="theme" v-bind="hoverProps" :style="{ opacity: opacityValue(isHovering) }">
				<v-main>
					<v-container class="pa-0 ma-0" fluid style="overflow-x: auto;">
						<router-view v-slot="{ Component }">
							<keep-alive>
								<component :is="Component" />
							</keep-alive>
						</router-view>
					</v-container>
				</v-main>
			</v-app>
		</template>
	</v-hover>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { TriggerRestart, Shutdown } from '../wailsjs/go/controller/Controller'
import { mapState } from 'pinia'
import { useConfigStore } from './store/config.ts'

export default defineComponent({
	computed: {
		...mapState(useConfigStore, ['config']),
		theme(): string {
			if (this.config.UI.Theme === 'system') {
				if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
					return 'dark'
				}
				return 'light'
			}

			return this.config.UI.Theme
		},
		opacity(): number {
			return this.config.UI.Window.BackgroundColor.A / 255
		}
	},
	methods: {
		handleGlobalKeydown(event: KeyboardEvent) {
			for (let i = 0; i < this.config.UI.QuitShortcut.Code.length; i++) {
				let code = event.code.toLowerCase() === this.config.UI.QuitShortcut.Code[i].toLowerCase()
				let ctrl = event.ctrlKey === this.config.UI.QuitShortcut.Ctrl[i]
				let shift = event.shiftKey === this.config.UI.QuitShortcut.Shift[i]
				let alt = event.altKey === this.config.UI.QuitShortcut.Alt[i]
				let meta = event.metaKey === this.config.UI.QuitShortcut.Meta[i]

				if (code && ctrl && shift && alt && meta) {
					Shutdown()
				}
			}

			for (let i = 0; i < this.config.Debug.RestartShortcut.Code.length; i++) {
				let code = event.code.toLowerCase() === this.config.Debug.RestartShortcut.Code[i].toLowerCase()
				let ctrl = event.ctrlKey === this.config.Debug.RestartShortcut.Ctrl[i]
				let shift = event.shiftKey === this.config.Debug.RestartShortcut.Shift[i]
				let alt = event.altKey === this.config.Debug.RestartShortcut.Alt[i]
				let meta = event.metaKey === this.config.Debug.RestartShortcut.Meta[i]

				if (code && ctrl && shift && alt && meta) {
					TriggerRestart()
				}
			}
		},
		opacityValue(isHovering: boolean | null): number {
			switch (this.config.UI.Window.Translucent) {
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
