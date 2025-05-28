<template>
	<ToolCall :tc="tc" icon="mdi-file-document-plus">
		<template v-slot:title>
			<template v-if="parsedResult">
				{{ parsedResult.path }}
			</template>
			<template v-else>
				{{ parsedArguments.path }}
			</template>
		</template>

		<template v-slot:content>
			<Markdown :content="contentAsMarkdown" />
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, file } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import FileCreationArguments = file.FileCreationArguments
import FileCreationResult = file.FileCreationResult
import ToolCall from './ToolCall.vue'
import Markdown from '../Markdown.vue'

export default defineComponent({
	name: 'BuiltinToolCallFileCreation',
	components: { Markdown, ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): FileCreationArguments {
			return JSON.parse(this.tc.Arguments) as FileCreationArguments
		},
		parsedResult(): FileCreationResult | null {
			if(this.tc.Result) {
				try {
					return JSON.parse(this.tc.Result.Content) as FileCreationResult
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
		contentAsMarkdown(): string {
			return '```\n' + this.parsedArguments.content + '\n```'
		}
	},
})
</script>

<style scoped></style>
