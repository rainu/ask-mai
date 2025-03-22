<template>
	<ToolCall :tc="tc" icon="mdi-clock-time-one" without-title>
		<template v-slot:content>
			<vue-markdown :source="contentAsMarkdown"></vue-markdown>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import ToolCall from './ToolCall.vue'
import VueMarkdown from 'vue-markdown-render'

export default defineComponent({
	name: 'BuiltinToolCallSystemTime',
	components: { ToolCall, VueMarkdown },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		contentAsMarkdown(){
			return '```\n' + this.tc.Result?.Content + '\n```'
		}
	}
})
</script>

<style scoped></style>
