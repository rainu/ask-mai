<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />

		<!-- header -->
		<template v-if="config.UI.Prompt.PinTop">
			<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<ToolBar />
				</div>
			</v-app-bar>

			<!-- this will trigger an redraw if attachments are added or removed -->
			<div :style="{ height: `${appbarHeight}px` }"></div>
		</template>

		<v-container>
			<v-row dense no-gutters>
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
		</v-container>

		<!-- footer -->
		<template v-if="!config.UI.Prompt.PinTop">
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
import { WindowSetSize } from '../../wailsjs/runtime'
import ZoomDetector from '../components/ZoomDetector.vue'
import { mapActions, mapState } from 'pinia'
import { useConfigStore } from '../store/config.ts'
import ToolBar from '../components/bar/ToolBar.vue'

export default defineComponent({
	name: 'Tool',
	components: { ToolBar, ZoomDetector },
	data() {
		return {
			appbarHeight: 0,
			zoom: this.$appConfig.UI.Window.InitialZoom.Value,
		}
	},
	computed: {
		...mapState(useConfigStore, ['config']),
		SystemInfo: {
			get() { return !this.config.LLM.Tools.BuiltInTools.SystemInfo.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.SystemInfo.Disable = !value },
		},
		Environment: {
			get() { return !this.config.LLM.Tools.BuiltInTools.Environment.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.Environment.Disable = !value },
		},
		SystemTime: {
			get() { return !this.config.LLM.Tools.BuiltInTools.SystemTime.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.SystemTime.Disable = !value },
		},
		Stats: {
			get() { return !this.config.LLM.Tools.BuiltInTools.Stats.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.Stats.Disable = !value },
		},
		ChangeMode: {
			get() { return !this.config.LLM.Tools.BuiltInTools.ChangeMode.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.ChangeMode.Disable = !value },
		},
		ChangeOwner: {
			get() { return !this.config.LLM.Tools.BuiltInTools.ChangeOwner.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.ChangeOwner.Disable = !value },
		},
		ChangeTimes: {
			get() { return !this.config.LLM.Tools.BuiltInTools.ChangeTimes.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.ChangeTimes.Disable = !value },
		},
		FileCreation: {
			get() { return !this.config.LLM.Tools.BuiltInTools.FileCreation.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.FileCreation.Disable = !value },
		},
		FileTempCreation: {
			get() { return !this.config.LLM.Tools.BuiltInTools.FileTempCreation.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.FileTempCreation.Disable = !value },
		},
		FileAppending: {
			get() { return !this.config.LLM.Tools.BuiltInTools.FileAppending.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.FileAppending.Disable = !value },
		},
		FileReading: {
			get() { return !this.config.LLM.Tools.BuiltInTools.FileReading.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.FileReading.Disable = !value },
		},
		FileDeletion: {
			get() { return !this.config.LLM.Tools.BuiltInTools.FileDeletion.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.FileDeletion.Disable = !value },
		},
		DirectoryCreation: {
			get() { return !this.config.LLM.Tools.BuiltInTools.DirectoryCreation.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.DirectoryCreation.Disable = !value },
		},
		DirectoryTempCreation: {
			get() { return !this.config.LLM.Tools.BuiltInTools.DirectoryTempCreation.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.DirectoryTempCreation.Disable = !value },
		},
		DirectoryDeletion: {
			get() { return !this.config.LLM.Tools.BuiltInTools.DirectoryDeletion.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.DirectoryDeletion.Disable = !value },
		},
		CommandExec: {
			get() { return !this.config.LLM.Tools.BuiltInTools.CommandExec.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.CommandExec.Disable = !value },
		},
		Http: {
			get() { return !this.config.LLM.Tools.BuiltInTools.Http.Disable },
			set(value: boolean) { this.config.LLM.Tools.BuiltInTools.Http.Disable = !value },
		},
	},
	methods: {
		...mapActions(useConfigStore, ['applyBuiltinToolsConfig']),
		onZoom(factor: number) {
			this.zoom = factor
			this.adjustHeight()
		},
		async adjustHeight() {
			this.appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0

			const pageHeight = (this.$refs.page as HTMLElement).clientHeight

			//the titlebar can not be manipulated while application lifecycle - so here we use the "initial" config
			const titleBarHeight = this.$appConfig.UI.Window.ShowTitleBar ? this.$appConfig.UI.Window.TitleBarHeight : 0
			const combinedHeight = Math.ceil(pageHeight * this.zoom) + titleBarHeight
			const width = this.config.UI.Window.InitialWidth.Value

			await WindowSetSize(width, combinedHeight)
		},
	},
	beforeRouteLeave() {
		this.applyBuiltinToolsConfig()
	},
	updated() {
		this.$nextTick(() => {
			this.adjustHeight()
		})
	},
})
</script>

<style scoped></style>