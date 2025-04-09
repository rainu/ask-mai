<template>
	<InputRow :minimized="minimized">
		<template v-slot:prepend>
			<v-btn icon density="compact" to="/history">
				<v-icon>mdi-history</v-icon>
			</v-btn>
			<v-btn icon density="compact" @click="onAddFile">
				<v-icon>mdi-paperclip</v-icon>
			</v-btn>
		</template>

		<v-textarea
			v-model="value.prompt"
			:rows="rows"
			:disabled="progress"
			@keyup="onKeyup"
			:hide-details="value.attachments.length === 0"
			autofocus
			:placeholder="$t('prompt.placeholder')"
		>
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

		<template v-slot:options>
			<v-btn-toggle class="h-100">
				<v-btn active-color="primary" @click="toggleVisibilityMode" :active="visibilityMode">
					<v-icon size="x-large">mdi-eye-settings-outline</v-icon>
				</v-btn>
				<v-btn @click="onClear">
					<v-icon size="x-large">mdi-chat-remove-outline</v-icon>
				</v-btn>
			</v-btn-toggle>
		</template>

		<template v-slot:append >
			<v-btn block variant="flat" class="h-100" v-show="!progress && isSubmitable" @click="onSubmit">
				<v-icon size="x-large">mdi-send</v-icon>
			</v-btn>
			<v-btn block variant="flat" class="h-100" color="error" v-show="progress" @click="onStop">
				<v-icon size="x-large">mdi-stop-circle</v-icon>
			</v-btn>
		</template>
	</InputRow>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { OpenFileDialog } from '../../wailsjs/go/controller/Controller'
import { controller } from '../../wailsjs/go/models.ts'
import { PathSeparator } from '../common/platform.ts'
import OpenFileDialogArgs = controller.OpenFileDialogArgs
import InputRow from './InputRow.vue'

export type ChatInputType = { prompt: string; attachments: string[] }

export default defineComponent({
	name: 'ChatInput',
	components: { InputRow },
	emits: ['update:modelValue', 'submit', 'interrupt', 'clear', 'changeVisibilityMode'],
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
		showClear: {
			type: Boolean,
			required: false,
			default: false,
		},
		minimized: {
			type: Boolean,
			required: false,
			default: false,
		},
	},
	data(){
		return {
			visibilityMode: false,
		}
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
		onClear() {
			this.$emit('clear')
		},
		toggleVisibilityMode(){
			this.visibilityMode = !this.visibilityMode
			this.$emit('changeVisibilityMode', this.visibilityMode)
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
