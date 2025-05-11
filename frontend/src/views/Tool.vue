<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />

		<!-- header -->
		<template v-if="profile.UI.Prompt.PinTop">
			<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<ToolBar />
				</div>
			</v-app-bar>

			<!-- this will trigger an redraw if attachments are added or removed -->
			<div :style="{ height: `${appbarHeight}px` }"></div>
		</template>

		<v-card class="ma-2">
			<v-card-title>{{ $t('tool.builtin.title') }}</v-card-title>
			<v-card-text>
				<v-row no-gutters>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="SystemInfo" label="SystemInfo" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="Environment" label="Environment" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="SystemTime" label="SystemTime" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="Stats" label="Stats" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="ChangeMode" label="ChangeMode" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="ChangeOwner" label="ChangeOwner" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="ChangeTimes" label="ChangeTimes" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="FileCreation" label="FileCreation" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="FileTempCreation" label="FileTempCreation" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="FileAppending" label="FileAppending" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="FileReading" label="FileReading" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="FileDeletion" label="FileDeletion" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="DirectoryCreation" label="DirectoryCreation" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="DirectoryTempCreation" label="DirectoryTempCreation" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="DirectoryDeletion" label="DirectoryDeletion" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="CommandExec" label="CommandExec" density="compact" hide-details></v-checkbox></v-col>
					<v-col cols="6" sm="3" lg="2"><v-checkbox v-model="Http" label="Http" density="compact" hide-details></v-checkbox></v-col>
				</v-row>
			</v-card-text>
		</v-card>

		<v-card class="ma-2">
			<v-card-title>{{ $t('tool.mcp.title') }}</v-card-title>
			<v-card-text>
				<v-progress-circular indeterminate v-if="mcp.loading" />

				<v-row dense>
					<v-col cols="12" v-for="(server, name) in mcp.server" :key="name">
						<v-card>
							<v-card-title>
								<span>
									{{ name }}
									<v-tooltip activator="parent" location="right" v-if="server.config.Name">
										{{ server.config.Name }} {{ server.config.Arguments.join(' ') }}
									</v-tooltip>
									<v-tooltip activator="parent" location="right" v-else>
										{{ server.config.BaseUrl }}{{ server.config.Endpoint }}
									</v-tooltip>
								</span>
							</v-card-title>
							<v-card-text>
								<v-row no-gutters>
									<v-col cols="6" sm="4" lg="3" v-for="(tool, ti) in server.tools" :key="ti">
										<v-checkbox density="compact" hide-details
																:model-value="!isMcpToolExcluded(tool.name, server.config.Exclude)"
																@update:model-value="(val) => setMcpToolExclusion(server.config, tool.name, val)">
											<template v-slot:label>
												<span>
													{{ tool.name }}
													<v-tooltip activator="parent" location="top">{{ tool.description }}</v-tooltip>
												</span>
											</template>
										</v-checkbox>
									</v-col>
								</v-row>
							</v-card-text>
						</v-card>
					</v-col>
				</v-row>
			</v-card-text>
		</v-card>

		<!-- footer -->
		<template v-if="!profile.UI.Prompt.PinTop">
			<v-footer app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<ToolBar />
				</div>
			</v-footer>
		</template>
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { mapActions, mapState } from 'pinia'
import { WindowSetSize } from '../../wailsjs/runtime'
import ZoomDetector from '../components/ZoomDetector.vue'
import { useConfigStore } from '../store/config.ts'
import ToolBar from '../components/bar/ToolBar.vue'
import { ListMcpTools } from '../../wailsjs/go/controller/Controller'
import { mcp, mcp_golang } from '../../wailsjs/go/models.ts'

export default defineComponent({
	name: 'Tool',
	components: { ToolBar, ZoomDetector },
	data() {
		return {
			appbarHeight: 0,
			zoom: this.$appProfile.UI.Window.InitialZoom.Value ?? 1,

			mcp: {
				loading: true,

				server: {} as Record<string, {config: mcp.Server, tools: mcp_golang.ToolRetType[]}>,
			}
		}
	},
	computed: {
		...mapState(useConfigStore, ['profile']),
		SystemInfo: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.SystemInfo.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.SystemInfo.Disable = !value },
		},
		Environment: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.Environment.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.Environment.Disable = !value },
		},
		SystemTime: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.SystemTime.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.SystemTime.Disable = !value },
		},
		Stats: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.Stats.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.Stats.Disable = !value },
		},
		ChangeMode: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.ChangeMode.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.ChangeMode.Disable = !value },
		},
		ChangeOwner: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.ChangeOwner.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.ChangeOwner.Disable = !value },
		},
		ChangeTimes: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.ChangeTimes.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.ChangeTimes.Disable = !value },
		},
		FileCreation: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.FileCreation.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.FileCreation.Disable = !value },
		},
		FileTempCreation: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.FileTempCreation.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.FileTempCreation.Disable = !value },
		},
		FileAppending: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.FileAppending.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.FileAppending.Disable = !value },
		},
		FileReading: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.FileReading.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.FileReading.Disable = !value },
		},
		FileDeletion: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.FileDeletion.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.FileDeletion.Disable = !value },
		},
		DirectoryCreation: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.DirectoryCreation.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.DirectoryCreation.Disable = !value },
		},
		DirectoryTempCreation: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.DirectoryTempCreation.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.DirectoryTempCreation.Disable = !value },
		},
		DirectoryDeletion: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.DirectoryDeletion.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.DirectoryDeletion.Disable = !value },
		},
		CommandExec: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.CommandExec.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.CommandExec.Disable = !value },
		},
		Http: {
			get() { return this.profile.LLM.Tools.BuiltInTools ? !this.profile.LLM.Tools.BuiltInTools.Http.Disable : false },
			set(value: boolean) { if(this.profile.LLM.Tools.BuiltInTools) this.profile.LLM.Tools.BuiltInTools.Http.Disable = !value },
		},
	},
	methods: {
		...mapActions(useConfigStore, ['applyToolsConfig']),
		onZoom(factor: number) {
			this.zoom = factor
			this.adjustHeight()
		},
		async adjustHeight() {
			this.appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0

			const pageHeight = (this.$refs.page as HTMLElement).clientHeight

			//the titlebar can not be manipulated while application lifecycle - so here we use the "initial" config
			const titleBarHeight = this.profile.UI.Window.ShowTitleBar ? (this.profile.UI.Window.TitleBarHeight ?? 0) : 0
			const combinedHeight = Math.ceil(pageHeight * this.zoom) + titleBarHeight
			const width = this.profile.UI.Window.InitialWidth.Value ?? 0

			await WindowSetSize(width, combinedHeight)
		},
		isMcpToolExcluded(tool: string, exclusion: string[] | null): boolean {
			if(!exclusion) {
				return false
			}
			return exclusion.findIndex((e) => e === tool) !== -1
		},
		setMcpToolExclusion(config: { Exclude: string[] }, toolName: string, exclusion: boolean | null) {
			if(!config.Exclude) {
				config.Exclude = []
			}

			if(exclusion) {
				config.Exclude = config.Exclude.filter(e => e !== toolName)
			} else {
				config.Exclude.push(toolName)
			}
		},
	},
	mounted() {
		ListMcpTools().then((tools) => {
			for (let name in tools) {
				this.mcp.server[name] = {
					config: this.profile.LLM.McpServer[name],
					tools: tools[name]
				}
			}
		})
		.finally(() => this.mcp.loading = false)
	},
	beforeRouteLeave() {
		this.applyToolsConfig()
	},
	updated() {
		this.$nextTick(() => {
			this.adjustHeight()
		})
	},
})
</script>

<style scoped></style>