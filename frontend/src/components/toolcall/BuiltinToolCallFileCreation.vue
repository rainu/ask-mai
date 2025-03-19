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
			<div>{{ parsedArguments.content }}</div>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import FileCreationArguments = tools.FileCreationArguments
import FileCreationResult = tools.FileCreationResult
import ToolCall from './ToolCall.vue'

export default defineComponent({
	name: 'BuiltinToolCallFileCreation',
	components: { ToolCall },
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
	},
})
</script>

<style scoped></style>
