<template>
	<div ref="page">
		<template v-if="output">
			<v-app-bar app class="pa-0 ma-0" density="compact">
				<div style="width: 100%" ref="appbar">
					<v-row dense class="pa-0 ma-0">
						<v-col :cols="progress ? 11 : 12" class="pa-0 ma-0">
							<v-text-field
								v-model="input"
								:disabled="progress"
								@keyup.enter="onSubmit"
								hide-details
								autofocus
								:placeholder="$t('prompt.placeholder')"
							>
							</v-text-field>
						</v-col>
						<v-col :cols="progress ? 1 : 0" v-show="progress" class="pa-0 ma-0">
							<v-btn @click="onStop" variant="flat" color="error" block style="height: 100%">
								<v-icon size="x-large">mdi-stop-circle</v-icon>
							</v-btn>
						</v-col>
					</v-row>
				</div>
			</v-app-bar>
			<ChatMessage :message="output" />
		</template>

		<template v-else>
			<v-row dense class="pa-0 ma-0">
				<v-col :cols="progress ? 11 : 12" class="pa-0 ma-0">
					<v-text-field
						v-model="input"
						:disabled="progress"
						@keyup.enter="onSubmit"
						hide-details
						autofocus
						:placeholder="$t('prompt.placeholder')"
					>
					</v-text-field>
				</v-col>
				<v-col :cols="progress ? 1 : 0" v-show="progress" class="pa-0 ma-0">
					<v-btn @click="onStop" variant="flat" color="error" block style="height: 100%">
						<v-icon size="x-large">mdi-stop-circle</v-icon>
					</v-btn>
				</v-col>
			</v-row>
		</template>
	</div>
</template>

<script lang="ts">
import { AppMounted, GetInitialPrompt, LLMAsk, LLMInterrupt } from '../../wailsjs/go/controller/Controller'
import { WindowGetSize, WindowSetSize } from '../../wailsjs/runtime'
import ChatMessage from '../components/ChatMessage.vue'

export default {
	name: 'Home',
	components: { ChatMessage },
	data() {
		return {
			progress: false,
			input: '' as string,
			output: '' as string,
		}
	},
	methods: {
		async adjustHeight() {
			const currentSize = await WindowGetSize()
			const pageHeight = (this.$refs.page as HTMLElement).clientHeight
			const appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0

			await WindowSetSize(currentSize.w, pageHeight + appbarHeight)
		},
		async onSubmit() {
			try {
				this.progress = true
				this.output = await LLMAsk({ Content: this.input })
			} catch (err) {
				console.error(err)
			} finally {
				this.progress = false
			}
		},
		async onStop() {
			await LLMInterrupt()
		},
	},
	mounted() {
		GetInitialPrompt().then((prompt) => {
			if (prompt) {
				this.input = prompt
				this.onSubmit()
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
