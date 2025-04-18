import { defineStore } from 'pinia'
import { controller, history } from '../../wailsjs/go/models.ts'
import LLMMessage = controller.LLMMessage
import MessageContentPart = history.MessageContentPart
import LLMMessageContentPart = controller.LLMMessageContentPart
import LLMMessageCall = controller.LLMMessageCall
import LLMMessageCallResult = controller.LLMMessageCallResult
import { ContentType, Role } from '../components/ChatMessage.vue'
import { useConfigStore } from './config.ts'

export type HistoryEntry = {
	Interrupted: boolean
	Hidden: boolean
	Message: controller.LLMMessage
}

const buildSystemMessage = (): HistoryEntry => ({
	Interrupted: false,
	Hidden: false,
	Message: LLMMessage.createFrom({
		Role: Role.System,
		ContentParts: [{
			Type: ContentType.Text,
			Content: useConfigStore().config.LLM.CallOptions.SystemPrompt,
		}],
		Created: Math.floor(new Date().getTime() / 1000),
	}),
})

export const useHistoryStore = defineStore('history', {
	state: () => ({
		chatHistory: [
			buildSystemMessage(),
		] as HistoryEntry[],
	}),
	actions: {
		setHistory(history: HistoryEntry[]) {
			this.chatHistory = history
		},
		pushHistory(entry: HistoryEntry) {
			this.chatHistory.push(entry)
		},
		replaceHistory(index: number, entry: HistoryEntry) {
			this.chatHistory[index] = entry
		},
		updateHistoryMessage(message: LLMMessage) {
			const i = this.chatHistory.findIndex((entry) => entry.Message.Id === message.Id)
			if (i >= 0) {
				this.chatHistory[i].Message = message
			} else {
				console.error('llm:message:update: message not found', message)
			}
		},
		clearHistory() {
			this.chatHistory = [
				buildSystemMessage(), // preserve system message
			]
		},
		loadConversation(conversation: history.Entry) {
			this.chatHistory = conversation.c.m.map(msg => {
				let entry: HistoryEntry = {
					Interrupted: false,
					Hidden: false,
					Message: LLMMessage.createFrom({
						Id: msg.i,
						Role: msg.r,
						Created: msg.t,
					})
				}
				
				entry.Message.ContentParts = (msg.p ?? []).map((msgPart: MessageContentPart) => {
					let entryPart = LLMMessageContentPart.createFrom({
						Type: msgPart.t,
						Content: msgPart.c,
					})

					if(msgPart.ca) {
						entryPart.Call = LLMMessageCall.createFrom({
							Id: msgPart.ca.i,
							Function: msgPart.ca.f,
							Arguments: msgPart.ca.a,
							BuiltIn: msgPart.ca.f?.startsWith("__")
						})
						if(msgPart.ca.r) {
							entryPart.Call.Result = LLMMessageCallResult.createFrom({
								Content: msgPart.ca.r.c,
								Error: msgPart.ca.r.e,
								DurationMs: msgPart.ca.r.d,
							})
						}
					}

					return entryPart
				})

				return entry
			})
		}
	}
})