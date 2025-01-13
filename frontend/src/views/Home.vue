<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />
		<UserScrollDetector @onScroll="onUserScroll" />

		<!-- app start state -->
		<template v-if="!(chatHistory.length > 0 || outputStream[0].Content || error)">
			<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<ChatInput v-model="input" :progress="progress" @submit="onSubmit" @interrupt="onInterrupt" />
				</div>
			</v-app-bar>

			<!-- this will trigger an redraw if attachments are added or removed -->
			<div :style="{ height: `${appbarHeight}px` }">{{ input }}</div>
		</template>

		<!-- after first prompt -->
		<template v-else>
			<!-- header -->
			<template v-if="$appConfig.UI.Prompt.PinTop">
				<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
					<div style="width: 100%" ref="appbar">
						<ChatInput v-model="input" :progress="progress" @submit="onSubmit" @interrupt="onInterrupt" />
					</div>
				</v-app-bar>

				<!-- div which is exactly high as the app-bar which is behind the appbar -->
				<div :style="{ height: `${appbarHeight}px` }">{{ input }}</div>
			</template>

			<!-- message section -->
			<template v-for="(entry, index) in chatHistory" :key="index">
				<ChatMessage :message="entry.ContentParts" :role="entry.Role" />
			</template>
			<template v-if="outputStream[0].Content">
				<ChatMessage :message="outputStream" :role="outputStreamRole" />
			</template>

			<v-alert v-if="error" type="error" :title="error.title" :text="error.message" />

			<!-- footer -->
			<template v-if="!$appConfig.UI.Prompt.PinTop">
				<v-footer app class="pa-0 ma-0" density="compact" height="auto">
					<div style="width: 100%" ref="appbar">
						<ChatInput v-model="input" :progress="progress" @submit="onSubmit" @interrupt="onInterrupt" />
					</div>
				</v-footer>
			</template>
		</template>
	</div>
</template>

<script lang="ts">
import {
	AppMounted,
	LLMAsk,
	LLMInterrupt,
	LLMWait,
} from '../../wailsjs/go/controller/Controller'
import { EventsOn, WindowGetSize, WindowSetPosition, WindowSetSize } from '../../wailsjs/runtime'
import ChatMessage, { ContentType, Role } from '../components/ChatMessage.vue'
import ChatInput, { ChatInputType } from '../components/ChatInput.vue'
import ZoomDetector from '../components/ZoomDetector.vue'
import UserScrollDetector from '../components/UserScrollDetector.vue'
import { controller } from '../../wailsjs/go/models.ts'
import LLMAskArgs = controller.LLMAskArgs
import LLMMessageContentPart = controller.LLMMessageContentPart
import LLMMessage = controller.LLMMessage

export default {
	name: 'Home',
	components: { UserScrollDetector, ZoomDetector, ChatInput, ChatMessage },
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
			chatHistory: [] as controller.LLMMessage[],
			userScroll: false,
			zoom: this.$appConfig.UI.Window.InitialZoom.Value,
		}
	},
	methods: {
		onZoom(factor: number) {
			this.zoom = factor
			this.adjustHeight()
		},
		async adjustHeight() {
			this.appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0

			const currentSize = await WindowGetSize()
			const pageHeight = (this.$refs.page as HTMLElement).clientHeight
			const combinedHeight = Math.ceil(pageHeight * this.zoom)
			const heightDiff = Math.min(combinedHeight, this.$appConfig.UI.Window.MaxHeight.Value) - currentSize.h

			await WindowSetSize(currentSize.w, combinedHeight)

			if (this.$appConfig.UI.Window.GrowTop && heightDiff > 0) {
				// move the window
				const offset = Math.min(combinedHeight, this.$appConfig.UI.Window.MaxHeight.Value)
				await WindowSetPosition(
					this.$appConfig.UI.Window.InitialPositionX.Value,
					this.$appConfig.UI.Window.InitialPositionY.Value - offset,
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
			})
		},
		async processLLM(input: ChatInputType, processFn: () => Promise<string>) {
			try {
				this.progress = true
				this.userScroll = false

				const setInput = () => {
					this.chatHistory.push(this.convertChatInputToLLMMessage(input))

					this.input.prompt = ''
					this.input.attachments = []
				}
				if (this.$appConfig.UI.Stream) {
					setInput()
				}

				const output = await processFn()

				if (!this.$appConfig.UI.Stream) {
					setInput()
				}

				this.chatHistory.push(
					LLMMessage.createFrom({
						Role: Role.Bot,
						ContentParts: [LLMMessageContentPart.createFrom({ Type: ContentType.Text, Content: output })],
					}),
				)
			} catch (err) {
				this.error = {
					title: 'Error while processing LLM',
					message: `${err}`,
				}
				console.error(err)
			} finally {
				this.progress = false
				this.outputStream[0].Content = ''
			}
		},
		async onSubmit(input: ChatInputType) {
			const args = LLMAskArgs.createFrom({
				History: [...this.chatHistory, this.convertChatInputToLLMMessage(input)] as LLMMessage[],
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
	},
	mounted() {
		EventsOn('llm:stream:chunk', (chunk: string) => {
			this.outputStream[0].Content += chunk
		})

		AppMounted()
			.then(() => {
				if (this.$appConfig.UI.Prompt.InitValue) {
					this.waitForLLM()
				}
			})
			.then(() => this.adjustHeight())
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
