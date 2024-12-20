<template>
	<div ref="page" :style="{ zoom }">
		<template v-if="chatHistory.length > 0 || error">
			<v-app-bar app class="pa-0 ma-0" density="compact">
				<div style="width: 100%" ref="appbar">
					<ChatInput v-model="input" :progress="progress" @submit="onSubmit" @interrupt="onInterrupt" />
				</div>
			</v-app-bar>

			<template v-for="(entry, index) in chatHistory" :key="index">
				<ChatMessage :message="entry.Content" :role="entry.Role" />
			</template>

			<v-alert v-if="error" type="error" :title="error.title" :text="error.message" />
		</template>

		<template v-else>
			<ChatInput v-model="input" :progress="progress" @submit="onSubmit" @interrupt="onInterrupt" />
		</template>

		<a ref="bottom"></a>
	</div>
</template>

<script lang="ts">
import { AppMounted, LLMAsk, LLMInterrupt, LLMWait } from '../../wailsjs/go/controller/Controller'
import { WindowGetSize, WindowSetSize } from '../../wailsjs/runtime'
import ChatMessage, { Role } from '../components/ChatMessage.vue'
import ChatInput from '../components/ChatInput.vue'
import { controller } from '../../wailsjs/go/models.ts'
import LLMAskArgs = controller.LLMAskArgs

export default {
	name: 'Home',
	components: { ChatInput, ChatMessage },
	data() {
		return {
			progress: false,
			input: '',
			error: null as { title: string; message: string } | null,
			chatHistory: [] as controller.LLMMessage[],
			zoom: this.$appConfig.UI.Window.InitialZoom,
		}
	},
	methods: {
		zoomIn() {
			if (this.zoom < 5) {
				this.zoom += 0.1
			}
		},
		zoomOut() {
			if (this.zoom > 0.1) {
				this.zoom -= 0.1
			}
		},
		handleKeyup(event: KeyboardEvent) {
			if (event.ctrlKey && event.key === '+') {
				this.zoomIn()
				this.adjustHeight()
			}
			if (event.ctrlKey && event.key === '-') {
				this.zoomOut()
				this.adjustHeight()
			}
			if (event.ctrlKey && event.key === '0') {
				this.zoom = 1
				this.adjustHeight()
			}
		},
		async adjustHeight() {
			const currentSize = await WindowGetSize()
			const pageHeight = (this.$refs.page as HTMLElement).clientHeight
			const appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0
			const combinedHeight = Math.ceil((pageHeight + appbarHeight) * this.zoom)

			await WindowSetSize(currentSize.w, combinedHeight)
		},
		scrollToBottom() {
			const bottomEl = this.$refs.bottom as HTMLElement
			bottomEl.scrollIntoView({ block: 'end', behavior: 'smooth' })
		},
		async processLLM(input: string, processFn: () => Promise<string>) {
			try {
				this.progress = true
				const output = await processFn()
				this.chatHistory.push({ Content: input, Role: Role.User }, { Content: output, Role: Role.Bot })
				this.input = ''
			} catch (err) {
				this.error = {
					title: 'Error while processing LLM',
					message: `${err}`,
				}
				console.error(err)
			} finally {
				this.progress = false
			}
		},
		async onSubmit(input: string) {
			const args = LLMAskArgs.createFrom({
				History: [
					...this.chatHistory,
					{
						Content: input,
						Role: Role.User,
					},
				],
			})
			await this.processLLM(input, () => LLMAsk(args))
		},
		async waitForLLM() {
			this.input = this.$appConfig.UI.Prompt
			await this.processLLM(this.input, () => LLMWait())
		},
		async onInterrupt() {
			await LLMInterrupt()
		},
	},
	mounted() {
		window.addEventListener('keyup', this.handleKeyup)

		if (this.$appConfig.UI.Prompt) {
			this.waitForLLM()
		}
		this.adjustHeight().then(() => AppMounted())
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
