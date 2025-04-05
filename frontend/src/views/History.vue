<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />

		<!-- header -->
		<template v-if="$appConfig.UI.Prompt.PinTop">
			<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<HistoryInput @queryChanged="onQueryChanged" />
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
		<template v-if="!$appConfig.UI.Prompt.PinTop">
			<v-footer app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<HistoryInput @queryChanged="onQueryChanged" />
				</div>
			</v-footer>
		</template>
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, history } from '../../wailsjs/go/models.ts'
import { HistoryGetCount, HistoryGetLast, HistorySearch } from '../../wailsjs/go/controller/Controller'
import { WindowSetSize } from '../../wailsjs/runtime'
import ZoomDetector from '../components/ZoomDetector.vue'
import HistoryEntry from '../components/HistoryEntry.vue'
import InputRow from '../components/InputRow.vue'
import HistoryInput from '../components/HistoryInput.vue'
import LLMMessage = controller.LLMMessage
import LLMMessageContentPart = controller.LLMMessageContentPart
import LLMMessageCall = controller.LLMMessageCall
import LLMMessageCallResult = controller.LLMMessageCallResult
import MessageContentPart = history.MessageContentPart

export default defineComponent({
	name: 'History',
	components: { HistoryInput, InputRow, HistoryEntry, ZoomDetector },
	data() {
		return {
			appbarHeight: 0,
			zoom: this.$appConfig.UI.Window.InitialZoom.Value,

			queried: false,
			total: 0,
			history: [] as history.Entry[],
		}
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
			//transform history entry to LLMMessage array
			window.transitiveState.lastConversation = entry.c.m.map(msg => {
				let entry: HistoryEntry = {
					Message: LLMMessage.createFrom({
						Role: msg.r,
					})
				}
				entry.Message.ContentParts = msg.p.map((msgPart: MessageContentPart) => {
					let entryPart = LLMMessageContentPart.createFrom({
						Type: msgPart.t,
						Content: msgPart.c,
					})

					if(msgPart.ca) {
						entryPart.Call = LLMMessageCall.createFrom({
							Id: msgPart.ca.i,
							Function: msgPart.ca.f,
							Arguments: msgPart.ca.a,
							BuiltIn: msgPart.ca.f?.startsWith("__")
						})
						if(msgPart.ca.r) {
							entryPart.Call.Result = LLMMessageCallResult.createFrom({
								Content: msgPart.ca.r.c,
								Error: msgPart.ca.r.e,
								DurationMs: msgPart.ca.r.d,
							})
						}
					}

					return entryPart
				})

				return entry
			})

			this.$router.push({ name: 'Home' })
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