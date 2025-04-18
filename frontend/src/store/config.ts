import { defineStore } from 'pinia'
import { model } from '../../wailsjs/go/models.ts'
import {
	SetActiveProfile,
} from '../../wailsjs/go/controller/Controller'

export const useConfigStore = defineStore('config', {
	state: () => ({
		config: window.$appConfig as model.Config,
		activeProfileName: window.$appConfig.Profile.Active,
	}),
	actions: {
		async setActiveProfile(profileName: string) {
			this.config = await SetActiveProfile(profileName)
			this.activeProfileName = profileName
		}
	}
})