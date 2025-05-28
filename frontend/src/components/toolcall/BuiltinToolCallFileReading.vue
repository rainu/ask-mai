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
			<Markdown v-if="parsedResult" :content="contentAsMarkdown" />
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
import Markdown from '../Markdown.vue'

export default defineComponent({
	name: 'BuiltinToolCallFileReading',
	components: { Markdown, ToolCall },
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
		contentAsMarkdown(): string {
			if(this.parsedResult) {
				return '```\n' + this.parsedResult.content + '\n```'
			}
			return ''
		}
	},
})
</script>

<style scoped></style>
