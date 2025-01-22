<template>
	<div style="width: 100%" ref="bar">
		<v-row dense>
			<v-col :cols="minimized ? 12 : 1" class="pa-0 ma-0" v-if="isMinMaxLeft">
				<v-btn block class="h-100" @click="onMinMaximize" :style="{ minHeight: `${height}px` }">
					<v-icon size="x-large" v-if="minimized">mdi-chevron-right-box</v-icon>
					<v-icon size="x-large" v-else>mdi-chevron-left-box</v-icon>
				</v-btn>
			</v-col>
			<v-col :cols="isMinMaxEnabled ? 11 : 12" class="pa-0 ma-0" v-show="!minimized">
				<ChatInput v-model="value" :progress="progress" @submit="onSubmit" @interrupt="onInterrupt" />
			</v-col>
			<v-col :cols="minimized ? 12 : 1" class="pa-0 ma-0" v-if="isMinMaxRight">
				<v-btn block class="h-100" @click="onMinMaximize" :style="{ minHeight: `${height}px` }">
					<v-icon size="x-large" v-if="minimized">mdi-chevron-left-box</v-icon>
					<v-icon size="x-large" v-else>mdi-chevron-right-box</v-icon>
				</v-btn>
			</v-col>
		</v-row>
	</div>
</template>

<script lang="ts">
import ChatInput, { ChatInputType } from './ChatInput.vue'

export default {
	name: 'ChatInputBar',
	components: { ChatInput },
	emits: ['update:modelValue', 'submit', 'interrupt', 'minMax'],
	props: {
		modelValue: {
			type: Object as () => ChatInputType,
			required: false,
			default: () => ({ prompt: '', attachments: [] }),
		},
		progress: {
			type: Boolean,
			required: true,
		},
		minimized: {
			type: Boolean,
			required: false,
			default: false,
		},
	},
	data() {
		return {
			height: 0,
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
		isMinMaxLeft(){
			return this.$appConfig.UI.MinMaxPosition === 'left'
		},
		isMinMaxRight(){
			return this.$appConfig.UI.MinMaxPosition === 'right'
		},
		isMinMaxEnabled(){
			return this.isMinMaxLeft || this.isMinMaxRight
		}
	},
	methods: {
		onSubmit(input: string) {
			this.$emit('submit', input)
		},
		onInterrupt() {
			this.$emit('interrupt')
		},
		onMinMaximize() {
			this.$emit('minMax')
		},
	},
	mounted() {
		this.height = this.$refs.bar ? (this.$refs.bar as HTMLElement).clientHeight : 0
	},
}
</script>
