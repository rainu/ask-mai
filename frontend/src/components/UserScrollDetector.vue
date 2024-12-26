<template></template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
	name: 'UserScrollDetector',
	emits: ['onScroll'],
	data() {
		return {
			ctrl: false,
		}
	},
	methods: {
		handleKeyup(event: KeyboardEvent) {
			if (event.key === 'Control') {
				this.ctrl = false
			}
		},
		handleKeydown(event: KeyboardEvent) {
			if (event.key === 'Control') {
				this.ctrl = true
			} else if (event.key === 'PageUp' || event.key === 'PageDown') {
				this.onUserScroll()
			}
		},
		handleWheel(event: WheelEvent) {
			if (!this.ctrl && event.deltaY !== 0) {
				this.onUserScroll()
			}
		},
		onUserScroll() {
			this.$emit('onScroll')
		},
	},
	created() {
		window.addEventListener('keyup', this.handleKeyup)
		window.addEventListener('keydown', this.handleKeydown)
		window.addEventListener('wheel', this.handleWheel)
	},
})
</script>

<style scoped></style>
