<template>
	<v-btn size="x-small" variant="tonal" @click="onToggleVisibility">
		<v-icon v-if="visibility">mdi-eye-outline</v-icon>
		<v-icon v-else>mdi-eye-closed</v-icon>
	</v-btn>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
	name: 'ChatMessageActions',
	emits: ['toggleVisibility'],
	data(){
		return {
			visibility: true,
		}
	},
	methods: {
		onToggleVisibility() {
			this.visibility = !this.visibility
			this.$emit('toggleVisibility')
		},
	}
})
</script>

<style scoped>

</style>