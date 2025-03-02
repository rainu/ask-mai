<template>
	<template v-if="isUserMessage">
		<v-row class="justify-end pa-2 mb-0 mt-1 mx-1 ml-15">
			<v-sheet color="green-accent-2" class="pa-2" rounded>
				<vue-markdown :source="textMessage" :options="options" />

				<!-- at the moment, only user messages can have attachments -->
				<template v-for="attachmentMeta of attachmentsMeta" :key="attachmentMeta.Path">
					<template v-if="isImage(attachmentMeta)">
						<div class="text-end">
							<img :src="attachmentMeta.Url" :alt="attachmentMeta.Url" :style="{ 'max-width': `${imageWidth}px` }" />
						</div>
					</template>
					<template v-else>
						<v-chip prepend-icon="mdi-file" color="primary" variant="flat">
							{{ fileName(attachmentMeta) }}
						</v-chip>
					</template>
				</template>
			</v-sheet>
		</v-row>
	</template>
	<template v-else-if="isToolMessage">
		<v-row class="pa-2 mb-0 mt-1 mx-1 mr-15" v-for="tc of toolCalls" :key="tc.Id">
			<v-sheet color="grey-lighten-2" rounded>
				<v-expansion-panels variant="accordion" bg-color="grey-lighten-2">
					<v-expansion-panel>
						<v-expansion-panel-title disable-icon-rotate class="pa-2">
							<v-icon icon="mdi-function"></v-icon>
							<span class="mr-2" v-html="highlight(`${tc.Function}(${tc.Arguments})`)"></span>
							<template v-slot:actions>
								<template v-if="tc.Result">
									<v-icon color="error" icon="mdi-alert-circle" v-if="tc.Result.Error"></v-icon>
									<v-icon color="success" icon="mdi-check-circle" v-else></v-icon>
								</template>
							</template>
						</v-expansion-panel-title>
						<v-expansion-panel-text v-if="tc.Result">
							<pre>{{ tc.Result.Content }}</pre>
							<v-alert type="error" v-if="tc.Result.Error" density="compact">{{ tc.Result.Error }}</v-alert>
						</v-expansion-panel-text>

						<v-progress-linear indeterminate size="small" v-if="!tc.NeedsApproval"></v-progress-linear>

						<template v-if="tc.NeedsApproval && !tc.Result">
							<v-row dense>
								<v-col cols="6" class="pr-0">
									<v-btn block color="success" @click="setToolCallApproval(tc, true)">
										<v-icon icon="mdi-check"></v-icon>
									</v-btn>
								</v-col>
								<v-col cols="6" class="pl-0">
									<v-btn block color="error" @click="setToolCallApproval(tc, false)">
										<v-icon icon="mdi-close"></v-icon>
									</v-btn>
								</v-col>
							</v-row>
						</template>
					</v-expansion-panel>
				</v-expansion-panels>
			</v-sheet>
		</v-row>
	</template>
	<template v-else>
		<v-row class="pa-2 mb-0 mt-1 mx-1 mr-15">
			<v-sheet color="grey-lighten-2" class="pa-2" rounded>
				<vue-markdown :source="textMessage" :options="options" />
			</v-sheet>
		</v-row>
	</template>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import VueMarkdown from 'vue-markdown-render'

import { GetAssetMeta, LLMApproveToolCall, LLMRejectToolCall } from '../../wailsjs/go/controller/Controller'
import { controller } from '../../wailsjs/go/models.ts'
import { PathSeparator } from '../common/platform.ts'
import AssetMeta = controller.AssetMeta
import LLMMessageContentPart = controller.LLMMessageContentPart
import LLMMessageCall = controller.LLMMessageCall
import hljs from 'highlight.js'
import { type Options as MarkdownItOptions } from 'markdown-it'
import { UseCodeStyle } from './code-style.ts'
import { ClipboardSetText } from '../../wailsjs/runtime'

export enum Role {
	User = 'human',
	Bot = 'ai',
	Tool = 'tool',
}

export enum ContentType {
	Text = 'text',
	Attachment = 'attachment',
	ToolCall = 'tool',
}

export default defineComponent({
	name: 'ChatMessage',
	components: { VueMarkdown },
	props: {
		message: {
			type: Array as PropType<LLMMessageContentPart[]>,
			required: true,
		},
		role: {
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
							const result = hljs.highlight(code, { language })
							return result.value
						} catch (__) {}
					}
					return '' // use external default escaping
				},
			} as MarkdownItOptions,
			attachmentsMeta: [] as AssetMeta[],
		}
	},
	computed: {
		isUserMessage() {
			return this.role === Role.User
		},
		isToolMessage() {
			return this.role === Role.Tool
		},
		textMessage() {
			return this.message
				.filter((part) => part.Type === ContentType.Text)
				.map((part) => part.Content)
				.join('')
		},
		toolCalls(): LLMMessageCall[] {
			return this.message.filter((part) => part.Type === ContentType.ToolCall).map((part) => part.Call)
		},
		attachments() {
			return this.message.filter((part) => part.Type === ContentType.Attachment).map((part) => part.Content)
		},
		imageWidth() {
			return this.$appConfig.UI.Window.InitialWidth.Value * 0.9
		},
	},
	methods: {
		enrichCopyButtons() {
			const codeBlocks = document.querySelectorAll('pre:not(.code-container pre)').values()
			for (const pre of codeBlocks) {
				const div = document.createElement('div')
				div.className = 'code-container'

				const button = document.createElement('button')
				button.className = 'copy-button mdi-clipboard-text-outline mdi v-icon notranslate v-icon--size-small'
				button.addEventListener('click', this.onCopyButtonClicked)

				div.appendChild(button)
				div.appendChild(pre.cloneNode(true))
				pre.replaceWith(div)
			}
		},
		onCopyButtonClicked(event: MouseEvent) {
			const preElement = (event.target as HTMLButtonElement)?.nextElementSibling as HTMLElement
			if (preElement && preElement.tagName === 'PRE') {
				let code = preElement.innerText
				const newlineCount = (code.match(/\n/g) || []).length
				if (newlineCount === 1) {
					//prevent that shell-code-statements will be executed directly
					code = code.trim()
				}

				ClipboardSetText(code).then(() => {
					preElement.classList.add('copied')
					setTimeout(() => {
						preElement.classList.remove('copied')
					}, 1000)
				})
			}
		},
		fileName(asset: AssetMeta) {
			return asset.Path.split(PathSeparator).pop() || ''
		},
		isImage(asset: AssetMeta) {
			return asset.MimeType.startsWith('image/')
		},
		highlight(code: string) {
			return hljs.highlight(code, { language: 'JavaScript' }).value
		},
		setToolCallApproval(call: LLMMessageCall, approved: boolean) {
			call.NeedsApproval = false
			if (approved) {
				LLMApproveToolCall(call.Id)
			} else {
				LLMRejectToolCall(call.Id)
			}
		},
	},
	watch: {
		textMessage() {
			this.$nextTick(() => this.enrichCopyButtons())
		},
		attachments: {
			async handler() {
				const promises = this.attachments.map((path) => GetAssetMeta(path))
				try {
					this.attachmentsMeta = await Promise.all(promises)
				} catch (e) {
					console.error(e)
				}
			},
			immediate: true,
		},
	},
	mounted() {
		UseCodeStyle(this.$appConfig.UI.CodeStyle)

		this.$nextTick(() => this.enrichCopyButtons())
	},
})
</script>

<style>
pre code {
	background-color: #f5f5f5;
	border: 1px solid #ccc;
	margin: 0.5em 0;
	padding: 0.5em 2em 0.5em 1em;
	border-radius: 5px;
	display: block;
	overflow-x: auto;
	position: relative;
}

/* blue block when hover over code-blocks */
pre code::before {
	content: '';
	position: absolute;
	top: 0;
	left: 0;
	width: 0;
	height: 100%;
	background-color: #007bff;
	transition: width 0.3s;
}

pre code:hover::before {
	width: 5px;
}

/* Add styles for the copy button */

div.code-container {
	position: relative;
}

.copy-button {
	position: absolute;
	right: 1px;
	padding-top: 1em;
	padding-right: 0.5em;
	z-index: 1;
	float: right;
}

pre.copied code {
	border: 2px solid #4db6ac;
	border-radius: 5px;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.5); /* Add elevation effect */
}

/* Inline code blocks inside text (not in headers) */
code:not(pre code):not(h1 code):not(h2 code):not(h3 code):not(h4 code):not(h5 code):not(h6 code) {
	background-color: #f5f5f5;
	border: 1px solid #ccc;
	padding: 2px 4px;
	border-radius: 3px;
}

/* Quote-Blocks */

blockquote {
	border-left: 5px solid #ccc;
	margin: 0.5em 0;
	padding: 0.5em 1em;
	color: #555;
	background: none;
	border-radius: 0;
}

blockquote:hover {
	border-color: #007bff;
}

/* Header Padding increase */
h1,
h2,
h3,
h4,
h5,
h6 {
	padding-top: 1em;
}

/* make list items visible */
ol,
ul li {
	margin-left: 2em;
}
</style>
