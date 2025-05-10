<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />

		<!-- header -->
		<template v-if="profile.UI.Prompt.PinTop">
			<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<HistoryBar @queryChanged="onQueryChanged" />
				</div>
			</v-app-bar>

			<!-- this will trigger an redraw if attachments are added or removed -->
			<div :style="{ height: `${appbarHeight}px` }"></div>
		</template>

		<v-container>
			<v-row dense>
				<v-col cols="12" v-for="entry in history" :key="entry.m.t">
					<HistoryEntry :model="entry" @onImport="onImport" />
				</v-col>

				<v-col cols="12" v-if="!queried && history.length < total">
					<v-row dense>
						<v-col cols="3" v-for="l of [25, 50, 75, 100]">
							<v-btn block @click="onLoadNext(l)">
								<v-icon icon="mdi-dots-horizontal"></v-icon>
								&nbsp;
								{{l}}
							</v-btn>
						</v-col>
					</v-row>
				</v-col>
			</v-row>
		</v-container>

		<!-- footer -->
		<template v-if="!profile.UI.Prompt.PinTop">
			<v-footer app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<HistoryBar @queryChanged="onQueryChanged" />
				</div>
			</v-footer>
		</template>
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { history } from '../../wailsjs/go/models.ts'
import { HistoryGetCount, HistoryGetLast, HistorySearch } from '../../wailsjs/go/controller/Controller'
import { WindowSetSize } from '../../wailsjs/runtime'
import ZoomDetector from '../components/ZoomDetector.vue'
import HistoryEntry from '../components/HistoryEntry.vue'
import HistoryBar from '../components/bar/HistoryBar.vue'
import { mapActions, mapState } from 'pinia'
import { useHistoryStore } from '../store/history.ts'
import { useConfigStore } from '../store/config.ts'

export default defineComponent({
	name: 'History',
	components: { HistoryBar, HistoryEntry, ZoomDetector },
	data() {
		return {
			appbarHeight: 0,
			zoom: this.$appProfile.UI.Window.InitialZoom.Value,

			queried: false,
			total: 0,
			history: [] as history.Entry[],
		}
	},
	computed: {
		...mapState(useConfigStore, ['profile'])
	},
	methods: {
		...mapActions(useHistoryStore, ['loadConversation']),
		onZoom(factor: number) {
			this.zoom = factor
			this.adjustHeight()
		},
		async adjustHeight() {
			this.appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0

			const pageHeight = (this.$refs.page as HTMLElement).clientHeight

			//the titlebar can not be manipulated while application lifecycle - so here we use the "initial" config
			const titleBarHeight = this.profile.UI.Window.ShowTitleBar ? this.profile.UI.Window.TitleBarHeight : 0
			const combinedHeight = Math.ceil(pageHeight * this.zoom) + titleBarHeight
			const width = this.profile.UI.Window.InitialWidth.Value

			await WindowSetSize(width, combinedHeight)
		},
		onLoadNext(limit: number) {
			HistoryGetLast(this.history.length, limit).then(entries => {
				this.history.push(...entries)
			})
		},
		onQueryChanged(query: string) {
			if(query) {
				HistorySearch(query).then(entries => {
					this.queried = true
					this.history = entries
				}).then(() => {
					this.adjustHeight()
				})
			} else {
				HistoryGetLast(0, 25).then((entries) => {
					this.queried = false
					this.history = entries
				}).then(() => {
					this.adjustHeight()
				})
			}
		},
		onImport(entry: history.Entry) {
			this.loadConversation(entry)
			this.$router.push({ name: 'Chat' })
		}
	},
	mounted() {
		HistoryGetCount().then((count) => {
			this.total = count
		})
		this.onQueryChanged('')
	},
	updated() {
		this.$nextTick(() => {
			this.adjustHeight()
		})
	}
})
</script>

<style scoped></style>