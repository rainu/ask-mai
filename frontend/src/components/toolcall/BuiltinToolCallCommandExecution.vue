<template>
	<ToolCall :tc="tc" icon="mdi-tools">
		<template v-slot:title>
			{{ title }}
		</template>

		<template v-slot:content v-if="resultAsMarkdown">
			<Markdown :content="resultAsMarkdown" />
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, command } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import CommandExecutionArguments = command.CommandExecutionArguments
import ToolCall from './ToolCall.vue'
import Markdown from '../Markdown.vue'

export default defineComponent({
	name: 'BuiltinToolCallCommandExecution',
	components: { Markdown, ToolCall },
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
		resultAsMarkdown(): string {
			if(!this.tc.Result) return '';

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
