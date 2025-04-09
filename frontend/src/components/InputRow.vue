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
					<v-row dense class="h-100">
						<v-col class="ma-0 pa-0">
							<slot name="append"></slot>
						</v-col>
						<v-col class="ma-0 pa-0" v-if="showOptions">
							<slot name="options"></slot>
						</v-col>
						<v-col class="ma-0 pa-0">
							<v-btn block variant="flat" class="h-100" @click="showOptions = !showOptions">
								<v-icon size="x-large" v-if="showOptions">mdi-chevron-right</v-icon>
								<v-icon size="x-large" v-else>mdi-chevron-left</v-icon>
							</v-btn>
						</v-col>
					</v-row>
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
