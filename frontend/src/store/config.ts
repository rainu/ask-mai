import { defineStore } from 'pinia'
import { model } from '../../wailsjs/go/models.ts'
import {
	SetActiveProfile, SetBuiltinTools, SetMcpTools
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
			await this.applyToolsConfig()
		},
		async setActiveProfile(profileName: string) {
			this.config = await SetActiveProfile(profileName)
			this.activeProfileName = profileName
		},
		async applyToolsConfig() {
			await SetBuiltinTools(this.config.LLM.Tools.BuiltInTools)
			await SetMcpTools(this.config.LLM.McpServer)
		}
	}
})