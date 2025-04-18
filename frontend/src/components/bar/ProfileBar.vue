<template>
	<GeneralBar :minimizable="false">
		<v-row no-gutters>
			<v-col cols="12">
				<v-btn block disabled class="h-100 text-none">
					{{ $t('profile.title') }}
				</v-btn>
			</v-col>
		</v-row>

		<!-- text-field is only here for height -->
		<v-text-field hide-details disabled></v-text-field>
	</GeneralBar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import GeneralBar from './GeneralBar.vue'
import { mapState } from 'pinia'
import { useConfigStore } from '../../store/config.ts'

export default defineComponent({
	name: 'ProfileBar',
	components: { GeneralBar },
	emits: ['onSave'],
	computed: {
		...mapState(useConfigStore, ['config']),
	},
	methods: {
		onNavigateBack() {
			this.$router.back()
		},
	},
})
</script>

<style scoped></style>