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

				<small class="d-flex justify-space-between align-center">
					<span class="opacity-50 pr-2">{{ createdAt }}</span>
					<ChatMessageActions reverse @toggleVisibility="onToggleVisibility" :hide-edit="hideEdit" @on-edit="onEdit" />
				</small>
			</v-sheet>
		</v-row>
	</template>
	<template v-else-if="isToolMessage">
		<v-row class="pa-2 mb-0 mt-1 mx-1 mr-15" v-for="tc of toolCalls" :key="tc.Id">
			<v-sheet color="grey-lighten-2" rounded>
				<BuiltinToolCallSystemInfo :tc="tc" v-if="tc.Meta.BuiltIn && tc.Function.endsWith('getSystemInformation')" />
				<BuiltinToolCallEnvironment :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('getEnvironment')" />
				<BuiltinToolCallSystemTime :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('getSystemTime')" />

				<BuiltinToolCallChangeMode :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('changeMode')" />
				<BuiltinToolCallChangeOwner :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('changeOwner')" />
				<BuiltinToolCallChangeTimes :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('changeTimes')" />

				<BuiltinToolCallStats :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('getStats')" />

				<BuiltinToolCallFileCreation
					:tc="tc"
					v-else-if="tc.Meta.BuiltIn && (tc.Function.endsWith('createFile') || tc.Function.endsWith('createTempFile'))"
				/>
				<BuiltinToolCallFileAppending :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('appendFile')" />
				<BuiltinToolCallFileReading :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('readTextFile')" />
				<BuiltinToolCallFileDeletion :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('deleteFile')" />

				<BuiltinToolCallDirectoryCreation
					:tc="tc"
					v-else-if="
						tc.Meta.BuiltIn && (tc.Function.endsWith('createDirectory') || tc.Function.endsWith('createTempDirectory'))
					"
				/>
				<BuiltinToolCallDirectoryDeletion :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('deleteDirectory')" />

				<BuiltinToolCallCommandExecution :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('executeCommand')" />

				<BuiltinToolCallHttpCall :tc="tc" v-else-if="tc.Meta.BuiltIn && tc.Function.endsWith('callHttp')" />

				<McpToolCall :tc="tc" v-else-if="tc.Meta.Mcp" />

				<GeneralToolCall :tc="tc" v-else />

				<small class="d-flex justify-space-between px-2 pb-2 align-center">
					<ChatMessageActions @toggleVisibility="onToggleVisibility" hide-edit />
					<span class="opacity-50 pl-2">{{ createdAt }}</span>
				</small>
			</v-sheet>
		</v-row>
	</template>
	<template v-else-if="isSystemMessage">
		<v-row class="pa-2 mb-0 mt-1 mx-1" v-if="textMessage">
			<v-col>
				<v-sheet class="pa-2" rounded>
					<vue-markdown :source="textMessage" :options="options" />
					<ChatMessageActions @toggleVisibility="onToggleVisibility" :hide-edit="hideEdit" @onEdit="onEdit" />
				</v-sheet>
			</v-col>
		</v-row>
	</template>
	<template v-else>
		<v-row class="pa-2 mb-0 mt-1 mx-1 mr-15">
			<v-sheet color="grey-lighten-2" class="pa-2" rounded>
				<vue-markdown :source="textMessage" :options="options" />
				<small class="d-flex justify-space-between align-center">
					<ChatMessageActions @toggleVisibility="onToggleVisibility" :hide-edit="hideEdit" @on-edit="onEdit" />
					<template v-if="consumption">
						<v-btn size="x-small" variant="tonal" @click="showConsumption = !showConsumption">
							<v-icon>mdi-chart-pie</v-icon>
						</v-btn>
					</template>
					<span class="opacity-50">{{ createdAt }}</span>
				</small>
				<template v-if="consumption">
					<Consumption :model="consumption" class="mt-2" v-show="showConsumption"/>
				</template>
			</v-sheet>
		</v-row>
	</template>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import VueMarkdown from 'vue-markdown-render'

import { GetAssetMeta } from '../../wailsjs/go/controller/Controller'
import { controller } from '../../wailsjs/go/models.ts'
import { PathSeparator } from '../common/platform.ts'
import AssetMeta = controller.AssetMeta
import LLMMessageContentPart = controller.LLMMessageContentPart
import LLMMessageCall = controller.LLMMessageCall
import hljs from 'highlight.js'
import { type Options as MarkdownItOptions } from 'markdown-it'
import { UseCodeStyle } from './code-style.ts'
import { ClipboardSetText } from '../../wailsjs/runtime'
import GeneralToolCall from './toolcall/GeneralToolCall.vue'
import BuiltinToolCallFileCreation from './toolcall/BuiltinToolCallFileCreation.vue'
import BuiltinToolCallCommandExecution from './toolcall/BuiltinToolCallCommandExecution.vue'
import BuiltinToolCallFileReading from './toolcall/BuiltinToolCallFileReading.vue'
import BuiltinToolCallFileDeletion from './toolcall/BuiltinToolCallFileDeletion.vue'
import BuiltinToolCallFileAppending from './toolcall/BuiltinToolCallFileAppending.vue'
import BuiltinToolCallSystemInfo from './toolcall/BuiltinToolCallSystemInfo.vue'
import BuiltinToolCallSystemTime from './toolcall/BuiltinToolCallSystemTime.vue'
import BuiltinToolCallDirectoryDeletion from './toolcall/BuiltinToolCallDirectoryDeletion.vue'
import BuiltinToolCallDirectoryCreation from './toolcall/BuiltinToolCallDirectoryCreation.vue'
import BuiltinToolCallStats from './toolcall/BuiltinToolCallStats.vue'
import BuiltinToolCallEnvironment from './toolcall/BuiltinToolCallEnvironment.vue'
import BuiltinToolCallChangeMode from './toolcall/BuiltinToolCallChangeMode.vue'
import BuiltinToolCallChangeOwner from './toolcall/BuiltinToolCallChangeOwner.vue'
import BuiltinToolCallChangeTimes from './toolcall/BuiltinToolCallChangeTimes.vue'
import BuiltinToolCallHttpCall from './toolcall/BuiltinToolCallHttpCall.vue'
import ChatMessageActions from './ChatMessageActions.vue'
import { mapState } from 'pinia'
import { useConfigStore } from '../store/config.ts'
import McpToolCall from './toolcall/McpToolCall.vue'
import { HistoryEntryConsumption } from '../store/history.ts'
import Consumption from './Consumption.vue'

export enum Role {
	System = 'system',
	Bot = 'ai',
	Tool = 'tool',
	User = 'human',
}

export enum ContentType {
	Text = 'text',
	Attachment = 'attachment',
	ToolCall = 'tool',
}

export default defineComponent({
	name: 'ChatMessage',
	components: {
		Consumption,
		McpToolCall,
		ChatMessageActions,
		BuiltinToolCallHttpCall,
		BuiltinToolCallChangeTimes,
		BuiltinToolCallChangeOwner,
		BuiltinToolCallChangeMode,
		BuiltinToolCallEnvironment,
		BuiltinToolCallStats,
		BuiltinToolCallDirectoryCreation,
		BuiltinToolCallDirectoryDeletion,
		BuiltinToolCallSystemTime,
		BuiltinToolCallSystemInfo,
		BuiltinToolCallFileAppending,
		BuiltinToolCallFileDeletion,
		BuiltinToolCallFileReading,
		GeneralToolCall,
		BuiltinToolCallCommandExecution,
		BuiltinToolCallFileCreation,
		VueMarkdown,
	},
	emits: ['toggleVisibility', 'onEdit'],
	props: {
		message: {
			type: Array as PropType<LLMMessageContentPart[]>,
			required: true,
		},
		consumption: {
			type: Object as () => HistoryEntryConsumption,
			required: false,
		},
		role: {
			type: String,
			required: true,
		},
		date: {
			type: Number,
			required: false,
		},
		hideEdit: {
			type: Boolean,
			required: false,
			default: false,
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
			showConsumption: false,
		}
	},
	computed: {
		...mapState(useConfigStore, ['profile']),
		isUserMessage() {
			return this.role === Role.User
		},
		isToolMessage() {
			return this.role === Role.Tool
		},
		isSystemMessage() {
			return this.role === Role.System
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
			return (this.profile.UI.Window.InitialWidth.Value ?? 0) * 0.9
		},
		createdAt() {
			if (!this.date) return null

			const d = new Date(this.date * 1000)
			return d.toLocaleTimeString()
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
		onEdit() {
			this.$emit('onEdit')
		},
		onToggleVisibility() {
			this.$emit('toggleVisibility')
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
		UseCodeStyle(this.profile.UI.CodeStyle)

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
