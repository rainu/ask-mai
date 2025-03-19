<template>
	<ToolCall :tc="tc" icon="mdi-tools">
		<template v-slot:title>
			{{ parsedArguments.name }} {{ parsedArguments.arguments.join(' ') }}
		</template>

		<template v-slot:content v-if="tc.Result">
			<code>{{ code }}</code>
			<pre>{{ tc.Result.Content }}</pre>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import CommandExecutionArguments = tools.CommandExecutionArguments
import ToolCall from './ToolCall.vue'

export default defineComponent({
	name: 'BuiltinToolCallCommandExecution',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): CommandExecutionArguments {
			return JSON.parse(this.tc.Arguments) as CommandExecutionArguments
		},
		code(){
			let line = ""

			if(this.parsedArguments.working_directory){
				line += "cd " + this.parsedArguments.working_directory + ";\n"
			}

			if(this.parsedArguments.environment) {
				Object.entries(this.parsedArguments.environment).forEach(([key, value]) => {
					line += key + "=" + value + "\n"
				})
			}

			line += this.parsedArguments.name + " " + this.parsedArguments.arguments.join(' ')

			return line
		}
	},
})
</script>

<style scoped></style>
