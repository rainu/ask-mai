<template>
	<ToolCall :tc="tc" icon="mdi-brain">
		<template v-slot:title>
			{{ tc.Meta.ToolName }}
		</template>

		<template v-slot:content>
			<v-row no-gutters>
				<v-col cols="12">
					<v-text-field :model-value="tc.Meta.ToolDescription"
												:label="$t('toolCall.mcp.description')"
												readonly hide-details density="compact"></v-text-field>
				</v-col>
				<v-col cols="12">
					<Markdown :content="textAsMarkdown(tc.Arguments)" />
				</v-col>

				<template v-if="content" >
					<v-col cols="12" v-for="c in content">
						<template v-if="c.type === 'text'">
							<Markdown :content="textAsMarkdown(c.text ?? '')" />
						</template>
						<template v-else-if="c.type === 'image'">
							<img :src="c.image" alt="Image" />
						</template>
					</v-col>
				</template>
			</v-row>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import ToolCall from './ToolCall.vue'
import { controller } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import Markdown from '../Markdown.vue'
import { ToolCallResult } from './types.ts'

type Content = {
	type: string
	text?: string
	image?: string
}

export default defineComponent({
	name: 'McpToolCall',
	components: { Markdown, ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	methods: {
		textAsMarkdown(text: string): string {
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
					const parsed = JSON.parse(this.tc.Result.Content) as ToolCallResult
					return parsed.content.map(c => {
						switch (c.type) {
							case 'text':
								return { type: 'text', text: c.text } as Content
							case 'image':
								return { type: 'image', image: c.image } as Content
							default:
								return { type: c.type } as Content // handle other types generically
						}
					})
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