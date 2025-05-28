<template>
	<ToolCall :tc="tc" icon="mdi-clock-time-one" without-title>
		<template v-slot:content>
			<Markdown :content="contentAsMarkdown" />
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import ToolCall from './ToolCall.vue'
import Markdown from '../Markdown.vue'

export default defineComponent({
	name: 'BuiltinToolCallSystemTime',
	components: { Markdown, ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		contentAsMarkdown(): string {
			return '```\n' + this.tc.Result?.Content + '\n```'
		}
	}
})
</script>

<style scoped></style>
