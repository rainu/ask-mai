import { defineStore } from 'pinia'
import { HistoryEntry } from '../views/Home.vue'
import { controller, history } from '../../wailsjs/go/models.ts'
import LLMMessage = controller.LLMMessage
import MessageContentPart = history.MessageContentPart
import LLMMessageContentPart = controller.LLMMessageContentPart
import LLMMessageCall = controller.LLMMessageCall
import LLMMessageCallResult = controller.LLMMessageCallResult

export const useHistoryStore = defineStore('history', {
	state: () => ({
		conversationToImport: null as HistoryEntry[] | null,
	}),
	actions: {
		loadConversation(conversation: history.Entry) {
			this.conversationToImport = conversation.c.m.map(msg => {
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
		},
		popConversation(): (HistoryEntry[] | null) {
			const conversation = this.conversationToImport
			this.conversationToImport = null
			return conversation
		}
	}
})