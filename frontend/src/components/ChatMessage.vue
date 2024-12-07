<template>
	<v-card>
		<v-card-text>
			<vue-markdown :source="message" :options="options" />
		</v-card-text>
	</v-card>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import VueMarkdown from 'vue-markdown-render'

import 'highlight.js/styles/github.css'
import hljs from 'highlight.js'

export default defineComponent({
	name: 'ChatMessage',
	components: { VueMarkdown },
	props: {
		message: {
			type: String,
			required: true,
		},
	},
	data() {
		return {
			options: {
				highlight: (code: string, language: string) => {
					if (language && hljs.getLanguage(language)) {
						try {
							return hljs.highlight(code, { language }).value
						} catch (__) {}
					}
					return '' // use external default escaping
				},
			},
		}
	},
})
</script>

<style scoped></style>
