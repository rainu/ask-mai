<template>
	<v-sheet color="grey-lighten-2" rounded>
		<v-expansion-panels variant="accordion" bg-color="grey-lighten-2">
			<v-expansion-panel>
				<v-expansion-panel-title disable-icon-rotate class="pa-2">
					<template v-if="tc.Result">
						<v-icon color="error" icon="mdi-alert-circle" v-if="tc.Result.Error"></v-icon>
						<v-icon color="success" icon="mdi-check-circle" v-else></v-icon>
					</template>
					<v-icon icon="mdi-function"></v-icon>
					<span class="mr-2" v-html="functionHeader"></span>
				</v-expansion-panel-title>
				<v-expansion-panel-text v-if="tc.Result">
					<pre>{{ tc.Result.Content }}</pre>
					<v-alert type="error" v-if="tc.Result.Error" density="compact">{{ tc.Result.Error }}</v-alert>
				</v-expansion-panel-text>

				<v-progress-linear indeterminate size="small" v-if="!tc.Result"></v-progress-linear>

				<template v-if="tc.NeedsApproval && !tc.Result">
					<v-row dense>
						<v-col cols="6" class="pr-0">
							<v-btn block color="success" @click="setToolCallApproval(tc, true)">
								<v-icon icon="mdi-check"></v-icon>
							</v-btn>
						</v-col>
						<v-col cols="6" class="pl-0">
							<v-btn block color="error" @click="setToolCallApproval(tc, false)">
								<v-icon icon="mdi-close"></v-icon>
							</v-btn>
						</v-col>
					</v-row>
				</template>
			</v-expansion-panel>
		</v-expansion-panels>
	</v-sheet>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import hljs from 'highlight.js'
import { LLMApproveToolCall, LLMRejectToolCall } from '../../../wailsjs/go/controller/Controller'

export default defineComponent({
	name: 'GeneralToolCall',
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		functionHeader() {
			let code = this.tc.Function + "(" + this.tc.Arguments + ")"
			return this.highlight(code)
		}
	},
	methods: {
		highlight(code: string) {
			return hljs.highlight(code, { language: 'JavaScript' }).value
		},
		setToolCallApproval(call: LLMMessageCall, approved: boolean) {
			call.NeedsApproval = false
			if (approved) {
				LLMApproveToolCall(call.Id)
			} else {
				LLMRejectToolCall(call.Id)
			}
		},
	}
})
</script>

<style scoped></style>
