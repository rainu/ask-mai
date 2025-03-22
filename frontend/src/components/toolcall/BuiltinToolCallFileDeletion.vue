<template>
	<ToolCall :tc="tc" icon="mdi-file-document-remove">
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
import FileDeletionArguments = tools.FileDeletionArguments
import FileDeletionResult = tools.FileDeletionResult
import ToolCall from './ToolCall.vue'

export default defineComponent({
	name: 'BuiltinToolCallFileDeletion',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): FileDeletionArguments {
			return JSON.parse(this.tc.Arguments) as FileDeletionArguments
		},
		parsedResult(): FileDeletionResult | null {
			if(this.tc.Result) {
				try {
					return JSON.parse(this.tc.Result.Content) as FileDeletionResult
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
