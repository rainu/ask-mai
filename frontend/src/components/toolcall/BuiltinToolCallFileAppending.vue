<template>
	<ToolCall :tc="tc" icon="mdi-file-document-edit">
		<template v-slot:title>
			<template v-if="parsedResult">
				{{ parsedResult.path }}
			</template>
			<template v-else>
				{{ parsedArguments.path }}
			</template>
		</template>

		<template v-slot:content>
			<vue-markdown :source="contentAsMarkdown"></vue-markdown>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, file } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import FileAppendingArguments = file.FileAppendingArguments
import FileAppendingResult = file.FileAppendingResult
import ToolCall from './ToolCall.vue'
import VueMarkdown from 'vue-markdown-render'

export default defineComponent({
	name: 'BuiltinToolCallFileAppending',
	components: { ToolCall, VueMarkdown },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): FileAppendingArguments {
			return JSON.parse(this.tc.Arguments) as FileAppendingArguments
		},
		parsedResult(): FileAppendingResult | null {
			if(this.tc.Result) {
				try {
					return JSON.parse(this.tc.Result.Content) as FileAppendingResult
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
		contentAsMarkdown(){
			return '```\n' + this.parsedArguments.content + '\n```'
		}
	},
})
</script>

<style scoped></style>
