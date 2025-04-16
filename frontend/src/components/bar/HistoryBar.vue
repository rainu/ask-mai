<template>
	<GeneralBar :minimizable="false">
		<v-text-field
			v-model="query"
			@change="onQueryChanged"
			hide-details
			autofocus
			:placeholder="$t('history.placeholder')"
		>
		</v-text-field>

		<template v-slot:option-buttons>
			<v-btn @click="onNavigateBack">
				<v-icon size="x-large">mdi-chat-processing-outline</v-icon>
			</v-btn>
		</template>
	</GeneralBar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import GeneralBar from './GeneralBar.vue'

export default defineComponent({
	name: 'HistoryBar',
	components: { GeneralBar },
	emits: ['queryChanged'],
	data(){
		return {
			query: '',
		}
	},
	methods: {
		onNavigateBack() {
			this.$router.push({ name: 'Chat' })
		},
		onQueryChanged() {
			this.$emit('queryChanged', this.query)
		}
	}
})
</script>

<style scoped></style>