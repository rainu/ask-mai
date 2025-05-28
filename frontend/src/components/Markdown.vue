<template>
	<!-- Templates -->
	<div style="display: none;">
		<span ref="tmpl-link" class="link-container">
			<v-btn size="x-small" variant="outlined" class="href text-none">LINK</v-btn>
			<v-btn size="x-small" variant="tonal" class="open">
				<v-icon>mdi-open-in-new</v-icon>
			</v-btn>
		</span>
	</div>

	<vue-markdown :source="content" :options="options" ref="markdown" />
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import hljs from 'highlight.js'
import { type Options as MarkdownItOptions } from 'markdown-it'
import VueMarkdown from 'vue-markdown-render'
import { BrowserOpenURL, ClipboardSetText } from '../../wailsjs/runtime'
import { UseCodeStyle } from './code-style.ts'
import { mapState } from 'pinia'
import { useConfigStore } from '../store/config.ts'

export default defineComponent({
	name: 'Markdown',
	components: { VueMarkdown },
	props: {
		content: {
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
			observer: null as MutationObserver | null,
		}
	},
	computed: {
		...mapState(useConfigStore, ['profile']),
	},
	methods: {
		patchMarkdown() {
			this.enrichCopyButtons()
			this.replaceLinks()
		},
		enrichCopyButtons() {
			const codeBlocks = this.$refs.markdown.$el.querySelectorAll('pre:not(.code-container pre)').values()
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
				let code = preElement.innerText
				const newlineCount = (code.match(/\n/g) || []).length
				if (newlineCount === 1) {
					//prevent that shell-code-statements will be executed directly
					code = code.trim()
				}

				ClipboardSetText(code).then(() => {
					preElement.classList.add('copied')
					setTimeout(() => {
						preElement.classList.remove('copied')
					}, 1000)
				})
			}
		},

		replaceLinks() {
			const links = this.$refs.markdown.$el.querySelectorAll('a').values()
			for (const link of links) {
				const span = this.$refs['tmpl-link'].cloneNode(true)
				const href = span.querySelector('.href')
				const openBtn = span.querySelector('button.open')

				href.textContent = link.textContent
				href.addEventListener('click', this.onLinkClicked)
				openBtn.addEventListener('click', this.onOpenLinkClicked)

				link.replaceWith(span)
			}
		},
		onLinkClicked(event: MouseEvent) {
			const container = (event.target as HTMLButtonElement)?.closest('.link-container');

			if (container) {
				const hrefElement = container.querySelector('.href');
				if (hrefElement && hrefElement.textContent) {
					ClipboardSetText(hrefElement.textContent).then(() => {
						hrefElement.classList.add('copied')
						setTimeout(() => {
							hrefElement.classList.remove('copied')
						}, 1000)
					})
				}
			}
		},
		onOpenLinkClicked(event: MouseEvent) {
			const hrefElement = (event.target as HTMLButtonElement)?.closest('.link-container')?.querySelector('.href')
			const href = hrefElement?.textContent

			if(href) {
				BrowserOpenURL(href)
			}
		}
	},
	mounted() {
		UseCodeStyle(this.profile.UI.CodeStyle)

		this.observer = new MutationObserver(() => this.patchMarkdown())
		this.observer.observe(this.$el, { childList: true, subtree: true })
		this.patchMarkdown()
	},
	beforeUnmount() {
		if(this.observer) {
			this.observer.disconnect()
		}
	},
})
</script>

<style scoped>
:deep(pre code) {
	background-color: #f5f5f5;
	border: 1px solid #ccc;
	margin: 0.5em 0;
	padding: 0.5em 2em 0.5em 1em;
	border-radius: 5px;
	display: block;
	overflow-x: auto;
	position: relative;
}

/* blue block when hover over code-blocks */
:deep(pre code::before) {
	content: '';
	position: absolute;
	top: 0;
	left: 0;
	width: 0;
	height: 100%;
	background-color: #007bff;
	transition: width 0.3s;
}

:deep(pre code:hover::before) {
	width: 5px;
}

/* Add styles for the copy button */

:deep(div.code-container) {
	position: relative;
}

:deep(.copy-button) {
	position: absolute;
	right: 1px;
	padding-top: 1em;
	padding-right: 0.5em;
	z-index: 1;
	float: right;
}

:deep(pre.copied code) {
	border: 2px solid #4db6ac;
	border-radius: 5px;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.5); /* Add elevation effect */
}

:deep(button.copied) {
	border: 2px solid #4db6ac;
	border-radius: 5px;
}

/* Inline code blocks inside text (not in headers) */
:deep(code:not(pre code):not(h1 code):not(h2 code):not(h3 code):not(h4 code):not(h5 code):not(h6 code)) {
	background-color: #f5f5f5;
	border: 1px solid #ccc;
	padding: 2px 4px;
	border-radius: 3px;
}

/* Quote-Blocks */

:deep(blockquote) {
	border-left: 5px solid #ccc;
	margin: 0.5em 0;
	padding: 0.5em 1em;
	color: #555;
	background: none;
	border-radius: 0;
}

:deep(blockquote:hover) {
	border-color: #007bff;
}

/* Header Padding increase */
:deep(h1),
:deep(h2),
:deep(h3),
:deep(h4),
:deep(h5),
:deep(h6) {
	padding-top: 1em;
}

/* make list items visible */
:deep(ol),
:deep(ul li) {
	margin-left: 2em;
}
</style>

