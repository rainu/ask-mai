<template>
  <v-row class="justify-end pa-2 mb-0 mt-1 mx-1" :class="classes">
    <v-sheet :color="color" class="pa-2" rounded>
      <vue-markdown :source="message" :options="options" />
    </v-sheet>
  </v-row>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import VueMarkdown from 'vue-markdown-render'

import 'highlight.js/styles/github.css'
import hljs from 'highlight.js'

export enum Role {
  User = 'user',
  Bot = 'bot'
}

export default defineComponent({
	name: 'ChatMessage',
	components: { VueMarkdown },
	props: {
		message: {
			type: String,
			required: true,
		},
    role: {
      type: String,
      required: true,
    },
	},
	data() {
		return {
			options: {
				highlight: (code: string, language: string) => {
					if (language && hljs.getLanguage(language)) {
						try {
							return hljs.highlight(code, { language }).value
						} catch (__) {}
					}
					return '' // use external default escaping
				},
			},
		}
	},
  computed: {
    color() {
      return this.role === Role.Bot ? 'grey-lighten-2' : 'green-accent-2'
    },
    classes() {
      return {
        'mr-15': this.role === Role.Bot,
        'ml-15': this.role === Role.User,
      }
    }
  }
})
</script>

<style scoped>
</style>
