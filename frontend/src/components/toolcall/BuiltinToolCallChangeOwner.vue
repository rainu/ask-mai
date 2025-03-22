<template>
	<ToolCall :tc="tc" icon="mdi-account-switch" without-title>
		<template v-slot:title>
			{{ parsedArguments.user_id }}:{{ parsedArguments.group_id }} -
			{{ parsedArguments.path }}
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import ChangeOwnerArguments = tools.ChangeOwnerArguments
import ToolCall from './ToolCall.vue'

export default defineComponent({
	name: 'BuiltinToolCallChangeOwner',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): ChangeOwnerArguments {
			return JSON.parse(this.tc.Arguments) as ChangeOwnerArguments
		}
	}
})
</script>

<style scoped></style>
