<template>
	<v-row dense class="ma-0">
		<v-col cols="12" class="pa-0 ma-0">
			<v-input hide-details class="input">
				<template v-slot:prepend>
					<v-btn icon density="compact" class="ml-1" @click="onMinMaximize">
						<v-icon v-if="minimized">mdi-chevron-up-box</v-icon>
						<v-icon v-else>mdi-chevron-down-box</v-icon>
					</v-btn>
					<slot name="prepend"></slot>
				</template>

				<slot></slot>

				<template v-slot:append>
					<slot name="append"></slot>
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
