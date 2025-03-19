<template>
	<v-card color="grey-lighten-2">
		<template v-slot:prepend>
			<v-btn variant="flat"
						 size="x-small"
						 class="mb-1"
						 :color="tc.Result ? (tc.Result.Error ? 'error' : 'success') : 'grey-lighten-2'"
						 @click="expanded = !expanded">
				<v-icon :icon="icon"></v-icon>
			</v-btn>
		</template>

		<template v-slot:subtitle>
			<span class="ml-2 mr-2">
				<slot name="title"></slot>
			</span>
		</template>

		<v-card-text v-if="tc.Result" v-show="expanded">
			<v-divider opacity="75" color="info" class="pb-4" />
			<slot name="content"></slot>
			<v-alert type="error" v-if="tc.Result.Error" density="compact">{{ tc.Result.Error }}</v-alert>
		</v-card-text>
		<v-progress-linear indeterminate size="small" v-else></v-progress-linear>

		<v-card-actions v-if="tc.NeedsApproval && !tc.Result" class="pa-0">
			<v-row dense>
				<v-col cols="6" class="pr-0">
					<v-btn block variant="flat" color="success" @click="setToolCallApproval(tc, true)">
						<v-icon icon="mdi-check"></v-icon>
					</v-btn>
				</v-col>
				<v-col cols="6" class="pl-0">
					<v-btn block variant="flat" color="error" @click="setToolCallApproval(tc, false)">
						<v-icon icon="mdi-close"></v-icon>
					</v-btn>
				</v-col>
			</v-row>
		</v-card-actions>
	</v-card>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import { LLMApproveToolCall, LLMRejectToolCall } from '../../../wailsjs/go/controller/Controller'

export default defineComponent({
	name: 'ToolCall',
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
		icon: {
			type: String,
			required: false,
			default: 'mdi-function',
		}
	},
	data(){
		return {
			expanded: false,
		}
	},
	methods: {
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
