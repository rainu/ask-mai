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
			<vue-markdown v-if="parsedResult" :source="contentAsMarkdown"></vue-markdown>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, file } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import FileReadingArguments = file.FileReadingArguments
import FileReadingResult = file.FileReadingResult
import ToolCall from './ToolCall.vue'
import VueMarkdown from 'vue-markdown-render'

export default defineComponent({
	name: 'BuiltinToolCallFileReading',
	components: { ToolCall, VueMarkdown },
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
		},
		contentAsMarkdown(){
			if(this.parsedResult) {
				return '```\n' + this.parsedResult.content + '\n```'
			}
			return null
		}
	},
})
</script>

<style scoped></style>
