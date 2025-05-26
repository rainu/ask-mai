<template>
	<ToolCall :tc="tc" icon="mdi-information-outline" without-title>
		<template v-slot:content>
			<vue-markdown :source="contentAsMarkdown"></vue-markdown>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, system } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import EnvironmentResult = system.EnvironmentResult
import ToolCall from './ToolCall.vue'
import VueMarkdown from 'vue-markdown-render'

export default defineComponent({
	name: 'BuiltinToolCallEnvironment',
	components: { ToolCall, VueMarkdown },
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
					return JSON.parse(this.tc.Result.Content) as EnvironmentResult
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
		contentAsMarkdown(){
			if(this.parsedResult) {
				return '```\n' + JSON.stringify(this.parsedResult, null, 3) + '\n```'
			}
			return null
		}
	}
})
</script>

<style scoped></style>
