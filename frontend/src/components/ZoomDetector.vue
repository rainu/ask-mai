<template></template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
	name: 'ZoomDetector',
	emits: ['onZoom'],
	data() {
		return {
			zoom: this.$appProfile.UI.Window.InitialZoom.Value ?? 1,
			ctrl: false,
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
		handleKeydown(event: KeyboardEvent) {
			if (event.key === 'Control') {
				this.ctrl = true
			}
		},
		handleKeyup(event: KeyboardEvent) {
			if (event.key === 'Control') {
				this.ctrl = false
			} else if (event.ctrlKey && event.key === '+') {
				this.zoomIn()
			} else if (event.ctrlKey && event.key === '-') {
				this.zoomOut()
			} else if (event.ctrlKey && event.key === '0') {
				this.zoom = 1
			}
		},
		handleWheel(event: WheelEvent) {
			if (this.ctrl && event.deltaY > 0) {
				this.zoomOut()
			} else if (this.ctrl && event.deltaY < 0) {
				this.zoomIn()
			}
		},
	},
	watch: {
		zoom() {
			this.$emit('onZoom', this.zoom)
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
