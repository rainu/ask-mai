<template>
	<v-row dense class="ma-0">
		<v-col cols="12" class="pa-0 ma-0">
			<v-input hide-details class="input">
				<template v-slot:prepend>
					<v-btn-toggle class="h-100">
						<v-btn @click="onMinMaximize" v-show="minimizable">
							<v-icon size="x-large" v-if="minimized">mdi-chevron-right</v-icon>
							<v-icon size="x-large" v-else>mdi-chevron-down</v-icon>
						</v-btn>
					</v-btn-toggle>
				</template>

				<slot></slot>

				<template v-slot:append>
					<slot name="append"></slot>

					<v-btn-toggle class="h-100" v-show="showOptions">
						<slot name="option-buttons"></slot>
						<v-btn @click="onNavigateProfile" v-if="availableProfilesCount > 1">
							<v-icon size="x-large">mdi-application-cog-outline</v-icon>
						</v-btn>
						<v-btn @click="onNavigateTool">
							<v-icon size="x-large">mdi-tools</v-icon>
						</v-btn>
					</v-btn-toggle>
					<v-btn-toggle class="h-100">
						<v-btn @click="toggleOptions">
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
import { mapActions, mapState } from 'pinia'
import { useGlobalStore } from '../../store/global.ts'
import { GetAvailableProfiles } from '../../../wailsjs/go/controller/Controller'

export default defineComponent({
	name: 'GeneralBar',
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
			availableProfilesCount: 0
		}
	},
	computed: {
		...mapState(useGlobalStore, ['showOptions'])
	},
	methods: {
		...mapActions(useGlobalStore, ['toggleOptions']),
		onNavigateProfile(){
			if(this.$route.name === 'Profile'){
				this.$router.back()
			} else {
				this.$router.push({ name: 'Profile' })
			}
		},
		onNavigateTool(){
			if(this.$route.name === 'Tool'){
				this.$router.back()
			} else {
				this.$router.push({ name: 'Tool' })
			}
		},
		onMinMaximize() {
			this.$emit('minMax')
		},
	},
	created() {
		GetAvailableProfiles().then(profiles => {
			this.availableProfilesCount = Object.entries(profiles).length
		})
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
