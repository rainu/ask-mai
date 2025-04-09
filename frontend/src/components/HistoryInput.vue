<template>
	<InputRow :minimizable="false">
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
	</InputRow>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import InputRow from './InputRow.vue'

export default defineComponent({
	name: 'HistoryInput',
	components: { InputRow },
	emits: ['queryChanged'],
	data(){
		return {
			query: '',
		}
	},
	methods: {
		onNavigateBack() {
			this.$router.push({ name: 'Home' })
		},
		onQueryChanged() {
			this.$emit('queryChanged', this.query)
		}
	}
})
</script>

<style scoped></style>