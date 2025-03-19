<template>
	<v-sheet color="grey-lighten-2" rounded>
		<v-expansion-panels variant="accordion" bg-color="grey-lighten-2">
			<v-expansion-panel>
				<v-expansion-panel-title disable-icon-rotate class="pa-2">
					<template v-if="tc.Result">
						<v-icon color="error" icon="mdi-alert-circle" v-if="tc.Result.Error"></v-icon>
						<v-icon color="success" icon="mdi-check-circle" v-else></v-icon>
					</template>
					<v-icon icon="mdi-file-eye"></v-icon>
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
					<div v-if="parsedArguments.limits">
						<span>{{ parsedArguments.limits.m }}: {{ offset }} - {{ limit }}</span>
					</div>
					<pre v-if="parsedResult">{{ parsedResult.content }}</pre>
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
import FileReadingArguments = tools.FileReadingArguments
import FileReadingResult = tools.FileReadingResult

export default defineComponent({
	name: 'BuiltinToolCallFileReading',
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): FileReadingArguments {
			return JSON.parse(this.tc.Arguments) as FileReadingArguments
		},
		parsedResult(): FileReadingResult | null {
			if(this.tc.Result) {
				try {
					return JSON.parse(this.tc.Result.Content) as FileReadingResult
				} catch (e) {
					// ignore JSON-Parse error
				}
			}
			return null
		},
		offset(): number {
			return this.parsedArguments.limits?.o || 0
		},
		limit(): number {
			return this.parsedArguments.limits?.l || 0
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
