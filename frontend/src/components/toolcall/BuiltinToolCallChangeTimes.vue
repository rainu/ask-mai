<template>
	<ToolCall :tc="tc" icon="mdi-timer-edit">
		<template v-slot:title>
			{{ parsedArguments.path }}
		</template>

		<template v-slot:content>
			<v-row dense>
				<v-col cols="6" >
					<strong>{{$t('toolCall.builtin.changeTimes.access')}}:</strong>
				</v-col>
				<v-col cols="6">
					{{ parsedArguments.access_time }}
				</v-col>

				<v-col cols="6" >
					<strong>{{$t('toolCall.builtin.changeTimes.mod')}}:</strong>
				</v-col>
				<v-col cols="6">
					{{ parsedArguments.modification_time }}
				</v-col>
			</v-row>
		</template>
	</ToolCall>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { controller, file } from '../../../wailsjs/go/models.ts'
import LLMMessageCall = controller.LLMMessageCall
import ChangeTimesArguments = file.ChangeTimesArguments
import ToolCall from './ToolCall.vue'

export default defineComponent({
	name: 'BuiltinToolCallChangeTimes',
	components: { ToolCall },
	props: {
		tc: {
			type: Object as () => LLMMessageCall,
			required: true,
		},
	},
	computed: {
		parsedArguments(): ChangeTimesArguments {
			return JSON.parse(this.tc.Arguments) as ChangeTimesArguments
		}
	}
})
</script>

<style scoped></style>
