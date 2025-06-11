<template>
	<ToolCall :tc="tc" icon="mdi-folder-remove">
		<template v-slot:title>
			<template v-if="parsedResult">
				{{ parsedResult.path }}
			</template>
			<template v-else>
				{{ parsedArguments.path }}
			</template>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, file, mcp } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import DirectoryDeletionArguments = file.DirectoryDeletionArguments
import DirectoryDeletionResult = file.DirectoryDeletionResult
import ToolCall from './ToolCall.vue'
import { ToolCallResult } from './types.ts'

export default defineComponent({
	name: 'BuiltinToolCallDirectoryDeletion',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): DirectoryDeletionArguments {
			return JSON.parse(this.tc.Arguments) as DirectoryDeletionArguments
		},
		parsedResult(): DirectoryDeletionResult | null {
			if(this.tc.Result) {
				try {
					const tcr = JSON.parse(this.tc.Result.Content) as ToolCallResult
					const firstTextContent = tcr.content.find(c => c.type === 'text') as mcp.TextContent
					if(firstTextContent) {
						return JSON.parse(firstTextContent.text) as DirectoryDeletionResult
					}
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
	},
})
</script>

<style scoped></style>
