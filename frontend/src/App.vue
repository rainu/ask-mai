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
		return {
			theme: 'light',
		}
	},
})
</script>

<style scoped></style>
