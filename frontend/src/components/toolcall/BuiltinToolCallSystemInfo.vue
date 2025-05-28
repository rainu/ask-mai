<template>
	<ToolCall :tc="tc" icon="mdi-information-outline" without-title>
		<template v-slot:content>
			<Markdown :content="contentAsMarkdown" />
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, system } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import SystemInfoResult = system.SystemInfoResult
import ToolCall from './ToolCall.vue'
import Markdown from '../Markdown.vue'

export default defineComponent({
	name: 'BuiltinToolCallSystemInfo',
	components: { Markdown, ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedResult(): SystemInfoResult | null {
			if(this.tc.Result) {
				try {
					return JSON.parse(this.tc.Result.Content) as SystemInfoResult
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
