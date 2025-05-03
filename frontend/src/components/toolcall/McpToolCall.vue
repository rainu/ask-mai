<template>
	<ToolCall :tc="tc" icon="mdi-brain">
		<template v-slot:title>
			{{ tc.McpToolName }}
		</template>

		<template v-slot:content>
			<v-row no-gutters>
				<v-text-field :model-value="tc.McpToolDescription"
											:label="$t('toolCall.mcp.description')"
											readonly hide-details density="compact"></v-text-field>
			</v-row>
			<v-row no-gutters v-if="content" >
				<v-col cols="12" v-for="c in content">
					<template v-if="c.type === 'text'">
						<vue-markdown :source="textAsMarkdown(c.text)" />
					</template>
					<template v-else-if="c.type === 'image'">
						<img :src="c.image" alt="Image" />
					</template>
				</v-col>
			</v-row>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import ToolCall from './ToolCall.vue'
import VueMarkdown from 'vue-markdown-render'
import { controller } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall

type Content = {
	type: string
	text?: string
	image?: string
}

export default defineComponent({
	name: 'McpToolCall',
	components: { ToolCall, VueMarkdown },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	methods: {
		textAsMarkdown(text: string) {
			// try to prettify JSON
			try {
				text = JSON.stringify(JSON.parse(text), null, 2)
			} catch (e) {}

			return '```\n' + text + '\n```'
		}
	},
	computed: {
		content(): null | Content[] {
			if(this.tc.Result) {
				try {
					const parsed = JSON.parse(this.tc.Result.Content)
					return parsed.content
				} catch (e) {
					console.warn("Failed to parse JSON", e)
				}
			}
			return null
		}
	}
})
</script>

<style scoped></style>