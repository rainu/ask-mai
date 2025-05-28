<template>
	<ToolCall :tc="tc" icon="mdi-web">
		<template v-slot:title>
			<strong>{{ httpMethod }}&nbsp;</strong>
			<template v-if="httpMethod === 'GET'">
				<a :href="parsedArguments.url" target="_blank" rel="noopener noreferrer">
					{{ parsedArguments.url }}
				</a>
			</template>
			<template v-else>
				{{ parsedArguments.url }}
			</template>
		</template>

		<template v-slot:content v-if="parsedResult">
			<div class="request">
				<v-row dense v-for="(values, name) of parsedArguments.header" >
					<v-col cols="4" md="3">
						<strong>{{ name }}</strong>
					</v-col>
					<v-col cols="8" md="9">
						{{values}}
					</v-col>
				</v-row>

				<Markdown :content="requestBodyAsMarkdown" v-if="parsedArguments.body" />
			</div>

			<v-divider thickness="5" inset class="my-4"/>

			<div class="response">
				<strong>{{ parsedResult.status }}</strong>

				<v-row dense v-for="(values, name) of parsedResult.header" >
					<v-col cols="4" md="3">
						<strong>{{ name }}</strong>
					</v-col>
					<v-col cols="8" md="9">
						{{values.join(', ')}}
					</v-col>
				</v-row>

				<Markdown :content="responseBodyAsMarkdown" />
			</div>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, http } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import CallArguments = http.CallArguments
import CallResult = http.CallResult
import ToolCall from './ToolCall.vue'
import Markdown from '../Markdown.vue'

export default defineComponent({
	name: 'BuiltinToolCallHttpCall',
	components: { Markdown, ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): CallArguments {
			return JSON.parse(this.tc.Arguments) as CallArguments
		},
		parsedResult(): CallResult | null {
			if(this.tc.Result) {
				try {
					return JSON.parse(this.tc.Result.Content) as CallResult
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
		requestBodyAsMarkdown(): string {
			if(!this.tc.Arguments) return '';

			let line = '```\n'
			if(this.parsedArguments.body) {
				line += this.parsedArguments.body
			}
			line += '\n```'

			return line
		},
		responseBodyAsMarkdown(): string {
			if(!this.tc.Result) return '';

			let line = '```\n'
			line += this.parsedResult?.body
			line += '\n```'

			return line
		},
		httpMethod() {
			return this.parsedArguments.method ? this.parsedArguments.method.toUpperCase() : "GET"
		},
	},
})
</script>

<style scoped></style>
