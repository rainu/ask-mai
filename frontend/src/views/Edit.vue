<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />

		<!-- header -->
		<template v-if="config.UI.Prompt.PinTop">
			<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<EditBar @on-save="onSave" />
				</div>
			</v-app-bar>

			<!-- this will trigger an redraw if attachments are added or removed -->
			<div :style="{ height: `${appbarHeight}px` }"></div>
		</template>

		<v-container v-if="entryCopy">
			<v-card>
				<v-card-title>
					<div class="d-flex justify-space-between align-center">
						<v-select v-model="entryCopy.Message.Role"
											:items="roleItems" item-title="role"
											density="compact" hide-details max-width="100">
							<template v-slot:selection="{ item }">
								<v-icon>{{item.raw.icon}}</v-icon>
							</template>
							<template v-slot:item="{ props: itemProps, item }">
								<v-list-item v-bind="itemProps" :prepend-icon="item.raw.icon"></v-list-item>
							</template>
						</v-select>

						<v-spacer />

						<span class="opacity-50">{{ createdAt }}</span>
					</div>
				</v-card-title>

				<v-card-text v-for="cp of entryCopy.Message.ContentParts">
					<v-textarea v-model="cp.Content" v-if="cp.Type === 'text'"></v-textarea>
				</v-card-text>
			</v-card>
		</v-container>

		<!-- footer -->
		<template v-if="!config.UI.Prompt.PinTop">
			<v-footer app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<EditBar @on-save="onSave" />
				</div>
			</v-footer>
		</template>
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { WindowSetSize } from '../../wailsjs/runtime'
import ZoomDetector from '../components/ZoomDetector.vue'
import { HistoryEntry, useHistoryStore } from '../store/history.ts'
import { mapActions, mapState } from 'pinia'
import EditBar from '../components/bar/EditBar.vue'
import ChatMessageActions from '../components/ChatMessageActions.vue'
import { Role } from '../components/ChatMessage.vue'
import { useConfigStore } from '../store/config.ts'

export default defineComponent({
	name: 'Edit',
	components: { ChatMessageActions, EditBar, ZoomDetector },
	data() {
		return {
			appbarHeight: 0,
			zoom: this.$appConfig.UI.Window.InitialZoom.Value,

			entryCopy: null as HistoryEntry | null,
		}
	},
	computed: {
		...mapState(useConfigStore, ['config']),
		...mapState(useHistoryStore, ['chatHistory']),
		entryIdx(): number {
			return this.$route.params.idx
		},
		entry(): HistoryEntry {
			return this.chatHistory[this.entryIdx]
		},
		createdAt() {
			if(!this.entry) return null

			const d = new Date(this.entry.Message.Created * 1000)
			return d.toLocaleTimeString()
		},
		roleItems() {
			return [
				{ role: Role.System, icon: "mdi-cog" },
				{ role: Role.Bot, icon: "mdi-robot" },
				{ role: Role.User, icon: "mdi-account" }
			]
		},
	},
	methods: {
		...mapActions(useHistoryStore, ['replaceHistory']),
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
		onSave() {
			if(!this.entryCopy) return

			this.replaceHistory(this.entryIdx, this.entryCopy)
			this.$router.push({ name: 'Chat' })
		}
	},
	activated() {
		try {
			const asString = JSON.stringify(this.entry)
			this.entryCopy = JSON.parse(asString)
		} catch (e) {
			console.error('Failed to copy entry nr. ' + this.entryIdx + ' :', e)
		}
	},
	updated() {
		this.$nextTick(() => {
			this.adjustHeight()
		})
	},
})
</script>

<style scoped></style>