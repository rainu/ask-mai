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
		let theme = 'light'
		if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
			theme = 'dark'
		}

		return {
			theme,
		}
	},
	methods: {
		handleGlobalKeydown(event: KeyboardEvent) {
			if (event.key === 'Escape') {
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
