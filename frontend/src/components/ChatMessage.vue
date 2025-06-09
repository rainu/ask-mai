<template>
	<template v-if="isUserMessage">
		<v-row class="justify-end pa-2 mb-0 mt-1 mx-1 ml-15">
			<v-sheet color="chat-user-message" class="pa-2" rounded>
				<Markdown :content="textMessage" />

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
			<v-sheet color="chat-tool-call" rounded>
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
				<v-sheet color="chat-system-message" class="pa-2" rounded>
					<Markdown :content="textMessage" />
					<ChatMessageActions @toggleVisibility="onToggleVisibility" :hide-edit="hideEdit" @onEdit="onEdit" />
				</v-sheet>
			</v-col>
		</v-row>
	</template>
	<template v-else>
		<v-row class="pa-2 mb-0 mt-1 mx-1 mr-15">
			<v-sheet color="chat-assistant-message" class="pa-2" rounded>
				<Markdown :content="textMessage" />
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

import { GetAssetMeta } from '../../wailsjs/go/controller/Controller'
import { controller } from '../../wailsjs/go/models.ts'
import { PathSeparator } from '../common/platform.ts'
import AssetMeta = controller.AssetMeta
import LLMMessageContentPart = controller.LLMMessageContentPart
import LLMMessageCall = controller.LLMMessageCall
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
import Markdown from './Markdown.vue'

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
		Markdown,
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
})
</script>

<style scoped>
</style>
