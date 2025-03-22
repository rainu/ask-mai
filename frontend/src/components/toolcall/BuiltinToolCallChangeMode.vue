<template>
	<ToolCall :tc="tc" icon="mdi-shield-edit" without-title>
		<template v-slot:title>
			{{ parsedArguments.permission }} - {{ parsedArguments.path }}
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import ChangeModeArguments = tools.ChangeModeArguments
import ToolCall from './ToolCall.vue'

export default defineComponent({
	name: 'BuiltinToolCallChangeMode',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): ChangeModeArguments {
			return JSON.parse(this.tc.Arguments) as ChangeModeArguments
		}
	}
})
</script>

<style scoped></style>
