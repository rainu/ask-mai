<template>
	<template v-if="isUserMessage">
		<v-row class="justify-end pa-2 mb-0 mt-1 mx-1 ml-15">
			<v-sheet color="green-accent-2" class="pa-2" rounded>
				<vue-markdown :source="message" :options="options" />
			</v-sheet>
		</v-row>
	</template>
	<template v-else>
		<v-row class="pa-2 mb-0 mt-1 mx-1 mr-15">
			<v-sheet color="grey-lighten-2" class="pa-2" rounded>
				<vue-markdown :source="message" :options="options" />
			</v-sheet>
		</v-row>
	</template>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import VueMarkdown from 'vue-markdown-render'

import hljs from 'highlight.js'
import { type Options as MarkdownItOptions } from 'markdown-it'
import { UseCodeStyle } from './code-style.ts'

export enum Role {
	User = 'human',
	Bot = 'ai',
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
							const result = hljs.highlight(code, { language })
							return result.value
						} catch (__) {}
					}
					return '' // use external default escaping
				},
			} as MarkdownItOptions,
		}
	},
	computed: {
		isUserMessage() {
			return this.role === Role.User
		},
	},
	methods: {
		enrichCopyButtons() {
			const codeBlocks = document.querySelectorAll('pre:not(.code-container pre)').values()
			for (const pre of codeBlocks) {
				const div = document.createElement('div')
				div.className = 'code-container'

				const button = document.createElement('button')
				button.className = 'copy-button mdi-clipboard-text-outline mdi v-icon notranslate v-icon--size-small'
				button.addEventListener('click', this.onCopyButtonClicked)

				div.appendChild(button)
				div.appendChild(pre.cloneNode(true))
				pre.replaceWith(div)
			}
		},
		onCopyButtonClicked(event: MouseEvent) {
			const preElement = (event.target as HTMLButtonElement)?.nextElementSibling as HTMLElement
			if (preElement && preElement.tagName === 'PRE') {
				const code = preElement.innerText
				navigator.clipboard.writeText(code).then(() => {
					preElement.classList.add('copied')
					setTimeout(() => {
						preElement.classList.remove('copied')
					}, 1000)
				})
			}
		},
	},
	watch: {
		message() {
			this.$nextTick(() => this.enrichCopyButtons())
		},
	},
	mounted() {
		UseCodeStyle(this.$appConfig.UI.CodeStyle)

		this.$nextTick(() => this.enrichCopyButtons())
	},
})
</script>

<style>
pre code {
	background-color: #f5f5f5;
	border: 1px solid #ccc;
	margin: 0.5em 0;
	padding: 0.5em 1em;
	border-radius: 5px;
	display: block;
	overflow-x: auto;
	position: relative;
}

/* blue block when hover over code-blocks */
pre code::before {
	content: '';
	position: absolute;
	top: 0;
	left: 0;
	width: 0;
	height: 100%;
	background-color: #007bff;
	transition: width 0.3s;
}

pre code:hover::before {
	width: 5px;
}

/* Add styles for the copy button */

div.code-container {
	position: relative;
}

.copy-button {
	position: absolute;
	right: 1px;
	padding-top: 0.5em;
	padding-right: 0.5em;
	z-index: 1;
	float: right;
}

pre.copied code {
	border: 2px solid #4db6ac;
	border-radius: 5px;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.5); /* Add elevation effect */
}

/* Inline code blocks inside text (not in headers) */
code:not(pre code):not(h1 code):not(h2 code):not(h3 code):not(h4 code):not(h5 code):not(h6 code) {
	background-color: #f5f5f5;
	border: 1px solid #ccc;
	padding: 2px 4px;
	border-radius: 3px;
}

/* Quote-Blocks */

blockquote {
	border-left: 5px solid #ccc;
	margin: 0.5em 0;
	padding: 0.5em 1em;
	color: #555;
	background: none;
	border-radius: 0;
}

blockquote:hover {
	border-color: #007bff;
}

/* Header Padding increase */
h1,
h2,
h3,
h4,
h5,
h6 {
	padding-top: 1em;
}

/* make list items visible */
ol,
ul li {
	margin-left: 2em;
}
</style>
