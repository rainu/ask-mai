<template>
	<div class="text-center">
		<v-chip prepend-icon="mdi-import" class="mx-2" variant="outlined" label v-if="inputToken !== undefined">
			{{inputToken.toLocaleString()}}
		</v-chip>
		<v-chip prepend-icon="mdi-cached" class="mx-2" variant="outlined" label v-if="inputToken !== undefined && cachedToken">
			{{cachedToken.toLocaleString()}}
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
			if('input' in this.model) {
				return this.model["input"]
			}
			return undefined
		},
		cachedToken(): number {
			if('cached' in this.model) {
				return this.model["cached"]
			}
			return 0
		},
		outputToken(): number | undefined {
			if('output' in this.model) {
				return this.model["output"]
			}
			return undefined
		},
		additional(): HistoryEntryConsumption {
			const result: HistoryEntryConsumption = {}
			for (const [key, value] of Object.entries(this.model)) {
				if (key !== 'input' && key !== 'output' && key !== 'cached' && value !== 0) {
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