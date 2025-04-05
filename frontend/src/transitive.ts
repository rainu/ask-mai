import HistoryEntry from './components/HistoryEntry.vue'

interface TransitiveState {
	lastConversation: HistoryEntry | null
}

declare global {
	interface Window {
		transitiveState: TransitiveState
	}
}

window.transitiveState = {
	lastConversation: null as HistoryEntry | null,
}
