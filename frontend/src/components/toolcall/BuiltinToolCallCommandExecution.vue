<template>
	<ToolCall :tc="tc" icon="mdi-tools">
		<template v-slot:title>
			{{ title }}
		</template>

		<template v-slot:content v-if="tc.Result">
			<vue-markdown :source="resultAsMarkdown"></vue-markdown>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, command } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import CommandExecutionArguments = command.CommandExecutionArguments
import ToolCall from './ToolCall.vue'
import VueMarkdown from 'vue-markdown-render'

export default defineComponent({
	name: 'BuiltinToolCallCommandExecution',
	components: { ToolCall, VueMarkdown },
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
		title(){
			let line = this.parsedArguments.name
			if(this.parsedArguments.arguments) {
				line += " " + this.parsedArguments.arguments.join(' ')
			}
			return line
		},
		resultAsMarkdown() {
			if(!this.tc.Result) return null;

			let line = "```\n$> "
			if(this.parsedArguments.working_directory){
				line += "cd " + this.parsedArguments.working_directory + "\n"
			}

			if(this.parsedArguments.environment) {
				Object.entries(this.parsedArguments.environment).forEach(([key, value]) => {
					line += key + "=" + value + "\n"
				})
			}

			line += this.parsedArguments.name
			if(this.parsedArguments.arguments) {
				line += " " + this.parsedArguments.arguments.join(' ')
			}
			line += '\n' + this.tc.Result.Content.trim() + '\n```'

			return line
		}
	},
})
</script>

<style scoped></style>
