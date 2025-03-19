<template>
	<ToolCall :tc="tc" icon="mdi-file-eye">
		<template v-slot:title>
			<template v-if="parsedResult">
				{{ parsedResult.path }}
			</template>
			<template v-else>
				{{ parsedArguments.path }}
			</template>
		</template>

		<template v-slot:content>
			<div v-if="parsedArguments.limits">
				<span>{{ parsedArguments.limits.m }}: {{ offset }} - {{ limit }}</span>
			</div>
			<div v-if="parsedResult">{{ parsedResult.content }}</div>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import FileReadingArguments = tools.FileReadingArguments
import FileReadingResult = tools.FileReadingResult
import ToolCall from './ToolCall.vue'

export default defineComponent({
	name: 'BuiltinToolCallFileReading',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): FileReadingArguments {
			return JSON.parse(this.tc.Arguments) as FileReadingArguments
		},
		parsedResult(): FileReadingResult | null {
			if(this.tc.Result) {
				try {
					return JSON.parse(this.tc.Result.Content) as FileReadingResult
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
		offset(): number {
			return this.parsedArguments.limits?.o || 0
		},
		limit(): number {
			return this.parsedArguments.limits?.l || 0
		}
	},
})
</script>

<style scoped></style>
