<template>
	<div ref="page">
		<template v-if="chatHistory.length > 0">
			<v-app-bar app class="pa-0 ma-0" density="compact">
				<div style="width: 100%" ref="appbar">
					<ChatInput v-model="input" :progress="progress" @submit="onSubmit" @interrupt="onInterrupt" />
				</div>
			</v-app-bar>

      <template v-for="(entry, index) in chatHistory" :key="index">
        <ChatMessage :message="entry.message" :role="entry.role" />
      </template>
		</template>

		<template v-else>
			<ChatInput v-model="input" :progress="progress" @submit="onSubmit" @interrupt="onInterrupt" />
		</template>
	</div>
</template>

<script lang="ts">
import { AppMounted, GetInitialPrompt, LLMAsk, LLMInterrupt } from '../../wailsjs/go/controller/Controller'
import { WindowGetSize, WindowSetSize } from '../../wailsjs/runtime'
import ChatMessage, { Role } from '../components/ChatMessage.vue'
import ChatInput from '../components/ChatInput.vue'

type ChatHistoryEntry = {
	message: string
	role: Role
}

export default {
	name: 'Home',
	components: { ChatInput, ChatMessage },
	data() {
		return {
			progress: false,
      input: '',
			chatHistory: [] as ChatHistoryEntry[],
		}
	},
	methods: {
		async adjustHeight() {
			const currentSize = await WindowGetSize()
			const pageHeight = (this.$refs.page as HTMLElement).clientHeight
			const appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0

			await WindowSetSize(currentSize.w, pageHeight + appbarHeight)
		},
		async onSubmit(input: string) {
			try {
				this.progress = true
				const output = await LLMAsk({ Content: input })
        this.chatHistory.push(
          { message: input, role: Role.User },
          { message: output, role: Role.Bot },
        )
        this.input = ""

				console.log(output)
			} catch (err) {
				console.error(err)
			} finally {
				this.progress = false
			}
		},
		async onInterrupt() {
			await LLMInterrupt()
		},
	},
	mounted() {
		GetInitialPrompt().then((prompt) => {
			if (prompt) {
        this.input = prompt
				this.onSubmit(prompt)
			}
		})
		this.adjustHeight().then(() => AppMounted())
	},
	updated() {
		this.$nextTick(() => this.adjustHeight())
	},
}
</script>

<style>
.prompt .v-input__append {
	margin-left: 0;
	margin-right: 0;
}
</style>
