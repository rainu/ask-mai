<template>
	<ToolCall :tc="tc" icon="mdi-information-outline" without-title>
		<template v-slot:content>
			<Markdown :content="contentAsMarkdown" />
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, mcp, system } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import EnvironmentResult = system.EnvironmentResult
import ToolCall from './ToolCall.vue'
import Markdown from '../Markdown.vue'
import { ToolCallResult } from './types.ts'

export default defineComponent({
	name: 'BuiltinToolCallEnvironment',
	components: { Markdown, ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedResult(): EnvironmentResult | null {
			if(this.tc.Result) {
				try {
					const tcr = JSON.parse(this.tc.Result.Content) as ToolCallResult
					const firstTextContent = tcr.content.find(c => c.type === 'text') as mcp.TextContent
					if(firstTextContent) {
						return JSON.parse(firstTextContent.text) as EnvironmentResult
					}
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
		contentAsMarkdown(): string {
			if(this.parsedResult) {
				return '```\n' + JSON.stringify(this.parsedResult, null, 3) + '\n```'
			}
			return ''
		}
	}
})
</script>

<style scoped></style>
