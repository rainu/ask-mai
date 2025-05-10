import { defineStore } from 'pinia'
import { model } from '../../wailsjs/go/models.ts'
import {
	SetActiveProfile, SetBuiltinTools, SetMcpTools
} from '../../wailsjs/go/controller/Controller'

export const useConfigStore = defineStore('config', {
	state: () => ({
		profile: window.$appProfile as model.Profile,
		activeProfileName: window.$appConfig.ActiveProfile,
	}),
	actions: {
		async setConfig(profile: model.Profile) {
			this.profile = profile
			window.$appProfile = profile
			await this.applyToolsConfig()
		},
		async setActiveProfile(profileName: string) {
			this.profile = await SetActiveProfile(profileName)
			this.activeProfileName = profileName
		},
		async applyToolsConfig() {
			await SetBuiltinTools(this.profile.LLM.Tools.BuiltInTools)
			await SetMcpTools(this.profile.LLM.McpServer)
		}
	}
})