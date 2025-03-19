<template>
	<v-sheet color="grey-lighten-2" rounded>
		<v-expansion-panels variant="accordion" bg-color="grey-lighten-2">
			<v-expansion-panel>
				<v-expansion-panel-title disable-icon-rotate class="pa-2">
					<template v-if="tc.Result">
						<v-icon color="error" icon="mdi-alert-circle" v-if="tc.Result.Error"></v-icon>
						<v-icon color="success" icon="mdi-check-circle" v-else></v-icon>
					</template>
					<v-icon icon="mdi-file-document-plus"></v-icon>
					<span class="mx-2">
						<template v-if="parsedResult">
							{{ parsedResult.path }}
						</template>
						<template v-else>
							{{ parsedArguments.path }}
						</template>
					</span>
				</v-expansion-panel-title>
				<v-expansion-panel-text v-if="tc.Result">
					<pre>{{ parsedArguments.content }}</pre>
					<v-alert type="error" v-if="tc.Result.Error" density="compact">{{ tc.Result.Error }}</v-alert>
				</v-expansion-panel-text>
				<v-progress-linear indeterminate size="small" v-else></v-progress-linear>

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
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import { LLMApproveToolCall, LLMRejectToolCall } from '../../../wailsjs/go/controller/Controller'
import FileCreationArguments = tools.FileCreationArguments
import FileCreationResult = tools.FileCreationResult

export default defineComponent({
	name: 'BuiltinToolCallFileCreation',
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): FileCreationArguments {
			return JSON.parse(this.tc.Arguments) as FileCreationArguments
		},
		parsedResult(): FileCreationResult | null {
			if(this.tc.Result) {
				try {
					return JSON.parse(this.tc.Result.Content) as FileCreationResult
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
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
