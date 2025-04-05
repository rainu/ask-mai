<template>
	<InputRow>
		<template v-slot:prepend>
			<v-btn icon density="compact" @click="onNavigateBack">
				<v-icon>mdi-chat-processing-outline</v-icon>
			</v-btn>
		</template>
		<v-text-field
			v-model="query"
			@change="onQueryChanged"
			hide-details
			autofocus
			:placeholder="$t('history.placeholder')"
		>
		</v-text-field>
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