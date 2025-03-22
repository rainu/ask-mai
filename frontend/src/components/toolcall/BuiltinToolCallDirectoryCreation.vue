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
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import DirectoryCreationArguments = tools.DirectoryCreationArguments
import DirectoryCreationResult = tools.DirectoryCreationResult
import ToolCall from './ToolCall.vue'

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
					return JSON.parse(this.tc.Result.Content) as DirectoryCreationResult
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
