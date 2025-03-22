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
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import DirectoryDeletionArguments = tools.DirectoryDeletionArguments
import DirectoryDeletionResult = tools.DirectoryDeletionResult
import ToolCall from './ToolCall.vue'

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
					return JSON.parse(this.tc.Result.Content) as DirectoryDeletionResult
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
