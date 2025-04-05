<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />

		<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
			<div style="width: 100%" ref="appbar">
				<v-row dense>
					<v-col :cols="isMinMaxEnabled ? 1 : 0" class="pa-0 ma-0"></v-col>
					<v-col :cols="isMinMaxEnabled ? 11 : 12" class="pa-0 ma-0">
						<v-text-field
							v-model="query"
							@change="onQueryChanged"
							hide-details
							autofocus
							:placeholder="$t('history.placeholder')"
						>

							<template v-slot:prepend-inner>
								<v-btn icon density="compact" @click="onNavigateBack">
									<v-icon>mdi-chat-processing-outline</v-icon>
								</v-btn>
							</template>
						</v-text-field>
					</v-col>
				</v-row>

			</div>
		</v-app-bar>

		<!-- this will trigger an redraw if attachments are added or removed -->
		<div :style="{ height: `${appbarHeight}px` }"></div>

		<v-container>
			<v-row dense>
				<v-col cols="12" v-for="entry in history" :key="entry.m.t">
					<HistoryEntry :model="entry" />
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
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { history } from '../../wailsjs/go/models.ts'
import { HistoryGetCount, HistoryGetLast, HistorySearch } from '../../wailsjs/go/controller/Controller'
import ChatInputBar from '../components/ChatInputBar.vue'
import { WindowSetSize } from '../../wailsjs/runtime'
import ZoomDetector from '../components/ZoomDetector.vue'
import HistoryEntry from '../components/HistoryEntry.vue'

export default defineComponent({
	name: 'History',
	components: { HistoryEntry, ZoomDetector, ChatInputBar },
	data() {
		return {
			appbarHeight: 0,
			zoom: this.$appConfig.UI.Window.InitialZoom.Value,

			query: '',
			queried: false,
			total: 0,
			history: [] as history.Entry[],
		}
	},
	computed: {
		isMinMaxLeft(){
			return this.$appConfig.UI.MinMaxPosition === 'left'
		},
		isMinMaxRight(){
			return this.$appConfig.UI.MinMaxPosition === 'right'
		},
		isMinMaxEnabled(){
			return this.isMinMaxLeft || this.isMinMaxRight
		},
	},
	methods: {
		onZoom(factor: number) {
			this.zoom = factor
			this.adjustHeight()
		},
		async adjustHeight() {
			this.appbarHeight = this.$refs.appbar ? (this.$refs.appbar as HTMLElement).clientHeight : 0

			const pageHeight = (this.$refs.page as HTMLElement).clientHeight
			const combinedHeight = Math.ceil(pageHeight * this.zoom)
			const width = this.$appConfig.UI.Window.InitialWidth.Value

			await WindowSetSize(width, combinedHeight)
		},
		onNavigateBack() {
			this.$router.push({ name: 'Home' })
		},
		onLoadNext(limit: number) {
			HistoryGetLast(this.history.length, limit).then(entries => {
				this.history.push(...entries)
			})
		},
		onQueryChanged() {
			if(this.query) {
				HistorySearch(this.query).then(entries => {
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
		}
	},
	mounted() {
		HistoryGetCount().then((count) => {
			this.total = count
		})
		this.onQueryChanged()
	},
	updated() {
		this.$nextTick(() => {
			this.adjustHeight()
		})
	}
})
</script>

<style scoped></style>