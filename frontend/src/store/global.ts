import { defineStore } from 'pinia'

export const useGlobalStore = defineStore('global', {
	state: () => ({
		showOptions: false,
	}),
	actions: {
		toggleOptions() {
			this.showOptions = !this.showOptions
		}
	}
})