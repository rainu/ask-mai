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
import { controller, command, mcp } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import CommandExecutionArguments = command.CommandExecutionArguments
import ToolCall from './ToolCall.vue'
import Markdown from '../Markdown.vue'
import { ToolCallResult } from './types.ts'

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
		parsedResult(): string | null {
			if(this.tc.Result) {
				try {
					const tcr = JSON.parse(this.tc.Result.Content) as ToolCallResult
					const firstTextContent = tcr.content.find(c => c.type === 'text') as mcp.TextContent
					if(firstTextContent) {
						return firstTextContent.text
					}
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
		title(){
			let line = this.parsedArguments.name
			if(this.parsedArguments.arguments) {
				if (Array.isArray(this.parsedArguments.arguments)) {
					line += " " + this.parsedArguments.arguments.join(' ')
				} else {
					line += " " + this.parsedArguments.arguments
				}
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
				if (Array.isArray(this.parsedArguments.arguments)) {
					line += " " + this.parsedArguments.arguments.join(' ')
				} else {
					line += " " + this.parsedArguments.arguments
				}
			}

			if(this.parsedResult) {
				line += '\n' + this.parsedResult.trim() + '\n```'
			}

			return line
		}
	},
})
</script>

<style scoped></style>
