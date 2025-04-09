<template>
	<v-row dense class="ma-0">
		<v-col cols="12" class="pa-0 ma-0">
			<v-input hide-details class="input">
				<template v-slot:prepend>
					<v-btn-toggle class="h-100">
						<v-btn @click="onMinMaximize" v-show="minimizable">
							<v-icon size="x-large" v-if="minimized">mdi-chevron-up</v-icon>
							<v-icon size="x-large" v-else>mdi-chevron-right</v-icon>
						</v-btn>
					</v-btn-toggle>
				</template>

				<slot></slot>

				<template v-slot:append>
					<v-btn-toggle class="h-100" v-show="showOptions">
						<slot name="option-buttons"></slot>
					</v-btn-toggle>
					<v-btn-toggle class="h-100">
						<v-btn @click="showOptions = !showOptions">
							<v-icon size="x-large" v-if="showOptions">mdi-chevron-down</v-icon>
							<v-icon size="x-large" v-else>mdi-dots-vertical</v-icon>
						</v-btn>
					</v-btn-toggle>
				</template>
			</v-input>
		</v-col>
	</v-row>

</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
	name: 'InputRow',
	emits: ['minMax'],
	props: {
		minimized: {
			type: Boolean,
			required: false,
			default: false,
		},
		minimizable: {
			type: Boolean,
			required: false,
			default: true,
		}
	},
	data(){
		return {
			showOptions: false
		}
	},
	methods: {
		onMinMaximize() {
			this.$emit('minMax')
		},
	}
})
</script>

<style scoped>
.input :deep(.v-input__prepend) {
	margin-right: 0;
}
.input :deep(.v-input__append) {
	margin-left: 0;
}
</style>
