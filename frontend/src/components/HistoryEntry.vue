<template>
	<v-card>
		<v-card-title>{{ firstUserPrompt }}</v-card-title>
		<v-card-subtitle>
			<v-icon icon="mdi-account" size="x-small" class="mr-1" />{{ humanCount }} /
			<v-icon icon="mdi-robot" size="x-small" class="mr-1" />{{ aiCount }} /
			<v-icon icon="mdi-tools" size="x-small" class="mr-1" />{{ toolCount }}
		</v-card-subtitle>
		<v-card-subtitle>{{ date }}</v-card-subtitle>

		<v-card-actions>
			<v-row dense>
				<v-col cols="6" md="4" >
					<v-btn block @click="showConversation = !showConversation">
						<v-icon size="x-large" v-if="showConversation">mdi-chevron-down</v-icon>
						<v-icon size="x-large" v-else>mdi-chevron-right</v-icon>
					</v-btn>
				</v-col>
				<v-col cols="0" md="4" />
				<v-col cols="6" md="4" >
					<v-btn block color="primary" @click="onImport">
						<v-icon size="x-large">mdi-import</v-icon>
					</v-btn>
				</v-col>
			</v-row>
		</v-card-actions>

		<v-card-text v-show="showConversation">
			<v-list density="compact">
				<v-list-item v-for="c of conversation">
					<template v-slot:prepend>
						<v-icon :icon="c.icon"></v-icon>
					</template>
					<v-list-item-title>{{c.content}}</v-list-item-title>
				</v-list-item>
			</v-list>
		</v-card-text>
	</v-card>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { history } from '../../wailsjs/go/models.ts'
import { ContentType, Role } from './ChatMessage.vue'

export default defineComponent({
	name: 'HistoryEntry',
	emits: ['onImport'],
	props: {
		model: {
			type: Object as () => history.Entry,
			required: true,
		},
	},
	data(){
		return {
			showConversation: false,
		}
	},
	computed: {
		firstUserMessage(): history.Message | undefined {
			return this.model.c.m.find((m) => m.r === Role.User)
		},
		firstUserPrompt(): string | null {
			if (!this.firstUserMessage) return null
			if (!this.firstUserMessage.p) return null

			return this.firstUserMessage.p
				.filter((p) => p.t === ContentType.Text)
				.map((p) => p.c)
				.join('')
		},
		humanCount(): number {
			return this.model.c.m.filter((m) => m.r === Role.User).length
		},
		aiCount(): number {
			return this.model.c.m.filter((m) => m.r === Role.Bot).length
		},
		toolCount(): number {
			return this.model.c.m.filter((m) => m.r === Role.Tool).length
		},
		date(): Date {
			return new Date(this.model.m.t)
		},
		conversation() {
			return this.model.c.m.map(m => {
				let content = ""

				if(m.p) {
					for (let part of m.p) {
						if(part.c) {
							content += part.c
						}

						if (part.ca?.f) {
							content += part.ca.f
						}
					}
				}

				let icon = ""
				switch (m.r) {
					case Role.User:
						icon = "mdi-account"
						break
					case Role.Bot:
						icon = "mdi-robot"
						break
					case Role.Tool:
						icon = "mdi-tools"
						break
				}

				return {
					icon,
					content
				}
			})
		},
	},
	methods: {
		onImport(){
			this.$emit('onImport', this.model)
		}
	}
})
</script>

<style scoped></style>