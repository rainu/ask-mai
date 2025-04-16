<template>
	<div v-if="!reverse">
		<v-btn size="x-small" variant="tonal" @click="onToggleVisibility">
			<v-icon v-if="visibility">mdi-eye-outline</v-icon>
			<v-icon v-else>mdi-eye-closed</v-icon>
		</v-btn>
		<v-btn size="x-small" variant="tonal" @click="onEdit" v-if="!hideEdit">
			<v-icon>mdi-pencil</v-icon>
		</v-btn>
	</div>

	<div v-else>
		<v-btn size="x-small" variant="tonal" @click="onEdit" v-if="!hideEdit">
			<v-icon>mdi-pencil</v-icon>
		</v-btn>
		<v-btn size="x-small" variant="tonal" @click="onToggleVisibility">
			<v-icon v-if="visibility">mdi-eye-outline</v-icon>
			<v-icon v-else>mdi-eye-closed</v-icon>
		</v-btn>
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
	name: 'ChatMessageActions',
	emits: ['toggleVisibility', 'onEdit'],
	props: {
		reverse: {
			type: Boolean,
			required: false,
			default: false,
		},
		hideEdit: {
			type: Boolean,
			required: false,
			default: false,
		}
	},
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
		onEdit() {
			this.$emit('onEdit')
		}
	}
})
</script>

<style scoped>

</style>