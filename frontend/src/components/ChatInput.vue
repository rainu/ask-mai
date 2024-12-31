<template>
	<v-row dense class="pa-0 ma-0">
		<v-col :cols="progress ? 11 : 12" class="pa-0 ma-0">
			<v-textarea
				v-model="value"
				:rows="rows"
				:disabled="progress"
				@keyup="onKeyup"
				hide-details
				autofocus
				:placeholder="$t('prompt.placeholder')"
			>
				<template v-slot:append-inner>
					<v-btn icon v-show="!progress && isSubmitable" @click="onSubmit">
						<v-icon>mdi-send</v-icon>
					</v-btn>
				</template>
			</v-textarea>
		</v-col>
		<v-col :cols="progress ? 1 : 0" v-show="progress" class="pa-0 ma-0">
			<v-btn @click="onStop" variant="flat" color="error" block style="height: 100%">
				<v-icon size="x-large">mdi-stop-circle</v-icon>
			</v-btn>
		</v-col>
	</v-row>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

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
			type: String,
			required: false,
			default: '',
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
		isSubmitable(){
			return this.modelValue.trim() !== ''
		},
		rows() {
			return Math.max(
				Math.min(this.modelValue.split('\n').length, this.$appConfig.UI.Prompt.MaxRows),
				this.$appConfig.UI.Prompt.MinRows,
			)
		},
	},
	methods: {
		onKeyup(event: KeyboardEvent) {
			const code = event.code.toLowerCase() === this.$appConfig.UI.Prompt.SubmitShortcut.Code.toLowerCase()
			const ctrl = event.ctrlKey === this.$appConfig.UI.Prompt.SubmitShortcut.Ctrl
			const shift = event.shiftKey === this.$appConfig.UI.Prompt.SubmitShortcut.Shift
			const alt = event.altKey === this.$appConfig.UI.Prompt.SubmitShortcut.Alt
			const meta = event.metaKey === this.$appConfig.UI.Prompt.SubmitShortcut.Meta

			if (code && ctrl && shift && alt && meta) {
				this.onSubmit()
			}
		},
		onSubmit() {
			if (!this.isSubmitable) return

			this.$emit('submit', this.modelValue.trim())
		},
		onStop() {
			this.$emit('interrupt')
		},
	},
})
</script>

<style scoped></style>
