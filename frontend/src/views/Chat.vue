<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />
		<UserScrollDetector @onScroll="onUserScroll" />

		<!-- app start state -->
		<template v-if="!(chatHistory.length > 0 || outputStream[0].Content || error)">
			<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<ChatBar
						v-model="input"
						:progress="progress"
						:show-clear="chatHistory.length > 0 && !progress"
						:show-visibility-mode="chatHistory.length > 0 && !progress"
						:minimized="minimized"
						@submit="onSubmit"
						@interrupt="onInterrupt"
						@clear="onClear"
						@min-max="onMinMax"
					/>
				</div>
			</v-app-bar>

			<!-- this will trigger an redraw if attachments are added or removed -->
			<div :style="{ height: `${appbarHeight}px` }">{{ input }}</div>
		</template>

		<!-- after first prompt -->
		<template v-else>
			<!-- header -->
			<template v-if="config.UI.Prompt.PinTop">
				<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
					<div style="width: 100%" ref="appbar">
						<ChatBar
							v-model="input"
							:progress="progress"
							:show-clear="chatHistory.length > 0 && !progress"
							:show-visibility-mode="chatHistory.length > 0 && !progress"
							:minimized="minimized"
							@submit="onSubmit"
							@interrupt="onInterrupt"
							@clear="onClear"
							@min-max="onMinMax"
						/>
					</div>
				</v-app-bar>

				<!-- div which is exactly high as the app-bar which is behind the appbar -->
				<div :style="{ height: `${appbarHeight}px` }">{{ input }}</div>
			</template>

			<!-- message section -->
			<div v-show="!minimized">
				<template v-for="(e, index) in chatHistoryEntries" :key="index">
					<template v-if="e.Date">
						<v-sheet class="text-center">
							{{ e.Date }}
						</v-sheet>
					</template>

					<v-row v-if="e.Entry" no-gutters dense>
						<v-col cols="12" :class="e.Entry.Hidden ? 'opacity-30' : '' ">
							<ChatMessage
								:message="e.Entry.Message.ContentParts"
								:role="e.Entry.Message.Role"
								:date="e.Entry.Message.Created"
								@toggle-visibility="onToggle(e.Entry)"
								@on-edit="onEdit(e.EntryIdx)"
							/>
						</v-col>
					</v-row>
				</template>
				<template v-if="outputStream[0].Content">
					<ChatMessage :message="outputStream" :role="outputStreamRole" hide-edit />
				</template>

				<v-alert v-if="error" type="error" :title="error.title" :text="error.message" />
			</div>

			<!-- footer -->
			<template v-if="!config.UI.Prompt.PinTop">
				<v-footer app class="pa-0 ma-0" density="compact" height="auto">
					<div style="width: 100%" ref="appbar">
						<ChatBar
							v-model="input"
							:progress="progress"
							:show-clear="chatHistory.length > 0 && !progress"
							:show-visibility-mode="chatHistory.length > 0 && !progress"
							:minimized="minimized"
							@submit="onSubmit"
							@interrupt="onInterrupt"
							@clear="onClear"
							@min-max="onMinMax"
						/>
					</div>
				</v-footer>
			</template>
		</template>
	</div>
</template>

<script lang="ts">
import {
	AppMounted,
	GetLastState,
	LLMAsk,
	LLMInterrupt,
	LLMWait,
	Restart,
} from '../../wailsjs/go/controller/Controller'
import { EventsOn, WindowGetSize, WindowSetPosition, WindowSetSize } from '../../wailsjs/runtime'
import ChatMessage, { ContentType, Role } from '../components/ChatMessage.vue'
import ChatBar, { ChatInputType } from '../components/bar/ChatBar.vue'
import ZoomDetector from '../components/ZoomDetector.vue'
import UserScrollDetector from '../components/UserScrollDetector.vue'
import { controller, model } from '../../wailsjs/go/models.ts'
import LLMAskArgs = controller.LLMAskArgs
import LLMMessageContentPart = controller.LLMMessageContentPart
import LLMMessage = controller.LLMMessage
import { mapActions, mapState } from 'pinia'
import { HistoryEntry, useHistoryStore } from '../store/history.ts'
import { useConfigStore } from '../store/config.ts'

type State = {
	input: ChatInputType
	chatHistory: HistoryEntry[]
	config: model.Config
}

type HistoryEntryOrDate = {
	Date?: string
	Entry?: HistoryEntry
	EntryIdx?: number
}

export default {
	name: 'Chat',
	components: { UserScrollDetector, ZoomDetector, ChatBar, ChatMessage },
	data() {
		return {
			appbarHeight: 0,
			progress: false,
			input: {
				prompt: this.$appConfig.UI.Prompt.InitValue,
				attachments: this.$appConfig.UI.Prompt.InitAttachments,
			} as ChatInputType,
			outputStream: [
				{
					Type: ContentType.Text,
					Content: '',
				},
			] as LLMMessageContentPart[],
			outputStreamRole: Role.Bot,
			error: null as { title: string; message: string } | null,
			userScroll: false,
			minimized: false,
			zoom: this.$appConfig.UI.Window.InitialZoom.Value,
		}
	},
	computed: {
		...mapState(useConfigStore, ['config']),
		...mapState(useHistoryStore, ['chatHistory']),
		purgedChatHistory(): controller.LLMMessage[] {
			return this.chatHistory
				.filter((entry) => !entry.Interrupted)
				.filter((entry) => !entry.Hidden)
				.map((entry) => entry.Message)
		},
		chatHistoryEntries(): HistoryEntryOrDate[] {
			let result = [] as HistoryEntryOrDate[]
			const today = new Date().toLocaleDateString()
			const alreadyAddedDates = new Set<string>()

			for (let index = 0; index < this.chatHistory.length; index++) {
				const historyEntry = this.chatHistory[index]

				if(historyEntry.Message.Created) {
					const date = new Date(historyEntry.Message.Created * 1000).toLocaleDateString()

					if (!alreadyAddedDates.has(date)) {
						if(!(date === today && result.length === 0)) {
							// do not add date if it is today and the first entry
							result.push({ Date: date })
						}
					}
					alreadyAddedDates.add(date)
				}

				result.push({ Entry: historyEntry, EntryIdx: index })
			}

			return result
		}
	},
	methods: {
		...mapActions(useHistoryStore, ['setHistory', 'pushHistory', 'updateHistoryMessage', 'clearHistory']),
		...mapActions(useConfigStore, ['setConfig']),
		onZoom(factor: number) {
			this.zoom = factor
			this.adjustHeight()
		},
		async adjustHeight() {
			this.appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0

			const currentSize = await WindowGetSize()
			const pageHeight = (this.$refs.page as HTMLElement).clientHeight

			//the titlebar can not be manipulated while application lifecycle - so here we use the "initial" config
			const titleBarHeight = this.$appConfig.UI.Window.ShowTitleBar ? this.$appConfig.UI.Window.TitleBarHeight : 0
			const combinedHeight = Math.ceil(pageHeight * this.zoom) + titleBarHeight
			const heightDiff = Math.min(combinedHeight, this.config.UI.Window.MaxHeight.Value) - currentSize.h
			const width = this.config.UI.Window.InitialWidth.Value

			await WindowSetSize(width, combinedHeight)

			if (this.config.UI.Window.GrowTop && heightDiff > 0) {
				// move the window
				const offset = Math.min(combinedHeight, this.config.UI.Window.MaxHeight.Value)
				await WindowSetPosition(
					this.config.UI.Window.InitialPositionX.Value,
					this.config.UI.Window.InitialPositionY.Value - offset,
				)
			}
		},
		scrollToBottom() {
			if (this.userScroll) {
				// do not automatic scroll if the user is scrolling!
				return
			}

			window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' })
		},
		onUserScroll() {
			this.userScroll = true
		},
		async onMinMax() {
			this.minimized = !this.minimized
			await this.adjustHeight()
		},
		async onClear() {
			this.clearHistory()
			this.error = null
			await this.adjustHeight()
		},
		convertChatInputToLLMMessage(input: ChatInputType): LLMMessage {
			return LLMMessage.createFrom({
				Role: Role.User,
				ContentParts: [
					LLMMessageContentPart.createFrom({ Type: ContentType.Text, Content: input.prompt }),
					...input.attachments.map((a) =>
						LLMMessageContentPart.createFrom({
							Type: ContentType.Attachment,
							Content: a,
						}),
					),
				],
				Created: Math.floor(new Date().getTime() / 1000),
			})
		},
		async processLLM(input: ChatInputType, processFn: () => Promise<string>) {
			try {
				this.progress = true
				this.error = null
				this.userScroll = false

				const setInput = () => {
					this.pushHistory({
						Interrupted: false,
						Hidden: false,
						Message: this.convertChatInputToLLMMessage(input),
					})

					this.input.prompt = ''
					this.input.attachments = []
				}
				if (this.config.UI.Stream) {
					setInput()
				}

				const output = await processFn()

				if (!this.config.UI.Stream) {
					setInput()
				}

				this.pushHistory({
					Interrupted: false,
					Hidden: false,
					Message: LLMMessage.createFrom({
						Role: Role.Bot,
						ContentParts: [LLMMessageContentPart.createFrom({ Type: ContentType.Text, Content: output })],
						Created: Math.floor(new Date().getTime() / 1000),
					}),
				})
			} catch (err) {
				console.error(err)
				this.error = {
					title: 'Error while processing LLM',
					message: `${err}`,
				}

				if (this.config.UI.Stream) {
					// mark last input as "interrupted"
					this.chatHistory[this.chatHistory.length - 1].Interrupted = true

					this.pushHistory({
						Interrupted: true,
						Hidden: false,
						Message: LLMMessage.createFrom({
							Role: Role.Bot,
							ContentParts: this.outputStream,
							Created: Math.floor(new Date().getTime() / 1000),
						}),
					})
				}
			} finally {
				this.progress = false
				this.outputStream[0].Content = ''
			}
		},
		async onSubmit(input: ChatInputType) {
			const args = LLMAskArgs.createFrom({
				History: [...this.purgedChatHistory, this.convertChatInputToLLMMessage(input)] as LLMMessage[],
			})
			await this.processLLM(input, () => LLMAsk(args))
		},
		async waitForLLM() {
			this.input.prompt = this.$appConfig.UI.Prompt.InitValue
			this.input.attachments = this.$appConfig.UI.Prompt.InitAttachments
			await this.processLLM(this.input, () => LLMWait())
		},
		async onInterrupt() {
			await LLMInterrupt()
		},
		onToggle(entry: HistoryEntry) {
			entry.Hidden = !entry.Hidden
		},
		onEdit(idx: number | undefined) {
			if(idx === undefined) return

			this.$router.push({ name: 'Edit', params: { idx } })
		}
	},
	mounted() {
		EventsOn('llm:stream:chunk', (chunk: string) => {
			this.outputStream[0].Content += chunk
		})
		EventsOn('llm:message:add', (message: LLMMessage) => {
			this.outputStream[0].Content = ''
			this.pushHistory({
				Interrupted: false,
				Hidden: false,
				Message: message
			})
		})
		EventsOn('llm:message:update', (message: LLMMessage) => {
			this.updateHistoryMessage(message)
		})
		EventsOn('system:restart', () => {
			// backend requested a restart
			// so we have to save the current state and restart the app
			// but we have to wait until the progress is done (if any)
			const restartAfterProgress = () => {
				if (this.progress) {
					setTimeout(restartAfterProgress, 50)
				} else {
					const state = {
						input: this.input,
						chatHistory: this.chatHistory,
						config: this.config,
					} as State

					Restart(JSON.stringify(state))
				}
			}
			restartAfterProgress()
		})

		GetLastState().then((stateAsString) => {
			if (stateAsString) {
				const state = JSON.parse(stateAsString) as State
				this.input = state.input
				this.setHistory(state.chatHistory)
				this.setConfig(state.config)
			}

			AppMounted()
				.then(() => {
					if (this.$appConfig.UI.Prompt.InitValue && !stateAsString) {
						this.waitForLLM()
					}
				})
				.then(() => this.adjustHeight())
		})
	},
	updated() {
		this.$nextTick(() => {
			this.adjustHeight()

			setTimeout(this.scrollToBottom, 250)
		})
	},
}
</script>

<style>
.prompt .v-input__append {
	margin-left: 0;
	margin-right: 0;
}
</style>
