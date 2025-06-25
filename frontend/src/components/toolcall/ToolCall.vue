<template>
	<v-dialog v-model="rejectionDialog" persistent>
		<v-card :title="$t('toolCall.rejection.dialog.title')">
			<template v-slot:text>
				<v-text-field
					v-model="rejectionMessage"
					@keyup.enter="setToolCallApproval(tc, false)"
					:placeholder="$t('toolCall.rejection.dialog.message')"
					clearable
					hide-details
					autofocus />
			</template>

			<template v-slot:actions>
				<v-btn block variant="flat" color="error" @click="setToolCallApproval(tc, false)">
					{{ $t('toolCall.rejection.dialog.ok') }}
				</v-btn>
			</template>
		</v-card>
	</v-dialog>

	<v-card color="chat-tool-call" elevation="0">
		<template v-slot:prepend>
			<v-btn variant="flat"
						 size="x-small"
						 class="mb-1"
						 :color="tc.Result ? (tc.Result.Error ? 'error' : 'success') : 'chat-tool-call'"
						 @click="expanded = !expanded">
				<v-icon :icon="icon"></v-icon>
			</v-btn>
		</template>

		<template v-slot:subtitle>
			<span :class="withoutTitle ? '' : 'ml-2 mr-2'">
				<slot name="title"></slot>
			</span>
		</template>

		<v-card-text v-if="tc.Result" v-show="expanded">
			<v-divider opacity="75" color="info" class="pb-4" />
			<slot name="content"></slot>
			<v-alert type="error" v-if="tc.Result.Error" density="compact">{{ tc.Result.Error }}</v-alert>
		</v-card-text>
		<v-progress-linear indeterminate size="small" v-else></v-progress-linear>

		<v-card-actions v-if="tc.Meta.NeedsApproval && !tc.Result" class="pa-0">
			<v-row dense>
				<v-col cols="6" class="pr-0">
					<v-btn block variant="flat" color="success" @click="setToolCallApproval(tc, true)">
						<v-icon icon="mdi-check"></v-icon>
					</v-btn>
				</v-col>
				<v-col cols="6" class="pl-0">
					<v-btn block variant="flat" color="error" @click="rejectionDialog = true">
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
		},
		withoutTitle: {
			type: Boolean,
			required: false,
			default: false,
		}
	},
	data(){
		return {
			expanded: false,
			rejectionDialog: false,
			rejectionMessage: '',
		}
	},
	methods: {
		setToolCallApproval(call: LLMMessageCall, approved: boolean) {
			call.Meta.NeedsApproval = false
			if (approved) {
				LLMApproveToolCall(call.Id, '')
			} else {
				LLMRejectToolCall(call.Id, this.rejectionMessage)
				this.rejectionMessage = ''
				this.rejectionDialog = false
			}
		},
	}
})
</script>

<style scoped></style>
