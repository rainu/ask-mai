<template>
	<ToolCall :tc="tc" icon="mdi-magnify" without-title>
		<template v-slot:title>
			<template v-if="parsedResult">
				{{ parsedResult.path }}
			</template>
			<template v-else>
				{{ parsedArguments.path }}
			</template>
		</template>
		<template v-slot:content>
			<Markdown :content="contentAsMarkdown" />
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, file, mcp } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import StatsArguments = file.StatsArguments
import StatsResult = file.StatsResult
import ToolCall from './ToolCall.vue'
import Markdown from '../Markdown.vue'
import { ToolCallResult } from './types.ts'

export default defineComponent({
	name: 'BuiltinToolCallStats',
	components: { Markdown, ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): StatsArguments {
			return JSON.parse(this.tc.Arguments) as StatsArguments
		},
		parsedResult(): StatsResult | null {
			if(this.tc.Result) {
				try {
					const tcr = JSON.parse(this.tc.Result.Content) as ToolCallResult
					const firstTextContent = tcr.content.find(c => c.type === 'text') as mcp.TextContent
					if(firstTextContent) {
						return JSON.parse(firstTextContent.text) as StatsResult
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
