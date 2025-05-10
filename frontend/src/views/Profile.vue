<template>
	<div ref="page" :style="{ zoom }">
		<ZoomDetector @onZoom="onZoom" />

		<!-- header -->
		<template v-if="profile.UI.Prompt.PinTop">
			<v-app-bar app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<ProfileBar />
				</div>
			</v-app-bar>

			<!-- this will trigger an redraw if attachments are added or removed -->
			<div :style="{ height: `${appbarHeight}px` }"></div>
		</template>

		<v-container>
			<v-list lines="three" density="compact" activatable class="pa-0">
				<v-list-item v-for="p of profiles" :key="p.name" :active="p.name === activeProfileName" @click="onChooseProfile(p.name)">
					<template v-slot:prepend>
						<v-avatar>
							<v-icon v-if="p.icon">{{ p.icon }}</v-icon>
							<v-icon v-else>mdi-cog</v-icon>
						</v-avatar>
					</template>
					<template v-slot:title>
						<span v-if="p.name">{{ p.name }}</span>
						<span v-else class="font-italic">{{ $t('profile.default') }}</span>
					</template>
					<template v-slot:subtitle>
						<span v-if="p.description">{{ p.description }}</span>
						<span v-else>-</span>
					</template>
				</v-list-item>
			</v-list>
		</v-container>

		<!-- footer -->
		<template v-if="!profile.UI.Prompt.PinTop">
			<v-footer app class="pa-0 ma-0" density="compact" height="auto">
				<div style="width: 100%" ref="appbar">
					<ProfileBar />
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
import ProfileBar from '../components/bar/ProfileBar.vue'
import { GetAvailableProfiles } from '../../wailsjs/go/controller/Controller'

type Profile = {
	name: string
	description: string
	icon: string
}

export default defineComponent({
	name: 'Profile',
	components: { ProfileBar, ZoomDetector },
	data() {
		return {
			appbarHeight: 0,
			zoom: this.$appProfile.UI.Window.InitialZoom.Value,

			profiles: [] as Profile[],
		}
	},
	computed: {
		...mapState(useConfigStore, ['profile', 'activeProfileName']),
	},
	methods: {
		...mapActions(useConfigStore, ['setActiveProfile']),
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
		onChooseProfile(name: string) {
			this.setActiveProfile(name)
		},
	},
	created() {
		GetAvailableProfiles().then((profiles) => {
			this.profiles = Object.entries(profiles)
				.map(([name, profile]) => {
					return {
						name,
						description: profile.Description,
						icon: profile.Icon,
					}
				})
				.sort((a, b) => a.name.localeCompare(b.name))
		})
	},
	updated() {
		this.$nextTick(() => {
			this.adjustHeight()
		})
	},
})
</script>

<style scoped></style>