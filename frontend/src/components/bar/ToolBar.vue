<template>
	<GeneralBar :minimizable="false">
		<v-row no-gutters>
			<v-col cols="12">
				<v-btn block disabled class="h-100 text-none">
					{{ $t('tool.title') }}
				</v-btn>
			</v-col>
		</v-row>

		<!-- text-field is only here for height -->
		<v-text-field hide-details disabled></v-text-field>

		<template v-slot:option-buttons>
			<v-btn @click="onNavigateChat">
				<v-icon size="x-large">mdi-chat-processing-outline</v-icon>
			</v-btn>
		</template>
	</GeneralBar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import GeneralBar from './GeneralBar.vue'

export default defineComponent({
	name: 'ToolBar',
	components: { GeneralBar },
	methods: {
		onNavigateChat() {
			this.$router.push({ name: 'Chat' })
		},
	}
})
</script>

<style scoped></style>