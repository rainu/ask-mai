<template>
	<div class="text-center">
		<v-chip prepend-icon="mdi-import" class="mx-2" variant="outlined" label v-if="inputToken !== undefined">
			{{inputToken.toLocaleString()}}
		</v-chip>
		<v-chip prepend-icon="mdi-export" class="mx-2" variant="outlined" label v-if="outputToken !== undefined">
			{{outputToken.toLocaleString()}}
		</v-chip>

		<v-chip class="mx-2" variant="outlined" label v-for="(value, key) in additional" :key="key">
			{{key}}: {{value.toLocaleString()}}
		</v-chip>
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { HistoryEntryConsumption } from '../store/history.ts'

export default defineComponent({
	name: 'Consumption',
	props: {
		model: {
			type: Object as () => HistoryEntryConsumption,
			required: true,
		},
	},
	computed: {
		inputToken(): number | undefined {
			if('prompt' in this.model) {
				return this.model["prompt"]
			}
			return undefined
		},
		outputToken(): number | undefined {
			if('completion' in this.model) {
				return this.model["completion"]
			}
			return undefined
		},
		additional(): HistoryEntryConsumption {
			const result: HistoryEntryConsumption = {}
			for (const [key, value] of Object.entries(this.model)) {
				if (key !== 'completion' && key !== 'prompt' && value !== 0) {
					result[key] = value
				}
			}
			return result
		}
	}
})
</script>

<style scoped>

</style>