<template>
	<GeneralBar :minimizable="false">
		<!-- text-field is only here for height -->
		<v-text-field hide-details disabled></v-text-field>

		<v-row no-gutters>
			<v-col cols="6">
				<v-btn block variant="tonal" class="h-100" color="warning" @click="onNavigateBack">
					<template v-slot:prepend>
						<v-icon>mdi-cancel</v-icon>
					</template>
					{{ $t('edit.cancel') }}
				</v-btn>
			</v-col>
			<v-col cols="6">
				<v-btn block variant="tonal" class="h-100" color="primary" @click="$emit('onSave')">
					<template v-slot:prepend>
						<v-icon>mdi-content-save-check</v-icon>
					</template>
					{{ $t('edit.apply') }}
				</v-btn>
			</v-col>
		</v-row>
	</GeneralBar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import GeneralBar from './GeneralBar.vue'

export default defineComponent({
	name: 'EditBar',
	components: { GeneralBar },
	emits: ['onSave'],
	methods: {
		onNavigateBack() {
			this.$router.push({ name: 'Chat' })
		},
	},
})
</script>

<style scoped></style>