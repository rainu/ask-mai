<template>
	<ToolCall :tc="tc" icon="mdi-function">
		<template v-slot:title>
			<span v-html="functionHeader"></span>
		</template>

		<template v-slot:content v-if="tc.Result">
			<pre>{{ tc.Result.Content }}</pre>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import hljs from 'highlight.js'
import ToolCall from './ToolCall.vue'

export default defineComponent({
	name: 'GeneralToolCall',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		functionHeader() {
			let code = this.tc.Function + "(" + this.tc.Arguments + ")"
			return this.highlight(code)
		}
	},
	methods: {
		highlight(code: string) {
			return hljs.highlight(code, { language: 'JavaScript' }).value
		},
	}
})
</script>

<style scoped></style>
