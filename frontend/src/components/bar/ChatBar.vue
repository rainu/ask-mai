<template>
	<GeneralBar :minimized="minimized">
		<v-textarea
			v-model="value.prompt"
			:rows="rows"
			:disabled="progress"
			@keyup="onKeyup"
			:hide-details="!hasAttachments"
			autofocus
			:placeholder="$t('prompt.placeholder')"
		>
			<template v-slot:prepend-inner>
				<v-btn icon density="compact" @click="onAddFile">
					<v-icon>mdi-paperclip</v-icon>
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

			<template v-slot:append-inner>
				<v-btn v-show="!progress && isSubmitable" @click="onSubmit">
					<v-icon size="x-large">mdi-send</v-icon>
				</v-btn>
			</template>
		</v-textarea>

		<template v-slot:append>
			<v-btn color="error" class="h-100" v-show="progress" @click="onStop">
				<v-icon size="x-large">mdi-stop-circle-outline</v-icon>
			</v-btn>
		</template>

		<template v-slot:option-buttons>
			<v-btn @click="onClear" v-show="showClear">
				<v-icon size="x-large">mdi-chat-remove-outline</v-icon>
			</v-btn>
			<v-btn to="/history">
				<v-icon size="x-large">mdi-history</v-icon>
			</v-btn>
		</template>
	</GeneralBar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { OpenFileDialog } from '../../../wailsjs/go/controller/Controller'
import { controller } from '../../../wailsjs/go/models.ts'
import { PathSeparator } from '../../common/platform.ts'
import OpenFileDialogArgs = controller.OpenFileDialogArgs
import GeneralBar from './GeneralBar.vue'
import { mapState } from 'pinia'
import { useConfigStore } from '../../store/config.ts'

export type ChatInputType = { prompt: string; attachments: string[] | null }

export default defineComponent({
	name: 'ChatBar',
	components: { GeneralBar },
	emits: ['update:modelValue', 'submit', 'interrupt', 'clear'],
	props: {
		progress: {
			type: Boolean,
			required: false,
			default: false,
		},
		modelValue: {
			type: Object as () => ChatInputType,
			required: false,
			default: () => ({ prompt: '' }),
		},
		showClear: {
			type: Boolean,
			required: false,
			default: false,
		},
		showVisibilityMode: {
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
	computed: {
		...mapState(useConfigStore, ['profile']),
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
				Math.min(this.modelValue.prompt.split('\n').length, (this.profile.UI.Prompt.MaxRows ?? 4)),
				(this.profile.UI.Prompt.MinRows ?? 1),
			)
		},
		hasAttachments(): boolean {
			return this.modelValue.attachments !== null && this.modelValue.attachments.length > 0
		}
	},
	methods: {
		onKeyup(event: KeyboardEvent) {
			for (let i = 0; i < this.profile.UI.Prompt.SubmitShortcut.Code.length; i++) {
				const code = event.code.toLowerCase() === this.profile.UI.Prompt.SubmitShortcut.Code[i].toLowerCase()
				const ctrl = event.ctrlKey === this.profile.UI.Prompt.SubmitShortcut.Ctrl[i]
				const shift = event.shiftKey === this.profile.UI.Prompt.SubmitShortcut.Shift[i]
				const alt = event.altKey === this.profile.UI.Prompt.SubmitShortcut.Alt[i]
				const meta = event.metaKey === this.profile.UI.Prompt.SubmitShortcut.Meta[i]

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
		onAddFile() {
			OpenFileDialog(
				OpenFileDialogArgs.createFrom({
					Title: this.$t('dialog.files.title'),
				}),
			).then((results) => {
				if(this.value.attachments === null) {
					this.value.attachments = []
				}
				this.value.attachments.push(...results)
			})
		},
		onRemoveFile(index: number) {
			if(this.value.attachments === null) {
				return
			}
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
