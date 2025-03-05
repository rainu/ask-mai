<template>
	<v-sheet color="grey-lighten-2" rounded>
		<v-expansion-panels variant="accordion" bg-color="grey-lighten-2">
			<v-expansion-panel>
				<v-expansion-panel-title disable-icon-rotate class="pa-2">
					<template v-if="tc.Result">
						<v-icon color="error" icon="mdi-alert-circle" v-if="tc.Result.Error"></v-icon>
						<v-icon color="success" icon="mdi-check-circle" v-else></v-icon>
					</template>
					<v-icon icon="mdi-tools"></v-icon>
					<span class="mx-2">
						{{ parsedArguments.name }} {{ parsedArguments.arguments.join(' ') }}
					</span>
				</v-expansion-panel-title>
				<v-expansion-panel-text v-if="tc.Result">
					<code>{{ code }}</code>
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
import { controller, tools } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import { LLMApproveToolCall, LLMRejectToolCall } from '../../../wailsjs/go/controller/Controller'
import CommandExecutionArguments = tools.CommandExecutionArguments

export default defineComponent({
	name: 'BuiltinToolCallCommandExecution',
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): CommandExecutionArguments {
			return JSON.parse(this.tc.Arguments) as CommandExecutionArguments
		},
		code(){
			let line = ""

			if(this.parsedArguments.working_directory){
				line += "cd " + this.parsedArguments.working_directory + ";\n"
			}

			if(this.parsedArguments.environment) {
				Object.entries(this.parsedArguments.environment).forEach(([key, value]) => {
					line += key + "=" + value + "\n"
				})
			}

			line += this.parsedArguments.name + " " + this.parsedArguments.arguments.join(' ')

			return line
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
