<template>
	<v-row dense class="pa-0 ma-0">
		<v-col id="chat-input-text-area" :cols="progress ? 8 : 12" :sm="progress ? 11 : 12" class="pa-0 ma-0">
			<v-textarea
				v-model="value.prompt"
				:rows="rows"
				:disabled="progress"
				@keyup="onKeyup"
				:hide-details="value.attachments.length === 0"
				autofocus
				:placeholder="$t('prompt.placeholder')"
			>
				<template v-slot:prepend-inner>
					<v-btn icon density="compact" to="/history">
						<v-icon>mdi-history</v-icon>
					</v-btn>
					<v-btn icon density="compact" @click="onAddFile">
						<v-icon>mdi-paperclip</v-icon>
					</v-btn>
				</template>
				<template v-slot:append-inner>
					<v-btn icon v-show="!progress && isSubmitable" @click="onSubmit">
						<v-icon>mdi-send</v-icon>
					</v-btn>
				</template>
				<template v-slot:details>
					<v-container class="pa-0 ma-0" style="overflow-y: auto;">
						<v-chip
							v-for="(attachment, i) in value.attachments"
							:key="attachment"
							:title="attachment"
							class="ma-1"
							prepend-icon="mdi-file"
							color="primary"
							variant="flat"
							@click="onRemoveFile(i)"
						>
							{{ shortFileName(attachment) }}
						</v-chip>
					</v-container>
				</template>
			</v-textarea>
		</v-col>
		<v-col :cols="progress ? 4 : 0" :sm="progress ? 1 : 0" v-show="progress" class="pa-0 ma-0">
			<v-btn @click="onStop" variant="flat" color="error" block style="height: 100%">
				<v-icon size="x-large">mdi-stop-circle</v-icon>
			</v-btn>
		</v-col>
	</v-row>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { OpenFileDialog } from '../../wailsjs/go/controller/Controller'
import { controller } from '../../wailsjs/go/models.ts'
import { PathSeparator } from '../common/platform.ts'
import OpenFileDialogArgs = controller.OpenFileDialogArgs

export type ChatInputType = { prompt: string; attachments: string[] }

export default defineComponent({
	name: 'ChatInput',
	emits: ['update:modelValue', 'submit', 'interrupt'],
	props: {
		progress: {
			type: Boolean,
			required: false,
			default: false,
		},
		modelValue: {
			type: Object as () => ChatInputType,
			required: false,
			default: () => ({ prompt: '', attachments: [] }),
		},
	},
	computed: {
		value: {
			get() {
				return this.modelValue
			},
			set(value: string) {
				this.$emit('update:modelValue', value)
			},
		},
		isSubmitable() {
			return this.modelValue.prompt.trim() !== ''
		},
		rows() {
			return Math.max(
				Math.min(this.modelValue.prompt.split('\n').length, this.$appConfig.UI.Prompt.MaxRows),
				this.$appConfig.UI.Prompt.MinRows,
			)
		},
	},
	methods: {
		onKeyup(event: KeyboardEvent) {
			for (let i = 0; i < this.$appConfig.UI.Prompt.SubmitShortcut.Code.length; i++) {
				const code = event.code.toLowerCase() === this.$appConfig.UI.Prompt.SubmitShortcut.Code[i].toLowerCase()
				const ctrl = event.ctrlKey === this.$appConfig.UI.Prompt.SubmitShortcut.Ctrl[i]
				const shift = event.shiftKey === this.$appConfig.UI.Prompt.SubmitShortcut.Shift[i]
				const alt = event.altKey === this.$appConfig.UI.Prompt.SubmitShortcut.Alt[i]
				const meta = event.metaKey === this.$appConfig.UI.Prompt.SubmitShortcut.Meta[i]

				if (code && ctrl && shift && alt && meta) {
					this.onSubmit()
				}
			}
		},
		onSubmit() {
			if (!this.isSubmitable) return

			this.$emit('submit', this.modelValue)
		},
		onStop() {
			this.$emit('interrupt')
		},
		onAddFile() {
			OpenFileDialog(
				OpenFileDialogArgs.createFrom({
					Title: this.$t('dialog.files.title'),
				}),
			).then((results) => {
				this.value.attachments.push(...results)
			})
		},
		onRemoveFile(index: number) {
			this.value.attachments.splice(index, 1)
		},
		shortFileName(path: string) {
			let name = path.split(PathSeparator).pop() || ''
			if (name.length > 10) {
				name = name.slice(0, 10) + '...'
			}
			return name
		},
	},
})
</script>

<style>
#chat-input-text-area .v-messages {
	display: none;
}
</style>
