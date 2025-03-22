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
			<vue-markdown :source="contentAsMarkdown"></vue-markdown>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import StatsArguments = tools.StatsArguments
import StatsResult = tools.StatsResult
import ToolCall from './ToolCall.vue'
import VueMarkdown from 'vue-markdown-render'

export default defineComponent({
	name: 'BuiltinToolCallStats',
	components: { ToolCall, VueMarkdown },
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
					return JSON.parse(this.tc.Result.Content) as StatsResult
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
