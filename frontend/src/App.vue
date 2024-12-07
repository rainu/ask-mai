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
})
</script>

<style scoped></style>
