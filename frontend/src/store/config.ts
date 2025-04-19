import { defineStore } from 'pinia'
import { model } from '../../wailsjs/go/models.ts'
import {
	SetActiveProfile, SetBuiltinTools
} from '../../wailsjs/go/controller/Controller'

export const useConfigStore = defineStore('config', {
	state: () => ({
		config: window.$appConfig as model.Config,
		activeProfileName: window.$appConfig.Profile.Active,
	}),
	actions: {
		async setConfig(config: model.Config) {
			this.config = config
			window.$appConfig = config
			await this.applyBuiltinToolsConfig()
		},
		async setActiveProfile(profileName: string) {
			this.config = await SetActiveProfile(profileName)
			this.activeProfileName = profileName
		},
		async applyBuiltinToolsConfig() {
			await SetBuiltinTools(this.config.LLM.Tools.BuiltInTools)
		}
	}
})