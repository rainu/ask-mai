<template>
	<ToolCall :tc="tc" icon="mdi-folder-plus">
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
import DirectoryCreationArguments = file.DirectoryCreationArguments
import DirectoryCreationResult = file.DirectoryCreationResult
import ToolCall from './ToolCall.vue'
import { ToolCallResult } from './types.ts'

export default defineComponent({
	name: 'BuiltinToolCallDirectoryCreation',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): DirectoryCreationArguments {
			return JSON.parse(this.tc.Arguments) as DirectoryCreationArguments
		},
		parsedResult(): DirectoryCreationResult | null {
			if(this.tc.Result) {
				try {
					const tcr = JSON.parse(this.tc.Result.Content) as ToolCallResult
					const firstTextContent = tcr.content.find(c => c.type === 'text') as mcp.TextContent
					if(firstTextContent) {
						return JSON.parse(firstTextContent.text) as DirectoryCreationResult
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
